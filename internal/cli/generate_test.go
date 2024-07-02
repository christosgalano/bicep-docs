package cli

import (
	"strings"
	"testing"

	"github.com/christosgalano/bicep-docs/internal/types"
)

func TestGenerateDocs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		output   string
		verbose  bool
		sections []types.Section
		expected string
	}{
		{
			name:     "directory input",
			input:    "./testdata",
			output:   "",
			verbose:  true,
			sections: []types.Section{types.DescriptionSection, types.UsageSection},
			expected: "",
		},
		{
			name:     "Bicep file input",
			input:    "./testdata/main.bicep",
			output:   "README.md",
			verbose:  false,
			sections: []types.Section{types.ModulesSection, types.ParametersSection},
			expected: "",
		},
		{
			name:     "non-existent input",
			input:    "./path/to/non-existent",
			output:   "",
			verbose:  true,
			sections: []types.Section{types.DescriptionSection, types.UsageSection},
			expected: "no such file or directory \"./path/to/non-existent\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generateDocs(tt.input, tt.output, tt.verbose, tt.sections)
			if tt.expected != "" {
				if err == nil {
					t.Errorf("generateDocs() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.expected) {
					t.Errorf("generateDocs() error = %v, expected to contain = %s", err, tt.expected)
				}
			} else if err != nil {
				t.Errorf("generateDocs() unexpected error = %v", err)
			}
		})
	}
}

func TestGenerateDocsFromDirectory(t *testing.T) {
	tests := []struct {
		name     string
		dirPath  string
		verbose  bool
		sections []types.Section
		expected string
	}{
		{
			name:     "valid directory",
			dirPath:  "./testdata",
			verbose:  true,
			sections: []types.Section{types.DescriptionSection, types.UsageSection},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generateDocsFromDirectory(tt.dirPath, tt.verbose, tt.sections)
			if tt.expected != "" {
				if err == nil {
					t.Errorf("generateDocsFromDirectory() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.expected) {
					t.Errorf("generateDocsFromDirectory() error = %v, expected to contain = %s", err, tt.expected)
				}
			} else if err != nil {
				t.Errorf("generateDocsFromDirectory() unexpected error = %v", err)
			}
		})
	}
}

func TestGenerateDocsFromBicepFile(t *testing.T) {
	tests := []struct {
		name         string
		bicepFile    string
		markdownFile string
		verbose      bool
		sections     []types.Section
		expected     string
	}{
		{
			name:         "valid Bicep file",
			bicepFile:    "./testdata/main.bicep",
			markdownFile: "./testdata/README.md",
			verbose:      true,
			sections:     []types.Section{types.ModulesSection, types.ParametersSection},
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generateDocsFromBicepFile(tt.bicepFile, tt.markdownFile, tt.verbose, tt.sections)
			if tt.expected != "" {
				if err == nil {
					t.Errorf("generateDocsFromBicepFile() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.expected) {
					t.Errorf("generateDocsFromBicepFile() error = %v, expected to contain = %s", err, tt.expected)
				}
			} else if err != nil {
				t.Errorf("generateDocsFromBicepFile() unexpected error = %v", err)
			}
		})
	}
}
