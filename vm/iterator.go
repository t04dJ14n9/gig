package vm

import (
	"reflect"
	"unicode/utf8"

	"github.com/t04dJ14n9/gig/value"
)

// iterator is a helper for range iteration over slices, arrays, maps, and strings.
type iterator struct {
	// collection is the value being iterated.
	collection value.Value

	// index is the current position.
	// For slices/arrays: element index (0-based).
	// For strings: current byte offset into the string.
	index int

	// mapIter is the map iterator (for maps).
	mapIter *reflect.MapIter
}

// next returns the next (key, value) pair and whether iteration should continue.
// It advances the iterator internally; the caller must not increment the index.
//
// For slices/arrays: key is the element index, value is the element.
// For strings: key is the byte offset, value is the rune (Unicode code point).
// For maps: key is the map key, value is the map value.
func (it *iterator) next() (key, val value.Value, ok bool) {
	switch it.collection.Kind() {
	case value.KindString:
		s := it.collection.String()
		if it.index >= len(s) {
			return value.MakeNil(), value.MakeNil(), false
		}
		r, size := utf8.DecodeRuneInString(s[it.index:])
		key = value.MakeInt(int64(it.index))
		val = value.MakeInt(int64(r))
		it.index += size
		return key, val, true
	case value.KindSlice, value.KindArray:
		if it.index >= it.collection.Len() {
			return value.MakeNil(), value.MakeNil(), false
		}
		key = value.MakeInt(int64(it.index))
		val = it.collection.Index(it.index)
		it.index++
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
			case reflect.String:
				s := rv.String()
				if it.index >= len(s) {
					return value.MakeNil(), value.MakeNil(), false
				}
				r, size := utf8.DecodeRuneInString(s[it.index:])
				key = value.MakeInt(int64(it.index))
				val = value.MakeInt(int64(r))
				it.index += size
				return key, val, true
			case reflect.Slice, reflect.Array:
				if it.index >= rv.Len() {
					return value.MakeNil(), value.MakeNil(), false
				}
				key = value.MakeInt(int64(it.index))
				val = value.MakeFromReflect(rv.Index(it.index))
				it.index++
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
