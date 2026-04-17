// ops_container.go handles slice/map/chan creation, index, append, copy, delete, range, len, and cap.
package vm

import (
	"go/types"
	"reflect"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
)

// pushCommaOk pushes a (value, ok) tuple onto the operand stack.
// Used by OpIndexOk, OpRecvOk, OpTypeAssert, etc.
func (v *vm) pushCommaOk(val value.Value, ok bool) {
	tuple := []value.Value{val, value.MakeBool(ok)}
	v.push(value.FromInterface(tuple))
}

// resolveType resolves a type index from a popped value into a types.Type.
// Returns the type and true if found, or nil and false otherwise.
func (v *vm) resolveType(typeIdxVal value.Value) (types.Type, bool) {
	typeIdx := uint16(typeIdxVal.Int())
	if int(typeIdx) < len(v.program.Types) {
		return v.program.Types[typeIdx], true
	}
	return nil, false
}

// mustReflectValue extracts a reflect.Value from a value.Value or returns an
// invalid reflect.Value if the value doesn't contain a reflect.Value.
// This helper reduces repetitive if rv, ok := val.ReflectValue() patterns.
func (v *vm) mustReflectValue(val value.Value) reflect.Value {
	if rv, ok := val.ReflectValue(); ok {
		return rv
	}
	return reflect.Value{}
}

// executeContainer handles slice, map, channel creation, index, append,
// copy, delete, range, len, and cap opcodes.
func (v *vm) executeContainer(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen,maintidx,unparam // frame: uniform dispatch signature
	switch op {
	case bytecode.OpMakeSlice:
		capVal := v.pop()
		lenVal := v.pop()
		typeIdxVal := v.pop()
		typ, ok := v.resolveType(typeIdxVal)
		if !ok {
			v.push(value.MakeNil())
			break
		}
		made := false
		if sliceType, ok := typ.(*types.Slice); ok {
			elemType := sliceType.Elem()
			// Native []int64 fast path for integer slice types
			if basic, isBasic := elemType.(*types.Basic); isBasic {
				switch basic.Kind() {
				case types.Int, types.Int64:
					v.push(value.MakeIntSlice(make([]int64, int(lenVal.Int()), int(capVal.Int()))))
					made = true
				}
			}
			// Function slice: use reflect path to create proper typed slice (e.g. []func() int)
			// instead of []value.Value, so it can be assigned to typed struct fields.
		}
		if !made {
			if rt := typeToReflect(typ, v.program); rt != nil {
				slice := reflect.MakeSlice(rt, int(lenVal.Int()), int(capVal.Int()))
				v.push(value.MakeFromReflect(slice))
			} else {
				v.push(value.MakeNil())
			}
		}

	case bytecode.OpMakeMap:
		sizeVal := v.pop()
		typeIdxVal := v.pop()
		typ, ok := v.resolveType(typeIdxVal)
		if !ok {
			v.push(value.MakeNil())
			break
		}
		if rt := typeToReflect(typ, v.program); rt != nil {
			m := reflect.MakeMap(rt)
			_ = sizeVal // Size hint ignored for simplicity
			v.push(value.MakeFromReflect(m))
		} else {
			v.push(value.MakeNil())
		}

	case bytecode.OpMakeChan:
		sizeVal := v.pop()
		typeIdxVal := v.pop()
		typ, ok := v.resolveType(typeIdxVal)
		if !ok {
			v.push(value.MakeNil())
			break
		}
		if rt := typeToReflect(typ, v.program); rt != nil {
			ch := reflect.MakeChan(rt, int(sizeVal.Int()))
			v.push(value.MakeFromReflect(ch))
		} else {
			v.push(value.MakeNil())
		}

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
					rv.Index(idx).Set(val.ToReflectValue(rv.Type().Elem()))
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
			if (rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice) && rv.Type().Elem().Kind() == reflect.Int {
				n := rv.Len()
				high := int(highVal.Int())
				if high == sliceEndSentinel {
					high = n
				}
				s := make([]int64, n)
				for i := 0; i < n; i++ {
					s[i] = rv.Index(i).Int()
				}
				if maxVal.Kind() != value.KindNil && maxVal.Int() != sliceEndSentinel {
					v.push(value.MakeIntSlice(s[low:high:int(maxVal.Int())]))
				} else {
					v.push(value.MakeIntSlice(s[low:high]))
				}
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

	// Range operations
	case bytecode.OpRange:
		// Create an iterator for the collection
		collection := v.pop()
		v.push(value.FromInterface(&iterator{collection: collection, index: 0}))

	case bytecode.OpRangeNext:
		// Advance iterator and push a tuple (ok, key, value)
		iterVal := v.pop()
		iter, ok := iterVal.Interface().(*iterator)
		if !ok {
			// Return tuple (false, nil, nil)
			tuple := []value.Value{value.MakeBool(false), value.MakeNil(), value.MakeNil()}
			v.push(value.FromInterface(tuple))
			return nil
		}
		key, val, iterOk := iter.next()
		// SSA Next returns (ok, key, value) as a tuple
		tuple := []value.Value{value.MakeBool(iterOk), key, val}
		v.push(value.FromInterface(tuple))

	// Builtins
	case bytecode.OpLen:
		obj := v.pop()
		switch obj.Kind() {
		case value.KindString:
			v.push(value.MakeInt(int64(len(obj.String()))))
		case value.KindBytes:
			if b, ok := obj.Bytes(); ok {
				v.push(value.MakeInt(int64(len(b))))
			} else {
				v.push(value.MakeInt(0))
			}
		case value.KindSlice:
			v.push(value.MakeInt(int64(obj.Len())))
		case value.KindArray, value.KindMap, value.KindChan:
			v.push(value.MakeInt(int64(obj.Len())))
		case value.KindInterface, value.KindReflect:
			// Handle both interface values and reflect-wrapped values
			rv := v.mustReflectValue(obj)
			if rv.IsValid() {
				kind := rv.Kind()
				if kind == reflect.Interface {
					// Unwrap interface to get underlying value
					if !rv.IsNil() {
						rv = rv.Elem()
						kind = rv.Kind()
					}
				}
				switch kind {
				case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
					v.push(value.MakeInt(int64(rv.Len())))
				default:
					v.push(value.MakeInt(0))
				}
			} else {
				v.push(value.MakeInt(0))
			}
		default:
			v.push(value.MakeInt(0))
		}

	case bytecode.OpCap:
		obj := v.pop()
		switch obj.Kind() {
		case value.KindSlice, value.KindArray, value.KindChan:
			v.push(value.MakeInt(int64(obj.Cap())))
		case value.KindReflect:
			rv := v.mustReflectValue(obj)
			if rv.IsValid() {
				v.push(value.MakeInt(int64(rv.Cap())))
			} else {
				v.push(value.MakeInt(0))
			}
		default:
			v.push(value.MakeInt(0))
		}

	case bytecode.OpAppend:
		elem := v.pop()
		slice := v.pop()
		v.push(appendValue(slice, elem))

	case bytecode.OpCopy:
		src := v.pop()
		dst := v.pop()
		// Native int slice fast path
		if ds, ok := dst.IntSlice(); ok {
			if ss, ok2 := src.IntSlice(); ok2 {
				v.push(value.MakeInt(int64(copy(ds, ss))))
				break
			}
			// Cross-type: dst is native []int64, src is reflect slice (e.g. []int)
			if srcRV := v.mustReflectValue(src); srcRV.IsValid() && srcRV.Kind() == reflect.Slice {
				n := len(ds)
				if srcRV.Len() < n {
					n = srcRV.Len()
				}
				for i := 0; i < n; i++ {
					ds[i] = srcRV.Index(i).Int()
				}
				v.push(value.MakeInt(int64(n)))
				break
			}
		}
		// Copy slice
		if dstRV := v.mustReflectValue(dst); dstRV.IsValid() {
			if srcRV := v.mustReflectValue(src); srcRV.IsValid() {
				n := reflect.Copy(dstRV, srcRV)
				v.push(value.MakeInt(int64(n)))
			} else if ss, ok2 := src.IntSlice(); ok2 {
				// Cross-type: dst is reflect slice, src is native []int64
				n := dstRV.Len()
				if len(ss) < n {
					n = len(ss)
				}
				for i := 0; i < n; i++ {
					dstRV.Index(i).SetInt(ss[i])
				}
				v.push(value.MakeInt(int64(n)))
			} else {
				v.push(value.MakeInt(0))
			}
		} else {
			v.push(value.MakeInt(0))
		}

	case bytecode.OpDelete:
		key := v.pop()
		m := v.pop()
		// In Go, delete on a nil map is a no-op
		if m.IsNil() {
			break
		}
		if rv := v.mustReflectValue(m); rv.IsValid() && rv.IsNil() {
			break
		}
		// For OpDelete, we want to delete the entry (deleteIfNil=true)
		m.SetMapIndexWithDelete(key, value.MakeNil(), true)
	}

	return nil
}

// intSliceToReflect converts a native []int64 to a reflect []int slice.
func intSliceToReflect(s []int64) reflect.Value {
	rs := reflect.MakeSlice(reflect.TypeOf([]int{}), len(s), cap(s))
	for i, v := range s {
		rs.Index(i).SetInt(v)
	}
	return rs
}

// appendValue implements the append builtin for the VM.
// It handles native int slices, byte slices, reflect slices, and nil slices.
func appendValue(slice, elem value.Value) value.Value {
	// Fast path: native []int64 slice
	if s, ok := slice.IntSlice(); ok {
		return appendToIntSlice(s, elem)
	}

	// Byte slice ([]byte / KindBytes)
	if slice.Kind() == value.KindBytes {
		if b, ok := slice.Bytes(); ok {
			if elem.Kind() == value.KindUint || elem.Kind() == value.KindInt {
				return value.MakeBytes(append(b, byte(elem.Uint())))
			}
			// If elem is a byte slice (spread append)
			if elem.Kind() == value.KindBytes {
				if eb, ok := elem.Bytes(); ok {
					return value.MakeBytes(append(b, eb...))
				}
			}
			// Fallback: convert via interface
			if v := elem.Interface(); v != nil {
				if bv, ok := v.(byte); ok {
					return value.MakeBytes(append(b, bv))
				}
				if bv, ok := v.(uint8); ok {
					return value.MakeBytes(append(b, bv))
				}
			}
		}
		return slice
	}

	// Native []int64 that needs reflect conversion (e.g., stored in [][]int)
	if slice.Kind() == value.KindSlice {
		if intSlice, ok := slice.IntSlice(); ok {
			return appendIntSliceViaReflect(intSlice, elem)
		}
	}

	// Reflect-based slice
	if rv, ok := slice.ReflectValue(); ok {
		return appendToReflectSlice(rv, elem)
	}

	// Nil slice: create a new slice
	if slice.IsNil() || slice.Kind() == value.KindInvalid {
		return appendToNilSlice(elem)
	}

	return slice
}

// appendToIntSlice appends to a native []int64.
func appendToIntSlice(s []int64, elem value.Value) value.Value {
	if es, ok := elem.IntSlice(); ok {
		return value.MakeIntSlice(append(s, es...))
	}
	if elemRV, ok := elem.ReflectValue(); ok && elemRV.Kind() == reflect.Slice {
		// elem is a reflect-based integer slice (e.g. []int from a [][]int range)
		for i := 0; i < elemRV.Len(); i++ {
			s = append(s, elemRV.Index(i).Int())
		}
		return value.MakeIntSlice(s)
	}
	return value.MakeIntSlice(append(s, elem.RawInt()))
}

// appendIntSliceViaReflect converts []int64 to reflect []int and appends.
func appendIntSliceViaReflect(intSlice []int64, elem value.Value) value.Value {
	rv := intSliceToReflect(intSlice)
	if elem.Kind() == value.KindInt {
		return value.MakeFromReflect(reflect.Append(rv, reflect.ValueOf(int(elem.RawInt()))))
	}
	if elem.Kind() == value.KindSlice {
		if elemIntSlice, ok := elem.IntSlice(); ok {
			return value.MakeFromReflect(reflect.AppendSlice(rv, intSliceToReflect(elemIntSlice)))
		}
	}
	return value.MakeFromReflect(reflect.Append(rv, elem.ToReflectValue(reflect.TypeOf(int(0)))))
}

var valueValueType = reflect.TypeOf(value.Value{})

// appendToReflectSlice appends to a reflect.Value slice.
func appendToReflectSlice(rv reflect.Value, elem value.Value) value.Value {
	sliceElemType := rv.Type().Elem()

	// []value.Value slices (function slices)
	if sliceElemType == valueValueType {
		if elemRV, ok := elem.ReflectValue(); ok && elemRV.Kind() == reflect.Slice && elemRV.Type().Elem() == sliceElemType {
			return value.MakeFromReflect(reflect.AppendSlice(rv, elemRV))
		}
		return value.MakeFromReflect(reflect.Append(rv, reflect.ValueOf(elem)))
	}

	// Check if elem is native []int64 needing spread-append
	if elem.Kind() == value.KindSlice {
		if elemIntSlice, ok := elem.IntSlice(); ok {
			for _, v := range elemIntSlice {
				rv = reflect.Append(rv, reflect.ValueOf(int(v)))
			}
			return value.MakeFromReflect(rv)
		}
	}

	// SSA-packed variadic slice spread
	if elemRV, ok := elem.ReflectValue(); ok && elemRV.Kind() == reflect.Slice {
		return value.MakeFromReflect(reflect.AppendSlice(rv, elemRV))
	}

	return value.MakeFromReflect(reflect.Append(rv, elem.ToReflectValue(sliceElemType)))
}

// appendToNilSlice creates a new slice from a nil/zero slice and appends.
func appendToNilSlice(elem value.Value) value.Value {
	// Native []int64 spread
	if es, ok := elem.IntSlice(); ok {
		return value.MakeIntSlice(append([]int64(nil), es...))
	}

	elemRV, ok := elem.ReflectValue()
	if ok && elemRV.Kind() == reflect.Slice {
		sliceType := reflect.SliceOf(elemRV.Type().Elem())
		newSlice := reflect.MakeSlice(sliceType, 0, 0)
		return value.MakeFromReflect(reflect.AppendSlice(newSlice, elemRV))
	}

	// Single-element append: infer type from value
	elemIface := elem.Interface()
	if elemIface != nil {
		elemRV2 := reflect.ValueOf(elemIface)
		sliceType := reflect.SliceOf(elemRV2.Type())
		newSlice := reflect.MakeSlice(sliceType, 0, 0)
		return value.MakeFromReflect(reflect.Append(newSlice, elemRV2))
	}
	return value.MakeNil()
}
