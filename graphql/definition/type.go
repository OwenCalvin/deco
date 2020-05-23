package definition

import "reflect"

// GqlType represents a GraphQL type
type GqlType struct {
	Name   string
	Fields []Field
}

// GqlObject represents a GraphQL object type
type GqlObject struct {
	GqlType
}

// GqlInterface represents a GraphQL interface type
type GqlInterface struct {
	GqlType
}

// GqlInput represents a GraphQL input type
type GqlInput struct {
	GqlType
}

func (g GqlType) FromReflection(rType reflect.Type, fields []reflect.StructField, methods []reflect.Method) GqlType {
	return GqlType{
		Name: rType.Name(),
		Fields: []Field{{
			Name:    "yo",
			GqlType: "string",
		}},
	}
}
