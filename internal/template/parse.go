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
	funcRegex                      = regexp.MustCompile(`^func\s+(\S+)\s*\(`)
	outputRegex                    = regexp.MustCompile(`^output\s+(\S+)\s+`)
	variableRegex                  = regexp.MustCompile(`^var\s+(\S+)\s+`)
	parameterRegex                 = regexp.MustCompile(`^param\s+(\S+)\s+`)
	inlineDescriptionRegex         = regexp.MustCompile(`^@(description|sys.description)\(('''|')(.*?)('''|')\)`)
	multilineDescriptionStartRegex = regexp.MustCompile(`^@(description|sys.description)\('''(.*)`)
	conditionStartRegex            = regexp.MustCompile(`[=:]\s*if\s*\(`)
	retryOnStartRegex              = regexp.MustCompile(`^@(sys\.)?retryOn\(`)
	onlyIfNotExistsRegex           = regexp.MustCompile(`^@(sys\.)?onlyIfNotExists\(\s*\)`)
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

			// Apply descriptions to matching ARM template variables.
			// Empty Bicep descriptions are skipped so that descriptions sourced
			// from the ARM template (e.g. exported variables) are preserved.
			for i := range template.Variables {
				if desc, ok := varMap[template.Variables[i].Name]; ok && desc != "" {
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
	var line string
	var pending pendingDeclaration
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
			pending.description = *description
			continue
		}

		// Parse resource decorators (@retryOn, @onlyIfNotExists)
		if pending.applyDecorator(line) {
			continue
		}

		// Ignore the description of parameters, outputs, types, and variables
		if ignoreDescription(line) {
			pending = pendingDeclaration{}
			continue
		}

		// Parse module, resource, and variable declarations
		if parseDeclarationLine(line, pending, &modules, &resources, &variables) {
			pending = pendingDeclaration{}
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

// pendingDeclaration accumulates the description and decorators seen above a declaration line.
// The state is applied to the next module, resource, or variable declaration and then reset.
type pendingDeclaration struct {
	description     string
	retryOn         string
	onlyIfNotExists bool
}

// applyDecorator updates the pending state if the line is a supported resource decorator.
// It reports whether the line was consumed as a decorator.
func (p *pendingDeclaration) applyDecorator(line string) bool {
	if retryOn := parseRetryOn(line); retryOn != "" {
		p.retryOn = retryOn
		return true
	}
	if onlyIfNotExistsRegex.MatchString(line) {
		p.onlyIfNotExists = true
		return true
	}
	return false
}

// parseDeclarationLine parses a module, resource, or variable declaration line,
// applies the pending description and decorators, and appends the result to the corresponding slice.
// It reports whether the line matched a declaration.
func parseDeclarationLine(line string, pending pendingDeclaration, modules *[]types.Module, resources *[]types.Resource, variables *[]types.Variable) bool {
	if module := parseModule(line); module != nil {
		module.Description = pending.description
		*modules = append(*modules, *module)
		return true
	}
	if resource := parseResource(line); resource != nil {
		resource.Description = pending.description
		resource.RetryOn = pending.retryOn
		resource.OnlyIfNotExists = pending.onlyIfNotExists
		*resources = append(*resources, *resource)
		return true
	}
	if variable := parseVariable(line); variable != nil {
		variable.Description = pending.description
		*variables = append(*variables, *variable)
		return true
	}
	return false
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
			Condition:    parseCondition(line),
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
			Condition:    parseCondition(line),
		}
	}
	return nil
}

// parseCondition extracts the condition expression of a conditional deployment
// from a declaration line. It supports both the plain form (resource ... = if (condition) {...})
// and the loop form (resource ... = [for item in items: if (condition) {...]).
// It returns the expression inside the balanced parentheses following "= if" or ": if".
// If the line has no condition, or the condition spans multiple lines, it returns an empty string.
func parseCondition(line string) string {
	loc := conditionStartRegex.FindStringIndex(line)
	if loc == nil {
		return ""
	}
	return extractBalancedParens(line, loc[1])
}

// parseRetryOn extracts the arguments of the @retryOn decorator from a line.
// It returns the arguments inside the balanced parentheses (e.g. "['ServerError'], 3").
// If the line is not a @retryOn decorator, or the arguments span multiple lines,
// it returns an empty string.
func parseRetryOn(line string) string {
	loc := retryOnStartRegex.FindStringIndex(line)
	if loc == nil {
		return ""
	}
	return extractBalancedParens(line, loc[1])
}

// extractBalancedParens returns the trimmed substring of line from start (an index
// just after an opening parenthesis) up to the matching closing parenthesis,
// ignoring parentheses that appear inside single-quoted strings (escaped quotes \' are handled).
// Quotes nested inside string interpolation (e.g. '${func('x')}') are not tracked and may
// cause the expression to be truncated. If no matching closing parenthesis is found on the
// line, it returns an empty string.
func extractBalancedParens(line string, start int) string {
	depth := 1
	inString := false
	for i := start; i < len(line); i++ {
		switch line[i] {
		case '\'':
			if i == 0 || line[i-1] != '\\' {
				inString = !inString
			}
		case '(':
			if !inString {
				depth++
			}
		case ')':
			if !inString {
				depth--
				if depth == 0 {
					return strings.TrimSpace(line[start:i])
				}
			}
		}
	}
	return ""
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
// It returns true if the line matches any of the type, function, output, or parameter patterns;
// otherwise, it returns false. This prevents descriptions (and decorators) of declarations
// that are parsed from the ARM template from leaking onto the next Bicep-parsed declaration.
func ignoreDescription(line string) bool {
	matchType := typeRegex.FindStringSubmatch(line)
	matchFunc := funcRegex.FindStringSubmatch(line)
	matchOutput := outputRegex.FindStringSubmatch(line)
	matchParameter := parameterRegex.FindStringSubmatch(line)
	return matchType != nil || matchFunc != nil || matchOutput != nil || matchParameter != nil
}
