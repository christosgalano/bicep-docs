package types

import (
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// unmarshalTypeOrRef unmarshals a JSON object into a type or a $ref.
func unmarshalTypeOrRef(data []byte) (string, error) {
	var result struct {
		Type string `json:"type"`
		Ref  string `json:"$ref"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	if result.Type != "" {
		return result.Type, nil
	}
	if result.Ref != "" {
		return result.Ref, nil
	}
	return "", fmt.Errorf("neither type nor $ref field found")
}

// UnmarshalJSON unmarshals a JSON object into a Parameter.
// The type field can be either a type or a $ref.
func (p *Parameter) UnmarshalJSON(data []byte) error {
	type Alias Parameter
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	tr, err := unmarshalTypeOrRef(data)
	if err != nil {
		return err
	}
	p.Type = tr

	return nil
}

// UnmarshalJSON unmarshals a JSON object into an Output.
// The type field can be either a type or a $ref.
func (o *Output) UnmarshalJSON(data []byte) error {
	type Alias Output
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Process type or $ref
	tr, err := unmarshalTypeOrRef(data)
	if err != nil {
		return err
	}
	o.Type = tr

	return nil
}

// UnmarshalJSON unmarshals a JSON object into a UserDefinedDataType.
// The type field can be either a type or a $ref.
func (u *UserDefinedDataType) UnmarshalJSON(data []byte) error {
	type Alias UserDefinedDataType
	aux := &struct {
		Properties map[string]UserDefinedDataTypeProperty `json:"properties"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Process type or $ref
	tr, err := unmarshalTypeOrRef(data)
	if err != nil {
		return err
	}
	u.Type = tr

	// Process properties
	for name, property := range aux.Properties {
		property.Name = name
		u.Properties = append(u.Properties, property)
	}

	return nil
}

// UnmarshalJSON unmarshals a JSON object into a UserDefinedDataTypeProperty.
// The type field can be either a type or a $ref.
func (p *UserDefinedDataTypeProperty) UnmarshalJSON(data []byte) error {
	type Alias UserDefinedDataTypeProperty
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	tr, err := unmarshalTypeOrRef(data)
	if err != nil {
		return err
	}
	p.Type = tr

	return nil
}

// UnmarshalJSON unmarshals a JSON object into a Template.
//
// The parameters, data types, variables, outputs, and functions are unmarshalled
// into slices of Parameter, UserDefinedDataType, Variable, Output, and
// UserDefinedFunction, respectively.
//
// The slices are then sorted by name.
func (t *Template) UnmarshalJSON(data []byte) error { //nolint:gocyclo // Complexity due to handling multiple template sections (parameters, variables, functions, etc.)
	type function struct {
		Namespace string                         `json:"namespace"`
		Functions map[string]UserDefinedFunction `json:"members"`
	}

	// copyOperation represents a copy operation in the variables section of a Bicep file
	type copyOperation struct {
		Name  string `json:"name"`
		Count string `json:"count"`
		Input any    `json:"input"`
	}

	type Alias Template
	aux := &struct {
		Parameters map[string]Parameter           `json:"parameters"`
		DataTypes  map[string]UserDefinedDataType `json:"definitions"`
		Variables  map[string]any                 `json:"variables"`
		Outputs    map[string]Output              `json:"outputs"`
		Functions  []function                     `json:"functions"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Process variables section
	if aux.Variables != nil {
		// Check if variables contains a copy array
		if copyData, ok := aux.Variables["copy"]; ok {
			// Unmarshal the copy array
			var copyOps []copyOperation
			if copyBytes, err := json.Marshal(copyData); err == nil {
				if err := json.Unmarshal(copyBytes, &copyOps); err == nil {
					// Process each copy operation
					for _, copyOp := range copyOps {
						variable := Variable{
							Name:  copyOp.Name,
							Value: copyOp.Input,
						}
						t.Variables = append(t.Variables, variable)
					}
				}
			}
		}

		// Process regular variables
		for name, value := range aux.Variables {
			if name == "copy" || strings.HasPrefix(name, "$fxv#") {
				continue
			}
			variable := Variable{
				Name:  name,
				Value: value,
			}
			t.Variables = append(t.Variables, variable)
		}
	}

	// Process parameters
	for name, param := range aux.Parameters {
		param.Name = name
		t.Parameters = append(t.Parameters, param)
	}

	// Process data types
	for name, dataType := range aux.DataTypes {
		dataType.Name = name
		t.UserDefinedDataTypes = append(t.UserDefinedDataTypes, dataType)
	}

	// Process functions
	for _, function := range aux.Functions {
		for name, userDefinedFunction := range function.Functions {
			userDefinedFunction.Name = name
			t.UserDefinedFunctions = append(t.UserDefinedFunctions, userDefinedFunction)
		}
	}

	// Process outputs
	for name, output := range aux.Outputs {
		output.Name = name
		t.Outputs = append(t.Outputs, output)
	}

	// Sort all fields
	t.Sort()

	return nil
}
