package common

// BlockchainHandler defines the interface for handling blockchain data fetching.
type BlockchainHandler interface {
    FetchData(blockNumber *big.Int) ([]EventData, error)
}

// EventData outlines the structure for blockchain event data.
type EventData struct {
    // Common fields...
    BlockNumber uint64
    Data        []byte
    ChainID     string
    Extra       interface{} // Use for blockchain-specific additional data
}

// Additional common handler functionalities...



// let's say we have a blockchain package that contains the following files:
//
// blockchain/
// ├── common/
// │   ├── decoding.go
// │   └── encoding.go
// ├── crawler/
// │   └── crawler.go
// ├── eth/
// │   ├── decoding.go
// │   └── encoding.go
// ├── starknet/
// │   ├── abi.go
// │   ├── compression.go
// │   └── starknet.go
// ├── tezos/
// │   ├── decoding.go
// │   └── encoding.go
// └── zilliqa/
//     ├── decoding.go
//     └── encoding.go
//
// We can use the following command to generate the documentation:
//	godoc2md github.com/NethermindEth/blockchain-crawler > doc.md
//	godoc2md github.com/NethermindEth/blockchain-crawler/common > common.md
//	godoc2md github.com/NethermindEth/blockchain-crawler/crawler > crawler.md
//	godoc2md github.com/NethermindEth/blockchain-crawler/eth > eth.md
//	godoc2md github.com/NethermindEth/blockchain-crawler/starknet > starknet.md
//	godoc2md github.com/NethermindEth/blockchain-crawler/tezos > tezos.md
//	godoc2md github.com/NethermindEth/blockchain-crawler/zilliqa > zilliqa.md
/