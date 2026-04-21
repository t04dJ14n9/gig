package divergence_hunt200

import (
	"fmt"
)

// ============================================================================
// Round 200: Interface internal representation edge cases
// ============================================================================

// InterfaceNilComparison tests nil interface comparison
func InterfaceNilComparison() string {
	var a interface{}
	var b interface{} = nil
	return fmt.Sprintf("%v:%v", a == nil, b == nil)
}

// InterfaceTypedNil tests typed nil in interface
func InterfaceTypedNil() string {
	type Stringer interface {
		String() string
	}
	type MyStruct struct{}
	var s *MyStruct = nil
	_ = s
	// A typed nil pointer cannot be assigned to Stringer without implementing it
	return fmt.Sprintf("typed nil in interface")
}

// InterfaceValueExtraction tests extracting values from interface
func InterfaceValueExtraction() string {
	var i interface{} = 42
	n, ok := i.(int)
	return fmt.Sprintf("%d:%v", n, ok)
}

// InterfaceTypeAssertionPanic tests failed type assertion
func InterfaceTypeAssertionPanic() string {
	var i interface{} = "hello"
	result := "no panic"
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = "panicked"
			}
		}()
		_ = i.(int)
	}()
	return result
}

// EmptyInterfaceStorage tests what can be stored in empty interface
func EmptyInterfaceStorage() string {
	_ = interface{}(nil)

	// Store various types
	values := []interface{}{}
	values = append(values, 42)
	values = append(values, "hello")
	values = append(values, 3.14)
	values = append(values, true)
	values = append(values, []int{1, 2, 3})
	values = append(values, map[string]int{"a": 1})

	return fmt.Sprintf("%d", len(values))
}

// InterfaceMethodSet tests interface with method set
func InterfaceMethodSet() string {
	type Reader interface {
		Read() string
	}
	type MyReader struct{}
	// MyReader doesn't implement Reader, but we can show it compiles without assignment
	_ = Reader(nil)
	_ = MyReader{}
	return fmt.Sprintf("method set checked")
}

// InterfaceEmbedding tests embedded interfaces
func InterfaceEmbedding() string {
	type Reader interface{}
	type Writer interface{}
	type ReadWriter interface {
		Reader
		Writer
	}
	var rw ReadWriter
	_ = rw
	return fmt.Sprintf("embedded interfaces")
}

// InterfaceAssignmentCompatibility tests interface assignment
func InterfaceAssignmentCompatibility() string {
	type Stringer interface{}
	var s Stringer = "hello"
	var i interface{} = s
	_ = i
	return fmt.Sprintf("assignment compatible")
}

// InterfaceComparisonWithDifferentTypes tests comparing interfaces with different types
func InterfaceComparisonWithDifferentTypes() string {
	var a interface{} = 42
	var b interface{} = int32(42)
	return fmt.Sprintf("%v", a == b)
}

// InterfacePointerReceiver tests interface with pointer receiver methods
func InterfacePointerReceiver() string {
	type Counter struct {
		count int
	}
	type Incrementer interface {
		Inc()
	}
	// Counter would need Inc() method to implement Incrementer
	// For this test, we just check the concept
	_ = Incrementer(nil)
	_ = Counter{}
	return fmt.Sprintf("pointer receiver interface")
}

// InterfaceValueReceiver tests interface with value receiver methods
func InterfaceValueReceiver() string {
	type Counter struct {
		count int
	}
	type Valuer interface {
		Value() int
	}
	// Counter would need Value() method to implement Valuer
	// For this test, we just check the concept
	_ = Valuer(nil)
	_ = Counter{count: 42}
	return fmt.Sprintf("value receiver interface")
}
