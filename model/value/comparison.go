package value

import (
	"fmt"
	"reflect"
)

// Cmp compares v with other. Returns -1, 0, or 1.
func (v Value) Cmp(other Value) int {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		return cmpBool(v.Bool(), other.Bool())
	case KindInt:
		return cmpOrdered(v.num, other.Int())
	case KindUint:
		return cmpOrdered(uint64(v.num), other.Uint())
	case KindFloat:
		return cmpOrdered(v.Float(), other.Float())
	case KindString:
		return cmpOrdered(v.obj.(string), other.obj.(string))
	default:
		panic(fmt.Sprintf("cannot compare %v", v.kind))
	}
}

func cmpBool(a, b bool) int {
	if a == b {
		return 0
	}
	if !a {
		return -1
	}
	return 1
}

// cmpOrdered centralizes the -1/0/1 contract used by comparison opcodes.
// For floats this intentionally preserves the previous NaN behavior: both
// ordering checks are false, so Cmp returns 0.
func cmpOrdered[T ~int64 | ~uint64 | ~float64 | ~string](a, b T) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// unwrapForComparison converts a KindReflect/KindInterface value to its underlying
// primitive Value if possible. This enables correct equality comparisons between
// interface-wrapped values and primitive values (e.g., ctx.Value("key") == "value").
func unwrapForComparison(v Value) Value {
	rv, ok := v.obj.(reflect.Value)
	if !ok {
		return v
	}
	// Unwrap interface
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	return MakeFromReflect(rv)
}

func normalizedEqualitySize(kind Kind, size Size) Size {
	switch kind {
	case KindInt, KindUint:
		if size == Size0 {
			return SizePtr
		}
	case KindFloat, KindComplex:
		if size == Size0 {
			return Size64
		}
	}
	return size
}

func sameNumericEqualityType(a, b Value) bool {
	return normalizedEqualitySize(a.kind, a.size) == normalizedEqualitySize(b.kind, b.size)
}

// Equal returns v == other.
func (v Value) Equal(other Value) bool {
	a, b := normalizeEqualityValues(v, other)
	if result, handled := equalNilInterface(a, b); handled {
		return result
	}
	if result, handled := equalInterpretedInterfaces(a, b); handled {
		return result
	}
	a, b = unwrapEqualityValues(a, b)
	if a.kind != b.kind {
		return equalDifferentKinds(a, b)
	}
	return equalSameKind(a, b)
}

func normalizeEqualityValues(a, b Value) (Value, Value) {
	// SSA can leave zero stores implicit for nil-able globals. Treat the
	// interpreter zero Value the same way Go treats an unstored nil reference.
	return normalizeInvalidForEquality(a), normalizeInvalidForEquality(b)
}

func normalizeInvalidForEquality(v Value) Value {
	if v.kind == KindInvalid {
		return MakeNil()
	}
	return v
}

func equalNilInterface(a, b Value) (bool, bool) {
	if b.kind == KindNil {
		if isIface, isNil := interfaceNilState(a); isIface {
			return isNil, true
		}
	}
	if a.kind == KindNil {
		if isIface, isNil := interfaceNilState(b); isIface {
			return isNil, true
		}
	}
	return false, false
}

func interfaceNilState(val Value) (bool, bool) {
	// Go distinguishes a nil interface from an interface holding a typed nil.
	// Script-defined dynamic values are typed interface contents, even when the
	// carried value is nil.
	if val.kind == KindInterface {
		if _, ok := val.InterpretedInterface(); ok {
			return true, false
		}
		if rv, ok := val.obj.(reflect.Value); ok && rv.Kind() == reflect.Interface {
			return true, rv.IsNil()
		}
		return true, true
	}
	if val.kind == KindReflect {
		if rv, ok := val.obj.(reflect.Value); ok && rv.Kind() == reflect.Interface {
			return true, rv.IsNil()
		}
	}
	return false, false
}

func equalInterpretedInterfaces(a, b Value) (bool, bool) {
	dynA, okA := a.InterpretedInterface()
	dynB, okB := b.InterpretedInterface()
	if okA && okB {
		return dynA.TypeName == dynB.TypeName && dynA.Value.Equal(dynB.Value), true
	}
	if okA || okB {
		return false, true
	}
	return false, false
}

func unwrapEqualityValues(a, b Value) (Value, Value) {
	// External calls often return interface{} contents as KindReflect; unwrap
	// those before comparing against native Gig primitive values.
	return unwrapEqualityValue(a), unwrapEqualityValue(b)
}

func unwrapEqualityValue(v Value) Value {
	if v.kind == KindReflect || v.kind == KindInterface {
		return unwrapForComparison(v)
	}
	return v
}

func equalDifferentKinds(a, b Value) bool {
	if a.kind == KindNil || b.kind == KindNil {
		return a.IsNil() && b.IsNil()
	}
	return false
}

func equalSameKind(a, b Value) bool {
	switch a.kind {
	case KindNil:
		return true
	case KindBool:
		return a.num == b.num
	case KindInt:
		return sameNumericEqualityType(a, b) && a.num == b.num
	case KindUint:
		return sameNumericEqualityType(a, b) && a.num == b.num
	case KindFloat:
		return sameNumericEqualityType(a, b) && a.Float() == b.Float()
	case KindString:
		return a.obj.(string) == b.obj.(string)
	case KindComplex:
		return sameNumericEqualityType(a, b) && a.obj.(complex128) == b.obj.(complex128)
	default:
		return equalReferenceOrComposite(a, b)
	}
}

func equalReferenceOrComposite(a, b Value) bool {
	if result, handled := equalValuePointerIdentity(a, b); handled {
		return result
	}
	if result, handled := equalReflectPointerIdentity(a, b); handled {
		return result
	}
	return reflect.DeepEqual(a.Interface(), b.Interface())
}

func equalValuePointerIdentity(a, b Value) (bool, bool) {
	// Interpreter pointers compare by identity, matching Go pointer equality.
	vp, ok := a.obj.(*Value)
	if !ok {
		return false, false
	}
	op, ok := b.obj.(*Value)
	return ok && vp == op, true
}

func equalReflectPointerIdentity(a, b Value) (bool, bool) {
	rv, ok := a.obj.(reflect.Value)
	if !ok || rv.Kind() != reflect.Ptr {
		return false, false
	}
	orv, ok := b.obj.(reflect.Value)
	if !ok || orv.Kind() != reflect.Ptr {
		return false, true
	}
	if rv.IsNil() && orv.IsNil() {
		return true, true
	}
	if rv.IsNil() || orv.IsNil() {
		return false, true
	}
	return rv.Pointer() == orv.Pointer(), true
}
