package compiler

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
)

// CompiledFunction represents a compiled function.
type CompiledFunction struct {
	Name         string
	Instructions []byte
	NumLocals    int
	NumParams    int
	NumFreeVars  int
	MaxStack     int
	Source       *ssa.Function // for debugging
}

// Program represents a compiled program.
type Program struct {
	Functions map[string]*CompiledFunction
	Constants []any          // constant pool
	Globals   map[string]int // global name -> index
	MainPkg   *ssa.Package
	Types     []types.Type            // type pool (indexed by addType)
	FuncIndex map[*ssa.Function]int   // SSA function -> index for calls
}

// Compiler compiles SSA IR to bytecode.
type Compiler struct {
	program     *Program
	constants   []any
	types       []types.Type
	globals     map[string]int
	funcs       map[string]*CompiledFunction
	funcIndex   map[*ssa.Function]int
	currentFunc *CompiledFunction
	symbolTable *SymbolTable
	jumps       []jumpInfo // tracks jumps that need patching
	phiSlots    map[*ssa.Phi]int // Phi nodes -> local slots
}

// jumpInfo tracks a jump instruction that needs its target patched.
type jumpInfo struct {
	offset      int             // offset of the jump instruction in bytecode
	targetBlock *ssa.BasicBlock // the target basic block
}

// phiMove represents a move instruction for Phi elimination.
type phiMove struct {
	sourceValue ssa.Value // the value to copy
	targetSlot  int       // the local slot to copy to
}

// SymbolTable tracks SSA values to local slots.
type SymbolTable struct {
	locals    map[ssa.Value]int // SSA value -> local slot
	freeVars  map[ssa.Value]int // free var -> index
	numLocals int
}

// NewSymbolTable creates a new symbol table.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		locals:   make(map[ssa.Value]int),
		freeVars: make(map[ssa.Value]int),
	}
}

// AllocLocal allocates a new local slot for a value.
func (s *SymbolTable) AllocLocal(v ssa.Value) int {
	if idx, ok := s.locals[v]; ok {
		return idx
	}
	idx := s.numLocals
	s.locals[v] = idx
	s.numLocals++
	return idx
}

// GetLocal returns the local slot for a value.
func (s *SymbolTable) GetLocal(v ssa.Value) (int, bool) {
	idx, ok := s.locals[v]
	return idx, ok
}

// NumLocals returns the number of allocated locals.
func (s *SymbolTable) NumLocals() int {
	return s.numLocals
}

// NewCompiler creates a new compiler.
func NewCompiler() *Compiler {
	return &Compiler{
		constants: make([]any, 0),
		types:     make([]types.Type, 0),
		globals:   make(map[string]int),
		funcs:     make(map[string]*CompiledFunction),
		funcIndex: make(map[*ssa.Function]int),
	}
}

// Compile compiles an SSA package to bytecode.
func Compile(mainPkg *ssa.Package) (*Program, error) {
	c := NewCompiler()
	c.program = &Program{
		Functions: make(map[string]*CompiledFunction),
		Globals:   make(map[string]int),
		MainPkg:   mainPkg,
		Types:     make([]types.Type, 0),
		FuncIndex: make(map[*ssa.Function]int),
	}

	// Collect all functions (including anonymous/nested)
	var allFuncs []*ssa.Function
	var collectFuncs func(fn *ssa.Function)
	collectFuncs = func(fn *ssa.Function) {
		allFuncs = append(allFuncs, fn)
		for _, anon := range fn.AnonFuncs {
			collectFuncs(anon)
		}
	}
	for _, member := range mainPkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			collectFuncs(fn)
		}
	}

	// First pass: assign indices to all functions
	for idx, fn := range allFuncs {
		c.funcIndex[fn] = idx
		c.program.FuncIndex[fn] = idx
	}

	// Second pass: compile each function
	for _, fn := range allFuncs {
		compiled, err := c.compileFunction(fn)
		if err != nil {
			return nil, fmt.Errorf("compile function %s: %w", fn.Name(), err)
		}
		c.funcs[fn.Name()] = compiled
		c.program.Functions[fn.Name()] = compiled
	}

	c.program.Constants = c.constants
	c.program.Types = c.types

	return c.program, nil
}

// compileFunction compiles a single SSA function.
func (c *Compiler) compileFunction(fn *ssa.Function) (*CompiledFunction, error) {
	c.currentFunc = &CompiledFunction{
		Name:         fn.Name(),
		Instructions: make([]byte, 0),
		Source:       fn,
		NumParams:    len(fn.Params),
	}

	c.symbolTable = NewSymbolTable()
	c.jumps = nil // reset jumps for this function
	c.phiSlots = make(map[*ssa.Phi]int) // initialize phi slots

	// Allocate locals for parameters
	for _, param := range fn.Params {
		c.symbolTable.AllocLocal(param)
	}

	// Allocate locals for free variables (for closures)
	for i, freeVar := range fn.FreeVars {
		c.symbolTable.freeVars[freeVar] = i
	}

	// First pass: collect Phi nodes and allocate slots for them
	// Phi nodes are always at the beginning of blocks
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if phi, ok := instr.(*ssa.Phi); ok {
				slot := c.symbolTable.AllocLocal(phi)
				c.phiSlots[phi] = slot
			}
		}
	}

	// Allocate locals for all other values in the function
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if val, ok := instr.(ssa.Value); ok {
				// Skip Phi (already allocated) and Alloc (handled separately)
				if _, isPhi := instr.(*ssa.Phi); !isPhi {
					if _, isAlloc := instr.(*ssa.Alloc); !isAlloc {
						c.symbolTable.AllocLocal(val)
					}
				}
			}
		}
	}

	// Pre-allocate slots for Alloc instructions too
	// (they were skipped in the pass above)
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if _, isAlloc := instr.(*ssa.Alloc); isAlloc {
				c.symbolTable.AllocLocal(instr.(ssa.Value))
			}
		}
	}

	c.currentFunc.NumLocals = c.symbolTable.NumLocals()
	c.currentFunc.NumFreeVars = len(fn.FreeVars)

	// Compile basic blocks
	// Use reverse postorder for better jump optimization
	blocks := reversePostorder(fn)

	// Map blocks to instruction offsets
	blockOffsets := make(map[*ssa.BasicBlock]int)

	// First pass: compile blocks and record offsets
	for _, block := range blocks {
		blockOffsets[block] = len(c.currentFunc.Instructions)
		c.compileBlock(fn, block)
	}

	// Second pass: patch jump targets
	c.patchJumps(blockOffsets)

	return c.currentFunc, nil
}

// compileBlock compiles a single basic block.
func (c *Compiler) compileBlock(fn *ssa.Function, block *ssa.BasicBlock) {
	for _, instr := range block.Instrs {
		c.compileInstruction(fn, instr)
	}

	// Handle block terminator
	if block.Instrs != nil {
		last := block.Instrs[len(block.Instrs)-1]
		switch term := last.(type) {
		case *ssa.Return:
			// Already handled in compileInstruction
		case *ssa.Jump:
			// Jump to successor - emit Phi moves first
			c.emitPhiMoves(block, block.Succs[0])
			c.emitJump(block.Succs[0])
		case *ssa.If:
			// Conditional jump - compile the condition value
			c.compileValue(term.Cond)
			// Emit Phi moves for each successor
			c.emitPhiMoves(block, block.Succs[1]) // false branch
			c.emitJumpFalse(block.Succs[1])
			c.emitPhiMoves(block, block.Succs[0]) // true branch
			c.emitJump(block.Succs[0])
		case *ssa.Panic:
			// Panic instruction - banned, but we compile it for error handling
			c.compileValue(term.X)
			c.emit(OpPanic)
		}
	}
}

// emitPhiMoves emits move instructions for Phi nodes in the target block.
// This is called before jumping to the target block.
func (c *Compiler) emitPhiMoves(predBlock, targetBlock *ssa.BasicBlock) {
	// Find the index of predBlock in targetBlock's predecessors
	predIndex := -1
	for i, pred := range targetBlock.Preds {
		if pred == predBlock {
			predIndex = i
			break
		}
	}
	if predIndex < 0 {
		return // shouldn't happen, but be safe
	}

	// For each Phi instruction in the target block, emit a move
	for _, instr := range targetBlock.Instrs {
		phi, ok := instr.(*ssa.Phi)
		if !ok {
			break // Phi nodes are always at the beginning
		}

		// Get the value for this predecessor
		if predIndex < len(phi.Edges) {
			sourceValue := phi.Edges[predIndex]
			targetSlot := c.phiSlots[phi]

			// Emit: compile sourceValue, then store to targetSlot
			c.compileValue(sourceValue)
			c.emit(OpSetLocal, uint16(targetSlot))
		}
	}
}

// compileInstruction compiles a single SSA instruction.
func (c *Compiler) compileInstruction(fn *ssa.Function, instr ssa.Instruction) {
	switch i := instr.(type) {
	// Value-producing instructions
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

	// Non-value instructions
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
		c.emit(OpRecover) // Run deferred functions
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

// compileAlloc compiles an Alloc instruction.
func (c *Compiler) compileAlloc(i *ssa.Alloc) {
	// Allocate a local slot for the address
	addrIdx := c.symbolTable.AllocLocal(i)

	// Both heap and stack allocs need a real pointer
	// i.Type() is *T, so Elem() gives us T
	typeIdx := c.addType(i.Type().(*types.Pointer).Elem())
	c.emit(OpNew, uint16(typeIdx))
	c.emit(OpSetLocal, uint16(addrIdx))
}

// compileBinOp compiles a BinOp instruction.
func (c *Compiler) compileBinOp(i *ssa.BinOp) {
	// Compile operands
	c.compileValue(i.X)
	c.compileValue(i.Y)

	var op OpCode
	switch i.Op {
	case token.ADD:
		op = OpAdd
	case token.SUB:
		op = OpSub
	case token.MUL:
		op = OpMul
	case token.QUO:
		op = OpDiv
	case token.REM:
		op = OpMod
	case token.AND:
		op = OpAnd
	case token.OR:
		op = OpOr
	case token.XOR:
		op = OpXor
	case token.AND_NOT:
		op = OpAndNot
	case token.SHL:
		op = OpLsh
	case token.SHR:
		op = OpRsh
	case token.EQL:
		op = OpEqual
	case token.NEQ:
		op = OpNotEqual
	case token.LSS:
		op = OpLess
	case token.LEQ:
		op = OpLessEq
	case token.GTR:
		op = OpGreater
	case token.GEQ:
		op = OpGreaterEq
	case token.LAND:
		// Logical AND - should be decomposed by SSA, but handle it
		op = OpAnd
	case token.LOR:
		// Logical OR - should be decomposed by SSA, but handle it
		op = OpOr
	}

	c.emit(op)

	// Store result in local slot
	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileUnOp compiles a UnOp instruction.
func (c *Compiler) compileUnOp(i *ssa.UnOp) {
	// Compile operand
	c.compileValue(i.X)

	var op OpCode
	switch i.Op {
	case token.ADD:
		// Unary + - just pass through
		// Still need to store result
		resultIdx := c.symbolTable.AllocLocal(i)
		c.emit(OpSetLocal, uint16(resultIdx))
		return
	case token.SUB:
		op = OpNeg
	case token.NOT:
		op = OpNot
	case token.XOR:
		op = OpXor
	case token.ARROW:
		// Channel receive
		op = OpRecv
	case token.MUL:
		// Pointer dereference
		op = OpDeref
	}

	c.emit(op)

	// Store result
	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileCall compiles a Call instruction.
func (c *Compiler) compileCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	// Check if it's a builtin
	if builtin, ok := i.Call.Value.(*ssa.Builtin); ok {
		c.compileBuiltinCall(builtin, i.Call.Args, resultIdx)
		return
	}

	// Check if it's a static call
	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		// Push arguments
		for _, arg := range i.Call.Args {
			c.compileValue(arg)
		}

		// Call function - manually encode with proper operand sizes
		funcIdx := c.funcIndex[fn]
		numArgs := len(i.Call.Args)
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(OpCall),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		c.emit(OpSetLocal, uint16(resultIdx))
		return
	}

	// Method call or indirect call (closures, function values)
	c.compileIndirectCall(i)
}

// compileBuiltinCall compiles a builtin function call.
func (c *Compiler) compileBuiltinCall(builtin *ssa.Builtin, args []ssa.Value, resultIdx int) {
	name := builtin.Name()

	switch name {
	case "len":
		c.compileValue(args[0])
		c.emit(OpLen)
	case "cap":
		c.compileValue(args[0])
		c.emit(OpCap)
	case "append":
		for i, arg := range args {
			c.compileValue(arg)
			if i > 0 {
				c.emit(OpAppend)
			}
		}
	case "copy":
		c.compileValue(args[0])
		c.compileValue(args[1])
		c.emit(OpCopy)
	case "delete":
		c.compileValue(args[0])
		c.compileValue(args[1])
		c.emit(OpDelete)
	case "panic":
		c.compileValue(args[0])
		c.emit(OpPanic)
	case "print":
		for _, arg := range args {
			c.compileValue(arg)
		}
		c.emit(OpPrint, uint16(len(args)))
	case "println":
		for _, arg := range args {
			c.compileValue(arg)
		}
		c.emit(OpPrintln, uint16(len(args)))
	case "new":
		typeIdx := c.addType(args[0].Type())
		c.emit(OpNew, uint16(typeIdx))
	case "make":
		c.compileMakeBuiltin(args)
	default:
		// Unknown builtin
		c.emit(OpNil)
	}

	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileMakeBuiltin compiles the make builtin.
func (c *Compiler) compileMakeBuiltin(args []ssa.Value) {
	// make(T) or make(T, size) or make(T, len, cap)
	// First arg is always the type
	t := args[0].Type()
	typeIdx := c.addType(t)
	typeIdxConst := c.addConstant(int64(typeIdx))
	zeroIdx := c.addConstant(int64(0))

	switch t.(type) {
	case *types.Slice:
		c.emit(OpConst, typeIdxConst)
		if len(args) >= 2 {
			c.compileValue(args[1])
		} else {
			c.emit(OpConst, zeroIdx)
		}
		if len(args) >= 3 {
			c.compileValue(args[2])
		} else {
			c.emit(OpConst, zeroIdx)
		}
		c.emit(OpMakeSlice)
	case *types.Map:
		c.emit(OpConst, typeIdxConst)
		if len(args) > 1 {
			c.compileValue(args[1])
		} else {
			c.emit(OpConst, zeroIdx)
		}
		c.emit(OpMakeMap)
	case *types.Chan:
		c.emit(OpConst, typeIdxConst)
		if len(args) > 1 {
			c.compileValue(args[1])
		} else {
			c.emit(OpConst, zeroIdx)
		}
		c.emit(OpMakeChan)
	}
}

// compileExternalCall compiles an external function call.
func (c *Compiler) compileExternalCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	// Push arguments
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	// The method should be an external function
	funcIdx := c.addConstant(i.Call.Value)
	numArgs := len(i.Call.Args)
	// Manually encode OpCallExternal: [opcode(1)] [funcIdx(2)] [numArgs(1)]
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(OpCallExternal),
		byte(funcIdx>>8), byte(funcIdx),
		byte(numArgs))
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileIndirectCall compiles an indirect call (closure, function value).
func (c *Compiler) compileIndirectCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	// Push the callee (function value / closure) onto the stack
	c.compileValue(i.Call.Value)

	// Push arguments
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	numArgs := len(i.Call.Args)
	// Emit OpCallIndirect: [opcode(1)] [numArgs(1)]
	c.emit(OpCallIndirect, uint16(numArgs))
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileReturn compiles a Return instruction.
func (c *Compiler) compileReturn(i *ssa.Return) {
	if len(i.Results) == 0 {
		c.emit(OpReturn)
		return
	}

	// For single return value, just push and return
	if len(i.Results) == 1 {
		c.compileValue(i.Results[0])
		c.emit(OpReturnVal)
		return
	}

	// For multiple return values, push them all then pack into a slice
	for _, result := range i.Results {
		c.compileValue(result)
	}
	
	// Pack the values into a slice
	c.emit(OpPack, uint16(len(i.Results)))
	
	c.emit(OpReturnVal)
}

// compileValue compiles a value to push it onto the stack.
func (c *Compiler) compileValue(v ssa.Value) {
	switch val := v.(type) {
	case *ssa.Const:
		// Handle constants
		c.compileConst(val)
	case *ssa.Function:
		// Handle function references
		fnIdx := c.funcIndex[val]
		c.emit(OpConst, uint16(fnIdx))
	case *ssa.Phi:
		// Phi nodes - use the pre-allocated local slot
		if slot, ok := c.phiSlots[val]; ok {
			c.emit(OpLocal, uint16(slot))
		} else {
			c.emit(OpNil)
		}
	case *ssa.FreeVar:
		// Free variables (captured by closures)
		if idx, ok := c.symbolTable.freeVars[val]; ok {
			c.emit(OpFree, uint16(idx))
		} else {
			c.emit(OpNil)
		}
	default:
		// Handle values stored in locals
		if idx, ok := c.symbolTable.GetLocal(v); ok {
			c.emit(OpLocal, uint16(idx))
		} else {
			// Check if it's a free var by identity
			if idx, ok := c.symbolTable.freeVars[v]; ok {
				c.emit(OpFree, uint16(idx))
			} else {
				// Fallback: try to allocate and use
				c.emit(OpNil)
			}
		}
	}
}

// compileConst compiles a constant value.
func (c *Compiler) compileConst(cnst *ssa.Const) {
	// Convert constant to Value and add to constant pool
	var v any
	switch t := cnst.Type().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			v = cnst.Value != nil && cnst.Value.Kind() == constant.Bool && constant.BoolVal(cnst.Value)
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
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
		case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
			if cnst.Value != nil {
				v, _ = constant.Uint64Val(cnst.Value)
			} else {
				v = uint64(0)
			}
		case types.Float32, types.Float64:
			if cnst.Value != nil {
				f, _ := constant.Float64Val(cnst.Value)
				v = f
			} else {
				v = 0.0
			}
		case types.String:
			if cnst.Value != nil {
				v = constant.StringVal(cnst.Value)
			} else {
				v = ""
			}
		default:
			v = nil
		}
	default:
		v = nil
	}

	idx := c.addConstant(v)
	c.emit(OpConst, idx)
}

// compileField compiles a Field instruction.
func (c *Compiler) compileField(i *ssa.Field) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)
	c.emit(OpField, uint16(i.Field))
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileFieldAddr compiles a FieldAddr instruction.
func (c *Compiler) compileFieldAddr(i *ssa.FieldAddr) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)
	c.emit(OpAddr, uint16(i.Field))
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileIndex compiles an Index instruction.
func (c *Compiler) compileIndex(i *ssa.Index) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)
	c.compileValue(i.Index)
	c.emit(OpIndex)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileIndexAddr compiles an IndexAddr instruction.
func (c *Compiler) compileIndexAddr(i *ssa.IndexAddr) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)
	c.compileValue(i.Index)
	c.emit(OpIndexAddr)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileLookup compiles a Lookup instruction.
func (c *Compiler) compileLookup(i *ssa.Lookup) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)
	c.compileValue(i.Index)
	c.emit(OpIndex)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileStore compiles a Store instruction.
func (c *Compiler) compileStore(i *ssa.Store) {
	c.compileValue(i.Addr)
	c.compileValue(i.Val)
	c.emit(OpSetDeref)
}

// compileMakeSlice compiles a MakeSlice instruction.
func (c *Compiler) compileMakeSlice(i *ssa.MakeSlice) {
	typeIdx := c.addType(i.Type())
	resultIdx := c.symbolTable.AllocLocal(i)

	// Push typeIdx as an integer constant on the stack
	typeIdxConst := c.addConstant(int64(typeIdx))
	c.emit(OpConst, typeIdxConst)
	c.compileValue(i.Len)
	c.compileValue(i.Cap)
	c.emit(OpMakeSlice)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileMakeMap compiles a MakeMap instruction.
func (c *Compiler) compileMakeMap(i *ssa.MakeMap) {
	typeIdx := c.addType(i.Type())
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdxConst := c.addConstant(int64(typeIdx))
	c.emit(OpConst, typeIdxConst)

	if i.Reserve != nil {
		c.compileValue(i.Reserve)
	} else {
		c.emit(OpConst, uint16(c.addConstant(int64(0))))
	}

	c.emit(OpMakeMap)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileMakeChan compiles a MakeChan instruction.
func (c *Compiler) compileMakeChan(i *ssa.MakeChan) {
	typeIdx := c.addType(i.Type())
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdxConst := c.addConstant(int64(typeIdx))
	c.emit(OpConst, typeIdxConst)

	if i.Size != nil {
		c.compileValue(i.Size)
	} else {
		c.emit(OpConst, uint16(c.addConstant(int64(0))))
	}

	c.emit(OpMakeChan)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileMakeInterface compiles a MakeInterface instruction.
func (c *Compiler) compileMakeInterface(i *ssa.MakeInterface) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileMakeClosure compiles a MakeClosure instruction.
func (c *Compiler) compileMakeClosure(i *ssa.MakeClosure) {
	fnIdx := c.funcIndex[i.Fn.(*ssa.Function)]
	resultIdx := c.symbolTable.AllocLocal(i)

	// Push free variables
	for _, binding := range i.Bindings {
		c.compileValue(binding)
	}

	// Manually encode OpClosure: [opcode(1)] [funcIdx(2)] [numFree(1)]
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(OpClosure),
		byte(fnIdx>>8), byte(fnIdx),
		byte(len(i.Bindings)))
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileMapUpdate compiles a MapUpdate instruction.
func (c *Compiler) compileMapUpdate(i *ssa.MapUpdate) {
	c.compileValue(i.Map)
	c.compileValue(i.Key)
	c.compileValue(i.Value)
	c.emit(OpSetIndex)
}

// compileRange compiles a Range instruction.
func (c *Compiler) compileRange(i *ssa.Range) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)
	c.emit(OpRange)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileNext compiles a Next instruction.
func (c *Compiler) compileNext(i *ssa.Next) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.Iter)
	c.emit(OpRangeNext)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileSelect compiles a Select instruction.
func (c *Compiler) compileSelect(i *ssa.Select) {
	// Build select cases
	// This is complex - for now, use reflect.Select
	c.emit(OpSelect)
}

// compileSlice compiles a Slice instruction.
func (c *Compiler) compileSlice(i *ssa.Slice) {
	resultIdx := c.symbolTable.AllocLocal(i)

	// Compile the underlying value
	c.compileValue(i.X)

	if i.Low != nil {
		c.compileValue(i.Low)
	} else {
		// Push integer 0 using constant pool
		c.emit(OpConst, uint16(c.addConstant(int64(0))))
	}

	if i.High != nil {
		c.compileValue(i.High)
	} else {
		// Push max uint16 as marker for "no high"
		c.emit(OpConst, uint16(c.addConstant(int64(0xFFFF))))
	}

	if i.Max != nil {
		c.compileValue(i.Max)
	} else {
		// Push max uint16 as marker for "no max"
		c.emit(OpConst, uint16(c.addConstant(int64(0xFFFF))))
	}

	c.emit(OpSlice)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileTypeAssert compiles a TypeAssert instruction.
func (c *Compiler) compileTypeAssert(i *ssa.TypeAssert) {
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdx := c.addType(i.AssertedType)

	c.compileValue(i.X)
	c.emit(OpAssert, uint16(typeIdx))
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileChangeInterface compiles a ChangeInterface instruction.
func (c *Compiler) compileChangeInterface(i *ssa.ChangeInterface) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileChangeType compiles a ChangeType instruction.
func (c *Compiler) compileChangeType(i *ssa.ChangeType) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileConvert compiles a Convert instruction.
func (c *Compiler) compileConvert(i *ssa.Convert) {
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdx := c.addType(i.Type())

	c.compileValue(i.X)
	c.emit(OpConvert, uint16(typeIdx))
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileExtract compiles an Extract instruction.
func (c *Compiler) compileExtract(i *ssa.Extract) {
	// Get the tuple value
	c.compileValue(i.Tuple)

	// Push the index as a constant
	c.emit(OpConst, uint16(c.addConstant(i.Index)))

	// Extract the value at the given index
	c.emit(OpIndex)

	// Store result
	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileSend compiles a Send instruction.
func (c *Compiler) compileSend(i *ssa.Send) {
	c.compileValue(i.Chan)
	c.compileValue(i.X)
	c.emit(OpSend)
}

// compileDefer compiles a Defer instruction.
func (c *Compiler) compileDefer(i *ssa.Defer) {
	// Push arguments
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	// Defer the call
	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		fnIdx := c.funcIndex[fn]
		c.emit(OpDefer, uint16(fnIdx))
	}
}

// compileGo compiles a Go instruction.
func (c *Compiler) compileGo(i *ssa.Go) {
	// Push arguments
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	// Start goroutine
	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		fnIdx := c.funcIndex[fn]
		c.emit(OpCall, uint16(fnIdx), uint16(len(i.Call.Args)))
		c.emit(OpGo)
	}
}

// emit appends an opcode to the current function.
func (c *Compiler) emit(op OpCode, operands ...uint16) {
	c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(op))

	// Get operand width for this opcode
	width := OperandWidths[op]

	for _, operand := range operands {
		switch width {
		case 2:
			// Always write 2 bytes for 2-byte operands
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(operand>>8), byte(operand))
		case 1:
			// Write 1 byte for 1-byte operands
			c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(operand))
		default:
			// Variable width based on value
			if operand > 0xFF {
				c.currentFunc.Instructions = append(c.currentFunc.Instructions,
					byte(operand>>8), byte(operand))
			} else {
				c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(operand))
			}
		}
	}
}

// emitJump emits a jump instruction (placeholder, patched later).
func (c *Compiler) emitJump(target *ssa.BasicBlock) {
	offset := len(c.currentFunc.Instructions)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(OpJump), 0, 0)
	c.jumps = append(c.jumps, jumpInfo{offset: offset, targetBlock: target})
}

// emitJumpFalse emits a conditional jump (placeholder, patched later).
func (c *Compiler) emitJumpFalse(target *ssa.BasicBlock) {
	offset := len(c.currentFunc.Instructions)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(OpJumpFalse), 0, 0)
	c.jumps = append(c.jumps, jumpInfo{offset: offset, targetBlock: target})
}

// patchJumps patches jump targets with actual offsets.
func (c *Compiler) patchJumps(blockOffsets map[*ssa.BasicBlock]int) {
	for _, jump := range c.jumps {
		targetOffset := blockOffsets[jump.targetBlock]
		// The operand is at offset+1 (2 bytes)
		c.currentFunc.Instructions[jump.offset+1] = byte(targetOffset >> 8)
		c.currentFunc.Instructions[jump.offset+2] = byte(targetOffset)
	}
}

// addConstant adds a constant to the pool and returns its index.
func (c *Compiler) addConstant(val any) uint16 {
	idx := len(c.constants)
	c.constants = append(c.constants, val)
	return uint16(idx)
}

// addType adds a type to the pool and returns its index.
func (c *Compiler) addType(t types.Type) uint16 {
	idx := len(c.types)
	c.types = append(c.types, t)
	return uint16(idx)
}

// reversePostorder returns blocks in reverse postorder.
func reversePostorder(fn *ssa.Function) []*ssa.BasicBlock {
	if len(fn.Blocks) == 0 {
		return nil
	}

	visited := make(map[*ssa.BasicBlock]bool)
	var order []*ssa.BasicBlock

	var visit func(b *ssa.BasicBlock)
	visit = func(b *ssa.BasicBlock) {
		if visited[b] {
			return
		}
		visited[b] = true
		for _, succ := range b.Succs {
			visit(succ)
		}
		order = append(order, b)
	}

	visit(fn.Blocks[0])

	// Reverse to get postorder
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	return order
}
