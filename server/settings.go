package server

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	bugout "github.com/bugout-dev/bugout-go/pkg"
	"github.com/google/uuid"
)

var (
	BUGOUT_BROOD_URL     = "https://auth.bugout.dev"
	BUGOUT_BROOD_URL_RAW = os.Getenv("BUGOUT_BROOD_URL")

	BUGOUT_AUTH_CALL_TIMEOUT         = time.Second * 5
	BUGOUT_RESOURCE_TYPE_SEER_ACCESS = "seer-access"

	MOONSTREAM_APPLICATION_ID = os.Getenv("MOONSTREAM_APPLICATION_ID")

	SEER_BUGOUT_TIMEOUT_SECONDS     = 10
	SEER_BUGOUT_TIMEOUT_SECONDS_RAW = os.Getenv("SEER_BUGOUT_TIMEOUT_SECONDS")
)

func CheckVariablesForServer() error {
	if BUGOUT_BROOD_URL_RAW != "" {
		BUGOUT_BROOD_URL = BUGOUT_BROOD_URL_RAW
		log.Printf("Bugout URL modified and set to %s", BUGOUT_BROOD_URL)
	}

	uuidErr := uuid.Validate(MOONSTREAM_APPLICATION_ID)
	if uuidErr != nil {
		return fmt.Errorf("MOONSTREAM_APPLICATION_ID environment variable is required in UUID format")
	}

	return nil
}

func InitBugoutClient() (*bugout.BugoutClient, error) {
	bugoutTimeoutSeconds, atoiErr := strconv.Atoi(SEER_BUGOUT_TIMEOUT_SECONDS_RAW)
	if atoiErr == nil {
		log.Printf("Bugout timeout modified and set to %d", SEER_BUGOUT_TIMEOUT_SECONDS)
		SEER_BUGOUT_TIMEOUT_SECONDS = bugoutTimeoutSeconds
	}

	SEER_BUGOUT_TIMEOUT := time.Duration(SEER_BUGOUT_TIMEOUT_SECONDS) * time.Second

	bugoutClient := bugout.ClientBrood(BUGOUT_BROOD_URL, SEER_BUGOUT_TIMEOUT)
	return &bugoutClient, nil
}
