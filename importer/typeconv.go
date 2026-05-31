// typeconv.go converts reflect.Type to go/types.Type for external package types.
package importer

import (
	"go/types"
	"reflect"
	"sync"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// typeOf is a function that converts reflect.Type to types.Type.
// It is initialized by init() to break the circular dependency between
// register.go (which uses typeOf in AddVariable/AddConstant) and typeconv.go.
var typeOf func(reflect.Type) types.Type

func init() {
	typeOf = convertReflectType
}

// typeCache caches converted types to prevent infinite recursion for self-referential types.
var typeCache sync.Map // map[reflect.Type]types.Type

// convertReflectType converts a reflect.Type to types.Type.
// This is the main entry point for type conversion, handling all Go types.
// Uses a cache to prevent infinite recursion for self-referential types.
func convertReflectType(rt reflect.Type) types.Type {
	if rt == nil {
		return types.Typ[types.Invalid]
	}

	// Check cache first to prevent infinite recursion
	if cached, ok := typeCache.Load(rt); ok {
		return cached.(types.Type)
	}

	// Map well-known universe interfaces to their named types.
	// Without this, reflect sees "error" as an anonymous interface{Error() string}
	// which doesn't match the named "error" type from Go's universe scope.
	// This causes type check failures like: cannot use []error as []interface{Error() string}.
	if rt.Kind() == reflect.Interface && rt.PkgPath() == "" {
		if obj := types.Universe.Lookup(rt.Name()); obj != nil {
			t := obj.Type()
			typeCache.Store(rt, t)
			return t
		}
	}

	// For named types, handle specially to support self-referential types.
	// Named types include types like "time.Duration" which has Kind() == reflect.Int64
	// but has a distinct name. We must preserve the named type, not collapse to the
	// underlying basic type.
	// However, basic type aliases like "byte" (alias for uint8) have Name() != "" but
	// are NOT true named types - they should be treated as basic types to preserve
	// byte/uint8 compatibility. We detect this by checking if the type name matches
	// the underlying kind's name in types.Typ.
	if rt.Name() != "" {
		// Check if this is a basic type with a name alias (e.g., uint8, byte, string).
		// For basic kinds, types.Typ[kind] gives the canonical type with the kind's name.
		// If the reflect type's name matches the canonical kind name, it's a basic type alias.
		isBasicAlias := false
		if bt := bytecode.BasicTypeFromReflectKind(rt.Kind()); bt != nil && bt.Name() == rt.Name() {
			isBasicAlias = true
		}
		// If it's a basic type alias, don't treat as named type - fall through to basic handling.
		if !isBasicAlias {
			// Create a placeholder named type to break recursion.
			// Use rt.PkgPath() to attach the correct package, so the type checker
			// and compiler can distinguish types with the same name from different
			// packages (e.g., encoding/json.Encoder vs encoding/xml.Encoder).
			pkg := getOrCreateTypesPackage(rt.PkgPath())
			typeName := types.NewTypeName(0, pkg, rt.Name(), nil)
			named := types.NewNamed(typeName, types.Typ[types.Invalid], nil)
			typeCache.Store(rt, named)

			// Now compute the actual underlying type
			underlying := convertReflectTypeForUnderlying(rt)
			named.SetUnderlying(underlying)

			// Add methods from reflect.Type to concrete named types. Named interfaces
			// get their method set from the underlying interface; adding the same
			// methods to the Named type would duplicate them.
			if rt.Kind() != reflect.Interface {
				addMethodsToNamed(named, rt)
			}

			return named
		}
	}

	switch rt.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.String, reflect.UnsafePointer:
		return bytecode.BasicTypeFromReflectKind(rt.Kind())

	case reflect.Array:
		elem := convertReflectType(rt.Elem())
		result := types.NewArray(elem, int64(rt.Len()))
		typeCache.Store(rt, result)
		return result

	case reflect.Slice:
		elem := convertReflectType(rt.Elem())
		result := types.NewSlice(elem)
		typeCache.Store(rt, result)
		return result

	case reflect.Chan:
		elem := convertReflectType(rt.Elem())
		var dir types.ChanDir
		switch rt.ChanDir() {
		case reflect.SendDir:
			dir = types.SendOnly
		case reflect.RecvDir:
			dir = types.RecvOnly
		default:
			dir = types.SendRecv
		}
		result := types.NewChan(dir, elem)
		typeCache.Store(rt, result)
		return result

	case reflect.Func:
		result := convertFuncType(rt)
		typeCache.Store(rt, result)
		return result

	case reflect.Interface:
		result := convertInterfaceType(rt)
		typeCache.Store(rt, result)
		return result

	case reflect.Map:
		key := convertReflectType(rt.Key())
		elem := convertReflectType(rt.Elem())
		result := types.NewMap(key, elem)
		typeCache.Store(rt, result)
		return result

	case reflect.Ptr:
		elem := convertReflectType(rt.Elem())
		result := types.NewPointer(elem)
		typeCache.Store(rt, result)
		return result

	case reflect.Struct:
		result := convertStructType(rt)
		typeCache.Store(rt, result)
		return result

	default:
		// For other named types, create a TypeName
		if rt.Name() != "" {
			pkg := getOrCreateTypesPackage(rt.PkgPath())
			typeName := types.NewTypeName(0, pkg, rt.Name(), nil)
			underlying := convertReflectTypeForUnderlying(rt)
			result := types.NewNamed(typeName, underlying, nil)
			typeCache.Store(rt, result)
			return result
		}
		return types.Typ[types.Invalid]
	}
}

// convertReflectTypeForUnderlying converts a reflect.Type to its underlying types.Type.
// This is used for named types to avoid wrapping in another Named type.
// Unlike convertReflectType, it does NOT cache results in typeCache — doing so
// would overwrite the Named type entry that the caller stored for recursion breaking.
func convertReflectTypeForUnderlying(rt reflect.Type) types.Type {
	// For basic named types, return the corresponding basic type
	if bt := bytecode.BasicTypeFromReflectKind(rt.Kind()); bt != nil {
		return bt
	}

	switch rt.Kind() {
	case reflect.Struct:
		return convertStructType(rt)
	case reflect.Ptr:
		return types.NewPointer(convertReflectType(rt.Elem()))
	case reflect.Slice:
		return types.NewSlice(convertReflectType(rt.Elem()))
	case reflect.Array:
		return types.NewArray(convertReflectType(rt.Elem()), int64(rt.Len()))
	case reflect.Map:
		return types.NewMap(convertReflectType(rt.Key()), convertReflectType(rt.Elem()))
	case reflect.Func:
		return convertFuncType(rt)
	case reflect.Interface:
		return convertInterfaceTypeNoCache(rt)
	case reflect.Chan:
		elem := convertReflectType(rt.Elem())
		var dir types.ChanDir
		switch rt.ChanDir() {
		case reflect.SendDir:
			dir = types.SendOnly
		case reflect.RecvDir:
			dir = types.RecvOnly
		default:
			dir = types.SendRecv
		}
		return types.NewChan(dir, elem)
	default:
		return types.Typ[types.Invalid]
	}
}
