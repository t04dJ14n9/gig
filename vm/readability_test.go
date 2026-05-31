package vm

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"
)

func TestRunLoopFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "run.go", 900, "move cold runtime paths to focused files")
}

func TestVMFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "vm.go", 360, "move execution entry and argument preparation to focused files")
}

func TestInterfaceAdapterFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "interface_adapter.go", 180, "move adapter call mechanics to focused files")
}

func TestInterfaceBoundaryFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "interface_boundary.go", 170, "move proxy lookup and host-interface classification to focused files")
}

func TestCallBoundaryFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "call_boundary.go", 140, "move reflect scanning and callable policy to focused files")
}

func TestKindMatchesTypeStaysShallow(t *testing.T) {
	count := directSwitchCaseCount(t, "ops_dispatch.go", "kindMatchesType")
	if count > 8 {
		t.Fatalf("kindMatchesType has %d direct switch cases, want <= 8; split primitive and composite type matching", count)
	}
}

func TestExecuteTypeConvertStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "ops_type_convert.go", "executeTypeConvert")
	if count > 12 {
		t.Fatalf("executeTypeConvert has %d branch points, want <= 12; split conversion domains into named helpers", count)
	}
}

func assertFileLineLimit(t *testing.T, path string, maxLines int, hint string) {
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

func directSwitchCaseCount(t *testing.T, path, funcName string) int {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	fn := findFuncDecl(file, funcName)
	if fn == nil || fn.Body == nil {
		t.Fatalf("find function %s in %s", funcName, path)
	}

	count := 0
	for _, stmt := range fn.Body.List {
		switch s := stmt.(type) {
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

	fn := findFuncDecl(file, funcName)
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
		}
		return true
	})
	return count
}

func findFuncDecl(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}
