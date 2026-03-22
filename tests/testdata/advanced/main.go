package advanced

// Type Conversions

func TypeConvertIntIdentity() int {
	x := 42
	f := float64(x)
	return int(f)
}

// Nested Function Calls

func DeepCallChain() int {
	return a(0)
}

func a(x int) int { return b(x + 1) }

func b(x int) int { return c(x + 1) }

func c(x int) int { return d(x + 1) }

func d(x int) int { return e(x + 1) }

func e(x int) int { return x + 1 }

// Complex Control Flow

func EarlyReturn() int {
	s := make([]int, 0)
	s = append(s, 10)
	s = append(s, 20)
	s = append(s, 30)
	return findFirst(s, 20)
}

func findFirst(s []int, target int) int {
	for i := 0; i < len(s); i++ {
		if s[i] == target {
			return i
		}
	}
	return -1
}

func NestedIfInLoop() int {
	count := 0
	for i := 1; i <= 100; i++ {
		if i%3 == 0 && i%5 == 0 {
			count = count + 1
		}
	}
	return count
}

// Algorithm Tests

func BubbleSort() int {
	s := make([]int, 0)
	s = append(s, 5)
	s = append(s, 3)
	s = append(s, 8)
	s = append(s, 1)
	s = append(s, 9)
	s = append(s, 2)
	n := len(s)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if s[j] > s[j+1] {
				tmp := s[j]
				s[j] = s[j+1]
				s[j+1] = tmp
			}
		}
	}
	return s[0]*100000 + s[1]*10000 + s[2]*1000 + s[3]*100 + s[4]*10 + s[5]
}

func BinarySearch() int {
	s := make([]int, 0)
	for i := 0; i < 10; i++ {
		s = append(s, i*10)
	}
	return bsearch(s, 50)
}

func bsearch(s []int, target int) int {
	lo := 0
	hi := len(s) - 1
	for lo <= hi {
		mid := (lo + hi) / 2
		if s[mid] == target {
			return mid
		} else if s[mid] < target {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return -1
}

func GCD() int {
	return gcd(48, 18)
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func SieveOfEratosthenes() int {
	n := 100
	sieve := make([]int, n+1)
	for i := 2; i <= n; i++ {
		sieve[i] = 1
	}
	for i := 2; i*i <= n; i++ {
		if sieve[i] == 1 {
			for j := i * i; j <= n; j = j + i {
				sieve[j] = 0
			}
		}
	}
	count := 0
	for i := 2; i <= n; i++ {
		if sieve[i] == 1 {
			count = count + 1
		}
	}
	return count
}

func MatrixMultiply() int {
	a := make([]int, 4)
	a[0] = 1
	a[1] = 2
	a[2] = 3
	a[3] = 4
	b := make([]int, 4)
	b[0] = 5
	b[1] = 6
	b[2] = 7
	b[3] = 8
	c := make([]int, 4)
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			sum := 0
			for k := 0; k < 2; k++ {
				sum = sum + a[i*2+k]*b[k*2+j]
			}
			c[i*2+j] = sum
		}
	}
	return c[0]*1000 + c[1]*100 + c[2]*10 + c[3]
}

// Edge Cases

func EmptyFunctionReturn() int {
	return identity(42)
}

func identity(x int) int { return x }

func SingleReturnValue() int { return 1 }

func ZeroIteration() int {
	sum := 0
	for i := 0; i < 0; i++ {
		sum = sum + 1
	}
	return sum
}

func LargeLoop() int {
	sum := 0
	for i := 0; i < 10000; i++ {
		sum = sum + 1
	}
	return sum
}

func DeepRecursion() int {
	return countdown(50)
}

func countdown(n int) int {
	if n <= 0 {
		return 0
	}
	return 1 + countdown(n-1)
}

// Multi-Feature Combination

func MapWithClosure() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	transform := func(x int) int { return x * 10 }
	sum := 0
	for _, v := range m {
		sum = sum + transform(v)
	}
	return sum
}

func SliceWithMultiReturn() int {
	s := make([]int, 0)
	s = append(s, 3)
	s = append(s, 1)
	s = append(s, 4)
	s = append(s, 1)
	s = append(s, 5)
	mn, mx := minmax(s)
	return mn*10 + mx
}

func minmax(s []int) (int, int) {
	mn := s[0]
	mx := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] < mn {
			mn = s[i]
		}
		if s[i] > mx {
			mx = s[i]
		}
	}
	return mn, mx
}

func RecursiveDataBuild() int {
	return buildAndSum(10)
}

func buildAndSum(n int) int {
	s := make([]int, 0)
	for i := 0; i < n; i++ {
		s = append(s, i*i)
	}
	sum := 0
	for _, v := range s {
		sum = sum + v
	}
	return sum
}

func FunctionChain() int {
	return sub(mul(add(1, 2), add(3, 4)), 5)
}

func add(a, b int) int { return a + b }

func mul(a, b int) int { return a * b }

func sub(a, b int) int { return a - b }

func ComplexExpressions() int {
	a := 10
	b := 20
	c := 30
	return (a+b)*(c-a) + b/a - c%b
}

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// FindFirst returns the index of target in s, or -1
func FindFirst(s []int, target int) int { return findFirst(s, target) }

// Bsearch returns the index of target in sorted s, or -1
func Bsearch(s []int, target int) int { return bsearch(s, target) }

// Gcd computes the greatest common divisor of a and b
func Gcd(a, b int) int { return gcd(a, b) }

// Identity returns x unchanged
func Identity(x int) int { return identity(x) }

// Minmax returns (min, max) of slice s
func Minmax(s []int) (int, int) { return minmax(s) }

// Countdown returns n (counting down recursively)
func Countdown(n int) int { return countdown(n) }

// Add returns a + b
func Add(a, b int) int { return add(a, b) }

// Mul returns a * b
func Mul(a, b int) int { return a * b }

// Sub returns a - b
func Sub(a, b int) int { return a - b }
