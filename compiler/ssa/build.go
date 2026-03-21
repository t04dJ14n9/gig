// Package ssa provides SSA (Static Single Assignment) construction for the
// Gig interpreter compiler pipeline.
//
// It encapsulates golang.org/x/tools/go/ssa, providing a clean API for
// the compiler aggregate root.
package ssa

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
)

// BuildResult holds the output of SSA construction.
type BuildResult struct {
	Prog *ssa.Program
	Pkg  *ssa.Package
}

// Build creates an SSA program from a type-checked package.
// It builds SSA representations for all imported packages and the main package.
func Build(fset *token.FileSet, pkg *types.Package, file *ast.File, info *types.Info) (*BuildResult, error) {
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

	return &BuildResult{
		Prog: prog,
		Pkg:  ssaPkg,
	}, nil
}
