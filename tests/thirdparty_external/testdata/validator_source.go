package main

import (
	validator "github.com/go-playground/validator/v10"
)

// ValidatorVarValid tests variable-level validation — uses init-constructed validator
var v *validator.Validate

func init() {
	v = validator.New()
}

// ValidatorVarValid tests that a valid email passes
func ValidatorVarValid() bool {
	err := v.Var("test@example.com", "required,email")
	return err == nil
}

// ValidatorVarInvalid tests that invalid input is caught
func ValidatorVarInvalid() bool {
	err := v.Var("not-an-email", "email")
	return err != nil
}
