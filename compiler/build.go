// build.go implements the full compilation pipeline: source → parse → SSA → bytecode.
package compiler

import (
	"fmt"

	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/compiler/parser"
	ssabuilder "git.woa.com/youngjin/gig/compiler/ssa"
	"git.woa.com/youngjin/gig/importer"
)

// BuildResult holds the output of the full compilation pipeline.
type BuildResult struct {
	Program *bytecode.Program
	SSAPkg  *ssa.Package
}

// buildConfig holds internal configuration parsed from BuildOption values.
type buildConfig struct {
	allowPanic bool
}

// BuildOption configures the behaviour of Build.
type BuildOption func(*buildConfig)

// WithAllowPanic allows the use of panic() in compiled code.
func WithAllowPanic() BuildOption {
	return func(c *buildConfig) {
		c.allowPanic = true
	}
}

// Build compiles Go source code into bytecode through the full pipeline:
//
//  1. Parse source → typed AST (compiler/parser)
//  2. Build SSA from typed AST (compiler/ssa)
//  3. Compile SSA to bytecode (codegen)
//
// PackageRegistry is the compiler's primary dependency — it provides package
// resolution for both the type checker and the codegen phase (via PackageLookup).
func Build(source string, reg importer.PackageRegistry, opts ...BuildOption) (*BuildResult, error) {
	cfg := buildConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	// 1. Parse + type-check + validate
	var parseOpts []parser.ParseOption
	if cfg.allowPanic {
		parseOpts = append(parseOpts, parser.WithAllowPanic())
	}
	parseResult, err := parser.Parse(source, reg, parseOpts...)
	if err != nil {
		return nil, err
	}

	// 2. Build SSA
	ssaResult, err := ssabuilder.Build(parseResult.FSet, parseResult.Pkg, parseResult.File, parseResult.Info)
	if err != nil {
		return nil, err
	}

	// 3. Compile SSA to bytecode — PackageLookup is derived from PackageRegistry
	// because resolving external functions/methods is a compiler responsibility.
	lookup := importer.NewPackageLookup(reg)
	compiled, err := NewCompiler(lookup).Compile(ssaResult.Pkg)
	if err != nil {
		return nil, fmt.Errorf("compile error: %w", err)
	}

	return &BuildResult{
		Program: compiled,
		SSAPkg:  ssaResult.Pkg,
	}, nil
}
