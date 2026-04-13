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
		v.push(value.MakeFloat(real(c.Complex())))

	case bytecode.OpImag:
		c := v.pop()
		v.push(value.MakeFloat(imag(c.Complex())))

	case bytecode.OpComplex:
		im := v.pop().Float()
		re := v.pop().Float()
		v.push(value.MakeComplex(re, im))

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
