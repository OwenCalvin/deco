package definition

import "reflect"

// AnyType represents a GraphQL type (object, input or interface)
type AnyType interface {
}

// Type represents a GraphQL type (object or interface)
type Type struct {
	Description string
	Name        string
	Fields      Fields
}

// OutpoutType represents a GraphQL type (object, input or interface)
type OutpoutType struct {
	Type
	AnyType
	Directives []Directive
}

// Object represents a GraphQL object type
type Object struct {
	OutpoutType
	ImplementsInterfaces []string
}

// Interface represents a GraphQL interface type
type Interface struct {
	OutpoutType
}

// Input represents a GraphQL input type
type Input struct {
	Type
	FieldsDefinition []Field
}

// FromReflection creates a type from reflection
func FromReflection(rType reflect.Type, fields []reflect.StructField, methods []reflect.Method) {
}
