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
    test_parameter1: 'test'
    test_parameter2: {}
    test_parameter3: {
      key1: 'value1'
    }
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
| test_parameter1 | string | This is a test parameter. | "test" |
| test_parameter2 | object | This is a test parameter. | {} |
| test_parameter3 | object | This is a test parameter. | {"key1": "value1"} |

## Variables

| Name | Description |
| --- | --- |
| test_variable | This is a test variable. |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| test_output | string | This is a test output. |
