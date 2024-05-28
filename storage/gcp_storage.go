package storage

import (
	"bufio"
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

// NewGCSStorage initializes a GCS storage with the provided client
func NewGCSStorage(client *storage.Client, basePath string) *GCS {
	return &GCS{
		client:   client,
		basePath: basePath,
	}
}

// GCS implements the Storer interface for Google Cloud Storage

type GCS struct {
	client   *storage.Client
	basePath string
}

func (g *GCS) Save(key string, data []string) error {

	ctx := context.Background()

	bucket := g.client.Bucket(SeerCrawlerStorageBucket)

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

	bucket := g.client.Bucket(SeerCrawlerStorageBucket)

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

	bucket := g.client.Bucket(SeerCrawlerStorageBucket)

	obj := bucket.Object(key)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object from bucket: %v", err)
	}

	return nil

}

func (g *GCS) ReadBatch(readItems []ReadItem) (map[string][]string, error) {
	ctx := context.Background()

	bucket := g.client.Bucket(SeerCrawlerStorageBucket)

	result := make(map[string][]string)

	for _, item := range readItems {

		obj := bucket.Object(item.Key)

		r, err := obj.NewReader(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create reader: %v", err)
		}
		defer r.Close()

		scanner := bufio.NewScanner(r)

		rowMap := make(map[uint64]bool)
		for _, id := range item.RowIds {
			rowMap[id] = true
		}

		var currentRow uint64 = 0

		for scanner.Scan() {
			if rowMap[currentRow+1] {
				if result[item.Key] == nil {
					result[item.Key] = make([]string, 0)
				}
				result[item.Key] = append(result[item.Key], scanner.Text())
			}
			currentRow++
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read object from bucket: %v", err)
		}
	}

	return result, nil
}
