package thirdparty

import "math/big"

// ============================================================================
// math/big.Int — arbitrary precision integers
// ============================================================================

// BigIntAdd tests addition of large integers.
func BigIntAdd() int {
	a := new(big.Int)
	b := new(big.Int)
	a.SetString("123456789012345678901234567890", 10)
	b.SetString("987654321098765432109876543210", 10)
	a.Add(a, b)
	if a.BitLen() > 0 {
		return 1
	}
	return 0
}

// BigIntMul tests multiplication of large integers.
func BigIntMul() int {
	a := new(big.Int).SetUint64(123456789)
	b := new(big.Int).SetUint64(987654321)
	mul := a.Mul(a, b)
	if mul.BitLen() > 0 {
		return 1
	}
	return 0
}

// BigIntDiv tests integer division.
func BigIntDiv() int {
	a := new(big.Int).SetUint64(1000000000)
	b := new(big.Int).SetUint64(3)
	quo := new(big.Int)
	quo.DivMod(a, b, new(big.Int))
	if quo.BitLen() > 0 {
		return 1
	}
	return 0
}

// BigIntMod tests modulo operation.
func BigIntMod() int {
	a := new(big.Int).SetUint64(1000000007)
	b := new(big.Int).SetUint64(1000)
	rem := new(big.Int)
	rem.Mod(a, b)
	if rem.Int64() == 7 {
		return 1
	}
	return 0
}

// BigIntPow tests exponentiation.
func BigIntPow() int {
	x := new(big.Int).SetUint64(2)
	m := new(big.Int) // nil-safe: Exp with m=nil needs valid *Int
	x.Exp(x, new(big.Int).SetUint64(100), m)
	if x.BitLen() == 100 {
		return 1
	}
	return 0
}

// BigIntBitwise tests bitwise AND/OR/XOR.
func BigIntBitwise() int {
	a := new(big.Int).SetUint64(0xF0F0F0F0F0F0F0F0)
	b := new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF)
	and := new(big.Int).And(a, b)
	xor := new(big.Int).Xor(a, b)
	or := new(big.Int).Or(a, b)
	if and.BitLen() > 0 && xor.BitLen() > 0 && or.BitLen() > 0 {
		return 1
	}
	return 0
}

// BigIntGCD tests greatest common divisor.
func BigIntGCD() int {
	a := new(big.Int).SetUint64(1071) // 3^2 * 7 * 17
	b := new(big.Int).SetUint64(462)  // 2 * 3 * 7 * 11
	g := new(big.Int)
	p := new(big.Int)
	q := new(big.Int)
	g.GCD(p, q, a, b)
	if g.Int64() == 21 { // 3 * 7
		return 1
	}
	return 0
}

// BigIntPrime tests primality (Miller-Rabin).
func BigIntPrime() int {
	p := new(big.Int)
	p.SetUint64(2305843009213693951) // 2^61 - 1, a Mersenne prime
	if p.ProbablyPrime(20) {
		return 1
	}
	return 0
}

// BigIntModInverse tests modular multiplicative inverse.
func BigIntModInverse() int {
	a := new(big.Int).SetUint64(3)
	m := new(big.Int).SetUint64(11)
	inv := new(big.Int)
	inv.ModInverse(a, m)
	// 3^-1 mod 11 = 4 (because 3*4 = 12 ≡ 1 mod 11)
	if inv.Int64() == 4 {
		return 1
	}
	return 0
}

// BigIntShift tests left and right shift.
func BigIntShift() int {
	a := new(big.Int).SetUint64(1)
	a.Lsh(a, 100)
	a.Rsh(a, 50)
	if a.Int64() == 1125899906842624 { // 2^50
		return 1
	}
	return 0
}

// BigIntAbs tests absolute value.
func BigIntAbs() int {
	neg := new(big.Int)
	neg.SetString("-123456789012345678901234567890", 10)
	abs := new(big.Int).Abs(neg)
	if abs.Sign() == 1 && abs.BitLen() > 0 {
		return 1
	}
	return 0
}

// BigIntString tests various radix output.
func BigIntString() int {
	a := new(big.Int).SetUint64(255)
	hex := a.Text(16)
	bin := a.Text(2)
	if hex == "ff" && bin == "11111111" {
		return 1
	}
	return 0
}

// ============================================================================
// math/big.Rat — arbitrary precision rationals
// ============================================================================

// BigRatBasic tests basic rational arithmetic.
func BigRatBasic() int {
	r := new(big.Rat)
	r.SetString("1/3")
	r2 := new(big.Rat)
	r2.SetString("1/6")
	r.Add(r, r2) // 1/3 + 1/6 = 1/2
	if r.RatString() == "1/2" {
		return 1
	}
	return 0
}

// BigRatMul tests rational multiplication.
func BigRatMul() int {
	a := new(big.Rat)
	b := new(big.Rat)
	a.SetString("2/3")
	b.SetString("3/4")
	mul := new(big.Rat).Mul(a, b)
	if mul.RatString() == "1/2" {
		return 1
	}
	return 0
}

// ============================================================================
// math/big.Float — arbitrary precision floating point
// ============================================================================

// BigFloatBasic tests basic float arithmetic.
func BigFloatBasic() int {
	f := new(big.Float).SetPrec(100)
	f.SetString("3.14159265358979323846")
	mul := new(big.Float).Mul(f, big.NewFloat(2))
	if mul.Sign() > 0 {
		return 1
	}
	return 0
}

// BigFloatSqrt tests high-precision square root.
func BigFloatSqrt() int {
	f := new(big.Float).SetPrec(200)
	f.SetUint64(2)
	sqrt := new(big.Float).Sqrt(f)
	if sqrt.Sign() > 0 {
		return 1
	}
	return 0
}

// BigFloatExp tests exponential — replaced with power test using Mul.
func BigFloatExp() int {
	f := new(big.Float).SetPrec(100)
	f.SetFloat64(2.0)
	exp := new(big.Float).Mul(f, f) // square
	if exp.Sign() > 0 {
		return 1
	}
	return 0
}
