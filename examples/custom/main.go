// Package main demonstrates using gig with custom dependencies.
//
// Setup steps (new workflow):
//
//	# Step 1: Initialize a dependency package
//	gig init -package mydep
//
//	# Step 2: Edit mydep/pkgs.go to add third-party libraries
//
//	# Step 3: Generate registration code
//	gig gen ./mydep
//
//	# Step 4: Import the package in your program
//	import _ "myapp/mydep/packages"
package main

import (
	"fmt"

	_ "myapp/mydep/packages"

	"git.woa.com/youngjin/gig"
)

func main() {
	fmt.Println("=== Custom Dependency Example with gjson ===")
	fmt.Println()

	// This example uses github.com/tidwall/gjson - a third-party library
	source := `
package main

import "fmt"
import "github.com/tidwall/gjson"

// GetUserName extracts the "name" field from a JSON string
func GetUserName(jsonStr string) string {
	return gjson.Get(jsonStr, "name").String()
}

// GetNestedValue extracts a nested value using gjson path syntax
func GetNestedValue(jsonStr string, path string) string {
	return gjson.Get(jsonStr, path).String()
}

// GetMultipleValues extracts multiple values and formats them
func GetMultipleValues(jsonStr string) string {
	name := gjson.Get(jsonStr, "user.name").String()
	age := gjson.Get(jsonStr, "user.age").Int()
	active := gjson.Get(jsonStr, "user.active").Bool()
	return fmt.Sprintf("User: %s, Age: %d, Active: %v", name, age, active)
}
`

	prog, err := gig.Build(source)
	if err != nil {
		panic(err)
	}

	// Test 1: Simple JSON parsing
	json1 := `{"name": "Alice", "age": 30}`
	result, err := prog.Run("GetUserName", json1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON: %s\n", json1)
	fmt.Printf("GetUserName: %v\n\n", result)

	// Test 2: Nested JSON with gjson path
	json2 := `{
		"user": {
			"name": "Bob",
			"age": 25,
			"active": true
		},
		"meta": {
			"created": "2024-01-01"
		}
	}`
	result, err = prog.Run("GetNestedValue", json2, "user.name")
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetNestedValue(json, \"user.name\"): %v\n\n", result)

	// Test 3: Multiple values with formatting
	result, err = prog.Run("GetMultipleValues", json2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetMultipleValues: %v\n", result)
}
