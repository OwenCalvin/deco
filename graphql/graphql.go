package graphql

import (
	"dego/graphql/definition"
	"dego/graphql/queries"
	"dego/reflector"
	"fmt"
	"reflect"
	"strings"
)

func LoadTypes(r ...interface{}) (schema definition.Schema) {
	types := r
	types = append(types, &queries.IntrospectionQuery{})
	reflections := reflector.ReflectTypes(types...)

	schema = definition.Schema{
		TypeMap: make(map[string]definition.Type),
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
			queryName := method.Name

			args := definition.Arguments{}

			if method.Type.NumOut() <= 0 {
				panic(fmt.Errorf("You must specify a return type for %v.%v", item.Type.Elem().Name, queryName))
			}

			returnType := method.Type.Out(0)
			fieldValue := definition.Field{
				Name:    queryName,
				Type:    returnType.Name(),
				TypeRef: item.Original,
				Args:    args,
			}

			if method.Type.NumIn() > 1 {
				firstArg := method.Type.In(1)

				for i := 0; i < firstArg.NumField(); i++ {
					arg := firstArg.Field(i)
					argValue := definition.Argument{
						Name:         arg.Name,
						Description:  "",
						Type:         arg.Type.Name(),
						DefaultValue: nil,
					}
					args[arg.Name] = argValue
				}

				fieldValue.ArgStructType = firstArg
			}

			fieldValue.Resolve = func(ref interface{}, args interface{}, infos definition.Infos) interface{} {
				res := method.Func.Call([]reflect.Value{
					reflect.ValueOf(ref),
					reflect.ValueOf(args),
					reflect.ValueOf(infos),
				})

				return schema.Send(res[0].Interface(), infos)
			}

			switch {
			case strings.HasSuffix(queryName, definition.QUERY):
				addFieldToTypeMap(&schema, definition.QUERY, queryName, fieldValue)
			case strings.HasSuffix(queryName, definition.MUTATION):
				addFieldToTypeMap(&schema, definition.MUTATION, queryName, fieldValue)
			default:
				fields[queryName] = fieldValue
			}
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

func addFieldToTypeMap(schema *definition.Schema, key string, name string, field definition.Field) {
	c := strings.LastIndex(name, key)
	finalQueryName := name[:c]
	schema.TypeMap[key].Fields[finalQueryName] = field
}
