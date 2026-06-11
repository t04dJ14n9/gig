package main

import "github.com/asaskevich/govalidator"

func GovalidatorIsEmail() bool {
	return govalidator.IsEmail("test@example.com")
}

func GovalidatorIsURL() bool {
	return govalidator.IsURL("https://example.com")
}

func GovalidatorIsAlpha() bool {
	return govalidator.IsAlpha("hello")
}

func GovalidatorToString() string {
	return govalidator.ToString(42)
}
