package markdown

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/christosgalano/bicep-docs/internal/types"
)

var (
	defaultSections = []types.Section{
		types.DescriptionSection,
		types.UsageSection,
		types.ModulesSection,
		types.ResourcesSection,
		types.ParametersSection,
		types.UserDefinedDataTypesSection,
		types.UserDefinedFunctionsSection,
		types.VariablesSection,
		types.OutputsSection,
	}
)

func TestCreateFile(t *testing.T) {
	templateName := "test"
	templateDescription := "This is a test template."
	parameterDescription := "This is a test parameter."
	stringType := "string"
	positiveIntType := "#/definitions/positive_int"

	basicTemplate := &types.Template{
		FileName: "test.bicep",
		Modules: []types.Module{
			{
				SymbolicName: "test_module",
				Source:       "./modules/test_module/main.bicep",
				Description:  "This is a test module.",
			},
		},
		Resources: []types.Resource{
			{
				SymbolicName: "test_resource",
				Type:         "Microsoft.Storage/storageAccounts",
				Description:  "This is a test resource.",
			},
		},
		Parameters: []types.Parameter{
			{
				Name:         "test_parameter1",
				Type:         "string",
				DefaultValue: "test",
				Metadata: &types.Metadata{
					Description: &parameterDescription,
				},
			},
			{
				Name:         "test_parameter2",
				Type:         "object",
				DefaultValue: map[string]any{},
				Metadata: &types.Metadata{
					Description: &parameterDescription,
				},
			},
			{
				Name:         "test_parameter3",
				Type:         "object",
				DefaultValue: map[string]any{"key1": "value1"},
				Metadata: &types.Metadata{
					Description: &parameterDescription,
				},
			},
		},
		Variables: []types.Variable{
			{
				Name:        "test_variable",
				Description: "This is a test variable.",
			},
		},
		Outputs: []types.Output{
			{
				Name: "test_output",
				Type: "string",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a test output."; return &s }(),
				},
			},
		},
		Metadata: &types.Metadata{
			Name:        &templateName,
			Description: &templateDescription,
		},
	}
	extendedTemplate := &types.Template{
		FileName: "test.bicep",
		Modules: []types.Module{
			{
				SymbolicName: "test_module",
				Source:       "./modules/test_module/main.bicep",
				Description:  "This is a test module.",
			},
		},
		Resources: []types.Resource{
			{
				SymbolicName: "test_resource",
				Type:         "Microsoft.Storage/storageAccounts",
				Description:  "This is a test resource.",
			},
		},
		Parameters: []types.Parameter{
			{
				Name: "required",
				Type: "string",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a required parameter."; return &s }(),
				},
			},
			{
				Name:     "nullable",
				Type:     "string",
				Nullable: true,
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a nullable parameter."; return &s }(),
				},
			},
			{
				Name:         "optional",
				Type:         "string",
				DefaultValue: "test",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is an optional parameter."; return &s }(),
				},
			},
			{
				Name:         "string_array",
				Type:         "array",
				Items:        func() *types.Items { i := &types.Items{Type: &stringType}; return i }(),
				DefaultValue: []string{},
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a string array parameter."; return &s }(),
				},
			},
			{
				Name:  "pint_array",
				Type:  "array",
				Items: func() *types.Items { i := &types.Items{Ref: &positiveIntType}; return i }(),
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a positive int custom type array parameter."; return &s }(),
				},
			},
			{
				Name: "simple_array",
				Type: "array",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is simple array parameter."; return &s }(),
				},
			},
		},
		UserDefinedDataTypes: []types.UserDefinedDataType{
			{
				Name: "pint",
				Type: positiveIntType,
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined type (alias)."; return &s }(),
				},
			},
			{
				Name: "positive_int",
				Type: "int",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined type."; return &s }(),
				},
			},
			{
				Name: "string_array",
				Type: "array",
				Items: &types.Items{
					Type: &stringType,
				},
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined type with array items."; return &s }(),
				},
			},
			{
				Name: "custom_type",
				Type: "object",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined type with properties."; return &s }(),
				},
				Properties: []types.UserDefinedDataTypeProperty{
					{
						Name: "property1",
						Type: "string",
						Metadata: &types.Metadata{
							Description: func() *string { s := "This is a property of a user defined type."; return &s }(),
						},
					},
					{
						Name: "property_2",
						Type: positiveIntType,
						Metadata: &types.Metadata{
							Description: func() *string { s := "This is another property of a user defined type which uses ref."; return &s }(),
						},
					},
					{
						Name: "property_3",
						Type: "array",
						Items: &types.Items{
							Ref: &positiveIntType,
						},
						Metadata: &types.Metadata{
							Description: func() *string { s := "This is a property of a user defined type with array items."; return &s }(),
						},
					},
				},
			},
		},
		UserDefinedFunctions: []types.UserDefinedFunction{
			{
				Name: "build_url",
				Output: types.Output{
					Type: stringType,
				},
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined function."; return &s }(),
				},
			},
			{
				Name: "double",
				Output: types.Output{
					Type: positiveIntType,
				},
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined function with uddts."; return &s }(),
				},
			},
			{
				Name: "get_string_array",
				Output: types.Output{
					Type: "array",
					Items: &types.Items{
						Type: &stringType,
					},
				},
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined function with array items as output."; return &s }(),
				},
			},
		},
		Variables: []types.Variable{
			{
				Name:        "test_variable",
				Description: "This is a test variable.",
			},
		},
		Outputs: []types.Output{
			{
				Name: "pint",
				Type: positiveIntType,
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is an output with uddt."; return &s }(),
				},
			},
			{
				Name: "pint_array",
				Type: "array",
				Items: &types.Items{
					Ref: &positiveIntType,
				},
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is an output with uddt array items."; return &s }(),
				},
			},
		},
		Metadata: &types.Metadata{
			Name:        &templateName,
			Description: &templateDescription,
		},
	}

	type args struct {
		filename          string
		template          *types.Template
		showAllDecorators bool
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		checkFile string
	}{
		{
			name: "basic template",
			args: args{
				filename:          "basic.md",
				template:          basicTemplate,
				showAllDecorators: false,
			},
			wantErr:   false,
			checkFile: "./testdata/basic.md",
		},
		{
			name: "extended template",
			args: args{
				filename:          "extended.md",
				template:          extendedTemplate,
				showAllDecorators: false,
			},
			wantErr:   false,
			checkFile: "./testdata/extended.md",
		},
		{
			name: "multiline markup",
			args: args{
				filename: "multiline_markup.md",
				template: &types.Template{
					FileName: "test.bicep",
					Parameters: []types.Parameter{
						{
							Name: "storageAccountName",
							Type: "string",
							Metadata: &types.Metadata{
								Description: func() *string {
									s := "Storage account name restrictions:\n" +
										"- Storage account names must be between 3 and 24 characters in length and may contain numbers and lowercase letters only.\n" +
										"- Your storage account name must be unique within Azure. No two storage accounts can have the same name.\n"
									return &s
								}(),
							},
						},
					},
				},
				showAllDecorators: false,
			},
			wantErr:   false,
			checkFile: "./testdata/multiline_markup.md",
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
				showAllDecorators: false,
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
				showAllDecorators: false,
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
				showAllDecorators: false,
			},
			wantErr:   false,
			checkFile: "./testdata/no_metadata.md",
		},
		{
			name: "given path is a directory",
			args: args{
				filename:          "testdata",
				template:          nil,
				showAllDecorators: false,
			},
			wantErr: true,
		},
		{
			name: "nil template",
			args: args{
				filename:          "nil_template.md",
				template:          nil,
				showAllDecorators: false,
			},
			wantErr: true,
		},
	}

	// Create a temporary directory for the output files
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call CreateFile with the filename in the temporary directory
			filename := filepath.Join(tempDir, tt.args.filename)
			if err := CreateFile(filename, tt.args.template, false, defaultSections, tt.args.showAllDecorators); (err != nil) != tt.wantErr {
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

func Test_checkFileExists(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "file exists",
			args:    args{filename: "./testdata/basic.md"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "file does not exist",
			args:    args{filename: "./testdata/does_not_exist.md"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "given path is a directory",
			args:    args{filename: "./testdata"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := checkFileExists(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkFileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkFileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

// compareFiles compares the contents of two files.
func compareFiles(file1, file2 string) error {
	// Read the contents of the first file
	generatedContent, err := os.ReadFile(file1)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file1, err)
	}

	// Read the contents of the second file
	expectedContent, err := os.ReadFile(file2)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file2, err)
	}

	// Normalize the newline characters
	generatedContent = bytes.ReplaceAll(generatedContent, []byte("\r\n"), []byte("\n"))
	expectedContent = bytes.ReplaceAll(expectedContent, []byte("\r\n"), []byte("\n"))

	// Compare the contents of the two files
	if !bytes.Equal(generatedContent, expectedContent) {
		return fmt.Errorf("contents of %s and %s are not the same:\n%s\n------\n%s", file1, file2, generatedContent, expectedContent)
	}

	return nil
}
