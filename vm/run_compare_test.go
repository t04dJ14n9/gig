package vm

import (
	"math"
	"testing"

	"github.com/t04dJ14n9/gig/model/value"
)

func TestRunCompareNaNOrderedComparisons(t *testing.T) {
	nan := value.MakeFloat(math.NaN())
	one := value.MakeFloat(1)

	if lessEqCmp(nan, one) {
		t.Fatal("NaN <= 1 returned true")
	}
	if lessEqCmp(one, nan) {
		t.Fatal("1 <= NaN returned true")
	}
	if greaterEqCmp(nan, one) {
		t.Fatal("NaN >= 1 returned true")
	}
	if greaterEqCmp(one, nan) {
		t.Fatal("1 >= NaN returned true")
	}
}
