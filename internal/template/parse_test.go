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
	variableOptimizationTemplate := &types.Template{
		FileName: "testdata/var_optimization.bicep",
		Parameters: []types.Parameter{
			{
				Name: "servicePlanName",
				Type: "string",
				Metadata: &types.Metadata{
					Description: strPtr("Name of the App Service Plan to host Web-App on"),
				},
			},
		},
		Resources: []types.Resource{
			{
				SymbolicName: "servicePlan",
				Type:         "Microsoft.Web/serverfarms",
				Description:  "Get App Service Plan Object",
			},
		},
		Variables: []types.Variable{
			{
				Name:        "isZoneRedundant",
				Description: "Get resilience options from Service Plan",
			},
			{
				Name:        "reserved",
				Description: "",
			},
		},
		Outputs: []types.Output{
			{
				Name: "isZoneRedundant",
				Type: "bool",
			},
			{
				Name: "reserved",
				Type: "bool",
			},
		},
		Metadata: &types.Metadata{
			Name:        strPtr("var_optimization"),
			Description: strPtr("Test template with variable optimization"),
		},
	}
	anyTemplate := &types.Template{
		FileName: "testdata/any.bicep",
		Parameters: []types.Parameter{
			{
				Name: "test",
				Type: "any",
			},
		},
		Metadata: &types.Metadata{},
	}
	conditionTemplate := &types.Template{
		FileName: "testdata/condition.bicep",
		Modules: []types.Module{
			{
				SymbolicName: "conditional_module",
				Source:       "./modules/test_module/main.bicep",
				Condition:    "deploy",
				Description:  "This is a conditional module.",
			},
		},
		Resources: []types.Resource{
			{
				SymbolicName: "conditional_resource",
				Type:         "Microsoft.Storage/storageAccounts",
				Condition:    "deploy && environment == 'prod'",
				Description:  "This is a conditional resource.",
			},
			{
				SymbolicName: "loop_resource",
				Type:         "Microsoft.Storage/storageAccounts",
				Condition:    "deploy",
				Description:  "This is a conditional loop resource.",
			},
			{
				SymbolicName: "plain_resource",
				Type:         "Microsoft.Storage/storageAccounts",
				Description:  "This is an unconditional resource.",
			},
		},
		Parameters: []types.Parameter{
			{
				Name:         "deploy",
				Type:         "bool",
				DefaultValue: true,
				Metadata: &types.Metadata{
					Description: strPtr("Whether to deploy the resources."),
				},
			},
			{
				Name:         "environment",
				Type:         "string",
				DefaultValue: "dev",
				Metadata: &types.Metadata{
					Description: strPtr("The environment name."),
				},
			},
		},
		Metadata: &types.Metadata{
			Name:        strPtr("condition-test"),
			Description: strPtr("Template that exercises conditional resources and modules."),
		},
	}

	exportVariablesTemplate := &types.Template{
		FileName: "testdata/export_variables.bicep",
		Variables: []types.Variable{
			{
				Name:       "defaultTags",
				Exportable: true,
			},
			{
				Name:        "internalValue",
				Description: "An internal variable (not exported).",
			},
			{
				Name:        "namePrefix",
				Exportable:  true,
				Description: "A shared prefix used across resources.",
			},
		},
		Outputs: []types.Output{
			{
				Name: "internal",
				Type: "string",
			},
		},
		Metadata: &types.Metadata{
			Name:        strPtr("export-variables-test"),
			Description: strPtr("Template that exercises exported variables."),
		},
	}

	decoratorsTemplate := &types.Template{
		FileName: "testdata/decorators.bicep",
		Resources: []types.Resource{
			{
				SymbolicName:    "idempotent_resource",
				Type:            "Microsoft.Storage/storageAccounts",
				OnlyIfNotExists: true,
				Description:     "This is a resource that is only created if it does not exist.",
			},
			{
				SymbolicName: "plain_resource",
				Type:         "Microsoft.Storage/storageAccounts",
				Description:  "This is a plain resource.",
			},
			{
				SymbolicName: "retry_resource",
				Type:         "Microsoft.Storage/storageAccounts",
				RetryOn:      "['ServerError', 'Conflict'], 3",
				Description:  "This is a resource with retry behavior.",
			},
		},
		Metadata: &types.Metadata{
			Name:        strPtr("decorators-test"),
			Description: strPtr("Template that exercises resource decorators."),
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
			name: "var_optimization_template",
			args: args{
				bicepFile: "testdata/var_optimization.bicep",
				armFile:   "testdata/var_optimization.json",
			},
			want:    variableOptimizationTemplate,
			wantErr: false,
		},
		{
			name: "any_type_template",
			args: args{
				bicepFile: "testdata/any.bicep",
				armFile:   "testdata/any.json",
			},
			want:    anyTemplate,
			wantErr: false,
		},
		{
			name: "condition_template",
			args: args{
				bicepFile: "testdata/condition.bicep",
				armFile:   "testdata/condition.json",
			},
			want:    conditionTemplate,
			wantErr: false,
		},
		{
			name: "export_variables_template",
			args: args{
				bicepFile: "testdata/export_variables.bicep",
				armFile:   "testdata/export_variables.json",
			},
			want:    exportVariablesTemplate,
			wantErr: false,
		},
		{
			name: "decorators_template",
			args: args{
				bicepFile: "testdata/decorators.bicep",
				armFile:   "testdata/decorators.json",
			},
			want:    decoratorsTemplate,
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
			name: "func_description_does_not_leak_to_variable",
			args: args{
				bicepFile: "testdata/func_description.bicep",
			},
			wantModules:   []types.Module{},
			wantResources: []types.Resource{},
			wantVariables: []types.Variable{
				{
					Name: "undescribedVariable",
				},
				{
					Name:        "describedVariable",
					Description: "This is a described variable.",
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
			name: "conditional_module",
			args: args{
				line:        "module test './modules/test_module/main.bicep' = if (deploy) {",
				description: "This is a conditional module",
			},
			want: &types.Module{
				SymbolicName: "test",
				Source:       "./modules/test_module/main.bicep",
				Condition:    "deploy",
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
			name: "conditional_storage_account",
			args: args{
				line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = if (deploy && environment == 'prod') {",
			},
			want: &types.Resource{
				SymbolicName: "test",
				Type:         "Microsoft.Storage/storageAccounts",
				Condition:    "deploy && environment == 'prod'",
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

func Test_parseCondition(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "no_condition",
			line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = {",
			want: "",
		},
		{
			name: "simple_condition",
			line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = if (deploy) {",
			want: "deploy",
		},
		{
			name: "nested_parentheses",
			line: "module test './main.bicep' = if (deploy && (env == 'prod' || env == 'staging')) {",
			want: "deploy && (env == 'prod' || env == 'staging')",
		},
		{
			name: "function_call_in_condition",
			line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = if (contains(names, 'test')) {",
			want: "contains(names, 'test')",
		},
		{
			name: "parenthesis_inside_string",
			line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = if (name == 'value (special)') {",
			want: "name == 'value (special)'",
		},
		{
			name: "unbalanced_condition",
			line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = if (deploy &&",
			want: "",
		},
		{
			name: "loop_with_condition",
			line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = [for i in range(0, 2): if (deploy) {",
			want: "deploy",
		},
		{
			name: "loop_without_condition",
			line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = [for i in range(0, 2): {",
			want: "",
		},
		{
			name: "escaped_quote_inside_string",
			line: `resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = if (name == 'it\'s (a) test') {`,
			want: `name == 'it\'s (a) test'`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseCondition(tt.line); got != tt.want {
				t.Errorf("parseCondition() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_parseRetryOn(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "no_retry_on",
			line: "resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = {",
			want: "",
		},
		{
			name: "retry_on_with_codes_and_count",
			line: "@retryOn(['ServerError', 'Conflict'], 3)",
			want: "['ServerError', 'Conflict'], 3",
		},
		{
			name: "retry_on_with_sys_prefix",
			line: "@sys.retryOn(['ServerError'], 2)",
			want: "['ServerError'], 2",
		},
		{
			name: "retry_on_codes_only",
			line: "@retryOn(['ResourceNotFound'])",
			want: "['ResourceNotFound']",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseRetryOn(tt.line); got != tt.want {
				t.Errorf("parseRetryOn() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_pendingDeclaration_applyDecorator(t *testing.T) {
	tests := []struct {
		name                    string
		lines                   []string
		expectedConsumed        bool
		expectedRetryOn         string
		expectedOnlyIfNotExists bool
	}{
		{
			name:                    "retry_on_only",
			lines:                   []string{"@retryOn(['ServerError'], 3)"},
			expectedConsumed:        true,
			expectedRetryOn:         "['ServerError'], 3",
			expectedOnlyIfNotExists: false,
		},
		{
			name:                    "only_if_not_exists_only",
			lines:                   []string{"@onlyIfNotExists()"},
			expectedConsumed:        true,
			expectedRetryOn:         "",
			expectedOnlyIfNotExists: true,
		},
		{
			name:                    "both_decorators",
			lines:                   []string{"@retryOn(['Conflict'], 2)", "@onlyIfNotExists()"},
			expectedConsumed:        true,
			expectedRetryOn:         "['Conflict'], 2",
			expectedOnlyIfNotExists: true,
		},
		{
			name:                    "not_a_decorator",
			lines:                   []string{"@minLength(3)"},
			expectedConsumed:        false,
			expectedRetryOn:         "",
			expectedOnlyIfNotExists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var pending pendingDeclaration
			consumed := false
			for _, line := range tt.lines {
				consumed = pending.applyDecorator(line)
			}
			if consumed != tt.expectedConsumed {
				t.Errorf("applyDecorator() = %v, want %v", consumed, tt.expectedConsumed)
			}
			if pending.retryOn != tt.expectedRetryOn {
				t.Errorf("retryOn = %q, want %q", pending.retryOn, tt.expectedRetryOn)
			}
			if pending.onlyIfNotExists != tt.expectedOnlyIfNotExists {
				t.Errorf("onlyIfNotExists = %v, want %v", pending.onlyIfNotExists, tt.expectedOnlyIfNotExists)
			}
		})
	}
}

func Test_parseDeclarationLine_decoratorAttribution(t *testing.T) {
	// Decorators pending before a module or variable declaration must not leak into it,
	// and must only be applied to resource declarations.
	pending := pendingDeclaration{description: "desc", retryOn: "['ServerError'], 3", onlyIfNotExists: true}

	var modules []types.Module
	var resources []types.Resource
	var variables []types.Variable

	if !parseDeclarationLine("module test './main.bicep'", pending, &modules, &resources, &variables) {
		t.Fatal("expected module declaration to match")
	}
	if !parseDeclarationLine("resource test 'Microsoft.Storage/storageAccounts@2023-01-01' = {", pending, &modules, &resources, &variables) {
		t.Fatal("expected resource declaration to match")
	}
	if !parseDeclarationLine("var test = 1", pending, &modules, &resources, &variables) {
		t.Fatal("expected variable declaration to match")
	}

	if len(modules) != 1 || len(resources) != 1 || len(variables) != 1 {
		t.Fatalf("declarations: got %d modules, %d resources, %d variables; want 1 of each", len(modules), len(resources), len(variables))
	}
	if resources[0].RetryOn != "['ServerError'], 3" || !resources[0].OnlyIfNotExists {
		t.Errorf("resource decorators not applied: %+v", resources[0])
	}
	if modules[0].Description != "desc" || variables[0].Description != "desc" {
		t.Errorf("description not applied: module %q, variable %q", modules[0].Description, variables[0].Description)
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
		if got[i].Condition != want[i].Condition {
			t.Errorf("Module[%d].Condition = %v, want %v", i, got[i].Condition, want[i].Condition)
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
		if got[i].Condition != want[i].Condition {
			t.Errorf("Resource[%d].Condition = %v, want %v", i, got[i].Condition, want[i].Condition)
		}
		if got[i].RetryOn != want[i].RetryOn {
			t.Errorf("Resource[%d].RetryOn = %v, want %v", i, got[i].RetryOn, want[i].RetryOn)
		}
		if got[i].OnlyIfNotExists != want[i].OnlyIfNotExists {
			t.Errorf("Resource[%d].OnlyIfNotExists = %v, want %v", i, got[i].OnlyIfNotExists, want[i].OnlyIfNotExists)
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
		if got[i].Exportable != want[i].Exportable {
			t.Errorf("Variable[%d].Exportable = %v, want %v", i, got[i].Exportable, want[i].Exportable)
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
