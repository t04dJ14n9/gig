package vm

import (
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/model/value"
)

type methodClosureReceiver struct{}

func (methodClosureReceiver) Apply(fn func(int) int) int { return fn(3) }

type methodClosureExecutor struct{}

var _ value.ClosureExecutor = methodClosureExecutor{}

func (methodClosureExecutor) Execute(args []reflect.Value, outTypes []reflect.Type) []reflect.Value {
	if len(args) == 0 || len(outTypes) == 0 {
		return nil
	}
	return []reflect.Value{args[0].Convert(outTypes[0])}
}

func TestConvertClosureArgsForMethodConvertsFuncParams(t *testing.T) {
	args := []value.Value{
		value.FromInterface(methodClosureReceiver{}),
		value.MakeFunc(methodClosureExecutor{}),
	}

	convertClosureArgsForMethod("Apply", args)

	rv, ok := args[1].ReflectValue()
	if !ok {
		t.Fatal("converted method closure argument is not reflect-backed")
	}
	if rv.Kind() != reflect.Func {
		t.Fatalf("converted method closure kind = %v, want Func", rv.Kind())
	}
	out := rv.Call([]reflect.Value{reflect.ValueOf(7)})
	if len(out) != 1 || out[0].Int() != 7 {
		t.Fatalf("converted closure returned %#v, want 7", out)
	}
}

func TestConvertClosureArgsForMethodLeavesNonFuncParamsAlone(t *testing.T) {
	args := []value.Value{
		value.FromInterface(methodClosureReceiver{}),
		value.MakeInt(1),
		value.MakeFunc(methodClosureExecutor{}),
	}

	convertClosureArgsForMethod("Apply", args)

	if args[2].Kind() != value.KindFunc {
		t.Fatalf("out-of-range closure argument kind = %v, want KindFunc", args[2].Kind())
	}
}

func TestConvertClosureArgsForMethodIgnoresMissingMethod(t *testing.T) {
	args := []value.Value{
		value.FromInterface(methodClosureReceiver{}),
		value.MakeFunc(methodClosureExecutor{}),
	}

	convertClosureArgsForMethod("Missing", args)

	if args[1].Kind() != value.KindFunc {
		t.Fatalf("missing method closure argument kind = %v, want KindFunc", args[1].Kind())
	}
}

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
