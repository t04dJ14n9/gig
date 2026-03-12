// Package importer provides type information and package registration for the interpreter.
//
// The importer package is responsible for:
//   - Registering external packages for use in interpreted code
//   - Providing type information for external packages to the type checker
//   - Converting between Go's reflect.Type and types.Type representations
//
// # Package Registration
//
// External packages are registered using RegisterPackage:
//
//	pkg := importer.RegisterPackage("fmt", "fmt")
//	pkg.AddFunction("Sprintf", fmt.Sprintf, "", directCallWrapper)
//	pkg.AddConstant("NoError", "", "")
//
// # Object Types
//
// The importer supports four types of external objects:
//   - Functions: Go functions callable from interpreted code
//   - Variables: Mutable package-level variables
//   - Constants: Compile-time constant values
//   - Types: Named types with methods
//
// # DirectCall Optimization
//
// For frequently-called functions, a DirectCall wrapper can be provided to bypass
// reflection overhead. The wrapper converts value.Value arguments to native Go types,
// calls the function, and wraps the result.
package importer

import (
	"fmt"
	"go/types"
	"reflect"
	"sync"

	"github.com/t04dJ14n9/gig/value"
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

	// DirectCall is an optional typed wrapper that bypasses reflect.Call.
	// If provided, the VM will use this for faster function dispatch.
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
}

// Global registry
var (
	packagesMutex   sync.RWMutex
	packagesByName  = make(map[string]*ExternalPackage) // keyed by package path
	packagesByAlias = make(map[string]*ExternalPackage) // keyed by package name (for auto-import)
	typesMutex      sync.RWMutex
	externalTypes   = make(map[types.Type]reflect.Type) // types.Type -> reflect.Type

	methodDirectCallsMutex sync.RWMutex
	methodDirectCalls      = make(map[string]func([]value.Value) value.Value) // "pkgPath.TypeName.MethodName" -> DirectCall
)

// RegisterPackage registers a new external package with the given import path and name.
// Returns an ExternalPackage for adding functions, variables, constants, and types.
// Packages are registered globally and can be looked up by path or name.
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
// Returns nil if no package with the given path is registered.
func GetPackageByPath(path string) *ExternalPackage {
	packagesMutex.RLock()
	defer packagesMutex.RUnlock()
	return packagesByName[path]
}

// GetPackageByName returns a registered package by its name (for auto-import).
// Returns nil if no package with the given name is registered.
func GetPackageByName(name string) *ExternalPackage {
	packagesMutex.RLock()
	defer packagesMutex.RUnlock()
	return packagesByAlias[name]
}

// GetAllPackages returns a copy of all registered packages, keyed by import path.
func GetAllPackages() map[string]*ExternalPackage {
	packagesMutex.RLock()
	defer packagesMutex.RUnlock()

	result := make(map[string]*ExternalPackage, len(packagesByName))
	for k, v := range packagesByName {
		result[k] = v
	}
	return result
}

// AddFunction adds a function to the package.
// The fn parameter must be a function value.
// The directCall parameter is an optional typed wrapper for fast dispatch.
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
// The ptr parameter must be a pointer to the variable.
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
// The val parameter is the constant value.
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
// The typ parameter is the reflect.Type representation of the type.
func (p *ExternalPackage) AddType(name string, typ reflect.Type, doc string) {
	if typ == nil {
		// Skip nil types (interface placeholders, etc.)
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
func (p *ExternalPackage) AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value) {
	methodDirectCallsMutex.Lock()
	defer methodDirectCallsMutex.Unlock()
	key := typeName + "." + methodName
	methodDirectCalls[key] = dc
}

// LookupMethodDirectCall looks up a method DirectCall wrapper by type name and method name.
func LookupMethodDirectCall(typeName, methodName string) (func([]value.Value) value.Value, bool) {
	methodDirectCallsMutex.RLock()
	defer methodDirectCallsMutex.RUnlock()
	key := typeName + "." + methodName
	dc, ok := methodDirectCalls[key]
	return dc, ok
}

// SetExternalType associates a types.Type with a reflect.Type.
// This mapping is used when the VM needs to allocate or manipulate external types.
func SetExternalType(t types.Type, rt reflect.Type) {
	typesMutex.Lock()
	defer typesMutex.Unlock()
	externalTypes[t] = rt
}

// GetExternalType returns the reflect.Type associated with a types.Type.
// Returns nil if no association exists.
func GetExternalType(t types.Type) reflect.Type {
	typesMutex.RLock()
	defer typesMutex.RUnlock()
	return externalTypes[t]
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
