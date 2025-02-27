# storage account

## Description

Create a storage account.

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters
    name:

    // Optional parameters
    kind: 'StorageV2'
    location: '[resourceGroup().location]'
    skuName: 'Standard_LRS'
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Resources

| Symbolic Name | Type | Description |
| --- | --- | --- |
| st | [Microsoft.Storage/storageAccounts](https://learn.microsoft.com/en-us/azure/templates/microsoft.storage/storageaccounts) | This is a test resource. |

## Parameters

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| kind | string | The kind of storage account. | "StorageV2" |
| location | string | Location to deploy the storage account. | "[resourceGroup().location]" |
| name | string | Name of the storage account. |  |
| skuName | string | Name of the storage account's sku. | "Standard_LRS" |

## User Defined Data Types (UDDTs)

| Name | Type | Description | Properties |
| --- | --- | --- | --- |
| customType | object | Custom type for the input parameter. | [View Properties](#customtype) |
| positiveInt | int | Positive integer (> 0). |  |

### customType

| Name | Type | Description |
| --- | --- | --- |
| property1 | string | This is a string property of the custom type. |
| property2 | positiveInt (uddt) | This is a positive integer property of the custom type. |

## User Defined Functions (UDFs)

| Name | Description | Output Type |
| --- | --- | --- |
| double | Doubles a positive integer. | positiveInt (uddt) |

## Variables

| Name | Description |
| --- | --- |
| test_number |  |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| doubled | positiveInt (uddt) | Double test_number. |
| resourceId | string | Resource ID of the storage account. |
