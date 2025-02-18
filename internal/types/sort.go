package types

import "sort"

// Sort sorts the template's parameters, user defined data types, variables,
// outputs and user defined functions by name.
func (t *Template) Sort() {
	sort.Slice(t.Parameters, func(i, j int) bool {
		return t.Parameters[i].Name < t.Parameters[j].Name
	})

	sort.Slice(t.UserDefinedDataTypes, func(i, j int) bool {
		return t.UserDefinedDataTypes[i].Name < t.UserDefinedDataTypes[j].Name
	})
	for i := range t.UserDefinedDataTypes {
		sort.Slice(t.UserDefinedDataTypes[i].Properties, func(j, k int) bool {
			return t.UserDefinedDataTypes[i].Properties[j].Name < t.UserDefinedDataTypes[i].Properties[k].Name
		})
	}

	sort.Slice(t.Variables, func(i, j int) bool {
		return t.Variables[i].Name < t.Variables[j].Name
	})

	sort.Slice(t.Outputs, func(i, j int) bool {
		return t.Outputs[i].Name < t.Outputs[j].Name
	})

	sort.Slice(t.UserDefinedFunctions, func(i, j int) bool {
		return t.UserDefinedFunctions[i].Name < t.UserDefinedFunctions[j].Name
	})
}
