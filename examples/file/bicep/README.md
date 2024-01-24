# storage account

## Description

Create a storage account.

## Resources

| Symbolic Name | Type | Description |
| --- | --- | --- |
| st | [Microsoft.Storage/storageAccounts](https://learn.microsoft.com/en-us/azure/templates/microsoft.storage/storageaccounts) | This is a test resource. |

## Parameters

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| kind | string | The kind of storage account. | StorageV2 |
| location | string | Location to deploy the storage account. | [resourceGroup().location] |
| name | string | Name of the storage account. |  |
| skuName | string | Name of the storage account's sku. | Standard_LRS |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| resourceId | string | Resource ID of the storage account. |
