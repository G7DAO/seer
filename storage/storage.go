package storage

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
)

type ListReturnFunc func(any) string

type Storer interface {
	Save(batchDir, filename string, bf bytes.Buffer) error
	Read(key string) (bytes.Buffer, error)
	ReadBatch(readItems []ReadItem) (map[string][]string, error)
	Delete(key string) error
	List(ctx context.Context, delim, blockBatch string, timeout int, returnFunc ListReturnFunc) ([]string, error)
}

type ReadItem struct {
	Key    string
	RowIds []uint64
}

func ReadFiles(keys []string, storageInstance Storer) ([]bytes.Buffer, error) {
	var result []bytes.Buffer

	for _, key := range keys {
		buf, err := storageInstance.Read(key)
		if err != nil {
			return nil, fmt.Errorf("failed to read object from bucket %s: %v", key, err)
		}

		result = append(result, buf)
	}

	return result, nil
}

func ReadFilesAsync(keys []string, threads int, storageInstance Storer) ([]bytes.Buffer, error) {
	var result []bytes.Buffer
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

			buf, err := storageInstance.Read(k)
			if err != nil {
				errChan <- fmt.Errorf("failed to read object from bucket %s: %v", k, err)
				return
			}

			mu.Lock()
			result = append(result, buf)
			mu.Unlock()
		}(key)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check if any errors occurred
	if len(errChan) > 0 {
		var errMsgs []string
		for err := range errChan {
			errMsgs = append(errMsgs, err.Error())
		}
		return result, fmt.Errorf("errors occurred during file reads:\n%s",
			strings.Join(errMsgs, "\n"))
	}

	return result, nil
}
