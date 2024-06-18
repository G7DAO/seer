# seer

`seer` is Moonstream's second generation blockchain-adjacent tooling for crawling data and performing
smart contract interactions.

It builds on what we have learned from our first generation of [`crawlers`](https://github.com/moonstream-to/api/tree/e69d81d1fb081cbddb0c8a1983af41e53d5a0f8f/crawlers) and from [`moonworm`](https://github.com/moonstream-to/moonworm).

## Build

You can use `make` to build `seer`. From the root of this project, run:

```bash
make
```

If you would like to force a build from scratch:

```bash
make rebuild
```

To see a list of available `make` targets, please look at the [`Makefile`](./Makefile).

## Running seer

### Go bindings for Starknet contracts

To generate the Go bindings to a Starknet contract from its ABI, run:

```bash
seer starknet generate --abi $ABI_FILE --package $GO_PACKAGE_NAME
```

1. `$ABI_FILE` should be the path to the JSON file containing the Starknet contract ABI.
2. `$GO_PACKAGE_NAME` should be the name of the Go package that the generated code will belong to. If specified,
the line `package $GO_PACKAGE_NAME` will be emitted at the top of the generated code. If not specified, no
such line is emitted.

You can also pipe the ABI JSON into this command rather than specifying the `--abi` argument. For example:

```bash
jq . $ABI_FILE | seer starknet generate --package $GO_PACKAGE_NAME
```

### Go bindings for Ethereum Virtual Machine (EVM) contracts

To generate the Go bindings to an EVM contract, run:

```bash
seer evm generate \
    --abi $ABI_FILE \
    --bytecode $BIN_FILE \
    --cli \
    --package $GO_PACKAGE_NAME \
    --struct $GO_STRUCT_NAME
```

1. `$ABI_FILE` should be the path to a JSON file containing the contract's ABI.
2. `$BIN_FILE` should be a path to the file containing the compiled contract bytecode. If the `--bytecode` is not provided,
the bindings are generated with no deployment method.
3. `$GO_PACKAGE_NAME` should be the name of the Go package that the generated code will fall under.
4. `$GO_STRUCT_NAME` should be the name of the struct that you would like to represent an instance of the contract with the given ABI.

If you want to write the output to a file, you can use the `--output` argument to do so. Or shell redirections.

#### Example: `OwnableERC721`

The code in [`examples/ownable-erc-721`](./examples/ownable-erc-721/OwnableERC721.go) was generated from the project root directory using:

```bash
seer evm generate \
    --abi fixtures/OwnableERC721.json \
    --bytecode fixtures/OwnableERC721.bin \
    --cli \
    --package main \
    --struct OwnableERC721 \
    --output examples/ownable-erc-721/OwnableERC721.go
```

To run this code, first build it:

```bash
go build ./examples/ownable-erc-721
```

This will create an executable file called `ownable-erc-721` (on Windows, you may want to rename it to `ownable-erc-721.exe` for convenience).

Try running it:

```
$ ./ownable-erc-721 -h

Interact with the OwnableERC721 contract

Usage:
  ownable-erc-721 [flags]
  ownable-erc-721 [command]

Commands which deploy contracts
  deploy               Deploy a new OwnableERC721 contract

Commands which view contract state
  balance-of           Call the BalanceOf view method on a OwnableERC721 contract
  get-approved         Call the GetApproved view method on a OwnableERC721 contract
  is-approved-for-all  Call the IsApprovedForAll view method on a OwnableERC721 contract
  name                 Call the Name view method on a OwnableERC721 contract
  owner                Call the Owner view method on a OwnableERC721 contract
  owner-of             Call the OwnerOf view method on a OwnableERC721 contract
  supports-interface   Call the SupportsInterface view method on a OwnableERC721 contract
  symbol               Call the Symbol view method on a OwnableERC721 contract
  token-uri            Call the TokenURI view method on a OwnableERC721 contract

Commands which submit transactions
  approve              Execute the Approve method on a OwnableERC721 contract
  mint                 Execute the Mint method on a OwnableERC721 contract
  renounce-ownership   Execute the RenounceOwnership method on a OwnableERC721 contract
  safe-transfer-from   Execute the SafeTransferFrom method on a OwnableERC721 contract
  safe-transfer-from-0 Execute the SafeTransferFrom0 method on a OwnableERC721 contract
  set-approval-for-all Execute the SetApprovalForAll method on a OwnableERC721 contract
  transfer-from        Execute the TransferFrom method on a OwnableERC721 contract
  transfer-ownership   Execute the TransferOwnership method on a OwnableERC721 contract

Additional Commands:
  completion           Generate the autocompletion script for the specified shell
  help                 Help about any command

Flags:
  -h, --help   help for ownable-erc-721

Use "ownable-erc-721 [command] --help" for more information about a command.
```

# Crawler

That part of seer responsible for crawling raw blocks,tx_calls and events from the blockchain.

List of supported blockchains:

- arbitrum_one
- arbitrum_sepolia
- ethereum
- game7_orbit_arbitrum_sepolia
- mantle
- mantle_sepolia
- polygon
- xai
- xai_sepolia

## Build

You can use `make` to build `crawler`. From the root of this project, run:

```bash
make
```

Or build with go tools:

```bash
go build -o seer .
```

Or use dev script:

```bash
./dev.sh --help
```

Set environment variables:

```bash
export MOONSTREAM_DB_V3_INDEXES_URI="driver://user:pass@localhost/dbname"
```

## Generate crawler {chain} interface

note: You need add the chain endpoint it will fetch the data from endpoints.

Blockchain structure:

```bash
├── blockchain
│   ├── polygon
│   │   ├── blocks_transactions_polygon.proto
│   │   ├── blocks_transactions_polygon.pb.go
│   │   ├── polygon.go
```

## Regenerate proto interface

1. Generate code with proto compiler, docs: https://protobuf.dev/reference/go/go-generated/

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    blockchain/ethereum/ethereum_index_types.proto
```

2. Rename generated file similar to chain package.
3. Generate interface with seer, if chain is L2, specify flag `--side-chain`:

```bash
./seer blockchain generate -n ethereum
```

Or use bash script to do all job. But be careful it by default generates interfaces for L1 chains with additional fields, for side chains this script requires modification:

```bash
./prepare_blockchains.sh
```

## Run crawler

Before running the crawler, you need initialize the database with the following command:

```bash
./seer crawler --chain polygon --start-block 53922484 --force
```

## Inspect database

Find first and last blocks indexed in database and verify it's batch at storage:

```bash
./seer inspector db --chain polygon --storage-verify
```
