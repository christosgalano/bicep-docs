# bicep-docs

[![ci](https://github.com/christosgalano/bicep-docs/actions/workflows/ci.yaml/badge.svg?branch=main&event=push)](https://github.com/christosgalano/bicep-docs/actions/workflows/ci.yaml)
[![Code Coverage](https://img.shields.io/badge/coverage-89.2%25-31C754)](https://img.shields.io/badge/coverage-89.2%25-31C754)
[![Go Report Card](https://goreportcard.com/badge/github.com/christosgalano/bicep-docs)](https://goreportcard.com/report/github.com/christosgalano/bicep-docs)
[![Go Reference](https://pkg.go.dev/badge/github.com/christosgalano/bicep-docs.svg)](https://pkg.go.dev/github.com/christosgalano/bicep-docs)
[![Github Downloads](https://img.shields.io/github/downloads/christosgalano/bicep-docs/total.svg)](https://github.com/christosgalano/bicep-docs/releases)

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

### Go

```bash
go install github.com/christosgalano/bicep-docs/cmd/bicep-docs@latest
```

### Binary

Download the latest binary from the [releases page](https://github.com/christosgalano/bicep-docs/releases/latest).

## Requirements

To run bicep-docs, either the Azure CLI or the Bicep CLI must be [installed](https://learn.microsoft.com/en-us/azure/azure-resource-manager/bicep/install).

| CLI | Minimum Required Version |
| --- | --- |
| Azure | 2.56.0 |
| Bicep | 0.24.24 |

Currently, you need to modify your `bicepconfig.json` to enable the *userDefinedFunctions* feature. More information [here](https://learn.microsoft.com/en-us/azure/azure-resource-manager/bicep/user-defined-functions#enable-the-preview-feature).

## Usage

bicep-docs is a command-line tool that generates documentation for Bicep templates.

Given an input Bicep file or directory, it parses the file(s) and generates corresponding documentation in Markdown format.

This can be used to automatically create and update documentation for your Bicep templates.

If the input is a directory, then for each `main.bicep` it will generate a `README.md` in the same directory. This happens recursively for all subdirectories.

If the input is a Bicep file, the output must be a file; otherwise, an error will be returned.

The default value for the output is `README.md`, relative to the directory where the command is executed.

**CAUTION:** If the Markdown file already exists, it will be **overwritten**.

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

More examples can be found [here](examples).

### Documentation format

The documentation has the following format:

```markdown
# module name | file name

## Description

...

## Modules

table of modules

## Resources

table of resources

## Parameters

table of parameters

## User Defined Data Types (UDDTs)

table of UDDTs

## User Defined Functions (UDFs)

table of UDFS

## Variables

table of variables

## Outputs

table of outputs

```

**NOTE:** The **Description** column of **variables** is left intentionally blank, as it is not currently possible to have metadata for variables.

<!-- markdownlint-disable -->
<details>
  <summary>Example documentation file:</summary>

# module-name

## Description

Module description

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters
    location:
    required_array:
    required_bool:
    required_int:
    required_object:
    required_string:

    // Optional parameters
    optional_array: [
      'value1'
      '[__bicep.double(2)]'
      true
      {
        key1: 'test'
        key2: 42
        key3: true
      }
    ]
    optional_bool: true
    optional_int: 42
    optional_object: {
      key1: 'value1'
      key2: '[__bicep.double(2)]'
      key3: true
      key4: [
        'a'
        'b'
        'c'
      ]
    }
    optional_string: 'default'
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Resources

| Symbolic Name | Type | Description |
| --- | --- | --- |
| identity | [Microsoft.ManagedIdentity/userAssignedIdentities](https://learn.microsoft.com/en-us/azure/templates/microsoft.managedidentity/userassignedidentities) |  |

## Parameters

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| location | string | Resource group location. |  |
| optional_array | array | Optional array parameter. | ["value1", "[__bicep.double(2)]", true, {"key1": "test", "key2": 42, "key3": true}] |
| optional_bool | bool | Optional boolean parameter. | true |
| optional_int | int | Optional integer parameter. | 42 |
| optional_object | object | Optional object parameter. | {"key1": "value1", "key2": "[__bicep.double(2)]", "key3": true, "key4": ["a", "b", "c"]} |
| optional_string | string | Optional string parameter. | "default" |
| required_array | array | Required array parameter. |  |
| required_bool | bool | Required boolean parameter. |  |
| required_int | int | Required integer parameter. |  |
| required_object | object | Required object parameter. |  |
| required_string | string | Required string parameter. |  |

## User Defined Data Types (UDDTs)

| Name | Type | Description |
| --- | --- | --- |
| positiveInt | int | Positive integer (> 0). |

## User Defined Functions (UDFs)

| Name | Description |
| --- | --- |
| double | Doubles a positive integer. |

## Variables

| Name | Description |
| --- | --- |
| doubleValue | Doubles the required integer. |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| clientId | string | Client ID of the identity. |
| principalId | string | Principal ID of the identity. |
| resourceId | string | Resource ID of the identity. |

</details>
<!-- markdownlint-enable -->

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

## Contributing

Information about contributing to this project can be found [here](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).
