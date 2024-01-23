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

// processInput processes the input file or directory.
func processInput(input string) error {
	// Invalid file
	fileInfo, err := os.Stat(input)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no such file or directory %q: %w", input, err)
		}
		return fmt.Errorf("failed to stat file %q: %w", input, err)
	}

	if fileInfo.IsDir() {
		return processDirectory(input)
	}
	return processBicepFile(input, output)
}

// processDirectory processes the directory and its subdirectories recursively.
//
// For each main.bicep file, it creates a README.md file in the same directory.
func processDirectory(dirPath string) error {
	// Create a new errgroup with a limit of 10 goroutines
	g := new(errgroup.Group)
	g.SetLimit(10)

	// Traverse the directory and process each main.bicep file
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "main.bicep" {
			markdownFile := filepath.Join(filepath.Dir(path), "README.md")
			g.Go(func() error {
				return processBicepFile(path, markdownFile)
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

// processBicepFile processes a Bicep template and creates
// a corresponding Markdown file.
//
// First it builds the Bicep template into an ARM template.
//
// Then it parses both Bicep and ARM templates gathering information about
// the template's resources, modules, parameters, outputs, and metadata.
//
// Finally it creates a corresponding Markdown file based on the gathered information
// and deletes the ARM template.
func processBicepFile(bicepFile, markdownFile string) error {
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
	if err := markdown.CreateFile(markdownFile, t); err != nil {
		return fmt.Errorf("error processing %s: %w", bicepFile, err)
	}

	return nil
}
