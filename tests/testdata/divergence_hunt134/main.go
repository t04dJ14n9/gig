package divergence_hunt134

import "fmt"

// ============================================================================
// Round 134: Named return values and their interaction with defer
// ============================================================================

func NamedReturnBasic() (x int) {
	x = 42
	return
}

func NamedReturnOverride() (x int) {
	return 99
}

func NamedReturnDefer() (x int) {
	defer func() { x++ }()
	return 10
}

func NamedReturnDeferDouble() (x int) {
	defer func() { x *= 2 }()
	return 5
}

func NamedReturnMulti() (a int, b string) {
	a = 10
	b = "hello"
	return
}

func NamedReturnDeferMulti() (a int, b string) {
	defer func() {
		a += 100
		b += "-modified"
	}()
	a = 1
	b = "test"
	return
}

func NamedReturnShadow() (x int) {
	x = 10
	{
		x := 999 // shadow
		_ = x
	}
	return // returns the named x, not the shadow
}

func NamedReturnZeroValue() (s string, n int, f float64) {
	return
}

func NamedReturnPanicRecover() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("recovered-%v", r)
		}
	}()
	panic("oops")
}

func NamedReturnDeferModify() (result string) {
	result = "initial"
	defer func() { result = "deferred" }()
	result = "before-return"
	return result
}
