package tests

// Package tests - known_issue_test.go
//
// Tests for known interpreter bugs. Each test compares interpreted execution
// with native Go execution. Tests that PANIC or produce wrong results are
// expected failures — they document bugs awaiting fixes.
//
// When a bug is fixed, move its function to testdata/resolved_issue/main.go
// and register it in correctness_test.go's resolved_issueTests map.

import (
	_ "embed"
	"reflect"
	"testing"

	"git.woa.com/youngjin/gig"
	ki "git.woa.com/youngjin/gig/tests/testdata/known_issues"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
)

//go:embed testdata/known_issues/main.go
var knownIssuesSrc string

// KnownIssue represents a test case for a known bug.
type KnownIssue struct {
	funcName string        // function name in embedded source
	native   func() any    // native Go function for comparison
	issue    string       // bug description
	panics   bool          // true if interpreter panics (vs. wrong result)
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

		if !reflect.DeepEqual(interpResult, nativeResult) {
			t.Errorf("BUG (mismatch): %s\n  interpreter: %v (%T)\n  native:      %v (%T)",
				tc.issue, interpResult, interpResult, nativeResult, nativeResult)
		}
	})
}

// TestKnownIssues_Tricky runs all known interpreter bugs.
// Every sub-test here is EXPECTED TO FAIL — they document real bugs.
func TestKnownIssues_Tricky(t *testing.T) {
	issues := map[string]KnownIssue{
		// Bug 8: Method dispatch type collision — json.Encoder.Encode vs xml.Encoder.Encode
		"JsonEncodeBug8": {
			funcName: "JsonEncodeBug8",
			native:   func() any { return ki.JsonEncodeBug8() },
			issue:    "Method dispatch picks xml.Encoder.Encode instead of json.Encoder.Encode — type collision in compiled method cache",
			panics:   true,
		},
	}

	if len(issues) == 0 {
		t.Skip("No known issues — all resolved!")
	}

	prog, err := gig.Build(knownIssuesSrc)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	for name, tc := range issues {
		runKnownIssueTest(t, prog, name, tc)
	}
}