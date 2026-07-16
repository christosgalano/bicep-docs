# storage account

## Parameters

| Name | Status | Type | Description | Default |
| --- | --- | --- | --- | --- |
| kind | Optional | string | The kind of storage account. | "StorageV2" |
| location | Optional | string | Location to deploy the storage account. | "[resourceGroup().location]" |
| name | Required | string | Name of the storage account. |  |
| skuName | Optional | string | Name of the storage account's sku. | "Standard_LRS" |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| doubled | positiveInt (uddt) | Double test_number. |
| resourceId | string | Resource ID of the storage account. |
