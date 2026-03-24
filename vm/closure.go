// closure.go defines the Closure type and ClosureExecutor for reflect.MakeFunc integration.
package vm

import (
	"context"
	"reflect"
	"sync"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/value"
)

// Closure represents a closure with captured free variables.
// When a closure is called, its free variables are bound to the calling context.
// Closure implements value.ClosureExecutor so that value.ToReflectValue can
// wrap it into a real Go function via reflect.MakeFunc without a global callback.
type Closure struct {
	// Fn is the compiled function bytecode.
	Fn *bytecode.CompiledFunction

	// FreeVars are pointers to captured variables.
	// They are stored as pointers to allow shared state between closures.
	FreeVars []*value.Value

	// Program is a reference to the compiled program, needed when the closure
	// is wrapped as a real Go function (via reflect.MakeFunc) for typed containers.
	Program *bytecode.Program

	// InitialGlobals is the post-init globals snapshot used to seed temporary VMs
	// when this closure is converted to a real Go function via Execute().
	InitialGlobals []value.Value
}

// Execute runs the closure in a temporary VM and returns the results as reflect.Values.
// This implements value.ClosureExecutor, allowing value.ToReflectValue to convert
// closures to real Go functions without a global callback.
func (c *Closure) Execute(args []reflect.Value, outTypes []reflect.Type) []reflect.Value {
	if c.Program == nil {
		return nil
	}
	// Create a temporary VM to execute the closure
	closureVM := &vm{
		program: c.Program,
		stack:   make([]value.Value, 256),
		sp:      0,
		frames:  make([]*Frame, initialFrameDepth),
		fp:      0,
		globals: make([]value.Value, len(c.Program.Globals)),
		ctx:     context.Background(),
		extCallCache: &externalCallCache{
			cache: make([]*extCallCacheEntry, len(c.Program.Constants)),
		},
	}
	if len(c.InitialGlobals) == len(closureVM.globals) {
		copy(closureVM.globals, c.InitialGlobals)
	}
	closureVM.initialGlobals = c.InitialGlobals
	// Convert reflect.Value args to value.Value args
	valArgs := make([]value.Value, len(args))
	for i, arg := range args {
		valArgs[i] = value.MakeFromReflect(arg)
	}
	// Call the closure function with its captured free variables
	closureVM.callFunction(c.Fn, valArgs, c.FreeVars)
	result, _ := closureVM.run()
	// Return the result as reflect.Value
	if result.Kind() == value.KindNil {
		return []reflect.Value{}
	}
	if len(outTypes) > 0 {
		return []reflect.Value{result.ToReflectValue(outTypes[0])}
	}
	iface := result.Interface()
	if iface == nil {
		return []reflect.Value{}
	}
	return []reflect.Value{reflect.ValueOf(iface)}
}

// closurePool pools Closure objects to reduce heap allocations.
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
