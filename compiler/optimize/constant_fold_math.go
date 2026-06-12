package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

func foldedBinaryBytes(op bytecode.OpCode, left, right any, constants *[]any) ([]byte, bool) {
	if folded, ok := foldInt64Binary(op, left, right); ok {
		return constBytes(appendConstant(constants, folded)), true
	}
	if folded, ok := foldIntBinary(op, left, right); ok {
		return constBytes(appendConstant(constants, folded)), true
	}
	if folded, ok := foldBoolBinary(op, left, right); ok {
		return boolBytes(folded), true
	}
	return nil, false
}

func foldInt64Binary(op bytecode.OpCode, left, right any) (any, bool) {
	l, ok := left.(int64)
	if !ok {
		return nil, false
	}
	r, ok := right.(int64)
	if !ok {
		return nil, false
	}
	return foldSignedBinary(op, l, r)
}

func foldIntBinary(op bytecode.OpCode, left, right any) (any, bool) {
	l, ok := left.(int)
	if !ok {
		return nil, false
	}
	r, ok := right.(int)
	if !ok {
		return nil, false
	}
	folded, ok := foldSignedBinary(op, int64(l), int64(r))
	if !ok {
		return nil, false
	}
	if asInt, ok := folded.(int64); ok {
		return int(asInt), true
	}
	return folded, true
}

func foldSignedBinary(op bytecode.OpCode, left, right int64) (any, bool) {
	switch op {
	case bytecode.OpAdd:
		return left + right, true
	case bytecode.OpSub:
		return left - right, true
	case bytecode.OpMul:
		return left * right, true
	case bytecode.OpDiv:
		return divideSigned(left, right)
	case bytecode.OpMod:
		return modSigned(left, right)
	case bytecode.OpEqual:
		return left == right, true
	case bytecode.OpNotEqual:
		return left != right, true
	case bytecode.OpLess:
		return left < right, true
	case bytecode.OpLessEq:
		return left <= right, true
	case bytecode.OpGreater:
		return left > right, true
	case bytecode.OpGreaterEq:
		return left >= right, true
	default:
		return nil, false
	}
}

func divideSigned(left, right int64) (any, bool) {
	if right == 0 {
		return nil, false
	}
	return left / right, true
}

func modSigned(left, right int64) (any, bool) {
	if right == 0 {
		return nil, false
	}
	return left % right, true
}

func foldBoolBinary(op bytecode.OpCode, left, right any) (bool, bool) {
	l, ok := left.(bool)
	if !ok {
		return false, false
	}
	r, ok := right.(bool)
	if !ok {
		return false, false
	}
	switch op {
	case bytecode.OpEqual:
		return l == r, true
	case bytecode.OpNotEqual:
		return l != r, true
	default:
		return false, false
	}
}

func isFoldableConstBinaryOp(op bytecode.OpCode) bool {
	switch op {
	case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod,
		bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpLess, bytecode.OpLessEq,
		bytecode.OpGreater, bytecode.OpGreaterEq:
		return true
	default:
		return false
	}
}
