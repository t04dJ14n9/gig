package importer

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"sync"

	"git.woa.com/youngjin/gig/value"
)

// ObjectKind represents the kind of external object (function, variable, constant, or type).
type ObjectKind int

const (
	// ObjectKindInvalid indicates an invalid or uninitialized object.
	ObjectKindInvalid ObjectKind = iota

	// ObjectKindFunction indicates a function object.
	ObjectKindFunction

	// ObjectKindVariable indicates a mutable variable object.
	ObjectKindVariable

	// ObjectKindConstant indicates an immutable constant object.
	ObjectKindConstant

	// ObjectKindType indicates a named type object.
	ObjectKindType
)

// ExternalObject represents a function, variable, constant, or type from an external package.
// It stores the Go value and type information needed for the interpreter.
type ExternalObject struct {
	// Name is the object's identifier (e.g., "Sprintf", "NoError").
	Name string

	// Kind indicates whether this is a function, variable, constant, or type.
	Kind ObjectKind

	// Value is the Go value:
	//   - Function: the function value
	//   - Variable: pointer to the variable
	//   - Constant: the constant value
	//   - Type: zero value of the type
	Value any

	// Type is the Go types.Type representation.
	Type types.Type

	// Doc is optional documentation text.
	Doc string

	DirectCall func([]value.Value) value.Value
}

// ExternalPackage represents a registered external package.
// It maps package-level objects by name for lookup during type checking and execution.
type ExternalPackage struct {
	// Path is the import path (e.g., "fmt", "encoding/json").
	Path string

	// Name is the package identifier (e.g., "fmt", "json").
	Name string

	// Objects maps object names to their ExternalObject entries.
	Objects map[string]*ExternalObject

	// Types maps type names to their reflect.Type representations.
	Types map[string]reflect.Type

	// registry is a back-reference to the owning registry.
	registry PackageRegistry
}

// PackageRegistry manages external package registration.
// It provides methods to register, lookup, and query packages, types, and method DirectCalls.
type PackageRegistry interface {
	// RegisterPackage registers a new external package with the given import path and name.
	RegisterPackage(path, name string) *ExternalPackage
	GetPackageByPath(path string) *ExternalPackage
	GetPackageByName(name string) *ExternalPackage
	GetAllPackages() map[string]*ExternalPackage

	// Type management
	SetExternalType(t types.Type, rt reflect.Type)
	GetExternalType(t types.Type) reflect.Type

	// Method DirectCall management
	AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value)
	LookupMethodDirectCall(typeName, methodName string) (func([]value.Value) value.Value, bool)

	// Package lookup helpers
	LookupPackage(name string) (*ExternalPackage, error)
	AutoImport(name string) (path string, pkg *ExternalPackage, ok bool)
}

// Registry is the concrete implementation of PackageRegistry.
// It stores all registered packages, external type mappings, and method DirectCall wrappers.
type Registry struct {
	mu              sync.RWMutex
	packagesByName  map[string]*ExternalPackage // keyed by package path
	packagesByAlias map[string]*ExternalPackage // keyed by package name (for auto-import)

	typesMu  sync.RWMutex
	extTypes map[types.Type]reflect.Type // types.Type -> reflect.Type

	methodsMu sync.RWMutex
	methods   map[string]func([]value.Value) value.Value // "pkgPath.TypeName.MethodName" -> DirectCall
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

func (r *Registry) RegisterPackage(path, name string) *ExternalPackage {
	pkg := &ExternalPackage{
		Path:     path,
		Name:     name,
		Objects:  make(map[string]*ExternalObject),
		Types:    make(map[string]reflect.Type),
		registry: r,
	}
	r.mu.Lock()
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
	r.typesMu.Lock()
	defer r.typesMu.Unlock()
	r.extTypes[t] = rt
}

func (r *Registry) GetExternalType(t types.Type) reflect.Type {
	r.typesMu.RLock()
	defer r.typesMu.RUnlock()
	return r.extTypes[t]
}

func (r *Registry) AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value) {
	r.methodsMu.Lock()
	defer r.methodsMu.Unlock()
	r.methods[typeName+"."+methodName] = dc
}

func (r *Registry) LookupMethodDirectCall(typeName, methodName string) (func([]value.Value) value.Value, bool) {
	r.methodsMu.RLock()
	defer r.methodsMu.RUnlock()
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
	if pkg := r.GetPackageByName(name); pkg != nil {
		return pkg.Path, pkg, true
	}
	for p, pkg := range r.GetAllPackages() {
		parts := strings.Split(p, "/")
		if parts[len(parts)-1] == name {
			return p, pkg, true
		}
	}
	return "", nil, false
}

// --- ExternalPackage methods ---

// AddFunction adds a function to the package.
func (p *ExternalPackage) AddFunction(name string, fn any, doc string, directCall func([]value.Value) value.Value) {
	sig := funcSignature(fn)
	p.Objects[name] = &ExternalObject{
		Name:       name,
		Kind:       ObjectKindFunction,
		Value:      fn,
		Type:       sig,
		Doc:        doc,
		DirectCall: directCall,
	}
}

// AddVariable adds a variable to the package.
func (p *ExternalPackage) AddVariable(name string, ptr any, doc string) {
	typ := typeOf(reflect.TypeOf(ptr).Elem())
	p.Objects[name] = &ExternalObject{
		Name:  name,
		Kind:  ObjectKindVariable,
		Value: ptr,
		Type:  typ,
		Doc:   doc,
	}
}

// AddConstant adds a constant to the package.
func (p *ExternalPackage) AddConstant(name string, val any, doc string) {
	typ := typeOf(reflect.TypeOf(val))
	p.Objects[name] = &ExternalObject{
		Name:  name,
		Kind:  ObjectKindConstant,
		Value: val,
		Type:  typ,
		Doc:   doc,
	}
}

// AddType adds a named type to the package.
func (p *ExternalPackage) AddType(name string, typ reflect.Type, doc string) {
	if typ == nil {
		return
	}
	p.Types[name] = typ
	p.Objects[name] = &ExternalObject{
		Name:  name,
		Kind:  ObjectKindType,
		Value: reflect.Zero(typ).Interface(),
		Type:  typeOf(typ),
		Doc:   doc,
	}
}

// AddMethodDirectCall registers a DirectCall wrapper for a method on a type in this package.
// It uses the package's owning registry instance.
func (p *ExternalPackage) AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value) {
	if p.registry != nil {
		p.registry.AddMethodDirectCall(p.Path+"."+typeName, methodName, dc)
	}
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

// typeOf is a function that converts reflect.Type to types.Type.
// It is initialized by the importer package's init function.
var typeOf func(reflect.Type) types.Type
