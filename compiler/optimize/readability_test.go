package optimize

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"
)

func TestFuseSliceOpsStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "slice_fusion.go", "FuseSliceOps")
	if count > 10 {
		t.Fatalf("FuseSliceOps has %d branch points, want <= 10; split scan, candidate decoding, and fusion construction domains", count)
	}
}

func TestIntSpecializeStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "int_specialize.go", "IntSpecialize")
	if count > 10 {
		t.Fatalf("IntSpecialize has %d branch points, want <= 10; split rule-table setup, int-use discovery, and rewrite passes", count)
	}
}

func TestConstantFoldingFileStaysFocused(t *testing.T) {
	assertOptimizerFileLineLimit(t, "constant_folding.go", 160, "split propagation, rewrite collection, and arithmetic semantics into focused files")
}

func assertOptimizerFileLineLimit(t *testing.T, path string, maxLines int, hint string) {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	lines := bytes.Count(src, []byte{'\n'})
	if lines > maxLines {
		t.Fatalf("%s has %d lines, want <= %d; %s", path, lines, maxLines, hint)
	}
}

func recursiveBranchCount(t *testing.T, path, funcName string) int {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	fn := findReadabilityFunc(file, funcName)
	if fn == nil || fn.Body == nil {
		t.Fatalf("find function %s in %s", funcName, path)
	}

	count := 0
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.IfStmt:
			count++
		case *ast.ForStmt:
			count++
		case *ast.RangeStmt:
			count++
		case *ast.SwitchStmt:
			count += len(s.Body.List)
		case *ast.TypeSwitchStmt:
			count += len(s.Body.List)
		}
		return true
	})
	return count
}

func findReadabilityFunc(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}
