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
	switch op { //nolint:exhaustive
	case bytecode.OpAddLocalLocal, bytecode.OpSubLocalLocal, bytecode.OpMulLocalLocal,
		bytecode.OpAddLocalConst, bytecode.OpSubLocalConst:
		return v.runArithmeticSuperinstruction(op, frame, sp, locals, prebaked)
	case bytecode.OpLessLocalLocalJumpTrue, bytecode.OpLessLocalConstJumpTrue,
		bytecode.OpLessEqLocalConstJumpTrue, bytecode.OpGreaterLocalLocalJumpTrue,
		bytecode.OpLessLocalLocalJumpFalse, bytecode.OpLessLocalConstJumpFalse,
		bytecode.OpLessEqLocalConstJumpFalse:
		return runComparisonJumpSuperinstruction(op, frame, sp, locals, prebaked)
	case bytecode.OpAddSetLocal, bytecode.OpSubSetLocal,
		bytecode.OpLocalLocalAddSetLocal, bytecode.OpLocalConstAddSetLocal,
		bytecode.OpLocalConstSubSetLocal, bytecode.OpLocalLocalSubSetLocal,
		bytecode.OpLocalLocalMulSetLocal, bytecode.OpLocalConstMulSetLocal:
		return v.runSetLocalSuperinstruction(op, frame, sp, locals, intLocals, prebaked)
	default:
		return sp
	}
}

func (v *vm) runArithmeticSuperinstruction(
	op bytecode.OpCode,
	frame *Frame,
	sp int,
	locals []value.Value,
	prebaked []value.Value,
) int {
	stack := v.stack
	switch op { //nolint:exhaustive
	case bytecode.OpAddLocalLocal:
		idxA, idxB := frame.readUint16(), frame.readUint16()
		stack[sp] = addValues(locals[idxA], locals[idxB])
	case bytecode.OpSubLocalLocal:
		idxA, idxB := frame.readUint16(), frame.readUint16()
		stack[sp] = subValues(locals[idxA], locals[idxB])
	case bytecode.OpMulLocalLocal:
		idxA, idxB := frame.readUint16(), frame.readUint16()
		stack[sp] = mulValues(locals[idxA], locals[idxB])
	case bytecode.OpAddLocalConst:
		idxA, idxB := frame.readUint16(), frame.readUint16()
		stack[sp] = addValues(locals[idxA], prebaked[idxB])
	case bytecode.OpSubLocalConst:
		idxA, idxB := frame.readUint16(), frame.readUint16()
		stack[sp] = subValues(locals[idxA], prebaked[idxB])
	}
	return sp + 1
}

func runComparisonJumpSuperinstruction(
	op bytecode.OpCode,
	frame *Frame,
	sp int,
	locals []value.Value,
	prebaked []value.Value,
) int {
	switch op { //nolint:exhaustive
	case bytecode.OpLessLocalLocalJumpTrue:
		idxA, idxB, offset := readBinaryJumpOperands(frame)
		jumpIf(frame, offset, lessValues(locals[idxA], locals[idxB]))
	case bytecode.OpLessLocalConstJumpTrue:
		idxA, idxB, offset := readBinaryJumpOperands(frame)
		jumpIf(frame, offset, lessValues(locals[idxA], prebaked[idxB]))
	case bytecode.OpLessEqLocalConstJumpTrue:
		idxA, idxB, offset := readBinaryJumpOperands(frame)
		jumpIf(frame, offset, lessEqValues(locals[idxA], prebaked[idxB]))
	case bytecode.OpGreaterLocalLocalJumpTrue:
		idxA, idxB, offset := readBinaryJumpOperands(frame)
		jumpIf(frame, offset, greaterValues(locals[idxA], locals[idxB]))
	case bytecode.OpLessLocalLocalJumpFalse:
		idxA, idxB, offset := readBinaryJumpOperands(frame)
		jumpIf(frame, offset, greaterEqValues(locals[idxA], locals[idxB]))
	case bytecode.OpLessLocalConstJumpFalse:
		idxA, idxB, offset := readBinaryJumpOperands(frame)
		jumpIf(frame, offset, greaterEqValues(locals[idxA], prebaked[idxB]))
	case bytecode.OpLessEqLocalConstJumpFalse:
		idxA, idxB, offset := readBinaryJumpOperands(frame)
		jumpIf(frame, offset, greaterValues(locals[idxA], prebaked[idxB]))
	}
	return sp
}

func (v *vm) runSetLocalSuperinstruction(
	op bytecode.OpCode,
	frame *Frame,
	sp int,
	locals []value.Value,
	intLocals []int64,
	prebaked []value.Value,
) int {
	stack := v.stack
	switch op { //nolint:exhaustive
	case bytecode.OpAddSetLocal:
		idx := frame.readUint16()
		var a, b value.Value
		sp, a, b = popBinaryOperands(stack, sp)
		locals[idx] = addValues(a, b)
		cacheIntLocal(intLocals, idx, locals[idx], a, b)
	case bytecode.OpSubSetLocal:
		idx := frame.readUint16()
		var a, b value.Value
		sp, a, b = popBinaryOperands(stack, sp)
		locals[idx] = subValues(a, b)
		cacheIntLocal(intLocals, idx, locals[idx], a, b)
	case bytecode.OpLocalLocalAddSetLocal:
		idxA, idxB, idxC := readTernaryLocalOperands(frame)
		locals[idxC] = addValues(locals[idxA], locals[idxB])
	case bytecode.OpLocalConstAddSetLocal:
		idxA, idxB, idxC := readTernaryLocalOperands(frame)
		locals[idxC] = addValues(locals[idxA], prebaked[idxB])
	case bytecode.OpLocalConstSubSetLocal:
		idxA, idxB, idxC := readTernaryLocalOperands(frame)
		locals[idxC] = subValues(locals[idxA], prebaked[idxB])
	case bytecode.OpLocalLocalSubSetLocal:
		idxA, idxB, idxC := readTernaryLocalOperands(frame)
		locals[idxC] = subValues(locals[idxA], locals[idxB])
	case bytecode.OpLocalLocalMulSetLocal:
		idxA, idxB, idxC := readTernaryLocalOperands(frame)
		locals[idxC] = mulValues(locals[idxA], locals[idxB])
	case bytecode.OpLocalConstMulSetLocal:
		idxA, idxB, idxC := readTernaryLocalOperands(frame)
		locals[idxC] = mulValues(locals[idxA], prebaked[idxB])
	}
	return sp
}

func readBinaryJumpOperands(frame *Frame) (uint16, uint16, uint16) {
	return frame.readUint16(), frame.readUint16(), frame.readUint16()
}

func readTernaryLocalOperands(frame *Frame) (uint16, uint16, uint16) {
	return frame.readUint16(), frame.readUint16(), frame.readUint16()
}

func popBinaryOperands(stack []value.Value, sp int) (int, value.Value, value.Value) {
	sp--
	b := stack[sp]
	sp--
	a := stack[sp]
	return sp, a, b
}

func cacheIntLocal(intLocals []int64, idx uint16, result, a, b value.Value) {
	if intLocals != nil && a.Kind() == value.KindInt && b.Kind() == value.KindInt {
		intLocals[idx] = result.RawInt()
	}
}

func jumpIf(frame *Frame, offset uint16, condition bool) {
	if condition {
		frame.ip = int(offset)
	}
}

func lessValues(a, b value.Value) bool {
	if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
		return a.RawInt() < b.RawInt()
	}
	return a.Cmp(b) < 0
}

func lessEqValues(a, b value.Value) bool {
	if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
		return a.RawInt() <= b.RawInt()
	}
	return lessEqCmp(a, b)
}

func greaterValues(a, b value.Value) bool {
	if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
		return a.RawInt() > b.RawInt()
	}
	return a.Cmp(b) > 0
}

func greaterEqValues(a, b value.Value) bool {
	if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
		return a.RawInt() >= b.RawInt()
	}
	return a.Cmp(b) >= 0
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
