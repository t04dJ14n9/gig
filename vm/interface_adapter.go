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
	globals          []value.Value
	initialGlobals   []value.Value
	shared           *SharedGlobals
	ctx              context.Context
	goroutines       *GoroutineTracker
}

type interpretedAdapterError struct {
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

func newInterpretedAdapterError(
	program *bytecode.CompiledProgram,
	receiver value.Value,
	receiverTypeName string,
	globals []value.Value,
	initialGlobals []value.Value,
	shared *SharedGlobals,
	ctx context.Context,
	goroutines *GoroutineTracker,
) *interpretedAdapterError {
	return &interpretedAdapterError{
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

func (a *interpretedAdapterError) Error() string {
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
