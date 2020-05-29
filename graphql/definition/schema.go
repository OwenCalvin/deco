package definition

import (
	"deco/graphql/language/ast"
	"deco/utils"
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
	var opExecuted *ast.OperationDefinition

	for _, def := range AST.Definitions {
		switch def.(type) {
		case *ast.OperationDefinition:
			opExecuted = def.(*ast.OperationDefinition)
			// TODO: VariableDefinition
			// TODO: Alias
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

	operationTypeName := utils.UpperFirstLetter(opExecuted.Operation)
	response := map[string]interface{}{}
	s.executeField(
		&executable,
		operationTypeName,
		response,
	)

	sendable := s.Send(
		response,
		executable,
		opExecuted.Name.Value,
	)

	return sendable, nil
}

func (s *Schema) executeField(
	executed *ast.Field,
	parent string,
	r interface{},
) {
	for _, selection := range executed.SelectionSet.Selections {
		var value interface{}
		field := selection.(*ast.Field)

		f, ok := s.TypeMap[parent].Fields[field.Name.Value]
		if ok {
			parsedArgs := parseArgs(&f, field.Arguments)
			infos := Infos{
				Field: field,
			}

			if f.DefaultValue != nil {
				value = reflect.ValueOf(f.DefaultValue).Elem().Interface()
			} else {
				value = f.Resolve(f.TypeRef, parsedArgs, infos)
			}

			mapped := utils.StructToMap(value)

			if it, ok := r.(map[string]interface{}); ok {
				s.executeField(field, f.Type, it)
				it[f.Name] = mapped
			} else {
				r = mapped
			}
		}
	}
}

func (s *Schema) Send(
	res map[string]interface{},
	field ast.Field,
	responseName string,
) map[string]interface{} {
	if responseName == "" {
		responseName = field.Name.Value
	}

	return map[string]interface{}{
		responseName: selectFields(
			res,
			field.SelectionSet.Selections,
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

// TODO: Find type definition for names
func selectFields(
	obj map[string]interface{},
	selections []ast.Selection,
) interface{} {
	res := map[string]interface{}{}

	for _, selection := range selections {
		astField := selection.(*ast.Field)
		value := obj[astField.Name.Value]
		ss := astField.SelectionSet
		if ss != nil {
			if mapped, ok := value.(map[string]interface{}); ok {
				value = selectFields(mapped, ss.Selections)
			}
		}
		res[astField.Name.Value] = value
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
