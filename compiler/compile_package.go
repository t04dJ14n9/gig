package compiler

import (
	"fmt"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func (c *compiler) indexFunctions(functions []*ssa.Function) {
	// Use the collected slice order as the single source of truth for function
	// indices. All later direct calls and closure ops refer back to this map.
	for idx, fn := range functions {
		c.funcIndex[fn] = idx
	}
}

func (c *compiler) compilePackageFunctions(functions []*ssa.Function) error {
	c.program.FuncByIndex = make([]*bytecode.CompiledFunction, len(functions))
	for _, fn := range functions {
		compiled, err := c.compileFunction(fn)
		if err != nil {
			return fmt.Errorf("compile function %s: %w", fn.Name(), err)
		}
		c.registerCompiledFunction(fn, compiled)
	}
	return nil
}

func (c *compiler) registerCompiledFunction(fn *ssa.Function, compiled *bytecode.CompiledFunction) {
	// Keep both lookup paths populated: Functions supports legacy name lookup,
	// while FuncByIndex is the collision-safe path for calls between compiled
	// functions with the same short name.
	c.funcs[fn.Name()] = compiled
	c.program.Functions[fn.Name()] = compiled
	idx := c.funcIndex[fn]
	c.program.FuncByIndex[idx] = compiled
}
