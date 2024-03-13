package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var Actions *PostgreSQLpgx
var err error

// init initializes the Actions variable with a new PostgreSQLpgx instance.
func InitDBConnection() {

	Actions, err = NewPostgreSQLpgx()

	if err != nil {
		fmt.Println("Error initializing Actions: ", err)
	}
}

func WriteIndicesToFile(indices []interface{}, filePath string) error {

	IndexType := strings.Split(filePath, "_")[0]

	existingIndices, err := ReadIndicesFromFile(filePath)
	if err != nil {
		return err
	}

	// Filter out duplicates
	newIndices := filterDuplicates(IndexType, existingIndices, indices)
	if len(newIndices) == 0 {
		return nil // No new indices to write
	}

	// Open or create the file for appending
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serialize and write each new index to the file
	encoder := json.NewEncoder(file)
	for _, index := range newIndices {
		if err := encoder.Encode(index); err != nil {
			return err
		}
	}
	return nil
}

// ReadIndicesFromFile reads indices from a JSON file into a slice of interface{}.
func ReadIndicesFromFile(filePath string) ([]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var indices []interface{}
	decoder := json.NewDecoder(file)
	for {
		var index interface{}
		if err := decoder.Decode(&index); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		indices = append(indices, index)
	}
	return indices, nil
}

// filterDuplicates removes duplicates from the newIndices slice  based on existingIndices.
func filterDuplicates(indexType string, existingIndices, newIndices []interface{}) []interface{} {
	// Implement logic to identify and remove duplicates

	// if indexType == "block" {
	if indexType == "block" {
		// Convert the existing indices to a map for faster lookups
		var existingBlockIndices map[uint64]BlockIndex
		for _, index := range existingIndices {
			blockIndex := index.(BlockIndex)
			existingBlockIndices[blockIndex.BlockNumber] = blockIndex
		}

		// Filter out new indices that already exist
		var filteredIndices []interface{}
		for _, index := range newIndices {
			blockIndex := index.(BlockIndex)
			if _, exists := existingBlockIndices[blockIndex.BlockNumber]; !exists {
				filteredIndices = append(filteredIndices, index)
			}

		}
		return filteredIndices
	}

	if indexType == "transaction" {
		// Convert the existing indices to a map for faster lookups
		var existingTransactionIndices map[string]TransactionIndex
		for _, index := range existingIndices {
			transactionIndex := index.(TransactionIndex)
			existingTransactionIndices[transactionIndex.TransactionHash] = transactionIndex
		}

		// Filter out new indices that already exist
		var filteredIndices []interface{}
		for _, index := range newIndices {
			transactionIndex := index.(TransactionIndex)
			if _, exists := existingTransactionIndices[transactionIndex.TransactionHash]; !exists {
				filteredIndices = append(filteredIndices, index)
			}
		}
		return filteredIndices
	}

	if indexType == "log" {
		// Convert the existing indices to a map for faster lookups
		// Unique identifier for log indices is the transaction hash and log index
		var existingLogIndices map[string]LogIndex
		for _, index := range existingIndices {
			logIndex := index.(LogIndex)
			uniqueKey := logIndex.TransactionHash + string(logIndex.LogIndex)
			existingLogIndices[uniqueKey] = logIndex
		}

		// Filter out new indices that already exist
		var filteredIndices []interface{}
		for _, index := range newIndices {
			logIndex := index.(LogIndex)
			uniqueKey := logIndex.TransactionHash + string(logIndex.LogIndex)
			if _, exists := existingLogIndices[uniqueKey]; !exists {
				filteredIndices = append(filteredIndices, index)
			}
		}

		return filteredIndices
	}

	return nil
}

func WriteIndexesToDatabase(blockchain string, indexes []interface{}, indexType string) error {
	// Implement logic to write indexes to the database

	if len(indexes) == 0 {
		return nil
	}

	switch indexType {
	case "block":
		var blockIndexes []BlockIndex
		for _, i := range indexes {
			index, ok := i.(BlockIndex)
			if !ok {
				return errors.New("invalid type for block index")
			}
			index.chain = blockchain
			blockIndexes = append(blockIndexes, index)
		}
		return Actions.writeBlockIndexToDB(blockchain+"_blocks", blockIndexes)

	case "transaction":
		var transactionIndexes []TransactionIndex
		for _, i := range indexes {
			index, ok := i.(TransactionIndex)
			if !ok {
				return errors.New("invalid type for transaction index")
			}
			index.chain = blockchain
			transactionIndexes = append(transactionIndexes, index)
		}
		return Actions.writeTransactionIndexToDB(blockchain+"_transactions", transactionIndexes)

	case "log":
		var logIndexes []LogIndex
		for _, i := range indexes {
			index, ok := i.(LogIndex)
			if !ok {
				return errors.New("invalid type for log index")
			}
			index.chain = blockchain
			logIndexes = append(logIndexes, index)
		}
		return Actions.writeLogIndexToDB(blockchain+"_logs", logIndexes)

	default:
		return errors.New("unsupported index type")
	}
}
