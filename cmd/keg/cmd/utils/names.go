package utils

import "github.com/samber/lo"

type Names struct {
	Name       string
	PascalCase string // HelloWorld
	CamelCase  string // helloWorld
	KebabCase  string // hello-world
	SnakeCase  string // hello_world
}

func NewNames(name string) Names {
	return Names{
		Name:       name,
		PascalCase: lo.PascalCase(name),
		CamelCase:  lo.CamelCase(name),
		KebabCase:  lo.KebabCase(name),
		SnakeCase:  lo.SnakeCase(name),
	}
}
