package synchronizer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	seer_blockchain "github.com/moonstream-to/seer/blockchain"
	"github.com/moonstream-to/seer/blockchain/common"
	"github.com/moonstream-to/seer/crawler"
	"github.com/moonstream-to/seer/indexer"
	"github.com/moonstream-to/seer/storage"
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

func GetDBConnection(uuid string) (string, error) {

	// Create the request
	req, err := http.NewRequest("GET", MOONSTREAM_DB_V3_CONTROLLER_API+"/customers/"+uuid+"/instances/1/creds/seer/url", nil)
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
		log.Println("Failed to get connection string for id", uuid, ":", MOONSTREAM_DB_V3_CONTROLLER_API+"/customers/"+uuid+"/instances/1/creds/seer/url")
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

func CustomersLatestBlocks(customerDBConnections map[string]CustomerDBConnection, blockchain string) (map[string]uint64, error) {
	latestBlocks := make(map[string]uint64)
	for id, customer := range customerDBConnections {
		pool := customer.Pgx.GetPool()
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Println("Error acquiring pool connection: ", err)
			return nil, err
		}
		defer conn.Release()

		latestLabelBlock, err := customer.Pgx.ReadLastLabel(blockchain)
		if err != nil {
			log.Println("Error reading latest block: ", err)
			return nil, err
		}
		latestBlocks[id] = latestLabelBlock
		log.Printf("Latest block for customer %s is: %d\n", id, latestLabelBlock)
	}
	return latestBlocks, nil
}

func (d *Synchronizer) StartBlockLookUp(customerDBConnections map[string]CustomerDBConnection, blockchain string, blockShift int64) (uint64, error) {
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
func (d *Synchronizer) getCustomers(customerDbUriFlag string, customerIds []string) (map[string]CustomerDBConnection, []string, error) {
	customerDBConnections := make(map[string]CustomerDBConnection)
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

		if customerDbUriFlag == "" {
			connectionString, dbConnErr = GetDBConnection(id)
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

		customerDBConnections[id] = CustomerDBConnection{
			Uri: connectionString,
			Pgx: pgx,
		}
		finalCustomerIds = append(finalCustomerIds, id)
	}

	log.Println("Customer IDs to sync:", finalCustomerIds)

	return customerDBConnections, finalCustomerIds, nil
}

func (d *Synchronizer) Start(customerDbUriFlag string) {
	var isEnd bool

	ticker := time.NewTicker(10 * time.Second)
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
	indexedLatestBlock, idxLatestErr := indexer.DBConnection.GetLatestDBBlockNumber(d.blockchain)
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
			return isEnd, fmt.Errorf("error reading events for customers %s: %w", readErr)
		}

		log.Printf("Read %d users updates from the indexer db in range of blocks %d-%d\n", len(updates), d.startBlock, lastBlockOfChank)
		var wg sync.WaitGroup

		sem := make(chan struct{}, 5)  // Semaphore to control concurrency
		errChan := make(chan error, 1) // Buffered channel for error handling

		for _, update := range updates {
			wg.Add(1)
			go d.processProtoCustomerUpdate(update, rawData, customerDBConnections, sem, errChan, &wg)
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

func (d *Synchronizer) HistoricalSyncRef(customerDbUriFlag string, addresses []string, customerIds []string, batchSize uint64, auto bool) error {
	var useRPC bool
	var isCycleFinished bool

	// Initialize start block if 0
	if d.startBlock == 0 {
		// Get the latest block from the indexer db
		indexedLatestBlock, err := indexer.DBConnection.GetLatestDBBlockNumber(d.blockchain)
		if err != nil {
			return fmt.Errorf("error getting latest block number: %w", err)
		}
		d.startBlock = indexedLatestBlock
		fmt.Printf("Start block is %d\n", d.startBlock)
	}

	// Automatically update ABI jobs as active if auto mode is enabled
	if auto {
		if err := indexer.DBConnection.UpdateAbiJobsStatus(d.blockchain); err != nil {
			return fmt.Errorf("error updating ABI: %w", err)
		}
	}

	// Retrieve customer updates and deployment blocks
	customerUpdates, addressesAbisInfo, err := indexer.DBConnection.SelectAbiJobs(d.blockchain, addresses, customerIds)
	if err != nil {
		return fmt.Errorf("error selecting ABI jobs: %w", err)
	}

	fmt.Printf("Found %d customer updates\n", len(customerUpdates))

	// Filter out blocks more
	for address, abisInfo := range addressesAbisInfo {
		fmt.Printf("Address %s, block %d\n", address, abisInfo.DeployedBlockNumber)
		if abisInfo.DeployedBlockNumber > d.startBlock {
			fmt.Printf("Deleting address %s\n", address, abisInfo.DeployedBlockNumber)
			delete(addressesAbisInfo, address)
		}
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

	// Main processing loop
	for {
		for address, abisInfo := range addressesAbisInfo {
			if abisInfo.DeployedBlockNumber > d.startBlock {

				// update the status of the address for the customer to done
				err := indexer.DBConnection.UpdateAbisAsDone(abisInfo.IDs)
				if err != nil {
					return err
				}

				// drop the address
				delete(addressesAbisInfo, address)
			}
		}

		if len(addressesAbisInfo) == 0 {
			break
		}

		// Determine the processing strategy (RPC or storage)
		var path string
		var firstBlockOfChunk uint64
		if !useRPC {

			path, firstBlockOfChunk, _, err = indexer.DBConnection.FindBatchPath(d.blockchain, d.startBlock)
			if err != nil {
				return fmt.Errorf("error finding batch path: %w", err)
			}

			if path == "" {
				useRPC = true
			} else {
				d.endBlock = firstBlockOfChunk
			}
		}

		if useRPC {
			// Calculate the end block
			if d.startBlock > batchSize {
				d.endBlock = d.startBlock - batchSize
			} else {
				d.endBlock = 1 // Prevent underflow by setting endBlock to 1 if startBlock < batchSize
			}

			fmt.Printf("Start block is %d, end block is %d\n", d.startBlock, d.endBlock)
		}

		// Read raw data from storage or via RPC
		var rawData bytes.Buffer
		if useRPC {
			// protoMessage, _, _, err := seer_blockchain.CrawlEntireBlocks(d.Client, big.NewInt(int64(d.endBlock)), big.NewInt(int64(d.startBlock)), false, 5)
			// if err != nil {
			// 	return fmt.Errorf("error reading events via RPC: %w", err)
			// }

			// blocksBatch, err := d.Client.ProcessBlocksToBatch(protoMessage)
			// if err != nil {
			// 	return fmt.Errorf("error processing blocks to batch: %w", err)
			// }

			// dataBytes, err := proto.Marshal(blocksBatch)
			// if err != nil {
			// 	return fmt.Errorf("error marshaling protoMessage: %w", err)
			// }

			// rawData = *bytes.NewBuffer(dataBytes)

			// Read transactions and events from the blockchain

		} else {
			rawData, err = d.StorageInstance.Read(path)
			if err != nil {
				return fmt.Errorf("error reading events from storage: %w", err)
			}
		}

		log.Printf("Processing %d customer updates for block range %d-%d", len(customerUpdates), d.startBlock, d.endBlock)

		// Process customer updates in parallel
		var wg sync.WaitGroup
		sem := make(chan struct{}, d.threads) // Semaphore to control concurrency
		errChan := make(chan error, 1)        // Buffered channel for error handling

		for _, update := range customerUpdates {
			wg.Add(1)
			if useRPC {
				go d.processRPCCustomerUpdate(update, customerDBConnections, sem, errChan, &wg)
			} else {
				go d.processProtoCustomerUpdate(update, rawData, customerDBConnections, sem, errChan, &wg)
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
			log.Println("Finished processing all customer updates")
			break
		}
	}

	return nil
}

func (d *Synchronizer) processProtoCustomerUpdate(
	update indexer.CustomerUpdates,
	rawData bytes.Buffer,
	customerDBConnections map[string]CustomerDBConnection,
	sem chan struct{},
	errChan chan error,
	wg *sync.WaitGroup,
) {
	// Decode input raw proto data using ABIs
	// Write decoded data to the user Database

	defer wg.Done()
	sem <- struct{}{} // Acquire semaphore

	customer, exists := customerDBConnections[update.CustomerID]
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

func (d *Synchronizer) processRPCCustomerUpdate(
	update indexer.CustomerUpdates,
	customerDBConnections map[string]CustomerDBConnection,
	sem chan struct{},
	errChan chan error,
	wg *sync.WaitGroup,
) {
	// Decode input raw proto data using ABIs
	// Write decoded data to the user Database

	defer wg.Done()
	sem <- struct{}{} // Acquire semaphore

	customer, exists := customerDBConnections[update.CustomerID]
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

	// split abis by the type of the task

	eventAbis := make(map[string]map[string]*indexer.AbiEntry)
	transactionAbis := make(map[string]map[string]*indexer.AbiEntry)

	for address, selectorMap := range update.Abis {

		for selector, abiInfo := range selectorMap {

			if abiInfo.AbiType == "event" {
				if _, ok := eventAbis[address]; !ok {
					eventAbis[address] = make(map[string]*indexer.AbiEntry)
				}
				eventAbis[address][selector] = abiInfo
			} else {
				if _, ok := transactionAbis[address]; !ok {
					transactionAbis[address] = make(map[string]*indexer.AbiEntry)
				}
				transactionAbis[address][selector] = abiInfo
			}
		}

	}

	// Transactions

	var transactions []indexer.TransactionLabel
	var blocksCache map[uint64]common.BlockWithTransactions

	if len(transactionAbis) != 0 {
		fmt.Printf("Getting transactions for customer %s\n", update.CustomerID)
		transactions, blocksCache, err = d.Client.GetTransactionsLabels(d.endBlock, d.startBlock, transactionAbis, d.threads)

		if err != nil {
			log.Println("Error getting transactions for customer", update.CustomerID, ":", err)
			errChan <- fmt.Errorf("error getting transactions for customer %s: %w", update.CustomerID, err)
			<-sem // Release semaphore
			return
		}
	} else {
		transactions = make([]indexer.TransactionLabel, 0)
	}

	fmt.Printf("Transactions length: %d\n", len(transactions))

	// Events

	var events []indexer.EventLabel

	if len(eventAbis) != 0 {

		events, err = d.Client.GetEventsLabels(d.endBlock, d.startBlock, eventAbis, blocksCache)

		if err != nil {

			log.Println("Error getting events for customer", update.CustomerID, ":", err)

			errChan <- fmt.Errorf("error getting events for customer %s: %w", update.CustomerID, err)
			<-sem // Release semaphore
			return

		}
	} else {
		events = make([]indexer.EventLabel, 0)
	}

	fmt.Printf("Events length: %d\n", len(events))

	err = customer.Pgx.WriteLabes(d.blockchain, transactions, events)

	if err != nil {
		errChan <- fmt.Errorf("error writing labels for customer %s: %w", update.CustomerID, err)
		<-sem // Release semaphore
		return
	}

	<-sem // Release semaphore
}
