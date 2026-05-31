package vm

import (
	"fmt"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

// callExternal calls an external (native Go) function or method.
//
// Redesigned: the constant pool is pre-resolved into program.ExternCalls[] at
// compile time (see bytecode.ResolveExternCalls). This eliminates the old
// per-VM extCallCache with its RWMutex — all VMs sharing the same program do
// a lock-free array lookup.
//
// Flow: OpCallExternal → callExternal → ExternCalls[funcIdx] → DirectCall | reflect
func (v *vm) callExternal(funcIdx, numArgs int) error {
	// Pop arguments first
	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = v.pop()
	}

	// Method calls use the ExternalMethodInfo constant directly (not ExternCalls).
	// This is because method dispatch needs the ReceiverTypeName hint and
	// GlobalRef/*value.Value resolution, which are method-specific.
	if funcIdx < len(v.program.Constants) {
		if methodInfo, ok := v.program.Constants[funcIdx].(*external.ExternalMethodInfo); ok {
			return v.callExternalMethod(methodInfo, args)
		}
	}

	// Direct array lookup into pre-resolved table (no locks!)
	rc := v.externCall(funcIdx)
	if rc == nil {
		return fmt.Errorf("unresolved external call at constant index %d", funcIdx)
	}
	return v.callResolvedExternal(rc, args)
}

func (v *vm) callResolvedExternal(rc *bytecode.ResolvedCall, args []value.Value) error {
	if rc == nil {
		return fmt.Errorf("unresolved external call")
	}
	if err := v.validateExternalBoundary(rc, args); err != nil {
		return err
	}

	// Fast path: DirectCall available
	if rc.DirectCall != nil {
		if rc.IsVariadic && len(args) == rc.NumIn {
			args = unpackVariadicArgs(args, len(args))
		}
		if rc.FnType != nil {
			convertClosureArgs(args, rc.FnType)
		}
		v.push(rc.DirectCall(args))
		return v.checkCtx()
	}

	// Slow path: reflect.Call
	return v.callExternalReflect(rc, args)
}

// externCall returns the pre-resolved call entry for the given constant pool index.
// Falls back to resolving on-the-fly for indices outside the pre-resolved table
// (shouldn't happen in normal operation, but provides safety).
func (v *vm) externCall(funcIdx int) *bytecode.ResolvedCall {
	if funcIdx < len(v.program.ExternCalls) {
		if rc := v.program.ExternCalls[funcIdx]; rc != nil {
			return rc
		}
	}
	// Fallback: resolve from constant pool (rare — only if pre-resolution missed it)
	return resolveConstantFallback(v.program, funcIdx)
}

// resolveConstantFallback resolves a constant pool entry on-the-fly.
// This is a safety net for cases where ResolveExternCalls wasn't called
// or the entry wasn't pre-resolved.
func resolveConstantFallback(prog *bytecode.CompiledProgram, funcIdx int) *bytecode.ResolvedCall {
	if funcIdx >= len(prog.Constants) {
		return nil
	}
	// Re-use the same resolution logic as ResolveExternCalls
	rc := bytecode.ResolveConstant(prog.Constants[funcIdx])
	if rc != nil {
		// Patch the program-level table so future calls skip this path
		if funcIdx < len(prog.ExternCalls) {
			prog.ExternCalls[funcIdx] = rc
		}
	}
	return rc
}

// buildReflectArgs converts []value.Value args to []reflect.Value for reflect.Call,
// handling SSA-packed variadic slices. fnType is the target function type.
func buildReflectArgs(args []value.Value, fnType reflect.Type) []reflect.Value {
	numIn := fnType.NumIn()
	isVariadic := fnType.IsVariadic()
	numArgs := len(args)

	if isVariadic && numArgs == numIn {
		lastArg := args[numArgs-1]
		if rv, ok := lastArg.ReflectValue(); ok && rv.Kind() == reflect.Slice {
			sliceLen := rv.Len()
			in := make([]reflect.Value, numIn-1+sliceLen)
			for i := 0; i < numArgs-1; i++ {
				in[i] = args[i].ToReflectValue(fnType.In(i))
			}
			elemType := fnType.In(numIn - 1).Elem()
			for i := 0; i < sliceLen; i++ {
				elem := rv.Index(i)
				if elem.Kind() == reflect.Interface && !elem.IsNil() {
					elem = elem.Elem()
				}
				if elem.Type().ConvertibleTo(elemType) {
					in[numIn-1+i] = elem.Convert(elemType)
				} else {
					in[numIn-1+i] = elem
				}
			}
			return in
		}
		in := make([]reflect.Value, numArgs)
		for i, arg := range args {
			if i >= numIn-1 {
				variadicType := fnType.In(numIn - 1).Elem()
				in[i] = arg.ToReflectValue(variadicType)
			} else {
				in[i] = arg.ToReflectValue(fnType.In(i))
			}
		}
		return in
	}

	in := make([]reflect.Value, numArgs)
	for i, arg := range args {
		if i < numIn {
			in[i] = arg.ToReflectValue(fnType.In(i))
		} else if isVariadic {
			variadicType := fnType.In(numIn - 1).Elem()
			in[i] = arg.ToReflectValue(variadicType)
		}
	}
	return in
}

// callExternalReflect executes an external function using reflect.Call.
// This is the slow path when no DirectCall wrapper is available.
func (v *vm) callExternalReflect(rc *bytecode.ResolvedCall, args []value.Value) error {
	if rc == nil {
		return fmt.Errorf("unresolved external call")
	}
	if !rc.Fn.IsValid() || rc.Fn.Kind() != reflect.Func {
		return fmt.Errorf("invalid external function: missing reflect.Func")
	}

	in := buildReflectArgs(args, rc.FnType)
	out := rc.Fn.Call(in)

	if err := v.checkCtx(); err != nil {
		return err
	}
	v.pushReflectResults(out)
	return nil
}
