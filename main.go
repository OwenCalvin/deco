package main

import (
	"dego/graphql"
	"dego/graphql/server"
)

type Admin struct {
	access string
}

type User struct {
	Admin
	Name string
}

func (u *User) GetQuery(
	l struct{ Name string },
) struct{ Name string } {
	return struct{ Name string }{Name: "yo"}
}

func main() {
	schema := graphql.LoadTypes(&User{}, &Admin{})

	server.Serve(
		schema,
		"/graphql",
		"8080",
	)
}
