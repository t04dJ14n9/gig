package functions

// Call tests basic function call
func Call() int { return add(5, 7) }

func add(a int, b int) int { return a + b }

// MultipleReturn tests multiple return values
func MultipleReturn() int {
	x, y := swap(3, 7)
	return x + y
}

func swap(a, b int) (int, int) { return b, a }

// MultipleReturnDivmod tests divmod pattern
func MultipleReturnDivmod() int {
	q, r := divmod(17, 5)
	return q*10 + r
}

func divmod(a, b int) (int, int) { return a / b, a % b }

// RecursionFactorial tests recursive factorial
func RecursionFactorial() int { return factorial(5) }

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// MutualRecursion tests mutual recursion
func MutualRecursion() int {
	if isEven(10) {
		return 1
	}
	return 0
}

func isEven(n int) bool {
	if n == 0 {
		return true
	}
	return isOdd(n - 1)
}

func isOdd(n int) bool {
	if n == 0 {
		return false
	}
	return isEven(n - 1)
}

// FibonacciIterative tests iterative fibonacci
func FibonacciIterative() int { return fibIter(20) }

func fibIter(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		c := a + b
		a = b
		b = c
	}
	return b
}

// FibonacciRecursive tests recursive fibonacci
func FibonacciRecursive() int { return fibRec(15) }

func fibRec(n int) int {
	if n <= 1 {
		return n
	}
	return fibRec(n-1) + fibRec(n-2)
}

// VariadicFunction tests variadic function
func VariadicFunction() int { return sum(1, 2, 3, 4, 5) }

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total = total + n
	}
	return total
}

// FunctionAsValue tests function as value
func FunctionAsValue() int {
	return apply(double, 5) + apply(triple, 5)
}

func apply(f func(int) int, x int) int { return f(x) }

func double(x int) int { return x * 2 }

func triple(x int) int { return x * 3 }

// HigherOrderMap tests higher order map function
func HigherOrderMap() int {
	s := make([]int, 3)
	s[0] = 1
	s[1] = 2
	s[2] = 3
	doubled := mapSlice(s, func(x int) int { return x * 2 })
	return doubled[0] + doubled[1] + doubled[2]
}

func mapSlice(s []int, f func(int) int) []int {
	result := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = f(s[i])
	}
	return result
}

// HigherOrderFilter tests higher order filter function
func HigherOrderFilter() int {
	s := make([]int, 0)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	s = append(s, 4)
	s = append(s, 5)
	return count(s, func(x int) bool { return x > 3 })
}

func count(s []int, pred func(int) bool) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if pred(s[i]) {
			n = n + 1
		}
	}
	return n
}

// HigherOrderReduce tests higher order reduce function
func HigherOrderReduce() int {
	s := make([]int, 0)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	s = append(s, 4)
	return reduce(s, 0, func(a, b int) int { return a + b })
}

func reduce(s []int, init int, f func(int, int) int) int {
	acc := init
	for i := 0; i < len(s); i++ {
		acc = f(acc, s[i])
	}
	return acc
}

// ThreeReturnValues tests function returning 3 values - the function itself returns 3 values!
func ThreeReturnValues() (int, int, int) {
	return 1, 2, 3
}

// FourReturnValues tests function returning 4 values - the function itself returns 4 values!
func FourReturnValues() (int, int, int, int) {
	return 1, 2, 3, 4
}

// FiveReturnValues tests function returning 5 values - the function itself returns 5 values!
func FiveReturnValues() (int, int, int, int, int) {
	return 1, 2, 3, 4, 5
}

// MixedTypeReturn tests function returning mixed types - the function itself returns 3 mixed-type values!
func MixedTypeReturn() (int, string, bool) {
	return 42, "hello", true
}

// PassMultiReturnToFunc tests passing multi-return directly to another function
func PassMultiReturnToFunc() int {
	return addPair(pair())
}

func pair() (int, int)     { return 10, 20 }
func addPair(a, b int) int { return a + b }

// ChainMultiReturn tests one function returning another's multi-return
func ChainMultiReturn() int {
	a, b := swap(1, 2)
	return a + b
}

// NestedMultiReturn tests using multi-return in expressions
func NestedMultiReturn() int {
	a, _ := swap(5, 10)
	b, _ := swap(20, 30)
	return a + b
}

// MultiReturnAsSliceIndex tests multi-return used as slice index
func MultiReturnAsSliceIndex() int {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	idx, _ := swap(3, 5)
	return s[idx]
}

// MultiReturnToMap tests storing multi-return in map
func MultiReturnToMap() int {
	m := make(map[string]int)
	k, v := keyValue()
	m[k] = v
	return m[k]
}

func keyValue() (string, int) { return "key", 100 }

// MultiReturnAsCondition tests multi-return in condition
func MultiReturnAsCondition() int {
	if _, ok := maybeValue(); ok {
		return 1
	}
	return 0
}

func maybeValue() (int, bool) { return 42, true }

// MultiReturnComplexTypes tests returning complex types
func MultiReturnComplexTypes() int {
	sl, m := complexTypes()
	return len(sl) + len(m)
}

func complexTypes() ([]int, map[string]int) {
	return []int{1, 2, 3}, map[string]int{"a": 1}
}

// MultiReturnInClosure tests multi-return inside closure
func MultiReturnInClosure() int {
	f := func() (int, int) { return 7, 8 }
	a, b := f()
	return a + b
}

// AssignMultiReturnToExistingVars tests assigning to existing variables
func AssignMultiReturnToExistingVars() int {
	var a, b int
	a, b = swap(100, 200)
	return a + b
}

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// Add returns a + b
func Add(a, b int) int { return add(a, b) }

// Swap returns (b, a)
func Swap(a, b int) (int, int) { return swap(a, b) }

// Divmod returns (a/b, a%b)
func Divmod(a, b int) (int, int) { return divmod(a, b) }

// FactorialN returns n!
func FactorialN(n int) int { return factorial(n) }

// FibIterN returns the nth Fibonacci number iteratively
func FibIterN(n int) int { return fibIter(n) }

// FibRecN returns the nth Fibonacci number recursively
func FibRecN(n int) int { return fibRec(n) }

// IsEvenN returns true if n is even
func IsEvenN(n int) bool { return isEven(n) }

// IsOddN returns true if n is odd
func IsOddN(n int) bool { return isOdd(n) }

// SumVariadic returns the sum of all arguments
func SumVariadic(nums ...int) int { return sum(nums...) }
