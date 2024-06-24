/*
Package markdown provides functionality to create a Markdown file from a Bicep template.
*/
package markdown

import (
	"fmt"
	"os"
	"strings"

	"github.com/christosgalano/bicep-docs/internal/types"
)

// CreateFile creates or updates a file with the specified filename using the provided template.
// If the file already exists and its content matches the generated Markdown string, no changes are made.
// If the file does not exist or its content differs from the generated Markdown string, the file is created or updated accordingly.
// The verbose parameter controls whether informational messages are printed to stdout.
// The sections parameter specifies the sections to include in the generated Markdown string.
// Returns an error if any operation fails.
func CreateFile(filename string, template *types.Template, verbose bool, sections []types.Section) error { //nolint:gocyclo // This function is complex by design.
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
	markdownString, err := buildMarkdownString(template, sections)
	if err != nil {
		return fmt.Errorf("failed to build Markdown string: %w", err)
	}

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

// buildMarkdownString generates a markdown string based on the provided template and sections.
// It constructs the markdown file by appending different sections in a specific order.
// The function returns the generated markdown string and an error if any occurred.
func buildMarkdownString(template *types.Template, sections []types.Section) (string, error) {
	var builder strings.Builder

	// Template metadata
	var title *string
	if template.Metadata == nil || template.Metadata.Name == nil || *template.Metadata.Name == "" {
		title = &template.FileName
	} else {
		title = template.Metadata.Name
	}
	builder.WriteString(fmt.Sprintf("# %s\n\n", *title))

	// Create a mapping between the section enum and the corresponding markdown function
	// Each function will be called in turn to generate a specific section of the markdown file.
	// The order of the functions in the slice determines the order of the sections in the markdown file.
	sectionMarkdownFunctions := map[types.Section]func(*types.Template) (string, error){
		types.DescriptionSection:          generateDescriptionSection,
		types.UsageSection:                generateUsageSection,
		types.ModulesSection:              generateModulesSection,
		types.ResourcesSection:            generateResourcesSection,
		types.ParametersSection:           generateParametersSection,
		types.UserDefinedDataTypesSection: generateUserDefinedDataTypesSection,
		types.UserDefinedFunctionsSection: generateUserDefinedFunctionsSection,
		types.VariablesSection:            generateVariablesSection,
		types.OutputsSection:              generateOutputsSection,
	}

	// Iterate over the sections slice and call the corresponding function for each section
	for _, section := range sections {
		if function, ok := sectionMarkdownFunctions[section]; ok {
			if markdownString, err := function(template); err == nil {
				builder.WriteString(markdownString)
				if markdownString != "" {
					builder.WriteString("\n")
				}
			} else {
				return "", err
			}
		} else {
			return "", fmt.Errorf("invalid section: %s", section)
		}
	}

	// Trim trailing newlines and add a single newline at the end
	result := strings.TrimRight(builder.String(), "\n") + "\n"

	return result, nil
}
