metadata name = 'identity'
metadata description = 'Creates a User Assigned Identity.'

@description('Name of the identity.')
param identityName string

@description('Location of the identity.')
param location string = resourceGroup().location

resource identity 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: identityName
  location: location
}

@description('Resource ID of the identity.')
output resourceId string = identity.id

@description('Client ID of the identity.')
output clientId string = identity.properties.clientId

@description('Principal ID of the identity.')
output principalId string = identity.properties.principalId
