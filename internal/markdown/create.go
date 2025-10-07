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
// The showAllDecorators parameter controls whether to include all decorator columns in the output.
// Returns an error if any operation fails.
func CreateFile(filename string, template *types.Template, verbose bool, sections []types.Section, showAllDecorators bool) error { //nolint:gocyclo // This function is complex by design.
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
	var builder strings.Builder
	builder.Grow(estimateMarkdownSize(template, sections))
	if err := buildMarkdownString(&builder, template, sections, showAllDecorators); err != nil {
		return fmt.Errorf("failed to build Markdown string: %w", err)
	}
	markdownString := builder.String()

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

func buildMarkdownString(builder *strings.Builder, template *types.Template, sections []types.Section, showAllDecorators bool) error {
	// Template metadata
	var title *string
	if template.Metadata == nil || template.Metadata.Name == nil || *template.Metadata.Name == "" {
		title = &template.FileName
	} else {
		title = template.Metadata.Name
	}
	fmt.Fprintf(builder, "# %s\n\n", *title)

	// Create a mapping between the section enum and the corresponding markdown function
	// Each function will be called in turn to generate a specific section of the markdown file.
	// The order of the functions in the slice determines the order of the sections in the markdown file.
	sectionMarkdownFunctions := map[types.Section]func(*types.Template, bool) (string, error){
		types.DescriptionSection: func(t *types.Template, _ bool) (string, error) {
			return generateDescriptionSection(t)
		},
		types.UsageSection: func(t *types.Template, _ bool) (string, error) {
			return generateUsageSection(t)
		},
		types.ModulesSection: func(t *types.Template, _ bool) (string, error) {
			return generateModulesSection(t)
		},
		types.ResourcesSection: func(t *types.Template, _ bool) (string, error) {
			return generateResourcesSection(t)
		},
		types.ParametersSection:           generateParametersSection,
		types.UserDefinedDataTypesSection: generateUserDefinedDataTypesSection,
		types.UserDefinedFunctionsSection: generateUserDefinedFunctionsSection,
		types.VariablesSection: func(t *types.Template, _ bool) (string, error) {
			return generateVariablesSection(t)
		},
		types.OutputsSection: generateOutputsSection,
	}

	// Iterate over the sections slice and call the corresponding function for each section
	for _, section := range sections {
		if function, ok := sectionMarkdownFunctions[section]; ok {
			sectionContent, err := function(template, showAllDecorators)
			if err != nil {
				return err
			}
			builder.WriteString(sectionContent)
			if sectionContent != "" {
				builder.WriteString("\n")
			}
		} else {
			return fmt.Errorf("invalid section: %s", section)
		}
	}

	// Trim trailing newlines and add a single newline at the end
	result := strings.TrimRight(builder.String(), "\n") + "\n"
	builder.Reset()
	builder.WriteString(result)

	return nil
}

// estimateMarkdownSize estimates the size of the markdown file based on the template content and sections.
//
//nolint:mnd,gocyclo // This function does some estimations based on the template content.
func estimateMarkdownSize(template *types.Template, sections []types.Section) int {
	baseSize := 500 // Start with a base size for the title and basic structure

	for _, section := range sections {
		switch section {
		case types.DescriptionSection:
			if template.Metadata != nil && template.Metadata.Description != nil {
				baseSize += len(*template.Metadata.Description) + 50 // Add description length plus some overhead
			}
		case types.UsageSection:
			baseSize += 150 + (len(template.Parameters) * 20) // Estimate based on number of parameters
		case types.ModulesSection:
			baseSize += len(template.Modules) * 100 // Estimate 100 characters per module
		case types.ResourcesSection:
			baseSize += len(template.Resources) * 150 // Estimate 150 characters per resource
		case types.ParametersSection:
			baseSize += len(template.Parameters) * 50 // Estimate 50 characters per parameter
		case types.UserDefinedDataTypesSection:
			baseSize += len(template.UserDefinedDataTypes) * 50 // Estimate 50 characters per user-defined type
		case types.UserDefinedFunctionsSection:
			baseSize += len(template.UserDefinedFunctions) * 40 // Estimate 40 characters per user-defined function
		case types.VariablesSection:
			baseSize += len(template.Variables) * 40 // Estimate 40 characters per variable
		case types.OutputsSection:
			baseSize += len(template.Outputs) * 50 // Estimate 50 characters per output
		}
	}

	return baseSize
}
