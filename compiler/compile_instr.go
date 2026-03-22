package compiler

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/bytecode"
)

// compileInstruction compiles a single SSA instruction to bytecode.
func (ctx *funcContext) compileInstruction(instr ssa.Instruction) {
	switch i := instr.(type) {
	case *ssa.Alloc:
		ctx.compileAlloc(i)
	case *ssa.BinOp:
		ctx.compileBinOp(i)
	case *ssa.UnOp:
		ctx.compileUnOp(i)
	case *ssa.Call:
		ctx.compileCall(i)
	case *ssa.ChangeInterface:
		ctx.compileChangeInterface(i)
	case *ssa.ChangeType:
		ctx.compileChangeType(i)
	case *ssa.Convert:
		ctx.compileConvert(i)
	case *ssa.Extract:
		ctx.compileExtract(i)
	case *ssa.Field:
		ctx.compileField(i)
	case *ssa.FieldAddr:
		ctx.compileFieldAddr(i)
	case *ssa.Index:
		ctx.compileIndex(i)
	case *ssa.IndexAddr:
		ctx.compileIndexAddr(i)
	case *ssa.Lookup:
		ctx.compileLookup(i)
	case *ssa.MakeInterface:
		ctx.compileMakeInterface(i)
	case *ssa.MakeClosure:
		ctx.compileMakeClosure(i)
	case *ssa.MakeChan:
		ctx.compileMakeChan(i)
	case *ssa.MakeMap:
		ctx.compileMakeMap(i)
	case *ssa.MakeSlice:
		ctx.compileMakeSlice(i)
	case *ssa.Next:
		ctx.compileNext(i)
	case *ssa.Phi:
		// Phi nodes are handled by inserting moves in predecessors
	case *ssa.Range:
		ctx.compileRange(i)
	case *ssa.Select:
		ctx.compileSelect(i)
	case *ssa.Slice:
		ctx.compileSlice(i)
	case *ssa.TypeAssert:
		ctx.compileTypeAssert(i)
	case *ssa.DebugRef:
		// Skip debug info
	case *ssa.Defer:
		ctx.compileDefer(i)
	case *ssa.Go:
		ctx.compileGo(i)
	case *ssa.MapUpdate:
		ctx.compileMapUpdate(i)
	case *ssa.Panic:
		// Handled in block terminator
	case *ssa.Return:
		ctx.compileReturn(i)
	case *ssa.RunDefers:
		ctx.emit(bytecode.OpRunDefers)
	case *ssa.Send:
		ctx.compileSend(i)
	case *ssa.Store:
		ctx.compileStore(i)
	case *ssa.Jump:
		// Handled in block terminator
	case *ssa.If:
		// Handled in block terminator
	}
}

// compileAlloc compiles an Alloc instruction (variable allocation).
func (ctx *funcContext) compileAlloc(i *ssa.Alloc) {
	addrIdx := ctx.symbolTable.AllocLocal(i)
	elemType := i.Type().(*types.Pointer).Elem()
	typeIdx := ctx.c.addType(elemType)
	ctx.emit(bytecode.OpNew, uint16(typeIdx))
	ctx.emit(bytecode.OpSetLocal, uint16(addrIdx))
}

// compileBinOp compiles a binary operation.
func (ctx *funcContext) compileBinOp(i *ssa.BinOp) {
	ctx.compileValue(i.X)
	ctx.compileValue(i.Y)

	var op bytecode.OpCode
	switch i.Op { //nolint:exhaustive
	case token.ADD:
		op = bytecode.OpAdd
	case token.SUB:
		op = bytecode.OpSub
	case token.MUL:
		op = bytecode.OpMul
	case token.QUO:
		op = bytecode.OpDiv
	case token.REM:
		op = bytecode.OpMod
	case token.AND:
		op = bytecode.OpAnd
	case token.OR:
		op = bytecode.OpOr
	case token.XOR:
		op = bytecode.OpXor
	case token.AND_NOT:
		op = bytecode.OpAndNot
	case token.SHL:
		op = bytecode.OpLsh
	case token.SHR:
		op = bytecode.OpRsh
	case token.EQL:
		op = bytecode.OpEqual
	case token.NEQ:
		op = bytecode.OpNotEqual
	case token.LSS:
		op = bytecode.OpLess
	case token.LEQ:
		op = bytecode.OpLessEq
	case token.GTR:
		op = bytecode.OpGreater
	case token.GEQ:
		op = bytecode.OpGreaterEq
	case token.LAND:
		op = bytecode.OpAnd
	case token.LOR:
		op = bytecode.OpOr
	}

	ctx.emit(op)

	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileUnOp compiles a UnOp instruction.
func (ctx *funcContext) compileUnOp(i *ssa.UnOp) {
	ctx.compileValue(i.X)

	var op bytecode.OpCode
	switch i.Op { //nolint:exhaustive
	case token.ADD:
		resultIdx := ctx.symbolTable.AllocLocal(i)
		ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	case token.SUB:
		op = bytecode.OpNeg
	case token.NOT:
		op = bytecode.OpNot
	case token.XOR:
		// Unary ^ (bitwise NOT) is not a single-stack op — it requires
		// pushing an all-ones constant then XORing. e.g. ^x = x ^ allOnes.
		// We must NOT emit OpXor directly here; that would pop 2 values
		// when only 1 (the operand) is on the stack, causing a panic.
		allOnes := allOnesConstant(i.X.Type())
		ctx.emit(bytecode.OpConst, uint16(ctx.c.addConstant(allOnes)))
		ctx.emit(bytecode.OpXor)
		resultIdx := ctx.symbolTable.AllocLocal(i)
		ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	case token.ARROW:
		if i.CommaOk {
			op = bytecode.OpRecvOk
		} else {
			op = bytecode.OpRecv
		}
	case token.MUL:
		op = bytecode.OpDeref
	}

	ctx.emit(op)

	resultIdx := ctx.symbolTable.AllocLocal(i)
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// allOnesConstant returns the "all ones" value for a given numeric type,
// used to implement unary ^ (bitwise NOT) as x ^ allOnes.
// Returns an any so it can be passed to addConstant, which calls FromInterface
// to create the correct value.Kind (KindUint for unsigned, KindInt for signed).
func allOnesConstant(t types.Type) any {
	switch u := t.Underlying().(type) {
	case *types.Basic:
		switch u.Kind() {
		case types.Uint8:
			return uint8(0xFF)
		case types.Uint16:
			return uint16(0xFFFF)
		case types.Uint32:
			return uint32(0xFFFFFFFF)
		case types.Uint:
			return uint(^uint(0))
		case types.Uint64:
			return uint64(^uint64(0))
		case types.Uintptr:
			return uintptr(^uintptr(0))
		case types.Int8:
			return int8(-1)
		case types.Int16:
			return int16(-1)
		case types.Int32:
			return int32(-1)
		case types.Int:
			return int(^int(0))
		case types.Int64:
			return int64(-1)
		}
	}
	return int64(-1) // fallback
}

// compileCall compiles a function call instruction.
func (ctx *funcContext) compileCall(i *ssa.Call) {
	resultIdx := ctx.symbolTable.AllocLocal(i)

	// Handle interface method invocation (e.g., iface.Method())
	// SSA represents this with IsInvoke() == true, where Call.Value is the interface value
	// and Call.Method is the method being invoked.
	if i.Call.IsInvoke() {
		// Push the receiver (interface value) as the first argument
		ctx.compileValue(i.Call.Value)
		// Push remaining arguments
		for _, arg := range i.Call.Args {
			ctx.compileValue(arg)
		}
		methodInfo := &bytecode.ExternalMethodInfo{
			MethodName: i.Call.Method.Name(),
		}
		// For invoke calls, try to extract the concrete receiver type from the
		// interface value. This helps callCompiledMethod disambiguate methods
		// when multiple types define the same method name (e.g., Get, Add).
		if recvType := i.Call.Value.Type(); recvType != nil {
			// The receiver is an interface type; the concrete type is unknown statically.
			// However, the interface's method set constrains which types are valid.
			// We store the interface type name as a hint for fallback dispatch.
			if iface, ok := recvType.Underlying().(*types.Interface); ok {
				_ = iface // interface type available for future use
			}
			// If the value itself has a known concrete type (rare for invoke), use it.
			if named := extractNamedType(recvType); named != nil {
				methodInfo.ReceiverTypeName = named.Obj().Name()
			}
		}
		funcIdx := ctx.c.addConstant(methodInfo)
		numArgs := len(i.Call.Args) + 1 // +1 for receiver
		ctx.cf.Instructions = append(ctx.cf.Instructions,
			byte(bytecode.OpCallExternal),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	}

	if builtin, ok := i.Call.Value.(*ssa.Builtin); ok {
		ctx.compileBuiltinCall(builtin, i.Call.Args, resultIdx)
		return
	}

	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		if _, known := ctx.c.funcIndex[fn]; known {
			for _, arg := range i.Call.Args {
				ctx.compileValue(arg)
			}

			funcIdx := ctx.c.funcIndex[fn]
			numArgs := len(i.Call.Args)
			ctx.cf.Instructions = append(ctx.cf.Instructions,
				byte(bytecode.OpCall),
				byte(funcIdx>>8), byte(funcIdx),
				byte(numArgs))
			ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
			return
		}

		ctx.compileExternalStaticCall(i, fn, resultIdx)
		return
	}

	ctx.compileIndirectCall(i)
}

// compileBuiltinCall compiles a call to a builtin function.
func (ctx *funcContext) compileBuiltinCall(builtin *ssa.Builtin, args []ssa.Value, resultIdx int) {
	name := builtin.Name()

	switch name {
	case "len":
		ctx.compileValue(args[0])
		ctx.emit(bytecode.OpLen)
	case "cap":
		ctx.compileValue(args[0])
		ctx.emit(bytecode.OpCap)
	case "append":
		for i, arg := range args {
			ctx.compileValue(arg)
			if i > 0 {
				ctx.emit(bytecode.OpAppend)
			}
		}
	case "copy":
		ctx.compileValue(args[0])
		ctx.compileValue(args[1])
		ctx.emit(bytecode.OpCopy)
	case "delete":
		ctx.compileValue(args[0])
		ctx.compileValue(args[1])
		ctx.emit(bytecode.OpDelete)
		return
	case "panic":
		ctx.compileValue(args[0])
		ctx.emit(bytecode.OpPanic)
		return
	case "print":
		for _, arg := range args {
			ctx.compileValue(arg)
		}
		ctx.emit(bytecode.OpPrint, uint16(len(args)))
		return
	case "println":
		for _, arg := range args {
			ctx.compileValue(arg)
		}
		ctx.emit(bytecode.OpPrintln, uint16(len(args)))
		return
	case "new":
		typeIdx := ctx.c.addType(args[0].Type())
		ctx.emit(bytecode.OpNew, uint16(typeIdx))
	case "make":
		ctx.compileMakeBuiltin(args)
	case "close":
		ctx.compileValue(args[0])
		ctx.emit(bytecode.OpClose)
		return
	case "ssa:wrapnilchk":
		// ssa:wrapnilchk checks that its argument is non-nil.
		// In our VM, we simply pass through the value (the nil check
		// is an optimization/safety check that we skip).
		ctx.compileValue(args[0])
	default:
		ctx.emit(bytecode.OpNil)
	}

	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeBuiltin compiles the make builtin.
func (ctx *funcContext) compileMakeBuiltin(args []ssa.Value) {
	t := args[0].Type()
	typeIdx := ctx.c.addType(t)
	typeIdxConst := ctx.c.addConstant(int64(typeIdx))
	zeroIdx := ctx.c.addConstant(int64(0))

	switch t.(type) {
	case *types.Slice:
		ctx.emit(bytecode.OpConst, typeIdxConst)
		if len(args) >= 2 {
			ctx.compileValue(args[1])
		} else {
			ctx.emit(bytecode.OpConst, zeroIdx)
		}
		if len(args) >= 3 {
			ctx.compileValue(args[2])
		} else {
			ctx.emit(bytecode.OpConst, zeroIdx)
		}
		ctx.emit(bytecode.OpMakeSlice)
	case *types.Map:
		ctx.emit(bytecode.OpConst, typeIdxConst)
		if len(args) > 1 {
			ctx.compileValue(args[1])
		} else {
			ctx.emit(bytecode.OpConst, zeroIdx)
		}
		ctx.emit(bytecode.OpMakeMap)
	case *types.Chan:
		ctx.emit(bytecode.OpConst, typeIdxConst)
		if len(args) > 1 {
			ctx.compileValue(args[1])
		} else {
			ctx.emit(bytecode.OpConst, zeroIdx)
		}
		ctx.emit(bytecode.OpMakeChan)
	}
}
