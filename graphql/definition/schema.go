package definition

import (
	"dego/graphql/language/ast"
	"fmt"
	"reflect"
)

// Schema represents a graphql schema
type Schema struct {
	Directives       []Directive
	TypeMap          map[string]Type
	Types            []Type
	SubscriptionType Type
	QueryType        Type
	MutationType     Type
}

func (s *Schema) Execute(
	operation string,
	field string,
	node *ast.Field,
	fragments map[string]*ast.FragmentDefinition,
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
	return selectFields(res, &infos.Requested.SelectionSet.Selections)
}

func selectFields(obj interface{}, selections *[]ast.Selection) (res map[string]interface{}) {
	if selections == nil {
		return nil
	}

	res = map[string]interface{}{}
	value := reflect.ValueOf(obj)

	for i := 0; i < value.NumField(); i++ {
		var selected *ast.Field = nil
		rField := value.Type().Field(i)

		for _, s := range *selections {
			f := s.(*ast.Field)
			if f.Name.Value == rField.Name {
				selected = f
				break
			}
		}

		if selected == nil {
			continue
		}

		var v interface{}
		value := value.Field(i).Interface()

		switch rField.Type.Kind() {
		case reflect.Struct:
			v = selectFields(value, &selected.SelectionSet.Selections)
		default:
			v = value
		}

		if v != nil {
			res[rField.Name] = v
		}
	}

	return res
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
