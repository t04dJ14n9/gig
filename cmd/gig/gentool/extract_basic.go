package gentool

import (
	"fmt"
	"go/types"
)

// extractBasic generates the Go expression to extract a basic typed value from a value.Value.
// Returns "" for unsupported types such as unsafe.Pointer.
func extractBasic(bt *types.Basic, valExpr string) string {
	info := bt.Info()
	kind := bt.Kind()

	switch {
	case kind == types.Bool:
		return fmt.Sprintf("%s.Bool()", valExpr)
	case info&types.IsInteger != 0 && info&types.IsUnsigned != 0:
		return extractUnsignedBasic(kind, valExpr)
	case info&types.IsInteger != 0:
		return extractSignedBasic(kind, valExpr)
	case info&types.IsFloat != 0:
		return extractFloatBasic(kind, valExpr)
	case info&types.IsComplex != 0:
		if kind == types.Complex64 {
			return fmt.Sprintf("complex64(%s.Complex())", valExpr)
		}
		return fmt.Sprintf("%s.Complex()", valExpr)
	case info&types.IsString != 0:
		return fmt.Sprintf("%s.String()", valExpr)
	case kind == types.UnsafePointer:
		return ""
	default:
		return fmt.Sprintf("%s.Interface()", valExpr)
	}
}

func extractUnsignedBasic(kind types.BasicKind, valExpr string) string {
	switch kind {
	case types.Uint8:
		return fmt.Sprintf("byte(%s.Uint())", valExpr)
	case types.Uint16:
		return fmt.Sprintf("uint16(%s.Uint())", valExpr)
	case types.Uint32:
		return fmt.Sprintf("uint32(%s.Uint())", valExpr)
	case types.Uint64:
		return fmt.Sprintf("%s.Uint()", valExpr)
	case types.Uint:
		return fmt.Sprintf("uint(%s.Uint())", valExpr)
	case types.Uintptr:
		return fmt.Sprintf("uintptr(%s.Uint())", valExpr)
	default:
		return fmt.Sprintf("uint(%s.Uint())", valExpr)
	}
}

func extractSignedBasic(kind types.BasicKind, valExpr string) string {
	switch kind {
	case types.Int8:
		return fmt.Sprintf("int8(%s.Int())", valExpr)
	case types.Int16:
		return fmt.Sprintf("int16(%s.Int())", valExpr)
	case types.Int32:
		return fmt.Sprintf("int32(%s.Int())", valExpr)
	case types.Int64:
		return fmt.Sprintf("%s.Int()", valExpr)
	case types.Int:
		return fmt.Sprintf("int(%s.Int())", valExpr)
	default:
		return fmt.Sprintf("int(%s.Int())", valExpr)
	}
}

func extractFloatBasic(kind types.BasicKind, valExpr string) string {
	switch kind {
	case types.Float32:
		return fmt.Sprintf("float32(%s.Float())", valExpr)
	case types.Float64:
		return fmt.Sprintf("%s.Float()", valExpr)
	default:
		return fmt.Sprintf("%s.Float()", valExpr)
	}
}

// extractSlice generates the Go expression to extract a slice from a value.Value.
func extractSlice(st *types.Slice, valExpr string, pkgRef string) string {
	if bt, ok := st.Elem().Underlying().(*types.Basic); ok {
		switch bt.Kind() {
		case types.Byte:
			return fmt.Sprintf("func() []byte { if b, ok := (%s).Bytes(); ok { return b }; v := (%s).Interface(); if v == nil { return nil }; return v.([]byte) }()", valExpr, valExpr)
		case types.String:
			return fmt.Sprintf("%s.Interface().([]string)", valExpr)
		default:
			elemName := resolveTypeName(st.Elem(), pkgRef)
			if elemName == "" {
				return fmt.Sprintf("%s.Interface()", valExpr)
			}
			return fmt.Sprintf("%s.Interface().([]%s)", valExpr, elemName)
		}
	}
	return fmt.Sprintf("%s.Interface()", valExpr)
}
