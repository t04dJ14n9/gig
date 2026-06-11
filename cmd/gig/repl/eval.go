package repl

import (
	"context"
	"fmt"
	"strings"

	"github.com/t04dJ14n9/gig"
)

// compileAndRun compiles source code and executes a function.
func (s *Session) compileAndRun(source, funcName string) (any, error) {
	prog, err := gig.Build(source)
	if err != nil {
		return nil, fmt.Errorf("compile error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	result, err := prog.RunWithContext(ctx, funcName)
	if err != nil {
		return nil, fmt.Errorf("runtime error: %w", err)
	}

	return result, nil
}

// handleExpression evaluates an expression.
func (s *Session) handleExpression(input string) {
	s.exprCounter++
	funcName := fmt.Sprintf("%s%d", exprFuncPrefix, s.exprCounter)

	s.autoImportFromInput(input)
	source := s.buildExpressionSource(input, funcName)

	result, err := s.compileAndRun(source, funcName)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "too many return values") ||
			strings.Contains(errMsg, "used as value") ||
			strings.Contains(errMsg, "no value") {
			s.handleVoidExpression(input, funcName)
			return
		}
		fmt.Printf("Error: %v\n", err)
		return
	}

	if result != nil {
		fmt.Printf("%v\n", formatResult(result))
	} else {
		fmt.Println("(nil)")
	}
}

// handleVoidExpression handles expressions that don't return a value
// (e.g., fmt.Println("hello"), print(x)).
func (s *Session) handleVoidExpression(input, funcName string) {
	source := s.buildVoidExpressionSource(input, funcName)

	_, err := s.compileAndRun(source, funcName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println()
}

// handleStatement executes a statement.
func (s *Session) handleStatement(input string) {
	s.exprCounter++
	funcName := fmt.Sprintf("%s%d", stmtFuncPrefix, s.exprCounter)

	s.autoImportFromInput(input)
	source := s.buildStatementSource(input, funcName)

	_, err := s.compileAndRun(source, funcName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	s.captureShortDeclarations(input)
	fmt.Println()
}

// handleDeclaration adds a declaration to the session.
func (s *Session) handleDeclaration(input string) {
	s.extractImports(input)
	s.declarations.WriteString(input)
	s.declarations.WriteString("\n\n")
	fmt.Println()
}
