package vm

import (
	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/value"
)

// executeArithmetic handles arithmetic, bitwise, comparison, and logical opcodes.
func (v *vm) executeArithmetic(op bytecode.OpCode, frame *Frame) error { //nolint:cyclop
	switch op {
	// Arithmetic
	case bytecode.OpAdd:
		b := v.pop()
		a := v.pop()
		// Fast path for int+int (most common case in loops)
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			v.push(value.MakeIntSized(a.RawInt()+b.RawInt(), a.RawSize()))
		} else {
			v.push(a.Add(b))
		}

	case bytecode.OpSub:
		b := v.pop()
		a := v.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			v.push(value.MakeIntSized(a.RawInt()-b.RawInt(), a.RawSize()))
		} else {
			v.push(a.Sub(b))
		}

	case bytecode.OpMul:
		b := v.pop()
		a := v.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			v.push(value.MakeIntSized(a.RawInt()*b.RawInt(), a.RawSize()))
		} else {
			v.push(a.Mul(b))
		}

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
		var n uint
		if shiftVal.Kind() == value.KindUint {
			n = uint(shiftVal.Uint())
		} else {
			n = uint(shiftVal.Int())
		}
		a := v.pop()
		v.push(a.Lsh(n))

	case bytecode.OpRsh:
		shiftVal := v.pop()
		var n uint
		if shiftVal.Kind() == value.KindUint {
			n = uint(shiftVal.Uint())
		} else {
			n = uint(shiftVal.Int())
		}
		a := v.pop()
		v.push(a.Rsh(n))

	// Comparison
	case bytecode.OpEqual:
		b := v.pop()
		a := v.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			v.push(value.MakeBool(a.RawInt() == b.RawInt()))
		} else {
			v.push(value.MakeBool(a.Equal(b)))
		}

	case bytecode.OpNotEqual:
		b := v.pop()
		a := v.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			v.push(value.MakeBool(a.RawInt() != b.RawInt()))
		} else {
			v.push(value.MakeBool(!a.Equal(b)))
		}

	case bytecode.OpLess:
		b := v.pop()
		a := v.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			v.push(value.MakeBool(a.RawInt() < b.RawInt()))
		} else {
			v.push(value.MakeBool(a.Cmp(b) < 0))
		}

	case bytecode.OpLessEq:
		b := v.pop()
		a := v.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			v.push(value.MakeBool(a.RawInt() <= b.RawInt()))
		} else {
			v.push(value.MakeBool(a.Cmp(b) <= 0))
		}

	case bytecode.OpGreater:
		b := v.pop()
		a := v.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			v.push(value.MakeBool(a.RawInt() > b.RawInt()))
		} else {
			v.push(value.MakeBool(a.Cmp(b) > 0))
		}

	case bytecode.OpGreaterEq:
		b := v.pop()
		a := v.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			v.push(value.MakeBool(a.RawInt() >= b.RawInt()))
		} else {
			v.push(value.MakeBool(a.Cmp(b) >= 0))
		}

	// Logical
	case bytecode.OpNot:
		a := v.pop()
		v.push(value.MakeBool(!a.Bool()))
	}

	return nil
}
