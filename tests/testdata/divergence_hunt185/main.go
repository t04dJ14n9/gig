package divergence_hunt185

import (
	"fmt"
	"math/big"
)

// ============================================================================
// Round 185: Math big integer operations
// ============================================================================

func NewInt() string {
	x := big.NewInt(42)
	return fmt.Sprintf("%s", x.String())
}

func IntAdd() string {
	x := big.NewInt(10)
	y := big.NewInt(20)
	z := new(big.Int).Add(x, y)
	return fmt.Sprintf("%s", z.String())
}

func IntSub() string {
	x := big.NewInt(50)
	y := big.NewInt(15)
	z := new(big.Int).Sub(x, y)
	return fmt.Sprintf("%s", z.String())
}

func IntMul() string {
	x := big.NewInt(6)
	y := big.NewInt(7)
	z := new(big.Int).Mul(x, y)
	return fmt.Sprintf("%s", z.String())
}

func IntDiv() string {
	x := big.NewInt(100)
	y := big.NewInt(4)
	z := new(big.Int).Div(x, y)
	return fmt.Sprintf("%s", z.String())
}

func IntMod() string {
	x := big.NewInt(17)
	y := big.NewInt(5)
	z := new(big.Int).Mod(x, y)
	return fmt.Sprintf("%s", z.String())
}

func IntAbs() string {
	x := big.NewInt(-42)
	z := new(big.Int).Abs(x)
	return fmt.Sprintf("%s", z.String())
}

func IntNeg() string {
	x := big.NewInt(42)
	z := new(big.Int).Neg(x)
	return fmt.Sprintf("%s", z.String())
}

func IntCmp() string {
	x := big.NewInt(10)
	y := big.NewInt(20)
	z := big.NewInt(10)
	return fmt.Sprintf("%d:%d:%d", x.Cmp(y), y.Cmp(x), x.Cmp(z))
}

func IntSet() string {
	x := big.NewInt(100)
	y := new(big.Int).Set(x)
	return fmt.Sprintf("%s", y.String())
}

func IntLargeNumber() string {
	x, _ := new(big.Int).SetString("12345678901234567890", 10)
	y := big.NewInt(2)
	z := new(big.Int).Mul(x, y)
	return fmt.Sprintf("%s", z.String())
}
