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

    // Optional parameters
    test_parameter: 'test'
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
| test_parameter | string | This is a test parameter. | "test" |

## User Defined Data Types (UDDTs)

| Name | Type | Description | Properties |
| --- | --- | --- | --- |
| pint | positiveInt (uddt) | This is a user defined type (alias). |  |
| positiveInt | int | This is a user defined type. |  |
| custom_type | object | This is a user defined type with properties. | [View Properties](#custom_type) |

### custom_type

| Name | Type | Description |
| --- | --- | --- |
| property1 | string | This is a property of a user defined type. |
| property_2 | positiveInt (uddt) | This is another property of a user defined type which uses ref. |

## User Defined Functions (UDFs)

| Name | Description |
| --- | --- |
| buildUrl | This is a user defined function. |
| double | This is a user defined function with uddts. |

## Variables

| Name | Description |
| --- | --- |
| test_variable | This is a test variable. |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| test_output | positiveInt (uddt) | This is a test output. |
