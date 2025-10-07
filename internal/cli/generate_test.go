package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/christosgalano/bicep-docs/internal/types"
)

func TestGenerateDocs(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		output            string
		verbose           bool
		sections          []types.Section
		showAllDecorators bool
		expected          string
	}{
		{
			name:              "directory_input",
			input:             "./testdata",
			output:            "",
			verbose:           true,
			sections:          []types.Section{types.DescriptionSection, types.ParametersSection, types.VariablesSection},
			showAllDecorators: false,
			expected:          "",
		},
		{
			name:              "file_input",
			input:             "./testdata/main.bicep",
			output:            "./testdata/README.md",
			verbose:           false,
			sections:          []types.Section{types.ModulesSection, types.ParametersSection},
			showAllDecorators: false,
			expected:          "",
		},
		{
			name:              "non_existent_input",
			input:             "./path/to/non-existent",
			output:            "",
			verbose:           true,
			sections:          []types.Section{types.DescriptionSection, types.ParametersSection, types.VariablesSection},
			showAllDecorators: false,
			expected:          "no such file or directory \"./path/to/non-existent\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateDocs(tt.input, tt.output, tt.verbose, tt.sections, tt.showAllDecorators)
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
		name              string
		dirPath           string
		verbose           bool
		sections          []types.Section
		showAllDecorators bool
		expected          string
	}{
		{
			name:              "valid_directory",
			dirPath:           "./testdata",
			verbose:           true,
			sections:          []types.Section{types.DescriptionSection, types.ParametersSection, types.VariablesSection},
			showAllDecorators: false,
			expected:          "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generateDocsFromDirectory(tt.dirPath, tt.verbose, tt.sections, tt.showAllDecorators)
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
		name              string
		bicepFile         string
		markdownFile      string
		verbose           bool
		sections          []types.Section
		showAllDecorators bool
		expected          string
	}{
		{
			name:              "valid_file",
			bicepFile:         "./testdata/main.bicep",
			markdownFile:      "./testdata/README.md",
			verbose:           true,
			sections:          []types.Section{types.ModulesSection, types.ParametersSection},
			showAllDecorators: false,
			expected:          "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generateDocsFromBicepFile(tt.bicepFile, tt.markdownFile, tt.verbose, tt.sections, tt.showAllDecorators)
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

// createTestDirectory creates a temporary directory with the specified number of main.bicep files.
func createTestDirectory(numFiles int) (string, error) {
	tempDir, err := os.MkdirTemp("", "bicep-docs-benchmark")
	if err != nil {
		return "", err
	}

	for i := 0; i < numFiles; i++ {
		dirPath := filepath.Join(tempDir, fmt.Sprintf("dir_%d", i))
		err := os.MkdirAll(dirPath, 0o755)
		if err != nil {
			os.RemoveAll(tempDir)
			return "", err
		}

		content := fmt.Sprintf(`metadata name = 'test %d'
metadata description = 'This is test template %d.'

@sys.description('This is a test parameter.')
param test_parameter string = 'test'

@sys.description('This is a test variable.')
var test_variable = '${test_parameter}'`, i, i)

		filePath := filepath.Join(dirPath, "main.bicep")
		err = os.WriteFile(filePath, []byte(content), 0o600)
		if err != nil {
			os.RemoveAll(tempDir)
			return "", err
		}
	}

	return tempDir, nil
}

func BenchmarkGenerateDocs(b *testing.B) {
	fileCounts := []int{50, 100, 200}
	sections := []types.Section{types.DescriptionSection, types.ParametersSection, types.VariablesSection}

	for _, count := range fileCounts {
		b.Run(fmt.Sprintf("Files-%d", count), func(b *testing.B) {
			tempDir, err := createTestDirectory(count)
			if err != nil {
				b.Fatalf("Failed to create test directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := GenerateDocs(tempDir, "", false, sections, false)
				if err != nil {
					b.Fatalf("GenerateDocs() failed: %v", err)
				}
			}
		})
	}
}
