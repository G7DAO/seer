package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"gorm.io/gorm/clause"
)

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
	fmt.Println("Writing index to database")
	// chain to table name
	// block -> ethereum_block_index
	// transaction -> ethereum_transaction_index
	// log -> ethereum_log_index

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
			fmt.Println(index.chain)
			index.chain = blockchain
			blockIndexes = append(blockIndexes, index)
		}
		return writeBlockIndexToDB(blockchain+"_block_index", blockIndexes)

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
		return writeTransactionIndexToDB(blockchain+"_transaction_index", transactionIndexes)

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
		return writeLogIndexToDB(blockchain+"_log_index", logIndexes)

	default:
		return errors.New("unsupported index type")
	}
}

func writeBlockIndexToDB(tableName string, indexes []BlockIndex) error {
	// Start a transaction

	fmt.Println("Writing block index to database")
	fmt.Println(indexes[0].TableName())
	tx := DB.Begin()
	if tx.Error != nil {
		fmt.Println("Error starting transaction: ", tx.Error)
		return fmt.Errorf("starting transaction: %w", tx.Error)

	}

	// Attempt to insert indexes, handling conflicts as needed
	result := tx.Table(indexes[0].TableName()).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&indexes)
	if result.Error != nil {
		tx.Rollback() // Rollback on error
		fmt.Println("Error inserting block indexes: ", result.Error)
		return fmt.Errorf("inserting block indexes: %w", result.Error)
	}

	// Commit the transaction if all went well
	if commitErr := tx.Commit().Error; commitErr != nil {
		fmt.Println("Error committing transaction: ", commitErr)
		return fmt.Errorf("committing transaction: %w", commitErr)
	}

	return nil // Return nil on success
}

func writeTransactionIndexToDB(tableName string, indexes []TransactionIndex) error {
	// Insert with dedication to the database

	tx := DB.Begin()
	if tx.Error != nil {
		fmt.Println("Error starting transaction: ", tx.Error)
		return fmt.Errorf("starting transaction: %w", tx.Error)
	}

	// Attempt to insert indexes, handling conflicts as needed
	result := tx.Table(tableName).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&indexes)
	if result.Error != nil {
		tx.Rollback() // Rollback on error
		fmt.Println("Error inserting transaction indexes: ", result.Error)
		return fmt.Errorf("inserting transaction indexes: %w", result.Error)
	}

	// Commit the transaction if all went well
	if commitErr := tx.Commit().Error; commitErr != nil {
		DB.Rollback()
		fmt.Println("Error committing transaction: ", commitErr)
		return fmt.Errorf("committing transaction: %w", commitErr)
	}

	return nil
}

func writeLogIndexToDB(tableName string, indexes []LogIndex) error {

	fmt.Println("Writing log index to database")

	tx := DB.Begin()
	if tx.Error != nil {
		fmt.Println("Error starting transaction: ", tx.Error)
		return fmt.Errorf("starting transaction: %w", tx.Error)
	}

	batchSize := 1000
	for i := 0; i < len(indexes); i += batchSize {
		end := i + batchSize
		if end > len(indexes) {
			end = len(indexes)
		}

		// Insert the current batch of indexes.
		result := tx.Table(tableName).Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(indexes[i:end])

		if result.Error != nil {
			tx.Rollback() // Rollback on error
			fmt.Println("Error inserting log indexes: ", result.Error)
			return fmt.Errorf("inserting log indexes: %w", result.Error)
		}
	}

	// Commit the transaction if all went well
	if commitErr := tx.Commit().Error; commitErr != nil {
		fmt.Println("Error committing transaction: ", commitErr)
		return fmt.Errorf("committing transaction: %w", commitErr)
	}

	return nil
}

// ReadIndexFromDatabase reads indices from the database based on the type and block range.
func ReadIndexFromDatabase(chain string, addresses, startBlock, endBlock uint64) ([]interface{}, error) {

	var indexRecords []interface{} // Assuming IndexRecord is the correct struct

	// Assuming DB is correctly set up and addresses is a slice that can be used in the query

	event_table_name := chain + "_log_indexes"

	records := DB.Table(event_table_name).
		Joins("JOIN transactions ON transactions.hash = (?).transaction_hash", event_table_name).
		Where("log_indexes.block_number >= ? AND log_indexes.block_number <= ? AND transactions.to IN (?)", startBlock, endBlock, addresses).
		Find(&indexRecords)

	if records.Error != nil {
		return nil, records.Error
	}

	var indices []interface{}
	for _, record := range indexRecords { // Now iterating over the slice of results
		indices = append(indices, record)
	}

	return indices, nil

}
