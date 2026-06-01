package tests

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestGetBenchmarksSrcStaysShallow(t *testing.T) {
	count := benchmarkTestDecisionCount(t, "getBenchmarksSrc")
	if count > 8 {
		t.Fatalf("getBenchmarksSrc has %d decision points, want <= 8; parse benchmark files with focused helpers", count)
	}
}

func benchmarkTestDecisionCount(t *testing.T, funcName string) int {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "benchmark_test.go", nil, 0)
	if err != nil {
		t.Fatalf("parse benchmark_test.go: %v", err)
	}

	fn := benchmarkTestFuncDecl(file, funcName)
	if fn == nil || fn.Body == nil {
		t.Fatalf("find function %s in benchmark_test.go", funcName)
	}

	count := 0
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.IfStmt:
			count++
		case *ast.SwitchStmt:
			count += len(s.Body.List)
		case *ast.TypeSwitchStmt:
			count += len(s.Body.List)
		case *ast.BinaryExpr:
			if s.Op == token.LAND || s.Op == token.LOR {
				count++
			}
		}
		return true
	})
	return count
}

func benchmarkTestFuncDecl(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}
