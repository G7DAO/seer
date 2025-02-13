package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bugout-dev/bugout-go/pkg/spire"
)

type BugoutAPIClient struct {
	BugoutSpireClient *spire.SpireClient
	JournalID         string
	Token             string
}

// NewBugoutAPIClient initializes and returns a BugoutAPIClient
func NewBugoutAPIClient() (*BugoutAPIClient, error) {
	client := spire.NewClient(BugoutSpireAPIURL, time.Duration(SeerBugoutAPITimeout)*time.Second)
	return &BugoutAPIClient{BugoutSpireClient: &client, JournalID: BugoutStoreJournalID, Token: BugoutAccessTokens}, nil
}

// findOrCreateAddressEntry searches for an existing entry for the given address
// If none is found, it creates a new one and returns its entry ID.
// If none is found, it creates a new one and returns its entry ID.
func (c *BugoutAPIClient) FindOrCreateAddressEntry(
	address, chain string,
) (string, error) {
	query := fmt.Sprintf("#%s:%s #%s:%s", BugoutAddressTag, address, BugoutChainTag, chain)
	params := map[string]string{"content": "false"}

	entries, err := c.BugoutSpireClient.SearchEntries(c.Token, c.JournalID, query, 1, 0, params)
	if err != nil {
		return "", fmt.Errorf("search error: %w", err)
	}
	if entries.TotalResults > 0 && len(entries.Results) > 0 {
		parts := strings.Split(entries.Results[0].Url, "/")
		if len(parts) == 0 {
			return "", errors.New("invalid entry URL in search result")
		}
		return parts[len(parts)-1], nil
	}

	bugoutCtx := spire.EntryContext{
		ContextType: "blockchain",
		ContextID:   chain,
	}

	// No address entry found -> Create a new one
	title := fmt.Sprintf("[%s] Address: %s", chain, address)
	tags := []string{fmt.Sprintf("%s:%s", BugoutAddressTag, address), fmt.Sprintf("%s:%s", BugoutChainTag, chain)}

	entry, err := c.BugoutSpireClient.CreateEntry(c.Token, c.JournalID, title, "", tags, bugoutCtx)
	if err != nil {
		return "", fmt.Errorf("failed to create new address entry: %w", err)
	}
	return entry.Id, nil
}

func (c *BugoutAPIClient) GetAllEntriesAndAddresses() (map[string]string, error) {
	EntryIDToAddress := make(map[string]string)

	query := fmt.Sprintf("#%s #%s", BugoutAddressTag, BugoutChainTag)
	params := map[string]string{"content": "false"}

	entries, err := c.BugoutSpireClient.SearchEntries(c.Token, c.JournalID, query, 1, 0, params)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}

	for _, entry := range entries.Results {
		parts := strings.Split(entry.Url, "/")
		if len(parts) == 0 {
			return nil, errors.New("invalid entry URL in search result")
		}
		EntryIDToAddress[parts[len(parts)-1]] = entry.Title
	}

	return EntryIDToAddress, nil
}

func (c *BugoutAPIClient) UpdateAddressesEntryContent(blockchain string, addresses []string) (string, error) {

	/// find entry with all addresses
	query := fmt.Sprintf("#%s:%s #addressesStorage #test", BugoutChainTag, blockchain)
	params := map[string]string{"content": "true"}

	entries, err := c.BugoutSpireClient.SearchEntries(c.Token, c.JournalID, query, 1, 0, params)
	if err != nil {
		fmt.Printf("search error: token: %s, journalID: %s", c.Token, c.JournalID)
		return "", fmt.Errorf("search error: %w", err)
	}

	log.Printf("Found %d entries", len(entries.Results))

	bugoutCtx := spire.EntryContext{
		ContextType: "blockchain",
		ContextID:   blockchain,
	}

	if len(entries.Results) == 0 {
		// create new entry
		title := fmt.Sprintf("[%s] Addresses Storage", blockchain)
		tags := []string{fmt.Sprintf("%s:%s", BugoutChainTag, blockchain), "addressesStorage"}
		content, err := json.Marshal(addresses)
		if err != nil {
			return "", fmt.Errorf("failed to marshal addresses: %w", err)
		}
		entry, err := c.BugoutSpireClient.CreateEntry(c.Token, c.JournalID, title, string(content), tags, bugoutCtx)
		if err != nil {
			return "", fmt.Errorf("failed to create new entry: %w", err)
		}
		log.Printf("Created new entry: %s", entry.Id)
		return entry.Id, nil
	}

	// update existing entry
	entryID := entries.Results[0].Id
	log.Printf("Updating existing entry: %s", entryID)
	content, err := json.Marshal(addresses)
	if err != nil {
		return entryID, fmt.Errorf("failed to marshal addresses: %w", err)
	}
	_, err = c.BugoutSpireClient.UpdateEntry(c.Token, c.JournalID, entryID, entries.Results[0].Title, string(content))
	if err != nil {
		return entryID, fmt.Errorf("failed to update entry: %w", err)
	}
	log.Printf("Updated entry: %s", entryID)
	return entryID, nil
}
