// compile_value.go compiles SSA values: constants, conversions, slicing, type asserts.
package compiler

import (
	"go/constant"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/model/bytecode"
)

// compileValue compiles an SSA value to push it onto the stack.
func (c *compiler) compileValue(v ssa.Value) {
	switch val := v.(type) {
	case *ssa.Const:
		c.compileConst(val)
	case *ssa.Function:
		if fnIdx, ok := c.funcIndex[val]; ok {
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(bytecode.OpClosure),
				byte(fnIdx>>8), byte(fnIdx),
				byte(0))
		} else {
			c.emit(bytecode.OpNil)
		}
	case *ssa.Phi:
		if slot, ok := c.phiSlots[val]; ok {
			c.emit(bytecode.OpLocal, uint16(slot))
		} else {
			c.emit(bytecode.OpNil)
		}
	case *ssa.FreeVar:
		if idx, ok := c.symbolTable.freeVars[val]; ok {
			c.emit(bytecode.OpFree, uint16(idx))
		} else {
			c.emit(bytecode.OpNil)
		}
	case *ssa.Global:
		// For external package globals, use qualified name (e.g., "time.UTC")
		// to distinguish from main package globals and enable lookup.
		globalName := val.Name()
		isExternal := val.Pkg != nil && val.Pkg.Pkg != nil && val.Pkg.Pkg.Path() != "main"
		if isExternal {
			globalName = val.Pkg.Pkg.Path() + "." + globalName
		}
		globalIdx, ok := c.globals[globalName]
		if !ok {
			globalIdx = len(c.globals)
			c.globals[globalName] = globalIdx

			// For external variables, look up the variable pointer and dereference it
			// to get the actual value. SSA represents the global as having an extra
			// level of indirection, so we need to store the value, not the pointer.
			if isExternal && c.lookup != nil {
				if ptr, found := c.lookup.LookupExternalVar(val.Pkg.Pkg.Path(), val.Name()); found {
					// ptr is a pointer to the external variable (e.g., &time.UTC)
					// We need to dereference it to get the actual value (e.g., time.UTC)
					rv := reflect.ValueOf(ptr)
					if rv.Kind() == reflect.Ptr && !rv.IsNil() {
						c.externalVarValues[globalIdx] = rv.Elem().Interface()
					}
				}
			}
		}
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(bytecode.OpGlobal),
			byte(globalIdx>>8), byte(globalIdx))
	default:
		if idx, ok := c.symbolTable.GetLocal(v); ok {
			c.emit(bytecode.OpLocal, uint16(idx))
		} else {
			if idx, ok := c.symbolTable.freeVars[v]; ok {
				c.emit(bytecode.OpFree, uint16(idx))
			} else {
				c.emit(bytecode.OpNil)
			}
		}
	}
}

// compileConst compiles a constant value.
func (c *compiler) compileConst(cnst *ssa.Const) {
	var v any
	switch t := cnst.Type().(type) {
	case *types.Basic:
		switch t.Kind() { //nolint:exhaustive
		case types.Bool, types.UntypedBool:
			v = cnst.Value != nil && cnst.Value.Kind() == constant.Bool && constant.BoolVal(cnst.Value)
		case types.Int, types.UntypedInt, types.UntypedRune:
			if cnst.Value != nil {
				i, exact := constant.Int64Val(cnst.Value)
				if exact {
					v = int(i)
				} else {
					v = int(0)
				}
			} else {
				v = int(0)
			}
		case types.Int8:
			if cnst.Value != nil {
				i, _ := constant.Int64Val(cnst.Value)
				v = int8(i)
			} else {
				v = int8(0)
			}
		case types.Int16:
			if cnst.Value != nil {
				i, _ := constant.Int64Val(cnst.Value)
				v = int16(i)
			} else {
				v = int16(0)
			}
		case types.Int32:
			if cnst.Value != nil {
				i, _ := constant.Int64Val(cnst.Value)
				v = int32(i)
			} else {
				v = int32(0)
			}
		case types.Int64:
			if cnst.Value != nil {
				i, exact := constant.Int64Val(cnst.Value)
				if exact {
					v = i
				} else {
					v = int64(0)
				}
			} else {
				v = int64(0)
			}
		case types.Uint:
			if cnst.Value != nil {
				u, _ := constant.Uint64Val(cnst.Value)
				v = uint(u)
			} else {
				v = uint(0)
			}
		case types.Uint8:
			if cnst.Value != nil {
				u, _ := constant.Uint64Val(cnst.Value)
				v = uint8(u)
			} else {
				v = uint8(0)
			}
		case types.Uint16:
			if cnst.Value != nil {
				u, _ := constant.Uint64Val(cnst.Value)
				v = uint16(u)
			} else {
				v = uint16(0)
			}
		case types.Uint32:
			if cnst.Value != nil {
				u, _ := constant.Uint64Val(cnst.Value)
				v = uint32(u)
			} else {
				v = uint32(0)
			}
		case types.Uint64:
			if cnst.Value != nil {
				u, _ := constant.Uint64Val(cnst.Value)
				v = u
			} else {
				v = uint64(0)
			}
		case types.Uintptr:
			if cnst.Value != nil {
				u, _ := constant.Uint64Val(cnst.Value)
				v = uint64(u)
			} else {
				v = uint64(0)
			}
		case types.Float32:
			if cnst.Value != nil {
				f, _ := constant.Float64Val(cnst.Value)
				v = float32(f)
			} else {
				v = float32(0)
			}
		case types.Float64, types.UntypedFloat:
			if cnst.Value != nil {
				f, _ := constant.Float64Val(cnst.Value)
				v = f
			} else {
				v = 0.0
			}
		case types.String, types.UntypedString:
			if cnst.Value != nil {
				v = constant.StringVal(cnst.Value)
			} else {
				v = ""
			}
		case types.Complex64:
			if cnst.Value != nil {
				re := constant.Real(cnst.Value)
				im := constant.Imag(cnst.Value)
				reVal, _ := constant.Float64Val(re)
				imVal, _ := constant.Float64Val(im)
				v = complex(float32(reVal), float32(imVal))
			} else {
				v = complex64(0)
			}
		case types.Complex128:
			if cnst.Value != nil {
				re := constant.Real(cnst.Value)
				im := constant.Imag(cnst.Value)
				reVal, _ := constant.Float64Val(re)
				imVal, _ := constant.Float64Val(im)
				v = complex(reVal, imVal)
			} else {
				v = complex128(0)
			}
		case types.UntypedComplex:
			if cnst.Value != nil {
				re := constant.Real(cnst.Value)
				im := constant.Imag(cnst.Value)
				reVal, _ := constant.Float64Val(re)
				imVal, _ := constant.Float64Val(im)
				v = complex(reVal, imVal)
			} else {
				v = complex128(0)
			}
		default:
			v = nil
		}
	case *types.Named:
		// Handle named types by extracting their underlying basic type
		if cnst.Value != nil {
			switch underlying := t.Underlying().(type) {
			case *types.Basic:
				switch underlying.Kind() { //nolint:exhaustive
				case types.Int, types.UntypedInt, types.UntypedRune:
					i, exact := constant.Int64Val(cnst.Value)
					if exact {
						v = int(i)
					} else {
						v = int(0)
					}
				case types.Int8:
					i, _ := constant.Int64Val(cnst.Value)
					v = int8(i)
				case types.Int16:
					i, _ := constant.Int64Val(cnst.Value)
					v = int16(i)
				case types.Int32:
					i, _ := constant.Int64Val(cnst.Value)
					v = int32(i)
				case types.Int64:
					i, _ := constant.Int64Val(cnst.Value)
					v = i
				case types.Uint:
					u, _ := constant.Uint64Val(cnst.Value)
					v = uint(u)
				case types.Uint8:
					u, _ := constant.Uint64Val(cnst.Value)
					v = uint8(u)
				case types.Uint16:
					u, _ := constant.Uint64Val(cnst.Value)
					v = uint16(u)
				case types.Uint32:
					u, _ := constant.Uint64Val(cnst.Value)
					v = uint32(u)
				case types.Uint64:
					u, _ := constant.Uint64Val(cnst.Value)
					v = u
				case types.Float32:
					f, _ := constant.Float64Val(cnst.Value)
					v = float32(f)
				case types.Float64, types.UntypedFloat:
					f, _ := constant.Float64Val(cnst.Value)
					v = f
				case types.String, types.UntypedString:
					v = constant.StringVal(cnst.Value)
				case types.Bool, types.UntypedBool:
					v = cnst.Value != nil && cnst.Value.Kind() == constant.Bool && constant.BoolVal(cnst.Value)
				}
			}
		} else {
			// cnst.Value == nil: emit zero value for named types
			if rt := constTypeToReflect(t); rt != nil {
				v = reflect.Zero(rt)
			}
		}
	case *types.Alias:
		// Handle type aliases (e.g., type MyInt = int) - they are identical to the aliased type
		// Use the underlying type for value handling
		if cnst.Value != nil {
			switch underlying := t.Underlying().(type) {
			case *types.Basic:
				switch underlying.Kind() { //nolint:exhaustive
				case types.Int, types.UntypedInt, types.UntypedRune:
					i, exact := constant.Int64Val(cnst.Value)
					if exact {
						v = int(i)
					} else {
						v = int(0)
					}
				case types.Int8:
					i, _ := constant.Int64Val(cnst.Value)
					v = int8(i)
				case types.Int16:
					i, _ := constant.Int64Val(cnst.Value)
					v = int16(i)
				case types.Int32:
					i, _ := constant.Int64Val(cnst.Value)
					v = int32(i)
				case types.Int64:
					i, exact := constant.Int64Val(cnst.Value)
					if exact {
						v = i
					} else {
						v = int64(0)
					}
				case types.Uint:
					u, _ := constant.Uint64Val(cnst.Value)
					v = uint(u)
				case types.Uint8:
					u, _ := constant.Uint64Val(cnst.Value)
					v = uint8(u)
				case types.Uint16:
					u, _ := constant.Uint64Val(cnst.Value)
					v = uint16(u)
				case types.Uint32:
					u, _ := constant.Uint64Val(cnst.Value)
					v = uint32(u)
				case types.Uint64:
					u, _ := constant.Uint64Val(cnst.Value)
					v = u
				case types.Float32:
					f, _ := constant.Float64Val(cnst.Value)
					v = float32(f)
				case types.Float64, types.UntypedFloat:
					f, _ := constant.Float64Val(cnst.Value)
					v = f
				case types.String, types.UntypedString:
					v = constant.StringVal(cnst.Value)
				case types.Bool, types.UntypedBool:
					v = cnst.Value != nil && cnst.Value.Kind() == constant.Bool && constant.BoolVal(cnst.Value)
				}
			}
		} else {
			// cnst.Value == nil: emit zero value for alias types
			if rt := constTypeToReflect(t); rt != nil {
				v = reflect.Zero(rt)
			}
		}
	default:
		// For nil constants of reference types (map, slice, chan, func, pointer),
		// emit a typed nil via reflect.Zero so the VM preserves type information.
		// This is critical for nil map access (returns zero value of element type)
		// and nil-typed interface returns.
		// Also handles struct types (including empty struct) for zero values.
		if cnst.Value == nil {
			if rt := constTypeToReflect(cnst.Type()); rt != nil {
				v = reflect.Zero(rt)
			}
		}
	case *types.Struct:
		// Handle struct zero values (including empty struct literal {})
		if cnst.Value == nil {
			if rt := constTypeToReflect(t); rt != nil {
				v = reflect.Zero(rt)
			}
		}
	}

	idx := c.addConstant(v)
	c.emit(bytecode.OpConst, idx)
}

// constTypeToReflect converts a go/types.Type to reflect.Type for nil constant emission.
// Only handles reference types (map, slice, chan, pointer, func) since those are the
// types that can have meaningful nil values with type information.
// Named empty structs now share reflect.TypeOf(struct{}{}) with the VM, so they
// can be emitted directly.
// Returns nil for types that cannot be converted (the caller falls back to untyped nil).
func constTypeToReflect(t types.Type) reflect.Type {
	// Named empty structs share struct{}{} with the VM (no phantom field),
	// so we can emit struct{} directly.
	if named, ok := t.(*types.Named); ok {
		if structType, ok := named.Underlying().(*types.Struct); ok && structType.NumFields() == 0 {
			return reflect.TypeFor[struct{}]()
		}
	}
	// Handle type aliases - they are identical to the aliased type
	if alias, ok := t.(*types.Alias); ok {
		if structType, ok := alias.Underlying().(*types.Struct); ok && structType.NumFields() == 0 {
			return reflect.TypeFor[struct{}]()
		}
	}

	switch typ := t.Underlying().(type) {
	case *types.Map:
		keyRT := constTypeToReflect(typ.Key())
		elemRT := constTypeToReflect(typ.Elem())
		if keyRT != nil && elemRT != nil {
			return reflect.MapOf(keyRT, elemRT)
		}
	case *types.Slice:
		elemRT := constTypeToReflect(typ.Elem())
		if elemRT != nil {
			return reflect.SliceOf(elemRT)
		}
	case *types.Pointer:
		elemRT := constTypeToReflect(typ.Elem())
		if elemRT != nil {
			return reflect.PointerTo(elemRT)
		}
	case *types.Chan:
		elemRT := constTypeToReflect(typ.Elem())
		if elemRT != nil {
			dir := reflect.BothDir
			switch typ.Dir() {
			case types.SendOnly:
				dir = reflect.SendDir
			case types.RecvOnly:
				dir = reflect.RecvDir
			}
			return reflect.ChanOf(dir, elemRT)
		}
	case *types.Basic:
		switch typ.Kind() {
		case types.Bool:
			return reflect.TypeFor[bool]()
		case types.Int:
			return reflect.TypeFor[int]()
		case types.Int8:
			return reflect.TypeFor[int8]()
		case types.Int16:
			return reflect.TypeFor[int16]()
		case types.Int32:
			return reflect.TypeFor[int32]()
		case types.Int64:
			return reflect.TypeFor[int64]()
		case types.Uint:
			return reflect.TypeFor[uint]()
		case types.Uint8:
			return reflect.TypeFor[uint8]()
		case types.Uint16:
			return reflect.TypeFor[uint16]()
		case types.Uint32:
			return reflect.TypeFor[uint32]()
		case types.Uint64:
			return reflect.TypeFor[uint64]()
		case types.Float32:
			return reflect.TypeFor[float32]()
		case types.Float64:
			return reflect.TypeFor[float64]()
		case types.String:
			return reflect.TypeFor[string]()
		}
	case *types.Interface:
		// Interface with no methods = any
		if typ.NumMethods() == 0 {
			return reflect.TypeFor[any]()
		}
	case *types.Signature:
		// Build reflect.FuncOf for function types (e.g., func() int, func(string) bool)
		params := make([]reflect.Type, typ.Params().Len())
		for i := 0; i < typ.Params().Len(); i++ {
			pt := constTypeToReflect(typ.Params().At(i).Type())
			if pt == nil {
				return nil
			}
			params[i] = pt
		}
		results := make([]reflect.Type, typ.Results().Len())
		for i := 0; i < typ.Results().Len(); i++ {
			rt := constTypeToReflect(typ.Results().At(i).Type())
			if rt == nil {
				return nil
			}
			results[i] = rt
		}
		return reflect.FuncOf(params, results, typ.Variadic())
	case *types.Struct:
		// Handle empty struct (struct{}) — the common case for map[K]struct{}.
		// Named empty structs also use struct{}{} in the VM (no phantom field).
		if typ.NumFields() == 0 {
			// Return the basic struct{} type for empty structs
			return reflect.TypeFor[struct{}]()
		}
	}
	return nil
}

// compileField compiles a Field instruction.
func (c *compiler) compileField(i *ssa.Field) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpField, uint16(i.Field))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileFieldAddr compiles a FieldAddr instruction.
func (c *compiler) compileFieldAddr(i *ssa.FieldAddr) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpFieldAddr, uint16(i.Field))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileIndex compiles an Index instruction.
func (c *compiler) compileIndex(i *ssa.Index) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.compileValue(i.Index)
	c.emit(bytecode.OpIndex)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileIndexAddr compiles an IndexAddr instruction.
func (c *compiler) compileIndexAddr(i *ssa.IndexAddr) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.compileValue(i.Index)
	c.emit(bytecode.OpIndexAddr)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileLookup compiles a Lookup instruction.
func (c *compiler) compileLookup(i *ssa.Lookup) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.compileValue(i.Index)

	if i.CommaOk {
		c.emit(bytecode.OpIndexOk)
	} else {
		c.emit(bytecode.OpIndex)
	}
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileStore compiles a Store instruction.
func (c *compiler) compileStore(i *ssa.Store) {
	c.compileValue(i.Addr)
	c.compileValue(i.Val)
	c.emit(bytecode.OpSetDeref)
}

// compileMakeSlice compiles a MakeSlice instruction.
func (c *compiler) compileMakeSlice(i *ssa.MakeSlice) {
	typeIdx := c.addType(i.Type())
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdxConst := c.addConstant(int64(typeIdx))
	c.emit(bytecode.OpConst, typeIdxConst)
	c.compileValue(i.Len)
	c.compileValue(i.Cap)
	c.emit(bytecode.OpMakeSlice)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeMap compiles a MakeMap instruction.
func (c *compiler) compileMakeMap(i *ssa.MakeMap) {
	typeIdx := c.addType(i.Type())
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdxConst := c.addConstant(int64(typeIdx))
	c.emit(bytecode.OpConst, typeIdxConst)

	if i.Reserve != nil {
		c.compileValue(i.Reserve)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0))))
	}

	c.emit(bytecode.OpMakeMap)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeChan compiles a MakeChan instruction.
func (c *compiler) compileMakeChan(i *ssa.MakeChan) {
	typeIdx := c.addType(i.Type())
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdxConst := c.addConstant(int64(typeIdx))
	c.emit(bytecode.OpConst, typeIdxConst)

	if i.Size != nil {
		c.compileValue(i.Size)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0))))
	}

	c.emit(bytecode.OpMakeChan)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeInterface compiles a MakeInterface instruction.
func (c *compiler) compileMakeInterface(i *ssa.MakeInterface) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeClosure compiles a MakeClosure instruction.
func (c *compiler) compileMakeClosure(i *ssa.MakeClosure) {
	fnIdx := c.funcIndex[i.Fn.(*ssa.Function)]
	resultIdx := c.symbolTable.AllocLocal(i)

	for _, binding := range i.Bindings {
		if alloc, ok := binding.(*ssa.Alloc); ok {
			if slotIdx, ok := c.symbolTable.GetLocal(alloc); ok {
				// For *ssa.Alloc (which is already a pointer type), we need to
				// get the pointer value itself (OpLocal), not the address of the slot (OpAddr).
				// Each Alloc creates a new pointer in heap/stack, and closures should
				// capture this pointer value, not the slot address which gets overwritten
				// in loop iterations.
				c.emit(bytecode.OpLocal, uint16(slotIdx))
				continue
			}
		}
		c.compileValue(binding)
	}

	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpClosure),
		byte(fnIdx>>8), byte(fnIdx),
		byte(len(i.Bindings)))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMapUpdate compiles a MapUpdate instruction.
func (c *compiler) compileMapUpdate(i *ssa.MapUpdate) {
	c.compileValue(i.Map)
	c.compileValue(i.Key)
	c.compileValue(i.Value)
	c.emit(bytecode.OpSetIndex)
}

// compileRange compiles a Range instruction.
func (c *compiler) compileRange(i *ssa.Range) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpRange)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileNext compiles a Next instruction.
func (c *compiler) compileNext(i *ssa.Next) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.Iter)
	c.emit(bytecode.OpRangeNext)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSelect compiles a Select instruction.
func (c *compiler) compileSelect(i *ssa.Select) {
	numRecv := 0
	for _, st := range i.States {
		if st.Dir == types.RecvOnly {
			numRecv++
		}
	}

	dirs := make([]bool, len(i.States))
	for idx, st := range i.States {
		dirs[idx] = (st.Dir == types.SendOnly)
	}

	meta := bytecode.SelectMeta{
		NumStates: len(i.States),
		Blocking:  i.Blocking,
		Dirs:      dirs,
		NumRecv:   numRecv,
	}

	for _, st := range i.States {
		c.compileValue(st.Chan)
		if st.Dir == types.SendOnly {
			c.compileValue(st.Send)
		}
	}

	metaIdx := c.addConstant(meta)
	c.emit(bytecode.OpSelect, metaIdx)

	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSlice compiles a Slice instruction.
func (c *compiler) compileSlice(i *ssa.Slice) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)

	if i.Low != nil {
		c.compileValue(i.Low)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0))))
	}

	if i.High != nil {
		c.compileValue(i.High)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0xFFFF))))
	}

	if i.Max != nil {
		c.compileValue(i.Max)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0xFFFF))))
	}

	c.emit(bytecode.OpSlice)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileTypeAssert compiles a TypeAssert instruction.
func (c *compiler) compileTypeAssert(i *ssa.TypeAssert) {
	typeIdx := c.addType(i.AssertedType)
	c.compileValue(i.X)
	c.emit(bytecode.OpAssert, uint16(typeIdx))

	if !i.CommaOk {
		// Non-comma-ok assertion: extract just the value from the [result, ok] tuple.
		// SSA's `typeassert t.(T)` (without comma-ok) returns a single value and
		// panics on failure. OpAssert always produces a tuple, so we extract #0.
		c.emit(bytecode.OpConst, uint16(c.addConstant(0)))
		c.emit(bytecode.OpIndex)
	}

	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileChangeInterface compiles a ChangeInterface instruction.
func (c *compiler) compileChangeInterface(i *ssa.ChangeInterface) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileChangeType compiles a ChangeType instruction.
// ChangeType converts between types with identical underlying types (e.g., []int -> sort.IntSlice).
// We emit OpChangeType which carries both the target type and the source local index,
// so the VM can update the source variable to share the same backing array after conversion.
func (c *compiler) compileChangeType(i *ssa.ChangeType) {
	resultIdx := c.symbolTable.AllocLocal(i)
	typeIdx := c.addType(i.Type())

	// Try to find the source local index. If the source is a local variable,
	// we pass its index so the VM can update it for slice aliasing.
	srcLocalIdx := uint16(0xFFFF) // sentinel: no source local
	if srcIdx, ok := c.symbolTable.GetLocal(i.X); ok {
		srcLocalIdx = uint16(srcIdx)
	}

	c.compileValue(i.X)
	// Emit OpChangeType with 4 bytes of operands: type_idx(2) + src_local(2)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpChangeType),
		byte(typeIdx>>8), byte(typeIdx),
		byte(srcLocalIdx>>8), byte(srcLocalIdx))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileConvert compiles a Convert instruction.
func (c *compiler) compileConvert(i *ssa.Convert) {
	resultIdx := c.symbolTable.AllocLocal(i)
	typeIdx := c.addType(i.Type())
	c.compileValue(i.X)
	c.emit(bytecode.OpConvert, uint16(typeIdx))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileExtract compiles an Extract instruction.
func (c *compiler) compileExtract(i *ssa.Extract) {
	c.compileValue(i.Tuple)
	c.emit(bytecode.OpConst, uint16(c.addConstant(i.Index)))
	c.emit(bytecode.OpIndex)
	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSend compiles a Send instruction.
func (c *compiler) compileSend(i *ssa.Send) {
	c.compileValue(i.Chan)
	c.compileValue(i.X)
	c.emit(bytecode.OpSend)
}

// compileDefer compiles a Defer instruction.
func (c *compiler) compileDefer(i *ssa.Defer) {
	switch val := i.Call.Value.(type) {
	case *ssa.Function:
		// If the function has free variables, we need to create a closure
		if len(val.FreeVars) > 0 {
			// First create the closure (OpDeferIndirect expects: closure, args... on stack)
			// Push the free variable bindings
			for _, fv := range val.FreeVars {
				c.compileValue(fv)
			}
			// Create the closure (leaves it on stack, no SETLOCAL)
			fnIdx := c.funcIndex[val]
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(bytecode.OpClosure),
				byte(fnIdx>>8), byte(fnIdx),
				byte(len(val.FreeVars)))
			// Push arguments AFTER closure
			for _, arg := range i.Call.Args {
				c.compileValue(arg)
			}
			numArgs := len(i.Call.Args)
			c.emit(bytecode.OpDeferIndirect, uint16(numArgs))
			return
		}
		// No free variables - push args then use OpDefer directly
		for _, arg := range i.Call.Args {
			c.compileValue(arg)
		}
		fnIdx := c.funcIndex[val]
		c.emit(bytecode.OpDefer, uint16(fnIdx))

	case *ssa.MakeClosure:
		// Check if this MakeClosure was already compiled (has a local slot)
		if idx, ok := c.symbolTable.GetLocal(val); ok {
			// Already compiled - load the closure from local FIRST
			c.emit(bytecode.OpLocal, uint16(idx))
			// Then push arguments
			for _, arg := range i.Call.Args {
				c.compileValue(arg)
			}
			numArgs := len(i.Call.Args)
			c.emit(bytecode.OpDeferIndirect, uint16(numArgs))
			return
		}
		// Not yet compiled - create the closure now FIRST
		// Compile bindings - need to handle FreeVar specially
		for _, binding := range val.Bindings {
			// Check if binding is a FreeVar (captured from enclosing function)
			if fv, ok := binding.(*ssa.FreeVar); ok {
				if idx, ok := c.symbolTable.freeVars[fv]; ok {
					c.emit(bytecode.OpFree, uint16(idx))
					continue
				}
			}
			// Handle Alloc (pointer variable)
			if alloc, ok := binding.(*ssa.Alloc); ok {
				if slotIdx, ok := c.symbolTable.GetLocal(alloc); ok {
					c.emit(bytecode.OpLocal, uint16(slotIdx))
					continue
				}
			}
			c.compileValue(binding)
		}
		// Create closure (on stack, no SETLOCAL)
		fnIdx := c.funcIndex[val.Fn.(*ssa.Function)]
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(bytecode.OpClosure),
			byte(fnIdx>>8), byte(fnIdx),
			byte(len(val.Bindings)))
		// Push arguments AFTER closure
		for _, arg := range i.Call.Args {
			c.compileValue(arg)
		}
		numArgs := len(i.Call.Args)
		c.emit(bytecode.OpDeferIndirect, uint16(numArgs))

	default:
		// Other cases - first compile the callable, then push args
		c.compileValue(i.Call.Value)
		for _, arg := range i.Call.Args {
			c.compileValue(arg)
		}
		numArgs := len(i.Call.Args)
		c.emit(bytecode.OpDeferIndirect, uint16(numArgs))
	}
}

// compileGo compiles a Go instruction.
func (c *compiler) compileGo(i *ssa.Go) {
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		funcIdx := c.funcIndex[fn]
		numArgs := len(i.Call.Args)
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(bytecode.OpGoCall),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		return
	}

	c.compileValue(i.Call.Value)

	numArgs := len(i.Call.Args)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpGoCallIndirect),
		byte(numArgs))
}
