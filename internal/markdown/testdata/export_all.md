# test.bicep

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters
    config:
    settings:

    // Optional parameters
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Parameters

| Name | Status | Type | Description | Default | Allowed Values | Min Length | Max Length | Min Value | Max Value |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| config | Required | ExportedConfig (uddt) | Parameter using the exported type. |  |  |  |  |  |  |
| settings | Required | InternalConfig (uddt) | Parameter using the internal type. |  |  |  |  |  |  |

## User Defined Data Types (UDDTs)

| Name | Type | Description | Sealed | Exportable | Properties | Min Length | Max Length | Min Value | Max Value |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ExportedConfig | object | An exported configuration type. |  | Yes | [View Properties](#exportedconfig) |  |  |  |  |
| InternalConfig | object | An internal type (not exported). |  |  | [View Properties](#internalconfig) |  |  |  |  |

### ExportedConfig

| Name | Type | Description | Allowed Values | Min Length | Max Length | Min Value | Max Value |
| --- | --- | --- | --- | --- | --- | --- | --- |
| name | string | The resource name. |  |  |  |  |  |
| value | int | The numeric value. |  |  |  |  |  |

### InternalConfig

| Name | Type | Description | Allowed Values | Min Length | Max Length | Min Value | Max Value |
| --- | --- | --- | --- | --- | --- | --- | --- |
| host | string | The host address. |  |  |  |  |  |
| port | int | The port number. |  |  |  |  |  |

## User Defined Functions (UDFs)

| Name | Description | Exportable | Output Type |
| --- | --- | --- | --- |
| buildUrl | Exported function to build a URL. | Yes | string |
| double | Internal function to double a value. |  | int |
