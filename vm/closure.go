package vm

import (
	"sync"

	"github.com/t04dJ14n9/gig/bytecode"
	"github.com/t04dJ14n9/gig/value"
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

// closurePool pools Closure objects to reduce heap allocations.
// Closures are returned to the pool after they finish executing (callIndirect/callFunction).
var closurePool = sync.Pool{
	New: func() any {
		return &Closure{}
	},
}

// getClosure returns a Closure from the pool, resized for numFree.
func getClosure(fn *bytecode.CompiledFunction, numFree int) *Closure {
	c := closurePool.Get().(*Closure)
	c.Fn = fn
	if numFree == 0 {
		c.FreeVars = nil
	} else if cap(c.FreeVars) >= numFree {
		c.FreeVars = c.FreeVars[:numFree]
	} else {
		c.FreeVars = make([]*value.Value, numFree)
	}
	return c
}

// putClosure returns a Closure to the pool.
func putClosure(c *Closure) {
	c.Fn = nil
	// Don't clear FreeVars pointers — the frame.freeVars may still reference them.
	// Just nil the slice header so GC can reclaim if no frame references it.
	c.FreeVars = c.FreeVars[:0]
	closurePool.Put(c)
}
