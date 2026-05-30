// register.go implements the global registry: package registration, type mapping,
// method DirectCall lookup, and auto-import.
package importer

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"sync"

	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

// ExternalPackage represents a registered external package.
// It maps package-level objects by name for lookup during type checking and execution.
type ExternalPackage struct {
	// Path is the import path (e.g., "fmt", "encoding/json").
	Path string

	// Name is the package identifier (e.g., "fmt", "json").
	Name string

	// Objects maps object names to their ExternalObject entries.
	Objects map[string]*external.ExternalObject

	// Types maps type names to their reflect.Type representations.
	Types map[string]reflect.Type

	// InterfaceProxies maps interface type names to native proxy metadata.
	InterfaceProxies map[string]*external.InterfaceProxyInfo

	// registry is a back-reference to the owning registry.
	registry PackageRegistry
}

// PackageRegistry manages external package registration and lookup.
// It provides methods to register, lookup, and query packages, types, and method DirectCalls.
type PackageRegistry interface {
	// Read operations
	GetPackageByPath(path string) *ExternalPackage
	GetPackageByName(name string) *ExternalPackage
	GetAllPackages() map[string]*ExternalPackage
	LookupPackage(name string) (*ExternalPackage, error)
	AutoImport(name string) (path string, pkg *ExternalPackage, ok bool)
	LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool)
	LookupMethodDirectCall(typeName, methodName string) (directCall func([]value.Value) value.Value, ok bool)
	LookupExternalVar(pkgPath, varName string) (ptr any, ok bool)
	LookupExternalType(t types.Type) (reflect.Type, bool)
	LookupExternalTypeByName(pkgPath, typeName string) (reflect.Type, bool)

	// Write operations
	RegisterPackage(path, name string) *ExternalPackage
	SetExternalType(t types.Type, rt reflect.Type)
	AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value)
}

// Registry is the concrete implementation of PackageRegistry.
// It stores all registered packages, external type mappings, and method DirectCall wrappers.
//
// All mutable state is protected by a single RWMutex. Registration happens at init
// time; reads happen during compilation — there's no contention benefit from separate locks.
type Registry struct {
	mu                     sync.RWMutex
	packagesByName         map[string]*ExternalPackage                // keyed by package path
	packagesByAlias        map[string]*ExternalPackage                // keyed by package name (for auto-import)
	extTypes               map[types.Type]reflect.Type                // types.Type -> reflect.Type
	methods                map[string]func([]value.Value) value.Value // "pkgPath.TypeName.MethodName" -> DirectCall
	interfaceProxiesByName map[string]*external.InterfaceProxyInfo    // "pkgPath.TypeName" -> proxy metadata
	interfaceProxiesByType map[reflect.Type]*external.InterfaceProxyInfo
}

// NewRegistry creates a new empty package registry.
func NewRegistry() *Registry {
	return &Registry{
		packagesByName:         make(map[string]*ExternalPackage),
		packagesByAlias:        make(map[string]*ExternalPackage),
		extTypes:               make(map[types.Type]reflect.Type),
		methods:                make(map[string]func([]value.Value) value.Value),
		interfaceProxiesByName: make(map[string]*external.InterfaceProxyInfo),
		interfaceProxiesByType: make(map[reflect.Type]*external.InterfaceProxyInfo),
	}
}

func (r *Registry) RegisterPackage(path, name string) *ExternalPackage {
	pkg := &ExternalPackage{
		Path:             path,
		Name:             name,
		Objects:          make(map[string]*external.ExternalObject),
		Types:            make(map[string]reflect.Type),
		InterfaceProxies: make(map[string]*external.InterfaceProxyInfo),
		registry:         r,
	}
	r.mu.Lock()
	for _, info := range r.interfaceProxiesByName {
		if info.PkgPath == path {
			pkg.InterfaceProxies[info.Name] = info
		}
	}
	r.packagesByName[path] = pkg
	r.packagesByAlias[name] = pkg
	r.mu.Unlock()
	return pkg
}

func (r *Registry) GetPackageByPath(path string) *ExternalPackage {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.packagesByName[path]
}

func (r *Registry) GetPackageByName(name string) *ExternalPackage {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.packagesByAlias[name]
}

func (r *Registry) GetAllPackages() map[string]*ExternalPackage {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]*ExternalPackage, len(r.packagesByName))
	for k, v := range r.packagesByName {
		result[k] = v
	}
	return result
}

func (r *Registry) SetExternalType(t types.Type, rt reflect.Type) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.extTypes[t] = rt
}

func (r *Registry) GetExternalType(t types.Type) reflect.Type {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.extTypes[t]
}

func (r *Registry) AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.methods[typeName+"."+methodName] = dc
}

func (r *Registry) LookupMethodDirectCall(typeName, methodName string) (func([]value.Value) value.Value, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	dc, ok := r.methods[typeName+"."+methodName]
	return dc, ok
}

func (r *Registry) AddInterfaceProxy(pkgPath, typeName string, ifaceType reflect.Type, requiredMethods []string, factory external.InterfaceProxyFactory) {
	if ifaceType == nil || factory == nil {
		return
	}
	info := &external.InterfaceProxyInfo{
		PkgPath:         pkgPath,
		Name:            typeName,
		InterfaceType:   ifaceType,
		RequiredMethods: append([]string(nil), requiredMethods...),
		Factory:         factory,
	}
	key := pkgPath + "." + typeName

	r.mu.Lock()
	defer r.mu.Unlock()
	r.interfaceProxiesByName[key] = info
	r.interfaceProxiesByType[ifaceType] = info
	if pkg := r.packagesByName[pkgPath]; pkg != nil {
		pkg.InterfaceProxies[typeName] = info
	}
}

func (r *Registry) LookupInterfaceProxy(pkgPath, typeName string) (*external.InterfaceProxyInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.interfaceProxiesByName[pkgPath+"."+typeName]
	return info, ok
}

func (r *Registry) LookupInterfaceProxyByType(ifaceType reflect.Type) (*external.InterfaceProxyInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.interfaceProxiesByType[ifaceType]
	return info, ok
}

func (r *Registry) LookupInterfaceProxyByInterface(iface *types.Interface) (*external.InterfaceProxyInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, info := range r.interfaceProxiesByName {
		if sameInterfaceMethodSet(info.InterfaceType, iface) {
			return info, true
		}
	}
	return nil, false
}

func sameInterfaceMethodSet(rt reflect.Type, iface *types.Interface) bool {
	if rt == nil || rt.Kind() != reflect.Interface || iface == nil || rt.NumMethod() != iface.NumMethods() {
		return false
	}

	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		reflectMethod, ok := rt.MethodByName(method.Name())
		if !ok {
			return false
		}
		methodType, ok := convertReflectType(reflectMethod.Type).(*types.Signature)
		if !ok || !types.Identical(method.Type(), methodType) {
			return false
		}
	}
	return true
}

func (r *Registry) LookupPackage(name string) (*ExternalPackage, error) {
	if pkg := r.GetPackageByPath(name); pkg != nil {
		return pkg, nil
	}
	if pkg := r.GetPackageByName(name); pkg != nil {
		return pkg, nil
	}
	return nil, fmt.Errorf("package %q not found", name)
}

func (r *Registry) AutoImport(name string) (path string, pkg *ExternalPackage, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p := r.packagesByAlias[name]; p != nil {
		return p.Path, p, true
	}
	for p, pkg := range r.packagesByName {
		if idx := strings.LastIndex(p, "/"); idx >= 0 {
			if p[idx+1:] == name {
				return p, pkg, true
			}
		} else if p == name {
			return p, pkg, true
		}
	}
	return "", nil, false
}

// LookupExternalFunc looks up an external function by package path and function name.
func (r *Registry) LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool) {
	pkg := r.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, nil, false
	}
	obj, exists := pkg.Objects[funcName]
	if !exists || obj.Kind != external.ObjectKindFunction {
		return nil, nil, false
	}
	return obj.Value, obj.DirectCall, true
}

// LookupExternalVar looks up an external variable by package path and variable name.
func (r *Registry) LookupExternalVar(pkgPath, varName string) (ptr any, ok bool) {
	pkg := r.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, false
	}
	obj, exists := pkg.Objects[varName]
	if !exists || obj.Kind != external.ObjectKindVariable {
		return nil, false
	}
	return obj.Value, true
}

// LookupExternalType looks up an external type by types.Type.
func (r *Registry) LookupExternalType(t types.Type) (reflect.Type, bool) {
	rt := r.GetExternalType(t)
	if rt != nil {
		return rt, true
	}
	return nil, false
}

// LookupExternalTypeByName looks up an external type by package path and type name.
func (r *Registry) LookupExternalTypeByName(pkgPath, typeName string) (reflect.Type, bool) {
	pkg := r.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, false
	}
	rt, ok := pkg.Types[typeName]
	return rt, ok
}

// --- Global registry (backward compatibility) ---

// globalRegistry is the default registry used by global convenience functions.
var globalRegistry = NewRegistry() //nolint:gochecknoglobals // backward compatibility

// GlobalRegistry returns the default global registry.
// This registry is pre-populated by init() functions in generated package wrappers.
func GlobalRegistry() PackageRegistry {
	return globalRegistry
}

// RegisterPackage registers a new external package with the global registry.
// This is a convenience function for init() functions in generated package wrappers.
func RegisterPackage(path, name string) *ExternalPackage {
	return globalRegistry.RegisterPackage(path, name)
}

// GetPackageByPath returns a registered package by import path from the global registry.
func GetPackageByPath(path string) *ExternalPackage {
	return globalRegistry.GetPackageByPath(path)
}

// GetPackageByName returns a registered package by name from the global registry.
func GetPackageByName(name string) *ExternalPackage {
	return globalRegistry.GetPackageByName(name)
}

// GetAllPackages returns all registered packages from the global registry.
func GetAllPackages() map[string]*ExternalPackage {
	return globalRegistry.GetAllPackages()
}

// LookupMethodDirectCall looks up a method DirectCall wrapper from the global registry.
func LookupMethodDirectCall(typeName, methodName string) (func([]value.Value) value.Value, bool) {
	return globalRegistry.LookupMethodDirectCall(typeName, methodName)
}

// AddInterfaceProxy registers native proxy metadata on the global registry.
func AddInterfaceProxy(pkgPath, typeName string, ifaceType reflect.Type, requiredMethods []string, factory external.InterfaceProxyFactory) {
	globalRegistry.AddInterfaceProxy(pkgPath, typeName, ifaceType, requiredMethods, factory)
}

// LookupInterfaceProxy looks up native proxy metadata on the global registry.
func LookupInterfaceProxy(pkgPath, typeName string) (*external.InterfaceProxyInfo, bool) {
	return globalRegistry.LookupInterfaceProxy(pkgPath, typeName)
}

// LookupInterfaceProxyByType looks up native proxy metadata by interface reflect.Type.
func LookupInterfaceProxyByType(ifaceType reflect.Type) (*external.InterfaceProxyInfo, bool) {
	return globalRegistry.LookupInterfaceProxyByType(ifaceType)
}

// LookupInterfaceProxyByInterface looks up native proxy metadata by exact method set.
func LookupInterfaceProxyByInterface(iface *types.Interface) (*external.InterfaceProxyInfo, bool) {
	return globalRegistry.LookupInterfaceProxyByInterface(iface)
}

// funcSignature creates a types.Signature from a function value using reflection.
func funcSignature(fn any) *types.Signature {
	rt := reflect.TypeOf(fn)
	if rt.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected function, got %v", rt.Kind()))
	}
	return typeOf(rt).(*types.Signature)
}
