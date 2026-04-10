// register.go implements the global registry: package registration, type mapping,
// method DirectCall lookup, and auto-import.
package importer

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"sync"

	"git.woa.com/youngjin/gig/model/external"
	"git.woa.com/youngjin/gig/model/value"
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

	// registry is a back-reference to the owning registry.
	registry PackageRegistry
}

// PackageRegistry manages external package registration.
// It provides methods to register, lookup, and query packages, types, and method DirectCalls.
// PackageRegistry embeds PackageLookup for read operations and adds write operations.
type PackageRegistry interface {
	PackageLookup // Embed read interface

	// Registration (write operations)
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
	mu              sync.RWMutex
	packagesByName  map[string]*ExternalPackage                // keyed by package path
	packagesByAlias map[string]*ExternalPackage                // keyed by package name (for auto-import)
	extTypes        map[types.Type]reflect.Type                // types.Type -> reflect.Type
	methods         map[string]func([]value.Value) value.Value // "pkgPath.TypeName.MethodName" -> DirectCall

	frozen bool // if true, mutations panic
}

// NewRegistry creates a new empty package registry.
func NewRegistry() *Registry {
	return &Registry{
		packagesByName:  make(map[string]*ExternalPackage),
		packagesByAlias: make(map[string]*ExternalPackage),
		extTypes:        make(map[types.Type]reflect.Type),
		methods:         make(map[string]func([]value.Value) value.Value),
	}
}

// Freeze makes the registry read-only. Any subsequent mutation (RegisterPackage,
// SetExternalType, AddMethodDirectCall) will panic. This is called after all stdlib
// init() registrations are complete to prevent cross-program pollution.
func (r *Registry) Freeze() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.frozen = true
}

// IsFrozen reports whether the registry has been frozen.
func (r *Registry) IsFrozen() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.frozen
}

// checkFrozen panics if the registry is frozen. Safe to call under any lock
// since frozen transitions from false→true exactly once and never back.
func (r *Registry) checkFrozen() {
	if r.frozen {
		panic("importer: mutation on frozen registry; the global registry is read-only after init")
	}
}

func (r *Registry) RegisterPackage(path, name string) *ExternalPackage {
	pkg := &ExternalPackage{
		Path:     path,
		Name:     name,
		Objects:  make(map[string]*external.ExternalObject),
		Types:    make(map[string]reflect.Type),
		registry: r,
	}
	r.mu.Lock()
	r.checkFrozen()
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
	r.checkFrozen()
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
	r.checkFrozen()
	r.methods[typeName+"."+methodName] = dc
}

func (r *Registry) LookupMethodDirectCall(typeName, methodName string) (func([]value.Value) value.Value, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	dc, ok := r.methods[typeName+"."+methodName]
	return dc, ok
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

// SetExternalType associates a types.Type with a reflect.Type in the global registry.
func SetExternalType(t types.Type, rt reflect.Type) {
	globalRegistry.SetExternalType(t, rt)
}

// GetExternalType returns the reflect.Type associated with a types.Type from the global registry.
func GetExternalType(t types.Type) reflect.Type {
	return globalRegistry.GetExternalType(t)
}

// funcSignature creates a types.Signature from a function value using reflection.
func funcSignature(fn any) *types.Signature {
	rt := reflect.TypeOf(fn)
	if rt.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected function, got %v", rt.Kind()))
	}
	return typeOf(rt).(*types.Signature)
}
