// iterator.go implements range iteration over slices, arrays, maps, and strings.
package vm

import (
	"reflect"
	"unicode/utf8"

	"github.com/t04dJ14n9/gig/model/value"
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
//
// For maps, reflect.MapRange().Next() is used, which correctly observes keys
// added during iteration (matching native Go range behavior where new keys
// "may or may not" be visited).
func (it *iterator) next() (key, val value.Value, ok bool) {
	switch it.collection.Kind() {
	case value.KindString:
		s := it.collection.String()
		if it.index >= len(s) {
			return noNext()
		}
		r, size := utf8.DecodeRuneInString(s[it.index:])
		key = value.MakeInt(int64(it.index))
		val = value.MakeInt32(r)
		it.index += size
		return key, val, true
	case value.KindSlice, value.KindArray:
		if it.index >= it.collection.Len() {
			return noNext()
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
				return noNext()
			}
		}
		if !it.mapIter.Next() {
			return noNext()
		}
		key = value.MakeFromReflect(it.mapIter.Key())
		val = value.MakeFromReflect(it.mapIter.Value())
		return key, val, true
	case value.KindChan:
		val, ok = it.collection.Recv()
		return value.MakeNil(), val, ok
	default:
		return it.nextReflect()
	}
}

func (it *iterator) nextReflectString(s string) (key, val value.Value, ok bool) {
	if it.index >= len(s) {
		return noNext()
	}
	r, size := utf8.DecodeRuneInString(s[it.index:])
	key = value.MakeInt(int64(it.index))
	val = value.MakeInt32(r)
	it.index += size
	return key, val, true
}

func (it *iterator) nextReflect() (key, val value.Value, ok bool) {
	rv, ok := it.collection.ReflectValue()
	if !ok {
		return noNext()
	}
	switch rv.Kind() {
	case reflect.String:
		return it.nextReflectString(rv.String())
	case reflect.Slice, reflect.Array:
		return it.nextReflectIndex(rv)
	case reflect.Map:
		return it.nextReflectMap(rv)
	case reflect.Chan:
		return nextReflectChannel(rv)
	default:
		return noNext()
	}
}

func (it *iterator) nextReflectIndex(rv reflect.Value) (key, val value.Value, ok bool) {
	if it.index >= rv.Len() {
		return noNext()
	}
	key = value.MakeInt(int64(it.index))
	val = value.MakeFromReflect(rv.Index(it.index))
	it.index++
	return key, val, true
}

func (it *iterator) nextReflectMap(rv reflect.Value) (key, val value.Value, ok bool) {
	if it.mapIter == nil {
		it.mapIter = rv.MapRange()
	}
	if !it.mapIter.Next() {
		return noNext()
	}
	key = value.MakeFromReflect(it.mapIter.Key())
	val = value.MakeFromReflect(it.mapIter.Value())
	return key, val, true
}

func nextReflectChannel(rv reflect.Value) (key, val value.Value, ok bool) {
	v, ok := rv.Recv()
	if !ok {
		return noNext()
	}
	return value.MakeNil(), value.MakeFromReflect(v), true
}

func noNext() (key, val value.Value, ok bool) {
	return value.MakeNil(), value.MakeNil(), false
}
