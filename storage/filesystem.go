package storage

import (
	"bufio"
	"log"
	"os"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{basePath: basePath}
}

func (fs *FileStorage) Save(key string, data []string) error {

	// Check if the directory exists
	// If not, create it
	if _, err := os.Stat(fs.basePath); os.IsNotExist(err) {
		os.MkdirAll(fs.basePath, os.ModePerm)
	}

	file, err := os.OpenFile(key, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // that cool

	if err != nil {
		log.Fatalf("Failed to open file %s: %v", key, err)
	}

	defer file.Close()

	for _, item := range data {
		_, err := file.WriteString(item + "\n")
		if err != nil {
			return err
		}
	}

	return nil

}

func (fs *FileStorage) Read(key string) ([]string, error) {

	file, err := os.Open(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result []string

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (fs *FileStorage) ReadBatch(readItems []ReadItem) (map[string][]string, error) {

	result := make(map[string][]string)

	for _, item := range readItems {
		if result[item.Key] == nil {
			result[item.Key] = make([]string, 0)
		}

		if len(item.RowIds) == 0 {
			data, err := fs.Read(item.Key)
			if err != nil {
				return nil, err
			}
			result[item.Key] = append(result[item.Key], data...)
		} else {
			file, err := os.Open(item.Key)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			reader := bufio.NewReader(file)
			rowMap := make(map[uint64]bool)
			for _, id := range item.RowIds {
				rowMap[id] = true
			}
			var currentRow uint64 = 0 // assuming rows are 1-indexed
			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err.Error() != "EOF" {
						return nil, err
					}
					break
				}
				if rowMap[currentRow] {
					result[item.Key] = append(result[item.Key], string(line))
				}
				currentRow++
			}
			// Append the output if
		}
	}

	return result, nil
}

func (fs *FileStorage) Delete(key string) error {

	// Implement the Delete method
	return nil
}
