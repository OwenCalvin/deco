package definition

// Schema represents a graphql schema
type Schema struct {
	Directives []Directive
	TypeMap    map[string]Type
}

func (s *Schema) Execute(operation string, field string, variables interface{}) interface{} {
	f, ok := s.TypeMap[operation].FieldsDefinition[field]
	if ok {
		res := f.Resolve(f.TypeRef, variables)
		return res
	}
	return struct{}{}
}
