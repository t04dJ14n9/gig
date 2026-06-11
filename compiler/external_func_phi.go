package compiler

import "golang.org/x/tools/go/ssa"

func (c *compiler) externalFuncOriginsFromPhi(phi *ssa.Phi, seen map[ssa.Value]bool) []externalFuncOrigin {
	// Phi nodes merge values from multiple predecessors. If any edge carries an
	// external callable, the eventual call site must be checked against that
	// origin, while duplicates from equivalent branches should be collapsed.
	var origins []externalFuncOrigin
	for _, edge := range phi.Edges {
		origins = appendExternalFuncOrigins(origins, c.externalFuncOriginsSeen(edge, seen))
	}
	return origins
}
