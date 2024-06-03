# Crawler

That part of seer responsible for crawling raw blocks,tx_calls and events from the blockchain.

## Build

You can use `make` to build `crawler`. From the root of this project, run:

```bash
make
```

Or build with go tools:

```bash
go build -o seer .
```

Set environment variables:

```bash
export MOONSTREAM_DB_V3_INDEXES_URI="driver://user:pass@localhost/dbname"
```

## Generate crawler {chain} interface

note: You need add the chain endpoint it will fetch the data from endpoints.

```bash
./seer crawler generate --chain polygon
```

Will generate the following files:

```bash
├── blockchain
│   ├── polygon
│   │   ├── blocks_transactions_polygon.proto
│   │   ├── blocks_transactions_polygon.pb.go
│   │   ├── polygon.go
│   │   ├── types.go
```

## Regenerate proto interface

1. Go to chain directory, for example:

```bash
cd blockchain/ethereum
```

2. Generate code with proto compiler, docs: https://protobuf.dev/reference/go/go-generated/

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    ethereum_index_types.proto
```

3. Rename generated file similar to chain package.
4. Generate interface with seer, if chain is L2, specify flag `--side-chain`:

```bash
cd ../..
./seer blockchain generate -n ethereum
```

## Run crawler

Before running the crawler, you need initialize the database with the following command:

```bash
./seer crawler --chain polygon --start-block 53922484 --force
```
