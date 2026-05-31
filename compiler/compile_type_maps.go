package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"
)

type localTypeMaps struct {
	IsInt      []bool
	IsIntSlice []bool
}

func buildLocalTypeMaps(values map[ssa.Value]int) localTypeMaps {
	// The optimizer consumes dense local-indexed slices, not SSA maps. Build
	// both maps from the same local table snapshot so specialized int and slice
	// instructions agree with the final local-slot numbering.
	return localTypeMaps{
		IsInt:      buildTypeMap(values, isIntType),
		IsIntSlice: buildTypeMap(values, isIntSliceType),
	}
}

// isIntType returns true if the type is int or int64 (full-width signed integers).
// Only full-width types are eligible for OpInt* superinstructions because
// those instructions use int64-backed intLocals without truncation wrapping.
// Sized types (int8, int16, int32) must use the generic path which preserves
// the size tag via MakeIntSized for correct overflow wrapping semantics.
func isIntType(t types.Type) bool {
	if t == nil {
		return false
	}
	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Int, types.Int64:
		return true
	}
	return false
}

// isIntSliceType returns true if the type is []int or []int64 (matches native int slice fast path).
func isIntSliceType(t types.Type) bool {
	if t == nil {
		return false
	}
	sl, ok := t.Underlying().(*types.Slice)
	if !ok {
		return false
	}
	elem := sl.Elem()
	if elem == nil {
		return false
	}
	basic, ok := elem.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Int, types.Int64:
		return true
	}
	return false
}

func buildTypeMap(values map[ssa.Value]int, predicate func(types.Type) bool) []bool {
	result := make([]bool, len(values))
	for v, idx := range values {
		if predicate(v.Type()) {
			result[idx] = true
		}
	}
	return result
}

func buildConstIntMap(constants []any) []bool {
	// compileConst stores Go integer widths as concrete Go values. The
	// optimizer only needs to know whether a constant can feed int-specialized
	// opcodes; width-preserving wrapping happens in value construction.
	constIsInt := make([]bool, len(constants))
	for i, k := range constants {
		switch k.(type) {
		case int, int8, int16, int32, int64:
			constIsInt[i] = true
		}
	}
	return constIsInt
}
