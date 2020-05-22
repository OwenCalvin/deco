package reflector

import (
	"reflect"
)

// ReflectAll returns the types with their fields
func ReflectAll(r ...interface{}) (types [][]interface{}) {
	for _, item := range r {
		name := reflect.TypeOf(item).Elem().Name()
		typeDef := []interface{}{name, Reflect(item)}
		types = append(types, typeDef)
	}

	return types
}

// Reflect returns the fields of a struct
func Reflect(r interface{}) (fields [][]string) {
	e := reflect.ValueOf(r).Elem()

	for i := 0; i < e.NumField(); i++ {
		field := e.Type().Field(i)
		infos := []string{field.Name, field.Type.Name()}
		fields = append(fields, infos)
	}

	return fields
}
