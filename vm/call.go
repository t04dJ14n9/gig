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

	// Fast path: DirectCall available
	if rc.DirectCall != nil {
		if rc.IsVariadic && numArgs == rc.NumIn {
			args = unpackVariadicArgs(args, numArgs)
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
	if !rc.Fn.IsValid() || rc.Fn.Kind() != reflect.Func {
		v.push(value.MakeNil())
		return nil
	}

	in := buildReflectArgs(args, rc.FnType)
	out := rc.Fn.Call(in)

	if err := v.checkCtx(); err != nil {
		return err
	}
	v.pushReflectResults(out)
	return nil
}

// callExternalMethod dispatches a method call on an external type.
// args[0] is the receiver, args[1:] are the method arguments.
func (v *vm) callExternalMethod(methodInfo *external.ExternalMethodInfo, args []value.Value) error {
	if len(args) == 0 {
		v.push(value.MakeNil())
		return nil
	}

	// Resolve GlobalRef / *value.Value receivers.
	if iface0 := args[0].Interface(); iface0 != nil {
		switch ref := iface0.(type) {
		case *GlobalRef:
			args[0] = ref.Load()
		case *value.Value:
			args[0] = *ref
		}
	}

	// Fast path: DirectCall wrapper resolved at compile time
	if methodInfo.DirectCall != nil {
		convertClosureArgsForMethod(methodInfo.MethodName, args)
		v.push(methodInfo.DirectCall(args))
		return v.checkCtx()
	}

	// Slow path: use reflect.MethodByName + reflect.Call
	return v.callExternalMethodReflect(methodInfo, args)
}

// callExternalMethodReflect dispatches a method call using reflection.
func (v *vm) callExternalMethodReflect(methodInfo *external.ExternalMethodInfo, args []value.Value) error {
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

	// For interface method dispatch: unwrap to concrete type
	if rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			v.push(value.MakeNil())
			return nil
		}
		concrete := rv.Elem()
		if concrete.Kind() == reflect.Ptr && concrete.IsNil() {
			v.panicking = true
			v.panicVal = value.FromInterface("runtime error: invalid memory address or nil pointer dereference")
			return nil
		}
		rv = concrete
		args[0] = value.MakeFromReflect(rv)
	}

	// Look up the method by name
	method, found := findMethod(rv, methodInfo.MethodName, args)
	if !found {
		return v.callCompiledMethod(methodInfo.MethodName, methodInfo.ReceiverTypeName, args)
	}

	// Build arguments (skip the receiver at args[0])
	methodType := method.Type()
	methodArgs := args[1:]
	in := buildReflectArgs(methodArgs, methodType)

	out := method.Call(in)

	if err := v.checkCtx(); err != nil {
		return err
	}
	v.pushReflectResults(out)
	return nil
}

// findMethod resolves a method by name on a reflect.Value, trying (in order):
// 1. Direct method lookup on the value
// 2. Pointer receiver method via Addr()
// 3. Pointer receiver via addressable copy (for non-addressable structs)
// 4. Methods on concrete values inside embedded interface fields
func findMethod(rv reflect.Value, methodName string, args []value.Value) (reflect.Value, bool) {
	method := rv.MethodByName(methodName)
	if method.IsValid() {
		return method, true
	}

	if rv.CanAddr() {
		method = rv.Addr().MethodByName(methodName)
		if method.IsValid() {
			return method, true
		}
	}

	if !rv.CanAddr() && rv.Kind() == reflect.Struct {
		addrCopy := reflect.New(rv.Type()).Elem()
		addrCopy.Set(rv)
		method = addrCopy.Addr().MethodByName(methodName)
		if method.IsValid() {
			return method, true
		}
	}

	if rv.Kind() == reflect.Struct {
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			if field.Kind() != reflect.Interface || field.IsNil() {
				continue
			}
			concrete := field.Elem()
			if m := concrete.MethodByName(methodName); m.IsValid() {
				args[0] = value.MakeFromReflect(concrete)
				return m, true
			}
			if concrete.CanAddr() {
				if m := concrete.Addr().MethodByName(methodName); m.IsValid() {
					args[0] = value.MakeFromReflect(concrete.Addr())
					return m, true
				}
			}
		}
	}

	return reflect.Value{}, false
}

// callCompiledMethod searches the compiled function table for a method with the
// given name and calls it. This is the fallback path for invoke (interface method)
// calls when reflection-based MethodByName fails.
func (v *vm) callCompiledMethod(methodName string, receiverTypeName string, args []value.Value) error {
	candidates := v.program.MethodsByName[methodName]

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

	if len(candidates) > 0 {
		fn := candidates[0]
		for _, arg := range args {
			v.push(arg)
		}
		v.callCompiledFunction(fn.FuncIdx, len(args))
		return nil
	}

	v.push(value.MakeNil())
	return nil
}

// inferReceiverTypeName tries to extract a type name from a runtime value.Value receiver.
func inferReceiverTypeName(receiver value.Value, prog *bytecode.CompiledProgram) string {
	rv, ok := receiver.ReflectValue()
	if !ok {
		return ""
	}
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	t := rv.Type()
	if t.Name() != "" {
		return t.Name()
	}
	if prog != nil {
		if name := prog.LookupTypeName(t); name != "" {
			return name
		}
	}
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			pkgPath := f.PkgPath
			if idx := strings.LastIndex(pkgPath, "#"); idx >= 0 {
				qualName := pkgPath[idx+1:]
				if dotIdx := strings.LastIndex(qualName, "."); dotIdx >= 0 {
					return qualName[dotIdx+1:]
				}
				return qualName
			}
		}
	}
	return ""
}
