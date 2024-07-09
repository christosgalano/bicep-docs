/*
Package template provides functions to build and parse Bicep and the corresponding ARM templates.
*/
package template

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/christosgalano/bicep-docs/internal/types"
)

func TestParseTemplates(t *testing.T) {
	templateName := "test"
	templateDescription := "This is a test template."
	parameterDescription := "This is a test parameter."
	outputDescription := "This is a test output."
	basicTemplate := &types.Template{
		FileName: "testdata/basic.bicep",
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
				Name:         "test_parameter",
				Type:         "string",
				DefaultValue: "test",
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
					Description: &outputDescription,
				},
			},
		},
		Metadata: &types.Metadata{
			Name:        &templateName,
			Description: &templateDescription,
		},
	}
	extendedTemplate := &types.Template{
		FileName: "testdata/extended.bicep",
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
				Name:         "test_parameter",
				Type:         "string",
				DefaultValue: "test",
				Metadata: &types.Metadata{
					Description: &parameterDescription,
				},
			},
		},
		UserDefinedDataTypes: []types.UserDefinedDataType{
			{
				Name: "pint",
				Type: "#/definitions/positiveInt",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined type (alias)."; return &s }(),
				},
			},
			{
				Name: "positiveInt",
				Type: "int",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined type."; return &s }(),
				},
			},
		},
		UserDefinedFunctions: []types.UserDefinedFunction{
			{
				Name: "buildUrl",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined function."; return &s }(),
				},
			},
			{
				Name: "double",
				Metadata: &types.Metadata{
					Description: func() *string { s := "This is a user defined function with uddts."; return &s }(),
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
				Type: "#/definitions/positiveInt",
				Metadata: &types.Metadata{
					Description: &outputDescription,
				},
			},
		},
		Metadata: &types.Metadata{
			Name:        &templateName,
			Description: &templateDescription,
		},
	}

	type args struct {
		bicepFile string
		armFile   string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Template
		wantErr bool
	}{
		{
			name: "basic_template",
			args: args{
				bicepFile: "testdata/basic.bicep",
				armFile:   "testdata/basic.json",
			},
			want:    basicTemplate,
			wantErr: false,
		},
		{
			name: "extended_template",
			args: args{
				bicepFile: "testdata/extended.bicep",
				armFile:   "testdata/extended.json",
			},
			want:    extendedTemplate,
			wantErr: false,
		},
		{
			name: "non_existent_template",
			args: args{
				bicepFile: "testdata/non-existent.bicep",
				armFile:   "testdata/non-existent.json",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseTemplates(tt.args.bicepFile, tt.args.armFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTemplates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got.FileName != tt.want.FileName {
				t.Errorf("ParseTemplates() = %v, want %v", got.FileName, tt.want.FileName)
			}
			if !reflect.DeepEqual(got.Modules, tt.want.Modules) {
				t.Errorf("ParseTemplates() = %v, want %v", got.Modules, tt.want.Modules)
			}
			if !reflect.DeepEqual(got.Resources, tt.want.Resources) {
				t.Errorf("ParseTemplates() = %v, want %v", got.Resources, tt.want.Resources)
			}

			for name := range got.Parameters {
				if !reflect.DeepEqual(got.Parameters[name], tt.want.Parameters[name]) {
					t.Errorf("ParseTemplates() = %v, want %v", got.Parameters[name], tt.want.Parameters[name])
				}
			}
			for name := range got.Outputs {
				if !reflect.DeepEqual(got.Outputs[name], tt.want.Outputs[name]) {
					t.Errorf("ParseTemplates() = %v, want %v", got.Outputs[name], tt.want.Outputs[name])
				}
			}
			if !reflect.DeepEqual(got.Metadata, tt.want.Metadata) {
				t.Errorf("ParseTemplates() = %v, want %v", got.Metadata, tt.want.Metadata)
			}
		})
	}
}

func BenchmarkParseTemplates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseTemplates("testdata/extended.bicep", "testdata/extended.json")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Test_parseBicepTemplate(t *testing.T) {
	type args struct {
		bicepFile string
	}
	tests := []struct {
		name          string
		args          args
		wantModules   []types.Module
		wantResources []types.Resource
		wantVariables []types.Variable
		wantErr       bool
	}{
		{
			name: "basic",
			args: args{
				bicepFile: "testdata/basic.bicep",
			},
			wantModules: []types.Module{
				{
					SymbolicName: "test_module",
					Source:       "./modules/test_module/main.bicep",
					Description:  "This is a test module.",
				},
			},
			wantResources: []types.Resource{
				{
					SymbolicName: "test_resource",
					Type:         "Microsoft.Storage/storageAccounts",
					Description:  "This is a test resource.",
				},
			},
			wantVariables: []types.Variable{
				{
					Name:        "test_variable",
					Description: "This is a test variable.",
				},
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				bicepFile: "testdata/empty.bicep",
			},
			wantModules:   []types.Module{},
			wantResources: []types.Resource{},
			wantVariables: []types.Variable{},
			wantErr:       false,
		},
		{
			name: "non_existent",
			args: args{
				bicepFile: "testdata/non-existent.bicep",
			},
			wantModules:   []types.Module{},
			wantResources: []types.Resource{},
			wantVariables: []types.Variable{},
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			modules, resources, variables, err := parseBicepTemplate(tt.args.bicepFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBicepTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(modules, tt.wantModules) {
				t.Errorf("parseBicepTemplate() got modules = %v, want %v", modules, tt.wantModules)
			}
			if !reflect.DeepEqual(resources, tt.wantResources) {
				t.Errorf("parseBicepTemplate() got resources = %v, want %v", resources, tt.wantResources)
			}
			if !reflect.DeepEqual(variables, tt.wantVariables) {
				t.Errorf("parseBicepTemplate() got variables = %v, want %v", variables, tt.wantVariables)
			}
		})
	}
}

func Test_parseDescription(t *testing.T) {
	type args struct {
		line    string
		scanner *bufio.Scanner
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "inline_description",
			args: args{
				line:    "@description('This is a description')",
				scanner: nil,
			},
			want: func() *string { s := "This is a description"; return &s }(),
		},
		{
			name: "multiline_description",
			args: args{
				line:    "@sys.description('''This is a multiline ",
				scanner: bufio.NewScanner(strings.NewReader("\ndescription\n.''' )")),
			},
			want: func() *string { s := "This is a multiline description."; return &s }(),
		},
		{
			name: "no_description",
			args: args{
				line:    "This is not a description",
				scanner: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseDescription(tt.args.line, tt.args.scanner); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDescription() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func Test_parseModule(t *testing.T) {
	type args struct {
		line        string
		description string
	}
	tests := []struct {
		name string
		args args
		want *types.Module
	}{
		{
			name: "registry_source",
			args: args{
				line:        "module test 'br:exampleregistry.azurecr.io/bicep/modules/storage:v1'",
				description: "This is a module",
			},
			want: &types.Module{
				SymbolicName: "test",
				Source:       "br:exampleregistry.azurecr.io/bicep/modules/storage:v1",
			},
		},
		{
			name: "local_source",
			args: args{
				line:        "module test './modules/test_module/main.bicep'",
				description: "This is a module",
			},
			want: &types.Module{
				SymbolicName: "test",
				Source:       "./modules/test_module/main.bicep",
			},
		},
		{
			name: "invalid_module",
			args: args{
				line:        "invalid line",
				description: "This is not a module",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseModule(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseModule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseResource(t *testing.T) {
	type args struct {
		line        string
		description string
	}
	tests := []struct {
		name string
		args args
		want *types.Resource
	}{
		{
			name: "storage_account",
			args: args{
				line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01'",
			},
			want: &types.Resource{
				SymbolicName: "test",
				Type:         "Microsoft.Storage/storageAccounts",
			},
		},
		{
			name: "invalid_resource",
			args: args{
				line:        "invalid line",
				description: "This is not a resource",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseResource(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_skipComment(t *testing.T) {
	type args struct {
		line     string
		scanner  *bufio.Scanner
		expected string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "single_line_comment",
			args: args{
				line:     "// This is a single line comment",
				scanner:  nil,
				expected: "",
			},
		},
		{
			name: "multiline comment",
			args: args{
				line:     "/* This is a",
				scanner:  bufio.NewScanner(strings.NewReader("\n multiline\ncomment */\n")),
				expected: "",
			},
		},
		{
			name: "normal_line",
			args: args{
				line:     "This is a normal line",
				scanner:  nil,
				expected: "This is a normal line",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			skipped, err := skipComment(tt.args.line, tt.args.scanner)
			if err != nil {
				t.Errorf("skipComment() error = %v", err)
			}
			if skipped {
				if tt.args.scanner != nil && tt.args.scanner.Scan() {
					t.Errorf("got %q, want %q", tt.args.scanner.Text(), tt.args.expected)
				}
			} else if tt.args.line != tt.args.expected {
				t.Errorf("got %q, want %q", tt.args.line, tt.args.expected)
			}
		})
	}
}
