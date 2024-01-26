# test

## Description

This is a test template.

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
| test_parameter | string | This is a test parameter. | test |

## User Defined Data Types (UDDTs)

| Name | Type | Description |
| --- | --- | --- |
| pint | positiveInt (uddt) | This is a user defined type (alias). |
| positiveInt | int | This is a user defined type. |

## User Defined Functions (UDFs)

| Name | Description |
| --- | --- |
| buildUrl | This is a user defined function. |
| double | This is a user defined function with uddts. |

## Variables

| Name | Description |
| --- | --- |
| test_variable | |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| test_output | positiveInt (uddt) | This is a test output. |
