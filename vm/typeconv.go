package vm

import (
	"go/types"
	"reflect"
	"strings"

	"git.woa.com/youngjin/gig/bytecode"
)

// typeToReflect converts a go/types.Type to reflect.Type using the program-level
// cache to ensure the same types.Type always maps to the same reflect.Type.
// This prevents reflect.StructOf from returning different reflect.Type objects
// across multiple VM executions, which would cause "reflect.Set: value not assignable" panics.
func typeToReflect(t types.Type, prog *bytecode.Program) reflect.Type {
	if t == nil {
		return nil
	}
	// Check program-level cache first (fast path, lock-free read)
	if rt, ok := prog.CachedReflectType(t); ok {
		return rt
	}
	// Compute with a local cycle-detection cache
	localCache := make(map[types.Type]reflect.Type)
	rt := typeToReflectWithCache(t, localCache, "", prog)
	if rt != nil {
		// Store in program-level cache (uses LoadOrStore for thread safety)
		rt = prog.CacheReflectType(t, rt)
	}
	return rt
}

// typeToReflectWithCache is the internal recursive helper that carries a local cache
// to detect and break cycles caused by self-referencing types.
// When a *types.Named is encountered a second time (cycle), the pointer field
// that caused the recursion is replaced with unsafe.Pointer, which has the same
// size and alignment as any Go pointer.
// The uniqueSuffix parameter is used to create unique reflect.Types for named structs
// to prevent reflect.StructOf from deduplicating different types with same field layout.
//
// NOTE: This function does NOT use the program-level cache internally because the same
// types.Type (e.g., *types.Struct for struct{v int}) may be reached through different
// named types with different uniqueSuffix values. Caching at this level would cause
// suffix-insensitive collisions. Program-level caching is done only at the top-level
// typeToReflect entry point, which caches the final result keyed by the original type.
func typeToReflectWithCache(t types.Type, cache map[types.Type]reflect.Type, uniqueSuffix string, prog *bytecode.Program) reflect.Type {
	if t == nil {
		return nil
	}

	// Check local cache first (for cycle detection)
	if cached, ok := cache[t]; ok {
		return cached
	}

	result := typeToReflectInner(t, cache, uniqueSuffix, prog)
	if result != nil {
		cache[t] = result
	}

	return result
}

// typeToReflectInner does the actual conversion without caching logic.
func typeToReflectInner(t types.Type, cache map[types.Type]reflect.Type, uniqueSuffix string, prog *bytecode.Program) reflect.Type {
	switch tt := t.(type) {
	case *types.Basic:
		switch tt.Kind() {
		case types.Bool:
			return reflect.TypeFor[bool]()
		case types.Int:
			return reflect.TypeFor[int]()
		case types.Int8:
			return reflect.TypeFor[int8]()
		case types.Int16:
			return reflect.TypeFor[int16]()
		case types.Int32:
			return reflect.TypeFor[int32]()
		case types.Int64:
			return reflect.TypeFor[int64]()
		case types.Uint:
			return reflect.TypeFor[uint]()
		case types.Uint8:
			return reflect.TypeFor[uint8]()
		case types.Uint16:
			return reflect.TypeFor[uint16]()
		case types.Uint32:
			return reflect.TypeFor[uint32]()
		case types.Uint64:
			return reflect.TypeFor[uint64]()
		case types.Uintptr:
			return reflect.TypeFor[uintptr]()
		case types.Float32:
			return reflect.TypeFor[float32]()
		case types.Float64:
			return reflect.TypeFor[float64]()
		case types.Complex64:
			return reflect.TypeFor[complex64]()
		case types.Complex128:
			return reflect.TypeFor[complex128]()
		case types.String:
			return reflect.TypeFor[string]()
		default:
			return nil
		}
	case *types.Slice:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog)
		if elem != nil {
			return reflect.SliceOf(elem)
		}
		return nil
	case *types.Array:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog)
		if elem != nil {
			return reflect.ArrayOf(int(tt.Len()), elem)
		}
		return nil
	case *types.Map:
		key := typeToReflectWithCache(tt.Key(), cache, "", prog)
		val := typeToReflectWithCache(tt.Elem(), cache, "", prog)
		if key != nil && val != nil {
			return reflect.MapOf(key, val)
		}
		return nil
	case *types.Chan:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog)
		if elem != nil {
			return reflect.ChanOf(reflect.BothDir, elem)
		}
		return nil
	case *types.Pointer:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog)
		if elem != nil {
			return reflect.PointerTo(elem)
		}
		// If elem is nil due to a cycle (self-referencing struct pointer),
		// use interface{} as a placeholder. The VM stores such values as
		// reflect.Value internally, and interface{} can hold any pointer value.
		return reflect.TypeFor[any]()
	case *types.Interface:
		// Interface type — use the empty interface (any) type
		// For the VM, all interfaces are represented as interface{}
		return reflect.TypeFor[any]()
	case *types.Named:
		// Check if this is a registered external type (e.g., bytes.Buffer, strings.Builder).
		// If so, use the real reflect.Type instead of synthesizing a struct type.
		if prog != nil && prog.Lookup != nil {
			if rt, ok := prog.Lookup.LookupExternalType(tt); ok {
				cache[tt] = rt
				return rt
			}
		}
		// Mark this named type as being processed BEFORE recursing into the
		// underlying type. If we encounter it again via a pointer field, the
		// cache check at the top returns nil, and the *types.Pointer case
		// falls back to unsafe.Pointer.
		cache[tt] = nil
		// Pass a unique suffix based on the type name to prevent reflect.StructOf
		// from deduplicating different named types with the same field layout.
		// This is needed because reflect.StructOf caches types internally by
		// (fields, PkgPath), so two structs like GetterImpl{v int} and AdderStruct{v int}
		// would otherwise get the same reflect.Type.
		typeName := tt.Obj().Name()
		// Build uniqueSuffix with package-qualified name for _gig_id field:
		// "#PkgName.TypeName" (e.g., "#known_issues.point").
		// The _gig_id field uses this full suffix for fmt.Sprintf(%T) support.
		// Regular unexported fields strip the package prefix to keep type identity stable.
		qualSuffix := "#" + typeName
		if pkg := tt.Obj().Pkg(); pkg != nil {
			qualSuffix = "#" + pkg.Name() + "." + typeName
		}
		result := typeToReflectWithCache(tt.Underlying(), cache, qualSuffix, prog)
		if result != nil {
			cache[tt] = result
		}
		return result
	case *types.Struct:
		// Build struct type dynamically using reflect
		numFields := tt.NumFields()
		fields := make([]reflect.StructField, 0, numFields+1)
		hasUnexported := false
		for i := 0; i < numFields; i++ {
			f := tt.Field(i)
			// For named field types, use the type's own unique suffix
			// This ensures embedded structs maintain their type identity
			fieldSuffix := ""
			if named, ok := f.Type().(*types.Named); ok {
				fieldSuffix = "#" + named.Obj().Name()
			}
			ft := typeToReflectWithCache(f.Type(), cache, fieldSuffix, prog)
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
			//
			// CRITICAL: We append uniqueSuffix to PkgPath to prevent reflect.StructOf from
			// deduplicating different named types with the same field layout.
			// e.g., GetterImpl{v int} and AdderStruct{v int} should have different types.
			// The uniqueSuffix is passed from the *types.Named case and contains "#TypeName".
			if !f.Exported() {
				hasUnexported = true
				if sf.Anonymous {
					sf.Anonymous = false
				}
				// For regular unexported fields, use only the bare type suffix
				// (e.g., "#TypeName") to maintain type identity stability.
				// The full qualified suffix ("#PkgName.TypeName") is reserved
				// for the _gig_id sentinel field only.
				bareSuffix := uniqueSuffix
				if idx := strings.LastIndex(bareSuffix, "."); idx > 0 && bareSuffix[0] == '#' {
					bareSuffix = "#" + bareSuffix[idx+1:]
				}
				pkg := f.Pkg()
				if pkg != nil {
					pkgPath := pkg.Path()
					if bareSuffix != "" {
						pkgPath += bareSuffix
					}
					sf.PkgPath = pkgPath
				} else if bareSuffix != "" {
					sf.PkgPath = "gig/internal" + bareSuffix
				}
			}
			if tag := tt.Tag(i); tag != "" {
				sf.Tag = reflect.StructTag(tag)
			}
			fields = append(fields, sf)
		}
		// If the struct has only exported fields and a uniqueSuffix is provided,
		// we must add a phantom unexported field to force reflect.StructOf to create
		// a distinct type. Without this, structs like GetterHolder{Getter interface}
		// and any other struct{SomeInterface interface} would collide because
		// all interface fields become interface{} after conversion.
		if !hasUnexported && uniqueSuffix != "" {
			fields = append(fields, reflect.StructField{
				Name:    "_gig_id",
				Type:    reflect.TypeFor[struct{}](),
				PkgPath: "gig/internal" + uniqueSuffix,
			})
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
			pt := typeToReflectWithCache(params.At(i).Type(), cache, "", prog)
			if pt == nil {
				return nil
			}
			paramTypes[i] = pt
		}
		// Get result types
		results := tt.Results()
		resultTypes := make([]reflect.Type, results.Len())
		for i := 0; i < results.Len(); i++ {
			rt := typeToReflectWithCache(results.At(i).Type(), cache, "", prog)
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
