package indexer

import "time"

// gorm is a Go ORM library for working with databases

// Define the interface for handling general index types such as blocks, transactions, and logs for fast access to blockchain data.

type BlockIndex struct {
	chain          string
	BlockNumber    uint64
	BlockHash      string
	BlockTimestamp uint64
	ParentHash     string
	RowID          uint64
	Path           string
	L1BlockNumber  uint64
}

func (b *BlockIndex) SetChain(chain string) {
	b.chain = chain
}

// NewBlockIndex creates a new instance of BlockIndex with the chain set.
func NewBlockIndex(chain string, blockNumber uint64, blockHash string, blockTimestamp uint64, parentHash string, row_id uint64, path string, l1BlockNumber uint64) BlockIndex {
	return BlockIndex{
		chain:          chain,
		BlockNumber:    blockNumber,
		BlockHash:      blockHash,
		BlockTimestamp: blockTimestamp,
		ParentHash:     parentHash,
		RowID:          row_id,
		Path:           path,
		L1BlockNumber:  l1BlockNumber,
	}
}

type IndexType string

const (
	BlockIndexType       IndexType = "blocks"
	TransactionIndexType IndexType = "transactions"
	LogIndexType         IndexType = "logs"
)

type BlockCache struct {
	BlockNumber    uint64
	BlockTimestamp uint64
	BlockHash      string
}

type AbiJob struct {
	ID                    string
	Address               []byte
	UserID                string
	CustomerID            string
	AbiSelector           string
	Chain                 string
	AbiName               string
	Status                string
	HistoricalCrawlStatus string
	Progress              int
	TaskPickedup          bool
	Abi                   string
	AbiType               string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeploymentBlockNumber *uint64
}

type CustomerUpdates struct {
	CustomerID string                                  `json:"customer_id"`
	Abis       map[string]map[string]map[string]string `json:"abis"`
	LastBlock  uint64                                  `json:"last_block"`
	Path       string                                  `json:"path"`
}

type TaskForTransaction struct {
	Address  string `json:"address"`
	Selector string `json:"selector"`
	ABI      string `json:"abi"`
}

type TaskForLog struct { // Assuming structure similar to TransactionIndex
	Address  string `json:"address"`
	Selector string `json:"selector"`
	ABI      string `json:"abi"`
}

type EventLabel struct {
	Address         string
	BlockNumber     uint64
	BlockHash       string
	CallerAddress   string
	Label           string
	LabelName       string
	LabelType       string
	OriginAddress   string
	TransactionHash string
	LabelData       string
	BlockTimestamp  uint64
	LogIndex        uint64
}

type TransactionLabel struct {
	Address         string
	BlockNumber     uint64
	BlockHash       string
	CallerAddress   string
	Label           string
	LabelName       string
	LabelType       string
	OriginAddress   string
	TransactionHash string
	LabelData       string
	BlockTimestamp  uint64
}

type AbiJobsDeployInfo struct {
	DeployedBlockNumber uint64
	IDs                 []string
}
