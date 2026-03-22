package compiler

import (
	"fmt"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/bytecode"
	"github.com/t04dJ14n9/gig/compiler/parser"
	ssabuilder "github.com/t04dJ14n9/gig/compiler/ssa"
	"github.com/t04dJ14n9/gig/importer"
)

// BuildResult holds the output of the full compilation pipeline.
type BuildResult struct {
	Program *bytecode.Program
	SSAPkg  *ssa.Package
}

// Build compiles Go source code into bytecode through the full pipeline:
//
//  1. Parse source → typed AST (compiler/parser)
//  2. Build SSA from typed AST (compiler/ssa)
//  3. Compile SSA to bytecode (codegen)
//
// PackageRegistry is the compiler's primary dependency — it provides package
// resolution for both the type checker and the codegen phase (via PackageLookup).
func Build(source string, reg importer.PackageRegistry) (*BuildResult, error) {
	// 1. Parse + type-check + validate
	parseResult, err := parser.Parse(source, reg)
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
