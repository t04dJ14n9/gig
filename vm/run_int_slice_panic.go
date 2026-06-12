package vm

import (
	"fmt"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) setIntSliceIndexPanic(idx int64, length int) {
	// Match Go's slice-index panic text so native-vs-interpreted parity checks
	// continue to compare the same observable failure.
	v.panicking = true
	v.panicVal = value.FromInterface(fmt.Sprintf("runtime error: index out of range [%d] with length %d", idx, length))
}
