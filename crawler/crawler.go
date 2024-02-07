package crawler

import (
	"log"
	"math/big"
	"os"
	"path/filepath"

	"github.com/moonstream-to/seer/blockchain/common"
	"google.golang.org/protobuf/proto"
)

// Crawler defines the crawler structure.
type Crawler struct {
	blockchain string
}

// NewCrawler creates a new crawler instance with the given blockchain handler.
func NewCrawler(bh string) *Crawler {
	return &Crawler{
		blockchain: bh,
	}

}

func writeProtoMessagesToFile(messages []proto.Message, filePath string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	for _, msg := range messages {
		data, err := proto.Marshal(msg)
		if err != nil {
			log.Fatalf("Failed to marshal proto message: %v", err)
		}
		if _, err := file.Write(data); err != nil {
			log.Fatalf("Failed to write data to file: %v", err)
		}
	}
}

// Start initiates the crawling process for the configured blockchain.
func (c *Crawler) Start() {

	chainType := "ethereum"
	log.Println("Starting crawler...")
	client, err := common.NewClient(chainType, "http://localhost:8545")

	if err != nil {
		log.Fatal(err)
	}

	// get the latest block number

	latestBlockNumber, err := client.GetLatestBlockNumber()

	if err != nil {
		log.Fatal(err)
	}

	// Get the latest block number
	if err != nil {
		log.Fatalf("Failed to get latest block number: %v", err)
	}

	// Assuming latestBlockNumber is *big.Int and operations for big.Int
	startBlock := new(big.Int).Sub(latestBlockNumber, big.NewInt(100)) // Start 100 blocks behind the latest
	crawlBatchSize := big.NewInt(10)                                   // Crawl in batches of 10
	confirmations := big.NewInt(10)                                    // Consider 10 confirmations

	for latestBlockNumber.Int64()-startBlock.Int64() > confirmations.Int64()+crawlBatchSize.Int64() {

		endBlock := new(big.Int).Add(startBlock, crawlBatchSize)

		// Fetch the blocks and transactions

		blocks, transactions, blockIndex, transactionIndex, err := common.CrawlBlocks(client, startBlock, endBlock)

		if err != nil {
			log.Fatalf("Failed to get blocks: %v", err)
		}

		// Process the blocks and transactions

		// write the blocks and transactions to disk

		batchDir := filepath.Join("data", startBlock.String()+"-"+endBlock.String())
		if err := os.MkdirAll(batchDir, 0755); err != nil {
			log.Fatal(err)
		}

		if err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}

		// write the blocks and transactions to disk

		writeProtoMessagesToFile(blocks, filepath.Join(batchDir, "blocks.proto"))

		writeProtoMessagesToFile(transactions, filepath.Join(batchDir, "transactions.proto"))

		// write the block and transaction indices to disk

		// open index general file for writing

		file_out := "index.db"

		file, err := os.Create(file_out)

		if err != nil {
			panic(err)
		}

		defer file.Close()

		// Update the startBlock for the next batch
		startBlock = endBlock
	}
}

// You can add more methods here for additional functionalities
// such as handling reconnection logic, managing crawler state, etc.
