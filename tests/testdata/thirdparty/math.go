package thirdparty

import "math"

// MathAbs tests math.Abs.
func MathAbs() float64 {
	return math.Abs(-123.45)
}

// MathMax tests math.Max.
func MathMax() float64 {
	return math.Max(10.5, 20.3)
}

// MathMin tests math.Min.
func MathMin() float64 {
	return math.Min(10.5, 20.3)
}

// MathFloor tests math.Floor.
func MathFloor() float64 {
	return math.Floor(123.7)
}

// MathCeil tests math.Ceil.
func MathCeil() float64 {
	return math.Ceil(123.3)
}

// MathRound tests math.Round.
func MathRound() float64 {
	return math.Round(123.5)
}

// MathPow tests math.Pow.
func MathPow() float64 {
	return math.Pow(2, 10)
}

// MathSqrt tests math.Sqrt.
func MathSqrt() float64 {
	return math.Sqrt(144)
}

// MathMod tests math.Mod.
func MathMod() float64 {
	return math.Mod(10, 3)
}

// MathSin tests math.Sin.
func MathSin() float64 {
	return math.Sin(math.Pi / 2)
}

// MathCos tests math.Cos.
func MathCos() float64 {
	return math.Cos(0)
}

// MathTan tests math.Tan.
func MathTan() float64 {
	return math.Tan(0)
}

// MathLog tests math.Log.
func MathLog() float64 {
	return math.Log(math.E)
}

// MathLog10 tests math.Log10.
func MathLog10() float64 {
	return math.Log10(100)
}

// MathExp tests math.Exp.
func MathExp() float64 {
	return math.Exp(2)
}

// MathInf tests math.IsInf.
func MathInf() int {
	if math.IsInf(math.Inf(1), 1) {
		return 1
	}
	return 0
}

// MathNaN tests math.IsNaN.
func MathNaN() int {
	if math.IsNaN(math.NaN()) {
		return 1
	}
	return 0
}

// MathCopysign tests math.Copysign.
func MathCopysign() float64 {
	return math.Copysign(-5, 1)
}
