package vm

import (
	"math"

	"github.com/t04dJ14n9/gig/model/value"
)

// lessEqCmp returns a <= b, correctly handling IEEE 754 NaN (NaN <= x is always false).
func lessEqCmp(a, b value.Value) bool {
	cmp := a.Cmp(b)
	if cmp < 0 {
		return true
	}
	if cmp > 0 {
		return false
	}
	// cmp == 0: could be a == b, or one/both are NaN
	if a.Kind() == value.KindFloat && math.IsNaN(a.Float()) {
		return false
	}
	if b.Kind() == value.KindFloat && math.IsNaN(b.Float()) {
		return false
	}
	return true
}

// greaterEqCmp returns a >= b, correctly handling IEEE 754 NaN (NaN >= x is always false).
func greaterEqCmp(a, b value.Value) bool {
	cmp := a.Cmp(b)
	if cmp > 0 {
		return true
	}
	if cmp < 0 {
		return false
	}
	// cmp == 0: could be a == b, or one/both are NaN
	if a.Kind() == value.KindFloat && math.IsNaN(a.Float()) {
		return false
	}
	if b.Kind() == value.KindFloat && math.IsNaN(b.Float()) {
		return false
	}
	return true
}
