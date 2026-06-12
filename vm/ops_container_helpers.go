package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

// pushCommaOk pushes a (value, ok) tuple onto the operand stack.
// Used by OpIndexOk, OpRecvOk, OpTypeAssert, etc.
func (v *vm) pushCommaOk(val value.Value, ok bool) {
	tuple := []value.Value{val, value.MakeBool(ok)}
	v.push(value.FromInterface(tuple))
}

// resolveType resolves a type index from a popped value into a types.Type.
// Returns the type and true if found, or nil and false otherwise.
func (v *vm) resolveType(typeIdxVal value.Value) (types.Type, bool) {
	typeIdx := uint16(typeIdxVal.Int())
	if int(typeIdx) < len(v.program.Types) {
		return v.program.Types[typeIdx], true
	}
	return nil, false
}

// mustReflectValue extracts a reflect.Value from a value.Value or returns an
// invalid reflect.Value if the value doesn't contain a reflect.Value.
// This helper reduces repetitive if rv, ok := val.ReflectValue() patterns.
func (v *vm) mustReflectValue(val value.Value) reflect.Value {
	if rv, ok := val.ReflectValue(); ok {
		return rv
	}
	return reflect.Value{}
}

func (v *vm) valueForReflectSet(val value.Value, target reflect.Type) reflect.Value {
	return value.ReflectValueForSet(val, target)
}
