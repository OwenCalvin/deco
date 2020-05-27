package queries

import (
	"deco/graphql/definition"
	"deco/graphql/server"
)

var Schema = definition.Field{
	DefaultValue: &server.Schema,
	Name:         "__schema",
	TypeRef:      server.Schema,
	Type:         "Schema",
}
