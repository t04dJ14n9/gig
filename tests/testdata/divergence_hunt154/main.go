package divergence_hunt154

import "fmt"

// ============================================================================
// Round 154: Anonymous structs and function types
// ============================================================================

// AnonymousStructLiteral tests anonymous struct literals
func AnonymousStructLiteral() string {
	p := struct {
		Name string
		Age  int
	}{"Alice", 30}
	return fmt.Sprintf("%s-%d", p.Name, p.Age)
}

// AnonymousStructSlice tests slice of anonymous structs
func AnonymousStructSlice() string {
	people := []struct {
		Name string
		Age  int
	}{
		{"Bob", 25},
		{"Carol", 28},
	}
	return fmt.Sprintf("count=%d", len(people))
}

// AnonymousStructMap tests map with anonymous struct values
func AnonymousStructMap() string {
	scores := map[string]struct {
		Math  int
		Sci   int
	}{
		"Alice": {90, 85},
		"Bob":   {75, 80},
	}
	return fmt.Sprintf("alice-math=%d", scores["Alice"].Math)
}

// FunctionTypeComparison tests comparing functions (they're only comparable to nil)
func FunctionTypeComparison() string {
	var f func()
	isNil := f == nil
	g := func() {}
	isNotNil := g != nil
	return fmt.Sprintf("nil=%t-notnil=%t", isNil, isNotNil)
}

// FunctionTypeAssignment tests assigning functions to variables
func FunctionTypeAssignment() string {
	var op func(int, int) int
	op = func(a, b int) int { return a + b }
	result := op(3, 4)
	return fmt.Sprintf("result=%d", result)
}

// FunctionTypeReturn tests returning function types
func FunctionTypeReturn() string {
	makeMultiplier := func(n int) func(int) int {
		return func(x int) int { return x * n }
	}
	double := makeMultiplier(2)
	triple := makeMultiplier(3)
	return fmt.Sprintf("d=%d-t=%d", double(5), triple(5))
}

// FunctionTypeParam tests function as parameter
func FunctionTypeParam() string {
	apply := func(x int, f func(int) int) int {
		return f(x)
	}
	result := apply(5, func(x int) int { return x * x })
	return fmt.Sprintf("result=%d", result)
}

// FunctionTypeSlice tests slice of functions
func FunctionTypeSlice() string {
	ops := []func(int) int{
		func(x int) int { return x + 1 },
		func(x int) int { return x * 2 },
		func(x int) int { return x * x },
	}
	result := 3
	for _, op := range ops {
		result = op(result)
	}
	return fmt.Sprintf("result=%d", result)
}

// FunctionTypeMap tests map of functions
func FunctionTypeMap() string {
	ops := map[string]func(int, int) int{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
	}
	return fmt.Sprintf("add=%d-mul=%d", ops["add"](3, 4), ops["mul"](3, 4))
}

// AnonymousStructNested tests nested anonymous structs
func AnonymousStructNested() string {
	data := struct {
		Name    string
		Address struct {
			City    string
			Country string
		}
	}{
		Name: "Test",
		Address: struct {
			City    string
			Country string
		}{
			City:    "NYC",
			Country: "USA",
		},
	}
	return fmt.Sprintf("%s-%s", data.Name, data.Address.City)
}
