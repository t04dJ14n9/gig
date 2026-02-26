// Package math registers the Go standard library math package.
package math

import (
	"math"

	"gig/importer"
	"gig/value"
)

func init() {
	pkg := importer.RegisterPackage("math", "math")

	// Constants
	pkg.AddConstant("E", math.E, "e")
	pkg.AddConstant("Pi", math.Pi, "π")
	pkg.AddConstant("Phi", math.Phi, "golden ratio")
	pkg.AddConstant("Sqrt2", math.Sqrt2, "√2")
	pkg.AddConstant("SqrtE", math.SqrtE, "√e")
	pkg.AddConstant("SqrtPi", math.SqrtPi, "√π")
	pkg.AddConstant("SqrtPhi", math.SqrtPhi, "√φ")
	pkg.AddConstant("Ln2", math.Ln2, "ln(2)")
	pkg.AddConstant("Log2E", math.Log2E, "log₂(e)")
	pkg.AddConstant("Ln10", math.Ln10, "ln(10)")
	pkg.AddConstant("Log10E", math.Log10E, "log₁₀(e)")
	pkg.AddConstant("MaxFloat32", math.MaxFloat32, "max float32")
	pkg.AddConstant("SmallestNonzeroFloat32", math.SmallestNonzeroFloat32, "smallest non-zero float32")
	pkg.AddConstant("MaxFloat64", math.MaxFloat64, "max float64")
	pkg.AddConstant("SmallestNonzeroFloat64", math.SmallestNonzeroFloat64, "smallest non-zero float64")
	pkg.AddConstant("MaxInt", math.MaxInt, "max int")
	pkg.AddConstant("MinInt", math.MinInt, "min int")
	pkg.AddConstant("MaxInt8", math.MaxInt8, "max int8")
	pkg.AddConstant("MinInt8", math.MinInt8, "min int8")
	pkg.AddConstant("MaxInt16", math.MaxInt16, "max int16")
	pkg.AddConstant("MinInt16", math.MinInt16, "min int16")
	pkg.AddConstant("MaxInt32", math.MaxInt32, "max int32")
	pkg.AddConstant("MinInt32", math.MinInt32, "min int32")
	pkg.AddConstant("MaxInt64", math.MaxInt64, "max int64")
	pkg.AddConstant("MinInt64", math.MinInt64, "min int64")
	// MaxUint and MaxUint64 are too large for int, use uint64 directly
	pkg.AddConstant("MaxUint8", math.MaxUint8, "max uint8")
	pkg.AddConstant("MaxUint16", math.MaxUint16, "max uint16")
	pkg.AddConstant("MaxUint32", math.MaxUint32, "max uint32")
	// MaxUint64 would overflow

	// Basic functions
	pkg.AddFunction("Abs", math.Abs, "", directAbs)
	pkg.AddFunction("Signbit", math.Signbit, "", directSignbit)
	pkg.AddFunction("Copysign", math.Copysign, "", directCopysign)

	// Trigonometric functions
	pkg.AddFunction("Sin", math.Sin, "", directSin)
	pkg.AddFunction("Cos", math.Cos, "", directCos)
	pkg.AddFunction("Tan", math.Tan, "", directTan)
	pkg.AddFunction("Asin", math.Asin, "", directAsin)
	pkg.AddFunction("Acos", math.Acos, "", directAcos)
	pkg.AddFunction("Atan", math.Atan, "", directAtan)
	pkg.AddFunction("Atan2", math.Atan2, "", directAtan2)
	pkg.AddFunction("Sinh", math.Sinh, "", directSinh)
	pkg.AddFunction("Cosh", math.Cosh, "", directCosh)
	pkg.AddFunction("Tanh", math.Tanh, "", directTanh)
	pkg.AddFunction("Asinh", math.Asinh, "", directAsinh)
	pkg.AddFunction("Acosh", math.Acosh, "", directAcosh)
	pkg.AddFunction("Atanh", math.Atanh, "", directAtanh)

	// Exponential and logarithmic functions
	pkg.AddFunction("Exp", math.Exp, "", directExp)
	pkg.AddFunction("Exp2", math.Exp2, "", directExp2)
	pkg.AddFunction("Expm1", math.Expm1, "", directExpm1)
	pkg.AddFunction("Log", math.Log, "", directLog)
	pkg.AddFunction("Log10", math.Log10, "", directLog10)
	pkg.AddFunction("Log2", math.Log2, "", directLog2)
	pkg.AddFunction("Log1p", math.Log1p, "", directLog1p)
	pkg.AddFunction("Logb", math.Logb, "", directLogb)
	pkg.AddFunction("Ilogb", math.Ilogb, "", directIlogb)

	// Power functions
	pkg.AddFunction("Pow", math.Pow, "", directPow)
	pkg.AddFunction("Pow10", math.Pow10, "", directPow10)
	pkg.AddFunction("Sqrt", math.Sqrt, "", directSqrt)
	pkg.AddFunction("Cbrt", math.Cbrt, "", directCbrt)

	// Rounding functions
	pkg.AddFunction("Ceil", math.Ceil, "", directCeil)
	pkg.AddFunction("Floor", math.Floor, "", directFloor)
	pkg.AddFunction("Trunc", math.Trunc, "", directTrunc)
	pkg.AddFunction("Round", math.Round, "", directRound)
	pkg.AddFunction("RoundToEven", math.RoundToEven, "", directRoundToEven)
	pkg.AddFunction("Mod", math.Mod, "", directMod)
	pkg.AddFunction("Modf", math.Modf, "", directModf)

	// Min/Max
	pkg.AddFunction("Min", math.Min, "", directMin)
	pkg.AddFunction("Max", math.Max, "", directMax)

	// Dim and hypot
	pkg.AddFunction("Dim", math.Dim, "", directDim)
	pkg.AddFunction("Hypot", math.Hypot, "", directHypot)

	// Special values
	pkg.AddFunction("Inf", math.Inf, "", directInf)
	pkg.AddFunction("NaN", math.NaN, "", directNaN)
	pkg.AddFunction("IsInf", math.IsInf, "", directIsInf)
	pkg.AddFunction("IsNaN", math.IsNaN, "", directIsNaN)

	// Nextafter
	pkg.AddFunction("Nextafter", math.Nextafter, "", directNextafter)
	pkg.AddFunction("Nextafter32", math.Nextafter32, "", nil)
}

// Direct wrappers for common functions

func directAbs(args []value.Value) value.Value {
	return value.MakeFloat(math.Abs(args[0].Float()))
}

func directSignbit(args []value.Value) value.Value {
	return value.MakeBool(math.Signbit(args[0].Float()))
}

func directCopysign(args []value.Value) value.Value {
	return value.MakeFloat(math.Copysign(args[0].Float(), args[1].Float()))
}

func directSin(args []value.Value) value.Value {
	return value.MakeFloat(math.Sin(args[0].Float()))
}

func directCos(args []value.Value) value.Value {
	return value.MakeFloat(math.Cos(args[0].Float()))
}

func directTan(args []value.Value) value.Value {
	return value.MakeFloat(math.Tan(args[0].Float()))
}

func directAsin(args []value.Value) value.Value {
	return value.MakeFloat(math.Asin(args[0].Float()))
}

func directAcos(args []value.Value) value.Value {
	return value.MakeFloat(math.Acos(args[0].Float()))
}

func directAtan(args []value.Value) value.Value {
	return value.MakeFloat(math.Atan(args[0].Float()))
}

func directAtan2(args []value.Value) value.Value {
	return value.MakeFloat(math.Atan2(args[0].Float(), args[1].Float()))
}

func directSinh(args []value.Value) value.Value {
	return value.MakeFloat(math.Sinh(args[0].Float()))
}

func directCosh(args []value.Value) value.Value {
	return value.MakeFloat(math.Cosh(args[0].Float()))
}

func directTanh(args []value.Value) value.Value {
	return value.MakeFloat(math.Tanh(args[0].Float()))
}

func directAsinh(args []value.Value) value.Value {
	return value.MakeFloat(math.Asinh(args[0].Float()))
}

func directAcosh(args []value.Value) value.Value {
	return value.MakeFloat(math.Acosh(args[0].Float()))
}

func directAtanh(args []value.Value) value.Value {
	return value.MakeFloat(math.Atanh(args[0].Float()))
}

func directExp(args []value.Value) value.Value {
	return value.MakeFloat(math.Exp(args[0].Float()))
}

func directExp2(args []value.Value) value.Value {
	return value.MakeFloat(math.Exp2(args[0].Float()))
}

func directExpm1(args []value.Value) value.Value {
	return value.MakeFloat(math.Expm1(args[0].Float()))
}

func directLog(args []value.Value) value.Value {
	return value.MakeFloat(math.Log(args[0].Float()))
}

func directLog10(args []value.Value) value.Value {
	return value.MakeFloat(math.Log10(args[0].Float()))
}

func directLog2(args []value.Value) value.Value {
	return value.MakeFloat(math.Log2(args[0].Float()))
}

func directLog1p(args []value.Value) value.Value {
	return value.MakeFloat(math.Log1p(args[0].Float()))
}

func directLogb(args []value.Value) value.Value {
	return value.MakeFloat(math.Logb(args[0].Float()))
}

func directIlogb(args []value.Value) value.Value {
	return value.MakeInt(int64(math.Ilogb(args[0].Float())))
}

func directPow(args []value.Value) value.Value {
	return value.MakeFloat(math.Pow(args[0].Float(), args[1].Float()))
}

func directPow10(args []value.Value) value.Value {
	return value.MakeFloat(math.Pow10(int(args[0].Int())))
}

func directSqrt(args []value.Value) value.Value {
	return value.MakeFloat(math.Sqrt(args[0].Float()))
}

func directCbrt(args []value.Value) value.Value {
	return value.MakeFloat(math.Cbrt(args[0].Float()))
}

func directCeil(args []value.Value) value.Value {
	return value.MakeFloat(math.Ceil(args[0].Float()))
}

func directFloor(args []value.Value) value.Value {
	return value.MakeFloat(math.Floor(args[0].Float()))
}

func directTrunc(args []value.Value) value.Value {
	return value.MakeFloat(math.Trunc(args[0].Float()))
}

func directRound(args []value.Value) value.Value {
	return value.MakeFloat(math.Round(args[0].Float()))
}

func directRoundToEven(args []value.Value) value.Value {
	return value.MakeFloat(math.RoundToEven(args[0].Float()))
}

func directMod(args []value.Value) value.Value {
	return value.MakeFloat(math.Mod(args[0].Float(), args[1].Float()))
}

func directModf(args []value.Value) value.Value {
	intPart, fracPart := math.Modf(args[0].Float())
	return value.FromInterface([]float64{intPart, fracPart})
}

func directMin(args []value.Value) value.Value {
	return value.MakeFloat(math.Min(args[0].Float(), args[1].Float()))
}

func directMax(args []value.Value) value.Value {
	return value.MakeFloat(math.Max(args[0].Float(), args[1].Float()))
}

func directDim(args []value.Value) value.Value {
	return value.MakeFloat(math.Dim(args[0].Float(), args[1].Float()))
}

func directHypot(args []value.Value) value.Value {
	return value.MakeFloat(math.Hypot(args[0].Float(), args[1].Float()))
}

func directInf(args []value.Value) value.Value {
	return value.MakeFloat(math.Inf(int(args[0].Int())))
}

func directNaN(args []value.Value) value.Value {
	return value.MakeFloat(math.NaN())
}

func directIsInf(args []value.Value) value.Value {
	return value.MakeBool(math.IsInf(args[0].Float(), int(args[1].Int())))
}

func directIsNaN(args []value.Value) value.Value {
	return value.MakeBool(math.IsNaN(args[0].Float()))
}

func directNextafter(args []value.Value) value.Value {
	return value.MakeFloat(math.Nextafter(args[0].Float(), args[1].Float()))
}
