package divergence_hunt47

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 47: Error handling patterns - custom errors, wrapping, validation
// ============================================================================

func SimpleErrorCheck() int {
	mightFail := func(ok bool) (int, error) {
		if !ok { return 0, fmt.Errorf("failed") }
		return 42, nil
	}
	v, err := mightFail(true)
	if err != nil { return -1 }
	return v
}

func ErrorPropagation() int {
	step1 := func() error { return nil }
	step2 := func() (int, error) { return 42, nil }
	if err := step1(); err != nil { return -1 }
	v, err := step2()
	if err != nil { return -2 }
	return v
}

func ErrorInClosure() int {
	process := func(data []int) (int, error) {
		if len(data) == 0 { return 0, fmt.Errorf("empty") }
		sum := 0
		for _, v := range data { sum += v }
		return sum, nil
	}
	v, err := process([]int{1, 2, 3})
	if err != nil { return -1 }
	return v
}

func ErrorChain() string {
	inner := func() error { return fmt.Errorf("inner error") }
	outer := func() error {
		if err := inner(); err != nil {
			return fmt.Errorf("outer: %w", err)
		}
		return nil
	}
	return outer().Error()
}

func ValidationError() int {
	validate := func(age int) error {
		if age < 0 { return fmt.Errorf("negative age") }
		if age > 150 { return fmt.Errorf("age too high") }
		return nil
	}
	if err := validate(25); err != nil { return -1 }
	if err := validate(-1); err != nil { return -2 }
	return 0
}

func MultiErrorCollect() int {
	errors := []error{}
	check := func(s string) error {
		if s == "" { return fmt.Errorf("empty string") }
		return nil
	}
	for _, s := range []string{"hello", "", "world", ""} {
		if err := check(s); err != nil {
			errors = append(errors, err)
		}
	}
	return len(errors)
}

func PanicInsteadOfError() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = -1
		}
	}()
	mustSucceed := func(ok bool) {
		if !ok { panic("assertion failed") }
	}
	mustSucceed(false)
	return 0
}

func ErrorTypeAssertion() int {
	customErr := fmt.Errorf("code %d", 404)
	var err error = customErr
	if v, ok := err.(error); ok {
		_ = v
		return 1
	}
	return 0
}

func JSONUnmarshalError() int {
	var x int
	err := json.Unmarshal([]byte("not json"), &x)
	if err != nil { return -1 }
	return x
}

func FmtErrorfWrap() string {
	inner := fmt.Errorf("inner")
	outer := fmt.Errorf("outer: %w", inner)
	return outer.Error()
}

func ErrorStringMethod() string {
	type MyError struct{ Code int }
	err := &MyError{Code: 404}
	// Can't add methods in tests, use fmt.Errorf
	return fmt.Errorf("error code %d", err.Code).Error()
}

func SortWithValidation() int {
	data := []int{5, 3, 1, 4, 2}
	if len(data) == 0 { return -1 }
	sort.Ints(data)
	return data[0] + data[4]
}

func StringsErrorCheck() int {
	result := strings.TrimPrefix("hello world", "hello ")
	if result == "" { return -1 }
	return len(result)
}
