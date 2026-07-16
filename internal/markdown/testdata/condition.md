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

## Modules

| Symbolic Name | Source | Condition | Description |
| --- | --- | --- | --- |
| conditional_module | ./modules/test_module/main.bicep | `deploy` | This is a conditional module. |
| plain_module | ./modules/other/main.bicep |  | This is a plain module. |

## Resources

| Symbolic Name | Type | Condition | Description |
| --- | --- | --- | --- |
| conditional_resource | [Microsoft.Storage/storageAccounts](https://learn.microsoft.com/en-us/azure/templates/microsoft.storage/storageaccounts) | `deploy \|\| force` | This is a conditional resource. |
| plain_resource | [Microsoft.Web/serverfarms](https://learn.microsoft.com/en-us/azure/templates/microsoft.web/serverfarms) |  | This is a plain resource. |
