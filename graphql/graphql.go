package graphql

import (
	"dego/graphql/definition"
	"dego/reflector"
	"reflect"
	"strings"
)

func LoadTypes(r ...interface{}) (schema definition.Schema) {
	schema = definition.Schema{
		TypeMap: make(map[string]definition.Type),
	}
	reflections := reflector.ReflectTypes(r...)

	schema.TypeMap[definition.QUERY] = definition.Type{
		Name:             definition.QUERY,
		Description:      "",
		FieldsDefinition: make(definition.Fields),
	}

	schema.TypeMap[definition.MUTATION] = definition.Type{
		Name:             definition.MUTATION,
		Description:      "",
		FieldsDefinition: make(definition.Fields),
	}

	for _, item := range reflections {
		typeName := item.Type.Name()
		fields := make(definition.Fields)

		for _, field := range item.Fields {
			fieldKey := field.Name
			fieldValue := definition.Field{
				Type: field.Type.Name(),
			}
			fields[fieldKey] = fieldValue
		}

		for _, method := range item.Methods {
			isQueryOrMutation := false
			queryType := ""

			switch {
			case strings.HasSuffix(method.Name, definition.QUERY):
				isQueryOrMutation = true
				queryType = definition.QUERY
			case strings.HasSuffix(method.Name, definition.MUTATION):
				isQueryOrMutation = true
				queryType = definition.MUTATION
			}

			if !isQueryOrMutation {
				continue
			}

			args := definition.Arguments{}

			firstArg := method.Type.In(1)

			for i := 0; i < firstArg.NumField(); i++ {
				arg := firstArg.Field(i)
				argKey := definition.ArgumentKey{
					Name:        arg.Name,
					Description: "",
				}
				argValue := definition.ArgumentValue{
					Type:         arg.Type.Name(),
					DefaultValue: nil,
				}
				args[argKey] = argValue
			}

			returnType := method.Type.Out(0)
			fieldValue := definition.Field{
				Name:                method.Name,
				Type:                returnType.Name(),
				TypeRef:             item.Original,
				ArgumentsDefinition: args,
			}

			fieldValue.Resolve = func(v ...interface{}) interface{} {
				res := method.Func.Call([]reflect.Value{
					reflect.ValueOf(v[0]),
					reflect.ValueOf(v[1]),
				})
				return res[0].Interface()
			}

			schema.TypeMap[queryType].FieldsDefinition[method.Name] = fieldValue
		}

		t := definition.Type{
			Name:             typeName,
			Description:      "",
			FieldsDefinition: fields,
		}

		schema.TypeMap[typeName] = t
	}

	return schema
}
