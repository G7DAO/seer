package crawler

import (
	"fmt"
	"os"
)

var (
	SeerCrawlerStoragePrefix string = "dev"
)

func CheckVariablesForCrawler() error {
	SeerCrawlerStoragePrefixEnvVar := os.Getenv("SEER_CRAWLER_STORAGE_PREFIX")
	switch SeerCrawlerStoragePrefixEnvVar {
	case "dev":
		SeerCrawlerStoragePrefix = "dev"
	case "prod":
		SeerCrawlerStoragePrefix = "prod"
	default:
		return fmt.Errorf("unknown storage prefix set: %s", SeerCrawlerStoragePrefixEnvVar)
	}

	return nil
}
