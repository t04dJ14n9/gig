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
	"time"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/compiler"
	"github.com/t04dJ14n9/gig/importer"
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
	"github.com/t04dJ14n9/gig/runner"
)

// DefaultTimeout is the default execution timeout.
const DefaultTimeout = 10 * time.Second

// ErrTimeout is returned when execution times out.
var ErrTimeout = context.DeadlineExceeded

// buildConfig holds internal configuration parsed from BuildOption values.
type buildConfig struct {
	registry        importer.PackageRegistry
	statefulGlobals bool
	allowPanic      bool
}

// BuildOption configures the behaviour of Build.
// Use the With* functions to obtain concrete options.
type BuildOption func(*buildConfig)

// WithRegistry sets a custom PackageRegistry for resolving external packages.
// If not provided, Build uses the global registry (pre-populated by init() functions).
func WithRegistry(r importer.PackageRegistry) BuildOption {
	return func(c *buildConfig) {
		c.registry = r
	}
}

// WithStatefulGlobals enables persistent package-level globals across Run calls.
// When enabled, mutations to package-level variables in one Run call are visible
// to subsequent Run calls on the same Program. Multiple concurrent Run calls are
// supported — global variable access is protected by a sync.RWMutex.
//
// By default (when this option is not passed), each Run starts from the
// post-init() global state snapshot and mutations are discarded after the call.
func WithStatefulGlobals() BuildOption {
	return func(c *buildConfig) {
		c.statefulGlobals = true
	}
}

// WithAllowPanic allows the use of panic() in interpreted code.
// By default, panic() is banned at compile time for sandbox safety.
// When enabled, panic/recover/defer work as in standard Go, and unrecovered
// panics are returned as errors rather than crashing the host process.
func WithAllowPanic() BuildOption {
	return func(c *buildConfig) {
		c.allowPanic = true
	}
}

// Program represents a compiled Go program ready for execution.
// It delegates execution to a runner.Runner for VM pool management and global state handling.
type Program struct {
	runner *runner.Runner // execution orchestration (VM pool, stateful globals)
	ssaPkg *ssa.Package   // SSA package for debugging/inspection
}

// InternalProgram exposes the compiled bytecode program for testing/debugging.
func (p *Program) InternalProgram() *bytecode.CompiledProgram { return p.runner.InternalProgram() }

// Close releases resources associated with the Program.
// It unregisters the per-program method resolver to prevent memory leaks.
// Callers should defer prog.Close() after calling Build().
func (p *Program) Close() {
	p.runner.Close()
}

// Build compiles Go source code into a Program.
//
// The source must define a function that can be called via Run/RunWithContext.
// If the source does not start with a package declaration, "package main" is prepended automatically.
//
// Options control runtime behaviour; see WithStatefulGlobals and other With* functions.
//
// The compilation process:
//  1. Parse source code into AST
//  2. Check for banned imports (unsafe, reflect)
//  3. Check for banned panic usage (unless WithAllowPanic is set)
//  4. Type-check with custom importer for external packages
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
func Build(sourceCode string, opts ...BuildOption) (*Program, error) {
	// Parse options
	cfg := buildConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.registry == nil {
		cfg.registry = importer.GlobalRegistry()
	}

	// Compile: parse → SSA → bytecode (full pipeline owned by compiler package)
	var compilerOpts []compiler.BuildOption
	if cfg.allowPanic {
		compilerOpts = append(compilerOpts, compiler.WithAllowPanic())
	}
	result, err := compiler.Build(sourceCode, cfg.registry, compilerOpts...)
	if err != nil {
		return nil, err
	}

	// Run init() and snapshot globals if present
	initialGlobals, err := runner.ExecuteInit(result.Program)
	if err != nil {
		return nil, fmt.Errorf("executing init(): %w", err)
	}

	var runnerOpts []runner.RunnerOption
	if cfg.statefulGlobals {
		runnerOpts = append(runnerOpts, runner.WithStatefulGlobals())
	}

	r := runner.New(result.Program, initialGlobals, runnerOpts...)

	return &Program{
		runner: r,
		ssaPkg: result.SSAPkg,
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
	return p.runner.Run(funcName, params...)
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
	return p.runner.RunWithContext(ctx, funcName, params...)
}

// RunWithValues executes a function with pre-converted Value arguments.
// This is more efficient than Run/RunWithContext when you need to call the same function
// multiple times with the same parameter types, as it avoids repeated type conversion.
// Context is the first parameter following Go idioms.
func (p *Program) RunWithValues(ctx context.Context, funcName string, args []value.Value) (value.Value, error) {
	return p.runner.RunWithValues(ctx, funcName, args)
}

// NewSandboxRegistry creates a fresh, empty PackageRegistry for sandboxed execution.
// Unlike the global registry (which is pre-populated by stdlib init() functions),
// a sandbox registry starts empty, allowing the caller to register only the
// packages they want to expose to interpreted code.
//
// Example:
//
//	reg := gig.NewSandboxRegistry()
//	gig.RegisterPackage("fmt", "fmt") // registers to global registry, not sandbox
//	prog, err := gig.Build(source, gig.WithRegistry(reg))
func NewSandboxRegistry() importer.PackageRegistry {
	return importer.NewRegistry()
}

// RegisterPackage registers an external package for use in interpreted code.
// Delegates to the global importer registry.
func RegisterPackage(path, name string) *importer.ExternalPackage {
	return importer.RegisterPackage(path, name)
}

// GetPackageByPath returns a registered package by import path.
// Delegates to the global importer registry.
func GetPackageByPath(path string) *importer.ExternalPackage {
	return importer.GetPackageByPath(path)
}

// GetPackageByName returns a registered package by name.
// Delegates to the global importer registry.
func GetPackageByName(name string) *importer.ExternalPackage {
	return importer.GetPackageByName(name)
}

// GetAllPackages returns all registered packages.
// Delegates to the global importer registry.
func GetAllPackages() map[string]*importer.ExternalPackage {
	return importer.GetAllPackages()
}
