package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeGoCall(frame *Frame) error {
	// OpGoCall consumes its arguments immediately; the child VM receives a
	// private args slice while sharing the parent's global state and context.
	funcIdx := frame.readUint16()
	args := v.popCallArgs(frame.readByte())
	fn := v.compiledFunctionByIndex(funcIdx)
	if fn == nil {
		return nil
	}
	return v.startCompiledGoroutine(fn, args, nil)
}

func (v *vm) executeGoCallExternal(frame *Frame) error {
	funcIdx := frame.readUint16()
	args := v.popCallArgs(frame.readByte())
	if int(funcIdx) >= len(v.program.Constants) {
		return nil
	}

	switch info := v.program.Constants[funcIdx].(type) {
	case *external.ExternalMethodInfo:
		return v.startExternalMethodGoroutine(info, args)
	case *external.ExternalFuncInfo:
		return v.startExternalFuncGoroutine(info, args)
	default:
		return nil
	}
}

func (v *vm) executeGoCallIndirect(frame *Frame) error {
	// Indirect goroutine calls only support interpreted closures here. Reflect
	// function values are handled by the direct hot path, not this cold opcode.
	args := v.popCallArgs(frame.readByte())
	callee := v.pop()
	closure, ok := callee.RawObj().(*Closure)
	if !ok {
		return nil
	}
	return v.startCompiledGoroutine(closure.Fn, args, closure.FreeVars)
}

func (v *vm) popCallArgs(numArgs byte) []value.Value {
	args := make([]value.Value, numArgs)
	for i := int(numArgs) - 1; i >= 0; i-- {
		args[i] = v.pop()
	}
	return args
}

func (v *vm) startCompiledGoroutine(fn *bytecode.CompiledFunction, args []value.Value, freeVars []*value.Value) error {
	childVM := v.newChildVM()
	capturedFn := fn
	capturedArgs := args
	capturedFreeVars := freeVars
	return v.goroutines.Start(func() {
		childFrame := newFrame(capturedFn, capturedArgs, capturedFreeVars)
		childVM.frames[0] = childFrame
		childVM.fp = 1

		// Goroutine calls intentionally discard return values, matching Go's
		// source-level `go f()` semantics.
		_, _ = childVM.run()
	})
}

func (v *vm) startExternalMethodGoroutine(info *external.ExternalMethodInfo, args []value.Value) error {
	if err := v.validateExternalMethodBoundary(info, args); err != nil {
		return err
	}
	childVM := v.newChildVM()
	capturedInfo := info
	capturedArgs := args
	return v.goroutines.Start(func() {
		_ = childVM.callExternalMethod(capturedInfo, capturedArgs)
	})
}

func (v *vm) startExternalFuncGoroutine(info *external.ExternalFuncInfo, args []value.Value) error {
	rc := bytecode.ResolveConstant(info)
	if err := v.validateExternalBoundary(rc, args); err != nil {
		return err
	}
	childVM := v.newChildVM()
	capturedRC := rc
	capturedArgs := args
	return v.goroutines.Start(func() {
		_ = childVM.callResolvedExternal(capturedRC, capturedArgs)
	})
}
