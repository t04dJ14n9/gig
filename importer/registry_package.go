package importer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/t04dJ14n9/gig/model/external"
)

func (r *Registry) RegisterPackage(path, name string) *ExternalPackage {
	pkg := &ExternalPackage{
		Path:             path,
		Name:             name,
		Objects:          make(map[string]*external.ExternalObject),
		Types:            make(map[string]reflect.Type),
		InterfaceProxies: make(map[string]*external.InterfaceProxyInfo),
		registry:         r,
	}
	r.mu.Lock()
	for _, info := range r.interfaceProxiesByName {
		if info.PkgPath == path {
			pkg.InterfaceProxies[info.Name] = info
		}
	}
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
