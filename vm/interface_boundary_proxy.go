package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

type interfaceProxyLookup interface {
	LookupInterfaceProxy(pkgPath, typeName string) (*external.InterfaceProxyInfo, bool)
	LookupInterfaceProxyByType(ifaceType reflect.Type) (*external.InterfaceProxyInfo, bool)
}

func (v *vm) makeInterpretedInterfaceAdapter(targetType, concreteType types.Type, receiver value.Value) (any, bool) {
	receiverTypeName := namedTypeName(concreteType)
	if receiverTypeName == "" {
		return nil, false
	}

	if proxy, ok := v.makeRegisteredInterfaceProxy(targetType, receiverTypeName, receiver); ok {
		return proxy, true
	}

	if isErrorInterface(targetType) && !isStructLikeType(concreteType) {
		return newInterpretedErrorAdapter(
			v.program, receiver, receiverTypeName,
			v.getGlobals(), v.initialGlobals, v.shared, v.ctx, v.goroutines,
		), true
	}

	if !isHostCallbackInterface(targetType) {
		return nil, false
	}
	return newInterpretedInterfaceAdapter(
		v.program, receiver, receiverTypeName,
		v.getGlobals(), v.initialGlobals, v.shared, v.ctx, v.goroutines,
	), true
}

func (v *vm) makeRegisteredInterfaceProxy(targetType types.Type, receiverTypeName string, receiver value.Value) (any, bool) {
	info, ok := v.lookupRegisteredInterfaceProxy(targetType)
	if !ok {
		return nil, false
	}
	// Registered proxies are the explicit third-party escape hatch: the native
	// proxy satisfies the real Go interface while calls are routed back into the
	// interpreted receiver methods.
	call := func(methodName string, args ...value.Value) (value.Value, bool) {
		return callInterfaceMethodValue(
			v.program, methodName, receiverTypeName, receiver, args,
			v.getGlobals(), v.initialGlobals, v.shared, v.ctx, v.goroutines,
		)
	}
	return info.Factory(receiver, receiverTypeName, call)
}

func (v *vm) lookupRegisteredInterfaceProxy(targetType types.Type) (*external.InterfaceProxyInfo, bool) {
	if v == nil || v.program == nil || v.program.TypeResolver == nil {
		return nil, false
	}
	lookup, ok := v.program.TypeResolver.(interfaceProxyLookup)
	if !ok {
		return nil, false
	}
	if named, ok := targetType.(*types.Named); ok {
		obj := named.Obj()
		if obj != nil && obj.Pkg() != nil {
			if info, ok := lookup.LookupInterfaceProxy(obj.Pkg().Path(), obj.Name()); ok {
				return info, true
			}
		}
	}
	if ifaceRT := typeToReflect(targetType, v.program); ifaceRT != nil {
		if info, ok := lookup.LookupInterfaceProxyByType(ifaceRT); ok {
			return info, true
		}
	}
	return nil, false
}
