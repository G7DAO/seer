package ethereum

/// Block represents an Ethereum block.

// type Block struct {
// 	Number           *big.Int       `json:"number"`           // The block number. nil when its pending block.
// 	Hash             common.Hash    `json:"hash"`             // Hash of the block. nil when its pending block.
// 	ParentHash       common.Hash    `json:"parentHash"`       // Hash of the parent block.
// 	Nonce            *big.Int       `json:"nonce"`            // Hash of the generated proof-of-work. nil when its pending block.
// 	Sha3Uncles       common.Hash    `json:"sha3Uncles"`       // SHA3 of the uncles data in the block.
// 	LogsBloom        common.Hash    `json:"logsBloom"`        // The bloom filter for the logs of the block. nil when its pending block.
// 	TransactionsRoot common.Hash    `json:"transactionsRoot"` // The root of the transaction trie of the block.
// 	StateRoot        common.Hash    `json:"stateRoot"`        // The root of the final state trie of the block.
// 	ReceiptsRoot     common.Hash    `json:"receiptsRoot"`     // The root of the receipts trie of the block.
// 	Miner            common.Address `json:"miner"`            // The address of the beneficiary to whom the mining rewards were given.
// 	Difficulty       *big.Int       `json:"difficulty"`       // Integer of the difficulty for this block.
// 	TotalDifficulty  *big.Int       `json:"totalDifficulty"`  // Integer of the total difficulty of the chain until this block.
// 	ExtraData        []byte         `json:"extraData"`        // The "extra data" field of this block.
// 	Size             *big.Int       `json:"size"`             // Integer the size of this block in bytes.
// 	GasLimit         uint64         `json:"gasLimit"`         // The maximum gas allowed in this block.
// 	GasUsed          uint64         `json:"gasUsed"`          // The total used gas by all transactions in this block.
// 	Timestamp        *big.Int       `json:"timestamp"`        // The unix timestamp for when the block was collated.
// 	Transactions     []Transaction  `json:"transactions"`     // Array of transaction objects, or 32 Bytes transaction hashes depending on the last given parameter.
// 	Uncles           []common.Hash  `json:"uncles"`           // Array of uncle hashes.
// }

// /// Transaction represents an Ethereum transaction.

// type Transaction struct {
// 	BlockHash        common.Hash     `json:"blockHash"`        // Hash of the block where this transaction was in.
// 	BlockNumber      *big.Int        `json:"blockNumber"`      // Block number where this transaction was in.
// 	From             common.Address  `json:"from"`             // Address of the sender.
// 	Gas              *big.Int        `json:"gas"`              // Gas provided by the sender.
// 	GasPrice         *big.Int        `json:"gasPrice"`         // Gas price provided by the sender in wei.
// 	Hash             common.Hash     `json:"hash"`             // Hash of the transaction.
// 	Input            []byte          `json:"input"`            // Input data.
// 	Nonce            uint64          `json:"nonce"`            // Nonce of the sender.
// 	To               *common.Address `json:"to"`               // Address of the receiver. nil when it's a contract creation transaction.
// 	TransactionIndex uint64          `json:"transactionIndex"` // Index of the transaction in the block.
// 	Value            *big.Int        `json:"value"`            // Value transferred in wei.
// 	V                *big.Int        `json:"v"`                // ECDSA recovery id.
// 	R                *big.Int        `json:"r"`                // ECDSA signature r.
// 	S                *big.Int        `json:"s"`                // ECDSA signature s.
// }

// /// EventLog represents a log of an Ethereum event.

// type EventLog struct {
// 	Address          common.Address `json:"address"`          // The address of the contract that generated the log.
// 	Topics           []common.Hash  `json:"topics"`           // Topics are indexed parameters during log generation.
// 	Data             []byte         `json:"data"`             // The data field from the log.
// 	BlockNumber      *big.Int       `json:"blockNumber"`      // The block number where this log was in.
// 	TransactionHash  common.Hash    `json:"transactionHash"`  // The hash of the transaction that generated this log.
// 	BlockHash        common.Hash    `json:"blockHash"`        // The hash of the block where this log was in.
// 	Removed          bool           `json:"removed"`          // True if the log was reverted due to a chain reorganization.
// 	LogIndex         uint64         `json:"logIndex"`         // The index of the log in the block.
// 	TransactionIndex uint64         `json:"transactionIndex"` // The index of the transaction in the block.
// }

// /// Trace represents

// type Trace struct {
// 	Action              json.RawMessage `json:"action"`
// 	BlockHash           common.Hash     `json:"blockHash"`
// 	BlockNumber         *big.Int        `json:"blockNumber"`
// 	Result              *TraceResult    `json:"result,omitempty"`
// 	SubTraces           int             `json:"subtraces"`
// 	TraceAddress        []int           `json:"traceAddress"`
// 	TransactionHash     common.Hash     `json:"transactionHash"`
// 	TransactionPosition uint64          `json:"transactionPosition"`
// 	Type                string          `json:"type"`
// }

// /// TraceActionCall represents a call action in a trace.

// type TraceActionCall struct {
// 	From     common.Address `json:"from"`
// 	CallType string         `json:"callType"`
// 	Gas      *big.Int       `json:"gas"`
// 	Input    string         `json:"input"`
// 	To       common.Address `json:"to"`
// 	Value    *big.Int       `json:"value"`
// }

// /// TraceActionReward represents a reward action in a trace.

// type TraceActionReward struct {
// 	Author     common.Address `json:"author"`
// 	RewardType string         `json:"rewardType"`
// 	Value      *big.Int       `json:"value"`
// }

// /// TraceResult represents the result of a trace.

// type TraceResult struct {
// 	GasUsed *big.Int `json:"gasUsed"`
// 	Output  string   `json:"output"`
// }

// /// Struct for store write opcodes and their results storage value and corresponding gas costs.

// type WriteOpcodes struct {
// 	Opcode      string
// 	StorageCost uint64
// 	Pointer     string
// 	Value       string
// }

// Block represents a single blockchain block
/*
   {
       "difficulty": "0",
       "extraData": "0x4070656e6775696e6275696c642e6f7267",
       "gasLimit": 30000000,
       "gasUsed": 17127420,
       "hash": "0xe4a1232ad6aed6dccaa49b9924e1fca0b86f3631bab445475a3cfb3e95fbe61f",
       "logsBloom": "niFoSk64c1NeIHMlovzRqLkBAGg8ZQdlG/nhHzOs5HEXid3AzwYjK3YQOJf20tUVXmHVooqiPvgWJ1WvGX73H+lj/o1Zi4q5+o5fb7qYkfjDm3AZxWTY3x0sHVC6cPpHukjO4SrmusPIytSzyUguT8g2jqppnlyh+ktDMQbsWna6aqtbL9HByAB0eVUSBooefEOc9c93SOp0LjlCQLIOLfrj/VFPQ+36PNdN3py3nBUNJ8weZT7oh/ClO8SXYo82kMMoip3nnptvOE0EVClSxSV4TyK2uLdXZA1FJ2Ir8B2h+bqvpD5lHVhszoMzsASQ9fOfqD6fIXBL+YkrO759TQ==",
       "miner": "0xf15689636571dba322b48E9EC9bA6cFB3DF818e1",
       "mixHash": "0xe988b72d2414c2f939e3ae5e3d3b072853723f76c9e7359e4f422acf387addeb",
       "nonce": "0x0",
       "number": "19034382",
       "parentHash": "0x94d207044c5242dfd8a29312239dfeca3d9b6466566320e308a528b973865876",
       "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
       "stateRoot": "0x67571966e515b08151b69cf4d99fb088c1bda61afb68adb6ce581f1a24922278",
       "timestamp": 1705588583,
       "transactions": null,
       "transactionsRoot": "0x7314eecfe4432fdbe6018ddab237afed7a96d32885dc64486dc5632cde350ee6"
   },
*/

type BlockJson struct {
	// difficulty originally a hex string
	Difficulty       string                   `json:"difficulty"`
	ExtraData        string                   `json:"extraData"`
	GasLimit         string                   `json:"gasLimit"`
	GasUsed          string                   `json:"gasUsed"`
	Hash             string                   `json:"hash"`
	LogsBloom        string                   `json:"logsBloom"`
	Miner            string                   `json:"miner"`
	MixHash          string                   `json:"mixHash"`
	Nonce            string                   `json:"nonce"`
	BlockNumber      string                   `json:"number"`
	ParentHash       string                   `json:"parentHash"`
	ReceiptRoot      string                   `json:"receiptRoot"`
	Sha3Uncles       string                   `json:"sha3Uncles"`
	StateRoot        string                   `json:"stateRoot"`
	Timestamp        string                   `json:"timestamp"`
	TotalDifficulty  string                   `json:"totalDifficulty"`
	TransactionsRoot string                   `json:"transactionsRoot"`
	Uncles           []string                 `json:"uncles"`
	Size             string                   `json:"size"`
	BaseFeePerGas    string                   `json:"baseFeePerGas"`
	IndexedAt        uint64                   `json:"indexed_at"`
	Transactions     []*SingleTransactionJson `json:"transactions"`
}

type SingleTransactionJson struct {
	AccessList           []AccessList `json:"accessList"`
	BlockHash            string       `json:"blockHash"`
	BlockNumber          string       `json:"blockNumber"`
	ChainId              string       `json:"chainId"`
	FromAddress          string       `json:"from"`
	Gas                  string       `json:"gas"`
	GasPrice             string       `json:"gasPrice"`
	Hash                 string       `json:"hash"`
	Input                string       `json:"input"`
	MaxFeePerGas         string       `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string       `json:"maxPriorityFeePerGas"`
	Nonce                string       `json:"nonce"`
	R                    string       `json:"r"`
	S                    string       `json:"s"`
	ToAddress            string       `json:"to"`
	TransactionIndex     string       `json:"transactionIndex"`
	TransactionType      string       `json:"type"`
	Value                string       `json:"value"`
	YParity              string       `json:"yParity"`
	IndexedAt            string       `json:"indexed_at"`
	BlockTimestamp       string       `json:"block_timestamp"`
}

type AccessList struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

// Address          common.Address `json:"address"`          // The address of the contract that generated the log.
// Topics           []common.Hash  `json:"topics"`           // Topics are indexed parameters during log generation.
// Data             []byte         `json:"data"`             // The data field from the log.
// BlockNumber      *big.Int       `json:"blockNumber"`      // The block number where this log was in.
// TransactionHash  common.Hash    `json:"transactionHash"`  // The hash of the transaction that generated this log.
// BlockHash        common.Hash    `json:"blockHash"`        // The hash of the block where this log was in.
// Removed          bool           `json:"removed"`          // True if the log was reverted due to a chain reorganization.
// LogIndex         uint64         `json:"logIndex"`         // The index of the log in the block.
// TransactionIndex uint64         `json:"transactionIndex"` // The index of the transaction in the block.

// SingleEvent represents a single event within a transaction
type SingleEventJson struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	TransactionHash  string   `json:"transactionHash"`
	BlockNumber      string   `json:"blockNumber"`
	BlockHash        string   `json:"blockHash"`
	Removed          bool     `json:"removed"`
	LogIndex         string   `json:"logIndex"`
	TransactionIndex string   `json:"transactionIndex"`
}

type QueryFilter struct {
	BlockHash string     `json:"blockHash"`
	FromBlock string     `json:"fromBlock"`
	ToBlock   string     `json:"toBlock"`
	Address   []string   `json:"address"`
	Topics    [][]string `json:"topics"`
}
