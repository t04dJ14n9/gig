// typeconv.go converts go/types.Type to reflect.Type with cycle detection and caching.
package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
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
	if rt, ok := compositeTypeToReflect(t, cache, prog, depth); ok {
		return rt
	}

	switch tt := t.(type) {
	case *types.Basic:
		return bytecode.BasicKindToReflectType[tt.Kind()] // nil for unsupported kinds
	case *types.Named:
		return namedToReflect(tt, cache, prog, depth)
	case *types.Alias:
		// Type aliases (e.g., type MyInt = int) are identical to the aliased type.
		// Just use the underlying type.
		return typeToReflectWithCache(tt.Underlying(), cache, uniqueSuffix, prog, depth+1)
	case *types.Struct:
		return structToReflect(tt, cache, uniqueSuffix, prog, depth)
	case *types.Signature:
		return signatureToReflect(tt, cache, prog, depth)
	default:
		return nil
	}
}

// compositeTypeToReflect owns recursive container and indirection shapes.
// Keeping these away from the named/struct/signature router makes it easier to
// reason about cycle fallbacks, because every branch here advances depth through
// typeToReflectWithCache before constructing a reflect.Type.
func compositeTypeToReflect(t types.Type, cache map[types.Type]reflect.Type, prog *bytecode.CompiledProgram, depth int) (reflect.Type, bool) {
	switch tt := t.(type) {
	case *types.Slice:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if elem != nil {
			return reflect.SliceOf(elem), true
		}
		// If elem is nil due to a cycle (e.g., []*TreeNode where TreeNode has cycle),
		// use []any as a placeholder. The VM will convert slice elements
		// at assignment time when the concrete type is known.
		return reflect.SliceOf(reflect.TypeFor[any]()), true
	case *types.Array:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if elem != nil {
			return reflect.ArrayOf(int(tt.Len()), elem), true
		}
		return nil, true
	case *types.Map:
		key := typeToReflectWithCache(tt.Key(), cache, "", prog, depth+1)
		val := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if key != nil && val != nil {
			return reflect.MapOf(key, val), true
		}
		return nil, true
	case *types.Chan:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if elem != nil {
			return reflect.ChanOf(reflect.BothDir, elem), true
		}
		return nil, true
	case *types.Pointer:
		elem := typeToReflectWithCache(tt.Elem(), cache, "", prog, depth+1)
		if elem != nil {
			return reflect.PointerTo(elem), true
		}
		// If elem is nil due to a cycle (self-referencing struct pointer),
		// use any as a placeholder. The VM stores such values as
		// reflect.Value internally, and any can hold any pointer value.
		return reflect.TypeFor[any](), true
	case *types.Interface:
		// Interface type — use the empty interface (any) type
		// For the VM, all interfaces are represented as any
		return reflect.TypeFor[any](), true
	default:
		return nil, false
	}
}
