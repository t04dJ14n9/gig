package compiler

import (
	"bytes"
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

func TestCompileConstFileStaysFocused(t *testing.T) {
	assertCompilerFileLineLimit(t, "compile_const.go", 120, "move constant conversion helpers to focused files")
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
