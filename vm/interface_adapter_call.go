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
	fn, methodReceiver, ok := selectInterfaceMethodCandidate(program, methodName, receiverTypeName, receiver)
	if !ok {
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

func selectInterfaceMethodCandidate(
	program *bytecode.CompiledProgram,
	methodName string,
	receiverTypeName string,
	receiver value.Value,
) (*bytecode.CompiledFunction, value.Value, bool) {
	if program == nil {
		return nil, value.MakeNil(), false
	}
	for _, fn := range program.MethodsByName[methodName] {
		if fn.ReceiverTypeName == "" {
			continue
		}
		if receiverTypeName != "" && fn.ReceiverTypeName != receiverTypeName {
			continue
		}
		if dyn, ok := receiver.InterpretedInterface(); ok && fn.ReceiverIsPointer && !dyn.IsPointer {
			continue
		}
		return fn, methodReceiverForCompiledFunction(receiver, fn), true
	}
	if receiverTypeName == "" {
		inferred := inferReceiverTypeName(receiver, program)
		for _, fn := range program.MethodsByName[methodName] {
			if fn.ReceiverTypeName == inferred {
				return fn, methodReceiverForCompiledFunction(receiver, fn), true
			}
		}
	}
	return nil, value.MakeNil(), false
}

func runInterfaceCallbackVM(tempVM *vm, methodName string) (result value.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("compiled method %q panicked: %v", methodName, r)
		}
	}()
	return tempVM.run()
}

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
		extCallCache: &externalCallCache{
			cache: make([]*extCallCacheEntry, len(program.Constants)),
		},
	}
	if shared != nil {
		v.shared = shared
	}
	return v
}

func callReceiverSliceSwap(receiver value.Value, i, j int) bool {
	rv, ok := receiver.ReflectValue()
	if !ok {
		return false
	}
	rv, ok = derefInterfaceOrPointer(rv)
	if !ok || rv.Kind() != reflect.Slice {
		return false
	}
	tmp := reflect.New(rv.Type().Elem()).Elem()
	tmp.Set(rv.Index(i))
	rv.Index(i).Set(rv.Index(j))
	rv.Index(j).Set(tmp)
	return true
}

func callReflectInterfaceMethod(methodName string, receiver value.Value, args []value.Value) (value.Value, bool) {
	rv, ok := reflectValueForInterfaceMethod(receiver)
	if !ok {
		return value.MakeNil(), false
	}
	method, found := findMethod(rv, methodName, nil)
	if !found {
		return value.MakeNil(), false
	}
	in := buildReflectArgs(args, method.Type())
	out := method.Call(in)
	if len(out) == 0 {
		return value.MakeNil(), true
	}
	return value.MakeFromReflect(out[0]), true
}

func reflectValueForInterfaceMethod(receiver value.Value) (reflect.Value, bool) {
	rv, ok := receiver.ReflectValue()
	if !ok {
		iface := receiver.Interface()
		if iface == nil {
			return reflect.Value{}, false
		}
		rv = reflect.ValueOf(iface)
	}
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	return rv, true
}

func derefInterfaceOrPointer(rv reflect.Value) (reflect.Value, bool) {
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, false
		}
		rv = rv.Elem()
	}
	return rv, true
}
