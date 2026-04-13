// basickinds.go provides canonical mappings between Go's various type systems:
// go/types.BasicKind ↔ reflect.Kind ↔ reflect.Type
//
// This is the single source of truth for basic type conversions used by both
// the VM and the importer packages. By centralizing these mappings, we ensure
// consistency and prevent divergence between packages.
package bytecode

import (
	"go/types"
	"reflect"
)

// BasicKindToReflectType maps go/types.BasicKind to reflect.Type.
// Used for efficient O(1) lookup during go/types → reflect.Type conversion.
var BasicKindToReflectType = map[types.BasicKind]reflect.Type{
	types.Bool:       reflect.TypeFor[bool](),
	types.Int:        reflect.TypeFor[int](),
	types.Int8:       reflect.TypeFor[int8](),
	types.Int16:      reflect.TypeFor[int16](),
	types.Int32:      reflect.TypeFor[int32](),
	types.Int64:      reflect.TypeFor[int64](),
	types.Uint:       reflect.TypeFor[uint](),
	types.Uint8:      reflect.TypeFor[uint8](),
	types.Uint16:     reflect.TypeFor[uint16](),
	types.Uint32:     reflect.TypeFor[uint32](),
	types.Uint64:     reflect.TypeFor[uint64](),
	types.Uintptr:    reflect.TypeFor[uintptr](),
	types.Float32:    reflect.TypeFor[float32](),
	types.Float64:    reflect.TypeFor[float64](),
	types.Complex64:  reflect.TypeFor[complex64](),
	types.Complex128: reflect.TypeFor[complex128](),
	types.String:     reflect.TypeFor[string](),
}

// ReflectKindToBasicKind maps reflect.Kind to go/types.BasicKind.
// Used for efficient O(1) lookup during reflect.Type → go/types conversion.
var ReflectKindToBasicKind = map[reflect.Kind]types.BasicKind{
	reflect.Bool:          types.Bool,
	reflect.Int:           types.Int,
	reflect.Int8:          types.Int8,
	reflect.Int16:         types.Int16,
	reflect.Int32:         types.Int32,
	reflect.Int64:         types.Int64,
	reflect.Uint:          types.Uint,
	reflect.Uint8:         types.Uint8,
	reflect.Uint16:        types.Uint16,
	reflect.Uint32:        types.Uint32,
	reflect.Uint64:        types.Uint64,
	reflect.Uintptr:       types.Uintptr,
	reflect.Float32:       types.Float32,
	reflect.Float64:       types.Float64,
	reflect.Complex64:     types.Complex64,
	reflect.Complex128:    types.Complex128,
	reflect.String:        types.String,
	reflect.UnsafePointer: types.UnsafePointer,
}

// BasicKindFromReflectKind returns the types.BasicKind for a reflect.Kind,
// or types.Invalid if the kind is not a basic type.
func BasicKindFromReflectKind(k reflect.Kind) types.BasicKind {
	if bk, ok := ReflectKindToBasicKind[k]; ok {
		return bk
	}
	return types.Invalid
}

// BasicTypeFromReflectKind returns the types.Typ (canonical type) for a reflect.Kind,
// or nil if the kind is not a basic type.
func BasicTypeFromReflectKind(k reflect.Kind) *types.Basic {
	if bk, ok := ReflectKindToBasicKind[k]; ok {
		return types.Typ[bk]
	}
	return nil
}
