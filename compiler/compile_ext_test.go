package compiler

import (
	"go/types"
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

func TestAttachExternalFuncReflectMetadataRecordsVariadicShape(t *testing.T) {
	info := &external.ExternalFuncInfo{}
	fn := func(prefix string, parts ...string) string { return prefix }

	attachExternalFuncReflectMetadata(info, fn)

	if !info.IsVariadic {
		t.Fatal("IsVariadic = false, want true")
	}
	if info.NumIn != 2 {
		t.Fatalf("NumIn = %d, want 2", info.NumIn)
	}
}

func TestAttachExternalFuncReflectMetadataIgnoresNonFunctions(t *testing.T) {
	info := &external.ExternalFuncInfo{IsVariadic: true, NumIn: 3}

	attachExternalFuncReflectMetadata(info, "not a function")

	if !info.IsVariadic || info.NumIn != 3 {
		t.Fatalf("metadata changed for non-function: IsVariadic=%v NumIn=%d", info.IsVariadic, info.NumIn)
	}
}

func TestExternalFuncValueFromLookupReturnsRegisteredFunction(t *testing.T) {
	fn := func() int { return 42 }
	lookup := &externalFuncValueLookup{
		pkgPath:  "example.com/host",
		funcName: "Answer",
		fn:       fn,
	}

	got := externalFuncValueFromLookup(lookup, "example.com/host", "Answer")

	gotFn, ok := got.(func() int)
	if !ok {
		t.Fatalf("lookup returned %T, want func() int", got)
	}
	if gotFn() != 42 {
		t.Fatalf("lookup function returned %d, want 42", gotFn())
	}
}

func TestExternalFuncValueFromLookupReturnsNilWithoutMatch(t *testing.T) {
	lookup := &externalFuncValueLookup{
		pkgPath:  "example.com/host",
		funcName: "Answer",
		fn:       func() int { return 42 },
	}

	if got := externalFuncValueFromLookup(nil, "example.com/host", "Answer"); got != nil {
		t.Fatalf("nil lookup returned %T, want nil", got)
	}
	if got := externalFuncValueFromLookup(lookup, "", "Answer"); got != nil {
		t.Fatalf("empty package returned %T, want nil", got)
	}
	if got := externalFuncValueFromLookup(lookup, "example.com/host", "Missing"); got != nil {
		t.Fatalf("missing function returned %T, want nil", got)
	}
}

func TestAttachExternalMethodDirectCallUsesQualifiedReceiverName(t *testing.T) {
	directCall := func([]value.Value) value.Value { return value.MakeInt(7) }
	lookup := &methodDirectCallLookup{
		typeName:   "example.com/host.Widget",
		methodName: "Close",
		directCall: directCall,
	}
	info := &external.ExternalMethodInfo{MethodName: "Close"}
	recvType := namedReceiverType("example.com/host", "host", "Widget")

	attachExternalMethodDirectCall(info, lookup, types.NewPointer(recvType))

	if info.DirectCall == nil {
		t.Fatal("DirectCall was not attached")
	}
	if got := info.DirectCall(nil).RawInt(); got != 7 {
		t.Fatalf("DirectCall result = %d, want 7", got)
	}
}

func TestAttachExternalMethodDirectCallLeavesInfoUnchangedWithoutLookup(t *testing.T) {
	info := &external.ExternalMethodInfo{MethodName: "Close"}
	recvType := namedReceiverType("example.com/host", "host", "Widget")

	attachExternalMethodDirectCall(info, nil, recvType)

	if info.DirectCall != nil {
		t.Fatal("DirectCall changed without lookup")
	}
}

func TestShouldSkipUnresolvedExternalFunctionOnlySkipsImportedInitStubs(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		pkgPath  string
		want     bool
	}{
		{name: "imported init", funcName: "init", pkgPath: "example.com/dep", want: true},
		{name: "main init", funcName: "init", pkgPath: "main", want: false},
		{name: "command line init", funcName: "init", pkgPath: "command-line-arguments", want: false},
		{name: "empty package init", funcName: "init", pkgPath: "", want: false},
		{name: "normal function", funcName: "F", pkgPath: "example.com/dep", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSkipUnresolvedExternalFunction(tt.funcName, tt.pkgPath)
			if got != tt.want {
				t.Fatalf("shouldSkipUnresolvedExternalFunction(%q, %q) = %v, want %v", tt.funcName, tt.pkgPath, got, tt.want)
			}
		})
	}
}

type externalFuncValueLookup struct {
	pkgPath  string
	funcName string
	fn       any
}

func (m *externalFuncValueLookup) LookupExternalFunc(pkgPath, funcName string) (any, func([]value.Value) value.Value, bool) {
	if pkgPath == m.pkgPath && funcName == m.funcName {
		return m.fn, nil, true
	}
	return nil, nil, false
}

func (m *externalFuncValueLookup) LookupMethodDirectCall(string, string) (func([]value.Value) value.Value, bool) {
	return nil, false
}

func (m *externalFuncValueLookup) LookupExternalVar(string, string) (any, bool) {
	return nil, false
}

func (m *externalFuncValueLookup) LookupExternalType(types.Type) (reflect.Type, bool) {
	return nil, false
}

func (m *externalFuncValueLookup) LookupExternalTypeByName(string, string) (reflect.Type, bool) {
	return nil, false
}

type methodDirectCallLookup struct {
	typeName   string
	methodName string
	directCall func([]value.Value) value.Value
}

func (m *methodDirectCallLookup) LookupExternalFunc(string, string) (any, func([]value.Value) value.Value, bool) {
	return nil, nil, false
}

func (m *methodDirectCallLookup) LookupMethodDirectCall(typeName, methodName string) (func([]value.Value) value.Value, bool) {
	if typeName == m.typeName && methodName == m.methodName {
		return m.directCall, true
	}
	return nil, false
}

func (m *methodDirectCallLookup) LookupExternalVar(string, string) (any, bool) {
	return nil, false
}

func (m *methodDirectCallLookup) LookupExternalType(types.Type) (reflect.Type, bool) {
	return nil, false
}

func (m *methodDirectCallLookup) LookupExternalTypeByName(string, string) (reflect.Type, bool) {
	return nil, false
}

func namedReceiverType(pkgPath, pkgName, typeName string) *types.Named {
	pkg := types.NewPackage(pkgPath, pkgName)
	obj := types.NewTypeName(0, pkg, typeName, nil)
	return types.NewNamed(obj, types.NewStruct(nil, nil), nil)
}

var _ PackageLookup = (*externalFuncValueLookup)(nil)
var _ PackageLookup = (*methodDirectCallLookup)(nil)
