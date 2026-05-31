// compile_func.go handles the top-level per-function SSA→bytecode workflow.
package compiler

import (
	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/compiler/optimize"
	"github.com/t04dJ14n9/gig/model/bytecode"
)

// compileFunction compiles a single SSA function to bytecode.
func (c *compiler) compileFunction(fn *ssa.Function) (*bytecode.CompiledFunction, error) {
	// Keep the order here stable:
	// 1. locals are allocated before local type maps are built;
	// 2. blocks are emitted before constant specialization maps are built;
	// 3. optimization runs after jump patching so instruction offsets are final.
	c.beginFunction(fn)
	localTypes := buildLocalTypeMaps(c.symbolTable.locals)
	c.compileBlocks(fn)
	constIsInt := buildConstIntMap(c.constants)
	c.currentFunc.Instructions, c.currentFunc.HasIntLocals = optimize.Optimize(
		c.currentFunc.Instructions,
		localTypes.IsInt,
		constIsInt,
		localTypes.IsIntSlice,
	)

	return c.currentFunc, nil
}
