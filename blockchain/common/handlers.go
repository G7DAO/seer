package common

import (
	"errors"
	"math/big"

	"github.com/moonstream-to/seer/blockchain/ethereum"
	"github.com/moonstream-to/seer/blockchain/polygon"

	// "github.com/moonstream-to/seer/blockchain/polygon"
	// "github.com/moonstream-to/seer/blockchain/starknet"
	"github.com/moonstream-to/seer/indexer"
	"google.golang.org/protobuf/proto"
)

func wrapClient(url, chain string) (BlockchainClient, error) {
	if chain == "ethereum" {
		client, err := ethereum.NewClient(url)
		return client, err
	} else if chain == "polygon" {
		client, err := polygon.NewClient(url)
		return client, err
	} else {
		return nil, errors.New("unsupported chain type")
	} // Ensure ethereum.Client implements BlockchainClient.
}

type BlockData struct {
	BlockNumber    uint64
	BlockHash      string
	BlockTimestamp int64
	ChainID        string
	Data           map[string]interface{}
}

type BlockchainClient interface {
	// GetBlockByNumber(ctx context.Context, number *big.Int) (*proto.Message, error)
	GetLatestBlockNumber() (*big.Int, error)
	FetchAsProtoEvents(*big.Int, *big.Int, map[uint64]indexer.BlockCahche) ([]proto.Message, []indexer.LogIndex, error)
	FetchAsProtoBlocks(*big.Int, *big.Int) ([]proto.Message, []proto.Message, []indexer.BlockIndex, []indexer.TransactionIndex, map[uint64]indexer.BlockCahche, error)
	ChainType() string
}

// return client depends on the blockchain type
func NewClient(chainType string, url string) (BlockchainClient, error) {
	client, err := wrapClient(url, chainType)
	if err != nil {
		return nil, errors.New("unsupported chain type")
	}

	return client, nil
}

// crawl blocks
func CrawlBlocks(client BlockchainClient, startBlock *big.Int, endBlock *big.Int) ([]proto.Message, []proto.Message, []indexer.BlockIndex, []indexer.TransactionIndex, map[uint64]indexer.BlockCahche, error) {
	blocks, transactions, blocksIndex, transactionsIndex, blocksCache, err := client.FetchAsProtoBlocks(startBlock, endBlock)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// index blocks and transactions

	return blocks, transactions, blocksIndex, transactionsIndex, blocksCache, nil

}

func CrawlEvents(client BlockchainClient, startBlock *big.Int, endBlock *big.Int, BlockCahche map[uint64]indexer.BlockCahche) ([]proto.Message, []indexer.LogIndex, error) {

	events, eventsIndex, err := client.FetchAsProtoEvents(startBlock, endBlock, BlockCahche)
	if err != nil {
		return nil, nil, err
	}

	return events, eventsIndex, nil
}
