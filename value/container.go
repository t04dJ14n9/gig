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
	default:
		panic(fmt.Sprintf("cannot take cap of %v", v.kind))
	}
}

// Index returns element at index i for slice, array, or string.
func (v Value) Index(i int) Value {
	switch v.kind {
	case KindString:
		// s[i] returns a byte (uint8), not a string
		return MakeUint(uint64(v.obj.(string)[i]))
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
		// For function element types, store the Value directly (closures are *Closure)
		if elemType.Kind() == reflect.Func {
			rv.Index(i).Set(reflect.ValueOf(val))
			return
		}
		// For []value.Value slices (used for function slices)
		if elemType == reflect.TypeOf(Value{}) {
			rv.Index(i).Set(reflect.ValueOf(val))
			return
		}
		rv.Index(i).Set(val.ToReflectValue(elemType))
		return
	}
	// Handle native []value.Value slice
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

// SetMapIndex sets value at key k for map.
func (v Value) SetMapIndex(k, val Value) {
	if rv, ok := v.obj.(reflect.Value); ok {
		key := k.ToReflectValue(rv.Type().Key())
		if val.IsNil() {
			rv.SetMapIndex(key, reflect.Value{})
		} else {
			rv.SetMapIndex(key, val.ToReflectValue(rv.Type().Elem()))
		}
		return
	}
	panic("invalid reflect.Value in SetMapIndex()")
}

// MapIter iterates over a map.
func (v Value) MapIter(f func(key, val Value) bool) {
	if rv, ok := v.obj.(reflect.Value); ok {
		iter := rv.MapRange()
		for iter.Next() {
			key := MakeFromReflect(iter.Key())
			val := MakeFromReflect(iter.Value())
			if !f(key, val) {
				break
			}
		}
		return
	}
	panic("invalid reflect.Value in MapIter()")
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
		rv.Field(i).Set(val.ToReflectValue(rv.Type().Field(i).Type))
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
	// Fast path: *Value pointer (from OpGlobal / OpAddr on value.Value locals)
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
	if rv, ok := v.obj.(reflect.Value); ok {
		// Handle different reflect.Value kinds
		kind := rv.Kind()
		if kind == reflect.Ptr {
			// Handle pointer case
			elemType := rv.Type().Elem()
			if elemType.Kind() == reflect.Func {
				rv.Elem().Set(reflect.ValueOf(val))
				return
			}
			if elemType.Name() == "Value" && elemType.PkgPath() == "github.com/t04dJ14n9/gig/value" {
				ptr := rv.Interface().(*Value)
				*ptr = val
				return
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
				targetRV.Set(val.ToReflectValue(elemType))
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

// Pointer returns the underlying pointer value.
func (v Value) Pointer() uintptr {
	if rv, ok := v.obj.(reflect.Value); ok {
		return rv.Pointer()
	}
	return 0
}

// Package packs multiple values into a slice.
func Package(vals ...Value) []Value {
	return vals
}

// Unpackage unpacks a slice of values.
func Unpackage(vals []Value) []Value {
	return vals
}
