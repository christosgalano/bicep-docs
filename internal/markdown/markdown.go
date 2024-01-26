/*
Package markdown provides functionality to create a Markdown file from a Bicep template.
*/
package markdown

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/christosgalano/bicep-docs/internal/types"
)

// CreateFile creates or overwrites a Markdown file with the information from a Bicep template.
//
// If verbose is true, it prints a message to stdout.
func CreateFile(filename string, template *types.Template, verbose bool) error {
	// Check if template is nil
	if template == nil {
		return fmt.Errorf("invalid template (nil)")
	}

	// Check if file exists and is not a directory
	fileExists, err := checkFileExists(filename)
	if err != nil {
		return err
	}

	// Read file content if it exists
	fileContent := ""
	if fileExists {
		fileContent, err = readFileContent(filename)
		if err != nil {
			return err
		}
	}

	// Build Markdown string
	markdownString := buildMarkdownString(template)

	// Check if file needs to be updated
	if fileExists && fileContent == markdownString {
		if verbose {
			fmt.Printf("No changes to %s\n", filename)
		}
		return nil
	}

	// Create/Truncate file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write to file the Markdown string
	_, err = file.WriteString(markdownString)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	// Print message to stdout
	if verbose {
		if fileExists {
			fmt.Printf("Updated %s\n", filename)
		} else {
			fmt.Printf("Created %s\n", filename)
		}
	}

	return nil
}

// checkFileExists checks if a file exists and is not a directory.
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

// readFileContent reads the content of a file.
func readFileContent(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", filename, err)
	}
	return string(bytes), nil
}

// buildMarkdownString builds the Markdown string from a Bicep template.
func buildMarkdownString(template *types.Template) string {
	var builder strings.Builder
	builder.WriteString(templateMetadataToMarkdown(template))
	if len(template.Modules) > 0 {
		builder.WriteString(modulesToMarkdown(template))
		builder.WriteString("\n")
	}
	if len(template.Resources) > 0 {
		builder.WriteString(resourcesToMarkdown(template))
		builder.WriteString("\n")
	}
	if len(template.Parameters) > 0 {
		builder.WriteString(parametersToMarkdown(template))
		builder.WriteString("\n")
	}
	if len(template.UserDefinedDataTypes) > 0 {
		builder.WriteString(userDefinedDataTypesToMarkdown(template))
		builder.WriteString("\n")
	}
	if len(template.UserDefinedFunctions) > 0 {
		builder.WriteString(userDefinedFunctionsToMarkdown(template))
		builder.WriteString("\n")
	}
	if len(template.Variables) > 0 {
		builder.WriteString(variablesToMarkdown(template))
		builder.WriteString("\n")
	}
	if len(template.Outputs) > 0 {
		builder.WriteString(outputsToMarkdown(template))
	}
	return builder.String()
}

// generateTableHeaders generates the Markdown table headers.
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

// templateMetadataToMarkdown converts the metadata part of a template to Markdown.
func templateMetadataToMarkdown(template *types.Template) string {
	var builder strings.Builder
	var title *string
	if template.Metadata == nil || template.Metadata.Name == nil || *template.Metadata.Name == "" {
		title = &template.FileName
	} else {
		title = template.Metadata.Name
	}
	builder.WriteString(fmt.Sprintf("# %s\n", *title))

	// Add description if it exists
	if template.Metadata != nil && template.Metadata.Description != nil && *template.Metadata.Description != "" {
		builder.WriteString(fmt.Sprintf("\n## Description\n\n%s\n", *template.Metadata.Description))
	}

	// Add a newline if there are any modules, resources, parameters or outputs
	if len(template.Modules)+len(template.Resources)+len(template.Parameters)+len(template.Outputs) > 0 {
		builder.WriteString("\n")
	}

	return builder.String()
}

// modulesToMarkdown converts the modules a template to Markdown.
func modulesToMarkdown(template *types.Template) string {
	var builder strings.Builder
	moduleHeaders := []string{"Symbolic Name", "Source", "Description"}
	if len(template.Modules) > 0 {
		builder.WriteString("## Modules\n\n")
		builder.WriteString(generateTableHeaders(moduleHeaders))
		for _, module := range template.Modules {
			builder.WriteString(
				fmt.Sprintf("| %s | %s | %s |\n",
					module.SymbolicName,
					module.Source,
					strings.ReplaceAll(module.Description, "\n", "<br>"),
				),
			)
		}
	}
	return builder.String()
}

// resourcesToMarkdown converts the resources a template to Markdown.
func resourcesToMarkdown(template *types.Template) string {
	var builder strings.Builder
	resourceHeaders := []string{"Symbolic Name", "Type", "Description"}
	if len(template.Resources) > 0 {
		builder.WriteString("## Resources\n\n")
		builder.WriteString(generateTableHeaders(resourceHeaders))
		for _, resource := range template.Resources {
			typeLink := fmt.Sprintf("[%s](https://learn.microsoft.com/en-us/azure/templates/%s)", resource.Type, strings.ToLower(resource.Type))
			builder.WriteString(
				fmt.Sprintf("| %s | %s | %s |\n",
					resource.SymbolicName,
					typeLink,
					strings.ReplaceAll(resource.Description, "\n", "<br>"),
				),
			)
		}
	}
	return builder.String()
}

// parametersToMarkdown converts the parameters a template to Markdown.
func parametersToMarkdown(template *types.Template) string {
	var builder strings.Builder
	parameterHeaders := []string{"Name", "Type", "Description", "Default"}
	if len(template.Parameters) > 0 {
		builder.WriteString("## Parameters\n\n")
		builder.WriteString(generateTableHeaders(parameterHeaders))
		for _, parameter := range template.Parameters {
			defaultValue := ""
			if parameter.DefaultValue != nil {
				if dv, ok := parameter.DefaultValue.(map[string]any); ok {
					if len(dv) == 0 {
						defaultValue = "{}"
					} else {
						defaultValue = "{"
						i := 0
						for k, v := range dv {
							if i == len(dv)-1 {
								defaultValue += fmt.Sprintf("%s: %v}", k, v)
								break
							}
							defaultValue += fmt.Sprintf("%s: %v, ", k, v)
							i++
						}
					}
				} else {
					defaultValue = fmt.Sprintf("%v", parameter.DefaultValue)
				}
			}
			builder.WriteString(
				fmt.Sprintf("| %s | %s | %s | %s |\n",
					parameter.Name,
					extractType(parameter.Type),
					extractDescription(parameter.Metadata),
					defaultValue,
				),
			)
		}
	}
	return builder.String()
}

// outputsToMarkdown converts the outputs a template to Markdown.
func outputsToMarkdown(template *types.Template) string {
	var builder strings.Builder
	outputHeaders := []string{"Name", "Type", "Description"}
	if len(template.Outputs) > 0 {
		builder.WriteString("## Outputs\n\n")
		builder.WriteString(generateTableHeaders(outputHeaders))
		for _, output := range template.Outputs {
			builder.WriteString(
				fmt.Sprintf("| %s | %s | %s |\n",
					output.Name,
					extractType(output.Type),
					extractDescription(output.Metadata),
				),
			)
		}
	}
	return builder.String()
}

// userDefinedDataTypesToMarkdown converts the user defined data types a template to Markdown.
func userDefinedDataTypesToMarkdown(template *types.Template) string {
	var builder strings.Builder
	userDefinedDataTypeHeaders := []string{"Name", "Type", "Description"}
	if len(template.UserDefinedDataTypes) > 0 {
		builder.WriteString("## User Defined Data Types (UDDTs)\n\n")
		builder.WriteString(generateTableHeaders(userDefinedDataTypeHeaders))
		for _, userDefinedDataType := range template.UserDefinedDataTypes {
			builder.WriteString(
				fmt.Sprintf("| %s | %s | %s |\n",
					userDefinedDataType.Name,
					extractType(userDefinedDataType.Type),
					extractDescription(userDefinedDataType.Metadata),
				),
			)
		}
	}
	return builder.String()
}

// userDefinedFunctionsToMarkdown converts the user defined functions a template to Markdown.
func userDefinedFunctionsToMarkdown(template *types.Template) string {
	var builder strings.Builder
	userDefinedFunctionHeaders := []string{"Name", "Description"}
	if len(template.UserDefinedFunctions) > 0 {
		builder.WriteString("## User Defined Functions (UDFs)\n\n")
		builder.WriteString(generateTableHeaders(userDefinedFunctionHeaders))
		for _, userDefinedFunction := range template.UserDefinedFunctions {
			builder.WriteString(
				fmt.Sprintf("| %s | %s |\n",
					userDefinedFunction.Name,
					extractDescription(userDefinedFunction.Metadata),
				),
			)
		}
	}
	return builder.String()
}

// variablesToMarkdown converts the variables a template to Markdown.
func variablesToMarkdown(template *types.Template) string {
	var builder strings.Builder
	variableHeaders := []string{"Name"}
	if len(template.Variables) > 0 {
		builder.WriteString("## Variables\n\n")
		builder.WriteString(generateTableHeaders(variableHeaders))
		for _, variable := range template.Variables {
			builder.WriteString(fmt.Sprintf("| %s |\n", variable.Name))
		}
	}
	return builder.String()
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

// extractDescription extracts the description from an entity's metadata.
func extractDescription(metadata *types.Metadata) string {
	description := ""
	if metadata != nil && metadata.Description != nil {
		description = strings.ReplaceAll(*metadata.Description, "\n", "<br>")
	}
	return description
}
