package compiler

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/t04dJ14n9/gig/importer"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
	"github.com/t04dJ14n9/gig/vm"

	"github.com/t04dJ14n9/gig/compiler/testdata/thirdpartyiface"
)

type boundaryStringer interface {
	String() string
}

func boundaryAcceptAny(any) int {
	return 1
}

func boundaryAcceptAnys(...any) int {
	return 1
}

func boundaryAcceptIntFunc(fn func(int) int) int {
	return fn(41)
}

func boundaryAcceptAnyResultFunc(fn func() any) int {
	if fn() == nil {
		return 0
	}
	return 1
}

func boundaryAcceptStringer(boundaryStringer) int {
	return 1
}

func boundaryAcceptStringers(items ...thirdpartyiface.Stringer) int {
	return len(items)
}

type boundaryMismatchSorter interface {
	Len() int
	Less(string, string) bool
	Swap(int, int)
}

func boundaryAcceptMismatchSorter(boundaryMismatchSorter) int {
	return 1
}

type boundarySorterLike interface {
	Len() int
	Less(int, int) bool
	Swap(int, int)
}

type boundaryHost struct{}

func boundaryNewHost() boundaryHost {
	return boundaryHost{}
}

func (boundaryHost) AcceptStringer(thirdpartyiface.Stringer) int {
	return 1
}

func (boundaryHost) AcceptAny(any) int {
	return 1
}

func (boundaryHost) AcceptAnyThenStringer(any, thirdpartyiface.Stringer) int {
	return 1
}

type BoundaryExpressionHost struct{}

func boundaryNewExpressionHost() BoundaryExpressionHost {
	return BoundaryExpressionHost{}
}

func (BoundaryExpressionHost) AcceptAny(any) int {
	return 1
}

func boundaryAcceptSorterLike(boundarySorterLike) int {
	return 1
}

func newBoundaryRegistry() importer.PackageRegistry {
	reg := importer.NewRegistry()
	pkg := reg.RegisterPackage("example.com/host", "host")
	pkg.AddFunction("AcceptAny", boundaryAcceptAny, "", nil)
	pkg.AddFunction("AcceptAnys", boundaryAcceptAnys, "", nil)
	pkg.AddFunction("AcceptIntFunc", boundaryAcceptIntFunc, "", nil)
	pkg.AddFunction("AcceptAnyResultFunc", boundaryAcceptAnyResultFunc, "", nil)
	pkg.AddFunction("AcceptStringer", boundaryAcceptStringer, "", nil)
	pkg.AddFunction("AcceptMismatchSorter", boundaryAcceptMismatchSorter, "", nil)
	pkg.AddFunction("AcceptSorterLike", boundaryAcceptSorterLike, "", nil)
	return reg
}

func newBoundaryRegistryWithSortProxy() importer.PackageRegistry {
	reg := newBoundaryRegistry().(*importer.Registry)
	reg.AddInterfaceProxy("sort", "Interface", reflect.TypeOf((*sort.Interface)(nil)).Elem(), []string{"Len", "Less", "Swap"}, func(value.Value, string, external.InterfaceMethodCaller) (any, bool) {
		return sort.IntSlice{}, true
	})
	return reg
}

func newBoundaryRegistryWithMethodProxy() importer.PackageRegistry {
	reg := newBoundaryRegistryWithProxy().(*importer.Registry)
	pkg := reg.RegisterPackage("example.com/host", "host")
	pkg.AddFunction("NewHost", boundaryNewHost, "", nil)
	pkg.AddFunction("NewExpressionHost", boundaryNewExpressionHost, "", nil)
	pkg.AddFunction("AcceptStringers", boundaryAcceptStringers, "", nil)
	pkg.AddType("Host", reflect.TypeOf(boundaryHost{}), "")
	pkg.AddType("BoundaryExpressionHost", reflect.TypeOf(BoundaryExpressionHost{}), "")
	return reg
}

type boundaryStringerProxy struct {
	call external.InterfaceMethodCaller
}

func (p *boundaryStringerProxy) String() string {
	result, ok := p.call("String")
	if !ok {
		return ""
	}
	return result.String()
}

func newBoundaryRegistryWithProxy() importer.PackageRegistry {
	reg := importer.NewRegistry()
	pkgPath := "github.com/t04dJ14n9/gig/compiler/testdata/thirdpartyiface"
	pkg := reg.RegisterPackage(pkgPath, "thirdpartyiface")
	ifaceType := reflect.TypeOf((*thirdpartyiface.Stringer)(nil)).Elem()
	pkg.AddFunction("AcceptStringer", thirdpartyiface.AcceptStringer, "", nil)
	pkg.AddType("Stringer", ifaceType, "")
	pkg.AddInterfaceProxy("Stringer", ifaceType, []string{"String"}, func(_ value.Value, _ string, call external.InterfaceMethodCaller) (any, bool) {
		return &boundaryStringerProxy{call: call}, true
	})
	return reg
}

func TestThirdPartyBoundaryRejectsScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() int {
	return host.AcceptAny(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistry())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to third-party any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsIndirectScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() int {
	fns := []func(any) int{host.AcceptAny}
	return fns[0](ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistry())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed indirectly to third-party any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsReturnedExternalFuncScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func getHost() func(any) int {
	return host.AcceptAny
}

func F() int {
	return getHost()(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistry())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to returned third-party function value any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsPassthroughReturnedExternalFuncScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func passthrough(fn func(any) int) func(any) int {
	return fn
}

func F() int {
	return passthrough(host.AcceptAny)(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistry())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to passthrough returned third-party function value any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsClosureReturnedExternalFuncScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() int {
	getHost := func() func(any) int {
		return host.AcceptAny
	}
	return getHost()(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistry())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to closure-returned third-party function value any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsClosurePassthroughReturnedExternalFuncScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() int {
	passthrough := func(fn func(any) int) func(any) int {
		return fn
	}
	return passthrough(host.AcceptAny)(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistry())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to closure-passthrough returned third-party function value any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsDeferredScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() {
	defer host.AcceptAny(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistry())
	if err == nil {
		t.Fatal("expected Build to reject script struct deferred to third-party any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeHiddenDeferredScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func deferHost(v any) {
	defer host.AcceptAny(v)
}

func F() int {
	deferHost(ScriptStruct{Name: "x"})
	return 1
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject deferred script struct hidden behind any at third-party boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsGoScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() {
	go host.AcceptAny(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistry())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to third-party any in go statement")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeHiddenGoScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func callHost(v any) {
	go host.AcceptAny(v)
}

func F() {
	callHost(ScriptStruct{Name: "x"})
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject go script struct hidden behind any at third-party boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsGoMethodScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() {
	h := host.NewHost()
	go h.AcceptAny(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to third-party method any in go statement")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeHiddenGoMethodScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func callHost(v any) {
	h := host.NewHost()
	go h.AcceptAny(v)
}

func F() {
	callHost(ScriptStruct{Name: "x"})
}
`
	result, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject go script struct hidden behind any at third-party method boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeHiddenScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func callHost(v any) int {
	return host.AcceptAny(v)
}

func F() int {
	return callHost(ScriptStruct{Name: "x"})
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject script struct hidden behind any at third-party boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeHiddenScriptStructSliceToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func callHost(v any) int {
	return host.AcceptAny(v)
}

func F() int {
	return callHost([]ScriptStruct{{Name: "x"}})
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject script struct slice hidden behind any at third-party boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeHiddenScriptStructVariadicSpreadToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func callHost(v []any) int {
	return host.AcceptAnys(v...)
}

func F() int {
	return callHost([]any{ScriptStruct{Name: "x"}})
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject script struct hidden in variadic spread to third-party any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeHiddenScriptStructMapToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func callHost(v any) int {
	return host.AcceptAny(v)
}

func F() int {
	return callHost(map[string]any{"x": ScriptStruct{Name: "x"}})
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject script struct map hidden behind any at third-party boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeDeeplyHiddenScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func callHost(v any) int {
	return host.AcceptAny(v)
}

func F() int {
	var v any = ScriptStruct{Name: "x"}
	for i := 0; i < 70; i++ {
		v = []any{v}
	}
	return callHost(v)
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject deeply nested script struct hidden behind any at third-party boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeHiddenScriptFuncToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func takesScript(ScriptStruct) {
}

func callHost(v any) int {
	return host.AcceptAny(v)
}

func F() int {
	return callHost(takesScript)
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject script function hidden behind any at third-party boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryAllowsTypedScriptFunc(t *testing.T) {
	src := `
package main

import "example.com/host"

func inc(v int) int {
	return v + 1
}

func F() int {
	return host.AcceptIntFunc(inc)
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build rejected typed script func passed to third-party func parameter: %v", err)
	}
	got, err := vm.New(result.Program).Execute("F", context.Background())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got.Int() != 42 {
		t.Fatalf("F() = %d, want 42", got.Int())
	}
}

func TestThirdPartyBoundaryRejectsTypedScriptFuncWithInterfaceResult(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func makeScript() any {
	return ScriptStruct{Name: "x"}
}

func F() int {
	return host.AcceptAnyResultFunc(makeScript)
}
`
	result, err := Build(src, newBoundaryRegistry())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject script function with interface result at third-party boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryUnsafeOptionAllowsRuntimeHiddenScriptStruct(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func callHost(v any) int {
	return host.AcceptAny(v)
}

func F() int {
	return callHost(ScriptStruct{Name: "x"})
}
`
	result, err := Build(src, newBoundaryRegistry(), WithAllowUnsafeTypePass())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	got, err := vm.New(result.Program).Execute("F", context.Background())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got.Int() != 1 {
		t.Fatalf("F() = %d, want 1", got.Int())
	}
}

func TestThirdPartyBoundaryRejectsScriptStructToInterfaceWithoutProxy(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func (s ScriptStruct) String() string {
	return s.Name
}

func F() int {
	return host.AcceptStringer(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistry())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to third-party interface without proxy")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryAllowsPrimitiveToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

func F() int {
	return host.AcceptAny(42)
}
`
	if _, err := Build(src, newBoundaryRegistry()); err != nil {
		t.Fatalf("Build rejected primitive passed to third-party any: %v", err)
	}
}

func TestThirdPartyBoundaryUnsafeOptionAllowsScriptStruct(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() int {
	return host.AcceptAny(ScriptStruct{Name: "x"})
}
`
	if _, err := Build(src, newBoundaryRegistry(), WithAllowUnsafeTypePass()); err != nil {
		t.Fatalf("Build rejected unsafe script struct pass: %v", err)
	}
}

func TestThirdPartyBoundaryAllowsScriptStructToInterfaceWithProxy(t *testing.T) {
	src := `
package main

import "github.com/t04dJ14n9/gig/compiler/testdata/thirdpartyiface"

type ScriptStruct struct {
	Name string
}

func (s ScriptStruct) String() string {
	return s.Name
}

func F() int {
	return thirdpartyiface.AcceptStringer(ScriptStruct{Name: "x"})
}
`
	result, err := Build(src, newBoundaryRegistryWithProxy())
	if err != nil {
		t.Fatalf("Build rejected proxied third-party interface: %v", err)
	}
	got, err := vm.New(result.Program).Execute("F", context.Background())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got.Int() != 1 {
		t.Fatalf("F() = %d, want 1", got.Int())
	}
}

func TestThirdPartyBoundaryAllowsScriptStructToMethodInterfaceWithProxy(t *testing.T) {
	src := `
package main

import (
	"example.com/host"
	"github.com/t04dJ14n9/gig/compiler/testdata/thirdpartyiface"
)

var _ thirdpartyiface.Stringer

type ScriptStruct struct {
	Name string
}

func (s ScriptStruct) String() string {
	return s.Name
}

func F() int {
	h := host.NewHost()
	return h.AcceptStringer(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err != nil {
		t.Fatalf("Build rejected proxied third-party method interface: %v", err)
	}
}

func TestThirdPartyBoundaryAllowsMethodValueScriptStructToInterfaceWithProxy(t *testing.T) {
	src := `
package main

import (
	"example.com/host"
	"github.com/t04dJ14n9/gig/compiler/testdata/thirdpartyiface"
)

var _ thirdpartyiface.Stringer

type ScriptStruct struct {
	Name string
}

func (s ScriptStruct) String() string {
	return s.Name
}

func F() int {
	h := host.NewHost()
	f := h.AcceptStringer
	return f(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err != nil {
		t.Fatalf("Build rejected proxied third-party method value interface: %v", err)
	}
}

func TestThirdPartyBoundaryAllowsVariadicScriptStructToInterfaceWithProxy(t *testing.T) {
	src := `
package main

import (
	"example.com/host"
	"github.com/t04dJ14n9/gig/compiler/testdata/thirdpartyiface"
)

var _ thirdpartyiface.Stringer

type ScriptStruct struct {
	Name string
}

func (s ScriptStruct) String() string {
	return s.Name
}

func F() int {
	return host.AcceptStringers(ScriptStruct{Name: "x"}, ScriptStruct{Name: "y"})
}
`
	result, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err != nil {
		t.Fatalf("Build rejected variadic proxied third-party interface: %v", err)
	}
	got, err := vm.New(result.Program).Execute("F", context.Background())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got.Int() != 2 {
		t.Fatalf("F() = %d, want 2", got.Int())
	}
}

func TestThirdPartyBoundaryAllowsVariadicInterfaceSliceSpreadWithProxy(t *testing.T) {
	src := `
package main

import (
	"example.com/host"
	"github.com/t04dJ14n9/gig/compiler/testdata/thirdpartyiface"
)

var _ thirdpartyiface.Stringer

type ScriptStruct struct {
	Name string
}

func (s ScriptStruct) String() string {
	return s.Name
}

func callHost(items []thirdpartyiface.Stringer) int {
	return host.AcceptStringers(items...)
}

func F() int {
	return callHost([]thirdpartyiface.Stringer{ScriptStruct{Name: "x"}, ScriptStruct{Name: "y"}})
}
`
	result, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err != nil {
		t.Fatalf("Build rejected variadic proxied third-party interface spread: %v", err)
	}
	got, err := vm.New(result.Program).Execute("F", context.Background())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got.Int() != 2 {
		t.Fatalf("F() = %d, want 2", got.Int())
	}
}

func TestThirdPartyBoundaryAllowsRuntimeHiddenScriptStructToInterfaceWithProxy(t *testing.T) {
	src := `
package main

import (
	"github.com/t04dJ14n9/gig/compiler/testdata/thirdpartyiface"
)

type ScriptStruct struct {
	Name string
}

func (s ScriptStruct) String() string {
	return s.Name
}

func callHost(v thirdpartyiface.Stringer) int {
	return thirdpartyiface.AcceptStringer(v)
}

func F() int {
	return callHost(ScriptStruct{Name: "x"})
}
`
	result, err := Build(src, newBoundaryRegistryWithProxy())
	if err != nil {
		t.Fatalf("Build rejected runtime-hidden proxied third-party interface: %v", err)
	}
	got, err := vm.New(result.Program).Execute("F", context.Background())
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got.Int() != 1 {
		t.Fatalf("F() = %d, want 1", got.Int())
	}
}

func TestThirdPartyBoundaryRejectsMethodAnyBeforeProxiedInterface(t *testing.T) {
	src := `
package main

import (
	"example.com/host"
	"github.com/t04dJ14n9/gig/compiler/testdata/thirdpartyiface"
)

var _ thirdpartyiface.Stringer

type ScriptStruct struct {
	Name string
}

type ScriptStringer struct {
	Name string
}

func (s ScriptStringer) String() string {
	return s.Name
}

func F() int {
	h := host.NewHost()
	return h.AcceptAnyThenStringer(ScriptStruct{Name: "x"}, ScriptStringer{Name: "ok"})
}
`
	_, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to method any before proxied interface")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsMethodValueScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() int {
	h := host.NewHost()
	f := h.AcceptAny
	return f(ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to third-party method value any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsIndirectMethodValueScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() int {
	h := host.NewHost()
	fns := []func(any) int{h.AcceptAny}
	return fns[0](ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed indirectly to third-party method value any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsMethodExpressionScriptStructToAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func F() int {
	h := host.NewExpressionHost()
	f := host.BoundaryExpressionHost.AcceptAny
	return f(h, ScriptStruct{Name: "x"})
}
`
	_, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err == nil {
		t.Fatal("expected Build to reject script struct passed to third-party method expression any")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsRuntimeHiddenScriptStructToMethodAny(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptStruct struct {
	Name string
}

func callHost(v any) int {
	h := host.NewHost()
	return h.AcceptAny(v)
}

func F() int {
	return callHost(ScriptStruct{Name: "x"})
}
`
	result, err := Build(src, newBoundaryRegistryWithMethodProxy())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	_, err = vm.New(result.Program).Execute("F", context.Background())
	if err == nil {
		t.Fatal("expected Execute to reject script struct hidden behind any at third-party method boundary")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsInterfaceProxyWithMismatchedSignature(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptSorter []int

func (s ScriptSorter) Len() int {
	return len(s)
}

func (s ScriptSorter) Less(i, j string) bool {
	return i < j
}

func (s ScriptSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func F() int {
	return host.AcceptMismatchSorter(ScriptSorter{1, 2})
}
`
	_, err := Build(src, newBoundaryRegistryWithSortProxy())
	if err == nil {
		t.Fatal("expected Build to reject mismatched interface signature despite registered sort proxy")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}

func TestThirdPartyBoundaryRejectsStructuralProxyMatchWithoutRegisteredTarget(t *testing.T) {
	src := `
package main

import "example.com/host"

type ScriptSorter []int

func (s ScriptSorter) Len() int {
	return len(s)
}

func (s ScriptSorter) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s ScriptSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func F() int {
	return host.AcceptSorterLike(ScriptSorter{1, 2})
}
`
	_, err := Build(src, newBoundaryRegistryWithSortProxy())
	if err == nil {
		t.Fatal("expected Build to reject structural interface match without registered target proxy")
	}
	if !strings.Contains(err.Error(), "cannot pass interpreter-defined type") {
		t.Fatalf("error = %v, want interpreter-defined type diagnostic", err)
	}
}
