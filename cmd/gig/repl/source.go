// Package repl provides the core REPL session management.
//
// This file contains source code generation and import detection for the REPL session.
package repl

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	"git.woa.com/youngjin/gig"
)

// buildExpressionSource creates source code for expression evaluation.
func (s *Session) buildExpressionSource(expr, funcName string) string {
	var sb strings.Builder

	// Transform variable references to getter calls
	transformed := s.transformVars(expr)

	// Package declaration
	sb.WriteString("package main\n\n")

	// Imports - include fmt for main() plus those needed for the expression
	imports := s.detectUsedPackages(expr)
	imports["fmt"] = true // always need fmt for main()
	sb.WriteString(s.buildImportsFromSet(imports))

	// Inject captured variables as getter functions
	sb.WriteString(s.buildCapturedVars())

	// Accumulated declarations (funcs, types)
	sb.WriteString(s.declarations.String())
	sb.WriteString("\n")

	// Main function that prints the result
	sb.WriteString("func main() {\n")
	fmt.Fprintf(&sb, "\tfmt.Println(%s)\n", transformed)
	sb.WriteString("}\n\n")

	// Result function for programmatic access
	fmt.Fprintf(&sb, "func %s() any {\n", funcName)
	fmt.Fprintf(&sb, "\treturn %s\n", transformed)
	sb.WriteString("}\n")

	return sb.String()
}

// buildVoidExpressionSource creates source code for void expression evaluation
// (expressions that don't return a value, like fmt.Println or multi-return calls).
func (s *Session) buildVoidExpressionSource(expr, funcName string) string {
	var sb strings.Builder

	// Transform variable references to getter calls
	transformed := s.transformVars(expr)

	// Package declaration
	sb.WriteString("package main\n\n")

	// Imports
	imports := s.detectUsedPackages(expr)
	sb.WriteString(s.buildImportsFromSet(imports))

	// Inject captured variables as getter functions
	sb.WriteString(s.buildCapturedVars())

	// Accumulated declarations (funcs, types)
	sb.WriteString(s.declarations.String())
	sb.WriteString("\n")

	// Main function that executes the expression
	sb.WriteString("func main() {\n")
	fmt.Fprintf(&sb, "\t%s\n", transformed)
	sb.WriteString("}\n\n")

	// Void wrapper function (no return type, just execute)
	fmt.Fprintf(&sb, "func %s() {\n", funcName)
	fmt.Fprintf(&sb, "\t%s\n", transformed)
	sb.WriteString("}\n")

	return sb.String()
}

// buildStatementSource creates source code for statement execution.
func (s *Session) buildStatementSource(stmt, funcName string) string {
	var sb strings.Builder

	// Transform variable references to getter calls (but not on LHS of :=)
	transformed := s.transformVarsInStatement(stmt)

	// Package declaration
	sb.WriteString("package main\n\n")

	// Imports - only those needed for the current statement
	sb.WriteString(s.buildImportsForInput(stmt))

	// Inject captured variables as getter functions
	sb.WriteString(s.buildCapturedVars())

	// Accumulated declarations (funcs, types)
	sb.WriteString(s.declarations.String())
	sb.WriteString("\n")

	// Main function that executes the statement
	sb.WriteString("func main() {\n")
	fmt.Fprintf(&sb, "\t%s\n", transformed)

	// Suppress "declared and not used" errors by referencing variables
	varNames := s.extractVarNamesFromStatement(stmt)
	for _, name := range varNames {
		fmt.Fprintf(&sb, "\t_ = %s // suppress unused variable error\n", name)
	}

	sb.WriteString("}\n\n")

	// Wrapper function for programmatic access (same as main)
	fmt.Fprintf(&sb, "func %s() {\n", funcName)
	fmt.Fprintf(&sb, "\t%s\n", transformed)
	for _, name := range varNames {
		fmt.Fprintf(&sb, "\t_ = %s\n", name)
	}
	sb.WriteString("}\n")

	return sb.String()
}

// buildCapturedVars generates getter functions for captured variables.
func (s *Session) buildCapturedVars() string {
	if len(s.vars) == 0 {
		return ""
	}

	var sb strings.Builder
	for name, info := range s.vars {
		literal := s.formatValue(info.value)
		// Generate a getter function that returns the value
		fmt.Fprintf(&sb, "func %s%s() %s { return %s }\n", varFuncPrefix, name, info.typeName, literal)
	}
	sb.WriteString("\n")
	return sb.String()
}

// buildImportsForInput generates import declarations only for packages used in the input.
func (s *Session) buildImportsForInput(input string) string {
	usedImports := s.detectUsedPackages(input)
	return s.buildImportsFromSet(usedImports)
}

// buildImportsFromSet generates import declarations from a set of package paths.
func (s *Session) buildImportsFromSet(imports map[string]bool) string {
	if len(imports) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("import (\n")
	for pkg := range imports {
		fmt.Fprintf(&sb, "\t%q\n", pkg)
	}
	sb.WriteString(")\n\n")
	return sb.String()
}

// detectUsedPackages detects which known packages are used in the input.
func (s *Session) detectUsedPackages(input string) map[string]bool {
	used := make(map[string]bool)

	file, _, ok := tryParse(input)
	if !ok {
		return used
	}

	for pkgName, path := range findSelectorPackages(file) {
		_ = pkgName
		used[path] = true
	}

	return used
}

// autoImportFromInput detects package usage in input and auto-imports known packages.
func (s *Session) autoImportFromInput(input string) {
	file, _, ok := tryParse(input)
	if !ok {
		return
	}

	for _, path := range findSelectorPackages(file) {
		s.imports[path] = true
	}
}

// extractImports extracts import paths from the input.
func (s *Session) extractImports(input string) {
	file, _, ok := tryParseAsTopLevel(input)
	if !ok {
		return
	}

	for _, imp := range file.Imports {
		path, _ := strconv.Unquote(imp.Path.Value)
		s.imports[path] = true
	}
}

// findSelectorPackages walks an AST and returns import paths for
// any selector expressions whose X matches a known or registered package.
// Returns a map of pkgName → importPath.
func findSelectorPackages(file *ast.File) map[string]string {
	result := make(map[string]string)

	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		pkgName := ident.Name
		if path, ok := knownPackages[pkgName]; ok {
			result[pkgName] = path
			return true
		}

		// Also check registered packages (including hot-loaded ones)
		if pkg := gig.GetPackageByName(pkgName); pkg != nil {
			result[pkgName] = pkg.Path
		}

		return true
	})

	return result
}
