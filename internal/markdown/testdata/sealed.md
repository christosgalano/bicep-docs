# test.bicep

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters
    config:
    openConfig:

    // Optional parameters
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Parameters

| Name | Status | Type | Description | Default |
| --- | --- | --- | --- | --- |
| config | Required | SealedConfig (uddt) | Parameter using the sealed type. |  |
| openConfig | Required | OpenConfig (uddt) | Parameter using the open type. |  |

## User Defined Data Types (UDDTs)

| Name | Type | Description | Properties |
| --- | --- | --- | --- |
| SealedConfig | object | A sealed configuration type that does not allow extra properties. | [View Properties](#sealedconfig) |
| OpenConfig | object | An unsealed type for comparison. | [View Properties](#openconfig) |

### SealedConfig

| Name | Type | Description |
| --- | --- | --- |
| name | string | The name of the resource. |
| value | int | The numeric value. |

### OpenConfig

| Name | Type | Description |
| --- | --- | --- |
| host | string | The host address. |
| port | int | The port number. |
