package value

import "reflect"

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
				rv.Elem().Set(ReflectValueForSet(val, elemType))
				return
			}
			if elemType.Name() == "Value" && elemType.PkgPath() == "github.com/t04dJ14n9/gig/model/value" {
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
				// Native int slice -> reflect conversion when assigning to *[]int
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
				valRV := ReflectValueForSet(val, elemType)
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
			rv.Set(ReflectValueForSet(val, rv.Type()))
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
