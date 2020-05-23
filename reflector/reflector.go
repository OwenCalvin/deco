package reflector

import (
	"reflect"
)

type Reflected struct {
	rType   reflect.Type
	fields  []reflect.StructField
	methods []reflect.Method
}

// ReflectTypes returns the types with their fields
func ReflectTypes(r ...interface{}) (reflections []Reflected) {
	for _, item := range r {
		r := Reflected{}
		r.rType = reflect.TypeOf(item).Elem()
		r.fields = ReflectFields(item)
		r.methods = ReflectMethods(item)

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
