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
// 2. TestDeferInGoroutineKnownIssue - defer in goroutine does not execute
//    Issue: defer statements inside goroutines (e.g., defer wg.Done()) do not
//           execute, causing WaitGroup.Wait() to hang forever.
//    Root cause: Gig's goroutine implementation does not run deferred functions
//                when the goroutine body completes.
//
// 3. TestDeferUnlockKnownIssue - defer mu.Unlock() may not execute reliably
//    Issue: Using "defer mu.Unlock()" after mu.Lock() may leave the mutex
//           permanently locked, causing deadlocks in subsequent Lock() calls.
//    Workaround: Use explicit mu.Unlock() instead of defer.
//    Root cause: Related to the defer execution reliability issue — if the
//                defer stack is not fully drained on function return, Unlock()
//                is skipped.

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

// TestDeferInGoroutineKnownIssue documents that defer in goroutines does not execute.
//
// Pattern: go func() { defer wg.Done(); ... }(); wg.Wait()
// Expected: wg.Done() called via defer, wg.Wait() returns.
// Actual: defer never runs, wg.Wait() hangs forever.
func TestDeferInGoroutineKnownIssue(t *testing.T) {
	t.Skip("KNOWN ISSUE: defer in goroutine does not execute — wg.Done() never called, WaitGroup.Wait() hangs")
}

// TestDeferUnlockKnownIssue documents that "defer mu.Unlock()" may not execute.
//
// Pattern: mu.Lock(); defer mu.Unlock(); ...; return
// Expected: Unlock() called via defer, mutex released.
// Actual: Unlock() may not execute, mutex stays locked, subsequent Lock() deadlocks.
// Workaround: Use explicit mu.Unlock() instead of defer.
func TestDeferUnlockKnownIssue(t *testing.T) {
	t.Skip("KNOWN ISSUE: defer mu.Unlock() may not execute — use explicit mu.Unlock() instead")
}
