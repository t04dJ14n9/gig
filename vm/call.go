package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/bytecode"
	"github.com/t04dJ14n9/gig/value"
)

// ExternalCallCancelledError is returned when a context is cancelled before/after an external call.
// This wraps the original context error to provide more context.
type ExternalCallCancelledError struct {
	Cause error
}

func (e *ExternalCallCancelledError) Error() string {
	return "external call cancelled: " + e.Cause.Error()
}

func (e *ExternalCallCancelledError) Unwrap() error {
	return e.Cause
}

// unpackVariadicArgs unpacks the last argument (a packed variadic slice from SSA) into
// individual value.Value elements. This avoids reflection-based rv.Len()/rv.Index(i) calls
// for native slice types ([]value.Value, []int64, []byte).
func unpackVariadicArgs(args []value.Value, numArgs int) []value.Value {
	lastArg := args[numArgs-1]

	// Fast path 1: native []value.Value slice (function slices, etc.)
	if lastArg.Kind() == value.KindReflect || lastArg.Kind() == value.KindSlice {
		// Try native []value.Value first (zero reflection)
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
func (vm *VM) callCompiledFunction(funcIdx, numArgs int) {
	// O(1) function lookup via direct index table
	if funcIdx < 0 || funcIdx >= len(vm.program.FuncByIndex) {
		vm.push(value.MakeNil())
		return
	}
	fn := vm.program.FuncByIndex[funcIdx]
	if fn == nil {
		vm.push(value.MakeNil())
		return
	}

	// Get a pooled frame (avoids Frame + locals allocation)
	frame := vm.fpool.get(fn, vm.sp, nil)

	// Pop arguments directly into the frame's locals (avoids temporary args slice)
	intL := frame.intLocals
	for i := numArgs - 1; i >= 0; i-- {
		if i < fn.NumLocals {
			v := vm.pop()
			frame.locals[i] = v
			// Mirror int parameters into intLocals for OpInt* opcodes
			if intL != nil {
				intL[i] = v.RawInt()
			}
		} else {
			vm.pop()
		}
	}

	vm.frames[vm.fp] = frame
	vm.fp++
}

// callFunction calls a function with the given arguments and free variables.
// Used for calling closures.
func (vm *VM) callFunction(fn *bytecode.CompiledFunction, args []value.Value, freeVars []*value.Value) {
	frame := vm.fpool.get(fn, vm.sp, freeVars)
	// Copy arguments into locals
	for i, arg := range args {
		if i < fn.NumLocals {
			frame.locals[i] = arg
			// Mirror int parameters into intLocals for OpInt* opcodes
			if frame.intLocals != nil {
				frame.intLocals[i] = arg.RawInt()
			}
		}
	}
	vm.frames[vm.fp] = frame
	vm.fp++
}

// callExternal calls an external (native Go) function.
// It first checks the inline cache, then resolves the function if not cached.
// Supports both DirectCall (fast path) and reflect.Call (slow path).
// Returns an error if the context is cancelled during or immediately after the call.
func (vm *VM) callExternal(funcIdx, numArgs int) error {
	// Pop arguments first (before any cache lookup)
	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = vm.pop()
	}

	// Check if this is a method call (ExternalMethodInfo)
	if funcIdx < len(vm.program.Constants) {
		if methodInfo, ok := vm.program.Constants[funcIdx].(*bytecode.ExternalMethodInfo); ok {
			return vm.callExternalMethod(methodInfo, args)
		}
	}

	// Check inline cache with lock for concurrent safety
	var cacheEntry *extCallCacheEntry
	vm.extCallCache.mu.Lock()
	if funcIdx < len(vm.extCallCache.cache) {
		cacheEntry = vm.extCallCache.cache[funcIdx]
		if cacheEntry == nil {
			// Resolve and cache while holding the lock
			cacheEntry = vm.resolveExternalFunc(funcIdx)
			vm.extCallCache.cache[funcIdx] = cacheEntry
		}
	}
	vm.extCallCache.mu.Unlock()

	// Fast path: DirectCall available
	if cacheEntry.directCall != nil {
		// For variadic functions, SSA packs variadic args into a slice.
		// We need to unpack them for DirectCall wrappers.
		if cacheEntry.isVariadic && numArgs == cacheEntry.numIn {
			args = unpackVariadicArgs(args, numArgs)
		}
		result := cacheEntry.directCall(args)
		vm.push(result)
		// Check context after the call for post-call cancellation
		select {
		case <-vm.ctx.Done():
			return &ExternalCallCancelledError{Cause: vm.ctx.Err()}
		default:
		}
		return nil
	}

	// Slow path: use reflect.Call
	return vm.callExternalReflect(cacheEntry, args)
}

// resolveExternalFunc resolves an external function from the constant pool.
// It creates a cache entry for future calls.
func (vm *VM) resolveExternalFunc(funcIdx int) *extCallCacheEntry {
	entry := &extCallCacheEntry{}

	// Check if constant is ExternalFuncInfo (new optimized path)
	if funcIdx < len(vm.program.Constants) {
		if extInfo, ok := vm.program.Constants[funcIdx].(*bytecode.ExternalFuncInfo); ok {
			entry.directCall = extInfo.DirectCall
			if extInfo.Func != nil {
				entry.fn = reflect.ValueOf(extInfo.Func)
				// Check if it's actually a function (not an *ssa.Function or other type)
				if entry.fn.Kind() == reflect.Func {
					entry.fnType = entry.fn.Type()
					entry.isVariadic = entry.fnType.IsVariadic()
					entry.numIn = entry.fnType.NumIn()
				}
			}
			return entry
		}
		// Fallback: old-style constant (just the function value)
		extFunc := vm.program.Constants[funcIdx]
		if extFunc != nil {
			entry.fn = reflect.ValueOf(extFunc)
			if entry.fn.Kind() == reflect.Func {
				entry.fnType = entry.fn.Type()
				entry.isVariadic = entry.fnType.IsVariadic()
				entry.numIn = entry.fnType.NumIn()
			}
		}
	}

	return entry
}

// callExternalReflect executes an external function using reflect.Call.
// This is the slow path when no DirectCall wrapper is available.
func (vm *VM) callExternalReflect(entry *extCallCacheEntry, args []value.Value) error {
	if !entry.fn.IsValid() || entry.fn.Kind() != reflect.Func {
		vm.push(value.MakeNil())
		return nil
	}

	numArgs := len(args)

	// Build reflect.Value arguments
	var in []reflect.Value

	// For variadic calls where SSA passes the variadic slice as the last arg,
	// we need to unpack it for reflect.Call
	if entry.isVariadic && numArgs == entry.numIn {
		// The last arg might be the variadic slice packed by SSA
		lastArg := args[numArgs-1]
		if rv, ok := lastArg.ReflectValue(); ok && rv.Kind() == reflect.Slice {
			// Unpack: use first N-1 args normally, then spread the slice elements
			sliceLen := rv.Len()
			in = make([]reflect.Value, entry.numIn-1+sliceLen)
			for i := 0; i < numArgs-1; i++ {
				in[i] = args[i].ToReflectValue(entry.fnType.In(i))
			}
			elemType := entry.fnType.In(entry.numIn - 1).Elem()
			for i := 0; i < sliceLen; i++ {
				elem := rv.Index(i)
				// If elem is interface{}, unwrap it
				if elem.Kind() == reflect.Interface && !elem.IsNil() {
					elem = elem.Elem()
				}
				if elem.Type().ConvertibleTo(elemType) {
					in[entry.numIn-1+i] = elem.Convert(elemType)
				} else {
					in[entry.numIn-1+i] = elem
				}
			}
		} else {
			// Last arg is not a slice, treat normally
			in = make([]reflect.Value, numArgs)
			for i, arg := range args {
				if i >= entry.numIn-1 {
					variadicType := entry.fnType.In(entry.numIn - 1).Elem()
					in[i] = arg.ToReflectValue(variadicType)
				} else {
					in[i] = arg.ToReflectValue(entry.fnType.In(i))
				}
			}
		}
	} else {
		in = make([]reflect.Value, numArgs)
		for i, arg := range args {
			if i < entry.numIn {
				in[i] = arg.ToReflectValue(entry.fnType.In(i))
			}
		}
	}

	// Call the function
	out := entry.fn.Call(in)

	// Check context after the call for post-call cancellation
	select {
	case <-vm.ctx.Done():
		return &ExternalCallCancelledError{Cause: vm.ctx.Err()}
	default:
	}

	// Convert result
	if len(out) == 0 {
		vm.push(value.MakeNil())
	} else if len(out) == 1 {
		vm.push(value.MakeFromReflect(out[0]))
	} else {
		// Multiple return values - pack as slice
		results := make([]value.Value, len(out))
		for i, v := range out {
			results[i] = value.MakeFromReflect(v)
		}
		vm.push(value.FromInterface(results))
	}
	return nil
}

// callExternalMethod dispatches a method call on an external type.
// args[0] is the receiver, args[1:] are the method arguments.
// Uses DirectCall fast path if available, otherwise falls back to reflection.
func (vm *VM) callExternalMethod(methodInfo *bytecode.ExternalMethodInfo, args []value.Value) error {
	if len(args) == 0 {
		vm.push(value.MakeNil())
		return nil
	}

	// Fast path: DirectCall wrapper resolved at compile time
	if methodInfo.DirectCall != nil {
		result := methodInfo.DirectCall(args)
		vm.push(result)
		// Check context after the call
		select {
		case <-vm.ctx.Done():
			return &ExternalCallCancelledError{Cause: vm.ctx.Err()}
		default:
		}
		return nil
	}

	// Slow path: use reflect.MethodByName + reflect.Call
	return vm.callExternalMethodReflect(methodInfo, args)
}

// callExternalMethodReflect dispatches a method call using reflection.
func (vm *VM) callExternalMethodReflect(methodInfo *bytecode.ExternalMethodInfo, args []value.Value) error {
	// Get the receiver as a reflect.Value
	receiver := args[0]
	var rv reflect.Value
	if reflectVal, ok := receiver.ReflectValue(); ok {
		rv = reflectVal
	} else {
		iface := receiver.Interface()
		if iface == nil {
			vm.push(value.MakeNil())
			return nil
		}
		rv = reflect.ValueOf(iface)
	}

	if !rv.IsValid() {
		vm.push(value.MakeNil())
		return nil
	}

	// For interface method dispatch: if the receiver is an interface, unwrap it
	// to get the concrete type, then look up the method on the concrete type.
	if rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			vm.push(value.MakeNil())
			return nil
		}
		rv = rv.Elem()
	}

	// Look up the method by name
	method := rv.MethodByName(methodInfo.MethodName)
	if !method.IsValid() {
		// Try pointer receiver
		if rv.CanAddr() {
			method = rv.Addr().MethodByName(methodInfo.MethodName)
		}
		if !method.IsValid() {
			// Reflection-based lookup failed. This happens when typeToReflect
			// strips named types to anonymous structs (e.g., Impl → struct{val int}).
			// Fall back to calling a compiled method from the function table.
			return vm.callCompiledMethod(methodInfo.MethodName, args)
		}
	}

	// Build arguments (skip the receiver at args[0])
	methodType := method.Type()
	numIn := methodType.NumIn()
	isVariadic := methodType.IsVariadic()
	methodArgs := args[1:]

	var in []reflect.Value

	if isVariadic && len(methodArgs) == numIn {
		// Check if the last arg is a packed variadic slice from SSA
		lastArg := methodArgs[len(methodArgs)-1]
		if lastRV, ok := lastArg.ReflectValue(); ok && lastRV.Kind() == reflect.Slice {
			sliceLen := lastRV.Len()
			in = make([]reflect.Value, numIn-1+sliceLen)
			for i := 0; i < len(methodArgs)-1; i++ {
				in[i] = methodArgs[i].ToReflectValue(methodType.In(i))
			}
			elemType := methodType.In(numIn - 1).Elem()
			for i := 0; i < sliceLen; i++ {
				elem := lastRV.Index(i)
				if elem.Kind() == reflect.Interface && !elem.IsNil() {
					elem = elem.Elem()
				}
				if elem.Type().ConvertibleTo(elemType) {
					in[numIn-1+i] = elem.Convert(elemType)
				} else {
					in[numIn-1+i] = elem
				}
			}
		} else {
			in = make([]reflect.Value, len(methodArgs))
			for i, arg := range methodArgs {
				if i < numIn {
					in[i] = arg.ToReflectValue(methodType.In(i))
				}
			}
		}
	} else {
		in = make([]reflect.Value, len(methodArgs))
		for i, arg := range methodArgs {
			if i < numIn {
				in[i] = arg.ToReflectValue(methodType.In(i))
			} else if isVariadic {
				variadicType := methodType.In(numIn - 1).Elem()
				in[i] = arg.ToReflectValue(variadicType)
			}
		}
	}

	// Call the method
	out := method.Call(in)

	// Check context after the call for post-call cancellation
	select {
	case <-vm.ctx.Done():
		return &ExternalCallCancelledError{Cause: vm.ctx.Err()}
	default:
	}

	// Convert result
	if len(out) == 0 {
		vm.push(value.MakeNil())
	} else if len(out) == 1 {
		vm.push(value.MakeFromReflect(out[0]))
	} else {
		results := make([]value.Value, len(out))
		for i, v := range out {
			results[i] = value.MakeFromReflect(v)
		}
		vm.push(value.FromInterface(results))
	}
	return nil
}

// callCompiledMethod searches the compiled function table for a method with the
// given name and calls it. This is the fallback path for invoke (interface method)
// calls when reflection-based MethodByName fails — typically because typeToReflect
// converts named types to anonymous struct types that lack the named type's methods.
//
// args[0] is the receiver, args[1:] are the method arguments.
func (vm *VM) callCompiledMethod(methodName string, args []value.Value) error {
	// Search for a compiled function whose SSA name matches the method name
	// and that has a receiver (i.e., is actually a method, not a plain function).
	for _, fn := range vm.program.FuncByIndex {
		if fn == nil || fn.Source == nil {
			continue
		}
		if fn.Source.Name() != methodName {
			continue
		}
		sig := fn.Source.Signature
		if sig.Recv() == nil {
			continue
		}
		// Found a matching compiled method — call it with args as a compiled function.
		// Push args onto the stack (receiver first, then method args).
		for _, arg := range args {
			vm.push(arg)
		}
		vm.callCompiledFunction(vm.program.FuncIndex[fn.Source], len(args))
		return nil
	}

	// No compiled method found — push nil
	vm.push(value.MakeNil())
	return nil
}
