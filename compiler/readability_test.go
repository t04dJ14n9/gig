package compiler

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"
)

func TestCompilerFileStaysFocused(t *testing.T) {
	assertCompilerFileLineLimit(t, "compiler.go", 190, "move compile-stage helpers to focused files")
}

func TestCompileFuncFileStaysFocused(t *testing.T) {
	assertCompilerFileLineLimit(t, "compile_func.go", 160, "move function setup and block helpers to focused files")
}

func TestExternalFuncOriginFileStaysFocused(t *testing.T) {
	assertCompilerFileLineLimit(t, "external_func_origin.go", 160, "move origin tracing helpers to focused files")
}

func TestExternalFuncOriginsSeenStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "external_func_origin.go", "externalFuncOriginsSeen")
	if count > 8 {
		t.Fatalf("externalFuncOriginsSeen has %d branch points, want <= 8; split SSA routing and cycle guards", count)
	}
}

func TestExternalFuncOriginsFromSSAValueStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "external_func_origin.go", "externalFuncOriginsFromSSAValue")
	if count > 8 {
		t.Fatalf("externalFuncOriginsFromSSAValue has %d branch points, want <= 8; split alias, container, and producer routing domains", count)
	}
}

func TestCompileConstFileStaysFocused(t *testing.T) {
	assertCompilerFileLineLimit(t, "compile_const.go", 120, "move constant conversion helpers to focused files")
}

func TestCompileBuiltinCallStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "compile_builtin.go", "compileBuiltinCall")
	if count > 10 {
		t.Fatalf("compileBuiltinCall has %d branch points, want <= 10; split no-result, simple-result, and custom-result builtin domains", count)
	}
}

func TestPackedVarargsValuesStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "compile_builtin.go", "packedVarargsValues")
	if count > 10 {
		t.Fatalf("packedVarargsValues has %d branch points, want <= 10; split varargs shape validation, store collection, and completeness checks", count)
	}
}

func TestCompileInstructionDispatchStaysShallow(t *testing.T) {
	count := directTypeSwitchCaseCount(t, "compile_instr.go", "compileInstruction")
	if count > 8 {
		t.Fatalf("compileInstruction has %d direct type-switch cases, want <= 8; route through focused instruction families", count)
	}
}

func TestContainsUserDefinedTypeSeenStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "typecheck.go", "containsUserDefinedTypeSeen")
	if count > 10 {
		t.Fatalf("containsUserDefinedTypeSeen has %d branch points, want <= 10; split composite type traversal helpers", count)
	}
}

func TestValidateExternalCallBoundaryStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "typecheck.go", "validateExternalCallBoundary")
	if count > 3 {
		t.Fatalf(
			"validateExternalCallBoundary has %d branch points, want <= 3; split package trust, argument policy, and diagnostics",
			count,
		)
	}
}

func TestBasicConstValueStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "compile_const_basic.go", "basicConstValue")
	if count > 5 {
		t.Fatalf("basicConstValue has %d branch points, want <= 5; move kind-specific conversion into table-driven helpers", count)
	}
}

func TestAllOnesConstantStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "compile_operator.go", "allOnesConstant")
	if count > 5 {
		t.Fatalf("allOnesConstant has %d branch points, want <= 5; move per-width masks into a lookup table", count)
	}
}

func assertCompilerFileLineLimit(t *testing.T, path string, maxLines int, hint string) {
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

func directTypeSwitchCaseCount(t *testing.T, path, funcName string) int {
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
		// The readability guard intentionally measures only direct routing
		// in this function. Focused helpers may still use small type switches.
		if sw, ok := stmt.(*ast.TypeSwitchStmt); ok {
			count += len(sw.Body.List)
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

func findFuncDecl(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}
