package cli

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"github.com/christosgalano/bicep-docs/internal/markdown"
	"github.com/christosgalano/bicep-docs/internal/template"
	"github.com/christosgalano/bicep-docs/internal/types"
)

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

	if err != nil {
		return err
	}

	// Wait for all goroutines to finish and return the first non-nil error
	return g.Wait()
}

func processBicepFile(bicepFile, markdownFile string) error {
	// Build Bicep template into ARM template
	armFile, err := template.BuildBicepTemplate(bicepFile)
	if err != nil {
		return fmt.Errorf("error processing %s: %w", bicepFile, err)
	}

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
