# development

## Description

Deploy to the development environment.

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters

    // Optional parameters
    location: 'westeurope'
    rgName: 'rgtest123'
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

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
| location | string | Location to deploy the resources. | "westeurope" |
| rgName | string | Name of the resource group to deploy the resources. | "rgtest123" |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| identityClientId | string | Client ID of the User Assigned Identity. |
