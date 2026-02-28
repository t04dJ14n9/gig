package value

import (
	"fmt"
	"math"
	"reflect"
)

// --- Arithmetic Operations ---

// Add returns v + other.
func (v Value) Add(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeInt(v.num + other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) + other.Uint())
	case KindFloat:
		return MakeFloat(v.Float() + other.Float())
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
		return MakeInt(v.num - other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) - other.Uint())
	case KindFloat:
		return MakeFloat(v.Float() - other.Float())
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
		return MakeInt(v.num * other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) * other.Uint())
	case KindFloat:
		return MakeFloat(v.Float() * other.Float())
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
		return MakeInt(v.num / other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) / other.Uint())
	case KindFloat:
		return MakeFloat(v.Float() / other.Float())
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
		return MakeInt(v.num % other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) % other.Uint())
	case KindFloat:
		return MakeFloat(math.Mod(v.Float(), other.Float()))
	default:
		panic(fmt.Sprintf("cannot mod %v", v.kind))
	}
}

// Neg returns -v.
func (v Value) Neg() Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeInt(-v.num)
	case KindFloat:
		return MakeFloat(-v.Float())
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

// Equal returns v == other.
func (v Value) Equal(other Value) bool {
	if v.kind != other.kind {
		// Handle nil comparison
		if v.kind == KindNil || other.kind == KindNil {
			return v.IsNil() && other.IsNil()
		}
		return false
	}
	switch v.kind {
	case KindNil:
		return true
	case KindBool:
		return v.num == other.num
	case KindInt:
		return v.num == other.num
	case KindUint:
		return v.num == other.num
	case KindFloat:
		return v.Float() == other.Float()
	case KindString:
		return v.obj.(string) == other.obj.(string)
	case KindComplex:
		return v.obj.(complex128) == other.obj.(complex128)
	default:
		// For complex types, use reflect.DeepEqual
		return reflect.DeepEqual(v.Interface(), other.Interface())
	}
}

// --- Bitwise Operations ---

// And returns v & other.
func (v Value) And(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num & other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) & other.Uint())
	default:
		panic(fmt.Sprintf("cannot and %v", v.kind))
	}
}

// Or returns v | other.
func (v Value) Or(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num | other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) | other.Uint())
	default:
		panic(fmt.Sprintf("cannot or %v", v.kind))
	}
}

// Xor returns v ^ other.
func (v Value) Xor(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num ^ other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) ^ other.Uint())
	default:
		panic(fmt.Sprintf("cannot xor %v", v.kind))
	}
}

// AndNot returns v &^ other.
func (v Value) AndNot(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num &^ other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) &^ other.Uint())
	default:
		panic(fmt.Sprintf("cannot andnot %v", v.kind))
	}
}

// Lsh returns v << n.
func (v Value) Lsh(n uint) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num << n)
	case KindUint:
		return MakeUint(uint64(v.num) << n)
	default:
		panic(fmt.Sprintf("cannot lsh %v", v.kind))
	}
}

// Rsh returns v >> n.
func (v Value) Rsh(n uint) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num >> n)
	case KindUint:
		return MakeUint(uint64(v.num) >> n)
	default:
		panic(fmt.Sprintf("cannot rsh %v", v.kind))
	}
}
