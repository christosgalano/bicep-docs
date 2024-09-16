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
func generateDocsFromDirectory(dirPath string, verbose bool, sections []types.Section) error {
	const maxGoRoutines = 10
	g := new(errgroup.Group)
	g.SetLimit(maxGoRoutines)

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

// func generateDocsFromDirectory(dirPath string, verbose bool, sections []types.Section) error {
// 	numWorkers := getOptimalWorkerCount()
// 	bufferSize := getOptimalBufferSize(numWorkers)

// 	jobs := make(chan string, bufferSize)
// 	results := make(chan error, bufferSize)

// 	// Start worker pool
// 	for w := 1; w <= numWorkers; w++ {
// 		go processFile(jobs, results, verbose, sections)
// 	}

// 	// Walk the directory and enqueue the paths of all 'main.bicep' files
// 	go func() {
// 		err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
// 			if err != nil {
// 				return err
// 			}
// 			if !d.IsDir() && d.Name() == "main.bicep" {
// 				jobs <- path
// 			}
// 			return nil
// 		})

// 		// If an error occurs while walking the directory, send it to the results channel
// 		if err != nil {
// 			results <- err
// 		}

// 		// Close the jobs channel to signal workers that no more jobs will be enqueued
// 		close(jobs)
// 	}()

// 	// Collect results
// 	for i := 0; i < numWorkers; i++ {
// 		if err := <-results; err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// processFile is a worker function that processes a single Bicep file and generates its documentation.
// It reads file paths from the jobs channel and sends any errors to the results channel.
//
// Note: This function assumes that the paths provided in the jobs channel are valid Bicep files.
func processFile(jobs <-chan string, results chan<- error, verbose bool, sections []types.Section) {
	for path := range jobs {
		markdownFile := filepath.Join(filepath.Dir(path), "README.md")
		err := generateDocsFromBicepFile(path, markdownFile, verbose, sections)
		if err != nil {
			results <- err
			return
		}
	}
	results <- nil
}

// getOptimalWorkerCount returns the optimal number of workers based on the number of CPU cores.
//
//nolint:mnd // Default values for the number of workers.
func getOptimalWorkerCount() int {
	numCPU := runtime.GOMAXPROCS(0)
	switch {
	case numCPU <= 2:
		return 4
	case numCPU <= 8:
		return numCPU
	default:
		workers := numCPU * 3 / 4
		if workers > 16 {
			workers = 16
		}
		return workers
	}
}

// getOptimalBufferSize returns the optimal buffer size for the job channel based on the number of workers.
//
//nolint:mnd // Default values for the buffer size.
func getOptimalBufferSize(numWorkers int) int {
	bufferSize := 100 * numWorkers / 2

	// Ensure a minimum buffer size
	if bufferSize < 100 {
		return 100
	}

	// Cap the buffer size at 500
	if bufferSize > 500 {
		return 500
	}

	return bufferSize
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
