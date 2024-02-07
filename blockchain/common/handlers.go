package common

import (
	"errors"
	"math/big"
	"path/filepath"

	"github.com/moonstream-to/seer/blockchain/ethereum"
	"github.com/moonstream-to/seer/blockchain/polygon"
	"github.com/moonstream-to/seer/blockchain/starknet"
	"google.golang.org/protobuf/proto"
)

func wrapEthereumClient(url string) (BlockchainClient, error) {
	client, err := ethereum.NewClient(url)
	return client, err // Ensure ethereum.Client implements BlockchainClient.
}

var (
	ErrUnsupportedChainType = errors.New("unsupported chain type")
	blockchainClients       = map[string]func(string) (BlockchainClient, error){
		"ethereum": ethereum.NewClient,
		"polygon":  polygon.NewClient,
		"starknet": starknet.NewClient,
	}
	Blocks = make(map[string][]BlockData)
)

type BlockIndex struct {
	BlockNumber    uint64
	BlockHash      string
	BlockTimestamp int64
	ParentHash     string
	filepath       string
}

type TransactionIndex struct {
	BlockNumber          uint64
	BlockHash            string
	BlockTimestamp       int64
	TransactionHash      string
	TransactionIndex     uint64
	TransactionTimestamp int64
	filepath             string
}

type LogIndex struct {
	BlockNumber     uint64
	BlockHash       string
	BlockTimestamp  int64
	TransactionHash string
	LogIndex        uint64
	filepath        string
}

type BlockData struct {
	BlockNumber    uint64
	BlockHash      string
	BlockTimestamp int64
	ChainID        string
	Data           map[string]interface{}
}

type BlockchainClient interface {
	GetBlockByNumber(blockNumber uint64) ([]BlockData, error)
	GetLatestBlockNumber() (*big.Int, error)
	ParseBlocksAndTransactions(*big.Int, *big.Int) (proto.Message, proto.Message, error)
	ChainType() string
}

func indexBlocks(client BlockchainClient, blocks []proto.Message, startBlock, endBlock *big.Int) []BlockIndex {
	var blocksIndex []BlockIndex
	for _, block := range blocks {
		// block is a protobuf message of client.Block
		blocksIndex = append(blocksIndex, BlockIndex{
			BlockNumber:    block.BlockNumber,
			BlockHash:      block.BlockHash,
			BlockTimestamp: block.BlockTimestamp,
			ParentHash:     block.ParentHash,
			filepath:       filepath.Join(startBlock.String()+"-"+endBlock.String(), "blocks.proto"),
		})

	}
	return blocksIndex
}

func indexTransactions(client BlockchainClient, transactions []proto.Message, startBlock, endBlock *big.Int) []TransactionIndex {
	var transactionsIndex []TransactionIndex

	for _, transaction := range transactions {
		transactionsIndex = append(transactionsIndex, TransactionIndex{
			BlockNumber:          transaction.BlockNumber,
			BlockHash:            transaction.BlockHash,
			BlockTimestamp:       transaction.BlockTimestamp,
			TransactionHash:      transaction.TransactionHash,
			TransactionIndex:     transaction.TransactionIndex,
			TransactionTimestamp: transaction.TransactionTimestamp,
			filepath:             filepath.Join(startBlock.String()+"-"+endBlock.String(), "transactions.proto"),
		})
	}
	return transactionsIndex
}

// return client depends on the blockchain type
func NewClient(chainType string, url string) (BlockchainClient, error) {
	clientFunc, ok := blockchainClients[chainType]
	if !ok {
		return nil, ErrUnsupportedChainType
	}

	return clientFunc(url)
}

// crawl blocks
func CrawlBlocks(client BlockchainClient, startBlock *big.Int, endBlock *big.Int) ([]*proto.Message, []*proto.Message, []BlockIndex, []TransactionIndex, error) {
	blocks, transactions, err := client.ParseBlocksAndTransactions(startBlock, endBlock)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// index blocks and transactions

	blocksIndex := indexBlocks(client, blocks)

	transactionsIndex := indexTransactions(client, transactions)

	return blocks, transactions, blocksIndex, transactionsIndex, nil

}
