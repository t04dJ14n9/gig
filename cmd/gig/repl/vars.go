// Package repl provides the core REPL session management.
//
// This file contains variable capture, transformation, and formatting
// for the REPL session.
package repl

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

// replacement represents a text replacement in source code.
type replacement struct {
	start, end int
	newText    string
}

// applyReplacements applies text replacements to input, processing from end to start
// to preserve positions. Replacements are sorted descending by start position.
func applyReplacements(input string, replacements []replacement) string {
	if len(replacements) == 0 {
		return input
	}

	// Sort replacements by start position descending
	sort.Slice(replacements, func(i, j int) bool {
		return replacements[i].start > replacements[j].start
	})

	result := input
	for _, r := range replacements {
		if r.start >= 0 && r.end <= len(result) {
			result = result[:r.start] + r.newText + result[r.end:]
		}
	}

	return result
}

// transformVars replaces variable references with getter calls in the input.
func (s *Session) transformVars(input string) string {
	if len(s.vars) == 0 {
		return input
	}

	// Parse the input to find and replace variable references
	file, fset, ok := tryParse(input)
	if !ok {
		return input
	}

	// Determine which wrapper was used to calculate offset
	wrapperLen := stmtWrapperLen
	if _, _, stmtOk := tryParseAsStmt(input); !stmtOk {
		wrapperLen = exprWrapperLen
	}

	var replacements []replacement

	// Find identifiers that match captured variable names
	ast.Inspect(file, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}

		// Check if this identifier is a captured variable
		if _, exists := s.vars[ident.Name]; exists {
			pos := fset.Position(ident.Pos())
			offset := pos.Offset

			if offset >= wrapperLen {
				actualStart := offset - wrapperLen
				actualEnd := actualStart + len(ident.Name)
				replacements = append(replacements, replacement{
					start:   actualStart,
					end:     actualEnd,
					newText: varFuncPrefix + ident.Name + "()",
				})
			}
		}
		return true
	})

	return applyReplacements(input, replacements)
}

// transformVarsInStatement transforms variable references but skips LHS of short declarations.
func (s *Session) transformVarsInStatement(input string) string {
	if len(s.vars) == 0 {
		return input
	}

	// First, find positions that should NOT be transformed (LHS of :=)
	file, fset, ok := tryParseAsStmt(input)
	if !ok {
		return input
	}

	wrapperLen := stmtWrapperLen

	// Track LHS positions of short declarations
	skipPositions := make(map[int]bool)

	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		if assign.Tok == token.DEFINE {
			for _, lhs := range assign.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					pos := fset.Position(ident.Pos())
					offset := pos.Offset - wrapperLen
					if offset >= 0 {
						skipPositions[offset] = true
					}
				}
			}
		}
		return true
	})

	// Now do the transformation, skipping LHS positions
	var replacements []replacement

	ast.Inspect(file, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}

		if _, exists := s.vars[ident.Name]; exists {
			pos := fset.Position(ident.Pos())
			offset := pos.Offset - wrapperLen

			// Skip if this is LHS of :=
			if skipPositions[offset] {
				return true
			}

			if offset >= 0 {
				replacements = append(replacements, replacement{
					start:   offset,
					end:     offset + len(ident.Name),
					newText: varFuncPrefix + ident.Name + "()",
				})
			}
		}
		return true
	})

	return applyReplacements(input, replacements)
}

// captureShortDeclarations extracts short variable declarations (:=)
// and creates persistent var declarations with captured values.
func (s *Session) captureShortDeclarations(input string) {
	// Parse the statement to find short declarations
	varNames := s.extractVarNamesFromStatement(input)
	if len(varNames) == 0 {
		return
	}

	// Capture each variable individually to avoid slice issues
	for _, name := range varNames {
		s.captureSingleVar(input, name)
	}
}

// captureSingleVar captures a single variable's value.
func (s *Session) captureSingleVar(input, varName string) {
	s.exprCounter++
	captureFunc := fmt.Sprintf("%s%d", capFuncPrefix, s.exprCounter)

	// Transform variable references in the statement (but not on LHS)
	transformed := s.transformVarsInStatement(input)

	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString(s.buildImportsForInput(input))

	// Inject previously captured variables
	sb.WriteString(s.buildCapturedVars())

	sb.WriteString(s.declarations.String())
	sb.WriteString("\n")

	// Build capture function that returns the single value
	fmt.Fprintf(&sb, "func %s() any {\n", captureFunc)
	fmt.Fprintf(&sb, "\t%s\n", transformed)
	fmt.Fprintf(&sb, "\treturn %s\n", varName)
	sb.WriteString("}\n")

	// Compile and run capture function
	result, err := s.compileAndRun(sb.String(), captureFunc)
	if err != nil {
		return
	}

	// Store variable info in session
	typeName := s.getTypeName(result)
	s.vars[varName] = varInfo{
		typeName: typeName,
		value:    result,
	}
}

// extractVarNamesFromStatement extracts variable names from top-level short declarations.
// Variables declared inside for/if blocks are not included as they have block scope.
func (s *Session) extractVarNamesFromStatement(stmt string) []string {
	var names []string

	file, _, ok := tryParseAsStmt(stmt)
	if !ok {
		return names
	}

	// Find all short declarations that are direct children of the function body
	ast.Inspect(file, func(n ast.Node) bool {
		// Look for function declarations
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Process statements in the function body
		if fn.Body == nil {
			return true
		}

		for _, s := range fn.Body.List {
			// Only process assign statements at the top level of the function
			assign, ok := s.(*ast.AssignStmt)
			if !ok {
				continue
			}

			// Check if it's a short declaration (:=)
			if assign.Tok == token.DEFINE {
				for _, lhs := range assign.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						names = append(names, ident.Name)
					}
				}
			}
		}

		return false // Don't recurse further
	})

	return names
}

// getTypeName returns the Go type name for a value.
// Normalizes integer types to int for compatibility with for loops.
func (s *Session) getTypeName(v any) string {
	if v == nil {
		return typeAny
	}

	// Use fmt type formatting
	typeStr := fmt.Sprintf("%T", v)

	// Normalize integer types to int for better compatibility
	// This allows captured variables to work with for loop counters
	switch typeStr {
	case typeInt, "int8", "int16", "int32", "int64":
		return typeInt
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return typeInt // use int for unsigned too, for simplicity
	case "float32", typeFloat64:
		return typeFloat64
	case typeString, typeBool, "rune", "byte":
		return typeStr
	default:
		// For complex types, use any
		return typeAny
	}
}

// formatValue formats a value as a Go literal.
func (s *Session) formatValue(v any) string {
	if v == nil {
		return "nil"
	}

	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case bool:
		return fmt.Sprintf("%v", val)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", val)
	default:
		// For complex types, use the default representation
		return fmt.Sprintf("%#v", val)
	}
}

// formatResult formats a result for display.
func formatResult(v any) string {
	if v == nil {
		return "nil"
	}

	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case []byte:
		return fmt.Sprintf("[]byte(%q)", string(val))
	case []any:
		elems := make([]string, len(val))
		for i, e := range val {
			elems[i] = formatResult(e)
		}
		return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
	case map[string]any:
		pairs := make([]string, 0, len(val))
		for k, e := range val {
			pairs = append(pairs, fmt.Sprintf("%q: %s", k, formatResult(e)))
		}
		return fmt.Sprintf("map[string]any{%s}", strings.Join(pairs, ", "))
	default:
		return fmt.Sprintf("%v", val)
	}
}
