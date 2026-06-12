package value

import "fmt"

// bitwiseOp applies a bitwise operation to int/uint values.
func (v Value) bitwiseOp(other Value, intOp func(int64, int64) int64, uintOp func(uint64, uint64) uint64, name string) Value {
	switch v.kind {
	case KindInt:
		return MakeIntSized(intOp(v.num, other.Int()), v.size)
	case KindUint:
		return makeUintSized(uintOp(uint64(v.num), other.Uint()), v.size)
	default:
		panic(fmt.Sprintf("cannot %s %v", name, v.kind))
	}
}

// And returns v & other.
func (v Value) And(other Value) Value {
	return v.bitwiseOp(other, func(a, b int64) int64 { return a & b }, func(a, b uint64) uint64 { return a & b }, "and")
}

// Or returns v | other.
func (v Value) Or(other Value) Value {
	return v.bitwiseOp(other, func(a, b int64) int64 { return a | b }, func(a, b uint64) uint64 { return a | b }, "or")
}

// Xor returns v ^ other.
func (v Value) Xor(other Value) Value {
	return v.bitwiseOp(other, func(a, b int64) int64 { return a ^ b }, func(a, b uint64) uint64 { return a ^ b }, "xor")
}

// AndNot returns v &^ other.
func (v Value) AndNot(other Value) Value {
	return v.bitwiseOp(other, func(a, b int64) int64 { return a &^ b }, func(a, b uint64) uint64 { return a &^ b }, "andnot")
}

// Lsh returns v << n.
func (v Value) Lsh(n uint) Value {
	switch v.kind {
	case KindInt:
		return MakeIntSized(v.num<<n, v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)<<n, v.size)
	default:
		panic(fmt.Sprintf("cannot lsh %v", v.kind))
	}
}

// Rsh returns v >> n.
func (v Value) Rsh(n uint) Value {
	switch v.kind {
	case KindInt:
		return MakeIntSized(v.num>>n, v.size)
	case KindUint:
		return makeUintSized(uint64(v.num)>>n, v.size)
	default:
		panic(fmt.Sprintf("cannot rsh %v", v.kind))
	}
}
