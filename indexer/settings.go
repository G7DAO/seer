package indexer

import (
	"fmt"
	"os"
)

var (
	InsertBatchSize  = 1000
	SeerCrawlerLabel string
)

func CheckVariablesForIndexer() error {
	SeerCrawlerLabel = os.Getenv("SEER_CRAWLER_INDEXER_LABEL")
	if SeerCrawlerLabel == "" {
		return fmt.Errorf("SEER_CRAWLER_INDEXER_LABEL environment variable is required")
	}

	return nil
}
