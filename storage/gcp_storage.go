package storage

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

// GCS implements the Storer interface for Google Cloud Storage
type GCS struct {
	Client   *storage.Client
	BasePath string
}

// NewGCSStorage initializes a GCS storage with the provided client
func NewGCSStorage(client *storage.Client, basePath string) *GCS {
	return &GCS{
		Client:   client,
		BasePath: basePath,
	}
}

func (g *GCS) Save(batchDir, filename string, data []string) error {
	key := filepath.Join(g.BasePath, batchDir, filename)

	ctx := context.Background()

	bucket := g.Client.Bucket(SeerCrawlerStorageBucket)

	obj := bucket.Object(key)

	// Write the data as a string
	w := obj.NewWriter(ctx)

	for _, item := range data {
		if _, err := w.Write([]byte(item + "\n")); err != nil {
			return fmt.Errorf("failed to write object to bucket: %v", err)
		}
	}

	wc := obj.NewWriter(ctx)
	wc.Metadata = map[string]string{
		"rows": fmt.Sprintf("%d", len(data)),
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}

func (g *GCS) Read(key string) ([]string, error) {

	ctx := context.Background()

	bucket := g.Client.Bucket(SeerCrawlerStorageBucket)

	obj := bucket.Object(key)

	r, err := obj.NewReader(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %v", err)
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)

	var result []string

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read object from bucket: %v", err)
	}

	return result, nil
}

func (g *GCS) Delete(key string) error {

	ctx := context.Background()

	bucket := g.Client.Bucket(SeerCrawlerStorageBucket)

	obj := bucket.Object(key)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object from bucket: %v", err)
	}

	return nil

}

func (g *GCS) ReadBatch(readItems []ReadItem) (map[string][]string, error) {
	ctx := context.Background()

	bucket := g.Client.Bucket(SeerCrawlerStorageBucket)

	result := make(map[string][]string)

	for _, item := range readItems {
		obj := bucket.Object(item.Key)

		r, err := obj.NewReader(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create reader: %v", err)
		}
		defer r.Close()

		reader := bufio.NewReader(r)

		rowMap := make(map[uint64]bool)
		for _, id := range item.RowIds {
			rowMap[id] = true
		}

		var currentRow uint64 = 0

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("failed to read object from bucket: %v", err)
			}

			// Remove the newline character from the end of the line if it exists
			line = strings.TrimSuffix(line, "\n")

			if rowMap[currentRow] {
				if result[item.Key] == nil {
					result[item.Key] = make([]string, 0)
				}
				result[item.Key] = append(result[item.Key], line)
			}
			currentRow++
		}
	}

	return result, nil
}
