// ops_memory.go handles stack ops, constants, locals/globals/free vars, fields, addresses, and new.
package vm

import (
	"go/types"
	"reflect"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
)

// executeMemory handles stack, constant, local/global/free variable, field,
// address, dereference, new, and make opcodes.
func (v *vm) executeMemory(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen,maintidx
	switch op {
	// Stack operations
	case bytecode.OpNop:
		// No operation

	case bytecode.OpPop:
		v.pop()

	case bytecode.OpDup:
		val := v.peek()
		v.push(val)

	// Constants and locals
	case bytecode.OpConst:
		idx := frame.readUint16()
		if int(idx) < len(v.program.PrebakedConstants) {
			v.push(v.program.PrebakedConstants[idx])
		} else if int(idx) < len(v.program.Constants) {
			v.push(value.FromInterface(v.program.Constants[idx]))
		}

	case bytecode.OpNil:
		v.push(value.MakeNil())

	case bytecode.OpTrue:
		v.push(value.MakeBool(true))

	case bytecode.OpFalse:
		v.push(value.MakeBool(false))

	case bytecode.OpLocal:
		idx := frame.readUint16()
		if int(idx) < len(frame.locals) {
			v.push(frame.locals[idx])
		}

	case bytecode.OpSetLocal:
		idx := frame.readUint16()
		val := v.pop()
		if int(idx) < len(frame.locals) {
			frame.locals[idx] = val
		}

	case bytecode.OpGlobal:
		idx := frame.readUint16()
		globals := v.getGlobals()
		if int(idx) < len(globals) {
			// Push a pointer to the global slot
			// This allows OpDeref/OpSetDeref to work correctly
			ptr := &globals[idx]
			v.push(value.FromInterface(ptr))
		}

	case bytecode.OpSetGlobal:
		idx := frame.readUint16()
		val := v.pop()
		globals := v.getGlobals()
		if int(idx) < len(globals) {
			globals[idx] = val
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
		ptr.SetElem(val)

	case bytecode.OpNew:
		typeIdx := frame.readUint16()
		// Allocate new pointer to value of the given type
		if int(typeIdx) < len(v.program.Types) {
			typ := v.program.Types[typeIdx]
			// For function types, create a pointer to a Value (to store closures)
			switch t := typ.(type) {
			case *types.Signature:
				_ = t // Function signature not needed for allocation
				// Create a pointer to a nil Value
				var nilVal value.Value
				newPtr := reflect.ValueOf(&nilVal)
				v.push(value.MakeFromReflect(newPtr))
			case *types.Slice:
				// Use typeToReflect for proper typed slices (including function slices).
				// This avoids creating []value.Value which can't be assigned to typed fields.
				if rt := typeToReflect(typ, v.program); rt != nil {
					newPtr := reflect.New(rt)
					v.push(value.MakeFromReflect(newPtr))
				} else {
					v.push(value.MakeNil())
				}
			case *types.Array:
				// Use typeToReflect for proper typed arrays (including function arrays).
				if rt := typeToReflect(typ, v.program); rt != nil {
					newPtr := reflect.New(rt)
					v.push(value.MakeFromReflect(newPtr))
				} else {
					v.push(value.MakeNil())
				}
			default:
				if rt := typeToReflect(typ, v.program); rt != nil {
					// Create a new pointer to zero value of the type
					newPtr := reflect.New(rt)
					v.push(value.MakeFromReflect(newPtr))
				} else {
					v.push(value.MakeNil())
				}
			}
		} else {
			v.push(value.MakeNil())
		}
	}

	return nil
}
