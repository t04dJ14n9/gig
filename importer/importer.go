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
			typesObj = types.NewFunc(0, pkg, name, obj.Type.(*types.Signature))
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

	// Check cache first
	if ext := GetExternalType(nil); ext != nil {
		// This is just to ensure the cache is initialized
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
		return types.NewArray(elem, int64(rt.Len()))

	case reflect.Slice:
		elem := convertReflectType(rt.Elem())
		return types.NewSlice(elem)

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

	case reflect.Func:
		return convertFuncType(rt)

	case reflect.Interface:
		return convertInterfaceType(rt)

	case reflect.Map:
		key := convertReflectType(rt.Key())
		elem := convertReflectType(rt.Elem())
		return types.NewMap(key, elem)

	case reflect.Ptr:
		elem := convertReflectType(rt.Elem())
		return types.NewPointer(elem)

	case reflect.Struct:
		return convertStructType(rt)

	default:
		// For named types, create a TypeName
		if rt.Name() != "" {
			typeName := types.NewTypeName(0, nil, rt.Name(), nil)
			// For named types, use the underlying type
			underlying := convertReflectTypeForUnderlying(rt)
			return types.NewNamed(typeName, underlying, nil)
		}
		return types.Typ[types.Invalid]
	}
}

// convertReflectTypeForUnderlying handles underlying types for named types.
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
	default:
		return convertReflectType(rt)
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
