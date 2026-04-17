// compile_value.go compiles SSA values: constants, conversions, slicing, type asserts.
package compiler

import (
	"go/constant"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/external"
)

// compileValue compiles an SSA value to push it onto the stack.
func (c *compiler) compileValue(v ssa.Value) {
	switch val := v.(type) {
	case *ssa.Const:
		c.compileConst(val)
	case *ssa.Function:
		if fnIdx, ok := c.funcIndex[val]; ok {
			c.emitClosure(fnIdx, 0)
		} else {
			// External function not in funcIndex — look up the actual Go function
			// and store it as a constant so it can be used as a value (e.g., passed
			// as a callback argument). OpCallIndirect handles reflect.Func values.
			c.compileExternalFuncValue(val)
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

			// Record the zero value for this global. SSA globals have type *T.
			// For external named struct types (e.g., sync.Mutex), allocate a
			// heap object via reflect.New(T) and store the POINTER in the global
			// slot. This ensures all method calls (including concurrent ones)
			// operate on the same underlying object — no copy, no write-back.
			// For basic types (int, string, etc.), store proper zero values.
			if ptrType, ok := val.Type().(*types.Pointer); ok {
				elemType := ptrType.Elem()
				switch t := elemType.(type) {
				case *types.Named:
					obj := t.Obj()
					if obj != nil && obj.Pkg() != nil && c.lookup != nil {
						pkgPath := obj.Pkg().Path()
						typeName := obj.Name()
						if rt, found := c.lookup.LookupExternalTypeByName(pkgPath, typeName); found {
							if rt.Kind() == reflect.Struct {
								// Store pointer *T, not value T. All method calls
								// will use this same heap-allocated object.
								c.globalZeroValues[globalIdx] = reflect.New(rt)
							}
						}
					}
				case *types.Basic:
					if rt, ok := bytecode.BasicKindToReflectType[t.Kind()]; ok {
						c.globalZeroValues[globalIdx] = reflect.Zero(rt)
					}
				}
			}

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
		c.emit(bytecode.OpGlobal, uint16(globalIdx))
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
		v = basicConstValue(t.Kind(), cnst.Value)
	case *types.Named, *types.Alias:
		// Named types and type aliases share the same compilation logic:
		// extract the underlying basic type for non-nil values.
		if cnst.Value != nil {
			if underlying, ok := t.Underlying().(*types.Basic); ok {
				v = basicConstValue(underlying.Kind(), cnst.Value)
			}
		} else {
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

// basicConstValue extracts a Go value from a constant.Value based on the basic type kind.
// Returns nil for unsupported kinds.
func basicConstValue(kind types.BasicKind, val constant.Value) any { //nolint:gocyclo,cyclop
	if val == nil {
		return basicZeroValue(kind)
	}

	switch kind { //nolint:exhaustive
	case types.Bool, types.UntypedBool:
		return val.Kind() == constant.Bool && constant.BoolVal(val)
	case types.Int, types.UntypedInt, types.UntypedRune:
		i, exact := constant.Int64Val(val)
		if exact {
			return int(i)
		}
		return int(0)
	case types.Int8:
		i, _ := constant.Int64Val(val)
		return int8(i)
	case types.Int16:
		i, _ := constant.Int64Val(val)
		return int16(i)
	case types.Int32:
		i, _ := constant.Int64Val(val)
		return int32(i)
	case types.Int64:
		i, exact := constant.Int64Val(val)
		if exact {
			return i
		}
		return int64(0)
	case types.Uint:
		u, _ := constant.Uint64Val(val)
		return uint(u)
	case types.Uint8:
		u, _ := constant.Uint64Val(val)
		return uint8(u)
	case types.Uint16:
		u, _ := constant.Uint64Val(val)
		return uint16(u)
	case types.Uint32:
		u, _ := constant.Uint64Val(val)
		return uint32(u)
	case types.Uint64:
		u, _ := constant.Uint64Val(val)
		return u
	case types.Uintptr:
		u, _ := constant.Uint64Val(val)
		return uint64(u)
	case types.Float32:
		f, _ := constant.Float64Val(val)
		return float32(f)
	case types.Float64, types.UntypedFloat:
		f, _ := constant.Float64Val(val)
		return f
	case types.String, types.UntypedString:
		return constant.StringVal(val)
	case types.Complex64:
		re := constant.Real(val)
		im := constant.Imag(val)
		reVal, _ := constant.Float64Val(re)
		imVal, _ := constant.Float64Val(im)
		return complex(float32(reVal), float32(imVal))
	case types.Complex128, types.UntypedComplex:
		re := constant.Real(val)
		im := constant.Imag(val)
		reVal, _ := constant.Float64Val(re)
		imVal, _ := constant.Float64Val(im)
		return complex(reVal, imVal)
	default:
		return nil
	}
}

// basicZeroValue returns the zero value for a basic type kind, or nil if unsupported.
var basicZeroValues = map[types.BasicKind]any{
	types.Bool: false, types.UntypedBool: false,
	types.Int: int(0), types.UntypedInt: int(0), types.UntypedRune: int(0),
	types.Int8: int8(0), types.Int16: int16(0), types.Int32: int32(0), types.Int64: int64(0),
	types.Uint: uint(0), types.Uint8: uint8(0), types.Uint16: uint16(0),
	types.Uint32: uint32(0), types.Uint64: uint64(0), types.Uintptr: uint64(0),
	types.Float32: float32(0), types.Float64: 0.0, types.UntypedFloat: 0.0,
	types.String: "", types.UntypedString: "",
	types.Complex64: complex64(0), types.Complex128: complex128(0), types.UntypedComplex: complex128(0),
}

func basicZeroValue(kind types.BasicKind) any {
	return basicZeroValues[kind] // nil for unsupported kinds
}

// isEmptyStruct checks if a type is an empty struct (struct{}).
func isEmptyStruct(t types.Type) bool {
	// Named or Alias wrapper
	switch u := t.(type) {
	case *types.Named:
		t = u.Underlying()
	case *types.Alias:
		t = u.Underlying()
	}
	// Check the underlying Struct type
	if st, ok := t.(*types.Struct); ok {
		return st.NumFields() == 0
	}
	return false
}

func constTypeToReflect(t types.Type) reflect.Type {
	// Handle empty structs early (Named, Alias, and direct Struct types)
	if isEmptyStruct(t) {
		return reflect.TypeFor[struct{}]()
	}

	switch typ := t.Underlying().(type) {
	case *types.Basic:
		return bytecode.BasicKindToReflectType[typ.Kind()]
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
			return reflect.ChanOf(chanDirection(typ), elemRT)
		}
	case *types.Interface:
		if typ.NumMethods() == 0 {
			return reflect.TypeFor[any]()
		}
	case *types.Signature:
		return buildFuncType(typ)
	}
	return nil
}

// chanDirection returns the reflect.ChanDir for a types.Chan.
func chanDirection(typ *types.Chan) reflect.ChanDir {
	switch typ.Dir() {
	case types.SendOnly:
		return reflect.SendDir
	case types.RecvOnly:
		return reflect.RecvDir
	default:
		return reflect.BothDir
	}
}

// buildFuncType builds a reflect.Type for a function signature.
func buildFuncType(sig *types.Signature) reflect.Type {
	params := make([]reflect.Type, sig.Params().Len())
	for i := 0; i < sig.Params().Len(); i++ {
		pt := constTypeToReflect(sig.Params().At(i).Type())
		if pt == nil {
			return nil
		}
		params[i] = pt
	}
	results := make([]reflect.Type, sig.Results().Len())
	for i := 0; i < sig.Results().Len(); i++ {
		rt := constTypeToReflect(sig.Results().At(i).Type())
		if rt == nil {
			return nil
		}
		results[i] = rt
	}
	return reflect.FuncOf(params, results, sig.Variadic())
}

// compileField compiles a Field instruction.
func (c *compiler) compileField(i *ssa.Field) {
	c.compileSimpleUnaryOp(i, i.X, bytecode.OpField, uint16(i.Field))
}

// compileFieldAddr compiles a FieldAddr instruction.
func (c *compiler) compileFieldAddr(i *ssa.FieldAddr) {
	c.compileSimpleUnaryOp(i, i.X, bytecode.OpFieldAddr, uint16(i.Field))
}

// compileIndex compiles an Index instruction.
func (c *compiler) compileIndex(i *ssa.Index) {
	c.compileSimpleBinaryOp(i, i.X, i.Index, bytecode.OpIndex)
}

// compileIndexAddr compiles an IndexAddr instruction.
func (c *compiler) compileIndexAddr(i *ssa.IndexAddr) {
	c.compileSimpleBinaryOp(i, i.X, i.Index, bytecode.OpIndexAddr)
}

// compileLookup compiles a Lookup instruction.
func (c *compiler) compileLookup(i *ssa.Lookup) {
	opcode := bytecode.OpIndex
	if i.CommaOk {
		opcode = bytecode.OpIndexOk
	}
	c.compileSimpleBinaryOp(i, i.X, i.Index, opcode)
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
	ifaceTypeIdx := c.addType(i.Type())     // interface type (e.g., error)
	concreteTypeIdx := c.addType(i.X.Type()) // concrete type (e.g., *MyError3)
	c.compileValue(i.X)
	// Emit: OpMakeInterface, ifaceTypeIdx_hi, ifaceTypeIdx_lo, concreteTypeIdx_hi, concreteTypeIdx_lo
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpMakeInterface),
		byte(ifaceTypeIdx>>8), byte(ifaceTypeIdx),
		byte(concreteTypeIdx>>8), byte(concreteTypeIdx))
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

	c.emitClosure(fnIdx, len(i.Bindings))
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
	c.compileSimpleUnaryOp(i, i.X, bytecode.OpRange)
}

// compileNext compiles a Next instruction.
func (c *compiler) compileNext(i *ssa.Next) {
	c.compileSimpleUnaryOp(i, i.Iter, bytecode.OpRangeNext)
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
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(bytecode.SliceEndSentinel))))
	}

	if i.Max != nil {
		c.compileValue(i.Max)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(bytecode.SliceEndSentinel))))
	}

	c.emit(bytecode.OpSlice)

	// If the result type is a named slice type (e.g., sort.IntSlice from [5]int[:]),
	// emit OpChangeType to convert the underlying []int to the named type.
	if named, ok := i.Type().(*types.Named); ok {
		if _, isSlice := named.Underlying().(*types.Slice); isSlice {
			typeIdx := c.addType(named)
			srcLocalIdx := uint16(bytecode.NoSourceLocal)
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(bytecode.OpChangeType),
				byte(typeIdx>>8), byte(typeIdx),
				byte(srcLocalIdx>>8), byte(srcLocalIdx))
		}
	}

	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileTypeAssert compiles a TypeAssert instruction.
func (c *compiler) compileTypeAssert(i *ssa.TypeAssert) {
	typeIdx := c.addType(i.AssertedType)
	c.compileValue(i.X)
	c.emit(bytecode.OpAssert, uint16(typeIdx))

	if !i.CommaOk {
		// Non-comma-ok assertion: SSA's `typeassert t.(T)` (without comma-ok)
		// panics on failure by branching to the recover block. We must check
		// the ok value and emit OpPanic if the assertion fails.
		// Stack has: [result, ok] tuple
		// Duplicate the tuple, extract ok (index 1), check if false → panic.
		c.emit(bytecode.OpDup)                             // [tuple, tuple]
		c.emit(bytecode.OpConst, uint16(c.addConstant(1))) // [tuple, tuple, 1]
		c.emit(bytecode.OpIndex)                           // [tuple, ok]
		// Emit OpJumpTrue with placeholder offset (3 bytes: opcode + u16 offset)
		jumpTrueOffset := len(c.currentFunc.Instructions)
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(bytecode.OpJumpTrue), 0, 0)
		// ok was false — panic with a type assertion error
		c.emit(bytecode.OpConst, uint16(c.addConstant("interface conversion: type assertion failed")))
		c.emit(bytecode.OpPanic)
		// Patch the JumpTrue to land here (ok case)
		skipOffset := len(c.currentFunc.Instructions)
		c.currentFunc.Instructions[jumpTrueOffset+1] = byte(skipOffset >> 8)
		c.currentFunc.Instructions[jumpTrueOffset+2] = byte(skipOffset)
		// Stack still has: [tuple] — extract the value
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
	srcLocalIdx := uint16(bytecode.NoSourceLocal)
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
	// Interface method invocation (e.g., defer iface.Method())
	if i.Call.IsInvoke() {
		c.compileDeferInvoke(i)
		return
	}

	switch val := i.Call.Value.(type) {
	case *ssa.Function:
		c.compileDeferFunction(i, val)
	case *ssa.MakeClosure:
		c.compileDeferMakeClosure(i, val)
	default:
		// Other cases: compile the callable, then push args
		c.compileValue(i.Call.Value)
		c.compileDeferCallArgs(i.Call.Args)
		c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
	}
}

// compileDeferInvoke handles defer of an interface method call.
func (c *compiler) compileDeferInvoke(i *ssa.Defer) {
	c.compileValue(i.Call.Value)
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}
	methodInfo := &external.ExternalMethodInfo{
		MethodName: i.Call.Method.Name(),
	}
	if recvType := i.Call.Value.Type(); recvType != nil {
		if named := extractNamedType(recvType); named != nil {
			methodInfo.ReceiverTypeName = named.Obj().Name()
		}
	}
	funcIdx := c.addConstant(methodInfo)
	c.emitCallOp(bytecode.OpDeferExternal, uint16(funcIdx), len(i.Call.Args)+1)
}

// compileDeferFunction handles defer of a static function call.
func (c *compiler) compileDeferFunction(i *ssa.Defer, val *ssa.Function) {
	// Known internal function
	if _, known := c.funcIndex[val]; known {
		if len(val.FreeVars) > 0 {
			// Has free variables — create closure, then push args
			fnIdx := c.funcIndex[val]
			c.compileAndEmitClosureFromFreeVars(val.FreeVars, fnIdx)
			c.compileDeferCallArgs(i.Call.Args)
			c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
			return
		}
		// No free variables — push args, use OpDefer directly
		c.compileDeferCallArgs(i.Call.Args)
		c.emit(bytecode.OpDefer, uint16(c.funcIndex[val]))
		return
	}

	// External method wrapper (not in funcIndex)
	if val.Signature.Recv() != nil {
		c.compileDeferCallArgs(i.Call.Args)
		methodName := extractMethodName(val.Name())
		methodInfo := &external.ExternalMethodInfo{MethodName: methodName}
		if c.lookup != nil {
			typeName := extractReceiverTypeName(val.Signature.Recv().Type())
			if typeName != "" {
				if dc, ok := c.lookup.LookupMethodDirectCall(typeName, methodName); ok {
					methodInfo.DirectCall = dc
				}
			}
		}
		funcIdx := c.addConstant(methodInfo)
		c.emitCallOp(bytecode.OpDeferExternal, uint16(funcIdx), len(i.Call.Args))
		return
	}

	// External package function
	c.compileValue(i.Call.Value)
	c.compileDeferCallArgs(i.Call.Args)
	c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
}

// compileDeferMakeClosure handles defer of a closure expression.
func (c *compiler) compileDeferMakeClosure(i *ssa.Defer, val *ssa.MakeClosure) {
	// Already compiled — load from local
	if idx, ok := c.symbolTable.GetLocal(val); ok {
		c.emit(bytecode.OpLocal, uint16(idx))
		c.compileDeferCallArgs(i.Call.Args)
		c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
		return
	}
	// Not yet compiled — create the closure now
	for _, binding := range val.Bindings {
		if fv, ok := binding.(*ssa.FreeVar); ok {
			if idx, ok := c.symbolTable.freeVars[fv]; ok {
				c.emit(bytecode.OpFree, uint16(idx))
				continue
			}
		}
		if alloc, ok := binding.(*ssa.Alloc); ok {
			if slotIdx, ok := c.symbolTable.GetLocal(alloc); ok {
				c.emit(bytecode.OpLocal, uint16(slotIdx))
				continue
			}
		}
		c.compileValue(binding)
	}
	c.emitClosure(c.funcIndex[val.Fn.(*ssa.Function)], len(val.Bindings))
	c.compileDeferCallArgs(i.Call.Args)
	c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
}

// compileGo compiles a Go instruction.
func (c *compiler) compileGo(i *ssa.Go) {
	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		// Check if this is a known internal function
		if _, known := c.funcIndex[fn]; known {
			// If the function has free variables, we need to create a closure
			// so the child goroutine can access captured variables (e.g., channels,
			// mutexes from the enclosing scope). Without this, OpGoCall passes
			// nil for freeVars and the child VM cannot access them.
			if len(fn.FreeVars) > 0 {
				// Create closure, then push arguments
				fnIdx := c.funcIndex[fn]
				c.compileAndEmitClosureFromFreeVars(fn.FreeVars, fnIdx)
				// Push arguments AFTER closure
				for _, arg := range i.Call.Args {
					c.compileValue(arg)
				}
				c.emit(bytecode.OpGoCallIndirect, uint16(len(i.Call.Args)))
				return
			}
			// No free variables — use OpGoCall directly
			for _, arg := range i.Call.Args {
				c.compileValue(arg)
			}
			c.emitCallOp(bytecode.OpGoCall, uint16(c.funcIndex[fn]), len(i.Call.Args))
			return
		}

		// External function/method wrapper (not in funcIndex).
		// Fall through to OpGoCallIndirect path which handles external callables.
	}

	// Indirect call (closure or MakeClosure result): push callee FIRST,
	// then args. OpGoCallIndirect pops args first, then callee.
	c.compileValue(i.Call.Value)

	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	c.emit(bytecode.OpGoCallIndirect, uint16(len(i.Call.Args)))
}
