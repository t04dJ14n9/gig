package value

import (
	"fmt"
	"math"
	"reflect"
)

// --- Arithmetic Operations ---

// MakeIntSized creates an int value preserving the given size tag.
func MakeIntSized(i int64, s Size) Value {
	return Value{kind: KindInt, size: s, num: i}
}

// makeUintSized creates a uint value preserving the given size tag.
func makeUintSized(u uint64, s Size) Value {
	return Value{kind: KindUint, size: s, num: int64(u)}
}

// makeFloatSized creates a float value preserving the given size tag.
func makeFloatSized(f float64, s Size) Value {
	return Value{kind: KindFloat, size: s, num: int64(math.Float64bits(f))}
}

// Add returns v + other.
func (v Value) Add(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeIntSized(v.num+other.Int(), v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)+other.Uint(), v.size)
	case KindFloat:
		return makeFloatSized(v.Float()+other.Float(), v.size)
	case KindString:
		return MakeString(v.obj.(string) + other.obj.(string))
	case KindComplex:
		a := v.obj.(complex128)
		b := other.obj.(complex128)
		return MakeComplex(real(a)+real(b), imag(a)+imag(b))
	default:
		panic(fmt.Sprintf("cannot add %v", v.kind))
	}
}

// Sub returns v - other.
func (v Value) Sub(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeIntSized(v.num-other.Int(), v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)-other.Uint(), v.size)
	case KindFloat:
		return makeFloatSized(v.Float()-other.Float(), v.size)
	case KindComplex:
		a := v.obj.(complex128)
		b := other.obj.(complex128)
		return MakeComplex(real(a)-real(b), imag(a)-imag(b))
	default:
		panic(fmt.Sprintf("cannot sub %v", v.kind))
	}
}

// Mul returns v * other.
func (v Value) Mul(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeIntSized(v.num*other.Int(), v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)*other.Uint(), v.size)
	case KindFloat:
		return makeFloatSized(v.Float()*other.Float(), v.size)
	case KindComplex:
		a := v.obj.(complex128)
		b := other.obj.(complex128)
		return MakeComplex(real(a)*real(b)-imag(a)*imag(b), real(a)*imag(b)+real(b)*imag(a))
	default:
		panic(fmt.Sprintf("cannot mul %v", v.kind))
	}
}

// Div returns v / other.
func (v Value) Div(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeIntSized(v.num/other.Int(), v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)/other.Uint(), v.size)
	case KindFloat:
		return makeFloatSized(v.Float()/other.Float(), v.size)
	case KindComplex:
		a := v.obj.(complex128)
		b := other.obj.(complex128)
		denom := real(b)*real(b) + imag(b)*imag(b)
		return MakeComplex((real(a)*real(b)+imag(a)*imag(b))/denom, (imag(a)*real(b)-real(a)*imag(b))/denom)
	default:
		panic(fmt.Sprintf("cannot div %v", v.kind))
	}
}

// Mod returns v % other.
func (v Value) Mod(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeIntSized(v.num%other.Int(), v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)%other.Uint(), v.size)
	case KindFloat:
		return makeFloatSized(math.Mod(v.Float(), other.Float()), v.size)
	default:
		panic(fmt.Sprintf("cannot mod %v", v.kind))
	}
}

// Neg returns -v.
func (v Value) Neg() Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeIntSized(-v.num, v.size)
	case KindFloat:
		return makeFloatSized(-v.Float(), v.size)
	case KindComplex:
		c := v.obj.(complex128)
		return MakeComplex(-real(c), -imag(c))
	default:
		panic(fmt.Sprintf("cannot neg %v", v.kind))
	}
}

// --- Comparison Operations ---

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

// Equal returns v == other.
func (v Value) Equal(other Value) bool {
	// Unwrap interface/reflect values for comparison.
	// When an external function returns an interface{} holding e.g. a string,
	// it comes back as KindReflect. We need to compare the underlying value.
	a, b := v, other
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
		return a.num == b.num
	case KindUint:
		return a.num == b.num
	case KindFloat:
		return a.Float() == b.Float()
	case KindString:
		return a.obj.(string) == b.obj.(string)
	case KindComplex:
		return a.obj.(complex128) == b.obj.(complex128)
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

// --- Bitwise Operations ---

// And returns v & other.
func (v Value) And(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeIntSized(v.num&other.Int(), v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)&other.Uint(), v.size)
	default:
		panic(fmt.Sprintf("cannot and %v", v.kind))
	}
}

// Or returns v | other.
func (v Value) Or(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeIntSized(v.num|other.Int(), v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)|other.Uint(), v.size)
	default:
		panic(fmt.Sprintf("cannot or %v", v.kind))
	}
}

// Xor returns v ^ other.
func (v Value) Xor(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeIntSized(v.num^other.Int(), v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)^other.Uint(), v.size)
	default:
		panic(fmt.Sprintf("cannot xor %v", v.kind))
	}
}

// AndNot returns v &^ other.
func (v Value) AndNot(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeIntSized(v.num&^other.Int(), v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)&^other.Uint(), v.size)
	default:
		panic(fmt.Sprintf("cannot andnot %v", v.kind))
	}
}

// Lsh returns v << n.
func (v Value) Lsh(n uint) Value {
	switch v.kind {
	case KindInt:
		return MakeIntSized(v.num<<n, v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)<<n, v.size)
	default:
		panic(fmt.Sprintf("cannot lsh %v", v.kind))
	}
}

// Rsh returns v >> n.
func (v Value) Rsh(n uint) Value {
	switch v.kind {
	case KindInt:
		return MakeIntSized(v.num>>n, v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)>>n, v.size)
	default:
		panic(fmt.Sprintf("cannot rsh %v", v.kind))
	}
}
