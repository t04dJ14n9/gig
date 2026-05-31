package vm

import "github.com/t04dJ14n9/gig/model/bytecode"

func (v *vm) executeFieldMemory(op bytecode.OpCode, frame *Frame) {
	switch op {
	case bytecode.OpField:
		fieldIdx := frame.readUint16()
		obj := v.pop()
		v.push(obj.Field(int(fieldIdx)))
	case bytecode.OpSetField:
		fieldIdx := frame.readUint16()
		val := v.pop()
		obj := v.pop()
		obj.SetField(int(fieldIdx), val)
	}
}
