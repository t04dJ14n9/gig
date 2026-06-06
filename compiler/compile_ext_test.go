package compiler

import (
	"testing"

	"github.com/t04dJ14n9/gig/model/external"
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
