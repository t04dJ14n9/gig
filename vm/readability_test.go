package vm

import (
	"bytes"
	"os"
	"testing"
)

func TestRunLoopFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "run.go", 900, "move cold runtime paths to focused files")
}

func TestVMFileStaysFocused(t *testing.T) {
	assertFileLineLimit(t, "vm.go", 360, "move execution entry and argument preparation to focused files")
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
