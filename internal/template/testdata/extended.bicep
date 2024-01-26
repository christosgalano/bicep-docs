metadata name = 'test'
metadata description = 'This is a test template.'

@sys.description('This is a test parameter.')
param test_parameter string = 'test'

var test_variable = '${test_parameter}'

@minValue(1)
@sys.description('This is a user defined type.')
type positiveInt = int

@sys.description('This is a user defined type (alias).')
type pint = positiveInt

@sys.description('This is a user defined function.')
func buildUrl(https bool, hostname string, path string) string =>
  '${https ? 'https' : 'http'}://${hostname}${empty(path) ? '' : '/${path}'}'

@sys.description('This is a user defined function with uddts.')
func double(input pint) positiveInt => input * 2

@sys.description('This is a test module.')
module test_module './modules/test_module/main.bicep' = {
  name: 'test'
}

@sys.description('This is a test resource.')
resource test_resource 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: 'test'
  location: 'westus'
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
}

@sys.description('This is a test output.')
output test_output positiveInt = 1
