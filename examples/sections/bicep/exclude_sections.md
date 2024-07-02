# storage account

## Parameters

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| kind | string | The kind of storage account. | "StorageV2" |
| location | string | Location to deploy the storage account. | "[resourceGroup().location]" |
| name | string | Name of the storage account. |  |
| skuName | string | Name of the storage account's sku. | "Standard_LRS" |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| doubled | positiveInt (uddt) | Double test_number. |
| resourceId | string | Resource ID of the storage account. |
