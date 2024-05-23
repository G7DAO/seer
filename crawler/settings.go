package crawler

import "os"

var (
	InfuraURLs               map[string]string
	infuraKey                string
	SeerCrawlerStoragePrefix string
)

func init() {

	SeerCrawlerStoragePrefix = os.Getenv("SEER_CRAWLER_STORAGE_PREFIX")
	if SeerCrawlerStoragePrefix == "" {
		SeerCrawlerStoragePrefix = "dev"
	}

	infuraKey = os.Getenv("INFURA_KEY")

	InfuraURLs = map[string]string{
		"ethereum": "https://mainnet.infura.io/v3/" + infuraKey,
		"polygon":  "https://polygon-mainnet.infura.io/v3/" + infuraKey,
	}
}
