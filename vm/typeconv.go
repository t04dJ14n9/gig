// typeconv.go converts go/types.Type to reflect.Type with cycle detection and caching.
package vm

import (
	"go/types"
	"reflect"
	"strings"

	"git.woa.com/youngjin/gig/model/bytecode"
)

// maxTypeRecursionDepth is the hard limit on recursive type conversion depth.
// Prevents stack overflow on deeply nested but acyclic types (e.g., [][][]...[]int).
const maxTypeRecursionDepth = 256

// typeToReflect converts a go/types.Type to reflect.Type using the program-level
// cache to ensure the same types.Type always maps to the same reflect.Type.
// This prevents reflect.StructOf from returning different reflect.Type objects
// across multiple VM executions, which would cause "reflect.Set: value not assignable" panics.
func typeToReflect(t types.Type, prog *bytecode.CompiledProgram) reflect.Type {
	if t == nil {
		return nil
	}
	// Check program-level cache first (fast path, lock-free read)
	if rt, ok := prog.CachedReflectType(t); ok {
		return rt
	}
	// Compute with a local cycle-detection cache
	localCache := make(map[types.Type]reflect.Type)
	rt := typeToReflectWithCache(t, localCache, "", prog, 0)
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
// The depth parameter prevents stack overflow on deeply nested acyclic types.
//
// NOTE: This function does NOT use the program-level cache internally because the same
// types.Type (e.g., *types.Struct for struct{v int}) may be reached through different
// named types with different uniqueSuffix values. Caching at this level would cause
// suffix-insensitive collisions. Program-level caching is done only at the top-level
// typeToReflect entry point, which caches the final result keyed by the original type.
func typeToReflectWithCache(t types.Type, cache map[types.Type]reflect.Type, uniqueSuffix string, prog *bytecode.CompiledProgram, depth int) reflect.Type {
	if t == nil {
		return nil
	}
	if depth > maxTypeRecursionDepth {
		return nil
	}

	// Check local cache first (for cycle detection)
	if cached, ok := cache[t]; ok {
		return cached
	}

	result := typeToReflectInner(t, cache, uniqueSuffix, prog, depth)
	if result != nil {
		cache[t] = result
	}

	return result
}

// typeToReflectInner does the actual conversion without caching logic.
func typeToReflectInner(t types.Type, cache map[types.Type]reflect.Type, uniqueSuffix string, prog *bytecode.CompiledProgram, depth int) reflect.Type {
	switch tt := t.(type) {
	case *types.Basic:
		return bytecode.BasicKindToReflectType[tt.Kind()] // nil for unsupported kinds
	case *types.Slice:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if elem != nil {
			return reflect.SliceOf(elem)
		}
		// If elem is nil due to a cycle (e.g., []*TreeNode where TreeNode has cycle),
		// use []any as a placeholder. The VM will convert slice elements
		// at assignment time when the concrete type is known.
		return reflect.SliceOf(reflect.TypeFor[any]())
	case *types.Array:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if elem != nil {
			return reflect.ArrayOf(int(tt.Len()), elem)
		}
		return nil
	case *types.Map:
		key := typeToReflectWithCache(tt.Key(), cache, "", prog, depth+1)
		val := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if key != nil && val != nil {
			return reflect.MapOf(key, val)
		}
		return nil
	case *types.Chan:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if elem != nil {
			return reflect.ChanOf(reflect.BothDir, elem)
		}
		return nil
	case *types.Pointer:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if elem != nil {
			return reflect.PointerTo(elem)
		}
		// If elem is nil due to a cycle (self-referencing struct pointer),
		// use any as a placeholder. The VM stores such values as
		// reflect.Value internally, and any can hold any pointer value.
		return reflect.TypeFor[any]()
	case *types.Interface:
		// Interface type — use the empty interface (any) type
		// For the VM, all interfaces are represented as any
		return reflect.TypeFor[any]()
	case *types.Named:
		// Check if this is a registered external type (e.g., bytes.Buffer, strings.Builder).
		// If so, use the real reflect.Type instead of synthesizing a struct type.
		if prog != nil && prog.TypeResolver != nil {
			if rt, ok := prog.TypeResolver.LookupExternalType(tt); ok {
				cache[tt] = rt
				return rt
			}
			// Fallback: lookup by package path + type name. The SSA/type-checker creates
			// its own *types.Named objects that differ by pointer identity from the ones
			// stored via AddType/SetExternalType. A name-based lookup resolves this
			// mismatch for external named types like sort.IntSlice.
			obj := tt.Obj()
			if pkg := obj.Pkg(); pkg != nil {
				if rt, ok := prog.TypeResolver.LookupExternalTypeByName(pkg.Path(), obj.Name()); ok {
					cache[tt] = rt
					return rt
				}
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
		// Build uniqueSuffix for struct tag uniqueness and type name registry.
		// Format: "#PkgName.TypeName" (e.g., "#known_issues.point").
		qualSuffix := "#" + typeName
		if pkg := tt.Obj().Pkg(); pkg != nil {
			qualSuffix = "#" + pkg.Name() + "." + typeName
		}
		result := typeToReflectWithCache(tt.Underlying(), cache, qualSuffix, prog, depth+1)
		if result != nil {
			cache[tt] = result
			// Register the type name in the program-level registry for method dispatch.
			// This replaces the old _gig_id phantom field approach.
			if prog != nil && result.Kind() == reflect.Struct {
				prog.RegisterTypeName(result, typeName)
			}
		}
		return result
	case *types.Alias:
		// Type aliases (e.g., type MyInt = int) are identical to the aliased type.
		// Just use the underlying type.
		return typeToReflectWithCache(tt.Underlying(), cache, uniqueSuffix, prog, depth+1)
	case *types.Struct:
		// Build struct type dynamically using reflect
		numFields := tt.NumFields()
		fields := make([]reflect.StructField, 0, numFields)
		hasUnexported := false
		for i := 0; i < numFields; i++ {
			f := tt.Field(i)
			// For named field types, use the type's own unique suffix
			// This ensures embedded structs maintain their type identity
			fieldSuffix := ""
			if named, ok := f.Type().(*types.Named); ok {
				fieldSuffix = "#" + named.Obj().Name()
			}
			ft := typeToReflectWithCache(f.Type(), cache, fieldSuffix, prog, depth+1)
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
				// For unexported fields, use the bare type suffix
				// (e.g., "#TypeName") to maintain type identity stability.
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
		// For structs with only exported fields and a uniqueSuffix, we must make
		// each field distinguishable via a struct tag so that reflect.StructOf
		// produces a unique reflect.Type. Without this, structs like
		// GetterHolder{Getter interface} and any other struct{SomeInterface interface}
		// would collide because all interface fields become any after conversion.
		//
		// For empty named structs (no fields, uniqueSuffix != ""), we share
		// reflect.TypeOf(struct{}{}) and rely on the ReflectTypeNames registry
		// for method dispatch.
		if !hasUnexported && uniqueSuffix != "" && len(fields) > 0 {
			gigTag := reflect.StructTag(`gig:"` + uniqueSuffix + `"`)
			for i := range fields {
				if fields[i].Tag == "" {
					fields[i].Tag = gigTag
				} else {
					fields[i].Tag = fields[i].Tag + " " + gigTag
				}
			}
		}
		if len(fields) == 0 {
			// Empty struct (struct{} or named empty struct) — return the real Go type directly.
			// For named empty structs, the ReflectTypeNames registry handles identification.
			return reflect.TypeOf(struct{}{})
		}
		result := reflect.StructOf(fields)

		return result
	case *types.Signature:
		// Function type - need to build the function type dynamically
		// Get parameter types
		params := tt.Params()
		paramTypes := make([]reflect.Type, params.Len())
		for i := 0; i < params.Len(); i++ {
			pt := typeToReflectWithCache(params.At(i).Type(), cache, "", prog, depth+1)
			if pt == nil {
				return nil
			}
			paramTypes[i] = pt
		}
		// Get result types
		results := tt.Results()
		resultTypes := make([]reflect.Type, results.Len())
		for i := 0; i < results.Len(); i++ {
			rt := typeToReflectWithCache(results.At(i).Type(), cache, "", prog, depth+1)
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
