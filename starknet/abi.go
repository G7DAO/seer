package starknet

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Represents a particular value in a Starknet ABI enum.
type EnumVariant struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

// Represents an enum in a Starknet ABI.
type Enum struct {
	Type     string         `json:"type"`
	Name     string         `json:"name"`
	Variants []*EnumVariant `json:"variants"`
}

// Represents a member of a Starknet ABI struct.
type StructMember struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Kind string `json:"kind,omitempty"`
}

// Represents a struct in a Starknet ABI.
type Struct struct {
	Type    string          `json:"type"`
	Name    string          `json:"name"`
	Members []*StructMember `json:"members"`
}

// Represents a struct in a Starknet ABI which is used as an event.
type EventStruct struct {
	Type    string          `json:"type"`
	Name    string          `json:"name"`
	Kind    string          `json:"kind"`
	Members []*StructMember `json:"members"`
}

// Represents a single item in a Starknet ABI.
type ABIItemType struct {
	Type string `json:"type,omitempty"`
	Kind string `json:"kind,omitempty"`
}

// Represents a parsed Starknet ABI.
type ParsedABI struct {
	Enums   []*Enum        `json:"enums"`
	Structs []*Struct      `json:"structs"`
	Events  []*EventStruct `json:"events"`
}

// Internal representation of a Starknet ABI used while parsing the ABI into its Go representation as a
// ParsedABI.
type IntermediateABI struct {
	Types    []ABIItemType
	Messages []json.RawMessage
}

// Parses a Starknet ABI from a JSON byte array.
func ParseABI(rawABI []byte) (*ParsedABI, error) {
	parsedABI := &ParsedABI{}

	var itemTypes []ABIItemType
	var rawMessages []json.RawMessage
	intermediateUnmarshalErr := json.Unmarshal(rawABI, &itemTypes)
	if intermediateUnmarshalErr != nil {
		return parsedABI, intermediateUnmarshalErr
	}

	messagesUnmarshalErr := json.Unmarshal(rawABI, &rawMessages)
	if messagesUnmarshalErr != nil {
		return parsedABI, messagesUnmarshalErr
	}

	numEnums := 0
	numStructs := 0
	numEvents := 0
	for _, item := range itemTypes {
		switch item.Type {
		case "enum":
			numEnums++
		case "struct":
			numStructs++
		case "event":
			if item.Kind == "struct" {
				numEvents++
			}
		}
	}

	parsedABI.Enums = make([]*Enum, numEnums)
	currentEnum := 0
	for i, item := range itemTypes {
		if item.Type == "enum" {
			var enum *Enum
			enumUnmarshalErr := json.Unmarshal(rawMessages[i], &enum)
			if enumUnmarshalErr != nil {
				return parsedABI, enumUnmarshalErr
			}

			parsedABI.Enums[currentEnum] = enum
			currentEnum++
		}
	}

	for _, item := range parsedABI.Enums {
		for i, variant := range item.Variants {
			variant.Index = i
		}
	}

	parsedABI.Structs = make([]*Struct, numStructs)
	currentStruct := 0
	for i, item := range itemTypes {
		if item.Type == "struct" {
			var structItem *Struct
			structUnmarshalErr := json.Unmarshal(rawMessages[i], &structItem)
			if structUnmarshalErr != nil {
				return parsedABI, structUnmarshalErr
			}

			parsedABI.Structs[currentStruct] = structItem
			currentStruct++
		}
	}

	parsedABI.Events = make([]*EventStruct, numEvents)
	currentEvent := 0
	for i, item := range itemTypes {
		if item.Type == "event" && item.Kind == "struct" {
			var event *EventStruct
			eventUnmarshalErr := json.Unmarshal(rawMessages[i], &event)
			if eventUnmarshalErr != nil {
				return parsedABI, eventUnmarshalErr
			}

			parsedABI.Events[currentEvent] = event
			currentEvent++
		}
	}

	return parsedABI, nil
}

// Calculates the Starknet hash corresponding to the name of a Starknet ABI event. This hash is how
// the event is represented in Starknet event logs.
func HashFromName(name string) (string, error) {
	x := big.NewInt(0)
	mask := big.NewInt(0)

	x.Exp(big.NewInt(2), big.NewInt(250), nil)
	mask.Sub(x, big.NewInt(1))

	components := strings.Split(name, "::")
	eventName := components[len(components)-1]

	// Very important to use the LegacyKeccak256 here - to match Ethereum:
	// https://pkg.go.dev/golang.org/x/crypto/sha3#NewLegacyKeccak256
	hash := sha3.NewLegacyKeccak256()
	_, hashErr := hash.Write([]byte(eventName))
	if hashErr != nil {
		return "", hashErr
	}

	b := make([]byte, 0)
	hashedNameBytes := hash.Sum(b)

	hashedEncodedName := big.NewInt(0).SetBytes(hashedNameBytes)

	starknetHashedEncodedName := big.NewInt(0).And(hashedEncodedName, mask)
	return hex.EncodeToString(starknetHashedEncodedName.Bytes()), nil
}
