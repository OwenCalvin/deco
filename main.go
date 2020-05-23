package main

import (
	"dego/graphql"
)

type Admin struct {
	access string
}

type User struct {
	Admin
	Name string
}

func (u *User) Get() {
}

func main() {
	graphql.New(
		"/graphql",
		"8080",
		&User{},
		&Admin{},
	)
}
