package definition

import (
	"reflect"
)

// FieldResolver represents implements the resolve func for a field
type FieldResolver interface {
	Resolve() interface{}
}

// Field represents a GraphQL field value
type Field struct {
	Resolve       func(ref interface{}, args interface{}, infos Infos) interface{}
	Name          string
	Args          Arguments
	Type          string
	TypeRef       interface{}
	Directives    []Directive
	ArgStructType reflect.Type
}

// Fields represents a GraphQL list of fields
type Fields map[string]Field
