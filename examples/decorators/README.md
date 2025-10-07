# Decorators Example

This example demonstrates the `--show-all-decorators` feature of bicep-docs, which showcases all parameter decorators, constraints, and exportable status in the generated documentation.

## Structure

```
decorators/
├── README.md          # This file
└── bicep/
    ├── main.bicep     # Bicep template with various decorators
    ├── README.md      # Generated without --show-all-decorators (clean view)
    └── README_decorators.md  # Generated with --show-all-decorators (detailed view)
```

## Usage

### Generate clean documentation (default)

```bash
bicep-docs -i bicep/main.bicep -o bicep/README.md
```

### Generate detailed documentation with all decorators

```bash
bicep-docs -i bicep/main.bicep -o bicep/README_decorators.md --show-all-decorators
```

## Key Features Demonstrated

- **Parameter Constraints**: `@minLength`, `@maxLength`, `@minValue`, `@maxValue`, `@allowed`
- **Secure Parameters**: `@secure()` decorator
- **Exportable Types**: `@export()` decorator on custom types and functions
- **Output Constraints**: Decorators applied to outputs
- **Custom Types**: User-defined data types with constraints and exportable status

## Comparison

Compare the generated documentation files to see the difference:

- `README.md` - Clean, essential view (default)
- `README_decorators.md` - Detailed view with all constraint and exportable information
