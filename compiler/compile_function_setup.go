package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func (c *compiler) beginFunction(fn *ssa.Function) {
	// Function compilation state is intentionally reset per SSA function.
	// The compiler object owns cross-function pools, but jumps, phi slots, and
	// locals are only meaningful for the currently emitted bytecode function.
	c.currentFunc = newCompiledFunction(fn, c.funcIndex[fn])
	c.symbolTable = NewSymbolTable()
	c.jumps = nil
	c.phiSlots = make(map[*ssa.Phi]int)
	c.allocateFunctionLocals(fn)
	c.currentFunc.NumLocals = c.symbolTable.NumLocals()
	c.currentFunc.NumFreeVars = len(fn.FreeVars)
	c.currentFunc.ResultAllocSlots = detectResultAllocSlots(fn, c.symbolTable)
}

func newCompiledFunction(fn *ssa.Function, funcIndex int) *bytecode.CompiledFunction {
	// FuncIdx is assigned from the SSA pointer catalog, not from function name.
	// Methods on different receiver types can share a name, so name-based
	// indexing would corrupt direct calls and method dispatch.
	cf := &bytecode.CompiledFunction{
		Name:         fn.Name(),
		Instructions: make([]byte, 0),
		NumParams:    len(fn.Params),
		ParamTypes:   make([]types.Type, len(fn.Params)),
		FuncIdx:      funcIndex,
	}
	for i, param := range fn.Params {
		cf.ParamTypes[i] = param.Type()
	}
	populateFunctionSignatureMetadata(cf, fn.Signature)
	return cf
}

func populateFunctionSignatureMetadata(cf *bytecode.CompiledFunction, sig *types.Signature) {
	// Receiver metadata is consumed by VM method lookup. Store the short name
	// plus pointer bit here so runtime dispatch can avoid repeated type parsing.
	cf.IsVariadic = sig.Variadic()
	if cf.IsVariadic && sig.Params().Len() > 0 {
		cf.VariadicParamType = sig.Params().At(sig.Params().Len() - 1).Type()
	}
	if sig.Recv() == nil {
		return
	}
	cf.HasReceiver = true
	cf.ReceiverTypeName = extractReceiverShortName(sig.Recv().Type())
	_, cf.ReceiverIsPointer = sig.Recv().Type().(*types.Pointer)
}

func (c *compiler) allocateFunctionLocals(fn *ssa.Function) {
	// Parameters occupy stable local slots first because call setup pushes them
	// in parameter order. Free variables are tracked separately: closure reads
	// address them via the free-var table, not via normal locals.
	for _, param := range fn.Params {
		c.symbolTable.AllocLocal(param)
	}
	for i, freeVar := range fn.FreeVars {
		c.symbolTable.freeVars[freeVar] = i
	}
	for _, block := range fn.Blocks {
		c.allocateInstructionLocals(block)
	}
}

func (c *compiler) allocateInstructionLocals(block *ssa.BasicBlock) {
	// Allocate every value-producing instruction before emitting bytecode.
	// This lets forward references, phi moves, and later optimized instruction
	// sequences all agree on a single local-slot numbering.
	for _, instr := range block.Instrs {
		switch instr := instr.(type) {
		case *ssa.Phi:
			slot := c.symbolTable.AllocLocal(instr)
			c.phiSlots[instr] = slot
		case *ssa.Alloc:
			c.symbolTable.AllocLocal(instr)
		case ssa.Value:
			c.symbolTable.AllocLocal(instr)
		}
	}
}
