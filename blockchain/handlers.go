package blockchain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
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
	"google.golang.org/protobuf/proto"
)

func NewClient(chain, url string, timeout int) (BlockchainClient, error) {
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

func CommonClient(url string) (*seer_common.EvmClient, error) {
	client, err := seer_common.NewEvmClient(url, 4)
	return client, err
}

type ChainInfo struct {
	ChainType string                 `json:"chainType"` // L1 or L2
	ChainID   *big.Int               `json:"chainID"`   // Chain ID
	Block     *seer_common.BlockJson `json:"block"`
}

type BlockchainTemplateData struct {
	BlockchainName      string
	BlockchainNameLower string
	ChainDashName       string
	IsSideChain         bool
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
	DecodeProtoEntireBlockToLabels(*bytes.Buffer, map[string]map[string]*indexer.AbiEntry, int) ([]indexer.EventLabel, []indexer.TransactionLabel, error)
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

func DeployBlocksLookUpAndUpdate(blockchain string) error {

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

				client, err := NewClient(chain, BlockchainURLs[chain], 4)

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

func CollectChainInfo(client seer_common.EvmClient) (*ChainInfo, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	var chainInfo ChainInfo

	// Get chain id
	chainID, err := client.GetChainID(ctx)

	if err != nil {
		log.Printf("Failed to get chain id: %v", err)
		return nil, err
	}

	chainInfo.ChainID = chainID

	// Get latest block number
	latestBlockNumber, err := client.GetLatestBlockNumber()

	if err != nil {
		log.Printf("Failed to get latest block number: %v", err)
		return nil, err
	}

	// Get block by number

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	block, err := client.GetBlockByNumber(ctx, latestBlockNumber, true)

	if err != nil {
		log.Printf("Failed to get block by number: %v", err)
		return nil, err
	}

	chainInfo.Block = block

	// Initialize a score
	l2Score := 0

	// Method 1: Check for L2-specific fields

	// Check for L2-specific fields
	if block.L1BlockNumber != "" {
		l2Score++
	}
	if block.SendRoot != "" {
		l2Score++
	}

	if l2Score >= 1 {
		chainInfo.ChainType = "L2"
	} else {
		chainInfo.ChainType = "L1"
	}

	return &chainInfo, nil

}

func GenerateProtoFile(chainName string, goPackage string, isL2 bool, outputPath string) error {
	tmplData := seer_common.ProtoTemplateData{
		GoPackage: goPackage,
		ChainName: chainName,
		IsL2:      isL2,
	}

	// Read the template file
	tmplContent, err := ioutil.ReadFile("./blockchain/common/evm_proto_template.proto.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read template file: %v", err)
	}

	// Parse the template
	tmpl, err := template.New("proto").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// Execute the template with data
	var output bytes.Buffer
	err = tmpl.Execute(&output, tmplData)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Write the output to the specified path
	err = ioutil.WriteFile(outputPath, output.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write proto file: %v", err)
	}

	fmt.Printf("Protobuf definitions generated successfully at %s.\n", outputPath)
	return nil
}

func ProtoKeys(msg proto.Message) []string {
	m := make([]string, 0)
	messageReflect := msg.ProtoReflect()
	fields := messageReflect.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fieldDesc := fields.Get(i)
		fieldName := string(fieldDesc.Name())
		m = append(m, fieldName)

	}
	return m
}

func GenerateModelsFiles(data BlockchainTemplateData, modelsPath string) error {

	// Define the directory path
	//dirPath := filepath.Join(modelsPath, data.BlockchainNameLower)

	// Check if the directory exists
	// if _, err := os.Stat(dirPath); os.IsNotExist(err) {
	// 	// Create the directory along with any necessary parents
	// 	if err := os.MkdirAll(dirPath, 0755); err != nil {
	// 		return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
	// 	}
	// }
	templateIndexFile := fmt.Sprintf("%s/models_indexes.py.tmpl", modelsPath)
	templateLabelsFile := fmt.Sprintf("%s/models_labels.py.tmpl", modelsPath)

	// Read the template file
	tmplIndexContent, err := ioutil.ReadFile(templateIndexFile)
	if err != nil {
		return fmt.Errorf("failed to read template file: %v", err)
	}

	tmplLabelsContent, err := ioutil.ReadFile(templateLabelsFile)
	if err != nil {
		return fmt.Errorf("failed to read template file: %v", err)
	}

	// Parse the template
	tmplIndex, err := template.New("indexes").Parse(string(tmplIndexContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	tmplLabels, err := template.New("labels").Parse(string(tmplLabelsContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// Execute the template with data
	var outputIndex bytes.Buffer
	err = tmplIndex.Execute(&outputIndex, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	var outputLabels bytes.Buffer
	err = tmplLabels.Execute(&outputLabels, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	outputPath := fmt.Sprintf("%s/models_indexes/%s_indexes.py", modelsPath, data.BlockchainNameLower)

	// Write the output to the specified path
	err = ioutil.WriteFile(outputPath, outputIndex.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write models file: %v", err)
	}

	outputPath = fmt.Sprintf("%s/models/%s_labels.py", modelsPath, data.BlockchainNameLower)
	err = ioutil.WriteFile(outputPath, outputLabels.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write models file: %v", err)
	}

	fmt.Printf("Models generated successfully at %s.\n", outputPath)
	return nil
}

func GenerateDeploy(chain BlockchainTemplateData, deployPath string) error {
	funcMap := template.FuncMap{
		"replaceUnderscoreWithDash": func(s string) string {
			return strings.ReplaceAll(s, "_", "-")
		},
		"replaceUnderscoreWithSpace": func(s string) string {
			return strings.ReplaceAll(s, "_", " ")
		},
	}

	templateDir := "deploy/templates"

	// Read all entries in the template directory
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories
		}

		templateFileName := entry.Name()
		// Only process files with .tmpl extension
		if !strings.HasSuffix(templateFileName, ".tmpl") {
			continue
		}

		templateFilePath := filepath.Join(templateDir, templateFileName)

		// Parse the template file
		tmpl, parseErr := template.New(templateFileName).Funcs(funcMap).ParseFiles(templateFilePath)
		if parseErr != nil {
			return parseErr
		}

		// Generate the output file name
		outputFileName := generateOutputFileName(templateFileName, chain.ChainDashName)
		outputFilePath := filepath.Join(deployPath, outputFileName)

		// Create the output file
		outputFile, createErr := os.Create(outputFilePath)
		if createErr != nil {
			return createErr
		}
		defer outputFile.Close()

		// Execute the template and write to the output file
		if err := tmpl.Execute(outputFile, chain); err != nil {
			return err
		}

		log.Printf("File generated successfully: %s", outputFilePath)
	}

	return nil
}

// Helper function to generate output file names dynamically
func generateOutputFileName(templateFileName, chainName string) string {
	// Remove the .tmpl extension
	baseFileName := strings.TrimSuffix(templateFileName, ".tmpl")

	// Split the file name into name and extension
	ext := filepath.Ext(baseFileName)
	name := strings.TrimSuffix(baseFileName, ext)

	// Insert the chain name before the extension
	outputFileName := fmt.Sprintf("%s-%s%s", name, chainName, ext)
	return outputFileName
}
