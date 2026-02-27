package vm

import (
	"go/types"
	"reflect"
)

// typeToReflect converts a go/types.Type to reflect.Type.
// This is used for runtime type operations like allocations and type assertions.
// It handles basic types, slices, arrays, maps, channels, pointers, structs, and functions.
func typeToReflect(t types.Type) reflect.Type {
	if t == nil {
		return nil
	}

	switch tt := t.(type) {
	case *types.Basic:
		switch tt.Kind() {
		case types.Bool:
			return reflect.TypeOf(false)
		case types.Int:
			return reflect.TypeOf(int(0))
		case types.Int8:
			return reflect.TypeOf(int8(0))
		case types.Int16:
			return reflect.TypeOf(int16(0))
		case types.Int32:
			return reflect.TypeOf(int32(0))
		case types.Int64:
			return reflect.TypeOf(int64(0))
		case types.Uint:
			return reflect.TypeOf(uint(0))
		case types.Uint8:
			return reflect.TypeOf(uint8(0))
		case types.Uint16:
			return reflect.TypeOf(uint16(0))
		case types.Uint32:
			return reflect.TypeOf(uint32(0))
		case types.Uint64:
			return reflect.TypeOf(uint64(0))
		case types.Uintptr:
			return reflect.TypeOf(uintptr(0))
		case types.Float32:
			return reflect.TypeOf(float32(0))
		case types.Float64:
			return reflect.TypeOf(float64(0))
		case types.Complex64:
			return reflect.TypeOf(complex64(0))
		case types.Complex128:
			return reflect.TypeOf(complex128(0))
		case types.String:
			return reflect.TypeOf("")
		default:
			return nil
		}
	case *types.Slice:
		elem := typeToReflect(tt.Elem())
		if elem != nil {
			return reflect.SliceOf(elem)
		}
		return nil
	case *types.Array:
		elem := typeToReflect(tt.Elem())
		if elem != nil {
			return reflect.ArrayOf(int(tt.Len()), elem)
		}
		return nil
	case *types.Map:
		key := typeToReflect(tt.Key())
		val := typeToReflect(tt.Elem())
		if key != nil && val != nil {
			return reflect.MapOf(key, val)
		}
		return nil
	case *types.Chan:
		elem := typeToReflect(tt.Elem())
		if elem != nil {
			return reflect.ChanOf(reflect.BothDir, elem)
		}
		return nil
	case *types.Pointer:
		elem := typeToReflect(tt.Elem())
		if elem != nil {
			return reflect.PointerTo(elem)
		}
		return nil
	case *types.Interface:
		// Interface type — use the empty interface (any) type
		// For the VM, all interfaces are represented as interface{}
		var emptyIface any
		return reflect.TypeOf(&emptyIface).Elem()
	case *types.Named:
		// For named types, try to get the underlying type
		return typeToReflect(tt.Underlying())
	case *types.Struct:
		// Build struct type dynamically using reflect
		numFields := tt.NumFields()
		fields := make([]reflect.StructField, 0, numFields)
		for i := 0; i < numFields; i++ {
			f := tt.Field(i)
			ft := typeToReflect(f.Type())
			if ft == nil {
				return nil
			}
			sf := reflect.StructField{
				Name:      f.Name(),
				Type:      ft,
				Anonymous: f.Anonymous(),
			}
			// For unexported fields, we must set PkgPath
			// Check if the field is exported (starts with uppercase)
			if len(f.Name()) > 0 && f.Name()[0] >= 'a' && f.Name()[0] <= 'z' {
				// Unexported field - need to use package path
				// Use empty string for anonymous unexported, or the package path
				if !f.Anonymous() {
					sf.PkgPath = f.Pkg().Path()
				}
			}
			if tag := tt.Tag(i); tag != "" {
				sf.Tag = reflect.StructTag(tag)
			}
			fields = append(fields, sf)
		}
		if len(fields) == 0 {
			return nil
		}
		return reflect.StructOf(fields)
	case *types.Signature:
		// Function type - need to build the function type dynamically
		// Get parameter types
		params := tt.Params()
		paramTypes := make([]reflect.Type, params.Len())
		for i := 0; i < params.Len(); i++ {
			pt := typeToReflect(params.At(i).Type())
			if pt == nil {
				return nil
			}
			paramTypes[i] = pt
		}
		// Get result types
		results := tt.Results()
		resultTypes := make([]reflect.Type, results.Len())
		for i := 0; i < results.Len(); i++ {
			rt := typeToReflect(results.At(i).Type())
			if rt == nil {
				return nil
			}
			resultTypes[i] = rt
		}
		// Create function type using reflect.FuncOf
		return reflect.FuncOf(paramTypes, resultTypes, tt.Variadic())
	default:
		return nil
	}
}
