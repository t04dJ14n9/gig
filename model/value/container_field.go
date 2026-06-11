package value

import "reflect"

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

		field.Set(ReflectValueForSet(val, fieldType))
		return
	}
	panic("invalid reflect.Value in SetField()")
}
