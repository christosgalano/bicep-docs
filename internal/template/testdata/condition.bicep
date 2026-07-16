metadata name = 'condition-test'
metadata description = 'Template that exercises conditional resources and modules.'

@sys.description('Whether to deploy the resources.')
param deploy bool = true

@sys.description('The environment name.')
param environment string = 'dev'

@sys.description('This is a conditional module.')
module conditional_module './modules/test_module/main.bicep' = if (deploy) {
  name: 'test'
}

@sys.description('This is a conditional resource.')
resource conditional_resource 'Microsoft.Storage/storageAccounts@2023-01-01' = if (deploy && environment == 'prod') {
  name: 'test'
  location: 'westus'
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
}

@sys.description('This is a conditional loop resource.')
resource loop_resource 'Microsoft.Storage/storageAccounts@2023-01-01' = [for i in range(0, 2): if (deploy) {
  name: 'test${i}'
  location: 'westus'
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
}]

@sys.description('This is an unconditional resource.')
resource plain_resource 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: 'test2'
  location: 'westus'
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
}
