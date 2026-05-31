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

func TestNewVMStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "vm.go", "newVM")
	if count > 8 {
		t.Fatalf("newVM has %d branch points, want <= 8; split global initialization domains from VM allocation", count)
	}
}

func TestInterfaceAdapterFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "interface_adapter.go", 180, "move adapter call mechanics to focused files")
}

func TestInterfaceBoundaryFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "interface_boundary.go", 170, "move proxy lookup and host-interface classification to focused files")
}

func TestInterpretedTypeSatisfiesInterfaceStaysShallow(t *testing.T) {
	count := recursiveDecisionCount(t, "interface_boundary.go", "interpretedTypeSatisfiesInterface")
	if count > 8 {
		t.Fatalf("interpretedTypeSatisfiesInterface has %d decision points, want <= 8; split method lookup, receiver eligibility, and dynamic receiver matching", count)
	}
}

func TestCallBoundaryFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "call_boundary.go", 140, "move reflect scanning and callable policy to focused files")
}

func TestReflectTypeContainsInterfaceStaysShallow(t *testing.T) {
	count := recursiveDecisionCount(t, "call_boundary_func.go", "reflectTypeContainsInterface")
	if count > 8 {
		t.Fatalf("reflectTypeContainsInterface has %d decision points, want <= 8; split recursion guard, kind routing, and composite scans", count)
	}
}

func TestBuildReflectArgsStaysShallow(t *testing.T) {
	count := recursiveDecisionCount(t, "call_external.go", "buildReflectArgs")
	if count > 8 {
		t.Fatalf("buildReflectArgs has %d decision points, want <= 8; split packed variadic, positional, and element conversion paths", count)
	}
}

func TestUnpackVariadicArgsStaysShallow(t *testing.T) {
	count := recursiveDecisionCount(t, "call_runtime.go", "unpackVariadicArgs")
	if count > 8 {
		t.Fatalf("unpackVariadicArgs has %d decision points, want <= 8; split value slice, int slice, byte slice, and reflect slice unpacking", count)
	}
}

func TestReflectMethodTypeForBoundaryStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "call_external_method.go", "reflectMethodTypeForBoundary")
	if count > 8 {
		t.Fatalf(
			"reflectMethodTypeForBoundary has %d branch points, want <= 8; reuse receiver normalization and method lookup helpers",
			count,
		)
	}
}

func TestKindMatchesTypeStaysShallow(t *testing.T) {
	count := directSwitchCaseCount(t, "ops_dispatch.go", "kindMatchesType")
	if count > 8 {
		t.Fatalf("kindMatchesType has %d direct switch cases, want <= 8; split primitive and composite type matching", count)
	}
}

func TestSameReflectKindFamilyStaysShallow(t *testing.T) {
	count := recursiveDecisionCount(t, "ops_dispatch.go", "sameReflectKindFamily")
	if count > 8 {
		t.Fatalf("sameReflectKindFamily has %d decision points, want <= 8; split reflect.Kind family predicates", count)
	}
}

func TestExecuteTypeConvertStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "ops_type_convert.go", "executeTypeConvert")
	if count > 12 {
		t.Fatalf("executeTypeConvert has %d branch points, want <= 12; split conversion domains into named helpers", count)
	}
}

func TestConvertBasicValueStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "ops_type_convert.go", "convertBasicValue")
	if count > 8 {
		t.Fatalf("convertBasicValue has %d branch points, want <= 8; split string, signed, unsigned, and float conversion domains", count)
	}
}

func TestTypeToReflectInnerStaysGrouped(t *testing.T) {
	count := directSwitchCaseCount(t, "typeconv.go", "typeToReflectInner")
	if count > 8 {
		t.Fatalf("typeToReflectInner has %d direct type-switch cases, want <= 8; split composite and named reflect-type domains", count)
	}
}

func TestGenericSuperinstructionStaysGrouped(t *testing.T) {
	count := directSwitchCaseCount(t, "run_super_generic.go", "runGenericSuperinstruction")
	if count > 8 {
		t.Fatalf("runGenericSuperinstruction has %d direct switch cases, want <= 8; route by superinstruction family", count)
	}
}

func TestIteratorNextStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "iterator.go", "next")
	if count > 10 {
		t.Fatalf("iterator.next has %d branch points, want <= 10; split native and reflected range domains", count)
	}
}

func TestExecuteSliceStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "ops_slice.go", "executeSlice")
	if count > 10 {
		t.Fatalf("executeSlice has %d branch points, want <= 10; split nil, native, and reflect slicing paths", count)
	}
}

func TestExecuteAssertStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "ops_assert.go", "executeAssert")
	if count > 10 {
		t.Fatalf("executeAssert has %d branch points, want <= 10; split assertion target, interface, reflect, and primitive paths", count)
	}
}

func TestExecuteMakeInterfaceStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "ops_make_interface.go", "executeMakeInterface")
	if count > 10 {
		t.Fatalf("executeMakeInterface has %d branch points, want <= 10; split adapter, interpreted, pass-through, and typed-nil paths", count)
	}
}

func TestDereferenceValueStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "reference.go", "dereferenceValue")
	if count > 8 {
		t.Fatalf("dereferenceValue has %d branch points, want <= 8; split reference, reflect, pointer, and nil fallback domains", count)
	}
}

func TestAppendValueStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "ops_append.go", "appendValue")
	if count > 10 {
		t.Fatalf("appendValue has %d branch points, want <= 10; split native int, byte, reflect, and nil append domains", count)
	}
}

func TestExecuteCopyStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "ops_copy_delete.go", "executeCopy")
	if count > 8 {
		t.Fatalf("executeCopy has %d branch points, want <= 8; split byte, native int, and reflect copy domains", count)
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

func recursiveDecisionCount(t *testing.T, path, funcName string) int {
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
		case *ast.BinaryExpr:
			if s.Op == token.LAND || s.Op == token.LOR {
				count++
			}
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
