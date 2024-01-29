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
	"sort"
	"strings"

	"github.com/christosgalano/bicep-docs/internal/types"
)

// Regular expressions to parse Bicep templates.
var (
	moduleRegex                    = regexp.MustCompile(`^module\s+(\S+)\s+'(\S+)'`)
	resourceRegex                  = regexp.MustCompile(`^resource\s+(\S+)\s+'(\S+)'`)
	typeRegex                      = regexp.MustCompile(`^type\s+(\S+)\s+`)
	outputRegex                    = regexp.MustCompile(`^output\s+(\S+)\s+`)
	variableRegex                  = regexp.MustCompile(`^var\s+(\S+)\s+`)
	parameterRegex                 = regexp.MustCompile(`^param\s+(\S+)\s+`)
	inlineDescriptionRegex         = regexp.MustCompile(`^@(description|sys.description)\(('''|')(.*?)('''|')\)`)
	multilineDescriptionStartRegex = regexp.MustCompile(`^@(description|sys.description)\('''(.*)`)
)

// ParseTemplates parses a Bicep and its corresponding ARM template.
// It returns a Template struct that contains the information about the Bicep template.
func ParseTemplates(bicepFile, armFile string) (*types.Template, error) {
	var err error
	var template types.Template
	template.FileName = bicepFile

	// Parse Bicep template
	template.Modules, template.Resources, err = parseBicepTemplate(bicepFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Bicep modules: %w", err)
	}

	// Parse ARM template
	err = parseArmTemplate(armFile, &template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ARM template: %w", err)
	}

	return &template, nil
}

// parseArmTemplate decodes an ARM template into a Template struct.
func parseArmTemplate(armFile string, template *types.Template) error {
	// Open JSON file
	file, err := os.Open(armFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode ARM template into Template struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&template)
	if err != nil {
		return err
	}

	return nil
}

// parseBicepTemplate extracts information about any existing modules or resources from a Bicep template.
func parseBicepTemplate(bicepFile string) ([]types.Module, []types.Resource, error) {
	file, err := os.Open(bicepFile)
	if err != nil {
		return []types.Module{}, []types.Resource{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	modules := []types.Module{}
	resources := []types.Resource{}

	var description *string
	var line, currentDescription string
	for scanner.Scan() {
		line = scanner.Text()

		// Skip empty line
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Skip comment
		skipped, err := skipComment(line, scanner)
		if err != nil {
			return nil, nil, err
		}
		if skipped {
			continue
		}

		// Parse description
		description = parseDescription(line, scanner)
		if description != nil {
			currentDescription = *description
			continue
		}

		// Ignore the description of parameters, outputs, types, and variables
		if ignoreDescription(line) {
			currentDescription = ""
			continue
		}

		// Parse module
		module := parseModule(line)
		if module != nil {
			module.Description = currentDescription
			modules = append(modules, *module)
			currentDescription = ""
			continue
		}

		// Parse resource
		resource := parseResource(line)
		if resource != nil {
			resource.Description = currentDescription
			resources = append(resources, *resource)
			currentDescription = ""
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	// Sort the resource symbolic names
	sort.SliceStable(resources, func(i, j int) bool {
		return resources[i].SymbolicName < resources[j].SymbolicName
	})

	// Sort the module symbolic names
	sort.SliceStable(modules, func(i, j int) bool {
		return modules[i].SymbolicName < modules[j].SymbolicName
	})

	return modules, resources, err
}

// parseDescription extracts the description of a module or resource from a line.
func parseDescription(line string, scanner *bufio.Scanner) *string {
	// Parse inline description
	matches := inlineDescriptionRegex.FindStringSubmatch(line)
	if matches != nil {
		return &matches[3]
	}

	// Parse multiline description
	matches = multilineDescriptionStartRegex.FindStringSubmatch(line)
	if matches != nil {
		description := matches[2]
		afterMultilineTicks := false

		// Consume lines until the multiline description ends => '''[\n\r\s]*\).
		// If the line contains the multiline description's end ticks ('''), add the text before them.
		// If the line ends with a closing parenthesis, then the multiline description is over;
		// otherwise, keep consuming lines until a closing parenthesis is found.
		for scanner.Scan() {
			line = scanner.Text()
			switch {
			case strings.Contains(line, "'''") && strings.HasSuffix(line, ")"):
				description += strings.Split(line, "'''")[0]
				return &description
			case strings.Contains(line, "'''"):
				description += strings.Split(line, "'''")[0]
				afterMultilineTicks = true
			case afterMultilineTicks && strings.HasSuffix(line, ")"):
				return &description
			default:
				description += strings.TrimSpace(line)
			}
		}
	}
	return nil
}

// parseModule extracts information about a module from a line.
func parseModule(line string) *types.Module {
	matches := moduleRegex.FindStringSubmatch(line)
	if matches != nil {
		moduleSource := strings.ReplaceAll(matches[2], "'", "")
		return &types.Module{
			SymbolicName: matches[1],
			Source:       moduleSource,
		}
	}
	return nil
}

// parseResource extracts information about a resource from a line.
func parseResource(line string) *types.Resource {
	matches := resourceRegex.FindStringSubmatch(line)
	if matches != nil {
		resourceType := strings.Split(matches[2], "@")[0]
		resourceType = strings.ReplaceAll(resourceType, "'", "")
		return &types.Resource{
			SymbolicName: matches[1],
			Type:         resourceType,
		}
	}
	return nil
}

// skipComment skips single line and multiline comments.
// It returns true if a comment was skipped, false otherwise.
func skipComment(line string, scanner *bufio.Scanner) (bool, error) {
	// Skip single line comments
	if strings.HasPrefix(strings.TrimSpace(line), "//") {
		return true, nil
	}

	// Skip multiline comments
	if strings.HasPrefix(strings.TrimSpace(line), "/*") {
		for scanner.Scan() {
			line = scanner.Text()
			if strings.HasSuffix(strings.TrimSpace(line), "*/") {
				break
			}
		}

		// If we've reached here without breaking, the comment was not properly closed
		if scanner.Err() != nil {
			return false, fmt.Errorf("multiline comment was not closed")
		}

		return true, nil
	}

	return false, nil
}

// ignoreDescription returns true if the provided line defines a type, output, variable, or parameter.
// That is, if the line starts with "type", "output", "var", or "param".
// Then, the description of the type, output, variable, or parameter should be ignored.
func ignoreDescription(line string) bool {
	matchType := typeRegex.FindStringSubmatch(line)
	matchOutput := outputRegex.FindStringSubmatch(line)
	matchVariable := variableRegex.FindStringSubmatch(line)
	matchParameter := parameterRegex.FindStringSubmatch(line)
	return matchType != nil || matchOutput != nil || matchVariable != nil || matchParameter != nil
}
