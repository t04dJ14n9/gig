package vm

import "github.com/t04dJ14n9/gig/model/value"

func (v *vm) runSlowConst(idx uint16) value.Value {
	// Most constants are prebaked and loaded directly in run.go. This fallback
	// keeps legacy interface-backed constants correct without making the common
	// constant load carry the extra bounds branch.
	if int(idx) < len(v.program.Constants) {
		return value.FromInterface(v.program.Constants[idx])
	}
	return value.Value{}
}
