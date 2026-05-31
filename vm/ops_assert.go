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

	// Type assertion - check if obj can be asserted to targetType
	// Returns (value, ok) tuple on stack
	var result value.Value
	var assertionOk bool

	// Special case: any value can be asserted to interface{} (empty interface).
	// This handles nested interface assertions like: outer.(interface{})
	// Use Underlying() to also handle *types.Named wrapping interface{} (e.g., "any").
	if iface, isInterface := targetType.Underlying().(*types.Interface); isInterface && iface.NumMethods() == 0 {
		// Any value can be asserted to interface{}, including nil
		result = obj
		assertionOk = true
		v.pushCommaOk(result, assertionOk)
		return
	}

	if dyn, ok := obj.InterpretedInterface(); ok {
		result, assertionOk = v.assertInterpretedInterfaceValue(dyn, targetType, obj)
		v.pushCommaOk(result, assertionOk)
		return
	}

	if obj.Kind() == value.KindInterface {
		// Get the underlying interface
		if rv, isReflect := obj.ReflectValue(); isReflect && rv.Kind() == reflect.Interface {
			if rv.IsNil() {
				// Interface is nil, assertion fails
				result = value.MakeNil()
				assertionOk = false
			} else {
				// Get the underlying value
				underlying := rv.Elem()
				targetReflectType := typeToReflect(targetType, v.program)
				if targetReflectType == nil {
					result = value.MakeNil()
					assertionOk = false
				} else if underlying.Type().AssignableTo(targetReflectType) {
					// Successful assertion - create a value from the underlying
					result = value.MakeFromReflect(underlying)
					assertionOk = true
				} else {
					// Failed assertion — return zero value of target type, not nil
					result = value.MakeFromReflect(reflect.Zero(targetReflectType))
					assertionOk = false
				}
			}
		} else {
			result = obj
			assertionOk = true
		}
	} else if obj.Kind() == value.KindReflect {
		// Already a reflect value — may be an interface wrapping a concrete type
		if rv, isReflect := obj.ReflectValue(); isReflect {
			// If the reflect value is an interface, unwrap it to the concrete type
			concreteRV := rv
			if rv.Kind() == reflect.Interface && !rv.IsNil() {
				concreteRV = rv.Elem()
			}
			targetReflectType := typeToReflect(targetType, v.program)
			if targetReflectType != nil && concreteRV.Type().AssignableTo(targetReflectType) {
				result = value.MakeFromReflect(concreteRV)
				assertionOk = true
			} else if targetReflectType != nil && sameReflectKindFamily(concreteRV.Type(), targetReflectType) {
				// Gig stores numeric values with internal types (int64 for int, float64 for float32, etc.).
				// When these are stored in interface{}, the concrete type may differ from what Go
				// would use (e.g., int64 instead of int). For type switch correctness, we match
				// by reflect.Kind family and convert to the target type.
				result = value.MakeFromReflect(concreteRV.Convert(targetReflectType))
				assertionOk = true
			} else {
				result = value.MakeNil()
				assertionOk = false
			}
		} else {
			targetReflectType := typeToReflect(targetType, v.program)
			if targetReflectType != nil {
				result = value.MakeFromReflect(reflect.Zero(targetReflectType))
			} else {
				result = value.MakeNil()
			}
			assertionOk = false
		}
	} else {
		// For primitive kinds (KindInt, KindString, KindBool, KindFloat, etc.),
		// check whether the concrete value kind actually matches the target type.
		// This is critical for type switches on interface slice elements.
		assertionOk = kindMatchesType(obj.Kind(), obj.RawSize(), targetType)
		if assertionOk {
			result = obj
		} else {
			// Failed assertion — return zero value of target type, not nil
			targetReflectType := typeToReflect(targetType, v.program)
			if targetReflectType != nil {
				result = value.MakeFromReflect(reflect.Zero(targetReflectType))
			} else {
				result = value.MakeNil()
			}
		}
	}

	// Push result as a tuple [result, ok]
	v.pushCommaOk(result, assertionOk)
}
