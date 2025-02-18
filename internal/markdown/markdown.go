package markdown

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/christosgalano/bicep-docs/internal/types"
)

// HeaderType represents the type of header in a Markdown document.
type HeaderType int

const (
	H1 HeaderType = iota + 1 // H1 represents a level 1 header (#)
	H2                       // H2 represents a level 2 header (##)
	H3                       // H3 represents a level 3 header (###)
)

// String returns the string representation of the HeaderType.
func (h HeaderType) String() string {
	switch h {
	case H1:
		return "#"
	case H2:
		return "##" //nolint:goconst // Ignore the duplication of the return value.
	case H3:
		return "###"
	default:
		return "##" // Default to H2 if invalid
	}
}

// MarkdownTable represents a table in a Markdown document.
// It contains a title, header type, headers, and rows.
type MarkdownTable struct {
	Title      string
	HeaderType HeaderType
	Headers    []string
	Rows       [][]string
}

// NewMarkdownTable creates a new MarkdownTable with the specified title, header type, headers, and rows.
// If headerType is invalid, it defaults to H2.
// If headers is nil or empty, it returns nil.
// If rows is nil, it defaults to an empty slice.
func NewMarkdownTable(title string, headerType HeaderType, headers []string, rows [][]string) *MarkdownTable {
	if len(headers) == 0 {
		return nil
	}

	if headerType < H1 || headerType > H3 {
		headerType = H2
	}

	if rows == nil {
		rows = [][]string{}
	}

	return &MarkdownTable{
		Title:      title,
		HeaderType: headerType,
		Headers:    headers,
		Rows:       rows,
	}
}

// String returns the string representation of the MarkdownTable.
// It generates the table headers and rows and returns them as a single string.
func (table *MarkdownTable) String() string {
	var builder strings.Builder

	// Add title with proper header level
	builder.WriteString(fmt.Sprintf("%s %s\n\n", table.HeaderType.String(), table.Title))

	// Add table headers
	builder.WriteString(generateTableHeaders(table.Headers))

	// Add table rows
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

// generateModulesSection converts a template's modules into a markdown table.
// If the template has no modules, it returns an empty string.
// The table headers are "Symbolic Name", "Source", and "Description".
// If an error occurs, it is returned along with an empty string.
func generateModulesSection(template *types.Template) (string, error) { //nolint:unparam // Ignore the error return value; it is there for consistency.
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
	return NewMarkdownTable("Modules", H2, headers, rows).String(), nil
}

// generateResourcesSection converts a template's resources into a markdown table.
// If the template has no resources, it returns an empty string.
// The table headers are "Symbolic Name", "Type", and "Description".
// If an error occurs, it is returned along with an empty string.
func generateResourcesSection(template *types.Template) (string, error) { //nolint:unparam // Ignore the error return value; it is there for consistency.
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
	return NewMarkdownTable("Resources", H2, headers, rows).String(), nil
}

// generateParametersSection generates the parameters section of a template in markdown format.
// It takes a pointer to a types.Template as input and returns the generated markdown string and an error, if any.
// If the template has no parameters, it returns an empty string and a nil error.
func generateParametersSection(template *types.Template) (string, error) {
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
	return NewMarkdownTable("Parameters", H2, headers, rows).String(), nil
}

// generateOutputsSection generates the outputs section of the template markdown.
// It takes a pointer to a types.Template and returns a string representation of the outputs section and an error, if any.
// If the template has no outputs, it returns an empty string and no error.
func generateOutputsSection(template *types.Template) (string, error) { //nolint:unparam // Ignore the error return value; it is there for consistency.
	if len(template.Outputs) == 0 {
		return "", nil
	}
	headers := []string{"Name", "Type", "Description"}
	rows := make([][]string, len(template.Outputs))
	for i, output := range template.Outputs {
		rows[i] = []string{output.Name, extractType(output.Type), extractDescription(output.Metadata)}
	}
	return NewMarkdownTable("Outputs", H2, headers, rows).String(), nil
}

// generateUserDefinedDataTypesSection generates a markdown table section for user-defined data types (UDDTs) based on the provided template.
// If there are no user-defined data types in the template, an empty string is returned.
// The table includes columns for Name, Type, and Description.
// Each row in the table represents a user-defined data type, with the corresponding values extracted from the template.
// The function returns the generated markdown table as a string and any error encountered during the process.
func generateUserDefinedDataTypesSection(template *types.Template) (string, error) { //nolint:unparam // Ignore the error return value; it is there for consistency.
	if len(template.UserDefinedDataTypes) == 0 {
		return "", nil
	}

	// Main table
	headers := []string{"Name", "Type", "Description", "Properties"}
	rows := make([][]string, len(template.UserDefinedDataTypes))
	for i, dataType := range template.UserDefinedDataTypes {
		propertiesColumn := ""
		if len(dataType.Properties) > 0 {
			propertiesColumn = fmt.Sprintf("[View Properties](#%s)", strings.ToLower(dataType.Name))
		}

		rows[i] = []string{
			dataType.Name,
			extractType(dataType.Type),
			extractDescription(dataType.Metadata),
			propertiesColumn,
		}
	}
	table := NewMarkdownTable("User Defined Data Types (UDDTs)", H2, headers, rows).String()

	// Sub-tables for properties
	propertyHeaders := []string{"Name", "Type", "Description"}
	for _, dataType := range template.UserDefinedDataTypes {
		if len(dataType.Properties) == 0 {
			continue
		}
		propertyRows := make([][]string, len(dataType.Properties))
		for i, property := range dataType.Properties {
			propertyRows[i] = []string{
				property.Name,
				extractType(property.Type),
				extractDescription(property.Metadata),
			}
		}
		table += "\n" + NewMarkdownTable(dataType.Name, H3, propertyHeaders, propertyRows).String()
	}

	return table, nil
}

// generateUserDefinedFunctionsSection converts a template's user-defined functions into a markdown table.
// If the template has no user-defined functions, it returns an empty string.
// The table headers are "Name" and "Description".
// If an error occurs, it is returned along with an empty string.
func generateUserDefinedFunctionsSection(template *types.Template) (string, error) { //nolint:unparam // Ignore the error return value; it is there for consistency.
	if len(template.UserDefinedFunctions) == 0 {
		return "", nil
	}
	headers := []string{"Name", "Description"}
	rows := make([][]string, len(template.UserDefinedFunctions))
	for i, function := range template.UserDefinedFunctions {
		rows[i] = []string{function.Name, extractDescription(function.Metadata)}
	}
	return NewMarkdownTable("User Defined Functions (UDFs)", H2, headers, rows).String(), nil
}

// generateVariablesSection generates the variables section of the markdown document based on the provided template.
// If the template has no variables, it returns an empty string.
// Otherwise, it creates a markdown table with the variable names and descriptions.
func generateVariablesSection(template *types.Template) (string, error) { //nolint:unparam // Ignore the error return value; it is there for consistency.
	if len(template.Variables) == 0 {
		return "", nil
	}
	headers := []string{"Name", "Description"}
	rows := make([][]string, len(template.Variables))
	for i, variable := range template.Variables {
		rows[i] = []string{variable.Name, variable.Description}
	}
	return NewMarkdownTable("Variables", H2, headers, rows).String(), nil
}

// generateUsageSection generates the usage section for the Bicep module.
// It takes a pointer to a types.Template object as input and returns a string containing the generated usage section.
// The usage section includes a basic example of how to use the Bicep module, including both required and optional parameters.
// The default values for the optional parameters are also included in the example.
// Note: The default values may contain function calls or references enclosed in square brackets.
// The function returns an error if there is any issue with marshaling the default values.
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

// generateDescriptionSection generates the description section of a template.
// It takes a pointer to a types.Template as input and returns the generated description section as a string.
// If the template has a non-empty description in its metadata, it will be included in the generated section.
// Otherwise, an empty string will be returned.
func generateDescriptionSection(template *types.Template) (string, error) { //nolint:unparam // Ignore the error return value; it is there for consistency.
	var builder strings.Builder
	if template.Metadata != nil && template.Metadata.Description != nil && *template.Metadata.Description != "" {
		builder.WriteString(fmt.Sprintf("## Description\n\n%s\n", *template.Metadata.Description))
	} else {
		builder.WriteString("")
	}
	return builder.String(), nil
}

// checkFileExists checks if a file exists and is not a directory.
// It returns true if the file exists, false otherwise, along with any error encountered.
func checkFileExists(filename string) (bool, error) {
	f, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to stat file %q: %w", filename, err)
	}
	if f.IsDir() {
		return false, fmt.Errorf("output %q is a directory", filename)
	}
	return true, nil
}

// readFileContent reads the content of a file and returns it as a string.
// It takes a filename as input and returns the file content and any error encountered.
func readFileContent(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", filename, err)
	}
	return string(bytes), nil
}

// extractType extracts the type from a type string.
// If the type is a user defined data type, it returns the name of it.
func extractType(t string) string {
	if strings.HasPrefix(t, "#/definitions/") {
		split := strings.Split(t, "/")
		return split[len(split)-1] + " (uddt)"
	}
	return t
}

// extractDescription extracts the description from the given metadata and returns it.
// If the metadata or the description is nil, an empty string is returned.
func extractDescription(metadata *types.Metadata) string {
	description := ""
	if metadata != nil && metadata.Description != nil {
		description = strings.ReplaceAll(*metadata.Description, "\r\n", "<br>")
		description = strings.ReplaceAll(description, "\n", "<br>")
	}
	return description
}
