package gentool

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"
)

func TestDirectCallFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "directcall.go", 180, "move direct-call helpers to focused files")
}

func TestResolveFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "resolve.go", 180, "move type-resolution helpers to focused files")
}

func TestExtractFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "extract.go", 180, "move argument extraction helpers to focused files")
}

func TestExtractBasicLinesStayReadable(t *testing.T) {
	assertMaxLineLength(t, "extract_basic.go", 160)
}

func TestGeneratorFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "generator.go", 180, "move package generation helpers to focused files")
}

func TestGenerateSingleMethodDirectCallStaysShallow(t *testing.T) {
	count := cyclomaticBranchCount(t, "directcall_method.go", "generateSingleMethodDirectCall")
	if count > 20 {
		t.Fatalf("generateSingleMethodDirectCall has complexity %d, want <= 20; split eligibility, receiver extraction, arguments, variadics, and result emission", count)
	}
}

func TestGenerateDirectCallStaysShallow(t *testing.T) {
	count := cyclomaticBranchCount(t, "directcall.go", "generateDirectCall")
	if count > 18 {
		t.Fatalf("generateDirectCall has complexity %d, want <= 18; split eligibility, argument emission, variadics, call selection, and result emission", count)
	}
}

func TestWrapBasicReturnStaysShallow(t *testing.T) {
	count := cyclomaticBranchCount(t, "wrap.go", "wrapBasicReturn")
	if count > 18 {
		t.Fatalf("wrapBasicReturn has complexity %d, want <= 18; split unsigned, signed, float, complex, and string wrapping", count)
	}
}

func TestCollectTypeImportsStaysShallow(t *testing.T) {
	count := cyclomaticBranchCount(t, "resolve_imports.go", "collectTypeImports")
	if count > 12 {
		t.Fatalf("collectTypeImports has complexity %d, want <= 12; split object, element, tuple, and interface import walking", count)
	}
}

func TestGeneratorEmitAvoidsNestedFormattedWriteString(t *testing.T) {
	assertNoNestedFormattedWriteString(t, "generator_emit.go")
}

func TestGeneratorBarrelAvoidsNestedFormattedWriteString(t *testing.T) {
	assertNoNestedFormattedWriteString(t, "generator_barrel.go")
}

func TestInterfaceProxyAvoidsNestedFormattedWriteString(t *testing.T) {
	assertNoNestedFormattedWriteString(t, "interface_proxy.go")
}

func assertNoNestedFormattedWriteString(t *testing.T, path string) {
	t.Helper()

	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if count := bytes.Count(src, []byte("WriteString(fmt.Sprintf(")); count > 0 {
		t.Fatalf("%s has %d WriteString(fmt.Sprintf(...)) calls; use fmt.Fprintf on the builder instead", path, count)
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

func assertMaxLineLength(t *testing.T, path string, maxColumns int) {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	for i, line := range bytes.Split(src, []byte{'\n'}) {
		if len(line) > maxColumns {
			t.Fatalf("%s:%d has %d columns, want <= %d", path, i+1, len(line), maxColumns)
		}
	}
}

func cyclomaticBranchCount(t *testing.T, path, funcName string) int {
	t.Helper()

	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	fn := findReadabilityFunc(file, funcName)
	if fn == nil || fn.Body == nil {
		t.Fatalf("find function %s in %s", funcName, path)
	}

	count := 1
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
			count++
		case *ast.CaseClause:
			count++
		case *ast.BinaryExpr:
			if x.Op.String() == "&&" || x.Op.String() == "||" {
				count++
			}
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
