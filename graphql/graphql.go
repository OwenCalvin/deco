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
		Name:        definition.QUERY,
		Description: "",
		Fields:      make(definition.Fields),
	}

	schema.TypeMap[definition.MUTATION] = definition.Type{
		Name:        definition.MUTATION,
		Description: "",
		Fields:      make(definition.Fields),
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
			queryName := method.Name

			switch {
			case strings.HasSuffix(method.Name, definition.QUERY):
				isQueryOrMutation = true
				queryType = definition.QUERY
				queryName = strings.ReplaceAll(queryName, definition.QUERY, "")
			case strings.HasSuffix(method.Name, definition.MUTATION):
				isQueryOrMutation = true
				queryType = definition.MUTATION
				queryName = strings.ReplaceAll(queryName, definition.MUTATION, "")
			}

			if !isQueryOrMutation {
				continue
			}

			args := definition.Arguments{}

			firstArg := method.Type.In(1)

			for i := 0; i < firstArg.NumField(); i++ {
				arg := firstArg.Field(i)
				argValue := definition.Argument{
					Name:         arg.Name,
					Description:  "",
					Type:         arg.Type.Name(),
					DefaultValue: nil,
				}
				args[argValue.Name] = argValue
			}

			returnType := method.Type.Out(0)
			fieldValue := definition.Field{
				Name:          queryName,
				Type:          returnType.Name(),
				TypeRef:       item.Original,
				Args:          args,
				ArgStructType: firstArg,
			}

			fieldValue.Resolve = func(ref interface{}, args interface{}, infos definition.Infos) interface{} {
				res := method.Func.Call([]reflect.Value{
					reflect.ValueOf(ref),
					reflect.ValueOf(args),
					reflect.ValueOf(infos),
				})

				return schema.Send(res[0].Interface(), infos)
			}

			schema.TypeMap[queryType].Fields[queryName] = fieldValue
		}

		t := definition.Type{
			Name:        typeName,
			Description: "",
			Fields:      fields,
		}

		schema.TypeMap[typeName] = t
	}

	return schema
}
