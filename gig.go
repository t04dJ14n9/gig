// Package gig provides a Go interpreter with SSA-to-bytecode compilation and VM execution.
//
// Gig (Go Interpreter in Go) is designed for high-performance interpretation of Go code
// within a Go application, suitable for rule engines, scripting, and embedded logic.
//
// # Overview
//
// Gig compiles Go source code to SSA (Static Single Assignment) form using golang.org/x/tools/go/ssa,
// then translates SSA to a custom bytecode format. The bytecode is executed by a stack-based
// virtual machine with a tagged-union value system for efficient primitive operations.
//
// # Architecture
//
// The interpreter consists of three main components:
//
//  1. Compiler (gig/compiler) - Translates SSA IR to bytecode instructions
//  2. VM (gig/vm) - Stack-based virtual machine for bytecode execution
//  3. Value (gig/value) - Tagged-union value system for efficient type handling
//
// # Security Model
//
// For safety in embedded contexts, Gig bans:
//   - "unsafe" package - prevents raw memory access
//   - "reflect" package - prevents type introspection bypass
//   - "panic" builtin - prevents uncontrolled control flow
//
// # Example Usage
//
// Basic usage with built-in standard library:
//
//	prog, err := gig.Build(`
//		package main
//
//		import "fmt"
//
//		func Greet(name string) string {
//			return fmt.Sprintf("Hello, %s!", name)
//		}
//	`)
//	if err != nil {
//		panic(err)
//	}
//
//	result, err := prog.Run("Greet", "World")
//	fmt.Println(result) // Output: Hello, World!
//
// # External Packages
//
// Gig supports calling external Go packages by registering them before compilation.
// See gig/stdlib for built-in standard library packages, or use the gig CLI tool
// to generate wrappers for third-party libraries.
package gig

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"time"

	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/compiler"
	"git.woa.com/youngjin/gig/importer"
	"git.woa.com/youngjin/gig/value"
	"git.woa.com/youngjin/gig/vm"
)

// DefaultTimeout is the default execution timeout.
const DefaultTimeout = 10 * time.Second

// ErrTimeout is returned when execution times out.
var ErrTimeout = context.DeadlineExceeded

// Program represents a compiled Go program ready for execution.
// It contains the compiled bytecode, constant pool, type information, and SSA package reference.
type Program struct {
	program *bytecode.Program // compiled bytecode and metadata
	ssaPkg  *ssa.Package      // SSA package for debugging/inspection
	vmPool  *vm.VMPool        // reusable VM pool (eliminates 32KB alloc per Run)
}

// Build compiles Go source code into a Program.
//
// The source must define a function that can be called via Run/RunWithContext.
// If the source does not start with a package declaration, "package main" is prepended automatically.
//
// The compilation process:
//  1. Parse source code into AST
//  2. Check for banned imports (unsafe, reflect)
//  3. Type-check with custom importer for external packages
//  4. Check for banned panic usage
//  5. Build SSA intermediate representation
//  6. Compile SSA to bytecode
//
// Example:
//
//	prog, err := gig.Build(`
//		func Add(a, b int) int {
//			return a + b
//		}
//	`)
//	result, _ := prog.Run("Add", 1, 2) // result = 3
func Build(sourceCode string, packages ...string) (*Program, error) {
	// Auto-wrap with "package main" if no package declaration
	sourceCode = strings.TrimSpace(sourceCode)
	if !strings.HasPrefix(sourceCode, "package ") {
		sourceCode = "package main\n\n" + sourceCode
	}

	// Create file set
	fset := token.NewFileSet()

	// Parse source code
	file, err := parser.ParseFile(fset, "main.go", sourceCode, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	// Check for banned imports (unsafe, reflect)
	if err := checkBannedImports(file); err != nil {
		return nil, err
	}

	// Auto-import registered packages if needed
	autoImport(file)

	// Create type checker with custom importer
	imp := importer.NewImporter()
	typeConfig := &types.Config{
		Importer: imp,
		Sizes:    &types.StdSizes{WordSize: 8, MaxAlign: 8},
	}

	// Type check
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	pkg, err := typeConfig.Check("main", fset, []*ast.File{file}, info)
	if err != nil {
		return nil, fmt.Errorf("type check error: %w", err)
	}

	// Check for panic usage (banned)
	if err := checkPanicUsage(file, info); err != nil {
		return nil, err
	}

	// Build SSA
	// Create a new SSA program
	prog := ssa.NewProgram(fset, ssa.SanityCheckFunctions|ssa.GlobalDebug)

	// Create SSA packages for all imported packages
	for _, imp := range pkg.Imports() {
		prog.CreatePackage(imp, nil, nil, true)
	}

	// Create the main package from the type-checked package
	ssaPkg := prog.CreatePackage(pkg, []*ast.File{file}, info, true)
	if ssaPkg == nil {
		return nil, fmt.Errorf("failed to create SSA package")
	}

	// Build the package
	ssaPkg.Build()

	// Compile to bytecode
	lookup := newPackageLookupAdapter()
	compiled, err := compiler.Compile(lookup, ssaPkg)
	if err != nil {
		return nil, fmt.Errorf("compile error: %w", err)
	}

	return &Program{
		program: compiled,
		ssaPkg:  ssaPkg,
		vmPool:  vm.NewVMPool(compiled),
	}, nil
}

// Run executes a function in the program with the given arguments.
// It uses the default timeout (DefaultTimeout = 10 seconds).
// Parameters are automatically converted to value.Value using FromInterface.
//
// Example:
//
//	result, err := prog.Run("Add", 1, 2)
func (p *Program) Run(funcName string, params ...any) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return p.RunWithContext(ctx, funcName, params...)
}

// RunWithContext executes a function in the program with context for timeout control.
// This allows custom timeout values and cancellation.
// Context is the first parameter following Go idioms.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	result, err := prog.RunWithContext(ctx, "LongRunningTask", input)
func (p *Program) RunWithContext(ctx context.Context, funcName string, params ...any) (any, error) {
	// Convert params to Value
	args := make([]value.Value, len(params))
	for i, param := range params {
		args[i] = value.FromInterface(param)
	}

	// Get VM from pool, execute, and return to pool
	virtualMachine := p.vmPool.Get()
	result, err := virtualMachine.ExecuteWithValues(funcName, ctx, args)
	p.vmPool.Put(virtualMachine)
	if err != nil {
		return nil, err
	}

	return result.Interface(), nil
}

// RunWithValues executes a function with pre-converted Value arguments.
// This is more efficient than Run/RunWithContext when you need to call the same function
// multiple times with the same parameter types, as it avoids repeated type conversion.
// Context is the first parameter following Go idioms.
func (p *Program) RunWithValues(ctx context.Context, funcName string, args []value.Value) (value.Value, error) {
	virtualMachine := p.vmPool.Get()
	result, err := virtualMachine.ExecuteWithValues(funcName, ctx, args)
	p.vmPool.Put(virtualMachine)
	return result, err
}

// checkBannedImports checks for banned imports (unsafe, reflect).
// These packages are banned because they can bypass the interpreter's safety guarantees.
// Returns an error if any banned import is found.
func checkBannedImports(file *ast.File) error {
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if path == "unsafe" {
			return fmt.Errorf("import of \"unsafe\" is not allowed")
		}
		if path == "reflect" {
			return fmt.Errorf("import of \"reflect\" is not allowed")
		}
	}
	return nil
}

// autoImport automatically adds imports for registered packages if used.
// This allows users to reference package names without explicit import declarations.
// It scans the AST for selector expressions (e.g., fmt.Println) and checks if
// the identifier matches a registered package name.
func autoImport(file *ast.File) {
	// Collect already-imported paths to avoid duplicates
	alreadyImported := make(map[string]bool)
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		alreadyImported[path] = true
	}

	// Scan for identifiers that match package names
	usedPackages := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := node.X.(*ast.Ident); ok {
				// Check if this is a known package
				if pkgPath, _, ok := importer.AutoImport(ident.Name); ok {
					usedPackages[pkgPath] = true
				}
			}
		}
		return true
	})

	// Inject import declarations into the AST for packages not already imported
	for pkgPath := range usedPackages {
		if alreadyImported[pkgPath] {
			continue
		}

		// Build a new import spec: import "pkgPath"
		importSpec := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"` + pkgPath + `"`,
			},
		}

		// Wrap it in a GenDecl and prepend to file.Decls
		importDecl := &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{importSpec},
		}
		file.Decls = append([]ast.Decl{importDecl}, file.Decls...)
		file.Imports = append(file.Imports, importSpec)

		alreadyImported[pkgPath] = true
	}
}

// checkPanicUsage checks for panic usage in the code.
// Panic is banned because it can cause uncontrolled control flow that bypasses
// the interpreter's execution model.
// Returns an error if the builtin panic function is called.
func checkPanicUsage(file *ast.File, info *types.Info) error {
	found := false

	ast.Inspect(file, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				if ident.Name == "panic" {
					// Check if it's the builtin panic
					if obj := info.Uses[ident]; obj != nil {
						if _, isBuiltin := obj.(*types.Builtin); isBuiltin {
							found = true
							return false
						}
					}
				}
			}
		}
		return true
	})

	if found {
		return fmt.Errorf("use of \"panic\" is not allowed")
	}
	return nil
}

// RegisterPackage
// RegisterPackage registers an external package for use in interpreted code.
// This is typically called by init() functions in generated package wrappers.
// The path is the import path (e.g., "fmt"), and name is the package identifier (e.g., "fmt").
func RegisterPackage(path, name string) *importer.ExternalPackage {
	return importer.RegisterPackage(path, name)
}

// GetPackageByPath returns a registered package by import path.
// Returns nil if no package with the given path is registered.
func GetPackageByPath(path string) *importer.ExternalPackage {
	return importer.GetPackageByPath(path)
}

// GetPackageByName returns a registered package by name.
// This is used for auto-import functionality.
// Returns nil if no package with the given name is registered.
func GetPackageByName(name string) *importer.ExternalPackage {
	return importer.GetPackageByName(name)
}

// GetAllPackages returns all registered packages.
// The returned map is keyed by import path.
func GetAllPackages() map[string]*importer.ExternalPackage {
	return importer.GetAllPackages()
}

// packageLookupAdapter implements bytecode.PackageLookup using the importer registry.
// It serves as the DI bridge between the compiler and the importer.
type packageLookupAdapter struct{}

// newPackageLookupAdapter creates a new PackageLookup adapter.
func newPackageLookupAdapter() bytecode.PackageLookup {
	return &packageLookupAdapter{}
}

// LookupExternalFunc resolves an external function by package path and function name.
func (a *packageLookupAdapter) LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool) {
	pkg := importer.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, nil, false
	}
	obj, exists := pkg.Objects[funcName]
	if !exists || obj.Kind != importer.ObjectKindFunction {
		return nil, nil, false
	}
	return obj.Value, obj.DirectCall, true
}

// LookupMethodDirectCall resolves a method DirectCall wrapper by type name and method name.
func (a *packageLookupAdapter) LookupMethodDirectCall(typeName, methodName string) (directCall func([]value.Value) value.Value, ok bool) {
	return importer.LookupMethodDirectCall(typeName, methodName)
}
