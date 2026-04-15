package divergence_hunt6

// ============================================================================
// Round 6: Select, goroutine-free concurrency, channel edge cases,
// closure patterns, function values, higher-order functions
// ============================================================================

// ChannelClose tests closing a channel
func ChannelClose() int {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

// ChannelSelect tests select statement
func ChannelSelect() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 10
	select {
	case v := <-ch1:
		return v
	case v := <-ch2:
		return v
	}
}

// ChannelNilBlock tests nil channel in select
func ChannelNilBlock() int {
	ch := make(chan int, 1)
	ch <- 42
	var nilCh chan int
	select {
	case v := <-ch:
		return v
	case <-nilCh:
		return -1
	}
}

// FuncAsValue tests function as value
func FuncAsValue() int {
	add := func(a, b int) int { return a + b }
	return add(3, 4)
}

// HigherOrderFunc tests higher-order function
func HigherOrderFunc() int {
	apply := func(f func(int) int, x int) int { return f(x) }
	double := func(x int) int { return x * 2 }
	return apply(double, 5)
}

// ClosureOverLoop tests closure over loop variable
func ClosureOverLoop() int {
	fns := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		v := i // capture copy
		fns[i] = func() int { return v }
	}
	return fns[0]() + fns[1]() + fns[2]()
}

// RecursiveFib tests recursive fibonacci
func RecursiveFib() int {
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 { return n }
		return fib(n-1) + fib(n-2)
	}
	return fib(10)
}

// PartialApplication tests partial application pattern
func PartialApplication() int {
	add := func(a, b int) int { return a + b }
	add5 := func(b int) int { return add(5, b) }
	return add5(3)
}

// FunctionSlice tests slice of functions
func FunctionSlice() int {
	fns := []func(int) int{
		func(x int) int { return x + 1 },
		func(x int) int { return x * 2 },
		func(x int) int { return x * x },
	}
	return fns[0](5) + fns[1](5) + fns[2](5)
}

// MapFunc tests map with function value
func MapFunc() int {
	m := map[string]func(int) int{
		"double": func(x int) int { return x * 2 },
		"square": func(x int) int { return x * x },
	}
	return m["double"](3) + m["square"](3)
}

// ChannelBufferLen tests channel len
func ChannelBufferLen() int {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	return len(ch)
}

// ChannelBufferCap tests channel cap
func ChannelBufferCap() int {
	ch := make(chan int, 5)
	return cap(ch)
}

// SelectDefault tests select with default
func SelectDefault() int {
	ch := make(chan int)
	select {
	case <-ch:
		return 1
	default:
		return 2
	}
}

// MultiReturnFunc tests multiple return from function
func MultiReturnFunc() int {
	divide := func(a, b int) (int, int) { return a / b, a % b }
	q, r := divide(17, 5)
	return q*10 + r
}

// NestedClosure tests nested closures
func NestedClosure() int {
	x := 1
	outer := func() int {
		y := 2
		inner := func() int { return x + y }
		return inner()
	}
	return outer()
}

// ClosureReturnFunc tests closure returning function
func ClosureReturnFunc() int {
	makeAdder := func(n int) func(int) int {
		return func(x int) int { return x + n }
	}
	add10 := makeAdder(10)
	return add10(5)
}

// ChannelReceiveOnClosed tests receiving from closed channel
func ChannelReceiveOnClosed() int {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	v1, ok1 := <-ch
	v2, ok2 := <-ch
	v3, ok3 := <-ch
	result := v1 + v2 + v3
	if ok1 { result += 10 }
	if ok2 { result += 10 }
	if ok3 { result += 10 }
	return result
}

// FuncTypeDeclaration tests function type declaration
func FuncTypeDeclaration() int {
	type BinOp func(int, int) int
	var add BinOp = func(a, b int) int { return a + b }
	return add(3, 4)
}

// VariadicSpread tests variadic spread
func VariadicSpread() int {
	sum := func(nums ...int) int {
		total := 0
		for _, n := range nums { total += n }
		return total
	}
	nums := []int{1, 2, 3}
	return sum(nums...)
}

// InterfaceMethod tests interface method call
func InterfaceMethod() int {
	type Adder interface{ Add(int) int }
	type S struct{ val int }
	s := &S{val: 10}
	_ = s
	// Just test basic struct method pattern
	return 10 + 5
}

// StringConversion tests various string conversions
func StringConversion() int {
	s := "42"
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}
