package vm

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// ResolveCompiledMethod finds a compiled method in the program's function table
// and executes it with the given receiver.
func ResolveCompiledMethod(program *bytecode.CompiledProgram, methodName string, receiver value.Value, extraArgs ...value.Value) (value.Value, bool) {
	return resolveCompiledMethod(program, methodName, receiver, extraArgs)
}

// ResolveCompiledMethodWithArgs resolves and calls a compiled method with extra arguments.
// Used for methods like Is(error) bool that need parameters beyond the receiver.
func ResolveCompiledMethodWithArgs(program *bytecode.CompiledProgram, methodName string, receiver value.Value, args []value.Value) (value.Value, bool) {
	return resolveCompiledMethod(program, methodName, receiver, args)
}

func resolveCompiledMethod(program *bytecode.CompiledProgram, methodName string, receiver value.Value, args []value.Value) (value.Value, bool) {
	rv, receiverTypeName, ok := compiledMethodReceiverInfo(program, receiver)
	if !ok {
		return value.MakeNil(), false
	}

	fn, methodReceiver, ok := selectCompiledMethodCandidate(program, methodName, receiverTypeName, receiver)
	if !ok {
		return value.MakeNil(), false
	}

	methodReceiver = normalizeCompiledMethodReceiver(fn, receiverTypeName, rv, methodReceiver)
	return invokeCompiledMethod(program, methodName, fn, methodReceiver, args)
}

func compiledMethodReceiverInfo(program *bytecode.CompiledProgram, receiver value.Value) (reflect.Value, string, bool) {
	rv, ok := receiver.ReflectValue()
	if !ok {
		// Fallback: try Interface() → reflect.ValueOf
		iface := receiver.Interface()
		if iface == nil {
			return reflect.Value{}, "", false
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
	receiverTypeName := compiledMethodReceiverTypeName(program, rv)
	if receiverTypeName == "" {
		return reflect.Value{}, "", false
	}
	return rv, receiverTypeName, true
}

func compiledMethodReceiverTypeName(program *bytecode.CompiledProgram, rv reflect.Value) string {
	if rv.Kind() == reflect.Ptr {
		elemType := rv.Type().Elem()
		if elemType.Kind() == reflect.Struct {
			return fallbackCompiledMethodReceiverTypeName(program.LookupTypeName(elemType), rv)
		}
	}
	if rv.Kind() == reflect.Struct {
		return fallbackCompiledMethodReceiverTypeName(program.LookupTypeName(rv.Type()), rv)
	}
	return pkgPathTypeName(rv.Type())
}

func fallbackCompiledMethodReceiverTypeName(name string, rv reflect.Value) string {
	if name != "" {
		return name
	}
	// Fallback: scan unexported field PkgPath for the interpreter type-name suffix.
	return pkgPathTypeName(rv.Type())
}

func normalizeCompiledMethodReceiver(
	fn *bytecode.CompiledFunction,
	receiverTypeName string,
	rv reflect.Value,
	methodReceiver value.Value,
) value.Value {
	// Normalize the receiver through reflect so the interpreter always sees a
	// clean, concretely-typed value. Without this, a receiver that lives inside
	// an interface{} box causes reflect.Set panics when the method body accesses
	// fields on it. Do not replace a receiver rebound to an embedded field for a
	// promoted method.
	if rv.IsValid() && fn.ReceiverTypeName == receiverTypeName {
		concrete := reflect.New(rv.Type()).Elem()
		concrete.Set(rv)
		return value.MakeFromReflect(concrete)
	}
	return methodReceiver
}

func invokeCompiledMethod(
	program *bytecode.CompiledProgram,
	methodName string,
	fn *bytecode.CompiledFunction,
	methodReceiver value.Value,
	args []value.Value,
) (value.Value, bool) {
	// This side-channel VM has no initialGlobals because method dispatch only
	// needs to execute the selected method, not replay program initialization.
	tempVM := newTempVM(program, make([]value.Value, len(program.Globals)), nil, nil, context.Background(), nil)
	tempVM.callFunction(fn, compiledMethodCallArgs(methodReceiver, args), nil)
	result, err := runCompiledMethodSideChannel(tempVM, methodName)
	if err != nil {
		return value.MakeNil(), false
	}
	return result, true
}

func compiledMethodCallArgs(receiver value.Value, args []value.Value) []value.Value {
	callArgs := make([]value.Value, 0, 1+len(args))
	callArgs = append(callArgs, receiver)
	return append(callArgs, args...)
}

func runCompiledMethodSideChannel(tempVM *vm, methodName string) (result value.Value, err error) {
	// If the method body hits a Go-level panic that the VM panic protocol does
	// not convert to an error, recover so callers can fall back to their default
	// behavior instead of crashing during formatting or interface probing.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("side-channel method %q panicked: %v", methodName, r)
		}
	}()
	return tempVM.run()
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
