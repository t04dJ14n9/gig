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

func TestValueCmpStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "comparison.go", "Cmp")
	if count > 8 {
		t.Fatalf("Value.Cmp has %d branch points, want <= 8; split kind routing from primitive comparison mechanics", count)
	}
}

func TestMakeFromReflectStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "reflect.go", "MakeFromReflect")
	if count > 8 {
		t.Fatalf("MakeFromReflect has %d branch points, want <= 8; split validity, scalar, native slice, and reflect fallback domains", count)
	}
}

func TestFromInterfaceStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "reflect.go", "FromInterface")
	if count > 8 {
		t.Fatalf("FromInterface has %d branch points, want <= 8; split nil handling, fast interface cases, and reflect fallback domains", count)
	}
}

func TestReflectPrimitiveValueStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "reflect.go", "reflectPrimitiveValue")
	if count > 8 {
		t.Fatalf("reflectPrimitiveValue has %d branch points, want <= 8; split reflected scalar conversion by kind family", count)
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

func TestGigErrorsAsMatchValueStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "extern_errors_as.go", "gigAsMatchValue")
	if count > 10 {
		t.Fatalf("gigAsMatchValue has %d branch points, want <= 10; split native and gig-wrapper target matching", count)
	}
}

func TestValueInterfaceStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "accessor.go", "Interface")
	if count > 14 {
		t.Fatalf("Value.Interface has %d branch points, want <= 14; split primitive widths, interface unwrap, and reflected payloads", count)
	}
}

func TestSetElemStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "container_elem.go", "SetElem")
	if count > 10 {
		t.Fatalf("SetElem has %d branch points, want <= 10; split direct pointers, reflect pointers, and reflect containers", count)
	}
}

func TestExternWrapperFormatStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "extern_wrapper_format.go", "Format")
	if count > 10 {
		t.Fatalf("gigStructWrapper.Format has %d branch points, want <= 10; split verb dispatch and struct rendering", count)
	}
}

func TestGoStringValueStaysShallow(t *testing.T) {
	count := recursiveBranchCount(t, "extern_wrapper_gostring.go", "goStringValue")
	if count > 8 {
		t.Fatalf("goStringValue has %d branch points, want <= 8; split gig struct and gig slice rendering domains", count)
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
