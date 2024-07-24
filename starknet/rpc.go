package starknet

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
)

func NewProvider(web3ProviderUri string) (*rpc.Provider, error) {
	provider, providerErr := rpc.NewProvider(web3ProviderUri)
	if providerErr != nil {
		return nil, providerErr
	}

	return provider, nil
}

func ParseAddress(contractAddress string) (*felt.Felt, error) {
	fieldAdditiveIdentity := fp.NewElement(0)
	if contractAddress[:2] == "0x" {
		contractAddress = contractAddress[2:]
	}
	decodedAddress, decodeErr := hex.DecodeString(contractAddress)
	if decodeErr != nil {
		return nil, decodeErr
	}
	address := felt.NewFelt(&fieldAdditiveIdentity)
	address.SetBytes(decodedAddress)

	return address, nil
}

func ContractExistsAtBlock(ctx context.Context, provider *rpc.Provider, address *felt.Felt, blockNumber uint64) (bool, error) {
	_, err := provider.ClassHashAt(ctx, rpc.BlockID{Number: &blockNumber}, address)
	if err != nil {
		// Note: No other comparison (e.g. using errors.Is) is working.
		if err.Error() == rpc.ErrContractNotFound.Error() {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Perform a binary search to determine the block number at which the contract at the given address
// was deployed.
// Since the starknet_getCode method has been deprecated, this uses starknet_getClassHashAt in order
// to conduct the search. If the contract has not been deployed at a given block, calling
// starknet_getClassHashAt at that block will result in an error with code 20.
func DeploymentBlock(ctx context.Context, provider *rpc.Provider, address *felt.Felt) (uint64, error) {
	maxBlock, blockNumberErr := provider.BlockNumber(ctx)
	if blockNumberErr != nil {
		return 0, blockNumberErr
	}

	var minBlock uint64 = 0

	midBlock := (minBlock + maxBlock) / 2

	var isDeployed map[uint64]bool = make(map[uint64]bool)

	isDeployedAtBlock, blockErr := ContractExistsAtBlock(ctx, provider, address, maxBlock)
	if blockErr != nil {
		return 0, blockErr
	}
	if !isDeployedAtBlock {
		return 0, errors.New("address is not a contract")
	}
	isDeployed[maxBlock] = isDeployedAtBlock

	isDeployed[minBlock], blockErr = ContractExistsAtBlock(ctx, provider, address, minBlock)
	if blockErr != nil {
		return 0, blockErr
	}

	isDeployed[midBlock], blockErr = ContractExistsAtBlock(ctx, provider, address, midBlock)
	if blockErr != nil {
		return 0, blockErr
	}

	for (maxBlock - minBlock) >= 2 {
		if !isDeployed[minBlock] && !isDeployed[midBlock] {
			minBlock = midBlock
		} else {
			maxBlock = midBlock
		}

		midBlock = (minBlock + maxBlock) / 2

		isDeployed[midBlock], blockErr = ContractExistsAtBlock(ctx, provider, address, midBlock)
		if blockErr != nil {
			return 0, blockErr
		}
	}

	if isDeployed[minBlock] {
		return minBlock, nil
	}
	return maxBlock, nil
}
