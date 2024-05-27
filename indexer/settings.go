package indexer

import (
	"log"
	"os"
)

var (
	InsertBatchSize  int
	SeerCrawlerLabel string
)

func init() {
	SeerCrawlerLabel = os.Getenv("SEER_CRAWLER_LABEL")

	if SeerCrawlerLabel == "" {
		log.Fatal("SEER_CRAWLER_LABEL environment variable is required")
	}

	InsertBatchSize = 1000

}
