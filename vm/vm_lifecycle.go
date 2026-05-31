package vm

import "github.com/t04dJ14n9/gig/model/value"

// Reset prepares the VM for reuse by clearing execution state.
func (v *vm) Reset() {
	v.sp = 0
	v.fp = 0
	v.panicking = false
	v.panicVal = value.MakeNil()
	v.panicStack = v.panicStack[:0]
	v.deferDepth = 0
	v.ctx = nil
	// Clear all frames (prevents stale frame references from previous execution).
	for i := range v.frames {
		v.frames[i] = nil
	}
	// If shared globals are set (stateful mode or goroutine),
	// do not restore the local globals copy; the caller manages the shared state.
	if v.shared != nil {
		v.shared = nil
		return
	}
	// Stateless mode: restore globals to post-init snapshot, or zero them.
	if len(v.initialGlobals) == len(v.globals) {
		copy(v.globals, v.initialGlobals)
	} else {
		for i := range v.globals {
			v.globals[i] = value.Value{}
		}
	}
	// Restore external variable values (they should always be the same).
	for idx, ptr := range v.program.ExternalVarValues {
		if idx < len(v.globals) {
			v.globals[idx] = value.FromInterface(ptr)
		}
	}
	// Re-apply zero-valued struct globals (may have been overwritten by SSA init nil stores).
	for idx, zeroRV := range v.program.GlobalZeroValues {
		if idx < len(v.globals) {
			g := v.globals[idx]
			if !g.IsValid() || g.IsNil() {
				v.globals[idx] = value.MakeFromReflect(zeroRV)
			}
		}
	}
}

// growFrames doubles the frame stack capacity up to maxFrameDepth.
// Called when fp reaches the current slice length.
// Returns false if the stack is already at maximum capacity (stack overflow).
func (v *vm) growFrames() bool {
	cur := len(v.frames)
	if cur >= maxFrameDepth {
		return false
	}
	newCap := cur * 2
	if newCap > maxFrameDepth {
		newCap = maxFrameDepth
	}
	grown := make([]*Frame, newCap)
	copy(grown, v.frames)
	v.frames = grown
	return true
}

// BindSharedGlobals makes this VM execute against the provided SharedGlobals.
// All global loads/stores will go through the shared (locked) globals.
func (v *vm) BindSharedGlobals(sg *SharedGlobals) {
	v.shared = sg
}

// UnbindSharedGlobals detaches the VM from shared globals so that Reset (called
// when the VM is returned to the pool) does not clobber the shared state.
func (v *vm) UnbindSharedGlobals() {
	v.shared = nil
}

// Globals returns the VM's global variable slice.
func (v *vm) Globals() []value.Value {
	return v.globals
}

// getGlobals returns the globals slice for non-locked access.
// For shared mode, returns the raw slice from SharedGlobals.
// Individual OpGlobal/OpSetGlobal use locked methods directly.
func (v *vm) getGlobals() []value.Value {
	if v.shared != nil {
		return v.shared.Globals()
	}
	return v.globals
}
