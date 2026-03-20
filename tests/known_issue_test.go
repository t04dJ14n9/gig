package tests

// Package tests - known_issue_test.go
//
// Tests for known interpreter bugs. Each test compares interpreted execution
// with native Go execution. Tests that PANIC or produce wrong results are
// expected failures — they document bugs awaiting fixes.

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
	issue    string        // bug description
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
	prog, err := gig.Build(knownIssuesSrc)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	issues := map[string]KnownIssue{
		// Bug 1: Named-type conversion to external sort types panics
		"SortIntSlice": {
			funcName: "SortIntSlice",
			native:   func() any { return ki.SortIntSlice() },
			issue:    "sort.IntSlice([]int) conversion not recognized — VM keeps []int, sort.Sort panics",
			panics:   true,
		},
		"SortFloat64Slice": {
			funcName: "SortFloat64Slice",
			native:   func() any { return ki.SortFloat64Slice() },
			issue:    "sort.Float64Slice([]float64) conversion not recognized — VM keeps []float64",
			panics:   true,
		},
		"SortStringSlice": {
			funcName: "SortStringSlice",
			native:   func() any { return ki.SortStringSlice() },
			issue:    "sort.StringSlice([]string) conversion not recognized — VM keeps []string",
			panics:   true,
		},
		"SortReverse": {
			funcName: "SortReverse",
			native:   func() any { return ki.SortReverse() },
			issue:    "sort.Reverse(sort.IntSlice(s)) panics — same root cause as SortIntSlice",
			panics:   true,
		},
		"SortIntsInPlace": {
			funcName: "SortIntsInPlace",
			native:   func() any { return ki.SortIntsInPlace() },
			issue:    "sort.Ints(s) mutates []int copy from []int64 conversion — original unchanged",
		},

		// Bug 2: FIXED — time.Duration DirectCall wrappers now use cast instead of type assertion
		// (context.WithTimeout, etc. now work correctly)

		// Bug 3: fmt.Stringer interface not honored on interpreted types
		"FmtStringerNotCalled": {
			funcName: "FmtStringerNotCalled",
			native:   func() any { return ki.FmtStringerNotCalled() },
			issue:    "fmt.Sprintf(\"%v\") ignores String() method on interpreted struct, prints raw fields",
		},

		// Bug 4: %T reports synthesized struct type instead of declared name
		"FmtSprintfTypeWrong": {
			funcName: "FmtSprintfTypeWrong",
			native:   func() any { return ki.FmtSprintfTypeWrong() },
			issue:    "fmt.Sprintf(\"%T\") reports synthesized struct type with _gig_id, not declared name",
		},

		// Bug 5: Extra _gig_id field in %v output
		"FmtSprintfExtraField": {
			funcName: "FmtSprintfExtraField",
			native:   func() any { return ki.FmtSprintfExtraField() },
			issue:    "fmt.Sprintf(\"%v\") includes extra _gig_id sentinel field: \"{1 2 {}}\" vs \"{1 2}\"",
		},

		// Bug 6: prog.Run() narrows int64→int and uint64→uint
		"StrconvParseIntNarrowed": {
			funcName: "StrconvParseIntNarrowed",
			native:   func() any { return ki.StrconvParseIntNarrowed() },
			issue:    "prog.Run() returns int instead of int64 for declared int64 return type",
		},
		"StrconvParseUintNarrowed": {
			funcName: "StrconvParseUintNarrowed",
			native:   func() any { return ki.StrconvParseUintNarrowed() },
			issue:    "prog.Run() returns uint instead of uint64 for declared uint64 return type",
		},
	}

	if len(issues) == 0 {
		t.Skip("No known issues — all resolved!")
	}

	for name, tc := range issues {
		runKnownIssueTest(t, prog, name, tc)
	}
}
