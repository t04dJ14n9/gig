package vm

import (
	"fmt"
	"reflect"

	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

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

	if err := v.validateExternalMethodBoundary(methodInfo, args); err != nil {
		return err
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

func (v *vm) validateExternalMethodBoundary(methodInfo *external.ExternalMethodInfo, args []value.Value) error {
	if methodInfo == nil || v.program.AllowUnsafeTypePass || isStdlibExternalPath(methodInfo.PkgPath) {
		return nil
	}
	if len(args) == 0 {
		return nil
	}
	methodType := reflectMethodTypeForBoundary(args[0], methodInfo.MethodName)
	for i, arg := range args[1:] {
		targetType := externalBoundaryReflectArgType(methodType, i)
		if typeName, ok := v.interpreterDefinedBoundaryType(arg, targetType); ok {
			funcName := methodInfo.FuncName
			if funcName == "" {
				funcName = methodInfo.MethodName
			}
			return fmt.Errorf(
				"cannot pass interpreter-defined type %q to third-party function %s.%s (argument %d): "+
					"value crossed the boundary through an interface. "+
					"Use primitive types, slices, maps, types from registered packages, or a registered interface proxy instead",
				typeName, methodInfo.PkgPath, funcName, i+1,
			)
		}
	}
	return nil
}

func reflectMethodTypeForBoundary(receiver value.Value, methodName string) reflect.Type {
	rv, ok := receiver.ReflectValue()
	if !ok {
		iface := receiver.Interface()
		if iface == nil {
			return nil
		}
		rv = reflect.ValueOf(iface)
	}
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if method := rv.MethodByName(methodName); method.IsValid() {
		return method.Type()
	}
	if rv.CanAddr() {
		if method := rv.Addr().MethodByName(methodName); method.IsValid() {
			return method.Type()
		}
	}
	if !rv.CanAddr() && rv.Kind() == reflect.Struct {
		addrCopy := reflect.New(rv.Type()).Elem()
		addrCopy.Set(rv)
		if method := addrCopy.Addr().MethodByName(methodName); method.IsValid() {
			return method.Type()
		}
	}
	if rv.Kind() == reflect.Struct {
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			if field.Kind() != reflect.Interface || field.IsNil() {
				continue
			}
			concrete := field.Elem()
			if method := concrete.MethodByName(methodName); method.IsValid() {
				return method.Type()
			}
			if concrete.CanAddr() {
				if method := concrete.Addr().MethodByName(methodName); method.IsValid() {
					return method.Type()
				}
			}
		}
	}
	return nil
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
			v.panicking = true
			v.panicVal = value.FromInterface("runtime error: invalid memory address or nil pointer dereference")
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
			v.panicking = true
			v.panicVal = value.FromInterface("runtime error: invalid memory address or nil pointer dereference")
			return nil
		}
		concrete := rv.Elem()
		rv = concrete
		args[0] = value.MakeFromReflect(rv)
	}

	// Nil pointer check for non-interface dispatch (e.g., free variables
	// that lost the interface wrapper). In Go, calling any method on a
	// nil *T through an interface panics because the runtime dereferences
	// the pointer. When the value arrives here as a raw nil pointer (not
	// wrapped in an interface), we must also panic to match Go semantics.
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
				if len(args) > 0 {
					args[0] = value.MakeFromReflect(concrete)
				}
				return m, true
			}
			if concrete.CanAddr() {
				if m := concrete.Addr().MethodByName(methodName); m.IsValid() {
					if len(args) > 0 {
						args[0] = value.MakeFromReflect(concrete.Addr())
					}
					return m, true
				}
			}
		}
	}

	return reflect.Value{}, false
}
