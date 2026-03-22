// Package parser provides source code parsing, type checking, and validation
// for the Gig interpreter compiler pipeline.
//
// It encapsulates go/parser, go/ast, and go/types, providing a clean API for
// the compiler aggregate root.
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"

	"github.com/t04dJ14n9/gig/importer"
)

// ParseResult holds the output of parsing and type-checking a source file.
type ParseResult struct {
	File *ast.File
	FSet *token.FileSet
	Info *types.Info
	Pkg  *types.Package
}

// Parse parses Go source code, performs type checking, and validates it.
// It handles auto-import of registered packages and checks for banned constructs.
//
// The src parameter is the source code string.
// The reg parameter provides package registry for import resolution.
func Parse(src string, reg importer.PackageRegistry) (*ParseResult, error) {
	// Auto-wrap with "package main" if no package declaration
	src = strings.TrimSpace(src)
	if !strings.HasPrefix(src, "package ") {
		src = "package main\n\n" + src
	}

	fset := token.NewFileSet()

	// Parse source code
	file, err := parser.ParseFile(fset, "main.go", src, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	// Auto-import registered packages if needed
	autoImport(file, reg)

	// Create type checker with custom importer
	imp := importer.NewImporter(reg)
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

	return &ParseResult{
		File: file,
		FSet: fset,
		Info: info,
		Pkg:  pkg,
	}, nil
}

// autoImport automatically adds imports for registered packages if used.
func autoImport(file *ast.File, reg importer.PackageRegistry) {
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
				if pkgPath, _, ok := reg.AutoImport(ident.Name); ok {
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

		importSpec := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"` + pkgPath + `"`,
			},
		}

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
func checkPanicUsage(file *ast.File, info *types.Info) error {
	found := false

	ast.Inspect(file, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				if ident.Name == "panic" {
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
