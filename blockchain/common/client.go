package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func NewClient(url string, timeout int) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	rpcClient, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Client{rpcClient: rpcClient}, nil
}

// Client is a wrapper around the Ethereum JSON-RPC client.

type Client struct {
	rpcClient *rpc.Client
}

// GetLatestBlockNumber returns the latest block number.
func (c *Client) GetLatestBlockNumber() (*big.Int, error) {
	var result string
	if err := c.rpcClient.CallContext(context.Background(), &result, "eth_blockNumber"); err != nil {
		return nil, err
	}

	// Convert the hex string to a big.Int
	blockNumber, ok := new(big.Int).SetString(result, 0) // The 0 base lets the function infer the base from the string prefix.
	if !ok {
		return nil, fmt.Errorf("invalid block number format: %s", result)
	}

	return blockNumber, nil
}

// BlockByNumber returns the block with the given number.
func (c *Client) GetBlockByNumber(ctx context.Context, number *big.Int, withTransactions bool) (*BlockJson, error) {

	var rawResponse json.RawMessage // Use RawMessage to capture the entire JSON response
	err := c.rpcClient.CallContext(ctx, &rawResponse, "eth_getBlockByNumber", "0x"+number.Text(16), withTransactions)
	if err != nil {
		fmt.Println("Error calling eth_getBlockByNumber: ", err)
		return nil, err
	}

	var response_json map[string]interface{}

	err = json.Unmarshal(rawResponse, &response_json)
	if err != nil {
		fmt.Println("Error unmarshalling response: ", err)
		return nil, err
	}

	return &block, nil // Return a pointer to the block
}

// BlockByHash returns the block with the given hash.
func (c *Client) BlockByHash(ctx context.Context, hash common.Hash) (BlockJson, error) {
	var block BlockJson
	err := c.rpcClient.CallContext(ctx, &block, "eth_getBlockByHash", hash, true) // true to include transactions
	return block, err
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
func (c *Client) TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	err := c.rpcClient.CallContext(ctx, &receipt, "eth_getTransactionReceipt", hash)
	return receipt, err
}

// Get bytecode of a contract by address.
func (c *Client) GetCode(ctx context.Context, address common.Address, blockNumber uint64) ([]byte, error) {
	var code hexutil.Bytes
	if blockNumber == 0 {
		latestBlockNumber, err := c.GetLatestBlockNumber()
		if err != nil {
			return nil, err
		}
		blockNumber = latestBlockNumber.Uint64()
	}
	err := c.rpcClient.CallContext(ctx, &code, "eth_getCode", address, "0x"+fmt.Sprintf("%x", blockNumber))
	if err != nil {
		log.Printf("Failed to get code for address %s at block %d: %v", address.Hex(), blockNumber, err)
		return nil, err
	}

	if len(code) == 0 {
		return nil, nil
	}
	return code, nil
}
