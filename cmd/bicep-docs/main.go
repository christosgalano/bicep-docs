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
bicep-docs is a command-line tool that generates documentation for Bicep templates.

Given an input Bicep file or directory, it parses the file(s) and generates a corresponding Markdown file with the extracted information.
This can be used to automatically generate and update documentation for your Bicep modules and resources.

If the input is a directory, it will recursively parse all main.bicep files in the directory and its subdirectories.
The output will be a corresponding README.md file in the same directory as the main.bicep file.

CAUTION: If the Markdown file already exists, it will be overwritten.

Example usage:

Parse a Bicep file and generate a Markdown file:

	bicep-docs -i main.bicep -o readme.md

Parse a Bicep file and generate a README.md file in the same directory:

	bicep-docs -i main.bicep

Parse a directory and generate a README.md file for each main.bicep file:

	bicep-docs -i ./bicep

For full usage details, run `bicep-docs --help`.
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
