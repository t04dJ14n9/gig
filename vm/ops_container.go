// ops_container.go dispatches VM container opcodes to focused handlers.
package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// executeContainer handles slice, map, channel creation, index, append,
// copy, delete, range, len, and cap opcodes.
func (v *vm) executeContainer(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen,maintidx,unparam // frame: uniform dispatch signature
	switch op {
	case bytecode.OpMakeSlice:
		v.executeMakeSlice()

	case bytecode.OpMakeMap:
		v.executeMakeMap()

	case bytecode.OpMakeChan:
		v.executeMakeChan()

	// Index operations
	case bytecode.OpIndex:
		key := v.pop()
		container := v.pop()
		switch container.Kind() {
		case value.KindSlice:
			// Native int slice fast path
			if s, ok := container.IntSlice(); ok {
				v.push(value.MakeInt(s[int(key.RawInt())]))
			} else {
				v.push(container.Index(int(key.Int())))
			}
		case value.KindArray:
			idx := int(key.Int())
			v.push(container.Index(idx))
		case value.KindMap:
			v.push(container.MapIndex(key))
		case value.KindString:
			idx := int(key.Int())
			v.push(container.Index(idx))
		case value.KindBytes:
			// Native []byte indexing — return uint8 as KindUint
			if b, ok := container.Bytes(); ok {
				v.push(value.MakeUint8(b[int(key.RawInt())]))
			} else {
				v.push(value.MakeNil())
			}
		case value.KindReflect:
			// Handle reflect.Value containing a slice, array, or map
			rv := v.mustReflectValue(container)
			if rv.IsValid() {
				switch rv.Kind() {
				case reflect.Slice, reflect.Array:
					idx := int(key.Int())
					v.push(value.MakeFromReflect(rv.Index(idx)))
				case reflect.Map:
					k := key.ToReflectValue(rv.Type().Key())
					elem := rv.MapIndex(k)
					if !elem.IsValid() {
						// Return zero value of element type, not nil (Go semantics)
						v.push(value.MakeFromReflect(reflect.Zero(rv.Type().Elem())))
					} else {
						v.push(value.MakeFromReflect(elem))
					}
				default:
					v.push(value.MakeNil())
				}
			} else {
				v.push(value.MakeNil())
			}
		default:
			v.push(value.MakeNil())
		}

	case bytecode.OpIndexOk:
		// Index with comma-ok: returns (value, ok) tuple for maps
		key := v.pop()
		container := v.pop()
		switch container.Kind() {
		case value.KindMap:
			// For maps, check if key exists
			rv := v.mustReflectValue(container)
			if rv.IsValid() {
				k := key.ToReflectValue(rv.Type().Key())
				elem := rv.MapIndex(k)
				if !elem.IsValid() {
					v.pushCommaOk(value.MakeFromReflect(reflect.Zero(rv.Type().Elem())), false)
				} else {
					v.pushCommaOk(value.MakeFromReflect(elem), true)
				}
			} else {
				v.pushCommaOk(value.MakeNil(), false)
			}
		case value.KindReflect:
			rv := v.mustReflectValue(container)
			if rv.IsValid() {
				switch rv.Kind() {
				case reflect.Map:
					if rv.Type().Key() == nil {
						v.pushCommaOk(value.MakeNil(), false)
						break
					}
					k := key.ToReflectValue(rv.Type().Key())
					if !k.IsValid() {
						v.pushCommaOk(value.MakeNil(), false)
						break
					}
					elem := rv.MapIndex(k)
					if !elem.IsValid() {
						v.pushCommaOk(value.MakeFromReflect(reflect.Zero(rv.Type().Elem())), false)
					} else {
						v.pushCommaOk(value.MakeFromReflect(elem), true)
					}
				case reflect.Slice, reflect.Array:
					idx := int(key.Int())
					if idx < 0 || idx >= rv.Len() {
						v.pushCommaOk(value.MakeNil(), false)
					} else {
						v.pushCommaOk(value.MakeFromReflect(rv.Index(idx)), true)
					}
				default:
					v.pushCommaOk(value.MakeNil(), false)
				}
			} else {
				v.pushCommaOk(value.MakeNil(), false)
			}
		default:
			v.pushCommaOk(value.MakeNil(), false)
		}

	case bytecode.OpSetIndex:
		val := v.pop()
		key := v.pop()
		container := v.pop()
		switch container.Kind() {
		case value.KindSlice:
			// Native int slice fast path
			if s, ok := container.IntSlice(); ok {
				s[int(key.RawInt())] = val.RawInt()
			} else {
				container.SetIndex(int(key.Int()), val)
			}
		case value.KindArray:
			idx := int(key.Int())
			container.SetIndex(idx, val)
		case value.KindMap:
			// For OpSetIndex, nil value means set to typed nil (not delete)
			container.SetMapIndexWithDelete(key, val, false)
		case value.KindReflect:
			rv := v.mustReflectValue(container)
			if rv.IsValid() {
				switch rv.Kind() {
				case reflect.Slice, reflect.Array:
					idx := int(key.Int())
					rv.Index(idx).Set(v.valueForReflectSet(val, rv.Type().Elem()))
				case reflect.Map:
					// For OpSetIndex, nil value means set to typed nil (not delete)
					container.SetMapIndexWithDelete(key, val, false)
				}
			}
		}

	case bytecode.OpSlice:
		// Slice operation: container[low:high:max]
		maxVal := v.pop()
		highVal := v.pop()
		lowVal := v.pop()
		container := v.pop()

		low := int(lowVal.Int())

		// Handle nil container: nil[0:0] returns nil, nil[0:n] panics for n > 0
		if container.Kind() == value.KindNil {
			high := int(highVal.Int())
			if high == sliceEndSentinel {
				high = 0
			}
			if low != 0 || high != 0 {
				panic("runtime error: slice bounds out of range")
			}
			// nil[0:0] returns a nil slice with the same type (Go semantics)
			// Use container's original type if available, otherwise just push nil
			v.push(container)
			break
		}

		// Handle string slicing specially
		if container.Kind() == value.KindString {
			high := int(highVal.Int())
			if high == sliceEndSentinel {
				high = len(container.String())
			}
			v.push(value.MakeString(container.String()[low:high]))
			break
		}

		// Native []byte slice fast path
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
				break
			}
		}

		// Native []int64 slice fast path
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
			break
		}

		rv := v.mustReflectValue(container)
		if rv.IsValid() {
			// Nil slice check: in Go, nil[0:0] returns nil, not an empty non-nil slice.
			if rv.Kind() == reflect.Slice && rv.IsNil() {
				high := int(highVal.Int())
				if high == sliceEndSentinel {
					high = 0
				}
				if low != 0 || high != 0 {
					panic("runtime error: slice bounds out of range")
				}
				// Return nil slice with the correct type
				v.push(container)
				break
			}
			// If it's a pointer to an array or slice, dereference it first
			if rv.Kind() == reflect.Ptr {
				elemKind := rv.Elem().Kind()
				if elemKind == reflect.Array || elemKind == reflect.Slice {
					rv = rv.Elem()
				}
			}

			// Native int array/slice → []int64 fast path
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
				break
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
				break
			}

			// Handle native []value.Value slices (used for function slices)
			if rv.Kind() == reflect.Slice && rv.Type().Elem() == reflect.TypeOf(value.Value{}) {
				high := int(highVal.Int())
				if high == sliceEndSentinel {
					high = rv.Len()
				}
				sliced := rv.Slice(low, high)
				v.push(value.MakeFromReflect(sliced))
				break
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
		} else {
			// Nil slice subslicing: in Go, nil[0:0] returns nil (not an empty non-nil slice).
			// Since rv is invalid, the container is nil. Return nil to match Go semantics.
			v.push(value.MakeNil())
		}

	case bytecode.OpRange:
		v.executeRange()

	case bytecode.OpRangeNext:
		v.executeRangeNext()

	case bytecode.OpLen:
		v.executeLen()

	case bytecode.OpCap:
		v.executeCap()

	case bytecode.OpAppend:
		elem := v.pop()
		slice := v.pop()
		v.push(appendValue(slice, elem))

	case bytecode.OpCopy:
		v.executeCopy()

	case bytecode.OpDelete:
		v.executeDelete()

	}

	return nil
}
