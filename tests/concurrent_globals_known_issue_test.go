package tests

// This file contains known issues with concurrent global variables.
// These tests document current limitations and should be fixed in future updates.
// Each test uses t.Skip() to avoid blocking the test suite.
//
// Known Issues:
//
// 1. TestSyncCondKnownIssue - sync.Cond.Wait() + Broadcast
//    Issue: cond.Wait() hangs indefinitely, never wakes up after Broadcast().
//    Root cause: sync.NewCond(&mu) requires the Locker interface. The value
//                passed to NewCond may not reference the same mutex object
//                that the guest code locks/unlocks, or the internal
//                Unlock/Lock cycle within cond.Wait() doesn't work correctly
//                on the interpreted mutex.
//
// Fixed Issues (removed from skip list):
// - defer mu.Unlock() — fixed by handling external method wrappers in compileDefer
// - defer in goroutines — fixed by the same compileDefer fix

import (
	"testing"
)

// TestSyncCondKnownIssue documents the sync.Cond Wait/Broadcast issue.
//
// Pattern: Waiter calls cond.Wait() in a loop; producer calls cond.Broadcast().
// Expected: Waiter wakes up after Broadcast.
// Actual: cond.Wait() hangs forever.
func TestSyncCondKnownIssue(t *testing.T) {
	t.Skip("KNOWN ISSUE: sync.Cond.Wait() hangs — cond.Wait() internal Unlock/Lock cycle does not work on interpreted mutex")
}
