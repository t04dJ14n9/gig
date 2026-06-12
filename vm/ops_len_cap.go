package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeLen() {
	obj := v.pop()
	switch obj.Kind() {
	case value.KindString:
		v.push(value.MakeInt(int64(len(obj.String()))))
	case value.KindBytes:
		if b, ok := obj.Bytes(); ok {
			v.push(value.MakeInt(int64(len(b))))
		} else {
			v.push(value.MakeInt(0))
		}
	case value.KindSlice:
		v.push(value.MakeInt(int64(obj.Len())))
	case value.KindArray, value.KindMap, value.KindChan:
		v.push(value.MakeInt(int64(obj.Len())))
	case value.KindInterface, value.KindReflect:
		// Handle both interface values and reflect-wrapped values
		rv := v.mustReflectValue(obj)
		if rv.IsValid() {
			kind := rv.Kind()
			if kind == reflect.Interface {
				// Unwrap interface to get underlying value
				if !rv.IsNil() {
					rv = rv.Elem()
					kind = rv.Kind()
				}
			}
			switch kind {
			case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
				v.push(value.MakeInt(int64(rv.Len())))
			default:
				v.push(value.MakeInt(0))
			}
		} else {
			v.push(value.MakeInt(0))
		}
	default:
		v.push(value.MakeInt(0))
	}
}

func (v *vm) executeCap() {
	obj := v.pop()
	switch obj.Kind() {
	case value.KindSlice, value.KindArray, value.KindChan:
		v.push(value.MakeInt(int64(obj.Cap())))
	case value.KindBytes:
		if b, ok := obj.Bytes(); ok {
			v.push(value.MakeInt(int64(cap(b))))
		} else {
			v.push(value.MakeInt(0))
		}
	case value.KindReflect:
		rv := v.mustReflectValue(obj)
		if rv.IsValid() {
			v.push(value.MakeInt(int64(rv.Cap())))
		} else {
			v.push(value.MakeInt(0))
		}
	default:
		v.push(value.MakeInt(0))
	}
}
