package executor

import (
	"dego/graphql/definition"
	"dego/graphql/language/ast"
)

type ExecutionParams struct {
	Schema definition.Schema
	AST    *ast.Document
}

func Execute(execParams *ExecutionParams) interface{} {
	for _, def := range execParams.AST.Definitions {
		switch def.(type) {
		case *ast.OperationDefinition:
			return execParams.Schema.Execute("Query", "GetQuery", struct{ Name string }{Name: "yo"})
		case *ast.FragmentDefinition:
		}
	}
	return struct{}{}
}
