package storage

import (
	"context"
	"fmt"

	gcp_storage "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func NewStorage(storageType ...string) (Storer, error) {
	var stype string
	if len(storageType) > 0 && storageType[0] != "" {
		stype = storageType[0]
	} else {
		stype = SeerCrawlerStorageType
	}

	baseDir := SeerCrawlerStoragePath

	if baseDir == "" {
		baseDir = "data"
	}

	switch stype {
	case "filesystem":
		return NewFileStorage(baseDir), nil
	case "gcp-storage":
		// Google Cloud Storage
		ctx := context.Background()

		var client *gcp_storage.Client
		var clientErr error
		if GCPStorageServiceAccountCredsPath != "" {
			client, clientErr = gcp_storage.NewClient(ctx, option.WithCredentialsFile(GCPStorageServiceAccountCredsPath))
		} else {
			client, clientErr = gcp_storage.NewClient(ctx)
		}
		if clientErr != nil {
			return nil, fmt.Errorf("failed to create GCS client: %v", clientErr)
		}

		return NewGCSStorage(client, baseDir), nil
	case "aws-bucket":
		// Amazon S3
		return NewS3Storage(), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", stype)
	}
}
