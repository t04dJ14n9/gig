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
