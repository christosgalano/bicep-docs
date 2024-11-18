metadata name = 'var_optimization'
metadata description = 'Test template with variable optimization'

@description('Name of the App Service Plan to host Web-App on')
param servicePlanName string

@description('Get resilience options from Service Plan')
var isZoneRedundant = servicePlan.properties.zoneRedundant
var reserved = servicePlan.properties.reserved

@description('Get App Service Plan Object')
resource servicePlan 'Microsoft.Web/serverfarms@2024-04-01' existing = {
  name: servicePlanName
}

output isZoneRedundant bool = isZoneRedundant
output reserved bool = reserved
