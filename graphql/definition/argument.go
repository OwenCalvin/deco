package definition

// ArgumentValue represents a GraphQL argument value
type Argument struct {
	Description  string
	Name         string
	Type         string
	DefaultValue interface{}
	Directives   []Directive
}

// Arguments represents a GraphQL argument list
type Arguments map[string]Argument
