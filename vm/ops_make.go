package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeMakeSlice() {
	capVal := v.pop()
	lenVal := v.pop()
	typeIdxVal := v.pop()
	typ, ok := v.resolveType(typeIdxVal)
	if !ok {
		v.push(value.MakeNil())
		return
	}
	made := false
	if sliceType, ok := typ.(*types.Slice); ok {
		elemType := sliceType.Elem()
		// Native []int64 fast path for integer slice types
		if basic, isBasic := elemType.(*types.Basic); isBasic {
			switch basic.Kind() {
			case types.Int, types.Int64:
				v.push(value.MakeIntSlice(make([]int64, int(lenVal.Int()), int(capVal.Int()))))
				made = true
			}
		}
		// Function slice: use reflect path to create proper typed slice (e.g. []func() int)
		// instead of []value.Value, so it can be assigned to typed struct fields.
	}
	if !made {
		if rt := typeToReflect(typ, v.program); rt != nil {
			slice := reflect.MakeSlice(rt, int(lenVal.Int()), int(capVal.Int()))
			v.push(value.MakeFromReflect(slice))
		} else {
			v.push(value.MakeNil())
		}
	}
}

func (v *vm) executeMakeMap() {
	sizeVal := v.pop()
	typeIdxVal := v.pop()
	typ, ok := v.resolveType(typeIdxVal)
	if !ok {
		v.push(value.MakeNil())
		return
	}
	if rt := typeToReflect(typ, v.program); rt != nil {
		m := reflect.MakeMap(rt)
		_ = sizeVal // Size hint ignored for simplicity
		v.push(value.MakeFromReflect(m))
	} else {
		v.push(value.MakeNil())
	}
}

func (v *vm) executeMakeChan() {
	sizeVal := v.pop()
	typeIdxVal := v.pop()
	typ, ok := v.resolveType(typeIdxVal)
	if !ok {
		v.push(value.MakeNil())
		return
	}
	if rt := typeToReflect(typ, v.program); rt != nil {
		ch := reflect.MakeChan(rt, int(sizeVal.Int()))
		v.push(value.MakeFromReflect(ch))
	} else {
		v.push(value.MakeNil())
	}
}
