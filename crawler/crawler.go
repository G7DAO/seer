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
	"google.golang.org/protobuf/proto"
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

	blockchain     string
	startBlock     int64
	endBlock       int64
	confirmations  int64
	batchSize      int64
	force          bool
	baseDir        string
	basePath       string
	protoSizeLimit uint64
	protoTimeLimit int
}

// NewCrawler creates a new crawler instance with the given blockchain handler.
func NewCrawler(blockchain string, startBlock, endBlock, confirmations int64, batchSize int64, timeout int, baseDir string, force bool, protoSizeLimit uint64, protoTimeLimit int) (*Crawler, error) {
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

		blockchain:     blockchain,
		startBlock:     startBlock,
		endBlock:       endBlock,
		confirmations:  confirmations,
		batchSize:      batchSize,
		force:          force,
		baseDir:        baseDir,
		basePath:       basePath,
		protoSizeLimit: protoSizeLimit,
		protoTimeLimit: protoTimeLimit,
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

func SetDefaultStartBlock(confirmations int64, latestBlockNumber *big.Int) int64 {
	startBlock := latestBlockNumber.Int64() - confirmations - 100
	log.Printf("Start block set with shift: %d\n", startBlock)
	return startBlock
}

type BlocksBufferBatch struct {
	StartBlock int64
	EndBlock   int64

	Buffer bytes.Buffer
}

func (c *Crawler) PushPackOfData(blocksBufferPack *bytes.Buffer, blocksIndexPack []indexer.BlockIndex, txsIndexPack []indexer.TransactionIndex, eventsIndexPack []indexer.LogIndex, packStartBlock, packEndBlock int64) error {
	packRange := fmt.Sprintf("%d-%d", packStartBlock, packEndBlock)

	// Save proto data
	if err := c.StorageInstance.Save(packRange, "data.proto", *blocksBufferPack); err != nil {
		return fmt.Errorf("failed to save data.proto: %w", err)
	}
	log.Printf("Saved .proto blocks with transactions and events to %s", packRange)

	// Save indexes data
	var interfaceBlocksIndexPack []indexer.BlockIndex
	for _, v := range blocksIndexPack {
		v.Path = filepath.Join(c.basePath, packRange, "data.proto")
		interfaceBlocksIndexPack = append(interfaceBlocksIndexPack, v)
	}

	var interfaceTxsIndexPack []indexer.TransactionIndex
	for _, v := range txsIndexPack {
		v.Path = filepath.Join(c.basePath, packRange, "data.proto")
		interfaceTxsIndexPack = append(interfaceTxsIndexPack, v)
	}

	var interfaceEventsIndexPack []indexer.LogIndex
	for _, v := range eventsIndexPack {
		v.Path = filepath.Join(c.basePath, packRange, "data.proto")
		interfaceEventsIndexPack = append(interfaceEventsIndexPack, v)
	}

	// Write indexes to database
	err := indexer.WriteIndicesToDatabase(c.blockchain, interfaceBlocksIndexPack, interfaceTxsIndexPack, interfaceEventsIndexPack)

	if err != nil {
		return fmt.Errorf("failed to write indices to database: %w", err)
	}

	return nil
}

// Start initiates the crawling process for the configured blockchain.
func (c *Crawler) Start(threads int) {
	protoBufferSizeLimit := c.protoSizeLimit * 1024 * 1024 // In Mb
	protoDurationTimeLimit := time.Duration(c.protoTimeLimit) * time.Second

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

	// Variables to accumulate packs before write
	var blocksPack []proto.Message
	var blocksPackSize uint64
	packCrawlStartTs := time.Now()
	var blocksIndexPack []indexer.BlockIndex
	var txsIndexPack []indexer.TransactionIndex
	var eventsIndexPack []indexer.LogIndex
	packStartBlock := c.startBlock

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
			// Before wait check if there is something to push
			if packCrawlStartTs.Add(protoDurationTimeLimit).Before(time.Now()) && len(blocksPack) > 0 {
				blocksBatch, batchErr := c.Client.ProcessBlocksToBatch(blocksPack)
				if batchErr != nil {
					log.Printf("Unable to process blocks to batch, err: %v", batchErr)
				}

				dataBytes, err := proto.Marshal(blocksBatch)
				if err != nil {
					log.Fatalf("Failed to marshal blocks: %v", err)
				}

				if pushEr := c.PushPackOfData(bytes.NewBuffer(dataBytes), blocksIndexPack, txsIndexPack, eventsIndexPack, packStartBlock, tempEndBlock); err != nil {
					log.Printf("Unable to push data correctly, err: %v", pushEr)
				}

				blocksPackSize = 0
				blocksPack = []proto.Message{}
				blocksIndexPack = []indexer.BlockIndex{}
				txsIndexPack = []indexer.TransactionIndex{}
				eventsIndexPack = []indexer.LogIndex{}

				packStartBlock = tempEndBlock + 1
				packCrawlStartTs = time.Now()
			}

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
			log.Printf("Operates with batch of blocks: %d-%d", c.startBlock, tempEndBlock)

			// Fetch blocks with transactions
			blocks, blocksIndex, txsIndex, eventsIndex, blocksSize, crawlErr := seer_blockchain.CrawlEntireBlocks(c.Client, big.NewInt(c.startBlock), big.NewInt(tempEndBlock), SEER_CRAWLER_DEBUG, threads)
			if crawlErr != nil {
				return fmt.Errorf("failed to crawl blocks, txs and events: %w", err)
			}

			blocksPackSize += blocksSize
			blocksPack = append(blocksPack, blocks...)

			blocksIndexPack = append(blocksIndexPack, blocksIndex...)
			txsIndexPack = append(txsIndexPack, txsIndex...)
			eventsIndexPack = append(eventsIndexPack, eventsIndex...)

			if packCrawlStartTs.Add(protoDurationTimeLimit).Before(time.Now()) || blocksPackSize >= protoBufferSizeLimit {
				blocksBatch, batchErr := c.Client.ProcessBlocksToBatch(blocksPack)
				if batchErr != nil {
					return fmt.Errorf("unable to process blocks to batch: %w", batchErr)
				}

				dataBytes, err := proto.Marshal(blocksBatch)
				if err != nil {
					log.Fatalf("Failed to marshal blocks: %v", err)
				}

				if pushEr := c.PushPackOfData(bytes.NewBuffer(dataBytes), blocksIndexPack, txsIndexPack, eventsIndexPack, packStartBlock, tempEndBlock); err != nil {
					return fmt.Errorf("unable to push data correctly: %w", pushEr)
				}

				blocksPackSize = 0
				blocksPack = []proto.Message{}
				blocksIndexPack = []indexer.BlockIndex{}
				txsIndexPack = []indexer.TransactionIndex{}
				eventsIndexPack = []indexer.LogIndex{}

				packStartBlock = tempEndBlock + 1
				packCrawlStartTs = time.Now()
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

	if len(blocksPack) > 0 {
		blocksBatch, batchErr := c.Client.ProcessBlocksToBatch(blocksPack)
		if batchErr != nil {
			log.Printf("Unable to process blocks to batch, err: %v", batchErr)
		}

		dataBytes, err := proto.Marshal(blocksBatch)
		if err != nil {
			log.Fatalf("Failed to marshal blocks: %v", err)
		}

		if pushEr := c.PushPackOfData(bytes.NewBuffer(dataBytes), blocksIndexPack, txsIndexPack, eventsIndexPack, packStartBlock, tempEndBlock); err != nil {
			log.Printf("Unable to push last pack of data correctly, err: %v", pushEr)
		}

		blocksPackSize = 0
		blocksPack = []proto.Message{}
		blocksIndexPack = []indexer.BlockIndex{}
		txsIndexPack = []indexer.TransactionIndex{}
		eventsIndexPack = []indexer.LogIndex{}

		packStartBlock = tempEndBlock + 1
		packCrawlStartTs = time.Now()
	}
}

// TODO: methods here for additional functionalities
// such as handling reconnection logic, managing crawler state, etc.
