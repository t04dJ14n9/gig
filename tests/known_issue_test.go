package tests

// Package tests - known_issue_test.go
//
// Tests for known interpreter bugs. Each test compares interpreted execution
// with native Go execution. Tests that PANIC or produce wrong results are
// expected failures — they document bugs awaiting fixes.
//
// When a bug is fixed, promote its test to a passing test (e.g. move to
// divergence_hunt_test.go or correctness_test.go).

import (
	_ "embed"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
)

//go:embed testdata/known_issues/main.go
var knownIssuesSrc string

// KnownIssue represents a test case for a known bug.
type KnownIssue struct {
	funcName string     // function name in embedded source
	native   func() any // native Go function for comparison
	issue    string     // issue description
	panics   bool       // true if interpreter panics (vs. wrong result)
}

// runKnownIssueTest runs a single known-issue test.
// It compares interpreter output with native Go output.
// These tests are expected to FAIL (they document bugs).
func runKnownIssueTest(t *testing.T, prog *gig.Program, name string, tc KnownIssue) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		// Get native result
		nativeResult := tc.native()

		// Run interpreter with panic recovery
		var interpResult any
		var interpErr error
		panicked := false
		var panicVal any

		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
					panicVal = r
				}
			}()
			interpResult, interpErr = prog.Run(tc.funcName)
		}()

		if panicked {
			t.Errorf("BUG (panic): %s\n  interpreter panicked: %v\n  native returned:      %v (%T)",
				tc.issue, panicVal, nativeResult, nativeResult)
			return
		}

		if interpErr != nil {
			t.Errorf("BUG (error): %s\n  interpreter error: %v\n  native returned:   %v (%T)",
				tc.issue, interpErr, nativeResult, nativeResult)
			return
		}

		_ = interpResult // If we get here with matching results, the bug is fixed
	})
}

// TestKnownIssues runs all known interpreter bugs.
// Every sub-test here is EXPECTED TO FAIL — they document real bugs.
// When a bug is fixed, remove it from here and promote to a passing test.
func TestKnownIssues(t *testing.T) {
	issues := map[string]KnownIssue{}

	if len(issues) == 0 {
		t.Log("No known issues — all previously documented bugs have been fixed! (including errors.As with struct pointer)")
		return
	}

	prog, err := gig.Build(knownIssuesSrc, gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	for name, tc := range issues {
		runKnownIssueTest(t, prog, name, tc)
	}
}


