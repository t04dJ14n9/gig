// compile_instr.go routes SSA instructions to focused instruction-lowering helpers.
package compiler

import (
	"golang.org/x/tools/go/ssa"
)

// compileInstruction compiles a single SSA instruction to bytecode.
func (c *compiler) compileInstruction(instr ssa.Instruction) {
	if c.compileExpressionInstruction(instr) {
		return
	}
	if c.compileAddressingInstruction(instr) {
		return
	}
	if c.compileConstructionInstruction(instr) {
		return
	}
	if c.compileControlInstruction(instr) {
		return
	}
	if c.compileEffectInstruction(instr) {
		return
	}
	c.compileIgnoredInstruction(instr)
}
