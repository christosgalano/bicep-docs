# identity

## Description

Creates a User Assigned Identity.

## Resources

| Symbolic Name | Type | Description |
| --- | --- | --- |
| identity | [Microsoft.ManagedIdentity/userAssignedIdentities](https://learn.microsoft.com/en-us/azure/templates/microsoft.managedidentity/userassignedidentities) |  |

## Parameters

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| identityName | string | Name of the identity. |  |
| location | string | Location of the identity. | [resourceGroup().location] |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| clientId | string | Client ID of the identity. |
| principalId | string | Principal ID of the identity. |
| resourceId | string | Resource ID of the identity. |
