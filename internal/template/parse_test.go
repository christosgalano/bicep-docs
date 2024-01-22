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
			name: "testdata/main.bicep",
			args: args{
				bicepFile: "testdata/main.bicep",
				armFile:   "testdata/main.json",
			},
			want: &types.Template{
				FileName: "testdata/main.bicep",
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
				Parameters: map[string]types.Parameter{
					"test_parameter": {
						Type:         "string",
						DefaultValue: "test",
						Metadata: &types.Metadata{
							Description: func() *string { s := "This is a test parameter."; return &s }(),
						},
					},
				},
				Outputs: map[string]types.Output{
					"test_output": {
						Type: "string",
						Metadata: &types.Metadata{
							Description: func() *string { s := "This is a test output."; return &s }(),
						},
					},
				},
				Metadata: &types.Metadata{
					Name:        func() *string { s := "test"; return &s }(),
					Description: func() *string { s := "This is a test template."; return &s }(),
				},
			},
			wantErr: false,
		},
		{
			name: "non-existent",
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

func Test_parseBicepTemplate(t *testing.T) {
	type args struct {
		bicepFile string
	}
	tests := []struct {
		name          string
		args          args
		wantModules   []types.Module
		wantResources []types.Resource
		wantErr       bool
	}{
		{
			name: "testdata/main.bicep",
			args: args{
				bicepFile: "testdata/main.bicep",
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
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				bicepFile: "testdata/empty.bicep",
			},
			wantModules:   []types.Module{},
			wantResources: []types.Resource{},
			wantErr:       false,
		},
		{
			name: "non-existent",
			args: args{
				bicepFile: "testdata/non-existent.bicep",
			},
			wantModules:   []types.Module{},
			wantResources: []types.Resource{},
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modules, resources, err := parseBicepTemplate(tt.args.bicepFile)
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
			name: "inline description",
			args: args{
				line:    "@description('This is a description')",
				scanner: nil,
			},
			want: func() *string { s := "This is a description"; return &s }(),
		},
		{
			name: "multiline description",
			args: args{
				line:    "@sys.description('''This is a multiline ",
				scanner: bufio.NewScanner(strings.NewReader("\ndescription\n.''' )")),
			},
			want: func() *string { s := "This is a multiline description."; return &s }(),
		},
		{
			name: "no description",
			args: args{
				line:    "This is not a description",
				scanner: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			name: "valid module",
			args: args{
				line:        "module test 'br:exampleregistry.azurecr.io/bicep/modules/storage:v1'",
				description: "This is a module",
			},
			want: &types.Module{
				SymbolicName: "test",
				Source:       "br:exampleregistry.azurecr.io/bicep/modules/storage:v1",
				Description:  "This is a module",
			},
		},
		{
			name: "invalid module",
			args: args{
				line:        "invalid line",
				description: "This is not a module",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseModule(tt.args.line, tt.args.description); !reflect.DeepEqual(got, tt.want) {
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
			name: "valid resource",
			args: args{
				line:        "resource test 'Microsoft.Storage/storageAccounts@2023-01-01'",
				description: "This is a resource",
			},
			want: &types.Resource{
				SymbolicName: "test",
				Type:         "Microsoft.Storage/storageAccounts",
				Description:  "This is a resource",
			},
		},
		{
			name: "invalid resource",
			args: args{
				line:        "invalid line",
				description: "This is not a resource",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseResource(tt.args.line, tt.args.description); !reflect.DeepEqual(got, tt.want) {
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
			name: "single line comment",
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
			name: "normal line",
			args: args{
				line:     "This is a normal line",
				scanner:  nil,
				expected: "This is a normal line",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if skipComment(tt.args.line, tt.args.scanner) {
				if tt.args.scanner != nil && tt.args.scanner.Scan() {
					t.Errorf("got %q, want %q", tt.args.scanner.Text(), tt.args.expected)
				}
			} else {
				if tt.args.line != tt.args.expected {
					t.Errorf("got %q, want %q", tt.args.line, tt.args.expected)
				}
			}
		})
	}
}
