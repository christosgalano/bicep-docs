# test

## Description

This is a test template.

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters
    required:
    pint_array:
    simple_array:

    // Optional parameters
    nullable: null
    optional: 'test'
    string_array: []
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Modules

| Symbolic Name | Source | Description |
| --- | --- | --- |
| test_module | ./modules/test_module/main.bicep | This is a test module. |

## Resources

| Symbolic Name | Type | Description |
| --- | --- | --- |
| test_resource | [Microsoft.Storage/storageAccounts](https://learn.microsoft.com/en-us/azure/templates/microsoft.storage/storageaccounts) | This is a test resource. |

## Parameters

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| required | string | This is a required parameter. |  |
| nullable | string | This is a nullable parameter. | null |
| optional | string | This is an optional parameter. | "test" |
| string_array | string[] | This is a string array parameter. | [] |
| pint_array | positive_int[] (uddt) | This is a positive int custom type array parameter. |  |
| simple_array | array | This is simple array parameter. |  |

## User Defined Data Types (UDDTs)

| Name | Type | Description | Properties |
| --- | --- | --- | --- |
| pint | positive_int (uddt) | This is a user defined type (alias). |  |
| positive_int | int | This is a user defined type. |  |
| string_array | string[] | This is a user defined type with array items. |  |
| custom_type | object | This is a user defined type with properties. | [View Properties](#custom_type) |

### custom_type

| Name | Type | Description |
| --- | --- | --- |
| property1 | string | This is a property of a user defined type. |
| property_2 | positive_int (uddt) | This is another property of a user defined type which uses ref. |
| property_3 | positive_int[] (uddt) | This is a property of a user defined type with array items. |

## User Defined Functions (UDFs)

| Name | Description | Output Type |
| --- | --- | --- |
| build_url | This is a user defined function. | string |
| double | This is a user defined function with uddts. | positive_int (uddt) |
| get_string_array | This is a user defined function with array items as output. | string[] |

## Variables

| Name | Description |
| --- | --- |
| test_variable | This is a test variable. |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| pint | positive_int (uddt) | This is an output with uddt. |
| pint_array | positive_int[] (uddt) | This is an output with uddt array items. |
