package compiler

import (
	"go/types"
	"testing"
)

func TestAllOnesConstantPreservesMaskWidths(t *testing.T) {
	tests := []struct {
		name string
		typ  types.Type
		want any
	}{
		{"uint8", types.Typ[types.Uint8], uint8(0xFF)},
		{"uint16", types.Typ[types.Uint16], uint16(0xFFFF)},
		{"uint32", types.Typ[types.Uint32], uint32(0xFFFFFFFF)},
		{"uint", types.Typ[types.Uint], uint(^uint(0))},
		{"uint64", types.Typ[types.Uint64], uint64(^uint64(0))},
		{"uintptr", types.Typ[types.Uintptr], uintptr(^uintptr(0))},
		{"int8", types.Typ[types.Int8], int8(-1)},
		{"int16", types.Typ[types.Int16], int16(-1)},
		{"int32", types.Typ[types.Int32], int32(-1)},
		{"int", types.Typ[types.Int], int(^int(0))},
		{"int64", types.Typ[types.Int64], int64(-1)},
		{"unsupported basic", types.Typ[types.String], int64(-1)},
		{"non-basic", types.NewSlice(types.Typ[types.Int]), int64(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allOnesConstant(tt.typ)
			if got != tt.want {
				t.Fatalf("allOnesConstant(%s) = %#v (%T), want %#v (%T)", tt.typ.String(), got, got, tt.want, tt.want)
			}
		})
	}
}
