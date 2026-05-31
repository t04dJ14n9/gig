package importer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestConvertReflectTypeStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "typeconv.go", "convertReflectType")
	if count > 8 {
		t.Fatalf("convertReflectType has %d branch points, want <= 8; keep cache, named-type, and kind conversion domains separate", count)
	}
}

func TestConvertReflectTypeForUnderlyingStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "typeconv.go", "convertReflectTypeForUnderlying")
	if count > 6 {
		t.Fatalf("convertReflectTypeForUnderlying has %d branch points, want <= 6; share composite kind conversion with the public converter", count)
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
