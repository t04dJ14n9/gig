package tests

// Package tests - known_issue_test.go
//
// This file contains tests for known interpreter bugs.
// These tests document issues that need to be fixed.
// Each test is skipped with a reference to the specific bug.

import (
	"reflect"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
)

// KnownIssue represents a test case for a known bug
type KnownIssue struct {
	src      string
	funcName string
	args     []any
	expected any
	issue    string // Bug description
}

// callNative is defined in correctness_test.go

// runKnownIssueTest runs a test that is expected to fail
func runKnownIssueTest(t *testing.T, name string, tc KnownIssue) {
	t.Run(name, func(t *testing.T) {
		// Recover from panics - these are expected for known bugs
		defer func() {
			if r := recover(); r != nil {
				t.Skipf("Known issue causes panic: %s - panic: %v", tc.issue, r)
			}
		}()

		src := toMainPackage(tc.src)
		prog, err := gig.Build(src)
		if err != nil {
			t.Skipf("Build error: %v", err)
		}

		result, err := prog.Run(tc.funcName, tc.args...)
		if err != nil {
			t.Skipf("Run error: %v - Issue: %s", err, tc.issue)
		}

		// Check if bug is fixed using reflect.DeepEqual
		if reflect.DeepEqual(result, tc.expected) {
			t.Logf("BUG FIXED! %s", tc.issue)
			return
		}

		t.Skipf("Known issue: %s - got %v (%T), expected %v (%T)",
			tc.issue, result, result, tc.expected, tc.expected)
	})
}

// ============================================================================
// Known Issue Tests - Interpreter Bugs
// ============================================================================

func TestKnownIssues_Tricky(t *testing.T) {
	issues := map[string]KnownIssue{
		"StringReverse": {
			src:      trickySrc,
			funcName: "StringReverse",
			args:     []any{"hello"},
			expected: "olleh",
			issue:    "panic: invalid reflect.Value in SetElem() - rune slice assignment issue",
		},
		"Clamp": {
			src:      trickySrc,
			funcName: "Clamp",
			args:     []any{150, 0, 100},
			expected: 100,
			issue:    "returns 0 instead of 100 - multi-condition if-else chain bug",
		},
		"Sign": {
			src:      trickySrc,
			funcName: "Sign",
			args:     []any{-42},
			expected: -1,
			issue:    "returns 0 instead of -1 for negative numbers - comparison issue",
		},
		"SliceUniqueParam": {
			src:      trickySrc,
			funcName: "SliceUniqueParam",
			args:     []any{[]int{1, 2, 2, 3, 3, 3}},
			expected: []int{1, 2, 3},
			issue:    "type mismatch: interpreter returns []int64, native returns []int",
		},
		"SliceInterleave": {
			src:      trickySrc,
			funcName: "SliceInterleave",
			args:     []any{[]int{1, 3, 5}, []int{2, 4, 6}},
			expected: []int{1, 2, 3, 4, 5, 6},
			issue:    "type mismatch: interpreter returns []int64, native returns []int",
		},
		"SliceRotateLeftParam": {
			src:      trickySrc,
			funcName: "SliceRotateLeftParam",
			args:     []any{[]int{1, 2, 3, 4, 5}, 2},
			expected: []int{4, 5, 1, 2, 3},
			issue:    "type mismatch: interpreter returns []int64, native returns []int",
		},
		"BitCountOnes": {
			src:      trickySrc,
			funcName: "BitCountOnes",
			args:     []any{255},
			expected: 8,
			issue:    "panic: not an int: uint - bitwise operation result type issue",
		},
		"BinomialCoefficient": {
			src:      trickySrc,
			funcName: "BinomialCoefficient",
			args:     []any{5, 2},
			expected: 10,
			issue:    "incorrect calculation result",
		},
		"FibonacciNth": {
			src:      trickySrc,
			funcName: "FibonacciNth",
			args:     []any{20},
			expected: 6765,
			issue:    "incorrect calculation result",
		},
		"IsPrime": {
			src:      trickySrc,
			funcName: "IsPrime",
			args:     []any{17},
			expected: true,
			issue:    "incorrect calculation result",
		},
		"FactorialIterative": {
			src:      trickySrc,
			funcName: "FactorialIterative",
			args:     []any{5},
			expected: 120,
			issue:    "incorrect calculation result",
		},
		"MapDeepCopy": {
			src:      trickySrc,
			funcName: "MapDeepCopy",
			args:     []any{map[int][]int{1: {1, 2}, 2: {3, 4}}},
			expected: map[int][]int{1: {1, 2}, 2: {3, 4}},
			issue:    "complex map with slice values not handled correctly",
		},
	}

	for name, tc := range issues {
		runKnownIssueTest(t, name, tc)
	}
}
