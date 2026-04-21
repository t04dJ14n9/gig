package vm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
)

// ResolveCompiledMethod finds a compiled method in the program's function table
// and executes it with the given receiver.
func ResolveCompiledMethod(program *bytecode.CompiledProgram, methodName string, receiver value.Value) (value.Value, bool) {
	rv, ok := receiver.ReflectValue()
	if !ok {
		// Fallback: try Interface() → reflect.ValueOf
		iface := receiver.Interface()
		if iface == nil {
			return value.MakeNil(), false
		}
		rv = reflect.ValueOf(iface)
	}
	// Unwrap interface layers — the receiver may be stored as an interface{} reflect.Value
	// (e.g., when passed through fmt.Sprint → FmtWrap → resolveErrorer → callMethod).
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}

	// Extract the type name using the program-level ReflectTypeNames registry,
	// falling back to scanning unexported field PkgPath for the # suffix.
	receiverTypeName := ""
	if rv.Kind() == reflect.Ptr {
		elem := rv.Elem()
		if elem.Kind() == reflect.Struct {
			receiverTypeName = program.LookupTypeName(elem.Type())
		}
	} else if rv.Kind() == reflect.Struct {
		receiverTypeName = program.LookupTypeName(rv.Type())
	}
	// Fallback: scan unexported field PkgPath for # suffix
	if receiverTypeName == "" && rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			sf := rt.Field(i)
			if idx := strings.LastIndex(sf.PkgPath, "#"); idx >= 0 {
				qualName := sf.PkgPath[idx+1:]
				if dotIdx := strings.LastIndex(qualName, "."); dotIdx >= 0 {
					receiverTypeName = qualName[dotIdx+1:]
				} else {
					receiverTypeName = qualName
				}
				break
			}
		}
	}

	if receiverTypeName == "" {
		return value.MakeNil(), false
	}

	// Search the compiled function table for a method with matching name and receiver type
	for _, fn := range program.MethodsByName[methodName] {
		if fn.ReceiverTypeName != receiverTypeName {
			continue
		}
		// Found the method! Execute it with a temporary VM.
		tempVM := &vm{
			program: program,
			stack:   make([]value.Value, deferVMStackSize),
			sp:      0,
			frames:  make([]*Frame, initialFrameDepth),
			fp:      0,
			globals: make([]value.Value, len(program.Globals)),
			ctx:     context.Background(),
		}
		// Note: tempVM does not have initialGlobals since resolveCompiledMethod
		// is called without a VM context. This is acceptable because method resolution
		// only needs to execute the method, not full program init.
		//
		// Normalize the receiver through reflect so the interpreter always sees
		// a clean, concretely-typed value. Without this, a receiver that lives
		// inside an interface{} box causes reflect.Set panics when the method
		// body accesses fields on it.
		methodReceiver := receiver
		// Use the already-unwrapped rv (not receiver.ReflectValue()) to avoid
		// re-wrapping an interface{} layer.
		if rv.IsValid() {
			concrete := reflect.New(rv.Type()).Elem()
			concrete.Set(rv)
			methodReceiver = value.MakeFromReflect(concrete)
		}
		tempVM.callFunction(fn, []value.Value{methodReceiver}, nil)
		// Side-channel method invocation: if the method body encounters a
		// Go-level panic (e.g., reflect operations on mismatched types) that
		// the VM's own panic protocol doesn't convert to an error, recover
		// so the caller gets (nil, false) and can fall back to the default
		// formatting instead of crashing.
		var result value.Value
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("side-channel method %q panicked: %v", methodName, r)
				}
			}()
			result, err = tempVM.run()
		}()
		if err != nil {
			return value.MakeNil(), false
		}
		return result, true
	}
	return value.MakeNil(), false
}

// VMPool is a lock-free pool of VMs for a given program using sync.Pool.
// This provides better performance under high concurrency compared to mutex-based pools.
type VMPool struct {
	pool sync.Pool
}

// NewVMPool creates a VM pool for the given program.
func NewVMPool(program *bytecode.CompiledProgram, initialGlobals []value.Value, goroutines *GoroutineTracker) *VMPool {
	return &VMPool{
		pool: sync.Pool{
			New: func() any {
				return newVM(program, initialGlobals, goroutines)
			},
		},
	}
}

// Get returns an idle VM from the pool, or creates a new one if the pool is empty.
func (p *VMPool) Get() VM {
	return p.pool.Get().(*vm)
}

// Put returns a VM to the pool for reuse.
// The VM is reset before being returned to the pool.
func (p *VMPool) Put(x VM) {
	x.Reset()
	p.pool.Put(x.(*vm))
}
