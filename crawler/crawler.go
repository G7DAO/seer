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
	LatestUpdateTs    time.Time
	LatestBlockNumber *big.Int

	mux sync.RWMutex
}

func (bs *BlockchainState) RaiseLatestBlockNumber(blockNumber *big.Int) {
	bs.mux.Lock()
	if bs.LatestBlockNumber == nil || bs.LatestBlockNumber.Cmp(blockNumber) == -1 {
		bs.LatestUpdateTs = time.Now()
		bs.LatestBlockNumber = blockNumber
	}
	bs.mux.Unlock()
}

func (bs *BlockchainState) GetLatestBlockNumber() *big.Int {
	bs.mux.RLock()
	blockNumber := bs.LatestBlockNumber
	bs.mux.RUnlock()
	return blockNumber
}

func (bs *BlockchainState) GetLatestUpdateTs() time.Time {
	bs.mux.RLock()
	latestUpdateTs := bs.LatestUpdateTs
	bs.mux.RUnlock()
	return latestUpdateTs
}

type DynamicBatch struct {
	Size int64
}

func (db *DynamicBatch) GetSize() int64 {
	return db.Size
}

func (db *DynamicBatch) DynamicDecreaseSize(diff int64) {
	size := db.Size

	step := size / 2
	size -= step

	if diff >= size {
		size = diff
	}

	if size <= 1 {
		size = 1
	}

	db.Size = size
}

// Crawler defines the crawler structure.
type Crawler struct {
	Client          seer_blockchain.BlockchainClient
	StorageInstance storage.Storer

	blockchain     string
	startBlock     int64
	finalBlock     int64
	confirmations  int64
	batchSize      int64
	baseDir        string
	basePath       string
	protoSizeLimit uint64
	protoTimeLimit int
}

// NewCrawler creates a new crawler instance with the given blockchain handler.
func NewCrawler(blockchain string, startBlock, finalBlock, confirmations, batchSize int64, timeout int, baseDir string, protoSizeLimit uint64, protoTimeLimit int) (*Crawler, error) {
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

	log.Printf("Initialized new crawler at blockchain: %s, startBlock: %d, finalBlock: %d", blockchain, startBlock, finalBlock)
	crawler = Crawler{
		Client:          client,
		StorageInstance: storageInstance,

		blockchain:     blockchain,
		startBlock:     startBlock,
		finalBlock:     finalBlock,
		confirmations:  confirmations,
		batchSize:      batchSize,
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

func (c *Crawler) PushPackOfData(blocksBufferPack *bytes.Buffer, blocksIndexPack []indexer.BlockIndex, txsIndexPack []indexer.TransactionIndex, eventsIndexPack []indexer.LogIndex, packStartBlock, packEndBlock uint64) error {
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

type CrawlPack struct {
	PackSize         int64
	PackCrawlStartTs time.Time

	BlocksPack      []proto.Message
	BlocksIndexPack []indexer.BlockIndex
	TxsIndexPack    []indexer.TransactionIndex
	EventsIndexPack []indexer.LogIndex

	PackStartBlock int64
	PackEndBlock   int64

	isInitialized bool
}

func (cp *CrawlPack) Initialize(startBlock int64) {
	cp.PackStartBlock = startBlock
	cp.PackCrawlStartTs = time.Now()
	cp.isInitialized = true
}

func (cp *CrawlPack) ProcessAndPush(client seer_blockchain.BlockchainClient, crawler *Crawler) error {
	blocksBatch, batchErr := client.ProcessBlocksToBatch(cp.BlocksPack)
	if batchErr != nil {
		return fmt.Errorf("unable to process blocks to batch: %w", batchErr)
	}
	dataBytes, marshalErr := proto.Marshal(blocksBatch)
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal blocks: %v", marshalErr)
	}
	if pushEr := crawler.PushPackOfData(bytes.NewBuffer(dataBytes), cp.BlocksIndexPack, cp.TxsIndexPack, cp.EventsIndexPack, uint64(cp.PackStartBlock), uint64(cp.PackEndBlock)); pushEr != nil {
		return fmt.Errorf("unable to push data correctly: %w", pushEr)
	}
	return nil
}

func (c *Crawler) Start(threads int) {
	protoBufferSizeLimit := int64(c.protoSizeLimit * 1024 * 1024) // In Mb
	protoDurationTimeLimit := time.Duration(c.protoTimeLimit) * time.Second

	dynamicBatch := DynamicBatch{
		Size: c.batchSize,
	}

	retryAttempts := 3
	retryWaitTime := 5 * time.Second
	waitForBlocksTime := retryWaitTime
	maxWaitForBlocksTime := 24 * retryWaitTime

	// If Start block is not set, using last crawled block from indexes database
	if c.startBlock == 0 {
		latestIndexedBlock, latestErr := indexer.DBConnection.GetLatestDBBlockNumber(c.blockchain)

		// If there are no rows in result then set startBlock with SetDefaultStartBlock()
		if latestErr != nil {
			if latestErr.Error() != "no rows in result set" {
				log.Fatalf("Failed to get latest indexed block: %v", latestErr)
			}

			latestIndexedBlock = uint64(SetDefaultStartBlock(c.confirmations, CurrentBlockchainState.GetLatestBlockNumber()))
		}

		c.startBlock = int64(latestIndexedBlock) + 1
		log.Printf("Start block fetched from indexes database and set to: %d\n", c.startBlock)
	}

	var endBlock int64
	var crawlPack CrawlPack

	for {
		endBlock = c.startBlock + dynamicBatch.GetSize()

		if SEER_CRAWLER_DEBUG {
			log.Printf("[DEBUG] [crawler.Start.1] latestBlock: %d, c.endBlock: %d, c.startBlock: %d, dynamicBatchSize: %d, PackStartBlock: %d", CurrentBlockchainState.GetLatestBlockNumber().Uint64(), endBlock, c.startBlock, dynamicBatch.GetSize(), crawlPack.PackStartBlock)
		}

		// Check if final block specified at start reached then stop
		if c.finalBlock != 0 && endBlock >= c.finalBlock {
			break
		}

		// Push pack to database and storage if time comes
		if crawlPack.PackCrawlStartTs.Add(protoDurationTimeLimit).Before(time.Now()) || crawlPack.PackSize >= protoBufferSizeLimit {
			if papErr := crawlPack.ProcessAndPush(c.Client, c); papErr != nil {
				log.Fatalf("failed to process and push: %v", papErr)
			}
			crawlPack = CrawlPack{}
		}

		// Check if next iteration will overtake blockchain latest block minus confirmation
		safeBlock := CurrentBlockchainState.GetLatestBlockNumber().Int64() - c.confirmations
		if endBlock >= safeBlock {
			latestBlockNumber, latestErr := seer_blockchain.GetLatestBlockNumberWithRetry(c.Client, retryAttempts, retryWaitTime)
			if latestErr != nil {
				log.Fatalf("failed to fetch latest block from blockchain: %v", latestErr)
			}
			CurrentBlockchainState.RaiseLatestBlockNumber(latestBlockNumber)
		}

		// Check if next iteration is again overtake blockchain latest block minus confirmation
		if endBlock > safeBlock {
			sinceLatestUpdateTs := time.Since(CurrentBlockchainState.GetLatestUpdateTs())
			// Identified slow blockchain, reducing batch size
			if dynamicBatch.GetSize() != 1 && sinceLatestUpdateTs >= retryWaitTime {
				dynamicBatch.DynamicDecreaseSize(safeBlock - c.startBlock)
			}

			log.Printf("Waiting %d seconds for new blocks to be mined. Current blockchain latest block number: %d, crawler end block: %d and dynamic batch size: %d", int(waitForBlocksTime.Seconds()), CurrentBlockchainState.GetLatestBlockNumber().Int64(), endBlock, dynamicBatch.GetSize())

			time.Sleep(waitForBlocksTime)
			if waitForBlocksTime < maxWaitForBlocksTime {
				waitForBlocksTime = waitForBlocksTime * 2
			}

			continue
		}
		waitForBlocksTime = retryWaitTime
		dynamicBatch.Size = c.batchSize

		if SEER_CRAWLER_DEBUG {
			log.Printf("[DEBUG] [crawler.Start.1] latestBlock: %d, c.endBlock: %d, c.startBlock: %d, dynamicBatchSize: %d, PackStartBlock: %d", CurrentBlockchainState.GetLatestBlockNumber().Uint64(), endBlock, c.startBlock, dynamicBatch.GetSize(), crawlPack.PackStartBlock)
		}

		// Check if crawlPack not yet initialized
		if !crawlPack.isInitialized {
			crawlPack.Initialize(c.startBlock)
		}

		if retryErr := retryOperation(retryAttempts, retryWaitTime, func() error {
			// Fetch blocks with transactions
			blocks, blocksIndex, txsIndex, eventsIndex, blocksSize, crawlErr := seer_blockchain.CrawlEntireBlocks(c.Client, new(big.Int).SetInt64(c.startBlock), new(big.Int).SetInt64(endBlock), SEER_CRAWLER_DEBUG, threads)
			if crawlErr != nil {
				return fmt.Errorf("failed to crawl blocks, txs and events: %w", crawlErr)
			}

			crawlPack.PackSize += int64(blocksSize)
			crawlPack.BlocksPack = append(crawlPack.BlocksPack, blocks...)
			crawlPack.BlocksIndexPack = append(crawlPack.BlocksIndexPack, blocksIndex...)
			crawlPack.TxsIndexPack = append(crawlPack.TxsIndexPack, txsIndex...)
			crawlPack.EventsIndexPack = append(crawlPack.EventsIndexPack, eventsIndex...)
			crawlPack.PackEndBlock = endBlock

			// Push pack to database and storage if time comes
			if crawlPack.PackCrawlStartTs.Add(protoDurationTimeLimit).Before(time.Now()) || crawlPack.PackSize >= protoBufferSizeLimit {
				if papErr := crawlPack.ProcessAndPush(c.Client, c); papErr != nil {
					return papErr
				}
				crawlPack = CrawlPack{}
			}

			if SEER_CRAWLER_DEBUG {
				log.Printf("[DEBUG] [crawler.Start.1] latestBlock: %d, c.endBlock: %d, c.startBlock: %d, dynamicBatchSize: %d, PackStartBlock: %d", CurrentBlockchainState.GetLatestBlockNumber().Uint64(), endBlock, c.startBlock, dynamicBatch.GetSize(), crawlPack.PackStartBlock)
			}

			return nil
		}); retryErr != nil {
			log.Fatalf("Crawl retry operation failed: %v", retryErr)
		}

		c.startBlock = endBlock + 1
	}

	// If after stop there are some blocks, do not leave it
	if len(crawlPack.BlocksPack) > 0 {
		if papErr := crawlPack.ProcessAndPush(c.Client, c); papErr != nil {
			log.Fatalf("failed to process and push: %v", papErr)
		}
		crawlPack = CrawlPack{}
	}
}

// TODO: methods here for additional functionalities
// such as handling reconnection logic, managing crawler state, etc.
