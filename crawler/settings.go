package crawler

import (
	"fmt"
	"os"
	"strconv"
)

var (
	SeerCrawlerStoragePrefix string = "dev"

	BlockchainURLs map[string]string

	SEER_CRAWLER_DEBUG = false
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
	MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI environment variable is required")
	}
	MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI environment variable is required")
	}
	MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI environment variable is required")
	}
	MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI environment variable is required")
	}
	MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI environment variable is required")
	}
	MOONSTREAM_NODE_XAI_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_XAI_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_XAI_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_XAI_A_EXTERNAL_URI environment variable is required")
	}
	MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI environment variable is required")
	}
	MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI environment variable is required")
	}
	MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI := os.Getenv("MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI")
	if MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI == "" {
		return fmt.Errorf("MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI environment variable is required")
	}

	SEER_CRAWLER_DEBUG_RAW := os.Getenv("SEER_CRAWLER_DEBUG")
	SEER_CRAWLER_DEBUG, _ = strconv.ParseBool(SEER_CRAWLER_DEBUG_RAW)

	BlockchainURLs = map[string]string{
		"ethereum":                     MOONSTREAM_NODE_ETHEREUM_A_EXTERNAL_URI,
		"sepolia":                      MOONSTREAM_NODE_SEPOLIA_A_EXTERNAL_URI,
		"polygon":                      MOONSTREAM_NODE_POLYGON_A_EXTERNAL_URI,
		"arbitrum_one":                 MOONSTREAM_NODE_ARBITRUM_ONE_A_EXTERNAL_URI,
		"arbitrum_sepolia":             MOONSTREAM_NODE_ARBITRUM_SEPOLIA_A_EXTERNAL_URI,
		"game7_orbit_arbitrum_sepolia": MOONSTREAM_NODE_GAME7_ORBIT_ARBITRUM_SEPOLIA_A_EXTERNAL_URI,
		"xai":                          MOONSTREAM_NODE_XAI_A_EXTERNAL_URI,
		"xai_sepolia":                  MOONSTREAM_NODE_XAI_SEPOLIA_A_EXTERNAL_URI,
		"mantle":                       MOONSTREAM_NODE_MANTLE_A_EXTERNAL_URI,
		"mantle_sepolia":               MOONSTREAM_NODE_MANTLE_SEPOLIA_A_EXTERNAL_URI,
	}

	return nil
}
