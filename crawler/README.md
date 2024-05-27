### Crawler

That part of seer responsible for crawling raw blocks,tx_calls and events from the blockchain.

#### Build

You can use `make` to build `crawler`. From the root of this project, run:

```bash
make
```


Crawler variables

```env
export MOONSTREAM_INDEX_URI="driver://user:pass@localhost/dbname"
export INFURA_KEY="infura_key"
```


### Generate crawler {chain} interface

note: You need add the chain endpoint it will fetch the data from endpoints.

```bash
go run *.go crawler generate --chain polygon
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



### Regenerate proto interface
    
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    blocks_transactions_<chain>.proto
```



### Usage run crawler


Before running the crawler, you need initialize the database with the following command:


```bash
go run *.go crawler --chain polygon --start-block 53922484 --force --provider-uri
```



###


