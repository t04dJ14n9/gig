package importer

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"reflect"
	"strings"
	"sync"
)

// typeCache caches converted types to prevent infinite recursion for self-referential types.
var typeCache sync.Map // map[reflect.Type]types.Type

func init() {
	// Initialize typeOf function
	typeOf = convertReflectType
}

// Importer implements types.Importer for registered packages.
type Importer struct {
	packages map[string]*types.Package
	mutex    sync.RWMutex
}

// NewImporter creates a new Importer.
func NewImporter() *Importer {
	return &Importer{
		packages: make(map[string]*types.Package),
	}
}

// Import returns the types.Package for the given import path.
func (i *Importer) Import(path string) (*types.Package, error) {
	i.mutex.RLock()
	if pkg, ok := i.packages[path]; ok {
		i.mutex.RUnlock()
		if pkg == nil {
			return nil, fmt.Errorf("package %q not found", path)
		}
		return pkg, nil
	}
	i.mutex.RUnlock()

	// Check if it's a registered external package
	extPkg := GetPackageByPath(path)
	if extPkg == nil {
		// Try to find by name (for auto-imported packages)
		extPkg = GetPackageByName(path)
		if extPkg == nil {
			return nil, fmt.Errorf("package %q not registered", path)
		}
	}

	// Build types.Package from external package
	pkg := i.buildPackage(extPkg)

	i.mutex.Lock()
	i.packages[path] = pkg
	i.mutex.Unlock()

	return pkg, nil
}

// buildPackage creates a types.Package from an ExternalPackage.
func (i *Importer) buildPackage(extPkg *ExternalPackage) *types.Package {
	pkg := types.NewPackage(extPkg.Path, extPkg.Name)

	// Add all objects to the package scope
	for name, obj := range extPkg.Objects {
		var typesObj types.Object

		switch obj.Kind {
		case ObjectKindFunction:
			if sig, ok := obj.Type.(*types.Signature); ok {
				typesObj = types.NewFunc(0, pkg, name, sig)
			}
		case ObjectKindVariable:
			typesObj = types.NewVar(0, pkg, name, obj.Type)
		case ObjectKindConstant:
			// Create a constant with the appropriate value
			val := convertToConstantValue(obj.Value)
			typesObj = types.NewConst(0, pkg, name, obj.Type, val)
		case ObjectKindType:
			// Type names are handled separately
			if named, ok := obj.Type.(*types.Named); ok {
				typesObj = named.Obj()
			} else {
				typeName := types.NewTypeName(0, pkg, name, obj.Type)
				typesObj = typeName
			}
		}

		if typesObj != nil {
			pkg.Scope().Insert(typesObj)
		}
	}

	// Add types
	for name, rt := range extPkg.Types {
		t := convertReflectType(rt)
		var typeName *types.TypeName

		if named, ok := t.(*types.Named); ok {
			typeName = named.Obj()
		} else {
			typeName = types.NewTypeName(0, pkg, name, t)
			// Create a new named type
			t = types.NewNamed(typeName, t, nil)
		}

		pkg.Scope().Insert(typeName)
		SetExternalType(t, rt)
	}

	// Mark the package as complete so the type checker can use it
	pkg.MarkComplete()

	return pkg
}

// convertToConstantValue converts a Go value to a types.Const value.
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
func convertReflectType(rt reflect.Type) types.Type {
	if rt == nil {
		return types.Typ[types.Invalid]
	}

	// Check cache first to prevent infinite recursion
	if cached, ok := typeCache.Load(rt); ok {
		return cached.(types.Type)
	}

	// For named types, handle specially to support self-referential types
	if rt.Name() != "" && rt.Kind() != reflect.Bool && rt.Kind() != reflect.Int &&
		rt.Kind() != reflect.Int8 && rt.Kind() != reflect.Int16 && rt.Kind() != reflect.Int32 &&
		rt.Kind() != reflect.Int64 && rt.Kind() != reflect.Uint && rt.Kind() != reflect.Uint8 &&
		rt.Kind() != reflect.Uint16 && rt.Kind() != reflect.Uint32 && rt.Kind() != reflect.Uint64 &&
		rt.Kind() != reflect.Uintptr && rt.Kind() != reflect.Float32 && rt.Kind() != reflect.Float64 &&
		rt.Kind() != reflect.Complex64 && rt.Kind() != reflect.Complex128 && rt.Kind() != reflect.String {
		// Create a placeholder named type to break recursion
		typeName := types.NewTypeName(0, nil, rt.Name(), nil)
		named := types.NewNamed(typeName, types.Typ[types.Invalid], nil)
		typeCache.Store(rt, named)

		// Now compute the actual underlying type
		underlying := convertReflectTypeForUnderlying(rt)
		named.SetUnderlying(underlying)

		// Add methods from reflect.Type to the Named type
		addMethodsToNamed(named, rt)

		return named
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
			typeName := types.NewTypeName(0, nil, rt.Name(), nil)
			underlying := convertReflectTypeForUnderlying(rt)
			result := types.NewNamed(typeName, underlying, nil)
			typeCache.Store(rt, result)
			return result
		}
		return types.Typ[types.Invalid]
	}
}

// convertReflectTypeForUnderlying handles underlying types for named types.
// It returns the actual underlying type (not a Named type).
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

// convertFuncType converts a reflect.Func type to types.Signature.
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

// convertInterfaceType converts a reflect.Interface type to types.Interface.
func convertInterfaceType(rt reflect.Type) *types.Interface {
	if rt.NumMethod() == 0 {
		// Empty interface (any)
		return types.NewInterfaceType(nil, nil)
	}

	var methods []*types.Func
	for i := 0; i < rt.NumMethod(); i++ {
		method := rt.Method(i)
		sig := convertFuncType(method.Type)
		methods = append(methods, types.NewFunc(0, nil, method.Name, sig))
	}

	return types.NewInterfaceType(methods, nil)
}

// convertStructType converts a reflect.Struct type to types.Struct.
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

// LookupPackage looks up a package by path or name.
func LookupPackage(name string) (*ExternalPackage, error) {
	// Try by path first
	if pkg := GetPackageByPath(name); pkg != nil {
		return pkg, nil
	}

	// Try by name
	if pkg := GetPackageByName(name); pkg != nil {
		return pkg, nil
	}

	return nil, fmt.Errorf("package %q not found", name)
}

// addMethodsToNamed adds methods from a reflect.Type to a types.Named type.
// This allows the type checker to find methods on external types (e.g., gjson.Result.String()).
func addMethodsToNamed(named *types.Named, rt reflect.Type) {
	// Enumerate all exported methods on the value receiver
	for i := 0; i < rt.NumMethod(); i++ {
		method := rt.Method(i)
		if !method.IsExported() {
			continue
		}
		// method.Type is func(ReceiverType, params...) (results...)
		// We need to create a *types.Signature with the receiver set.
		methodType := method.Type
		if methodType.Kind() != reflect.Func || methodType.NumIn() < 1 {
			continue
		}

		// Build the parameter list (excluding the receiver which is In(0))
		var params []*types.Var
		for j := 1; j < methodType.NumIn(); j++ {
			paramType := convertReflectType(methodType.In(j))
			params = append(params, types.NewVar(0, nil, "", paramType))
		}

		// Build the result list
		var results []*types.Var
		for j := 0; j < methodType.NumOut(); j++ {
			resultType := convertReflectType(methodType.Out(j))
			results = append(results, types.NewVar(0, nil, "", resultType))
		}

		// The receiver
		recv := types.NewVar(0, nil, "", named)

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

// AutoImport tries to find and import a package by name.
func AutoImport(name string) (path string, pkg *ExternalPackage, ok bool) {
	// Try exact match first
	if pkg := GetPackageByName(name); pkg != nil {
		return pkg.Path, pkg, true
	}

	// Try path suffix match (e.g., "json" matches "encoding/json")
	for path, pkg := range GetAllPackages() {
		parts := strings.Split(path, "/")
		if parts[len(parts)-1] == name {
			return path, pkg, true
		}
	}

	return "", nil, false
}
