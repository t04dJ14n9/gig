package divergence_hunt130

import "fmt"

// ============================================================================
// Round 130: Defer and recover patterns — more edge cases
// ============================================================================

func DeferStackOrder() string {
	var result string
	defer func() { result += "A" }()
	defer func() { result += "B" }()
	defer func() { result += "C" }()
	_ = result // use before defers
	return "done"
}

func DeferModifyReturn() (result string) {
	defer func() { result = "modified" }()
	return "original"
}

func DeferNamedReturn() (x int) {
	defer func() { x++ }()
	return 10
}

func DeferCaptureValue() string {
	x := 10
	defer func() {
		_ = fmt.Sprintf("x=%d", x)
	}()
	x = 20
	return fmt.Sprintf("x=%d", x)
}

func DeferCapturePointer() string {
	x := 10
	ptr := &x
	defer func() {
		*ptr = 99
	}()
	x = 20
	return fmt.Sprintf("ptr=%d", *ptr)
}

func RecoverBasic() string {
	defer func() {
		if r := recover(); r != nil {
			_ = fmt.Sprintf("recovered: %v", r)
		}
	}()
	panic("test-panic")
}

func RecoverInDefer() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("caught-%v", r)
		}
	}()
	panic("hello")
}

func RecoverNoPanic() string {
	defer func() {
		if r := recover(); r != nil {
			_ = fmt.Sprintf("unexpected: %v", r)
		}
	}()
	return "normal"
}

func DeferMultipleRecovers() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("first-%v", r)
		}
	}()
	defer func() {
		// This runs first (LIFO), but there's no panic yet
		if r := recover(); r != nil {
			result += fmt.Sprintf("second-%v", r)
		}
	}()
	panic("boom")
}

func DeferPanicInDefer() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("outer-%v", r)
		}
	}()
	defer func() {
		panic("inner-panic")
	}()
	return "normal"
}
