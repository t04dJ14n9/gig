// interfaces.go defines PackageLookup and PackageRegistry interfaces for dependency injection.
package importer

import (
	"go/types"
	"reflect"

	"git.woa.com/youngjin/gig/model/external"
	"git.woa.com/youngjin/gig/model/value"
)

// PackageLookup provides read-only access to registered packages.
// This interface is used by the compiler and VM for lookups.
type PackageLookup interface {
	// Package queries
	GetPackageByPath(path string) *ExternalPackage
	GetPackageByName(name string) *ExternalPackage
	GetAllPackages() map[string]*ExternalPackage
	LookupPackage(name string) (*ExternalPackage, error)
	AutoImport(name string) (path string, pkg *ExternalPackage, ok bool)

	// External function/method lookup
	LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool)
	LookupMethodDirectCall(typeName, methodName string) (directCall func([]value.Value) value.Value, ok bool)

	// External variable lookup
	// The returned value should be a pointer to the variable (e.g., &time.UTC).
	LookupExternalVar(pkgPath, varName string) (ptr any, ok bool)

	// External type lookup
	// Returns the real reflect.Type for an external named type.
	LookupExternalType(t types.Type) (reflect.Type, bool)

	// LookupExternalTypeByName resolves an external type by package path and type name.
	// Returns the real reflect.Type and whether it was found.
	LookupExternalTypeByName(pkgPath, typeName string) (reflect.Type, bool)
}

// NewPackageLookup creates a PackageLookup backed by the given registry.
func NewPackageLookup(reg PackageRegistry) PackageLookup {
	return &registryLookup{reg: reg}
}

// registryLookup implements PackageLookup by delegating to a PackageRegistry.
type registryLookup struct {
	reg PackageRegistry
}

func (l *registryLookup) GetPackageByPath(path string) *ExternalPackage {
	return l.reg.GetPackageByPath(path)
}

func (l *registryLookup) GetPackageByName(name string) *ExternalPackage {
	return l.reg.GetPackageByName(name)
}

func (l *registryLookup) GetAllPackages() map[string]*ExternalPackage {
	return l.reg.GetAllPackages()
}

func (l *registryLookup) LookupPackage(name string) (*ExternalPackage, error) {
	return l.reg.LookupPackage(name)
}

func (l *registryLookup) AutoImport(name string) (path string, pkg *ExternalPackage, ok bool) {
	return l.reg.AutoImport(name)
}

func (l *registryLookup) LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool) {
	pkg := l.reg.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, nil, false
	}
	obj, exists := pkg.Objects[funcName]
	if !exists || obj.Kind != external.ObjectKindFunction {
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
	if !exists || obj.Kind != external.ObjectKindVariable {
		return nil, false
	}
	return obj.Value, true
}

func (l *registryLookup) LookupExternalType(t types.Type) (reflect.Type, bool) {
	return l.reg.LookupExternalType(t)
}

func (l *registryLookup) LookupExternalTypeByName(pkgPath, typeName string) (reflect.Type, bool) {
	pkg := l.reg.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, false
	}
	rt, ok := pkg.Types[typeName]
	return rt, ok
}
