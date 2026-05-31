package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeCopy() {
	src := v.pop()
	dst := v.pop()
	// Native []byte copy fast path (handles copy([]byte, string))
	if dst.Kind() == value.KindBytes {
		if db, ok := dst.Bytes(); ok {
			if src.Kind() == value.KindString {
				ss := src.String()
				n := copy(db, ss)
				v.push(value.MakeInt(int64(n)))
				return
			}
			if sb, ok2 := src.Bytes(); ok2 {
				n := copy(db, sb)
				v.push(value.MakeInt(int64(n)))
				return
			}
		}
	}
	// Native int slice fast path
	if ds, ok := dst.IntSlice(); ok {
		if ss, ok2 := src.IntSlice(); ok2 {
			v.push(value.MakeInt(int64(copy(ds, ss))))
			return
		}
		// Cross-type: dst is native []int64, src is reflect slice (e.g. []int)
		if srcRV := v.mustReflectValue(src); srcRV.IsValid() && srcRV.Kind() == reflect.Slice {
			n := len(ds)
			if srcRV.Len() < n {
				n = srcRV.Len()
			}
			for i := 0; i < n; i++ {
				ds[i] = srcRV.Index(i).Int()
			}
			v.push(value.MakeInt(int64(n)))
			return
		}
	}
	// Copy slice
	if dstRV := v.mustReflectValue(dst); dstRV.IsValid() {
		if srcRV := v.mustReflectValue(src); srcRV.IsValid() {
			n := reflect.Copy(dstRV, srcRV)
			v.push(value.MakeInt(int64(n)))
		} else if ss, ok2 := src.IntSlice(); ok2 {
			// Cross-type: dst is reflect slice, src is native []int64
			n := dstRV.Len()
			if len(ss) < n {
				n = len(ss)
			}
			for i := 0; i < n; i++ {
				dstRV.Index(i).SetInt(ss[i])
			}
			v.push(value.MakeInt(int64(n)))
		} else {
			v.push(value.MakeInt(0))
		}
	} else {
		v.push(value.MakeInt(0))
	}
}

func (v *vm) executeDelete() {
	key := v.pop()
	m := v.pop()
	// In Go, delete on a nil map is a no-op
	if m.IsNil() {
		return
	}
	if rv := v.mustReflectValue(m); rv.IsValid() && rv.IsNil() {
		return
	}
	// For OpDelete, we want to delete the entry (deleteIfNil=true)
	m.SetMapIndexWithDelete(key, value.MakeNil(), true)
}
