// Package gig provides a Go interpreter with SSA-to-bytecode compilation and VM execution.
// Gig (Go Interpreter in Go) is designed for high-performance interpretation of Go code
// within a Go application, suitable for rule engines, scripting, and embedded logic.
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

	"gig/compiler"
	"gig/importer"
	"gig/value"
	"gig/vm"
)

// DefaultTimeout is the default execution timeout.
const DefaultTimeout = 10 * time.Second

// ErrTimeout is returned when execution times out.
var ErrTimeout = context.DeadlineExceeded

// Program represents a compiled Go program.
type Program struct {
	program *compiler.Program
	ssaPkg  *ssa.Package
}

// Build compiles Go source code into a Program.
// The source must define a function that can be called via Run/RunWithContext.
func Build(sourceCode string, packages ...string) (*Program, error) {
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

	// Create the main package from the type-checked package
	ssaPkg := prog.CreatePackage(pkg, []*ast.File{file}, info, true)
	if ssaPkg == nil {
		return nil, fmt.Errorf("failed to create SSA package")
	}

	// Build the package
	ssaPkg.Build()

	// Wrap external values
	externalValueWrap(ssaPkg, imp)

	// Compile to bytecode
	compiled, err := compiler.Compile(ssaPkg)
	if err != nil {
		return nil, fmt.Errorf("compile error: %w", err)
	}

	return &Program{
		program: compiled,
		ssaPkg:  ssaPkg,
	}, nil
}

// Run executes a function in the program with the given arguments.
func (p *Program) Run(funcName string, params ...interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return p.RunWithContext(funcName, ctx, params...)
}

// RunWithContext executes a function in the program with context for timeout control.
func (p *Program) RunWithContext(funcName string, ctx context.Context, params ...interface{}) (interface{}, error) {
	// Convert params to Value
	args := make([]value.Value, len(params))
	for i, param := range params {
		args[i] = value.FromInterface(param)
	}

	// Create VM and execute
	virtualMachine := vm.New(p.program)
	result, err := virtualMachine.ExecuteWithValues(funcName, ctx, args)
	if err != nil {
		return nil, err
	}

	return result.Interface(), nil
}

// RunWithValues executes a function with pre-converted Value arguments.
func (p *Program) RunWithValues(funcName string, ctx context.Context, args []value.Value) (value.Value, error) {
	virtualMachine := vm.New(p.program)
	return virtualMachine.ExecuteWithValues(funcName, ctx, args)
}

// checkBannedImports checks for banned imports (unsafe, reflect).
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
func autoImport(file *ast.File) {
	// Scan for identifiers that match package names
	usedPackages := make(map[string]bool)

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := node.X.(*ast.Ident); ok {
				// Check if this is a known package
				if _, pkg, ok := importer.AutoImport(ident.Name); ok {
					usedPackages[pkg.Path] = true
				}
			}
		}
		return true
	})

	// Add import declarations for used packages
	// (In a real implementation, we'd modify the AST)
	_ = usedPackages
}

// checkPanicUsage checks for panic usage in the code.
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

// externalValueWrap wraps external package references in SSA.
func externalValueWrap(ssaPkg *ssa.Package, imp *importer.Importer) {
	// Iterate through all functions in the package
	for _, member := range ssaPkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			// Process each function
			for _, block := range fn.Blocks {
				for _, instr := range block.Instrs {
					// Look for external function calls
					if call, ok := instr.(*ssa.Call); ok {
						if fn, ok := call.Call.Value.(*ssa.Function); ok {
							if fn.Pkg == nil || fn.Pkg != ssaPkg {
								// External function - would wrap here
							}
						}
					}
				}
			}
		}
	}
}

// RegisterPackage registers an external package for use in interpreted code.
func RegisterPackage(path, name string) *importer.ExternalPackage {
	return importer.RegisterPackage(path, name)
}

// GetPackageByPath returns a registered package by import path.
func GetPackageByPath(path string) *importer.ExternalPackage {
	return importer.GetPackageByPath(path)
}

// GetPackageByName returns a registered package by name.
func GetPackageByName(name string) *importer.ExternalPackage {
	return importer.GetPackageByName(name)
}

// GetAllPackages returns all registered packages.
func GetAllPackages() map[string]*importer.ExternalPackage {
	return importer.GetAllPackages()
}
