package starknet

import (
	"bytes"
	"math/big"

	"github.com/NethermindEth/starknet.go/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/protobuf/proto"

	seer_common "github.com/moonstream-to/seer/blockchain/common"
	"github.com/moonstream-to/seer/indexer"
)

func NewClient(url string, timeout int) (*Client, error) {
	rpcClient, clientErr := rpc.NewClient(url)
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

func (c *Client) GetLatestBlockNumber() (*big.Int, error) {
	return nil, nil
}

func (c *Client) ProcessBlocksToBatch(msgs []proto.Message) (proto.Message, error) {
	return nil, nil
}
