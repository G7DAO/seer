package synchronizer

import (
	"bytes"
	"context"
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
)

type Synchronizer struct {
	Client          seer_blockchain.BlockchainClient
	StorageInstance storage.Storer

	blockchain         string
	startBlock         uint64
	endBlock           uint64
	batchSize          uint64
	baseDir            string
	basePath           string
	threads            int
	minBlocksToSync    int
	addRawTransactions bool
}

// NewSynchronizer creates a new synchronizer instance with the given blockchain handler.
func NewSynchronizer(blockchain, rpcUrl, baseDir string, startBlock, endBlock, batchSize uint64, timeout int, threads int, minBlocksToSync int, addRawTransactions bool) (*Synchronizer, error) {
	var synchronizer Synchronizer

	basePath := filepath.Join(baseDir, crawler.SeerCrawlerStoragePrefix, "data", blockchain)
	storageInstance, err := storage.NewStorage(storage.SeerCrawlerStorageType, basePath)
	if err != nil {
		log.Fatalf("Failed to create storage instance: %v", err)
		panic(err)
	}

	client, err := seer_blockchain.NewClient(blockchain, rpcUrl, timeout)
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

		blockchain:         blockchain,
		startBlock:         startBlock,
		endBlock:           endBlock,
		batchSize:          batchSize,
		baseDir:            baseDir,
		basePath:           basePath,
		threads:            threads,
		minBlocksToSync:    minBlocksToSync,
		addRawTransactions: addRawTransactions,
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

func GetDBConnection(uuid string, id int, dbUsername string) (string, error) {
	// Define timeout duration
	ctx, cancel := context.WithTimeout(context.Background(), SEER_MDB_V3_CONTROLLER_API_TIMEOUT)
	defer cancel()

	uri := fmt.Sprintf("%s/customers/%s/instances/%s/creds/%s/url", MOONSTREAM_DB_V3_CONTROLLER_API, uuid, strconv.Itoa(id), dbUsername)

	// Create the request with the context
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return "", err
	}

	// Set the authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN))

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("request timed out: %w", err)
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get connection from uri %s", uri)
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
	// Define timeout duration
	ctx, cancel := context.WithTimeout(context.Background(), SEER_MDB_V3_CONTROLLER_API_TIMEOUT)
	defer cancel()

	// Create the request with the context
	req, err := http.NewRequestWithContext(ctx, "GET", MOONSTREAM_DB_V3_CONTROLLER_API+"/customers/"+uuid, nil)
	if err != nil {
		return nil, err
	}

	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN)

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("request timed out: %w", err)
		}
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

func (d *Synchronizer) CloseIndexerDBConnections(customerDBConnections map[string]map[int]CustomerDBConnection) {
	for _, customer := range customerDBConnections {
		for _, instance := range customer {
			instance.Pgx.Close()
		}
	}
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
			log.Println("Skipping customer", id)
			continue
		}

		for _, instance := range instances {

			if customerDbUriFlag == "" {
				connectionString, dbConnErr = GetDBConnection(id, instance, "seer")
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

	// close the indexer db connection
	defer d.CloseIndexerDBConnections(customerDBConnections)

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

	if d.startBlock > indexedLatestBlock && d.endBlock == 0 {
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
		_, lastBlockOfChank, paths, updates, err := indexer.DBConnection.ReadUpdates(d.blockchain, d.startBlock, customerIds, d.minBlocksToSync)
		if err != nil {
			return isEnd, fmt.Errorf("error reading updates: %w", err)
		}

		if len(updates) == 0 {
			log.Printf("No updates found for block %d\n", d.startBlock)
			return isEnd, nil
		}

		if crawler.SEER_CRAWLER_DEBUG {
			log.Printf("Read batch key: %s", paths)
		}
		log.Println("Last block of current chank: ", lastBlockOfChank)

		// Read the raw data from the storage for current path
		rawData, readErr := storage.ReadFilesAsync(paths, d.threads, d.StorageInstance)
		if readErr != nil {
			return isEnd, fmt.Errorf("error reading raw data: %w", readErr)
		}

		log.Printf("Read %d users updates from the indexer db in range of blocks %d-%d\n", len(updates), d.startBlock, lastBlockOfChank)
		// Process customer updates in parallel
		var wg sync.WaitGroup

		// count the number of goroutines that will be running
		var totalGoroutines int
		for _, update := range updates {
			totalGoroutines += len(customerDBConnections[update.CustomerID])
		}

		sem := make(chan struct{}, d.threads)        // Semaphore to control concurrency
		errChan := make(chan error, totalGoroutines) // Channel to collect errors from goroutines

		var errs []error
		for _, update := range updates {
			for instanceId := range customerDBConnections[update.CustomerID] {
				wg.Add(1)
				go d.processProtoCustomerUpdate(update, rawData, customerDBConnections, instanceId, sem, errChan, &wg, d.addRawTransactions)
			}
		}

		wg.Wait()
		close(errChan) // Close the channel to signal that all goroutines have finished

		// Check if there were any errors
		for err := range errChan {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			var errMsg string
			for _, e := range errs {
				errMsg += e.Error() + "\n"
			}

			return isEnd, fmt.Errorf("errors processing customer updates:\n%s", errMsg)
		}

		d.startBlock = lastBlockOfChank + 1

		if isCycleFinished {
			break
		}
		// close the indexer db connection
	}

	return isEnd, nil
}

func (d *Synchronizer) HistoricalSyncRef(customerDbUriFlag string, addresses []string, customerIds []string, batchSize uint64, autoJobs bool) error {
	var isCycleFinished bool
	var updateDeadline time.Time
	var initialStartBlock uint64
	var initialEndBlock uint64

	initialEndBlock = d.endBlock

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

	fmt.Printf("Found %d customer updates\n", len(customerUpdates))

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

	// close the indexer db connection
	defer d.CloseIndexerDBConnections(customerDBConnections)

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
			log.Println("No addresses to crawl")
			break
		}

		if d.startBlock <= initialEndBlock {
			log.Println("Targeted end block reached finish history sync")
			break
		}

		// Determine the processing strategy (RPC or storage)
		var paths []string
		var firstBlockOfChunk uint64

		for {
			paths, firstBlockOfChunk, _, err = indexer.DBConnection.RetrievePathsAndBlockBounds(d.blockchain, d.startBlock, d.minBlocksToSync)
			if err != nil {
				return fmt.Errorf("error finding batch path: %w", err)
			}

			if paths != nil {
				if d.endBlock <= firstBlockOfChunk {
					break
				}
			}

			log.Printf("No batch path found for block %d, finish history sync", d.startBlock)
			// as older block is not available, we finish history sync
			return nil
		}

		// Read raw data from storage or via RPC
		var rawData []bytes.Buffer
		rawData, err = storage.ReadFilesAsync(paths, d.threads, d.StorageInstance)
		if err != nil {
			return fmt.Errorf("error reading events from storage: %w", err)
		}

		log.Printf("Processing %d customer updates for block range %d-%d", len(customerUpdates), d.startBlock, d.endBlock)

		// Process customer updates in parallel
		var wg sync.WaitGroup

		// count the number of goroutines that will be running
		var totalGoroutines int
		for _, update := range customerUpdates {
			totalGoroutines += len(customerDBConnections[update.CustomerID])
		}

		sem := make(chan struct{}, d.threads)        // Semaphore to control concurrency
		errChan := make(chan error, totalGoroutines) // Channel to collect errors from goroutines

		var errs []error

		for _, update := range customerUpdates {

			for instanceId := range customerDBConnections[update.CustomerID] {
				wg.Add(1)
				go d.processProtoCustomerUpdate(update, rawData, customerDBConnections, instanceId, sem, errChan, &wg, d.addRawTransactions)
			}

		}

		wg.Wait()
		close(errChan) // Close the channel to signal that all goroutines have finished

		// Check if there were any errors
		for err := range errChan {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			var errMsg string
			for _, e := range errs {
				errMsg += e.Error() + "\n"
			}
			return fmt.Errorf("errors processing customer updates:\n%s", errMsg)
		}

		d.startBlock = d.endBlock - 1

		fmt.Printf("Processed %d customer updates for block range %d-%d\n", len(customerUpdates), d.startBlock, d.endBlock)

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
	rawDataList []bytes.Buffer,
	customerDBConnections map[string]map[int]CustomerDBConnection,
	id int,
	sem chan struct{},
	errChan chan error,
	wg *sync.WaitGroup,
	addRawTransactions bool,
) {
	// Decode input raw proto data using ABIs
	// Write decoded data to the user Database

	defer func() {
		if r := recover(); r != nil {
			errChan <- fmt.Errorf("panic in goroutine for customer %s, instance %d: %v", update.CustomerID, id, r)
		}
		wg.Done()
	}()

	sem <- struct{}{}        // Acquire semaphore
	defer func() { <-sem }() // Release semaphore

	customer, exists := customerDBConnections[update.CustomerID][id]
	if !exists {
		errChan <- fmt.Errorf("no DB connection for customer %s", update.CustomerID)
		return
	}

	conn, err := customer.Pgx.GetPool().Acquire(context.Background())
	if err != nil {
		errChan <- fmt.Errorf("error acquiring connection for customer %s: %w", update.CustomerID, err)
		return
	}
	defer conn.Release()

	var listDecodedEvents []indexer.EventLabel
	var listDecodedTransactions []indexer.TransactionLabel
	var listDecodedRawTransactions []indexer.RawTransaction
	for _, rawData := range rawDataList {
		// Decode the raw data to transactions
		decodedEvents, decodedTransactions, decodedRawTransactions, err := d.Client.DecodeProtoEntireBlockToLabels(&rawData, update.Abis, addRawTransactions, d.threads)

		listDecodedEvents = append(listDecodedEvents, decodedEvents...)
		listDecodedTransactions = append(listDecodedTransactions, decodedTransactions...)
		listDecodedRawTransactions = append(listDecodedRawTransactions, decodedRawTransactions...)
		if err != nil {
			errChan <- fmt.Errorf("error decoding data for customer %s: %w", update.CustomerID, err)
			return
		}

	}
	// make retrying
	retry := 0
	for {
		err = customer.Pgx.WriteDataToCustomerDB(d.blockchain, listDecodedTransactions, listDecodedEvents, listDecodedRawTransactions)

		if err != nil {
			retry++
			if retry > 3 {
				errChan <- fmt.Errorf("error writing labels for customer %s: %w", update.CustomerID, err)
				return
			}
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
}
