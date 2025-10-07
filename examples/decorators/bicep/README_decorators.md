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

| Name | Status | Type | Description | Default | Allowed Values | Min Length | Max Length | Min Value | Max Value |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| additionalSubnets | Optional | array | Optional array of additional subnets | [] |  |  |  |  |  |
| adminPassword | Required | securestring | Administrator password |  |  | 12 | 128 |  |  |
| adminUsername | Required | string | Administrator username |  |  | 3 | 20 |  |  |
| customResourceName | Required | resourceName (uddt) | Custom resource name using exported type |  |  |  |  |  |  |
| diskSizeGB | Optional | int | Data disk size in GB | 128 |  |  |  | 32 | 1024 |
| enableMonitoring | Optional | bool | Enable monitoring and diagnostics | true |  |  |  |  |  |
| environment | Required | string | Application environment (dev, test, prod) |  | `dev`, `test`, `prod` |  |  |  |  |
| instanceCount | Optional | int | Number of instances to deploy | 2 |  |  |  | 1 | 10 |
| location | Optional | string | Azure region for resource deployment | "eastus" | `eastus`, `westus`, `centralus`, `westeurope`, `northeurope` |  |  |  |  |
| networkSettings | Required | networkConfig (uddt) | Network configuration settings |  |  |  |  |  |  |
| resourcePrefix | Required | string | Resource name prefix |  |  | 2 | 10 |  |  |
| storageSettings | Required | storageConfig (uddt) | Storage account configuration |  |  |  |  |  |  |
| tags | Optional | object | Resource tags as key-value pairs | {"Environment": "[parameters('environment')]", "ManagedBy": "bicep-docs"} |  |  |  |  |  |
| vmSize | Optional | string | Virtual machine size | "Standard_B2s" | `Standard_B2s`, `Standard_D2s_v3`, `Standard_D4s_v3` |  |  |  |  |

## User Defined Data Types (UDDTs)

| Name | Type | Description | Exportable | Properties | Min Length | Max Length | Min Value | Max Value |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| networkConfig | object | Non-exportable custom type for network settings | False | [View Properties](#networkconfig) |  |  |  |  |
| resourceName | string | Exportable custom string type for resource names | True |  | 5 | 50 |  |  |
| storageConfig | object | Exportable custom type for storage account configuration | True | [View Properties](#storageconfig) |  |  |  |  |
| tagValue | string | Non-exportable custom string type for tags | False |  | 1 | 10 |  |  |

### networkConfig

| Name | Type | Description | Allowed Values | Min Length | Max Length | Min Value | Max Value |
| --- | --- | --- | --- | --- | --- | --- | --- |
| enableDdosProtection | bool | Enable DDoS protection |  |  |  |  |  |
| vnetName | string | Virtual network name |  | 1 | 80 |  |  |

### storageConfig

| Name | Type | Description | Allowed Values | Min Length | Max Length | Min Value | Max Value |
| --- | --- | --- | --- | --- | --- | --- | --- |
| enableHierarchicalNamespace | bool | Enable hierarchical namespace |  |  |  |  |  |
| name | string | Storage account name |  | 3 | 24 |  |  |
| sku | string | Storage account SKU | `Premium_LRS`, `Standard_GRS`, `Standard_LRS` |  |  |  |  |

## User Defined Functions (UDFs)

| Name | Description | Exportable | Output Type |
| --- | --- | --- | --- |
| calculateStorageSize | Non-exportable function to calculate storage size | False | int |
| generateResourceName | Exportable function to generate unique resource names | True | string |
| getDefaultTags | Non-exportable function to get default tags | False | object |
| isValidResourceName | Exportable function to validate resource naming convention | True | bool |

## Variables

| Name | Description |
| --- | --- |
| fullResourceName | Non-exportable function to get default tags |
| mergedTags |  |
| totalStorageNeeded |  |

## Outputs

| Name | Type | Description | Min Length | Max Length | Min Value | Max Value |
| --- | --- | --- | --- | --- | --- | --- |
| adminUser | string | Admin username provided |  |  |  |  |
| appliedTags | object | Resource tags applied |  |  |  |  |
| deploymentLocation | string | Resource deployment location |  |  |  |  |
| generatedResourceName | string | Generated unique resource name | 10 | 100 |  |  |
| monitoringEnabled | bool | Monitoring enabled status |  |  |  |  |
| nameValidation | bool | Function validation result |  |  |  |  |
| selectedVmSize | string | VM size selected |  |  |  |  |
| storageEndpoint | string | Storage account primary endpoint | 10 | 200 |  |  |
| totalStorageGB | int | Total storage capacity calculated |  |  | 32 | 10240 |
| vnetId | string | Virtual network resource ID |  |  |  |  |
