package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"github.com/christosgalano/bicep-docs/internal/markdown"
	"github.com/christosgalano/bicep-docs/internal/template"
	"github.com/christosgalano/bicep-docs/internal/types"
)

// generateDocs creates/updates Markdown documentation for the given input.
func generateDocs(input, output string, verbose bool) error {
	// Non-existing file or directory
	f1, err := os.Stat(input)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no such file or directory %q: %w", input, err)
		}
		return fmt.Errorf("failed to stat file %q: %w", input, err)
	}

	// output is a directory
	f2, err := os.Stat(output)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to stat file %q: %w", output, err)
		}
	} else if f2.IsDir() {
		return fmt.Errorf("output %q is not a file", output)
	}

	if f1.IsDir() {
		return generateDocsFromDirectory(input, verbose)
	}
	return generateDocsFromBicepFile(input, output, verbose)
}

// generateDocsFromDirectory processes the directory and its subdirectories recursively.
//
// For each main.bicep file, it creates/updates a README.md file in the same directory.
func generateDocsFromDirectory(dirPath string, verbose bool) error {
	// Create a new errgroup with a limit of 10 goroutines
	g := new(errgroup.Group)
	g.SetLimit(10)

	// Traverse the directory and process each main.bicep file
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "main.bicep" {
			// Create a README.md file in the same directory as the main.bicep file
			markdownFile := filepath.Join(filepath.Dir(path), "README.md")
			g.Go(func() error {
				return generateDocsFromBicepFile(path, markdownFile, verbose)
			})
		}
		return nil
	})

	// WalkDir error
	if err != nil {
		return err
	}

	// Wait for all goroutines to finish and return the first non-nil error
	return g.Wait()
}

// generateDocsFromBicepFile processes a Bicep template and creates/updates
// a/the corresponding Markdown file.
//
// First it builds the Bicep template into an ARM template.
//
// Then it parses both Bicep and ARM templates gathering information about
// the template's resources, modules, parameters, outputs, and metadata.
//
// Finally it creates a corresponding Markdown file based on the gathered information
// and deletes the ARM template.
//
// If the Markdown file already exists, it will be overwritten.
func generateDocsFromBicepFile(bicepFile, markdownFile string, verbose bool) error {
	// If output is a directory, set the output file name to "README.md"
	f, err := os.Stat(markdownFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to stat file %q: %w", markdownFile, err)
	}
	if err == nil && f.IsDir() {
		return fmt.Errorf("output %q is a directory", markdownFile)
	}

	// Build Bicep template into ARM template
	armFile, err := template.BuildBicepTemplate(bicepFile)
	if err != nil {
		return fmt.Errorf("error processing %s: %w", bicepFile, err)
	}
	defer os.Remove(armFile)

	// Parse both Bicep and ARM templates
	var t *types.Template
	t, err = template.ParseTemplates(bicepFile, armFile)
	if err != nil {
		return fmt.Errorf("error processing %s: %w", bicepFile, err)
	}

	// Generate Markdown file
	if err := markdown.CreateFile(markdownFile, t, verbose); err != nil {
		return fmt.Errorf("error processing %s: %w", bicepFile, err)
	}

	return nil
}