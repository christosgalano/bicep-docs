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

// ParseTemplates parses the Bicep and ARM templates and returns a populated types.Template struct.
// It takes the paths to the Bicep file and ARM file as input parameters.
// The function returns a pointer to the types.Template struct and an error, if any.
func ParseTemplates(bicepFile, armFile string) (*types.Template, error) {
	var err error
	var template types.Template
	template.FileName = bicepFile

	// Parse Bicep template
	var variables []types.Variable
	template.Modules, template.Resources, variables, err = parseBicepTemplate(bicepFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Bicep modules: %w", err)
	}

	// Parse ARM template
	err = parseArmTemplate(armFile, &template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ARM template: %w", err)
	}

	// Handle variables that might be optimized away in ARM template
	if len(variables) > 0 {
		// If we found variables in Bicep but none in ARM, use the Bicep ones
		if len(template.Variables) == 0 {
			template.Variables = variables
		} else {
			// If we have variables in both, try to match descriptions
			varMap := make(map[string]string)
			for _, v := range variables {
				varMap[v.Name] = v.Description
			}

			// Apply descriptions to matching ARM template variables
			for i := range template.Variables {
				if desc, ok := varMap[template.Variables[i].Name]; ok {
					template.Variables[i].Description = desc
				}
			}
		}
	}

	return &template, nil
}

// parseArmTemplate parses the specified ARM template file and populates the provided template struct.
// It opens the JSON file, decodes the ARM template into the template struct, and returns any errors encountered.
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

// parseBicepTemplate parses a Bicep template file and extracts the modules, resources, and variables defined in the file.
// It returns the parsed modules, resources, variables, and any error encountered during parsing.
func parseBicepTemplate(bicepFile string) ([]types.Module, []types.Resource, []types.Variable, error) {
	file, err := os.Open(bicepFile)
	if err != nil {
		return []types.Module{}, []types.Resource{}, []types.Variable{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	modules := []types.Module{}
	resources := []types.Resource{}
	variables := []types.Variable{}

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
			return nil, nil, nil, err
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

		// Parse variable
		variable := parseVariable(line)
		if variable != nil {
			variable.Description = currentDescription
			variables = append(variables, *variable)
			currentDescription = ""
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, nil, err
	}

	// Sort the resource symbolic names
	sort.SliceStable(resources, func(i, j int) bool {
		return resources[i].SymbolicName < resources[j].SymbolicName
	})

	// Sort the module symbolic names
	sort.SliceStable(modules, func(i, j int) bool {
		return modules[i].SymbolicName < modules[j].SymbolicName
	})

	return modules, resources, variables, err
}

// parseDescription parses a line of text and returns a pointer to the description.
// It supports both inline and multiline descriptions.
// If the line does not match the regex pattern, it returns nil.
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

// parseModule parses a line of text and returns a pointer to a types.Module struct.
// If the line does not match the regex pattern, it returns nil.
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

// parseResource parses a line of text and returns a pointer to a types.Resource struct.
// If the line does not match the regex pattern, it returns nil.
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

// parseVariable parses a line of text and returns a pointer to a types.Variable struct.
// If the line does not match the regex pattern, it returns nil.
func parseVariable(line string) *types.Variable {
	matches := variableRegex.FindStringSubmatch(line)
	if matches != nil {
		return &types.Variable{
			Name: matches[1],
		}
	}
	return nil
}

// skipComment checks if the given line is a comment and skips it.
// It supports both single-line and multi-line comments.
// If the comment is multi-line, it continues scanning until the closing "*/" is found.
// If the comment is not properly closed, it returns an error.
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

// ignoreDescription checks if a given line should be ignored based on certain patterns.
// It returns true if the line matches any of the type, output, variable or parameter patterns; otherwise, it returns false.
func ignoreDescription(line string) bool {
	matchType := typeRegex.FindStringSubmatch(line)
	matchOutput := outputRegex.FindStringSubmatch(line)
	matchParameter := parameterRegex.FindStringSubmatch(line)
	return matchType != nil || matchOutput != nil || matchParameter != nil
}
