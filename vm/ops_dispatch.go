package vm

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/bytecode"
	"github.com/t04dJ14n9/gig/value"
)

// executeOp executes a single bytecode instruction.
// This is the heart of the VM - it dispatches to the appropriate handler
// for each opcode type.
func (vm *VM) executeOp(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen,maintidx
	switch op {
	// Stack operations
	case bytecode.OpNop:
		// No operation

	case bytecode.OpPop:
		vm.pop()

	case bytecode.OpDup:
		val := vm.peek()
		vm.push(val)

	// Constants and locals
	case bytecode.OpConst:
		idx := frame.readUint16()
		if int(idx) < len(vm.program.PrebakedConstants) {
			vm.push(vm.program.PrebakedConstants[idx])
		} else if int(idx) < len(vm.program.Constants) {
			vm.push(value.FromInterface(vm.program.Constants[idx]))
		}

	case bytecode.OpNil:
		vm.push(value.MakeNil())

	case bytecode.OpTrue:
		vm.push(value.MakeBool(true))

	case bytecode.OpFalse:
		vm.push(value.MakeBool(false))

	case bytecode.OpLocal:
		idx := frame.readUint16()
		if int(idx) < len(frame.locals) {
			vm.push(frame.locals[idx])
		}

	case bytecode.OpSetLocal:
		idx := frame.readUint16()
		val := vm.pop()
		if int(idx) < len(frame.locals) {
			frame.locals[idx] = val
		}

	case bytecode.OpGlobal:
		idx := frame.readUint16()
		globals := vm.getGlobals()
		if int(idx) < len(globals) {
			// Push a pointer to the global slot
			// This allows OpDeref/OpSetDeref to work correctly
			ptr := &globals[idx]
			vm.push(value.FromInterface(ptr))
		}

	case bytecode.OpSetGlobal:
		idx := frame.readUint16()
		val := vm.pop()
		globals := vm.getGlobals()
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
			vm.push(*frame.freeVars[idx])
		} else {
			vm.push(value.MakeNil())
		}

	case bytecode.OpSetFree:
		idx := frame.readByte()
		val := vm.pop()
		if int(idx) < len(frame.freeVars) && frame.freeVars[idx] != nil {
			*frame.freeVars[idx] = val
		}

	// Arithmetic
	case bytecode.OpAdd:
		b := vm.pop()
		a := vm.pop()
		// Fast path for int+int (most common case in loops)
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			vm.push(value.MakeInt(a.RawInt() + b.RawInt()))
		} else {
			vm.push(a.Add(b))
		}

	case bytecode.OpSub:
		b := vm.pop()
		a := vm.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			vm.push(value.MakeInt(a.RawInt() - b.RawInt()))
		} else {
			vm.push(a.Sub(b))
		}

	case bytecode.OpMul:
		b := vm.pop()
		a := vm.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			vm.push(value.MakeInt(a.RawInt() * b.RawInt()))
		} else {
			vm.push(a.Mul(b))
		}

	case bytecode.OpDiv:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Div(b))

	case bytecode.OpMod:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Mod(b))

	case bytecode.OpNeg:
		a := vm.pop()
		vm.push(a.Neg())

	// Bitwise
	case bytecode.OpAnd:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.And(b))

	case bytecode.OpOr:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Or(b))

	case bytecode.OpXor:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Xor(b))

	case bytecode.OpAndNot:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.AndNot(b))

	case bytecode.OpLsh:
		n := uint(vm.pop().Int())
		a := vm.pop()
		vm.push(a.Lsh(n))

	case bytecode.OpRsh:
		n := uint(vm.pop().Int())
		a := vm.pop()
		vm.push(a.Rsh(n))

	// Comparison
	case bytecode.OpEqual:
		b := vm.pop()
		a := vm.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			vm.push(value.MakeBool(a.RawInt() == b.RawInt()))
		} else {
			vm.push(value.MakeBool(a.Equal(b)))
		}

	case bytecode.OpNotEqual:
		b := vm.pop()
		a := vm.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			vm.push(value.MakeBool(a.RawInt() != b.RawInt()))
		} else {
			vm.push(value.MakeBool(!a.Equal(b)))
		}

	case bytecode.OpLess:
		b := vm.pop()
		a := vm.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			vm.push(value.MakeBool(a.RawInt() < b.RawInt()))
		} else {
			vm.push(value.MakeBool(a.Cmp(b) < 0))
		}

	case bytecode.OpLessEq:
		b := vm.pop()
		a := vm.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			vm.push(value.MakeBool(a.RawInt() <= b.RawInt()))
		} else {
			vm.push(value.MakeBool(a.Cmp(b) <= 0))
		}

	case bytecode.OpGreater:
		b := vm.pop()
		a := vm.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			vm.push(value.MakeBool(a.RawInt() > b.RawInt()))
		} else {
			vm.push(value.MakeBool(a.Cmp(b) > 0))
		}

	case bytecode.OpGreaterEq:
		b := vm.pop()
		a := vm.pop()
		if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
			vm.push(value.MakeBool(a.RawInt() >= b.RawInt()))
		} else {
			vm.push(value.MakeBool(a.Cmp(b) >= 0))
		}

	// Logical
	case bytecode.OpNot:
		a := vm.pop()
		vm.push(value.MakeBool(!a.Bool()))

	// Control flow
	case bytecode.OpJump:
		offset := frame.readUint16()
		frame.ip = int(offset)

	case bytecode.OpJumpTrue:
		offset := frame.readUint16()
		cond := vm.pop()
		if cond.Bool() {
			frame.ip = int(offset)
		}

	case bytecode.OpJumpFalse:
		offset := frame.readUint16()
		cond := vm.pop()
		if !cond.Bool() {
			frame.ip = int(offset)
		}

	case bytecode.OpCall:
		funcIdx := frame.readUint16()
		numArgs := frame.readByte()
		vm.callCompiledFunction(int(funcIdx), int(numArgs))

	case bytecode.OpReturn:
		// Pop frame and return it to pool
		vm.fpool.put(frame)
		vm.fp--
		if vm.fp > 0 {
			// Restore stack pointer
			prevFrame := vm.frames[vm.fp-1]
			vm.sp = prevFrame.basePtr
		}
		// Push nil for void returns (SSA may expect a value)
		vm.push(value.MakeNil())

	case bytecode.OpReturnVal:
		// Save return value
		retVal := vm.pop()
		// Pop frame and return it to pool
		vm.fpool.put(frame)
		vm.fp--
		if vm.fp > 0 {
			// Restore stack pointer
			prevFrame := vm.frames[vm.fp-1]
			vm.sp = prevFrame.basePtr
		}
		// Push return value
		vm.push(retVal)

	// Container operations
	case bytecode.OpMakeSlice:
		capVal := vm.pop()
		lenVal := vm.pop()
		typeIdxVal := vm.pop()
		typeIdx := uint16(typeIdxVal.Int())
		// Create slice using the type from the type pool
		if int(typeIdx) < len(vm.program.Types) {
			typ := vm.program.Types[typeIdx]
			made := false
			if sliceType, ok := typ.(*types.Slice); ok {
				elemType := sliceType.Elem()
				// Native []int64 fast path for integer slice types
				if basic, isBasic := elemType.(*types.Basic); isBasic {
					switch basic.Kind() {
					case types.Int, types.Int64:
						vm.push(value.MakeIntSlice(make([]int64, int(lenVal.Int()), int(capVal.Int()))))
						made = true
					}
				}
				// Function slice special case
				if !made {
					if _, isFunc := elemType.(*types.Signature); isFunc {
						slice := make([]value.Value, int(lenVal.Int()), int(capVal.Int()))
						vm.push(value.FromInterface(slice))
						made = true
					}
				}
			}
			if !made {
				if rt := typeToReflect(typ); rt != nil {
					slice := reflect.MakeSlice(rt, int(lenVal.Int()), int(capVal.Int()))
					vm.push(value.MakeFromReflect(slice))
				} else {
					vm.push(value.MakeNil())
				}
			}
		} else {
			vm.push(value.MakeNil())
		}

	case bytecode.OpMakeMap:
		sizeVal := vm.pop()
		typeIdxVal := vm.pop()
		typeIdx := uint16(typeIdxVal.Int())
		if int(typeIdx) < len(vm.program.Types) {
			typ := vm.program.Types[typeIdx]
			if rt := typeToReflect(typ); rt != nil {
				m := reflect.MakeMap(rt)
				_ = sizeVal // Size hint ignored for simplicity
				vm.push(value.MakeFromReflect(m))
			} else {
				vm.push(value.MakeNil())
			}
		} else {
			vm.push(value.MakeNil())
		}

	case bytecode.OpMakeChan:
		sizeVal := vm.pop()
		typeIdxVal := vm.pop()
		typeIdx := uint16(typeIdxVal.Int())
		if int(typeIdx) < len(vm.program.Types) {
			typ := vm.program.Types[typeIdx]
			if rt := typeToReflect(typ); rt != nil {
				ch := reflect.MakeChan(rt, int(sizeVal.Int()))
				vm.push(value.MakeFromReflect(ch))
			} else {
				vm.push(value.MakeNil())
			}
		} else {
			vm.push(value.MakeNil())
		}

	// Index operations
	case bytecode.OpIndex:
		key := vm.pop()
		container := vm.pop()
		switch container.Kind() {
		case value.KindSlice:
			// Native int slice fast path
			if s, ok := container.IntSlice(); ok {
				vm.push(value.MakeInt(s[int(key.RawInt())]))
			} else {
				vm.push(container.Index(int(key.Int())))
			}
		case value.KindArray:
			idx := int(key.Int())
			vm.push(container.Index(idx))
		case value.KindMap:
			vm.push(container.MapIndex(key))
		case value.KindString:
			idx := int(key.Int())
			vm.push(container.Index(idx))
		case value.KindReflect:
			// Handle reflect.Value containing a slice, array, or map
			if rv, ok := container.ReflectValue(); ok {
				switch rv.Kind() {
				case reflect.Slice, reflect.Array:
					idx := int(key.Int())
					vm.push(value.MakeFromReflect(rv.Index(idx)))
				case reflect.Map:
					k := key.ToReflectValue(rv.Type().Key())
					elem := rv.MapIndex(k)
					if !elem.IsValid() {
						// Return zero value of element type, not nil (Go semantics)
						vm.push(value.MakeFromReflect(reflect.Zero(rv.Type().Elem())))
					} else {
						vm.push(value.MakeFromReflect(elem))
					}
				default:
					vm.push(value.MakeNil())
				}
			} else {
				vm.push(value.MakeNil())
			}
		default:
			vm.push(value.MakeNil())
		}

	case bytecode.OpIndexOk:
		// Index with comma-ok: returns (value, ok) tuple for maps
		key := vm.pop()
		container := vm.pop()
		switch container.Kind() {
		case value.KindMap:
			// For maps, check if key exists
			if rv, ok := container.ReflectValue(); ok {
				k := key.ToReflectValue(rv.Type().Key())
				elem := rv.MapIndex(k)
				if !elem.IsValid() {
					// Key doesn't exist: return (zero_value, false)
					zeroVal := value.MakeFromReflect(reflect.Zero(rv.Type().Elem()))
					tuple := []value.Value{zeroVal, value.MakeBool(false)}
					vm.push(value.FromInterface(tuple))
				} else {
					// Key exists: return (value, true)
					vm.push(value.FromInterface([]value.Value{value.MakeFromReflect(elem), value.MakeBool(true)}))
				}
			} else {
				tuple := []value.Value{value.MakeNil(), value.MakeBool(false)}
				vm.push(value.FromInterface(tuple))
			}
		case value.KindReflect:
			if rv, ok := container.ReflectValue(); ok {
				switch rv.Kind() {
				case reflect.Map:
					k := key.ToReflectValue(rv.Type().Key())
					elem := rv.MapIndex(k)
					if !elem.IsValid() {
						zeroVal := value.MakeFromReflect(reflect.Zero(rv.Type().Elem()))
						tuple := []value.Value{zeroVal, value.MakeBool(false)}
						vm.push(value.FromInterface(tuple))
					} else {
						vm.push(value.FromInterface([]value.Value{value.MakeFromReflect(elem), value.MakeBool(true)}))
					}
				case reflect.Slice, reflect.Array:
					// For slices/arrays, always return true for ok
					idx := int(key.Int())
					vm.push(value.FromInterface([]value.Value{value.MakeFromReflect(rv.Index(idx)), value.MakeBool(true)}))
				default:
					tuple := []value.Value{value.MakeNil(), value.MakeBool(false)}
					vm.push(value.FromInterface(tuple))
				}
			} else {
				tuple := []value.Value{value.MakeNil(), value.MakeBool(false)}
				vm.push(value.FromInterface(tuple))
			}
		default:
			tuple := []value.Value{value.MakeNil(), value.MakeBool(false)}
			vm.push(value.FromInterface(tuple))
		}

	case bytecode.OpSetIndex:
		val := vm.pop()
		key := vm.pop()
		container := vm.pop()
		switch container.Kind() {
		case value.KindSlice:
			// Native int slice fast path
			if s, ok := container.IntSlice(); ok {
				s[int(key.RawInt())] = val.RawInt()
			} else {
				container.SetIndex(int(key.Int()), val)
			}
		case value.KindArray:
			idx := int(key.Int())
			container.SetIndex(idx, val)
		case value.KindMap:
			container.SetMapIndex(key, val)
		case value.KindReflect:
			if rv, ok := container.ReflectValue(); ok {
				switch rv.Kind() {
				case reflect.Slice, reflect.Array:
					idx := int(key.Int())
					rv.Index(idx).Set(val.ToReflectValue(rv.Type().Elem()))
				case reflect.Map:
					container.SetMapIndex(key, val)
				}
			}
		}

	case bytecode.OpSlice:
		// Slice operation: container[low:high:max]
		maxVal := vm.pop()
		highVal := vm.pop()
		lowVal := vm.pop()
		container := vm.pop()

		low := int(lowVal.Int())

		// Handle string slicing specially
		if container.Kind() == value.KindString {
			high := int(highVal.Int())
			if high == 0xFFFF {
				high = len(container.String())
			}
			vm.push(value.MakeString(container.String()[low:high]))
			break
		}

		// Native []int64 slice fast path
		if s, ok := container.IntSlice(); ok {
			high := int(highVal.Int())
			if high == 0xFFFF {
				high = len(s)
			}
			if maxVal.Kind() != value.KindNil && maxVal.Int() != 0xFFFF {
				vm.push(value.MakeIntSlice(s[low:high:int(maxVal.Int())]))
			} else {
				vm.push(value.MakeIntSlice(s[low:high]))
			}
			break
		}

		if rv, ok := container.ReflectValue(); ok {
			// If it's a pointer to an array or slice, dereference it first
			if rv.Kind() == reflect.Ptr {
				elemKind := rv.Elem().Kind()
				if elemKind == reflect.Array || elemKind == reflect.Slice {
					rv = rv.Elem()
				}
			}

			// Native int array/slice → []int64 fast path
			// SSA compiles make([]int, N) with constant N as Alloc([N]int) + Slice,
			// so we intercept it here to produce a native []int64.
			if (rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice) && rv.Type().Elem().Kind() == reflect.Int {
				n := rv.Len()
				high := int(highVal.Int())
				if high == 0xFFFF {
					high = n
				}
				s := make([]int64, n)
				for i := 0; i < n; i++ {
					s[i] = rv.Index(i).Int()
				}
				if maxVal.Kind() != value.KindNil && maxVal.Int() != 0xFFFF {
					vm.push(value.MakeIntSlice(s[low:high:int(maxVal.Int())]))
				} else {
					vm.push(value.MakeIntSlice(s[low:high]))
				}
				break
			}

			// Handle native []value.Value slices (used for function slices)
			if rv.Kind() == reflect.Slice && rv.Type().Elem() == reflect.TypeOf(value.Value{}) {
				high := int(highVal.Int())
				if high == 0xFFFF {
					high = rv.Len()
				}
				sliced := rv.Slice(low, high)
				vm.push(value.MakeFromReflect(sliced))
				break
			}

			high := int(highVal.Int())
			if high == 0xFFFF {
				high = rv.Len()
			}

			var sliced reflect.Value
			if maxVal.Kind() != value.KindNil && maxVal.Int() != 0xFFFF {
				// 3-index slice: container[low:high:max]
				max := int(maxVal.Int())
				sliced = rv.Slice3(low, high, max)
			} else {
				// 2-index slice: container[low:high]
				sliced = rv.Slice(low, high)
			}
			vm.push(value.MakeFromReflect(sliced))
		} else {
			vm.push(value.MakeNil())
		}

	case bytecode.OpField:
		fieldIdx := frame.readUint16()
		obj := vm.pop()
		vm.push(obj.Field(int(fieldIdx)))

	case bytecode.OpSetField:
		fieldIdx := frame.readUint16()
		val := vm.pop()
		obj := vm.pop()
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
			vm.push(value.FromInterface(ptr))
		} else {
			vm.push(value.MakeNil())
		}

	case bytecode.OpFieldAddr:
		// Get address of a struct field: &struct.field
		fieldIdx := frame.readUint16()
		structPtr := vm.pop()
		if rv, ok := structPtr.ReflectValue(); ok {
			// Dereference pointer to get struct
			s := rv
			if s.Kind() == reflect.Ptr {
				s = s.Elem()
			}
			if s.Kind() == reflect.Struct {
				field := s.Field(int(fieldIdx))
				if field.CanAddr() {
					// Use reflect.NewAt to get a settable pointer even for unexported fields.
					// This allows the VM to mutate unexported struct fields (pointer-receiver methods).
					fieldPtr := reflect.NewAt(field.Type(), value.UnsafeAddrOf(field))
					vm.push(value.MakeFromReflect(fieldPtr))
				} else {
					vm.push(value.MakeFromReflect(field))
				}
			} else {
				vm.push(value.MakeNil())
			}
		} else {
			vm.push(value.MakeNil())
		}

	case bytecode.OpIndexAddr:
		// Get address of slice/array element: &slice[index]
		index := vm.pop()
		container := vm.pop()
		idx := int(index.Int())

		// Native int slice: return *int64 pointer directly (avoids reflect)
		if s, ok := container.IntSlice(); ok {
			vm.push(value.MakeIntPtr(&s[idx]))
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
					vm.push(value.MakeFromReflect(elem.Addr()))
				} else {
					vm.push(value.MakeFromReflect(elem))
				}
			} else {
				// Get element address using reflect
				elem := rv.Index(idx)
				if elem.CanAddr() {
					elemPtr := elem.Addr()
					vm.push(value.MakeFromReflect(elemPtr))
				} else {
					// Can't address - set directly
					vm.push(value.MakeFromReflect(elem))
				}
			}
		} else {
			vm.push(value.MakeNil())
		}

	case bytecode.OpDeref:
		ptr := vm.pop()
		switch ptr.Kind() {
		case value.KindPointer:
			vm.push(ptr.Elem())
		case value.KindInterface:
			// For interface values, just pass through (interfaces are already dereferenced)
			vm.push(ptr)
		case value.KindReflect:
			if rv, ok := ptr.ReflectValue(); ok && rv.Kind() == reflect.Ptr {
				if !rv.IsNil() {
					// Fast path: *value.Value pointer — unwrap directly.
					if rv.CanInterface() {
						if vp, ok2 := rv.Interface().(*value.Value); ok2 {
							vm.push(*vp)
							break
						}
					}
					vm.push(value.MakeFromReflect(rv.Elem()))
				} else {
					vm.push(value.MakeNil())
				}
			} else {
				vm.push(ptr)
			}
		default:
			vm.push(ptr)
		}

	case bytecode.OpSetDeref:
		val := vm.pop()
		ptr := vm.pop()
		ptr.SetElem(val)

	// Type operations
	case bytecode.OpAssert:
		typeIdx := frame.readUint16()
		targetType := vm.program.Types[typeIdx]
		obj := vm.pop()

		// Type assertion - check if obj can be asserted to targetType
		// Returns (value, ok) tuple on stack
		var result value.Value
		var assertionOk bool

		if obj.Kind() == value.KindInterface {
			// Get the underlying interface
			if rv, isReflect := obj.ReflectValue(); isReflect && rv.Kind() == reflect.Interface {
				if rv.IsNil() {
					// Interface is nil, assertion fails
					result = value.MakeNil()
					assertionOk = false
				} else {
					// Get the underlying value
					underlying := rv.Elem()
					targetReflectType := typeToReflect(targetType)
					if targetReflectType == nil {
						result = value.MakeNil()
						assertionOk = false
					} else if underlying.Type().AssignableTo(targetReflectType) {
						// Successful assertion - create a value from the underlying
						result = value.MakeFromReflect(underlying)
						assertionOk = true
					} else {
						result = value.MakeNil()
						assertionOk = false
					}
				}
			} else {
				result = obj
				assertionOk = true
			}
		} else if obj.Kind() == value.KindReflect {
			// Already a reflect value
			if rv, isReflect := obj.ReflectValue(); isReflect {
				targetReflectType := typeToReflect(targetType)
				if targetReflectType != nil && rv.Type().AssignableTo(targetReflectType) {
					result = obj
					assertionOk = true
				} else {
					result = value.MakeNil()
					assertionOk = false
				}
			} else {
				result = obj
				assertionOk = false
			}
		} else {
			// For other kinds, assume success
			result = obj
			assertionOk = true
		}

		// Push result as a tuple [result, ok]
		// Use a slice to represent the tuple
		tuple := []value.Value{result, value.MakeBool(assertionOk)}
		vm.push(value.FromInterface(tuple))

	case bytecode.OpConvert:
		typeIdx := frame.readUint16()
		targetType := vm.program.Types[typeIdx]
		val := vm.pop()

		// Handle type conversion
		switch t := targetType.(type) {
		case *types.Basic:
			switch t.Kind() {
			case types.String:
				// Convert to string
				switch val.Kind() {
				case value.KindInt:
					// int -> string: convert rune to string
					vm.push(value.MakeString(string(rune(val.Int()))))
				case value.KindUint:
					// byte/uint8 -> string: convert byte to string
					vm.push(value.MakeString(string(byte(val.Uint()))))
				case value.KindString:
					vm.push(val)
				case value.KindBytes:
					if b, ok := val.Bytes(); ok {
						vm.push(value.MakeString(string(b)))
					} else {
						vm.push(value.MakeString(""))
					}
				default:
					// Use reflection for other types
					vm.push(value.MakeString(fmt.Sprintf("%v", val.Interface())))
				}
			case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
				// Handle conversion from various types to int
				switch val.Kind() {
				case value.KindInt:
					vm.push(val)
				case value.KindUint:
					vm.push(value.MakeInt(int64(val.Uint())))
				case value.KindFloat:
					vm.push(value.MakeInt(int64(val.Float())))
				default:
					vm.push(value.MakeInt(val.Int()))
				}
			case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
				// Handle conversion from various types to uint
				switch val.Kind() {
				case value.KindInt:
					vm.push(value.MakeUint(uint64(val.Int())))
				case value.KindUint:
					vm.push(val)
				case value.KindFloat:
					vm.push(value.MakeUint(uint64(val.Float())))
				default:
					vm.push(value.MakeUint(val.Uint()))
				}
			case types.Float32, types.Float64:
				// Handle conversion from int/uint to float
				switch val.Kind() {
				case value.KindInt:
					vm.push(value.MakeFloat(float64(val.Int())))
				case value.KindUint:
					vm.push(value.MakeFloat(float64(val.Uint())))
				case value.KindFloat:
					vm.push(val)
				default:
					vm.push(value.MakeFloat(val.Float()))
				}
			default:
				vm.push(val)
			}
		default:
			// For non-basic types, just pass through for now
			vm.push(val)
		}

	// Function operations
	case bytecode.OpClosure:
		funcIdx := frame.readUint16()
		numFree := frame.readByte()
		// Look up the function by index (O(1))
		var fn *bytecode.CompiledFunction
		if int(funcIdx) < len(vm.program.FuncByIndex) {
			fn = vm.program.FuncByIndex[funcIdx]
		}
		if fn != nil {
			closure := getClosure(fn, int(numFree))
			// Get free variables (popped in reverse order)
			for i := int(numFree) - 1; i >= 0; i-- {
				v := vm.pop()
				// Create a new *value.Value slot holding the captured value.
				// This allows the closure to read/write the slot via OpFree/OpSetFree.
				// If the captured value is a reflect pointer (e.g., *int from Alloc),
				// all closures sharing that pointer will see each other's modifications.
				slot := new(value.Value)
				*slot = v
				closure.FreeVars[i] = slot
			}
			vm.push(value.MakeFunc(closure))
		} else {
			// Still need to pop free vars to keep stack balanced
			for i := 0; i < int(numFree); i++ {
				vm.pop()
			}
			vm.push(value.MakeNil())
		}

	// Concurrency
	case bytecode.OpGoCall:
		// OpGoCall spawns a new goroutine to execute a function call.
		// Operands: [func_idx:2, num_args:1]
		// Stack: [... args] -> [...] (arguments consumed)
		funcIdx := frame.readUint16()
		numArgs := frame.readByte()

		// Pop arguments from current goroutine's stack
		args := make([]value.Value, numArgs)
		for i := int(numArgs) - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}

		// Get the function to call (O(1))
		var goFn *bytecode.CompiledFunction
		if int(funcIdx) < len(vm.program.FuncByIndex) {
			goFn = vm.program.FuncByIndex[funcIdx]
		}

		if goFn != nil {
			// Create a child VM with shared globals
			childVM := vm.newChildVM()

			// Capture for closure
			capturedFn := goFn
			capturedArgs := args

			// Track the goroutine
			StartGoroutine(func() {
				// Create initial frame for the child goroutine
				childFrame := newFrame(capturedFn, 0, capturedArgs, nil)
				childVM.frames[0] = childFrame
				childVM.fp = 1

				// Run the child VM (ignore return value - goroutine result is discarded)
				_, _ = childVM.run()
			})
		}

	case bytecode.OpGoCallIndirect:
		// OpGoCallIndirect spawns a new goroutine to execute a closure call.
		// Operands: [num_args:1]
		// Stack: [... closure args...] -> [...] (closure and arguments consumed)
		numArgs := frame.readByte()

		// Pop arguments from current goroutine's stack
		args := make([]value.Value, numArgs)
		for i := int(numArgs) - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}

		// Pop the closure
		callee := vm.pop()

		if closure, ok := callee.RawObj().(*Closure); ok {
			// Create a child VM with shared globals
			childVM := vm.newChildVM()

			// Capture for closure
			capturedClosure := closure
			capturedArgs := args

			// Track the goroutine
			StartGoroutine(func() {
				// Create initial frame for the child goroutine with free vars
				childFrame := newFrame(capturedClosure.Fn, 0, capturedArgs, capturedClosure.FreeVars)
				childVM.frames[0] = childFrame
				childVM.fp = 1

				// Run the child VM (ignore return value - goroutine result is discarded)
				_, _ = childVM.run()
			})
		}

	case bytecode.OpSend:
		val := vm.pop()
		ch := vm.pop()
		if err := ch.SendContext(vm.ctx, val); err != nil {
			return err
		}

	case bytecode.OpRecv:
		ch := vm.pop()
		val, _, err := ch.RecvContext(vm.ctx)
		if err != nil {
			return err
		}
		vm.push(val)

	case bytecode.OpRecvOk:
		// Receive with comma-ok: returns (value, ok) tuple
		ch := vm.pop()
		val, recvOK, err := ch.RecvContext(vm.ctx)
		if err != nil {
			return err
		}
		// Push as tuple (value, ok)
		tuple := []value.Value{val, value.MakeBool(recvOK)}
		vm.push(value.FromInterface(tuple))

	case bytecode.OpClose:
		ch := vm.pop()
		ch.Close()

	case bytecode.OpSelect:
		// OpSelect performs a select statement using reflect.Select.
		// Operands: [meta_idx:2]
		// Stack (bottom to top): for each state, Chan; if send, also SendVal.
		// Result pushed: tuple (index, recvOk, recv_0, ..., recv_{n-1})
		metaIdx := frame.readUint16()
		meta, ok := vm.program.Constants[metaIdx].(bytecode.SelectMeta)
		if !ok {
			return fmt.Errorf("OpSelect: invalid meta at index %d", metaIdx)
		}

		// Pop channels and send values from stack (they were pushed in order,
		// so we need to pop in reverse).
		type stateData struct {
			ch      value.Value
			sendVal value.Value
			isSend  bool
		}
		states := make([]stateData, meta.NumStates)
		// Pop in reverse order
		for i := meta.NumStates - 1; i >= 0; i-- {
			if meta.Dirs[i] { // send
				states[i].sendVal = vm.pop()
				states[i].ch = vm.pop()
				states[i].isSend = true
			} else { // recv
				states[i].ch = vm.pop()
			}
		}

		// Build reflect.SelectCase slice
		// Add 1 for default case (non-blocking) or context cancellation case (blocking)
		numCases := meta.NumStates + 1
		cases := make([]reflect.SelectCase, numCases)
		for i := 0; i < meta.NumStates; i++ {
			rv, _ := states[i].ch.ReflectValue()
			if states[i].isSend {
				sendRV := states[i].sendVal.ToReflectValue(rv.Type().Elem())
				cases[i] = reflect.SelectCase{
					Dir:  reflect.SelectSend,
					Chan: rv,
					Send: sendRV,
				}
			} else {
				cases[i] = reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: rv,
				}
			}
		}
		if !meta.Blocking {
			cases[meta.NumStates] = reflect.SelectCase{Dir: reflect.SelectDefault}
		} else {
			// Inject context cancellation case for blocking select
			cases[meta.NumStates] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(vm.ctx.Done()),
			}
		}

		// Perform the select
		chosen, recv, recvOK := reflect.Select(cases)

		// Check if context was cancelled (chosen == meta.NumStates in blocking mode)
		if meta.Blocking && chosen == meta.NumStates {
			return vm.ctx.Err()
		}

		// Adjust chosen index: if default was selected, chosen == meta.NumStates → map to -1
		if !meta.Blocking && chosen == meta.NumStates {
			chosen = -1
		}

		// Build result tuple: (index, recvOk, recv_0, ..., recv_{n-1})
		tupleLen := 2 + meta.NumRecv
		tuple := make([]value.Value, tupleLen)
		tuple[0] = value.MakeInt(int64(chosen))
		tuple[1] = value.MakeBool(recvOK)

		// Fill recv values: for each recv state (in order), if it was the chosen one, set the value
		recvIdx := 0
		for i := 0; i < meta.NumStates; i++ {
			if !meta.Dirs[i] { // recv state
				if i == chosen {
					tuple[2+recvIdx] = value.MakeFromReflect(recv)
				} else {
					tuple[2+recvIdx] = value.MakeNil()
				}
				recvIdx++
			}
		}

		vm.push(value.FromInterface(tuple))

	// Defer/recover
	case bytecode.OpDefer:
		funcIdx := frame.readUint16()
		// Look up the function by index
		var fn *bytecode.CompiledFunction
		if int(funcIdx) < len(vm.program.FuncByIndex) {
			fn = vm.program.FuncByIndex[funcIdx]
		}
		if fn == nil {
			return fmt.Errorf("defer: function index %d not found", funcIdx)
		}
		// Pop arguments from stack (they were pushed before OpDefer)
		args := make([]value.Value, fn.NumParams)
		for i := fn.NumParams - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}
		// Add to defer list (will be executed in LIFO order)
		frame.defers = append(frame.defers, DeferInfo{
			fn:   fn,
			args: args,
		})

	case bytecode.OpDeferIndirect:
		numArgs := int(frame.readUint16())
		// The closure is on the stack after the arguments
		// Stack layout: [... args... closure]
		// Pop closure first (it was pushed last)
		closureVal := vm.pop()
		closure, ok := closureVal.Interface().(*Closure)
		if !ok {
			return fmt.Errorf("defer indirect: expected closure, got %v", closureVal.Kind())
		}
		// Pop arguments
		args := make([]value.Value, numArgs)
		for i := numArgs - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}
		// Add to defer list with closure info
		frame.defers = append(frame.defers, DeferInfo{
			fn:      closure.Fn,
			args:    args,
			closure: closure,
		})

	case bytecode.OpRunDefers:
		// Execute all pending deferred calls synchronously in LIFO order.
		// This is critical for named return values: the code after RunDefers
		// reads the (potentially modified) return values.
		for len(frame.defers) > 0 {
			// Pop the last defer (LIFO)
			d := frame.defers[len(frame.defers)-1]
			frame.defers = frame.defers[:len(frame.defers)-1]

			// Get free variables from closure if present
			var freeVars []*value.Value
			if d.closure != nil {
				freeVars = d.closure.FreeVars
			}

			// Execute the deferred function synchronously using a child VM
			// that shares the same globals/context/program. This avoids
			// interference with the parent frame stack.
			childVM := &VM{
				program:      vm.program,
				stack:        make([]value.Value, 256),
				sp:           0,
				frames:       make([]*Frame, 64),
				fp:           0,
				globals:      vm.globals,
				globalsPtr:   vm.globalsPtr,
				ctx:          vm.ctx,
				extCallCache: vm.extCallCache,
			}
			deferFrame := newFrame(d.fn, 0, d.args, freeVars)
			childVM.frames[0] = deferFrame
			childVM.fp = 1
			_, _ = childVM.run()
		}

	case bytecode.OpRecover:
		// Recover from panic
		if vm.panicking {
			vm.push(vm.panicVal)
			vm.panicking = false
			vm.panicVal = value.MakeNil()
		} else {
			vm.push(value.MakeNil())
		}

	// Range operations
	case bytecode.OpRange:
		// Create an iterator for the collection
		collection := vm.pop()
		vm.push(value.FromInterface(&iterator{collection: collection, index: 0}))

	case bytecode.OpRangeNext:
		// Advance iterator and push a tuple (ok, key, value)
		iterVal := vm.pop()
		iter, ok := iterVal.Interface().(*iterator)
		if !ok {
			// Return tuple (false, nil, nil)
			tuple := []value.Value{value.MakeBool(false), value.MakeNil(), value.MakeNil()}
			vm.push(value.FromInterface(tuple))
			return nil
		}
		key, val, iterOk := iter.next()
		// SSA Next returns (ok, key, value) as a tuple
		tuple := []value.Value{value.MakeBool(iterOk), key, val}
		vm.push(value.FromInterface(tuple))

	// Builtins
	case bytecode.OpLen:
		obj := vm.pop()
		switch obj.Kind() {
		case value.KindString:
			vm.push(value.MakeInt(int64(len(obj.String()))))
		case value.KindSlice:
			vm.push(value.MakeInt(int64(obj.Len())))
		case value.KindArray, value.KindMap, value.KindChan:
			vm.push(value.MakeInt(int64(obj.Len())))
		case value.KindInterface, value.KindReflect:
			// Handle both interface values and reflect-wrapped values
			if rv, ok := obj.ReflectValue(); ok {
				kind := rv.Kind()
				if kind == reflect.Interface {
					// Unwrap interface to get underlying value
					if !rv.IsNil() {
						rv = rv.Elem()
						kind = rv.Kind()
					}
				}
				switch kind {
				case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
					vm.push(value.MakeInt(int64(rv.Len())))
				default:
					vm.push(value.MakeInt(0))
				}
			} else {
				vm.push(value.MakeInt(0))
			}
		default:
			vm.push(value.MakeInt(0))
		}

	case bytecode.OpCap:
		obj := vm.pop()
		switch obj.Kind() {
		case value.KindSlice, value.KindArray, value.KindChan:
			vm.push(value.MakeInt(int64(obj.Cap())))
		case value.KindReflect:
			if rv, ok := obj.ReflectValue(); ok {
				vm.push(value.MakeInt(int64(rv.Cap())))
			} else {
				vm.push(value.MakeInt(0))
			}
		default:
			vm.push(value.MakeInt(0))
		}

	case bytecode.OpAppend:
		elem := vm.pop()
		slice := vm.pop()
		// Native int slice fast path
		if s, ok := slice.IntSlice(); ok {
			if es, ok2 := elem.IntSlice(); ok2 {
				vm.push(value.MakeIntSlice(append(s, es...)))
			} else {
				vm.push(value.MakeIntSlice(append(s, elem.RawInt())))
			}
			break
		}
		// Nil slice with native int element: create []int64
		if slice.IsNil() || slice.Kind() == value.KindInvalid {
			if elem.Kind() == value.KindInt {
				// Single int append to nil → create native []int64
				intReflectSlice := reflect.MakeSlice(reflect.TypeOf([]int{}), 0, 1)
				intReflectSlice = reflect.Append(intReflectSlice, reflect.ValueOf(int(elem.RawInt())))
				vm.push(value.MakeFromReflect(intReflectSlice))
				break
			}
			if elem.Kind() == value.KindSlice {
				if es, ok2 := elem.IntSlice(); ok2 {
					// Spread-append native []int64 to nil
					intReflectSlice := reflect.MakeSlice(reflect.TypeOf([]int{}), len(es), cap(es))
					for i, v := range es {
						intReflectSlice.Index(i).SetInt(v)
					}
					vm.push(value.MakeFromReflect(intReflectSlice))
					break
				}
			}
		}
		// Append element to slice
		if rv, ok := slice.ReflectValue(); ok {
			sliceElemType := rv.Type().Elem()
			// Handle []value.Value slices (used for function slices)
			if sliceElemType == reflect.TypeOf(value.Value{}) {
				// Append value.Value element(s) to []value.Value
				if elemRV, ok2 := elem.ReflectValue(); ok2 && elemRV.Kind() == reflect.Slice && elemRV.Type().Elem() == sliceElemType {
					newSlice := reflect.AppendSlice(rv, elemRV)
					vm.push(value.MakeFromReflect(newSlice))
				} else {
					newSlice := reflect.Append(rv, reflect.ValueOf(elem))
					vm.push(value.MakeFromReflect(newSlice))
				}
			} else {
				// Check if elem is a native []int64 that needs spread-append
				if elem.Kind() == value.KindSlice {
					if elemIntSlice, ok2 := elem.IntSlice(); ok2 {
						// Element is []int64, spread append each element
						for _, v := range elemIntSlice {
							rv = reflect.Append(rv, reflect.ValueOf(int(v)))
						}
						vm.push(value.MakeFromReflect(rv))
						break
					}
				}
				// Check if SSA packed variadic args into a slice (e.g., append(s, elems...))
				if elemRV, ok2 := elem.ReflectValue(); ok2 && elemRV.Kind() == reflect.Slice {
					// The element is a slice of the same element type — spread it
					newSlice := reflect.AppendSlice(rv, elemRV)
					vm.push(value.MakeFromReflect(newSlice))
				} else {
					newSlice := reflect.Append(rv, elem.ToReflectValue(sliceElemType))
					vm.push(value.MakeFromReflect(newSlice))
				}
			}
		} else if slice.IsNil() || slice.Kind() == value.KindInvalid {
			// Nil slice: create a new slice and append the element.
			// Fast path: elem is a native []int64 (spread-append to create new []int64).
			if es, ok2 := elem.IntSlice(); ok2 {
				vm.push(value.MakeIntSlice(append([]int64(nil), es...)))
				break
			}
			elemRV, ok2 := elem.ReflectValue()
			if ok2 && elemRV.Kind() == reflect.Slice {
				if elemRV.Type().Elem() == reflect.TypeOf(value.Value{}) {
					// Create new []value.Value from the element slice
					newSlice := make([]value.Value, 0)
					newRV := reflect.ValueOf(newSlice)
					result := reflect.AppendSlice(newRV, elemRV)
					vm.push(value.MakeFromReflect(result))
				} else {
					// Create new slice of the element's type and spread-append
					sliceType := reflect.SliceOf(elemRV.Type().Elem())
					newSlice := reflect.MakeSlice(sliceType, 0, 0)
					result := reflect.AppendSlice(newSlice, elemRV)
					vm.push(value.MakeFromReflect(result))
				}
			} else {
				// Single-element append to a nil slice: infer element type from the value.
				elemIface := elem.Interface()
				if elemIface != nil {
					elemRV2 := reflect.ValueOf(elemIface)
					sliceType := reflect.SliceOf(elemRV2.Type())
					newSlice := reflect.MakeSlice(sliceType, 0, 0)
					result := reflect.Append(newSlice, elemRV2)
					vm.push(value.MakeFromReflect(result))
				} else {
					vm.push(slice)
				}
			}
		} else {
			vm.push(slice)
		}

	case bytecode.OpCopy:
		src := vm.pop()
		dst := vm.pop()
		// Native int slice fast path
		if ds, ok := dst.IntSlice(); ok {
			if ss, ok2 := src.IntSlice(); ok2 {
				vm.push(value.MakeInt(int64(copy(ds, ss))))
				break
			}
		}
		// Copy slice
		if dstRV, ok := dst.ReflectValue(); ok {
			if srcRV, ok := src.ReflectValue(); ok {
				n := reflect.Copy(dstRV, srcRV)
				vm.push(value.MakeInt(int64(n)))
			}
		} else {
			vm.push(value.MakeInt(0))
		}

	case bytecode.OpDelete:
		key := vm.pop()
		m := vm.pop()
		m.SetMapIndex(key, value.MakeNil())

	case bytecode.OpPanic:
		msg := vm.pop()
		vm.panicking = true
		vm.panicVal = msg

	case bytecode.OpPrint:
		n := frame.readByte()
		for i := 0; i < int(n); i++ {
			val := vm.pop()
			fmt.Print(val.Interface())
		}

	case bytecode.OpPrintln:
		n := frame.readByte()
		args := make([]any, n)
		for i := int(n) - 1; i >= 0; i-- {
			args[i] = vm.pop().Interface()
		}
		fmt.Println(args...)

	case bytecode.OpNew:
		typeIdx := frame.readUint16()
		// Allocate new pointer to value of the given type
		if int(typeIdx) < len(vm.program.Types) {
			typ := vm.program.Types[typeIdx]
			// For function types, create a pointer to a Value (to store closures)
			switch t := typ.(type) {
			case *types.Signature:
				_ = t // Function signature not needed for allocation
				// Create a pointer to a nil Value
				var nilVal value.Value
				newPtr := reflect.ValueOf(&nilVal)
				vm.push(value.MakeFromReflect(newPtr))
			case *types.Slice:
				// Check if slice element type is function
				if _, isFunc := t.Elem().(*types.Signature); isFunc {
					// Create []value.Value for function slices
					var slice []value.Value
					newPtr := reflect.ValueOf(&slice)
					vm.push(value.MakeFromReflect(newPtr))
				} else if rt := typeToReflect(typ); rt != nil {
					newPtr := reflect.New(rt)
					vm.push(value.MakeFromReflect(newPtr))
				} else {
					vm.push(value.MakeNil())
				}
			case *types.Array:
				// Check if array element type is function
				if _, isFunc := t.Elem().(*types.Signature); isFunc {
					// Create array of value.Value for function arrays
					arrLen := int(t.Len())
					array := make([]value.Value, arrLen)
					newPtr := reflect.ValueOf(&array)
					vm.push(value.MakeFromReflect(newPtr))
				} else if rt := typeToReflect(typ); rt != nil {
					newPtr := reflect.New(rt)
					vm.push(value.MakeFromReflect(newPtr))
				} else {
					vm.push(value.MakeNil())
				}
			default:
				if rt := typeToReflect(typ); rt != nil {
					// Create a new pointer to zero value of the type
					newPtr := reflect.New(rt)
					vm.push(value.MakeFromReflect(newPtr))
				} else {
					vm.push(value.MakeNil())
				}
			}
		} else {
			vm.push(value.MakeNil())
		}

	case bytecode.OpMake:
		_ = frame.readUint16() // typeIdx
		_ = frame.readUint16() // sizeIdx
		// Make operation handled by specific OpMakeSlice/Map/Chan

	// External call
	case bytecode.OpCallExternal:
		funcIdx := frame.readUint16()
		numArgs := frame.readByte()
		if err := vm.callExternal(int(funcIdx), int(numArgs)); err != nil {
			return err
		}

	// Indirect call (closures, function values, method expressions)
	case bytecode.OpCallIndirect:
		numArgs := frame.readByte()
		// Pop arguments using stack-allocated buffer
		var argsBuf [8]value.Value
		var args []value.Value
		if int(numArgs) <= len(argsBuf) {
			args = argsBuf[:numArgs]
		} else {
			args = make([]value.Value, numArgs)
		}
		for i := int(numArgs) - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}
		// Pop the callee
		callee := vm.pop()
		switch fn := callee.RawObj().(type) {
		case *Closure:
			// Call closure: create new frame with free vars
			vm.callFunction(fn.Fn, args, fn.FreeVars)
		case *bytecode.CompiledFunction:
			// Call compiled function
			vm.callFunction(fn, args, nil)
		default:
			// Not a known callable — push nil
			vm.push(value.MakeNil())
		}

	case bytecode.OpPack:
		count := frame.readUint16()
		// Pop 'count' values from stack and pack into a slice
		values := make([]value.Value, count)
		for i := int(count) - 1; i >= 0; i-- {
			values[i] = vm.pop()
		}
		vm.push(value.FromInterface(values))

	case bytecode.OpUnpack:
		// Pop a slice and push each element onto the stack
		slice := vm.pop()
		// Fast path: native []value.Value (produced by DirectCall multi-return wrappers)
		if vals, ok := slice.ValueSlice(); ok {
			for _, v := range vals {
				vm.push(v)
			}
			break
		}
		if slice.Kind() == value.KindSlice || slice.Kind() == value.KindReflect {
			if rv, ok := slice.ReflectValue(); ok {
				for i := 0; i < rv.Len(); i++ {
					vm.push(value.MakeFromReflect(rv.Index(i)))
				}
			}
		}

	case bytecode.OpHalt:
		return fmt.Errorf("halt")

	default:
		return fmt.Errorf("unknown opcode: %v", op)
	}

	return nil
}
