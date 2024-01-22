package markdown

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/christosgalano/bicep-docs/internal/types"
)

func TestCreateFile(t *testing.T) {
	templateName := "test"
	templateDescription := "This is a test template."
	parameterDescription := "This is a test parameter."
	outputDescription := "This is a test output."

	type args struct {
		filename string
		template *types.Template
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		checkFile string
	}{
		{
			name: "full template",
			args: args{
				filename: "full_template.md",
				template: &types.Template{
					FileName:  "test.bicep",
					Modules:   []types.Module{{SymbolicName: "test_module", Source: "./modules/test_module/main.bicep", Description: "This is a test module."}},
					Resources: []types.Resource{{SymbolicName: "test_resource", Type: "Microsoft.Storage/storageAccounts", Description: "This is a test resource."}},
					Parameters: map[string]types.Parameter{
						"test_parameter1": {
							Type:         "string",
							DefaultValue: "test",
							Metadata: &types.Metadata{
								Description: &parameterDescription,
							},
						},
						"test_parameter2": {
							Type:         "object",
							DefaultValue: map[string]any{},
							Metadata: &types.Metadata{
								Description: &parameterDescription,
							},
						},
						"test_parameter3": {
							Type:         "object",
							DefaultValue: map[string]any{"key1": "value1", "key2": "value2"},
							Metadata: &types.Metadata{
								Description: &parameterDescription,
							},
						},
					},
					Outputs: map[string]types.Output{
						"test_output": {
							Type: "string",
							Metadata: &types.Metadata{
								Description: &outputDescription,
							},
						},
					},
					Metadata: &types.Metadata{
						Name:        &templateName,
						Description: &templateDescription,
					},
				},
			},
			wantErr:   false,
			checkFile: "./testdata/full_template.md",
		},
		{
			name: "no name",
			args: args{
				filename: "no_name.md",
				template: &types.Template{
					FileName: "test.bicep",
					Metadata: &types.Metadata{
						Description: &templateDescription,
					},
				},
			},
			wantErr:   false,
			checkFile: "./testdata/no_name.md",
		},
		{
			name: "no description",
			args: args{
				filename: "no_description.md",
				template: &types.Template{
					FileName: "test.bicep",
					Metadata: &types.Metadata{
						Name: &templateName,
					},
				},
			},
			wantErr:   false,
			checkFile: "./testdata/no_description.md",
		},
		{
			name: "no metadata",
			args: args{
				filename: "no_metadata.md",
				template: &types.Template{
					FileName: "test.bicep",
				},
			},
			wantErr:   false,
			checkFile: "./testdata/no_metadata.md",
		},
		{
			name: "given path is a directory",
			args: args{
				filename: "testdata",
				template: nil,
			},
			wantErr: true,
		},
		{
			name: "nil template",
			args: args{
				filename: "nil_template.md",
				template: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for the output files
			tempDir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			// Call CreateFile with the filename in the temporary directory
			filename := filepath.Join(tempDir, tt.args.filename)
			if err := CreateFile(filename, tt.args.template); (err != nil) != tt.wantErr {
				t.Errorf("CreateFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			// Compare the contents of the generated file with the expected file
			if err := compareFiles(filename, tt.checkFile); err != nil {
				t.Errorf("CreateFile() = %v", err)
			}
		})
	}
}

// compareFiles compares the contents of two files.
func compareFiles(file1, file2 string) error {
	// Read the contents of the first file
	bytes1, err := os.ReadFile(file1)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file1, err)
	}

	// Read the contents of the second file
	bytes2, err := os.ReadFile(file2)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file2, err)
	}

	// Compare the contents of the two files
	if !bytes.Equal(bytes1, bytes2) {
		return fmt.Errorf("contents of %s and %s are not the same:\n%s\n------\n%s", file1, file2, string(bytes1), string(bytes2))
	}

	return nil
}
