package definition

import "deco/graphql/language/ast"

type Infos struct {
	Field     Field
	Requested ast.Field
	Schema    Schema
}
