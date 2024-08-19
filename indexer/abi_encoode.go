package indexer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// Print contract function and event signatures from the ABI
func PrintContractSignatures(abiString string) (interface{}, error) {
	abiObj, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}

	// Print the interface header
	fmt.Println("// Interface generated seer")
	fmt.Println("// Interface ID: ___")
	fmt.Println("interface IGeneratedContract {")

	// Print events
	events := make([]string, 0, len(abiObj.Events))
	for _, event := range abiObj.Events {
		events = append(events, event.Name)
	}
	sort.Strings(events)

	fmt.Println("\t// events")
	for _, eventName := range events {
		event := abiObj.Events[eventName]
		fmt.Printf("\tevent %s signature %s;\n", event.Sig, event.ID.String())
	}

	// Print functions
	methods := make([]string, 0, len(abiObj.Methods))
	for methodName := range abiObj.Methods {
		methods = append(methods, methodName)
	}
	sort.Strings(methods)

	fmt.Println("\n\t// functions")
	for _, methodName := range methods {
		method := abiObj.Methods[methodName]
		fmt.Printf("\t// Selector: 0x%s\n", fmt.Sprintf("%x", method.ID))
		fmt.Printf("\tfunction %s;\n", method.Sig)
	}

	// Close the interface
	fmt.Println("\n\t// errors")
	fmt.Println("}")

	return abiObj, nil
}
