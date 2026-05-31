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
	high := int(highVal.Int())
	// Normalize the optional max bound once at the opcode boundary. The
	// helper paths below stay on primitive ints, which keeps the hot bound
	// normalization small enough for the compiler to inline.
	max, hasMax := sliceMaxIndex(maxVal)

	// Keep the common native slice cases in the opcode body. These are the hot
	// paths emitted for ordinary Gig source, while nil and reflect handling stay
	// below as small named domains.
	switch container.Kind() {
	case value.KindNil:
		panicIfNonZeroNilSliceBounds(low, high)
		v.push(container)
		return
	case value.KindString:
		s := container.String()
		v.push(value.MakeString(s[low:sliceHighIndex(high, len(s))]))
		return
	case value.KindBytes:
		if b, ok := container.Bytes(); ok {
			high := sliceHighIndex(high, len(b))
			if hasMax {
				v.push(value.MakeBytes(b[low:high:max]))
			} else {
				v.push(value.MakeBytes(b[low:high]))
			}
			return
		}
	}

	if s, ok := container.IntSlice(); ok {
		high := sliceHighIndex(high, len(s))
		if hasMax {
			v.push(value.MakeIntSlice(s[low:high:max]))
		} else {
			v.push(value.MakeIntSlice(s[low:high]))
		}
		return
	}

	rv := v.mustReflectValue(container)
	if v.executeNilReflectSlice(container, rv, low, high) {
		return
	}

	rv = dereferenceSliceTarget(rv)

	if v.executeReflectIntSlice(rv, low, high, max, hasMax) {
		return
	}

	v.push(value.MakeFromReflect(sliceReflectValue(rv, low, high, max, hasMax)))
}

func (v *vm) executeNilReflectSlice(container value.Value, rv reflect.Value, low, high int) bool {
	if !rv.IsValid() {
		v.push(value.MakeNil())
		return true
	}
	if rv.Kind() != reflect.Slice || !rv.IsNil() {
		return false
	}
	panicIfNonZeroNilSliceBounds(low, high)
	v.push(container)
	return true
}

func panicIfNonZeroNilSliceBounds(low, high int) {
	if low != 0 || sliceHighIndex(high, 0) != 0 {
		panic("runtime error: slice bounds out of range")
	}
}

func dereferenceSliceTarget(rv reflect.Value) reflect.Value {
	if rv.Kind() != reflect.Ptr {
		return rv
	}
	elemKind := rv.Elem().Kind()
	if elemKind == reflect.Array || elemKind == reflect.Slice {
		return rv.Elem()
	}
	return rv
}

func sliceHighIndex(high, length int) int {
	if high == sliceEndSentinel {
		return length
	}
	return high
}

func sliceMaxIndex(maxVal value.Value) (int, bool) {
	if maxVal.Kind() == value.KindNil {
		return 0, false
	}
	max := int(maxVal.Int())
	return max, max != sliceEndSentinel
}

func sliceReflectValue(rv reflect.Value, low, high, max int, hasMax bool) reflect.Value {
	high = sliceHighIndex(high, rv.Len())
	if hasMax {
		return rv.Slice3(low, high, max)
	}
	return rv.Slice(low, high)
}

func (v *vm) executeReflectIntSlice(rv reflect.Value, low, high, max int, hasMax bool) bool {
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
	if isReflectIntArrayOrSlice(rv) {
		v.push(value.MakeFromReflect(sliceReflectValue(rv, low, high, max, hasMax)))
		return true
	}

	// Handle native []value.Value slices used for function slices.
	if rv.Kind() == reflect.Slice && rv.Type().Elem() == reflect.TypeOf(value.Value{}) {
		high := sliceHighIndex(high, rv.Len())
		sliced := rv.Slice(low, high)
		v.push(value.MakeFromReflect(sliced))
		return true
	}

	return false
}

func isReflectIntArrayOrSlice(rv reflect.Value) bool {
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return false
	}
	return rv.Type().Elem().Kind() == reflect.Int
}
