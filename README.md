# seer

`seer` is Moonstream's second generation EVM chain crawler.

It builds on what we have learned from our first generation of [`crawlers`](https://github.com/moonstream-to/api/tree/e69d81d1fb081cbddb0c8a1983af41e53d5a0f8f/crawlers).

`seer` is written in Go, and is designed to run in several modes:

- [ ] JSONRPC mode, in which it reads data from HTTP APIs confirming to the [Ethereum JSONRPC API specification](https://ethereum.org/en/developers/docs/apis/json-rpc/).
- [ ] `geth` filesystem mode, in which it reads data from a `geth` node's data directory.

`seer` is also written to support several storage backends:

- [ ] Raw files on a filesystem
- [ ] PostgreSQL
- [ ] S3-compliant blob storage

`seer` supports crawling of:

- [ ] Events
- [ ] Transactions
- [ ] Internal messages inside transactions
