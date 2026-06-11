package compiler

import (
	"strings"

	"golang.org/x/tools/go/ssa"
)

func (c *compiler) externalFuncOriginsFromSyntheticMethodWrapper(fn *ssa.Function, seen map[ssa.Value]bool) []externalFuncOrigin {
	// Method values and method expressions use SSA-generated $bound/$thunk
	// wrappers. Boundary validation should use the wrapper signature at the call
	// site, while the package/name still comes from the wrapped external method.
	if !isSyntheticMethodWrapper(fn) {
		return nil
	}

	var origins []externalFuncOrigin
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			origins = appendExternalFuncOrigins(origins, c.externalFuncOriginsFromSyntheticInstruction(fn, instr, seen))
		}
	}
	return origins
}

func (c *compiler) externalFuncOriginsFromSyntheticInstruction(fn *ssa.Function, instr ssa.Instruction, seen map[ssa.Value]bool) []externalFuncOrigin {
	call, ok := instr.(ssa.CallInstruction)
	if !ok {
		return nil
	}
	common := call.Common()
	if common == nil {
		return nil
	}
	callValue := common.Value
	if static := common.StaticCallee(); static != nil {
		callValue = static
	}

	var origins []externalFuncOrigin
	for _, origin := range c.externalFuncOriginsSeen(callValue, seen) {
		if fn.Signature != nil {
			origin.Sig = fn.Signature
		}
		origins = appendExternalFuncOrigins(origins, []externalFuncOrigin{origin})
	}
	return origins
}

func isSyntheticMethodWrapper(fn *ssa.Function) bool {
	if fn == nil || fn.Synthetic == "" {
		return false
	}
	return strings.HasSuffix(fn.Name(), "$bound") ||
		strings.HasSuffix(fn.Name(), "$thunk") ||
		strings.Contains(fn.Synthetic, "method wrapper")
}
