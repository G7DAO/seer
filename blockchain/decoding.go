package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func DecodeTransactionInputData(contractABI *abi.ABI, data []byte) (map[string]interface{}, error) {
	methodSigData := data[:4]
	inputsSigData := data[4:]
	method, err := contractABI.MethodById(methodSigData)
	if err != nil {
		log.Fatal(err)
	}
	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		return nil, err
	}

	// Prepare the extended map
	labelData := make(map[string]interface{})
	labelData["type"] = "tx_call"
	labelData["gas_used"] = 0
	labelData["args"] = inputsMap

	// check if labeData is valid json
	_, err = json.Marshal(labelData)
	if err != nil {
		fmt.Println("Error marshalling labelData: ", labelData)
		return nil, err
	}

	return labelData, nil
}

func DecodeLogArgsToLabelData(contractABI *abi.ABI, topics []string, data string) (map[string]interface{}, error) {

	topic0 := topics[0]

	// Convert the topic0 string to common.Hash
	topic0Hash := common.HexToHash(topic0)

	event, err := contractABI.EventByID(
		topic0Hash,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Decode the data string from hex to bytes
	dataBytes, err := hex.DecodeString(strings.TrimPrefix(data, "0x"))
	if err != nil {
		log.Fatalf("Failed to decode data string: %v", err)
	}

	// Prepare the map to hold the input data
	labelData := make(map[string]interface{})
	labelData["type"] = "event"
	labelData["name"] = event.Name
	labelData["args"] = make(map[string]interface{})

	i := 1
	// Extract indexed parameters from topics
	for _, input := range event.Inputs {
		var arg interface{}
		if input.Indexed {
			// Note: topic[0] is the event signature, so indexed params start from topic[1]
			switch input.Type.T {
			case abi.AddressTy:
				arg = common.HexToAddress(topics[i]).Hex()
			case abi.BytesTy:
				arg = common.HexToHash(topics[i]).Hex()
			case abi.FixedBytesTy:
				if input.Type.Size == 32 {
					arg = common.HexToHash(topics[i]).Hex()
				} else {
					arg = common.BytesToHash(common.Hex2Bytes(topics[i][2:])).Hex() // for other fixed sizes
				}
			case abi.UintTy:
				arg = new(big.Int).SetBytes(common.Hex2Bytes(topics[i][2:]))
			case abi.BoolTy:
				arg = new(big.Int).SetBytes(common.Hex2Bytes(topics[i][2:])).Cmp(big.NewInt(0)) != 0
			case abi.StringTy:
				argBytes, err := hex.DecodeString(strings.TrimPrefix(topics[i], "0x"))
				if err != nil {
					return nil, fmt.Errorf("failed to decode hex string to normal string: %v", err)
				}
				arg = string(argBytes)
			default:
				log.Fatalf("Unsupported indexed type: %s", input.Type.String())
			}
			i++
		} else {
			arg = "NON-INDEXED" // Placeholder for non-indexed arguments, which will be unpacked later
		}
		labelData["args"].(map[string]interface{})[input.Name] = arg
	}

	// Unpack the data bytes into the args map
	if err := event.Inputs.UnpackIntoMap(labelData["args"].(map[string]interface{}), dataBytes); err != nil {
		return nil, err
	}

	return labelData, nil
}
