targetScope = 'subscription'

metadata name = 'development'
metadata description = 'Deploy to the development environment.'

@sys.description('Location to deploy the resources.')
param location string = 'westeurope'

@sys.description('Name of the resource group to deploy the resources.')
param rgName string = 'rgtest123'

@sys.description('Resource Group to deploy the resources.')
resource rg 'Microsoft.Resources/resourceGroups@2023-07-01' = {
  name: rgName
  location: location
}

@sys.description('Create a User Assigned Identity.')
module identity '../../modules/identity/main.bicep' = {
  name: 'identity-deployment'
  scope: resourceGroup(rgName)
  params: {
    identityName: 'identity'
    location: location
  }
  dependsOn: [
    rg
  ]
}

@description('Client ID of the User Assigned Identity.')
output identityClientId string = identity.outputs.clientId
