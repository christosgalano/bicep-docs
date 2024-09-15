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
			name:     "directory_input",
			input:    "./testdata",
			output:   "",
			verbose:  true,
			sections: []types.Section{types.DescriptionSection, types.UsageSection},
			expected: "",
		},
		{
			name:     "file_input",
			input:    "./testdata/main.bicep",
			output:   "README.md",
			verbose:  false,
			sections: []types.Section{types.ModulesSection, types.ParametersSection},
			expected: "",
		},
		{
			name:     "non_existent_input",
			input:    "./path/to/non-existent",
			output:   "",
			verbose:  true,
			sections: []types.Section{types.DescriptionSection, types.UsageSection},
			expected: "no such file or directory \"./path/to/non-existent\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateDocs(tt.input, tt.output, tt.verbose, tt.sections)
			if tt.expected != "" {
				if err == nil {
					t.Errorf("GenerateDocs() expected error but got none")
				} else if !strings.Contains(err.Error(), tt.expected) {
					t.Errorf("GenerateDocs() error = %v, expected to contain = %s", err, tt.expected)
				}
			} else if err != nil {
				t.Errorf("GenerateDocs() unexpected error = %v", err)
			}
		})
	}
}

func Test_generateDocsFromDirectory(t *testing.T) {
	tests := []struct {
		name     string
		dirPath  string
		verbose  bool
		sections []types.Section
		expected string
	}{
		{
			name:     "valid_directory",
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

func Test_generateDocsFromBicepFile(t *testing.T) {
	tests := []struct {
		name         string
		bicepFile    string
		markdownFile string
		verbose      bool
		sections     []types.Section
		expected     string
	}{
		{
			name:         "valid_file",
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
