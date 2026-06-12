package gig

import (
	"strings"
	"testing"
)

func TestDebugDumpIncludesReadableSSAAndBytecode(t *testing.T) {
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
		"# Gig Debug Dump",
		"## SSA",
		"func Add",
		"func Answer",
		"## Bytecode",
		"### Function index",
		"### Function Add",
		"locals:",
		"0000",
		"; local[0], local[1]",
		"RETURNVAL",
		"### Function Answer",
		"CONST",
		"; const[",
		"## Constants",
		"## Globals",
		"## Types",
	}
	for _, part := range wantParts {
		if !strings.Contains(dump, part) {
			t.Fatalf("DebugDump missing %q\n%s", part, dump)
		}
	}

	addPos := strings.Index(dump, "### Function Add")
	answerPos := strings.Index(dump, "### Function Answer")
	if addPos < 0 || answerPos < 0 {
		t.Fatalf("missing function sections in dump:\n%s", dump)
	}
	if addPos > answerPos {
		t.Fatalf("functions are not sorted by name:\n%s", dump)
	}
}

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
	if !strings.Contains(dump, "func init") {
		t.Fatalf("DebugDump missing init SSA:\n%s", dump)
	}
	if !strings.Contains(dump, "### Function F") {
		t.Fatalf("DebugDump missing function bytecode:\n%s", dump)
	}
}
