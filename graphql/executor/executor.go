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
	var executable *ast.OperationDefinition
	fragments := make(map[string]*ast.FragmentDefinition)

	for _, def := range execParams.AST.Definitions {
		switch def.(type) {
		case *ast.OperationDefinition:
			executable = def.(*ast.OperationDefinition)
		case *ast.FragmentDefinition:
			f := def.(*ast.FragmentDefinition)
			fragments[f.Name.Value] = f
		}
	}

	executeFields(execParams, fragments, executable.SelectionSet)

	return nil, nil
}

func parseFragment(executable *ast.OperationDefinition, fragment *ast.FragmentDefinition) {
}

func executeFields(execParams *ExecutionParams, fragments map[string]*ast.FragmentDefinition, fieldsRoot *ast.SelectionSet) {
	for i, s := range fieldsRoot.Selections {
		field := s.(*ast.Field)
		operation := execParams.AST.Definitions[i].(*ast.OperationDefinition).Operation

		switch operation {
		case "query":
			operation = "Query"
		case "mutation":
			operation = "Mutation"
		case "subscription":
			operation = "Subscription"
		}

		execParams.Schema.Execute(
			operation,
			field.Name.Value,
			field,
			fragments,
		)
	}
}
