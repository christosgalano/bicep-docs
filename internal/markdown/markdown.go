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

// CreateFile creates or overwrites a Markdown file
// with the information from a Bicep template.
func CreateFile(filename string, template *types.Template) error {
	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Build Markdown string
	var builder strings.Builder
	builder.WriteString(templateMetadataToMarkdown(template))
	if len(template.Modules) > 0 {
		builder.WriteString("\n")
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
	if len(template.Outputs) > 0 {
		builder.WriteString(outputsToMarkdown(template))
	}
	markdownString := builder.String()

	// Write to file
	_, err = file.WriteString(markdownString)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
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

	if template.Metadata != nil && template.Metadata.Description != nil && *template.Metadata.Description != "" {
		builder.WriteString(fmt.Sprintf("\n## Description\n\n%s\n", *template.Metadata.Description))
	}
	return builder.String()
}

// modulesToMarkdown converts the modules part of a template to Markdown.
func modulesToMarkdown(template *types.Template) string {
	var builder strings.Builder
	moduleHeaders := []string{"Symbolic Name", "Source", "Description"}
	if len(template.Modules) > 0 {
		builder.WriteString("## Modules\n\n")
		builder.WriteString(generateTableHeaders(moduleHeaders))
		for _, module := range template.Modules {
			sourceLink := fmt.Sprintf("[%s](%s)", module.Source, module.Source)
			builder.WriteString(fmt.Sprintf("| %s | %s | %s |\n", module.SymbolicName, sourceLink, module.Description))
		}
	}
	return builder.String()
}

// resourcesToMarkdown converts the resources part of a template to Markdown.
func resourcesToMarkdown(template *types.Template) string {
	var builder strings.Builder
	resourceHeaders := []string{"Symbolic Name", "Type", "Description"}
	if len(template.Resources) > 0 {
		builder.WriteString("## Resources\n\n")
		builder.WriteString(generateTableHeaders(resourceHeaders))
		for _, resource := range template.Resources {
			typeLink := fmt.Sprintf("[%s](https://learn.microsoft.com/en-us/azure/templates/%s)", resource.Type, strings.ToLower(resource.Type))
			builder.WriteString(fmt.Sprintf("| %s | %s | %s |\n", resource.SymbolicName, typeLink, resource.Description))
		}
	}
	return builder.String()
}

// parametersToMarkdown converts the parameters part of a template to Markdown.
func parametersToMarkdown(template *types.Template) string {
	var builder strings.Builder
	parameterHeaders := []string{"Name", "Type", "Description", "Default"}
	if len(template.Parameters) > 0 {
		builder.WriteString("## Parameters\n\n")
		builder.WriteString(generateTableHeaders(parameterHeaders))
		for name, param := range template.Parameters {
			defaultValue := ""
			if param.DefaultValue != nil {
				if dv, ok := param.DefaultValue.(map[string]any); ok {
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
				}
			}
			builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", name, param.Type, extractDescription(param.Metadata), defaultValue))
		}
	}
	return builder.String()
}

// outputsToMarkdown converts the outputs part of a template to Markdown.
func outputsToMarkdown(template *types.Template) string {
	var builder strings.Builder
	outputHeaders := []string{"Name", "Type", "Description"}
	if len(template.Outputs) > 0 {
		builder.WriteString("## Outputs\n\n")
		builder.WriteString(generateTableHeaders(outputHeaders))
		for name, output := range template.Outputs {
			builder.WriteString(fmt.Sprintf("| %s | %s | %s |\n", name, output.Type, extractDescription(output.Metadata)))
		}
	}
	return builder.String()
}

// extractDescription extracts the description from the metadata part of an entity.
func extractDescription(metadata *types.Metadata) string {
	description := ""
	if metadata != nil && metadata.Description != nil {
		description = *metadata.Description
	}
	return description
}
