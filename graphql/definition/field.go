package definition

import (
	"reflect"
)

// Field represents a GraphQL field value
type Field struct {
	Resolve           func(ref interface{}, args interface{}, infos Infos) interface{}
	OriginalName      string
	Name              string
	ParentType        *Type
	ParentTypeRef     interface{}
	Args              Arguments
	DefaultValue      interface{}
	Type              string
	TypeRef           interface{}
	Nullable          bool
	ListOf            string
	Directives        []Directive
	ArgStructType     reflect.Type
	isDeprecated      bool
	deprecationReason string
}

// Fields represents a GraphQL list of fields
type Fields map[string]Field
