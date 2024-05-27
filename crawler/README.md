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
export WEB3_JSON_RPC_URI="http://localhost:8545"
```

## Generate crawler {chain} interface

note: You need add the chain endpoint it will fetch the data from endpoints.

```bash
./seer crawler generate --web3 $WEB3_JSON_RPC_URI
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

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    blocks_transactions_<chain>.proto
```

## Run crawler

Before running the crawler, you need initialize the database with the following command:

```bash
./seer crawler --web3 $WEB3_JSON_RPC_URI --start-block 53922484 --force
```
