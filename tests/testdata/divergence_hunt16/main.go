package divergence_hunt16

// ============================================================================
// Round 16: Control flow edge cases, switch patterns, loop patterns
// ============================================================================

// SwitchNoCase tests switch with no matching case
func SwitchNoCase() int {
	x := 99
	switch x {
	case 1: return 10
	case 2: return 20
	}
	return 0
}

// SwitchMultipleCases tests switch with multiple cases
func SwitchMultipleCases() int {
	x := 2
	switch x {
	case 1, 2, 3: return 10
	default: return 20
	}
}

// SwitchWithInit tests switch with init statement
func SwitchWithInit() int {
	switch x := 5; x {
	case 5: return 1
	default: return 0
	}
}

// NestedSwitch tests nested switch
func NestedSwitch() int {
	x, y := 1, 2
	switch x {
	case 1:
		switch y {
		case 2: return 3
		default: return 1
		}
	default: return 0
	}
}

// ForRangeWithIndex tests for range with index
func ForRangeWithIndex() int {
	s := []int{10, 20, 30}
	sum := 0
	for i := range s {
		sum += i
	}
	return sum
}

// ForRangeWithValue tests for range with value
func ForRangeWithValue() int {
	s := []int{10, 20, 30}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

// ForRangeMap tests for range over map
func ForRangeMap() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	count := 0
	for range m { count++ }
	return count
}

// IfElseChain tests if-else chain
func IfElseChain() int {
	x := 3
	if x == 1 { return 10 }
	if x == 2 { return 20 }
	if x == 3 { return 30 }
	return 0
}

// NestedIfElse tests nested if-else
func NestedIfElse() int {
	x, y := 5, 10
	if x > 3 {
		if y > 8 { return 1 }
		return 2
	}
	return 3
}

// InfiniteLoopBreak tests infinite loop with break
func InfiniteLoopBreak() int {
	i := 0
	for {
		i++
		if i >= 10 { break }
	}
	return i
}

// ForLoopContinue tests for loop with continue
func ForLoopContinue() int {
	sum := 0
	for i := 0; i < 10; i++ {
		if i%3 == 0 { continue }
		sum += i
	}
	return sum
}

// LoopWithMultipleBreaks tests loop with multiple break conditions
func LoopWithMultipleBreaks() int {
	s := []int{1, 3, 5, 7, 2, 4, 6}
	for i, v := range s {
		if v%2 == 0 { return i }
	}
	return -1
}

// SwitchExpression tests switch with expression
func SwitchExpression() string {
	x := 42
	switch {
	case x < 0: return "negative"
	case x == 0: return "zero"
	case x > 0: return "positive"
	}
	return "impossible"
}

// ForRangeString tests for range over string
func ForRangeString() int {
	s := "Hello"
	count := 0
	for range s { count++ }
	return count
}

// ForRangeEmptySlice tests for range over empty slice
func ForRangeEmptySlice() int {
	s := []int{}
	count := 0
	for range s { count++ }
	return count
}

// DoubleLoop tests double loop pattern
func DoubleLoop() int {
	result := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result += i*3 + j
		}
	}
	return result
}

// LoopAccumulator tests accumulator in loop
func LoopAccumulator() int {
	acc := 1
	for i := 1; i <= 5; i++ {
		acc *= i
	}
	return acc
}

// SwitchFallthroughSimulated tests simulated fallthrough with if
func SwitchFallthroughSimulated() int {
	x := 2
	result := 0
	switch x {
	case 1: result += 1; fallthrough
	case 2: result += 10; fallthrough
	case 3: result += 100
	}
	return result
}

// EarlyReturn tests early return pattern
func EarlyReturn() int {
	x := 5
	if x > 3 { return x * 2 }
	return x
}

// LoopWithEarlyReturn tests loop with early return
func LoopWithEarlyReturn() int {
	for i := 0; i < 10; i++ {
		if i*i > 20 { return i }
	}
	return -1
}
