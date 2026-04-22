package main

import "github.com/iancoleman/strcase"

func StrcaseToCamel() string {
	return strcase.ToCamel("hello_world")
}

func StrcaseToLowerCamel() string {
	return strcase.ToLowerCamel("hello_world")
}

func StrcaseToSnake() string {
	return strcase.ToSnake("HelloWorld")
}

func StrcaseToScreamingSnake() string {
	return strcase.ToScreamingSnake("HelloWorld")
}

func StrcaseToKebab() string {
	return strcase.ToKebab("HelloWorld")
}

func StrcaseToScreamingKebab() string {
	return strcase.ToScreamingKebab("HelloWorld")
}

func StrcaseToDelimited() string {
	return strcase.ToDelimited("HelloWorld", '.')
}

func StrcaseToScreamingDelimited() string {
	return strcase.ToScreamingDelimited("HelloWorld", '.', "", true)
}

func StrcaseToSnakeWithIgnore() string {
	return strcase.ToSnakeWithIgnore("HelloWorldAPI", "API")
}

func StrcaseConfigureAcronym() string {
	strcase.ConfigureAcronym("API", "api")
	return strcase.ToCamel("my_api_key")
}
