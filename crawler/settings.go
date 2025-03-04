package crawler

import (
	"fmt"
	"os"
	"strconv"
)

var (
	SeerDefaultBlockShift    int64  = 100
	SeerCrawlerStoragePrefix string = "dev"

	SEER_CRAWLER_DEBUG = false
)

// CheckVariablesForCrawler checks required environment variables but only for the specified chain.
func CheckVariablesForCrawler() error {
	// 1. Check the storage prefix environment
	SeerCrawlerStoragePrefixEnvVar := os.Getenv("SEER_CRAWLER_STORAGE_PREFIX")
	switch SeerCrawlerStoragePrefixEnvVar {
	case "dev":
		SeerCrawlerStoragePrefix = "dev"
	case "prod":
		SeerCrawlerStoragePrefix = "prod"
	default:
		return fmt.Errorf("unknown storage prefix set: %s", SeerCrawlerStoragePrefixEnvVar)
	}

	// 2. Set the debug mode if requested
	SEER_CRAWLER_DEBUG_RAW := os.Getenv("SEER_CRAWLER_DEBUG")
	SEER_CRAWLER_DEBUG, _ = strconv.ParseBool(SEER_CRAWLER_DEBUG_RAW)

	return nil
}

// checkEnv is a tiny helper verifying that an environment variable is non-empty
func checkEnv(envVar string) error {
	if os.Getenv(envVar) == "" {
		return fmt.Errorf("%s environment variable is required", envVar)
	}
	return nil
}
