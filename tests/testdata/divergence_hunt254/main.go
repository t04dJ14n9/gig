package divergence_hunt254

import (
	"fmt"
)

// ============================================================================
// Round 254: Const overflow
// ============================================================================

// OverflowInt8Max tests int8 max boundary
func OverflowInt8Max() string {
	const maxInt8 = 127
	return fmt.Sprintf("maxInt8=%d", maxInt8)
}

// OverflowInt8Min tests int8 min boundary
func OverflowInt8Min() string {
	const minInt8 = -128
	return fmt.Sprintf("minInt8=%d", minInt8)
}

// OverflowUint8Max tests uint8 max boundary
func OverflowUint8Max() string {
	const maxUint8 = 255
	return fmt.Sprintf("maxUint8=%d", maxUint8)
}

// OverflowInt16Max tests int16 max boundary
func OverflowInt16Max() string {
	const maxInt16 = 32767
	return fmt.Sprintf("maxInt16=%d", maxInt16)
}

// OverflowInt32Max tests int32 max boundary
func OverflowInt32Max() string {
	const maxInt32 = 2147483647
	return fmt.Sprintf("maxInt32=%d", maxInt32)
}

// OverflowUint32Max tests uint32 max boundary
func OverflowUint32Max() string {
	const maxUint32 = 4294967295
	return fmt.Sprintf("maxUint32=%d", maxUint32)
}

// OverflowInt64Max tests int64 max boundary
func OverflowInt64Max() string {
	const maxInt64 = 9223372036854775807
	return fmt.Sprintf("maxInt64=%d", maxInt64)
}

// OverflowUint64Max tests uint64 max boundary
func OverflowUint64Max() string {
	const maxUint64 uint64 = 18446744073709551615
	return fmt.Sprintf("maxUint64=%d", maxUint64)
}

// OverflowFloat32Max tests float32 max boundary
func OverflowFloat32Max() string {
	const maxFloat32 float32 = 3.4028235e+38
	return fmt.Sprintf("maxFloat32=%e", maxFloat32)
}

// LargeShift tests large bit shift
func LargeShift() string {
	const shifted = 1 << 62
	return fmt.Sprintf("shifted=%d", shifted)
}

// UintOverflowBoundary tests uint overflow
func UintOverflowBoundary() string {
	const maxUint16 = 65535
	return fmt.Sprintf("maxUint16=%d", maxUint16)
}
