# test.bicep

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters

    // Optional parameters
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Resources

| Symbolic Name | Type | Retry On | Only If Not Exists | Description |
| --- | --- | --- | --- | --- |
| retry_resource | [Microsoft.Storage/storageAccounts](https://learn.microsoft.com/en-us/azure/templates/microsoft.storage/storageaccounts) | `['ServerError', 'Conflict'], 3` |  | This is a resource with retry behavior. |
| idempotent_resource | [Microsoft.Storage/storageAccounts](https://learn.microsoft.com/en-us/azure/templates/microsoft.storage/storageaccounts) |  | Yes | This is a resource that is only created if it does not exist. |
| plain_resource | [Microsoft.Web/serverfarms](https://learn.microsoft.com/en-us/azure/templates/microsoft.web/serverfarms) |  |  | This is a plain resource. |
