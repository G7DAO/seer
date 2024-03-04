package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	// "time"
)

type BlockJson struct {
	Difficulty       int64  `json:"difficulty"`
	ExtraData        string `json:"extraData"`
	GasLimit         int64  `json:"gasLimit"`
	GasUsed          int64  `json:"gasUsed"`
	Hash             string `json:"hash"`
	LogsBloom        string `json:"logsBloom"`
	Miner            string `json:"miner"`
	MixHash          string `json:"mixHash"`
	Nonce            string `json:"nonce"`
	BlockNumber      int64  `json:"number"`
	ParentHash       string `json:"parentHash"`
	ReceiptRoot      string `json:"receiptRoot"`
	Sha3Uncles       string `json:"sha3Uncles"`
	StateRoot        string `json:"stateRoot"`
	Timestamp        int64  `json:"timestamp"`
	TotalDifficulty  string `json:"totalDifficulty"`
	TransactionsRoot string `json:"transactionsRoot"`
	Uncles           string `json:"uncles"`
	Size             int32  `json:"size"`
	BaseFeePerGas    string `json:"baseFeePerGas"`
	IndexedAt        int64  `json:"indexed_at"`
}

type SingleTransactionJson struct {
	AccessList           []AccessList `json:"accessList"`
	BlockHash            string       `json:"blockHash"`
	BlockNumber          int64        `json:"blockNumber"`
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
	TransactionIndex     int64        `json:"transactionIndex"`
	TransactionType      int32        `json:"type"`
	Value                string       `json:"value"`
	YParity              string       `json:"yParity"`
	IndexedAt            int64        `json:"indexed_at"`
	BlockTimestamp       int64        `json:"block_timestamp"`
}

type AccessList struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

// SingleEvent represents a single event within a transaction
type SingleEventJson struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      int64    `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	BlockHash        string   `json:"blockHash"`
	Removed          bool     `json:"removed"`
	LogIndex         int64    `json:"logIndex"`
	TransactionIndex int64    `json:"transactionIndex"`
}

// ReadJsonBlocks reads blocks from a JSON file
func ReadJsonBlocks() []*BlockJson {
	file, err := os.Open("data/blocks_go.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var blocks []*BlockJson
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&blocks)
	if err != nil {
		fmt.Println("error:", err)
	}

	return blocks
}

// ReadJsonTransactions reads transactions from a JSON file
func ReadJsonTransactions() []*SingleTransactionJson {

	file, err := os.Open("data/transactions_go.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var transactions []*SingleTransactionJson
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&transactions)
	if err != nil {
		fmt.Println("error:", err)
	}

	return transactions
}

// ReadJsonEventLogs reads event logs from a JSON file
func ReadJsonEventLogs() []*SingleEventJson {

	file, err := os.Open("data/event_logs_go.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var eventLogs []*SingleEventJson

	decoder := json.NewDecoder(file)

	fmt.Println("decoder:", decoder)

	err = decoder.Decode(&eventLogs)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("eventLogs:", eventLogs[0])

	return eventLogs
}
