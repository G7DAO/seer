package storage

import (
	"os"
)

var (
	SeerCrawlerStorageType   string
	SeerCrawlerStorageBucket string
	GCSServiceFilePath       string
	SeerCrawlerStoragePath   string
)

func init() {

	SeerCrawlerStorageType = os.Getenv("SEER_CRAWLER_STORAGE_TYPE")
	if SeerCrawlerStorageType == "" {
		SeerCrawlerStorageType = "filesystem"
	}

	SeerCrawlerStorageBucket = os.Getenv("SEER_CRAWLER_STORAGE_BUCKET")

	GCSServiceFilePath = os.Getenv("GCS_SERVICE_FILE_PATH")

	SeerCrawlerStoragePath = os.Getenv("SEER_CRAWLER_STORAGE_PATH")
	if SeerCrawlerStoragePath == "" {
		SeerCrawlerStoragePath = "data"
	}
}

// define a map
var Blokchain = map[string]string{
	"ethereum": "ethereum_smartcontract",
	"polygon":  "polygon_smartcontract",
}
