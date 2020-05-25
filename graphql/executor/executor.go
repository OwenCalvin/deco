package executor

import (
	"dego/graphql/definition"
	"dego/graphql/language/ast"
)

type ExecutionParams struct {
	Schema definition.Schema
	AST    *ast.Document
}

func Execute(execParams *ExecutionParams) (res interface{}, err error) {
	for _, def := range execParams.AST.Definitions {
		switch def.(type) {
		case *ast.OperationDefinition:
			fDef := def.(*ast.OperationDefinition)
			executeFields(execParams, fDef.SelectionSet)
		case *ast.FragmentDefinition:
		}
	}
	return nil, nil
}

func executeFields(execParams *ExecutionParams, fieldsRoot *ast.SelectionSet) {
	for _, s := range fieldsRoot.Selections {
		field := s.(*ast.Field)
		execParams.Schema.Execute("Query", field.Name.Value, field)
	}
}
