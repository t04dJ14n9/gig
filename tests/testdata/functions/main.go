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
func double(x int) int                  { return x * 2 }
func triple(x int) int                  { return x * 3 }

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
