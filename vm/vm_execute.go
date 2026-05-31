package vm

import (
	"context"
	"fmt"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// Execute runs the specified function with the given arguments.
// A Go-level recover() safety net catches host panics and converts them to errors.
func (v *vm) Execute(funcName string, ctx context.Context, args ...value.Value) (result value.Value, err error) {
	v.ctx = ctx

	fn, ok := v.program.Functions[funcName]
	if !ok {
		return value.MakeNil(), fmt.Errorf("function %q not found", funcName)
	}

	valArgs := make([]value.Value, len(args))
	copy(valArgs, args)

	valArgs, err = v.validateAndPrepareEntryArgs(fn, valArgs)
	if err != nil {
		return value.MakeNil(), err
	}

	v.startEntryFrame(fn, valArgs)

	defer func() {
		if r := recover(); r != nil {
			result = value.MakeNil()
			err = fmt.Errorf("runtime panic: %v", r)
		}
	}()

	result, err = v.run()
	return result, err
}

// ExecuteWithValues runs the specified function with pre-converted Value arguments.
// Includes the same Go-level panic safety net as Execute.
func (v *vm) ExecuteWithValues(funcName string, ctx context.Context, args []value.Value) (result value.Value, err error) {
	v.ctx = ctx

	fn, ok := v.program.Functions[funcName]
	if !ok {
		return value.MakeNil(), fmt.Errorf("function %q not found", funcName)
	}

	args, err = v.validateAndPrepareEntryArgs(fn, args)
	if err != nil {
		return value.MakeNil(), err
	}

	v.startEntryFrame(fn, args)

	defer func() {
		if r := recover(); r != nil {
			result = value.MakeNil()
			err = fmt.Errorf("runtime panic: %v", r)
		}
	}()

	return v.run()
}

func (v *vm) startEntryFrame(fn *bytecode.CompiledFunction, args []value.Value) {
	frame := v.fpool.get(fn, 0, nil)
	for i, arg := range args {
		if i < fn.NumLocals {
			frame.locals[i] = arg
			if frame.intLocals != nil {
				frame.intLocals[i] = arg.RawInt()
			}
		}
	}
	v.frames[0] = frame
	v.fp = 1
}
