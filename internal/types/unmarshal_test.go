package types

import (
	"testing"
)

func TestUnmarshalTypeOrRef(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "type field present",
			input:    []byte(`{"type": "string"}`),
			expected: "string",
			wantErr:  false,
		},
		{
			name:     "$ref field present",
			input:    []byte(`{"$ref": "#/definitions/myType"}`),
			expected: "#/definitions/myType",
			wantErr:  false,
		},
		{
			name:     "empty object - any type",
			input:    []byte(`{}`),
			expected: "any",
			wantErr:  false,
		},
		{
			name:     "empty object with whitespace - any type",
			input:    []byte(`  {}  `),
			expected: "any",
			wantErr:  false,
		},
		{
			name:     "object with metadata but no type or $ref - any type",
			input:    []byte(`{"metadata": {"description": "test"}}`),
			expected: "any",
			wantErr:  false,
		},
		{
			name:     "object with other fields but no type or $ref - any type",
			input:    []byte(`{"other": "value"}`),
			expected: "any",
			wantErr:  false,
		},
		{
			name:     "object with nullable but no type or $ref - any type",
			input:    []byte(`{"nullable": true}`),
			expected: "any",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unmarshalTypeOrRef(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalTypeOrRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("unmarshalTypeOrRef() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParameter_UnmarshalJSON_AnyType(t *testing.T) {
	// Test unmarshaling a parameter with 'any' type (empty object {})
	input := []byte(`{}`)
	var param Parameter
	err := param.UnmarshalJSON(input)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if param.Type != "any" { //nolint:goconst // "any" is a type name, not a reusable constant
		t.Errorf("Expected Type to be 'any', got %v", param.Type)
	}
}
