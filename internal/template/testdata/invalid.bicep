metadata name = 'test'
metadata description = 'This is a test template.'

@sys.description('This is a test parameter.')
param test_parameter string = 'test'

@sys.description('This is a test resource.')
resource test_resource 'Microsoft.Storage/storageAccounts@2023-01-01' = {
  name: 'test'
}

@sys.description('This is a test output')
output test_output  = 'test'
