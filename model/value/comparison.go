package value

import (
	"fmt"
	"reflect"
)

// Cmp compares v with other. Returns -1, 0, or 1.
func (v Value) Cmp(other Value) int {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		a, b := v.Bool(), other.Bool()
		if a == b {
			return 0
		}
		if !a {
			return -1
		}
		return 1
	case KindInt:
		a, b := v.num, other.Int()
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case KindUint:
		a, b := uint64(v.num), other.Uint()
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case KindFloat:
		a, b := v.Float(), other.Float()
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case KindString:
		a := v.obj.(string)
		b := other.obj.(string)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	default:
		panic(fmt.Sprintf("cannot compare %v", v.kind))
	}
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
	// Unwrap interface/reflect values for comparison.
	// When an external function returns an interface{} holding e.g. a string,
	// it comes back as KindReflect. We need to compare the underlying value.
	a, b := v, other

	// Treat KindInvalid (uninitialized zero Value) as KindNil.
	// Go globals with reference types (pointer, slice, map, chan, func, interface)
	// default to nil, but SSA may not emit explicit zero stores for them,
	// leaving them as KindInvalid in the interpreter. Comparing an uninitialized
	// global to nil should return true, just like in Go.
	if a.kind == KindInvalid {
		a = MakeNil()
	}
	if b.kind == KindInvalid {
		b = MakeNil()
	}

	// Handle nil comparison with interface values first.
	// In Go, a typed nil interface is NOT equal to nil.
	// e.g., var p *T = nil; var e error = p; e == nil -> false
	isInterfaceNil := func(val Value) (bool, bool) {
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
	if b.kind == KindNil {
		if isIface, isNil := isInterfaceNil(a); isIface {
			return isNil
		}
	}
	if a.kind == KindNil {
		if isIface, isNil := isInterfaceNil(b); isIface {
			return isNil
		}
	}

	if dynA, ok := a.InterpretedInterface(); ok {
		dynB, ok := b.InterpretedInterface()
		return ok && dynA.TypeName == dynB.TypeName && dynA.Value.Equal(dynB.Value)
	}
	if _, ok := b.InterpretedInterface(); ok {
		return false
	}

	if a.kind == KindReflect || a.kind == KindInterface {
		a = unwrapForComparison(a)
	}
	if b.kind == KindReflect || b.kind == KindInterface {
		b = unwrapForComparison(b)
	}
	if a.kind != b.kind {
		// Handle nil comparison
		if a.kind == KindNil || b.kind == KindNil {
			return a.IsNil() && b.IsNil()
		}
		return false
	}
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
		// For pointer types, compare by identity (address), not by value.
		// Go's == on pointers checks whether they point to the same location.
		if vp, ok := a.obj.(*Value); ok {
			if op, ok2 := b.obj.(*Value); ok2 {
				return vp == op // pointer identity
			}
			return false
		}
		if rv, ok := a.obj.(reflect.Value); ok {
			if rv.Kind() == reflect.Ptr {
				if orv, ok2 := b.obj.(reflect.Value); ok2 && orv.Kind() == reflect.Ptr {
					if rv.IsNil() && orv.IsNil() {
						return true
					}
					if rv.IsNil() || orv.IsNil() {
						return false
					}
					return rv.Pointer() == orv.Pointer()
				}
				return false
			}
		}
		// For other complex types, use reflect.DeepEqual
		return reflect.DeepEqual(a.Interface(), b.Interface())
	}
}
