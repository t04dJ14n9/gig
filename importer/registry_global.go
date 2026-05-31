package importer

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

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
