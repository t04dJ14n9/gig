package compiler

import (
	"go/constant"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// compileConst compiles a constant value.
func (c *compiler) compileConst(cnst *ssa.Const) {
	var v any
	switch t := cnst.Type().(type) {
	case *types.Basic:
		v = basicConstValue(t.Kind(), cnst.Value)
	case *types.Named, *types.Alias:
		// Named types and type aliases share the same compilation logic:
		// extract the underlying basic type for non-nil values.
		if cnst.Value != nil {
			if underlying, ok := t.Underlying().(*types.Basic); ok {
				v = basicConstValue(underlying.Kind(), cnst.Value)
			}
		} else {
			if rt := constTypeToReflect(t); rt != nil {
				v = reflect.Zero(rt)
			}
		}
	default:
		// For nil constants of reference types (map, slice, chan, func, pointer),
		// emit a typed nil via reflect.Zero so the VM preserves type information.
		// This is critical for nil map access (returns zero value of element type)
		// and nil-typed interface returns.
		// Also handles struct types (including empty struct) for zero values.
		if cnst.Value == nil {
			if rt := constTypeToReflect(cnst.Type()); rt != nil {
				v = reflect.Zero(rt)
			}
		}
	case *types.Struct:
		// Handle struct zero values (including empty struct literal {})
		if cnst.Value == nil {
			if rt := constTypeToReflect(t); rt != nil {
				v = reflect.Zero(rt)
			}
		}
	}

	idx := c.addConstant(v)
	c.emit(bytecode.OpConst, idx)
}

// basicConstValue extracts a Go value from a constant.Value based on the basic type kind.
// Returns nil for unsupported kinds.
func basicConstValue(kind types.BasicKind, val constant.Value) any { //nolint:gocyclo,cyclop
	if val == nil {
		return basicZeroValue(kind)
	}

	switch kind { //nolint:exhaustive
	case types.Bool, types.UntypedBool:
		return val.Kind() == constant.Bool && constant.BoolVal(val)
	case types.Int, types.UntypedInt, types.UntypedRune:
		i, exact := constant.Int64Val(val)
		if exact {
			return int(i)
		}
		return int(0)
	case types.Int8:
		i, _ := constant.Int64Val(val)
		return int8(i)
	case types.Int16:
		i, _ := constant.Int64Val(val)
		return int16(i)
	case types.Int32:
		i, _ := constant.Int64Val(val)
		return int32(i)
	case types.Int64:
		i, exact := constant.Int64Val(val)
		if exact {
			return i
		}
		return int64(0)
	case types.Uint:
		u, _ := constant.Uint64Val(val)
		return uint(u)
	case types.Uint8:
		u, _ := constant.Uint64Val(val)
		return uint8(u)
	case types.Uint16:
		u, _ := constant.Uint64Val(val)
		return uint16(u)
	case types.Uint32:
		u, _ := constant.Uint64Val(val)
		return uint32(u)
	case types.Uint64:
		u, _ := constant.Uint64Val(val)
		return u
	case types.Uintptr:
		u, _ := constant.Uint64Val(val)
		return uint64(u)
	case types.Float32:
		f, _ := constant.Float64Val(val)
		return float32(f)
	case types.Float64, types.UntypedFloat:
		f, _ := constant.Float64Val(val)
		return f
	case types.String, types.UntypedString:
		return constant.StringVal(val)
	case types.Complex64:
		re := constant.Real(val)
		im := constant.Imag(val)
		reVal, _ := constant.Float64Val(re)
		imVal, _ := constant.Float64Val(im)
		return complex(float32(reVal), float32(imVal))
	case types.Complex128, types.UntypedComplex:
		re := constant.Real(val)
		im := constant.Imag(val)
		reVal, _ := constant.Float64Val(re)
		imVal, _ := constant.Float64Val(im)
		return complex(reVal, imVal)
	default:
		return nil
	}
}

// basicZeroValue returns the zero value for a basic type kind, or nil if unsupported.
var basicZeroValues = map[types.BasicKind]any{
	types.Bool: false, types.UntypedBool: false,
	types.Int: int(0), types.UntypedInt: int(0), types.UntypedRune: int(0),
	types.Int8: int8(0), types.Int16: int16(0), types.Int32: int32(0), types.Int64: int64(0),
	types.Uint: uint(0), types.Uint8: uint8(0), types.Uint16: uint16(0),
	types.Uint32: uint32(0), types.Uint64: uint64(0), types.Uintptr: uint64(0),
	types.Float32: float32(0), types.Float64: 0.0, types.UntypedFloat: 0.0,
	types.String: "", types.UntypedString: "",
	types.Complex64: complex64(0), types.Complex128: complex128(0), types.UntypedComplex: complex128(0),
}

func basicZeroValue(kind types.BasicKind) any {
	return basicZeroValues[kind] // nil for unsupported kinds
}

// isEmptyStruct checks if a type is an empty struct (struct{}).
func isEmptyStruct(t types.Type) bool {
	// Named or Alias wrapper
	switch u := t.(type) {
	case *types.Named:
		t = u.Underlying()
	case *types.Alias:
		t = u.Underlying()
	}
	// Check the underlying Struct type
	if st, ok := t.(*types.Struct); ok {
		return st.NumFields() == 0
	}
	return false
}

func constTypeToReflect(t types.Type) reflect.Type {
	// Handle empty structs early (Named, Alias, and direct Struct types)
	if isEmptyStruct(t) {
		return emptyStructReflectType(t)
	}

	switch typ := t.Underlying().(type) {
	case *types.Basic:
		return bytecode.BasicKindToReflectType[typ.Kind()]
	case *types.Map:
		keyRT := constTypeToReflect(typ.Key())
		elemRT := constTypeToReflect(typ.Elem())
		if keyRT != nil && elemRT != nil {
			return reflect.MapOf(keyRT, elemRT)
		}
	case *types.Slice:
		elemRT := constTypeToReflect(typ.Elem())
		if elemRT != nil {
			return reflect.SliceOf(elemRT)
		}
	case *types.Pointer:
		elemRT := constTypeToReflect(typ.Elem())
		if elemRT != nil {
			return reflect.PointerTo(elemRT)
		}
	case *types.Chan:
		elemRT := constTypeToReflect(typ.Elem())
		if elemRT != nil {
			return reflect.ChanOf(chanDirection(typ), elemRT)
		}
	case *types.Interface:
		if typ.NumMethods() == 0 {
			return reflect.TypeFor[any]()
		}
	case *types.Signature:
		return buildFuncType(typ)
	}
	return nil
}

func emptyStructReflectType(t types.Type) reflect.Type {
	named, ok := t.(*types.Named)
	if !ok {
		return reflect.TypeFor[struct{}]()
	}
	obj := named.Obj()
	if obj == nil {
		return reflect.TypeFor[struct{}]()
	}
	typeName := obj.Name()
	qualName := "#" + typeName
	if pkg := obj.Pkg(); pkg != nil {
		qualName = "#" + pkg.Name() + "." + typeName
	}
	return reflect.StructOf([]reflect.StructField{{
		Name:    "gigType",
		Type:    reflect.TypeFor[struct{}](),
		PkgPath: "gig/internal",
		Tag:     reflect.StructTag(`gig:"` + qualName + `"`),
	}})
}

// chanDirection returns the reflect.ChanDir for a types.Chan.
func chanDirection(typ *types.Chan) reflect.ChanDir {
	switch typ.Dir() {
	case types.SendOnly:
		return reflect.SendDir
	case types.RecvOnly:
		return reflect.RecvDir
	default:
		return reflect.BothDir
	}
}

// buildFuncType builds a reflect.Type for a function signature.
func buildFuncType(sig *types.Signature) reflect.Type {
	params := make([]reflect.Type, sig.Params().Len())
	for i := 0; i < sig.Params().Len(); i++ {
		pt := constTypeToReflect(sig.Params().At(i).Type())
		if pt == nil {
			return nil
		}
		params[i] = pt
	}
	results := make([]reflect.Type, sig.Results().Len())
	for i := 0; i < sig.Results().Len(); i++ {
		rt := constTypeToReflect(sig.Results().At(i).Type())
		if rt == nil {
			return nil
		}
		results[i] = rt
	}
	return reflect.FuncOf(params, results, sig.Variadic())
}
