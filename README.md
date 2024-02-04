# seer

`seer` is Moonstream's second generation EVM chain crawler.

It builds on what we have learned from our first generation of [`crawlers`](https://github.com/moonstream-to/api/tree/e69d81d1fb081cbddb0c8a1983af41e53d5a0f8f/crawlers).

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

```bash
./ownable-erc-721 -h
```
