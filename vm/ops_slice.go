package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeSlice() {
	// Slice operation: container[low:high:max]
	maxVal := v.pop()
	highVal := v.pop()
	lowVal := v.pop()
	container := v.pop()

	low := int(lowVal.Int())

	// Handle nil container: nil[0:0] returns nil, nil[0:n] panics for n > 0.
	if container.Kind() == value.KindNil {
		high := int(highVal.Int())
		if high == sliceEndSentinel {
			high = 0
		}
		if low != 0 || high != 0 {
			panic("runtime error: slice bounds out of range")
		}
		// nil[0:0] returns a nil slice with the same type (Go semantics).
		v.push(container)
		return
	}

	if v.executeSpecialSlice(container, low, highVal, maxVal) {
		return
	}

	rv := v.mustReflectValue(container)
	if !rv.IsValid() {
		// Nil slice subslicing: in Go, nil[0:0] returns nil, not an empty non-nil slice.
		v.push(value.MakeNil())
		return
	}

	// Nil slice check: in Go, nil[0:0] returns nil, not an empty non-nil slice.
	if rv.Kind() == reflect.Slice && rv.IsNil() {
		high := int(highVal.Int())
		if high == sliceEndSentinel {
			high = 0
		}
		if low != 0 || high != 0 {
			panic("runtime error: slice bounds out of range")
		}
		v.push(container)
		return
	}

	// If it's a pointer to an array or slice, dereference it first.
	if rv.Kind() == reflect.Ptr {
		elemKind := rv.Elem().Kind()
		if elemKind == reflect.Array || elemKind == reflect.Slice {
			rv = rv.Elem()
		}
	}

	if v.executeReflectIntSlice(container, rv, low, highVal, maxVal) {
		return
	}

	high := int(highVal.Int())
	if high == sliceEndSentinel {
		high = rv.Len()
	}

	var sliced reflect.Value
	if maxVal.Kind() != value.KindNil && maxVal.Int() != sliceEndSentinel {
		// 3-index slice: container[low:high:max]
		max := int(maxVal.Int())
		sliced = rv.Slice3(low, high, max)
	} else {
		// 2-index slice: container[low:high]
		sliced = rv.Slice(low, high)
	}
	v.push(value.MakeFromReflect(sliced))
}

func (v *vm) executeSpecialSlice(container value.Value, low int, highVal, maxVal value.Value) bool {
	if container.Kind() == value.KindString {
		high := int(highVal.Int())
		if high == sliceEndSentinel {
			high = len(container.String())
		}
		v.push(value.MakeString(container.String()[low:high]))
		return true
	}

	if container.Kind() == value.KindBytes {
		if b, ok := container.Bytes(); ok {
			high := int(highVal.Int())
			if high == sliceEndSentinel {
				high = len(b)
			}
			if maxVal.Kind() != value.KindNil && maxVal.Int() != sliceEndSentinel {
				v.push(value.MakeBytes(b[low:high:int(maxVal.Int())]))
			} else {
				v.push(value.MakeBytes(b[low:high]))
			}
			return true
		}
	}

	if s, ok := container.IntSlice(); ok {
		high := int(highVal.Int())
		if high == sliceEndSentinel {
			high = len(s)
		}
		if maxVal.Kind() != value.KindNil && maxVal.Int() != sliceEndSentinel {
			v.push(value.MakeIntSlice(s[low:high:int(maxVal.Int())]))
		} else {
			v.push(value.MakeIntSlice(s[low:high]))
		}
		return true
	}

	return false
}

func (v *vm) executeReflectIntSlice(container value.Value, rv reflect.Value, low int, highVal, maxVal value.Value) bool {
	// Native int array/slice -> []int64 fast path
	// SSA compiles make([]int, N) with constant N as Alloc([N]int) + Slice,
	// so we intercept it here to produce a native []int64.
	//
	// IMPORTANT: For arrays, we must NOT copy to []int64 because that
	// breaks shared-underlying-array semantics (e.g., a[1:3] and a[2:4]
	// must share memory). Use rv.Slice() instead which preserves sharing.
	// For slices from reflect, also use rv.Slice() to preserve sharing.
	// The []int64 fast path is only used when the source is already
	// a native []int64 (handled by the IntSlice() check above).
	if rv.Kind() == reflect.Array && rv.Type().Elem().Kind() == reflect.Int {
		high := int(highVal.Int())
		if high == sliceEndSentinel {
			high = rv.Len()
		}
		var sliced reflect.Value
		if maxVal.Kind() != value.KindNil && maxVal.Int() != sliceEndSentinel {
			sliced = rv.Slice3(low, high, int(maxVal.Int()))
		} else {
			sliced = rv.Slice(low, high)
		}
		v.push(value.MakeFromReflect(sliced))
		return true
	}
	if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() == reflect.Int {
		high := int(highVal.Int())
		if high == sliceEndSentinel {
			high = rv.Len()
		}
		var sliced reflect.Value
		if maxVal.Kind() != value.KindNil && maxVal.Int() != sliceEndSentinel {
			sliced = rv.Slice3(low, high, int(maxVal.Int()))
		} else {
			sliced = rv.Slice(low, high)
		}
		v.push(value.MakeFromReflect(sliced))
		return true
	}

	// Handle native []value.Value slices used for function slices.
	if rv.Kind() == reflect.Slice && rv.Type().Elem() == reflect.TypeOf(value.Value{}) {
		high := int(highVal.Int())
		if high == sliceEndSentinel {
			high = rv.Len()
		}
		sliced := rv.Slice(low, high)
		v.push(value.MakeFromReflect(sliced))
		return true
	}

	_ = container
	return false
}
