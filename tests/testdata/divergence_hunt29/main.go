package divergence_hunt29

import "fmt"

// ============================================================================
// Round 29: Error handling, custom error types, validation patterns
// ============================================================================

func SimpleError() string {
	err := fmt.Errorf("something went wrong")
	return err.Error()
}

func ErrorWithFormat() string {
	err := fmt.Errorf("value %d is out of range [0, %d]", 100, 10)
	return err.Error()
}

func ValidatePositive() (result int) {
	validate := func(n int) error {
		if n < 0 { return fmt.Errorf("negative: %d", n) }
		return nil
	}
	if err := validate(5); err == nil { result += 1 }
	if err := validate(-1); err != nil { result += 10 }
	return result
}

func ValidateRange() (result int) {
	validate := func(n, min, max int) error {
		if n < min || n > max { return fmt.Errorf("%d not in [%d, %d]", n, min, max) }
		return nil
	}
	if err := validate(5, 0, 10); err == nil { result += 1 }
	if err := validate(15, 0, 10); err != nil { result += 10 }
	return result
}

func ErrorPropagation() (result int) {
	step1 := func() error { return nil }
	step2 := func() error { return fmt.Errorf("step2 failed") }
	step3 := func() error { return nil }
	if err := step1(); err != nil { return -1 }
	if err := step2(); err != nil { return -2 }
	if err := step3(); err != nil { return -3 }
	return 0
}

func ErrorInDefer() (result int) {
	defer func() {
		if r := recover(); r != nil { result = -1 }
	}()
	panic(fmt.Errorf("deferred error"))
}

func MultiErrorCollect() int {
	errors := []error{}
	validate := func(n int) error {
		if n < 0 { return fmt.Errorf("negative: %d", n) }
		return nil
	}
	for _, n := range []int{1, -2, 3, -4, 5} {
		if err := validate(n); err != nil {
			errors = append(errors, err)
		}
	}
	return len(errors)
}

func ErrorTypeAssertion() string {
	err := fmt.Errorf("test error")
	return err.Error()
}

func PanicWithFmtError() (result string) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				result = err.Error()
			} else {
				result = fmt.Sprintf("%v", r)
			}
		}
	}()
	panic(fmt.Errorf("custom error"))
}

func NilErrorCheck() int {
	var err error
	if err == nil { return 1 }
	return 0
}

func ErrorStringMethods() int {
	type MyError struct{ Code int }
	myErr := &MyError{Code: 404}
	_ = myErr
	return 404
}

func ValidateStruct() (result int) {
	type Config struct {
		Name  string
		Value int
	}
	validate := func(c Config) error {
		if c.Name == "" { return fmt.Errorf("name is empty") }
		if c.Value < 0 { return fmt.Errorf("value is negative") }
		return nil
	}
	c1 := Config{Name: "test", Value: 10}
	c2 := Config{Name: "", Value: 10}
	c3 := Config{Name: "test", Value: -1}
	if err := validate(c1); err == nil { result++ }
	if err := validate(c2); err != nil { result += 10 }
	if err := validate(c3); err != nil { result += 100 }
	return result
}

func ErrorInClosure() (result int) {
	fn := func() error { return fmt.Errorf("closure error") }
	if err := fn(); err != nil { result = 1 }
	return result
}

func FmtErrorfWrap() string {
	inner := fmt.Errorf("inner error")
	outer := fmt.Errorf("outer: %w", inner)
	return outer.Error()
}
