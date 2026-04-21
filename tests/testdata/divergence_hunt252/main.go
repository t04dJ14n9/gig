package divergence_hunt252

import (
	"fmt"
)

// ============================================================================
// Round 252: Const typed vs untyped
// ============================================================================

const untypedInt = 42
const untypedFloat = 3.14
const untypedString = "hello"
const untypedBool = true

const typedInt int = 42
const typedFloat64 float64 = 3.14
const typedString string = "hello"
const typedBool bool = true

// UntypedIntAssign tests untyped int assignment
func UntypedIntAssign() string {
	var a int = untypedInt
	var b int32 = untypedInt
	var c int64 = untypedInt
	var d float64 = untypedInt
	return fmt.Sprintf("a=%d,b=%d,c=%d,d=%.0f", a, b, c, d)
}

// UntypedFloatAssign tests untyped float assignment
func UntypedFloatAssign() string {
	var a float32 = untypedFloat
	var b float64 = untypedFloat
	return fmt.Sprintf("a=%.2f,b=%.2f", a, b)
}

// TypedIntAssign tests typed int restrictions
func TypedIntAssign() string {
	var a int = typedInt
	// var b int32 = typedInt  // Would not compile
	return fmt.Sprintf("a=%d", a)
}

// UntypedStringAssign tests untyped string
func UntypedStringAssign() string {
	var s string = untypedString
	return fmt.Sprintf("s=%s", s)
}

// UntypedBoolAssign tests untyped bool
func UntypedBoolAssign() string {
	var b bool = untypedBool
	return fmt.Sprintf("b=%v", b)
}

// UntypedArithmetic tests arithmetic with untyped constants
func UntypedArithmetic() string {
	const a = 10
	const b = 3
	const result = a / b  // integer division on untyped
	return fmt.Sprintf("result=%d,type=untyped", result)
}

// TypedArithmetic tests arithmetic with typed constants
func TypedArithmetic() string {
	const a int = 10
	const b int = 3
	const result = a / b
	return fmt.Sprintf("result=%d", result)
}

// MixedTypeArithmetic tests mixing typed and untyped
func MixedTypeArithmetic() string {
	const untyped = 100
	const typed int = 10
	result := untyped / typed
	return fmt.Sprintf("result=%d,type=int", result)
}

// DefaultTypeInference tests default type inference
func DefaultTypeInference() string {
	var a = untypedInt    // becomes int
	var b = untypedFloat  // becomes float64
	var c = untypedString // becomes string
	var d = untypedBool   // becomes bool
	return fmt.Sprintf("%T,%T,%T,%T", a, b, c, d)
}

// LargeUntypedInteger tests large untyped integer
func LargeUntypedInteger() string {
	const big = 1 << 60  // fits in untyped, would overflow int32
	var i64 int64 = big
	return fmt.Sprintf("i64=%d", i64)
}

// UntypedRune tests untyped rune constant
func UntypedRune() string {
	const r = 'A'
	var ru rune = r
	var i int = r
	return fmt.Sprintf("rune=%c,int=%d", ru, i)
}
