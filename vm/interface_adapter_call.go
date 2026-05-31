package vm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func callInterfaceMethodValue(
	program *bytecode.CompiledProgram,
	methodName string,
	receiverTypeName string,
	receiver value.Value,
	args []value.Value,
	globals []value.Value,
	initialGlobals []value.Value,
	shared *SharedGlobals,
	ctx context.Context,
	goroutines *GoroutineTracker,
) (value.Value, bool) {
	// Prefer compiled script methods. Reflect fallback is only for embedded or
	// host-provided receivers, preserving the boundary that script-defined custom
	// types do not masquerade as arbitrary third-party interface values.
	if result, ok := callCompiledMethodValue(program, methodName, receiverTypeName, receiver, args, globals, initialGlobals, shared, ctx, goroutines); ok {
		return result, true
	}
	return callReflectInterfaceMethod(methodName, receiver, args)
}

func callCompiledMethodValue(
	program *bytecode.CompiledProgram,
	methodName string,
	receiverTypeName string,
	receiver value.Value,
	args []value.Value,
	globals []value.Value,
	initialGlobals []value.Value,
	shared *SharedGlobals,
	ctx context.Context,
	goroutines *GoroutineTracker,
) (value.Value, bool) {
	// Adapter callbacks run while host code is executing. They need a temporary
	// VM with the caller's globals and shared state, but they must not reuse the
	// caller's stack or frames because host calls can be nested.
	if program == nil {
		return value.MakeNil(), false
	}

	fn, methodReceiver, ok := selectCompiledMethodCandidate(program, methodName, receiverTypeName, receiver)
	if !ok {
		return value.MakeNil(), false
	}
	if shouldPanicOnNilValueReceiver(receiver, fn) {
		return value.MakeNil(), false
	}

	callArgs := make([]value.Value, 0, len(args)+1)
	callArgs = append(callArgs, methodReceiver)
	callArgs = append(callArgs, args...)

	tempVM := newTempVM(program, globals, initialGlobals, shared, ctx, goroutines)
	tempVM.callFunction(fn, callArgs, nil)

	result, err := runInterfaceCallbackVM(tempVM, methodName)
	if err != nil {
		return value.MakeNil(), false
	}
	return result, true
}

func runInterfaceCallbackVM(tempVM *vm, methodName string) (result value.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("compiled method %q panicked: %v", methodName, r)
		}
	}()
	return tempVM.run()
}

// newTempVM creates a standalone VM for executing compiled methods from adapter
// callbacks. It copies the caller's globals and context so the callback sees
// correct program state. The VM is NOT pooled — adapter callbacks are short-lived
// and the pool lifecycle is managed by the runner, not reachable from here.
func newTempVM(
	program *bytecode.CompiledProgram,
	globals []value.Value,
	initialGlobals []value.Value,
	shared *SharedGlobals,
	ctx context.Context,
	goroutines *GoroutineTracker,
) *vm {
	callGlobals := make([]value.Value, len(program.Globals))
	if globals != nil {
		copy(callGlobals, globals)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	v := &vm{
		program:        program,
		stack:          make([]value.Value, deferVMStackSize),
		frames:         make([]*Frame, initialFrameDepth),
		globals:        callGlobals,
		initialGlobals: initialGlobals,
		ctx:            ctx,
		goroutines:     goroutines,
	}
	if shared != nil {
		v.shared = shared
	}
	return v
}

func receiverForCompiledMethod(methodName string, receiver value.Value) value.Value {
	// Sort methods on heap-like values often have value receivers, while Push
	// and Pop require a pointer receiver. Normalize only the value-receiver sort
	// methods so mutations still happen through the compiled method body.
	if dyn, ok := receiver.InterpretedInterface(); ok {
		return dyn.Value
	}

	rv, ok := receiver.ReflectValue()
	if !ok || !rv.IsValid() || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return receiver
	}

	// heap.Interface mixes value-receiver sort methods with pointer-receiver Push/Pop.
	// Value-receiver methods should see the pointed-to named slice value; Push/Pop need the pointer.
	switch methodName {
	case "Len", "Less", "Swap":
		elem := rv.Elem()
		if elem.IsValid() && elem.Kind() == reflect.Slice {
			return value.MakeFromReflect(elem)
		}
	}

	return receiver
}
