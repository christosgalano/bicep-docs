metadata name = 'decorators-test'
metadata description = 'Template that exercises resource decorators.'

@sys.description('This is a resource with retry behavior.')
@retryOn(['ServerError', 'Conflict'], 3)
resource retry_resource 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: 'test'
  location: 'westus'
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
}

@onlyIfNotExists()
@sys.description('This is a resource that is only created if it does not exist.')
resource idempotent_resource 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: 'test2'
  location: 'westus'
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
}

@sys.description('This is a plain resource.')
resource plain_resource 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: 'test3'
  location: 'westus'
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
}
