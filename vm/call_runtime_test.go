package vm

import (
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/model/value"
)

func TestUnpackVariadicArgsExpandsPackedSlices(t *testing.T) {
	tests := []struct {
		name string
		last value.Value
		want []any
	}{
		{
			name: "native value slice",
			last: value.MakeValueSlice([]value.Value{value.MakeInt(1), value.MakeString("x")}),
			want: []any{"prefix", int(1), "x"},
		},
		{
			name: "native int slice",
			last: value.MakeIntSlice([]int64{2, 3}),
			want: []any{"prefix", int(2), int(3)},
		},
		{
			name: "native bytes",
			last: value.MakeBytes([]byte{4, 5}),
			want: []any{"prefix", uint(4), uint(5)},
		},
		{
			name: "reflect slice",
			last: value.FromInterface([]string{"a", "b"}),
			want: []any{"prefix", "a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unpackVariadicArgs([]value.Value{value.MakeString("prefix"), tt.last}, 2)

			assertValueInterfaces(t, got, tt.want)
		})
	}
}

func TestUnpackVariadicArgsReturnsOriginalArgsForNonSlice(t *testing.T) {
	args := []value.Value{value.MakeString("prefix"), value.MakeString("tail")}

	got := unpackVariadicArgs(args, 2)

	if len(got) != 2 {
		t.Fatalf("unpacked len = %d, want 2", len(got))
	}
	if got[0].String() != "prefix" || got[1].String() != "tail" {
		t.Fatalf("unpacked values = %#v", got)
	}
}

func assertValueInterfaces(t *testing.T, got []value.Value, want []any) {
	t.Helper()

	gotInterfaces := make([]any, len(got))
	for i, val := range got {
		gotInterfaces[i] = val.Interface()
	}
	if !reflect.DeepEqual(gotInterfaces, want) {
		t.Fatalf("unpacked values = %#v, want %#v", gotInterfaces, want)
	}
}
