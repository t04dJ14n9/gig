// ops_convert.go dispatches VM conversion and interface opcodes.
package vm

import "github.com/t04dJ14n9/gig/model/bytecode"

// executeConvert handles type assertion, conversion, and change-type opcodes.
func (v *vm) executeConvert(op bytecode.OpCode, frame *Frame) error {
	switch op {
	case bytecode.OpAssert:
		v.executeAssert(frame)
	case bytecode.OpConvert:
		v.executeTypeConvert(frame)
	case bytecode.OpChangeType:
		v.executeChangeType(frame)
	case bytecode.OpMakeInterface:
		v.executeMakeInterface(frame)
	}
	return nil
}
