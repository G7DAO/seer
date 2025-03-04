package blockchain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/G7DAO/seer/blockchain/arbitrum_one"
	"github.com/G7DAO/seer/blockchain/arbitrum_sepolia"
	"github.com/G7DAO/seer/blockchain/b3"
	"github.com/G7DAO/seer/blockchain/b3_sepolia"
	seer_common "github.com/G7DAO/seer/blockchain/common"
	"github.com/G7DAO/seer/blockchain/ethereum"
	"github.com/G7DAO/seer/blockchain/game7"
	"github.com/G7DAO/seer/blockchain/game7_orbit_arbitrum_sepolia"
	"github.com/G7DAO/seer/blockchain/game7_testnet"
	"github.com/G7DAO/seer/blockchain/imx_zkevm"
	"github.com/G7DAO/seer/blockchain/imx_zkevm_sepolia"
	"github.com/G7DAO/seer/blockchain/mantle"
	"github.com/G7DAO/seer/blockchain/mantle_sepolia"
	"github.com/G7DAO/seer/blockchain/polygon"
	"github.com/G7DAO/seer/blockchain/ronin"
	"github.com/G7DAO/seer/blockchain/ronin_saigon"
	"github.com/G7DAO/seer/blockchain/sepolia"
	"github.com/G7DAO/seer/blockchain/xai"
	"github.com/G7DAO/seer/blockchain/xai_sepolia"
	"github.com/G7DAO/seer/indexer"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/protobuf/proto"
)

var BlockchainChainIDs = map[string]int64{
	"ethereum":                     1,
	"sepolia":                      11155111,
	"polygon":                      137,
	"arbitrum_one":                 42161,
	"arbitrum_sepolia":             421614,
	"game7":                        2187,
	"game7_orbit_arbitrum_sepolia": 32321,
	"game7_testnet":                13746,
	"mantle":                       5000,
	"mantle_sepolia":               5003,
	"xai":                          660279,
	"xai_sepolia":                  37714555429,
	"imx_zkevm":                    13371,
	"imx_zkevm_sepolia":            13473,
	"b3":                           8333,
	"b3_sepolia":                   1993,
	"ronin":                        2020,
	"ronin_saigon":                 2021,
}

func NewClient(chain, url string, timeout int) (BlockchainClient, error) {
	verifyErr := VerifyChainID(chain, url)
	if verifyErr != nil {
		return nil, verifyErr
	}

	fmt.Printf("Chain: %s\n", chain)
	fmt.Printf("URL: %s\n", url)
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
	} else if chain == "game7" {
		client, err := game7.NewClient(url, timeout)
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
		client, err := b3.NewClient(url, timeout)
		return client, err
	} else if chain == "b3_sepolia" {
		client, err := b3_sepolia.NewClient(url, timeout)
		return client, err
	} else if chain == "ronin" {
		client, err := ronin.NewClient(url, timeout)
		return client, err
	} else if chain == "ronin_saigon" {
		client, err := ronin_saigon.NewClient(url, timeout)
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
	DecodeProtoEntireBlockToLabels(*bytes.Buffer, map[string]map[string]*indexer.AbiEntry, bool, int) ([]indexer.EventLabel, []indexer.TransactionLabel, []indexer.RawTransaction, error)
	DecodeProtoTransactionsToLabels([]string, map[uint64]uint64, map[string]map[string]*indexer.AbiEntry) ([]indexer.TransactionLabel, error)
	ChainType() string
	GetCode(context.Context, common.Address, uint64) ([]byte, error)
	GetTransactionsLabels(uint64, uint64, map[string]map[string]*indexer.AbiEntry, int) ([]indexer.TransactionLabel, map[uint64]seer_common.BlockWithTransactions, error)
	GetEventsLabels(uint64, uint64, map[string]map[string]*indexer.AbiEntry, map[uint64]seer_common.BlockWithTransactions) ([]indexer.EventLabel, error)
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

func DeployBlocksLookUpAndUpdate(blockchain string, rpcUrl string, rpcTimeout int) error {

	// get all abi jobs without deployed block

	chainsAddresses, err := indexer.DBConnection.GetAbiJobsWithoutDeployBlocks(blockchain)

	if err != nil {
		log.Printf("Failed to get abi jobs without deployed blocks: %v", err)
		return err
	}

	if len(chainsAddresses) == 0 {
		log.Printf("No abi jobs without deployed blocks")
		return nil
	}

	for chain, addresses := range chainsAddresses {

		var wg sync.WaitGroup

		sem := make(chan struct{}, 5)  // Semaphore to control
		errChan := make(chan error, 1) // Buffered channel for error handling

		log.Printf("Processing chain: %s with amount of addresses: %d\n", chain, len(addresses))

		for address, ids := range addresses {

			wg.Add(1)
			go func(address string, chain string, ids []string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				client, err := NewClient(chain, rpcUrl, rpcTimeout)

				if err != nil {
					errChan <- err
					return
				}

				// get all abi jobs without deployed block
				deployedBlock, err := FindDeployedBlock(client, address)

				if err != nil {
					errChan <- err
					return
				}

				log.Printf("Deployed block: %d for address: %s in chain: %s\n", deployedBlock, address, chain)

				if deployedBlock != 0 {
					// update abi job with deployed block
					err := indexer.DBConnection.UpdateAbiJobsDeployBlock(deployedBlock, ids)
					if err != nil {
						errChan <- err
						return
					}
				}

			}(address, chain, ids)

		}

		go func() {
			wg.Wait()
			close(errChan)
		}()

		for err := range errChan {
			if err != nil {
				log.Printf("Failed to get deployed block: %v", err)
				return err
			}
		}

	}

	return nil

}

func FindDeployedBlock(client BlockchainClient, address string) (uint64, error) {

	// Binary search by get code

	ctx := context.Background()

	latestBlockNumber, err := client.GetLatestBlockNumber()

	if err != nil {
		return 0, err
	}

	var left uint64 = 1
	var right uint64 = latestBlockNumber.Uint64()
	var code []byte

	for left < right {
		mid := (left + right) / 2

		// with timeout

		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)

		defer cancel()

		code, err = client.GetCode(ctx, common.HexToAddress(address), mid)

		if err != nil {
			log.Printf("Failed to get code: %v", err)
			return 0, err
		}

		if len(code) == 0 {
			left = mid + 1
		} else {
			right = mid
		}

	}

	return left, nil
}

func VerifyChainID(chainName string, rpcURL string) error {
	consent := "y"
	expectedChainID, exists := BlockchainChainIDs[chainName]
	if !exists {
		log.Printf("Unknown blockchain: %s", chainName)
		fmt.Printf("Do you want to continue? (y/n): ")
		fmt.Scanln(&consent)
		if consent != "y" {
			return fmt.Errorf("unknown blockchain: %s", chainName)
		}
		return nil
	}

	// Create a temporary client to query chain ID if it possible
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Printf("Failed to connect to RPC URL: %v", err)
		return fmt.Errorf("failed to connect to RPC URL: %w", err)
	}
	defer client.Close()

	RPCChainID, err := client.ChainID(ctx)
	if err != nil {
		log.Printf("Failed to retrieve chain ID: %v", err)
		return fmt.Errorf("failed to retrieve chain ID: %w", err)
	}

	log.Printf("RPC chain ID: %d", RPCChainID.Int64())

	if RPCChainID.Int64() != expectedChainID {
		log.Printf("Chain ID mismatch: expected %d for %s but got %d from RPC endpoint",
			expectedChainID, chainName, RPCChainID.Int64())
		fmt.Printf("Do you want to continue? (y/n): ")
		fmt.Scanln(&consent)
		if consent != "y" {
			return fmt.Errorf("chain ID mismatch: expected %d for %s but got %d from RPC endpoint",
				expectedChainID, chainName, RPCChainID.Int64())
		}
	}

	return nil
}
