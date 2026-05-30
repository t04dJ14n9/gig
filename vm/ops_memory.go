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
		// Unwrap *value.Value from OpGlobal slot pointer.
		if structPtr.IsValid() && !structPtr.IsNil() && structPtr.CanInterface() {
			if iface := structPtr.Interface(); iface != nil {
				if vp, ok := iface.(*value.Value); ok {
					structPtr = *vp
				}
			}
		}

		if rv, ok := structPtr.ReflectValue(); ok {
			// Dereference pointer to get struct
			s := rv
			if s.Kind() == reflect.Ptr {
				// Check if this is a *value.Value (pointer to a global slot).
				// If so, unwrap it to get the underlying reflect.Value stored
				// inside the value.Value, which is the actual struct pointer.
				if s.CanInterface() {
					if vp, ok2 := s.Interface().(*value.Value); ok2 {
						// Dereference *value.Value to get the Value
						if innerRV, ok3 := vp.ReflectValue(); ok3 {
							s = innerRV
							if s.Kind() == reflect.Ptr {
								s = s.Elem()
							}
						} else {
							v.push(value.MakeNil())
							break
						}
					} else {
						s = s.Elem()
					}
				} else {
					s = s.Elem()
				}
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
				// If s is a value.Value struct (from OpGlobal slot pointer),
				// extract the wrapped value before accessing fields.
				if s.Type() == reflect.TypeOf(value.Value{}) {
					if innerRV, ok2 := s.Interface().(value.Value).ReflectValue(); ok2 && innerRV.IsValid() {
						s = innerRV
						if s.Kind() == reflect.Ptr {
							s = s.Elem()
						}
					}
				}
				if s.Kind() == reflect.Struct {
					field := s.Field(int(fieldIdx))
					if field.CanAddr() {
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
		} else {
			v.push(value.MakeNil())
		}

	case bytecode.OpIndexAddr:
		// Get address of slice/array element: &slice[index]
		index := v.pop()
		container := v.pop()

		// Unwrap *value.Value from OpGlobal slot pointer.
		if container.IsValid() && !container.IsNil() && container.CanInterface() {
			if iface := container.Interface(); iface != nil {
				if vp, ok := iface.(*value.Value); ok {
					container = *vp
				}
			}
		}
		idx := int(index.Int())

		// Native int slice: return *int64 pointer directly (avoids reflect)
		if s, ok := container.IntSlice(); ok {
			v.push(value.MakeIntPtr(&s[idx]))
			break
		}

		// Native []byte: convert to reflect.Value so the reflect path can handle it
		// KindBytes stores []byte as obj (not reflect.Value), so ReflectValue() fails.
		if container.Kind() == value.KindBytes {
			if b, ok := container.Bytes(); ok {
				rv := reflect.ValueOf(b)
				if idx >= 0 && idx < len(b) {
					elem := rv.Index(idx)
					if elem.CanAddr() {
						v.push(value.MakeFromReflect(elem.Addr()))
					} else {
						// Non-addressable: create a settable copy pointer
						elemPtr := reflect.New(elem.Type())
						elemPtr.Elem().Set(elem)
						v.push(value.MakeFromReflect(elemPtr))
					}
				} else {
					v.push(value.MakeNil())
				}
			} else {
				v.push(value.MakeNil())
			}
			break
		}

		if rv, ok := container.ReflectValue(); ok {
			// Dereference pointer if needed
			if rv.Kind() == reflect.Ptr {
				// Check if this is a *value.Value (pointer to a global slot).
				if rv.CanInterface() {
					if vp, ok2 := rv.Interface().(*value.Value); ok2 {
						if innerRV, ok3 := vp.ReflectValue(); ok3 {
							rv = innerRV
							if rv.Kind() == reflect.Ptr {
								rv = rv.Elem()
							}
						} else {
							v.push(value.MakeNil())
							break
						}
					} else {
						rv = rv.Elem()
					}
				} else {
					rv = rv.Elem()
				}
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
		if (ptr.Kind() == value.KindReflect || ptr.Kind() == value.KindInterface) && ptr.CanInterface() {
			if iface := ptr.Interface(); iface != nil {
				if ref, ok := iface.(*GlobalRef); ok {
					v.push(ref.Load())
					break
				}
			}
		}
		switch ptr.Kind() {
		case value.KindPointer:
			if ptr.Elem().IsValid() {
				v.push(ptr.Elem())
			} else {
				// Nil pointer dereference
				panic("runtime error: invalid memory address or nil pointer dereference")
			}
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
					// When dereferencing a pointer that points into addressable
					// memory (slice element, struct field, array element), the
					// resulting reflect.Value aliases the original storage.
					// We must create an independent copy so subsequent Store
					// operations don't silently mutate this loaded value
					// (critical for swap patterns like a[i],a[j]=a[j],a[i]).
					if elem.Kind() == reflect.Ptr && elem.CanSet() {
						v.push(value.MakeFromReflect(reflect.ValueOf(elem.Interface())))
					} else if elem.Kind() == reflect.Interface && elem.CanSet() && elem.Type().NumMethod() == 0 {
						if elem.IsNil() {
							v.push(value.MakeNil())
						} else {
							concrete := elem.Elem()
							v.push(value.MakeFromReflect(reflect.ValueOf(concrete.Interface())))
						}
					} else if elem.CanAddr() {
						v.push(value.MakeFromReflect(cloneReflectValue(elem)))
					} else {
						v.push(value.MakeFromReflect(elem))
					}
				} else {
					// Nil pointer dereference — panic, matching Go semantics.
					panic("runtime error: invalid memory address or nil pointer dereference")
				}
			} else {
				v.push(ptr)
			}
		default:
			// Nil/invalid pointer dereference — panic, matching Go semantics.
			// This catches KindNil and KindInvalid being dereferenced.
			if ptr.IsNil() || !ptr.IsValid() {
				panic("runtime error: invalid memory address or nil pointer dereference")
			}
			v.push(ptr)
		}

	case bytecode.OpSetDeref:
		val := v.pop()
		ptr := v.pop()
		// Nil pointer dereference check
		if ptr.IsNil() || !ptr.IsValid() {
			panic("runtime error: invalid memory address or nil pointer dereference")
		}
		// Fast path: GlobalRef from shared-mode OpGlobal — use locked write.
		if ptr.CanInterface() {
			if iface := ptr.Interface(); iface != nil {
				if ref, ok := iface.(*GlobalRef); ok {
					ref.Store(val)
					break
				}
			}
		}
		if rv, ok := ptr.ReflectValue(); ok && rv.Kind() == reflect.Ptr && !rv.IsNil() {
			elem := rv.Elem()
			if elem.IsValid() && elem.CanSet() && elem.Kind() == reflect.Interface {
				elem.Set(v.valueForReflectSet(val, elem.Type()))
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

// cloneReflectValue creates an independent copy of a reflect.Value that
// references addressable memory (slice element, struct field, etc.).
// This breaks the alias so subsequent writes through the original
// pointer don't corrupt the copy.
func cloneReflectValue(rv reflect.Value) reflect.Value {
	copy := reflect.New(rv.Type()).Elem()
	copy.Set(rv)
	return copy
}
