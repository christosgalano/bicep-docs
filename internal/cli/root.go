/*
Package cli provides a command-line interface (CLI) for the bicep-docs tool, utilizing cobra-cli.
*/
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/christosgalano/bicep-docs/internal/types"
)

// CLI flags.
var (
	input           string
	output          string
	verbose         bool
	includeSections string
	excludeSections string
)

// CLI variables.
var (
	sections []types.Section
)

// CLI constants.
const (
	defaultSections = "description,usage,modules,resources,parameters,uddts,udfs,variables,outputs"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Version: "v1.3.1",
	Use:     "bicep-docs",
	Short:   "bicep-docs is a command-line tool that generates documentation for Bicep templates.",
	Long: `bicep-docs is a command-line tool that generates documentation for Bicep templates.

It parses Bicep files or directories to produce Markdown documentation. For directories,
it processes all main.bicep files, creating README.md in each directory containing a main.bicep file.
For single Bicep files, it generates a README.md in the same directory unless an output path is specified.
Existing README.md files will be overwritten.

Azure CLI or Bicep CLI need to be installed.
`,
	//revive:disable:unused-parameter
	Run: func(cmd *cobra.Command, args []string) {
		if err := GenerateDocs(input, output, verbose, sections); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// init initializes the root command.
func init() {
	// input - required
	rootCmd.Flags().StringVarP(
		&input,
		"input",
		"i",
		"",
		"input Bicep file or directory",
	)
	if err := rootCmd.MarkFlagRequired("input"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// output - optional
	rootCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"README.md",
		"output Markdown file; ignored if input is a directory",
	)

	// verbose - optional
	rootCmd.Flags().BoolVarP(
		&verbose,
		"verbose",
		"V",
		false,
		"verbose output",
	)

	// include-sections - optional
	rootCmd.Flags().StringVarP(
		&includeSections,
		"include-sections",
		"I",
		defaultSections,
		"comma-separated list of sections to include in the output, order matters",
		// "available sections: description, usage, modules, resources, parameters, uddts, "+
		// "udfs, variables, outputs",
	)

	// exclude-sections - optional
	rootCmd.Flags().StringVarP(
		&excludeSections,
		"exclude-sections",
		"E",
		"",
		"comma-separated list of sections to exclude from the default output; "+
			"available sections: description, usage, modules, resources, parameters, uddts, udfs, variables, outputs",
	)

	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Check for mutual exclusivity of include and exclude flags
		if includeSections != defaultSections && excludeSections != "" {
			return fmt.Errorf("include and exclude arguments cannot be provided simultaneously (even with default input for include)")
		}

		// Find the difference between include and exclude sections.
		// This will be the final list of sections to include in the output.
		var err error
		sections, err = computeSectionDifference(includeSections, excludeSections)
		if err != nil {
			return err
		}

		return nil
	}
}
