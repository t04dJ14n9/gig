package thirdparty

import "sync"

// ============================================================================
// CLOSURE AND DEFER PATTERNS
// ============================================================================

// ClosureWithDeferAndPanicRecovery tests closure with defer.
func ClosureWithDeferAndPanicRecovery() int {
	result := 0
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = -1
			}
		}()
		panic("test")
	}()
	if result == 0 {
		result = 1
	}
	return result
}

// ClosureCapturingExternalVar tests closure capturing external var.
func ClosureCapturingExternalVar() int {
	counter := 0
	increment := func() int {
		counter++
		return counter
	}

	sum := 0
	for i := 0; i < 5; i++ {
		sum += increment()
	}
	return sum
}

// ============================================================================
// DEFER IN COMPLEX SCENARIOS
// ============================================================================

// MultipleDefersWithExternalCalls tests multiple defers.
func MultipleDefersWithExternalCalls() int {
	mu := sync.Mutex{}
	counter := 0

	func() {
		mu.Lock()
		defer mu.Unlock()
		counter++

		mu.Lock()
		defer mu.Unlock()
		counter++

		mu.Lock()
		defer mu.Unlock()
		counter++
	}()

	return counter
}

// ============================================================================
// TYPE ASSERTION WITH SWITCH
// ============================================================================

// TypeSwitchWithExternalTypes tests type switch on external types.
func TypeSwitchWithExternalTypes() int {
	var v interface{} = "hello"

	switch val := v.(type) {
	case string:
		return len(val)
	case int:
		return val
	default:
		return -1
	}
}

// TypeSwitchWithMultipleExternalTypes tests type switch with multiple types.
func TypeSwitchWithMultipleExternalTypes() int {
	sum := 0
	for _, v := range []interface{}{1, "hello", 3.14, int64(42)} {
		switch val := v.(type) {
		case int:
			sum += val
		case string:
			sum += len(val)
		case float64:
			sum += int(val)
		case int64:
			sum += int(val)
		}
	}
	return sum
}

// ============================================================================
// CHAINED BUILDER PATTERN WITH INTERFACES
// ============================================================================

// Builder interface.
type Builder interface {
	Build() string
}

// StringBuilderImpl for builder pattern.
type StringBuilderImpl struct {
	parts []string
}

func (s *StringBuilderImpl) Add(part string) *StringBuilderImpl {
	s.parts = append(s.parts, part)
	return s
}

func (s *StringBuilderImpl) Build() string {
	result := ""
	for _, p := range s.parts {
		result += p
	}
	return result
}

// ChainOnInterface tests method chaining on interface.
func ChainOnInterface() int {
	var builder Builder = &StringBuilderImpl{}
	impl := builder.(*StringBuilderImpl)
	impl.Add("hello")
	impl.Add("world")
	return len(impl.Build())
}
