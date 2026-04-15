// ops_arithmetic.go implements arithmetic, bitwise, and complex number operations.
// Note: OpAdd, OpSub, OpMul, OpEqual, OpNotEqual, OpLess, OpLessEq, OpGreater,
// OpGreaterEq, and OpNot are inlined in run.go's hot path and never reach this handler.
package vm

import (
	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
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
		c := v.pop()
		if c.RawSize() == value.Size32 {
			v.push(value.MakeFloat32(float32(real(c.Complex()))))
		} else {
			v.push(value.MakeFloat(real(c.Complex())))
		}

	case bytecode.OpImag:
		c := v.pop()
		if c.RawSize() == value.Size32 {
			v.push(value.MakeFloat32(float32(imag(c.Complex()))))
		} else {
			v.push(value.MakeFloat(imag(c.Complex())))
		}

	case bytecode.OpComplex:
		imVal := v.pop()
		reVal := v.pop()
		re := reVal.Float()
		im := imVal.Float()
		// If either operand is float32-sized, create complex64
		if reVal.RawSize() == value.Size32 || imVal.RawSize() == value.Size32 {
			v.push(value.MakeComplex64(float32(re), float32(im)))
		} else {
			v.push(value.MakeComplex(re, im))
		}

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
