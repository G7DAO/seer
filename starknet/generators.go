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
	ParseFunctionName string
	Definition        *Enum
	Code              string
}

// The output of the code generation process for struct items in a Starknet ABI.
type GeneratedStruct struct {
	GenerationParameters
	Definition   *Struct
	ParamsLength int
	Code         string
}

type GeneratedEvent struct {
	GenerationParameters
	Definition       *EventStruct
	EventNameVar     string
	EventHashVar     string
	EventHash        string
	ParamFeltLengths map[string]int
	Code             string
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
		parseFunctionName := fmt.Sprintf("Parse%s", goName)

		generated := GeneratedEnum{
			GenerationParameters: GenerationParameters{
				OriginalName: enum.Name,
				GoName:       goName,
			},
			ParseFunctionName: parseFunctionName,
			Definition:        enum,
			Code:              "",
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
		if ShouldGenerateStructType(goName) {
			generated := GeneratedStruct{
				GenerationParameters: GenerationParameters{
					OriginalName: structItem.Name,
					GoName:       goName,
				},
				Definition:   structItem,
				ParamsLength: len(structItem.Members),
				Code:         "",
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

			eventHash, hashErr := HashFromName(event.Name)
			if hashErr != nil {
				return result, hashErr
			}

			generated := GeneratedEvent{
				GenerationParameters: GenerationParameters{
					OriginalName: event.Name,
					GoName:       goName,
				},
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

	commonCode := ""

	if len(parsed.Events) > 0 {
		commonCode = EventsCommonCode
	}

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

// {{.GoName}} is an alias for string
type {{.GoName}} = string

// This function maps a Felt corresponding to the index of an enum variant to the name of that variant.
func {{.ParseFunctionName}}(parameter *felt.Felt) {{.GoName}} {
	parameterInt := parameter.Uint64()
	switch parameterInt {
	{{range .Definition.Variants}}case {{.Index}}:
		return "{{.Name}}"
	{{end}}
	}
	return "UNKNOWN"
}`

// This is the Go template which is used to generate the Go definition of a Starknet ABI struct.
// This template should be applied to a GeneratedStruct struct.
var StructTemplate string = `// ABI: {{.OriginalName}}

// {{.GoName}} is the Go struct corresponding to the {{.OriginalName}} struct.
type {{.GoName}} struct {
	{{range .Definition.Members}}
	{{(CamelCase .Name)}} {{(GenerateGoNameForType .Type)}}
	{{- end}}
}
`

// Common code used in the code generated for events.
var EventsCommonCode string = `var ErrIncorrectEventKey error = errors.New("incorrect event key")
var ErrIncorrectParameters error = errors.New("incorrect parameters")

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
