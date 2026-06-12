// arithmetic.go implements arithmetic operations on Value.
package value

import (
	"fmt"
	"math"
)

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
		return MakeComplexSized(real(a)+real(b), imag(a)+imag(b), v.size)
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
		return MakeComplexSized(real(a)-real(b), imag(a)-imag(b), v.size)
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
		return MakeComplexSized(real(a)*real(b)-imag(a)*imag(b), real(a)*imag(b)+real(b)*imag(a), v.size)
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
		return MakeComplexSized((real(a)*real(b)+imag(a)*imag(b))/denom, (imag(a)*real(b)-real(a)*imag(b))/denom, v.size)
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
		return MakeComplexSized(-real(c), -imag(c), v.size)
	default:
		panic(fmt.Sprintf("cannot neg %v", v.kind))
	}
}
