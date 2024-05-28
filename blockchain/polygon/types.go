package polygon

type PolygonBlockJson struct {
	Difficulty       string                          `json:"difficulty"`
	ExtraData        string                          `json:"extraData"`
	GasLimit         string                          `json:"gasLimit"`
	GasUsed          string                          `json:"gasUsed"`
	Hash             string                          `json:"hash"`
	LogsBloom        string                          `json:"logsBloom"`
	Miner            string                          `json:"miner"`
	MixHash          string                          `json:"mixHash"`
	Nonce            string                          `json:"nonce"`
	BlockNumber      string                          `json:"number"`
	ParentHash       string                          `json:"parentHash"`
	ReceiptRoot      string                          `json:"receiptRoot"`
	Sha3Uncles       string                          `json:"sha3Uncles"`
	StateRoot        string                          `json:"stateRoot"`
	Timestamp        string                          `json:"timestamp"`
	TotalDifficulty  string                          `json:"totalDifficulty"`
	TransactionsRoot string                          `json:"transactionsRoot"`
	Uncles           []string                        `json:"uncles"`
	Size             string                          `json:"size"`
	BaseFeePerGas    string                          `json:"baseFeePerGas"`
	IndexedAt        uint64                          `json:"indexed_at"`
	Transactions     []*PolygonSingleTransactionJson `json:"transactions"`
}

type PolygonSingleTransactionJson struct {
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

type PolygonSingleEventJson struct {
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
