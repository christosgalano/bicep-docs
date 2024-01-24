# Directory Example

This example demonstrates how to use the bicep-docs CLI to generate documentation for a directory of Bicep files.

The directory structure is as follows:

```text
bicep
├── modules
│   └── identity
│       └── main.bicep
└── environments
    └── development
        └── main.bicep
```

## Running the Example

1. Navigate to the `directory/bicep` directory.
2. Run the following command:

```bash
bicep-docs --input ./ --verbose
```
