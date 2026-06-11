package importer

import (
	"go/types"
	"strings"
	"sync"
)

// typePkgCache caches *types.Package objects by package path to ensure
// that the same package path always maps to the same *types.Package instance.
var typePkgCache sync.Map // map[string]*types.Package

// getOrCreateTypesPackage returns a cached *types.Package for the given
// package path, creating one if it doesn't exist yet. The package name is
// derived from the last path segment (e.g., "encoding/json" -> "json").
func getOrCreateTypesPackage(pkgPath string) *types.Package {
	if pkgPath == "" {
		return nil
	}
	if cached, ok := typePkgCache.Load(pkgPath); ok {
		return cached.(*types.Package)
	}
	// Derive package name from path (last segment)
	name := pkgPath
	if idx := strings.LastIndex(pkgPath, "/"); idx >= 0 {
		name = pkgPath[idx+1:]
	}
	pkg := types.NewPackage(pkgPath, name)
	actual, _ := typePkgCache.LoadOrStore(pkgPath, pkg)
	return actual.(*types.Package)
}
