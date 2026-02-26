package scope

// IfInitShortVar tests if init with short var
func IfInitShortVar() int {
	if v := abs(-42); v > 0 {
		return v
	}
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// IfInitMultiCondition tests if init with multiple conditions
func IfInitMultiCondition() int {
	result := 0
	for i := 0; i < 10; i++ {
		if rem := i % 3; rem == 0 {
			result = result + i
		}
	}
	return result
}

// NestedScopes tests nested scopes
func NestedScopes() int {
	x := 1
	y := 0
	if x > 0 {
		x := 10
		y = x
	}
	return x + y
}

// ForScopeIsolation tests for scope isolation
func ForScopeIsolation() int {
	sum := 0
	for i := 0; i < 3; i++ {
		x := i * 10
		sum = sum + x
	}
	return sum
}

// MultipleBlockScopes tests multiple block scopes
func MultipleBlockScopes() int {
	result := 0
	x := 1
	if x > 0 {
		a := 10
		result = result + a
	}
	if x > 0 {
		b := 20
		result = result + b
	}
	return result
}

// ClosureCapturesOuterScope tests closure captures outer scope
func ClosureCapturesOuterScope() int {
	x := 100
	y := 200
	add := func() int { return x + y }
	x = 150
	return add()
}
