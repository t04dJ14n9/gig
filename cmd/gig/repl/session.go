// Package repl provides the core REPL session management.
//
// This file owns the interactive session lifecycle and input routing.
package repl

import (
	"fmt"
	"go/parser"
	"strings"
	"time"

	_ "github.com/t04dJ14n9/gig/stdlib/packages" // Import built-in stdlib
	"github.com/peterh/liner"

	"github.com/t04dJ14n9/gig/cmd/gig/pluginmgr"
)

// Type name constants to avoid repetition.
const (
	typeInt     = "int"
	typeFloat64 = "float64"
	typeString  = "string"
	typeBool    = "bool"
	typeAny     = "any"
)

// Generated function name prefixes.
const (
	exprFuncPrefix = "_expr_"
	stmtFuncPrefix = "_stmt_"
	capFuncPrefix  = "_cap_"
	varFuncPrefix  = "_var_"
)

// Session maintains the REPL state.
type Session struct {
	declarations strings.Builder    // Accumulated declarations (var/const/type/func)
	imports      map[string]bool    // Imported packages
	vars         map[string]varInfo // Captured variables with their values
	exprCounter  int                // Counter for generated function names
	timeout      time.Duration      // Execution timeout
	pluginMgr    *pluginmgr.Manager // Plugin manager for hot-loading
}

// varInfo stores variable type and value information.
type varInfo struct {
	typeName string
	value    any
}

// NewSession creates a new REPL session.
func NewSession() *Session {
	return &Session{
		imports: make(map[string]bool),
		vars:    make(map[string]varInfo),
		timeout: 10 * time.Second,
	}
}

// Run starts the REPL main loop.
func (s *Session) Run() {
	fmt.Println("Gig REPL - Interactive Go Interpreter")
	fmt.Println("Type :help for commands, :quit to exit")
	fmt.Println()

	// Create liner for advanced line editing with tab completion
	line := liner.NewLiner()
	defer func() { _ = line.Close() }()

	// Configure liner
	line.SetCtrlCAborts(true)
	line.SetTabCompletionStyle(liner.TabCircular)

	// Set up tab completion
	line.SetCompleter(s.completer)

	var multiline strings.Builder
	var inMultiline bool

	for {
		// Show prompt
		var prompt string
		if inMultiline {
			prompt = "... "
		} else {
			prompt = ">>> "
		}

		// Read input line
		input, err := line.Prompt(prompt)
		if err != nil {
			break // EOF or error
		}

		// Add to history
		line.AppendHistory(input)

		// Handle multiline input
		if inMultiline {
			if input == "" {
				// Empty line ends multiline input
				code := multiline.String()
				multiline.Reset()
				inMultiline = false
				s.processInput(code)
			} else {
				multiline.WriteString("\n")
				multiline.WriteString(input)
				// Check if input is now complete
				if !needsMoreInput(multiline.String()) {
					code := multiline.String()
					multiline.Reset()
					inMultiline = false
					s.processInput(code)
				}
			}
			continue
		}

		// Check for multiline start (ends with { or unclosed brackets)
		if needsMoreInput(input) {
			multiline.WriteString(input)
			inMultiline = true
			continue
		}

		// Process single-line input
		s.processInput(input)
	}
}

// processInput handles a complete input.
func (s *Session) processInput(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// Classify input type
	inputType := s.classifyInput(input)

	switch inputType {
	case TypeCommand:
		s.handleCommand(input)
	case TypeImport:
		s.handleImport(input)
	case TypeExpression:
		s.handleExpression(input)
	case TypeStatement:
		s.handleStatement(input)
	case TypeDeclaration:
		s.handleDeclaration(input)
	case TypeInvalid:
		fmt.Printf("Error: invalid syntax\n")
	}
}

// classifyInput determines the type of user input.
func (s *Session) classifyInput(input string) InputType {
	input = strings.TrimSpace(input)

	// 1. Check for commands
	if strings.HasPrefix(input, ":") {
		return TypeCommand
	}

	// 2. Check for import statements
	if isImport(input) {
		return TypeImport
	}

	// 3. Check for declarations
	if isDeclaration(input) {
		return TypeDeclaration
	}

	// 4. Check for expression
	if _, err := parser.ParseExpr(input); err == nil {
		return TypeExpression
	}

	// 5. Check for statement (wrap in function and try to parse)
	if isStatement(input) {
		return TypeStatement
	}

	return TypeInvalid
}
