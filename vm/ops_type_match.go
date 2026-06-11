// ops_type_match.go contains type-switch matching helpers for OpAssert.
package vm

import (
	"go/types"

	"github.com/t04dJ14n9/gig/model/value"
)

func basicKindMatchesValue(k value.Kind, sz value.Size, basic types.BasicKind) bool {
	switch k {
	case value.KindInt:
		return signedIntegerSizeMatches(sz, basic)
	case value.KindUint:
		return unsignedIntegerSizeMatches(sz, basic)
	case value.KindFloat:
		return floatSizeMatches(sz, basic)
	case value.KindBool:
		return basic == types.Bool
	case value.KindString:
		return basic == types.String
	case value.KindComplex:
		return complexSizeMatches(sz, basic)
	default:
		return false
	}
}

func signedIntegerSizeMatches(sz value.Size, basic types.BasicKind) bool {
	switch basic {
	case types.Int:
		return sz == value.SizePtr
	case types.Int8:
		return sz == value.Size8
	case types.Int16:
		return sz == value.Size16
	case types.Int32:
		return sz == value.Size32
	case types.Int64:
		return sz == value.Size64
	default:
		return false
	}
}

func unsignedIntegerSizeMatches(sz value.Size, basic types.BasicKind) bool {
	switch basic {
	case types.Uint:
		return sz == value.SizePtr
	case types.Uint8:
		return sz == value.Size8
	case types.Uint16:
		return sz == value.Size16
	case types.Uint32:
		return sz == value.Size32
	case types.Uint64, types.Uintptr:
		return sz == value.Size64
	default:
		return false
	}
}

func floatSizeMatches(sz value.Size, basic types.BasicKind) bool {
	switch basic {
	case types.Float32:
		return sz == value.Size32
	case types.Float64:
		return sz == value.Size64
	default:
		return false
	}
}

func complexSizeMatches(sz value.Size, basic types.BasicKind) bool {
	switch basic {
	case types.Complex64:
		return sz == value.Size32
	case types.Complex128:
		return sz == value.Size64
	default:
		return false
	}
}

func compositeKindMatchesValue(k value.Kind, t types.Type) bool {
	switch k {
	case value.KindInt, value.KindUint, value.KindFloat, value.KindBool, value.KindString, value.KindComplex:
		return false
	case value.KindSlice:
		_, ok := t.(*types.Slice)
		return ok
	case value.KindMap:
		_, ok := t.(*types.Map)
		return ok
	case value.KindFunc:
		_, ok := t.(*types.Signature)
		return ok
	case value.KindBytes:
		return typeIsByteSlice(t)
	case value.KindNil:
		return false
	case value.KindInterface:
		_, ok := t.(*types.Interface)
		return ok
	default:
		// Reflect, pointer, struct, and other non-primitive values are normally
		// handled by the reflect path before this fallback.
		return true
	}
}

func typeIsByteSlice(t types.Type) bool {
	s, ok := t.(*types.Slice)
	if !ok {
		return false
	}
	basic, ok := s.Elem().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.Uint8 || basic.Kind() == types.Byte
}
