# decorators-showcase

## Description

Comprehensive example showcasing all Bicep parameter decorators and exportable features

## Usage

Here is a basic example of how to use this Bicep module:

```bicep
module reference_name 'path_to_module | container_registry_reference' = {
  name: 'deployment_name'
  params: {
    // Required parameters
    adminPassword:
    adminUsername:
    customResourceName:
    environment:
    networkSettings:
    resourcePrefix:
    storageSettings:

    // Optional parameters
    additionalSubnets: []
    diskSizeGB: 128
    enableMonitoring: true
    instanceCount: 2
    location: 'eastus'
    tags: {
      Environment: '[parameters('environment')]'
      ManagedBy: 'bicep-docs'
    }
    vmSize: 'Standard_B2s'
  }
}
```

> Note: In the default values, strings enclosed in square brackets (e.g. '[resourceGroup().location]' or '[__bicep.function_name(args...)']) represent function calls or references.

## Resources

| Symbolic Name | Type | Description |
| --- | --- | --- |
| applicationInsights | [Microsoft.Insights/components](https://learn.microsoft.com/en-us/azure/templates/microsoft.insights/components) |  |
| storageAccount | [Microsoft.Storage/storageAccounts](https://learn.microsoft.com/en-us/azure/templates/microsoft.storage/storageaccounts) |  |
| virtualNetwork | [Microsoft.Network/virtualNetworks](https://learn.microsoft.com/en-us/azure/templates/microsoft.network/virtualnetworks) |  |

## Parameters

| Name | Status | Type | Description | Default |
| --- | --- | --- | --- | --- |
| additionalSubnets | Optional | array | Optional array of additional subnets | [] |
| adminPassword | Required | securestring | Administrator password |  |
| adminUsername | Required | string | Administrator username |  |
| customResourceName | Required | resourceName (uddt) | Custom resource name using exported type |  |
| diskSizeGB | Optional | int | Data disk size in GB | 128 |
| enableMonitoring | Optional | bool | Enable monitoring and diagnostics | true |
| environment | Required | string | Application environment (dev, test, prod) |  |
| instanceCount | Optional | int | Number of instances to deploy | 2 |
| location | Optional | string | Azure region for resource deployment | "eastus" |
| networkSettings | Required | networkConfig (uddt) | Network configuration settings |  |
| resourcePrefix | Required | string | Resource name prefix |  |
| storageSettings | Required | storageConfig (uddt) | Storage account configuration |  |
| tags | Optional | object | Resource tags as key-value pairs | {"Environment": "[parameters('environment')]", "ManagedBy": "bicep-docs"} |
| vmSize | Optional | string | Virtual machine size | "Standard_B2s" |

## User Defined Data Types (UDDTs)

| Name | Type | Description | Properties |
| --- | --- | --- | --- |
| networkConfig | object | Non-exportable custom type for network settings | [View Properties](#networkconfig) |
| resourceName | string | Exportable custom string type for resource names |  |
| storageConfig | object | Exportable custom type for storage account configuration | [View Properties](#storageconfig) |
| tagValue | string | Non-exportable custom string type for tags |  |

### networkConfig

| Name | Type | Description |
| --- | --- | --- |
| enableDdosProtection | bool | Enable DDoS protection |
| vnetName | string | Virtual network name |

### storageConfig

| Name | Type | Description |
| --- | --- | --- |
| enableHierarchicalNamespace | bool | Enable hierarchical namespace |
| name | string | Storage account name |
| sku | string | Storage account SKU |

## User Defined Functions (UDFs)

| Name | Description | Output Type |
| --- | --- | --- |
| calculateStorageSize | Non-exportable function to calculate storage size | int |
| generateResourceName | Exportable function to generate unique resource names | string |
| getDefaultTags | Non-exportable function to get default tags | object |
| isValidResourceName | Exportable function to validate resource naming convention | bool |

## Variables

| Name | Description |
| --- | --- |
| fullResourceName | Non-exportable function to get default tags |
| mergedTags |  |
| totalStorageNeeded |  |

## Outputs

| Name | Type | Description |
| --- | --- | --- |
| adminUser | string | Admin username provided |
| appliedTags | object | Resource tags applied |
| deploymentLocation | string | Resource deployment location |
| generatedResourceName | string | Generated unique resource name |
| monitoringEnabled | bool | Monitoring enabled status |
| nameValidation | bool | Function validation result |
| selectedVmSize | string | VM size selected |
| storageEndpoint | string | Storage account primary endpoint |
| totalStorageGB | int | Total storage capacity calculated |
| vnetId | string | Virtual network resource ID |
