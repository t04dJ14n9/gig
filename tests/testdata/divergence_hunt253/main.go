package divergence_hunt253

import (
	"fmt"
)

// ============================================================================
// Round 253: Const arithmetic
// ============================================================================

const (
	Base     = 10
	Add      = Base + 5
	Sub      = Base - 3
	Mul      = Base * 2
	Div      = Base / 2
	Mod      = Base % 3
	Neg      = -Base
	BitAnd   = Base & 7
	BitOr    = Base | 5
	BitXor   = Base ^ 3
	BitLeft  = Base << 2
	BitRight = Base >> 1
)

// ConstAdd tests constant addition
func ConstAdd() string {
	return fmt.Sprintf("Add=%d", Add)
}

// ConstSub tests constant subtraction
func ConstSub() string {
	return fmt.Sprintf("Sub=%d", Sub)
}

// ConstMul tests constant multiplication
func ConstMul() string {
	return fmt.Sprintf("Mul=%d", Mul)
}

// ConstDiv tests constant division
func ConstDiv() string {
	return fmt.Sprintf("Div=%d", Div)
}

// ConstMod tests constant modulo
func ConstMod() string {
	return fmt.Sprintf("Mod=%d", Mod)
}

// ConstNeg tests constant negation
func ConstNeg() string {
	return fmt.Sprintf("Neg=%d", Neg)
}

// ConstBitwiseAnd tests constant bitwise AND
func ConstBitwiseAnd() string {
	return fmt.Sprintf("BitAnd=%d", BitAnd)
}

// ConstBitwiseOr tests constant bitwise OR
func ConstBitwiseOr() string {
	return fmt.Sprintf("BitOr=%d", BitOr)
}

// ConstBitwiseXor tests constant bitwise XOR
func ConstBitwiseXor() string {
	return fmt.Sprintf("BitXor=%d", BitXor)
}

// ConstBitShiftLeft tests constant left shift
func ConstBitShiftLeft() string {
	return fmt.Sprintf("BitLeft=%d", BitLeft)
}

// ConstBitShiftRight tests constant right shift
func ConstBitShiftRight() string {
	return fmt.Sprintf("BitRight=%d", BitRight)
}

// ConstComplexExpr tests complex constant expression
func ConstComplexExpr() string {
	const result = (10 + 5) * 2 - 8 / 4
	return fmt.Sprintf("result=%d", result)
}
