package definition

import (
	"deco/graphql/language/ast"
	"deco/utils"
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

func (s *Schema) Execute(AST *ast.Document) (res interface{}, err error) {
	fragments := make(map[string]*ast.FragmentDefinition)
	var executed *ast.Field

	for _, def := range AST.Definitions {
		switch def.(type) {
		case *ast.OperationDefinition:
			opExecuted := def.(*ast.OperationDefinition)
			executed = ast.NewField(&ast.Field{
				Name:         opExecuted.Name,
				SelectionSet: opExecuted.SelectionSet,
			})
		case *ast.FragmentDefinition:
			f := def.(*ast.FragmentDefinition)
			fragments[f.Name.Value] = f
		}
	}

	executable := *executed

	parseFragments(
		&executable,
		fragments,
	)

	s.executeFields(&executable, AST, fragments)

	return nil, nil
}

func (s *Schema) ExecuteField(
	operation string,
	field string,
	node *ast.Field,
	fragments map[string]*ast.FragmentDefinition,
	responseName string,
) (res interface{}, executedField *Field, err error) {
	f, ok := s.TypeMap[operation].Fields[field]
	if ok {
		parsedArgs := parseArgs(&f, node.Arguments)
		var res interface{}
		infos := Infos{
			Field:     f,
			Requested: *node,
		}

		if f.DefaultValue != nil {
			res = reflect.ValueOf(f.DefaultValue).Elem().Interface()
		} else {
			res = f.Resolve(f.TypeRef, parsedArgs, infos)
		}

		toSend := s.Send(res, infos, responseName)
		return toSend, &f, nil
	}
	return nil, nil, fmt.Errorf("Operation not found")
}

// TODO: Recursive
func (s *Schema) executeFields(
	executed *ast.Field,
	AST *ast.Document,
	fragments map[string]*ast.FragmentDefinition,
) {
	for i, selection := range executed.SelectionSet.Selections {
		field := selection.(*ast.Field)
		operation := AST.Definitions[i].(*ast.OperationDefinition).Operation
		operation = utils.UpperFirstLetter(operation)

		// Choose correct operation
		// Only if selected
		s.ExecuteField(
			operation,
			field.Name.Value,
			field,
			fragments,
			executed.Name.Value,
		)
	}
}

func (s *Schema) Send(
	res interface{},
	infos Infos,
	responseName string,
) map[string]map[string]interface{} {
	if responseName == "" {
		responseName = infos.Field.Name
	}

	return map[string]map[string]interface{}{
		responseName: selectFields(
			res,
			&infos.Requested.SelectionSet.Selections,
			&infos,
		),
	}
}

func parseFragments(
	fieldRoot *ast.Field,
	fragments map[string]*ast.FragmentDefinition,
) {
	selections := (*fieldRoot).GetSelectionSet().Selections
	for i, s := range selections {
		switch s.(type) {
		case *ast.FragmentSpread:
			f := *s.(*ast.FragmentSpread)
			fragmentValue := fragments[f.Name.Value]
			newField := ast.NewField(&ast.Field{
				Name:         fragmentValue.Name,
				SelectionSet: fragmentValue.SelectionSet,
				Directives:   fragmentValue.Directives,
			})

			parseFragments(
				newField,
				fragments,
			)

			selections[i] = newField
		case *ast.Field:
			f := *s.(*ast.Field)
			if f.SelectionSet != nil {
				parseFragments(
					&f,
					fragments,
				)
			}
		}
	}
}

func selectFields(
	obj interface{},
	selections *[]ast.Selection,
	infos *Infos,
) (res map[string]interface{}) {
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
			name := utils.LowerFirstLetter(rField.Name)
			tagName := rField.Tag.Get("name")
			if tagName != "" {
				name = tagName
			}

			if f.Name.Value == name {
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
			v = selectFields(value, selections, infos)
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
	if field.ArgStructType != nil {
		structValue := reflect.New(field.ArgStructType)

		for _, item := range args {
			structValue.Elem().FieldByName(item.Name.Value).SetString(
				item.Value.GetValue().(string),
			)
		}

		return structValue.Elem().Interface()
	}
	return []interface{}{}
}
