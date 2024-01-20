/*
MIT License

Copyright (c) 2024 Christos Galanopoulos

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
bicep-docs is a command-line tool for scanning and updating the API version of Azure resources in bicep files.

It offers two main commands: scan and update.

# Scan

The scan command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and prints the results to stdout.

Example usage:

Scan a bicep file and print the results using the normal format:

	bruh scan --path ./main.bicep

Scan a directory and print only outdated resources using the table format:

	bruh scan --path ./bicep/modules --output table --outdated

Scan a directory including preview API versions and print the results using the markdown format:

	bruh scan --path ./bicep/modules --output markdown --include-preview

For full usage details, run `bruh scan --help` or `bruh help scan`.

# Update

The update command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and updates the file(s) in place or creates new ones with the "_updated.bicep" extension.

Example usage:

Update a bicep file in place:

	bruh update --path ./bicep/main.bicep --in-place

Create a new bicep file with the "_updated.bicep" extension:

	bruh update --path ./bicep/main.bicep

Update a directory in place including preview API versions:

	bruh update --path ./bicep/modules --in-place --include-preview

Use silent mode:

	bruh update --path ./bicep/main.bicep --silent

For full usage details, run `bruh update --help` or `bruh help update`.

Note: all the API versions are fetched from the official Microsoft Learn website (https://learn.microsoft.com/en-us/azure/templates/).
*/
package main

import (
	"os"

	"github.com/christosgalano/bicep-docs/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
