package compiler

import (
	"errors"
	"fmt"
	"go/types"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (c *compiler) newProgram() *bytecode.CompiledProgram {
	return &bytecode.CompiledProgram{
		Functions:           make(map[string]*bytecode.CompiledFunction),
		Globals:             make(map[string]int),
		Types:               make([]types.Type, 0),
		AllowUnsafeTypePass: c.allowUnsafeTypePass,
	}
}

func (c *compiler) finalizeProgram() {
	// The compiler accumulates pools and global metadata incrementally while
	// emitting functions. Finalization freezes those mutable compiler fields into
	// the immutable runtime program shape expected by the VM.
	c.program.Constants = c.constants
	c.program.Types = c.types
	c.program.Globals = c.globals
	c.program.GlobalZeroValues = c.globalZeroValues
	c.program.GlobalElemTypes = c.globalElemTypes
	c.program.ExternalVarValues = c.externalVarValues
	c.program.TypeResolver = c.lookup
	c.program.MethodsByName = buildMethodsByName(c.program.FuncByIndex)
	c.program.PrebakedConstants = prebakeConstants(c.constants)
	c.program.IntConstants = buildIntConstants(c.constants)
}

func buildMethodsByName(functions []*bytecode.CompiledFunction) map[string][]*bytecode.CompiledFunction {
	// Runtime method dispatch narrows candidates by method name first, then by
	// receiver metadata. This avoids repeatedly scanning the full function table.
	methodsByName := make(map[string][]*bytecode.CompiledFunction)
	for _, fn := range functions {
		if fn != nil && fn.HasReceiver {
			methodsByName[fn.Name] = append(methodsByName[fn.Name], fn)
		}
	}
	return methodsByName
}

func prebakeConstants(constants []any) []value.Value {
	// OpConst is hot in the VM. Pre-wrapping constants avoids allocating or
	// reflecting through value.FromInterface for every constant load.
	prebaked := make([]value.Value, len(constants))
	for i, k := range constants {
		prebaked[i] = value.FromInterface(k)
	}
	return prebaked
}

func buildIntConstants(constants []any) []int64 {
	// Int superinstructions read a compact int64 pool. Non-integer constants
	// intentionally keep the zero value because their matching const-is-int bit
	// will be false.
	intConstants := make([]int64, len(constants))
	for i, k := range constants {
		switch v := k.(type) {
		case int:
			intConstants[i] = int64(v)
		case int8:
			intConstants[i] = int64(v)
		case int16:
			intConstants[i] = int64(v)
		case int32:
			intConstants[i] = int64(v)
		case int64:
			intConstants[i] = v
		}
	}
	return intConstants
}

func (c *compiler) compilationError() error {
	if len(c.errors) == 0 {
		return nil
	}
	return fmt.Errorf("compilation errors:\n%w", errors.Join(c.errors...))
}
