// ops_memory.go handles stack ops, constants, locals/globals/free vars, fields, addresses, and new.
package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// executeMemory handles global/free variable, field, address, dereference, new, and make opcodes.
// Note: OpNop, OpPop, OpDup, OpConst, OpNil, OpTrue, OpFalse, OpLocal, OpSetLocal
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
		if rv, ok := structPtr.ReflectValue(); ok {
			// Dereference pointer to get struct
			s := rv
			if s.Kind() == reflect.Ptr {
				s = s.Elem()
			}
			// For self-referencing struct types, the recursive pointer field is
			// stored as interface{} by typeToReflect. When we later access fields
			// through this pointer, the reflect.Value will be an interface wrapping
			// the actual struct pointer. Unwrap it here.
			if s.Kind() == reflect.Interface && !s.IsNil() {
				s = s.Elem()
				if s.Kind() == reflect.Ptr {
					s = s.Elem()
				}
			}
			if s.Kind() == reflect.Struct {
				field := s.Field(int(fieldIdx))
				if field.CanAddr() {
					// Use reflect.NewAt to get a settable pointer even for unexported fields.
					// This allows the VM to mutate unexported struct fields (pointer-receiver methods).
					fieldPtr := reflect.NewAt(field.Type(), value.UnsafeAddrOf(field))
					v.push(value.MakeFromReflect(fieldPtr))
				} else {
					v.push(value.MakeFromReflect(field))
				}
			} else {
				v.push(value.MakeNil())
			}
		} else {
			v.push(value.MakeNil())
		}

	case bytecode.OpIndexAddr:
		// Get address of slice/array element: &slice[index]
		index := v.pop()
		container := v.pop()
		idx := int(index.Int())

		// Native int slice: return *int64 pointer directly (avoids reflect)
		if s, ok := container.IntSlice(); ok {
			v.push(value.MakeIntPtr(&s[idx]))
			break
		}

		if rv, ok := container.ReflectValue(); ok {
			// Dereference pointer if needed
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			// Handle []value.Value slices (used for function slices)
			if rv.Kind() == reflect.Slice && rv.Type().Elem() == reflect.TypeOf(value.Value{}) {
				elem := rv.Index(idx)
				if elem.CanAddr() {
					v.push(value.MakeFromReflect(elem.Addr()))
				} else {
					v.push(value.MakeFromReflect(elem))
				}
			} else {
				// Get element address using reflect
				elem := rv.Index(idx)
				if elem.CanAddr() {
					elemPtr := elem.Addr()
					v.push(value.MakeFromReflect(elemPtr))
				} else {
					// Can't address - set directly
					v.push(value.MakeFromReflect(elem))
				}
			}
		} else {
			v.push(value.MakeNil())
		}

	case bytecode.OpDeref:
		ptr := v.pop()
		// Fast path: GlobalRef from shared-mode OpGlobal — use locked read.
		if ptr.Kind() == value.KindReflect || ptr.Kind() == value.KindInterface {
			if iface := ptr.Interface(); iface != nil {
				if ref, ok := iface.(*GlobalRef); ok {
					v.push(ref.Load())
					break
				}
			}
		}
		switch ptr.Kind() {
		case value.KindPointer:
			v.push(ptr.Elem())
		case value.KindInterface:
			// For interface values, just pass through (interfaces are already dereferenced)
			v.push(ptr)
		case value.KindReflect:
			if rv, ok := ptr.ReflectValue(); ok && rv.Kind() == reflect.Ptr {
				if !rv.IsNil() {
					// Fast path: *value.Value pointer — unwrap directly.
					if rv.CanInterface() {
						if vp, ok2 := rv.Interface().(*value.Value); ok2 {
							v.push(*vp)
							break
						}
					}
					elem := rv.Elem()
					// When dereferencing a pointer-to-pointer (e.g. **int from FieldAddr),
					// the inner pointer value is addressable and references the struct field
					// directly. We must create an independent copy so that subsequent Store
					// operations on the struct field don't silently mutate this loaded value
					// (critical for swap patterns like p.a, p.b = p.b, p.a).
					if elem.Kind() == reflect.Ptr && elem.CanSet() {
						v.push(value.MakeFromReflect(reflect.ValueOf(elem.Interface())))
					} else {
						v.push(value.MakeFromReflect(elem))
					}
				} else {
					v.push(value.MakeNil())
				}
			} else {
				v.push(ptr)
			}
		default:
			v.push(ptr)
		}

	case bytecode.OpSetDeref:
		val := v.pop()
		ptr := v.pop()
		// Fast path: GlobalRef from shared-mode OpGlobal — use locked write.
		if iface := ptr.Interface(); iface != nil {
			if ref, ok := iface.(*GlobalRef); ok {
				ref.Store(val)
				break
			}
		}
		ptr.SetElem(val)

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
