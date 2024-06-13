package crawler

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"path/filepath"
	"time"

	"github.com/moonstream-to/seer/blockchain"
	"github.com/moonstream-to/seer/indexer"
	"github.com/moonstream-to/seer/storage"
	"google.golang.org/protobuf/proto"
)

// Crawler defines the crawler structure.
type Crawler struct {
	blockchain    string
	startBlock    uint64
	blocksBatch   uint64
	confirmations uint64
	endBlock      uint64
	force         bool
	baseDir       string
}

// NewCrawler creates a new crawler instance with the given blockchain handler.
func NewCrawler(chain string, startBlock, endBlock, batchSize, confirmations uint64, timeout int, baseDir string, force bool) *Crawler {
	log.Printf("Initialized new crawler at chain: %s, startBlock: %d, endBlock: %d, force: %t", chain, startBlock, endBlock, force)
	return &Crawler{
		blockchain:    chain,
		startBlock:    startBlock,
		endBlock:      endBlock,
		blocksBatch:   batchSize,
		confirmations: confirmations,
		baseDir:       baseDir,
		force:         force,
	}

}

// Utility function to handle retries
func retryOperation(attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		if err := fn(); err != nil {
			if i == attempts-1 {
				return err
			}
			log.Printf("Attempt %d/%d failed: %v. Retrying in %s...", i+1, attempts, err, sleep)
			time.Sleep(sleep)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed after %d attempts", attempts)
}

func updateBlockIndexFilepaths(indices []indexer.BlockIndex, basePath, batchDir string) {
	for i := range indices {
		indices[i].Path = filepath.Join(basePath, batchDir, "blocks.proto")
	}
}

func updateTransactionIndexFilepaths(indices []indexer.TransactionIndex, basePath, batchDir string) {
	for i := range indices {
		indices[i].Path = filepath.Join(basePath, batchDir, "transactions.proto")
	}
}

func UpdateEventIndexFilepaths(indices []indexer.LogIndex, basePath, batchDir string) {
	for i := range indices {
		indices[i].Path = filepath.Join(basePath, batchDir, "logs.proto")
	}
}

func SetDefaultStartBlock(confirmations uint64, latestBlockNumber *big.Int) uint64 {
	startBlock := latestBlockNumber.Uint64() - confirmations - 100
	log.Printf("Start block set with shift: %d\n", startBlock)
	return startBlock
}

// Start initiates the crawling process for the configured blockchain.
func (c *Crawler) Start() {
	// Initialize the storage instance

	basePath := filepath.Join(c.baseDir, SeerCrawlerStoragePrefix, "data", c.blockchain)
	storageInstance, err := storage.NewStorage(storage.SeerCrawlerStorageType, basePath)
	if err != nil {
		log.Fatalf("Failed to create storage instance: %v", err)
		panic(err)
	}

	client, err := blockchain.NewClient(c.blockchain, BlockchainURLs[c.blockchain])
	if err != nil {
		log.Fatal(err)
	}

	latestBlockNumber, err := client.GetLatestBlockNumber()
	if err != nil {
		log.Fatalf("Failed to get latest block number: %v", err)
	}

	log.Printf("Latest block number at blockchain: %d\n", latestBlockNumber)

	if c.force {
		// Start form specified startBlock if it is set and not 0
		if c.startBlock == 0 {
			c.startBlock = SetDefaultStartBlock(c.confirmations, latestBlockNumber)
		}
	} else {
		latestIndexedBlock, err := indexer.DBConnection.GetLatestBlockNumber(c.blockchain)

		// If there are no rows in result then set startBlock with SetDefaultStartBlock()

		if err != nil {
			if err.Error() == "no rows in result set" {
				c.startBlock = SetDefaultStartBlock(c.confirmations, latestBlockNumber)
			} else {
				log.Fatalf("Failed to get latest indexed block: %v", err)
			}

		}

		if latestIndexedBlock != 0 {
			c.startBlock = latestIndexedBlock + 1
			log.Printf("Start block fetched from indexes database and set to: %d\n", c.startBlock)
		}
	}

	tempEndBlock := uint64(0)
	safeBlock := uint64(0)

	retryWaitTime := 10 * time.Second
	waitForBlocksTime := retryWaitTime
	maxWaitForBlocksTime := 12 * retryWaitTime
	retryAttempts := 3

	// Crawling loop
	for {

		latestBlockNumber, err = client.GetLatestBlockNumber()

		if err != nil {
			log.Fatalf("Failed to get latest block number: %v", err)
			// Retry the operation
			time.Sleep(retryWaitTime)
			retryAttempts--
			if retryAttempts == 0 {
				log.Fatalf("Failed to get latest block number after %d attempts", retryAttempts)
			}
			continue
		}

		safeBlock = latestBlockNumber.Uint64() - c.confirmations

		tempEndBlock = c.startBlock + c.blocksBatch

		// TODO: Rewrite logic to not wait until we will get batch of blocks, but crawl it even it lower then batch number
		if tempEndBlock > safeBlock {
			// Auto adjust time
			log.Printf("Waiting for new blocks to be mined. Current latestBlockNumber: %d, safeBlock: %d", latestBlockNumber, safeBlock)
			time.Sleep(waitForBlocksTime)
			if waitForBlocksTime < maxWaitForBlocksTime {
				waitForBlocksTime = waitForBlocksTime * 2
			}
			continue
		}
		waitForBlocksTime = retryWaitTime

		if c.endBlock != 0 && tempEndBlock > c.endBlock {
			log.Printf("End block reached at: %d", tempEndBlock)
			break
		}

		// Retry the operation in case of failure with cumulative attempts
		err = retryOperation(retryAttempts, retryWaitTime, func() error {
			blocks, transactions, blockIndex, transactionIndex, blocksCache, err := blockchain.CrawlBlocks(client, big.NewInt(int64(c.startBlock)), big.NewInt(int64(tempEndBlock)))
			if err != nil {
				return fmt.Errorf("failed to crawl blocks: %w", err)
			}

			batchDir := fmt.Sprintf("%d-%d", c.startBlock, tempEndBlock)
			log.Printf("Operates with batch of blocks: %s", batchDir)

			var encodedBytesBlocks []string
			for _, block := range blocks {
				bytesBlocks, err := proto.Marshal(block)
				if err != nil {
					return fmt.Errorf("failed to marshal proto message: %w", err)
				}
				// Encode the bytes to base64 for newline delimited storage
				base64Block := base64.StdEncoding.EncodeToString(bytesBlocks)
				encodedBytesBlocks = append(encodedBytesBlocks, base64Block)
			}

			if err := storageInstance.Save(batchDir, "blocks.proto", encodedBytesBlocks); err != nil {
				return fmt.Errorf("failed to save blocks: %w", err)
			}
			log.Printf("Saved blocks.proto to %s", batchDir)

			var encodedBytesTransactions []string
			for _, transaction := range transactions {
				bytesTransaction, err := proto.Marshal(transaction)
				if err != nil {
					return fmt.Errorf("failed to marshal proto message: %w", err)
				}
				// Encode the bytes to base64 for newline delimited storage
				base64Transaction := base64.StdEncoding.EncodeToString(bytesTransaction)
				encodedBytesTransactions = append(encodedBytesTransactions, base64Transaction)
			}

			if err := storageInstance.Save(batchDir, "transactions.proto", encodedBytesTransactions); err != nil {
				return fmt.Errorf("failed to save transactions: %w", err)
			}
			log.Printf("Saved transactions.proto to %s", batchDir)

			updateBlockIndexFilepaths(blockIndex, basePath, batchDir)
			interfaceBlockIndex := make([]interface{}, len(blockIndex))
			for i, v := range blockIndex {
				interfaceBlockIndex[i] = v
			}

			if err := indexer.WriteIndexesToDatabase(c.blockchain, interfaceBlockIndex, "block"); err != nil {
				return fmt.Errorf("failed to write block index to database: %w", err)
			}

			updateTransactionIndexFilepaths(transactionIndex, basePath, batchDir)
			interfaceTransactionIndex := make([]interface{}, len(transactionIndex))
			for i, v := range transactionIndex {
				interfaceTransactionIndex[i] = v
			}

			if err := indexer.WriteIndexesToDatabase(c.blockchain, interfaceTransactionIndex, "transaction"); err != nil {
				return fmt.Errorf("failed to write transaction index to database: %w", err)
			}

			events, eventsIndex, err := blockchain.CrawlEvents(client, big.NewInt(int64(c.startBlock)), big.NewInt(int64(tempEndBlock)), blocksCache)
			if err != nil {
				return fmt.Errorf("failed to crawl events: %w", err)
			}

			UpdateEventIndexFilepaths(eventsIndex, basePath, batchDir)
			interfaceEventsIndex := make([]interface{}, len(eventsIndex))
			for i, v := range eventsIndex {
				interfaceEventsIndex[i] = v
			}

			if err := indexer.WriteIndexesToDatabase(c.blockchain, interfaceEventsIndex, "log"); err != nil {
				return fmt.Errorf("failed to write event index to database: %w", err)
			}

			var encodedBytesEvents []string
			for _, event := range events {
				bytesEvent, err := proto.Marshal(event)
				if err != nil {
					return fmt.Errorf("failed to marshal proto message: %w", err)
				}
				// Encode the bytes to base64 for newline delimited storage
				base64Event := base64.StdEncoding.EncodeToString(bytesEvent)
				encodedBytesEvents = append(encodedBytesEvents, base64Event)
			}

			if err := storageInstance.Save(batchDir, "logs.proto", encodedBytesEvents); err != nil {
				return fmt.Errorf("failed to save events: %w", err)
			}
			log.Printf("Saved logs.proto to %s", batchDir)

			return nil
		})

		if err != nil {
			log.Fatalf("Operation failed: %v", err)
		}

		nextStartBlock := tempEndBlock + 1
		c.startBlock = nextStartBlock
	}
}

// TODO: methods here for additional functionalities
// such as handling reconnection logic, managing crawler state, etc.
