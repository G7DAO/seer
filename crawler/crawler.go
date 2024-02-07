package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"github.com/moonstream-to/seer/blockchain/common"
	"github.com/moonstream-to/seer/indexer"
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

func updateBlockIndexFilepaths(indices []indexer.BlockIndex, baseDir string) {
	for i, _ := range indices {
		indices[i].Filepath = filepath.Join(baseDir, "blocks.proto")
	}
}

func updateTransactionIndexFilepaths(indices []indexer.TransactionIndex, baseDir string) {
	for i, _ := range indices {
		indices[i].Filepath = filepath.Join(baseDir, "transactions.proto")
	}
}

// writeIndicesToFile serializes and writes a batch of indices to a specified JSON file.
func writeIndicesToFile(indices []interface{}, filePath string) { // small cheating
	// Open or create the file for appending
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // that cool
	if err != nil {
		log.Fatalf("Failed to open file %s for writing: %v", filePath, err)
	}
	defer file.Close()

	// Serialize and write each index in the batch
	encoder := json.NewEncoder(file)
	for _, index := range indices {
		if err := encoder.Encode(index); err != nil {
			log.Fatalf("Failed to encode index to %s: %v", filePath, err)
		}
	}
}

// Start initiates the crawling process for the configured blockchain.
func (c *Crawler) Start() {

	chainType := "ethereum"

	url := os.Getenv("INFURA_URL")

	fmt.Println("Using blockchain url:", url)

	log.Println("Starting crawler...")
	client, err := common.NewClient(chainType, url)

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

	fmt.Println("Latest block number:", latestBlockNumber)

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

		// write the blocks and transactions to disk as example 100001-100010/blocks.proto and 100001-100010/transactions.proto
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

		// Update filepaths for each index
		updateBlockIndexFilepaths(blockIndex, batchDir)
		updateTransactionIndexFilepaths(transactionIndex, batchDir)

		// Write updated indices to their respective files
		// Convert blockIndex to []interface{}
		interfaceBlockIndex := make([]interface{}, len(blockIndex))
		for i, v := range blockIndex {
			interfaceBlockIndex[i] = v
		}

		// Convert transactionIndex to []interface{}

		interfaceTransactionIndex := make([]interface{}, len(transactionIndex))
		for i, v := range transactionIndex {
			interfaceTransactionIndex[i] = v
		}

		writeIndicesToFile(interfaceBlockIndex, filepath.Join(batchDir, "blockIndex.json"))
		writeIndicesToFile(interfaceTransactionIndex, filepath.Join(batchDir, "transactionIndex.json"))

		// write the block and transaction indices to disk

		// open index general file for writing

		file_out := "index.db"

		file, err := os.Create(file_out)

		if err != nil {
			panic(err)
		}

		defer file.Close()

		// write the block and transaction indices to disk

		// Update the startBlock for the next batch
		startBlock = endBlock
	}
}

// You can add more methods here for additional functionalities
// such as handling reconnection logic, managing crawler state, etc.
