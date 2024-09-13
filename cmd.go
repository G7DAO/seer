package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	seer_blockchain "github.com/moonstream-to/seer/blockchain"
	"github.com/moonstream-to/seer/crawler"
	"github.com/moonstream-to/seer/evm"
	"github.com/moonstream-to/seer/indexer"
	"github.com/moonstream-to/seer/starknet"
	"github.com/moonstream-to/seer/storage"
	"github.com/moonstream-to/seer/synchronizer"
	"github.com/moonstream-to/seer/version"
)

func CreateRootCommand() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "seer",
		Short: "Seer: Generate interfaces and crawlers from various blockchains",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	completionCmd := CreateCompletionCommand(rootCmd)
	versionCmd := CreateVersionCommand()
	blockchainCmd := CreateBlockchainCommand()
	starknetCmd := CreateStarknetCommand()
	crawlerCmd := CreateCrawlerCommand()
	inspectorCmd := CreateInspectorCommand()
	evmCmd := CreateEVMCommand()
	synchronizerCmd := CreateSynchronizerCommand()
	abiCmd := CreateAbiCommand()
	dbCmd := CreateDatabaseOperationCommand()
	historicalSyncCmd := CreateHistoricalSyncCommand()
	rootCmd.AddCommand(completionCmd, versionCmd, blockchainCmd, starknetCmd, evmCmd, crawlerCmd, inspectorCmd, synchronizerCmd, abiCmd, dbCmd, historicalSyncCmd)

	// By default, cobra Command objects write to stderr. We have to forcibly set them to output to
	// stdout.
	rootCmd.SetOut(os.Stdout)

	return rootCmd
}

func CreateCompletionCommand(rootCmd *cobra.Command) *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts for seer",
		Long: `Generate shell completion scripts for seer.

The command for each shell will print a completion script to stdout. You can source this script to get
completions in your current shell session. You can add this script to the completion directory for your
shell to get completions for all future sessions.

For example, to activate bash completions in your current shell:
		$ . <(seer completion bash)

To add seer completions for all bash sessions:
		$ seer completion bash > /etc/bash_completion.d/seer_completions`,
	}

	bashCompletionCmd := &cobra.Command{
		Use:   "bash",
		Short: "bash completions for seer",
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(cmd.OutOrStdout())
		},
	}

	zshCompletionCmd := &cobra.Command{
		Use:   "zsh",
		Short: "zsh completions for seer",
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenZshCompletion(cmd.OutOrStdout())
		},
	}

	fishCompletionCmd := &cobra.Command{
		Use:   "fish",
		Short: "fish completions for seer",
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		},
	}

	powershellCompletionCmd := &cobra.Command{
		Use:   "powershell",
		Short: "powershell completions for seer",
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
		},
	}

	completionCmd.AddCommand(bashCompletionCmd, zshCompletionCmd, fishCompletionCmd, powershellCompletionCmd)

	return completionCmd
}

func CreateVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of seer that you are currently using",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(version.SeerVersion)
		},
	}

	return versionCmd
}

func CreateBlockchainCommand() *cobra.Command {
	blockchainCmd := &cobra.Command{
		Use:   "blockchain",
		Short: "Generate methods and types for different blockchains",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	blockchainGenerateCmd := CreateBlockchainGenerateCommand()
	blockchainCmd.AddCommand(blockchainGenerateCmd)

	return blockchainCmd
}

type BlockchainTemplateData struct {
	BlockchainName      string
	BlockchainNameLower string
	IsSideChain         bool
}

func CreateBlockchainGenerateCommand() *cobra.Command {
	var blockchainNameLower string
	var sideChain bool

	blockchainGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate methods and types for different blockchains from template",
		RunE: func(cmd *cobra.Command, args []string) error {
			dirPath := filepath.Join(".", "blockchain", blockchainNameLower)
			blockchainNameFilePath := filepath.Join(dirPath, fmt.Sprintf("%s.go", blockchainNameLower))

			var blockchainName string
			blockchainNameList := strings.Split(blockchainNameLower, "_")
			for _, w := range blockchainNameList {
				blockchainName += strings.Title(w)
			}

			// Read and parse the template file
			tmpl, parseErr := template.ParseFiles("blockchain/blockchain.go.tmpl")
			if parseErr != nil {
				return parseErr
			}

			// Create output file
			if _, statErr := os.Stat(dirPath); os.IsNotExist(statErr) {
				mkdirErr := os.Mkdir(dirPath, 0775)
				if mkdirErr != nil {
					return mkdirErr
				}
			}

			outputFile, createErr := os.Create(blockchainNameFilePath)
			if createErr != nil {
				return createErr
			}
			defer outputFile.Close()

			// Execute template and write to output file
			data := BlockchainTemplateData{
				BlockchainName:      blockchainName,
				BlockchainNameLower: blockchainNameLower,
				IsSideChain:         sideChain,
			}
			execErr := tmpl.Execute(outputFile, data)
			if execErr != nil {
				return execErr
			}

			log.Printf("Blockchain file generated successfully: %s", blockchainNameFilePath)

			return nil
		},
	}

	blockchainGenerateCmd.Flags().StringVarP(&blockchainNameLower, "name", "n", "", "The name of the blockchain to generate lowercase (example: 'arbitrum_one')")
	blockchainGenerateCmd.Flags().BoolVar(&sideChain, "side-chain", false, "Set this flag to extend Blocks and Transactions with additional fields for side chains (default: false)")

	return blockchainGenerateCmd
}

func CreateStarknetCommand() *cobra.Command {
	starknetCmd := &cobra.Command{
		Use:   "starknet",
		Short: "Generate interfaces and crawlers for Starknet contracts",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	starknetABIParseCmd := CreateStarknetParseCommand()
	starknetABIGenGoCmd := CreateStarknetGenerateCommand()
	starknetCmd.AddCommand(starknetABIParseCmd, starknetABIGenGoCmd)

	return starknetCmd
}

func CreateCrawlerCommand() *cobra.Command {
	var startBlock, finalBlock, confirmations, batchSize int64
	var timeout, threads, protoTimeLimit, retryWait, retryMultiplier int
	var protoSizeLimit uint64
	var chain, baseDir string

	crawlerCmd := &cobra.Command{
		Use:   "crawler",
		Short: "Start crawlers for various blockchains",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			indexerErr := indexer.CheckVariablesForIndexer()
			if indexerErr != nil {
				return indexerErr
			}

			storageErr := storage.CheckVariablesForStorage()
			if storageErr != nil {
				return storageErr
			}

			crawlerErr := crawler.CheckVariablesForCrawler()
			if crawlerErr != nil {
				return crawlerErr
			}

			blockchainErr := seer_blockchain.CheckVariablesForBlockchains()
			if blockchainErr != nil {
				return blockchainErr
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			indexer.InitDBConnection()

			newCrawler, crawlerError := crawler.NewCrawler(chain, startBlock, finalBlock, confirmations, batchSize, timeout, baseDir, protoSizeLimit, protoTimeLimit, retryWait, retryMultiplier)
			if crawlerError != nil {
				return crawlerError
			}

			latestBlockNumber, latestErr := newCrawler.Client.GetLatestBlockNumber()
			if latestErr != nil {
				return fmt.Errorf("Failed to get latest block number: %v", latestErr)
			}

			if startBlock > latestBlockNumber.Int64() {
				log.Fatalf("Start block could not be greater then latest block number at blockchain")
			}

			crawler.CurrentBlockchainState.RaiseLatestBlockNumber(latestBlockNumber)

			newCrawler.Start(threads)

			return nil
		},
	}

	crawlerCmd.Flags().StringVar(&chain, "chain", "ethereum", "The blockchain to crawl")
	crawlerCmd.Flags().Int64Var(&startBlock, "start-block", 0, "The block number to start crawling from (default: fetch from database, if it is empty, run from latestBlockNumber minus shift)")
	crawlerCmd.Flags().Int64Var(&finalBlock, "final-block", 0, "The block number to end crawling at")
	crawlerCmd.Flags().IntVar(&timeout, "timeout", 30, "The timeout for the crawler in seconds")
	crawlerCmd.Flags().IntVar(&threads, "threads", 1, "Number of go-routines for concurrent crawling")
	crawlerCmd.Flags().Int64Var(&confirmations, "confirmations", 10, "The number of confirmations to consider for block finality")
	crawlerCmd.Flags().Int64Var(&batchSize, "batch-size", 10, "Dynamically changed maximum number of blocks to crawl in each batch")
	crawlerCmd.Flags().StringVar(&baseDir, "base-dir", "", "The base directory to store the crawled data")
	crawlerCmd.Flags().Uint64Var(&protoSizeLimit, "proto-size-limit", 25, "Proto file size limit in Mb")
	crawlerCmd.Flags().IntVar(&protoTimeLimit, "proto-time-limit", 300, "Proto time limit in seconds")
	crawlerCmd.Flags().IntVar(&retryWait, "retry-wait", 5000, "The wait time for the crawler in milliseconds before it try to fetch new block")
	crawlerCmd.Flags().IntVar(&retryMultiplier, "retry-multiplier", 24, "Multiply wait time to get max waiting time before fetch new block")

	return crawlerCmd
}

func CreateSynchronizerCommand() *cobra.Command {
	var startBlock, endBlock, batchSize uint64
	var timeout int
	var chain, baseDir, customerDbUriFlag string

	synchronizerCmd := &cobra.Command{
		Use:   "synchronizer",
		Short: "Decode the crawled data from various blockchains",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			indexerErr := indexer.CheckVariablesForIndexer()
			if indexerErr != nil {
				return indexerErr
			}

			storageErr := storage.CheckVariablesForStorage()
			if storageErr != nil {
				return storageErr
			}

			crawlerErr := crawler.CheckVariablesForCrawler()
			if crawlerErr != nil {
				return crawlerErr
			}

			syncErr := synchronizer.CheckVariablesForSynchronizer()
			if syncErr != nil {
				return syncErr
			}

			blockchainErr := seer_blockchain.CheckVariablesForBlockchains()
			if blockchainErr != nil {
				return blockchainErr
			}

			if chain == "" {
				return fmt.Errorf("blockchain is required via --chain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			indexer.InitDBConnection()

			newSynchronizer, synchonizerErr := synchronizer.NewSynchronizer(chain, baseDir, startBlock, endBlock, batchSize, timeout)
			if synchonizerErr != nil {
				return synchonizerErr
			}

			latestBlockNumber, latestErr := newSynchronizer.Client.GetLatestBlockNumber()
			if latestErr != nil {
				return fmt.Errorf("Failed to get latest block number: %v", latestErr)
			}

			if startBlock > latestBlockNumber.Uint64() {
				log.Fatalf("Start block could not be greater then latest block number at blockchain")
			}

			crawler.CurrentBlockchainState.RaiseLatestBlockNumber(latestBlockNumber)

			newSynchronizer.Start(customerDbUriFlag)

			return nil
		},
	}

	synchronizerCmd.Flags().StringVar(&chain, "chain", "ethereum", "The blockchain to crawl (default: ethereum)")
	synchronizerCmd.Flags().Uint64Var(&startBlock, "start-block", 0, "The block number to start decoding from (default: latest block)")
	synchronizerCmd.Flags().Uint64Var(&endBlock, "end-block", 0, "The block number to end decoding at (default: latest block)")
	synchronizerCmd.Flags().StringVar(&baseDir, "base-dir", "", "The base directory to store the crawled data (default: '')")
	synchronizerCmd.Flags().IntVar(&timeout, "timeout", 30, "The timeout for the crawler in seconds (default: 30)")
	synchronizerCmd.Flags().Uint64Var(&batchSize, "batch-size", 100, "The number of blocks to crawl in each batch (default: 100)")
	synchronizerCmd.Flags().StringVar(&customerDbUriFlag, "customer-db-uri", "", "Set customer database URI for development. This workflow bypass fetching customer IDs and its database URL connection strings from mdb-v3-controller API")

	return synchronizerCmd
}

type BlockInspectItem struct {
	StartBlock int64
	EndBlock   int64
}

func CreateInspectorCommand() *cobra.Command {
	inspectorCmd := &cobra.Command{
		Use:   "inspector",
		Short: "Inspect storage and database consistency",
	}

	var chain, baseDir, delim, returnFunc, batch string
	var timeout int

	readCommand := &cobra.Command{
		Use:   "read",
		Short: "Read and decode indexed proto data from storage",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			storageErr := storage.CheckVariablesForStorage()
			if storageErr != nil {
				return storageErr
			}

			crawlerErr := crawler.CheckVariablesForCrawler()
			if crawlerErr != nil {
				return crawlerErr
			}

			if batch == "" {
				return errors.New("batch is required via --batch")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := filepath.Join(baseDir, crawler.SeerCrawlerStoragePrefix, "data", chain)
			storageInstance, newStorageErr := storage.NewStorage(storage.SeerCrawlerStorageType, basePath)
			if newStorageErr != nil {
				return newStorageErr
			}

			targetFilePath := filepath.Join(basePath, batch, "data.proto")
			rawData, readErr := storageInstance.Read(targetFilePath)
			if readErr != nil {
				return readErr
			}

			client, cleintErr := seer_blockchain.NewClient(chain, crawler.BlockchainURLs[chain], timeout)
			if cleintErr != nil {
				return cleintErr
			}

			output, decErr := client.DecodeProtoEntireBlockToJson(&rawData)
			if decErr != nil {
				return decErr
			}

			jsonOutput, marErr := json.Marshal(output)
			if marErr != nil {
				return marErr
			}

			fmt.Println(string(jsonOutput))

			return nil
		},
	}

	readCommand.Flags().StringVar(&chain, "chain", "ethereum", "The blockchain to crawl (default: ethereum)")
	readCommand.Flags().StringVar(&baseDir, "base-dir", "", "The base directory to store the crawled data (default: '')")
	readCommand.Flags().StringVar(&batch, "batch", "", "What batch to read")

	var storageVerify bool

	dbCommand := &cobra.Command{
		Use:   "db",
		Short: "Inspect database consistency",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if chain == "" {
				return fmt.Errorf("blockchain is required via --chain")
			}

			crawlerErr := crawler.CheckVariablesForCrawler()
			if crawlerErr != nil {
				return crawlerErr
			}

			storageErr := storage.CheckVariablesForStorage()
			if storageErr != nil {
				return storageErr
			}

			indexerErr := indexer.CheckVariablesForIndexer()
			if indexerErr != nil {
				return indexerErr
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			indexer.InitDBConnection()

			ctx := context.Background()
			firstBlock, firstErr := indexer.DBConnection.GetEdgeDBBlock(ctx, chain, "first")
			if firstErr != nil {
				return firstErr
			}
			lastBlock, lastErr := indexer.DBConnection.GetEdgeDBBlock(ctx, chain, "last")
			if firstErr != nil {
				return lastErr
			}

			firstPathSlice := strings.Split(firstBlock.Path, "/")
			firstPathBatch := firstPathSlice[len(firstPathSlice)-2]

			lastPathSlice := strings.Split(lastBlock.Path, "/")
			lastPathBatch := lastPathSlice[len(lastPathSlice)-2]

			fmt.Printf("First batch blocks in database: %s\n", firstPathBatch)
			fmt.Printf("Last batch blocks in database: %s\n", lastPathBatch)

			if storageVerify {
				basePath := filepath.Join(baseDir, crawler.SeerCrawlerStoragePrefix, "data", chain)
				storageInstance, newStorageErr := storage.NewStorage(storage.SeerCrawlerStorageType, basePath)
				if newStorageErr != nil {
					return newStorageErr
				}

				listReturnFunc := storage.GCSListReturnNameFunc

				firstItems, firstListErr := storageInstance.List(ctx, delim, firstPathBatch, timeout, listReturnFunc)
				if firstListErr != nil {
					return firstListErr
				}

				lastItems, lastListErr := storageInstance.List(ctx, delim, lastPathBatch, timeout, listReturnFunc)
				if lastListErr != nil {
					return lastListErr
				}

				fmt.Println("First batch in storage:")
				for _, item := range firstItems {
					fmt.Printf("- %s\n", item)
				}

				fmt.Println("Last batch in storage:")
				for _, item := range lastItems {
					fmt.Printf("- %s\n", item)
				}
			}

			// TODO(kompotkot): Write inspect of missing blocks in database

			return nil
		},
	}

	dbCommand.Flags().StringVar(&chain, "chain", "", "The blockchain to crawl")
	dbCommand.Flags().BoolVar(&storageVerify, "storage-verify", false, "Set this flag to verify storage data by path (default: false)")

	storageCommand := &cobra.Command{
		Use:   "storage",
		Short: "Inspect filesystem, gcp-storage, aws-bucket consistency",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			storageErr := storage.CheckVariablesForStorage()
			if storageErr != nil {
				return storageErr
			}

			crawlerErr := crawler.CheckVariablesForCrawler()
			if crawlerErr != nil {
				return crawlerErr
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			basePath := filepath.Join(baseDir, crawler.SeerCrawlerStoragePrefix, "data", chain)
			storageInstance, newStorageErr := storage.NewStorage(storage.SeerCrawlerStorageType, basePath)
			if newStorageErr != nil {
				return newStorageErr
			}

			// Only for gcp-storage type.
			// Created for different manipulations what requires to list,
			// if value set to prefix, required to set delim = '/'
			var listReturnFunc storage.ListReturnFunc
			switch storage.SeerCrawlerStorageType {
			case "gcp-storage":
				switch returnFunc {
				case "prefix":
					listReturnFunc = storage.GCSListReturnPrefixFunc
				default:
					listReturnFunc = storage.GCSListReturnNameFunc
				}
			default:
				listReturnFunc = func(item any) string { return fmt.Sprintf("%v", item) }
			}

			items, listErr := storageInstance.List(ctx, delim, "", timeout, listReturnFunc)
			if listErr != nil {
				return listErr
			}

			itemsMap := make(map[string]BlockInspectItem)
			previousMapItemKey := ""

			for _, item := range items {
				itemSlice := strings.Split(item, "/")
				blockNums := itemSlice[len(itemSlice)-2]

				blockNumsSlice := strings.Split(blockNums, "-")

				blockNumS, atoiErrS := strconv.ParseInt(blockNumsSlice[0], 10, 64)
				if atoiErrS != nil {
					log.Printf("Unable to parse blockNumS from %s", blockNumsSlice[0])
					continue
				}
				blockNumF, atoiErrF := strconv.ParseInt(blockNumsSlice[1], 10, 64)
				if atoiErrF != nil {
					log.Printf("Unable to parse blockNumS from %s", blockNumsSlice[1])
					continue
				}

				if previousMapItemKey != blockNums && previousMapItemKey != "" {
					diff := blockNumS - itemsMap[previousMapItemKey].EndBlock
					if diff <= 0 {
						fmt.Printf("Found incorrect blocks order between batches: %s -> %s\n", previousMapItemKey, blockNums)
					} else if diff > 1 {
						fmt.Printf("Found missing %d blocks during batches: %s -> %s\n", diff, previousMapItemKey, blockNums)
					}
				}

				previousMapItemKey = blockNums
				itemsMap[blockNums] = BlockInspectItem{StartBlock: blockNumS, EndBlock: blockNumF}
			}

			log.Printf("Processed %d items", len(itemsMap))

			return nil
		},
	}

	storageCommand.Flags().StringVar(&chain, "chain", "ethereum", "The blockchain to crawl (default: ethereum)")
	storageCommand.Flags().StringVar(&baseDir, "base-dir", "", "The base directory to store the crawled data (default: '')")
	storageCommand.Flags().StringVar(&delim, "delim", "", "Only for gcp-storage. The delimiter argument can be used to restrict the results to only the objects in the given 'directory'")
	storageCommand.Flags().StringVar(&returnFunc, "return-func", "", "Which function use for return")
	storageCommand.Flags().IntVar(&timeout, "timeout", 180, "List timeout (default: 180)")

	inspectorCmd.AddCommand(storageCommand, readCommand, dbCommand)

	return inspectorCmd
}

func CreateAbiCommand() *cobra.Command {
	abiCmd := &cobra.Command{
		Use:   "abi",
		Short: "General commands for working with ABIs",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	abiParseCmd := CreateAbiParseCommand()
	abiEnsureSelectorsCmd := CreateAbiEnsureSelectorsCommand()
	abiCmd.AddCommand(abiParseCmd)
	abiCmd.AddCommand(abiEnsureSelectorsCmd)

	return abiCmd
}

func CreateAbiParseCommand() *cobra.Command {

	var inFile, outFile string
	var rawABI []byte
	var readErr error

	abiParseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse an ABI and return seer's interal representation of that ABI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if inFile != "" {
				rawABI, readErr = os.ReadFile(inFile)
			} else {
				rawABI, readErr = io.ReadAll(os.Stdin)
			}
			return readErr
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			parsedABI, parseErr := indexer.PrintContractSignatures(string(rawABI))
			if parseErr != nil {
				return parseErr
			}

			content, marshalErr := json.Marshal(parsedABI)
			if marshalErr != nil {
				return marshalErr
			}

			if outFile != "" {
				writeErr := os.WriteFile(outFile, content, 0644)
				if writeErr != nil {
					return writeErr
				}
			} else {
				cmd.Println("ABI parsed successfully:")
			}

			return nil
		},
	}

	abiParseCmd.Flags().StringVarP(&inFile, "abi", "a", "", "Path to contract ABI (default stdin)")
	abiParseCmd.Flags().StringVarP(&outFile, "out", "o", "", "Path to write the output (default stdout)")

	return abiParseCmd
}

func CreateAbiEnsureSelectorsCommand() *cobra.Command {

	var chain, outFilePath string
	var WriteToDB bool

	abiEnsureSelectorsCmd := &cobra.Command{
		Use:   "ensure-selectors",
		Short: "Ensure that all ABI functions have selectors",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			indexerErr := indexer.CheckVariablesForIndexer()
			if indexerErr != nil {
				return indexerErr
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			indexer.InitDBConnection()

			updateErr := indexer.DBConnection.EnsureCorrectSelectors(chain, WriteToDB, outFilePath, []string{})
			if updateErr != nil {
				return updateErr
			}
			return nil
		},
	}

	abiEnsureSelectorsCmd.Flags().StringVarP(&chain, "chain", "c", "", "The blockchain to crawl")
	abiEnsureSelectorsCmd.Flags().BoolVar(&WriteToDB, "write-to-db", false, "Set this flag to write the correct selectors to the database (default: false)")
	abiEnsureSelectorsCmd.Flags().StringVarP(&outFilePath, "out-file", "o", "./missing-selectors.txt", "The file to write the output to (default: stdout)")
	return abiEnsureSelectorsCmd
}

func CreateDatabaseOperationCommand() *cobra.Command {
	databaseCmd := &cobra.Command{
		Use:   "databases",
		Short: "Operations for database",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	indexCommand := &cobra.Command{
		Use:   "index",
		Short: "Actions for index database",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var chain string
	var batchLimit uint64
	var sleepTime int

	cleanCommand := &cobra.Command{
		Use:   "clean",
		Short: "Clean the database transactions and logs indexes",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			indexerErr := indexer.CheckVariablesForIndexer()
			if indexerErr != nil {
				return indexerErr
			}

			indexer.InitDBConnection()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			cleanErr := indexer.DBConnection.CleanIndexes(chain, batchLimit, sleepTime)
			if cleanErr != nil {
				return cleanErr
			}

			return nil
		},
	}

	cleanCommand.Flags().StringVar(&chain, "chain", "ethereum", "The blockchain to crawl (default: ethereum)")
	cleanCommand.Flags().Uint64Var(&batchLimit, "batch-limit", 1000, "The number of rows to delete in each batch (default: 1000)")
	cleanCommand.Flags().IntVar(&sleepTime, "sleep-time", 1, "The time to sleep between batches in seconds (default: 1)")

	indexCommand.AddCommand(cleanCommand)

	deploymentBlocksCommand := &cobra.Command{
		Use:   "deployment-blocks",
		Short: "Get deployment blocks from address in abi jobs",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			indexerErr := indexer.CheckVariablesForIndexer()
			if indexerErr != nil {
				return indexerErr
			}
			blockchainErr := seer_blockchain.CheckVariablesForBlockchains()
			if blockchainErr != nil {
				return blockchainErr
			}
			indexer.InitDBConnection()

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			deploymentBlocksErr := seer_blockchain.DeployBlocksLookUpAndUpdate()
			if deploymentBlocksErr != nil {
				return deploymentBlocksErr
			}

			return nil
		},
	}

	var jobChain, address, abiFile, customerId, userId string

	createJobsCommand := &cobra.Command{
		Use:   "create-jobs",
		Short: "Create jobs for ABI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			indexerErr := indexer.CheckVariablesForIndexer()
			if indexerErr != nil {
				return indexerErr
			}

			indexer.InitDBConnection()

			blockchainErr := seer_blockchain.CheckVariablesForBlockchains()
			if blockchainErr != nil {
				return blockchainErr
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			// check if the chain is supported
			if _, ok := seer_blockchain.BlockchainURLs[jobChain]; !ok {
				return fmt.Errorf("chain %s is not supported", jobChain)
			}

			client, clientErr := seer_blockchain.NewClient(jobChain, seer_blockchain.BlockchainURLs[jobChain], 30)
			if clientErr != nil {
				return clientErr
			}

			// detect deploy block

			deploy_block, deployErr := seer_blockchain.FindDeployedBlock(client, address)

			if deployErr != nil {
				return deployErr
			}

			createJobsErr := indexer.DBConnection.CreateJobsFromAbi(jobChain, address, abiFile, customerId, userId, deploy_block)
			if createJobsErr != nil {
				return createJobsErr
			}

			return nil
		},
	}

	createJobsCommand.Flags().StringVar(&jobChain, "jobChain", "ethereum", "The blockchain to crawl (default: ethereum)")
	createJobsCommand.Flags().StringVar(&address, "address", "", "The address to create jobs for")
	createJobsCommand.Flags().StringVar(&abiFile, "abi-file", "", "The path to the ABI file")
	createJobsCommand.Flags().StringVar(&customerId, "customer-id", "", "The customer ID to create jobs for (default: '')")
	createJobsCommand.Flags().StringVar(&userId, "user-id", "00000000-0000-0000-0000-000000000000", "The user ID to create jobs for (default: '00000000-0000-0000-0000-000000000000')")

	indexCommand.AddCommand(deploymentBlocksCommand)
	indexCommand.AddCommand(createJobsCommand)
	databaseCmd.AddCommand(indexCommand)

	return databaseCmd
}

func CreateHistoricalSyncCommand() *cobra.Command {

	var chain, baseDir, customerDbUriFlag string
	var addresses, customerIds []string
	var startBlock, endBlock, batchSize uint64
	var timeout, threads int
	var auto bool

	historicalSyncCmd := &cobra.Command{
		Use:   "historical-sync",
		Short: "Decode the historical data from various blockchains",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			indexerErr := indexer.CheckVariablesForIndexer()
			if indexerErr != nil {
				return indexerErr
			}

			storageErr := storage.CheckVariablesForStorage()
			if storageErr != nil {
				return storageErr
			}

			crawlerErr := crawler.CheckVariablesForCrawler()
			if crawlerErr != nil {
				return crawlerErr
			}

			syncErr := synchronizer.CheckVariablesForSynchronizer()
			if syncErr != nil {
				return syncErr
			}

			if chain == "" {
				return fmt.Errorf("blockchain is required via --chain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			indexer.InitDBConnection()

			newSynchronizer, synchonizerErr := synchronizer.NewSynchronizer(chain, baseDir, startBlock, endBlock, batchSize, timeout)
			if synchonizerErr != nil {
				return synchonizerErr
			}

			err := newSynchronizer.HistoricalSyncRef(customerDbUriFlag, addresses, customerIds, batchSize, auto, threads)

			if err != nil {
				return err
			}

			return nil
		},
	}

	historicalSyncCmd.Flags().StringVar(&chain, "chain", "ethereum", "The blockchain to crawl (default: ethereum)")
	historicalSyncCmd.Flags().StringVar(&baseDir, "base-dir", "", "The base directory to store the crawled data (default: '')")
	historicalSyncCmd.Flags().Uint64Var(&startBlock, "start-block", 0, "The block number to start decoding from (default: latest block)")
	historicalSyncCmd.Flags().Uint64Var(&endBlock, "end-block", 0, "The block number to end decoding at (default: latest block)")
	historicalSyncCmd.Flags().IntVar(&timeout, "timeout", 30, "The timeout for the crawler in seconds (default: 30)")
	historicalSyncCmd.Flags().Uint64Var(&batchSize, "batch-size", 100, "The number of blocks to crawl in each batch (default: 100)")
	historicalSyncCmd.Flags().StringVar(&customerDbUriFlag, "customer-db-uri", "", "Set customer database URI for development. This workflow bypass fetching customer IDs and its database URL connection strings from mdb-v3-controller API")
	historicalSyncCmd.Flags().StringSliceVar(&customerIds, "customer-ids", []string{}, "The list of customer IDs to sync")
	historicalSyncCmd.Flags().StringSliceVar(&addresses, "addresses", []string{}, "The list of addresses to sync")
	historicalSyncCmd.Flags().BoolVar(&auto, "auto", false, "Set this flag to sync all unfinished historical crawl from the database (default: false)")
	historicalSyncCmd.Flags().IntVar(&threads, "threads", 5, "Number of go-routines for concurrent crawling (default: 5)")

	return historicalSyncCmd
}

func CreateStarknetParseCommand() *cobra.Command {
	var infile string
	var rawABI []byte
	var readErr error

	starknetParseCommand := &cobra.Command{
		Use:   "parse",
		Short: "Parse a Starknet contract's ABI and return seer's interal representation of that ABI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if infile != "" {
				rawABI, readErr = os.ReadFile(infile)
			} else {
				rawABI, readErr = io.ReadAll(os.Stdin)
			}

			return readErr
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			parsedABI, parseErr := starknet.ParseABI(rawABI)
			if parseErr != nil {
				return parseErr
			}

			content, marshalErr := json.Marshal(parsedABI)
			if marshalErr != nil {
				return marshalErr
			}

			cmd.Println(string(content))
			return nil
		},
	}

	starknetParseCommand.Flags().StringVarP(&infile, "abi", "a", "", "Path to contract ABI (default stdin)")

	return starknetParseCommand
}

func CreateStarknetGenerateCommand() *cobra.Command {
	var infile, packageName string
	var rawABI []byte
	var readErr error

	starknetGenerateCommand := &cobra.Command{
		Use:   "generate",
		Short: "Generate Go bindings for a Starknet contract from its ABI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if infile != "" {
				rawABI, readErr = os.ReadFile(infile)
			} else {
				rawABI, readErr = io.ReadAll(os.Stdin)
			}

			return readErr
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			header, headerErr := starknet.GenerateHeader(packageName)
			if headerErr != nil {
				return headerErr
			}

			sections := []string{header}

			parsedABI, parseErr := starknet.ParseABI(rawABI)
			if parseErr != nil {
				return parseErr
			}

			code, codegenErr := starknet.Generate(parsedABI)
			if codegenErr != nil {
				return codegenErr
			}

			sections = append(sections, code)

			formattedCode, formattingErr := format.Source([]byte(strings.Join(sections, "\n\n")))
			if formattingErr != nil {
				return formattingErr
			}
			cmd.Println(string(formattedCode))
			return nil
		},
	}

	starknetGenerateCommand.Flags().StringVarP(&packageName, "package", "p", "", "The name of the package to generate")
	starknetGenerateCommand.Flags().StringVarP(&infile, "abi", "a", "", "Path to contract ABI (default stdin)")

	return starknetGenerateCommand
}

func CreateEVMCommand() *cobra.Command {
	evmCmd := &cobra.Command{
		Use:   "evm",
		Short: "Generate interfaces and crawlers for Ethereum Virtual Machine (EVM) contracts",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	evmGenerateCmd := CreateEVMGenerateCommand()
	evmCmd.AddCommand(evmGenerateCmd)

	return evmCmd
}

func CreateEVMGenerateCommand() *cobra.Command {
	var cli, noformat, includemain bool
	var infile, packageName, structName, bytecodefile, outfile, foundryBuildFile, hardhatBuildFile string
	var rawABI, bytecode []byte
	var readErr error
	var aliases map[string]string

	evmGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate Go bindings for an EVM contract from its ABI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if packageName == "" {
				return errors.New("package name is required via --package/-p")
			}
			if structName == "" {
				return errors.New("struct name is required via --struct/-s")
			}

			if foundryBuildFile != "" {
				var contents []byte
				contents, readErr = os.ReadFile(foundryBuildFile)
				if readErr != nil {
					return readErr
				}

				type foundryBytecodeObject struct {
					Object string `json:"object"`
				}

				type foundryBuildArtifact struct {
					ABI      json.RawMessage       `json:"abi"`
					Bytecode foundryBytecodeObject `json:"bytecode"`
				}

				var artifact foundryBuildArtifact
				readErr = json.Unmarshal(contents, &artifact)
				rawABI = []byte(artifact.ABI)
				bytecode = []byte(artifact.Bytecode.Object)
			} else if hardhatBuildFile != "" {
				var contents []byte
				contents, readErr = os.ReadFile(hardhatBuildFile)
				if readErr != nil {
					return readErr
				}

				type hardhatBuildArtifact struct {
					ABI      json.RawMessage `json:"abi"`
					Bytecode string          `json:"bytecode"`
				}

				var artifact hardhatBuildArtifact
				readErr = json.Unmarshal(contents, &artifact)
				rawABI = []byte(artifact.ABI)
				bytecode = []byte(artifact.Bytecode)
			} else if infile != "" {
				rawABI, readErr = os.ReadFile(infile)
			} else {
				rawABI, readErr = io.ReadAll(os.Stdin)
			}

			if bytecodefile != "" {
				bytecode, readErr = os.ReadFile(bytecodefile)
			}

			return readErr
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			code, codeErr := evm.GenerateTypes(structName, rawABI, bytecode, packageName, aliases)
			if codeErr != nil {
				return codeErr
			}

			header, headerErr := evm.GenerateHeader(packageName, cli, includemain, foundryBuildFile, infile, bytecodefile, structName, outfile, noformat)
			if headerErr != nil {
				return headerErr
			}

			code = header + code

			if cli {
				code, readErr = evm.AddCLI(code, structName, noformat, includemain)
				if readErr != nil {
					return readErr
				}
			}

			if outfile != "" {
				writeErr := os.WriteFile(outfile, []byte(code), 0644)
				if writeErr != nil {
					return writeErr
				}
			} else {
				cmd.Println(code)
			}
			return nil
		},
	}

	evmGenerateCmd.Flags().StringVarP(&packageName, "package", "p", "", "The name of the package to generate")
	evmGenerateCmd.Flags().StringVarP(&structName, "struct", "s", "", "The name of the struct to generate")
	evmGenerateCmd.Flags().StringVarP(&infile, "abi", "a", "", "Path to contract ABI (default stdin)")
	evmGenerateCmd.Flags().StringVarP(&bytecodefile, "bytecode", "b", "", "Path to contract bytecode (default none - in this case, no deployment method is created)")
	evmGenerateCmd.Flags().BoolVarP(&cli, "cli", "c", false, "Add a CLI for interacting with the contract (default false)")
	evmGenerateCmd.Flags().BoolVar(&noformat, "noformat", false, "Set this flag if you do not want the generated code to be formatted (useful to debug errors)")
	evmGenerateCmd.Flags().BoolVar(&includemain, "includemain", false, "Set this flag if you want to generate a \"main\" function to execute the CLI and make the generated code self-contained - this option is ignored if --cli is not set")
	evmGenerateCmd.Flags().StringVarP(&outfile, "output", "o", "", "Path to output file (default stdout)")
	evmGenerateCmd.Flags().StringVar(&foundryBuildFile, "foundry", "", "If your contract is compiled using Foundry, you can specify a path to the build file here (typically \"<foundry project root>/out/<solidity filename>/<contract name>.json\") instead of specifying --abi and --bytecode separately")
	evmGenerateCmd.Flags().StringVar(&hardhatBuildFile, "hardhat", "", "If your contract is compiled using Hardhat, you can specify a path to the build file here (typically \"<path to solidity file in hardhat artifact directory>/<contract name>.json\") instead of specifying --abi and --bytecode separately")
	evmGenerateCmd.Flags().StringToStringVar(&aliases, "alias", nil, "A map of identifier aliases (e.g. --alias name=somename)")

	return evmGenerateCmd
}
