package gentool

import (
	"fmt"
	"go/types"
)

// --- Return value wrapping ---

// wrapReturn generates the Go expression to wrap a native Go value into a value.Value.
// It handles basic types, []byte, error, and general interface{} returns.
func wrapReturn(t types.Type, goExpr string) string {
	// Unwrap named/alias types to check their underlying type
	underlying := t.Underlying()

	// Basic type (or named type with basic underlying): use typed Make* constructors
	if bt, ok := underlying.(*types.Basic); ok {
		// If it's a named type wrapping a basic, cast to the basic type first
		if _, isNamed := t.(*types.Named); isNamed {
			basicName := bt.Name()
			return wrapBasicReturn(bt, fmt.Sprintf("%s(%s)", basicName, goExpr))
		}
		if _, isAlias := t.(*types.Alias); isAlias {
			basicName := bt.Name()
			return wrapBasicReturn(bt, fmt.Sprintf("%s(%s)", basicName, goExpr))
		}
		return wrapBasicReturn(bt, goExpr)
	}

	// []byte: use MakeBytes for zero-reflection
	if st, ok := underlying.(*types.Slice); ok {
		if bt, ok := st.Elem().Underlying().(*types.Basic); ok && bt.Kind() == types.Byte {
			return fmt.Sprintf("value.MakeBytes([]byte(%s))", goExpr)
		}
	}

	// error interface: use FromInterface (handles nil correctly)
	if named, ok := t.(*types.Named); ok {
		if named.Obj().Pkg() == nil && named.Obj().Name() == errorTypeName {
			return fmt.Sprintf("value.FromInterface(%s)", goExpr)
		}
	}

	return fmt.Sprintf("value.FromInterface(%s)", goExpr)
}

// wrapBasicReturn generates the Go expression to wrap a basic Go value into a value.Value.
// It uses the typed Make* constructors for zero-overhead wrapping.
func wrapBasicReturn(bt *types.Basic, goExpr string) string {
	info := bt.Info()
	kind := bt.Kind()

	switch {
	case kind == types.Bool:
		return fmt.Sprintf("value.MakeBool(%s)", goExpr)
	case info&types.IsInteger != 0 && info&types.IsUnsigned != 0:
		return wrapUnsignedBasicReturn(kind, goExpr)
	case info&types.IsInteger != 0:
		return wrapSignedBasicReturn(kind, goExpr)
	case info&types.IsFloat != 0:
		return wrapFloatBasicReturn(kind, goExpr)
	case info&types.IsComplex != 0:
		return fmt.Sprintf("value.FromInterface(%s)", goExpr)
	case info&types.IsString != 0:
		return fmt.Sprintf("value.MakeString(string(%s))", goExpr)
	default:
		return fmt.Sprintf("value.FromInterface(%s)", goExpr)
	}
}

func wrapUnsignedBasicReturn(kind types.BasicKind, goExpr string) string {
	// Preserve exact scalar width across the VM boundary. MakeUint uses SizePtr
	// while MakeUint64 uses Size64, so each sized unsigned type needs its own
	// constructor instead of a generic uint64 path.
	switch kind {
	case types.Uint8:
		return fmt.Sprintf("value.MakeUint8(%s)", goExpr)
	case types.Uint16:
		return fmt.Sprintf("value.MakeUint16(%s)", goExpr)
	case types.Uint32:
		return fmt.Sprintf("value.MakeUint32(%s)", goExpr)
	case types.Uint64:
		return fmt.Sprintf("value.MakeUint64(%s)", goExpr)
	default:
		return fmt.Sprintf("value.MakeUint(uint64(%s))", goExpr)
	}
}

func wrapSignedBasicReturn(kind types.BasicKind, goExpr string) string {
	// Preserve exact scalar width across the VM boundary. MakeInt uses SizePtr
	// while MakeInt64 uses Size64, so each sized signed type needs its own
	// constructor instead of a generic int64 path.
	switch kind {
	case types.Int8:
		return fmt.Sprintf("value.MakeInt8(%s)", goExpr)
	case types.Int16:
		return fmt.Sprintf("value.MakeInt16(%s)", goExpr)
	case types.Int32:
		return fmt.Sprintf("value.MakeInt32(%s)", goExpr)
	case types.Int64:
		return fmt.Sprintf("value.MakeInt64(%s)", goExpr)
	default:
		return fmt.Sprintf("value.MakeInt(int64(%s))", goExpr)
	}
}

func wrapFloatBasicReturn(kind types.BasicKind, goExpr string) string {
	if kind == types.Float32 {
		return fmt.Sprintf("value.MakeFloat32(%s)", goExpr)
	}
	return fmt.Sprintf("value.MakeFloat(float64(%s))", goExpr)
}
