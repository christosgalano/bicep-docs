package cli

import (
	"errors"
	"reflect"
	"testing"

	"github.com/christosgalano/bicep-docs/internal/types"
)

func TestConvertStringsToSections(t *testing.T) {
	tests := []struct {
		name           string
		sections       []string
		expectedResult []types.Section
		expectedError  error
	}{
		{
			name:           "valid_sections",
			sections:       []string{"description", "usage", "modules"},
			expectedResult: []types.Section{types.DescriptionSection, types.UsageSection, types.ModulesSection},
			expectedError:  nil,
		},
		{
			name:           "invalid_section",
			sections:       []string{"description", "invalid", "modules"},
			expectedResult: nil,
			expectedError:  errors.New("invalid section: \"invalid\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertStringsToSections(tt.sections)
			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("convertStringsToSections() result = %v, expected %v", result, tt.expectedResult)
			}
			if (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) ||
				(err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError == nil) {
				t.Errorf("convertStringsToSections() error = %v, expected %v", err, tt.expectedError)
			}
		})
	}
}

func TestComputeSectionDifference(t *testing.T) {
	tests := []struct {
		name            string
		includeSections string
		excludeSections string
		expectedResult  []types.Section
		expectedError   error
	}{
		{
			name:            "included_sections_only",
			includeSections: "description,usage,modules",
			excludeSections: "",
			expectedResult:  []types.Section{types.DescriptionSection, types.UsageSection, types.ModulesSection},
			expectedError:   nil,
		},
		{
			name:            "excluded_sections_only",
			includeSections: defaultSections,
			excludeSections: "description,usage,modules",
			expectedResult: []types.Section{
				types.ResourcesSection,
				types.ParametersSection,
				types.UserDefinedDataTypesSection,
				types.UserDefinedFunctionsSection,
				types.VariablesSection,
				types.OutputsSection,
			},
			expectedError: nil,
		},
		{
			name:            "invalid_section",
			includeSections: "description,invalid,modules",
			excludeSections: "",
			expectedResult:  nil,
			expectedError:   errors.New("invalid section: \"invalid\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := computeSectionDifference(tt.includeSections, tt.excludeSections)
			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("computeSectionDifference() result = %v, expected %v", result, tt.expectedResult)
			}
			if tt.expectedError == nil {
				if err != nil {
					t.Errorf("computeSectionDifference() error = %v, expected %v", err, tt.expectedError)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("computeSectionDifference() error = %v, expected %v", err, tt.expectedError)
				}
			}
		})
	}
}
