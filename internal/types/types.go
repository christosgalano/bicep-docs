/*
Package types provides shared types used by multiple packages in the "bicep-docs" application.
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
// A parameter has a type, an optional default value and an optional metadata part.
//
// The type is the type of the parameter (e.g. "string").
// The default value is the default value of the parameter.
// The metadata part is the optional metadata part of the parameter (just the description).
type Parameter struct {
	Type         string    `json:"type"`
	DefaultValue any       `json:"defaultValue"`
	Metadata     *Metadata `json:"metadata"`
}

// Output is a struct that contains the information about an output.
// An output has a type and an optional metadata part.
//
// The type is the type of the output (e.g. "string").
// The metadata part is the optional metadata part of the output (just the description).
type Output struct {
	Type     string    `json:"type"`
	Metadata *Metadata `json:"metadata"`
}

// Template is a struct that contains the information about a Bicep template.
//
// A template has a file name, a list of modules, resources, a map of parameters, outputs
// and an optional metadata part.
type Template struct {
	FileName   string               `json:"-"`
	Modules    []Module             `json:"-"`
	Resources  []Resource           `json:"-"`
	Parameters map[string]Parameter `json:"parameters"`
	Outputs    map[string]Output    `json:"outputs"`
	Metadata   *Metadata            `json:"metadata"`
}
