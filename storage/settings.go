package storage

import (
	"fmt"
	"log"
	"os"
)

var (
	SeerCrawlerStorageType            string
	SeerCrawlerStorageBucket          string
	GCPStorageServiceAccountCredsPath string
	SeerCrawlerStoragePath            string = "data"
)

func SetStorageBucketFromEnv() error {
	SeerCrawlerStorageBucket = os.Getenv("SEER_CRAWLER_STORAGE_BUCKET")
	if SeerCrawlerStorageBucket == "" {
		return fmt.Errorf("SEER_CRAWLER_STORAGE_BUCKET environment variable is required")
	}
	return nil
}

func CheckVariablesForStorage() error {
	SeerCrawlerStorageTypeEnvVar := os.Getenv("SEER_CRAWLER_STORAGE_TYPE")
	switch SeerCrawlerStorageTypeEnvVar {
	case "filesystem":
		SeerCrawlerStorageType = "filesystem"
	case "gcp-storage":
		SeerCrawlerStorageType = "gcp-storage"

		bucketError := SetStorageBucketFromEnv()
		if bucketError != nil {
			return bucketError
		}

		GCPStorageServiceAccountCredsPath = os.Getenv("MOONSTREAM_STORAGE_GCP_SERVICE_ACCOUNT_CREDS_PATH")
		if GCPStorageServiceAccountCredsPath != "" {
			log.Printf("MOONSTREAM_STORAGE_GCP_SERVICE_ACCOUNT_CREDS_PATH environment variable is set, using it for authentification")
		}
	case "aws-bucket":
		SeerCrawlerStorageType = "aws-bucket"
		bucketError := SetStorageBucketFromEnv()
		if bucketError != nil {
			return bucketError
		}
	default:
		SeerCrawlerStorageType = "filesystem"
		log.Printf("SEER_CRAWLER_STORAGE_TYPE environment variable is not set or unknown, using default: %s", SeerCrawlerStorageType)
	}

	SeerCrawlerStoragePathEnvVar := os.Getenv("SEER_CRAWLER_STORAGE_PATH")
	if SeerCrawlerStoragePathEnvVar != "" {
		SeerCrawlerStoragePath = SeerCrawlerStoragePathEnvVar
		log.Printf("Default seer crawler storage path set to '%s'", SeerCrawlerStoragePath)
	}

	return nil
}

// Blockchains map for storage or database models
var Blockchains = map[string]string{
	"ethereum": "ethereum_smartcontract",
	"polygon":  "polygon_smartcontract",
}
