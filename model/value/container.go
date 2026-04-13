// container.go implements collection operations: Index, SetIndex, MapIndex, SetMapIndex,
// Len, Cap, Field, SetField, Elem, SetElem, and Append.
package value

import (
	"fmt"
	"reflect"
	"unsafe"
)

// UnsafeAddrOf returns the unsafe.Pointer for a reflect.Value that is addressable.
// This is used internally by the VM to obtain settable pointers to unexported struct fields.
func UnsafeAddrOf(v reflect.Value) unsafe.Pointer {
	return v.Addr().UnsafePointer()
}

// Len returns the length of string, slice, array, map, or chan.
func (v Value) Len() int {
	switch v.kind {
	case KindString:
		return len(v.obj.(string))
	case KindSlice:
		// Native int slice fast path
		if s, ok := v.obj.([]int64); ok {
			return len(s)
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Len()
		}
		panic("invalid obj in Len()")
	case KindArray, KindMap, KindChan:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Len()
		}
		panic("invalid reflect.Value in Len()")
	case KindReflect:
		if rv, ok := v.obj.(reflect.Value); ok {
			switch rv.Kind() {
			case reflect.String, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
				return rv.Len()
			}
		}
		panic(fmt.Sprintf("cannot take len of reflect kind %v", v.obj))
	default:
		panic(fmt.Sprintf("cannot take len of %v", v.kind))
	}
}

// Cap returns the capacity of slice, array, or chan.
func (v Value) Cap() int {
	switch v.kind {
	case KindSlice:
		if s, ok := v.obj.([]int64); ok {
			return cap(s)
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Cap()
		}
		panic("invalid obj in Cap()")
	case KindArray, KindChan:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Cap()
		}
		panic("invalid reflect.Value in Cap()")
	case KindReflect:
		if rv, ok := v.obj.(reflect.Value); ok {
			switch rv.Kind() {
			case reflect.Slice, reflect.Array, reflect.Chan:
				return rv.Cap()
			}
		}
		panic(fmt.Sprintf("cannot take cap of reflect kind %v", v.obj))
	default:
		panic(fmt.Sprintf("cannot take cap of %v", v.kind))
	}
}

// Index returns element at index i for slice, array, or string.
func (v Value) Index(i int) Value {
	switch v.kind {
	case KindString:
		// s[i] returns a byte (uint8), not a string
		return MakeUint8(v.obj.(string)[i])
	case KindSlice:
		// Native int slice fast path
		if s, ok := v.obj.([]int64); ok {
			return MakeInt(s[i])
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return indexReflectSlice(rv, i)
		}
		if slice, ok := v.obj.([]Value); ok {
			return slice[i]
		}
		panic("invalid obj in Index()")
	case KindArray:
		if rv, ok := v.obj.(reflect.Value); ok {
			return indexReflectSlice(rv, i)
		}
		if slice, ok := v.obj.([]Value); ok {
			return slice[i]
		}
		panic("invalid reflect.Value in Index()")
	case KindReflect:
		// Handle reflect.Value containing a slice
		if rv, ok := v.obj.(reflect.Value); ok {
			return indexReflectSlice(rv, i)
		}
		// Handle native []value.Value slice
		if slice, ok := v.obj.([]Value); ok {
			return slice[i]
		}
		panic("invalid reflect.Value in Index()")
	default:
		panic(fmt.Sprintf("cannot index %v", v.kind))
	}
}

// indexReflectSlice handles indexing into a reflect.Value slice/array.
func indexReflectSlice(rv reflect.Value, i int) Value {
	elem := rv.Index(i)
	if rv.Type().Elem().Kind() == reflect.Func {
		if val, ok := elem.Interface().(Value); ok {
			return val
		}
	}
	if rv.Type().Elem() == reflect.TypeOf(Value{}) {
		return elem.Interface().(Value)
	}
	return MakeFromReflect(elem)
}

// SetIndex sets element at index i for slice or array.
func (v Value) SetIndex(i int, val Value) {
	// Native int slice fast path
	if v.kind == KindSlice {
		if s, ok := v.obj.([]int64); ok {
			s[i] = val.RawInt()
			return
		}
	}
	if rv, ok := v.obj.(reflect.Value); ok {
		elemType := rv.Type().Elem()
		if elemType.Kind() == reflect.Func {
			rv.Index(i).Set(val.ToReflectValue(elemType))
			return
		}
		if elemType == reflect.TypeOf(Value{}) {
			rv.Index(i).Set(reflect.ValueOf(val))
			return
		}
		rv.Index(i).Set(val.ToReflectValue(elemType))
		return
	}
	if slice, ok := v.obj.([]Value); ok {
		slice[i] = val
		return
	}
	panic("invalid reflect.Value in SetIndex()")
}

// MapIndex returns value at key k for map.
func (v Value) MapIndex(k Value) Value {
	if rv, ok := v.obj.(reflect.Value); ok {
		key := k.ToReflectValue(rv.Type().Key())
		elem := rv.MapIndex(key)
		if !elem.IsValid() {
			// Return zero value of element type, not nil (Go semantics)
			return MakeFromReflect(reflect.Zero(rv.Type().Elem()))
		}
		return MakeFromReflect(elem)
	}
	panic("invalid reflect.Value in MapIndex()")
}

// SetMapIndexWithDelete sets a map entry with control over nil value handling.
// When deleteIfNil is true and val is nil, the key is deleted from the map.
// When deleteIfNil is false and val is nil, the key is set to a typed nil value.
func (v Value) SetMapIndexWithDelete(k, val Value, deleteIfNil bool) {
	if rv, ok := v.obj.(reflect.Value); ok {
		key := k.ToReflectValue(rv.Type().Key())
		if val.IsNil() {
			if deleteIfNil {
				// Delete the entry
				rv.SetMapIndex(key, reflect.Value{})
			} else {
				// Set to typed nil value (e.g., nil interface{})
				rv.SetMapIndex(key, reflect.Zero(rv.Type().Elem()))
			}
		} else {
			rv.SetMapIndex(key, val.ToReflectValue(rv.Type().Elem()))
		}
		return
	}
	panic("invalid reflect.Value in SetMapIndex()")
}

// Field returns struct field at index i.
func (v Value) Field(i int) Value {
	if rv, ok := v.obj.(reflect.Value); ok {
		return MakeFromReflect(rv.Field(i))
	}
	panic("invalid reflect.Value in Field()")
}

// SetField sets struct field at index i.
func (v Value) SetField(i int, val Value) {
	if rv, ok := v.obj.(reflect.Value); ok {
		field := rv.Field(i)
		fieldType := rv.Type().Field(i).Type

		// Handle slice conversion: []*T -> []interface{} for cyclic struct fields
		if fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Interface {
			if valRV, ok := val.ReflectValue(); ok && valRV.Kind() == reflect.Slice {
				if valRV.Type().Elem().Kind() != reflect.Interface {
					// Convert slice of concrete types to slice of interface{}
					convertedSlice := convertSliceToInterface(valRV, fieldType)
					field.Set(convertedSlice)
					return
				}
			}
		}

		field.Set(val.ToReflectValue(fieldType))
		return
	}
	panic("invalid reflect.Value in SetField()")
}

// Elem dereferences a pointer or returns the underlying value of interface.
func (v Value) Elem() Value {
	// Fast path: *int64 pointer (from native int slice)
	if ptr, ok := v.obj.(*int64); ok {
		return MakeInt(*ptr)
	}
	// Fast path: *Value pointer (from OpGlobal / OpAddr on value.Value locals / OpFree)
	if ptr, ok := v.obj.(*Value); ok {
		return *ptr
	}
	if rv, ok := v.obj.(reflect.Value); ok {
		// If the reflect.Value points to a value.Value struct, unwrap it directly.
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			if vp, ok2 := rv.Interface().(*Value); ok2 {
				return *vp
			}
		}
		return MakeFromReflect(rv.Elem())
	}
	panic("invalid reflect.Value in Elem()")
}

// SetElem sets the value pointed to by a pointer.
func (v Value) SetElem(val Value) {
	// Fast path: *int64 pointer (from native int slice OpIndexAddr)
	if ptr, ok := v.obj.(*int64); ok {
		*ptr = val.num
		return
	}
	// Fast path: *Value pointer (from OpFree for closure free vars)
	if ptr, ok := v.obj.(*Value); ok {
		*ptr = val
		return
	}
	if rv, ok := v.obj.(reflect.Value); ok {
		// Handle different reflect.Value kinds
		kind := rv.Kind()
		if kind == reflect.Ptr {
			// Handle pointer case
			elemType := rv.Type().Elem()
			if elemType.Kind() == reflect.Func {
				rv.Elem().Set(val.ToReflectValue(elemType))
				return
			}
			if elemType.Name() == "Value" && elemType.PkgPath() == "git.woa.com/youngjin/gig/model/value" {
				ptr := rv.Interface().(*Value)
				*ptr = val
				return
			}
			// If val contains a pointer type and elemType is not, unwrap it
			// But NOT if elemType is interface{} - interfaces can hold pointers
			if val.Kind() == KindReflect {
				if valRV, ok := val.obj.(reflect.Value); ok && valRV.Kind() == reflect.Ptr {
					// Special case: if elemType is interface{}, we can assign any value to it
					if elemType.Kind() == reflect.Interface {
						rv.Elem().Set(valRV)
						return
					}
					// val is a pointer, check if it points to elemType
					if valRV.Type().Elem() == elemType {
						rv.Elem().Set(valRV.Elem())
						return
					}
				}
			}
			targetRV := rv.Elem()
			if targetRV.CanSet() {
				// Native int slice → reflect conversion when assigning to *[]int
				if val.kind == KindSlice {
					if s, isInt := val.obj.([]int64); isInt && elemType.Kind() == reflect.Slice {
						// Convert []int64 to the target slice type (e.g. []int)
						target := reflect.MakeSlice(elemType, len(s), cap(s))
						for i, n := range s {
							target.Index(i).SetInt(n)
						}
						targetRV.Set(target)
						return
					}
				}
				valRV := val.ToReflectValue(elemType)
				if !valRV.Type().AssignableTo(elemType) {
					// Handle slice conversion: []*T -> []interface{} for cyclic struct fields
					if elemType.Kind() == reflect.Slice && valRV.Kind() == reflect.Slice {
						if elemType.Elem().Kind() == reflect.Interface && valRV.Type().Elem().Kind() != reflect.Interface {
							// Convert slice of concrete types to slice of interface{}
							convertedSlice := convertSliceToInterface(valRV, elemType)
							if convertedSlice.Type().AssignableTo(elemType) {
								targetRV.Set(convertedSlice)
								return
							}
						}
					}

					// Auto-unwrap pointer if elem matches
					if valRV.Kind() == reflect.Ptr && !valRV.IsNil() && valRV.Type().Elem().AssignableTo(elemType) {
						targetRV.Set(valRV.Elem())
						return
					}
				}
				targetRV.Set(valRV)
			}
			return
		}
		if kind == reflect.Interface {
			// For interface, just set the underlying value
			rv.Set(val.ToReflectValue(rv.Type()))
			return
		}
		if kind == reflect.Struct {
			// For struct values, we can't set elements - this shouldn't happen
			// but handle gracefully
			return
		}
	}
	panic("invalid reflect.Value in SetElem()")
}

// convertSliceToInterface converts a slice of concrete types to slice of interface{}.
// This is needed for self-referential structs where []*T becomes []interface{} due to cycle breaking.
func convertSliceToInterface(slice reflect.Value, targetType reflect.Type) reflect.Value {
	if slice.Kind() != reflect.Slice || targetType.Kind() != reflect.Slice {
		return slice
	}
	if targetType.Elem().Kind() != reflect.Interface {
		return slice
	}
	if slice.Type().Elem().Kind() == reflect.Interface {
		return slice // Already interface slice
	}

	// Create new slice with target type
	newSlice := reflect.MakeSlice(targetType, slice.Len(), slice.Cap())
	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)
		// Wrap each element in interface{}
		// For pointers, we need to convert to interface{} explicitly
		if elem.CanInterface() {
			newSlice.Index(i).Set(reflect.ValueOf(elem.Interface()))
		} else {
			// Unexported field - this shouldn't happen but handle gracefully
			newSlice.Index(i).Set(elem)
		}
	}
	return newSlice
}
