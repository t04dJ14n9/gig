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

// Bug 8: Method dispatch type collision — json.Encoder.Encode vs xml.Encoder.Encode
//
// When a program uses both json.Encoder and xml.Encoder, the compiled program
// contains methods with the same name "Encode" on different receiver types.
// The method resolver picks xml.Encoder.Encode for json.Encoder.Encode calls,
// causing: panic: interface {} is *json.Encoder, not *xml.Encoder
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
