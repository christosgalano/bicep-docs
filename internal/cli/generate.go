package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sync/errgroup"

	"github.com/christosgalano/bicep-docs/internal/markdown"
	"github.com/christosgalano/bicep-docs/internal/template"
	"github.com/christosgalano/bicep-docs/internal/types"
)

// GenerateDocs generates documentation based on the input file or directory.
//
// If the input is a directory, it generates documentation for all 'main.bicep' files in the directory.
// If the input is a Bicep file, it generates documentation for that file only.
//
// The output is used only when the input is a Bicep file; in other cases it is always set to 'README.md'.
//
// The sections slice contains the sections that should be included in the documentation.
//
// If verbose is true, additional information will be printed during the generation process.
func GenerateDocs(input, output string, verbose bool, sections []types.Section) error {
	f, err := os.Stat(input)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no such file or directory %q", input)
		}
		return err
	}

	if f.IsDir() {
		return generateDocsFromDirectory(input, verbose, sections)
	}
	return generateDocsFromBicepFile(input, output, verbose, sections)
}

// generateDocsFromDirectory processes the directory and its subdirectories recursively.
//
// For each 'main.bicep' file, it creates/updates a 'README.md' file in the same directory.
//
//nolint:mnd // Sensible default.
func generateDocsFromDirectory(dirPath string, verbose bool, sections []types.Section) error {
	g := new(errgroup.Group)
	g.SetLimit(runtime.GOMAXPROCS(0) * 10)

	// Traverse the directory and process each main.bicep file
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "main.bicep" {
			// Create a README.md file in the same directory as the main.bicep file
			markdownFile := filepath.Join(filepath.Dir(path), "README.md")
			g.Go(func() error {
				return generateDocsFromBicepFile(path, markdownFile, verbose, sections)
			})
		}
		return nil
	})
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
// the template's resources, modules, parameters, user defined data types,
// user defined functions, variables, outputs, and metadata.
//
// Finally it creates a corresponding Markdown file based on the gathered information
// and the provided section, while also deleting the ARM template.
//
// If the Markdown file already exists, it will be overwritten.
func generateDocsFromBicepFile(bicepFile, markdownFile string, verbose bool, sections []types.Section) error {
	// Build Bicep template into ARM template
	armFile, err := template.BuildBicepTemplate(bicepFile)
	if err != nil {
		return err
	}
	defer os.Remove(armFile)

	// Parse both Bicep and ARM templates
	var tmpl *types.Template
	tmpl, err = template.ParseTemplates(bicepFile, armFile)
	if err != nil {
		return fmt.Errorf("error processing %s: %w", bicepFile, err)
	}

	// Create/Update Markdown file
	if err := markdown.CreateFile(markdownFile, tmpl, verbose, sections); err != nil {
		return fmt.Errorf("error processing %s: %w", bicepFile, err)
	}

	return nil
}
