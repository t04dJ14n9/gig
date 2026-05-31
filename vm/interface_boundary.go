package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

type interfaceProxyLookup interface {
	LookupInterfaceProxy(pkgPath, typeName string) (*external.InterfaceProxyInfo, bool)
	LookupInterfaceProxyByType(ifaceType reflect.Type) (*external.InterfaceProxyInfo, bool)
}

func (v *vm) assertInterpretedInterfaceValue(dyn *value.InterpretedInterfaceValue, targetType types.Type, original value.Value) (value.Value, bool) {
	if iface, ok := targetType.Underlying().(*types.Interface); ok {
		if iface.NumMethods() == 0 {
			return original, true
		}
		if v.interpretedTypeSatisfiesInterface(dyn, iface) {
			return original, true
		}
		return zeroValueForType(targetType, v.program), false
	}

	if targetName := namedTypeName(targetType); targetName == dyn.TypeName && isPointerType(targetType) == dyn.IsPointer {
		if dyn.IsPointer {
			return dyn.Value, true
		}
		if kindMatchesType(dyn.Value.Kind(), dyn.Value.RawSize(), targetType) {
			return dyn.Value, true
		}
	}

	return zeroValueForType(targetType, v.program), false
}

func (v *vm) interpretedTypeSatisfiesInterface(dyn *value.InterpretedInterfaceValue, iface *types.Interface) bool {
	if dyn == nil || iface == nil || v == nil || v.program == nil {
		return false
	}
	for i := 0; i < iface.NumMethods(); i++ {
		methodName := iface.Method(i).Name()
		found := false
		for _, fn := range v.program.MethodsByName[methodName] {
			if fn.ReceiverTypeName == dyn.TypeName && (!fn.ReceiverIsPointer || dyn.IsPointer) {
				found = true
				break
			}
			if !fn.ReceiverIsPointer || dyn.IsPointer {
				if _, ok := receiverForCompiledMethodTarget(methodName, dyn.Value, fn, v.program); ok {
					found = true
					break
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func zeroValueForType(t types.Type, program *bytecode.CompiledProgram) value.Value {
	if rt := typeToReflect(t, program); rt != nil {
		return value.MakeFromReflect(reflect.Zero(rt))
	}
	return value.MakeNil()
}

func shouldPreserveInterpretedNamedType(t types.Type, program *bytecode.CompiledProgram) bool {
	typeName := namedTypeName(t)
	if typeName == "" {
		return false
	}
	rt := typeToReflect(t, program)
	return rt == nil || rt.Name() != typeName
}

func isPointerType(t types.Type) bool {
	_, ok := t.(*types.Pointer)
	return ok
}

func (v *vm) makeInterpretedInterfaceAdapter(targetType, concreteType types.Type, receiver value.Value) (any, bool) {
	receiverTypeName := namedTypeName(concreteType)
	if receiverTypeName == "" {
		return nil, false
	}

	if info, ok := v.lookupRegisteredInterfaceProxy(targetType); ok {
		call := func(methodName string, args ...value.Value) (value.Value, bool) {
			return callInterfaceMethodValue(
				v.program, methodName, receiverTypeName, receiver, args,
				v.getGlobals(), v.initialGlobals, v.shared, v.ctx, v.goroutines,
			)
		}
		if proxy, ok := info.Factory(receiver, receiverTypeName, call); ok {
			return proxy, true
		}
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

func isErrorInterface(t types.Type) bool {
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		return obj != nil && obj.Pkg() == nil && obj.Name() == "error"
	}
	if iface, ok := t.Underlying().(*types.Interface); ok && iface.NumMethods() == 1 {
		method := iface.Method(0)
		sig, ok := method.Type().(*types.Signature)
		return ok &&
			method.Name() == "Error" &&
			sig.Params().Len() == 0 &&
			sig.Results().Len() == 1 &&
			types.Identical(sig.Results().At(0).Type(), types.Typ[types.String])
	}
	return false
}

func isStructLikeType(t types.Type) bool {
	for {
		if t == nil {
			return false
		}
		switch tt := t.(type) {
		case *types.Named:
			_, ok := tt.Underlying().(*types.Struct)
			return ok
		case *types.Pointer:
			t = tt.Elem()
		default:
			_, ok := tt.Underlying().(*types.Struct)
			return ok
		}
	}
}

// isHostCallbackInterface checks whether the target type is sort.Interface or
// container/heap.Interface — the two stdlib interfaces that receive callbacks
// and for which gig provides interpreted-to-native adapters.
//
// The target may arrive as *types.Named (when the compiler stores the named
// type) or as *types.Interface (when it stores the underlying interface).
// We handle both.
func isHostCallbackInterface(t types.Type) bool {
	// Fast path: *types.Named wrapping sort.Interface or heap.Interface.
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		if obj == nil || obj.Pkg() == nil {
			return false
		}
		pkgPath := obj.Pkg().Path()
		return obj.Name() == "Interface" && (pkgPath == "sort" || pkgPath == "container/heap")
	}
	// Slow path: *types.Interface directly. Match by method signature.
	if iface, ok := t.(*types.Interface); ok {
		return hasInterfaceMethods(iface, "Len", "Less", "Swap")
	}
	return false
}

// hasInterfaceMethods checks that an interface type has methods with exactly
// the given names (order-independent).
func hasInterfaceMethods(iface *types.Interface, names ...string) bool {
	if iface.NumMethods() < len(names) {
		return false
	}
	for _, name := range names {
		found := false
		for i := 0; i < iface.NumMethods(); i++ {
			if iface.Method(i).Name() == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func namedTypeName(t types.Type) string {
	for {
		switch tt := t.(type) {
		case *types.Named:
			return tt.Obj().Name()
		case *types.Pointer:
			t = tt.Elem()
		default:
			return ""
		}
	}
}
