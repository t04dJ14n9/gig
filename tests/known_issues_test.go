// Package tests — known_issues_test.go
//
// Tests for known issues in the interpreter.
// These tests document bugs that are expected to fail.
//
// All previously known issues have been resolved and migrated to
// resolved_issue_test.go. This file is kept as a template for
// documenting future bugs.
package tests

import (
	_ "embed"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
)

//go:embed testdata/known_issues/main.go
var knownIssuesSrc string

// Known issue tests - these document bugs and are not expected to pass
var knownIssueTests = []struct {
	name string
}{
	// All known issues have been resolved!
}

// TestKnownIssues runs all known issue tests
// These tests are expected to fail or have issues - they document bugs
func TestKnownIssues(t *testing.T) {
	if len(knownIssueTests) == 0 {
		t.Skip("No known issues remaining — all resolved!")
	}

	for _, tc := range knownIssueTests {
		t.Run("known_issues/"+tc.name, func(t *testing.T) {
			// Try to run interpreted version with panic recovery
			interpResult, err := safeRunInterpreter(t, knownIssuesSrc, tc.name)

			if err != nil {
				t.Logf("Interpreter error (expected for known issue): %v", err)
				return
			}

			t.Logf("Interpreter result: %v", interpResult)
		})
	}
}

// safeRunInterpreter runs the interpreter with panic recovery
func safeRunInterpreter(_ *testing.T, src, funcName string) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = panicError{r}
		}
	}()

	prog, buildErr := gig.Build(src)
	if buildErr != nil {
		return nil, buildErr
	}

	return prog.Run(funcName)
}

type panicError struct {
	value any
}

func (p panicError) Error() string {
	return "panic: " + panicString(p.value)
}

func panicString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case error:
		return val.Error()
	default:
		return "unknown panic"
	}
}
