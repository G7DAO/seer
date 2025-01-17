package heartbeat

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bugout-dev/bugout-go/pkg/spire"
)

type Status struct {
	Status                    string    `json:"status"`
	StartBlock                uint64    `json:"start_block"`
	LastBlock                 uint64    `json:"last_block"`
	CrawlStartTime           string    `json:"crawl_start_time"`
	CurrentTime              string    `json:"current_time"`
	CurrentEventJobsLength    int       `json:"current_event_jobs_length"`
	CurrentFunctionJobsLength int       `json:"current_function_call_jobs_length"`
	JobsLastRefetchedAt      string    `json:"jobs_last_refetched_at"`
}

type Manager struct {
	client            spire.SpireClient
	journalID         string
	serviceType       string
	blockchainType    string
	lastHeartbeat     time.Time
	heartbeatEntryID  string
	startTime         time.Time
	heartbeatInterval time.Duration
}

func NewManager(token, journalID, serviceType, blockchainType string) (*Manager, error) {
	client := spire.NewClient(spire.BugoutSpireURL, 10*time.Second)

	// Test connection
	_, err := client.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Spire: %w", err)
	}

	// Search for existing heartbeat entry
	searchQuery := fmt.Sprintf("#%s #heartbeat #%s !#dead", serviceType, blockchainType)
	results, err := client.SearchEntries(token, journalID, searchQuery, 1, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search for heartbeat entry: %w", err)
	}

	var entryID string
	if len(results.Results) > 0 {
		entryID = results.Results[0].Id
		log.Printf("Found existing heartbeat entry for %s - %s", serviceType, blockchainType)
	} else {
		log.Printf("No %s heartbeat entry found for %s, creating one", serviceType, blockchainType)
		entry, err := client.CreateEntry(
			token,
			journalID,
			fmt.Sprintf("%s Heartbeat - %s", serviceType, blockchainType),
			"",
			[]string{serviceType, "heartbeat", blockchainType},
			spire.EntryContext{},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create heartbeat entry: %w", err)
		}
		entryID = entry.Id
	}

	return &Manager{
		client:            client,
		journalID:         journalID,
		serviceType:       serviceType,
		blockchainType:    blockchainType,
		lastHeartbeat:     time.Now(),
		heartbeatEntryID:  entryID,
		startTime:         time.Now(),
		heartbeatInterval: 5 * time.Minute,
	}, nil
}

func (h *Manager) Send(token string, status string, startBlock, lastBlock uint64, eventJobs, functionJobs int) error {
	if time.Since(h.lastHeartbeat) < h.heartbeatInterval {
		return nil
	}

	currentTime := time.Now().UTC()
	heartbeatStatus := Status{
		Status:                    status,
		StartBlock:                startBlock,
		LastBlock:                 lastBlock,
		CrawlStartTime:            h.startTime.UTC().Format(time.RFC3339),
		CurrentTime:               currentTime.Format(time.RFC3339),
		CurrentEventJobsLength:    eventJobs,
		CurrentFunctionJobsLength: functionJobs,
		JobsLastRefetchedAt:       currentTime.Format(time.RFC3339),
	}

	content, err := json.MarshalIndent(heartbeatStatus, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal heartbeat status: %w", err)
	}

	title := fmt.Sprintf("%s Heartbeat - %s. Status: %s - %s",
		h.serviceType,
		h.blockchainType,
		heartbeatStatus.Status,
		heartbeatStatus.CurrentTime,
	)

	_, err = h.client.UpdateEntry(
		token,
		h.journalID,
		h.heartbeatEntryID,
		title,
		string(content),
	)

	if err != nil {
		return fmt.Errorf("failed to update heartbeat entry: %w", err)
	}

	h.lastHeartbeat = currentTime
	return nil
}
