package types

import (
	"testing"
)

const secureObjectARMType = "secureObject" // ARM JSON uses mixed-case "secureObject"

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
	if param.Type != "any" {
		t.Errorf("Expected Type to be 'any', got %v", param.Type)
	}
}

func TestParameter_UnmarshalJSON_Secure(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectedType   string
		expectedSecure bool
	}{
		{
			name:           "securestring param",
			input:          []byte(`{"type":"securestring","metadata":{"description":"admin password"}}`),
			expectedType:   secureStringType,
			expectedSecure: true,
		},
		{
			name:           "secureObject param",
			input:          []byte(`{"type":"secureObject","metadata":{"description":"secret config"}}`),
			expectedType:   secureObjectARMType,
			expectedSecure: true,
		},
		{
			name:           "regular string param",
			input:          []byte(`{"type":"string","defaultValue":"westus"}`),
			expectedType:   "string",
			expectedSecure: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var param Parameter
			if err := param.UnmarshalJSON(tt.input); err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}
			if param.Type != tt.expectedType {
				t.Errorf("Type: got %q, want %q", param.Type, tt.expectedType)
			}
			if param.Secure != tt.expectedSecure {
				t.Errorf("Secure: got %v, want %v", param.Secure, tt.expectedSecure)
			}
		})
	}
}

func TestOutput_UnmarshalJSON_Secure(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectedType   string
		expectedSecure bool
	}{
		{
			name:           "securestring output",
			input:          []byte(`{"type":"securestring","value":"[parameters('adminPassword')]"}`),
			expectedType:   secureStringType,
			expectedSecure: true,
		},
		{
			name:           "secureObject output",
			input:          []byte(`{"type":"secureObject","value":"[parameters('secretConfig')]"}`),
			expectedType:   secureObjectARMType,
			expectedSecure: true,
		},
		{
			name:           "regular string output",
			input:          []byte(`{"type":"string","value":"westus"}`),
			expectedType:   "string",
			expectedSecure: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output Output
			if err := output.UnmarshalJSON(tt.input); err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}
			if output.Type != tt.expectedType {
				t.Errorf("Type: got %q, want %q", output.Type, tt.expectedType)
			}
			if output.Secure != tt.expectedSecure {
				t.Errorf("Secure: got %v, want %v", output.Secure, tt.expectedSecure)
			}
		})
	}
}

func TestUserDefinedDataType_UnmarshalJSON_Sealed(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectedSealed bool
	}{
		{
			name: "sealed type",
			input: []byte(`{
				"type": "object",
				"properties": {"name": {"type": "string"}},
				"additionalProperties": false,
				"metadata": {"description": "a sealed type"}
			}`),
			expectedSealed: true,
		},
		{
			name: "open type - no additionalProperties",
			input: []byte(`{
				"type": "object",
				"properties": {"host": {"type": "string"}},
				"metadata": {"description": "an open type"}
			}`),
			expectedSealed: false,
		},
		{
			name: "additionalProperties as object - not sealed",
			input: []byte(`{
				"type": "object",
				"additionalProperties": {"type": "string"}
			}`),
			expectedSealed: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uddt UserDefinedDataType
			if err := uddt.UnmarshalJSON(tt.input); err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}
			if uddt.Sealed != tt.expectedSealed {
				t.Errorf("Sealed: got %v, want %v", uddt.Sealed, tt.expectedSealed)
			}
		})
	}
}
