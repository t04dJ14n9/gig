// ops_arithmetic.go implements arithmetic, bitwise, and complex number operations.
// Note: OpAdd, OpSub, OpMul, OpEqual, OpNotEqual, OpLess, OpLessEq, OpGreater,
// OpGreaterEq, and OpNot are inlined in run.go's hot path and never reach this handler.
package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// toShiftAmount extracts a uint shift amount from a Value of any numeric kind.
// Used by OpLsh and OpRsh to ensure shift amounts are consistently handled.
func toShiftAmount(shiftVal value.Value) uint {
	if shiftVal.Kind() == value.KindUint {
		return uint(shiftVal.Uint())
	}
	return uint(shiftVal.Int())
}

// executeArithmetic handles non-hot-path arithmetic, bitwise, and complex opcodes.
// Hot-path ops (Add, Sub, Mul, comparisons, Not) are inlined in run.go.
func (v *vm) executeArithmetic(op bytecode.OpCode, frame *Frame) error { //nolint:cyclop,unparam // frame: uniform dispatch signature
	switch op {
	case bytecode.OpDiv:
		b := v.pop()
		a := v.pop()
		v.push(a.Div(b))

	case bytecode.OpMod:
		b := v.pop()
		a := v.pop()
		v.push(a.Mod(b))

	case bytecode.OpNeg:
		a := v.pop()
		v.push(a.Neg())

	case bytecode.OpReal:
		v.push(realComponent(v.pop()))

	case bytecode.OpImag:
		v.push(imagComponent(v.pop()))

	case bytecode.OpComplex:
		imVal := v.pop()
		reVal := v.pop()
		v.push(makeComplexValue(reVal, imVal))

	// Bitwise
	case bytecode.OpAnd:
		b := v.pop()
		a := v.pop()
		v.push(a.And(b))

	case bytecode.OpOr:
		b := v.pop()
		a := v.pop()
		v.push(a.Or(b))

	case bytecode.OpXor:
		b := v.pop()
		a := v.pop()
		v.push(a.Xor(b))

	case bytecode.OpAndNot:
		b := v.pop()
		a := v.pop()
		v.push(a.AndNot(b))

	case bytecode.OpLsh:
		shiftVal := v.pop()
		n := toShiftAmount(shiftVal)
		a := v.pop()
		v.push(a.Lsh(n))

	case bytecode.OpRsh:
		shiftVal := v.pop()
		n := toShiftAmount(shiftVal)
		a := v.pop()
		v.push(a.Rsh(n))
	}

	return nil
}

func realComponent(c value.Value) value.Value {
	if c.RawSize() == value.Size32 {
		return value.MakeFloat32(float32(real(c.Complex())))
	}
	return value.MakeFloat(real(c.Complex()))
}

func imagComponent(c value.Value) value.Value {
	if c.RawSize() == value.Size32 {
		return value.MakeFloat32(float32(imag(c.Complex())))
	}
	return value.MakeFloat(imag(c.Complex()))
}

func makeComplexValue(reVal, imVal value.Value) value.Value {
	re := reVal.Float()
	im := imVal.Float()
	if reVal.RawSize() == value.Size32 || imVal.RawSize() == value.Size32 {
		return value.MakeComplex64(float32(re), float32(im))
	}
	return value.MakeComplex(re, im)
}
