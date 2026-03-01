package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/value"
)

// iterator is a helper for range iteration over slices, arrays, maps, and strings.
type iterator struct {
	// collection is the value being iterated.
	collection value.Value

	// index is the current position (for slices/arrays/strings).
	index int

	// mapIter is the map iterator (for maps).
	mapIter *reflect.MapIter
}

// next advances the iterator and returns the next key, value, and whether there are more elements.
// For slices/arrays/strings, key is the index.
// For maps, key is the map key.
func (it *iterator) next() (key, val value.Value, ok bool) {
	switch it.collection.Kind() {
	case value.KindSlice, value.KindArray, value.KindString:
		if it.index >= it.collection.Len() {
			return value.MakeNil(), value.MakeNil(), false
		}
		key = value.MakeInt(int64(it.index))
		val = it.collection.Index(it.index)
		return key, val, true
	case value.KindMap:
		if it.mapIter == nil {
			if rv, isValid := it.collection.ReflectValue(); isValid {
				it.mapIter = rv.MapRange()
			} else {
				return value.MakeNil(), value.MakeNil(), false
			}
		}
		if !it.mapIter.Next() {
			return value.MakeNil(), value.MakeNil(), false
		}
		key = value.MakeFromReflect(it.mapIter.Key())
		val = value.MakeFromReflect(it.mapIter.Value())
		return key, val, true
	case value.KindChan:
		val, ok = it.collection.Recv()
		return value.MakeNil(), val, ok
	default:
		// Try to use reflect for other types
		if rv, isValid := it.collection.ReflectValue(); isValid {
			switch rv.Kind() {
			case reflect.Slice, reflect.Array, reflect.String:
				if it.index >= rv.Len() {
					return value.MakeNil(), value.MakeNil(), false
				}
				key = value.MakeInt(int64(it.index))
				val = value.MakeFromReflect(rv.Index(it.index))
				return key, val, true
			case reflect.Map:
				if it.mapIter == nil {
					it.mapIter = rv.MapRange()
				}
				if !it.mapIter.Next() {
					return value.MakeNil(), value.MakeNil(), false
				}
				key = value.MakeFromReflect(it.mapIter.Key())
				val = value.MakeFromReflect(it.mapIter.Value())
				return key, val, true
			case reflect.Chan:
				v, ok := rv.Recv()
				if !ok {
					return value.MakeNil(), value.MakeNil(), false
				}
				return value.MakeNil(), value.MakeFromReflect(v), true
			}
		}
		return value.MakeNil(), value.MakeNil(), false
	}
}
