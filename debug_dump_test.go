package gig

import (
	"strings"
	"testing"
)

// TestDebugDumpIncludesSSA exercises the v2 DebugDump output. The
// legacy version dumped both SSA and bytecode; the v2 SSA pipeline has
// no bytecode, so the dump now reports SSA only.
func TestDebugDumpIncludesSSA(t *testing.T) {
	src := `package main

func Add(a, b int) int {
	return a + b
}

func Answer() int {
	return 42
}
`
	dump, err := DebugDump(src)
	if err != nil {
		t.Fatalf("DebugDump error: %v", err)
	}
	wantParts := []string{
		"# Package: main",
		"# Member Add",
		"func Add",
		"# Member Answer",
	}
	for _, part := range wantParts {
		if !strings.Contains(dump, part) {
			t.Fatalf("DebugDump missing %q\n%s", part, dump)
		}
	}
	addPos := strings.Index(dump, "# Member Add")
	answerPos := strings.Index(dump, "# Member Answer")
	if addPos < 0 || answerPos < 0 {
		t.Fatalf("missing member entries in dump:\n%s", dump)
	}
	if addPos > answerPos {
		t.Fatalf("members are not sorted by name:\n%s", dump)
	}
}

// TestDebugDumpDoesNotExecuteInit verifies init() is not invoked
// during DebugDump (only SSA build runs).
func TestDebugDumpDoesNotExecuteInit(t *testing.T) {
	src := `package main

func init() {
	panic("debug dump should not execute init")
}

func F() int {
	return 1
}
`
	dump, err := DebugDump(src, WithAllowPanic())
	if err != nil {
		t.Fatalf("DebugDump error: %v", err)
	}
	if !strings.Contains(dump, "# Member init") {
		t.Fatalf("DebugDump missing init member:\n%s", dump)
	}
	if !strings.Contains(dump, "# Member F") {
		t.Fatalf("DebugDump missing F member:\n%s", dump)
	}
}
