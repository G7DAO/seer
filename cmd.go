package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

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
	indexCmd := CreateIndexCommand()
	evmCmd := CreateEVMCommand()
	synchronizerCmd := CreateSynchronizerCommand()
	rootCmd.AddCommand(completionCmd, versionCmd, blockchainCmd, starknetCmd, evmCmd, crawlerCmd, indexCmd, synchronizerCmd)

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
}

func CreateBlockchainGenerateCommand() *cobra.Command {
	var blockchainNameLower string

	blockchainGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate methods and types for different blockchains from template",
		RunE: func(cmd *cobra.Command, args []string) error {
			blockchainNameFilePath := fmt.Sprintf("blockchain/%s/%s.go", blockchainNameLower, blockchainNameLower)

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
			mkdirErr := os.Mkdir(blockchainNameLower, 0775)
			if mkdirErr != nil {
				return mkdirErr
			}
			outputFile, createErr := os.Create(blockchainNameFilePath)
			if createErr != nil {
				return createErr
			}
			defer outputFile.Close()

			// Execute template and write to output file
			data := BlockchainTemplateData{BlockchainName: blockchainName, BlockchainNameLower: blockchainNameLower}
			execErr := tmpl.Execute(outputFile, data)
			if execErr != nil {
				return execErr
			}

			log.Printf("Blockchain file generated successfully: %s", blockchainNameFilePath)

			return nil
		},
	}

	blockchainGenerateCmd.Flags().StringVarP(&blockchainNameLower, "name", "n", "", "The name of the blockchain to generate lowercase (example: 'arbitrum_one')")

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
	var startBlock, endBlock uint64
	var batchSize, confirmations, timeout int
	var chain, baseDir string
	var force bool

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

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			indexer.InitDBConnection()

			// read the blockchain url from $INFURA_URL
			// if it is not set, use the default url
			crawler := crawler.NewCrawler(chain, startBlock, endBlock, timeout, batchSize, confirmations, baseDir, force)

			crawler.Start()

		},
	}

	crawlerCmd.Flags().StringVar(&chain, "chain", "ethereum", "The blockchain to crawl (default: ethereum)")
	crawlerCmd.Flags().Uint64Var(&startBlock, "start-block", 0, "The block number to start crawling from (default: latest block)")
	crawlerCmd.Flags().Uint64Var(&endBlock, "end-block", 0, "The block number to end crawling at (default: latest block)")
	crawlerCmd.Flags().IntVar(&timeout, "timeout", 0, "The timeout for the crawler in seconds (default: 0 - no timeout)")
	crawlerCmd.Flags().IntVar(&batchSize, "batch-size", 10, "The number of blocks to crawl in each batch (default: 10)")
	crawlerCmd.Flags().IntVar(&confirmations, "confirmations", 10, "The number of confirmations to consider for block finality (default: 10)")
	crawlerCmd.Flags().StringVar(&baseDir, "base-dir", "data", "The base directory to store the crawled data (default: data)")
	crawlerCmd.Flags().BoolVar(&force, "force", false, "Set this flag to force the crawler start from the start block (default: false)")

	return crawlerCmd
}

func CreateSynchronizerCommand() *cobra.Command {
	var startBlock, endBlock uint64
	// var startBlockBig, endBlockBig big.Int
	var baseDir, output, abi_source string

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

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			indexer.InitDBConnection()

			// convert the start and end block to big.Int

			// read the blockchain url from $INFURA_URL
			// if it is not set, use the default url
			synchronizer := synchronizer.NewSynchronizer("polygon", startBlock, endBlock)

			synchronizer.SyncCustomers()
		},
	}

	synchronizerCmd.Flags().Uint64Var(&startBlock, "start-block", 0, "The block number to start decoding from (default: latest block)")
	synchronizerCmd.Flags().Uint64Var(&endBlock, "end-block", 0, "The block number to end decoding at (default: latest block)")
	synchronizerCmd.Flags().StringVar(&baseDir, "base-dir", "data", "The base directory to store the crawled data (default: data)")
	synchronizerCmd.Flags().StringVar(&output, "output", "output", "The output directory to store the decoded data (default: output)")
	synchronizerCmd.Flags().StringVar(&abi_source, "abi-source", "abi", "The source of the ABI (default: abi)")

	return synchronizerCmd
}

func CreateIndexCommand() *cobra.Command {

	indexCommand := &cobra.Command{
		Use:   "index",
		Short: "Index storage of moonstream blockstore",
	}

	// subcommands

	initializeCommand := &cobra.Command{
		Use:   "initialize",
		Short: "Initialize the index storage",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing index storage from go")
		},
	}

	readCommand := &cobra.Command{
		Use:   "read",
		Short: "Read the index storage",
		Run: func(cmd *cobra.Command, args []string) {
			// index.Read()
		},
	}

	indexCommand.AddCommand(initializeCommand, readCommand)

	return indexCommand
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
