package vm

import (
	"reflect"

	"gig/bytecode"
	"gig/value"
)

// callCompiledFunction calls a compiled function by its index.
// It creates a new call frame with the function's local variables.
func (vm *VM) callCompiledFunction(funcIdx, numArgs int) {
	// Get function
	var fn *bytecode.CompiledFunction
	for _, f := range vm.program.Functions {
		if vm.program.FuncIndex[f.Source] == funcIdx {
			fn = f
			break
		}
	}
	if fn == nil {
		vm.push(value.MakeNil())
		return
	}

	// Pop arguments
	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = vm.pop()
	}

	// Create new frame
	frame := newFrame(fn, vm.sp, args, nil)
	vm.frames[vm.fp] = frame
	vm.fp++
}

// callFunction calls a function with the given arguments and free variables.
// Used for calling closures.
func (vm *VM) callFunction(fn *bytecode.CompiledFunction, args []value.Value, freeVars []*value.Value) {
	frame := newFrame(fn, vm.sp, args, freeVars)
	vm.frames[vm.fp] = frame
	vm.fp++
}

// callExternal calls an external (native Go) function.
// It first checks the inline cache, then resolves the function if not cached.
// Supports both DirectCall (fast path) and reflect.Call (slow path).
func (vm *VM) callExternal(funcIdx, numArgs int) {
	// Pop arguments first (before any cache lookup)
	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = vm.pop()
	}

	// Check if this is a method call (ExternalMethodInfo)
	if funcIdx < len(vm.program.Constants) {
		if methodInfo, ok := vm.program.Constants[funcIdx].(*bytecode.ExternalMethodInfo); ok {
			vm.callExternalMethod(methodInfo, args)
			return
		}
	}

	// Check inline cache
	cachedEntry, cached := vm.extCallCache.Load(funcIdx)
	if !cached {
		// Resolve and cache
		cacheEntry := vm.resolveExternalFunc(funcIdx)
		vm.extCallCache.Store(funcIdx, cacheEntry)
		cachedEntry = cacheEntry
	}
	cacheEntry := cachedEntry.(*extCallCacheEntry)

	// Fast path: DirectCall available
	if cacheEntry.directCall != nil {
		// For variadic functions, SSA packs variadic args into a slice.
		// We need to unpack them for DirectCall wrappers.
		if cacheEntry.isVariadic && numArgs == cacheEntry.numIn {
			lastArg := args[numArgs-1]
			if rv, ok := lastArg.ReflectValue(); ok && rv.Kind() == reflect.Slice {
				// Unpack the variadic slice
				sliceLen := rv.Len()
				unpackedArgs := make([]value.Value, numArgs-1+sliceLen)
				copy(unpackedArgs, args[:numArgs-1])
				for i := 0; i < sliceLen; i++ {
					unpackedArgs[numArgs-1+i] = value.MakeFromReflect(rv.Index(i))
				}
				args = unpackedArgs
			}
		}
		result := cacheEntry.directCall(args)
		vm.push(result)
		return
	}

	// Slow path: use reflect.Call
	vm.callExternalReflect(cacheEntry, args)
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
				entry.fnType = entry.fn.Type()
				entry.isVariadic = entry.fnType.IsVariadic()
				entry.numIn = entry.fnType.NumIn()
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
func (vm *VM) callExternalReflect(entry *extCallCacheEntry, args []value.Value) {
	if !entry.fn.IsValid() || entry.fn.Kind() != reflect.Func {
		vm.push(value.MakeNil())
		return
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
}

// callExternalMethod dispatches a method call on an external type using reflection.
// args[0] is the receiver, args[1:] are the method arguments.
func (vm *VM) callExternalMethod(methodInfo *bytecode.ExternalMethodInfo, args []value.Value) {
	if len(args) == 0 {
		vm.push(value.MakeNil())
		return
	}

	// Get the receiver as a reflect.Value
	receiver := args[0]
	var rv reflect.Value
	if reflectVal, ok := receiver.ReflectValue(); ok {
		rv = reflectVal
	} else {
		rv = reflect.ValueOf(receiver.Interface())
	}

	if !rv.IsValid() {
		vm.push(value.MakeNil())
		return
	}

	// Look up the method by name
	method := rv.MethodByName(methodInfo.MethodName)
	if !method.IsValid() {
		// Try pointer receiver
		if rv.CanAddr() {
			method = rv.Addr().MethodByName(methodInfo.MethodName)
		}
		if !method.IsValid() {
			vm.push(value.MakeNil())
			return
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
}
