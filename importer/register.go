// register.go defines the external package registry model.
package importer

import (
	"go/types"
	"reflect"
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
