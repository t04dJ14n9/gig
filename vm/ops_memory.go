// ops_memory.go handles stack ops, constants, locals/globals/free vars, fields, addresses, and new.
package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// executeMemory handles global/free variable, field, address, dereference, and new opcodes.
// Note: OpPop, OpDup, OpConst, OpNil, OpTrue, OpFalse, OpLocal, OpSetLocal
// are inlined in run.go's hot path and never reach this handler.
func (v *vm) executeMemory(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen,maintidx
	switch op {
	case bytecode.OpGlobal:
		idx := frame.readUint16()
		if sg := v.shared; sg != nil {
			// Shared mode: push a GlobalRef that uses locked access.
			// This prevents data races from raw pointer exposure.
			if int(idx) < sg.Len() {
				ref := &GlobalRef{sg: sg, idx: int(idx)}
				v.push(value.FromInterface(ref))
			}
		} else {
			globals := v.globals
			if int(idx) < len(globals) {
				ptr := &globals[idx]
				v.push(value.FromInterface(ptr))
			}
		}

	case bytecode.OpSetGlobal:
		idx := frame.readUint16()
		val := v.pop()
		if sg := v.shared; sg != nil {
			// Shared mode: use locked write
			if int(idx) < sg.Len() {
				sg.Set(int(idx), val)
			}
		} else {
			globals := v.globals
			if int(idx) < len(globals) {
				globals[idx] = val
			}
		}

	case bytecode.OpFree:
		idx := frame.readByte()
		if int(idx) < len(frame.freeVars) && frame.freeVars[idx] != nil {
			// Push the actual value stored in the free variable slot.
			// The freeVars[idx] is a *value.Value pointer to the slot;
			// dereferencing it gives the captured value (e.g., a reflect *int pointer
			// for named return values, or any other captured variable).
			v.push(*frame.freeVars[idx])
		} else {
			v.push(value.MakeNil())
		}

	case bytecode.OpSetFree:
		idx := frame.readByte()
		val := v.pop()
		if int(idx) < len(frame.freeVars) && frame.freeVars[idx] != nil {
			*frame.freeVars[idx] = val
		}

	case bytecode.OpField:
		fieldIdx := frame.readUint16()
		obj := v.pop()
		v.push(obj.Field(int(fieldIdx)))

	case bytecode.OpSetField:
		fieldIdx := frame.readUint16()
		val := v.pop()
		obj := v.pop()
		obj.SetField(int(fieldIdx), val)

	// Pointer operations
	case bytecode.OpAddr:
		// Get address of a local variable (for taking pointer)
		localIdx := frame.readUint16()
		if int(localIdx) < len(frame.locals) {
			// Mark frame as having its address taken — cannot be pooled
			frame.addrTaken = true
			// Create a pointer to the local
			ptr := &frame.locals[localIdx]
			v.push(value.FromInterface(ptr))
		} else {
			v.push(value.MakeNil())
		}

	case bytecode.OpFieldAddr:
		// Get address of a struct field: &struct.field
		fieldIdx := frame.readUint16()
		structPtr := v.pop()
		v.push(fieldAddressValue(structPtr, int(fieldIdx)))

	case bytecode.OpIndexAddr:
		// Get address of slice/array element: &slice[index]
		index := v.pop()
		container := v.pop()
		v.push(indexAddressValue(container, int(index.Int())))

	case bytecode.OpDeref:
		ptr := v.pop()
		v.push(dereferenceValue(ptr))

	case bytecode.OpSetDeref:
		val := v.pop()
		ptr := v.pop()
		v.setDereferenceValue(ptr, val)

	case bytecode.OpNew:
		typeIdx := frame.readUint16()
		if int(typeIdx) < len(v.program.Types) {
			typ := v.program.Types[typeIdx]
			// Function types need a pointer to a Value (to store closures)
			if _, isSig := typ.(*types.Signature); isSig {
				var nilVal value.Value
				v.push(value.MakeFromReflect(reflect.ValueOf(&nilVal)))
			} else if rt := typeToReflect(typ, v.program); rt != nil {
				v.push(value.MakeFromReflect(reflect.New(rt)))
			} else {
				v.push(value.MakeNil())
			}
		} else {
			v.push(value.MakeNil())
		}
	}

	return nil
}
