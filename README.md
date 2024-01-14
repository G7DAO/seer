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

1. `$ABI_FILE` should be the path to the JSON file containing the Starknet contract ABI
2. `$GO_PACKAGE_NAME` should be the name of the Go package that the generated code will belong to. If specified,
the line `package $GO_PACKAGE_NAME` will be emitted at the top of the generated code. If not specified, no
such line is emitted.

You can also pipe the ABI JSON into this command rather than specifying the `--abi` argument. For example:

```bash
jq . $ABI_FILE | seer starknet generate --package $GO_PACKAGE_NAME
```
