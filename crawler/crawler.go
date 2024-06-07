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
	blockchain     string
	startBlock     uint64
	maxBlocksBatch int
	minBlocksBatch int
	confirmations  int
	endBlock       uint64
	force          bool
	baseDir        string
}

// NewCrawler creates a new crawler instance with the given blockchain handler.
func NewCrawler(chain string, startBlock uint64, endBlock uint64, timeout int, batchSize int, confirmations int, baseDir string, force bool) *Crawler {
	fmt.Printf("Start block: %d\n", startBlock)
	return &Crawler{
		blockchain:     chain,
		startBlock:     startBlock,
		endBlock:       endBlock,
		maxBlocksBatch: batchSize,
		minBlocksBatch: 1,
		confirmations:  confirmations,
		baseDir:        baseDir,
		force:          force,
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

func updateBlockIndexFilepaths(indices []indexer.BlockIndex, baseDir string) {
	for i, _ := range indices {
		indices[i].Path = filepath.Join(baseDir, "blocks.proto")
	}
}

func updateTransactionIndexFilepaths(indices []indexer.TransactionIndex, baseDir string) {
	for i, _ := range indices {
		indices[i].Path = filepath.Join(baseDir, "transactions.proto")
	}
}

func UpdateEventIndexFilepaths(indices []indexer.LogIndex, baseDir string) {
	for i, _ := range indices {
		indices[i].Path = filepath.Join(baseDir, "logs.proto")
	}
}

// Start initiates the crawling process for the configured blockchain.
func (c *Crawler) Start() {

	chainType := c.blockchain

	// Initialize the storage instance

	storageInstance, err := storage.NewStorage()

	if err != nil {

		log.Fatalf("Failed to create storage instance: %v", err)
		panic(err)
	}

	client, err := blockchain.NewClient(chainType, BlockchainURLs[chainType])
	if err != nil {
		log.Fatal(err)
	}

	latestBlockNumber, err := client.GetLatestBlockNumber()

	if err != nil {
		log.Fatalf("Failed to get latest block number: %v", err)
	}

	crawlBatchSize := uint64(5) // Crawl in batches of 10
	confirmations := uint64(10) // Consider 10 confirmations to ensure block finality

	fmt.Printf("Latest block number: %d\n", latestBlockNumber)
	fmt.Printf("Start block: %d\n", c.startBlock)

	fmt.Printf("force: %v\n", c.force)

	if c.force {
		if c.startBlock == uint64(0) {
			startBlock := latestBlockNumber.Uint64() - uint64(confirmations) - 100
			c.startBlock = startBlock

		}

	} else {
		latestIndexedBlock, err := indexer.DBConnection.GetLatestBlockNumber(chainType)

		// if err no rows in result set then set startBlock to latestBlockNumber from the blockchain

		if err != nil {
			if err.Error() == "no rows in result set" {
				c.startBlock = latestBlockNumber.Uint64() - uint64(confirmations) - 100
			} else {
				log.Fatalf("Failed to get latest indexed block: %v", err)
			}

		}

		if latestIndexedBlock != uint64(0) {
			c.startBlock = latestIndexedBlock
		}
	}

	endBlock := uint64(0)
	safeBlock := uint64(0)

	fmt.Printf("Start block: %d\n", c.startBlock)

	retryWaitTime := 10 * time.Second
	waitForBlocksTime := retryWaitTime
	maxWaitForBlocksTime := 12 * retryWaitTime
	retryAttempts := 3

	c.baseDir = filepath.Join(SeerCrawlerStoragePrefix, c.baseDir, chainType)

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

		safeBlock = latestBlockNumber.Uint64() - confirmations

		endBlock = c.startBlock + crawlBatchSize

		if endBlock > safeBlock {
			// Auto adjust time
			log.Printf("Waiting for new blocks to be mined. Current block: %d, Safe block: %d", latestBlockNumber, safeBlock)
			time.Sleep(waitForBlocksTime)
			if waitForBlocksTime < maxWaitForBlocksTime {
				waitForBlocksTime = waitForBlocksTime * 2
			}
			continue
		}
		waitForBlocksTime = retryWaitTime

		if c.endBlock != 0 && endBlock > c.endBlock {
			log.Printf("End block reached: %d", endBlock)
			break
		}

		// Retry the operation in case of failure with cumulative attempts
		err = retryOperation(retryAttempts, retryWaitTime, func() error {
			blocks, transactions, blockIndex, transactionIndex, blocksCache, err := blockchain.CrawlBlocks(client, big.NewInt(int64(c.startBlock)), big.NewInt(int64(endBlock)))
			if err != nil {
				return fmt.Errorf("failed to crawl blocks: %w", err)
			}

			batchDir := filepath.Join(c.baseDir, fmt.Sprintf("%d-%d", c.startBlock, endBlock))

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

			if err := storageInstance.Save(filepath.Join(batchDir, "blocks.proto"), encodedBytesBlocks); err != nil {
				return fmt.Errorf("failed to save blocks: %w", err)
			}

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

			if err := storageInstance.Save(filepath.Join(batchDir, "transactions.proto"), encodedBytesTransactions); err != nil {
				return fmt.Errorf("failed to save transactions: %w", err)
			}

			updateBlockIndexFilepaths(blockIndex, batchDir)
			interfaceBlockIndex := make([]interface{}, len(blockIndex))
			for i, v := range blockIndex {
				interfaceBlockIndex[i] = v
			}

			if err := indexer.WriteIndexesToDatabase(c.blockchain, interfaceBlockIndex, "block"); err != nil {
				return fmt.Errorf("failed to write block index to database: %w", err)
			}

			updateTransactionIndexFilepaths(transactionIndex, batchDir)
			interfaceTransactionIndex := make([]interface{}, len(transactionIndex))
			for i, v := range transactionIndex {
				interfaceTransactionIndex[i] = v
			}

			if err := indexer.WriteIndexesToDatabase(c.blockchain, interfaceTransactionIndex, "transaction"); err != nil {
				return fmt.Errorf("failed to write transaction index to database: %w", err)
			}

			events, eventsIndex, err := blockchain.CrawlEvents(client, big.NewInt(int64(c.startBlock)), big.NewInt(int64(endBlock)), blocksCache)
			if err != nil {
				return fmt.Errorf("failed to crawl events: %w", err)
			}

			UpdateEventIndexFilepaths(eventsIndex, batchDir)
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

			if err := storageInstance.Save(filepath.Join(batchDir, "logs.proto"), encodedBytesEvents); err != nil {
				return fmt.Errorf("failed to save events: %w", err)
			}

			return nil
		})

		if err != nil {
			log.Fatalf("Operation failed: %v", err)
		}

		nextStartBlock := endBlock + 1
		c.startBlock = nextStartBlock
	}
}

// TODO: methods here for additional functionalities
// such as handling reconnection logic, managing crawler state, etc.
