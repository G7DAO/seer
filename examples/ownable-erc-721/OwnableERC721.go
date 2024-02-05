// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"errors"
	"math/big"
	"strings"

	"context"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	// Reference imports to suppress errors if they are not otherwise used.
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// OwnableERC721MetaData contains all meta data concerning the OwnableERC721 contract.
var OwnableERC721MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620019fe380380620019fe833981016040819052620000349162000332565b8251839083906200004d906000906020850190620001bf565b50805162000063906001906020840190620001bf565b505050620000806200007a6200009460201b60201c565b62000098565b6200008b81620000ea565b505050620003fc565b3390565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6006546001600160a01b031633146200014a5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b6001600160a01b038116620001b15760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840162000141565b620001bc8162000098565b50565b828054620001cd90620003bf565b90600052602060002090601f016020900481019282620001f157600085556200023c565b82601f106200020c57805160ff19168380011785556200023c565b828001600101855582156200023c579182015b828111156200023c5782518255916020019190600101906200021f565b506200024a9291506200024e565b5090565b5b808211156200024a57600081556001016200024f565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200028d57600080fd5b81516001600160401b0380821115620002aa57620002aa62000265565b604051601f8301601f19908116603f01168101908282118183101715620002d557620002d562000265565b81604052838152602092508683858801011115620002f257600080fd5b600091505b83821015620003165785820183015181830184015290820190620002f7565b83821115620003285760008385830101525b9695505050505050565b6000806000606084860312156200034857600080fd5b83516001600160401b03808211156200036057600080fd5b6200036e878388016200027b565b945060208601519150808211156200038557600080fd5b5062000394868287016200027b565b604086015190935090506001600160a01b0381168114620003b457600080fd5b809150509250925092565b600181811c90821680620003d457607f821691505b60208210811415620003f657634e487b7160e01b600052602260045260246000fd5b50919050565b6115f2806200040c6000396000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c806370a08231116100a2578063a22cb46511610071578063a22cb4651461021b578063b88d4fde1461022e578063c87b56dd14610241578063e985e9c514610254578063f2fde38b1461029057600080fd5b806370a08231146101d9578063715018a6146101fa5780638da5cb5b1461020257806395d89b411461021357600080fd5b806323b872dd116100de57806323b872dd1461018d57806340c10f19146101a057806342842e0e146101b35780636352211e146101c657600080fd5b806301ffc9a71461011057806306fdde0314610138578063081812fc1461014d578063095ea7b314610178575b600080fd5b61012361011e3660046110cd565b6102a3565b60405190151581526020015b60405180910390f35b6101406102f5565b60405161012f9190611142565b61016061015b366004611155565b610387565b6040516001600160a01b03909116815260200161012f565b61018b61018636600461118a565b610421565b005b61018b61019b3660046111b4565b610537565b61018b6101ae36600461118a565b610568565b61018b6101c13660046111b4565b6105a0565b6101606101d4366004611155565b6105bb565b6101ec6101e73660046111f0565b610632565b60405190815260200161012f565b61018b6106b9565b6006546001600160a01b0316610160565b6101406106ef565b61018b61022936600461120b565b6106fe565b61018b61023c36600461125d565b610709565b61014061024f366004611155565b610741565b610123610262366004611339565b6001600160a01b03918216600090815260056020908152604080832093909416825291909152205460ff1690565b61018b61029e3660046111f0565b610829565b60006001600160e01b031982166380ac58cd60e01b14806102d457506001600160e01b03198216635b5e139f60e01b145b806102ef57506301ffc9a760e01b6001600160e01b03198316145b92915050565b6060600080546103049061136c565b80601f01602080910402602001604051908101604052809291908181526020018280546103309061136c565b801561037d5780601f106103525761010080835404028352916020019161037d565b820191906000526020600020905b81548152906001019060200180831161036057829003601f168201915b5050505050905090565b6000818152600260205260408120546001600160a01b03166104055760405162461bcd60e51b815260206004820152602c60248201527f4552433732313a20617070726f76656420717565727920666f72206e6f6e657860448201526b34b9ba32b73a103a37b5b2b760a11b60648201526084015b60405180910390fd5b506000908152600460205260409020546001600160a01b031690565b600061042c826105bb565b9050806001600160a01b0316836001600160a01b0316141561049a5760405162461bcd60e51b815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e656044820152603960f91b60648201526084016103fc565b336001600160a01b03821614806104b657506104b68133610262565b6105285760405162461bcd60e51b815260206004820152603860248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f74206f7760448201527f6e6572206e6f7220617070726f76656420666f7220616c6c000000000000000060648201526084016103fc565b61053283836108c4565b505050565b6105413382610932565b61055d5760405162461bcd60e51b81526004016103fc906113a7565b610532838383610a29565b6006546001600160a01b031633146105925760405162461bcd60e51b81526004016103fc906113f8565b61059c8282610bc9565b5050565b61053283838360405180602001604052806000815250610709565b6000818152600260205260408120546001600160a01b0316806102ef5760405162461bcd60e51b815260206004820152602960248201527f4552433732313a206f776e657220717565727920666f72206e6f6e657869737460448201526832b73a103a37b5b2b760b91b60648201526084016103fc565b60006001600160a01b03821661069d5760405162461bcd60e51b815260206004820152602a60248201527f4552433732313a2062616c616e636520717565727920666f7220746865207a65604482015269726f206164647265737360b01b60648201526084016103fc565b506001600160a01b031660009081526003602052604090205490565b6006546001600160a01b031633146106e35760405162461bcd60e51b81526004016103fc906113f8565b6106ed6000610be3565b565b6060600180546103049061136c565b61059c338383610c35565b6107133383610932565b61072f5760405162461bcd60e51b81526004016103fc906113a7565b61073b84848484610d04565b50505050565b6000818152600260205260409020546060906001600160a01b03166107c05760405162461bcd60e51b815260206004820152602f60248201527f4552433732314d657461646174613a2055524920717565727920666f72206e6f60448201526e3732bc34b9ba32b73a103a37b5b2b760891b60648201526084016103fc565b60006107d760408051602081019091526000815290565b905060008151116107f75760405180602001604052806000815250610822565b8061080184610d37565b60405160200161081292919061142d565b6040516020818303038152906040525b9392505050565b6006546001600160a01b031633146108535760405162461bcd60e51b81526004016103fc906113f8565b6001600160a01b0381166108b85760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016103fc565b6108c181610be3565b50565b600081815260046020526040902080546001600160a01b0319166001600160a01b03841690811790915581906108f9826105bb565b6001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b6000818152600260205260408120546001600160a01b03166109ab5760405162461bcd60e51b815260206004820152602c60248201527f4552433732313a206f70657261746f7220717565727920666f72206e6f6e657860448201526b34b9ba32b73a103a37b5b2b760a11b60648201526084016103fc565b60006109b6836105bb565b9050806001600160a01b0316846001600160a01b031614806109f15750836001600160a01b03166109e684610387565b6001600160a01b0316145b80610a2157506001600160a01b0380821660009081526005602090815260408083209388168352929052205460ff165b949350505050565b826001600160a01b0316610a3c826105bb565b6001600160a01b031614610aa45760405162461bcd60e51b815260206004820152602960248201527f4552433732313a207472616e73666572206f6620746f6b656e2074686174206960448201526839903737ba1037bbb760b91b60648201526084016103fc565b6001600160a01b038216610b065760405162461bcd60e51b8152602060048201526024808201527f4552433732313a207472616e7366657220746f20746865207a65726f206164646044820152637265737360e01b60648201526084016103fc565b610b116000826108c4565b6001600160a01b0383166000908152600360205260408120805460019290610b3a908490611472565b90915550506001600160a01b0382166000908152600360205260408120805460019290610b68908490611489565b909155505060008181526002602052604080822080546001600160a01b0319166001600160a01b0386811691821790925591518493918716917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4505050565b61059c828260405180602001604052806000815250610e35565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b816001600160a01b0316836001600160a01b03161415610c975760405162461bcd60e51b815260206004820152601960248201527f4552433732313a20617070726f766520746f2063616c6c65720000000000000060448201526064016103fc565b6001600160a01b03838116600081815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b610d0f848484610a29565b610d1b84848484610e68565b61073b5760405162461bcd60e51b81526004016103fc906114a1565b606081610d5b5750506040805180820190915260018152600360fc1b602082015290565b8160005b8115610d855780610d6f816114f3565b9150610d7e9050600a83611524565b9150610d5f565b60008167ffffffffffffffff811115610da057610da0611247565b6040519080825280601f01601f191660200182016040528015610dca576020820181803683370190505b5090505b8415610a2157610ddf600183611472565b9150610dec600a86611538565b610df7906030611489565b60f81b818381518110610e0c57610e0c61154c565b60200101906001600160f81b031916908160001a905350610e2e600a86611524565b9450610dce565b610e3f8383610f75565b610e4c6000848484610e68565b6105325760405162461bcd60e51b81526004016103fc906114a1565b60006001600160a01b0384163b15610f6a57604051630a85bd0160e11b81526001600160a01b0385169063150b7a0290610eac903390899088908890600401611562565b602060405180830381600087803b158015610ec657600080fd5b505af1925050508015610ef6575060408051601f3d908101601f19168201909252610ef39181019061159f565b60015b610f50573d808015610f24576040519150601f19603f3d011682016040523d82523d6000602084013e610f29565b606091505b508051610f485760405162461bcd60e51b81526004016103fc906114a1565b805181602001fd5b6001600160e01b031916630a85bd0160e11b149050610a21565b506001949350505050565b6001600160a01b038216610fcb5760405162461bcd60e51b815260206004820181905260248201527f4552433732313a206d696e7420746f20746865207a65726f206164647265737360448201526064016103fc565b6000818152600260205260409020546001600160a01b0316156110305760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e7465640000000060448201526064016103fc565b6001600160a01b0382166000908152600360205260408120805460019290611059908490611489565b909155505060008181526002602052604080822080546001600160a01b0319166001600160a01b03861690811790915590518392907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a45050565b6001600160e01b0319811681146108c157600080fd5b6000602082840312156110df57600080fd5b8135610822816110b7565b60005b838110156111055781810151838201526020016110ed565b8381111561073b5750506000910152565b6000815180845261112e8160208601602086016110ea565b601f01601f19169290920160200192915050565b6020815260006108226020830184611116565b60006020828403121561116757600080fd5b5035919050565b80356001600160a01b038116811461118557600080fd5b919050565b6000806040838503121561119d57600080fd5b6111a68361116e565b946020939093013593505050565b6000806000606084860312156111c957600080fd5b6111d28461116e565b92506111e06020850161116e565b9150604084013590509250925092565b60006020828403121561120257600080fd5b6108228261116e565b6000806040838503121561121e57600080fd5b6112278361116e565b91506020830135801515811461123c57600080fd5b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b6000806000806080858703121561127357600080fd5b61127c8561116e565b935061128a6020860161116e565b925060408501359150606085013567ffffffffffffffff808211156112ae57600080fd5b818701915087601f8301126112c257600080fd5b8135818111156112d4576112d4611247565b604051601f8201601f19908116603f011681019083821181831017156112fc576112fc611247565b816040528281528a602084870101111561131557600080fd5b82602086016020830137600060208483010152809550505050505092959194509250565b6000806040838503121561134c57600080fd5b6113558361116e565b91506113636020840161116e565b90509250929050565b600181811c9082168061138057607f821691505b602082108114156113a157634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526031908201527f4552433732313a207472616e736665722063616c6c6572206973206e6f74206f6040820152701ddb995c881b9bdc88185c1c1c9bdd9959607a1b606082015260800190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6000835161143f8184602088016110ea565b8351908301906114538183602088016110ea565b01949350505050565b634e487b7160e01b600052601160045260246000fd5b6000828210156114845761148461145c565b500390565b6000821982111561149c5761149c61145c565b500190565b60208082526032908201527f4552433732313a207472616e7366657220746f206e6f6e20455243373231526560408201527131b2b4bb32b91034b6b83632b6b2b73a32b960711b606082015260800190565b60006000198214156115075761150761145c565b5060010190565b634e487b7160e01b600052601260045260246000fd5b6000826115335761153361150e565b500490565b6000826115475761154761150e565b500690565b634e487b7160e01b600052603260045260246000fd5b6001600160a01b038581168252841660208201526040810183905260806060820181905260009061159590830184611116565b9695505050505050565b6000602082840312156115b157600080fd5b8151610822816110b756fea264697066735822122061bd7069fbe4cf9497a08868a69ed81fec10d620962d72d34b6772153fc0186364736f6c63430008090033",
}

// OwnableERC721ABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableERC721MetaData.ABI instead.
var OwnableERC721ABI = OwnableERC721MetaData.ABI

// OwnableERC721Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OwnableERC721MetaData.Bin instead.
var OwnableERC721Bin = OwnableERC721MetaData.Bin

// DeployOwnableERC721 deploys a new Ethereum contract, binding an instance of OwnableERC721 to it.
func DeployOwnableERC721(auth *bind.TransactOpts, backend bind.ContractBackend, name_ string, symbol_ string, owner common.Address) (common.Address, *types.Transaction, *OwnableERC721, error) {
	parsed, err := OwnableERC721MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OwnableERC721Bin), backend, name_, symbol_, owner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OwnableERC721{OwnableERC721Caller: OwnableERC721Caller{contract: contract}, OwnableERC721Transactor: OwnableERC721Transactor{contract: contract}, OwnableERC721Filterer: OwnableERC721Filterer{contract: contract}}, nil
}

// OwnableERC721 is an auto generated Go binding around an Ethereum contract.
type OwnableERC721 struct {
	OwnableERC721Caller     // Read-only binding to the contract
	OwnableERC721Transactor // Write-only binding to the contract
	OwnableERC721Filterer   // Log filterer for contract events
}

// OwnableERC721Caller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableERC721Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableERC721Transactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableERC721Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableERC721Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableERC721Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableERC721Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableERC721Session struct {
	Contract     *OwnableERC721    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableERC721CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableERC721CallerSession struct {
	Contract *OwnableERC721Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// OwnableERC721TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableERC721TransactorSession struct {
	Contract     *OwnableERC721Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// OwnableERC721Raw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableERC721Raw struct {
	Contract *OwnableERC721 // Generic contract binding to access the raw methods on
}

// OwnableERC721CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableERC721CallerRaw struct {
	Contract *OwnableERC721Caller // Generic read-only contract binding to access the raw methods on
}

// OwnableERC721TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableERC721TransactorRaw struct {
	Contract *OwnableERC721Transactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnableERC721 creates a new instance of OwnableERC721, bound to a specific deployed contract.
func NewOwnableERC721(address common.Address, backend bind.ContractBackend) (*OwnableERC721, error) {
	contract, err := bindOwnableERC721(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnableERC721{OwnableERC721Caller: OwnableERC721Caller{contract: contract}, OwnableERC721Transactor: OwnableERC721Transactor{contract: contract}, OwnableERC721Filterer: OwnableERC721Filterer{contract: contract}}, nil
}

// NewOwnableERC721Caller creates a new read-only instance of OwnableERC721, bound to a specific deployed contract.
func NewOwnableERC721Caller(address common.Address, caller bind.ContractCaller) (*OwnableERC721Caller, error) {
	contract, err := bindOwnableERC721(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableERC721Caller{contract: contract}, nil
}

// NewOwnableERC721Transactor creates a new write-only instance of OwnableERC721, bound to a specific deployed contract.
func NewOwnableERC721Transactor(address common.Address, transactor bind.ContractTransactor) (*OwnableERC721Transactor, error) {
	contract, err := bindOwnableERC721(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableERC721Transactor{contract: contract}, nil
}

// NewOwnableERC721Filterer creates a new log filterer instance of OwnableERC721, bound to a specific deployed contract.
func NewOwnableERC721Filterer(address common.Address, filterer bind.ContractFilterer) (*OwnableERC721Filterer, error) {
	contract, err := bindOwnableERC721(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableERC721Filterer{contract: contract}, nil
}

// bindOwnableERC721 binds a generic wrapper to an already deployed contract.
func bindOwnableERC721(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OwnableERC721MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableERC721 *OwnableERC721Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableERC721.Contract.OwnableERC721Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableERC721 *OwnableERC721Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableERC721.Contract.OwnableERC721Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableERC721 *OwnableERC721Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableERC721.Contract.OwnableERC721Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableERC721 *OwnableERC721CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableERC721.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableERC721 *OwnableERC721TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableERC721.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableERC721 *OwnableERC721TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableERC721.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_OwnableERC721 *OwnableERC721Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OwnableERC721.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_OwnableERC721 *OwnableERC721Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _OwnableERC721.Contract.BalanceOf(&_OwnableERC721.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_OwnableERC721 *OwnableERC721CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _OwnableERC721.Contract.BalanceOf(&_OwnableERC721.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_OwnableERC721 *OwnableERC721Caller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _OwnableERC721.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_OwnableERC721 *OwnableERC721Session) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _OwnableERC721.Contract.GetApproved(&_OwnableERC721.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_OwnableERC721 *OwnableERC721CallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _OwnableERC721.Contract.GetApproved(&_OwnableERC721.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_OwnableERC721 *OwnableERC721Caller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _OwnableERC721.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_OwnableERC721 *OwnableERC721Session) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _OwnableERC721.Contract.IsApprovedForAll(&_OwnableERC721.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_OwnableERC721 *OwnableERC721CallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _OwnableERC721.Contract.IsApprovedForAll(&_OwnableERC721.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_OwnableERC721 *OwnableERC721Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OwnableERC721.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_OwnableERC721 *OwnableERC721Session) Name() (string, error) {
	return _OwnableERC721.Contract.Name(&_OwnableERC721.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_OwnableERC721 *OwnableERC721CallerSession) Name() (string, error) {
	return _OwnableERC721.Contract.Name(&_OwnableERC721.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnableERC721 *OwnableERC721Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OwnableERC721.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnableERC721 *OwnableERC721Session) Owner() (common.Address, error) {
	return _OwnableERC721.Contract.Owner(&_OwnableERC721.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnableERC721 *OwnableERC721CallerSession) Owner() (common.Address, error) {
	return _OwnableERC721.Contract.Owner(&_OwnableERC721.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_OwnableERC721 *OwnableERC721Caller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _OwnableERC721.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_OwnableERC721 *OwnableERC721Session) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _OwnableERC721.Contract.OwnerOf(&_OwnableERC721.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_OwnableERC721 *OwnableERC721CallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _OwnableERC721.Contract.OwnerOf(&_OwnableERC721.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_OwnableERC721 *OwnableERC721Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _OwnableERC721.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_OwnableERC721 *OwnableERC721Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _OwnableERC721.Contract.SupportsInterface(&_OwnableERC721.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_OwnableERC721 *OwnableERC721CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _OwnableERC721.Contract.SupportsInterface(&_OwnableERC721.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_OwnableERC721 *OwnableERC721Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OwnableERC721.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_OwnableERC721 *OwnableERC721Session) Symbol() (string, error) {
	return _OwnableERC721.Contract.Symbol(&_OwnableERC721.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_OwnableERC721 *OwnableERC721CallerSession) Symbol() (string, error) {
	return _OwnableERC721.Contract.Symbol(&_OwnableERC721.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_OwnableERC721 *OwnableERC721Caller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _OwnableERC721.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_OwnableERC721 *OwnableERC721Session) TokenURI(tokenId *big.Int) (string, error) {
	return _OwnableERC721.Contract.TokenURI(&_OwnableERC721.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_OwnableERC721 *OwnableERC721CallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _OwnableERC721.Contract.TokenURI(&_OwnableERC721.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721Transactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721Session) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.Contract.Approve(&_OwnableERC721.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721TransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.Contract.Approve(&_OwnableERC721.TransactOpts, to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721Transactor) Mint(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.contract.Transact(opts, "mint", to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721Session) Mint(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.Contract.Mint(&_OwnableERC721.TransactOpts, to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721TransactorSession) Mint(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.Contract.Mint(&_OwnableERC721.TransactOpts, to, tokenId)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnableERC721 *OwnableERC721Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableERC721.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnableERC721 *OwnableERC721Session) RenounceOwnership() (*types.Transaction, error) {
	return _OwnableERC721.Contract.RenounceOwnership(&_OwnableERC721.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnableERC721 *OwnableERC721TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _OwnableERC721.Contract.RenounceOwnership(&_OwnableERC721.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721Session) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.Contract.SafeTransferFrom(&_OwnableERC721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721TransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.Contract.SafeTransferFrom(&_OwnableERC721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_OwnableERC721 *OwnableERC721Transactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _OwnableERC721.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_OwnableERC721 *OwnableERC721Session) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _OwnableERC721.Contract.SafeTransferFrom0(&_OwnableERC721.TransactOpts, from, to, tokenId, _data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes _data) returns()
func (_OwnableERC721 *OwnableERC721TransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, _data []byte) (*types.Transaction, error) {
	return _OwnableERC721.Contract.SafeTransferFrom0(&_OwnableERC721.TransactOpts, from, to, tokenId, _data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_OwnableERC721 *OwnableERC721Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _OwnableERC721.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_OwnableERC721 *OwnableERC721Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _OwnableERC721.Contract.SetApprovalForAll(&_OwnableERC721.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_OwnableERC721 *OwnableERC721TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _OwnableERC721.Contract.SetApprovalForAll(&_OwnableERC721.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721Session) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.Contract.TransferFrom(&_OwnableERC721.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_OwnableERC721 *OwnableERC721TransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _OwnableERC721.Contract.TransferFrom(&_OwnableERC721.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnableERC721 *OwnableERC721Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _OwnableERC721.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnableERC721 *OwnableERC721Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OwnableERC721.Contract.TransferOwnership(&_OwnableERC721.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnableERC721 *OwnableERC721TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OwnableERC721.Contract.TransferOwnership(&_OwnableERC721.TransactOpts, newOwner)
}

// OwnableERC721ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the OwnableERC721 contract.
type OwnableERC721ApprovalIterator struct {
	Event *OwnableERC721Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnableERC721ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableERC721Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnableERC721Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnableERC721ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableERC721ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableERC721Approval represents a Approval event raised by the OwnableERC721 contract.
type OwnableERC721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_OwnableERC721 *OwnableERC721Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*OwnableERC721ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _OwnableERC721.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &OwnableERC721ApprovalIterator{contract: _OwnableERC721.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_OwnableERC721 *OwnableERC721Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *OwnableERC721Approval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _OwnableERC721.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableERC721Approval)
				if err := _OwnableERC721.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_OwnableERC721 *OwnableERC721Filterer) ParseApproval(log types.Log) (*OwnableERC721Approval, error) {
	event := new(OwnableERC721Approval)
	if err := _OwnableERC721.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnableERC721ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the OwnableERC721 contract.
type OwnableERC721ApprovalForAllIterator struct {
	Event *OwnableERC721ApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnableERC721ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableERC721ApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnableERC721ApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnableERC721ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableERC721ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableERC721ApprovalForAll represents a ApprovalForAll event raised by the OwnableERC721 contract.
type OwnableERC721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_OwnableERC721 *OwnableERC721Filterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*OwnableERC721ApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _OwnableERC721.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &OwnableERC721ApprovalForAllIterator{contract: _OwnableERC721.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_OwnableERC721 *OwnableERC721Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *OwnableERC721ApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _OwnableERC721.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableERC721ApprovalForAll)
				if err := _OwnableERC721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_OwnableERC721 *OwnableERC721Filterer) ParseApprovalForAll(log types.Log) (*OwnableERC721ApprovalForAll, error) {
	event := new(OwnableERC721ApprovalForAll)
	if err := _OwnableERC721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnableERC721OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OwnableERC721 contract.
type OwnableERC721OwnershipTransferredIterator struct {
	Event *OwnableERC721OwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnableERC721OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableERC721OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnableERC721OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnableERC721OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableERC721OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableERC721OwnershipTransferred represents a OwnershipTransferred event raised by the OwnableERC721 contract.
type OwnableERC721OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnableERC721 *OwnableERC721Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableERC721OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnableERC721.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableERC721OwnershipTransferredIterator{contract: _OwnableERC721.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnableERC721 *OwnableERC721Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableERC721OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnableERC721.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableERC721OwnershipTransferred)
				if err := _OwnableERC721.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnableERC721 *OwnableERC721Filterer) ParseOwnershipTransferred(log types.Log) (*OwnableERC721OwnershipTransferred, error) {
	event := new(OwnableERC721OwnershipTransferred)
	if err := _OwnableERC721.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnableERC721TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the OwnableERC721 contract.
type OwnableERC721TransferIterator struct {
	Event *OwnableERC721Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnableERC721TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableERC721Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnableERC721Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnableERC721TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableERC721TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableERC721Transfer represents a Transfer event raised by the OwnableERC721 contract.
type OwnableERC721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_OwnableERC721 *OwnableERC721Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*OwnableERC721TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _OwnableERC721.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &OwnableERC721TransferIterator{contract: _OwnableERC721.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_OwnableERC721 *OwnableERC721Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *OwnableERC721Transfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _OwnableERC721.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableERC721Transfer)
				if err := _OwnableERC721.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_OwnableERC721 *OwnableERC721Filterer) ParseTransfer(log types.Log) (*OwnableERC721Transfer, error) {
	event := new(OwnableERC721Transfer)
	if err := _OwnableERC721.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func CreateOwnableERC721DeploymentCommand() *cobra.Command {
	var keyfile, nonce, password, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, rpc string
	var gasLimit uint64
	var simulate bool
	var timeout uint

	var name string

	var symbol string

	var owner common.Address
	var ownerRaw string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a new OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keyfile == "" {
				return fmt.Errorf("--keystore not specified (this should be a path to an Ethereum account keystore file)")
			}

			if ownerRaw == "" {
				return fmt.Errorf("--owner argument not specified")
			} else if !common.IsHexAddress(ownerRaw) {
				return fmt.Errorf("--owner argument is not a valid Ethereum address")
			}
			owner = common.HexToAddress(ownerRaw)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			key, keyErr := KeyFromFile(keyfile, password)
			if keyErr != nil {
				return keyErr
			}

			chainIDCtx, cancelChainIDCtx := NewChainContext(timeout)
			defer cancelChainIDCtx()
			chainID, chainIDErr := client.ChainID(chainIDCtx)
			if chainIDErr != nil {
				return chainIDErr
			}

			transactionOpts, transactionOptsErr := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
			if transactionOptsErr != nil {
				return transactionOptsErr
			}

			SetTransactionParametersFromArgs(transactionOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, gasLimit, simulate)

			address, deploymentTransaction, _, deploymentErr := DeployOwnableERC721(
				transactionOpts,
				client,
				name,
				symbol,
				owner,
			)
			if deploymentErr != nil {
				return deploymentErr
			}

			cmd.Printf("Transaction hash: %s\nContract address: %s\n", deploymentTransaction.Hash().Hex(), address.Hex())
			if transactionOpts.NoSend {
				estimationMessage := ethereum.CallMsg{
					From: transactionOpts.From,
					Data: deploymentTransaction.Data(),
				}

				gasEstimationCtx, cancelGasEstimationCtx := NewChainContext(timeout)
				defer cancelGasEstimationCtx()

				gasEstimate, gasEstimateErr := client.EstimateGas(gasEstimationCtx, estimationMessage)
				if gasEstimateErr != nil {
					return gasEstimateErr
				}

				transactionBinary, transactionBinaryErr := deploymentTransaction.MarshalBinary()
				if transactionBinaryErr != nil {
					return transactionBinaryErr
				}
				transactionBinaryHex := hex.EncodeToString(transactionBinary)

				cmd.Printf("Transaction: %s\nEstimated gas: %d\n", transactionBinaryHex, gasEstimate)
			} else {
				cmd.Println("Transaction submitted")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the keystore file to use for the transaction")
	cmd.Flags().StringVar(&password, "password", "", "Password to use to unlock the keystore (if not specified, you will be prompted for the password when the command executes)")
	cmd.Flags().StringVar(&nonce, "nonce", "", "Nonce to use for the transaction")
	cmd.Flags().StringVar(&value, "value", "", "Value to send with the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price to use for the transaction")
	cmd.Flags().StringVar(&maxFeePerGas, "max-fee-per-gas", "", "Maximum fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().StringVar(&maxPriorityFeePerGas, "max-priority-fee-per-gas", "", "Maximum priority fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 0, "Gas limit for the transaction")
	cmd.Flags().BoolVar(&simulate, "simulate", false, "Simulate the transaction without sending it")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")

	cmd.Flags().StringVar(&name, "name", "", "name argument")
	cmd.Flags().StringVar(&symbol, "symbol", "", "symbol argument")
	cmd.Flags().StringVar(&ownerRaw, "owner", "", "owner argument")

	return cmd
}

func CreateIsApprovedForAllCommand() *cobra.Command {
	var contractAddressRaw, rpc string
	var contractAddress common.Address
	var timeout uint

	var blockNumberRaw, fromAddressRaw string
	var pending bool

	var owner common.Address
	var ownerRaw string
	var operator common.Address
	var operatorRaw string

	var capture0 bool

	cmd := &cobra.Command{
		Use:   "is-approved-for-all",
		Short: "Call the IsApprovedForAll view method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if ownerRaw == "" {
				return fmt.Errorf("--owner argument not specified")
			} else if !common.IsHexAddress(ownerRaw) {
				return fmt.Errorf("--owner argument is not a valid Ethereum address")
			}
			owner = common.HexToAddress(ownerRaw)

			if operatorRaw == "" {
				return fmt.Errorf("--operator argument not specified")
			} else if !common.IsHexAddress(operatorRaw) {
				return fmt.Errorf("--operator argument is not a valid Ethereum address")
			}
			operator = common.HexToAddress(operatorRaw)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			callOpts := bind.CallOpts{}
			SetCallParametersFromArgs(&callOpts, pending, fromAddressRaw, blockNumberRaw)

			session := OwnableERC721CallerSession{
				Contract: &contract.OwnableERC721Caller,
				CallOpts: callOpts,
			}

			var callErr error
			capture0, callErr = session.IsApprovedForAll(
				owner,
				operator,
			)
			if callErr != nil {
				return callErr
			}

			capture0JSON, capture0JSONMarshalErr := json.Marshal(capture0)
			if capture0JSONMarshalErr != nil {
				return capture0JSONMarshalErr
			}
			cmd.Printf("0: %s\\n", string(capture0JSON))

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&blockNumberRaw, "block", "", "Block number at which to call the view method")
	cmd.Flags().BoolVar(&pending, "pending", false, "Set this flag if it's ok to call the view method against pending state")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")
	cmd.Flags().StringVar(&fromAddressRaw, "from", "", "Optional address for caller of the view method")

	cmd.Flags().StringVar(&ownerRaw, "owner", "", "owner argument")
	cmd.Flags().StringVar(&operatorRaw, "operator", "", "operator argument")

	return cmd
}
func CreateNameCommand() *cobra.Command {
	var contractAddressRaw, rpc string
	var contractAddress common.Address
	var timeout uint

	var blockNumberRaw, fromAddressRaw string
	var pending bool

	var capture0 string

	cmd := &cobra.Command{
		Use:   "name",
		Short: "Call the Name view method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			callOpts := bind.CallOpts{}
			SetCallParametersFromArgs(&callOpts, pending, fromAddressRaw, blockNumberRaw)

			session := OwnableERC721CallerSession{
				Contract: &contract.OwnableERC721Caller,
				CallOpts: callOpts,
			}

			var callErr error
			capture0, callErr = session.Name()
			if callErr != nil {
				return callErr
			}

			cmd.Printf("0: %s\n", capture0)

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&blockNumberRaw, "block", "", "Block number at which to call the view method")
	cmd.Flags().BoolVar(&pending, "pending", false, "Set this flag if it's ok to call the view method against pending state")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")
	cmd.Flags().StringVar(&fromAddressRaw, "from", "", "Optional address for caller of the view method")

	return cmd
}
func CreateOwnerCommand() *cobra.Command {
	var contractAddressRaw, rpc string
	var contractAddress common.Address
	var timeout uint

	var blockNumberRaw, fromAddressRaw string
	var pending bool

	var capture0 common.Address

	cmd := &cobra.Command{
		Use:   "owner",
		Short: "Call the Owner view method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			callOpts := bind.CallOpts{}
			SetCallParametersFromArgs(&callOpts, pending, fromAddressRaw, blockNumberRaw)

			session := OwnableERC721CallerSession{
				Contract: &contract.OwnableERC721Caller,
				CallOpts: callOpts,
			}

			var callErr error
			capture0, callErr = session.Owner()
			if callErr != nil {
				return callErr
			}

			cmd.Printf("0: %s\n", capture0.Hex())

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&blockNumberRaw, "block", "", "Block number at which to call the view method")
	cmd.Flags().BoolVar(&pending, "pending", false, "Set this flag if it's ok to call the view method against pending state")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")
	cmd.Flags().StringVar(&fromAddressRaw, "from", "", "Optional address for caller of the view method")

	return cmd
}
func CreateSupportsInterfaceCommand() *cobra.Command {
	var contractAddressRaw, rpc string
	var contractAddress common.Address
	var timeout uint

	var blockNumberRaw, fromAddressRaw string
	var pending bool

	var interfaceId [4]byte
	var interfaceIdRaw string

	var capture0 bool

	cmd := &cobra.Command{
		Use:   "supports-interface",
		Short: "Call the SupportsInterface view method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if interfaceIdRaw == "" {
				return fmt.Errorf("--interface-id argument not specified")
			} else if strings.HasPrefix(interfaceIdRaw, "@") {
				filename := strings.TrimPrefix(interfaceIdRaw, "@")
				contents, readErr := os.ReadFile(filename)
				if readErr != nil {
					return readErr
				}
				unmarshalErr := json.Unmarshal(contents, &interfaceId)
				if unmarshalErr != nil {
					return unmarshalErr
				}
			} else {
				unmarshalErr := json.Unmarshal([]byte(interfaceIdRaw), &interfaceId)
				if unmarshalErr != nil {
					return unmarshalErr
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			callOpts := bind.CallOpts{}
			SetCallParametersFromArgs(&callOpts, pending, fromAddressRaw, blockNumberRaw)

			session := OwnableERC721CallerSession{
				Contract: &contract.OwnableERC721Caller,
				CallOpts: callOpts,
			}

			var callErr error
			capture0, callErr = session.SupportsInterface(
				interfaceId,
			)
			if callErr != nil {
				return callErr
			}

			capture0JSON, capture0JSONMarshalErr := json.Marshal(capture0)
			if capture0JSONMarshalErr != nil {
				return capture0JSONMarshalErr
			}
			cmd.Printf("0: %s\\n", string(capture0JSON))

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&blockNumberRaw, "block", "", "Block number at which to call the view method")
	cmd.Flags().BoolVar(&pending, "pending", false, "Set this flag if it's ok to call the view method against pending state")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")
	cmd.Flags().StringVar(&fromAddressRaw, "from", "", "Optional address for caller of the view method")

	cmd.Flags().StringVar(&interfaceIdRaw, "interface-id", "", "interface-id argument")

	return cmd
}
func CreateSymbolCommand() *cobra.Command {
	var contractAddressRaw, rpc string
	var contractAddress common.Address
	var timeout uint

	var blockNumberRaw, fromAddressRaw string
	var pending bool

	var capture0 string

	cmd := &cobra.Command{
		Use:   "symbol",
		Short: "Call the Symbol view method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			callOpts := bind.CallOpts{}
			SetCallParametersFromArgs(&callOpts, pending, fromAddressRaw, blockNumberRaw)

			session := OwnableERC721CallerSession{
				Contract: &contract.OwnableERC721Caller,
				CallOpts: callOpts,
			}

			var callErr error
			capture0, callErr = session.Symbol()
			if callErr != nil {
				return callErr
			}

			cmd.Printf("0: %s\n", capture0)

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&blockNumberRaw, "block", "", "Block number at which to call the view method")
	cmd.Flags().BoolVar(&pending, "pending", false, "Set this flag if it's ok to call the view method against pending state")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")
	cmd.Flags().StringVar(&fromAddressRaw, "from", "", "Optional address for caller of the view method")

	return cmd
}
func CreateTokenUriCommand() *cobra.Command {
	var contractAddressRaw, rpc string
	var contractAddress common.Address
	var timeout uint

	var blockNumberRaw, fromAddressRaw string
	var pending bool

	var tokenId *big.Int
	var tokenIdRaw string

	var capture0 string

	cmd := &cobra.Command{
		Use:   "token-uri",
		Short: "Call the TokenURI view method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if tokenIdRaw == "" {
				return fmt.Errorf("--token-id argument not specified")
			}
			tokenId = new(big.Int)
			tokenId.SetString(tokenIdRaw, 0)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			callOpts := bind.CallOpts{}
			SetCallParametersFromArgs(&callOpts, pending, fromAddressRaw, blockNumberRaw)

			session := OwnableERC721CallerSession{
				Contract: &contract.OwnableERC721Caller,
				CallOpts: callOpts,
			}

			var callErr error
			capture0, callErr = session.TokenURI(
				tokenId,
			)
			if callErr != nil {
				return callErr
			}

			cmd.Printf("0: %s\n", capture0)

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&blockNumberRaw, "block", "", "Block number at which to call the view method")
	cmd.Flags().BoolVar(&pending, "pending", false, "Set this flag if it's ok to call the view method against pending state")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")
	cmd.Flags().StringVar(&fromAddressRaw, "from", "", "Optional address for caller of the view method")

	cmd.Flags().StringVar(&tokenIdRaw, "token-id", "", "token-id argument")

	return cmd
}
func CreateBalanceOfCommand() *cobra.Command {
	var contractAddressRaw, rpc string
	var contractAddress common.Address
	var timeout uint

	var blockNumberRaw, fromAddressRaw string
	var pending bool

	var owner common.Address
	var ownerRaw string

	var capture0 *big.Int

	cmd := &cobra.Command{
		Use:   "balance-of",
		Short: "Call the BalanceOf view method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if ownerRaw == "" {
				return fmt.Errorf("--owner argument not specified")
			} else if !common.IsHexAddress(ownerRaw) {
				return fmt.Errorf("--owner argument is not a valid Ethereum address")
			}
			owner = common.HexToAddress(ownerRaw)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			callOpts := bind.CallOpts{}
			SetCallParametersFromArgs(&callOpts, pending, fromAddressRaw, blockNumberRaw)

			session := OwnableERC721CallerSession{
				Contract: &contract.OwnableERC721Caller,
				CallOpts: callOpts,
			}

			var callErr error
			capture0, callErr = session.BalanceOf(
				owner,
			)
			if callErr != nil {
				return callErr
			}

			cmd.Printf("0: %s\n", capture0.String())

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&blockNumberRaw, "block", "", "Block number at which to call the view method")
	cmd.Flags().BoolVar(&pending, "pending", false, "Set this flag if it's ok to call the view method against pending state")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")
	cmd.Flags().StringVar(&fromAddressRaw, "from", "", "Optional address for caller of the view method")

	cmd.Flags().StringVar(&ownerRaw, "owner", "", "owner argument")

	return cmd
}
func CreateGetApprovedCommand() *cobra.Command {
	var contractAddressRaw, rpc string
	var contractAddress common.Address
	var timeout uint

	var blockNumberRaw, fromAddressRaw string
	var pending bool

	var tokenId *big.Int
	var tokenIdRaw string

	var capture0 common.Address

	cmd := &cobra.Command{
		Use:   "get-approved",
		Short: "Call the GetApproved view method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if tokenIdRaw == "" {
				return fmt.Errorf("--token-id argument not specified")
			}
			tokenId = new(big.Int)
			tokenId.SetString(tokenIdRaw, 0)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			callOpts := bind.CallOpts{}
			SetCallParametersFromArgs(&callOpts, pending, fromAddressRaw, blockNumberRaw)

			session := OwnableERC721CallerSession{
				Contract: &contract.OwnableERC721Caller,
				CallOpts: callOpts,
			}

			var callErr error
			capture0, callErr = session.GetApproved(
				tokenId,
			)
			if callErr != nil {
				return callErr
			}

			cmd.Printf("0: %s\n", capture0.Hex())

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&blockNumberRaw, "block", "", "Block number at which to call the view method")
	cmd.Flags().BoolVar(&pending, "pending", false, "Set this flag if it's ok to call the view method against pending state")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")
	cmd.Flags().StringVar(&fromAddressRaw, "from", "", "Optional address for caller of the view method")

	cmd.Flags().StringVar(&tokenIdRaw, "token-id", "", "token-id argument")

	return cmd
}
func CreateOwnerOfCommand() *cobra.Command {
	var contractAddressRaw, rpc string
	var contractAddress common.Address
	var timeout uint

	var blockNumberRaw, fromAddressRaw string
	var pending bool

	var tokenId *big.Int
	var tokenIdRaw string

	var capture0 common.Address

	cmd := &cobra.Command{
		Use:   "owner-of",
		Short: "Call the OwnerOf view method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if tokenIdRaw == "" {
				return fmt.Errorf("--token-id argument not specified")
			}
			tokenId = new(big.Int)
			tokenId.SetString(tokenIdRaw, 0)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			callOpts := bind.CallOpts{}
			SetCallParametersFromArgs(&callOpts, pending, fromAddressRaw, blockNumberRaw)

			session := OwnableERC721CallerSession{
				Contract: &contract.OwnableERC721Caller,
				CallOpts: callOpts,
			}

			var callErr error
			capture0, callErr = session.OwnerOf(
				tokenId,
			)
			if callErr != nil {
				return callErr
			}

			cmd.Printf("0: %s\n", capture0.Hex())

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&blockNumberRaw, "block", "", "Block number at which to call the view method")
	cmd.Flags().BoolVar(&pending, "pending", false, "Set this flag if it's ok to call the view method against pending state")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")
	cmd.Flags().StringVar(&fromAddressRaw, "from", "", "Optional address for caller of the view method")

	cmd.Flags().StringVar(&tokenIdRaw, "token-id", "", "token-id argument")

	return cmd
}

func CreateRenounceOwnershipCommand() *cobra.Command {
	var keyfile, nonce, password, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, rpc, contractAddressRaw string
	var gasLimit uint64
	var simulate bool
	var timeout uint
	var contractAddress common.Address

	cmd := &cobra.Command{
		Use:   "renounce-ownership",
		Short: "Execute the RenounceOwnership method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keyfile == "" {
				return fmt.Errorf("--keystore not specified")
			}

			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			key, keyErr := KeyFromFile(keyfile, password)
			if keyErr != nil {
				return keyErr
			}

			chainIDCtx, cancelChainIDCtx := NewChainContext(timeout)
			defer cancelChainIDCtx()
			chainID, chainIDErr := client.ChainID(chainIDCtx)
			if chainIDErr != nil {
				return chainIDErr
			}

			transactionOpts, transactionOptsErr := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
			if transactionOptsErr != nil {
				return transactionOptsErr
			}

			SetTransactionParametersFromArgs(transactionOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, gasLimit, simulate)

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			session := OwnableERC721TransactorSession{
				Contract:     &contract.OwnableERC721Transactor,
				TransactOpts: *transactionOpts,
			}

			transaction, transactionErr := session.RenounceOwnership()
			if transactionErr != nil {
				return transactionErr
			}

			cmd.Printf("Transaction hash: %s\n", transaction.Hash().Hex())
			if transactionOpts.NoSend {
				estimationMessage := ethereum.CallMsg{
					From: transactionOpts.From,
					To:   &contractAddress,
					Data: transaction.Data(),
				}

				gasEstimationCtx, cancelGasEstimationCtx := NewChainContext(timeout)
				defer cancelGasEstimationCtx()

				gasEstimate, gasEstimateErr := client.EstimateGas(gasEstimationCtx, estimationMessage)
				if gasEstimateErr != nil {
					return gasEstimateErr
				}

				transactionBinary, transactionBinaryErr := transaction.MarshalBinary()
				if transactionBinaryErr != nil {
					return transactionBinaryErr
				}
				transactionBinaryHex := hex.EncodeToString(transactionBinary)

				cmd.Printf("Transaction: %s\nEstimated gas: %d\n", transactionBinaryHex, gasEstimate)
			} else {
				cmd.Println("Transaction submitted")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the keystore file to use for the transaction")
	cmd.Flags().StringVar(&password, "password", "", "Password to use to unlock the keystore (if not specified, you will be prompted for the password when the command executes)")
	cmd.Flags().StringVar(&nonce, "nonce", "", "Nonce to use for the transaction")
	cmd.Flags().StringVar(&value, "value", "", "Value to send with the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price to use for the transaction")
	cmd.Flags().StringVar(&maxFeePerGas, "max-fee-per-gas", "", "Maximum fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().StringVar(&maxPriorityFeePerGas, "max-priority-fee-per-gas", "", "Maximum priority fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 0, "Gas limit for the transaction")
	cmd.Flags().BoolVar(&simulate, "simulate", false, "Simulate the transaction without sending it")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")

	return cmd
}
func CreateSafeTransferFromCommand() *cobra.Command {
	var keyfile, nonce, password, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, rpc, contractAddressRaw string
	var gasLimit uint64
	var simulate bool
	var timeout uint
	var contractAddress common.Address

	var from0 common.Address
	var from0Raw string
	var to0 common.Address
	var to0Raw string
	var tokenId *big.Int
	var tokenIdRaw string

	cmd := &cobra.Command{
		Use:   "safe-transfer-from",
		Short: "Execute the SafeTransferFrom method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keyfile == "" {
				return fmt.Errorf("--keystore not specified")
			}

			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if from0Raw == "" {
				return fmt.Errorf("--from-0 argument not specified")
			} else if !common.IsHexAddress(from0Raw) {
				return fmt.Errorf("--from-0 argument is not a valid Ethereum address")
			}
			from0 = common.HexToAddress(from0Raw)

			if to0Raw == "" {
				return fmt.Errorf("--to-0 argument not specified")
			} else if !common.IsHexAddress(to0Raw) {
				return fmt.Errorf("--to-0 argument is not a valid Ethereum address")
			}
			to0 = common.HexToAddress(to0Raw)

			if tokenIdRaw == "" {
				return fmt.Errorf("--token-id argument not specified")
			}
			tokenId = new(big.Int)
			tokenId.SetString(tokenIdRaw, 0)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			key, keyErr := KeyFromFile(keyfile, password)
			if keyErr != nil {
				return keyErr
			}

			chainIDCtx, cancelChainIDCtx := NewChainContext(timeout)
			defer cancelChainIDCtx()
			chainID, chainIDErr := client.ChainID(chainIDCtx)
			if chainIDErr != nil {
				return chainIDErr
			}

			transactionOpts, transactionOptsErr := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
			if transactionOptsErr != nil {
				return transactionOptsErr
			}

			SetTransactionParametersFromArgs(transactionOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, gasLimit, simulate)

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			session := OwnableERC721TransactorSession{
				Contract:     &contract.OwnableERC721Transactor,
				TransactOpts: *transactionOpts,
			}

			transaction, transactionErr := session.SafeTransferFrom(
				from0,
				to0,
				tokenId,
			)
			if transactionErr != nil {
				return transactionErr
			}

			cmd.Printf("Transaction hash: %s\n", transaction.Hash().Hex())
			if transactionOpts.NoSend {
				estimationMessage := ethereum.CallMsg{
					From: transactionOpts.From,
					To:   &contractAddress,
					Data: transaction.Data(),
				}

				gasEstimationCtx, cancelGasEstimationCtx := NewChainContext(timeout)
				defer cancelGasEstimationCtx()

				gasEstimate, gasEstimateErr := client.EstimateGas(gasEstimationCtx, estimationMessage)
				if gasEstimateErr != nil {
					return gasEstimateErr
				}

				transactionBinary, transactionBinaryErr := transaction.MarshalBinary()
				if transactionBinaryErr != nil {
					return transactionBinaryErr
				}
				transactionBinaryHex := hex.EncodeToString(transactionBinary)

				cmd.Printf("Transaction: %s\nEstimated gas: %d\n", transactionBinaryHex, gasEstimate)
			} else {
				cmd.Println("Transaction submitted")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the keystore file to use for the transaction")
	cmd.Flags().StringVar(&password, "password", "", "Password to use to unlock the keystore (if not specified, you will be prompted for the password when the command executes)")
	cmd.Flags().StringVar(&nonce, "nonce", "", "Nonce to use for the transaction")
	cmd.Flags().StringVar(&value, "value", "", "Value to send with the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price to use for the transaction")
	cmd.Flags().StringVar(&maxFeePerGas, "max-fee-per-gas", "", "Maximum fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().StringVar(&maxPriorityFeePerGas, "max-priority-fee-per-gas", "", "Maximum priority fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 0, "Gas limit for the transaction")
	cmd.Flags().BoolVar(&simulate, "simulate", false, "Simulate the transaction without sending it")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")

	cmd.Flags().StringVar(&from0Raw, "from-0", "", "from-0 argument")
	cmd.Flags().StringVar(&to0Raw, "to-0", "", "to-0 argument")
	cmd.Flags().StringVar(&tokenIdRaw, "token-id", "", "token-id argument")

	return cmd
}
func CreateSafeTransferFrom0Command() *cobra.Command {
	var keyfile, nonce, password, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, rpc, contractAddressRaw string
	var gasLimit uint64
	var simulate bool
	var timeout uint
	var contractAddress common.Address

	var from0 common.Address
	var from0Raw string
	var to0 common.Address
	var to0Raw string
	var tokenId *big.Int
	var tokenIdRaw string
	var data []byte
	var dataRaw string

	cmd := &cobra.Command{
		Use:   "safe-transfer-from-0",
		Short: "Execute the SafeTransferFrom0 method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keyfile == "" {
				return fmt.Errorf("--keystore not specified")
			}

			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if from0Raw == "" {
				return fmt.Errorf("--from-0 argument not specified")
			} else if !common.IsHexAddress(from0Raw) {
				return fmt.Errorf("--from-0 argument is not a valid Ethereum address")
			}
			from0 = common.HexToAddress(from0Raw)

			if to0Raw == "" {
				return fmt.Errorf("--to-0 argument not specified")
			} else if !common.IsHexAddress(to0Raw) {
				return fmt.Errorf("--to-0 argument is not a valid Ethereum address")
			}
			to0 = common.HexToAddress(to0Raw)

			if tokenIdRaw == "" {
				return fmt.Errorf("--token-id argument not specified")
			}
			tokenId = new(big.Int)
			tokenId.SetString(tokenIdRaw, 0)

			if dataRaw == "" {
				return fmt.Errorf("--data argument not specified")
			} else if strings.HasPrefix(dataRaw, "@") {
				filename := strings.TrimPrefix(dataRaw, "@")
				contents, readErr := os.ReadFile(filename)
				if readErr != nil {
					return readErr
				}
				unmarshalErr := json.Unmarshal(contents, &data)
				if unmarshalErr != nil {
					return unmarshalErr
				}
			} else {
				unmarshalErr := json.Unmarshal([]byte(dataRaw), &data)
				if unmarshalErr != nil {
					return unmarshalErr
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			key, keyErr := KeyFromFile(keyfile, password)
			if keyErr != nil {
				return keyErr
			}

			chainIDCtx, cancelChainIDCtx := NewChainContext(timeout)
			defer cancelChainIDCtx()
			chainID, chainIDErr := client.ChainID(chainIDCtx)
			if chainIDErr != nil {
				return chainIDErr
			}

			transactionOpts, transactionOptsErr := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
			if transactionOptsErr != nil {
				return transactionOptsErr
			}

			SetTransactionParametersFromArgs(transactionOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, gasLimit, simulate)

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			session := OwnableERC721TransactorSession{
				Contract:     &contract.OwnableERC721Transactor,
				TransactOpts: *transactionOpts,
			}

			transaction, transactionErr := session.SafeTransferFrom0(
				from0,
				to0,
				tokenId,
				data,
			)
			if transactionErr != nil {
				return transactionErr
			}

			cmd.Printf("Transaction hash: %s\n", transaction.Hash().Hex())
			if transactionOpts.NoSend {
				estimationMessage := ethereum.CallMsg{
					From: transactionOpts.From,
					To:   &contractAddress,
					Data: transaction.Data(),
				}

				gasEstimationCtx, cancelGasEstimationCtx := NewChainContext(timeout)
				defer cancelGasEstimationCtx()

				gasEstimate, gasEstimateErr := client.EstimateGas(gasEstimationCtx, estimationMessage)
				if gasEstimateErr != nil {
					return gasEstimateErr
				}

				transactionBinary, transactionBinaryErr := transaction.MarshalBinary()
				if transactionBinaryErr != nil {
					return transactionBinaryErr
				}
				transactionBinaryHex := hex.EncodeToString(transactionBinary)

				cmd.Printf("Transaction: %s\nEstimated gas: %d\n", transactionBinaryHex, gasEstimate)
			} else {
				cmd.Println("Transaction submitted")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the keystore file to use for the transaction")
	cmd.Flags().StringVar(&password, "password", "", "Password to use to unlock the keystore (if not specified, you will be prompted for the password when the command executes)")
	cmd.Flags().StringVar(&nonce, "nonce", "", "Nonce to use for the transaction")
	cmd.Flags().StringVar(&value, "value", "", "Value to send with the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price to use for the transaction")
	cmd.Flags().StringVar(&maxFeePerGas, "max-fee-per-gas", "", "Maximum fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().StringVar(&maxPriorityFeePerGas, "max-priority-fee-per-gas", "", "Maximum priority fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 0, "Gas limit for the transaction")
	cmd.Flags().BoolVar(&simulate, "simulate", false, "Simulate the transaction without sending it")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")

	cmd.Flags().StringVar(&from0Raw, "from-0", "", "from-0 argument")
	cmd.Flags().StringVar(&to0Raw, "to-0", "", "to-0 argument")
	cmd.Flags().StringVar(&tokenIdRaw, "token-id", "", "token-id argument")
	cmd.Flags().StringVar(&dataRaw, "data", "", "data argument")

	return cmd
}
func CreateSetApprovalForAllCommand() *cobra.Command {
	var keyfile, nonce, password, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, rpc, contractAddressRaw string
	var gasLimit uint64
	var simulate bool
	var timeout uint
	var contractAddress common.Address

	var operator common.Address
	var operatorRaw string
	var approved bool
	var approvedRaw string

	cmd := &cobra.Command{
		Use:   "set-approval-for-all",
		Short: "Execute the SetApprovalForAll method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keyfile == "" {
				return fmt.Errorf("--keystore not specified")
			}

			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if operatorRaw == "" {
				return fmt.Errorf("--operator argument not specified")
			} else if !common.IsHexAddress(operatorRaw) {
				return fmt.Errorf("--operator argument is not a valid Ethereum address")
			}
			operator = common.HexToAddress(operatorRaw)

			approvedRawLower := strings.ToLower(approvedRaw)
			switch approvedRawLower {
			case "true", "t", "y", "yes", "1":
				approved = true
			case "false", "f", "n", "no", "0":
				approved = false
			default:
				return fmt.Errorf("--approved argument is not valid (value: %s)", approvedRaw)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			key, keyErr := KeyFromFile(keyfile, password)
			if keyErr != nil {
				return keyErr
			}

			chainIDCtx, cancelChainIDCtx := NewChainContext(timeout)
			defer cancelChainIDCtx()
			chainID, chainIDErr := client.ChainID(chainIDCtx)
			if chainIDErr != nil {
				return chainIDErr
			}

			transactionOpts, transactionOptsErr := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
			if transactionOptsErr != nil {
				return transactionOptsErr
			}

			SetTransactionParametersFromArgs(transactionOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, gasLimit, simulate)

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			session := OwnableERC721TransactorSession{
				Contract:     &contract.OwnableERC721Transactor,
				TransactOpts: *transactionOpts,
			}

			transaction, transactionErr := session.SetApprovalForAll(
				operator,
				approved,
			)
			if transactionErr != nil {
				return transactionErr
			}

			cmd.Printf("Transaction hash: %s\n", transaction.Hash().Hex())
			if transactionOpts.NoSend {
				estimationMessage := ethereum.CallMsg{
					From: transactionOpts.From,
					To:   &contractAddress,
					Data: transaction.Data(),
				}

				gasEstimationCtx, cancelGasEstimationCtx := NewChainContext(timeout)
				defer cancelGasEstimationCtx()

				gasEstimate, gasEstimateErr := client.EstimateGas(gasEstimationCtx, estimationMessage)
				if gasEstimateErr != nil {
					return gasEstimateErr
				}

				transactionBinary, transactionBinaryErr := transaction.MarshalBinary()
				if transactionBinaryErr != nil {
					return transactionBinaryErr
				}
				transactionBinaryHex := hex.EncodeToString(transactionBinary)

				cmd.Printf("Transaction: %s\nEstimated gas: %d\n", transactionBinaryHex, gasEstimate)
			} else {
				cmd.Println("Transaction submitted")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the keystore file to use for the transaction")
	cmd.Flags().StringVar(&password, "password", "", "Password to use to unlock the keystore (if not specified, you will be prompted for the password when the command executes)")
	cmd.Flags().StringVar(&nonce, "nonce", "", "Nonce to use for the transaction")
	cmd.Flags().StringVar(&value, "value", "", "Value to send with the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price to use for the transaction")
	cmd.Flags().StringVar(&maxFeePerGas, "max-fee-per-gas", "", "Maximum fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().StringVar(&maxPriorityFeePerGas, "max-priority-fee-per-gas", "", "Maximum priority fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 0, "Gas limit for the transaction")
	cmd.Flags().BoolVar(&simulate, "simulate", false, "Simulate the transaction without sending it")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")

	cmd.Flags().StringVar(&operatorRaw, "operator", "", "operator argument")
	cmd.Flags().StringVar(&approvedRaw, "approved", "", "approved argument (true, t, y, yes, 1 OR false, f, n, no, 0)")

	return cmd
}
func CreateTransferFromCommand() *cobra.Command {
	var keyfile, nonce, password, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, rpc, contractAddressRaw string
	var gasLimit uint64
	var simulate bool
	var timeout uint
	var contractAddress common.Address

	var from0 common.Address
	var from0Raw string
	var to0 common.Address
	var to0Raw string
	var tokenId *big.Int
	var tokenIdRaw string

	cmd := &cobra.Command{
		Use:   "transfer-from",
		Short: "Execute the TransferFrom method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keyfile == "" {
				return fmt.Errorf("--keystore not specified")
			}

			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if from0Raw == "" {
				return fmt.Errorf("--from-0 argument not specified")
			} else if !common.IsHexAddress(from0Raw) {
				return fmt.Errorf("--from-0 argument is not a valid Ethereum address")
			}
			from0 = common.HexToAddress(from0Raw)

			if to0Raw == "" {
				return fmt.Errorf("--to-0 argument not specified")
			} else if !common.IsHexAddress(to0Raw) {
				return fmt.Errorf("--to-0 argument is not a valid Ethereum address")
			}
			to0 = common.HexToAddress(to0Raw)

			if tokenIdRaw == "" {
				return fmt.Errorf("--token-id argument not specified")
			}
			tokenId = new(big.Int)
			tokenId.SetString(tokenIdRaw, 0)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			key, keyErr := KeyFromFile(keyfile, password)
			if keyErr != nil {
				return keyErr
			}

			chainIDCtx, cancelChainIDCtx := NewChainContext(timeout)
			defer cancelChainIDCtx()
			chainID, chainIDErr := client.ChainID(chainIDCtx)
			if chainIDErr != nil {
				return chainIDErr
			}

			transactionOpts, transactionOptsErr := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
			if transactionOptsErr != nil {
				return transactionOptsErr
			}

			SetTransactionParametersFromArgs(transactionOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, gasLimit, simulate)

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			session := OwnableERC721TransactorSession{
				Contract:     &contract.OwnableERC721Transactor,
				TransactOpts: *transactionOpts,
			}

			transaction, transactionErr := session.TransferFrom(
				from0,
				to0,
				tokenId,
			)
			if transactionErr != nil {
				return transactionErr
			}

			cmd.Printf("Transaction hash: %s\n", transaction.Hash().Hex())
			if transactionOpts.NoSend {
				estimationMessage := ethereum.CallMsg{
					From: transactionOpts.From,
					To:   &contractAddress,
					Data: transaction.Data(),
				}

				gasEstimationCtx, cancelGasEstimationCtx := NewChainContext(timeout)
				defer cancelGasEstimationCtx()

				gasEstimate, gasEstimateErr := client.EstimateGas(gasEstimationCtx, estimationMessage)
				if gasEstimateErr != nil {
					return gasEstimateErr
				}

				transactionBinary, transactionBinaryErr := transaction.MarshalBinary()
				if transactionBinaryErr != nil {
					return transactionBinaryErr
				}
				transactionBinaryHex := hex.EncodeToString(transactionBinary)

				cmd.Printf("Transaction: %s\nEstimated gas: %d\n", transactionBinaryHex, gasEstimate)
			} else {
				cmd.Println("Transaction submitted")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the keystore file to use for the transaction")
	cmd.Flags().StringVar(&password, "password", "", "Password to use to unlock the keystore (if not specified, you will be prompted for the password when the command executes)")
	cmd.Flags().StringVar(&nonce, "nonce", "", "Nonce to use for the transaction")
	cmd.Flags().StringVar(&value, "value", "", "Value to send with the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price to use for the transaction")
	cmd.Flags().StringVar(&maxFeePerGas, "max-fee-per-gas", "", "Maximum fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().StringVar(&maxPriorityFeePerGas, "max-priority-fee-per-gas", "", "Maximum priority fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 0, "Gas limit for the transaction")
	cmd.Flags().BoolVar(&simulate, "simulate", false, "Simulate the transaction without sending it")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")

	cmd.Flags().StringVar(&from0Raw, "from-0", "", "from-0 argument")
	cmd.Flags().StringVar(&to0Raw, "to-0", "", "to-0 argument")
	cmd.Flags().StringVar(&tokenIdRaw, "token-id", "", "token-id argument")

	return cmd
}
func CreateTransferOwnershipCommand() *cobra.Command {
	var keyfile, nonce, password, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, rpc, contractAddressRaw string
	var gasLimit uint64
	var simulate bool
	var timeout uint
	var contractAddress common.Address

	var newOwner common.Address
	var newOwnerRaw string

	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "Execute the TransferOwnership method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keyfile == "" {
				return fmt.Errorf("--keystore not specified")
			}

			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if newOwnerRaw == "" {
				return fmt.Errorf("--new-owner argument not specified")
			} else if !common.IsHexAddress(newOwnerRaw) {
				return fmt.Errorf("--new-owner argument is not a valid Ethereum address")
			}
			newOwner = common.HexToAddress(newOwnerRaw)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			key, keyErr := KeyFromFile(keyfile, password)
			if keyErr != nil {
				return keyErr
			}

			chainIDCtx, cancelChainIDCtx := NewChainContext(timeout)
			defer cancelChainIDCtx()
			chainID, chainIDErr := client.ChainID(chainIDCtx)
			if chainIDErr != nil {
				return chainIDErr
			}

			transactionOpts, transactionOptsErr := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
			if transactionOptsErr != nil {
				return transactionOptsErr
			}

			SetTransactionParametersFromArgs(transactionOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, gasLimit, simulate)

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			session := OwnableERC721TransactorSession{
				Contract:     &contract.OwnableERC721Transactor,
				TransactOpts: *transactionOpts,
			}

			transaction, transactionErr := session.TransferOwnership(
				newOwner,
			)
			if transactionErr != nil {
				return transactionErr
			}

			cmd.Printf("Transaction hash: %s\n", transaction.Hash().Hex())
			if transactionOpts.NoSend {
				estimationMessage := ethereum.CallMsg{
					From: transactionOpts.From,
					To:   &contractAddress,
					Data: transaction.Data(),
				}

				gasEstimationCtx, cancelGasEstimationCtx := NewChainContext(timeout)
				defer cancelGasEstimationCtx()

				gasEstimate, gasEstimateErr := client.EstimateGas(gasEstimationCtx, estimationMessage)
				if gasEstimateErr != nil {
					return gasEstimateErr
				}

				transactionBinary, transactionBinaryErr := transaction.MarshalBinary()
				if transactionBinaryErr != nil {
					return transactionBinaryErr
				}
				transactionBinaryHex := hex.EncodeToString(transactionBinary)

				cmd.Printf("Transaction: %s\nEstimated gas: %d\n", transactionBinaryHex, gasEstimate)
			} else {
				cmd.Println("Transaction submitted")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the keystore file to use for the transaction")
	cmd.Flags().StringVar(&password, "password", "", "Password to use to unlock the keystore (if not specified, you will be prompted for the password when the command executes)")
	cmd.Flags().StringVar(&nonce, "nonce", "", "Nonce to use for the transaction")
	cmd.Flags().StringVar(&value, "value", "", "Value to send with the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price to use for the transaction")
	cmd.Flags().StringVar(&maxFeePerGas, "max-fee-per-gas", "", "Maximum fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().StringVar(&maxPriorityFeePerGas, "max-priority-fee-per-gas", "", "Maximum priority fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 0, "Gas limit for the transaction")
	cmd.Flags().BoolVar(&simulate, "simulate", false, "Simulate the transaction without sending it")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")

	cmd.Flags().StringVar(&newOwnerRaw, "new-owner", "", "new-owner argument")

	return cmd
}
func CreateApproveCommand() *cobra.Command {
	var keyfile, nonce, password, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, rpc, contractAddressRaw string
	var gasLimit uint64
	var simulate bool
	var timeout uint
	var contractAddress common.Address

	var to0 common.Address
	var to0Raw string
	var tokenId *big.Int
	var tokenIdRaw string

	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Execute the Approve method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keyfile == "" {
				return fmt.Errorf("--keystore not specified")
			}

			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if to0Raw == "" {
				return fmt.Errorf("--to-0 argument not specified")
			} else if !common.IsHexAddress(to0Raw) {
				return fmt.Errorf("--to-0 argument is not a valid Ethereum address")
			}
			to0 = common.HexToAddress(to0Raw)

			if tokenIdRaw == "" {
				return fmt.Errorf("--token-id argument not specified")
			}
			tokenId = new(big.Int)
			tokenId.SetString(tokenIdRaw, 0)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			key, keyErr := KeyFromFile(keyfile, password)
			if keyErr != nil {
				return keyErr
			}

			chainIDCtx, cancelChainIDCtx := NewChainContext(timeout)
			defer cancelChainIDCtx()
			chainID, chainIDErr := client.ChainID(chainIDCtx)
			if chainIDErr != nil {
				return chainIDErr
			}

			transactionOpts, transactionOptsErr := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
			if transactionOptsErr != nil {
				return transactionOptsErr
			}

			SetTransactionParametersFromArgs(transactionOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, gasLimit, simulate)

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			session := OwnableERC721TransactorSession{
				Contract:     &contract.OwnableERC721Transactor,
				TransactOpts: *transactionOpts,
			}

			transaction, transactionErr := session.Approve(
				to0,
				tokenId,
			)
			if transactionErr != nil {
				return transactionErr
			}

			cmd.Printf("Transaction hash: %s\n", transaction.Hash().Hex())
			if transactionOpts.NoSend {
				estimationMessage := ethereum.CallMsg{
					From: transactionOpts.From,
					To:   &contractAddress,
					Data: transaction.Data(),
				}

				gasEstimationCtx, cancelGasEstimationCtx := NewChainContext(timeout)
				defer cancelGasEstimationCtx()

				gasEstimate, gasEstimateErr := client.EstimateGas(gasEstimationCtx, estimationMessage)
				if gasEstimateErr != nil {
					return gasEstimateErr
				}

				transactionBinary, transactionBinaryErr := transaction.MarshalBinary()
				if transactionBinaryErr != nil {
					return transactionBinaryErr
				}
				transactionBinaryHex := hex.EncodeToString(transactionBinary)

				cmd.Printf("Transaction: %s\nEstimated gas: %d\n", transactionBinaryHex, gasEstimate)
			} else {
				cmd.Println("Transaction submitted")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the keystore file to use for the transaction")
	cmd.Flags().StringVar(&password, "password", "", "Password to use to unlock the keystore (if not specified, you will be prompted for the password when the command executes)")
	cmd.Flags().StringVar(&nonce, "nonce", "", "Nonce to use for the transaction")
	cmd.Flags().StringVar(&value, "value", "", "Value to send with the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price to use for the transaction")
	cmd.Flags().StringVar(&maxFeePerGas, "max-fee-per-gas", "", "Maximum fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().StringVar(&maxPriorityFeePerGas, "max-priority-fee-per-gas", "", "Maximum priority fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 0, "Gas limit for the transaction")
	cmd.Flags().BoolVar(&simulate, "simulate", false, "Simulate the transaction without sending it")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")

	cmd.Flags().StringVar(&to0Raw, "to-0", "", "to-0 argument")
	cmd.Flags().StringVar(&tokenIdRaw, "token-id", "", "token-id argument")

	return cmd
}
func CreateMintCommand() *cobra.Command {
	var keyfile, nonce, password, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, rpc, contractAddressRaw string
	var gasLimit uint64
	var simulate bool
	var timeout uint
	var contractAddress common.Address

	var to0 common.Address
	var to0Raw string
	var tokenId *big.Int
	var tokenIdRaw string

	cmd := &cobra.Command{
		Use:   "mint",
		Short: "Execute the Mint method on a OwnableERC721 contract",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keyfile == "" {
				return fmt.Errorf("--keystore not specified")
			}

			if contractAddressRaw == "" {
				return fmt.Errorf("--contract not specified")
			} else if !common.IsHexAddress(contractAddressRaw) {
				return fmt.Errorf("--contract is not a valid Ethereum address")
			}
			contractAddress = common.HexToAddress(contractAddressRaw)

			if to0Raw == "" {
				return fmt.Errorf("--to-0 argument not specified")
			} else if !common.IsHexAddress(to0Raw) {
				return fmt.Errorf("--to-0 argument is not a valid Ethereum address")
			}
			to0 = common.HexToAddress(to0Raw)

			if tokenIdRaw == "" {
				return fmt.Errorf("--token-id argument not specified")
			}
			tokenId = new(big.Int)
			tokenId.SetString(tokenIdRaw, 0)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, clientErr := NewClient(rpc)
			if clientErr != nil {
				return clientErr
			}

			key, keyErr := KeyFromFile(keyfile, password)
			if keyErr != nil {
				return keyErr
			}

			chainIDCtx, cancelChainIDCtx := NewChainContext(timeout)
			defer cancelChainIDCtx()
			chainID, chainIDErr := client.ChainID(chainIDCtx)
			if chainIDErr != nil {
				return chainIDErr
			}

			transactionOpts, transactionOptsErr := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chainID)
			if transactionOptsErr != nil {
				return transactionOptsErr
			}

			SetTransactionParametersFromArgs(transactionOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas, gasLimit, simulate)

			contract, contractErr := NewOwnableERC721(contractAddress, client)
			if contractErr != nil {
				return contractErr
			}

			session := OwnableERC721TransactorSession{
				Contract:     &contract.OwnableERC721Transactor,
				TransactOpts: *transactionOpts,
			}

			transaction, transactionErr := session.Mint(
				to0,
				tokenId,
			)
			if transactionErr != nil {
				return transactionErr
			}

			cmd.Printf("Transaction hash: %s\n", transaction.Hash().Hex())
			if transactionOpts.NoSend {
				estimationMessage := ethereum.CallMsg{
					From: transactionOpts.From,
					To:   &contractAddress,
					Data: transaction.Data(),
				}

				gasEstimationCtx, cancelGasEstimationCtx := NewChainContext(timeout)
				defer cancelGasEstimationCtx()

				gasEstimate, gasEstimateErr := client.EstimateGas(gasEstimationCtx, estimationMessage)
				if gasEstimateErr != nil {
					return gasEstimateErr
				}

				transactionBinary, transactionBinaryErr := transaction.MarshalBinary()
				if transactionBinaryErr != nil {
					return transactionBinaryErr
				}
				transactionBinaryHex := hex.EncodeToString(transactionBinary)

				cmd.Printf("Transaction: %s\nEstimated gas: %d\n", transactionBinaryHex, gasEstimate)
			} else {
				cmd.Println("Transaction submitted")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rpc, "rpc", "", "URL of the JSONRPC API to use")
	cmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the keystore file to use for the transaction")
	cmd.Flags().StringVar(&password, "password", "", "Password to use to unlock the keystore (if not specified, you will be prompted for the password when the command executes)")
	cmd.Flags().StringVar(&nonce, "nonce", "", "Nonce to use for the transaction")
	cmd.Flags().StringVar(&value, "value", "", "Value to send with the transaction")
	cmd.Flags().StringVar(&gasPrice, "gas-price", "", "Gas price to use for the transaction")
	cmd.Flags().StringVar(&maxFeePerGas, "max-fee-per-gas", "", "Maximum fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().StringVar(&maxPriorityFeePerGas, "max-priority-fee-per-gas", "", "Maximum priority fee per gas to use for the (EIP-1559) transaction")
	cmd.Flags().Uint64Var(&gasLimit, "gas-limit", 0, "Gas limit for the transaction")
	cmd.Flags().BoolVar(&simulate, "simulate", false, "Simulate the transaction without sending it")
	cmd.Flags().UintVar(&timeout, "timeout", 60, "Timeout (in seconds) for interactions with the JSONRPC API")
	cmd.Flags().StringVar(&contractAddressRaw, "contract", "", "Address of the contract to interact with")

	cmd.Flags().StringVar(&to0Raw, "to-0", "", "to-0 argument")
	cmd.Flags().StringVar(&tokenIdRaw, "token-id", "", "token-id argument")

	return cmd
}

var ErrNoRPCURL error = errors.New("no RPC URL provided -- please pass an RPC URL from the command line or set the OWNABLE_ERC_721_RPC_URL environment variable")

// Generates an Ethereum client to the JSONRPC API at the given URL. If rpcURL is empty, then it
// attempts to read the RPC URL from the OWNABLE_ERC_721_RPC_URL environment variable. If that is empty,
// too, then it returns an error.
func NewClient(rpcURL string) (*ethclient.Client, error) {
	if rpcURL == "" {
		rpcURL = os.Getenv("OWNABLE_ERC_721_RPC_URL")
	}

	if rpcURL == "" {
		return nil, ErrNoRPCURL
	}

	client, err := ethclient.Dial(rpcURL)
	return client, err
}

// Creates a new context to be used when interacting with the chain client.
func NewChainContext(timeout uint) (context.Context, context.CancelFunc) {
	baseCtx := context.Background()
	parsedTimeout := time.Duration(timeout) * time.Second
	ctx, cancel := context.WithTimeout(baseCtx, parsedTimeout)
	return ctx, cancel
}

// Unlocks a key from a keystore (byte contents of a keystore file) with the given password.
func UnlockKeystore(keystoreData []byte, password string) (*keystore.Key, error) {
	key, err := keystore.DecryptKey(keystoreData, password)
	return key, err
}

// Loads a key from file, prompting the user for the password if it is not provided as a function argument.
func KeyFromFile(keystoreFile string, password string) (*keystore.Key, error) {
	var emptyKey *keystore.Key
	keystoreContent, readErr := os.ReadFile(keystoreFile)
	if readErr != nil {
		return emptyKey, readErr
	}

	// If password is "", prompt user for password.
	if password == "" {
		fmt.Printf("Please provide a password for keystore (%s): ", keystoreFile)
		passwordRaw, inputErr := term.ReadPassword(int(os.Stdin.Fd()))
		if inputErr != nil {
			return emptyKey, fmt.Errorf("error reading password: %s", inputErr.Error())
		}
		fmt.Print("\n")
		password = string(passwordRaw)
	}

	key, err := UnlockKeystore(keystoreContent, password)
	return key, err
}

// This method is used to set the parameters on a view call from command line arguments (represented mostly as
// strings).
func SetCallParametersFromArgs(opts *bind.CallOpts, pending bool, fromAddress, blockNumber string) {
	if pending {
		opts.Pending = true
	}

	if fromAddress != "" {
		opts.From = common.HexToAddress(fromAddress)
	}

	if blockNumber != "" {
		opts.BlockNumber = new(big.Int)
		opts.BlockNumber.SetString(blockNumber, 0)
	}
}

// This method is used to set the parameters on a transaction from command line arguments (represented mostly as
// strings).
func SetTransactionParametersFromArgs(opts *bind.TransactOpts, nonce, value, gasPrice, maxFeePerGas, maxPriorityFeePerGas string, gasLimit uint64, noSend bool) {
	if nonce != "" {
		opts.Nonce = new(big.Int)
		opts.Nonce.SetString(nonce, 0)
	}

	if value != "" {
		opts.Value = new(big.Int)
		opts.Value.SetString(value, 0)
	}

	if gasPrice != "" {
		opts.GasPrice = new(big.Int)
		opts.GasPrice.SetString(gasPrice, 0)
	}

	if maxFeePerGas != "" {
		opts.GasFeeCap = new(big.Int)
		opts.GasFeeCap.SetString(maxFeePerGas, 0)
	}

	if maxPriorityFeePerGas != "" {
		opts.GasTipCap = new(big.Int)
		opts.GasTipCap.SetString(maxPriorityFeePerGas, 0)
	}

	if gasLimit != 0 {
		opts.GasLimit = gasLimit
	}

	opts.NoSend = noSend
}

func CreateOwnableERC721Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ownable-erc-721",
		Short: "Interact with the OwnableERC721 contract",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.SetOut(os.Stdout)

	DeployGroup := &cobra.Group{
		ID: "deploy", Title: "Commands which deploy contracts",
	}
	cmd.AddGroup(DeployGroup)
	ViewGroup := &cobra.Group{
		ID: "view", Title: "Commands which view contract state",
	}
	TransactGroup := &cobra.Group{
		ID: "transact", Title: "Commands which submit transactions",
	}
	cmd.AddGroup(ViewGroup, TransactGroup)

	cmdDeployOwnableERC721 := CreateOwnableERC721DeploymentCommand()
	cmdDeployOwnableERC721.GroupID = DeployGroup.ID
	cmd.AddCommand(cmdDeployOwnableERC721)

	cmdViewIsApprovedForAll := CreateIsApprovedForAllCommand()
	cmdViewIsApprovedForAll.GroupID = ViewGroup.ID
	cmd.AddCommand(cmdViewIsApprovedForAll)
	cmdViewName := CreateNameCommand()
	cmdViewName.GroupID = ViewGroup.ID
	cmd.AddCommand(cmdViewName)
	cmdViewOwner := CreateOwnerCommand()
	cmdViewOwner.GroupID = ViewGroup.ID
	cmd.AddCommand(cmdViewOwner)
	cmdViewSupportsInterface := CreateSupportsInterfaceCommand()
	cmdViewSupportsInterface.GroupID = ViewGroup.ID
	cmd.AddCommand(cmdViewSupportsInterface)
	cmdViewSymbol := CreateSymbolCommand()
	cmdViewSymbol.GroupID = ViewGroup.ID
	cmd.AddCommand(cmdViewSymbol)
	cmdViewTokenURI := CreateTokenUriCommand()
	cmdViewTokenURI.GroupID = ViewGroup.ID
	cmd.AddCommand(cmdViewTokenURI)
	cmdViewBalanceOf := CreateBalanceOfCommand()
	cmdViewBalanceOf.GroupID = ViewGroup.ID
	cmd.AddCommand(cmdViewBalanceOf)
	cmdViewGetApproved := CreateGetApprovedCommand()
	cmdViewGetApproved.GroupID = ViewGroup.ID
	cmd.AddCommand(cmdViewGetApproved)
	cmdViewOwnerOf := CreateOwnerOfCommand()
	cmdViewOwnerOf.GroupID = ViewGroup.ID
	cmd.AddCommand(cmdViewOwnerOf)

	cmdTransactRenounceOwnership := CreateRenounceOwnershipCommand()
	cmdTransactRenounceOwnership.GroupID = TransactGroup.ID
	cmd.AddCommand(cmdTransactRenounceOwnership)
	cmdTransactSafeTransferFrom := CreateSafeTransferFromCommand()
	cmdTransactSafeTransferFrom.GroupID = TransactGroup.ID
	cmd.AddCommand(cmdTransactSafeTransferFrom)
	cmdTransactSafeTransferFrom0 := CreateSafeTransferFrom0Command()
	cmdTransactSafeTransferFrom0.GroupID = TransactGroup.ID
	cmd.AddCommand(cmdTransactSafeTransferFrom0)
	cmdTransactSetApprovalForAll := CreateSetApprovalForAllCommand()
	cmdTransactSetApprovalForAll.GroupID = TransactGroup.ID
	cmd.AddCommand(cmdTransactSetApprovalForAll)
	cmdTransactTransferFrom := CreateTransferFromCommand()
	cmdTransactTransferFrom.GroupID = TransactGroup.ID
	cmd.AddCommand(cmdTransactTransferFrom)
	cmdTransactTransferOwnership := CreateTransferOwnershipCommand()
	cmdTransactTransferOwnership.GroupID = TransactGroup.ID
	cmd.AddCommand(cmdTransactTransferOwnership)
	cmdTransactApprove := CreateApproveCommand()
	cmdTransactApprove.GroupID = TransactGroup.ID
	cmd.AddCommand(cmdTransactApprove)
	cmdTransactMint := CreateMintCommand()
	cmdTransactMint.GroupID = TransactGroup.ID
	cmd.AddCommand(cmdTransactMint)

	return cmd
}
