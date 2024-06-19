package synchronizer

import (
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
	"github.com/moonstream-to/seer/crawler"
	"github.com/moonstream-to/seer/indexer"
	"github.com/moonstream-to/seer/storage"
	"golang.org/x/exp/slices"
)

type Synchronizer struct {
	Client          seer_blockchain.BlockchainClient
	StorageInstance storage.Storer

	blockchain string
	startBlock uint64
	endBlock   uint64
	batchSize  uint64
	force      bool
	baseDir    string
	basePath   string
}

// NewSynchronizer creates a new synchronizer instance with the given blockchain handler.
func NewSynchronizer(blockchain, baseDir string, startBlock, endBlock, batchSize uint64, force bool, timeout int) (*Synchronizer, error) {
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

	synchronizer = Synchronizer{
		Client:          client,
		StorageInstance: storageInstance,

		blockchain: blockchain,
		startBlock: startBlock,
		endBlock:   endBlock,
		batchSize:  batchSize,
		baseDir:    baseDir,
		basePath:   basePath,
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

func GetDBConnections(uuids []string) (map[string]string, error) {
	connections := make(map[string]string) // Initialize the map

	for _, id := range uuids {

		connectionString, err := GetDBConnection(id)

		if err != nil {
			return nil, err
		}

		connections[id] = connectionString
	}
	return connections, nil
}

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
	// Simulate reading ABI jobs from the database for a given blockchain.
	// This function will need to interact with a real database or an internal API in the future.
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

// getCustomers fetch ABI jobs, customer IDs and database URLs
func (d *Synchronizer) getCustomers(customerDbUriFlag string) (map[string]string, []string, error) {
	rdsConnections := make(map[string]string)
	var customerIds []string

	// Read ABI jobs from database
	abiJobs, err := d.ReadAbiJobsFromDatabase(d.blockchain)
	if err != nil {
		return nil, customerIds, err
	}

	// Create a set of customer IDs from ABI jobs to remove duplicates
	customerIdsSet := make(map[string]struct{})
	for _, job := range abiJobs {
		customerIdsSet[job.CustomerID] = struct{}{}
	}

	// Convert set to slice
	for id := range customerIdsSet {
		customerIds = append(customerIds, id)
	}
	log.Println("Customer IDs to sync:", customerIds)

	if customerDbUriFlag == "" {
		// Get RDS connections for customer IDs
		rdsConnections, err = GetDBConnections(customerIds)
		if err != nil {
			return nil, customerIds, err
		}
	} else {
		customersLen := 0
		for _, id := range customerIds {
			rdsConnections[id] = customerDbUriFlag
			customersLen++
		}
		log.Printf("For %d customers set one specified db URI with flag", customersLen)
	}

	return rdsConnections, customerIds, nil
}

func (d *Synchronizer) Start(customerDbUriFlag string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			isEnd, err := d.SyncCycle(customerDbUriFlag)
			if err != nil {
				fmt.Println("Error during synchronization cycle:", err)
			}
			if isEnd {
				return
			}
		}
	}
}

func (d *Synchronizer) SyncCycle(customerDbUriFlag string) (bool, error) {
	var isEnd bool
	var wg sync.WaitGroup
	errChan := make(chan error, 1) // Buffered channel for error handling

	rdsConnections, customerIds, customersErr := d.getCustomers(customerDbUriFlag)
	if customersErr != nil {
		return isEnd, customersErr
	}

	if !d.force {
		var latestCustomerBlocks []uint64
		for _, id := range customerIds {
			uri := rdsConnections[id]

			// TODO(kompotkot): Rewrite to not initialize each time new psql conneciton
			pgx, err := indexer.NewPostgreSQLpgxWithCustomURI(uri)
			if err != nil {
				log.Println("Error creating RDS connection: ", err)
				return isEnd, err
			}
			pool := pgx.GetPool()

			conn, err := pool.Acquire(context.Background())
			if err != nil {
				log.Println("Error acquiring pool connection: ", err)
				return isEnd, err
			}
			defer conn.Release()

			latestBlock, err := pgx.ReadLastLabel(d.blockchain)
			if err != nil {
				log.Println("Error reading latest block: ", err)
				return isEnd, err
			}
			latestCustomerBlocks = append(latestCustomerBlocks, latestBlock)
			log.Printf("Latest block for customer %s is: %d\n", id, latestBlock)
		}

		// Determine the start block as the maximum of the latest blocks of all customers
		maxCustomerLatestBlock := slices.Max(latestCustomerBlocks)
		if maxCustomerLatestBlock != 0 {
			d.startBlock = maxCustomerLatestBlock
		}
	}

	// In case start block is still 0, get the latest block from the blockchain minus shift
	if d.startBlock == 0 {
		latestBlockNumber, latestErr := d.Client.GetLatestBlockNumber()
		if latestErr != nil {
			return isEnd, fmt.Errorf("failed to get latest block number: %v", latestErr)
		}
		d.startBlock = uint64(crawler.SetDefaultStartBlock(0, latestBlockNumber))
	}

	// Get the latest block from indexes database
	indexedLatestBlock, idxLatestErr := indexer.DBConnection.GetLatestDBBlockNumber(d.blockchain)
	if idxLatestErr != nil {
		return isEnd, idxLatestErr
	}

	if d.endBlock != 0 && indexedLatestBlock > d.endBlock {
		indexedLatestBlock = d.endBlock
	}

	if d.startBlock >= indexedLatestBlock {
		log.Printf("Value in startBlock %d greater or equal indexedLatestBlock %d, waiting next iteration..", d.startBlock, indexedLatestBlock)
		return isEnd, nil
	}

	// Main loop Steps:
	// 1. Read updates from the indexer db
	// 2. For each update, read the original event data from storage
	// 3. Decode input data using ABIs
	// 4. Write updates to the user RDS
	tempEndBlock := d.startBlock + d.batchSize
	var isCycleFinished bool
	for {
		tempEndBlock = d.startBlock + d.batchSize
		if d.endBlock != 0 {
			if tempEndBlock >= d.endBlock {
				tempEndBlock = d.endBlock
				isEnd = true
				isCycleFinished = true
				log.Printf("End block %d almost reached", tempEndBlock)
			}
		}
		if tempEndBlock >= indexedLatestBlock {
			tempEndBlock = indexedLatestBlock
			isCycleFinished = true
		}

		if crawler.SEER_CRAWLER_DEBUG {
			log.Printf("Syncing %d blocks from %d to %d\n", tempEndBlock-d.startBlock, d.startBlock, tempEndBlock)
		}

		// Read updates from the indexer db
		// This function will return a list of customer updates 1 update is 1 customer
		updates, err := indexer.DBConnection.ReadUpdates(d.blockchain, d.startBlock, tempEndBlock, customerIds)
		if err != nil {
			return isEnd, fmt.Errorf("error reading updates: %w", err)
		}

		log.Printf("Read %d users updates from the indexer db\n", len(updates))

		for _, update := range updates {
			wg.Add(1)
			go func(update indexer.CustomerUpdates) {
				defer wg.Done()

				// Get the RDS connection for the customer
				uri := rdsConnections[update.CustomerID]

				// Create a connection to the user RDS
				// TODO(kompotkot): Rewrite to not initialize each time new psql conneciton
				pgx, err := indexer.NewPostgreSQLpgxWithCustomURI(uri)
				if err != nil {
					errChan <- fmt.Errorf("error creating connection to RDS for customer %s: %w", update.CustomerID, err)
					return
				}

				pool := pgx.GetPool()
				conn, err := pool.Acquire(context.Background())
				if err != nil {
					errChan <- fmt.Errorf("error acquiring connection for customer %s: %w", update.CustomerID, err)
					return
				}
				defer conn.Release()

				// Events group by path for batch reading
				groupByPathEvents := make(map[string][]uint64)

				for _, event := range update.Data.Events {

					if _, ok := groupByPathEvents[event.Path]; !ok {
						groupByPathEvents[event.Path] = []uint64{}
					}

					groupByPathEvents[event.Path] = append(groupByPathEvents[event.Path], event.RowID)
				}

				eventsReadMap := []storage.ReadItem{}

				for path, rowIds := range groupByPathEvents {
					eventsReadMap = append(eventsReadMap, storage.ReadItem{
						Key:    path,
						RowIds: rowIds,
					})
				}
				// Read all rowIds for each path
				encodedEvents, err := d.StorageInstance.ReadBatch(eventsReadMap)

				if err != nil {
					fmt.Println("Error reading events: ", err)
					errChan <- fmt.Errorf("error reading events for customer %s: %w", update.CustomerID, err)
					return
				}

				// Union all keys values into a single slice
				var all_events []string

				for _, data := range encodedEvents {
					all_events = append(all_events, data...)
				}

				// Decode the events using ABIs
				decodedEvents, err := d.Client.DecodeProtoEventsToLabels(all_events, update.BlocksCache, update.Abis)

				if err != nil {
					fmt.Println("Error decoding events: ", err)
					errChan <- fmt.Errorf("error decoding events for customer %s: %w", update.CustomerID, err)
					return
				}

				// Write events to user RDS
				pgx.WriteEvents(
					d.blockchain,
					decodedEvents,
				)

				// Transactions
				groupByPathTransactions := make(map[string][]uint64)

				for _, transaction := range update.Data.Transactions {

					if _, ok := groupByPathTransactions[transaction.Path]; !ok {
						groupByPathTransactions[transaction.Path] = []uint64{}
					}

					groupByPathTransactions[transaction.Path] = append(groupByPathTransactions[transaction.Path], transaction.RowID)
				}

				transactionsReadMap := []storage.ReadItem{}

				for path, rowIds := range groupByPathTransactions {
					transactionsReadMap = append(transactionsReadMap, storage.ReadItem{
						Key:    path,
						RowIds: rowIds,
					})
				}
				encodedTransactions, err := d.StorageInstance.ReadBatch(transactionsReadMap)

				if err != nil {
					errChan <- fmt.Errorf("error reading transactions for customer %s: %w", update.CustomerID, err)
					return
				}

				// remap the transactions to a single slice

				var all_transactions []string

				for _, data := range encodedTransactions {
					all_transactions = append(all_transactions, data...)
				}

				decodedTransactions, err := d.Client.DecodeProtoTransactionsToLabels(all_transactions, update.BlocksCache, update.Abis)

				if err != nil {
					errChan <- fmt.Errorf("error decoding transactions for customer %s: %w", update.CustomerID, err)
					return
				}

				// Write transactions to user RDS
				pgx.WriteTransactions(
					d.blockchain,
					decodedTransactions,
				)

			}(update)
		}

		wg.Wait()

		if isCycleFinished {
			break
		}

		d.startBlock = tempEndBlock + 1
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errChan) // Close the channel to signal that all goroutines have finished
	}()

	// Check for errors from goroutines
	for err := range errChan {
		fmt.Println("Error during synchronization cycle:", err)
		if err != nil {
			return isEnd, err
		}
	}

	return isEnd, nil
}
