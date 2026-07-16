metadata name = 'export-variables-test'
metadata description = 'Template that exercises exported variables.'

@export()
@sys.description('A shared prefix used across resources.')
var namePrefix = 'app'

@export()
var defaultTags = {
  environment: 'dev'
}

@sys.description('An internal variable (not exported).')
var internalValue = '${namePrefix}-internal'

output internal string = internalValue
