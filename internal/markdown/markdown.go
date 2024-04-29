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
// It returns an error if any operation fails.
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
	markdownString, err := buildMarkdownString(template)
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

// buildMarkdownString takes a pointer to a Template and builds a markdown string
// representation of the template. It returns the markdown string and an error, if any.
func buildMarkdownString(template *types.Template) (string, error) {
	var builder strings.Builder

	// Template metadata
	var title *string
	if template.Metadata == nil || template.Metadata.Name == nil || *template.Metadata.Name == "" {
		title = &template.FileName
	} else {
		title = template.Metadata.Name
	}
	builder.WriteString(fmt.Sprintf("# %s\n", *title))
	if template.Metadata != nil && template.Metadata.Description != nil && *template.Metadata.Description != "" {
		builder.WriteString(fmt.Sprintf("\n## Description\n\n%s\n", *template.Metadata.Description))
	}
	builder.WriteString("\n")

	// Create a slice of functions that each return a markdown table string and an error.
	markdownFunctions := []func(*types.Template) (string, error){
		generateUsageSection,
		modulesToMarkdownTable,
		resourcesToMarkdownTable,
		parametersToMarkdownTable,
		userDefinedDataTypesToMarkdownTable,
		userDefinedFunctionsToMarkdownTable,
		variablesToMarkdownTable,
		outputsToMarkdownTable,
	}

	// Iterate over the slice and call each function in turn.
	for _, function := range markdownFunctions {
		if md, err := function(template); err == nil {
			builder.WriteString(md)
			if md != "" {
				builder.WriteString("\n")
			}
		} else {
			return "", err
		}
	}

	// Trim trailing newlines and add a single newline at the end.
	result := strings.TrimRight(builder.String(), "\n") + "\n"

	return result, nil
}
