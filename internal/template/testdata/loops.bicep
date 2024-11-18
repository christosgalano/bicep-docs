metadata name = 'loop_test'
metadata description = 'Test template with loop constructs'

@description('Array of deployment locations')
param locations array = [
  'eastus'
  'westus'
  'northeurope'
]

@description('Prefix for resource names')
param namePrefix string = 'storage'

@description('Array of storage account configurations')
var storageConfigs = [
  for i in range(0, length(locations)): {
    name: '${namePrefix}${i}'
    location: locations[i]
    sku: 'Standard_LRS'
    kind: 'StorageV2'
  }
]

@description('Array of generated storage account names')
var storageNames = [for config in storageConfigs: config.name]

@description('This is a test module.')
module test_module './modules/test_module/main.bicep' = {
  name: 'test'
}

@description('This is a storage account resource array.')
resource test_resource 'Microsoft.Storage/storageAccounts@2023-01-01' = [
  for config in storageConfigs: {
    name: config.name
    location: config.location
    sku: {
      name: config.sku
    }
    kind: config.kind
  }
]

@description('Array of created storage account resource IDs')
output resourceIds array = [for i in range(0, length(locations)): test_resource[i].id]

@description('Array of created storage account names')
output storageNames array = storageNames
