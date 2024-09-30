package common

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func GetABI(abistr string) (*abi.ABI, error) {
	// Retrieve or create AbiEntry
	// ...

	parsedABI, err := abi.JSON(strings.NewReader(abistr))
	if err != nil {
		fmt.Println("string: ", abistr)
		fmt.Println("Error parsing ABI: ", err)
		return nil, err
	}

	return &parsedABI, nil

}
