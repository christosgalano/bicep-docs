/*
Package cli provides a command-line interface (CLI) for the bicep-docs tool, utilizing cobra-cli.
*/
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CLI flags.
var (
	input  string
	output string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Version: "v1.0.0",
	Use:     "bicep-docs",
	Short:   "bicep-docs is a command-line tool that generates documentation for Bicep templates.",
	Long: `bicep-docs is a command-line tool that generates documentation for Bicep templates.

Given an input Bicep file or directory, it parses the file(s) and generates a corresponding Markdown file with the extracted information.
This can be used to automatically generate and update documentation for your Bicep modules and resources.

If the input is a directory, it will recursively parse all main.bicep files in the directory and its subdirectories.
The output will be a corresponding README.md file in the same directory as the main.bicep file.

If the Markdown file already exists, it will be overwritten.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := processInput(input)
		if err != nil {
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
	// input - optional
	rootCmd.Flags().StringVarP(
		&input,
		"input",
		"i",
		".",
		"input Bicep file or directory",
	)

	// output - optional
	rootCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"README.md",
		"output Markdown file; ignored if input is a directory",
	)
}
