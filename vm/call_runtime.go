package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// ExternalCallCancelledError is returned when a context is cancelled before/after an external call.
type ExternalCallCancelledError struct {
	Cause error
}

func (e *ExternalCallCancelledError) Error() string {
	return "external call cancelled: " + e.Cause.Error()
}

func (e *ExternalCallCancelledError) Unwrap() error {
	return e.Cause
}

// checkCtx returns an ExternalCallCancelledError if the VM's context is done.
// Used by all external call paths to avoid duplicating the select block.
func (v *vm) checkCtx() error {
	select {
	case <-v.ctx.Done():
		return &ExternalCallCancelledError{Cause: v.ctx.Err()}
	default:
		return nil
	}
}

// pushReflectResults converts reflect.Call output and pushes it onto the VM stack.
// Handles 0, 1, or multiple return values. Used by all reflect-based call paths.
func (v *vm) pushReflectResults(out []reflect.Value) {
	switch len(out) {
	case 0:
		v.push(value.MakeNil())
	case 1:
		v.push(value.MakeFromReflect(out[0]))
	default:
		results := make([]value.Value, len(out))
		for i, val := range out {
			results[i] = value.MakeFromReflect(val)
		}
		v.push(value.FromInterface(results))
	}
}

// convertClosureArgsForMethod converts interpreted closure arguments to real Go
// functions for a method DirectCall. It reflects on the receiver (args[0]) to
// look up the method signature and wraps any KindFunc arguments via reflect.MakeFunc.
// args[0] is the receiver, args[1:] are method arguments.
func convertClosureArgsForMethod(methodName string, args []value.Value) {
	hasClosure := false
	for i := 1; i < len(args); i++ {
		if args[i].Kind() == value.KindFunc {
			hasClosure = true
			break
		}
	}
	if !hasClosure {
		return
	}

	// Get the receiver's reflect type to look up method signature
	var rv reflect.Value
	if rr, ok := args[0].ReflectValue(); ok {
		rv = rr
	} else {
		iface := args[0].Interface()
		if iface == nil {
			return
		}
		rv = reflect.ValueOf(iface)
	}

	// Look up the method by name on the receiver
	method := rv.MethodByName(methodName)
	if !method.IsValid() {
		if rv.CanAddr() {
			method = rv.Addr().MethodByName(methodName)
		}
		if !method.IsValid() {
			return
		}
	}

	mt := method.Type()
	for i := 1; i < len(args); i++ {
		if args[i].Kind() == value.KindFunc {
			paramIdx := i - 1 // method args are 0-based (no receiver in method.Type())
			if paramIdx < mt.NumIn() {
				paramType := mt.In(paramIdx)
				if paramType.Kind() == reflect.Func {
					args[i] = value.MakeFromReflect(args[i].ToReflectValue(paramType))
				}
			}
		}
	}
}

// convertClosureArgs scans args for interpreted closures (KindFunc) and converts them
// to real Go functions using reflect.MakeFunc. This allows DirectCall wrappers to receive
// proper Go function values instead of *v.Closure pointers.
func convertClosureArgs(args []value.Value, fnType reflect.Type) {
	numIn := fnType.NumIn()
	for i, arg := range args {
		if arg.Kind() == value.KindFunc {
			var targetType reflect.Type
			if i < numIn {
				targetType = fnType.In(i)
			} else if fnType.IsVariadic() && numIn > 0 {
				targetType = fnType.In(numIn - 1).Elem()
			}
			if targetType != nil && targetType.Kind() == reflect.Func {
				args[i] = value.MakeFromReflect(arg.ToReflectValue(targetType))
			}
		}
	}
}

// unpackVariadicArgs unpacks the last argument (a packed variadic slice from SSA) into
// individual value.Value elements. This avoids reflection-based rv.Len()/rv.Index(i) calls
// for native slice types ([]value.Value, []int64, []byte).
func unpackVariadicArgs(args []value.Value, numArgs int) []value.Value {
	lastArg := args[numArgs-1]

	// Fast path 1: native []value.Value slice (function slices, etc.)
	if lastArg.Kind() == value.KindReflect || lastArg.Kind() == value.KindSlice {
		if rawObj := lastArg.RawObj(); rawObj != nil {
			if valSlice, ok := rawObj.([]value.Value); ok {
				unpackedArgs := make([]value.Value, numArgs-1+len(valSlice))
				copy(unpackedArgs, args[:numArgs-1])
				copy(unpackedArgs[numArgs-1:], valSlice)
				return unpackedArgs
			}
		}
	}

	// Fast path 2: native []int64 slice (int variadic args)
	if lastArg.Kind() == value.KindSlice {
		if intSlice, ok := lastArg.IntSlice(); ok {
			unpackedArgs := make([]value.Value, numArgs-1+len(intSlice))
			copy(unpackedArgs, args[:numArgs-1])
			for i, n := range intSlice {
				unpackedArgs[numArgs-1+i] = value.MakeInt(n)
			}
			return unpackedArgs
		}
	}

	// Fast path 3: KindBytes ([]byte variadic args)
	if lastArg.Kind() == value.KindBytes {
		if b, ok := lastArg.Bytes(); ok {
			unpackedArgs := make([]value.Value, numArgs-1+len(b))
			copy(unpackedArgs, args[:numArgs-1])
			for i, byt := range b {
				unpackedArgs[numArgs-1+i] = value.MakeUint(uint64(byt))
			}
			return unpackedArgs
		}
	}

	// Slow path: reflect.Value slice (most common for stdlib variadic functions)
	if rv, ok := lastArg.ReflectValue(); ok && rv.Kind() == reflect.Slice {
		sliceLen := rv.Len()
		unpackedArgs := make([]value.Value, numArgs-1+sliceLen)
		copy(unpackedArgs, args[:numArgs-1])
		for i := 0; i < sliceLen; i++ {
			unpackedArgs[numArgs-1+i] = value.MakeFromReflect(rv.Index(i))
		}
		return unpackedArgs
	}

	return args
}

// callCompiledFunction calls a compiled function by its index.
// It creates a new call frame with the function's local variables.
func (v *vm) callCompiledFunction(funcIdx, numArgs int) {
	if funcIdx < 0 || funcIdx >= len(v.program.FuncByIndex) {
		v.push(value.MakeNil())
		return
	}
	fn := v.program.FuncByIndex[funcIdx]
	if fn == nil {
		v.push(value.MakeNil())
		return
	}

	frame := v.fpool.get(fn, v.sp, nil)

	intL := frame.intLocals
	for i := numArgs - 1; i >= 0; i-- {
		if i < fn.NumLocals {
			val := v.pop()
			frame.locals[i] = val
			if intL != nil {
				intL[i] = val.RawInt()
			}
		} else {
			v.pop()
		}
	}

	if v.fp >= len(v.frames) {
		if !v.growFrames() {
			panic("gig: call stack overflow")
		}
	}
	v.frames[v.fp] = frame
	v.fp++
}

// callFunction calls a function with the given arguments and free variables.
// Used for calling closures.
func (v *vm) callFunction(fn *bytecode.CompiledFunction, args []value.Value, freeVars []*value.Value) {
	frame := v.fpool.get(fn, v.sp, freeVars)
	for i, arg := range args {
		if i < fn.NumLocals {
			frame.locals[i] = arg
			if frame.intLocals != nil {
				frame.intLocals[i] = arg.RawInt()
			}
		}
	}
	if v.fp >= len(v.frames) {
		if !v.growFrames() {
			panic("gig: call stack overflow")
		}
	}
	v.frames[v.fp] = frame
	v.fp++
}
