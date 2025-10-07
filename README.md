# bicep-docs

[![ci](https://github.com/christosgalano/bicep-docs/actions/workflows/ci.yaml/badge.svg?branch=main&event=push)](https://github.com/christosgalano/bicep-docs/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/christosgalano/bicep-docs)](https://goreportcard.com/report/github.com/christosgalano/bicep-docs)
[![Go Reference](https://pkg.go.dev/badge/github.com/christosgalano/bicep-docs.svg)](https://pkg.go.dev/github.com/christosgalano/bicep-docs)
[![Github Downloads](https://img.shields.io/github/downloads/christosgalano/bicep-docs/total.svg)](https://github.com/christosgalano/bicep-docs/releases)

![bicep-docs](assets/images/main-extra-small.png)

## Table of contents

- [Description](#description)
- [Installation](#installation)
- [Requirements](#requirements)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Description

**bicep-docs** is a command-line tool that generates documentation for Bicep templates.

## Installation

### Homebrew

```bash
brew tap christosgalano/christosgalano
brew install bicep-docs
```

### Chocolatey

```bash
choco install bicep-docs
```

### Go

```bash
go install github.com/christosgalano/bicep-docs/cmd/bicep-docs@latest
```

### Binary

Download the latest binary from the [releases page](https://github.com/christosgalano/bicep-docs/releases/latest).

## Requirements

To run bicep-docs, either the Azure CLI or the Bicep CLI must be [installed](https://learn.microsoft.com/en-us/azure/azure-resource-manager/bicep/install).

| CLI   | Minimum Required Version |
| ----- | ------------------------ |
| Azure | 2.77.0                   |
| Bicep | 0.38.0                   |

## Usage

bicep-docs is a command-line tool that generates documentation for Bicep templates.

Given an input Bicep file or directory, it parses the file(s) and generates corresponding documentation in Markdown format.

This can be used to automatically create and update documentation for your Bicep templates.

If the input is a directory, then for each `main.bicep` it will generate a `README.md` in the same directory. This happens recursively for all subdirectories.

If the input is a Bicep file, the output must be a file; otherwise, an error will be returned.

The default value for the output is `README.md`, relative to the directory where the command is executed.

**CAUTION:** If the Markdown file already exists, it will be **overwritten**.

### Arguments

Regarding the arguments `--include-sections` and `--exclude-sections`, the available sections are: `description`, `usage`, `modules`, `resources`, `parameters`, `udfs`, `uddts`, `variables`, `outputs`.

The default sections ordered are `description,usage,modules,resources,parameters,udfs,uddts,variables,outputs`. The default input for`--exclude-sections` is `''`.  This ensures backward compatibility with the previous version.

The order of the sections is respected when including them.

When excluding sections, the result will be the default sections minus the excluded ones (e.g. `--exclude-sections description,usage` will include `modules,resources,parameters,udfs,uddts,variables,outputs` in that order).

Both arguments cannot be provided at the same time, unless the `--include-sections` argument is the same as the default sections (e.g. `--include-sections description,usage,modules,resources,parameters,udfs,uddts,variables,outputs`).

The `--show-all-decorators` flag can be used to include additional columns in the documentation tables showing constraint information from Bicep decorators (allowed values, min/max constraints, exportable status, etc.). By default, these details are hidden to keep the documentation concise.

### Example usage

Parse a Bicep file and generate a Markdown file:

```bash
bicep-docs --input main.bicep --output readme.md
```

Parse a Bicep file and generate a README.md file in the same directory:

```bash
bicep-docs -i main.bicep
```

Parse a directory and generate a README.md file for each main.bicep file with verbose output:

```bash
bicep-docs -i ./bicep -V
```

Parse a Bicep file and generate a README.md excluding the user-defined sections:

```bash
bicep-docs --input main.bicep --exclude-sections udfs,uddts
```

Parse a Bicep file and generate a README.md including only the resources and modules in that order:

```bash
bicep-docs ---input main.bicep --include-sections resources,modules
```

Parse a Bicep file and generate comprehensive documentation with all decorator information:

```bash
bicep-docs --input main.bicep --show-all-decorators
```

More examples can be found [here](examples).

### Documentation format

The documentation has the following format:

```markdown
# module name | file name

## Description

...

## Usage

...

## Modules

table of modules

## Resources

table of resources

## Parameters

table of parameters

## User Defined Data Types (UDDTs)

table of UDDTs

For every UDDT u with properties, a sub-section is created:

### u

table of properties

...

## User Defined Functions (UDFs)

table of UDFS

## Variables

table of variables

## Outputs

table of outputs

```

### Handling of Loops

The tool follows these conventions when documenting resources, modules, variables, and outputs that use copy/loop constructs:

**Resources and Modules:**
- For resource/module arrays (using copy loops), only the base resource/module is documented with its symbolic name
- The documentation shows the type and description once, rather than documenting each iteration
- Example:
  ```bicep
  resource storageAccount 'Microsoft.Storage/storageAccounts@2023-01-01' = [for i in range(0, 5): {...}]
  module appService 'br/modules:app:v1' = [for env in environments: {...}]
  ```
  Each is documented as a single entry in their respective tables

**Variables and Outputs:**
- For array comprehension variables/outputs (using copy), the item is documented once with its description
- Both the array item and its transformed version are included in the documentation
- Example:
  ```bicep
  var storageConfigs = [for i in range(0, length(locations)): {...}]
  var ipFormatWithDuplicates = [for this in storageConfigs: {...}]

  output resourceIds array = [for i in range(0, length(locations)): storageAccount[i].id]
  output names array = [for config in storageConfigs: config.name]
  ```
  Each appears as a single entry in their respective tables

This approach keeps the documentation clean and focused on the logical structure rather than implementation details.

### Folder structure

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

## Contributing

Information about contributing to this project can be found [here](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).
