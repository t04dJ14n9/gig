package divergence_hunt42

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// ============================================================================
// Round 42: Numeric precision and conversion edge cases
// ============================================================================

func Float32AddPrecision() float32 {
	var a float32 = 0.1
	var b float32 = 0.2
	return a + b
}

func Float64ToIntTrunc() int {
	v := 9.99
	return int(v)
}

func Float64NegativeTrunc() int {
	v := -3.7
	return int(v)
}

func IntToFloatRoundTrip() float64 {
	x := 42
	return float64(x)
}

func LargeIntToFloat() float64 {
	x := int64(1<<53 + 1)
	return float64(x) // may lose precision
}

func Uint64ToFloat() float64 {
	var x uint64 = 9007199254740993
	return float64(x)
}

func Float32ToFloat64() float64 {
	var f32 float32 = 3.14
	return float64(f32)
}

func Float64ToFloat32() float32 {
	var f64 float64 = 3.14159265358979
	return float32(f64)
}

func StrconvParseFloat32() float64 {
	f, _ := strconv.ParseFloat("3.14", 32)
	return f
}

func StrconvParseFloat64() float64 {
	f, _ := strconv.ParseFloat("3.14159265358979", 64)
	return f
}

func MathRoundEven() float64 {
	return math.Round(2.5) + math.Round(3.5) // Go rounds half to even? No, Go rounds half away from zero
}

func MathRoundNegative() float64 {
	return math.Round(-2.5)
}

func FloatCompareNaN() bool {
	nan := math.NaN()
	return nan != nan
}

func FloatCompareInf() bool {
	inf := math.Inf(1)
	return inf > 1e308
}

func FloatNegativeZero() bool {
	return math.Signbit(-0.0) && !math.Signbit(0.0)
}

func FloatAddInf() bool {
	inf := math.Inf(1)
	return math.IsInf(inf+1, 1)
}

func FloatNaNCompare() bool {
	return !(math.NaN() < 0) && !(math.NaN() > 0) && !(math.NaN() == 0)
}

func Int8ToInt16Promotion() int16 {
	var a int8 = 10
	var b int8 = 20
	return int16(a) + int16(b)
}

func Uint8Addition() uint8 {
	var a uint8 = 200
	var b uint8 = 55
	return a + b // wraps to 255
}

func JSONFloatPrecision() float64 {
	type Data struct{ Value float64 }
	d := Data{Value: 3.14159265358979}
	b, _ := json.Marshal(d)
	var decoded Data
	json.Unmarshal(b, &decoded)
	return decoded.Value
}

func FmtFloatPrecision() string {
	return fmt.Sprintf("%.10f", math.Pi)
}
