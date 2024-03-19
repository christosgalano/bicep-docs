/*
Package types provides shared types used by multiple packages in the "bicep-docs" application.

It also provides custom UnmarshalJSON and Sort functions for some of those types.
*/
package types

// Metadata is a struct that contains the metadata part of a parameter, an output, or the template itself.
// The metadata part consists of two optional fields: name and description.
//
// A name can be either a metadata item (metadata name = '...') for the template, or a parameter/output name.
//
// A description can be either an annotation (@description('...') | @sys.description('...'))
// or a metadata item (metadata description = '...').
type Metadata struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
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

// Parameter is a struct that contains the information about a parameter.
// A parameter has a name, type, an optional default value and an optional metadata part.
//
// The name is the name of the parameter.
// The type is the type of the parameter (e.g. "string" or even a user defined type).
// The default value is the optional default value of the parameter.
// The metadata part is the optional metadata part of the parameter (just the description).
type Parameter struct {
	Name         string    `json:"-"`
	Type         string    `json:"-"`
	DefaultValue any       `json:"defaultValue"`
	Metadata     *Metadata `json:"metadata"`
}

// UserDefinedDataType (UDDT) is a struct that contains the information about a user defined data type.
// A user defined data type has a name, a type and an optional metadata part.
//
// The name is the name of the user defined data type.
// The type is the type of the user defined data type (e.g. "object" or even including other user defined types).
// The metadata part is the optional metadata part of the user defined data type (just the description).
type UserDefinedDataType struct {
	Name     string    `json:"-"`
	Type     string    `json:"-"`
	Metadata *Metadata `json:"metadata"`
}

// UserDefinedFunction (UDF) is a struct that contains the information about a user defined function.
// A user defined function has a name, a list of parameters, an output and an optional metadata part.
//
// The name is the name of the user defined function.
// The parameters are the parameters of the user defined function.
// The output is the output of the user defined function.
// The metadata part is the optional metadata part of the user defined function (just the description).
type UserDefinedFunction struct {
	Name       string      `json:"-"`
	Parameters []Parameter `json:"parameters"`
	Output     Output      `json:"output"`
	Metadata   *Metadata   `json:"metadata"`
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
// An output has a type and an optional metadata part.
//
// The type is the type of the output (e.g. "string").
// The metadata part is the optional metadata part of the output (just the description).
type Output struct {
	Name     string    `json:"-"`
	Type     string    `json:"-"`
	Metadata *Metadata `json:"metadata"`
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
