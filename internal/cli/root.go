/*
Package cli provides a command-line interface (CLI) for the bicep-docs tool, utilizing cobra-cli.
*/
package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	input  string
	output string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "bicep-docs",
	Short: "bicep-docs is a command-line tool ...",
	Long: `bicep-docs

bicep-docs is a command-line tool for ...

It can be used to ...`,
	Run: func(cmd *cobra.Command, args []string) {
		// Invalid file
		fs, err := os.Stat(input)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "Error: no such file or directory %q\n", input)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(1)
		}

		if fs.IsDir() {
			if err = processDirectory(input); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		} else {
			if err = processBicepFile(input, output); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
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
		"input Bicep file or directory (default: current directory))",
	)

	// output - optional
	rootCmd.Flags().StringVarP(
		&output,
		"output",
		"o",
		"README.md",
		"output Markdown file (default: README.md); ignored if input is a directory, where output is always README.md",
	)

	// version
	rootCmd.Version = "v1.0.0"
}
