// Package repl provides the core REPL session management including
// state persistence, command handling, and the main REPL loop.
package repl

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages" // Import built-in stdlib
	"github.com/peterh/liner"

	"git.woa.com/youngjin/gig/cmd/gig/pluginmgr"
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

// handleCommand processes REPL commands.
func (s *Session) handleCommand(input string) {
	cmd := strings.TrimSpace(strings.TrimPrefix(input, ":"))
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "help", "h", "?":
		s.showHelp()
	case "quit", "q", "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "clear":
		s.clear()
		fmt.Println("Session cleared.")
	case "imports":
		s.showImports()
	case "vars":
		s.showVars()
	case "source":
		s.showSource()
	case "timeout":
		s.setTimeout(parts)
	case "plugins":
		s.showPlugins()
	case "env":
		s.showEnv()
	default:
		fmt.Printf("Unknown command: :%s\n", parts[0])
		fmt.Println("Type :help for available commands.")
	}
}

// handleImport processes import statements.
func (s *Session) handleImport(input string) {
	// Extract import path(s) from the statement
	fset := token.NewFileSet()
	testSource := "package main\n" + input
	file, err := parser.ParseFile(fset, "test.go", testSource, parser.ImportsOnly)
	if err != nil {
		fmt.Printf("Error parsing import: %v\n", err)
		return
	}

	for _, imp := range file.Imports {
		path, _ := strconv.Unquote(imp.Path.Value)

		// Check if package is already registered in gig
		if gig.GetPackageByPath(path) != nil {
			s.imports[path] = true
			fmt.Printf("Imported: %s\n", path)
			continue
		}

		// Try to hot-load the package via plugin system
		if s.pluginMgr == nil {
			s.pluginMgr = pluginmgr.NewManager()
		}

		if err := s.pluginMgr.LoadPackage(path); err != nil {
			fmt.Printf("Error loading package %s: %v\n", path, err)
			continue
		}

		s.imports[path] = true
		fmt.Printf("Imported: %s (hot-loaded)\n", path)
	}
}

// showHelp displays available commands.
func (s *Session) showHelp() {
	fmt.Print(`
Commands:
  :help, :h, :?    Show this help message
  :quit, :q        Exit the REPL
  :clear           Clear the session (reset all state)
  :imports         Show imported packages
  :vars            Show captured variables
  :env             Show all variables and packages
  :source          Show accumulated declarations
  :timeout <dur>   Set execution timeout (e.g., :timeout 5s)
  :plugins         Show loaded plugins

Input Types:
  Import       -> import "fmt"
  Expression   -> 1+1, x + y
  Statement    -> x := 1, for i:=0; i<5; i++ {}
  Declaration  -> func foo() {}, type Point struct{}

Examples:
  >>> import "fmt"
  >>> x := 5
  >>> y := x + 3
  >>> x * y
  40
  >>> :vars
  Variables:
    x int = 5
    y int = 8
`)
}

// clear resets the session state.
func (s *Session) clear() {
	s.declarations.Reset()
	s.imports = make(map[string]bool)
	s.vars = make(map[string]varInfo)
	s.exprCounter = 0
}

// showImports displays imported packages.
func (s *Session) showImports() {
	if len(s.imports) == 0 {
		fmt.Println("No packages imported.")
		return
	}
	fmt.Println("Imported packages:")
	for pkg := range s.imports {
		fmt.Printf("  %s\n", pkg)
	}
}

// showVars displays captured variables.
func (s *Session) showVars() {
	if len(s.vars) == 0 {
		fmt.Println("No variables defined.")
		return
	}
	fmt.Println("Variables:")
	for name, info := range s.vars {
		fmt.Printf("  %s %s = %v\n", name, info.typeName, info.value)
	}
}

// showSource displays accumulated source code.
func (s *Session) showSource() {
	src := s.declarations.String()
	if src == "" {
		fmt.Println("No declarations in session.")
		return
	}
	fmt.Println("Accumulated declarations:")
	fmt.Println(src)
}

// setTimeout sets the execution timeout.
func (s *Session) setTimeout(parts []string) {
	if len(parts) < 2 {
		fmt.Printf("Current timeout: %v\n", s.timeout)
		return
	}
	dur, err := time.ParseDuration(parts[1])
	if err != nil {
		fmt.Printf("Invalid duration: %s\n", parts[1])
		return
	}
	s.timeout = dur
	fmt.Printf("Timeout set to: %v\n", dur)
}

// showPlugins displays loaded plugins.
func (s *Session) showPlugins() {
	if s.pluginMgr == nil {
		fmt.Println("No plugins loaded.")
		return
	}
	plugins := s.pluginMgr.ListLoaded()
	if len(plugins) == 0 {
		fmt.Println("No plugins loaded.")
		return
	}
	fmt.Println("Loaded plugins:")
	for _, p := range plugins {
		fmt.Printf("  %s\n", p)
	}
}

// showEnv displays all variables and packages in the current session.
func (s *Session) showEnv() {
	// Show packages
	fmt.Println("Packages:")
	if len(s.imports) == 0 {
		fmt.Println("  (none)")
	} else {
		// Sort packages for consistent output
		pkgs := make([]string, 0, len(s.imports))
		for pkg := range s.imports {
			pkgs = append(pkgs, pkg)
		}
		sort.Strings(pkgs)
		for _, pkg := range pkgs {
			// Get package name from path
			pkgName := pkg
			if idx := strings.LastIndex(pkg, "/"); idx >= 0 {
				pkgName = pkg[idx+1:]
			}
			// Check if it's a known package with a different name
			if name, ok := pkgNameFromPath(pkg); ok {
				pkgName = name
			}
			fmt.Printf("  %s (%s)\n", pkgName, pkg)
		}
	}

	fmt.Println()

	// Show variables
	fmt.Println("Variables:")
	if len(s.vars) == 0 {
		fmt.Println("  (none)")
	} else {
		// Sort variables for consistent output
		names := make([]string, 0, len(s.vars))
		for name := range s.vars {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			info := s.vars[name]
			fmt.Printf("  %s %s = %v\n", name, info.typeName, info.value)
		}
	}
}

// pkgNameFromPath returns the package name for a given import path.
func pkgNameFromPath(path string) (string, bool) {
	// Check known packages first
	for name, p := range knownPackages {
		if p == path {
			return name, true
		}
	}
	// Check registered packages
	if pkg := gig.GetPackageByPath(path); pkg != nil {
		return pkg.Name, true
	}
	return "", false
}

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

	// Auto-detect imports from the expression
	s.autoImportFromInput(input)

	// Build source: declarations + expression wrapper
	source := s.buildExpressionSource(input, funcName)

	// Compile and run
	result, err := s.compileAndRun(source, funcName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print result
	if result != nil {
		fmt.Printf("%v\n", formatResult(result))
	} else {
		fmt.Println("(nil)")
	}
}

// handleStatement executes a statement.
func (s *Session) handleStatement(input string) {
	s.exprCounter++
	funcName := fmt.Sprintf("%s%d", stmtFuncPrefix, s.exprCounter)

	// Auto-detect imports from the statement
	s.autoImportFromInput(input)

	// Build source: declarations + statement wrapper
	source := s.buildStatementSource(input, funcName)

	// Compile and run
	_, err := s.compileAndRun(source, funcName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Extract and capture short variable declarations
	s.captureShortDeclarations(input)

	fmt.Println()
}

// handleDeclaration adds a declaration to the session.
func (s *Session) handleDeclaration(input string) {
	// Extract imports from the declaration
	s.extractImports(input)

	// Add to declarations
	s.declarations.WriteString(input)
	s.declarations.WriteString("\n\n")

	fmt.Println()
}
