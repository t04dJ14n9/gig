package vm

import (
	"testing"

	"github.com/t04dJ14n9/gig/model/value"
)

func TestAppendValueAppendsIntToByteSlice(t *testing.T) {
	got := appendValue(value.MakeBytes([]byte("a")), value.MakeInt('b'))

	bytes, ok := got.Bytes()
	if !ok {
		t.Fatalf("appendValue returned %s, want byte slice", got.Kind())
	}
	if string(bytes) != "ab" {
		t.Fatalf("appendValue bytes = %q, want %q", string(bytes), "ab")
	}
}
