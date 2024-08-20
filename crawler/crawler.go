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

// BlockchainState represents the current state of the blockchain, including the latest block number and the time when it was fetched.
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

// DynamicBatch is used to dynamically calculate the batch size of blocks in order to fetch them.
type DynamicBatch struct {
	Size int64

	mux sync.RWMutex
}

func (db *DynamicBatch) GetSize() int64 {
	db.mux.RLock()
	size := db.Size
	db.mux.RUnlock()
	return size
}

// DynamicDecreaseSize by dividing it by 2 and checking if it is not below the diff argument.
func (db *DynamicBatch) DynamicDecreaseSize(diff int64) {
	db.mux.Lock()
	defer db.mux.Unlock()

	if db.Size == 0 {
		return
	}

	// Decrease the size dynamically by half, but ensure at least a reduction of 1
	step := max(db.Size/2, 1)
	db.Size -= step

	// Ensure size doesn't fall below the specified diff or go below 0
	db.Size = max(db.Size, diff)
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

// ProcessAndPush makes preparations for blocks, txs, and logs data and pushes them to the database with storage.
func (cp *CrawlPack) ProcessAndPush(client seer_blockchain.BlockchainClient, crawler *Crawler) error {
	packRange := fmt.Sprintf("%d-%d", cp.PackStartBlock, cp.PackEndBlock)

	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 1)

	go func() {
		defer wg.Done()

		// Prepare and save proto data
		blocksBatch, batchErr := client.ProcessBlocksToBatch(cp.BlocksPack)
		if batchErr != nil {
			errChan <- fmt.Errorf("unable to process blocks to batch: %w", batchErr)
			return
		}
		dataBytes, marshalErr := proto.Marshal(blocksBatch)
		if marshalErr != nil {
			errChan <- fmt.Errorf("failed to marshal blocks: %v", marshalErr)
			return
		}

		if err := crawler.StorageInstance.Save(packRange, "data.proto", *bytes.NewBuffer(dataBytes)); err != nil {
			errChan <- fmt.Errorf("failed to save data.proto: %w", err)
			return
		}
		log.Printf("Saved .proto blocks with transactions and events to %s", packRange)
	}()

	go func() {
		defer wg.Done()

		// Prepare and save indexes data
		var interfaceBlocksIndexPack []indexer.BlockIndex
		for _, v := range cp.BlocksIndexPack {
			v.Path = filepath.Join(crawler.basePath, packRange, "data.proto")
			interfaceBlocksIndexPack = append(interfaceBlocksIndexPack, v)
		}

		var interfaceTxsIndexPack []indexer.TransactionIndex
		for _, v := range cp.TxsIndexPack {
			v.Path = filepath.Join(crawler.basePath, packRange, "data.proto")
			interfaceTxsIndexPack = append(interfaceTxsIndexPack, v)
		}

		var interfaceEventsIndexPack []indexer.LogIndex
		for _, v := range cp.EventsIndexPack {
			v.Path = filepath.Join(crawler.basePath, packRange, "data.proto")
			interfaceEventsIndexPack = append(interfaceEventsIndexPack, v)
		}

		err := indexer.WriteIndicesToDatabase(crawler.blockchain, interfaceBlocksIndexPack, interfaceTxsIndexPack, interfaceEventsIndexPack)
		if err != nil {
			errChan <- fmt.Errorf("failed to write indices to database: %w", err)
			return
		}
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// Main crawler loop.
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

		// If there are no rows in result then set startBlock with shift
		if latestErr != nil {
			if latestErr.Error() != "no rows in result set" {
				log.Fatalf("Failed to get latest indexed block: %v", latestErr)
			}

			latestIndexedBlock = uint64(CurrentBlockchainState.GetLatestBlockNumber().Int64() - c.confirmations - SeerDefaultBlockShift)
			log.Printf("There are no records in database, applied shift %d to latest block number", SeerDefaultBlockShift)
		}

		c.startBlock = int64(latestIndexedBlock) + 1
		log.Printf("Start block set according with indexes database to: %d\n", c.startBlock)
	}

	var isFinal bool
	var endBlock int64
	var crawlPack CrawlPack

	for {
		endBlock = c.startBlock + dynamicBatch.GetSize()

		if SEER_CRAWLER_DEBUG {
			log.Printf("[DEBUG] [crawler.Start.1] latestBlock: %d, c.endBlock: %d, c.startBlock: %d, dynamicBatchSize: %d, PackStartBlock: %d", CurrentBlockchainState.GetLatestBlockNumber().Uint64(), endBlock, c.startBlock, dynamicBatch.GetSize(), crawlPack.PackStartBlock)
		}

		// Check if final block specified at trigger stop
		if c.finalBlock != 0 && endBlock >= c.finalBlock {
			endBlock = c.finalBlock
			isFinal = true
		}

		// Push pack to database and storage if time comes
		if len(crawlPack.BlocksPack) > 0 {
			if crawlPack.PackCrawlStartTs.Add(protoDurationTimeLimit).Before(time.Now()) || crawlPack.PackSize >= protoBufferSizeLimit {
				if papErr := crawlPack.ProcessAndPush(c.Client, c); papErr != nil {
					log.Fatalf("failed to process and push: %v", papErr)
				}
				crawlPack = CrawlPack{}
			}
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
			// Identified slow blockchain, reducing batch size
			if time.Since(CurrentBlockchainState.GetLatestUpdateTs()) >= retryWaitTime {
				dynamicBatch.DynamicDecreaseSize(safeBlock - c.startBlock)
			}

			log.Printf("Waiting %d seconds for new blocks to be mined. Current blockchain latest block number: %d, calculated crawler end block: %d and dynamic batch size set to: %d", int(waitForBlocksTime.Seconds()), CurrentBlockchainState.GetLatestBlockNumber().Int64(), endBlock, dynamicBatch.GetSize())

			time.Sleep(waitForBlocksTime)
			if waitForBlocksTime < maxWaitForBlocksTime {
				waitForBlocksTime = waitForBlocksTime * 2
			}

			continue
		}
		waitForBlocksTime = retryWaitTime
		dynamicBatch.Size = c.batchSize

		if SEER_CRAWLER_DEBUG {
			log.Printf("[DEBUG] [crawler.Start.2] latestBlock: %d, c.endBlock: %d, c.startBlock: %d, dynamicBatchSize: %d, PackStartBlock: %d", CurrentBlockchainState.GetLatestBlockNumber().Uint64(), endBlock, c.startBlock, dynamicBatch.GetSize(), crawlPack.PackStartBlock)
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
			if len(crawlPack.BlocksPack) > 0 {
				if crawlPack.PackCrawlStartTs.Add(protoDurationTimeLimit).Before(time.Now()) || crawlPack.PackSize >= protoBufferSizeLimit {
					if papErr := crawlPack.ProcessAndPush(c.Client, c); papErr != nil {
						return papErr
					}
					crawlPack = CrawlPack{}
				}
			}

			if SEER_CRAWLER_DEBUG {
				log.Printf("[DEBUG] [crawler.Start.3] latestBlock: %d, c.endBlock: %d, c.startBlock: %d, dynamicBatchSize: %d, PackStartBlock: %d", CurrentBlockchainState.GetLatestBlockNumber().Uint64(), endBlock, c.startBlock, dynamicBatch.GetSize(), crawlPack.PackStartBlock)
			}

			return nil
		}); retryErr != nil {
			log.Fatalf("Crawl retry operation failed: %v", retryErr)
		}

		if isFinal {
			break
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
