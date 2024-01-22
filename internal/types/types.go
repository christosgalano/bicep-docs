/*
Package types provides shared types used by multiple packages in the "bicep-docs" application.
*/
package types

// Metadata is a struct that contains the metadata part of a parameter,
// an output, or the template itself.
// The metadata part consists of two optional fields: name and description.
type Metadata struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// Module is a struct that contains the information about a module.
// A module has a symbolic name, a source, and an optional description.
type Module struct {
	SymbolicName   string
	Source         string
	Description    string
}

// Resource is a struct that contains the information about a resource.
// A resource has a symbolic name, a type, and an optional description.
type Resource struct {
	SymbolicName string
	Type         string
	Description    string
}

// Parameter is a struct that contains the information about a parameter.
// A parameter has a type, an optional default value and an optional metadata part.
type Parameter struct {
	Type         string    `json:"type"`
	DefaultValue any       `json:"defaultValue"`
	Metadata     *Metadata `json:"metadata"`
}

// Output is a struct that contains the information about an output.
// An output has a type and an optional metadata part.
type Output struct {
	Type     string    `json:"type"`
	Metadata *Metadata `json:"metadata"`
}

// Template is a struct that contains the information about a template (Bicep or ARM).
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
