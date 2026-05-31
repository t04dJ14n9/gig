package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeIndex() {
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
	case value.KindBytes:
		// Native []byte indexing returns uint8 as KindUint.
		if b, ok := container.Bytes(); ok {
			v.push(value.MakeUint8(b[int(key.RawInt())]))
		} else {
			v.push(value.MakeNil())
		}
	case value.KindReflect:
		v.executeReflectIndex(container, key)
	default:
		v.push(value.MakeNil())
	}
}

func (v *vm) executeReflectIndex(container, key value.Value) {
	rv := v.mustReflectValue(container)
	if !rv.IsValid() {
		v.push(value.MakeNil())
		return
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		idx := int(key.Int())
		v.push(value.MakeFromReflect(rv.Index(idx)))
	case reflect.Map:
		k := key.ToReflectValue(rv.Type().Key())
		elem := rv.MapIndex(k)
		if !elem.IsValid() {
			// Return zero value of element type, not nil (Go semantics).
			v.push(value.MakeFromReflect(reflect.Zero(rv.Type().Elem())))
		} else {
			v.push(value.MakeFromReflect(elem))
		}
	default:
		v.push(value.MakeNil())
	}
}

func (v *vm) executeIndexOk() {
	key := v.pop()
	container := v.pop()
	switch container.Kind() {
	case value.KindMap:
		v.executeMapIndexOk(container, key)
	case value.KindReflect:
		v.executeReflectIndexOk(container, key)
	default:
		v.pushCommaOk(value.MakeNil(), false)
	}
}

func (v *vm) executeMapIndexOk(container, key value.Value) {
	rv := v.mustReflectValue(container)
	if !rv.IsValid() {
		v.pushCommaOk(value.MakeNil(), false)
		return
	}
	k := key.ToReflectValue(rv.Type().Key())
	elem := rv.MapIndex(k)
	if !elem.IsValid() {
		v.pushCommaOk(value.MakeFromReflect(reflect.Zero(rv.Type().Elem())), false)
		return
	}
	v.pushCommaOk(value.MakeFromReflect(elem), true)
}

func (v *vm) executeReflectIndexOk(container, key value.Value) {
	rv := v.mustReflectValue(container)
	if !rv.IsValid() {
		v.pushCommaOk(value.MakeNil(), false)
		return
	}
	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key() == nil {
			v.pushCommaOk(value.MakeNil(), false)
			return
		}
		k := key.ToReflectValue(rv.Type().Key())
		if !k.IsValid() {
			v.pushCommaOk(value.MakeNil(), false)
			return
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
}

func (v *vm) executeSetIndex() {
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
		// For OpSetIndex, nil value means set to typed nil, not delete.
		container.SetMapIndexWithDelete(key, val, false)
	case value.KindReflect:
		rv := v.mustReflectValue(container)
		if !rv.IsValid() {
			return
		}
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			idx := int(key.Int())
			rv.Index(idx).Set(v.valueForReflectSet(val, rv.Type().Elem()))
		case reflect.Map:
			// For OpSetIndex, nil value means set to typed nil, not delete.
			container.SetMapIndexWithDelete(key, val, false)
		}
	}
}
