package starknet

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/moonstream-to/seer/version"
)

// Common parameters required for the generation of all types of artifacts.
type GenerationParameters struct {
	OriginalName string
	GoName       string
}

// The output of the code generation process for enum items in a Starknet ABI.
type GeneratedEnum struct {
	GenerationParameters
	ParserName    string
	EvaluatorName string
	Definition    *Enum
	Code          string
}

// The output of the code generation process for struct items in a Starknet ABI.
type GeneratedStruct struct {
	GenerationParameters
	ParserName string
	Definition *Struct
	Code       string
}

type GeneratedEvent struct {
	GenerationParameters
	ParserName   string
	Definition   *EventStruct
	EventNameVar string
	EventHashVar string
	EventHash    string
	Code         string
}

// Defines the parameters used to create the header information for the generated code.
type HeaderParameters struct {
	Version     string
	PackageName string
}

// Generates a Go name for a Starknet ABI item given its fully qualified ABI name.
// Qualified names for Starknet ABI items are of the form:
// `core::starknet::contract_address::ContractAddress`
func GenerateGoNameForType(qualifiedName string) string {
	if qualifiedName == "core::integer::u8" || qualifiedName == "core::integer::u16" || qualifiedName == "core::integer::u32" || qualifiedName == "core::integer::u64" {
		return "uint64"
	} else if strings.HasPrefix(qualifiedName, "core::integer::") {
		return `*big.Int`
	} else if qualifiedName == "core::starknet::contract_address::ContractAddress" {
		return "string"
	} else if qualifiedName == "core::felt252" {
		return "string"
	} else if strings.HasPrefix(qualifiedName, "core::array::Array::<") {
		s1, _ := strings.CutPrefix(qualifiedName, "core::array::Array::<")
		s2, _ := strings.CutSuffix(s1, ">")
		return fmt.Sprintf("[]%s", GenerateGoNameForType(s2))
	}

	components := strings.Split(qualifiedName, "::")
	camelComponents := make([]string, len(components))
	for i, component := range components {
		camelComponents[i] = strcase.ToCamel(component)
	}
	return strings.Join(camelComponents, "_")
}

// Returns the name of the function that parses the given Go type.
func ParserFunction(goType string) string {
	baseType := ""
	numWraps := 0
	for strings.HasPrefix(baseType, "[]") {
		baseType = strings.TrimPrefix(baseType, "[]")
		numWraps++
	}

	var parserFunction string

	if numWraps == 0 {
		switch goType {
		case "uint64":
			parserFunction = "ParseUint64"
		case "*big.Int":
			parserFunction = "ParseBigInt"
		case "string":
			parserFunction = "ParseString"
		default:
			parserFunction = fmt.Sprintf("Parse%s", goType)
		}
	} else {
		baseParser := ParserFunction(baseType)
		parserFunction := ""
		for i := numWraps; i > 0; i-- {
			arrayParser := fmt.Sprintf("ParseArray[%s%s](", strings.Repeat("[]", i), baseType)
			parserFunction = fmt.Sprintf("%s%s", parserFunction, arrayParser)
		}
		parserFunction = fmt.Sprintf("%s%s%s", parserFunction, baseParser, strings.Repeat(")", numWraps))
	}

	return parserFunction
}

func ShouldGenerateStructType(goName string) bool {
	if goName == "uint64" || goName == "*big.Int" || goName == "string" || strings.HasPrefix(goName, "[]") {
		return false
	}
	return true
}

// Generate generates Go code for each of the items in a Starknet contract ABI.
// Returns a mapping of the go name of each object to a specification of the generated artifact.
// Currently supports:
// - Enums
// - Structs
// - Events
//
// ABI names are used to depuplicate code snippets. The assumption is that the Starknet fully
// qualified name for a type uniquely determines that type across the entire ABI. This way
// even if the ABI passed into the code generator contains duplicate instances of an ABI item,
// the Go code will only contain one definition of that item.
func GenerateSnippets(parsed *ParsedABI) (map[string]string, error) {
	result := map[string]string{}

	enumTemplate, enumTemplateParseErr := template.New("enum").Parse(EnumTemplate)
	if enumTemplateParseErr != nil {
		return result, enumTemplateParseErr
	}

	templateFuncs := map[string]any{
		"CamelCase":             strcase.ToCamel,
		"GenerateGoNameForType": GenerateGoNameForType,
		"ParserFunction":        ParserFunction,
	}

	structTemplate, structTemplateParseErr := template.New("struct").Funcs(templateFuncs).Parse(StructTemplate)
	if structTemplateParseErr != nil {
		return result, structTemplateParseErr
	}

	eventTemplate, eventTemplateParseErr := template.New("event").Funcs(templateFuncs).Parse(EventTemplate)
	if structTemplateParseErr != nil {
		return result, eventTemplateParseErr
	}

	for _, enum := range parsed.Enums {
		goName := GenerateGoNameForType(enum.Name)
		parseFunctionName := ParserFunction(goName)
		evaluateFunctionName := fmt.Sprintf("Evaluate%s", goName)

		generated := GeneratedEnum{
			GenerationParameters: GenerationParameters{
				OriginalName: enum.Name,
				GoName:       goName,
			},
			ParserName:    parseFunctionName,
			EvaluatorName: evaluateFunctionName,
			Definition:    enum,
			Code:          "",
		}

		var b bytes.Buffer
		templateErr := enumTemplate.Execute(&b, generated)
		if templateErr != nil {
			return result, templateErr
		}

		generated.Code = b.String()

		result[enum.Name] = generated.Code
	}

	for _, structItem := range parsed.Structs {
		goName := GenerateGoNameForType(structItem.Name)
		parseFunctionName := ParserFunction(goName)
		if ShouldGenerateStructType(goName) {
			generated := GeneratedStruct{
				GenerationParameters: GenerationParameters{
					OriginalName: structItem.Name,
					GoName:       goName,
				},
				ParserName: parseFunctionName,
				Definition: structItem,
				Code:       "",
			}

			var b bytes.Buffer
			templateErr := structTemplate.Execute(&b, generated)
			if templateErr != nil {
				return result, templateErr
			}

			generated.Code = b.String()

			result[structItem.Name] = generated.Code
		}
	}

	for _, event := range parsed.Events {
		if event.Kind == "struct" {
			goName := GenerateGoNameForType(event.Name)
			parseFunctionName := ParserFunction(goName)

			eventHash, hashErr := HashFromName(event.Name)
			if hashErr != nil {
				return result, hashErr
			}

			generated := GeneratedEvent{
				GenerationParameters: GenerationParameters{
					OriginalName: event.Name,
					GoName:       goName,
				},
				ParserName:   parseFunctionName,
				Definition:   event,
				EventNameVar: fmt.Sprintf("Event_%s", goName),
				EventHashVar: fmt.Sprintf("Hash_%s", goName),
				EventHash:    eventHash,
				Code:         "",
			}

			var b bytes.Buffer
			templateErr := eventTemplate.Execute(&b, generated)
			if templateErr != nil {
				return result, templateErr
			}

			generated.Code = b.String()

			result[event.Name] = generated.Code
		}
	}

	return result, nil
}

// Generates the header for the output code.
func GenerateHeader(packageName string) (string, error) {
	headerTemplate, headerTemplateParseErr := template.New("struct").Parse(HeaderTemplate)
	if headerTemplateParseErr != nil {
		return "", headerTemplateParseErr
	}

	parameters := HeaderParameters{
		Version:     version.SeerVersion,
		PackageName: packageName,
	}

	var b bytes.Buffer
	templateErr := headerTemplate.Execute(&b, parameters)
	if templateErr != nil {
		return "", templateErr
	}

	return b.String(), nil
}

// Generates a single string consisting of the Go code for all the artifacts in a parsed Starknet ABI.
func Generate(parsed *ParsedABI) (string, error) {
	snippets, snippetsErr := GenerateSnippets(parsed)
	if snippetsErr != nil {
		return "", snippetsErr
	}

	commonCode := strings.Join([]string{StructCommonCode, EventsCommonCode}, "\n\n")

	sections := make([]string, len(snippets))
	currentSection := 0
	for _, section := range snippets {
		sections[currentSection] = section
		currentSection++
	}

	snippetsCat := strings.Join(sections, "\n\n")

	return fmt.Sprintf("%s%s", commonCode, snippetsCat), nil
}

// This is the Go template which is used to generate the function corresponding to an Enum.
// This template should be applied to a GeneratedEnum struct.
var EnumTemplate string = `// ABI: {{.OriginalName}}

// {{.GoName}} is an alias for uint64
type {{.GoName}} = uint64

// {{.ParserName}} parses a {{.GoName}} from a list of felts. This function returns a tuple of:
// 1. The parsed {{.GoName}}
// 2. The number of field elements consumed in the parse
// 3. An error if the parse failed, nil otherwise
func {{.ParserName}} (parameters []*felt.Felt) ({{.GoName}}, int, error) {
	if len(parameters) < 1 {
		return 0, 0, ErrIncorrectParameters
	}
	return {{.GoName}}(parameters[0].Uint64()), 1, nil
}

// This function returns the string representation of a {{.GoName}} enum. This is the enum value from the ABI definition of the enum.
func {{.EvaluatorName}}(raw {{.GoName}}) string {
	switch raw {
	{{range .Definition.Variants}}case {{.Index}}:
		return "{{.Name}}"
	{{end -}}
	}
	return "UNKNOWN"
}`

var StructCommonCode string = `var ErrIncorrectParameters error = errors.New("incorrect parameters")

func ParseUint64(parameters []*felt.Felt) (uint64, int, error) {
	if len(parameters) < 1 {
		return 0, 0, ErrIncorrectParameters
	}
	return parameters[0].Uint64(), 1, nil
}

func ParseBigInt(parameters []*felt.Felt) (*big.Int, int, error) {
	if len(parameters) < 1 {
		return nil, 0, ErrIncorrectParameters
	}
	var result *big.Int
	result = parameters[0].BigInt(result)
	return result, 1, nil
}

func ParseString(parameters []*felt.Felt) (string, int, error) {
	if len(parameters) < 1 {
		return "", 0, ErrIncorrectParameters
	}
	return parameters[0].String(), 1, nil
}

func ParseArray[T any](parser func(parameters []*felt.Felt) (T, int, error)) func(parameters []*felt.Felt) ([]T, int, error) {
	return func (parameters []*felt.Felt) ([]T, int, error) {
		if len(parameters) < 1 {
			return nil, 0, ErrIncorrectParameters
		}

		arrayLengthRaw := parameters[0].Uint64()
		arrayLength := int(arrayLengthRaw)
		if len(parameters) < arrayLength + 1 {
			return nil, 0, ErrIncorrectParameters
		}

		result := make([]T, arrayLength)
		currentIndex := 0
		for i := 0; i < arrayLength; i++ {
			parsed, consumed, err := parser(parameters[currentIndex + 1:])
			if err != nil {
				return nil, 0, err
			}
			result[i] = parsed
			currentIndex += consumed
		}

		return result, currentIndex + 1, nil
	}
}
`

// This is the Go template which is used to generate the Go definition of a Starknet ABI struct.
// This template should be applied to a GeneratedStruct struct.
var StructTemplate string = `// ABI: {{.OriginalName}}

// {{.GoName}} is the Go struct corresponding to the {{.OriginalName}} struct.
type {{.GoName}} struct {
	{{range .Definition.Members}}
	{{(CamelCase .Name)}} {{(GenerateGoNameForType .Type)}}
	{{- end}}
}

// {{.ParserName}} parses a {{.GoName}} struct from a list of felts. This function returns a tuple of:
// 1. The parsed {{.GoName}} struct
// 2. The number of field elements consumed in the parse
// 3. An error if the parse failed, nil otherwise
func {{.ParserName}}(parameters []*felt.Felt) ({{.GoName}}, int, error) {
	currentIndex := 0
	result := {{.GoName}}{}

	{{range $index, $element := .Definition.Members}}
	value{{$index}}, consumed, err := {{(ParserFunction (GenerateGoNameForType .Type))}}(parameters[currentIndex:])
	if err != nil {
		return result, 0, err
	}
	result.{{(CamelCase .Name)}} = value{{$index}}
	currentIndex += consumed

	{{end}}

	return result, currentIndex + 1, nil
}
`

// Common code used in the code generated for events.
var EventsCommonCode string = `var ErrIncorrectEventKey error = errors.New("incorrect event key")

type RawEvent struct {
	BlockNumber     uint64
	BlockHash       *felt.Felt
	TransactionHash *felt.Felt
	FromAddress     *felt.Felt
	PrimaryKey      *felt.Felt
	Keys            []*felt.Felt
	Parameters      []*felt.Felt
}


`

// This is the Go template which is used to generate the Go bindings to a Starknet ABI event.
// This template should be applied to a GeneratedEvent struct.
var EventTemplate string = `
// ABI: {{.OriginalName}}

// ABI name for event
var {{.EventNameVar}} string = "{{.OriginalName}}"

// Starknet hash for the event, as it appears in Starknet event logs.
var {{.EventHashVar}} string = "{{.EventHash}}"

// {{.GoName}} is the Go struct corresponding to the {{.OriginalName}} event.
type {{.GoName}} struct {
	{{range .Definition.Members}}
	{{(CamelCase .Name)}} {{(GenerateGoNameForType .Type)}}
	{{- end}}
}
`

// This is the Go template used to create header information at the top of the generated code.
// At a bare minimum, the header specifies the version of seer that was used to generate the code.
// This template should be applied to a HeaderParameters struct.
var HeaderTemplate string = `// This file was generated by seer: https://github.com/moonstream-to/seer.
// seer version: {{.Version}}
// seer command: seer starknet abigentypes {{if .PackageName}}--package {{.PackageName}}{{end}}
// Warning: Edit at your own risk. Any edits you make will NOT survive the next code generation.

{{if .PackageName}}package {{.PackageName}}{{end}}

import (
	"errors"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
)
`
