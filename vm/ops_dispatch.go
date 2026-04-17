// ops_dispatch.go routes opcodes: executeOp dispatches to category-specific handlers.
package vm

import (
	"go/types"
	"reflect"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
)

// executeOp executes a single bytecode instruction.
// It routes to category-specific handlers for each opcode group.
// Go runtime panics (nil deref, index out of range, etc.) are caught and
// converted to VM panics so that guest code's recover() can handle them.
// Note: hot-path opcodes (arithmetic, comparisons, stack ops, jumps, returns,
// calls) are inlined in run.go and never reach this dispatcher.
func (v *vm) executeOp(op bytecode.OpCode, frame *Frame) (retErr error) {
	// Catch Go runtime panics and convert them to VM panics.
	// This allows guest code's recover() to catch errors like:
	// - nil pointer dereference
	// - index out of range
	// - assignment to entry in nil map
	// - integer division by zero
	defer func() {
		if r := recover(); r != nil {
			v.panicking = true
			v.panicVal = value.FromInterface(r)
			retErr = nil
		}
	}()

	switch op {
	// Non-hot-path arithmetic & bitwise
	case bytecode.OpDiv, bytecode.OpMod,
		bytecode.OpNeg, bytecode.OpReal, bytecode.OpImag, bytecode.OpComplex,
		bytecode.OpAnd, bytecode.OpOr, bytecode.OpXor, bytecode.OpAndNot,
		bytecode.OpLsh, bytecode.OpRsh:
		return v.executeArithmetic(op, frame)

	// Memory: globals, free vars, fields, addresses, new
	case bytecode.OpGlobal, bytecode.OpSetGlobal,
		bytecode.OpFree, bytecode.OpSetFree,
		bytecode.OpField, bytecode.OpSetField, bytecode.OpAddr, bytecode.OpFieldAddr, bytecode.OpIndexAddr,
		bytecode.OpDeref, bytecode.OpSetDeref, bytecode.OpNew, bytecode.OpMake:
		return v.executeMemory(op, frame)

	// Closures & goroutines
	case bytecode.OpClosure, bytecode.OpGoCall, bytecode.OpGoCallIndirect,
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
	case bytecode.OpAssert, bytecode.OpConvert, bytecode.OpChangeType,
		bytecode.OpMakeInterface:
		return v.executeConvert(op, frame)

	// Channels, defer, panic, print, halt
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
	case value.KindReflect:
		if rv, ok := v.ReflectValue(); ok {
			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return rv.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return int64(rv.Uint())
			case reflect.Float32, reflect.Float64:
				return int64(rv.Float())
			}
		}
		return v.Int()
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
	case value.KindReflect:
		if rv, ok := v.ReflectValue(); ok {
			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return uint64(rv.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return rv.Uint()
			case reflect.Float32, reflect.Float64:
				return uint64(rv.Float())
			}
		}
		return v.Uint()
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
	case value.KindReflect:
		if rv, ok := v.ReflectValue(); ok {
			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return float64(rv.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return float64(rv.Uint())
			case reflect.Float32, reflect.Float64:
				return rv.Float()
			}
		}
		return v.Float()
	default:
		return v.Float()
	}
}

// kindMatchesType checks whether a value.Kind + value.Size matches a go/types.Type.
// This is used by OpAssert (type switch) to correctly match primitive values
// against target types. The size parameter enables exact type matching
// (e.g., int vs int64, complex64 vs complex128).
func kindMatchesType(k value.Kind, sz value.Size, t types.Type) bool {
	// Unwrap named types to get the underlying type
	t = t.Underlying()

	switch k {
	case value.KindInt:
		if basic, ok := t.(*types.Basic); ok {
			switch basic.Kind() {
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
			}
		}
		return false
	case value.KindUint:
		if basic, ok := t.(*types.Basic); ok {
			switch basic.Kind() {
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
			}
		}
		return false
	case value.KindFloat:
		if basic, ok := t.(*types.Basic); ok {
			switch basic.Kind() {
			case types.Float32:
				return sz == value.Size32
			case types.Float64:
				return sz == value.Size64
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
			case types.Complex64:
				return sz == value.Size32
			case types.Complex128:
				return sz == value.Size64
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
	case (ak >= reflect.Int && ak <= reflect.Int64) && (bk >= reflect.Int && bk <= reflect.Int64):
		return true
	case (ak >= reflect.Uint && ak <= reflect.Uintptr) && (bk >= reflect.Uint && bk <= reflect.Uintptr):
		return true
	case (ak == reflect.Float32 || ak == reflect.Float64) && (bk == reflect.Float32 || bk == reflect.Float64):
		return true
	case (ak == reflect.Complex64 || ak == reflect.Complex128) && (bk == reflect.Complex64 || bk == reflect.Complex128):
		return true
	default:
		return false
	}
}
