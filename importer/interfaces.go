package importer

import (
	"go/types"
	"reflect"

	"git.woa.com/youngjin/gig/value"
)

// PackageLookup resolves external package functions for the compiler.
// This interface enables dependency injection: the compiler depends on this
// abstraction rather than importing the importer package directly.
type PackageLookup interface {
	LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool)
	LookupMethodDirectCall(typeName, methodName string) (directCall func([]value.Value) value.Value, ok bool)
	// LookupExternalVar returns the value of an external package variable.
	// The returned value should be a pointer to the variable (e.g., &time.UTC).
	LookupExternalVar(pkgPath, varName string) (ptr any, ok bool)
	// LookupExternalType returns the real reflect.Type for an external named type.
	// This is used to allocate real Go types (e.g., bytes.Buffer) instead of
	// synthesized struct types from reflect.StructOf.
	LookupExternalType(t types.Type) (reflect.Type, bool)
}

// NewPackageLookup creates a PackageLookup backed by the given registry.
func NewPackageLookup(reg PackageRegistry) PackageLookup {
	return &registryLookup{reg: reg}
}

// registryLookup implements PackageLookup by delegating to a PackageRegistry.
type registryLookup struct {
	reg PackageRegistry
}

func (l *registryLookup) LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool) {
	pkg := l.reg.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, nil, false
	}
	obj, exists := pkg.Objects[funcName]
	if !exists || obj.Kind != ObjectKindFunction {
		return nil, nil, false
	}
	return obj.Value, obj.DirectCall, true
}

func (l *registryLookup) LookupMethodDirectCall(typeName, methodName string) (directCall func([]value.Value) value.Value, ok bool) {
	return l.reg.LookupMethodDirectCall(typeName, methodName)
}

func (l *registryLookup) LookupExternalVar(pkgPath, varName string) (ptr any, ok bool) {
	pkg := l.reg.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, false
	}
	obj, exists := pkg.Objects[varName]
	if !exists || obj.Kind != ObjectKindVariable {
		return nil, false
	}
	return obj.Value, true
}

func (l *registryLookup) LookupExternalType(t types.Type) (reflect.Type, bool) {
	rt := l.reg.GetExternalType(t)
	if rt != nil {
		return rt, true
	}
	return nil, false
}
