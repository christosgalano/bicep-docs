package cli

import (
	"strings"

	"github.com/christosgalano/bicep-docs/internal/types"
)

// convertStringsToSections converts a slice of strings to a slice of Section enums.
func convertStringsToSections(sections []string) ([]types.Section, error) {
	var convertedSections []types.Section
	for _, section := range sections {
		s, err := types.ParseSectionFromString(section)
		if err != nil {
			return nil, err
		}
		convertedSections = append(convertedSections, s)
	}
	return convertedSections, nil
}

// computeSectionDifference computes the difference between two sets of sections.
// It takes two comma-separated strings, includeSections and excludeSections,
// and returns a slice of sections that are included in includeSections but not
// in excludeSections.
func computeSectionDifference(includeSections, excludeSections string) ([]types.Section, error) {
	var includedSections, excludedSections []types.Section
	var err error

	if includeSections != "" {
		includedSections, err = convertStringsToSections(strings.Split(includeSections, ","))
		if err != nil {
			return nil, err
		}
	}
	if excludeSections != "" {
		excludedSections, err = convertStringsToSections(strings.Split(excludeSections, ","))
		if err != nil {
			return nil, err
		}
	}

	excludedSet := make(map[types.Section]struct{}, len(excludedSections))
	for _, section := range excludedSections {
		excludedSet[section] = struct{}{}
	}

	var difference []types.Section
	for _, section := range includedSections {
		if _, found := excludedSet[section]; !found {
			difference = append(difference, section)
		}
	}

	return difference, nil
}
