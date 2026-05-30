package vm

import (
	"reflect"
	"sort"
	"testing"

	"go/types"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

type interfaceProxyTestResolver struct {
	info *external.InterfaceProxyInfo
}

func (r interfaceProxyTestResolver) LookupExternalType(types.Type) (reflect.Type, bool) {
	return r.info.InterfaceType, true
}

func (r interfaceProxyTestResolver) LookupExternalTypeByName(pkgPath, typeName string) (reflect.Type, bool) {
	if pkgPath == r.info.PkgPath && typeName == r.info.Name {
		return r.info.InterfaceType, true
	}
	return nil, false
}

func (r interfaceProxyTestResolver) LookupInterfaceProxy(pkgPath, typeName string) (*external.InterfaceProxyInfo, bool) {
	if pkgPath == r.info.PkgPath && typeName == r.info.Name {
		return r.info, true
	}
	return nil, false
}

func (r interfaceProxyTestResolver) LookupInterfaceProxyByType(ifaceType reflect.Type) (*external.InterfaceProxyInfo, bool) {
	if ifaceType == r.info.InterfaceType {
		return r.info, true
	}
	return nil, false
}

func TestMakeInterpretedInterfaceAdapterUsesRegisteredProxy(t *testing.T) {
	called := false
	info := &external.InterfaceProxyInfo{
		PkgPath:       "sort",
		Name:          "Interface",
		InterfaceType: reflect.TypeOf((*sort.Interface)(nil)).Elem(),
		Factory: func(receiver value.Value, receiverTypeName string, call external.InterfaceMethodCaller) (any, bool) {
			called = true
			if receiverTypeName != "MySorter" {
				t.Fatalf("receiverTypeName = %q, want MySorter", receiverTypeName)
			}
			return sort.IntSlice{}, true
		},
	}

	sortPkg := types.NewPackage("sort", "sort")
	targetType := types.NewNamed(types.NewTypeName(0, sortPkg, "Interface", nil), types.NewInterfaceType(nil, nil), nil)
	userPkg := types.NewPackage("command-line-arguments", "main")
	concreteType := types.NewNamed(types.NewTypeName(0, userPkg, "MySorter", nil), types.NewSlice(types.Typ[types.Int]), nil)

	v := &vm{
		program: &bytecode.CompiledProgram{
			TypeResolver: interfaceProxyTestResolver{info: info},
		},
	}

	adapter, ok := v.makeInterpretedInterfaceAdapter(targetType, concreteType, value.MakeNil())
	if !ok {
		t.Fatal("makeInterpretedInterfaceAdapter returned false")
	}
	if !called {
		t.Fatal("registered proxy factory was not called")
	}
	if _, ok := any(adapter).(sort.Interface); !ok {
		t.Fatalf("adapter type %T does not implement sort.Interface", adapter)
	}
}
