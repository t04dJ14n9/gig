package value

import "math"

// truncateInt truncates an int64 to the given bit width, preserving sign.
// This implements Go's integer overflow wrapping semantics (two's complement).
func truncateInt(i int64, s Size) int64 {
	switch s {
	case Size8:
		return int64(int8(i))
	case Size16:
		return int64(int16(i))
	case Size32:
		return int64(int32(i))
	default:
		return i // Size64, SizePtr, Size0 -> no truncation
	}
}

// truncateUint truncates a uint64 to the given bit width.
func truncateUint(u uint64, s Size) uint64 {
	switch s {
	case Size8:
		return uint64(uint8(u))
	case Size16:
		return uint64(uint16(u))
	case Size32:
		return uint64(uint32(u))
	default:
		return u // Size64, SizePtr, Size0 -> no truncation
	}
}

// MakeIntSized creates an int value preserving the given size tag.
// The value is truncated to fit the signed integer range for the given size,
// implementing Go's two's complement wrapping semantics.
func MakeIntSized(i int64, s Size) Value {
	return Value{kind: KindInt, size: s, num: truncateInt(i, s)}
}

// makeUintSized creates a uint value preserving the given size tag.
// The value is truncated to fit the unsigned integer range for the given size.
func makeUintSized(u uint64, s Size) Value {
	return Value{kind: KindUint, size: s, num: int64(truncateUint(u, s))}
}

// makeFloatSized creates a float value preserving the given size tag.
func makeFloatSized(f float64, s Size) Value {
	return Value{kind: KindFloat, size: s, num: int64(math.Float64bits(f))}
}
