# test.bicep

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters
    adminPassword:
    secretConfig:

    // Optional parameters
    location: 'westus'
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Parameters

| Name | Status | Type | Description | Default |
| --- | --- | --- | --- | --- |
| adminPassword | Required | string (secure) | Admin password for the virtual machine. |  |
| secretConfig | Required | object (secure) | Secret configuration object. |  |
| location | Optional | string | Non-secure parameter for comparison. | "westus" |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| connectionString | string (secure) | The generated connection string. |
| secretObject | object (secure) | The secret configuration output. |
| region | string | Non-secure output for comparison. |
