# Sections Example

This example demonstrates how to use the bicep-docs CLI to generate documentation for a single Bicep file,
providing the output file name and utilizing the sections' arguments.

## Include Sections

1. Navigate to the `file/bicep` directory.
2. Run the following command:

```bash
bicep-docs --input main.bicep --output include_sections.md --verbose --include-sections parameters,outputs
```

> NOTE: the order is respected when including sections.

## Exclude Sections


1. Navigate to the `file/bicep` directory.
2. Run the following command:

```bash
bicep-docs --input main.bicep --output exclude_sections.md --verbose --exclude-sections description,usage,resources,modules,udfs,uddts,variables
```

> NOTE: the sections incuded will be: default - excuded, where the default sections are: description, usage, modules, resources, parameters, udfs, uddts, variables, outputs.
