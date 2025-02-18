package storage

import (
	"fmt"
	"log"
	"os"
	"strconv"
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

var (
	BugoutBroodAPIURL    string
	BugoutSpireAPIURL    string
	SeerBugoutAPITimeout int
	BugoutAddressTag     string
	BugoutAmountTag      string
	BugoutChainTag       string
	BugoutBlockTag       string
	BugoutAccessTokens   string
	BugoutStoreJournalID string
)

func SetBugoutVariablesFromEnv() error {
	// Environment variable keys
	BugoutBroodAPIURL = os.Getenv("BROOD_API_URL")
	if BugoutBroodAPIURL == "" {
		BugoutBroodAPIURL = "https://auth.bugout.dev"
		log.Printf("BUGOUT_BROOD_API_URL environment variable is not set, using default: %s", BugoutBroodAPIURL)
	}
	BugoutSpireAPIURL = os.Getenv("SPIRE_API_URL")
	if BugoutSpireAPIURL == "" {
		BugoutSpireAPIURL = "https://spire.bugout.dev"
		log.Printf("BUGOUT_SPIRE_API_URL environment variable is not set, using default: %s", BugoutSpireAPIURL)
	}
	SeerBugoutAPITimeoutRaw := os.Getenv("SEER_BUGOUT_API_TIMEOUT_SECONDS")
	if SeerBugoutAPITimeoutRaw == "" {
		SeerBugoutAPITimeout = 10
		log.Printf("SEER_BUGOUT_API_TIMEOUT_SECONDS environment variable is not set, using default: %d", SeerBugoutAPITimeout)
	} else {
		timeout, err := strconv.Atoi(SeerBugoutAPITimeoutRaw)
		if err != nil {
			SeerBugoutAPITimeout = 10
			log.Printf("SEER_BUGOUT_API_TIMEOUT_SECONDS environment variable is invalid, using default: %d", SeerBugoutAPITimeout)
		} else {
			SeerBugoutAPITimeout = timeout
		}
	}

	BugoutAccessTokens = os.Getenv("SEER_BUGOUT_ACCESS_TOKENS")
	if BugoutAccessTokens == "" {
		return fmt.Errorf("SEER_BUGOUT_ACCESS_TOKENS environment variable is required")
	}
	BugoutStoreJournalID = os.Getenv("SEER_BUGOUT_STORE_JOURNAL_ID")
	if BugoutStoreJournalID == "" {
		return fmt.Errorf("SEER_BUGOUT_STORE_JOURNAL_ID environment variable is required")
	}

	BugoutAddressTag = os.Getenv("SEER_BUGOUT_ADDRESS_TAG")
	if BugoutAddressTag == "" {
		BugoutAddressTag = "address"
		log.Printf("SEER_BUGOUT_ADDRESS_TAG environment variable is not set, using default: %s", BugoutAddressTag)
	}
	BugoutAmountTag = os.Getenv("SEER_BUGOUT_AMOUNT_TAG")
	if BugoutAmountTag == "" {
		BugoutAmountTag = "amount"
		log.Printf("SEER_BUGOUT_AMOUNT_TAG environment variable is not set, using default: %s", BugoutAmountTag)
	}
	BugoutChainTag = os.Getenv("SEER_BUGOUT_CHAIN_TAG")
	if BugoutChainTag == "" {
		BugoutChainTag = "chain"
		log.Printf("SEER_BUGOUT_CHAIN_TAG environment variable is not set, using default: %s", BugoutChainTag)
	}
	BugoutBlockTag = os.Getenv("SEER_BUGOUT_BLOCK_TAG")
	if BugoutBlockTag == "" {
		BugoutBlockTag = "block"
		log.Printf("SEER_BUGOUT_BLOCK_TAG environment variable is not set, using default: %s", BugoutBlockTag)
	}
	return nil
}

// Blockchains map for storage or database models
var Blockchains = map[string]string{
	"ethereum":                     "ethereum_smartcontract",
	"polygon":                      "polygon_smartcontract",
	"arbitrum_one":                 "arbitrum_one_smartcontract",
	"arbitrum_sepolia":             "arbitrum_sepolia_smartcontract",
	"game7_orbit_arbitrum_sepolia": "game7_orbit_arbitrum_sepolia_smartcontract",
	"xai":                          "xai_smartcontract",
	"xai_sepolia":                  "xai_sepolia_smartcontract",
	"mantle":                       "mantle_smartcontract",
	"mantle_sepolia":               "mantle_sepolia_smartcontract",
}
