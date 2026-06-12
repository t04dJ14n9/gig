package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeVariableMemory(op bytecode.OpCode, frame *Frame) {
	switch op {
	case bytecode.OpGlobal:
		v.pushGlobal(frame.readUint16())
	case bytecode.OpSetGlobal:
		v.setGlobal(frame.readUint16(), v.pop())
	case bytecode.OpFree:
		v.pushFree(frame, frame.readByte())
	case bytecode.OpSetFree:
		v.setFree(frame, frame.readByte(), v.pop())
	}
}

func (v *vm) pushGlobal(idx uint16) {
	if sg := v.shared; sg != nil {
		v.pushSharedGlobal(sg, int(idx))
		return
	}
	v.pushLocalGlobal(idx)
}

func (v *vm) pushSharedGlobal(sg *SharedGlobals, idx int) {
	// Shared mode exposes a GlobalRef instead of a raw pointer so goroutines use
	// locked access and cannot race on the global value slot.
	if idx < sg.Len() {
		v.push(value.FromInterface(&GlobalRef{sg: sg, idx: idx}))
	}
}

func (v *vm) pushLocalGlobal(idx uint16) {
	globals := v.globals
	if int(idx) < len(globals) {
		v.push(value.FromInterface(&globals[idx]))
	}
}

func (v *vm) setGlobal(idx uint16, val value.Value) {
	if sg := v.shared; sg != nil {
		v.setSharedGlobal(sg, int(idx), val)
		return
	}
	v.setLocalGlobal(idx, val)
}

func (v *vm) setSharedGlobal(sg *SharedGlobals, idx int, val value.Value) {
	if idx < sg.Len() {
		sg.Set(idx, val)
	}
}

func (v *vm) setLocalGlobal(idx uint16, val value.Value) {
	globals := v.globals
	if int(idx) < len(globals) {
		globals[idx] = val
	}
}

func (v *vm) pushFree(frame *Frame, idx byte) {
	slot := freeSlot(frame, idx)
	if slot == nil {
		v.push(value.MakeNil())
		return
	}
	// Free variables are mutable slots shared by closures; push the value in
	// the slot, not the slot pointer itself.
	v.push(*slot)
}

func (v *vm) setFree(frame *Frame, idx byte, val value.Value) {
	slot := freeSlot(frame, idx)
	if slot != nil {
		*slot = val
	}
}

func freeSlot(frame *Frame, idx byte) *value.Value {
	if int(idx) >= len(frame.freeVars) {
		return nil
	}
	return frame.freeVars[idx]
}
