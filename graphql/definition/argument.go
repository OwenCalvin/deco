package definition

// ArgumentKey represents a GraphQL argument key
type ArgumentKey struct {
	Description string
	Name        string
}

// ArgumentValue represents a GraphQL argument value
type ArgumentValue struct {
	Type         string
	DefaultValue interface{}
	Directives   []Directive
}

// Arguments represents a GraphQL argument list
type Arguments map[ArgumentKey]ArgumentValue
