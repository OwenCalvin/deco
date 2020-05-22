package main

import (
	"dry/reflector"
	"fmt"
)

type admin struct {
	access string
}

type user struct {
	admin
	Name string
}

func main() {
	allReflection := reflector.ReflectAll(
		&user{},
		&admin{},
	)

	fmt.Println(allReflection)
}
