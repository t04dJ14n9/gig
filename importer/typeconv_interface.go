package importer

import (
	"go/types"
	"reflect"
)

// convertInterfaceType converts a reflect.Interface type to a types.Interface.
// For empty interfaces (any), returns types.NewInterfaceType(nil, nil).
func convertInterfaceType(rt reflect.Type) *types.Interface {
	return convertInterfaceTypeWithCache(rt, true)
}

func convertInterfaceTypeNoCache(rt reflect.Type) *types.Interface {
	return convertInterfaceTypeWithCache(rt, false)
}

func convertInterfaceTypeWithCache(rt reflect.Type, cache bool) *types.Interface {
	if rt.NumMethod() == 0 {
		// Empty interface (any)
		return types.NewInterfaceType(nil, nil)
	}

	// Create a placeholder interface and cache it first to break recursion
	iface := types.NewInterfaceType(nil, nil)
	if cache {
		typeCache.Store(rt, iface)
	}

	var methods []*types.Func
	for i := 0; i < rt.NumMethod(); i++ {
		method := rt.Method(i)
		sig := convertFuncType(method.Type)
		methods = append(methods, types.NewFunc(0, nil, method.Name, sig))
	}

	// Update the interface with the methods
	// Note: types.NewInterfaceType creates a complete interface, so we need to
	// create a new one with the methods
	result := types.NewInterfaceType(methods, nil)
	if cache {
		typeCache.Store(rt, result)
	}
	return result
}
