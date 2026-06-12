package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeCopy() {
	src := v.pop()
	dst := v.pop()
	v.push(value.MakeInt(int64(v.copyValue(dst, src))))
}

func (v *vm) copyValue(dst, src value.Value) int {
	if n, ok := copyByteValue(dst, src); ok {
		return n
	}
	if n, ok := v.copyNativeIntValue(dst, src); ok {
		return n
	}
	return v.copyReflectValue(dst, src)
}

// copyByteValue handles the special Go builtin case copy([]byte, string) plus
// the native []byte-to-[]byte fast path without falling into reflect.
func copyByteValue(dst, src value.Value) (int, bool) {
	if dst.Kind() == value.KindBytes {
		if db, ok := dst.Bytes(); ok {
			if src.Kind() == value.KindString {
				return copy(db, src.String()), true
			}
			if sb, ok2 := src.Bytes(); ok2 {
				return copy(db, sb), true
			}
		}
		return 0, true
	}
	return 0, false
}

// copyNativeIntValue keeps Gig's []int fast path on native []int64 storage and
// handles reflected []int sources produced by mixed native/reflect operations.
func (v *vm) copyNativeIntValue(dst, src value.Value) (int, bool) {
	if ds, ok := dst.IntSlice(); ok {
		if ss, ok2 := src.IntSlice(); ok2 {
			return copy(ds, ss), true
		}
		if n, ok := v.copyReflectToNativeIntSlice(ds, src); ok {
			return n, true
		}
		return 0, true
	}
	return 0, false
}

func (v *vm) copyReflectToNativeIntSlice(dst []int64, src value.Value) (int, bool) {
	srcRV := v.mustReflectValue(src)
	if !srcRV.IsValid() || srcRV.Kind() != reflect.Slice {
		return 0, false
	}
	n := min(len(dst), srcRV.Len())
	for i := 0; i < n; i++ {
		dst[i] = srcRV.Index(i).Int()
	}
	return n, true
}

func (v *vm) copyReflectValue(dst, src value.Value) int {
	dstRV := v.mustReflectValue(dst)
	if !dstRV.IsValid() {
		return 0
	}
	srcRV := v.mustReflectValue(src)
	if srcRV.IsValid() {
		return reflect.Copy(dstRV, srcRV)
	}
	if ss, ok := src.IntSlice(); ok {
		return copyNativeIntToReflectSlice(dstRV, ss)
	}
	return 0
}

func copyNativeIntToReflectSlice(dstRV reflect.Value, src []int64) int {
	n := min(dstRV.Len(), len(src))
	for i := 0; i < n; i++ {
		dstRV.Index(i).SetInt(src[i])
	}
	return n
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
