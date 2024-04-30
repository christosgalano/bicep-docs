# test.bicep

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters
    storageAccountName:

    // Optional parameters
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Parameters

| Name | Type | Description | Default |
| --- | --- | --- | --- |
| storageAccountName | string | Storage account name restrictions:<br>- Storage account names must be between 3 and 24 characters in length and may contain numbers and lowercase letters only.<br>- Your storage account name must be unique within Azure. No two storage accounts can have the same name.<br> |  |
