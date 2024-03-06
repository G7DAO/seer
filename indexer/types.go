package indexer

// gorm is a Go ORM library for working with databases

// Define the interface for handling general index types such as blocks, transactions, and logs for fast access to blockchain data.

type BlockIndex struct {
	chain          string
	BlockNumber    uint64
	BlockHash      string
	BlockTimestamp uint64
	ParentHash     string
	Filepath       string
}

func (b BlockIndex) TableName() string {
	// This is the table name in the database for particular chain
	return b.chain + "_block_index"
}

func (b *BlockIndex) SetChain(chain string) {
	b.chain = chain
}

// NewBlockIndex creates a new instance of BlockIndex with the chain set.
func NewBlockIndex(chain string, blockNumber uint64, blockHash string, blockTimestamp uint64, parentHash string, filepath string) BlockIndex {
	return BlockIndex{
		chain:          chain,
		BlockNumber:    blockNumber,
		BlockHash:      blockHash,
		BlockTimestamp: blockTimestamp,
		ParentHash:     parentHash,
		Filepath:       filepath,
	}
}

type TransactionIndex struct {
	chain                string
	BlockNumber          uint64
	BlockHash            string
	BlockTimestamp       uint64
	TransactionHash      string // TODO: Rename this to Hash
	TransactionIndex     uint64 // TODO: Rename this to Index
	TransactionTimestamp uint64 // TODO: Remove this field
	Filepath             string
}

func (t TransactionIndex) TableName() string {
	// This is the table name in the database for particular chain
	return t.chain + "_transaction_index"
}

func (t TransactionIndex) SetChain(chain string) {
	t.chain = chain
}

func NewTransactionIndex(chain string, blockNumber uint64, blockHash string, blockTimestamp uint64, transactionHash string, transactionIndex uint64, transactionTimestamp uint64, filepath string) TransactionIndex {
	return TransactionIndex{
		chain:                chain,
		BlockNumber:          blockNumber,
		BlockHash:            blockHash,
		BlockTimestamp:       blockTimestamp,
		TransactionHash:      transactionHash,
		TransactionIndex:     transactionIndex,
		TransactionTimestamp: transactionTimestamp,
		Filepath:             filepath,
	}
}

type LogIndex struct {
	chain           string
	BlockNumber     uint64
	BlockHash       string
	BlockTimestamp  uint64
	TransactionHash string
	Topic0          string // TODO: 1) Add Topic1, Topic2. 2) Rename Topic0 to selector
	LogIndex        uint64
	Filepath        string
}

func (l LogIndex) TableName() string {
	// This is the table name in the database for particular chain
	return l.chain + "_log_index"
}

func (l LogIndex) SetChain(chain string) {
	l.chain = chain
}

func NewLogIndex(chain string, blockNumber uint64, blockHash string, transactionHash string, BlockTimestamp uint64, topic0 string, logIndex uint64, filepath string) LogIndex {
	return LogIndex{
		chain:           chain,
		BlockNumber:     blockNumber,
		BlockHash:       blockHash,
		BlockTimestamp:  BlockTimestamp,
		TransactionHash: transactionHash,
		Topic0:          topic0,
		LogIndex:        logIndex,
		Filepath:        filepath,
	}
}

type IndexType string

const (
	BlockIndexType       IndexType = "block"
	TransactionIndexType IndexType = "transaction"
	LogIndexType         IndexType = "log"
)

type BlockCahche struct {
	BlockNumber    uint64
	BlockTimestamp uint64
	BlockHash      string
}

type abiJob struct {
	ABI             string
	Chain           string
	ContractAddress string
}
