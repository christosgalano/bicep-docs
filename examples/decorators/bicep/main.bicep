metadata name = 'decorators-showcase'
metadata description = 'Comprehensive example showcasing all Bicep parameter decorators and exportable features'

// ============================================================================
// CUSTOM TYPES
// ============================================================================

@export()
@description('Exportable custom type for storage account configuration')
type storageConfig = {
  @minLength(3)
  @maxLength(24)
  @description('Storage account name')
  name: string

  @description('Storage account SKU')
  sku: ('Standard_LRS' | 'Standard_GRS' | 'Premium_LRS')

  @description('Enable hierarchical namespace')
  enableHierarchicalNamespace: bool
}

@description('Non-exportable custom type for network settings')
type networkConfig = {
  @minLength(1)
  @maxLength(80)
  @description('Virtual network name')
  vnetName: string

  @description('Enable DDoS protection')
  enableDdosProtection: bool
}

@export()
@minLength(5)
@maxLength(50)
@description('Exportable custom string type for resource names')
type resourceName = string

@minLength(1)
@maxLength(10)
@description('Non-exportable custom string type for tags')
type tagValue = string

// ============================================================================
// PARAMETERS
// ============================================================================

@description('Application environment (dev, test, prod)')
@allowed(['dev', 'test', 'prod'])
param environment string

@description('Azure region for resource deployment')
@allowed(['eastus', 'westus', 'centralus', 'westeurope', 'northeurope'])
param location string = 'eastus'

@description('Resource name prefix')
@minLength(2)
@maxLength(10)
param resourcePrefix string

@description('Administrator username')
@minLength(3)
@maxLength(20)
param adminUsername string

@description('Administrator password')
@secure()
@minLength(12)
@maxLength(128)
param adminPassword string

@description('Virtual machine size')
@allowed(['Standard_B2s', 'Standard_D2s_v3', 'Standard_D4s_v3'])
param vmSize string = 'Standard_B2s'

@description('Number of instances to deploy')
@minValue(1)
@maxValue(10)
param instanceCount int = 2

@description('Data disk size in GB')
@minValue(32)
@maxValue(1024)
param diskSizeGB int = 128

@description('Storage account configuration')
param storageSettings storageConfig

@description('Network configuration settings')
param networkSettings networkConfig

@description('Resource tags as key-value pairs')
param tags object = {
  Environment: environment
  ManagedBy: 'bicep-docs'
}

@description('Optional array of additional subnets')
param additionalSubnets array = []

@description('Enable monitoring and diagnostics')
param enableMonitoring bool = true

@description('Custom resource name using exported type')
param customResourceName resourceName

// ============================================================================
// USER-DEFINED FUNCTIONS
// ============================================================================

@export()
@description('Exportable function to generate unique resource names')
func generateResourceName(prefix string, suffix string, environment string) string =>
  '${prefix}-${suffix}-${environment}-${uniqueString(resourceGroup().id)}'

@description('Non-exportable function to calculate storage size')
func calculateStorageSize(instanceCount int, diskSize int) int => instanceCount * diskSize

@export()
@description('Exportable function to validate resource naming convention')
func isValidResourceName(name string) bool => length(name) >= 3 && length(name) <= 50 && !contains(name, ' ')

@description('Non-exportable function to get default tags')
func getDefaultTags(env string) object => {
  Environment: env
  CreatedBy: 'Bicep'
}

// ============================================================================
// VARIABLES
// ============================================================================

var fullResourceName = generateResourceName(resourcePrefix, 'vm', environment)
var totalStorageNeeded = calculateStorageSize(instanceCount, diskSizeGB)
var mergedTags = union(getDefaultTags(environment), tags)

// ============================================================================
// RESOURCES
// ============================================================================

resource virtualNetwork 'Microsoft.Network/virtualNetworks@2024-01-01' = {
  name: networkSettings.vnetName
  location: location
  tags: mergedTags
  properties: {
    addressSpace: {
      addressPrefixes: ['10.0.0.0/16']
    }
    enableDdosProtection: networkSettings.enableDdosProtection
    subnets: concat(
      [
        {
          name: 'default'
          properties: {
            addressPrefix: '10.0.1.0/24'
          }
        }
      ],
      additionalSubnets
    )
  }
}

resource storageAccount 'Microsoft.Storage/storageAccounts@2024-01-01' = {
  name: storageSettings.name
  location: location
  tags: mergedTags
  sku: {
    name: storageSettings.sku
  }
  kind: 'StorageV2'
  properties: {
    isHnsEnabled: storageSettings.enableHierarchicalNamespace
    encryption: {
      services: {
        blob: {
          enabled: true
        }
      }
      keySource: 'Microsoft.Storage'
    }
  }
}

// Conditional monitoring resource based on enableMonitoring parameter
resource applicationInsights 'Microsoft.Insights/components@2020-02-02' = if (enableMonitoring) {
  name: '${customResourceName}-insights'
  location: location
  tags: mergedTags
  kind: 'web'
  properties: {
    Application_Type: 'web'
  }
}

// ============================================================================
// OUTPUTS
// ============================================================================

@description('Generated unique resource name')
@minLength(10)
@maxLength(100)
output generatedResourceName string = fullResourceName

@description('Virtual network resource ID')
output vnetId string = virtualNetwork.id

@description('Storage account primary endpoint')
@minLength(10)
@maxLength(200)
output storageEndpoint string = storageAccount.properties.primaryEndpoints.blob

@description('Total storage capacity calculated')
@minValue(32)
@maxValue(10240)
output totalStorageGB int = totalStorageNeeded

@description('Resource deployment location')
output deploymentLocation string = location

@description('Resource tags applied')
output appliedTags object = mergedTags

@description('Function validation result')
output nameValidation bool = isValidResourceName(fullResourceName)

@description('VM size selected')
output selectedVmSize string = vmSize

@description('Admin username provided')
output adminUser string = adminUsername

@description('Monitoring enabled status')
output monitoringEnabled bool = enableMonitoring
