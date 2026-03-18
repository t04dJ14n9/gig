package tests

// Package tests - known_issue_test.go
//
// This file contains tests for known interpreter bugs.
// These tests compare interpreted execution with native Go execution.
// Each test is skipped with a reference to the specific bug.

import (
	"reflect"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
	"git.woa.com/youngjin/gig/tests/testdata/tricky"
)

// KnownIssue represents a test case for a known bug
type KnownIssue struct {
	src      string
	funcName string
	args     []any
	native   any // native function reference - called via reflection
	issue    string
}

// runKnownIssueTest runs a test that compares interpreter vs native.
// This test FAILS to show the actual difference between interpreter and native.
func runKnownIssueTest(t *testing.T, name string, tc KnownIssue) {
	t.Run(name, func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("INTERPRETER PANIC: %v\n  Issue: %s", r, tc.issue)
			}
		}()

		src := toMainPackage(tc.src)
		prog, err := gig.Build(src)
		if err != nil {
			t.Fatalf("Build error: %v", err)
		}

		// Run interpreted code
		interpResult, err := prog.Run(tc.funcName, tc.args...)
		if err != nil {
			t.Errorf("INTERPRETER ERROR: %v\n  Issue: %s", err, tc.issue)
			return
		}

		// Get native result using reflection
		nativeResult := callNative(tc.native, tc.args)

		// Compare results - FAIL if different
		if !reflect.DeepEqual(interpResult, nativeResult) {
			t.Errorf("MISMATCH: %s\n  interpreter: %v (%T)\n  native:      %v (%T)",
				tc.issue, interpResult, interpResult, nativeResult, nativeResult)
		}
	})
}

// ============================================================================
// Known Issue Tests - Interpreter Bugs
// ============================================================================

func TestKnownIssues_Tricky(t *testing.T) {
	issues := map[string]KnownIssue{
		"StructEmbeddedInterface": {
			src:      trickySrc,
			funcName: "StructEmbeddedInterface",
			args:     nil,
			native:   tricky.StructEmbeddedInterface,
			issue:    "flaky: passes in isolation but fails when run with other tests - possible reflect.StructOf type collision",
		},
	}

	for name, tc := range issues {
		runKnownIssueTest(t, name, tc)
	}
}
