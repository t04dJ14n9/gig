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
		// Issue 1  → Resolved Issue 28 (sort named-type conversion)
		// Issue 2  → fixed in gentool (time.Duration DirectCall)
		// Issue 3  → Resolved Issue 29 (fmt.Stringer)
		// Issue 4  → Resolved Issue 30 (fmt.Sprintf %T)
		// Issue 5  → Resolved Issue 31 (fmt.Sprintf %v _gig_id)
		// Issue 6  → Resolved Issue 32 (int64/uint64 narrowing)
		// Issue 7  → Resolved Issue 33 (bytes.Buffer.Cap)
		// Issue 8  → Resolved Issue 34 (json.Encoder method dispatch collision)

		// Strange Syntax Bugs (found 2026-03-28)
		"Bug_ConvertNilToInterface": {
			funcName: "StrangeSyntax_Bug1_ConvertNilToInterface",
			native:   func() any { var s []int; return s },
			issue:    "Nil slice to interface loses type information",
			panics:   false,
		},
		"Bug_NilMapAccess": {
			funcName: "StrangeSyntax_Bug2_NilMapAccess",
			native:   func() any { var m map[string]int; return m["key"] },
			issue:    "Nil map access returns nil instead of zero value",
			panics:   false,
		},
		"Bug_NilMapDelete": {
			funcName: "StrangeSyntax_Bug3_NilMapDelete",
			native:   func() any { var m map[string]int; delete(m, "key"); return 0 },
			issue:    "Delete on nil map causes panic instead of no-op",
			panics:   true,
		},
		"Bug_BlankExpression": {
			funcName: "StrangeSyntax_Bug4_BlankExpression",
			native:   func() any { _ = 42; var s []any; return s },
			issue:    "Blank identifier expression loses type in interface return",
			panics:   false,
		},
		"Bug_ChannelClosedSend": {
			funcName: "StrangeSyntax_Bug5_ChannelClosedSend",
			native: func() any {
				// After panic+recovery, function returns zero value of int = 0
				return 0
			},
			issue:  "Send on closed channel panic handling — returns zero value of int after recovery",
			panics: false,
		},
		"Bug_ClosureReturnNil": {
			funcName: "StrangeSyntax_Bug6_ClosureReturnNil",
			native:   func() any { var f func() int; return f },
			issue:    "Nil function return loses type information",
			panics:   false,
		},
	}

	if len(issues) == 0 {
		t.Skip("No known issues — all resolved!")
	}

	prog, err := gig.Build(knownIssuesSrc, gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	for name, tc := range issues {
		runKnownIssueTest(t, prog, name, tc)
	}
}

// TestKnownIssues_PanicRecover tests panic/recover/defer bugs.
// These tests document that defer closures cannot properly modify captured variables
// during panic recovery. Root cause: value.SetElem() fails in defer closure context.
func TestKnownIssues_PanicRecover(t *testing.T) {
	issues := map[string]KnownIssue{
		// IMPORTANT: In Go, only NAMED return values can be modified by defers
		// during panic recovery. Local variables modified by defers do NOT affect
		// the return value — the function returns the zero value of the return type.
		//
		// PanicRecoverBasic: returns 0 (zero value of int) — PASS
		// PanicRecoverWithValue: returns 0 (result is a local var, not named return) — PASS
		// DeferRunsOnPanic: returns 0 (result is a local var) — PASS
		// MultipleDefersOnPanic: returns 0 (result/order are local vars) — PASS
		// NamedReturnPanicRecover: returns 42 (result IS a named return) — PASS
		// NestedRecover: returns 0 (result is a local var) — needs defer-after-recover fix
		// PanicInDefer: returns 1 (result=1 set before panic, not a named return) — needs fix
		// PanicInClosure: returns 0 (zero value) — PASS
		// DeferPanicRecoverChain: returns 0 (result is a local var) — needs fix
	}

	if len(issues) == 0 {
		t.Skip("No known issues — all resolved!")
	}

	prog, err := gig.Build(knownIssuesSrc, gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	for name, tc := range issues {
		runKnownIssueTest(t, prog, name, tc)
	}
}
