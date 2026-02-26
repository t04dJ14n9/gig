package edgecases

// MaxInt64 tests max int64 value
func MaxInt64() int {
	return 9223372036854775807
}

// MinInt64 tests min int64 value
func MinInt64() int {
	return -9223372036854775807 - 1
}

// DivisionByMinusOne tests division by -1
func DivisionByMinusOne() int {
	x := 42
	return x / (-1)
}

// ModuloNegative tests modulo with negative
func ModuloNegative() int {
	return (-7) % 3
}

// EmptyString tests empty string
func EmptyString() string { return "" }

// LargeSlice tests large slice operations
func LargeSlice() int {
	s := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		s[i] = i
	}
	sum := 0
	for _, v := range s {
		sum = sum + v
	}
	return sum
}

// NestedMapLookup tests nested map lookup
func NestedMapLookup() int {
	m := make(map[string]int)
	keys := make([]string, 0)
	keys = append(keys, "a")
	keys = append(keys, "b")
	keys = append(keys, "c")
	for i, k := range keys {
		m[k] = (i + 1) * 10
	}
	sum := 0
	for _, k := range keys {
		sum = sum + m[k]
	}
	return sum
}

// ZeroDivisionGuard tests safe division
func ZeroDivisionGuard() int {
	return safeDivide(10, 2) + safeDivide(10, 0)
}

func safeDivide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

// BooleanComplexExpr tests complex boolean expression
func BooleanComplexExpr() int {
	a := 10
	b := 20
	c := 30
	result := 0
	if a < b && b < c {
		result = result + 1
	}
	if a > b || c > b {
		result = result + 10
	}
	if !(a > c) {
		result = result + 100
	}
	if a < b && (c > 20 || b < 10) {
		result = result + 1000
	}
	return result
}

// SingleElementSlice tests single element slice
func SingleElementSlice() int {
	s := make([]int, 0)
	s = append(s, 42)
	return s[0] + len(s)
}

// EmptyMap tests empty map
func EmptyMap() int {
	m := make(map[string]int)
	return len(m)
}

// TightLoop tests tight loop with complex operations
func TightLoop() int {
	result := 0
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			if (i+j)%2 == 0 {
				result = result + 1
			}
		}
	}
	return result
}
