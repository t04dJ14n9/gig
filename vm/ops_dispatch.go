// ops_dispatch.go routes opcodes: executeOp dispatches to category-specific handlers.
package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
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
		bytecode.OpDeref, bytecode.OpSetDeref, bytecode.OpNew:
		return v.executeMemory(op, frame)

	// Closures & goroutines
	case bytecode.OpClosure, bytecode.OpGoCall, bytecode.OpGoCallExternal, bytecode.OpGoCallIndirect,
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
	t = t.Underlying()
	if basic, ok := t.(*types.Basic); ok {
		return basicKindMatchesValue(k, sz, basic.Kind())
	}
	return compositeKindMatchesValue(k, t)
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
	family := reflectNumericFamily(a.Kind())
	return family != reflectFamilyNone && family == reflectNumericFamily(b.Kind())
}

type reflectKindFamily uint8

// reflectKindFamily gives sameReflectKindFamily a stable vocabulary: two
// reflect kinds match only when both classify into the same non-empty family.
const (
	reflectFamilyNone reflectKindFamily = iota
	reflectFamilySignedInt
	reflectFamilyUnsignedInt
	reflectFamilyFloat
	reflectFamilyComplex
)

func reflectNumericFamily(k reflect.Kind) reflectKindFamily {
	switch {
	case isSignedReflectKind(k):
		return reflectFamilySignedInt
	case isUnsignedReflectKind(k):
		return reflectFamilyUnsignedInt
	case isFloatReflectKind(k):
		return reflectFamilyFloat
	case isComplexReflectKind(k):
		return reflectFamilyComplex
	default:
		return reflectFamilyNone
	}
}

func isSignedReflectKind(k reflect.Kind) bool {
	return k >= reflect.Int && k <= reflect.Int64
}

func isUnsignedReflectKind(k reflect.Kind) bool {
	return k >= reflect.Uint && k <= reflect.Uintptr
}

func isFloatReflectKind(k reflect.Kind) bool {
	return k == reflect.Float32 || k == reflect.Float64
}

func isComplexReflectKind(k reflect.Kind) bool {
	return k == reflect.Complex64 || k == reflect.Complex128
}
