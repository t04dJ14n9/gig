package importer

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

func (r *Registry) SetExternalType(t types.Type, rt reflect.Type) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.extTypes[t] = rt
}

func (r *Registry) GetExternalType(t types.Type) reflect.Type {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.extTypes[t]
}

func (r *Registry) AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.methods[typeName+"."+methodName] = dc
}

func (r *Registry) LookupMethodDirectCall(typeName, methodName string) (func([]value.Value) value.Value, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	dc, ok := r.methods[typeName+"."+methodName]
	return dc, ok
}

// LookupExternalFunc looks up an external function by package path and function name.
func (r *Registry) LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool) {
	pkg := r.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, nil, false
	}
	obj, exists := pkg.Objects[funcName]
	if !exists || obj.Kind != external.ObjectKindFunction {
		return nil, nil, false
	}
	return obj.Value, obj.DirectCall, true
}

// LookupExternalVar looks up an external variable by package path and variable name.
func (r *Registry) LookupExternalVar(pkgPath, varName string) (ptr any, ok bool) {
	pkg := r.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, false
	}
	obj, exists := pkg.Objects[varName]
	if !exists || obj.Kind != external.ObjectKindVariable {
		return nil, false
	}
	return obj.Value, true
}

// LookupExternalType looks up an external type by types.Type.
func (r *Registry) LookupExternalType(t types.Type) (reflect.Type, bool) {
	rt := r.GetExternalType(t)
	if rt != nil {
		return rt, true
	}
	return nil, false
}

// LookupExternalTypeByName looks up an external type by package path and type name.
func (r *Registry) LookupExternalTypeByName(pkgPath, typeName string) (reflect.Type, bool) {
	pkg := r.GetPackageByPath(pkgPath)
	if pkg == nil {
		return nil, false
	}
	rt, ok := pkg.Types[typeName]
	return rt, ok
}
