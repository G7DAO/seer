package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/protobuf/proto"
)

func NewClient(url string) (*Client, error) {
	rpcClient, err := rpc.DialContext(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &Client{rpcClient: rpcClient}, nil
}

// Client is a wrapper around the Ethereum JSON-RPC client.

type Client struct {
	rpcClient *rpc.Client
}

// Client common

// ChainType returns the chain type.
func (c *Client) ChainType() string {
	return "ethereum"
}

// Close closes the underlying RPC client.
func (c *Client) Close() {
	c.rpcClient.Close()
}

// GetLatestBlockNumber returns the latest block number.
func (c *Client) GetLatestBlockNumber(ctx context.Context) (*big.Int, error) {
	var blockNumber *big.Int
	err := c.rpcClient.CallContext(ctx, &blockNumber, "eth_blockNumber")
	return blockNumber, err
}

// BlockByNumber returns the block with the given number.
func (c *Client) GetBlockByNumber(ctx context.Context, number *big.Int) (*BlockJson, error) {
	var block *BlockJson
	err := c.rpcClient.CallContext(ctx, &block, "eth_getBlockByNumber", number, true) // true to include transactions
	return block, err
}

// BlockByHash returns the block with the given hash.

func (c *Client) BlockByHash(ctx context.Context, hash common.Hash) (*BlockJson, error) {
	var block *BlockJson
	err := c.rpcClient.CallContext(ctx, &block, "eth_getBlockByHash", hash, true) // true to include transactions
	return block, err
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.

func (c *Client) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	err := c.rpcClient.CallContext(ctx, &receipt, "eth_getTransactionReceipt", hash)
	return receipt, err
}

// FilterLogs returns the logs that satisfy the specified filter query.

func (c *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	var logs []types.Log
	err := c.rpcClient.CallContext(ctx, &logs, "eth_getLogs", q)
	return logs, err
}

// fetchBlocks returns the blocks for a given range.

func (c *Client) fetchBlocks(ctx context.Context, from, to *big.Int) ([]*BlockJson, error) {
	var blocks []*BlockJson

	for i := from; i.Cmp(to) <= 0; i.Add(i, big.NewInt(1)) {
		block, err := c.GetBlockByNumber(ctx, i)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

// Utility function to convert big.Int to its hexadecimal representation.
func toHex(number *big.Int) string {
	return fmt.Sprintf("0x%x", number)
}

// FetchBlocksInRange fetches blocks within a specified range.
// This could be useful for batch processing or analysis.
func (c *Client) FetchBlocksInRange(from, to *big.Int) ([]*BlockJson, error) {
	var blocks []*BlockJson
	ctx := context.Background() // For simplicity, using a background context; consider timeouts for production.

	for i := new(big.Int).Set(from); i.Cmp(to) <= 0; i.Add(i, big.NewInt(1)) {
		block, err := c.GetBlockByNumber(ctx, i)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// ParseBlocksAndTransactions parses blocks and their transactions into custom data structures.
// This method showcases how to handle and transform detailed block and transaction data.
func (c *Client) ParseBlocksAndTransactions(from, to *big.Int) ([]*Block, []*SingleTransaction, error) {
	blocksJson, err := c.FetchBlocksInRange(from, to)
	if err != nil {
		return nil, nil, err
	}

	var parsedBlocks []*Block
	var parsedTransactions []*SingleTransaction
	for _, blockJson := range blocksJson {
		// Convert BlockJson to Block and Transactions as required.
		parsedBlock := ToProtoSingleBlock(blockJson)

		// Example: Parsing transactions within the block
		for _, txJson := range blockJson.Transactions {

			parsedTransaction := ToProtoSingleTransaction(txJson)
			parsedTransactions = append(parsedTransactions, parsedTransaction)
		}

		parsedBlocks = append(parsedBlocks, parsedBlock)
	}

	return parsedBlocks, parsedTransactions, nil
}

// EncodeToProtoBlocks demonstrates a hypothetical method to encode blocks to a protobuf format.
// You would replace this with actual logic based on your protobuf definitions.
func (c *Client) EncodeToProtoBlocks(from, to *big.Int) ([]proto.Message, []proto.Message, error) {
	parsedBlocks, parsedTransactions, err := c.ParseBlocksAndTransactions(from, to)

	if err != nil {
		return nil, nil, err
	}

	var blocksProto []proto.Message
	for _, block := range parsedBlocks {
		blocksProto = append(blocksProto, block) // Assuming block is already a proto.Message
	}

	var transactionsProto []proto.Message
	for _, transaction := range parsedTransactions {
		transactionsProto = append(transactionsProto, transaction) // Assuming transaction is also a proto.Message
	}

	return blocksProto, transactionsProto, nil
}

func ToProtoSingleBlock(obj *BlockJson) *Block {
	return &Block{
		BlockNumber:      obj.BlockNumber,
		Difficulty:       obj.Difficulty,
		ExtraData:        obj.ExtraData,
		GasLimit:         obj.GasLimit,
		GasUsed:          obj.GasUsed,
		BaseFeePerGas:    obj.BaseFeePerGas,
		Hash:             obj.Hash,
		LogsBloom:        obj.LogsBloom,
		Miner:            obj.Miner,
		Nonce:            obj.Nonce,
		ParentHash:       obj.ParentHash,
		ReceiptRoot:      obj.ReceiptRoot,
		Uncles:           obj.Uncles,
		Size:             obj.Size,
		StateRoot:        obj.StateRoot,
		Timestamp:        obj.Timestamp,
		TotalDifficulty:  obj.TotalDifficulty,
		TransactionsRoot: obj.TransactionsRoot,
		IndexedAt:        obj.IndexedAt,
	}
}

func ToProtoSingleTransaction(obj *SingleTransactionJson) *SingleTransaction {
	return &SingleTransaction{
		Hash:                 obj.Hash,
		BlockNumber:          obj.BlockNumber,
		FromAddress:          obj.FromAddress,
		ToAddress:            obj.ToAddress,
		Gas:                  obj.Gas,
		GasPrice:             obj.GasPrice,
		MaxFeePerGas:         obj.MaxFeePerGas,
		MaxPriorityFeePerGas: obj.MaxPriorityFeePerGas,
		Input:                obj.Input,
		Nonce:                obj.Nonce,
		TransactionIndex:     obj.TransactionIndex,
		TransactionType:      obj.TransactionType,
		Value:                obj.Value,
		IndexedAt:            obj.IndexedAt,
		BlockTimestamp:       obj.BlockTimestamp,
	}
}

func ToProtoSingleEventLog(obj *SingleEventJson) *EventLog {
	return &EventLog{
		Address:         obj.Address,
		Topics:          obj.Topics,
		Data:            obj.Data,
		BlockNumber:     obj.BlockNumber,
		TransactionHash: obj.TransactionHash,
		BlockHash:       obj.BlockHash,
		Removed:         obj.Removed,
	}
}
