package markdown

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/christosgalano/bicep-docs/internal/types"
)

// MarkdownTable represents a table in a Markdown document.
// It contains a title, headers, and rows.
type MarkdownTable struct {
	Title   string
	Headers []string
	Rows    [][]string
}

// NewMarkdownTable creates a new MarkdownTable with the specified title, headers, and rows.
func NewMarkdownTable(title string, headers []string, rows [][]string) *MarkdownTable {
	return &MarkdownTable{Title: title, Headers: headers, Rows: rows}
}

// String returns the string representation of the MarkdownTable.
// It generates the table headers and rows and returns them as a single string.
func (table *MarkdownTable) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("## %s\n\n", table.Title))
	builder.WriteString(generateTableHeaders(table.Headers))
	for _, row := range table.Rows {
		builder.WriteString(generateTableRow(row))
	}
	return builder.String()
}

// generateTableHeaders generates the markdown table headers based on the given slice of headers.
// It returns a string containing the markdown table headers.
func generateTableHeaders(headers []string) string {
	var builder strings.Builder
	for _, header := range headers {
		builder.WriteString(fmt.Sprintf("| %s ", header))
	}
	builder.WriteString("|\n")
	for i := 0; i < len(headers); i++ {
		builder.WriteString("| --- ")
	}
	builder.WriteString("|\n")
	return builder.String()
}

// generateTableRow generates a markdown table row based on the given slice of strings.
// Each string in the slice represents a column in the table row.
// The generated table row is returned as a string.
func generateTableRow(row []string) string {
	var builder strings.Builder
	for _, column := range row {
		builder.WriteString(fmt.Sprintf("| %s ", column))
	}
	builder.WriteString("|\n")
	return builder.String()
}

// modulesToMarkdownTable converts a template's modules into a markdown table.
// If the template has no modules, it returns an empty string.
// The table headers are "Symbolic Name", "Source", and "Description".
// If an error occurs, it is returned along with an empty string.
func modulesToMarkdownTable(template *types.Template) (string, error) {
	if len(template.Modules) == 0 {
		return "", nil
	}
	headers := []string{"Symbolic Name", "Source", "Description"}
	rows := make([][]string, len(template.Modules))
	for i, module := range template.Modules {
		description := strings.ReplaceAll(module.Description, "\r\n", "\n")
		description = strings.ReplaceAll(description, "\n", "<br>")
		rows[i] = []string{module.SymbolicName, module.Source, description}
	}
	return NewMarkdownTable("Modules", headers, rows).String(), nil
}

// resourcesToMarkdownTable converts a template's resources into a markdown table.
// If the template has no resources, it returns an empty string.
// The table headers are "Symbolic Name", "Type", and "Description".
// If an error occurs, it is returned along with an empty string.
func resourcesToMarkdownTable(template *types.Template) (string, error) {
	if len(template.Resources) == 0 {
		return "", nil
	}
	headers := []string{"Symbolic Name", "Type", "Description"}
	rows := make([][]string, len(template.Resources))
	for i, resource := range template.Resources {
		typeLink := fmt.Sprintf("[%s](https://learn.microsoft.com/en-us/azure/templates/%s)", resource.Type, strings.ToLower(resource.Type))
		description := strings.ReplaceAll(resource.Description, "\r\n", "\n")
		description = strings.ReplaceAll(description, "\n", "<br>")
		rows[i] = []string{resource.SymbolicName, typeLink, description}
	}
	return NewMarkdownTable("Resources", headers, rows).String(), nil
}

// parametersToMarkdownTable converts a template's parameters into a markdown table.
// If the template has no parameters, it returns an empty string.
// The table headers are "Name", "Type", "Description", and "Default".
// If an error occurs, it is returned along with an empty string.
func parametersToMarkdownTable(template *types.Template) (string, error) {
	if len(template.Parameters) == 0 {
		return "", nil
	}
	re := regexp.MustCompile(`([^ ]):([^ ])|([^ ]),([^ ])`)
	headers := []string{"Name", "Type", "Description", "Default"}
	rows := make([][]string, len(template.Parameters))
	for i, parameter := range template.Parameters {
		defaultValue := ""
		if parameter.DefaultValue != nil {
			jsonValue, err := json.Marshal(parameter.DefaultValue)
			if err != nil {
				return "", fmt.Errorf("failed to marshal default value: %w", err)
			}
			defaultValue = string(jsonValue)
			defaultValue = re.ReplaceAllStringFunc(defaultValue, func(s string) string {
				if strings.Contains(s, ":") {
					return strings.Replace(s, ":", ": ", 1)
				} else if strings.Contains(s, ",") {
					return strings.Replace(s, ",", ", ", 1)
				}
				return s
			})
			defaultValue = strings.ReplaceAll(defaultValue, "\r\n", "\n")
			defaultValue = strings.ReplaceAll(defaultValue, "\n", "<br>")
		}
		rows[i] = []string{parameter.Name, extractType(parameter.Type), extractDescription(parameter.Metadata), defaultValue}
	}
	return NewMarkdownTable("Parameters", headers, rows).String(), nil
}

// outputsToMarkdownTable converts a template's outputs into a markdown table.
// If the template has no outputs, it returns an empty string.
// The table headers are "Name", "Type", and "Description".
// If an error occurs, it is returned along with an empty string.
func outputsToMarkdownTable(template *types.Template) (string, error) {
	if len(template.Outputs) == 0 {
		return "", nil
	}
	headers := []string{"Name", "Type", "Description"}
	rows := make([][]string, len(template.Outputs))
	for i, output := range template.Outputs {
		rows[i] = []string{output.Name, extractType(output.Type), extractDescription(output.Metadata)}
	}
	return NewMarkdownTable("Outputs", headers, rows).String(), nil
}

// userDefinedDataTypesToMarkdownTable converts a template's user-defined data types into a markdown table.
// If the template has no user-defined data types, it returns an empty string.
// The table headers are "Name", "Type", and "Description".
// If an error occurs, it is returned along with an empty string.
func userDefinedDataTypesToMarkdownTable(template *types.Template) (string, error) {
	if len(template.UserDefinedDataTypes) == 0 {
		return "", nil
	}
	headers := []string{"Name", "Type", "Description"}
	rows := make([][]string, len(template.UserDefinedDataTypes))
	for i, dataType := range template.UserDefinedDataTypes {
		rows[i] = []string{dataType.Name, extractType(dataType.Type), extractDescription(dataType.Metadata)}
	}
	return NewMarkdownTable("User Defined Data Types (UDDTs)", headers, rows).String(), nil
}

// userDefinedFunctionsToMarkdownTable converts a template's user-defined functions into a markdown table.
// If the template has no user-defined functions, it returns an empty string.
// The table headers are "Name" and "Description".
// If an error occurs, it is returned along with an empty string.
func userDefinedFunctionsToMarkdownTable(template *types.Template) (string, error) {
	if len(template.UserDefinedFunctions) == 0 {
		return "", nil
	}
	headers := []string{"Name", "Description"}
	rows := make([][]string, len(template.UserDefinedFunctions))
	for i, function := range template.UserDefinedFunctions {
		rows[i] = []string{function.Name, extractDescription(function.Metadata)}
	}
	return NewMarkdownTable("User Defined Functions (UDFs)", headers, rows).String(), nil
}

// variablesToMarkdownTable converts a template's variables into a markdown table.
// If the template has no variables, it returns an empty string.
// The table headers are "Name" and "Description".
// If an error occurs, it is returned along with an empty string.
func variablesToMarkdownTable(template *types.Template) (string, error) {
	if len(template.Variables) == 0 {
		return "", nil
	}
	headers := []string{"Name", "Description"}
	rows := make([][]string, len(template.Variables))
	for i, variable := range template.Variables {
		rows[i] = []string{variable.Name, variable.Description}
	}
	return NewMarkdownTable("Variables", headers, rows).String(), nil
}

// generateUsageSection generates the "Usage" section for the Bicep module documentation.
// It returns a string containing the generated section.
// If an error occurs, it is returned along with an empty string.
func generateUsageSection(template *types.Template) (string, error) {
	var builder strings.Builder

	builder.WriteString("## Usage\n\n")
	builder.WriteString("Here is a basic example of how to use this Bicep module:\n\n")
	builder.WriteString("```bicep\n")
	builder.WriteString("module reference_name 'path_to_module | container_registry_reference' = {\n")
	builder.WriteString("  name: 'deployment_name'\n")
	builder.WriteString("  params: {\n")

	// Required parameters (without a default value).
	builder.WriteString("    // Required parameters\n")
	for _, parameter := range template.Parameters {
		if parameter.DefaultValue == nil {
			builder.WriteString(fmt.Sprintf("    %s:\n", parameter.Name))
		}
	}

	// Optional parameters (with a default value).
	builder.WriteString("\n    // Optional parameters\n")
	for _, parameter := range template.Parameters {
		if parameter.DefaultValue == nil {
			continue
		}
		jsonValue, err := json.MarshalIndent(parameter.DefaultValue, "    ", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal default value: %w", err)
		}
		defaultValue := string(jsonValue)

		// Remove quotes from keys.
		re := regexp.MustCompile(`"(\w+)":`)
		defaultValue = re.ReplaceAllString(defaultValue, "$1:")

		// Replace double quotes with single quotes.
		defaultValue = strings.ReplaceAll(defaultValue, "\"", "'")

		// Remove commas.
		defaultValue = strings.ReplaceAll(defaultValue, ",", "")

		builder.WriteString(fmt.Sprintf("    %s: %s\n", parameter.Name, defaultValue))
	}
	builder.WriteString("  }\n")
	builder.WriteString("}\n")
	builder.WriteString("```\n")
	builder.WriteString("\n> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.\n") //nolint:lll // Ignore long line length.

	return builder.String(), nil
}
