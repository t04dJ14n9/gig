package gentool

import (
	"bytes"
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
