// call.go handles external function calls (DirectCall + reflect), method dispatch,
// and variadic argument unpacking.
package vm

import (
	"reflect"
	"strings"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/external"
	"git.woa.com/youngjin/gig/model/value"
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
			// Determine the target function type from the external function's signature
			var targetType reflect.Type
			if i < numIn {
				targetType = fnType.In(i)
			} else if fnType.IsVariadic() && numIn > 0 {
				// For variadic args, use the element type of the variadic slice
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
func (v *vm) callCompiledFunction(funcIdx, numArgs int) {
	// O(1) function lookup via direct index table
	if funcIdx < 0 || funcIdx >= len(v.program.FuncByIndex) {
		v.push(value.MakeNil())
		return
	}
	fn := v.program.FuncByIndex[funcIdx]
	if fn == nil {
		v.push(value.MakeNil())
		return
	}

	// Get a pooled frame (avoids Frame + locals allocation)
	frame := v.fpool.get(fn, v.sp, nil)

	// Pop arguments directly into the frame's locals (avoids temporary args slice)
	intL := frame.intLocals
	for i := numArgs - 1; i >= 0; i-- {
		if i < fn.NumLocals {
			val := v.pop()
			frame.locals[i] = val
			// Mirror int parameters into intLocals for OpInt* opcodes
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
	if v.fp >= len(v.frames) {
		if !v.growFrames() {
			// Stack overflow — trigger a panic that the safety net will catch.
			panic("gig: call stack overflow")
		}
	}
	v.frames[v.fp] = frame
	v.fp++
}

// callExternal calls an external (native Go) function.
// It first checks the inline cache, then resolves the function if not cached.
// Supports both DirectCall (fast path) and reflect.Call (slow path).
// Returns an error if the context is cancelled during or immediately after the call.
func (v *vm) callExternal(funcIdx, numArgs int) error {
	// Pop arguments first (before any cache lookup)
	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = v.pop()
	}

	// Check if this is a method call (ExternalMethodInfo)
	if funcIdx < len(v.program.Constants) {
		if methodInfo, ok := v.program.Constants[funcIdx].(*external.ExternalMethodInfo); ok {
			return v.callExternalMethod(methodInfo, args)
		}
	}

	// Fast path: read-lock check (common case — cache already populated)
	var cacheEntry *extCallCacheEntry
	v.extCallCache.mu.RLock()
	if funcIdx < len(v.extCallCache.cache) {
		cacheEntry = v.extCallCache.cache[funcIdx]
	}
	v.extCallCache.mu.RUnlock()

	if cacheEntry == nil {
		// Slow path: write-lock to populate the entry (double-checked)
		v.extCallCache.mu.Lock()
		if funcIdx < len(v.extCallCache.cache) {
			cacheEntry = v.extCallCache.cache[funcIdx]
			if cacheEntry == nil {
				cacheEntry = v.resolveExternalFunc(funcIdx)
				v.extCallCache.cache[funcIdx] = cacheEntry
			}
		}
		v.extCallCache.mu.Unlock()
	}

	// Fast path: DirectCall available
	if cacheEntry.directCall != nil {
		// For variadic functions, SSA packs variadic args into a slice.
		// We need to unpack them for DirectCall wrappers.
		if cacheEntry.isVariadic && numArgs == cacheEntry.numIn {
			args = unpackVariadicArgs(args, numArgs)
		}
		// Convert interpreted closures to real Go functions before calling DirectCall.
		// DirectCall wrappers use Interface() which returns *v.Closure, not a Go func.
		if cacheEntry.fnType != nil {
			convertClosureArgs(args, cacheEntry.fnType)
		}
		result := cacheEntry.directCall(args)
		v.push(result)
		// Check context after the call for post-call cancellation
		select {
		case <-v.ctx.Done():
			return &ExternalCallCancelledError{Cause: v.ctx.Err()}
		default:
		}
		return nil
	}

	// Slow path: use reflect.Call
	return v.callExternalReflect(cacheEntry, args)
}

// resolveExternalFunc resolves an external function from the constant pool.
// It creates a cache entry for future calls.
func (v *vm) resolveExternalFunc(funcIdx int) *extCallCacheEntry {
	entry := &extCallCacheEntry{}

	// Check if constant is ExternalFuncInfo (new optimized path)
	if funcIdx < len(v.program.Constants) {
		if extInfo, ok := v.program.Constants[funcIdx].(*external.ExternalFuncInfo); ok {
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
		extFunc := v.program.Constants[funcIdx]
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

// buildReflectArgs converts []value.Value args to []reflect.Value for reflect.Call,
// handling SSA-packed variadic slices. fnType is the target function type.
func buildReflectArgs(args []value.Value, fnType reflect.Type) []reflect.Value {
	numIn := fnType.NumIn()
	isVariadic := fnType.IsVariadic()
	numArgs := len(args)

	if isVariadic && numArgs == numIn {
		// The last arg might be the variadic slice packed by SSA
		lastArg := args[numArgs-1]
		if rv, ok := lastArg.ReflectValue(); ok && rv.Kind() == reflect.Slice {
			// Unpack: use first N-1 args normally, then spread the slice elements
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
		// Last arg is not a slice, treat normally
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
func (v *vm) callExternalReflect(entry *extCallCacheEntry, args []value.Value) error {
	if !entry.fn.IsValid() || entry.fn.Kind() != reflect.Func {
		v.push(value.MakeNil())
		return nil
	}

	in := buildReflectArgs(args, entry.fnType)

	// Call the function
	out := entry.fn.Call(in)

	// Check context after the call for post-call cancellation
	select {
	case <-v.ctx.Done():
		return &ExternalCallCancelledError{Cause: v.ctx.Err()}
	default:
	}

	// Convert result
	if len(out) == 0 {
		v.push(value.MakeNil())
	} else if len(out) == 1 {
		v.push(value.MakeFromReflect(out[0]))
	} else {
		// Multiple return values - pack as slice
		results := make([]value.Value, len(out))
		for i, val := range out {
			results[i] = value.MakeFromReflect(val)
		}
		v.push(value.FromInterface(results))
	}
	return nil
}

// callExternalMethod dispatches a method call on an external type.
// args[0] is the receiver, args[1:] are the method arguments.
// Uses DirectCall fast path if available, otherwise falls back to reflection.
func (v *vm) callExternalMethod(methodInfo *external.ExternalMethodInfo, args []value.Value) error {
	if len(args) == 0 {
		v.push(value.MakeNil())
		return nil
	}

	// Fast path: DirectCall wrapper resolved at compile time
	if methodInfo.DirectCall != nil {
		// Convert interpreted closures in method args to real Go functions.
		convertClosureArgsForMethod(methodInfo.MethodName, args)
		result := methodInfo.DirectCall(args)
		v.push(result)
		// Check context after the call
		select {
		case <-v.ctx.Done():
			return &ExternalCallCancelledError{Cause: v.ctx.Err()}
		default:
		}
		return nil
	}

	// Slow path: use reflect.MethodByName + reflect.Call
	return v.callExternalMethodReflect(methodInfo, args)
}

// callExternalMethodReflect dispatches a method call using reflection.
func (v *vm) callExternalMethodReflect(methodInfo *external.ExternalMethodInfo, args []value.Value) error {
	// Get the receiver as a reflect.Value
	receiver := args[0]
	var rv reflect.Value
	if reflectVal, ok := receiver.ReflectValue(); ok {
		rv = reflectVal
	} else {
		iface := receiver.Interface()
		if iface == nil {
			v.push(value.MakeNil())
			return nil
		}
		rv = reflect.ValueOf(iface)
	}

	if !rv.IsValid() {
		v.push(value.MakeNil())
		return nil
	}

	// For interface method dispatch: if the receiver is an interface, unwrap it
	// to get the concrete type, then look up the method on the concrete type.
	if rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			v.push(value.MakeNil())
			return nil
		}
		rv = rv.Elem()
		// Update args[0] with the unwrapped value for callCompiledMethod
		args[0] = value.MakeFromReflect(rv)
	}

	// Look up the method by name
	method := rv.MethodByName(methodInfo.MethodName)
	if !method.IsValid() {
		// Try pointer receiver
		if rv.CanAddr() {
			method = rv.Addr().MethodByName(methodInfo.MethodName)
		}
		if !method.IsValid() {
			// For structs with embedded interface fields (e.g., GetterHolder{Getter}),
			// check if any field is an interface that contains the method.
			// This handles the case where typeToReflect converts interfaces to interface{}
			// and the method is actually on the concrete value stored in the interface field.
			if rv.Kind() == reflect.Struct {
				for i := 0; i < rv.NumField(); i++ {
					field := rv.Field(i)
					if field.Kind() == reflect.Interface && !field.IsNil() {
						concrete := field.Elem()
						m := concrete.MethodByName(methodInfo.MethodName)
						if m.IsValid() {
							// Found the method on the embedded interface's concrete value.
							// Replace receiver with the concrete value and dispatch.
							args[0] = value.MakeFromReflect(concrete)
							method = m
							break
						}
						// Also try pointer receiver on the concrete value
						if concrete.CanAddr() {
							m = concrete.Addr().MethodByName(methodInfo.MethodName)
							if m.IsValid() {
								args[0] = value.MakeFromReflect(concrete.Addr())
								method = m
								break
							}
						}
					}
				}
			}
		}
		if !method.IsValid() {
			// Reflection-based lookup failed. This happens when typeToReflect
			// strips named types to anonymous structs (e.g., Impl → struct{val int}).
			// Fall back to calling a compiled method from the function table.
			return v.callCompiledMethod(methodInfo.MethodName, methodInfo.ReceiverTypeName, args)
		}
	}

	// Build arguments (skip the receiver at args[0])
	methodType := method.Type()
	methodArgs := args[1:]

	in := buildReflectArgs(methodArgs, methodType)

	// Call the method
	out := method.Call(in)

	// Check context after the call for post-call cancellation
	select {
	case <-v.ctx.Done():
		return &ExternalCallCancelledError{Cause: v.ctx.Err()}
	default:
	}

	// Convert result
	if len(out) == 0 {
		v.push(value.MakeNil())
	} else if len(out) == 1 {
		v.push(value.MakeFromReflect(out[0]))
	} else {
		results := make([]value.Value, len(out))
		for i, v := range out {
			results[i] = value.MakeFromReflect(v)
		}
		v.push(value.FromInterface(results))
	}
	return nil
}

// callCompiledMethod searches the compiled function table for a method with the
// given name and calls it. This is the fallback path for invoke (interface method)
// calls when reflection-based MethodByName fails — typically because typeToReflect
// converts named types to anonymous struct types that lack the named type's methods.
//
// receiverTypeName is an optional hint from ExternalMethodInfo.ReceiverTypeName.
// When non-empty, it restricts matching to methods whose receiver type name matches.
// This prevents mis-dispatch when multiple types define methods with the same name
// (e.g., both GetterImpl.Get and AdderStruct.Add, or multiple types with Get).
//
// args[0] is the receiver, args[1:] are the method arguments.
func (v *vm) callCompiledMethod(methodName string, receiverTypeName string, args []value.Value) error {
	// Use MethodsByName index for O(k) lookup instead of O(n) scan.
	candidates := v.program.MethodsByName[methodName]

	// If we have a receiver type hint, first try to match both name and receiver type.
	if receiverTypeName != "" {
		for _, fn := range candidates {
			if fn.ReceiverTypeName == receiverTypeName {
				for _, arg := range args {
					v.push(arg)
				}
				v.callCompiledFunction(fn.FuncIdx, len(args))
				return nil
			}
		}
	}

	// Fallback: try to infer receiver type from the actual runtime value in args[0].
	if len(args) > 0 {
		if concreteTypeName := inferReceiverTypeName(args[0], v.program); concreteTypeName != "" {
			for _, fn := range candidates {
				if fn.ReceiverTypeName == concreteTypeName {
					for _, arg := range args {
						v.push(arg)
					}
					v.callCompiledFunction(fn.FuncIdx, len(args))
					return nil
				}
			}
		}
	}

	// Last resort: match by method name only (first candidate).
	if len(candidates) > 0 {
		fn := candidates[0]
		for _, arg := range args {
			v.push(arg)
		}
		v.callCompiledFunction(fn.FuncIdx, len(args))
		return nil
	}

	// No compiled method found — push nil
	v.push(value.MakeNil())
	return nil
}

// inferReceiverTypeName tries to extract a type name from a runtime value.Value
// receiver. This is used when callCompiledMethod doesn't have a static type hint
// but needs to disambiguate by the actual value being dispatched on.
// It first checks the program-level ReflectTypeNames registry, then falls back
// to scanning field PkgPath suffixes for unexported fields.
func inferReceiverTypeName(receiver value.Value, prog *bytecode.CompiledProgram) string {
	rv, ok := receiver.ReflectValue()
	if !ok {
		return ""
	}
	// Unwrap interface
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	// Unwrap pointer
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	t := rv.Type()
	// For real Go types passed through, the name is available directly.
	if t.Name() != "" {
		return t.Name()
	}
	// Check the program-level ReflectTypeNames registry (new approach).
	if prog != nil {
		if name := prog.LookupTypeName(t); name != "" {
			return name
		}
	}
	// Fallback: for synthesized struct types (from reflect.StructOf via typeToReflect),
	// the type name may be embedded in the field PkgPath as a "#PkgName.TypeName" suffix
	// (or "#TypeName"). Extract just the bare TypeName to match against the
	// *types.Named receiver type name.
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			pkgPath := f.PkgPath
			if idx := strings.LastIndex(pkgPath, "#"); idx >= 0 {
				qualName := pkgPath[idx+1:]
				// Strip package prefix if present (e.g., "thirdparty.uppercaseProcessor" → "uppercaseProcessor")
				if dotIdx := strings.LastIndex(qualName, "."); dotIdx >= 0 {
					return qualName[dotIdx+1:]
				}
				return qualName
			}
		}
	}
	return ""
}
