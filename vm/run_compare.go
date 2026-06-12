package vm

import (
	"math"

	"github.com/t04dJ14n9/gig/model/value"
)

// lessEqCmp returns a <= b, preserving IEEE 754 NaN semantics.
func lessEqCmp(a, b value.Value) bool {
	cmp := a.Cmp(b)
	if cmp < 0 {
		return true
	}
	if cmp > 0 {
		return false
	}
	return !isNaNFloat(a) && !isNaNFloat(b)
}

// greaterEqCmp returns a >= b, preserving IEEE 754 NaN semantics.
func greaterEqCmp(a, b value.Value) bool {
	cmp := a.Cmp(b)
	if cmp > 0 {
		return true
	}
	if cmp < 0 {
		return false
	}
	return !isNaNFloat(a) && !isNaNFloat(b)
}

func isNaNFloat(v value.Value) bool {
	return v.Kind() == value.KindFloat && math.IsNaN(v.Float())
}
