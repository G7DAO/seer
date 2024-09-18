package blockchain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/moonstream-to/seer/blockchain/arbitrum_one"
	"github.com/moonstream-to/seer/blockchain/arbitrum_sepolia"
	seer_common "github.com/moonstream-to/seer/blockchain/common"
	"github.com/moonstream-to/seer/blockchain/ethereum"
	"github.com/moonstream-to/seer/blockchain/game7_orbit_arbitrum_sepolia"
	"github.com/moonstream-to/seer/blockchain/game7_testnet"
	"github.com/moonstream-to/seer/blockchain/imx_zkevm"
	"github.com/moonstream-to/seer/blockchain/imx_zkevm_sepolia"
	"github.com/moonstream-to/seer/blockchain/mantle"
	"github.com/moonstream-to/seer/blockchain/mantle_sepolia"
	"github.com/moonstream-to/seer/blockchain/polygon"
	"github.com/moonstream-to/seer/blockchain/sepolia"
	"github.com/moonstream-to/seer/blockchain/xai"
	"github.com/moonstream-to/seer/blockchain/xai_sepolia"
	"github.com/moonstream-to/seer/indexer"
	"google.golang.org/protobuf/proto"
)

func NewClient(chain, url string, timeout int) (BlockchainClient, error) {
	if chain == "ethereum" {
		client, err := ethereum.NewClient(url, timeout)
		return client, err
	} else if chain == "sepolia" {
		client, err := sepolia.NewClient(url, timeout)
		return client, err
	} else if chain == "polygon" {
		client, err := polygon.NewClient(url, timeout)
		return client, err
	} else if chain == "arbitrum_one" {
		client, err := arbitrum_one.NewClient(url, timeout)
		return client, err
	} else if chain == "arbitrum_sepolia" {
		client, err := arbitrum_sepolia.NewClient(url, timeout)
		return client, err
	} else if chain == "game7_orbit_arbitrum_sepolia" {
		client, err := game7_orbit_arbitrum_sepolia.NewClient(url, timeout)
		return client, err
	} else if chain == "game7_testnet" {
		client, err := game7_testnet.NewClient(url, timeout)
		return client, err
	} else if chain == "mantle" {
		client, err := mantle.NewClient(url, timeout)
		return client, err
	} else if chain == "mantle_sepolia" {
		client, err := mantle_sepolia.NewClient(url, timeout)
		return client, err
	} else if chain == "xai" {
		client, err := xai.NewClient(url, timeout)
		return client, err
	} else if chain == "xai_sepolia" {
		client, err := xai_sepolia.NewClient(url, timeout)
		return client, err
	} else if chain == "imx_zkevm" {
		client, err := imx_zkevm.NewClient(url, timeout)
		return client, err
	} else if chain == "imx_zkevm_sepolia" {
		client, err := imx_zkevm_sepolia.NewClient(url, timeout)
		return client, err
	} else if chain == "b3" {
		client, err := imx_zkevm_sepolia.NewClient(url, timeout)
		return client, err
	} else if chain == "b3_sepolia" {
		client, err := imx_zkevm_sepolia.NewClient(url, timeout)
		return client, err
	} else {
		return nil, errors.New("unsupported chain type")
	}
}

type BlockData struct {
	BlockNumber    uint64
	BlockHash      string
	BlockTimestamp int64
	ChainID        string
	Data           map[string]interface{}
}

type BlockchainClient interface {
	GetLatestBlockNumber() (*big.Int, error)
	FetchAsProtoBlocksWithEvents(*big.Int, *big.Int, bool, int) ([]proto.Message, []indexer.BlockIndex, uint64, error)
	ProcessBlocksToBatch([]proto.Message) (proto.Message, error)
	DecodeProtoEntireBlockToJson(*bytes.Buffer) (*seer_common.BlocksBatchJson, error)
	DecodeProtoEntireBlockToLabels(*bytes.Buffer, map[string]map[string]map[string]string) ([]indexer.EventLabel, []indexer.TransactionLabel, error)
	DecodeProtoTransactionsToLabels([]string, map[uint64]uint64, map[string]map[string]map[string]string) ([]indexer.TransactionLabel, error)
	TransactionReceipt(context.Context, string) (*types.Receipt, error)
	ChainType() string
}

func GetLatestBlockNumberWithRetry(client BlockchainClient, retryAttempts int, retryWaitTime time.Duration) (*big.Int, error) {
	for {
		latestBlockNumber, latestErr := client.GetLatestBlockNumber()
		if latestErr != nil {
			log.Printf("Failed to get latest block number: %v", latestErr)

			// Retry the operation
			retryAttempts--
			if retryAttempts == 0 {
				return nil, fmt.Errorf("failed to get latest block number after several attempts")
			}

			time.Sleep(retryWaitTime)
			continue
		}
		return latestBlockNumber, nil
	}
}

func CrawlEntireBlocks(client BlockchainClient, startBlock *big.Int, endBlock *big.Int, debug bool, maxRequests int) ([]proto.Message, []indexer.BlockIndex, uint64, error) {
	log.Printf("Operates with batch of blocks: %d-%d", startBlock, endBlock)

	blocks, blocksIndex, blocksSize, pBlockErr := client.FetchAsProtoBlocksWithEvents(startBlock, endBlock, debug, maxRequests)
	if pBlockErr != nil {
		return nil, nil, 0, pBlockErr
	}

	return blocks, blocksIndex, blocksSize, nil
}

func DecodeTransactionInputData(contractABI *abi.ABI, data []byte) {
	methodSigData := data[:4]
	inputsSigData := data[4:]
	method, err := contractABI.MethodById(methodSigData)
	if err != nil {
		log.Fatal(err)
	}
	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(inputsMap)
	}
	fmt.Printf("Method Name: %s\n", method.Name)
	fmt.Printf("Method inputs: %v\n", inputsMap)
}

func GetDeployedContracts(client BlockchainClient, BlocksBatch *seer_common.BlocksBatchJson) (map[string]indexer.EvmContract, error) {
	deployedContracts := make(map[string]indexer.EvmContract)
	for _, block := range BlocksBatch.Blocks {
		for _, transaction := range block.Transactions {
			if (transaction.ToAddress == "" || transaction.ToAddress == "0x") && transaction.Input != "0x" {

				// make get receipt call

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				transactionReceipt, err := client.TransactionReceipt(ctx, transaction.Hash)
				if err != nil {
					log.Printf("Failed to get transaction receipt: %v", err)
					continue
				}

				if transactionReceipt.ContractAddress == (common.Address{}) {
					continue
				}

				// Extract contract address as string
				contractAddress := transactionReceipt.ContractAddress.String() // Assuming ContractAddress is of type common.Address

				timestamp, err := strconv.ParseUint(block.Timestamp, 10, 64)
				if err != nil {
					log.Printf("Failed to parse block timestamp: %v", err)
					// Handle the error appropriately (e.g., continue to the next block or return)
					return nil, err
				}

				transactionIndex, err := strconv.ParseUint(transaction.TransactionIndex, 10, 64)
				if err != nil {
					log.Printf("Failed to parse transaction index: %v", err)
					// Handle the error appropriately (e.g., continue to the next transaction or return)
					return nil, err
				}

				// Assign EvmContract to the map
				deployedContracts[contractAddress] = indexer.EvmContract{
					Address:                  contractAddress,
					Bytecode:                 nil,
					DeployedBytecode:         transaction.Input,
					Abi:                      nil,
					DeployedAtBlockNumber:    transactionReceipt.BlockNumber.Uint64(),
					DeployedAtBlockHash:      block.Hash,
					DeployedAtBlockTimestamp: timestamp,
					TransactionHash:          transaction.Hash,
					TransactionIndex:         transactionIndex,
					Name:                     nil,
					Statistics:               nil,
					SupportedStandards:       nil,
				}

				// Try to get bytecode
			}
		}
	}

	log.Printf("Found %d deployed contracts", len(deployedContracts))

	return deployedContracts, nil
}
