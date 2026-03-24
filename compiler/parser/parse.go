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

	"git.woa.com/youngjin/gig/importer"
)

// ParseResult holds the output of parsing and type-checking a source file.
type ParseResult struct {
	File *ast.File
	FSet *token.FileSet
	Info *types.Info
	Pkg  *types.Package
}

// parseConfig holds internal configuration parsed from ParseOption values.
type parseConfig struct {
	allowPanic bool
}

// ParseOption configures the behaviour of Parse.
type ParseOption func(*parseConfig)

// WithAllowPanic allows the use of panic() in source code.
// By default, panic() calls are rejected at parse time for sandbox safety.
func WithAllowPanic() ParseOption {
	return func(c *parseConfig) {
		c.allowPanic = true
	}
}

// Parse parses Go source code, performs type checking, and validates it.
// It handles auto-import of registered packages and checks for banned constructs.
//
// The src parameter is the source code string.
// The reg parameter provides package registry for import resolution.
func Parse(src string, reg importer.PackageRegistry, opts ...ParseOption) (*ParseResult, error) {
	cfg := parseConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

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

	// Check for banned panic usage (unless explicitly allowed)
	if !cfg.allowPanic {
		if err := checkBannedPanic(fset, file); err != nil {
			return nil, err
		}
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

// checkBannedPanic walks the AST and returns an error if any panic() call is found.
// This is a compile-time safety check for sandboxed execution.
// Both panic("msg") and panic(expr) forms are detected.
func checkBannedPanic(fset *token.FileSet, file *ast.File) error {
	var panicPos token.Pos
	ast.Inspect(file, func(n ast.Node) bool {
		if panicPos.IsValid() {
			return false // already found one, stop
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
			panicPos = ident.Pos()
			return false
		}
		return true
	})
	if panicPos.IsValid() {
		pos := fset.Position(panicPos)
		return fmt.Errorf("compile error: panic() is not allowed in sandboxed code (at %s); use gig.WithAllowPanic() to enable", pos)
	}
	return nil
}
