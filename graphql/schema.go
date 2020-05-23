package graphql

import "dego/graphql/definition"

// Schema represents a graphql schema
type Schema struct {
	definitions []definition.GqlType
}
