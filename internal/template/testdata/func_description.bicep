@description('This is a user defined function.')
func getDefaultTags(env string) object => {
  Environment: env
}

var undescribedVariable = getDefaultTags('dev')

@description('This is a described variable.')
var describedVariable = 'value'
