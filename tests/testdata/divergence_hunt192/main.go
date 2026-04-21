package divergence_hunt192

import (
	"fmt"
)

// ============================================================================
// Round 192: Complex number operations
// ============================================================================

// ComplexAddition tests complex number addition
func ComplexAddition() string {
	a := complex(1, 2)
	b := complex(3, 4)
	c := a + b
	return fmt.Sprintf("%v:%v", real(c), imag(c))
}

// ComplexSubtraction tests complex number subtraction
func ComplexSubtraction() string {
	a := complex(5, 7)
	b := complex(2, 3)
	c := a - b
	return fmt.Sprintf("%v:%v", real(c), imag(c))
}

// ComplexMultiplication tests complex number multiplication
func ComplexMultiplication() string {
	// (1+2i) * (3+4i) = 3 + 4i + 6i + 8i^2 = 3 + 10i - 8 = -5 + 10i
	a := complex(1, 2)
	b := complex(3, 4)
	c := a * b
	return fmt.Sprintf("%v:%v", real(c), imag(c))
}

// ComplexDivision tests complex number division
func ComplexDivision() string {
	// (5+10i) / (1+2i) = 5
	a := complex(5, 10)
	b := complex(1, 2)
	c := a / b
	return fmt.Sprintf("%v:%v", real(c), imag(c))
}

// ComplexConjugate tests complex conjugate via arithmetic
func ComplexConjugate() string {
	z := complex(3, 4)
	conj := complex(real(z), -imag(z))
	product := z * conj
	return fmt.Sprintf("%v", real(product))
}

// ComplexMagnitude tests complex magnitude calculation
func ComplexMagnitude() string {
	z := complex(3, 4)
	magnitude := real(z)*real(z) + imag(z)*imag(z)
	return fmt.Sprintf("%v", magnitude)
}

// Complex64Operations tests complex64 type
func Complex64Operations() string {
	var a complex64 = complex(1.5, 2.5)
	var b complex64 = complex(0.5, 0.5)
	c := a + b
	return fmt.Sprintf("%v:%v", real(c), imag(c))
}

// ComplexZero tests complex zero
func ComplexZero() string {
	z := complex(0, 0)
	return fmt.Sprintf("%v:%v", real(z) == 0, imag(z) == 0)
}

// ComplexPureReal tests purely real complex
func ComplexPureReal() string {
	z := complex(5, 0)
	return fmt.Sprintf("%v:%v", real(z), imag(z))
}

// ComplexPureImag tests purely imaginary complex
func ComplexPureImag() string {
	z := complex(0, 3)
	return fmt.Sprintf("%v:%v", real(z), imag(z))
}

// ComplexComparison tests that complex numbers cannot be compared
func ComplexComparison() string {
	// Note: This function returns a string representation of the comparison workaround
	a := complex(1, 2)
	b := complex(1, 2)
	equal := real(a) == real(b) && imag(a) == imag(b)
	return fmt.Sprintf("%v", equal)
}
