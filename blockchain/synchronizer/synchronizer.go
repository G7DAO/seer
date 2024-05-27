package synchronizer

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/moonstream-to/seer/blockchain/common"
	"github.com/moonstream-to/seer/indexer"
	"github.com/moonstream-to/seer/storage"
)

type Synchronizer struct {
	blockchain  string
	startBlock  uint64
	endBlock    uint64
	providerURI string
}

// NewSynchronizer creates a new synchronizer instance with the given blockchain handler.
func NewSynchronizer(chain string, startBlock uint64, endBlock uint64) *Synchronizer {

	return &Synchronizer{
		blockchain: chain,
		startBlock: startBlock,
	}
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
	req, err := http.NewRequest("GET", MOONSTREAM_DB_MANAGER_API+"/customers/"+uuid+"/instances/1/creds/seer/url", nil)
	if err != nil {
		return "", err
	}

	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+MOONSTREAM_DB_MANAGER_ACCESS_TOKEN)

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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
	abiJobs, err := indexer.Actions.ReadABIJobs(blockchain)
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

func (d *Synchronizer) SyncCustomers() error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Run synchronization cycle...")
			err := d.syncCycle()
			if err != nil {
				fmt.Println("Error during synchronization cycle:", err)
			}
		}
	}
}

func (d *Synchronizer) syncCycle() error {
	// Initialize a wait group to synchronize goroutines
	var wg sync.WaitGroup
	errChan := make(chan error, 1) // Buffered channel for error handling

	// Read client
	client, err := common.NewClient(d.blockchain, "")
	if err != nil {
		log.Println("Error initializing blockchain client:", err)
		log.Fatal(err)
	}

	// Initialize storage
	storageInstance, err := storage.NewStorage()

	if err != nil {
		log.Println("Error initializing storage:", err)
		return err
	}

	// Read ABI jobs from database
	abiJobs, err := d.ReadAbiJobsFromDatabase(d.blockchain)
	if err != nil {
		return err
	}

	// Create a set of customer IDs from ABI jobs to remove duplicates
	customerIdsSet := make(map[string]struct{})
	for _, job := range abiJobs {
		customerIdsSet[job.CustomerID] = struct{}{}
	}

	// Convert set to slice
	var customerIds []string
	for id := range customerIdsSet {
		customerIds = append(customerIds, id)
	}
	log.Println("Customer IDs to sync:", customerIds)

	// Get RDS connections for customer IDs
	rdsConnections, err := GetDBConnections(customerIds)
	if err != nil {
		return err
	}

	if d.startBlock == 0 {

		var latestCustomerBlocks []uint64
		for _, id := range customerIds {
			uri := rdsConnections[id]

			pgx, err := indexer.CreateCustomConnectionToURI(uri)
			if err != nil {
				log.Println("Error creating RDS connection: ", err)
				return err // Error creating RDS connection
			}
			pool := pgx.GetPool()
			fmt.Println("Acquiring pool connection...")

			conn, err := pool.Acquire(context.Background())
			if err != nil {
				log.Println("Error acquiring pool connection: ", err)
				return err // Error acquiring pool connection
			}
			defer conn.Release()

			latestBlock, err := pgx.ReadLastLabel(d.blockchain)
			if err != nil {
				log.Println("Error reading latest block: ", err)
				return err // Error reading the latest block
			}
			latestCustomerBlocks = append(latestCustomerBlocks, latestBlock)
			log.Printf("Latest block for customer %s: %d\n", id, latestBlock)
		}

		// Determine the start block as the maximum of the latest blocks of all customers
		for _, block := range latestCustomerBlocks {
			if block > d.startBlock {
				d.startBlock = block - 100
			}
		}
	}

	// In case start block is still 0, get the latest block from the blockchain
	if d.startBlock == 0 {
		log.Println("Start block is 0, RDS not contain any blocks yet. Sync indexers then.")
		latestBlock, err := indexer.Actions.GetLatestBlockNumber(d.blockchain)
		if err != nil {
			return err
		}
		d.startBlock = latestBlock - 100
		d.endBlock = latestBlock
	}

	// Get the latest block from indexer
	latestBlock, err := indexer.Actions.GetLatestBlockNumber(d.blockchain)

	if err != nil {
		return err
	}
	d.endBlock = latestBlock

	// Main loop Steps:
	// 1. Read updates from the indexer db
	// 2. For each update, read the original event data from storage
	// 3. Decode input data using ABIs
	// 4. Write updates to the user RDS

	for i := d.startBlock; i < d.endBlock; i += 100 {
		endBlock := i + 100

		// Read updates from the indexer db
		// This function will return a list of customer updates 1 update is 1 customer
		updates, err := indexer.Actions.ReadUpdates(d.blockchain, i, endBlock, customerIds)
		if err != nil {
			return fmt.Errorf("error reading updates: %w", err)
		}

		log.Printf("Read %d users updates from the indexer db\n", len(updates))
		log.Printf("Syncing blocks from %d to %d\n", i, endBlock)

		for _, update := range updates {
			wg.Add(1)
			go func(update indexer.CustomerUpdates) {
				defer wg.Done()

				// Get the RDS connection for the customer
				uri := rdsConnections[update.CustomerID]

				// Create a connection to the user RDS
				pgx, err := indexer.CreateCustomConnectionToURI(uri)
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
				encodedEvents, err := storageInstance.ReadBatch(eventsReadMap)

				if err != nil {
					errChan <- fmt.Errorf("error reading events for customer %s: %w", update.CustomerID, err)
					return
				}

				// Union all keys values into a single slice

				var all_events []string

				for _, data := range encodedEvents {
					all_events = append(all_events, data...)
				}

				// Decode the events using ABIs
				decodedEvents, err := client.DecodeProtoEventsToLabels(all_events, update.BlocksCache, update.Abis)

				if err != nil {
					errChan <- fmt.Errorf("error decoding events for customer %s: %w", update.CustomerID, err)
					return
				}

				log.Println("Decoded events amount: ", len(decodedEvents))

				// try to write to user RDS

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

				encodedTransactions, err := storageInstance.ReadBatch(transactionsReadMap)

				if err != nil {
					errChan <- fmt.Errorf("error reading transactions for customer %s: %w", update.CustomerID, err)
					return
				}

				// remap the transactions to a single slice

				var all_transactions []string

				for _, data := range encodedTransactions {
					all_transactions = append(all_transactions, data...)
				}

				decodedTransactions, err := client.DecodeProtoTransactionsToLabels(all_transactions, update.BlocksCache, update.Abis)

				if err != nil {
					errChan <- fmt.Errorf("error decoding transactions for customer %s: %w", update.CustomerID, err)
					return
				}

				log.Println("Decoded transactions amount: ", len(decodedTransactions))

				// Read all rowIds for each path

				// insert decoded labels into the user RDS

				pgx.WriteTransactions(
					d.blockchain,
					decodedTransactions,
				)

				if err != nil {
					errChan <- fmt.Errorf("error reading transactions for customer %s: %w", update.CustomerID, err)
					return
				}

			}(update)
		}
		d.startBlock = latestBlock
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errChan) // Close the channel to signal that all goroutines have finished
	}()

	// Check for errors from goroutines
	for err := range errChan {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	return nil // Return nil to indicate success if no errors occurred
}
