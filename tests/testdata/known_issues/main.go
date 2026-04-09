package known_issues

// ─────────────────────────────────────────────────────────────────────────────
// KNOWN ISSUES
//
// Resolved bugs are documented in testdata/resolved_issue/main.go as
// regression tests. This file tracks bugs awaiting fixes.
//
// Resolution mapping:
//   Bug 1 → Resolved Issue 28 (sort named-type conversion)
//   Bug 2 → fixed in gentool (time.Duration DirectCall)
//   Bug 3 → Resolved Issue 29 (fmt.Stringer)
//   Bug 4 → Resolved Issue 30 (fmt.Sprintf %T)
//   Bug 5 → Resolved Issue 31 (fmt.Sprintf %v _gig_id)
//   Bug 6 → Resolved Issue 32 (int64/uint64 narrowing)
//   Bug 7 → Resolved Issue 33 (bytes.Buffer.Cap)
//   Bug 8 → Resolved Issue 34 (json.Encoder method dispatch collision)
//
// ─────────────────────────────────────────────────────────────────────────────

import (
	"bytes"
	"encoding/json"
)

// ─────────────────────────────────────────────────────────────────────────────
// OTHER KNOWN BUGS
// ─────────────────────────────────────────────────────────────────────────────

// Bug 8: Method dispatch type collision — json.Encoder.Encode
//
// When a program uses multiple types with the same method name "Encode",
// the compiled program contains methods with the same name on different
// receiver types. The method resolver could pick the wrong one.
// This was fixed by including the full receiver type in the method cache key.
//
// This is the same reflect.StructOf type identity issue that was partially
// fixed for other cases. The fix likely requires making the method cache
// key include the full receiver type, not just the method name.

// JsonEncodeBug8 tests json.NewEncoder.Encode call.
func JsonEncodeBug8() int {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(map[string]int{"y": 20})
	return buf.Len()
}

// ─────────────────────────────────────────────────────────────────────────────
// PANIC/RECOVER/DEFER BUGS
// ─────────────────────────────────────────────────────────────────────────────

// PanicDeferBug: Variable assignment in defer closure fails when recovering from panic.
//
// Root cause: In vm/ops_control.go, defer closures are executed during panic unwind,
// but the closure's captured variable assignment (result = 42) fails because
// value.SetElem() cannot properly handle the closure's upvalue context.
//
// The interpreter sets v.panicking=true and runs defers, but defer closures
// cannot correctly modify captured variables.

// PanicRecoverBasic_Bug tests basic panic and recover.
func PanicRecoverBasic_Bug() int {
	defer func() {
		recover()
	}()
	panic("test panic")
	return 0 // never reached
}

// PanicRecoverWithValue_Bug tests recovering panic value and assigning to captured var.
func PanicRecoverWithValue_Bug() int {
	var result int
	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok && s == "expected" {
				result = 42
			}
		}
	}()
	panic("expected")
	return result
}

// DeferRunsOnPanic_Bug tests deferred function modifying captured var during panic.
func DeferRunsOnPanic_Bug() int {
	result := 0
	defer func() {
		result += 10
		recover()
	}()
	result += 1
	panic("test")
	return result // never reached
}

// MultipleDefersOnPanic_Bug tests LIFO order of deferred functions during panic.
func MultipleDefersOnPanic_Bug() int {
	order := 0
	result := 0
	defer func() {
		order++
		result = result*10 + order
		recover()
	}()
	defer func() {
		order++
		result = result*10 + order
	}()
	defer func() {
		order++
		result = result*10 + order
	}()
	panic("test")
	return result
}

// NamedReturnPanicRecover_Bug tests named return value with panic/recover.
func NamedReturnPanicRecover_Bug() (result int) {
	defer func() {
		if recover() != nil {
			result = 42
		}
	}()
	panic("test")
	return
}

// NestedRecover_Bug tests recover in nested defer during panic chain.
func NestedRecover_Bug() int {
	result := 0
	defer func() {
		defer func() {
			if r := recover(); r != nil {
				result = 100
			}
		}()
		panic("second panic")
	}()
	defer func() {
		recover()
	}()
	panic("first panic")
	return result
}

// PanicInDefer_Bug tests panic in deferred function.
func PanicInDefer_Bug() int {
	result := 0
	defer func() {
		if r := recover(); r != nil {
			result = 50
		}
	}()
	defer func() {
		panic("panic in defer")
	}()
	result = 1
	return result
}

// PanicInClosure_Bug tests panic inside closure with recover.
func PanicInClosure_Bug() int {
	fn := func() {
		panic("closure panic")
	}
	defer func() {
		recover()
	}()
	fn()
	return 0
}

// DeferPanicRecoverChain_Bug tests chain of defer/panic/recover.
func DeferPanicRecoverChain_Bug() int {
	result := 0
	defer func() {
		result += 1000
		recover()
	}()
	defer func() {
		result += 100
		defer func() {
			result += 10
			recover()
		}()
		panic("inner")
	}()
	panic("outer")
}

// ─────────────────────────────────────────────────────────────────────────────
// STRANGE SYNTAX BUGS (Found by comprehensive testing 2026-03-28)
// ─────────────────────────────────────────────────────────────────────────────

// StrangeSyntax_Bug1: Nil slice to interface conversion loses type information
//
// When a nil slice is returned as interface{}, the interpreter returns nil
// instead of preserving the slice type. In native Go, a nil slice retains
// its type information even when assigned to interface{}.
//
// Expected: [] ([]int) - empty slice with type info
// Got: <nil> (<nil>) - loses type
//
// Root cause: Likely in value conversion or interface boxing logic.
func StrangeSyntax_Bug1_ConvertNilToInterface() interface{} {
	var s []int
	return s // nil slice assigned to interface
}

// StrangeSyntax_Bug2: Nil map access returns nil instead of zero value
//
// Accessing a non-existent key in a nil map should return the zero value
// for the value type, but the interpreter returns nil instead.
//
// Expected: 0 (int) - zero value for int
// Got: <nil> (<nil>)
//
// Root cause: Map access on nil maps doesn't properly handle zero values.
func StrangeSyntax_Bug2_NilMapAccess() int {
	var m map[string]int
	return m["key"] // Returns zero value
}

// StrangeSyntax_Bug3: Delete on nil map causes panic
//
// Deleting from a nil map should be a no-op in Go, but the interpreter
// panics with "invalid reflect.Value in SetMapIndex()".
//
// Expected: No-op, silently succeed
// Got: panic: invalid reflect.Value in SetMapIndex()
//
// Root cause: VM map delete operation doesn't check for nil map.
func StrangeSyntax_Bug3_NilMapDelete() int {
	var m map[string]int
	delete(m, "key") // No-op
	return 0
}

// StrangeSyntax_Bug4: Blank identifier expression loses type in interface return
//
// Functions returning interface{} from a blank identifier assignment lose
// type information. The interpreter returns nil instead of preserving type.
//
// Expected: [] ([]interface {}) - empty slice with type
// Got: <nil> (<nil>)
//
// Root cause: Blank identifier expressions don't preserve type information.
func StrangeSyntax_Bug4_BlankExpression() interface{} {
	_ = 42
	var s []interface{}
	return s
}

// StrangeSyntax_Bug5: Send on closed channel panic handling
//
// Sending on a closed channel should panic, and that panic should be
// recoverable. The interpreter does panic, but there may be issues with
// the test framework's handling of this case.
//
// Expected: Panic should be caught by defer/recover
// Got: Test framework issue with panic handling
//
// Root cause: Needs investigation - may be test framework or panic handling.
func StrangeSyntax_Bug5_ChannelClosedSend() int {
	ch := make(chan int, 1)
	close(ch)
	defer func() {
		if r := recover(); r != nil {
			// Recovered from panic
		}
	}()
	ch <- 1 // Will panic
	return 0
}

// StrangeSyntax_Bug6: Nil function return loses type information
//
// Returning nil as a function type should preserve the function type
// information, but the interpreter returns untyped nil.
//
// Expected: <nil> (func() int) - nil with function type
// Got: <nil> (<nil>) - loses type
//
// Root cause: Similar to Bug1 - nil value to interface loses type info.
func StrangeSyntax_Bug6_ClosureReturnNil() func() int {
	if false {
		return func() int { return 1 }
	}
	return nil
}
