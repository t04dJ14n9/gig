package compiler

import "golang.org/x/tools/go/ssa"

func (c *compiler) externalFuncOriginsFromCallResult(call *ssa.Call, resultIndex int, seen map[ssa.Value]bool) []externalFuncOrigin {
	// Some callables are returned from helper functions before they are invoked.
	// Follow the selected return result back through the callee body so boundary
	// validation still sees the original external package function.
	if call == nil || resultIndex < 0 {
		return nil
	}
	callee := call.Call.StaticCallee()
	if callee == nil {
		callee, _ = call.Call.Value.(*ssa.Function)
	}
	if callee == nil || callee.Blocks == nil || seen[callee] {
		return nil
	}
	seen[callee] = true

	var origins []externalFuncOrigin
	for _, block := range callee.Blocks {
		for _, instr := range block.Instrs {
			ret, ok := instr.(*ssa.Return)
			if !ok || resultIndex >= len(ret.Results) {
				continue
			}
			origins = appendExternalFuncOrigins(origins, c.externalFuncOriginsFromReturnValue(ret.Results[resultIndex], callee, call, seen))
		}
	}
	return origins
}

func (c *compiler) externalFuncOriginsFromReturnValue(v ssa.Value, callee *ssa.Function, call *ssa.Call, seen map[ssa.Value]bool) []externalFuncOrigin {
	// If the helper returns one of its parameters, remap that parameter to the
	// original call argument. This preserves origin tracking through identity
	// wrappers like func passthrough(f any) any { return f }.
	if param, ok := v.(*ssa.Parameter); ok && callee != nil && call != nil {
		for idx, candidate := range callee.Params {
			if candidate == param {
				if idx < len(call.Call.Args) {
					return c.externalFuncOriginsSeen(call.Call.Args[idx], seen)
				}
				return nil
			}
		}
	}
	return c.externalFuncOriginsSeen(v, seen)
}
