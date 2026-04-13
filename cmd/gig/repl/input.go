// Package repl provides the core REPL session management.
package repl

import (
	"go/ast"
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

// Wrapper prefixes used when parsing input snippets.
const (
	stmtWrapper    = "package main\nfunc _() {\n"
	exprWrapper    = "package main\nvar _ = "
	topLvlWrapper  = "package main\n"
	stmtWrapperLen = len(stmtWrapper)
	exprWrapperLen = len(exprWrapper)
)

// tryParseAsStmt wraps input in `func _() { ... }` and parses it.
func tryParseAsStmt(input string) (*ast.File, *token.FileSet, bool) {
	fset := token.NewFileSet()
	src := stmtWrapper + input + "\n}"
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		return nil, nil, false
	}
	return file, fset, true
}

// tryParseAsExpr wraps input in `var _ = ...` and parses it.
func tryParseAsExpr(input string) (*ast.File, *token.FileSet, bool) {
	fset := token.NewFileSet()
	src := exprWrapper + input
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		return nil, nil, false
	}
	return file, fset, true
}

// tryParseAsTopLevel wraps input in `package main\n...` and parses it.
func tryParseAsTopLevel(input string) (*ast.File, *token.FileSet, bool) {
	fset := token.NewFileSet()
	src := topLvlWrapper + input
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		return nil, nil, false
	}
	return file, fset, true
}

// tryParse tries stmt then expr wrappers, returning first success.
func tryParse(input string) (*ast.File, *token.FileSet, bool) {
	if file, fset, ok := tryParseAsStmt(input); ok {
		return file, fset, true
	}
	return tryParseAsExpr(input)
}

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
		_, _, ok := tryParseAsTopLevel(input)
		return ok
	}
	return false
}

// isDeclaration checks if input is a declaration.
func isDeclaration(input string) bool {
	trimmed := strings.TrimSpace(input)

	// var, const, type declarations
	if strings.HasPrefix(trimmed, "var ") ||
		strings.HasPrefix(trimmed, "const ") ||
		strings.HasPrefix(trimmed, "type ") ||
		strings.HasPrefix(trimmed, "func ") {
		_, _, ok := tryParseAsTopLevel(input)
		return ok
	}

	return false
}

// isStatement checks if input is a valid statement.
func isStatement(input string) bool {
	_, _, ok := tryParseAsStmt(input)
	return ok
}
