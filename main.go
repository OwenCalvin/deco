package main

import (
	"deco/graphql"
	"deco/graphql/definition"
	"deco/graphql/server"
)

type Admin struct {
	access string
}

type User struct {
	Admin
	Name string
}

type Query struct{}

func (u *Query) Get(
	l struct{ Name string },
	infos definition.Infos,
) struct {
	Name string
	Z    struct{ F string }
} {
	return struct {
		Name string
		Z    struct{ F string }
	}{Name: "yo", Z: struct{ F string }{F: "d"}}
}

func main() {
	schema := graphql.LoadTypes(
		&Query{},
		&User{},
		&Admin{},
	)

	server.Serve(
		schema,
		"/graphql",
		"8080",
	)
}
