### Crawler

That part of seer responsible for crawling raw blocks,tx_calls and events from the blockchain.

#### Build

You can use `make` to build `crawler`. From the root of this project, run:

```bash
make
```


Crawler variables

```env
export MOONSTREAM_INDEX_URI="sqlite://~/seer/moonstream.db"
export INFURA_KEY="infura_key"
```

Depending of MOONSTREAM_INDEX_URI, you can use different databases.



### Usage run crawler


Before running the crawler, you need initialize the database with the following command:

```bash
go run *.go index initialize
```

```bash
go run *.go crawler --chain polygon --start-block 53922484 --force --provider-uri
```



