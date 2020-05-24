package reflector

import (
	"reflect"
)

type Reflected struct {
	Type     reflect.Type
	Original interface{}
	Fields   []reflect.StructField
	Methods  []reflect.Method
}

// ReflectTypes returns the types with their fields
func ReflectTypes(r ...interface{}) (reflections []Reflected) {
	for _, item := range r {
		r := Reflected{}
		r.Original = item
		r.Type = reflect.TypeOf(item).Elem()
		r.Fields = ReflectFields(item)
		r.Methods = ReflectMethods(item)

		reflections = append(reflections, r)
	}

	return reflections
}

// ReflectMethods returns the fields of a struct
func ReflectMethods(r interface{}) (methods []reflect.Method) {
	e := reflect.TypeOf(r)

	for i := 0; i < e.NumMethod(); i++ {
		methods = append(methods, e.Method(i))
	}

	return methods
}

// ReflectFields returns the fields of a struct
func ReflectFields(r interface{}) (fields []reflect.StructField) {
	e := reflect.ValueOf(r).Elem()

	for i := 0; i < e.NumField(); i++ {
		fields = append(fields, e.Type().Field(i))
	}

	return fields
}
