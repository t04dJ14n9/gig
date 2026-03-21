package vm

import (
"go/types"
"reflect"

"git.woa.com/youngjin/gig/bytecode"
"git.woa.com/youngjin/gig/value"
)

// executeOp executes a single bytecode instruction.
// It routes to category-specific handlers for each opcode group.
func (v *vm) executeOp(op bytecode.OpCode, frame *Frame) error {
	switch op {
	// Arithmetic & comparisons
	case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod,
		bytecode.OpNeg, bytecode.OpAnd, bytecode.OpOr, bytecode.OpXor, bytecode.OpAndNot,
		bytecode.OpLsh, bytecode.OpRsh,
		bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpLess, bytecode.OpLessEq,
		bytecode.OpGreater, bytecode.OpGreaterEq, bytecode.OpNot:
		return v.executeArithmetic(op, frame)

	// Memory & stack
	case bytecode.OpNop, bytecode.OpPop, bytecode.OpDup,
		bytecode.OpConst, bytecode.OpNil, bytecode.OpTrue, bytecode.OpFalse,
		bytecode.OpLocal, bytecode.OpSetLocal, bytecode.OpGlobal, bytecode.OpSetGlobal,
		bytecode.OpFree, bytecode.OpSetFree,
		bytecode.OpField, bytecode.OpSetField, bytecode.OpAddr, bytecode.OpFieldAddr, bytecode.OpIndexAddr,
		bytecode.OpDeref, bytecode.OpSetDeref,
		bytecode.OpNew, bytecode.OpMake:
		return v.executeMemory(op, frame)

	// Calls & closures
	case bytecode.OpCall, bytecode.OpCallExternal, bytecode.OpCallIndirect,
		bytecode.OpClosure, bytecode.OpGoCall, bytecode.OpGoCallIndirect,
		bytecode.OpPack, bytecode.OpUnpack:
		return v.executeCall(op, frame)

	// Containers
	case bytecode.OpMakeSlice, bytecode.OpMakeMap, bytecode.OpMakeChan,
		bytecode.OpIndex, bytecode.OpIndexOk, bytecode.OpSetIndex, bytecode.OpSlice,
		bytecode.OpRange, bytecode.OpRangeNext,
		bytecode.OpLen, bytecode.OpCap,
		bytecode.OpAppend, bytecode.OpCopy, bytecode.OpDelete:
		return v.executeContainer(op, frame)

	// Type conversions
	case bytecode.OpAssert, bytecode.OpConvert, bytecode.OpChangeType:
		return v.executeConvert(op, frame)

	// Control flow, channels, defer, panic, print, halt
	default:
		return v.executeControl(op, frame)
	}
}

// toInt64 extracts an int64 from a Value of any numeric kind.
func toInt64(v value.Value) int64 {
	switch v.Kind() {
	case value.KindInt:
		return v.Int()
	case value.KindUint:
		return int64(v.Uint())
	case value.KindFloat:
		return int64(v.Float())
	default:
		return v.Int()
	}
}

// toUint64 extracts a uint64 from a Value of any numeric kind.
func toUint64(v value.Value) uint64 {
	switch v.Kind() {
	case value.KindInt:
		return uint64(v.Int())
	case value.KindUint:
		return v.Uint()
	case value.KindFloat:
		return uint64(v.Float())
	default:
		return v.Uint()
	}
}

// toFloat64 extracts a float64 from a Value of any numeric kind.
func toFloat64(v value.Value) float64 {
	switch v.Kind() {
	case value.KindInt:
		return float64(v.Int())
	case value.KindUint:
		return float64(v.Uint())
	case value.KindFloat:
		return v.Float()
	default:
		return v.Float()
	}
}

// kindMatchesType checks whether a value.Kind matches a go/types.Type.
// This is used by OpAssert (type switch) to correctly match primitive values
// against target types, rather than blindly assuming success.
func kindMatchesType(k value.Kind, t types.Type) bool {
	// Unwrap named types to get the underlying type
	t = t.Underlying()

	switch k {
	case value.KindInt:
		if basic, ok := t.(*types.Basic); ok {
			switch basic.Kind() {
			case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
				return true
			}
		}
		return false
	case value.KindUint:
		if basic, ok := t.(*types.Basic); ok {
			switch basic.Kind() {
			case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
				return true
			}
		}
		return false
	case value.KindFloat:
		if basic, ok := t.(*types.Basic); ok {
			switch basic.Kind() {
			case types.Float32, types.Float64:
				return true
			}
		}
		return false
	case value.KindBool:
		if basic, ok := t.(*types.Basic); ok {
			return basic.Kind() == types.Bool
		}
		return false
	case value.KindString:
		if basic, ok := t.(*types.Basic); ok {
			return basic.Kind() == types.String
		}
		return false
	case value.KindComplex:
		if basic, ok := t.(*types.Basic); ok {
			switch basic.Kind() {
			case types.Complex64, types.Complex128:
				return true
			}
		}
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
		// []byte is a slice of uint8
		if s, ok := t.(*types.Slice); ok {
			if basic, ok2 := s.Elem().(*types.Basic); ok2 {
				return basic.Kind() == types.Uint8 || basic.Kind() == types.Byte
			}
		}
		return false
	case value.KindNil:
		return false
	case value.KindInterface:
		_, ok := t.(*types.Interface)
		return ok
	default:
		// For KindReflect, KindPointer, KindStruct, etc., fall through to true
		// (these should normally be handled by the reflect path above).
		return true
	}
}

// sameReflectKindFamily checks whether two reflect.Types belong to the same
// numeric kind family. This is needed because Gig internally stores all integers
// as int64 and all floats as float64. When these values are stored in interface{},
// the concrete reflect type may be int64 instead of int. A Go type switch
// "case int:" should still match the value even though its reflect type is int64.
//
// This function returns true only for numeric types within the same family:
//   - Signed integers: int, int8, int16, int32, int64
//   - Unsigned integers: uint, uint8, uint16, uint32, uint64, uintptr
//   - Floats: float32, float64
//   - Complex: complex64, complex128
func sameReflectKindFamily(a, b reflect.Type) bool {
	ak, bk := a.Kind(), b.Kind()
	switch {
	case isSignedInt(ak) && isSignedInt(bk):
		return true
	case isUnsignedInt(ak) && isUnsignedInt(bk):
		return true
	case isFloat(ak) && isFloat(bk):
		return true
	case isComplex(ak) && isComplex(bk):
		return true
	default:
		return false
	}
}

func isSignedInt(k reflect.Kind) bool {
	return k >= reflect.Int && k <= reflect.Int64
}

func isUnsignedInt(k reflect.Kind) bool {
	return k >= reflect.Uint && k <= reflect.Uintptr
}

func isFloat(k reflect.Kind) bool {
	return k == reflect.Float32 || k == reflect.Float64
}

func isComplex(k reflect.Kind) bool {
	return k == reflect.Complex64 || k == reflect.Complex128
}
