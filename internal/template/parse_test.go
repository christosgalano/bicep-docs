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
	// Helper to reduce string pointer boilerplate
	strPtr := func(s string) *string {
		return &s
	}

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
					Description: strPtr("This is a test parameter."),
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
					Description: strPtr("This is a test output."),
				},
			},
		},
		Metadata: &types.Metadata{
			Name:        strPtr("test"),
			Description: strPtr("This is a test template."),
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
					Description: strPtr("This is a test parameter."),
				},
			},
		},
		UserDefinedDataTypes: []types.UserDefinedDataType{
			{
				Name: "pint",
				Type: "#/definitions/positiveInt",
				Metadata: &types.Metadata{
					Description: strPtr("This is a user defined type (alias)."),
				},
			},
			{
				Name: "positiveInt",
				Type: "int",
				Metadata: &types.Metadata{
					Description: strPtr("This is a user defined type."),
				},
			},
		},
		UserDefinedFunctions: []types.UserDefinedFunction{
			{
				Name: "buildUrl",
				Metadata: &types.Metadata{
					Description: strPtr("This is a user defined function."),
				},
			},
			{
				Name: "double",
				Metadata: &types.Metadata{
					Description: strPtr("This is a user defined function with uddts."),
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
					Description: strPtr("This is a test output."),
				},
			},
		},
		Metadata: &types.Metadata{
			Name:        strPtr("test"),
			Description: strPtr("This is a test template."),
		},
	}
	loopsTemplate := &types.Template{
		FileName: "testdata/loops.bicep",
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
				Description:  "This is a storage account resource array.",
			},
		},
		Parameters: []types.Parameter{
			{
				Name:         "locations",
				Type:         "array",
				DefaultValue: []any{"eastus", "westus", "northeurope"},
				Metadata: &types.Metadata{
					Description: strPtr("Array of deployment locations"),
				},
			},
			{
				Name:         "namePrefix",
				Type:         "string",
				DefaultValue: "storage",
				Metadata: &types.Metadata{
					Description: strPtr("Prefix for resource names"),
				},
			},
		},
		Variables: []types.Variable{
			{
				Name:        "storageConfigs",
				Description: "Array of storage account configurations",
			},
			{
				Name:        "storageNames",
				Description: "Array of generated storage account names",
			},
		},
		Outputs: []types.Output{
			{
				Name: "resourceIds",
				Type: "array",
				Metadata: &types.Metadata{
					Description: strPtr("Array of created storage account resource IDs"),
				},
			},
			{
				Name: "storageNames",
				Type: "array",
				Metadata: &types.Metadata{
					Description: strPtr("Array of created storage account names"),
				},
			},
		},
		Metadata: &types.Metadata{
			Name:        strPtr("loop_test"),
			Description: strPtr("Test template with loop constructs"),
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
			name: "loops_template",
			args: args{
				bicepFile: "testdata/loops.bicep",
				armFile:   "testdata/loops.json",
			},
			want:    loopsTemplate,
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

			// Compare all fields properly
			if got.FileName != tt.want.FileName {
				t.Errorf("FileName = %v, want %v", got.FileName, tt.want.FileName)
			}

			// Compare slices with proper ordering
			compareModules(t, got.Modules, tt.want.Modules)
			compareResources(t, got.Resources, tt.want.Resources)
			compareParameters(t, got.Parameters, tt.want.Parameters)
			compareVariables(t, got.Variables, tt.want.Variables)
			compareOutputs(t, got.Outputs, tt.want.Outputs)
			compareMetadata(t, got.Metadata, tt.want.Metadata)
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotModules, gotResources, gotVariables, err := parseBicepTemplate(tt.args.bicepFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBicepTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			compareModules(t, gotModules, tt.wantModules)
			compareResources(t, gotResources, tt.wantResources)
			compareVariables(t, gotVariables, tt.wantVariables)
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

// Helper functions for comparing slices.
func compareModules(t *testing.T, got, want []types.Module) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("Modules length = %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		if got[i].SymbolicName != want[i].SymbolicName {
			t.Errorf("Module[%d].SymbolicName = %v, want %v", i, got[i].SymbolicName, want[i].SymbolicName)
		}
		if got[i].Source != want[i].Source {
			t.Errorf("Module[%d].Source = %v, want %v", i, got[i].Source, want[i].Source)
		}
		if got[i].Description != want[i].Description {
			t.Errorf("Module[%d].Description = %v, want %v", i, got[i].Description, want[i].Description)
		}
	}
}

func compareResources(t *testing.T, got, want []types.Resource) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("Resources length = %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		if got[i].SymbolicName != want[i].SymbolicName {
			t.Errorf("Resource[%d].SymbolicName = %v, want %v", i, got[i].SymbolicName, want[i].SymbolicName)
		}
		if got[i].Type != want[i].Type {
			t.Errorf("Resource[%d].Type = %v, want %v", i, got[i].Type, want[i].Type)
		}
		if got[i].Description != want[i].Description {
			t.Errorf("Resource[%d].Description = %v, want %v", i, got[i].Description, want[i].Description)
		}
	}
}

func compareParameters(t *testing.T, got, want []types.Parameter) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("Parameters length = %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		if got[i].Name != want[i].Name {
			t.Errorf("Parameter[%d].Name = %v, want %v", i, got[i].Name, want[i].Name)
		}
		if got[i].Type != want[i].Type {
			t.Errorf("Parameter[%d].Type = %v, want %v", i, got[i].Type, want[i].Type)
		}

		if !reflect.DeepEqual(got[i].DefaultValue, want[i].DefaultValue) {
			t.Errorf("Parameter[%d].DefaultValue = %v, want %v", i, got[i].DefaultValue, want[i].DefaultValue)
		}

		compareMetadata(t, got[i].Metadata, want[i].Metadata)
	}
}

func compareVariables(t *testing.T, got, want []types.Variable) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("Variables length = %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		if got[i].Name != want[i].Name {
			t.Errorf("Variable[%d].Name = %v, want %v", i, got[i].Name, want[i].Name)
		}
		if got[i].Description != want[i].Description {
			t.Errorf("Variable[%d].Description = %v, want %v", i, got[i].Description, want[i].Description)
		}
	}
}

func compareOutputs(t *testing.T, got, want []types.Output) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("Outputs length = %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		if got[i].Name != want[i].Name {
			t.Errorf("Output[%d].Name = %v, want %v", i, got[i].Name, want[i].Name)
		}
		if got[i].Type != want[i].Type {
			t.Errorf("Output[%d].Type = %v, want %v", i, got[i].Type, want[i].Type)
		}
		compareMetadata(t, got[i].Metadata, want[i].Metadata)
	}
}

func compareMetadata(t *testing.T, got, want *types.Metadata) {
	t.Helper()
	if (got == nil) != (want == nil) {
		t.Errorf("Metadata presence mismatch: got %v, want %v", got, want)
		return
	}
	if got != nil {
		if (got.Name == nil) != (want.Name == nil) {
			t.Errorf("Metadata.Name presence mismatch")
		} else if got.Name != nil && *got.Name != *want.Name {
			t.Errorf("Metadata.Name = %v, want %v", *got.Name, *want.Name)
		}
		if (got.Description == nil) != (want.Description == nil) {
			t.Errorf("Metadata.Description presence mismatch")
		} else if got.Description != nil && *got.Description != *want.Description {
			t.Errorf("Metadata.Description = %v, want %v", *got.Description, *want.Description)
		}
	}
}
