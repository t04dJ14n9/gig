package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeNew(frame *Frame) {
	typeIdx := frame.readUint16()
	if int(typeIdx) >= len(v.program.Types) {
		v.push(value.MakeNil())
		return
	}
	v.push(v.newValueForType(v.program.Types[typeIdx]))
}

func (v *vm) newValueForType(typ types.Type) value.Value {
	// Function types store closures in a *value.Value cell. Other types use the
	// reflect type generated for the compiled program's type universe.
	if _, isSig := typ.(*types.Signature); isSig {
		var nilVal value.Value
		return value.MakeFromReflect(reflect.ValueOf(&nilVal))
	}
	if rt := typeToReflect(typ, v.program); rt != nil {
		return value.MakeFromReflect(reflect.New(rt))
	}
	return value.MakeNil()
}
