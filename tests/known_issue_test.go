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
	"reflect"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
	"git.woa.com/youngjin/gig/tests/testdata/known_issues"
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

// TestKnownIssues runs all known interpreter bugs.
// Every sub-test here is EXPECTED TO FAIL — they document real bugs.
// When a bug is fixed, remove it from here and promote to a passing test.
func TestKnownIssues(t *testing.T) {
	issues := map[string]KnownIssue{
		"SortIntSliceCompositeLiteral": {
			funcName: "SortIntSliceCompositeLiteral",
			native:   func() any { return known_issues.SortIntSliceCompositeLiteral() },
			issue:    "sort.IntSlice{...} composite literal lacks sort.Interface methods — sort.Sort panics with 'missing method Len'",
			panics:   true,
		},
		"TypedNilInterface": {
			funcName: "TypedNilInterface",
			native:   func() any { return known_issues.TypedNilInterface() },
			issue:    "typed nil pointer assigned to interface is incorrectly treated as nil — interface should be non-nil",
		},
		"InterfaceMethodOnTypedNil": {
			funcName: "InterfaceMethodOnTypedNil",
			native:   func() any { return known_issues.InterfaceMethodOnTypedNil() },
			issue:    "calling method on typed nil interface should panic with nil pointer dereference, but Gig returns 'ok'",
		},
		"CToF": {
			funcName: "CToF",
			native:   func() any { return known_issues.CToF() },
			issue:    "named type arithmetic returns float64 instead of named type Fahrenheit",
		},
		"AssignTypeAssertion": {
			funcName: "AssignTypeAssertion",
			native:   func() any { return known_issues.AssignTypeAssertion() },
			issue:    "failed type assertion comma-ok returns nil instead of zero value of target type",
		},
		"SliceNilSubslice": {
			funcName: "SliceNilSubslice",
			native:   func() any { return known_issues.SliceNilSubslice() },
			issue:    "nil slice [0:0] returns empty non-nil slice instead of nil slice",
		},
		"LinkedListReverse": {
			funcName: "LinkedListReverse",
			native:   func() any { return known_issues.LinkedListReverse() },
			issue:    "linked list reverse with pointer reassignment fails — struct field modifications through pointer traversal not propagated correctly",
		},
		"GlobalSliceAccess": {
			funcName: "GlobalSliceAccess",
			native:   func() any { return known_issues.GlobalSliceAccess() },
			issue:    "package-level var with slice initializer not properly initialized",
		},
		"GlobalMapAccess": {
			funcName: "GlobalMapAccess",
			native:   func() any { return known_issues.GlobalMapAccess() },
			issue:    "package-level var with map initializer not properly initialized",
		},
		"GlobalStringAccess": {
			funcName: "GlobalStringAccess",
			native:   func() any { return known_issues.GlobalStringAccess() },
			issue:    "package-level var with string initializer not properly initialized",
		},
		"GlobalBoolAccess": {
			funcName: "GlobalBoolAccess",
			native:   func() any { return known_issues.GlobalBoolAccess() },
			issue:    "package-level var with bool initializer not properly initialized",
		},
		"GlobalFloatAccess": {
			funcName: "GlobalFloatAccess",
			native:   func() any { return known_issues.GlobalFloatAccess() },
			issue:    "package-level var with float initializer not properly initialized",
		},
		"GlobalPointerNil": {
			funcName: "GlobalPointerNil",
			native:   func() any { return known_issues.GlobalPointerNil() },
			issue:    "package-level var with nil pointer initializer not properly initialized",
		},
	}

	prog, err := gig.Build(knownIssuesSrc, gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	for name, tc := range issues {
		runKnownIssueTest(t, prog, name, tc)
	}
}
