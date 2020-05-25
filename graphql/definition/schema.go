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

func (s *Schema) Execute(
	operation string,
	field string,
	node *ast.Field,
) (res interface{}, executedField *Field, err error) {
	f, ok := s.TypeMap[operation].Fields[field]
	if ok {
		parsedArgs := parseArgs(&f, node.Arguments)
		res := f.Resolve(f.TypeRef, parsedArgs, Infos{
			Field:     f,
			Requested: *node,
		})
		return res, &f, nil
	}
	return nil, nil, fmt.Errorf("Operation not found")
}

func (s *Schema) Send(res interface{}, infos Infos) map[string]interface{} {
	sendable := make(map[string]interface{})
	return sendable
}

func visitSelections(selectionSet *ast.SelectionSet) {
	for _, s := range selectionSet.Selections {
		ss := s.GetSelectionSet()
	}
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
