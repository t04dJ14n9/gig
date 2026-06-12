package main

import (
	mapstructure "github.com/mitchellh/mapstructure"
)

// MapstructureDecode decodes a map into a struct, returns the name field
func MapstructureDecode(input map[string]interface{}) string {
	type Person struct {
		Name string
		Age  int
	}
	var result Person
	mapstructure.Decode(input, &result)
	return result.Name
}

// MapstructureWeakDecode tests weakly-typed decoding
func MapstructureWeakDecode(input map[string]interface{}) int {
	type Config struct {
		Port int
	}
	var result Config
	mapstructure.WeakDecode(input, &result)
	return result.Port
}
