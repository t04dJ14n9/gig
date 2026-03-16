package compiler

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/bytecode"
)

// compileInstruction compiles a single SSA instruction to bytecode.
func (c *compiler) compileInstruction(instr ssa.Instruction) {
	switch i := instr.(type) {
	case *ssa.Alloc:
		c.compileAlloc(i)
	case *ssa.BinOp:
		c.compileBinOp(i)
	case *ssa.UnOp:
		c.compileUnOp(i)
	case *ssa.Call:
		c.compileCall(i)
	case *ssa.ChangeInterface:
		c.compileChangeInterface(i)
	case *ssa.ChangeType:
		c.compileChangeType(i)
	case *ssa.Convert:
		c.compileConvert(i)
	case *ssa.Extract:
		c.compileExtract(i)
	case *ssa.Field:
		c.compileField(i)
	case *ssa.FieldAddr:
		c.compileFieldAddr(i)
	case *ssa.Index:
		c.compileIndex(i)
	case *ssa.IndexAddr:
		c.compileIndexAddr(i)
	case *ssa.Lookup:
		c.compileLookup(i)
	case *ssa.MakeInterface:
		c.compileMakeInterface(i)
	case *ssa.MakeClosure:
		c.compileMakeClosure(i)
	case *ssa.MakeChan:
		c.compileMakeChan(i)
	case *ssa.MakeMap:
		c.compileMakeMap(i)
	case *ssa.MakeSlice:
		c.compileMakeSlice(i)
	case *ssa.Next:
		c.compileNext(i)
	case *ssa.Phi:
		// Phi nodes are handled by inserting moves in predecessors
	case *ssa.Range:
		c.compileRange(i)
	case *ssa.Select:
		c.compileSelect(i)
	case *ssa.Slice:
		c.compileSlice(i)
	case *ssa.TypeAssert:
		c.compileTypeAssert(i)
	case *ssa.DebugRef:
		// Skip debug info
	case *ssa.Defer:
		c.compileDefer(i)
	case *ssa.Go:
		c.compileGo(i)
	case *ssa.MapUpdate:
		c.compileMapUpdate(i)
	case *ssa.Panic:
		// Handled in block terminator
	case *ssa.Return:
		c.compileReturn(i)
	case *ssa.RunDefers:
		c.emit(bytecode.OpRunDefers)
	case *ssa.Send:
		c.compileSend(i)
	case *ssa.Store:
		c.compileStore(i)
	case *ssa.Jump:
		// Handled in block terminator
	case *ssa.If:
		// Handled in block terminator
	}
}

// compileAlloc compiles an Alloc instruction (variable allocation).
func (c *compiler) compileAlloc(i *ssa.Alloc) {
	addrIdx := c.symbolTable.AllocLocal(i)
	elemType := i.Type().(*types.Pointer).Elem()
	typeIdx := c.addType(elemType)
	c.emit(bytecode.OpNew, uint16(typeIdx))
	c.emit(bytecode.OpSetLocal, uint16(addrIdx))
}

// compileBinOp compiles a binary operation.
func (c *compiler) compileBinOp(i *ssa.BinOp) {
	c.compileValue(i.X)
	c.compileValue(i.Y)

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

	c.emit(op)

	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileUnOp compiles a UnOp instruction.
func (c *compiler) compileUnOp(i *ssa.UnOp) {
	c.compileValue(i.X)

	var op bytecode.OpCode
	switch i.Op { //nolint:exhaustive
	case token.ADD:
		resultIdx := c.symbolTable.AllocLocal(i)
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	case token.SUB:
		op = bytecode.OpNeg
	case token.NOT:
		op = bytecode.OpNot
	case token.XOR:
		op = bytecode.OpXor
	case token.ARROW:
		if i.CommaOk {
			op = bytecode.OpRecvOk
		} else {
			op = bytecode.OpRecv
		}
	case token.MUL:
		op = bytecode.OpDeref
	}

	c.emit(op)

	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileCall compiles a function call instruction.
func (c *compiler) compileCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	// Handle interface method invocation (e.g., iface.Method())
	// SSA represents this with IsInvoke() == true, where Call.Value is the interface value
	// and Call.Method is the method being invoked.
	if i.Call.IsInvoke() {
		// Push the receiver (interface value) as the first argument
		c.compileValue(i.Call.Value)
		// Push remaining arguments
		for _, arg := range i.Call.Args {
			c.compileValue(arg)
		}
		methodInfo := &bytecode.ExternalMethodInfo{
			MethodName: i.Call.Method.Name(),
		}
		funcIdx := c.addConstant(methodInfo)
		numArgs := len(i.Call.Args) + 1 // +1 for receiver
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(bytecode.OpCallExternal),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	}

	if builtin, ok := i.Call.Value.(*ssa.Builtin); ok {
		c.compileBuiltinCall(builtin, i.Call.Args, resultIdx)
		return
	}

	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		if _, known := c.funcIndex[fn]; known {
			for _, arg := range i.Call.Args {
				c.compileValue(arg)
			}

			funcIdx := c.funcIndex[fn]
			numArgs := len(i.Call.Args)
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(bytecode.OpCall),
				byte(funcIdx>>8), byte(funcIdx),
				byte(numArgs))
			c.emit(bytecode.OpSetLocal, uint16(resultIdx))
			return
		}

		c.compileExternalStaticCall(i, fn, resultIdx)
		return
	}

	c.compileIndirectCall(i)
}

// compileBuiltinCall compiles a call to a builtin function.
func (c *compiler) compileBuiltinCall(builtin *ssa.Builtin, args []ssa.Value, resultIdx int) {
	name := builtin.Name()

	switch name {
	case "len":
		c.compileValue(args[0])
		c.emit(bytecode.OpLen)
	case "cap":
		c.compileValue(args[0])
		c.emit(bytecode.OpCap)
	case "append":
		for i, arg := range args {
			c.compileValue(arg)
			if i > 0 {
				c.emit(bytecode.OpAppend)
			}
		}
	case "copy":
		c.compileValue(args[0])
		c.compileValue(args[1])
		c.emit(bytecode.OpCopy)
	case "delete":
		c.compileValue(args[0])
		c.compileValue(args[1])
		c.emit(bytecode.OpDelete)
		return
	case "panic":
		c.compileValue(args[0])
		c.emit(bytecode.OpPanic)
		return
	case "print":
		for _, arg := range args {
			c.compileValue(arg)
		}
		c.emit(bytecode.OpPrint, uint16(len(args)))
		return
	case "println":
		for _, arg := range args {
			c.compileValue(arg)
		}
		c.emit(bytecode.OpPrintln, uint16(len(args)))
		return
	case "new":
		typeIdx := c.addType(args[0].Type())
		c.emit(bytecode.OpNew, uint16(typeIdx))
	case "make":
		c.compileMakeBuiltin(args)
	case "close":
		c.compileValue(args[0])
		c.emit(bytecode.OpClose)
		return
	case "ssa:wrapnilchk":
		// ssa:wrapnilchk checks that its argument is non-nil.
		// In our VM, we simply pass through the value (the nil check
		// is an optimization/safety check that we skip).
		c.compileValue(args[0])
	default:
		c.emit(bytecode.OpNil)
	}

	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeBuiltin compiles the make builtin.
func (c *compiler) compileMakeBuiltin(args []ssa.Value) {
	t := args[0].Type()
	typeIdx := c.addType(t)
	typeIdxConst := c.addConstant(int64(typeIdx))
	zeroIdx := c.addConstant(int64(0))

	switch t.(type) {
	case *types.Slice:
		c.emit(bytecode.OpConst, typeIdxConst)
		if len(args) >= 2 {
			c.compileValue(args[1])
		} else {
			c.emit(bytecode.OpConst, zeroIdx)
		}
		if len(args) >= 3 {
			c.compileValue(args[2])
		} else {
			c.emit(bytecode.OpConst, zeroIdx)
		}
		c.emit(bytecode.OpMakeSlice)
	case *types.Map:
		c.emit(bytecode.OpConst, typeIdxConst)
		if len(args) > 1 {
			c.compileValue(args[1])
		} else {
			c.emit(bytecode.OpConst, zeroIdx)
		}
		c.emit(bytecode.OpMakeMap)
	case *types.Chan:
		c.emit(bytecode.OpConst, typeIdxConst)
		if len(args) > 1 {
			c.compileValue(args[1])
		} else {
			c.emit(bytecode.OpConst, zeroIdx)
		}
		c.emit(bytecode.OpMakeChan)
	}
}
