package value

import "reflect"

// ToReflectValue converts to reflect.Value.
// toReflectInt converts Value to reflect.Value for integer types.
func (v Value) toReflectInt(typ reflect.Type) reflect.Value {
	var intRV reflect.Value
	switch v.size {
	case Size8:
		intRV = reflect.ValueOf(int8(v.num))
	case Size16:
		intRV = reflect.ValueOf(int16(v.num))
	case Size32:
		intRV = reflect.ValueOf(int32(v.num))
	case Size64:
		intRV = reflect.ValueOf(v.num) // int64
	default:
		intRV = reflect.ValueOf(int(v.num)) // SizePtr / Size0 → int
	}
	if intRV.Type().ConvertibleTo(typ) {
		return intRV.Convert(typ)
	}
	return intRV
}

// toReflectUint converts Value to reflect.Value for unsigned integer types.
func (v Value) toReflectUint(typ reflect.Type) reflect.Value {
	var uintRV reflect.Value
	switch v.size {
	case Size8:
		uintRV = reflect.ValueOf(uint8(v.num))
	case Size16:
		uintRV = reflect.ValueOf(uint16(v.num))
	case Size32:
		uintRV = reflect.ValueOf(uint32(v.num))
	case Size64:
		uintRV = reflect.ValueOf(uint64(v.num))
	default:
		uintRV = reflect.ValueOf(uint(v.num)) // SizePtr / Size0 → uint
	}
	if uintRV.Type().ConvertibleTo(typ) {
		return uintRV.Convert(typ)
	}
	return uintRV
}

// toReflectFunc converts Value to reflect.Value for function types.
func (v Value) toReflectFunc(typ reflect.Type) reflect.Value {
	if typ.Kind() == reflect.Func {
		if ce, ok := v.obj.(ClosureExecutor); ok {
			numOut := typ.NumOut()
			outTypes := make([]reflect.Type, numOut)
			for i := 0; i < numOut; i++ {
				outTypes[i] = typ.Out(i)
			}
			fn := reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
				results := ce.Execute(args, outTypes)
				out := make([]reflect.Value, numOut)
				for i := 0; i < numOut; i++ {
					if i < len(results) && results[i].IsValid() {
						if results[i].Type().ConvertibleTo(outTypes[i]) {
							out[i] = results[i].Convert(outTypes[i])
						} else {
							out[i] = results[i]
						}
					} else {
						out[i] = reflect.Zero(outTypes[i])
					}
				}
				return out
			})
			return fn
		}
	}
	return reflect.ValueOf(v.obj)
}

// toReflectSlice converts Value to reflect.Value for slice types.
func (v Value) toReflectSlice(typ reflect.Type) reflect.Value {
	if s, ok := v.obj.([]int64); ok {
		// For []int64 storage (KindSlice), convert to []int when target is interface{}
		// or when target is []int. This ensures %T reports []int instead of []int64.
		if typ.Kind() == reflect.Interface && typ.NumMethod() == 0 {
			// Target is interface{} - convert to []int for Go compatibility
			result := make([]int, len(s))
			for i, n := range s {
				result[i] = int(n)
			}
			return reflect.ValueOf(result)
		}
		if typ.Kind() == reflect.Slice {
			target := reflect.MakeSlice(typ, len(s), cap(s))
			for i, n := range s {
				target.Index(i).SetInt(n)
			}
			return target
		}
	}
	if s, ok := v.obj.([]Value); ok && typ.Kind() == reflect.Slice {
		target := reflect.MakeSlice(typ, len(s), cap(s))
		elemType := typ.Elem()
		for i, elem := range s {
			target.Index(i).Set(ReflectValueForSet(elem, elemType))
		}
		return target
	}
	if rv, ok := v.obj.(reflect.Value); ok {
		return rv
	}
	return reflect.ValueOf(v.obj)
}

// toReflectReflect handles KindReflect values with special pointer-to-function conversions.
func (v Value) toReflectReflect(typ reflect.Type) reflect.Value {
	if rv, ok := v.obj.(reflect.Value); ok {
		// Handle *func(...) target type
		if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Func {
			if rv.Kind() == reflect.Ptr && !rv.IsNil() {
				if vp, ok2 := rv.Interface().(*Value); ok2 {
					funcRV := vp.ToReflectValue(typ.Elem())
					ptr := reflect.New(typ.Elem())
					ptr.Elem().Set(funcRV)
					return ptr
				}
			}
		}

		// Handle slice conversion: []*T -> []interface{}
		if typ.Kind() == reflect.Slice && rv.Kind() == reflect.Slice {
			if typ.Elem().Kind() == reflect.Interface && rv.Type().Elem().Kind() != reflect.Interface {
				return convertSliceToInterface(rv, typ)
			}
		}

		return rv
	}
	return reflect.ValueOf(v.obj)
}

func ReflectValueForSet(v Value, target reflect.Type) reflect.Value {
	rv := v.ToReflectValue(target)
	if rv.IsValid() && rv.Type().AssignableTo(target) {
		return rv
	}
	if isFmtStringerReflectType(target) {
		wrapped := FmtWrap(v)
		if wrapped != nil {
			wrappedRV := reflect.ValueOf(wrapped)
			if wrappedRV.IsValid() && wrappedRV.Type().AssignableTo(target) {
				return wrappedRV
			}
		}
	}
	return rv
}

func isFmtStringerReflectType(t reflect.Type) bool {
	return t != nil && t.Kind() == reflect.Interface && t.PkgPath() == "fmt" && t.Name() == "Stringer"
}

func (v Value) ToReflectValue(typ reflect.Type) reflect.Value {
	switch v.kind {
	case KindNil:
		return reflect.Zero(typ)
	case KindBool:
		return reflect.ValueOf(v.Bool())
	case KindInt:
		return v.toReflectInt(typ)
	case KindUint:
		return v.toReflectUint(typ)
	case KindFloat:
		return reflect.ValueOf(v.Float()).Convert(typ)
	case KindString:
		rv := reflect.ValueOf(v.obj.(string))
		if rv.Type() != typ {
			rv = rv.Convert(typ)
		}
		return rv
	case KindComplex:
		return reflect.ValueOf(v.obj.(complex128))
	case KindInterface:
		if dyn, ok := v.InterpretedInterface(); ok {
			return dyn.Value.ToReflectValue(typ)
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv
		}
		return reflect.ValueOf(v.obj)
	case KindFunc:
		return v.toReflectFunc(typ)
	case KindBytes:
		return reflect.ValueOf(v.obj.([]byte))
	case KindSlice:
		return v.toReflectSlice(typ)
	case KindReflect:
		return v.toReflectReflect(typ)
	default:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv
		}
		return reflect.ValueOf(v.obj)
	}
}

// ReflectValue returns the internal reflect.Value if stored.
func (v Value) ReflectValue() (reflect.Value, bool) {
	rv, ok := v.obj.(reflect.Value)
	return rv, ok
}
