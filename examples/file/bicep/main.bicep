metadata name = 'storage account'
metadata description = 'Create a storage account.'

@sys.description('Name of the storage account.')
param name string

@sys.description('Location to deploy the storage account.')
param location string = resourceGroup().location

@sys.description('Name of the storage account\'s sku.')
param skuName string = 'Standard_LRS'

@sys.description('The kind of storage account.')
param kind string = 'StorageV2'

@sys.description('This is a test resource.')
resource st 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: name
  location: location
  sku: {
    name: skuName
  }
  kind: kind
}

@sys.description('Resource ID of the storage account.')
output resourceId string = st.id
