package importer

import (
	"go/constant"
	"go/token"
	"go/types"
	"reflect"
	"strings"
	"sync"
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

// typePkgCache caches *types.Package objects by package path to ensure
// that the same package path always maps to the same *types.Package instance.
var typePkgCache sync.Map // map[string]*types.Package

// getOrCreateTypesPackage returns a cached *types.Package for the given
// package path, creating one if it doesn't exist yet. The package name is
// derived from the last path segment (e.g., "encoding/json" → "json").
func getOrCreateTypesPackage(pkgPath string) *types.Package {
	if pkgPath == "" {
		return nil
	}
	if cached, ok := typePkgCache.Load(pkgPath); ok {
		return cached.(*types.Package)
	}
	// Derive package name from path (last segment)
	name := pkgPath
	if idx := strings.LastIndex(pkgPath, "/"); idx >= 0 {
		name = pkgPath[idx+1:]
	}
	pkg := types.NewPackage(pkgPath, name)
	actual, _ := typePkgCache.LoadOrStore(pkgPath, pkg)
	return actual.(*types.Package)
}

// convertToConstantValue converts a Go value to a constant.Value for use in types.Const.
// Handles basic types (bool, int, uint, float, complex, string) and falls back
// to reflection for other types.
func convertToConstantValue(val any) constant.Value {
	switch v := val.(type) {
	case bool:
		return constant.MakeBool(v)
	case int:
		return constant.MakeInt64(int64(v))
	case int8:
		return constant.MakeInt64(int64(v))
	case int16:
		return constant.MakeInt64(int64(v))
	case int32:
		return constant.MakeInt64(int64(v))
	case int64:
		return constant.MakeInt64(v)
	case uint:
		return constant.MakeUint64(uint64(v))
	case uint8:
		return constant.MakeUint64(uint64(v))
	case uint16:
		return constant.MakeUint64(uint64(v))
	case uint32:
		return constant.MakeUint64(uint64(v))
	case uint64:
		return constant.MakeUint64(v)
	case float32:
		return constant.MakeFloat64(float64(v))
	case float64:
		return constant.MakeFloat64(v)
	case complex64:
		// complex values are represented as binary operations
		re := constant.MakeFloat64(float64(real(v)))
		im := constant.MakeFloat64(float64(imag(v)))
		return constant.BinaryOp(re, token.ADD, constant.BinaryOp(im, token.MUL, constant.MakeImag(constant.MakeInt64(1))))
	case complex128:
		re := constant.MakeFloat64(real(v))
		im := constant.MakeFloat64(imag(v))
		return constant.BinaryOp(re, token.ADD, constant.BinaryOp(im, token.MUL, constant.MakeImag(constant.MakeInt64(1))))
	case string:
		return constant.MakeString(v)
	default:
		// For other types, try to convert via reflection
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Bool:
			return constant.MakeBool(rv.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return constant.MakeInt64(rv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return constant.MakeUint64(rv.Uint())
		case reflect.Float32, reflect.Float64:
			return constant.MakeFloat64(rv.Float())
		case reflect.Complex64, reflect.Complex128:
			c := rv.Complex()
			re := constant.MakeFloat64(real(c))
			im := constant.MakeFloat64(imag(c))
			return constant.BinaryOp(re, token.ADD, constant.BinaryOp(im, token.MUL, constant.MakeImag(constant.MakeInt64(1))))
		case reflect.String:
			return constant.MakeString(rv.String())
		default:
			return constant.MakeUnknown()
		}
	}
}

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
		// Convert reflect.Kind to types.Kind for indexing types.Typ
		var typeKind types.BasicKind
		switch rt.Kind() {
		case reflect.Bool:
			typeKind = types.Bool
		case reflect.Int:
			typeKind = types.Int
		case reflect.Int8:
			typeKind = types.Int8
		case reflect.Int16:
			typeKind = types.Int16
		case reflect.Int32:
			typeKind = types.Int32
		case reflect.Int64:
			typeKind = types.Int64
		case reflect.Uint:
			typeKind = types.Uint
		case reflect.Uint8:
			typeKind = types.Uint8
		case reflect.Uint16:
			typeKind = types.Uint16
		case reflect.Uint32:
			typeKind = types.Uint32
		case reflect.Uint64:
			typeKind = types.Uint64
		case reflect.Uintptr:
			typeKind = types.Uintptr
		case reflect.Float32:
			typeKind = types.Float32
		case reflect.Float64:
			typeKind = types.Float64
		case reflect.Complex64:
			typeKind = types.Complex64
		case reflect.Complex128:
			typeKind = types.Complex128
		case reflect.String:
			typeKind = types.String
		case reflect.UnsafePointer:
			typeKind = types.UnsafePointer
		default:
			// Not a basic kind, can't be a basic type alias
			typeKind = types.Invalid
		}
		if typeKind != types.Invalid {
			bt := types.Typ[typeKind]
			if bt != nil && bt.Name() == rt.Name() {
				isBasicAlias = true
			}
		}
		// If it's a basic type alias, don't treat as named type - fall through to basic handling
		// For interface types, don't wrap in Named - just return the interface directly
		// (named interfaces in Go are still interface types, not Named types)
		if !isBasicAlias && rt.Kind() != reflect.Interface {
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

			// Add methods from reflect.Type to the Named type
			addMethodsToNamed(named, rt)

			return named
		}
	}

	switch rt.Kind() {
	case reflect.Bool:
		return types.Typ[types.Bool]
	case reflect.Int:
		return types.Typ[types.Int]
	case reflect.Int8:
		return types.Typ[types.Int8]
	case reflect.Int16:
		return types.Typ[types.Int16]
	case reflect.Int32:
		return types.Typ[types.Int32]
	case reflect.Int64:
		return types.Typ[types.Int64]
	case reflect.Uint:
		return types.Typ[types.Uint]
	case reflect.Uint8:
		return types.Typ[types.Uint8]
	case reflect.Uint16:
		return types.Typ[types.Uint16]
	case reflect.Uint32:
		return types.Typ[types.Uint32]
	case reflect.Uint64:
		return types.Typ[types.Uint64]
	case reflect.Uintptr:
		return types.Typ[types.Uintptr]
	case reflect.Float32:
		return types.Typ[types.Float32]
	case reflect.Float64:
		return types.Typ[types.Float64]
	case reflect.Complex64:
		return types.Typ[types.Complex64]
	case reflect.Complex128:
		return types.Typ[types.Complex128]
	case reflect.String:
		return types.Typ[types.String]
	case reflect.UnsafePointer:
		return types.Typ[types.UnsafePointer]

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
// It returns the actual underlying type (e.g., struct, pointer, slice).
func convertReflectTypeForUnderlying(rt reflect.Type) types.Type {
	// For basic named types, return the corresponding basic type
	switch rt.Kind() {
	case reflect.Bool:
		return types.Typ[types.Bool]
	case reflect.Int:
		return types.Typ[types.Int]
	case reflect.Int8:
		return types.Typ[types.Int8]
	case reflect.Int16:
		return types.Typ[types.Int16]
	case reflect.Int32:
		return types.Typ[types.Int32]
	case reflect.Int64:
		return types.Typ[types.Int64]
	case reflect.Uint:
		return types.Typ[types.Uint]
	case reflect.Uint8:
		return types.Typ[types.Uint8]
	case reflect.Uint16:
		return types.Typ[types.Uint16]
	case reflect.Uint32:
		return types.Typ[types.Uint32]
	case reflect.Uint64:
		return types.Typ[types.Uint64]
	case reflect.Uintptr:
		return types.Typ[types.Uintptr]
	case reflect.Float32:
		return types.Typ[types.Float32]
	case reflect.Float64:
		return types.Typ[types.Float64]
	case reflect.Complex64:
		return types.Typ[types.Complex64]
	case reflect.Complex128:
		return types.Typ[types.Complex128]
	case reflect.String:
		return types.Typ[types.String]
	case reflect.Struct:
		return convertStructType(rt)
	case reflect.Ptr:
		elem := convertReflectType(rt.Elem())
		return types.NewPointer(elem)
	case reflect.Slice:
		elem := convertReflectType(rt.Elem())
		return types.NewSlice(elem)
	case reflect.Array:
		elem := convertReflectType(rt.Elem())
		return types.NewArray(elem, int64(rt.Len()))
	case reflect.Map:
		key := convertReflectType(rt.Key())
		elem := convertReflectType(rt.Elem())
		return types.NewMap(key, elem)
	case reflect.Func:
		return convertFuncType(rt)
	case reflect.Interface:
		return convertInterfaceType(rt)
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

// convertFuncType converts a reflect.Func type to a types.Signature.
// It builds parameter and result tuples from the function's type information.
func convertFuncType(rt reflect.Type) *types.Signature {
	// Build parameter types
	var params []*types.Var
	for i := 0; i < rt.NumIn(); i++ {
		paramType := convertReflectType(rt.In(i))
		params = append(params, types.NewVar(0, nil, "", paramType))
	}

	// Build result types
	var results []*types.Var
	for i := 0; i < rt.NumOut(); i++ {
		resultType := convertReflectType(rt.Out(i))
		results = append(results, types.NewVar(0, nil, "", resultType))
	}

	// Check if last param is variadic
	variadic := rt.IsVariadic()

	return types.NewSignatureType(
		nil,      // recv
		nil, nil, // type params
		types.NewTuple(params...),
		types.NewTuple(results...),
		variadic,
	)
}

// convertInterfaceType converts a reflect.Interface type to a types.Interface.
// For empty interfaces (any), returns types.NewInterfaceType(nil, nil).
func convertInterfaceType(rt reflect.Type) *types.Interface {
	if rt.NumMethod() == 0 {
		// Empty interface (any)
		return types.NewInterfaceType(nil, nil)
	}

	// Create a placeholder interface and cache it first to break recursion
	iface := types.NewInterfaceType(nil, nil)
	typeCache.Store(rt, iface)

	var methods []*types.Func
	for i := 0; i < rt.NumMethod(); i++ {
		method := rt.Method(i)
		sig := convertFuncType(method.Type)
		methods = append(methods, types.NewFunc(0, nil, method.Name, sig))
	}

	// Update the interface with the methods
	// Note: types.NewInterfaceType creates a complete interface, so we need to
	// create a new one with the methods
	result := types.NewInterfaceType(methods, nil)
	typeCache.Store(rt, result)
	return result
}

// convertStructType converts a reflect.Struct type to a types.Struct.
// It preserves field names, types, anonymous fields, and struct tags.
func convertStructType(rt reflect.Type) *types.Struct {
	var fields []*types.Var
	var tags []string

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldType := convertReflectType(field.Type)
		fields = append(fields, types.NewField(0, nil, field.Name, fieldType, field.Anonymous))
		tags = append(tags, string(field.Tag))
	}

	return types.NewStruct(fields, tags)
}

// addMethodsToNamed adds methods from a reflect.Type to a types.Named type.
// This allows the type checker to find methods on external types.
// Both value receiver and pointer receiver methods are added.
// For interface types, methods do NOT include a receiver parameter.
func addMethodsToNamed(named *types.Named, rt reflect.Type) {
	// Check if this is an interface type - interface methods don't have receiver params
	isInterface := rt.Kind() == reflect.Interface

	// Enumerate all exported methods on the value receiver
	for i := 0; i < rt.NumMethod(); i++ {
		method := rt.Method(i)
		if !method.IsExported() {
			continue
		}
		// method.Type is:
		//   - For concrete types: func(ReceiverType, params...) (results...)
		//   - For interface types: func(params...) (results...) - NO receiver
		methodType := method.Type
		if methodType.Kind() != reflect.Func {
			continue
		}
		// For non-interface types, we need at least 1 param (the receiver)
		if !isInterface && methodType.NumIn() < 1 {
			continue
		}

		// Build the parameter list
		// For interfaces, start at 0 (no receiver). For concrete types, skip receiver at 0.
		var params []*types.Var
		startIdx := 0
		if !isInterface {
			startIdx = 1
		}
		for j := startIdx; j < methodType.NumIn(); j++ {
			paramType := convertReflectType(methodType.In(j))
			params = append(params, types.NewVar(0, nil, "", paramType))
		}

		// Build the result list
		var results []*types.Var
		for j := 0; j < methodType.NumOut(); j++ {
			resultType := convertReflectType(methodType.Out(j))
			results = append(results, types.NewVar(0, nil, "", resultType))
		}

		// For non-interface types, the receiver is the named type itself
		// For interface types, there is no receiver (nil)
		var recv *types.Var
		if !isInterface {
			recv = types.NewVar(0, nil, "", named)
		}

		sig := types.NewSignatureType(
			recv,
			nil, nil,
			types.NewTuple(params...),
			types.NewTuple(results...),
			methodType.IsVariadic(),
		)

		fn := types.NewFunc(0, named.Obj().Pkg(), method.Name, sig)
		named.AddMethod(fn)
	}

	// For interface types, don't process pointer receiver methods
	// (interfaces don't have pointer receivers in the same way)
	if isInterface {
		return
	}

	// Also enumerate methods on the pointer receiver (*T)
	ptrType := reflect.PtrTo(rt)
	for i := 0; i < ptrType.NumMethod(); i++ {
		method := ptrType.Method(i)
		if !method.IsExported() {
			continue
		}
		// Skip methods that are already added (value receiver methods are a subset of pointer receiver methods)
		alreadyAdded := false
		for j := 0; j < named.NumMethods(); j++ {
			if named.Method(j).Name() == method.Name {
				alreadyAdded = true
				break
			}
		}
		if alreadyAdded {
			continue
		}

		methodType := method.Type
		if methodType.Kind() != reflect.Func || methodType.NumIn() < 1 {
			continue
		}

		var params []*types.Var
		for j := 1; j < methodType.NumIn(); j++ {
			paramType := convertReflectType(methodType.In(j))
			params = append(params, types.NewVar(0, nil, "", paramType))
		}

		var results []*types.Var
		for j := 0; j < methodType.NumOut(); j++ {
			resultType := convertReflectType(methodType.Out(j))
			results = append(results, types.NewVar(0, nil, "", resultType))
		}

		// Pointer receiver
		recv := types.NewVar(0, nil, "", types.NewPointer(named))

		sig := types.NewSignatureType(
			recv,
			nil, nil,
			types.NewTuple(params...),
			types.NewTuple(results...),
			methodType.IsVariadic(),
		)

		fn := types.NewFunc(0, named.Obj().Pkg(), method.Name, sig)
		named.AddMethod(fn)
	}
}
