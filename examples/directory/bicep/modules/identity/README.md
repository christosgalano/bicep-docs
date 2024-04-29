# identity

## Description

Creates a User Assigned Identity.

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters
    identityName:

    // Optional parameters
    location: '[resourceGroup().location]'
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Resources

| Symbolic Name | Type | Description |
| --- | --- | --- |
| identity | [Microsoft.ManagedIdentity/userAssignedIdentities](https://learn.microsoft.com/en-us/azure/templates/microsoft.managedidentity/userassignedidentities) |  |

## Parameters

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| identityName | string | Name of the identity. |  |
| location | string | Location of the identity. | "[resourceGroup().location]" |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| clientId | string | Client ID of the identity. |
| principalId | string | Principal ID of the identity. |
| resourceId | string | Resource ID of the identity. |
