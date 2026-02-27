package compiler

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/ssa"

	"gig/importer"
	"gig/value"
)

// CompiledFunction represents a function compiled to bytecode.
// It contains the bytecode instructions, local variable count, and metadata.
type CompiledFunction struct {
	// Name is the function name for debugging.
	Name string

	// Instructions is the compiled bytecode.
	Instructions []byte

	// NumLocals is the number of local variable slots.
	// This includes parameters and intermediate values.
	NumLocals int

	// NumParams is the number of function parameters.
	NumParams int

	// NumFreeVars is the number of free variables (for closures).
	NumFreeVars int

	// MaxStack is the maximum stack depth (for future optimization).
	MaxStack int

	// Source is the original SSA function for debugging.
	Source *ssa.Function
}

// ExternalFuncInfo contains pre-resolved external function info for fast calls.
// This allows the VM to bypass reflection when calling external functions.
type ExternalFuncInfo struct {
	// Func is the actual function value.
	Func any

	// DirectCall is a typed wrapper that avoids reflect.Call.
	// If nil, the VM will use reflection for the call.
	DirectCall func([]value.Value) value.Value
}

// Program represents a compiled program ready for execution.
// It contains all compiled functions, constants, types, and global variables.
type Program struct {
	// Functions maps function names to their compiled bytecode.
	Functions map[string]*CompiledFunction

	// Constants is the constant pool for literal values and external references.
	Constants []any

	// Globals maps global variable names to their indices.
	Globals map[string]int

	// MainPkg is the SSA package (for debugging/inspection).
	MainPkg *ssa.Package

	// Types is the type pool for runtime type operations.
	Types []types.Type

	// FuncIndex maps SSA functions to their indices for call instructions.
	FuncIndex map[*ssa.Function]int
}

// Compiler compiles SSA IR to bytecode.
// It maintains state during compilation including the current function,
// symbol table, and jump targets that need patching.
type Compiler struct {
	// program is the output program being compiled.
	program *Program

	// constants is the constant pool being built.
	constants []any

	// types is the type pool being built.
	types []types.Type

	// globals maps global names to indices.
	globals map[string]int

	// funcs maps function names to compiled functions.
	funcs map[string]*CompiledFunction

	// funcIndex maps SSA functions to call indices.
	funcIndex map[*ssa.Function]int

	// currentFunc is the function being compiled.
	currentFunc *CompiledFunction

	// symbolTable tracks SSA values to local slots.
	symbolTable *SymbolTable

	// jumps tracks jump instructions needing target patching.
	jumps []jumpInfo

	// phiSlots maps Phi nodes to their allocated local slots.
	phiSlots map[*ssa.Phi]int
}

// jumpInfo tracks a jump instruction that needs its target patched.
// Jumps are emitted with placeholder targets, then patched after
// all basic blocks have been compiled.
type jumpInfo struct {
	// offset is the bytecode offset of the jump instruction.
	offset int

	// targetBlock is the SSA basic block to jump to.
	targetBlock *ssa.BasicBlock
}

// phiMove represents a move instruction for Phi elimination.
// Phi nodes are eliminated by copying values in predecessor blocks.
type phiMove struct {
	// sourceValue is the SSA value to copy.
	sourceValue ssa.Value

	// targetSlot is the local slot to copy to.
	targetSlot int
}

// SymbolTable tracks SSA values to local slots.
// This maps SSA instructions and values to their storage locations
// in the VM's local variable array.
type SymbolTable struct {
	// locals maps SSA values to local slot indices.
	locals map[ssa.Value]int

	// freeVars maps free variables to their indices.
	freeVars map[ssa.Value]int

	// numLocals is the total number of allocated slots.
	numLocals int
}

// NewSymbolTable creates a new symbol table for tracking SSA values.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		locals:   make(map[ssa.Value]int),
		freeVars: make(map[ssa.Value]int),
	}
}

// AllocLocal allocates a new local slot for an SSA value.
// If the value already has a slot, returns the existing slot.
// Otherwise, allocates a new slot and returns its index.
func (s *SymbolTable) AllocLocal(v ssa.Value) int {
	if idx, ok := s.locals[v]; ok {
		return idx
	}
	idx := s.numLocals
	s.locals[v] = idx
	s.numLocals++
	return idx
}

// GetLocal returns the local slot index for an SSA value.
// Returns the index and true if found, or 0 and false if not.
func (s *SymbolTable) GetLocal(v ssa.Value) (int, bool) {
	idx, ok := s.locals[v]
	return idx, ok
}

// NumLocals returns the number of allocated local slots.
func (s *SymbolTable) NumLocals() int {
	return s.numLocals
}

// NewCompiler creates a new compiler with empty pools.
func NewCompiler() *Compiler {
	return &Compiler{
		constants: make([]any, 0),
		types:     make([]types.Type, 0),
		globals:   make(map[string]int),
		funcs:     make(map[string]*CompiledFunction),
		funcIndex: make(map[*ssa.Function]int),
	}
}

// Compile compiles an SSA package to a bytecode Program.
//
// The compilation process:
//  1. Collect all functions (including nested/anonymous)
//  2. Assign indices to functions for call instructions
//  3. Compile each function to bytecode
//  4. Assemble the final Program
//
// Returns the compiled Program or an error if compilation fails.
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
	c.program.Globals = c.globals

	return c.program, nil
}

// compileFunction compiles a single SSA function to bytecode.
//
// The process:
//  1. Allocate local slots for parameters and free variables
//  2. Allocate slots for Phi nodes (for SSA phi elimination)
//  3. Allocate slots for all other values
//  4. Compile basic blocks in reverse postorder
//  5. Patch jump targets
//
// Phi elimination is done by storing to the phi's slot in each predecessor block
// before jumping to the block containing the phi.
func (c *Compiler) compileFunction(fn *ssa.Function) (*CompiledFunction, error) {
	c.currentFunc = &CompiledFunction{
		Name:         fn.Name(),
		Instructions: make([]byte, 0),
		Source:       fn,
		NumParams:    len(fn.Params),
	}

	c.symbolTable = NewSymbolTable()
	c.jumps = nil                       // reset jumps for this function
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

// compileBlock compiles a single basic block to bytecode.
// It compiles each instruction in order, then handles the block terminator
// (Jump, If, Return, or Panic).
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

// emitPhiMoves emits move instructions for Phi nodes before jumping to a block.
//
// In SSA form, Phi nodes select a value based on which predecessor block
// was executed. In bytecode, we implement this by storing the appropriate
// value to the phi's slot in each predecessor before jumping.
//
// For each Phi in the target block, this function emits:
//  1. Compile the value for this predecessor
//  2. Store to the phi's local slot
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

// compileInstruction compiles a single SSA instruction to bytecode.
// It dispatches to the appropriate compile function based on instruction type.
//
// Value-producing instructions store their result in a local slot.
// Non-value instructions (Defer, Go, MapUpdate, etc.) produce side effects.
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

// compileAlloc compiles an Alloc instruction (variable allocation).
// Alloc creates a pointer to a zero value of the specified type.
// Both heap and stack allocs are handled the same way (heap promotion).
func (c *Compiler) compileAlloc(i *ssa.Alloc) {
	// Allocate a local slot for the address
	addrIdx := c.symbolTable.AllocLocal(i)

	// Both heap and stack allocs need a real pointer
	// i.Type() is *T, so Elem() gives us T
	typeIdx := c.addType(i.Type().(*types.Pointer).Elem())
	c.emit(OpNew, uint16(typeIdx))
	c.emit(OpSetLocal, uint16(addrIdx))
}

// compileBinOp compiles a binary operation (arithmetic, comparison, logical).
// It compiles both operands and emits the appropriate opcode.
func (c *Compiler) compileBinOp(i *ssa.BinOp) {
	// Compile operands
	c.compileValue(i.X)
	c.compileValue(i.Y)

	var op OpCode
	switch i.Op { //nolint:exhaustive
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
	switch i.Op { //nolint:exhaustive
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
		if i.CommaOk {
			// Receive with comma-ok: returns (value, ok) tuple
			op = OpRecvOk
		} else {
			op = OpRecv
		}
	case token.MUL:
		// Pointer dereference
		op = OpDeref
	}

	c.emit(op)

	// Store result
	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileCall compiles a function call instruction.
// It handles:
//   - Builtin functions (len, cap, append, etc.)
//   - Static calls to compiled functions
//   - Static calls to external functions
//   - Indirect calls (closures, function values)
func (c *Compiler) compileCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	// Check if it's a builtin
	if builtin, ok := i.Call.Value.(*ssa.Builtin); ok {
		c.compileBuiltinCall(builtin, i.Call.Args, resultIdx)
		return
	}

	// Check if it's a static call
	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		// Check if this is a function we compiled (from the main package)
		if _, known := c.funcIndex[fn]; known {
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

		// External function from another package
		c.compileExternalStaticCall(i, fn, resultIdx)
		return
	}

	// Method call or indirect call (closures, function values)
	c.compileIndirectCall(i)
}

// compileBuiltinCall compiles a call to a builtin function.
// Supported builtins: len, cap, append, copy, delete, close, panic, print, println, new, make.
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
		return // delete has no return value
	case "panic":
		c.compileValue(args[0])
		c.emit(OpPanic)
		return // panic has no return value
	case "print":
		for _, arg := range args {
			c.compileValue(arg)
		}
		c.emit(OpPrint, uint16(len(args)))
		return // print has no return value
	case "println":
		for _, arg := range args {
			c.compileValue(arg)
		}
		c.emit(OpPrintln, uint16(len(args)))
		return // println has no return value
	case "new":
		typeIdx := c.addType(args[0].Type())
		c.emit(OpNew, uint16(typeIdx))
	case "make":
		c.compileMakeBuiltin(args)
	case "close":
		c.compileValue(args[0])
		c.emit(OpClose)
		return // close has no return value
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

// compileExternalCall compiles an external function call (non-static / unknown callee).
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

// ExternalMethodInfo contains info for dispatching method calls on external types.
// This is used when calling methods on types from external packages via reflection.
type ExternalMethodInfo struct {
	// MethodName is the name of the method (e.g., "String", "Error").
	MethodName string
}

// compileExternalStaticCall compiles a call to an external package function.
// It looks up the function in the importer registry and emits OpCallExternal.
// For methods, it emits ExternalMethodInfo for VM dispatch.
func (c *Compiler) compileExternalStaticCall(i *ssa.Call, fn *ssa.Function, resultIdx int) {
	// Push arguments (for methods, the first arg is the receiver)
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	// Check if this is a method call (receiver != nil in the signature)
	sig := fn.Signature
	if sig.Recv() != nil {
		// This is a method call on an external type.
		// SSA passes the receiver as the first element in i.Call.Args.
		// We use OpCallExternal with ExternalMethodInfo so the VM can dispatch via reflect.
		methodName := fn.Name()
		// SSA method names may be qualified like "(*Type).Method" — extract the bare name
		if idx := strings.LastIndex(methodName, "."); idx >= 0 {
			methodName = methodName[idx+1:]
		}
		// Strip parenthesized receiver prefix if present, e.g., "(Result).String" -> "String"
		if idx := strings.LastIndex(methodName, ")"); idx >= 0 {
			rest := methodName[idx+1:]
			if len(rest) > 0 && rest[0] == '.' {
				methodName = rest[1:]
			}
		}

		methodInfo := &ExternalMethodInfo{
			MethodName: methodName,
		}
		funcIdx := c.addConstant(methodInfo)
		numArgs := len(i.Call.Args)
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(OpCallExternal),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		c.emit(OpSetLocal, uint16(resultIdx))
		return
	}

	// Look up the external function in the importer registry
	// fn.Pkg.Pkg.Path() gives the import path (e.g. "fmt")
	// fn.Name() gives the function name (e.g. "Sprintf")
	var extFuncInfo *ExternalFuncInfo
	if fn.Pkg != nil {
		pkgPath := fn.Pkg.Pkg.Path()
		extPkg := importer.GetPackageByPath(pkgPath)
		if extPkg != nil {
			if obj, ok := extPkg.Objects[fn.Name()]; ok {
				extFuncInfo = &ExternalFuncInfo{
					Func:       obj.Value,
					DirectCall: obj.DirectCall,
				}
			}
		}
	}

	if extFuncInfo == nil {
		// Fallback: store the SSA function itself
		extFuncInfo = &ExternalFuncInfo{
			Func:       fn,
			DirectCall: nil,
		}
	}

	funcIdx := c.addConstant(extFuncInfo)
	numArgs := len(i.Call.Args)
	// Manually encode OpCallExternal: [opcode(1)] [funcIdx(2)] [numArgs(1)]
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(OpCallExternal),
		byte(funcIdx>>8), byte(funcIdx),
		byte(numArgs))
	c.emit(OpSetLocal, uint16(resultIdx))
}

// compileIndirectCall compiles an indirect call (closure or function value).
// The callee is determined at runtime from a variable or closure.
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

// compileValue compiles an SSA value to push it onto the stack.
// It handles:
//   - Constants (literal values)
//   - Function references (for closures)
//   - Phi nodes (reads from pre-allocated slot)
//   - Free variables (closure captures)
//   - Local variables (reads from slot)
//   - Global variables (pushes global address)
func (c *Compiler) compileValue(v ssa.Value) {
	switch val := v.(type) {
	case *ssa.Const:
		// Handle constants
		c.compileConst(val)
	case *ssa.Function:
		// Handle function references — push as a closure with 0 free vars
		if fnIdx, ok := c.funcIndex[val]; ok {
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(OpClosure),
				byte(fnIdx>>8), byte(fnIdx),
				byte(0))
		} else {
			c.emit(OpNil)
		}
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
	case *ssa.Global:
		// Global variables - push the global address
		// Globals are represented as pointers in SSA
		globalName := val.Name()
		globalIdx, ok := c.globals[globalName]
		if !ok {
			globalIdx = len(c.globals)
			c.globals[globalName] = globalIdx
		}
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(OpGlobal),
			byte(globalIdx>>8), byte(globalIdx))
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
		switch t.Kind() { //nolint:exhaustive
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
	c.emit(OpFieldAddr, uint16(i.Field))
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

	if i.CommaOk {
		// For comma-ok map lookup, return a tuple (value, ok)
		c.emit(OpIndexOk)
	} else {
		c.emit(OpIndex)
	}
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
		// Check if binding is an Alloc (a local variable that needs reference semantics)
		// For recursive closures, we need to capture the variable slot, not its value
		if alloc, ok := binding.(*ssa.Alloc); ok {
			// Check if this Alloc has a local slot (it's a local variable being captured)
			if slotIdx, ok := c.symbolTable.GetLocal(alloc); ok {
				// Emit OpAddr to push a reference to the slot
				c.emit(OpAddr, uint16(slotIdx))
				continue
			}
		}
		// Default: compile the binding value (for non-Alloc bindings)
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
//
// SSA Select returns an n+2-tuple: (index int, recvOk bool, r_0, ..., r_{n-1})
// where n is the number of RECV states.
//
// The compiler pushes all channel values and send values onto the stack,
// stores a SelectMeta descriptor in the constant pool, and emits OpSelect.
func (c *Compiler) compileSelect(i *ssa.Select) {
	// Count receive states for the result tuple size
	numRecv := 0
	for _, st := range i.States {
		if st.Dir == types.RecvOnly {
			numRecv++
		}
	}

	// Build metadata: direction for each state (true=send, false=recv)
	dirs := make([]bool, len(i.States))
	for idx, st := range i.States {
		dirs[idx] = (st.Dir == types.SendOnly)
	}

	meta := SelectMeta{
		NumStates: len(i.States),
		Blocking:  i.Blocking,
		Dirs:      dirs,
		NumRecv:   numRecv,
	}

	// Push channels and (for sends) the send values onto the stack.
	// Order: for each state, push Chan; if send, also push SendVal.
	for _, st := range i.States {
		c.compileValue(st.Chan)
		if st.Dir == types.SendOnly {
			c.compileValue(st.Send)
		}
	}

	metaIdx := c.addConstant(meta)
	c.emit(OpSelect, metaIdx)

	// Store the result tuple into a local
	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(OpSetLocal, uint16(resultIdx))
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
	typeIdx := c.addType(i.AssertedType)

	c.compileValue(i.X)
	c.emit(OpAssert, uint16(typeIdx))

	// For both CommaOk and non-CommaOk forms, store the result
	// The CommaOk form returns a tuple (value, ok), non-CommaOk returns just value
	// The Extract instruction will get individual values when needed
	resultIdx := c.symbolTable.AllocLocal(i)
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
// It emits OpGoCall to spawn a new goroutine that executes the function call.
// Arguments are evaluated in the current goroutine, then passed to the spawned goroutine.
// Supports both direct function calls (go foo()) and closure calls (go func(){}()).
func (c *Compiler) compileGo(i *ssa.Go) {
	// Push arguments (evaluated in current goroutine)
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	// Check if it's a direct function call
	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		// Emit OpGoCall to spawn goroutine
		// Manually encode with proper operand sizes: [opcode(1)] [funcIdx(2)] [numArgs(1)]
		funcIdx := c.funcIndex[fn]
		numArgs := len(i.Call.Args)
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(OpGoCall),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		return
	}

	// It's a closure or indirect call
	// Push the closure value onto the stack
	c.compileValue(i.Call.Value)

	// Emit OpGoCallIndirect to spawn goroutine with closure
	// Manually encode: [opcode(1)] [numArgs(1)]
	numArgs := len(i.Call.Args)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(OpGoCallIndirect),
		byte(numArgs))
}

// emit appends an opcode and its operands to the current function's bytecode.
// Operand width is determined by OperandWidths map.
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

// emitJump emits an unconditional jump instruction.
// The target offset is a placeholder that will be patched by patchJumps.
func (c *Compiler) emitJump(target *ssa.BasicBlock) {
	offset := len(c.currentFunc.Instructions)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(OpJump), 0, 0)
	c.jumps = append(c.jumps, jumpInfo{offset: offset, targetBlock: target})
}

// emitJumpFalse emits a conditional jump that executes if the top of stack is false.
// The target offset is a placeholder that will be patched by patchJumps.
func (c *Compiler) emitJumpFalse(target *ssa.BasicBlock) {
	offset := len(c.currentFunc.Instructions)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(OpJumpFalse), 0, 0)
	c.jumps = append(c.jumps, jumpInfo{offset: offset, targetBlock: target})
}

// patchJumps resolves jump targets with actual bytecode offsets.
// This is called after all basic blocks have been compiled,
// when the offset of each block is known.
func (c *Compiler) patchJumps(blockOffsets map[*ssa.BasicBlock]int) {
	for _, jump := range c.jumps {
		targetOffset := blockOffsets[jump.targetBlock]
		// The operand is at offset+1 (2 bytes)
		c.currentFunc.Instructions[jump.offset+1] = byte(targetOffset >> 8)
		c.currentFunc.Instructions[jump.offset+2] = byte(targetOffset)
	}
}

// addConstant adds a value to the constant pool and returns its index.
// Constants include literals, function references, and external objects.
func (c *Compiler) addConstant(val any) uint16 {
	idx := len(c.constants)
	c.constants = append(c.constants, val)
	return uint16(idx)
}

// addType adds a types.Type to the type pool and returns its index.
// Types are used for type assertions, conversions, and allocations.
func (c *Compiler) addType(t types.Type) uint16 {
	idx := len(c.types)
	c.types = append(c.types, t)
	return uint16(idx)
}

// reversePostorder returns basic blocks in reverse postorder.
// This ordering ensures that each block is processed before its successors,
// which is useful for optimization and for consistent jump targets.
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
