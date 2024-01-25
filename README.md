# bicep-docs

[![ci](https://github.com/christosgalano/bicep-docs/actions/workflows/ci.yaml/badge.svg?branch=main&event=push)](https://github.com/christosgalano/bicep-docs/actions/workflows/ci.yaml)
[![Code Coverage](https://img.shields.io/badge/coverage-87.7%25-31C754)](https://img.shields.io/badge/coverage-87.7%25-31C754)
[![Go Report Card](https://goreportcard.com/badge/github.com/christosgalano/bicep-docs)](https://goreportcard.com/report/github.com/christosgalano/bicep-docs)
[![Go Reference](https://pkg.go.dev/badge/github.com/christosgalano/bicep-docs.svg)](https://pkg.go.dev/github.com/christosgalano/bicep-docs)

## Table of contents

- [Description](#description)
- [Installation](#installation)
- [Usage](#usage)
- [GitHub Action](#github-action)
- [License](#license)

## Description

**bicep-docs** is a command-line tool that generates documentation for Bicep templates.

## Installation

### Homebrew

```bash
brew tap christosgalano/christosgalano
brew install bicep-docs
```

### Go

```bash
go install github.com/christosgalano/bicep-docs/cmd/bicep-docs@latest
```

### Binary

Download the latest binary from the [releases page](https://github.com/christosgalano/bicep-docs/releases/latest).

## Usage

bicep-docs is a command-line tool that generates documentation for Bicep templates.

Given an input Bicep file or directory, it parses the file(s) and generates corresponding documentation in Markdown format.

This can be used to automatically create and update documentation for your Bicep templates.

If the input is a directory, then for each `main.bicep` it will generate a `README.md` in the same directory. This happens recursively for all subdirectories.

If the input is a Bicep file, the output must be a file; otherwise, an error will be returned.

The default value for the output is `README.md`, relative to the directory where the command is executed.

**CAUTION:** If the Markdown file already exists, it will be **overwritten**.

**NOTE:** To run bicep-docs, either the Azure CLI or the Bicep CLI must be [installed](https://learn.microsoft.com/en-us/azure/azure-resource-manager/bicep/install).

### Example usage

Parse a Bicep file and generate a Markdown file with verbose output:

```bash
bicep-docs --input main.bicep --output readme.md --verbose
```

Parse a Bicep file and generate a README.md file in the same directory:

```bash
bicep-docs -i main.bicep
```

Parse a directory and generate a README.md file for each main.bicep file:

```bash
bicep-docs -i ./bicep
```

Parse the current directory and generate a README.md file for each main.bicep file:

```bash
bicep-docs
```

### Bicep folder structure

This tool is extremely useful if you are following this structure for your Bicep projects:

```text
.
├── bicep
│   │
│   ├── modules
│   │   ├── compute
│   │   │   ├── main.bicep
│   │   │   └── README.md
│   │   └── ...
│   │
│   ├── environments
│   │   ├── development
│   │   │   ├── main.bicep
│   │   │   ├── main.bicepparam
│   │   │   └── README.md
│   │   └── ...
```

## GitHub Action

bicep-docs can also be used as a GitHub Action to generate the documentation for Bicep modules in a repository.

**NOTE:** It is a Container action, which is only supported on Linux runners.

### Syntax

```yaml
  uses: christosgalano/bicep-docs@v1.0.0
  with:
    input: ./bicep          # input file or directory
    output: README.md       # output Markdown file; if input is a directory, this parameter is ignored
    verbose: true | false   # verbose output (default: false)
```

### Examples

Generate the documentation for a Bicep module by providing the input and output parameters:

```yaml
- name: Generate documentation for main.bicep
  uses: christosgalano/bicep-docs@v1.0.0
  with:
    input: ./bicep/main.bicep
    output: ./bicep/readme.md
    verbose: true
```

Generate the documentation for all Bicep modules in the current directory:

```yaml
- name: Generate documentation for Bicep modules
  uses: christosgalano/bicep-docs@v1.0.0
```

A complete example can be found below. It consists of the following steps:

1. Checkout the repository
2. Generate the documentation for Bicep modules (a `README.md` for each `main.bicep` file)
3. Commit the changes - if any
4. Push the changes - if needed

```yaml
bicep-docs:
  runs-on: ubuntu-latest
  permissions:
    contents: write
  steps:
    - name: Checkout
      uses: actions/checkout@v4

    - name: Generate documentation for Bicep modules
      uses: christosgalano/bicep-docs@v1.0.0
      with:
        input: ./bicep # path relative to workspace
        verbose: true

    - name: Commit changes - if any
      id: check-changes
      run: |
        git config --local user.name "github-actions[bot]"
        git config --local user.email "github-actions[bot]@users.noreply.github.com"
        git diff --quiet --exit-code -- ./bicep || UPDATED=true
        if [[ "$UPDATED" == "true" ]]; then
          git add ./bicep
          echo "changed=true" >> $GITHUB_ENV
          git commit -m "Updated documentation for Bicep modules"
        fi

    - name: Push changes - if needed
      if: steps.check-changes.outputs.changed == 'true'
      uses: ad-m/github-push-action@master
      with:
        branch: ${{ github.ref }}
        github_token: ${{ secrets.GITHUB_TOKEN }}
```

## License

This project is licensed under the [MIT License](LICENSE).
