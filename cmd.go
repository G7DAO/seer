package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/moonstream-to/seer/starknet"
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
	starknetCmd := CreateStarknetCommand()
	rootCmd.AddCommand(completionCmd, versionCmd, starknetCmd)

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

func CreateStarknetCommand() *cobra.Command {
	starknetCmd := &cobra.Command{
		Use:   "starknet",
		Short: "Generate interfaces and crawlers for Starknet",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	starknetABIParseCmd := CreateStarknetABIParseCommand()
	starknetABIGenGoCmd := CreateStarknetABIGenGoCommand()
	starknetCmd.AddCommand(starknetABIParseCmd, starknetABIGenGoCmd)

	return starknetCmd
}

func CreateStarknetABIParseCommand() *cobra.Command {
	var infile string
	var rawABI []byte
	var readErr error

	starknetABIParseCommand := &cobra.Command{
		Use:   "abiparse",
		Short: "Parse a Starknet contract's ABI",
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

	starknetABIParseCommand.Flags().StringVarP(&infile, "abi", "a", "", "Path to contract ABI (default stdin)")

	return starknetABIParseCommand
}

func CreateStarknetABIGenGoCommand() *cobra.Command {
	var packageName string

	starknetABIGenGoCommand := &cobra.Command{
		Use:   "abigengo",
		Short: "Generate Go bindings for a Starknet contract from its ABI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// If we need to generalize this validator for other subcommands, use `cmd.Flag("package").Value.String()`
			// Or make a higher order function which takes the packageName as an argument and produces a validator as an output.

			if packageName == "" {
				return errors.New("you must provide a package name using --package/-p")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	starknetABIGenGoCommand.Flags().StringVarP(&packageName, "package", "p", "", "The name of the package to generate")

	return starknetABIGenGoCommand
}
