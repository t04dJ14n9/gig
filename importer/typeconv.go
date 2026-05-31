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

type reflectInterfaceConverter func(reflect.Type) types.Type

// convertReflectType converts a reflect.Type to types.Type.
// This entry point only coordinates the high-level policy decisions:
// cache lookup, universe aliases, named type identity, then structural shape.
// Keeping those boundaries explicit matters because named recursive types must
// install a cache placeholder before their fields or methods are converted.
func convertReflectType(rt reflect.Type) types.Type {
	if rt == nil {
		return types.Typ[types.Invalid]
	}

	if cached, ok := cachedReflectType(rt); ok {
		return cached
	}

	if t, ok := convertUniverseReflectInterface(rt); ok {
		return t
	}

	if named, ok := convertNamedReflectType(rt); ok {
		return named
	}

	return convertUnnamedReflectType(rt)
}

// convertReflectTypeForUnderlying converts a reflect.Type to its underlying types.Type.
// This is used for named types to avoid wrapping in another Named type.
// Unlike convertReflectType, it does NOT cache results in typeCache — doing so
// would overwrite the Named type entry that the caller stored for recursion breaking.
func convertReflectTypeForUnderlying(rt reflect.Type) types.Type {
	if bt := bytecode.BasicTypeFromReflectKind(rt.Kind()); bt != nil {
		return bt
	}

	return convertReflectCompositeType(rt, convertUncachedReflectInterface)
}

func cachedReflectType(rt reflect.Type) (types.Type, bool) {
	cached, ok := typeCache.Load(rt)
	if !ok {
		return nil, false
	}
	return cached.(types.Type), true
}

func cacheReflectType(rt reflect.Type, t types.Type) types.Type {
	typeCache.Store(rt, t)
	return t
}

// convertUniverseReflectInterface preserves the identity of predeclared
// interface types such as error. reflect.Type exposes error as an unnamed
// interface shape, but go/types requires the universe object for assignment
// compatibility with real Go code.
func convertUniverseReflectInterface(rt reflect.Type) (types.Type, bool) {
	if rt.Kind() != reflect.Interface || rt.PkgPath() != "" {
		return nil, false
	}

	obj := types.Universe.Lookup(rt.Name())
	if obj == nil {
		return nil, false
	}

	return cacheReflectType(rt, obj.Type()), true
}

// convertNamedReflectType preserves type identity for imported Go types.
// The placeholder is deliberately cached before computing the underlying type:
// a recursive struct like type Node struct { Next *Node } must resolve *Node
// back to the same named shell while the fields are still being built.
func convertNamedReflectType(rt reflect.Type) (types.Type, bool) {
	if rt.Name() == "" || isBasicReflectAlias(rt) {
		return nil, false
	}

	pkg := getOrCreateTypesPackage(rt.PkgPath())
	typeName := types.NewTypeName(0, pkg, rt.Name(), nil)
	named := types.NewNamed(typeName, types.Typ[types.Invalid], nil)
	typeCache.Store(rt, named)

	named.SetUnderlying(convertReflectTypeForUnderlying(rt))

	// Named interfaces already expose methods through their underlying
	// interface. Concrete named types need the reflected method set attached so
	// calls through imported third-party APIs type-check correctly.
	if rt.Kind() != reflect.Interface {
		addMethodsToNamed(named, rt)
	}

	return named, true
}

func isBasicReflectAlias(rt reflect.Type) bool {
	bt := bytecode.BasicTypeFromReflectKind(rt.Kind())
	return bt != nil && bt.Name() == rt.Name()
}

func convertUnnamedReflectType(rt reflect.Type) types.Type {
	if bt := bytecode.BasicTypeFromReflectKind(rt.Kind()); bt != nil {
		return bt
	}

	return cacheReflectType(rt, convertReflectCompositeType(rt, convertCachedReflectInterface))
}

func convertReflectCompositeType(rt reflect.Type, convertInterface reflectInterfaceConverter) types.Type {
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
		return convertInterface(rt)
	case reflect.Chan:
		return convertChanType(rt)
	default:
		return types.Typ[types.Invalid]
	}
}

func convertCachedReflectInterface(rt reflect.Type) types.Type {
	return convertInterfaceType(rt)
}

func convertUncachedReflectInterface(rt reflect.Type) types.Type {
	return convertInterfaceTypeNoCache(rt)
}

func convertChanType(rt reflect.Type) *types.Chan {
	return types.NewChan(chanDirFromReflect(rt.ChanDir()), convertReflectType(rt.Elem()))
}

func chanDirFromReflect(dir reflect.ChanDir) types.ChanDir {
	switch dir {
	case reflect.SendDir:
		return types.SendOnly
	case reflect.RecvDir:
		return types.RecvOnly
	default:
		return types.SendRecv
	}
}
