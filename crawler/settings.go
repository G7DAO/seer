package crawler

import (
	"fmt"
	"os"
)

var (
	SeerCrawlerStoragePrefix string = "dev"

	BlockchainURLs map[string]string
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

	MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI environment variable is required")
	}
	MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI environment variable is required")
	}

	BlockchainURLs = map[string]string{
		"ethereum": MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI,
		"polygon":  MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI,
	}

	return nil
}
