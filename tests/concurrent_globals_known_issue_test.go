package tests

// This file previously contained known issues with concurrent global variables.
// All issues have been fixed:
//
// Fixed Issues:
// - sync.Cond.Wait() + Broadcast — now works correctly
// - defer mu.Unlock() — fixed by handling external method wrappers in compileDefer
// - defer in goroutines — fixed by the same compileDefer fix
//
// No known issues remain. This file is kept for historical reference.
