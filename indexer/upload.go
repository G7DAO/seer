package indexer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

var DBConnection *PostgreSQLpgx
var err error

// init initializes the DBConnection variable with a new PostgreSQLpgx instance.
func InitDBConnection() {

	DBConnection, err = NewPostgreSQLpgx(MOONSTREAM_DB_V3_INDEXES_URI)

	if err != nil {
		fmt.Println("Error initializing DBConnection: ", err)
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

	return nil
}

// WriteIndicesToDatabase writes the given indices to the database
func WriteIndicesToDatabase(blockchain string, blocks []BlockIndex) error {
	// Write block indices

	return DBConnection.WriteIndexes(blockchain, blocks)
}
