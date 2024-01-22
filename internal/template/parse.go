/*
Package template provides functions to build and parse Bicep and the corresponding ARM templates.
*/
package template

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/christosgalano/bicep-docs/internal/types"
)

// ParseTemplates parses a Bicep and its corresponding ARM template.
// It returns a Template struct that contains the information about the Bicep template.
func ParseTemplates(bicepFile, armFile string) (*types.Template, error) {
	var err error
	var template types.Template
	template.FileName = bicepFile

	// Parse Bicep modules
	template.Modules, template.Resources, err = parseBicepTemplate(bicepFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Bicep modules: %w", err)
	}

	// Open JSON file
	file, err := os.Open(armFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode ARM template into Template struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&template)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ARM template: %w", err)
	}

	return &template, nil
}

// parseBicepTemplate extracts information about any existing modules or resources from a Bicep template.
func parseBicepTemplate(bicepFile string) ([]types.Module, []types.Resource, error) {
	file, err := os.Open(bicepFile)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	modules := []types.Module{}
	resources := []types.Resource{}

	// Regex for Bicep entities
	outputRegex := regexp.MustCompile(`^output\s+(\S+)\s+`)
	parameterRegex := regexp.MustCompile(`^param\s+(\S+)\s+`)
	moduleRegex := regexp.MustCompile(`^module\s+(\S+)\s+'(\S+)'`)
	resourceRegex := regexp.MustCompile(`^resource\s+(\S+)\s+'(\S+)'`)

	// Regex for annotations
	inlineDescriptionRegex := regexp.MustCompile(`^@(description|sys.description)\(('''|')(.*?)('''|')\)`)
	multilineDescriptionStartRegex := regexp.MustCompile(`^@(description|sys.description)\('''(.*)`)
	// multilineDescriptionEndRegex := regexp.MustCompile(`(.*?)'''\s*\)$`)

	currentDescription := ""
	insideMultilineComment := false
	insideMultilineDescription := false
	afterMultilineTicks := false

	// For each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Skip single line comments
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}

		// Skip multiline comments
		if strings.HasPrefix(strings.TrimSpace(line), "/*") {
			insideMultilineComment = true
		}
		if insideMultilineComment {
			if strings.HasSuffix(strings.TrimSpace(line), "*/") {
				insideMultilineComment = false
			}
			continue
		}

		// Handle inline descriptions
		matches := inlineDescriptionRegex.FindStringSubmatch(line)
    	if matches != nil {
			currentDescription = matches[3]
			continue
		}

		// Handle multiline descriptions
		matches = multilineDescriptionStartRegex.FindStringSubmatch(line)
		if matches != nil {
			insideMultilineDescription = true
			currentDescription = matches[2]
			continue
		}
		if insideMultilineDescription {
			// Consume lines until the multiline description ends '''[\n\r\s]*\).
			// If the line contains the multiline description's end ticks ('''), add the text before them.
			// If the line ends with a closing parenthesis, then the multiline description is over;
			// otherwise, keep consuming lines until a closing parenthesis is found.
			if strings.Contains(line, "'''") {
				currentDescription += strings.Split(line, "'''")[0]
				if !strings.HasSuffix(line, ")") {
					afterMultilineTicks = true
				} else {
					insideMultilineDescription = false
				}
			} else if afterMultilineTicks && strings.HasSuffix(line, ")") {
				afterMultilineTicks = false
				insideMultilineDescription = false
			} else {
				currentDescription += strings.TrimSpace(line)
			}
			continue
		}

		// Handle parameters and outputs
		matchesP, matchesO := parameterRegex.FindStringSubmatch(line), outputRegex.FindStringSubmatch(line)
		if matchesP != nil || matchesO != nil {
			currentDescription = ""
			continue
		}

		// Handle modules
		matches = moduleRegex.FindStringSubmatch(line)
		if matches != nil {
			moduleSource := strings.ReplaceAll(matches[2], "'", "")
			modules = append(modules, types.Module{
				SymbolicName: matches[1],
				Source:       moduleSource,
				Description:  currentDescription,
			})
			currentDescription = ""
			continue
		}

		// Handle resources
		matches = resourceRegex.FindStringSubmatch(line)
		if matches != nil {
			resourceType := strings.Split(matches[2], "@")[0]
			resourceType = strings.ReplaceAll(resourceType, "'", "")
			resources = append(resources, types.Resource{
				SymbolicName: matches[1],
				Type:         resourceType,
				Description:  currentDescription,
			})
			currentDescription = ""
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return modules, resources, err
}
