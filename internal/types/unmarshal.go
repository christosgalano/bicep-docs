package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

// UnmarshalJSON unmarshals a JSON object into a Parameter.
// The type field can be either a type or a $ref.
func (p *Parameter) UnmarshalJSON(data []byte) error {
	type Alias Parameter
	aux := &struct {
		Type json.RawMessage `json:"type"`
		Ref  json.RawMessage `json:"$ref"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var ref string
	if err := json.Unmarshal(aux.Ref, &ref); err == nil {
		p.Type = ref
		return nil
	}

	var typ string
	if err := json.Unmarshal(aux.Type, &typ); err == nil {
		p.Type = typ
		return nil
	}

	return fmt.Errorf("type field could not be unmarshalled")
}

// UnmarshalJSON unmarshals a JSON object into an Output.
// The type field can be either a type or a $ref.
func (o *Output) UnmarshalJSON(data []byte) error {
	type Alias Output
	aux := &struct {
		Type json.RawMessage `json:"type"`
		Ref  json.RawMessage `json:"$ref"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var ref string
	if err := json.Unmarshal(aux.Ref, &ref); err == nil {
		o.Type = ref
		return nil
	}

	var typ string
	if err := json.Unmarshal(aux.Type, &typ); err == nil {
		o.Type = typ
		return nil
	}

	return fmt.Errorf("type field could not be unmarshalled")
}

// UnmarshalJSON unmarshals a JSON object into a UserDefinedDataType.
// The type field can be either a type or a $ref.
func (u *UserDefinedDataType) UnmarshalJSON(data []byte) error {
	type Alias UserDefinedDataType
	aux := &struct {
		Type json.RawMessage `json:"type"`
		Ref  json.RawMessage `json:"$ref"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var ref string
	if err := json.Unmarshal(aux.Ref, &ref); err == nil {
		u.Type = ref
		return nil
	}

	var typ string
	if err := json.Unmarshal(aux.Type, &typ); err == nil {
		u.Type = typ
		return nil
	}

	return fmt.Errorf("type field could not be unmarshalled")
}

// UnmarshalJSON unmarshals a JSON object into a Template.
//
// The parameters, data types, variables, outputs, and functions are unmarshalled
// into slices of Parameter, UserDefinedDataType, Variable, Output, and
// UserDefinedFunction, respectively.
//
// The slices are then sorted by name.
func (t *Template) UnmarshalJSON(data []byte) error {
	type function struct {
		Namespace string                         `json:"namespace"`
		Functions map[string]UserDefinedFunction `json:"members"`
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

	// Use the key of the map to set the Name field for
	// Parameters, UserDefinedDataTypes, Variables, Outputs,
	// and UserDefinedFunctions.

	for name, param := range aux.Parameters {
		param.Name = name
		t.Parameters = append(t.Parameters, param)
	}

	for name, dataType := range aux.DataTypes {
		dataType.Name = name
		t.UserDefinedDataTypes = append(t.UserDefinedDataTypes, dataType)
	}

	for name, value := range aux.Variables {
		if strings.HasPrefix(name, "$fxv#") {
			continue
		}
		variable := Variable{
			Name:  name,
			Value: value,
		}
		t.Variables = append(t.Variables, variable)
	}

	for name, output := range aux.Outputs {
		output.Name = name
		t.Outputs = append(t.Outputs, output)
	}

	for _, function := range aux.Functions {
		for name, userDefinedFunction := range function.Functions {
			userDefinedFunction.Name = name
			t.UserDefinedFunctions = append(t.UserDefinedFunctions, userDefinedFunction)
		}
	}

	// Sort Parameters, UserDefinedDataTypes, Variables,
	// Outputs, and UserDefinedFunctions by name
	t.Sort()

	return nil
}
