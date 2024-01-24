# development

## Description

Deploy to the development environment.

## Modules

| Symbolic Name | Source | Description |
| --- | --- | --- |
| identity | ../../modules/identity/main.bicep | Create a User Assigned Identity. |

## Resources

| Symbolic Name | Type | Description |
| --- | --- | --- |
| rg | [Microsoft.Resources/resourceGroups](https://learn.microsoft.com/en-us/azure/templates/microsoft.resources/resourcegroups) | Resource Group to deploy the resources. |

## Parameters

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| location | string | Location to deploy the resources. | westeurope |
| rgName | string | Name of the resource group to deploy the resources. | rgtest123 |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| identityClientId | string | Client ID of the User Assigned Identity. |
