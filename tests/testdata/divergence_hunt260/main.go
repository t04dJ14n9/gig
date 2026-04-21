package divergence_hunt260

import (
	"fmt"
)

// ============================================================================
// Round 260: Package-level variable initialization order
// ============================================================================

var (
	a = b + 1 // 3
	b = c + 1 // 2
	c = 1     // 1
)

// InitOrderDependency tests initialization with dependencies
func InitOrderDependency() string {
	return fmt.Sprintf("a=%d,b=%d,c=%d", a, b, c)
}

// InitOrderReverse tests reverse initialization order
var (
	x = 1
	y = x + 1
	z = y + 1
)

func InitOrderReverse() string {
	return fmt.Sprintf("x=%d,y=%d,z=%d", x, y, z)
}

// InitWithFunctionCall tests init with function calls
func getValue() int {
	return 42
}

var funcInit = getValue()

func InitWithFunctionCall() string {
	return fmt.Sprintf("funcInit=%d", funcInit)
}

// InitWithExpression tests init with expressions
var (
	expr1 = 10 * 2
	expr2 = expr1 + 5
	expr3 = expr2 / 5
)

func InitWithExpression() string {
	return fmt.Sprintf("expr1=%d,expr2=%d,expr3=%d", expr1, expr2, expr3)
}

// InitWithStringConcat tests init with string operations
var (
	str1 = "Hello"
	str2 = str1 + " World"
	str3 = str2 + "!"
)

func InitWithStringConcat() string {
	return fmt.Sprintf("str3=%s", str3)
}

// InitWithSliceLiteral tests init with slice literals
var sliceInit = []int{1, 2, 3, 4, 5}

func InitWithSliceLiteral() string {
	return fmt.Sprintf("len=%d,sum=%d", len(sliceInit), sliceInit[0]+sliceInit[1]+sliceInit[2])
}

// InitWithMapLiteral tests init with map literals
var mapInit = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
}

func InitWithMapLiteral() string {
	return fmt.Sprintf("one=%d,two=%d", mapInit["one"], mapInit["two"])
}

// InitWithStructLiteral tests init with struct literals
type Point struct {
	X, Y int
}

var pointInit = Point{X: 10, Y: 20}

func InitWithStructLiteral() string {
	return fmt.Sprintf("Point(%d,%d)", pointInit.X, pointInit.Y)
}

// InitWithArrayLiteral tests init with array literals
var arrayInit = [3]int{100, 200, 300}

func InitWithArrayLiteral() string {
	return fmt.Sprintf("array=[%d,%d,%d]", arrayInit[0], arrayInit[1], arrayInit[2])
}

// InitMixedTypes tests init with mixed types
var (
	intVal    = 42
	floatVal  = 3.14
	stringVal = "test"
	boolVal   = true
)

func InitMixedTypes() string {
	return fmt.Sprintf("int=%d,float=%.2f,string=%s,bool=%v", intVal, floatVal, stringVal, boolVal)
}

// InitChainDependency tests chain of dependencies
var (
	chainA = 1
	chainB = chainA * 2
	chainC = chainB * 2
	chainD = chainC * 2
)

func InitChainDependency() string {
	return fmt.Sprintf("chain=%d,%d,%d,%d", chainA, chainB, chainC, chainD)
}
