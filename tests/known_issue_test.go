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
	"testing"
)

// TestKnownIssues_Tricky runs all known interpreter bugs.
// Every sub-test here is EXPECTED TO FAIL — they document real bugs.
func TestKnownIssues_Tricky(t *testing.T) {
	// All known issues have been resolved and moved to resolved_issue:
	// - Bug 1: sort named-type conversion (OpChangeType) → Resolved Issue 28
	// - Bug 2: time.Duration DirectCall wrappers → fixed in gentool
	// - Bug 3: fmt.Stringer on interpreted types → Resolved Issue 29
	// - Bug 4: fmt.Sprintf %T type name → Resolved Issue 30
	// - Bug 5: fmt.Sprintf %v _gig_id field → Resolved Issue 31
	// - Bug 6: prog.Run() int64/uint64 narrowing → Resolved Issue 32
	// - Bug 7: bytes.Buffer.Cap() string→[]byte cap → Resolved Issue 33
	t.Skip("No known issues — all resolved! See resolved_issue/main.go")
}
