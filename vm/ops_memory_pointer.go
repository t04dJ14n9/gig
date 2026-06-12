package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executePointerMemory(op bytecode.OpCode, frame *Frame) {
	switch op {
	case bytecode.OpAddr:
		v.pushLocalAddress(frame, frame.readUint16())
	case bytecode.OpFieldAddr:
		fieldIdx := frame.readUint16()
		v.push(fieldAddressValue(v.pop(), int(fieldIdx)))
	case bytecode.OpIndexAddr:
		v.pushIndexAddress()
	case bytecode.OpDeref:
		v.push(dereferenceValue(v.pop()))
	case bytecode.OpSetDeref:
		val := v.pop()
		ptr := v.pop()
		v.setDereferenceValue(ptr, val)
	}
}

func (v *vm) pushLocalAddress(frame *Frame, localIdx uint16) {
	// Once a local's address escapes, frame pooling is unsafe because other
	// values may hold a pointer into this frame after the call returns.
	if int(localIdx) >= len(frame.locals) {
		v.push(value.MakeNil())
		return
	}
	frame.addrTaken = true
	v.push(value.FromInterface(&frame.locals[localIdx]))
}

func (v *vm) pushIndexAddress() {
	index := v.pop()
	container := v.pop()
	v.push(indexAddressValue(container, int(index.Int())))
}
