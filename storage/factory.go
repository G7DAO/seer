package storage

import (
	"context"
	"fmt"
	"log"

	gcp_storage "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// NewStorage initialize storage placement for protobuf batch data.
func NewStorage(storageType, basePath string) (Storer, error) {
	switch storageType {
	case "filesystem":
		log.Println("Using filesystem storage")
		return NewFileStorage(basePath), nil
	case "gcp-storage":
		// Google Cloud Storage
		ctx := context.Background()
		log.Println("Creating GCS client")

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

		return NewGCSStorage(client, basePath), nil
	case "aws-bucket":
		// Amazon S3 Bucket
		// TODO: Add client initialization
		log.Println("AWS bucket support not implemented yet")
		return NewS3Storage(basePath), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
