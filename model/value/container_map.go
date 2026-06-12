package value

import "reflect"

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
