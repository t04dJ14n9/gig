package value

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestValueEqualStaysShallow(t *testing.T) {
	count := directBranchCount(t, "comparison.go", "Equal")
	if count > 8 {
		t.Fatalf("Value.Equal has %d direct branches, want <= 8; move equality domains to named helpers", count)
	}
}

func TestIsGigStructStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "extern_type.go", "isGigStruct")
	if count > 8 {
		t.Fatalf("isGigStruct has %d branch points, want <= 8; keep pointer unwrapping and gig-name extraction separate", count)
	}
}

func TestSprintfExternStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "extern_sprintf.go", "SprintfExtern")
	if count > 8 {
		t.Fatalf("SprintfExtern has %d branch points, want <= 8; split format scanning from argument formatting", count)
	}
}

func TestGigErrorsIsStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "extern_errors_is.go", "GigErrorsIs")
	if count > 10 {
		t.Fatalf("GigErrorsIs has %d branch points, want <= 10; split match, custom Is, and unwrap traversal", count)
	}
}

func directBranchCount(t *testing.T, path, funcName string) int {
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
	for _, stmt := range fn.Body.List {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			count++
		case *ast.SwitchStmt:
			count += len(s.Body.List)
		case *ast.TypeSwitchStmt:
			count += len(s.Body.List)
		}
	}
	return count
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
