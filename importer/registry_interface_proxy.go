package importer

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/external"
)

func (r *Registry) AddInterfaceProxy(pkgPath, typeName string, ifaceType reflect.Type, requiredMethods []string, factory external.InterfaceProxyFactory) {
	if ifaceType == nil || factory == nil {
		return
	}
	info := &external.InterfaceProxyInfo{
		PkgPath:         pkgPath,
		Name:            typeName,
		InterfaceType:   ifaceType,
		RequiredMethods: append([]string(nil), requiredMethods...),
		Factory:         factory,
	}
	key := pkgPath + "." + typeName

	r.mu.Lock()
	defer r.mu.Unlock()
	r.interfaceProxiesByName[key] = info
	r.interfaceProxiesByType[ifaceType] = info
	if pkg := r.packagesByName[pkgPath]; pkg != nil {
		pkg.InterfaceProxies[typeName] = info
	}
}

func (r *Registry) LookupInterfaceProxy(pkgPath, typeName string) (*external.InterfaceProxyInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.interfaceProxiesByName[pkgPath+"."+typeName]
	return info, ok
}

func (r *Registry) LookupInterfaceProxyByType(ifaceType reflect.Type) (*external.InterfaceProxyInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.interfaceProxiesByType[ifaceType]
	return info, ok
}

func (r *Registry) LookupInterfaceProxyByInterface(iface *types.Interface) (*external.InterfaceProxyInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, info := range r.interfaceProxiesByName {
		if sameInterfaceMethodSet(info.InterfaceType, iface) {
			return info, true
		}
	}
	return nil, false
}

func sameInterfaceMethodSet(rt reflect.Type, iface *types.Interface) bool {
	if rt == nil || rt.Kind() != reflect.Interface || iface == nil || rt.NumMethod() != iface.NumMethods() {
		return false
	}

	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		reflectMethod, ok := rt.MethodByName(method.Name())
		if !ok {
			return false
		}
		methodType, ok := convertReflectType(reflectMethod.Type).(*types.Signature)
		if !ok || !types.Identical(method.Type(), methodType) {
			return false
		}
	}
	return true
}
