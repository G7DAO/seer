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

### Usage run crawler


Before running the crawler, you need initialize the database with the following command:




```bash
./alembic.sh -c configs/alembic.dev.ini ensure_version
```

```bash
./alembic.sh -c configs/alembic.dev.ini upgrade head
```

```bash
go run *.go crawler --chain polygon --start-block 53922484 --force --provider-uri
```



