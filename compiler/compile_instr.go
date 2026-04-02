// compile_instr.go compiles SSA instructions: calls, binops, unops, field access.
package compiler

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/external"
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

// binOpMap maps Go token operators to bytecode opcodes.
var binOpMap = map[token.Token]bytecode.OpCode{
	token.ADD:     bytecode.OpAdd,
	token.SUB:     bytecode.OpSub,
	token.MUL:     bytecode.OpMul,
	token.QUO:     bytecode.OpDiv,
	token.REM:     bytecode.OpMod,
	token.AND:     bytecode.OpAnd,
	token.OR:      bytecode.OpOr,
	token.XOR:     bytecode.OpXor,
	token.AND_NOT: bytecode.OpAndNot,
	token.SHL:     bytecode.OpLsh,
	token.SHR:     bytecode.OpRsh,
	token.EQL:     bytecode.OpEqual,
	token.NEQ:     bytecode.OpNotEqual,
	token.LSS:     bytecode.OpLess,
	token.LEQ:     bytecode.OpLessEq,
	token.GTR:     bytecode.OpGreater,
	token.GEQ:     bytecode.OpGreaterEq,
	token.LAND:    bytecode.OpAnd,
	token.LOR:     bytecode.OpOr,
}

// compileBinOp compiles a binary operation.
func (c *compiler) compileBinOp(i *ssa.BinOp) {
	c.compileValue(i.X)
	c.compileValue(i.Y)

	c.emit(binOpMap[i.Op])

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
		// Unary ^ (bitwise NOT) is not a single-stack op — it requires
		// pushing an all-ones constant then XORing. e.g. ^x = x ^ allOnes.
		// We must NOT emit OpXor directly here; that would pop 2 values
		// when only 1 (the operand) is on the stack, causing a panic.
		allOnes := allOnesConstant(i.X.Type())
		c.emit(bytecode.OpConst, uint16(c.addConstant(allOnes)))
		c.emit(bytecode.OpXor)
		resultIdx := c.symbolTable.AllocLocal(i)
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
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

	c.emit(op)

	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
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
		methodInfo := &external.ExternalMethodInfo{
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
		funcIdx := c.addConstant(methodInfo)
		numArgs := len(i.Call.Args) + 1 // +1 for receiver
		c.emitCallOp(bytecode.OpCallExternal, funcIdx, numArgs)
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
			c.emitCallOp(bytecode.OpCall, uint16(funcIdx), numArgs)
			c.emit(bytecode.OpSetLocal, uint16(resultIdx))
			return
		}

		c.compileExternalStaticCall(i, fn, resultIdx)
		return
	}

	c.compileIndirectCall(i)
}

// builtinOps maps simple builtin names to their single-opcode implementations.
// Builtins not in this map have custom compilation logic in compileBuiltinCall.
var builtinOps = map[string]bytecode.OpCode{
	"len":     bytecode.OpLen,
	"cap":     bytecode.OpCap,
	"copy":    bytecode.OpCopy,
	"panic":   bytecode.OpPanic,
	"recover": bytecode.OpRecover,
	"close":   bytecode.OpClose,
}

// compileBuiltinCall compiles a call to a builtin function.
func (c *compiler) compileBuiltinCall(builtin *ssa.Builtin, args []ssa.Value, resultIdx int) {
	name := builtin.Name()

	// Fast path: single-opcode builtins with 1 arg that push a value
	if op, ok := builtinOps[name]; ok {
		switch name {
		case "len", "cap":
			c.compileValue(args[0])
			c.emit(op)
		case "copy":
			c.compileValue(args[0])
			c.compileValue(args[1])
			c.emit(op)
		case "panic":
			c.compileValue(args[0])
			c.emit(op)
			return
		case "recover":
			// recover() takes no arguments
			c.emit(op)
		case "close":
			c.compileValue(args[0])
			c.emit(op)
			return
		}
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	}

	switch name {
	case "append":
		for i, arg := range args {
			c.compileValue(arg)
			if i > 0 {
				c.emit(bytecode.OpAppend)
			}
		}
	case "delete":
		c.compileValue(args[0])
		c.compileValue(args[1])
		c.emit(bytecode.OpDelete)
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
	case "ssa:wrapnilchk":
		c.compileValue(args[0])
	case "real":
		// real(complex) -> float
		c.compileValue(args[0])
		c.emit(bytecode.OpReal)
	case "imag":
		// imag(complex) -> float
		c.compileValue(args[0])
		c.emit(bytecode.OpImag)
	case "complex":
		// complex(real, imag) -> complex
		c.compileValue(args[0]) // real part
		c.compileValue(args[1]) // imag part
		c.emit(bytecode.OpComplex)
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
