// Package importer provides type information and package registration for the interpreter.
package importer

import (
	"fmt"
	"go/types"
	"reflect"
	"sync"

	"gig/value"
)

// ObjectKind represents the kind of external object.
type ObjectKind int

const (
	ObjectKindInvalid ObjectKind = iota
	ObjectKindFunction
	ObjectKindVariable
	ObjectKindConstant
	ObjectKindType
)

// ExternalObject represents a function, variable, constant, or type from an external package.
type ExternalObject struct {
	Name       string                          // Object name
	Kind       ObjectKind                      // Kind of object
	Value      any                             // Go value (function, variable pointer, constant value)
	Type       types.Type                      // Go types.Type
	Doc        string                          // Documentation
	DirectCall func([]value.Value) value.Value // Typed wrapper for function calls (avoids reflect.Call)
}

// ExternalPackage represents a registered external package.
type ExternalPackage struct {
	Path    string                     // Import path
	Name    string                     // Package name
	Objects map[string]*ExternalObject // Objects by name
	Types   map[string]reflect.Type    // reflect.Type by type name
}

// Global registry
var (
	packagesMutex   sync.RWMutex
	packagesByName  = make(map[string]*ExternalPackage) // keyed by package path
	packagesByAlias = make(map[string]*ExternalPackage) // keyed by package name (for auto-import)
	typesMutex      sync.RWMutex
	externalTypes   = make(map[types.Type]reflect.Type) // types.Type -> reflect.Type
)

// RegisterPackage registers a new external package.
func RegisterPackage(path, name string) *ExternalPackage {
	pkg := &ExternalPackage{
		Path:    path,
		Name:    name,
		Objects: make(map[string]*ExternalObject),
		Types:   make(map[string]reflect.Type),
	}

	packagesMutex.Lock()
	packagesByName[path] = pkg
	packagesByAlias[name] = pkg
	packagesMutex.Unlock()

	return pkg
}

// GetPackageByPath returns a registered package by its import path.
func GetPackageByPath(path string) *ExternalPackage {
	packagesMutex.RLock()
	defer packagesMutex.RUnlock()
	return packagesByName[path]
}

// GetPackageByName returns a registered package by its name (for auto-import).
func GetPackageByName(name string) *ExternalPackage {
	packagesMutex.RLock()
	defer packagesMutex.RUnlock()
	return packagesByAlias[name]
}

// GetAllPackages returns all registered packages.
func GetAllPackages() map[string]*ExternalPackage {
	packagesMutex.RLock()
	defer packagesMutex.RUnlock()

	result := make(map[string]*ExternalPackage, len(packagesByName))
	for k, v := range packagesByName {
		result[k] = v
	}
	return result
}

// AddFunction adds a function to a package.
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

// AddVariable adds a variable to a package.
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

// AddConstant adds a constant to a package.
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

// AddType adds a type to a package.
func (p *ExternalPackage) AddType(name string, typ reflect.Type, doc string) {
	p.Types[name] = typ
	p.Objects[name] = &ExternalObject{
		Name:  name,
		Kind:  ObjectKindType,
		Value: reflect.Zero(typ).Interface(),
		Type:  typeOf(typ),
		Doc:   doc,
	}
}

// SetExternalType associates a types.Type with a reflect.Type.
func SetExternalType(t types.Type, rt reflect.Type) {
	typesMutex.Lock()
	defer typesMutex.Unlock()
	externalTypes[t] = rt
}

// GetExternalType returns the reflect.Type for a types.Type.
func GetExternalType(t types.Type) reflect.Type {
	typesMutex.RLock()
	defer typesMutex.RUnlock()
	return externalTypes[t]
}

// funcSignature creates a types.Signature from a function value.
func funcSignature(fn any) *types.Signature {
	rt := reflect.TypeOf(fn)
	if rt.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected function, got %v", rt.Kind()))
	}

	return typeOf(rt).(*types.Signature)
}

// typeOf converts a reflect.Type to types.Type.
// This is implemented in importer.go.
var typeOf func(reflect.Type) types.Type
