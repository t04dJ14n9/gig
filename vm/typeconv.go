package vm

import (
	"go/types"
	"reflect"
)

// typeToReflect converts a go/types.Type to reflect.Type.
// This is used for runtime type operations like allocations and type assertions.
// It handles basic types, slices, arrays, maps, channels, pointers, structs, and functions.
// It is safe for self-referencing struct types (e.g., type node struct { next *node }).
func typeToReflect(t types.Type) reflect.Type {
	return typeToReflectWithCache(t, make(map[types.Type]reflect.Type))
}

// TypeToReflect is the exported form of typeToReflect.
// It converts a go/types.Type to reflect.Type so the public API layer (gig.go)
// can map internal Value representations back to the exact Go types declared in
// the user's source code (e.g. int instead of int64).
func TypeToReflect(t types.Type) reflect.Type {
	return typeToReflect(t)
}

// typeToReflectWithCache is the internal recursive helper that carries a cache
// to detect and break cycles caused by self-referencing types.
// When a *types.Named is encountered a second time (cycle), the pointer field
// that caused the recursion is replaced with unsafe.Pointer, which has the same
// size and alignment as any Go pointer.
func typeToReflectWithCache(t types.Type, cache map[types.Type]reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}

	// Check cache first to break cycles
	if cached, ok := cache[t]; ok {
		return cached
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
		elem := typeToReflectWithCache(tt.Elem(), cache)
		if elem != nil {
			return reflect.SliceOf(elem)
		}
		return nil
	case *types.Array:
		elem := typeToReflectWithCache(tt.Elem(), cache)
		if elem != nil {
			return reflect.ArrayOf(int(tt.Len()), elem)
		}
		return nil
	case *types.Map:
		key := typeToReflectWithCache(tt.Key(), cache)
		val := typeToReflectWithCache(tt.Elem(), cache)
		if key != nil && val != nil {
			return reflect.MapOf(key, val)
		}
		return nil
	case *types.Chan:
		elem := typeToReflectWithCache(tt.Elem(), cache)
		if elem != nil {
			return reflect.ChanOf(reflect.BothDir, elem)
		}
		return nil
	case *types.Pointer:
		elem := typeToReflectWithCache(tt.Elem(), cache)
		if elem != nil {
			return reflect.PointerTo(elem)
		}
		// If elem is nil due to a cycle (self-referencing struct pointer),
		// use interface{} as a placeholder. The VM stores such values as
		// reflect.Value internally, and interface{} can hold any pointer value.
		var emptyIface any
		return reflect.TypeOf(&emptyIface).Elem()
	case *types.Interface:
		// Interface type — use the empty interface (any) type
		// For the VM, all interfaces are represented as interface{}
		var emptyIface any
		return reflect.TypeOf(&emptyIface).Elem()
	case *types.Named:
		// Mark this named type as being processed BEFORE recursing into the
		// underlying type. If we encounter it again via a pointer field, the
		// cache check at the top returns nil, and the *types.Pointer case
		// falls back to unsafe.Pointer.
		cache[tt] = nil
		result := typeToReflectWithCache(tt.Underlying(), cache)
		if result != nil {
			cache[tt] = result
		}
		return result
	case *types.Struct:
		// Build struct type dynamically using reflect
		numFields := tt.NumFields()
		fields := make([]reflect.StructField, 0, numFields)
		for i := 0; i < numFields; i++ {
			f := tt.Field(i)
			ft := typeToReflectWithCache(f.Type(), cache)
			if ft == nil {
				// Skip fields that could not be converted (shouldn't normally happen
				// unless there's a deep cycle on a non-pointer path).
				continue
			}
			sf := reflect.StructField{
				Name:      f.Name(),
				Type:      ft,
				Anonymous: f.Anonymous(),
			}
			// For unexported fields, we must set PkgPath (required by reflect.StructOf).
			// reflect.StructOf does NOT support anonymous unexported fields (it panics
			// with "is anonymous but has PkgPath set" or "is unexported but missing PkgPath").
			// Workaround: demote anonymous unexported fields to regular unexported fields.
			if !f.Exported() {
				if sf.Anonymous {
					sf.Anonymous = false
				}
				sf.PkgPath = f.Pkg().Path()
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
			pt := typeToReflectWithCache(params.At(i).Type(), cache)
			if pt == nil {
				return nil
			}
			paramTypes[i] = pt
		}
		// Get result types
		results := tt.Results()
		resultTypes := make([]reflect.Type, results.Len())
		for i := 0; i < results.Len(); i++ {
			rt := typeToReflectWithCache(results.At(i).Type(), cache)
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
