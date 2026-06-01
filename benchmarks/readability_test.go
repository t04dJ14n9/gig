package benchmarks

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestExtremeStressStaysShallow(t *testing.T) {
	count := benchmarkDecisionCount(t, "extreme_stress_test.go", "TestExtremeStress")
	if count > 8 {
		t.Fatalf("TestExtremeStress has %d decision points, want <= 8; split setup, runs, and reports", count)
	}
}

func TestConcurrentGlobalsComparisonStaysShallow(t *testing.T) {
	count := benchmarkDecisionCount(t, "extreme_stress_test.go", "TestConcurrentGlobals_GoNative_vs_Gig")
	if count > 8 {
		t.Fatalf("TestConcurrentGlobals_GoNative_vs_Gig has %d decision points, want <= 8; split native, gig, report, and assertion phases", count)
	}
}

func TestStatefulStressStaysShallow(t *testing.T) {
	count := benchmarkDecisionCount(t, "extreme_stress_test.go", "TestStatefulStress")
	if count > 8 {
		t.Fatalf("TestStatefulStress has %d decision points, want <= 8; split correctness, stress table, reporting, and counter verification", count)
	}
}

func benchmarkDecisionCount(t *testing.T, path, funcName string) int {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	fn := benchmarkFuncDecl(file, funcName)
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
		case *ast.ForStmt:
			count++
		case *ast.RangeStmt:
			count++
		case *ast.BinaryExpr:
			if s.Op == token.LAND || s.Op == token.LOR {
				count++
			}
		}
		return true
	})
	return count
}

func benchmarkFuncDecl(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}
