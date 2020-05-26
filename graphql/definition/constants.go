package definition

const (
	ID       = "ID"
	BOOL     = "Boolean"
	STRING   = "String"
	NUMBER   = "Float"
	NOTNULL  = "!"
	ARRAY    = "[%v]"
	QUERY    = "Query"
	MUTATION = "Mutation"
)

var specialFields = []string{
	"__schema",
	"__typename",
}
