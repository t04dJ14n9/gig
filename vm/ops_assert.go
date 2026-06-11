package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeAssert(frame *Frame) {
	typeIdx := frame.readUint16()
	targetType := v.program.Types[typeIdx]
	obj := v.pop()

	// Empty interface assertions are valid for every value, including nil and
	// script-defined interface wrappers. Keep this before payload-specific
	// handling so `x.(any)` preserves the original dynamic value.
	if iface, isInterface := targetType.Underlying().(*types.Interface); isInterface && iface.NumMethods() == 0 {
		v.pushCommaOk(obj, true)
		return
	}

	if dyn, ok := obj.InterpretedInterface(); ok {
		result, assertionOk := v.assertInterpretedInterfaceValue(dyn, targetType, obj)
		v.pushCommaOk(result, assertionOk)
		return
	}

	if obj.Kind() == value.KindInterface {
		rv, isReflect := obj.ReflectValue()
		if !isReflect || rv.Kind() != reflect.Interface {
			v.pushCommaOk(obj, true)
			return
		}
		if rv.IsNil() {
			v.pushCommaOk(value.MakeNil(), false)
			return
		}

		underlying := rv.Elem()
		targetReflectType := typeToReflect(targetType, v.program)
		if targetReflectType == nil {
			v.pushCommaOk(value.MakeNil(), false)
			return
		}
		if underlying.Type().AssignableTo(targetReflectType) {
			v.pushCommaOk(value.MakeFromReflect(underlying), true)
			return
		}
		v.pushCommaOk(value.MakeFromReflect(reflect.Zero(targetReflectType)), false)
		return
	}

	var result value.Value
	var assertionOk bool
	if obj.Kind() == value.KindReflect {
		result, assertionOk = v.assertReflectValue(obj, targetType)
	} else {
		result, assertionOk = v.assertPrimitiveValue(obj, targetType)
	}
	v.pushCommaOk(result, assertionOk)
}

func (v *vm) assertReflectValue(obj value.Value, targetType types.Type) (value.Value, bool) {
	rv, ok := obj.ReflectValue()
	if !ok {
		return zeroValueForType(targetType, v.program), false
	}

	// Reflect values can themselves be interface boxes. Unwrap only non-nil
	// boxes; a nil interface box must continue to fail against concrete targets.
	concreteRV := reflectConcreteAssertionValue(rv)
	targetReflectType := typeToReflect(targetType, v.program)
	if targetReflectType == nil {
		return value.MakeNil(), false
	}
	return assertReflectValueToType(concreteRV, targetReflectType)
}

func reflectConcreteAssertionValue(rv reflect.Value) reflect.Value {
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		return rv.Elem()
	}
	return rv
}

func assertReflectValueToType(rv reflect.Value, target reflect.Type) (value.Value, bool) {
	if rv.Type().AssignableTo(target) {
		return value.MakeFromReflect(rv), true
	}
	if sameReflectKindFamily(rv.Type(), target) {
		// Gig stores numeric values with internal types (int64 for int, float64
		// for float32, etc.). Type switches should still match the Go kind and
		// return the target's reflected type.
		return value.MakeFromReflect(rv.Convert(target)), true
	}
	return value.MakeNil(), false
}

func (v *vm) assertPrimitiveValue(obj value.Value, targetType types.Type) (value.Value, bool) {
	// Primitive Value kinds hold unboxed interpreter scalars. They can only
	// satisfy a type assertion when their normalized kind and raw size match
	// the target Go type used by the compiled type-switch/assertion.
	if kindMatchesType(obj.Kind(), obj.RawSize(), targetType) {
		return obj, true
	}
	return zeroValueForType(targetType, v.program), false
}
