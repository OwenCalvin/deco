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

	response := map[string]interface{}{}
	s.ExecuteField(
		&executable,
		utils.UpperFirstLetter(opExecuted.Operation),
		response,
	)

	return nil, nil
}

// TODO: Recursive
func (s *Schema) ExecuteField(
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
				Field:     f,
				Requested: *field,
			}

			if f.DefaultValue != nil {
				value = reflect.ValueOf(f.DefaultValue).Elem().Interface()
			} else {
				value = f.Resolve(f.TypeRef, parsedArgs, infos)
			}

			mapped := utils.StructToMap(value)

			if it, ok := r.(map[string]interface{}); ok {
				s.ExecuteField(field, f.Type, it)
				it[f.Name] = mapped
			} else {
				r = mapped
			}
		}
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
