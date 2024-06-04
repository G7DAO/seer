package synchronizer

import (
	"fmt"
	"log"
	"os"
)

var (
	MOONSTREAM_DB_V3_CONTROLLER_API            string = "https://mdb-v3-api.moonstream.to"
	MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN string
)

func CheckVariablesForSynchronizer() error {
	mdbV3ControllerEnvVar := os.Getenv("MOONSTREAM_DB_V3_CONTROLLER_API")

	if mdbV3ControllerEnvVar != "" {
		MOONSTREAM_DB_V3_CONTROLLER_API = mdbV3ControllerEnvVar
		log.Printf("Moonstread DB V3 controller API URL set to: '%s'", MOONSTREAM_DB_V3_CONTROLLER_API)
	}
	MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN = os.Getenv("MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN")
	if MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN == "" {
		return fmt.Errorf("MOONSTREAM_DB_V3_CONTROLLER_SEER_ACCESS_TOKEN environment variable is required")
	}

	return nil
}
