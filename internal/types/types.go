/*
Package types provides shared types used by multiple packages in the "bicep-docs" application.

It also provides custom UnmarshalJSON and Sort functions for some of those types.
*/
package types

import (
	"errors"
	"strings"
)

// Metadata is a struct that contains the metadata part of a parameter, an output, or the template itself.
// The metadata part consists of three optional fields: name, description, and export flag.
//
// A name can be either a metadata item (metadata name = '...') for the template, or a parameter/output name.
//
// A description can be either an annotation (@description('...') | @sys.description('...'))
// or a metadata item (metadata description = '...').
//
// The export flag can be set using the @export() annotation or the metadata item (metadata __bicep_export! = true).
type Metadata struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Export      *bool   `json:"__bicep_export!,omitempty"`
}

// Items represents the array item type information.
// The type field can be either a type or a $ref.
type Items struct {
	Type *string `json:"type"`
	Ref  *string `json:"$ref"`
}

// Module is a struct that contains the information about a module.
// A module has a symbolic name, a source, and an optional description.
//
// The symbolic name is the name of the module that is used to reference the module.
// The source is either the path to the module file or the URL of the remote module.
// The description is an optional description of the module.
//
// Example:
// module network './modules/network/main.bicep'
//
// In the example above, the symbolic name is "network" and the source is "./modules/network/main.bicep".
type Module struct {
	SymbolicName string
	Source       string
	Description  string
}

// Resource is a struct that contains the information about a resource.
// A resource has a symbolic name, a type, and an optional description.
//
// The symbolic name is the name of the resource that is used to reference the resource.
// The type is the type of the resource (e.g. "Microsoft.Network/virtualNetworks").
// The description is an optional description of the resource.
type Resource struct {
	SymbolicName string
	Type         string
	Description  string
}

// ParameterStatus is an enum that represents the status of a parameter.
// The status can be either "Required" or "Optional".
type ParameterStatus string

const (
	RequiredParameterStatus ParameterStatus = "Required" // Required indicates that the parameter is required
	OptionalParameterStatus ParameterStatus = "Optional" // Optional indicates that the parameter is optional
)

func (ps ParameterStatus) String() string {
	return string(ps)
}

// Parameter is a struct that contains the information about a parameter.
// A parameter has a name, type, an optional default value, items (for array types), nullable flag,
// optional constraints (allowed values, minLength, maxLength, minValue, maxValue)
// and an optional metadata part.
type Parameter struct {
	Name          string    `json:"-"`
	Type          string    `json:"-"`
	DefaultValue  any       `json:"defaultValue"`
	Items         *Items    `json:"items"`
	Nullable      bool      `json:"nullable"`
	AllowedValues []any     `json:"allowedValues,omitempty"`
	MinLength     *int      `json:"minLength,omitempty"`
	MaxLength     *int      `json:"maxLength,omitempty"`
	MinValue      *int      `json:"minValue,omitempty"`
	MaxValue      *int      `json:"maxValue,omitempty"`
	Metadata      *Metadata `json:"metadata"`
}

// IsRequired checks if the parameter is required.
// A parameter is considered required if it has no default value and is not nullable.
func (p *Parameter) IsRequired() bool {
	return p.DefaultValue == nil && !p.Nullable
}

// GetStatus returns the status of the parameter.
// If the parameter is required, it returns RequiredParameterStatus.
// If the parameter is optional, it returns OptionalParameterStatus.
func (p *Parameter) GetStatus() ParameterStatus {
	if p.IsRequired() {
		return RequiredParameterStatus
	}
	return OptionalParameterStatus
}

// HasConstraints returns true if the parameter has any constraint decorators.
func (p *Parameter) HasConstraints() bool {
	return len(p.AllowedValues) > 0 ||
		p.MinLength != nil ||
		p.MaxLength != nil ||
		p.MinValue != nil ||
		p.MaxValue != nil
}

// GetDescription returns the description from metadata.
func (p *Parameter) GetDescription() string {
	if p.Metadata != nil && p.Metadata.Description != nil {
		return *p.Metadata.Description
	}
	return ""
}

// UserDefinedDataType (UDDT) is a struct that contains the information about a user defined data type.
// A user defined data type has a name, type, items (for array types), properties, nullable flag,
// optional constraints (minLength, maxLength, minValue, maxValue), export flag,
// and an optional metadata part.
type UserDefinedDataType struct {
	Name       string                        `json:"-"`
	Type       string                        `json:"-"`
	Items      *Items                        `json:"items"`
	Properties []UserDefinedDataTypeProperty `json:"-"`
	Nullable   bool                          `json:"nullable"`
	MinLength  *int                          `json:"minLength,omitempty"`
	MaxLength  *int                          `json:"maxLength,omitempty"`
	MinValue   *int                          `json:"minValue,omitempty"`
	MaxValue   *int                          `json:"maxValue,omitempty"`
	Exportable bool                          `json:"-"`
	Metadata   *Metadata                     `json:"metadata"`
}

// IsExportable returns true if the user-defined data type is marked as exportable.
func (u *UserDefinedDataType) IsExportable() bool {
	return u.Metadata != nil && u.Metadata.Export != nil && *u.Metadata.Export
}

// UserDefinedDataTypeProperty is a struct that contains the information about a property of a user defined data type.
// A property has a name, type, an optional default value, items (for array types), nullable flag,
// optional constraints (allowed values, minLength, maxLength, minValue, maxValue)
// and an optional metadata part.
type UserDefinedDataTypeProperty struct {
	Name          string    `json:"-"`
	Type          string    `json:"-"`
	Items         *Items    `json:"items"`
	Nullable      bool      `json:"nullable"`
	AllowedValues []any     `json:"allowedValues,omitempty"`
	MinLength     *int      `json:"minLength,omitempty"`
	MaxLength     *int      `json:"maxLength,omitempty"`
	MinValue      *int      `json:"minValue,omitempty"`
	MaxValue      *int      `json:"maxValue,omitempty"`
	Metadata      *Metadata `json:"metadata"`
}

// UserDefinedFunction (UDF) is a struct that contains the information about a user defined function.
// A user defined function has a name, a list of parameters, an output, an export flag,
// and an optional metadata part.
type UserDefinedFunction struct {
	Name       string      `json:"-"`
	Parameters []Parameter `json:"parameters"`
	Output     Output      `json:"output"`
	Exportable bool        `json:"-"`
	Metadata   *Metadata   `json:"metadata"`
}

// IsExportable returns true if the user-defined function is marked as exportable.
func (u *UserDefinedFunction) IsExportable() bool {
	return u.Metadata != nil && u.Metadata.Export != nil && *u.Metadata.Export
}

// Variable is a struct that contains the information about a variable.
// A variable has a name and a value.
//
// The name is the name of the variable.
// The value is the value of the variable.
// The description is an optional description of the variable.
type Variable struct {
	Name        string `json:"-"`
	Value       any    `json:"-"`
	Description string `json:"-"`
}

// Output is a struct that contains the information about an output.
// An output has a name, type, items (for array types), nullable flag,
// optional constraints (minLength, maxLength, minValue, maxValue)
// and an optional metadata part.
type Output struct {
	Name      string    `json:"-"`
	Type      string    `json:"-"`
	Items     *Items    `json:"items"`
	Nullable  bool      `json:"nullable"`
	MinLength *int      `json:"minLength,omitempty"`
	MaxLength *int      `json:"maxLength,omitempty"`
	MinValue  *int      `json:"minValue,omitempty"`
	MaxValue  *int      `json:"maxValue,omitempty"`
	Metadata  *Metadata `json:"metadata"`
}

// Template is a struct that contains the information about a Bicep template.
//
// A template has a list of: modules, resources, parameters, user defined data types,
// user defined functions, variables, outputs and an optional metadata part.
type Template struct {
	FileName             string                `json:"-"`
	Modules              []Module              `json:"-"`
	Resources            []Resource            `json:"-"`
	Parameters           []Parameter           `json:"-"`
	UserDefinedDataTypes []UserDefinedDataType `json:"-"`
	UserDefinedFunctions []UserDefinedFunction `json:"-"`
	Variables            []Variable            `json:"-"`
	Outputs              []Output              `json:"-"`
	Metadata             *Metadata             `json:"metadata"`
}

// Section is an enum that represents the different sections of the generated Markdown file.
type Section string

const (
	DescriptionSection          Section = "description"
	UsageSection                Section = "usage"
	ModulesSection              Section = "modules"
	ResourcesSection            Section = "resources"
	ParametersSection           Section = "parameters"
	UserDefinedDataTypesSection Section = "uddts"
	UserDefinedFunctionsSection Section = "udfs"
	VariablesSection            Section = "variables"
	OutputsSection              Section = "outputs"
)

// ParseSectionFromString converts a string to its corresponding Section enum value.
func ParseSectionFromString(str string) (Section, error) {
	switch strings.ToLower(str) {
	case "description":
		return DescriptionSection, nil
	case "usage":
		return UsageSection, nil
	case "modules":
		return ModulesSection, nil
	case "resources":
		return ResourcesSection, nil
	case "parameters":
		return ParametersSection, nil
	case "uddts":
		return UserDefinedDataTypesSection, nil
	case "udfs":
		return UserDefinedFunctionsSection, nil
	case "variables":
		return VariablesSection, nil
	case "outputs":
		return OutputsSection, nil
	default:
		return "", errors.New("invalid section: \"" + str + "\"")
	}
}

// String returns the string representation of a Section.
func (s Section) String() string {
	return string(s)
}
