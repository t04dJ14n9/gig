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
		case value.KindReflect:
			// Handle reflect.Value containing a slice, array, or map
			if rv, ok := container.ReflectValue(); ok {
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
			if rv, ok := container.ReflectValue(); ok {
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
			if rv, ok := container.ReflectValue(); ok {
				if !rv.IsValid() {
					v.pushCommaOk(value.MakeNil(), false)
					break
				}
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
			container.SetMapIndex(key, val)
		case value.KindReflect:
			if rv, ok := container.ReflectValue(); ok {
				switch rv.Kind() {
				case reflect.Slice, reflect.Array:
					idx := int(key.Int())
					rv.Index(idx).Set(val.ToReflectValue(rv.Type().Elem()))
				case reflect.Map:
					container.SetMapIndex(key, val)
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

		if rv, ok := container.ReflectValue(); ok {
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
			if rv, ok := obj.ReflectValue(); ok {
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
			if rv, ok := obj.ReflectValue(); ok {
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
		// Native int slice fast path
		if s, ok := slice.IntSlice(); ok {
			if es, ok2 := elem.IntSlice(); ok2 {
				v.push(value.MakeIntSlice(append(s, es...)))
			} else if elemRV, ok2 := elem.ReflectValue(); ok2 && elemRV.Kind() == reflect.Slice {
				// elem is a reflect-based integer slice (e.g. []int from a [][]int range).
				// Convert each element to int64 and spread-append.
				for i := 0; i < elemRV.Len(); i++ {
					s = append(s, elemRV.Index(i).Int())
				}
				v.push(value.MakeIntSlice(s))
			} else {
				v.push(value.MakeIntSlice(append(s, elem.RawInt())))
			}
			break
		}
		// Handle KindSlice that contains native []int64 but needs to work with reflect
		// This happens when a []int was optimized to []int64 but is now stored in a [][]int
		if slice.Kind() == value.KindSlice {
			if intSlice, ok := slice.IntSlice(); ok {
				// This is a []int64 that needs to be converted to []int for reflect operations
				// Convert to []int for proper reflect handling
				intReflectSlice := reflect.MakeSlice(reflect.TypeOf([]int{}), len(intSlice), cap(intSlice))
				for i, v := range intSlice {
					intReflectSlice.Index(i).SetInt(v)
				}
				// Now handle the append with the converted slice
				if elem.Kind() == value.KindInt {
					newSlice := reflect.Append(intReflectSlice, reflect.ValueOf(int(elem.RawInt())))
					v.push(value.MakeFromReflect(newSlice))
				} else if elem.Kind() == value.KindSlice {
					if elemIntSlice, ok2 := elem.IntSlice(); ok2 {
						// Element is also []int64, convert to []int
						elemReflectSlice := reflect.MakeSlice(reflect.TypeOf([]int{}), len(elemIntSlice), cap(elemIntSlice))
						for i, v := range elemIntSlice {
							elemReflectSlice.Index(i).SetInt(v)
						}
						newSlice := reflect.AppendSlice(intReflectSlice, elemReflectSlice)
						v.push(value.MakeFromReflect(newSlice))
					} else {
						// Element is a different slice type
						newSlice := reflect.Append(intReflectSlice, elem.ToReflectValue(reflect.TypeOf(int(0))))
						v.push(value.MakeFromReflect(newSlice))
					}
				} else {
					newSlice := reflect.Append(intReflectSlice, elem.ToReflectValue(reflect.TypeOf(int(0))))
					v.push(value.MakeFromReflect(newSlice))
				}
				break
			}
		}
		// Append element to slice
		if rv, ok := slice.ReflectValue(); ok {
			sliceElemType := rv.Type().Elem()
			// Handle []value.Value slices (used for function slices)
			if sliceElemType == reflect.TypeOf(value.Value{}) {
				// Append value.Value element(s) to []value.Value
				if elemRV, ok2 := elem.ReflectValue(); ok2 && elemRV.Kind() == reflect.Slice && elemRV.Type().Elem() == sliceElemType {
					newSlice := reflect.AppendSlice(rv, elemRV)
					v.push(value.MakeFromReflect(newSlice))
				} else {
					newSlice := reflect.Append(rv, reflect.ValueOf(elem))
					v.push(value.MakeFromReflect(newSlice))
				}
			} else {
				// Check if elem is a native []int64 that needs spread-append
				if elem.Kind() == value.KindSlice {
					if elemIntSlice, ok2 := elem.IntSlice(); ok2 {
						// Element is []int64, spread append each element
						for _, v := range elemIntSlice {
							rv = reflect.Append(rv, reflect.ValueOf(int(v)))
						}
						v.push(value.MakeFromReflect(rv))
						break
					}
				}
				// Check if SSA packed variadic args into a slice (e.g., append(s, elems...))
				if elemRV, ok2 := elem.ReflectValue(); ok2 && elemRV.Kind() == reflect.Slice {
					// The element is a slice of the same element type — spread it
					newSlice := reflect.AppendSlice(rv, elemRV)
					v.push(value.MakeFromReflect(newSlice))
				} else {
					newSlice := reflect.Append(rv, elem.ToReflectValue(sliceElemType))
					v.push(value.MakeFromReflect(newSlice))
				}
			}
		} else if slice.IsNil() || slice.Kind() == value.KindInvalid {
			// Nil slice: create a new slice and append the element.
			// Fast path: elem is a native []int64 (spread-append to create new []int64).
			if es, ok2 := elem.IntSlice(); ok2 {
				v.push(value.MakeIntSlice(append([]int64(nil), es...)))
				break
			}
			elemRV, ok2 := elem.ReflectValue()
			if ok2 && elemRV.Kind() == reflect.Slice {
				if elemRV.Type().Elem() == reflect.TypeOf(value.Value{}) {
					// Create new []value.Value from the element slice
					newSlice := make([]value.Value, 0)
					newRV := reflect.ValueOf(newSlice)
					result := reflect.AppendSlice(newRV, elemRV)
					v.push(value.MakeFromReflect(result))
				} else {
					// Create new slice of the element's type and spread-append
					sliceType := reflect.SliceOf(elemRV.Type().Elem())
					newSlice := reflect.MakeSlice(sliceType, 0, 0)
					result := reflect.AppendSlice(newSlice, elemRV)
					v.push(value.MakeFromReflect(result))
				}
			} else {
				// Single-element append to a nil slice: infer element type from the value.
				elemIface := elem.Interface()
				if elemIface != nil {
					elemRV2 := reflect.ValueOf(elemIface)
					sliceType := reflect.SliceOf(elemRV2.Type())
					newSlice := reflect.MakeSlice(sliceType, 0, 0)
					result := reflect.Append(newSlice, elemRV2)
					v.push(value.MakeFromReflect(result))
				} else {
					v.push(slice)
				}
			}
		} else {
			v.push(slice)
		}

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
			if srcRV, ok2 := src.ReflectValue(); ok2 && srcRV.Kind() == reflect.Slice {
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
		if dstRV, ok := dst.ReflectValue(); ok {
			if srcRV, ok := src.ReflectValue(); ok {
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
		m.SetMapIndex(key, value.MakeNil())
	}

	return nil
}
