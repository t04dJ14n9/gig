// Package importer provides type information and package registration for the interpreter.
// This file implements the types.Importer interface and type conversion utilities.
package importer

import (
	"fmt"
	"go/types"
	"sync"

	"git.woa.com/youngjin/gig/model/external"
)

// Importer implements types.Importer for registered external packages.
// It allows the Go type checker to resolve imports to registered packages.
type Importer struct {
	// reg is the package registry to use for resolution.
	reg PackageRegistry

	// packages caches resolved types.Package values.
	packages map[string]*types.Package

	// mutex protects concurrent access to packages.
	mutex sync.RWMutex
}

// NewImporter creates a new Importer for resolving registered packages.
func NewImporter(reg PackageRegistry) *Importer {
	return &Importer{
		reg:      reg,
		packages: make(map[string]*types.Package),
	}
}

// Import returns the types.Package for the given import path.
// It first checks the cache, then looks up registered external packages.
// Returns an error if the package is not registered.
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
	extPkg := i.reg.GetPackageByPath(path)
	if extPkg == nil {
		// Try to find by name (for auto-imported packages)
		extPkg = i.reg.GetPackageByName(path)
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

// buildPackage creates a types.Package from a registered ExternalPackage.
// It converts all objects (functions, variables, constants, types) to types.Object.
func (i *Importer) buildPackage(extPkg *ExternalPackage) *types.Package {
	pkg := types.NewPackage(extPkg.Path, extPkg.Name)

	// Add all objects to the package scope
	for name, obj := range extPkg.Objects {
		var typesObj types.Object

		switch obj.Kind {
		case external.ObjectKindFunction:
			if sig, ok := obj.Type.(*types.Signature); ok {
				typesObj = types.NewFunc(0, pkg, name, sig)
			}
		case external.ObjectKindVariable:
			typesObj = types.NewVar(0, pkg, name, obj.Type)
		case external.ObjectKindConstant:
			// Create a constant with the appropriate value
			val := convertToConstantValue(obj.Value)
			typesObj = types.NewConst(0, pkg, name, obj.Type, val)
		case external.ObjectKindType:
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
		i.reg.SetExternalType(t, rt)
	}

	// Mark the package as complete so the type checker can use it
	pkg.MarkComplete()

	return pkg
}

