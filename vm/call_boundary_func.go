package vm

import "reflect"

func canPassInterpretedFuncToThirdParty(targetType reflect.Type) bool {
	if targetType == nil || targetType.Kind() != reflect.Func {
		return false
	}
	// Interpreted closures are safe only when the native callback cannot receive
	// or hide an interpreted value behind an interface-typed result.
	for i := 0; i < targetType.NumOut(); i++ {
		if reflectTypeContainsInterface(targetType.Out(i), make(map[reflect.Type]bool)) {
			return false
		}
	}
	return true
}

func reflectTypeContainsInterface(rt reflect.Type, seen map[reflect.Type]bool) bool {
	if !markReflectTypeForInterfaceScan(rt, seen) {
		return false
	}

	switch rt.Kind() {
	case reflect.Interface:
		return true
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Chan:
		return reflectTypeContainsInterface(rt.Elem(), seen)
	case reflect.Map:
		return reflectMapTypeContainsInterface(rt, seen)
	case reflect.Struct:
		return reflectStructTypeContainsInterface(rt, seen)
	case reflect.Func:
		return reflectFuncResultsContainInterface(rt, seen)
	}
	return false
}

func markReflectTypeForInterfaceScan(rt reflect.Type, seen map[reflect.Type]bool) bool {
	if rt == nil || seen[rt] {
		return false
	}
	// Mark before descending so recursive type graphs terminate without hiding
	// interface-typed leaves that appear through a different path.
	seen[rt] = true
	return true
}

func reflectMapTypeContainsInterface(rt reflect.Type, seen map[reflect.Type]bool) bool {
	return reflectTypeContainsInterface(rt.Key(), seen) ||
		reflectTypeContainsInterface(rt.Elem(), seen)
}

func reflectStructTypeContainsInterface(rt reflect.Type, seen map[reflect.Type]bool) bool {
	for i := 0; i < rt.NumField(); i++ {
		if reflectTypeContainsInterface(rt.Field(i).Type, seen) {
			return true
		}
	}
	return false
}

func reflectFuncResultsContainInterface(rt reflect.Type, seen map[reflect.Type]bool) bool {
	for i := 0; i < rt.NumOut(); i++ {
		if reflectTypeContainsInterface(rt.Out(i), seen) {
			return true
		}
	}
	return false
}
