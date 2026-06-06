package bytecode

import (
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

func TestResolveExternalFuncKeepsDirectMetadataWhenFuncIsNil(t *testing.T) {
	directCall := func([]value.Value) value.Value { return value.MakeInt(1) }
	rc := ResolveConstant(&external.ExternalFuncInfo{
		PkgPath:    "example.com/host",
		FuncName:   "F",
		DirectCall: directCall,
		IsVariadic: true,
		NumIn:      2,
	})

	if rc == nil {
		t.Fatal("ResolveConstant returned nil")
	}
	if rc.PkgPath != "example.com/host" || rc.FuncName != "F" {
		t.Fatalf("resolved identity = %s.%s, want example.com/host.F", rc.PkgPath, rc.FuncName)
	}
	if rc.DirectCall == nil {
		t.Fatal("DirectCall was not preserved")
	}
	if !rc.IsVariadic || rc.NumIn != 2 {
		t.Fatalf("variadic metadata = (%v, %d), want (true, 2)", rc.IsVariadic, rc.NumIn)
	}
	if rc.Fn.IsValid() || rc.FnType != nil {
		t.Fatalf("nil Func should not populate reflect metadata: Fn=%v FnType=%v", rc.Fn, rc.FnType)
	}
}

func TestResolveExternalFuncAttachesReflectMetadata(t *testing.T) {
	fn := func(prefix string, parts ...string) string { return prefix }
	rc := ResolveConstant(&external.ExternalFuncInfo{
		PkgPath:  "example.com/host",
		FuncName: "Join",
		Func:     fn,
	})

	if rc == nil {
		t.Fatal("ResolveConstant returned nil")
	}
	if !rc.Fn.IsValid() || rc.FnType == nil {
		t.Fatalf("reflect metadata was not populated: Fn=%v FnType=%v", rc.Fn, rc.FnType)
	}
	if rc.FnType != reflect.TypeOf(fn) {
		t.Fatalf("FnType = %v, want %v", rc.FnType, reflect.TypeOf(fn))
	}
	if !rc.IsVariadic || rc.NumIn != 2 {
		t.Fatalf("variadic metadata = (%v, %d), want (true, 2)", rc.IsVariadic, rc.NumIn)
	}
}

func TestResolveLegacyFuncDefaultsToTrustedStdlibCall(t *testing.T) {
	fn := func() int { return 1 }
	rc := ResolveConstant(fn)

	if rc == nil {
		t.Fatal("ResolveConstant returned nil")
	}
	if !rc.IsStdlib {
		t.Fatal("legacy function constants should remain trusted stdlib calls")
	}
	if rc.FnType != reflect.TypeOf(fn) {
		t.Fatalf("FnType = %v, want %v", rc.FnType, reflect.TypeOf(fn))
	}
}
