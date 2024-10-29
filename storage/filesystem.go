package storage

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileStorage struct {
	BasePath string
}

func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{BasePath: basePath}
}

func (fs *FileStorage) Save(batchDir, filename string, bf bytes.Buffer) error {
	keyDir := filepath.Join(fs.BasePath, batchDir)
	key := filepath.Join(keyDir, filename)

	// Check if the directory exists
	// If not, create it
	if _, err := os.Stat(keyDir); os.IsNotExist(err) {
		os.MkdirAll(keyDir, os.ModePerm)
	}

	file, err := os.OpenFile(key, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // that cool

	if err != nil {
		log.Fatalf("Failed to open file %s: %v", key, err)
	}

	defer file.Close()

	_, err = io.Copy(file, &bf)

	if err != nil {
		log.Fatalf("Failed to write to file %s: %v", key, err)
	}

	return nil

}

func (fs *FileStorage) Read(key string) (bytes.Buffer, error) {

	file, err := os.Open(key)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to open file %s: %v", key, err)
	}
	defer file.Close()

	var bf bytes.Buffer
	_, err = io.Copy(&bf, file)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to read file %s: %v", key, err)
	}

	return bf, nil
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
			lines := bytes.Split(data.Bytes(), []byte{'\n'})
			for _, line := range lines {
				if len(line) > 0 {
					result[item.Key] = append(result[item.Key], string(line))
				}
			}
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

func (fs *FileStorage) List(ctx context.Context, delim, blockBatch string, timeout int, returnFunc ListReturnFunc) ([]string, error) {
	prefix := fmt.Sprintf("%s/", fs.BasePath)
	log.Printf("Loading directory items with prefix: %s", prefix)

	dirs, readDirErr := os.ReadDir(prefix)
	if readDirErr != nil {
		return []string{}, readDirErr
	}

	var items []string
	itemsLen := 0

	for _, d := range dirs {
		items = append(items, fmt.Sprintf("%s%s/", prefix, d.Name()))
		itemsLen++
	}

	log.Printf("Listed %d items", itemsLen)

	return items, nil
}

func (fs *FileStorage) Delete(key string) error {

	// Implement the Delete method
	return nil
}

func (fs *FileStorage) ReadFiles(keys []string) (bytes.Buffer, error) {

	var data bytes.Buffer

	for _, key := range keys {
		data, err := fs.Read(key)
		if err != nil {
			return bytes.Buffer{}, err
		}

		if _, err := io.Copy(&data, &data); err != nil {
			return bytes.Buffer{}, fmt.Errorf("failed to read object data: %v", err)
		}

	}
	return data, nil
}

func (fs *FileStorage) ReadFilesAsync(keys []string, threads int) (bytes.Buffer, error) {
	var data bytes.Buffer
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(keys))

	// Semaphore to limit the number of concurrent reads
	sem := make(chan struct{}, threads)

	for _, key := range keys {
		wg.Add(1)
		sem <- struct{}{}
		go func(k string) {
			defer func() {
				<-sem
				wg.Done()
			}()
			bf, err := fs.Read(k)
			if err != nil {
				errChan <- fmt.Errorf("failed to read file %s: %v", k, err)
				return
			}
			mu.Lock()
			_, writeErr := data.Write(bf.Bytes())
			mu.Unlock()
			if writeErr != nil {
				errChan <- fmt.Errorf("failed to write data for file %s: %v", k, writeErr)
			}
		}(key)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check if any errors occurred
	if len(errChan) > 0 {
		// Collect all errors
		var errMsgs []string
		for err := range errChan {
			errMsgs = append(errMsgs, err.Error())
		}
		return data, fmt.Errorf("errors occurred during file reads:\n%s",
			strings.Join(errMsgs, "\n"))
	}

	return data, nil
}
