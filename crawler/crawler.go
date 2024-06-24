package crawler

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"path/filepath"
	"sync"
	"time"

	seer_blockchain "github.com/moonstream-to/seer/blockchain"
	"github.com/moonstream-to/seer/indexer"
	"github.com/moonstream-to/seer/storage"
	"google.golang.org/protobuf/encoding/protodelim"
)

var CurrentBlockchainState BlockchainState

type BlockchainState struct {
	LatestBlockNumber *big.Int

	mux sync.RWMutex
}

func (bs *BlockchainState) SetLatestBlockNumber(blockNumber *big.Int) {
	bs.mux.Lock()
	bs.LatestBlockNumber = blockNumber
	bs.mux.Unlock()
}

func (bs *BlockchainState) GetLatestBlockNumber() *big.Int {
	bs.mux.RLock()
	blockNumber := bs.LatestBlockNumber
	bs.mux.RUnlock()
	return blockNumber
}

// Crawler defines the crawler structure.
type Crawler struct {
	Client          seer_blockchain.BlockchainClient
	StorageInstance storage.Storer

	blockchain    string
	startBlock    int64
	endBlock      int64
	batchSize     int64
	confirmations int64
	force         bool
	baseDir       string
	basePath      string
}

// NewCrawler creates a new crawler instance with the given blockchain handler.
func NewCrawler(blockchain string, startBlock, endBlock, batchSize, confirmations int64, timeout int, baseDir string, force bool) (*Crawler, error) {
	var crawler Crawler

	basePath := filepath.Join(baseDir, SeerCrawlerStoragePrefix, "data", blockchain)
	storageInstance, err := storage.NewStorage(storage.SeerCrawlerStorageType, basePath)
	if err != nil {
		log.Fatalf("Failed to create storage instance: %v", err)
		panic(err)
	}

	client, err := seer_blockchain.NewClient(blockchain, BlockchainURLs[blockchain], timeout)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Initialized new crawler at blockchain: %s, startBlock: %d, endBlock: %d, force: %t", blockchain, startBlock, endBlock, force)
	crawler = Crawler{
		Client:          client,
		StorageInstance: storageInstance,

		blockchain:    blockchain,
		startBlock:    startBlock,
		endBlock:      endBlock,
		batchSize:     batchSize,
		confirmations: confirmations,
		force:         force,
		baseDir:       baseDir,
		basePath:      basePath,
	}

	return &crawler, nil
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

func SetDefaultStartBlock(confirmations int64, latestBlockNumber *big.Int) int64 {
	startBlock := latestBlockNumber.Int64() - confirmations - 100
	log.Printf("Start block set with shift: %d\n", startBlock)
	return startBlock
}

// Start initiates the crawling process for the configured blockchain.
func (c *Crawler) Start(threads int) {
	latestBlockNumber := CurrentBlockchainState.GetLatestBlockNumber()
	if c.force {
		if c.startBlock == 0 {
			c.startBlock = SetDefaultStartBlock(c.confirmations, latestBlockNumber)
		}
	} else {
		latestIndexedBlock, err := indexer.DBConnection.GetLatestDBBlockNumber(c.blockchain)

		// If there are no rows in result then set startBlock with SetDefaultStartBlock()

		if err != nil {
			if err.Error() == "no rows in result set" {
				c.startBlock = SetDefaultStartBlock(c.confirmations, latestBlockNumber)
			} else {
				log.Fatalf("Failed to get latest indexed block: %v", err)
			}

		}

		if latestIndexedBlock != 0 {
			c.startBlock = int64(latestIndexedBlock) + 1
			log.Printf("Start block fetched from indexes database and set to: %d\n", c.startBlock)
		}
	}

	tempEndBlock := c.startBlock + c.batchSize
	var safeBlock int64

	retryWaitTime := 10 * time.Second
	waitForBlocksTime := retryWaitTime
	maxWaitForBlocksTime := 12 * retryWaitTime
	retryAttempts := 3

	var err error
	var isEnd bool
	for {
		// Using CurrentBlockchainState (in future via mutex for async) to not fetch too often if there is a big difference
		if tempEndBlock+c.confirmations >= latestBlockNumber.Int64() {
			latestBlockNumber, err = c.Client.GetLatestBlockNumber()
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
		}

		safeBlock = latestBlockNumber.Int64() - c.confirmations

		tempEndBlock = c.startBlock + c.batchSize
		if c.endBlock != 0 {
			if c.endBlock <= tempEndBlock {
				tempEndBlock = c.endBlock
				isEnd = true
				log.Printf("End block %d almost reached", tempEndBlock)
			}
		}

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

		// Retry the operation in case of failure with cumulative attempts
		err = retryOperation(retryAttempts, retryWaitTime, func() error {
			batchDir := fmt.Sprintf("%d-%d", c.startBlock, tempEndBlock)
			log.Printf("Operates with batch of blocks: %s", batchDir)

			// Fetch blocks with transactions
			blocks, transactions, blockIndex, transactionIndex, blocksCache, err := seer_blockchain.CrawlBlocks(c.Client, big.NewInt(c.startBlock), big.NewInt(tempEndBlock), SEER_CRAWLER_DEBUG, threads)
			if err != nil {
				return fmt.Errorf("failed to crawl blocks: %w", err)
			}
			var blocksBuffer bytes.Buffer

			for _, block := range blocks {
				protodelim.MarshalTo(&blocksBuffer, block)
			}

			if err := c.StorageInstance.Save(batchDir, "blocks.proto", blocksBuffer); err != nil {
				return fmt.Errorf("failed to save blocks.proto: %w", err)
			}
			log.Printf("Saved blocks.proto to %s", batchDir)

			updateBlockIndexFilepaths(blockIndex, c.basePath, batchDir)
			interfaceBlockIndex := make([]interface{}, len(blockIndex))
			for i, v := range blockIndex {
				interfaceBlockIndex[i] = v
			}

			if err := indexer.WriteIndexesToDatabase(c.blockchain, interfaceBlockIndex, "block"); err != nil {
				return fmt.Errorf("failed to write block index to database: %w", err)
			}

			var transactionsBuffer bytes.Buffer

			for _, transaction := range transactions {
				protodelim.MarshalTo(&transactionsBuffer, transaction)
			}

			if err := c.StorageInstance.Save(batchDir, "transactions.proto", transactionsBuffer); err != nil {
				return fmt.Errorf("failed to save transactions.proto: %w", err)
			}
			log.Printf("Saved transactions.proto to %s", batchDir)

			updateTransactionIndexFilepaths(transactionIndex, c.basePath, batchDir)
			interfaceTransactionIndex := make([]interface{}, len(transactionIndex))
			for i, v := range transactionIndex {
				interfaceTransactionIndex[i] = v
			}

			if err := indexer.WriteIndexesToDatabase(c.blockchain, interfaceTransactionIndex, "transaction"); err != nil {
				return fmt.Errorf("failed to write transaction index to database: %w", err)
			}

			events, eventsIndex, err := seer_blockchain.CrawlEvents(c.Client, big.NewInt(int64(c.startBlock)), big.NewInt(int64(tempEndBlock)), blocksCache, SEER_CRAWLER_DEBUG)
			if err != nil {
				return fmt.Errorf("failed to crawl events: %w", err)
			}

			var eventsBuffer bytes.Buffer

			for _, event := range events {
				protodelim.MarshalTo(&eventsBuffer, event)
			}

			if err := c.StorageInstance.Save(batchDir, "logs.proto", eventsBuffer); err != nil {
				return fmt.Errorf("failed to save logs.proto: %w", err)
			}
			log.Printf("Saved logs.proto to %s", batchDir)

			UpdateEventIndexFilepaths(eventsIndex, c.basePath, batchDir)
			interfaceEventsIndex := make([]interface{}, len(eventsIndex))
			for i, v := range eventsIndex {
				interfaceEventsIndex[i] = v
			}

			if err := indexer.WriteIndexesToDatabase(c.blockchain, interfaceEventsIndex, "log"); err != nil {
				return fmt.Errorf("failed to write event index to database: %w", err)
			}

			return nil
		})
		if err != nil {
			log.Fatalf("Operation failed: %v", err)
		}

		if isEnd {
			break
		}

		c.startBlock = tempEndBlock + 1
	}
}

// TODO: methods here for additional functionalities
// such as handling reconnection logic, managing crawler state, etc.
