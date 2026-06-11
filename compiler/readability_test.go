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

func TestLookupExternalFuncInfoStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "compile_ext.go", "lookupExternalFuncInfo")
	if count > 3 {
		t.Fatalf(
			"lookupExternalFuncInfo has %d branch points, want <= 3; split lookup, descriptor construction, and reflect metadata",
			count,
		)
	}
}

func TestCompileExternalFuncValueStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "compile_ext.go", "compileExternalFuncValue")
	if count > 2 {
		t.Fatalf(
			"compileExternalFuncValue has %d branch points, want <= 2; split external value lookup from fallback and constant emission",
			count,
		)
	}
}

func TestLookupExternalMethodInfoStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "compile_ext.go", "lookupExternalMethodInfo")
	if count > 3 {
		t.Fatalf(
			"lookupExternalMethodInfo has %d branch points, want <= 3; split descriptor construction from direct-call lookup",
			count,
		)
	}
}

func TestCompileExternalStaticCallStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "compile_ext.go", "compileExternalStaticCall")
	if count > 4 {
		t.Fatalf(
			"compileExternalStaticCall has %d branch points, want <= 4; split argument emission, method dispatch, unresolved init stubs, and call emission",
			count,
		)
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

func TestCompileMakeSliceLoweringFileStaysFocused(t *testing.T) {
	assertCompilerFileLineLimit(t, "compile_make_slice_lowering.go", 160, "keep synthetic make-slice guards close to the lowering")
}

func TestSyntheticMakeSliceArrayAllocStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "compile_make_slice_lowering.go", "syntheticMakeSliceArrayAlloc")
	if count > 4 {
		t.Fatalf(
			"syntheticMakeSliceArrayAlloc has %d branch points, want <= 4; split referrer scanning from allocation shape checks",
			count,
		)
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

func TestIsUserDefinedNamedTypeStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "typecheck.go", "isUserDefinedNamedType")
	if count > 3 {
		t.Fatalf(
			"isUserDefinedNamedType has %d branch points, want <= 3; split pointer unwrapping, named-type extraction, and package classification",
			count,
		)
	}
}

func TestContainsUserDefinedUnderlyingStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "typecheck.go", "containsUserDefinedUnderlying")
	if count > 5 {
		t.Fatalf(
			"containsUserDefinedUnderlying has %d branch points, want <= 5; split element containers, aggregate containers, and callable/interface types",
			count,
		)
	}
}

func TestContainsUserDefinedSignatureDelegatesTupleScans(t *testing.T) {
	count := directCallCount(t, "typecheck.go", "containsUserDefinedSignature", "containsUserDefinedTuple")
	if count > 0 {
		t.Fatalf(
			"containsUserDefinedSignature calls containsUserDefinedTuple %d times; split receiver, parameter, and result scans",
			count,
		)
	}
}

func TestContainsUserDefinedInterfaceStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "typecheck.go", "containsUserDefinedInterface")
	if count > 2 {
		t.Fatalf(
			"containsUserDefinedInterface has %d branch points, want <= 2; split method and embedded interface scans",
			count,
		)
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

func directCallCount(t *testing.T, path, funcName string, callee string) int {
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
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		ident, ok := call.Fun.(*ast.Ident)
		if ok && ident.Name == callee {
			count++
		}
		return true
	})
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
