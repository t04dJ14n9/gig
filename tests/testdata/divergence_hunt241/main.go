package divergence_hunt241

import (
	"fmt"
)

// ============================================================================
// Round 241: Defer with named returns
// ============================================================================

// NamedReturnDeferModify tests defer modifying named return value
func NamedReturnDeferModify() string {
	result := func() (r int) {
		defer func() { r = 42 }()
		return 1
	}()
	return fmt.Sprintf("%d", result)
}

// NamedReturnDeferIncrement tests defer incrementing named return
func NamedReturnDeferIncrement() string {
	result := func() (r int) {
		defer func() { r++ }()
		return 10
	}()
	return fmt.Sprintf("%d", result)
}

// NamedReturnDeferAdd tests defer adding to named return
func NamedReturnDeferAdd() string {
	result := func() (sum int) {
		defer func() { sum += 100 }()
		return 5
	}()
	return fmt.Sprintf("%d", result)
}

// NamedReturnDeferMultiply tests defer multiplying named return
func NamedReturnDeferMultiply() string {
	result := func() (r int) {
		defer func() { r *= 3 }()
		return 7
	}()
	return fmt.Sprintf("%d", result)
}

// NamedReturnDeferChain tests multiple defers with named return
func NamedReturnDeferChain() string {
	result := func() (r int) {
		defer func() { r = r*10 + 1 }()
		defer func() { r = r*10 + 2 }()
		defer func() { r = r*10 + 3 }()
		return 0
	}()
	return fmt.Sprintf("%d", result)
}

// NamedReturnDeferString tests defer modifying named string return
func NamedReturnDeferString() string {
	result := func() (s string) {
		defer func() { s = s + " world" }()
		return "hello"
	}()
	return result
}

// NamedReturnDeferSlice tests defer appending to named slice return
func NamedReturnDeferSlice() string {
	result := func() (s []int) {
		defer func() { s = append(s, 4, 5) }()
		return []int{1, 2, 3}
	}()
	return fmt.Sprintf("%v", result)
}

// NamedReturnDeferMap tests defer modifying named map return
func NamedReturnDeferMap() string {
	result := func() (m map[string]int) {
		defer func() { m["key2"] = 200 }()
		return map[string]int{"key1": 100}
	}()
	return fmt.Sprintf("%d:%d", result["key1"], result["key2"])
}

// NamedReturnDeferStruct tests defer modifying struct field in named return
func NamedReturnDeferStruct() string {
	type Point struct{ X, Y int }
	result := func() (p Point) {
		defer func() { p.X = 10; p.Y = 20 }()
		return Point{X: 1, Y: 2}
	}()
	return fmt.Sprintf("%d,%d", result.X, result.Y)
}

// NamedReturnDeferPointer tests defer modifying pointer in named return
func NamedReturnDeferPointer() string {
	result := func() (p *int) {
		defer func() {
			if p != nil {
				*p = 99
			}
		}()
		v := 42
		return &v
	}()
	return fmt.Sprintf("%d", *result)
}

// NamedReturnDeferWithPanicRecover tests defer with named return after panic
func NamedReturnDeferWithPanicRecover() string {
	result := func() (r string) {
		defer func() {
			if rec := recover(); rec != nil {
				r = fmt.Sprintf("recovered:%v", rec)
			}
		}()
		panic("boom")
	}()
	return result
}
