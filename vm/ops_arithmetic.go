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
func (v *vm) executeArithmetic(op bytecode.OpCode, frame *Frame) error { //nolint:unparam // frame: uniform dispatch signature
	switch op {
	case bytecode.OpDiv, bytecode.OpMod, bytecode.OpNeg:
		v.executeBasicArithmetic(op)
	case bytecode.OpReal, bytecode.OpImag, bytecode.OpComplex:
		v.executeComplexArithmetic(op)
	case bytecode.OpAnd, bytecode.OpOr, bytecode.OpXor, bytecode.OpAndNot:
		v.executeBitwiseArithmetic(op)
	case bytecode.OpLsh, bytecode.OpRsh:
		v.executeShiftArithmetic(op)
	}

	return nil
}

func (v *vm) executeBasicArithmetic(op bytecode.OpCode) {
	switch op {
	case bytecode.OpDiv:
		a, b := v.popBinaryValues()
		v.push(a.Div(b))
	case bytecode.OpMod:
		a, b := v.popBinaryValues()
		v.push(a.Mod(b))
	case bytecode.OpNeg:
		a := v.pop()
		v.push(a.Neg())
	}
}

func (v *vm) executeComplexArithmetic(op bytecode.OpCode) {
	switch op {
	case bytecode.OpReal:
		v.push(realComponent(v.pop()))
	case bytecode.OpImag:
		v.push(imagComponent(v.pop()))
	case bytecode.OpComplex:
		imVal := v.pop()
		reVal := v.pop()
		v.push(makeComplexValue(reVal, imVal))
	}
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
	// If either operand is float32-sized, create complex64.
	if reVal.RawSize() == value.Size32 || imVal.RawSize() == value.Size32 {
		return value.MakeComplex64(float32(re), float32(im))
	}
	return value.MakeComplex(re, im)
}

func (v *vm) executeBitwiseArithmetic(op bytecode.OpCode) {
	a, b := v.popBinaryValues()
	switch op {
	case bytecode.OpAnd:
		v.push(a.And(b))
	case bytecode.OpOr:
		v.push(a.Or(b))
	case bytecode.OpXor:
		v.push(a.Xor(b))
	case bytecode.OpAndNot:
		v.push(a.AndNot(b))
	}
}

func (v *vm) executeShiftArithmetic(op bytecode.OpCode) {
	shiftVal := v.pop()
	n := toShiftAmount(shiftVal)
	a := v.pop()
	switch op {
	case bytecode.OpLsh:
		v.push(a.Lsh(n))
	case bytecode.OpRsh:
		v.push(a.Rsh(n))
	}
}

func (v *vm) popBinaryValues() (value.Value, value.Value) {
	b := v.pop()
	a := v.pop()
	return a, b
}
