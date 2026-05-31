package importer

import (
	"go/types"
	"reflect"
)

// addMethodsToNamed adds methods from a reflect.Type to a types.Named type.
// This allows the type checker to find methods on external types.
// Both value receiver and pointer receiver methods are added.
// For interface types, methods do NOT include a receiver parameter.
func addMethodsToNamed(named *types.Named, rt reflect.Type) {
	// Check if this is an interface type - interface methods don't have receiver params
	isInterface := rt.Kind() == reflect.Interface

	// Enumerate all exported methods on the value receiver
	addMethodsFromType(named, rt, isInterface, false)

	// For interface types, don't process pointer receiver methods
	if isInterface {
		return
	}

	// Also enumerate methods on the pointer receiver (*T)
	addMethodsFromType(named, reflect.PointerTo(rt), false, true)
}

// addMethodsFromType adds exported methods from methodSource to the named type.
// If skipReceiver is false (interface types), parameters start at index 0.
// If isPointerRecv is true, skips methods already present on named and uses pointer receiver.
func addMethodsFromType(named *types.Named, methodSource reflect.Type, isInterface, isPointerRecv bool) {
	for i := 0; i < methodSource.NumMethod(); i++ {
		method := methodSource.Method(i)
		if !method.IsExported() {
			continue
		}
		// For pointer receiver pass, skip methods already added from value receiver
		if isPointerRecv {
			alreadyAdded := false
			for j := 0; j < named.NumMethods(); j++ {
				if named.Method(j).Name() == method.Name {
					alreadyAdded = true
					break
				}
			}
			if alreadyAdded {
				continue
			}
		}

		methodType := method.Type
		if methodType.Kind() != reflect.Func {
			continue
		}
		// For non-interface types, we need at least 1 param (the receiver)
		if !isInterface && methodType.NumIn() < 1 {
			continue
		}

		// Build the parameter list
		// For interfaces, start at 0 (no receiver). For concrete types, skip receiver at 0.
		startIdx := 0
		if !isInterface {
			startIdx = 1
		}
		var params []*types.Var
		for j := startIdx; j < methodType.NumIn(); j++ {
			paramType := convertReflectType(methodType.In(j))
			params = append(params, types.NewVar(0, nil, "", paramType))
		}

		// Build the result list
		var results []*types.Var
		for j := 0; j < methodType.NumOut(); j++ {
			resultType := convertReflectType(methodType.Out(j))
			results = append(results, types.NewVar(0, nil, "", resultType))
		}

		// Build receiver: nil for interfaces, named for value, *named for pointer
		var recv *types.Var
		if !isInterface {
			recvType := types.Type(named)
			if isPointerRecv {
				recvType = types.NewPointer(named)
			}
			recv = types.NewVar(0, nil, "", recvType)
		}

		sig := types.NewSignatureType(
			recv,
			nil, nil,
			types.NewTuple(params...),
			types.NewTuple(results...),
			methodType.IsVariadic(),
		)

		fn := types.NewFunc(0, named.Obj().Pkg(), method.Name, sig)
		named.AddMethod(fn)
	}
}
