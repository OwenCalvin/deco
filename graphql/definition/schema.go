package definition

import (
	"dego/graphql/language/ast"
	"fmt"
	"reflect"
)

// Schema represents a graphql schema
type Schema struct {
	Directives []Directive
	TypeMap    map[string]Type
}

func (s *Schema) Execute(operation string, field string, arguments []*ast.Argument) (res interface{}, err error) {
	f, ok := s.TypeMap[operation].Fields[field]
	if ok {
		parsedArgs := parseArgs(&f, arguments)
		res := f.Resolve(f.TypeRef, parsedArgs)
		return res, nil
	}
	return nil, fmt.Errorf("Operation not found")
}

func parseArgs(field *Field, args []*ast.Argument) interface{} {
	structValue := reflect.New(field.ArgStructType)

	for _, item := range args {
		structValue.Elem().FieldByName(item.Name.Value).SetString(
			item.Value.GetValue().(string),
		)
	}

	return structValue.Elem().Interface()
}
