package definition

// FieldResolver represents implements the resolve func for a field
type FieldResolver interface {
	Resolve() interface{}
}

// Field represents a GraphQL field value
type Field struct {
	Resolve             func(v ...interface{}) interface{}
	Name                string
	ArgumentsDefinition Arguments
	Type                string
	Directives          []Directive
	TypeRef             interface{}
}

// Fields represents a GraphQL list of fields
type Fields map[string]Field
