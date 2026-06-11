package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func callReceiverSliceSwap(receiver value.Value, i, j int) bool {
	// This is a narrow compatibility fallback for sort.Interface values that do
	// not expose a compiled Swap method. It only mutates actual slice receivers.
	rv, ok := receiver.ReflectValue()
	if !ok {
		return false
	}
	rv, ok = derefInterfaceOrPointer(rv)
	if !ok || rv.Kind() != reflect.Slice {
		return false
	}
	tmp := reflect.New(rv.Type().Elem()).Elem()
	tmp.Set(rv.Index(i))
	rv.Index(i).Set(rv.Index(j))
	rv.Index(j).Set(tmp)
	return true
}

func callReflectInterfaceMethod(methodName string, receiver value.Value, args []value.Value) (value.Value, bool) {
	// Reflect calls are intentionally late fallback. If compiled script methods
	// exist, callInterfaceMethodValue has already used them.
	rv, ok := reflectValueForInterfaceMethod(receiver)
	if !ok {
		return value.MakeNil(), false
	}
	method, found := findMethod(rv, methodName, nil)
	if !found {
		return value.MakeNil(), false
	}
	in := buildReflectArgs(args, method.Type())
	out := method.Call(in)
	return firstReflectResult(out)
}

func callEmbeddedInterfaceMethod(receiver value.Value, methodName string, args ...value.Value) (value.Value, bool) {
	// Some external structs embed interface fields. This helper searches those
	// embedded interface values without treating the outer script value as a
	// general third-party interface implementation.
	rv, ok := receiver.ReflectValue()
	if !ok {
		return value.MakeNil(), false
	}
	rv, ok = derefInterfaceOrPointer(rv)
	if !ok || rv.Kind() != reflect.Struct {
		return value.MakeNil(), false
	}
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.Kind() != reflect.Interface || field.IsNil() {
			continue
		}
		method := field.Elem().MethodByName(methodName)
		if !method.IsValid() {
			continue
		}
		in := buildReflectArgs(args, method.Type())
		out := method.Call(in)
		return firstReflectResult(out)
	}
	return value.MakeNil(), false
}

func reflectValueForInterfaceMethod(receiver value.Value) (reflect.Value, bool) {
	rv, ok := receiver.ReflectValue()
	if !ok {
		iface := receiver.Interface()
		if iface == nil {
			return reflect.Value{}, false
		}
		rv = reflect.ValueOf(iface)
	}
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	return rv, true
}

func derefInterfaceOrPointer(rv reflect.Value) (reflect.Value, bool) {
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, false
		}
		rv = rv.Elem()
	}
	return rv, true
}

func firstReflectResult(out []reflect.Value) (value.Value, bool) {
	if len(out) == 0 {
		return value.MakeNil(), true
	}
	return value.MakeFromReflect(out[0]), true
}
