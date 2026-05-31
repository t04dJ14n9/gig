package vm

import (
	"context"

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

type interpretedErrorAdapter struct {
	program          *bytecode.CompiledProgram
	receiver         value.Value
	receiverTypeName string
	globals          []value.Value
	initialGlobals   []value.Value
	shared           *SharedGlobals
	ctx              context.Context
	goroutines       *GoroutineTracker
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

func newInterpretedErrorAdapter(
	program *bytecode.CompiledProgram,
	receiver value.Value,
	receiverTypeName string,
	globals []value.Value,
	initialGlobals []value.Value,
	shared *SharedGlobals,
	ctx context.Context,
	goroutines *GoroutineTracker,
) *interpretedErrorAdapter {
	return &interpretedErrorAdapter{
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

func (a *interpretedErrorAdapter) Error() string {
	result, ok := callInterfaceMethodValue(a.program, "Error", a.receiverTypeName, a.receiver, nil, a.globals, a.initialGlobals, a.shared, a.ctx, a.goroutines)
	if !ok {
		return ""
	}
	return result.String()
}

func (a *interpretedInterfaceAdapter) Len() int {
	result := a.call("Len")
	return int(result.Int())
}

func (a *interpretedInterfaceAdapter) Less(i, j int) bool {
	result := a.call("Less", value.MakeInt(int64(i)), value.MakeInt(int64(j)))
	return result.Bool()
}

func (a *interpretedInterfaceAdapter) Swap(i, j int) {
	// Always try the compiled method first — user-defined Swap may update
	// auxiliary state (counters, indexes, parallel slices) beyond just
	// swapping elements. Only fall back to direct slice swap if no
	// compiled method is available.
	result, ok := callInterfaceMethodValue(
		a.program,
		"Swap",
		a.receiverTypeName,
		a.receiver,
		[]value.Value{value.MakeInt(int64(i)), value.MakeInt(int64(j))},
		a.globals,
		a.initialGlobals,
		a.shared,
		a.ctx,
		a.goroutines,
	)
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
