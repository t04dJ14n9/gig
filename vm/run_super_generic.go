package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) runGenericSuperinstruction(
	op bytecode.OpCode,
	frame *Frame,
	sp int,
	locals []value.Value,
	intLocals []int64,
	prebaked []value.Value,
) int {
	stack := v.stack
	switch op { //nolint:exhaustive
	case bytecode.OpAddLocalLocal:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		stack[sp] = addValues(locals[idxA], locals[idxB])
		return sp + 1

	case bytecode.OpSubLocalLocal:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		stack[sp] = subValues(locals[idxA], locals[idxB])
		return sp + 1

	case bytecode.OpMulLocalLocal:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		stack[sp] = mulValues(locals[idxA], locals[idxB])
		return sp + 1

	case bytecode.OpAddLocalConst:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		stack[sp] = addValues(locals[idxA], prebaked[idxB])
		return sp + 1

	case bytecode.OpSubLocalConst:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		stack[sp] = subValues(locals[idxA], prebaked[idxB])
		return sp + 1

	case bytecode.OpLessLocalLocalJumpTrue:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		offset := frame.readUint16()
		a, b := locals[idxA], locals[idxB]
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			if a.RawInt() < b.RawInt() {
				frame.ip = int(offset)
			}
		} else if a.Cmp(b) < 0 {
			frame.ip = int(offset)
		}
		return sp

	case bytecode.OpLessLocalConstJumpTrue:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		offset := frame.readUint16()
		a, b := locals[idxA], prebaked[idxB]
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			if a.RawInt() < b.RawInt() {
				frame.ip = int(offset)
			}
		} else if a.Cmp(b) < 0 {
			frame.ip = int(offset)
		}
		return sp

	case bytecode.OpLessEqLocalConstJumpTrue:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		offset := frame.readUint16()
		a, b := locals[idxA], prebaked[idxB]
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			if a.RawInt() <= b.RawInt() {
				frame.ip = int(offset)
			}
		} else if lessEqCmp(a, b) {
			frame.ip = int(offset)
		}
		return sp

	case bytecode.OpGreaterLocalLocalJumpTrue:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		offset := frame.readUint16()
		a, b := locals[idxA], locals[idxB]
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			if a.RawInt() > b.RawInt() {
				frame.ip = int(offset)
			}
		} else if a.Cmp(b) > 0 {
			frame.ip = int(offset)
		}
		return sp

	case bytecode.OpLessLocalLocalJumpFalse:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		offset := frame.readUint16()
		a, b := locals[idxA], locals[idxB]
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			if a.RawInt() >= b.RawInt() {
				frame.ip = int(offset)
			}
		} else if a.Cmp(b) >= 0 {
			frame.ip = int(offset)
		}
		return sp

	case bytecode.OpLessLocalConstJumpFalse:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		offset := frame.readUint16()
		a, b := locals[idxA], prebaked[idxB]
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			if a.RawInt() >= b.RawInt() {
				frame.ip = int(offset)
			}
		} else if a.Cmp(b) >= 0 {
			frame.ip = int(offset)
		}
		return sp

	case bytecode.OpLessEqLocalConstJumpFalse:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		offset := frame.readUint16()
		a, b := locals[idxA], prebaked[idxB]
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			if a.RawInt() > b.RawInt() {
				frame.ip = int(offset)
			}
		} else if !lessEqCmp(a, b) {
			frame.ip = int(offset)
		}
		return sp

	case bytecode.OpAddSetLocal:
		idx := frame.readUint16()
		sp--
		b := stack[sp]
		sp--
		a := stack[sp]
		locals[idx] = addValues(a, b)
		if intLocals != nil && a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			intLocals[idx] = locals[idx].RawInt()
		}
		return sp

	case bytecode.OpSubSetLocal:
		idx := frame.readUint16()
		sp--
		b := stack[sp]
		sp--
		a := stack[sp]
		locals[idx] = subValues(a, b)
		if intLocals != nil && a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			intLocals[idx] = locals[idx].RawInt()
		}
		return sp

	case bytecode.OpLocalLocalAddSetLocal:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		idxC := frame.readUint16()
		locals[idxC] = addValues(locals[idxA], locals[idxB])
		return sp

	case bytecode.OpLocalConstAddSetLocal:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		idxC := frame.readUint16()
		locals[idxC] = addValues(locals[idxA], prebaked[idxB])
		return sp

	case bytecode.OpLocalConstSubSetLocal:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		idxC := frame.readUint16()
		locals[idxC] = subValues(locals[idxA], prebaked[idxB])
		return sp

	case bytecode.OpLocalLocalSubSetLocal:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		idxC := frame.readUint16()
		locals[idxC] = subValues(locals[idxA], locals[idxB])
		return sp

	case bytecode.OpLocalLocalMulSetLocal:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		idxC := frame.readUint16()
		locals[idxC] = mulValues(locals[idxA], locals[idxB])
		return sp

	case bytecode.OpLocalConstMulSetLocal:
		idxA := frame.readUint16()
		idxB := frame.readUint16()
		idxC := frame.readUint16()
		locals[idxC] = mulValues(locals[idxA], prebaked[idxB])
		return sp
	default:
		return sp
	}
}

func addValues(a, b value.Value) value.Value {
	if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
		return value.MakeIntSized(a.RawInt()+b.RawInt(), a.RawSize())
	}
	return a.Add(b)
}

func subValues(a, b value.Value) value.Value {
	if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
		return value.MakeIntSized(a.RawInt()-b.RawInt(), a.RawSize())
	}
	return a.Sub(b)
}

func mulValues(a, b value.Value) value.Value {
	if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
		return value.MakeIntSized(a.RawInt()*b.RawInt(), a.RawSize())
	}
	return a.Mul(b)
}
