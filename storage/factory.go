package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func NewStorage(storageType ...string) (Storer, error) {
	var stype string
	if len(storageType) > 0 && storageType[0] != "" {
		stype = storageType[0]
	} else {
		stype = SeerCrawlerStorageType
		if stype == "" {
			stype = "filesystem" // Default to filesystem
		}
	}

	baseDir := SeerCrawlerStoragePath

	if baseDir == "" {
		baseDir = "data"
	}

	switch stype {
	case "filesystem":
		return NewFileStorage(baseDir), nil
	case "gcs":
		// Google Cloud Storage
		// Initialize the GCS client once
		serviceAccountKeyPath := GCSServiceFilePath
		if serviceAccountKeyPath == "" {
			return nil, fmt.Errorf("missing GOOGLE_APPLICATION_CREDENTIALS environment variable")
		}

		ctx := context.Background()
		client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountKeyPath))
		if err != nil {
			return nil, fmt.Errorf("failed to create GCS client: %v", err)
		}

		return NewGCSStorage(client, baseDir), nil
	case "s3":
		// Amazon S3
		return NewS3Storage(), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", stype)
	}
}
