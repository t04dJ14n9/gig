package vm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

type interpretedInterfaceAdapter struct {
	program          *bytecode.CompiledProgram
	receiver         value.Value
	receiverTypeName string
	// Caller VM context for correct callback execution
	globals        []value.Value
	initialGlobals []value.Value
	shared         *SharedGlobals
	ctx            context.Context
	goroutines     *GoroutineTracker
}

func newInterpretedInterfaceAdapter(
	program *bytecode.CompiledProgram,
	receiver value.Value,
	receiverTypeName string,
	globals []value.Value,
	initialGlobals []value.Value,
	shared *SharedGlobals,
	ctx context.Context,
	goroutines *GoroutineTracker,
) *interpretedInterfaceAdapter {
	return &interpretedInterfaceAdapter{
		program:          program,
		receiver:         receiver,
		receiverTypeName: receiverTypeName,
		globals:          globals,
		initialGlobals:   initialGlobals,
		shared:           shared,
		ctx:              ctx,
		goroutines:       goroutines,
	}
}

func (a *interpretedInterfaceAdapter) Len() int {
	result := a.call("Len")
	return int(result.Int())
}

func (a *interpretedInterfaceAdapter) Less(i, j int) bool {
	if a.receiverTypeName == "Reverse" || strings.HasSuffix(a.receiverTypeName, ".Reverse") {
		if result, ok := callEmbeddedInterfaceMethod(a.receiver, "Less", value.MakeInt(int64(j)), value.MakeInt(int64(i))); ok {
			return result.Bool()
		}
	}
	result := a.call("Less", value.MakeInt(int64(i)), value.MakeInt(int64(j)))
	return result.Bool()
}

func (a *interpretedInterfaceAdapter) Swap(i, j int) {
	// Always try the compiled method first — user-defined Swap may update
	// auxiliary state (counters, indexes, parallel slices) beyond just
	// swapping elements. Only fall back to direct slice swap if no
	// compiled method is available.
	result, ok := callInterfaceMethodValue(a.program, "Swap", a.receiverTypeName, a.receiver, []value.Value{value.MakeInt(int64(i)), value.MakeInt(int64(j))}, a.globals, a.initialGlobals, a.shared, a.ctx, a.goroutines)
	_ = result
	if ok {
		return
	}
	// Fallback: direct slice element swap (for types without a compiled Swap)
	callReceiverSliceSwap(a.receiver, i, j)
}

func (a *interpretedInterfaceAdapter) Push(x any) {
	a.call("Push", value.FromInterface(x))
}

func (a *interpretedInterfaceAdapter) Pop() any {
	return a.call("Pop").Interface()
}

func (a *interpretedInterfaceAdapter) call(methodName string, args ...value.Value) value.Value {
	result, ok := callInterfaceMethodValue(a.program, methodName, a.receiverTypeName, a.receiver, args, a.globals, a.initialGlobals, a.shared, a.ctx, a.goroutines)
	if !ok {
		return value.MakeNil()
	}
	return result
}

func callInterfaceMethodValue(program *bytecode.CompiledProgram, methodName, receiverTypeName string, receiver value.Value, args []value.Value, globals []value.Value, initialGlobals []value.Value, shared *SharedGlobals, ctx context.Context, goroutines *GoroutineTracker) (value.Value, bool) {
	if result, ok := callCompiledMethodValue(program, methodName, receiverTypeName, receiver, args, globals, initialGlobals, shared, ctx, goroutines); ok {
		return result, true
	}
	return callReflectInterfaceMethod(methodName, receiver, args)
}

func callCompiledMethodValue(program *bytecode.CompiledProgram, methodName, receiverTypeName string, receiver value.Value, args []value.Value, globals []value.Value, initialGlobals []value.Value, shared *SharedGlobals, ctx context.Context, goroutines *GoroutineTracker) (value.Value, bool) {
	if program == nil {
		return value.MakeNil(), false
	}

	for _, fn := range program.MethodsByName[methodName] {
		if receiverTypeName != "" && fn.ReceiverTypeName != receiverTypeName {
			continue
		}
		if shouldPanicOnNilValueReceiver(receiver, fn) {
			return value.MakeNil(), false
		}

		methodReceiver := receiverForCompiledMethod(methodName, receiver)
		callArgs := make([]value.Value, 0, len(args)+1)
		callArgs = append(callArgs, methodReceiver)
		callArgs = append(callArgs, args...)

		tempVM := newTempVM(program, globals, initialGlobals, shared, ctx, goroutines)
		tempVM.callFunction(fn, callArgs, nil)

		var result value.Value
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("compiled method %q panicked: %v", methodName, r)
				}
			}()
			result, err = tempVM.run()
		}()
		if err != nil {
			continue
		}
		return result, true
	}

	return value.MakeNil(), false
}

// newTempVM creates a standalone VM for executing compiled methods from adapter
// callbacks. It copies the caller's globals and context so the callback sees
// correct program state. The VM is NOT pooled — adapter callbacks are short-lived
// and the pool lifecycle is managed by the runner, not reachable from here.
func newTempVM(program *bytecode.CompiledProgram, globals, initialGlobals []value.Value, shared *SharedGlobals, ctx context.Context, goroutines *GoroutineTracker) *vm {
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

func callReceiverSliceSwap(receiver value.Value, i, j int) bool {
	rv, ok := receiver.ReflectValue()
	if !ok {
		return false
	}
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return false
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		return false
	}
	tmp := reflect.New(rv.Type().Elem()).Elem()
	tmp.Set(rv.Index(i))
	rv.Index(i).Set(rv.Index(j))
	rv.Index(j).Set(tmp)
	return true
}

func callReflectInterfaceMethod(methodName string, receiver value.Value, args []value.Value) (value.Value, bool) {
	rv, ok := receiver.ReflectValue()
	if !ok {
		iface := receiver.Interface()
		if iface == nil {
			return value.MakeNil(), false
		}
		rv = reflect.ValueOf(iface)
	}
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
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

func callEmbeddedInterfaceMethod(receiver value.Value, methodName string, args ...value.Value) (value.Value, bool) {
	rv, ok := receiver.ReflectValue()
	if !ok {
		return value.MakeNil(), false
	}
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return value.MakeNil(), false
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return value.MakeNil(), false
	}
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.Kind() != reflect.Interface || field.IsNil() {
			continue
		}
		method := field.Elem().MethodByName(methodName)
		if !method.IsValid() {
			continue
		}
		in := buildReflectArgs(args, method.Type())
		out := method.Call(in)
		if len(out) == 0 {
			return value.MakeNil(), true
		}
		return value.MakeFromReflect(out[0]), true
	}
	return value.MakeNil(), false
}
