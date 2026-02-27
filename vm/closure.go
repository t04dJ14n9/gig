package vm

import (
	"gig/bytecode"
	"gig/value"
)

// Closure represents a closure with captured free variables.
// When a closure is called, its free variables are bound to the calling context.
type Closure struct {
	// Fn is the compiled function bytecode.
	Fn *bytecode.CompiledFunction

	// FreeVars are pointers to captured variables.
	// They are stored as pointers to allow shared state between closures.
	FreeVars []*value.Value
}
