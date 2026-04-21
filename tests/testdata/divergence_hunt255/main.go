package divergence_hunt255

import (
	"fmt"
)

// ============================================================================
// Round 255: Variable shadowing
// ============================================================================

var globalX = "global"

// ShadowSimple tests simple variable shadowing
func ShadowSimple() string {
	x := "outer"
	{
		x := "inner"
		_ = x
	}
	return fmt.Sprintf("x=%s", x)
}

// ShadowInLoop tests shadowing in loop
func ShadowInLoop() string {
	result := ""
	for i := 0; i < 3; i++ {
		i := i * 10
		result += fmt.Sprintf("%d,", i)
	}
	return result
}

// ShadowIfStatement tests shadowing in if statement
func ShadowIfStatement() string {
	x := 10
	if x := 20; x > 15 {
		_ = x
	}
	return fmt.Sprintf("x=%d", x)
}

// ShadowSwitchStatement tests shadowing in switch
func ShadowSwitchStatement() string {
	x := 1
	switch x := 5; x {
	case 5:
		_ = x
	}
	return fmt.Sprintf("x=%d", x)
}

// ShadowForRange tests shadowing in range loop
func ShadowForRange() string {
	nums := []int{1, 2, 3}
	result := ""
	for i, v := range nums {
		_, _ = i, v
		i, v := i*10, v*10
		result += fmt.Sprintf("(%d,%d)", i, v)
	}
	return result
}

// ShadowMultipleScopes tests multiple nested scopes
func ShadowMultipleScopes() string {
	x := "level1"
	{
		x := "level2"
		{
			x := "level3"
			_ = x
		}
		_ = x
	}
	return fmt.Sprintf("x=%s", x)
}

// ShadowFunctionParam tests shadowing function parameter
func ShadowFunctionParam(x int) string {
	x = x + 10
	return fmt.Sprintf("x=%d", x)
}

// ShadowWithShortDecl tests shadowing with short declaration
func ShadowWithShortDecl() string {
	x, y := 1, 2
	_ = y
	if true {
		x, z := 3, 4
		_ = x
		_ = z
	}
	return fmt.Sprintf("x=%d", x)
}

// ShadowGlobal tests shadowing global variable
func ShadowGlobal() string {
	globalX := "local"
	return fmt.Sprintf("local=%s", globalX)
}

// ShadowSameType tests shadowing with same type
func ShadowSameType() string {
	x := 100
	{
		x := 200
		_ = x
	}
	return fmt.Sprintf("x=%d", x)
}

// ShadowDifferentType tests shadowing with different type
func ShadowDifferentType() string {
	x := 100
	{
		x := "string"
		_ = x
	}
	return fmt.Sprintf("x=%d", x)
}
