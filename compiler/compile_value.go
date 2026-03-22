package compiler

import (
	"go/constant"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/bytecode"
)

// compileValue compiles an SSA value to push it onto the stack.
func (ctx *funcContext) compileValue(v ssa.Value) {
	switch val := v.(type) {
	case *ssa.Const:
		ctx.compileConst(val)
	case *ssa.Function:
		if fnIdx, ok := ctx.c.funcIndex[val]; ok {
			ctx.cf.Instructions = append(ctx.cf.Instructions,
				byte(bytecode.OpClosure),
				byte(fnIdx>>8), byte(fnIdx),
				byte(0))
		} else {
			ctx.emit(bytecode.OpNil)
		}
	case *ssa.Phi:
		if slot, ok := ctx.phiSlots[val]; ok {
			ctx.emit(bytecode.OpLocal, uint16(slot))
		} else {
			ctx.emit(bytecode.OpNil)
		}
	case *ssa.FreeVar:
		if idx, ok := ctx.symbolTable.freeVars[val]; ok {
			ctx.emit(bytecode.OpFree, uint16(idx))
		} else {
			ctx.emit(bytecode.OpNil)
		}
	case *ssa.Global:
		// For external package globals, use qualified name (e.g., "time.UTC")
		// to distinguish from main package globals and enable lookup.
		globalName := val.Name()
		isExternal := val.Pkg != nil && val.Pkg.Pkg != nil && val.Pkg.Pkg.Path() != "main"
		if isExternal {
			globalName = val.Pkg.Pkg.Path() + "." + globalName
		}
		globalIdx, ok := ctx.c.globals[globalName]
		if !ok {
			globalIdx = len(ctx.c.globals)
			ctx.c.globals[globalName] = globalIdx

			// For external variables, look up the variable pointer and dereference it
			// to get the actual value. SSA represents the global as having an extra
			// level of indirection, so we need to store the value, not the pointer.
			if isExternal && ctx.c.registry != nil {
				if ptr, found := ctx.c.registry.LookupExternalVar(val.Pkg.Pkg.Path(), val.Name()); found {
					// ptr is a pointer to the external variable (e.g., &time.UTC)
					// We need to dereference it to get the actual value (e.g., time.UTC)
					rv := reflect.ValueOf(ptr)
					if rv.Kind() == reflect.Ptr && !rv.IsNil() {
						ctx.c.externalVarValues[globalIdx] = rv.Elem().Interface()
					}
				}
			}
		}
		ctx.cf.Instructions = append(ctx.cf.Instructions,
			byte(bytecode.OpGlobal),
			byte(globalIdx>>8), byte(globalIdx))
	default:
		if idx, ok := ctx.symbolTable.GetLocal(v); ok {
			ctx.emit(bytecode.OpLocal, uint16(idx))
		} else {
			if idx, ok := ctx.symbolTable.freeVars[v]; ok {
				ctx.emit(bytecode.OpFree, uint16(idx))
			} else {
				ctx.emit(bytecode.OpNil)
			}
		}
	}
}

// compileConst compiles a constant value.
func (ctx *funcContext) compileConst(cnst *ssa.Const) {
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
		}
	default:
		v = nil
	}

	idx := ctx.c.addConstant(v)
	ctx.emit(bytecode.OpConst, idx)
}

// compileField compiles a Field instruction.
func (ctx *funcContext) compileField(i *ssa.Field) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.compileValue(i.X)
	ctx.emit(bytecode.OpField, uint16(i.Field))
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileFieldAddr compiles a FieldAddr instruction.
func (ctx *funcContext) compileFieldAddr(i *ssa.FieldAddr) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.compileValue(i.X)
	ctx.emit(bytecode.OpFieldAddr, uint16(i.Field))
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileIndex compiles an Index instruction.
func (ctx *funcContext) compileIndex(i *ssa.Index) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.compileValue(i.X)
	ctx.compileValue(i.Index)
	ctx.emit(bytecode.OpIndex)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileIndexAddr compiles an IndexAddr instruction.
func (ctx *funcContext) compileIndexAddr(i *ssa.IndexAddr) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.compileValue(i.X)
	ctx.compileValue(i.Index)
	ctx.emit(bytecode.OpIndexAddr)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileLookup compiles a Lookup instruction.
func (ctx *funcContext) compileLookup(i *ssa.Lookup) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.compileValue(i.X)
	ctx.compileValue(i.Index)

	if i.CommaOk {
		ctx.emit(bytecode.OpIndexOk)
	} else {
		ctx.emit(bytecode.OpIndex)
	}
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileStore compiles a Store instruction.
func (ctx *funcContext) compileStore(i *ssa.Store) {
	ctx.compileValue(i.Addr)
	ctx.compileValue(i.Val)
	ctx.emit(bytecode.OpSetDeref)
}

// compileMakeSlice compiles a MakeSlice instruction.
func (ctx *funcContext) compileMakeSlice(i *ssa.MakeSlice) {
	typeIdx := ctx.c.addType(i.Type())
	resultIdx := ctx.symbolTable.AllocLocal(i)

	typeIdxConst := ctx.c.addConstant(int64(typeIdx))
	ctx.emit(bytecode.OpConst, typeIdxConst)
	ctx.compileValue(i.Len)
	ctx.compileValue(i.Cap)
	ctx.emit(bytecode.OpMakeSlice)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeMap compiles a MakeMap instruction.
func (ctx *funcContext) compileMakeMap(i *ssa.MakeMap) {
	typeIdx := ctx.c.addType(i.Type())
	resultIdx := ctx.symbolTable.AllocLocal(i)

	typeIdxConst := ctx.c.addConstant(int64(typeIdx))
	ctx.emit(bytecode.OpConst, typeIdxConst)

	if i.Reserve != nil {
		ctx.compileValue(i.Reserve)
	} else {
		ctx.emit(bytecode.OpConst, uint16(ctx.c.addConstant(int64(0))))
	}

	ctx.emit(bytecode.OpMakeMap)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeChan compiles a MakeChan instruction.
func (ctx *funcContext) compileMakeChan(i *ssa.MakeChan) {
	typeIdx := ctx.c.addType(i.Type())
	resultIdx := ctx.symbolTable.AllocLocal(i)

	typeIdxConst := ctx.c.addConstant(int64(typeIdx))
	ctx.emit(bytecode.OpConst, typeIdxConst)

	if i.Size != nil {
		ctx.compileValue(i.Size)
	} else {
		ctx.emit(bytecode.OpConst, uint16(ctx.c.addConstant(int64(0))))
	}

	ctx.emit(bytecode.OpMakeChan)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeInterface compiles a MakeInterface instruction.
func (ctx *funcContext) compileMakeInterface(i *ssa.MakeInterface) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.compileValue(i.X)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeClosure compiles a MakeClosure instruction.
func (ctx *funcContext) compileMakeClosure(i *ssa.MakeClosure) {
	fnIdx := ctx.c.funcIndex[i.Fn.(*ssa.Function)]
	resultIdx := ctx.symbolTable.AllocLocal(i)

	for _, binding := range i.Bindings {
		if alloc, ok := binding.(*ssa.Alloc); ok {
			if slotIdx, ok := ctx.symbolTable.GetLocal(alloc); ok {
				// For *ssa.Alloc (which is already a pointer type), we need to
				// get the pointer value itself (OpLocal), not the address of the slot (OpAddr).
				// Each Alloc creates a new pointer in heap/stack, and closures should
				// capture this pointer value, not the slot address which gets overwritten
				// in loop iterations.
				ctx.emit(bytecode.OpLocal, uint16(slotIdx))
				continue
			}
		}
		ctx.compileValue(binding)
	}

	ctx.cf.Instructions = append(ctx.cf.Instructions,
		byte(bytecode.OpClosure),
		byte(fnIdx>>8), byte(fnIdx),
		byte(len(i.Bindings)))
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMapUpdate compiles a MapUpdate instruction.
func (ctx *funcContext) compileMapUpdate(i *ssa.MapUpdate) {
	ctx.compileValue(i.Map)
	ctx.compileValue(i.Key)
	ctx.compileValue(i.Value)
	ctx.emit(bytecode.OpSetIndex)
}

// compileRange compiles a Range instruction.
func (ctx *funcContext) compileRange(i *ssa.Range) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.compileValue(i.X)
	ctx.emit(bytecode.OpRange)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileNext compiles a Next instruction.
func (ctx *funcContext) compileNext(i *ssa.Next) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.compileValue(i.Iter)
	ctx.emit(bytecode.OpRangeNext)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSelect compiles a Select instruction.
func (ctx *funcContext) compileSelect(i *ssa.Select) {
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
		ctx.compileValue(st.Chan)
		if st.Dir == types.SendOnly {
			ctx.compileValue(st.Send)
		}
	}

	metaIdx := ctx.c.addConstant(meta)
	ctx.emit(bytecode.OpSelect, metaIdx)

	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSlice compiles a Slice instruction.
func (ctx *funcContext) compileSlice(i *ssa.Slice) {
	resultIdx := ctx.symbolTable.AllocLocal(i)

	ctx.compileValue(i.X)

	if i.Low != nil {
		ctx.compileValue(i.Low)
	} else {
		ctx.emit(bytecode.OpConst, uint16(ctx.c.addConstant(int64(0))))
	}

	if i.High != nil {
		ctx.compileValue(i.High)
	} else {
		ctx.emit(bytecode.OpConst, uint16(ctx.c.addConstant(int64(0xFFFF))))
	}

	if i.Max != nil {
		ctx.compileValue(i.Max)
	} else {
		ctx.emit(bytecode.OpConst, uint16(ctx.c.addConstant(int64(0xFFFF))))
	}

	ctx.emit(bytecode.OpSlice)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileTypeAssert compiles a TypeAssert instruction.
func (ctx *funcContext) compileTypeAssert(i *ssa.TypeAssert) {
	typeIdx := ctx.c.addType(i.AssertedType)
	ctx.compileValue(i.X)
	ctx.emit(bytecode.OpAssert, uint16(typeIdx))

	if !i.CommaOk {
		// Non-comma-ok assertion: extract just the value from the [result, ok] tuple.
		// SSA's `typeassert t.(T)` (without comma-ok) returns a single value and
		// panics on failure. OpAssert always produces a tuple, so we extract #0.
		ctx.emit(bytecode.OpConst, uint16(ctx.c.addConstant(0)))
		ctx.emit(bytecode.OpIndex)
	}

	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileChangeInterface compiles a ChangeInterface instruction.
func (ctx *funcContext) compileChangeInterface(i *ssa.ChangeInterface) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.compileValue(i.X)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileChangeType compiles a ChangeType instruction.
// ChangeType converts between types with identical underlying types (e.g., []int -> sort.IntSlice).
// We emit OpChangeType which carries both the target type and the source local index,
// so the VM can update the source variable to share the same backing array after conversion.
func (ctx *funcContext) compileChangeType(i *ssa.ChangeType) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	typeIdx := ctx.c.addType(i.Type())

	// Try to find the source local index. If the source is a local variable,
	// we pass its index so the VM can update it for slice aliasing.
	srcLocalIdx := uint16(0xFFFF) // sentinel: no source local
	if srcIdx, ok := ctx.symbolTable.GetLocal(i.X); ok {
		srcLocalIdx = uint16(srcIdx)
	}

	ctx.compileValue(i.X)
	// Emit OpChangeType with 4 bytes of operands: type_idx(2) + src_local(2)
	ctx.cf.Instructions = append(ctx.cf.Instructions,
		byte(bytecode.OpChangeType),
		byte(typeIdx>>8), byte(typeIdx),
		byte(srcLocalIdx>>8), byte(srcLocalIdx))
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileConvert compiles a Convert instruction.
func (ctx *funcContext) compileConvert(i *ssa.Convert) {
	resultIdx := ctx.symbolTable.AllocLocal(i)
	typeIdx := ctx.c.addType(i.Type())
	ctx.compileValue(i.X)
	ctx.emit(bytecode.OpConvert, uint16(typeIdx))
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileExtract compiles an Extract instruction.
func (ctx *funcContext) compileExtract(i *ssa.Extract) {
	ctx.compileValue(i.Tuple)
	ctx.emit(bytecode.OpConst, uint16(ctx.c.addConstant(i.Index)))
	ctx.emit(bytecode.OpIndex)
	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSend compiles a Send instruction.
func (ctx *funcContext) compileSend(i *ssa.Send) {
	ctx.compileValue(i.Chan)
	ctx.compileValue(i.X)
	ctx.emit(bytecode.OpSend)
}

// compileDefer compiles a Defer instruction.
func (ctx *funcContext) compileDefer(i *ssa.Defer) {
	switch val := i.Call.Value.(type) {
	case *ssa.Function:
		// If the function has free variables, we need to create a closure
		if len(val.FreeVars) > 0 {
			// First create the closure (OpDeferIndirect expects: closure, args... on stack)
			// Push the free variable bindings
			for _, fv := range val.FreeVars {
				ctx.compileValue(fv)
			}
			// Create the closure (leaves it on stack, no SETLOCAL)
			fnIdx := ctx.c.funcIndex[val]
			ctx.cf.Instructions = append(ctx.cf.Instructions,
				byte(bytecode.OpClosure),
				byte(fnIdx>>8), byte(fnIdx),
				byte(len(val.FreeVars)))
			// Push arguments AFTER closure
			for _, arg := range i.Call.Args {
				ctx.compileValue(arg)
			}
			numArgs := len(i.Call.Args)
			ctx.emit(bytecode.OpDeferIndirect, uint16(numArgs))
			return
		}
		// No free variables - push args then use OpDefer directly
		for _, arg := range i.Call.Args {
			ctx.compileValue(arg)
		}
		fnIdx := ctx.c.funcIndex[val]
		ctx.emit(bytecode.OpDefer, uint16(fnIdx))

	case *ssa.MakeClosure:
		// Check if this MakeClosure was already compiled (has a local slot)
		if idx, ok := ctx.symbolTable.GetLocal(val); ok {
			// Already compiled - load the closure from local FIRST
			ctx.emit(bytecode.OpLocal, uint16(idx))
			// Then push arguments
			for _, arg := range i.Call.Args {
				ctx.compileValue(arg)
			}
			numArgs := len(i.Call.Args)
			ctx.emit(bytecode.OpDeferIndirect, uint16(numArgs))
			return
		}
		// Not yet compiled - create the closure now FIRST
		// Compile bindings - need to handle FreeVar specially
		for _, binding := range val.Bindings {
			// Check if binding is a FreeVar (captured from enclosing function)
			if fv, ok := binding.(*ssa.FreeVar); ok {
				if idx, ok := ctx.symbolTable.freeVars[fv]; ok {
					ctx.emit(bytecode.OpFree, uint16(idx))
					continue
				}
			}
			// Handle Alloc (pointer variable)
			if alloc, ok := binding.(*ssa.Alloc); ok {
				if slotIdx, ok := ctx.symbolTable.GetLocal(alloc); ok {
					ctx.emit(bytecode.OpLocal, uint16(slotIdx))
					continue
				}
			}
			ctx.compileValue(binding)
		}
		// Create closure (on stack, no SETLOCAL)
		fnIdx := ctx.c.funcIndex[val.Fn.(*ssa.Function)]
		ctx.cf.Instructions = append(ctx.cf.Instructions,
			byte(bytecode.OpClosure),
			byte(fnIdx>>8), byte(fnIdx),
			byte(len(val.Bindings)))
		// Push arguments AFTER closure
		for _, arg := range i.Call.Args {
			ctx.compileValue(arg)
		}
		numArgs := len(i.Call.Args)
		ctx.emit(bytecode.OpDeferIndirect, uint16(numArgs))

	default:
		// Other cases - first compile the callable, then push args
		ctx.compileValue(i.Call.Value)
		for _, arg := range i.Call.Args {
			ctx.compileValue(arg)
		}
		numArgs := len(i.Call.Args)
		ctx.emit(bytecode.OpDeferIndirect, uint16(numArgs))
	}
}

// compileGo compiles a Go instruction.
func (ctx *funcContext) compileGo(i *ssa.Go) {
	for _, arg := range i.Call.Args {
		ctx.compileValue(arg)
	}

	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		funcIdx := ctx.c.funcIndex[fn]
		numArgs := len(i.Call.Args)
		ctx.cf.Instructions = append(ctx.cf.Instructions,
			byte(bytecode.OpGoCall),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		return
	}

	ctx.compileValue(i.Call.Value)

	numArgs := len(i.Call.Args)
	ctx.cf.Instructions = append(ctx.cf.Instructions,
		byte(bytecode.OpGoCallIndirect),
		byte(numArgs))
}
