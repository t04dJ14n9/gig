// Package main provides the gig CLI tool.
package main

import (
	"go/parser"
	"go/token"
	"strings"
)

// InputType represents the type of user input.
type InputType int

const (
	TypeInvalid     InputType = iota
	TypeCommand               // :help, :quit, etc.
	TypeImport                // import "fmt"
	TypeExpression            // 1+1, fmt.Sprintf("hi")
	TypeStatement             // x := 1, for {}
	TypeDeclaration           // var x = 1, func foo() {}
)

// needsMoreInput checks if the input is incomplete and needs more lines.
func needsMoreInput(input string) bool {
	// Count unclosed brackets
	openBraces := strings.Count(input, "{") - strings.Count(input, "}")
	openParens := strings.Count(input, "(") - strings.Count(input, ")")
	openBrackets := strings.Count(input, "[") - strings.Count(input, "]")

	// Check for raw string literals
	hasUnclosedRawString := false
	for i := 0; i < len(input); i++ {
		if input[i] == '`' {
			hasUnclosedRawString = !hasUnclosedRawString
		}
	}

	return openBraces > 0 || openParens > 0 || openBrackets > 0 || hasUnclosedRawString
}

// isImport checks if input is an import statement.
func isImport(input string) bool {
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "import ") {
		testSource := "package main\n" + input
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", testSource, parser.ImportsOnly)
		return err == nil
	}
	return false
}

// isDeclaration checks if input is a declaration.
func isDeclaration(input string) bool {
	// Check for declaration keywords
	trimmed := strings.TrimSpace(input)

	// var, const, type declarations
	if strings.HasPrefix(trimmed, "var ") ||
		strings.HasPrefix(trimmed, "const ") ||
		strings.HasPrefix(trimmed, "type ") ||
		strings.HasPrefix(trimmed, "func ") {
		// Try to parse as a declaration
		testSource := "package main\n" + input
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", testSource, 0)
		return err == nil
	}

	// Short variable declaration is a statement, not a declaration
	// We'll handle it as a statement

	return false
}

// isStatement checks if input is a valid statement.
func isStatement(input string) bool {
	// Wrap input in a function and try to parse
	testSource := "package main\nfunc _() {\n" + input + "\n}"
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "test.go", testSource, 0)
	return err == nil
}
