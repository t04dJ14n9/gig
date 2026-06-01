package tests

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestGetBenchmarksSrcStaysShallow(t *testing.T) {
	count := testFileDecisionCount(t, "benchmark_test.go", "getBenchmarksSrc")
	if count > 8 {
		t.Fatalf("getBenchmarksSrc has %d decision points, want <= 8; parse benchmark files with focused helpers", count)
	}
}

func TestBenchmarkSummaryStaysShallow(t *testing.T) {
	count := testFileDecisionCount(t, "benchmark_test.go", "TestBenchmarkSummary")
	if count > 8 {
		t.Fatalf("TestBenchmarkSummary has %d decision points, want <= 8; split report sections into focused helpers", count)
	}
}

func TestStressMemoryLeakRunnerStaysShallow(t *testing.T) {
	count := testFileDecisionCount(t, "stress_leak_test.go", "runStressMemoryLeak")
	if count > 8 {
		t.Fatalf("runStressMemoryLeak has %d decision points, want <= 8; split setup, monitoring, shutdown, and reporting", count)
	}
}

func testFileDecisionCount(t *testing.T, path, funcName string) int {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	fn := testFileFuncDecl(file, funcName)
	if fn == nil || fn.Body == nil {
		t.Fatalf("find function %s in %s", funcName, path)
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

func testFileFuncDecl(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}
