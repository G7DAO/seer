package synchronizer

import (
	"log"
	"os"
)

var (
	MOONSTREAM_DB_MANAGER_API          string
	MOONSTREAM_DB_MANAGER_ACCESS_TOKEN string
)

func init() {
	MOONSTREAM_DB_MANAGER_API = os.Getenv("MOONSTREAM_DB_MANAGER_API")
	if MOONSTREAM_DB_MANAGER_API == "" {
		log.Fatal("MOONSTREAM_DB_MANAGER_API environment variable is required")
	}
	MOONSTREAM_DB_MANAGER_ACCESS_TOKEN = os.Getenv("MOONSTREAM_DB_MANAGER_ACCESS_TOKEN")
	if MOONSTREAM_DB_MANAGER_ACCESS_TOKEN == "" {
		log.Fatal("MOONSTREAM_DB_MANAGER_ACCESS_TOKEN environment variable is required")
	}

}
