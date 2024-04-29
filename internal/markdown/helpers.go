package markdown

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/christosgalano/bicep-docs/internal/types"
)

// checkFileExists checks if a file exists and is not a directory.
// It returns true if the file exists, false otherwise, along with any error encountered.
func checkFileExists(filename string) (bool, error) {
	f, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to stat file %q: %w", filename, err)
	}
	if f.IsDir() {
		return false, fmt.Errorf("output %q is a directory", filename)
	}
	return true, nil
}

// readFileContent reads the content of a file.
// readFileContent reads the content of a file and returns it as a string.
// It takes a filename as input and returns the file content and any error encountered.
func readFileContent(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", filename, err)
	}
	return string(bytes), nil
}

// extractType extracts the type from a type string.
// If the type is a user defined data type, it returns the name of it.
func extractType(t string) string {
	if strings.HasPrefix(t, "#/definitions/") {
		split := strings.Split(t, "/")
		return split[len(split)-1] + " (uddt)"
	}
	return t
}

// extractDescription extracts the description from the given metadata and returns it.
// If the metadata or the description is nil, an empty string is returned.
func extractDescription(metadata *types.Metadata) string {
	description := ""
	if metadata != nil && metadata.Description != nil {
		description = strings.ReplaceAll(*metadata.Description, "\r\n", "<br>")
		description = strings.ReplaceAll(description, "\n", "<br>")
	}
	return description
}
