package synchronizer

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	seer_blockchain "github.com/G7DAO/seer/blockchain"
	"github.com/G7DAO/seer/crawler"
	"github.com/G7DAO/seer/indexer"
	"github.com/G7DAO/seer/storage"
	"github.com/ethereum/go-ethereum/common"
)

type Synchronizer struct {
	Client          seer_blockchain.BlockchainClient
	StorageInstance storage.Storer

	blockchain string
	startBlock uint64
	endBlock   uint64
	batchSize  uint64
	baseDir    string
	basePath   string
	threads    int
}

// NewSynchronizer creates a new synchronizer instance with the given blockchain handler.
func NewSynchronizer(blockchain, baseDir string, startBlock, endBlock, batchSize uint64, timeout int, threads int) (*Synchronizer, error) {
	var synchronizer Synchronizer

	basePath := filepath.Join(baseDir, crawler.SeerCrawlerStoragePrefix, "data", blockchain)
	storageInstance, err := storage.NewStorage(storage.SeerCrawlerStorageType, basePath)
	if err != nil {
		log.Fatalf("Failed to create storage instance: %v", err)
		panic(err)
	}

	client, err := seer_blockchain.NewClient(blockchain, crawler.BlockchainURLs[blockchain], timeout)
	if err != nil {
		log.Println("Error initializing blockchain client:", err)
		log.Fatal(err)
	}

	log.Printf("Initialized new synchronizer at blockchain: %s, startBlock: %d, endBlock: %d", blockchain, startBlock, endBlock)

	if threads <= 0 {
		threads = 1
	}

	synchronizer = Synchronizer{
		Client:          client,
		StorageInstance: storageInstance,

		blockchain: blockchain,
		startBlock: startBlock,
		endBlock:   endBlock,
		batchSize:  batchSize,
		baseDir:    baseDir,
		basePath:   basePath,
		threads:    threads,
	}

	return &synchronizer, nil
}

// Read index storage

// -------------------------------------------------------------------------------------------------------------------------------
//                 ________Moonstream blockstore________
// crawler 1 ---->|       Contain raw                   |
// 			      |        blocks                       |
// crawler 2 ---->|        transactions                 |
// 			      |        events                       |
// crawler 3 ---->|                                     |
// 			      |                                     |---->{synchonizer} --------> {user_RDS}
// crawler 4 ---->|                                     |     read from blockstore
// 			      |                                     |     via indexer
// crawler 5 ---->|                                     |     Decode transactions and events
// ....		      |									    |     using ABIs and addresses
// 			      |									    |     and write to user_RDS
// crawler n ---->|								     	|
// 			      |									    |
// 		    	  |_____________________________________|

// -------------------------------------------------------------------------------------------------------------------------------

type Customer struct {
	Id             string             `json:"id"`
	Name           string             `json:"name"`
	NormalizedName string             `json:"normalized_name"`
	Instances      []CustomerInstance `json:"instances"`
}

type CustomerInstance struct {
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	Ip          string        `json:"ip"`
	Is_running  bool          `json:"is_running"`
	Credentials []interface{} `json:"credentials"`
}

func GetDBConnection(uuid string, id int) (string, error) {

	// Create the request
	req, err := http.NewRequest("GET", MOONSTREAM_DB_V3_CONTROLLER_API+"/customers/"+uuid+"/instances/"+strconv.Itoa(id)+"/creds/seer/url", nil)
	if err != nil {
		return "", err
	}

	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN)

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to get connection string for id", uuid, ":", MOONSTREAM_DB_V3_CONTROLLER_API+"/customers/"+uuid+"/instances/"+strconv.Itoa(id)+"/creds/seer/url")
		return "", fmt.Errorf("failed to get connection string for id %s: %s", uuid, resp.Status)
	}

	// Read string from body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	connectionString := strings.Trim(string(bodyBytes), "\"")

	// Validate and possibly correct the connection string
	connectionString, err = ensurePortInConnectionString(connectionString)
	if err != nil {
		return "", fmt.Errorf("invalid connection string for id %s: %v", uuid, err)
	}

	return connectionString, nil
}

func GetCustomerInstances(uuid string) ([]int, error) {
	// Create the request
	req, err := http.NewRequest("GET", MOONSTREAM_DB_V3_CONTROLLER_API+"/customers/"+uuid, nil)
	if err != nil {
		return nil, err
	}

	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN)

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to get instances for id", uuid, ":", MOONSTREAM_DB_V3_CONTROLLER_API+"/customers/"+uuid+"/instances")
		return nil, fmt.Errorf("failed to get instances for id %s: %s", uuid, resp.Status)
	}

	// Read string from body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var customerInstances Customer

	err = json.Unmarshal(bodyBytes, &customerInstances)
	if err != nil {
		return nil, err
	}

	var instances []int

	for _, instance := range customerInstances.Instances {
		if instance.Is_running {
			instances = append(instances, instance.Id)
		}
	}

	return instances, nil
}

func (d *Synchronizer) ReadAbiJobsFromDatabase(blockchain string) ([]indexer.AbiJob, error) {
	abiJobs, err := indexer.DBConnection.ReadABIJobs(blockchain)
	if err != nil {
		return nil, err
	}
	return abiJobs, nil
}

func ensurePortInConnectionString(connStr string) (string, error) {
	parsedURL, err := url.Parse(connStr)
	if err != nil {
		return "", err
	}

	host := parsedURL.Host
	if !strings.Contains(host, ":") {
		host = host + ":5432" // Add the default PostgreSQL port
		parsedURL.Host = host
	}

	return parsedURL.String(), nil
}

type CustomerDBConnection struct {
	Uri string

	Pgx *indexer.PostgreSQLpgx
}

func CustomersLatestBlocks(customerDBConnections map[string]map[int]CustomerDBConnection, blockchain string) (map[string]uint64, error) {
	latestBlocks := make(map[string]uint64)
	for id, customerInstances := range customerDBConnections {
		for _, instanceConnection := range customerInstances {
			pool := instanceConnection.Pgx.GetPool()
			conn, err := pool.Acquire(context.Background())
			if err != nil {
				log.Println("Error acquiring pool connection: ", err)
				return nil, err
			}
			defer conn.Release()

			latestLabelBlock, err := instanceConnection.Pgx.ReadLastLabel(blockchain)
			if err != nil {
				log.Println("Error reading latest block: ", err)
				return nil, err
			}
			/// check if id exists

			if _, ok := latestBlocks[id]; ok {
				if latestBlocks[id] > latestLabelBlock {
					latestBlocks[id] = latestLabelBlock
				}
			} else {
				latestBlocks[id] = latestLabelBlock
			}
			log.Printf("Latest block for customer %s is: %d\n", id, latestLabelBlock)

		}
	}
	return latestBlocks, nil
}

func (d *Synchronizer) StartBlockLookUp(customerDBConnections map[string]map[int]CustomerDBConnection, blockchain string, blockShift int64) (uint64, error) {
	var startBlock uint64
	latestCustomerBlocks, err := CustomersLatestBlocks(customerDBConnections, d.blockchain)
	if err != nil {
		return startBlock, fmt.Errorf("error getting latest blocks for customers: %w", err)
	}

	// Determine the start block as the maximum of the latest blocks of all customers
	var maxCustomerLatestBlock uint64
	if len(latestCustomerBlocks) > 0 {
		for _, block := range latestCustomerBlocks {
			if block > maxCustomerLatestBlock {
				maxCustomerLatestBlock = block
			}
		}
	}

	if maxCustomerLatestBlock != 0 {
		startBlock = maxCustomerLatestBlock
	} else {
		// In case start block is still 0, get the latest block from the blockchain minus shift
		latestBlockNumber, latestErr := d.Client.GetLatestBlockNumber()
		if latestErr != nil {
			return startBlock, fmt.Errorf("error getting latest block number: %w", latestErr)
		}
		startBlock = uint64(latestBlockNumber.Int64() - blockShift)
	}

	return startBlock, nil

}

// getCustomers fetches ABI jobs, extracts customer IDs, and establishes database connections.
func (d *Synchronizer) getCustomers(customerDbUriFlag string, customerIds []string) (map[string]map[int]CustomerDBConnection, []string, error) {
	customerDBConnections := make(map[string]map[int]CustomerDBConnection)
	customerIdSet := make(map[string]struct{})

	// If no customer IDs are provided as arguments, fetch them from ABI jobs.
	if len(customerIds) == 0 {
		abiJobs, err := d.ReadAbiJobsFromDatabase(d.blockchain)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read ABI jobs: %w", err)
		}

		// Extract unique customer IDs from the ABI jobs.
		for _, job := range abiJobs {
			customerIdSet[job.CustomerID] = struct{}{}
		}
	} else {
		// Otherwise, use the provided customer IDs directly.
		for _, id := range customerIds {
			customerIdSet[id] = struct{}{}
		}
	}

	var finalCustomerIds []string
	for id := range customerIdSet {
		var connectionString string
		var dbConnErr error

		// get list of customer ids

		instances, instancesErr := GetCustomerInstances(id)

		if instancesErr != nil {
			log.Printf("Error getting customer instances for customer %s, err: %v", id, instancesErr)
			continue
		}

		for _, instance := range instances {

			if customerDbUriFlag == "" {
				connectionString, dbConnErr = GetDBConnection(id, instance)
				if dbConnErr != nil {
					log.Printf("Unable to get connection database URI for customer %s, err: %v", id, dbConnErr)
					continue
				}
			} else {
				connectionString = customerDbUriFlag
			}

			pgx, pgxErr := indexer.NewPostgreSQLpgxWithCustomURI(connectionString)
			if pgxErr != nil {
				log.Printf("Error creating RDS connection for customer %s, err: %v", id, pgxErr)
				continue
			}

			if _, ok := customerDBConnections[id]; !ok {
				customerDBConnections[id] = make(map[int]CustomerDBConnection)
			}
			customerDBConnections[id][instance] = CustomerDBConnection{
				Uri: connectionString,
				Pgx: pgx,
			}
		}
		finalCustomerIds = append(finalCustomerIds, id)
	}

	log.Println("Customer IDs to sync:", finalCustomerIds)

	return customerDBConnections, finalCustomerIds, nil
}

func (d *Synchronizer) Start(customerDbUriFlag string, cycleTickerWaitTime int) {
	var isEnd bool

	ticker := time.NewTicker(time.Duration(cycleTickerWaitTime) * time.Second)
	defer ticker.Stop()

	isEnd, err := d.SyncCycle(customerDbUriFlag)
	if err != nil {
		log.Println("Error during initial synchronization cycle:", err)
	}
	if isEnd {
		return
	}

	for {
		select {
		case <-ticker.C:
			isEnd, err := d.SyncCycle(customerDbUriFlag)
			if err != nil {
				log.Fatalf("Error during synchronization cycle:", err)
			}
			if isEnd {
				return
			}
		}
	}
}

func (d *Synchronizer) SyncCycle(customerDbUriFlag string) (bool, error) {
	var isEnd bool

	customerDBConnections, customerIds, customersErr := d.getCustomers(customerDbUriFlag, nil)
	if customersErr != nil {
		return isEnd, customersErr
	}

	// Set start block if 0
	if d.startBlock == 0 {
		startBlock, startErr := d.StartBlockLookUp(customerDBConnections, d.blockchain, crawler.SeerDefaultBlockShift)
		if startErr != nil {
			return isEnd, fmt.Errorf("error determining start block: %w", startErr)
		}
		d.startBlock = startBlock
	}

	// Get the latest block from the indexer db
	indexedLatestBlock, idxLatestErr := indexer.DBConnection.GetLatestDBBlockNumber(d.blockchain, false)
	if idxLatestErr != nil {
		return isEnd, idxLatestErr
	}

	if d.endBlock > 0 && indexedLatestBlock > d.endBlock {
		indexedLatestBlock = d.endBlock
	}

	if d.startBlock > indexedLatestBlock {
		log.Printf("Value in startBlock %d greater then indexedLatestBlock %d, waiting next iteration..", d.startBlock, indexedLatestBlock)
		return isEnd, nil
	}

	// Main loop Steps:
	// 1. Read updates from the indexer db
	// 2. For each update, read the original event data from storage
	// 3. Decode input data using ABIs
	// 4. Write updates to the user RDS
	//var noUpdatesFoundErr *indexer.NoUpdatesFoundError
	var isCycleFinished bool
	for {
		// Check if end block is reached or start block exceeds end block
		if d.endBlock > 0 && d.startBlock >= d.endBlock {
			isEnd = true
			isCycleFinished = true
			log.Printf("End block %d almost reached", d.endBlock)
			break
		}

		if d.endBlock >= indexedLatestBlock {
			isCycleFinished = true
		}

		// Read updates from the indexer db
		// This function will return a list of customer updates 1 update is 1 customer
		_, lastBlockOfChank, path, updates, err := indexer.DBConnection.ReadUpdates(d.blockchain, d.startBlock, customerIds)
		if err != nil {
			return isEnd, fmt.Errorf("error reading updates: %w", err)
		}

		if len(updates) == 0 {
			log.Printf("No updates found for block %d\n", d.startBlock)
			return isEnd, nil
		}

		if crawler.SEER_CRAWLER_DEBUG {
			log.Printf("Read batch key: %s", path)
		}
		log.Println("Last block of current chank: ", lastBlockOfChank)

		// Read the raw data from the storage for current path
		rawData, readErr := d.StorageInstance.Read(path)
		if readErr != nil {
			return isEnd, fmt.Errorf("error reading raw data: %w", readErr)
		}

		log.Printf("Read %d users updates from the indexer db in range of blocks %d-%d\n", len(updates), d.startBlock, lastBlockOfChank)
		var wg sync.WaitGroup

		sem := make(chan struct{}, 5)  // Semaphore to control concurrency
		errChan := make(chan error, 1) // Buffered channel for error handling

		for _, update := range updates {
			for instanceId := range customerDBConnections[update.CustomerID] {
				wg.Add(1)
				go d.processProtoCustomerUpdate(update, rawData, customerDBConnections, instanceId, sem, errChan, &wg)
			}
		}

		wg.Wait()

		close(sem)
		close(errChan) // Close the channel to signal that all goroutines have finished

		// Check for errors from goroutines
		if err := <-errChan; err != nil {
			return isEnd, err
		}

		d.startBlock = lastBlockOfChank + 1

		if isCycleFinished {
			break
		}
	}

	return isEnd, nil
}

func (d *Synchronizer) SyncContracts() {

	indexedLatestBlock, idxLatestErr := indexer.DBConnection.GetLatestDBBlockNumber(d.blockchain)
	if idxLatestErr != nil {
		log.Println("Error getting latest block from the indexer db:", idxLatestErr)
		return
	}

	if d.startBlock == 0 {
		// Get the latest block from the indexer db

		if indexedLatestBlock != 0 {
			d.startBlock = indexedLatestBlock
		} else {
			panic("No start block found in indexer db")
		}

	}

	if d.endBlock == 0 {
		d.endBlock = 1
	}

	log.Printf("Sync contracts for blockchain %s in range of blocks %d-%d\n", d.blockchain, d.startBlock, d.endBlock)
	// Read all paths from database
	paths, err := indexer.DBConnection.GetPaths(d.blockchain, d.startBlock, d.endBlock)
	if err != nil {
		log.Println("Error reading paths from database:", err)
		return
	}

	for _, path := range paths {

		// Read the raw data from the storage for current path
		rawData, readErr := d.StorageInstance.Read(path)

		if readErr != nil {
			log.Println("Error reading events for customers:", readErr)
			return
		}

		log.Printf("Read blocks for path: %s\n", path)

		// get all deployed contracts
		batch, err := d.Client.DecodeProtoEntireBlockToJson(&rawData)

		if err != nil {
			log.Println("Error getting deployed contracts:", err)
			return
		}

		log.Printf("Decoded %d blocks\n", len(batch.Blocks))

		// get all deployed contracts
		contracts, err := seer_blockchain.GetDeployedContracts(d.Client, batch)

		if err != nil {
			log.Println("Error getting deployed contracts:", err)
			return
		}

		if len(contracts) == 0 {
			log.Println("No contracts found in the batch")
			continue
		}

		// get current deployed bytcode

		ctx := context.Background()

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)

		defer cancel()

		for _, contract := range contracts {
			contractBytecode, err := d.Client.GetCode(ctxWithTimeout, common.HexToAddress(contract.Address), indexedLatestBlock)
			if err != nil {
				log.Println("Error getting contract bytecode:", err)
				return
			}

			hexstring := hex.EncodeToString(contractBytecode)

			// convert []bytes to md5 hash of the bytecode
			hashString := fmt.Sprintf("%x", md5.Sum(contractBytecode)) // Convert [16]byte to hex string

			// Assign the hash string to BytecodeHash
			contract.BytecodeHash = hashString

			// Convert bytes to string and assign to Bytecode as a pointer

			contract.Bytecode = &hexstring

		}

		// Write contracts to the bytecode storage

		err = indexer.DBConnection.WriteBytecodeStorage(contracts)

		for _, contract := range contracts {

			bytecodeId, err := indexer.DBConnection.BytecodeIdByHash(contract.BytecodeHash)
			if err != nil {
				log.Println("Error getting bytecode storage id by hash:", err)
				continue
			}

			contract.BytecodeStorageId = &bytecodeId

		}

		// Write contracts to the database

		err = indexer.DBConnection.WriteContracts(d.blockchain, contracts)
		if err != nil {
			log.Println("Error writing contracts to database:", err)
			return
		}

	}

}

func (d *Synchronizer) HistoricalSyncRef(customerDbUriFlag string, addresses []string, customerIds []string, batchSize uint64, autoJobs bool) error {
	var isCycleFinished bool
	var updateDeadline time.Time
	var initialStartBlock uint64

	// Initialize start block if 0
	if d.startBlock == 0 {
		// Get the latest block from the indexer db
		indexedLatestBlock, err := indexer.DBConnection.GetLatestDBBlockNumber(d.blockchain, false)
		if err != nil {
			return fmt.Errorf("error getting latest block number: %w", err)
		}
		d.startBlock = indexedLatestBlock
		fmt.Printf("Start block is %d\n", d.startBlock)
	}

	earlyIndexedBlock, err := indexer.DBConnection.GetLatestDBBlockNumber(d.blockchain, true)

	if err != nil {
		return fmt.Errorf("error getting early indexer block: %w", err)
	}

	// Automatically update ABI jobs as active if auto mode is enabled
	if autoJobs {
		if err := indexer.DBConnection.UpdateAbiJobsStatus(d.blockchain); err != nil {
			return fmt.Errorf("error updating ABI: %w", err)
		}
	}

	// Retrieve customer updates and deployment blocks
	abiJobs, selectJobsErr := indexer.DBConnection.SelectAbiJobs(d.blockchain, addresses, customerIds, autoJobs, true, []string{"function", "event"})
	if selectJobsErr != nil {
		return fmt.Errorf("error selecting ABI jobs: %w", selectJobsErr)
	}
	customerUpdates, addressesAbisInfo, err := indexer.ConvertToCustomerUpdatedAndDeployBlockDicts(abiJobs)
	if err != nil {
		return fmt.Errorf("error parsing ABI jobs: %w", err)
	}

	log.Printf("Found %d customer updates\n", len(customerUpdates))

	// Filter out blocks more
	if autoJobs {
		for address, abisInfo := range addressesAbisInfo {
			log.Printf("Address %s has deployed block %d\n", address, abisInfo.DeployedBlockNumber)

			if abisInfo.DeployedBlockNumber > d.startBlock {
				log.Printf("Finished crawling for address %s at block %d\n", address, abisInfo.DeployedBlockNumber)
				delete(addressesAbisInfo, address)
			}

			if abisInfo.DeployedBlockNumber < earlyIndexedBlock {
				log.Printf("Address %s has deployed block %d less than early indexed block %d\n", address, abisInfo.DeployedBlockNumber, earlyIndexedBlock)
				delete(addressesAbisInfo, address)
			}
		}
	}

	if len(addressesAbisInfo) == 0 {
		log.Println("No addresses to crawl")
		return nil
	}

	// Get customer database connections
	customerIdsMap := make(map[string]bool)
	for _, update := range customerUpdates {
		customerIdsMap[update.CustomerID] = true
	}
	var customerIdsList []string
	for id := range customerIdsMap {
		customerIdsList = append(customerIdsList, id)
	}

	customerDBConnections, _, err := d.getCustomers(customerDbUriFlag, customerIdsList)
	if err != nil {
		return fmt.Errorf("error getting customers: %w", err)
	}

	updateDeadline = time.Now()

	// Main processing loop
	for {

		if autoJobs {

			for address, abisInfo := range addressesAbisInfo {
				if abisInfo.DeployedBlockNumber > d.startBlock {

					log.Printf("Finished crawling for address %s at block %d\n", address, abisInfo.DeployedBlockNumber)

					// update the status of the address for the customer to done
					err := indexer.DBConnection.UpdateAbisAsDone(abisInfo.IDs)
					if err != nil {
						return err
					}

					// drop the address
					delete(addressesAbisInfo, address)
				}

				// Check if the deadline for the update has passed
				if updateDeadline.Add(1 * time.Minute).Before(time.Now()) {
					for address, abisInfo := range addressesAbisInfo {
						ids := abisInfo.IDs

						// calculate progress as a percentage between 0 and 100

						progress := 100 - int(100*(d.startBlock-abisInfo.DeployedBlockNumber)/(d.startBlock-initialStartBlock))

						err := indexer.DBConnection.UpdateAbisProgress(ids, progress)
						if err != nil {
							continue
						}

						log.Printf("Updated progress for address %s to %d%%\n", address, progress)

					}
					updateDeadline = time.Now()

				}
			}

		}

		if len(addressesAbisInfo) == 0 {
			break
		}

		// Determine the processing strategy (RPC or storage)
		var path string
		var firstBlockOfChunk uint64

		for {
			path, firstBlockOfChunk, _, err = indexer.DBConnection.FindBatchPath(d.blockchain, d.startBlock)
			if err != nil {
				return fmt.Errorf("error finding batch path: %w", err)
			}

			if path != "" {
				d.endBlock = firstBlockOfChunk
				break
			}

			log.Printf("No batch path found for block %d, retrying...\n", d.startBlock)
			time.Sleep(30 * time.Second) // Wait for 5 seconds before retrying (adjust the duration as needed)
		}

		// Read raw data from storage or via RPC
		var rawData bytes.Buffer
		rawData, err = d.StorageInstance.Read(path)
		if err != nil {
			return fmt.Errorf("error reading events from storage: %w", err)
		}

		log.Printf("Processing %d customer updates for block range %d-%d", len(customerUpdates), d.startBlock, d.endBlock)

		// Process customer updates in parallel
		var wg sync.WaitGroup
		sem := make(chan struct{}, d.threads) // Semaphore to control concurrency
		errChan := make(chan error, 1)        // Buffered channel for error handling

		for _, update := range customerUpdates {

			for instanceId := range customerDBConnections[update.CustomerID] {
				wg.Add(1)
				go d.processProtoCustomerUpdate(update, rawData, customerDBConnections, instanceId, sem, errChan, &wg)

			}

		}

		wg.Wait()
		close(sem)

		// Check for errors from goroutines
		select {
		case err := <-errChan:
			log.Printf("Error processing customer updates: %v", err)
			return err
		default:
		}

		d.startBlock = d.endBlock - 1

		if isCycleFinished || d.startBlock == 0 {
			if autoJobs {
				for address, abisInfo := range addressesAbisInfo {

					log.Printf("Finished crawling for address %s at block %d\n", address, abisInfo.DeployedBlockNumber)

					// update the status of the address for the customer to done
					err := indexer.DBConnection.UpdateAbisAsDone(abisInfo.IDs)
					if err != nil {
						return err
					}
				}
			}
			log.Println("Finished processing all customer updates")
			break
		}
	}

	return nil
}

func (d *Synchronizer) processProtoCustomerUpdate(
	update indexer.CustomerUpdates,
	rawData bytes.Buffer,
	customerDBConnections map[string]map[int]CustomerDBConnection,
	id int,
	sem chan struct{},
	errChan chan error,
	wg *sync.WaitGroup,
) {
	// Decode input raw proto data using ABIs
	// Write decoded data to the user Database

	defer wg.Done()
	sem <- struct{}{} // Acquire semaphore

	customer, exists := customerDBConnections[update.CustomerID][id]
	if !exists {
		errChan <- fmt.Errorf("no DB connection for customer %s", update.CustomerID)
		<-sem // Release semaphore
		return
	}

	conn, err := customer.Pgx.GetPool().Acquire(context.Background())
	if err != nil {
		errChan <- fmt.Errorf("error acquiring connection for customer %s: %w", update.CustomerID, err)
		<-sem // Release semaphore
		return
	}
	defer conn.Release()
	decodedEvents, decodedTransactions, err := d.Client.DecodeProtoEntireBlockToLabels(&rawData, update.Abis, d.threads)
	if err != nil {
		errChan <- fmt.Errorf("error  %s: %w", update.CustomerID, err)
		<-sem // Release semaphore
		return
	}

	err = customer.Pgx.WriteLabes(d.blockchain, decodedTransactions, decodedEvents)

	if err != nil {
		errChan <- fmt.Errorf("error writing labels for customer %s: %w", update.CustomerID, err)
		<-sem // Release semaphore
		return
	}

	<-sem // Release semaphore
}
