package starknet

import (
	"bytes"
	"context"
	"math/big"
	"net/http"
	"net/http/cookiejar"

	// "github.com/NethermindEth/starknet.go/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/publicsuffix"
	"google.golang.org/protobuf/proto"

	seer_common "github.com/moonstream-to/seer/blockchain/common"
	"github.com/moonstream-to/seer/indexer"
)

// Based on starknet rpc Provider
func NewClient(url string, timeout int) (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}
	rpcClient, clientErr := ethrpc.DialHTTPWithClient(url, client)
	if clientErr != nil {
		return nil, clientErr
	}

	return &Client{rpcClient: rpcClient}, nil
}

// Client is a wrapper around the Starknet JSON-RPC client.

type Client struct {
	rpcClient *ethrpc.Client
}

// Client common

// ChainType returns the chain type.
func (c *Client) ChainType() string {
	return "starknet"
}

func (c *Client) DecodeProtoEntireBlockToJson(rawData *bytes.Buffer) (*seer_common.BlocksBatchJson, error) {
	return nil, nil
}

func (c *Client) DecodeProtoEntireBlockToLabels(rawData *bytes.Buffer, blocksCache map[uint64]uint64, abiMap map[string]map[string]map[string]string) ([]indexer.EventLabel, []indexer.TransactionLabel, error) {
	return nil, nil, nil
}

func (c *Client) DecodeProtoTransactionsToLabels(transactions []string, blocksCache map[uint64]uint64, abiMap map[string]map[string]map[string]string) ([]indexer.TransactionLabel, error) {
	return nil, nil
}

func (c *Client) FetchAsProtoBlocksWithEvents(from, to *big.Int, debug bool, maxRequests int) ([]proto.Message, []indexer.BlockIndex, []indexer.TransactionIndex, []indexer.LogIndex, uint64, error) {
	return nil, nil, nil, nil, 0, nil
}

// GetLatestBlockNumber returns the latest block number.
func (c *Client) GetLatestBlockNumber() (*big.Int, error) {
	ctx := context.Background()
	var bn uint64
	if err := c.rpcClient.CallContext(ctx, &bn, "starknet_blockNumber"); err != nil {
		return nil, err
	}

	// Convert the uint64 to a big.Int
	blockNumber := new(big.Int).SetUint64(bn)

	return blockNumber, nil
}

func (c *Client) ProcessBlocksToBatch(msgs []proto.Message) (proto.Message, error) {
	return nil, nil
}
