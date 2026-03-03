// Package main provides the gig CLI tool.
//
// The session module provides the core REPL session management including
// state persistence, command handling, and the main REPL loop.
package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/peterh/liner"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages" // Import built-in stdlib
)

// Type name constants to avoid repetition.
const (
	typeInt     = "int"
	typeFloat64 = "float64"
	typeString  = "string"
	typeBool    = "bool"
	typeAny     = "any"
)

// Session maintains the REPL state.
type Session struct {
	declarations strings.Builder    // Accumulated declarations (var/const/type/func)
	imports      map[string]bool    // Imported packages
	vars         map[string]varInfo // Captured variables with their values
	exprCounter  int                // Counter for generated function names
	timeout      time.Duration      // Execution timeout
	pluginMgr    *PluginManager     // Plugin manager for hot-loading
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

// completer provides tab completion for the REPL.
func (s *Session) completer(line string) []string {
	var completions []string

	// Get the word being typed (last word in line)
	line = strings.TrimRight(line, " \t")
	lastSpace := strings.LastIndexAny(line, " \t(")
	var prefix string
	var wordStart int
	if lastSpace >= 0 {
		prefix = line[:lastSpace+1]
		wordStart = lastSpace + 1
	} else {
		wordStart = 0
	}

	word := line[wordStart:]

	// If word is empty, no completions
	if word == "" {
		return completions
	}

	// Check if completing a selector expression (e.g., fmt.)
	if strings.Contains(word, ".") {
		completions = s.completeSelector(word, prefix)
	} else {
		// Complete simple identifiers
		completions = s.completeIdentifier(word, prefix)
	}

	return completions
}

// completeSelector completes package.Symbol expressions.
func (s *Session) completeSelector(word, prefix string) []string {
	var completions []string

	// Split word into package and partial symbol
	dotIdx := strings.LastIndex(word, ".")
	pkgName := word[:dotIdx]
	partialSymbol := word[dotIdx+1:]

	// Get the package symbols
	symbols := s.getPackageSymbols(pkgName)

	// Filter by prefix and add completions
	for _, sym := range symbols {
		if strings.HasPrefix(sym, partialSymbol) {
			completion := prefix + pkgName + "." + sym
			completions = append(completions, completion)
		}
	}

	return completions
}

// completeIdentifier completes simple identifiers (variables, packages, commands).
func (s *Session) completeIdentifier(word, prefix string) []string {
	var completions []string

	// Check if this is a command (starts with :)
	if strings.HasPrefix(word, ":") {
		commands := []string{":help", ":quit", ":clear", ":imports", ":vars", ":env", ":source", ":timeout", ":plugins"}
		for _, cmd := range commands {
			if strings.HasPrefix(cmd, word) {
				completions = append(completions, prefix+cmd)
			}
		}
		return completions
	}

	// Add variable names
	for name := range s.vars {
		if strings.HasPrefix(name, word) {
			completions = append(completions, prefix+name)
		}
	}

	// Add package names (from imports)
	for path := range s.imports {
		pkgName := path
		if idx := strings.LastIndex(path, "/"); idx >= 0 {
			pkgName = path[idx+1:]
		}
		// Check for known packages with different names
		if name, ok := pkgNameFromPath(path); ok {
			pkgName = name
		}
		if strings.HasPrefix(pkgName, word) {
			completions = append(completions, prefix+pkgName)
		}
	}

	// Add known packages
	for name := range knownPackages {
		if strings.HasPrefix(name, word) {
			completions = append(completions, prefix+name)
		}
	}

	// Add registered packages
	for _, pkg := range gig.GetAllPackages() {
		if strings.HasPrefix(pkg.Name, word) {
			completions = append(completions, prefix+pkg.Name)
		}
	}

	return completions
}

// getPackageSymbols returns exported symbols for a package.
func (s *Session) getPackageSymbols(pkgName string) []string {
	var symbols []string

	// Check registered packages first
	if pkg := gig.GetPackageByName(pkgName); pkg != nil {
		// We need to get symbols from the package
		// Since gig doesn't expose symbols directly, we'll use the plugin manager
		if s.pluginMgr != nil {
			symbols = s.pluginMgr.GetSymbols(pkg.Path)
		}
	}

	// Check built-in stdlib packages
	for path := range s.imports {
		name, ok := pkgNameFromPath(path)
		if ok && name == pkgName {
			if s.pluginMgr != nil {
				syms := s.pluginMgr.GetSymbols(path)
				symbols = append(symbols, syms...)
			}
		}
	}

	return symbols
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
			s.pluginMgr = NewPluginManager()
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
	funcName := fmt.Sprintf("_expr_%d", s.exprCounter)

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
	funcName := fmt.Sprintf("_stmt_%d", s.exprCounter)

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

// autoImportFromInput detects package usage in input and auto-imports known packages.
func (s *Session) autoImportFromInput(input string) {
	// Parse the input to find selector expressions (e.g., fmt.Println)
	fset := token.NewFileSet()
	testSource := "package main\nfunc _() { " + input + " }"
	file, err := parser.ParseFile(fset, "test.go", testSource, 0)
	if err != nil {
		// Try as expression only
		testSource = "package main\nvar _ = " + input
		file, err = parser.ParseFile(fset, "test.go", testSource, 0)
		if err != nil {
			return
		}
	}

	// Find all selector expressions and check if the X part is a known package
	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		// Check if this identifier matches a known package
		pkgName := ident.Name
		if path, ok := knownPackages[pkgName]; ok {
			s.imports[path] = true
			return true
		}

		// Also check registered packages (including hot-loaded ones)
		if pkg := gig.GetPackageByName(pkgName); pkg != nil {
			s.imports[pkg.Path] = true
		}

		return true
	})
}

// extractImports extracts import paths from the input.
func (s *Session) extractImports(input string) {
	// Simple regex-free extraction of import paths
	// Look for patterns like "pkg" or `pkg`
	fset := token.NewFileSet()
	testSource := "package main\n" + input
	file, err := parser.ParseFile(fset, "test.go", testSource, parser.ImportsOnly)
	if err != nil {
		return
	}

	for _, imp := range file.Imports {
		path, _ := strconv.Unquote(imp.Path.Value)
		s.imports[path] = true
	}
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
	captureFunc := fmt.Sprintf("_cap_%d", s.exprCounter)

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

	fset := token.NewFileSet()
	testSource := "package main\nfunc _() {\n" + stmt + "\n}"
	file, err := parser.ParseFile(fset, "test.go", testSource, 0)
	if err != nil {
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

// --- Source Building Methods ---

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
		fmt.Fprintf(&sb, "func _var_%s() %s { return %s }\n", name, info.typeName, literal)
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

	// Parse the input to find selector expressions
	fset := token.NewFileSet()
	testSource := "package main\nfunc _() { " + input + " }"
	file, err := parser.ParseFile(fset, "test.go", testSource, 0)
	if err != nil {
		// Try as expression only
		testSource = "package main\nvar _ = " + input
		file, err = parser.ParseFile(fset, "test.go", testSource, 0)
		if err != nil {
			return used
		}
	}

	// Find all selector expressions
	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		// Check if this identifier matches a known package
		pkgName := ident.Name
		if path, ok := knownPackages[pkgName]; ok {
			used[path] = true
			return true
		}

		// Also check registered packages (including hot-loaded ones)
		if pkg := gig.GetPackageByName(pkgName); pkg != nil {
			used[pkg.Path] = true
		}

		return true
	})

	return used
}

// transformVars replaces variable references with getter calls in the input.
func (s *Session) transformVars(input string) string {
	if len(s.vars) == 0 {
		return input
	}

	// Parse the input to find and replace variable references
	fset := token.NewFileSet()
	// Try as statement first
	testSource := "package main\nfunc _() {\n" + input + "\n}"
	file, err := parser.ParseFile(fset, "test.go", testSource, 0)
	if err != nil {
		// Try as expression
		testSource = "package main\nvar _ = " + input
		file, err = parser.ParseFile(fset, "test.go", testSource, 0)
		if err != nil {
			return input
		}
	}

	// Track positions to replace
	type replacement struct {
		start, end int
		newText    string
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
			// Get position from token.File
			pos := fset.Position(ident.Pos())
			// Adjust for the wrapper we added
			offset := pos.Offset

			// Calculate the actual offset in the original input
			// We need to account for the wrapper prefix
			var wrapperLen int
			if strings.HasPrefix(testSource, "package main\nfunc _() {\n") {
				wrapperLen = len("package main\nfunc _() {\n")
			} else {
				wrapperLen = len("package main\nvar _ = ")
			}

			if offset >= wrapperLen {
				actualStart := offset - wrapperLen
				actualEnd := actualStart + len(ident.Name)
				replacements = append(replacements, replacement{
					start:   actualStart,
					end:     actualEnd,
					newText: "_var_" + ident.Name + "()",
				})
			}
		}
		return true
	})

	// Apply replacements from end to start to preserve positions
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

// transformVarsInStatement transforms variable references but skips LHS of short declarations.
func (s *Session) transformVarsInStatement(input string) string {
	if len(s.vars) == 0 {
		return input
	}

	// First, find positions that should NOT be transformed (LHS of :=)
	fset := token.NewFileSet()
	testSource := "package main\nfunc _() {\n" + input + "\n}"
	file, err := parser.ParseFile(fset, "test.go", testSource, 0)
	if err != nil {
		return input
	}

	wrapperLen := len("package main\nfunc _() {\n")

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
	type replacement struct {
		start, end int
		newText    string
	}
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
					newText: "_var_" + ident.Name + "()",
				})
			}
		}
		return true
	})

	if len(replacements) == 0 {
		return input
	}

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
