package compiler

import "golang.org/x/tools/go/ssa"

func (c *compiler) externalFuncOriginsThroughContainer(container ssa.Value, seen map[ssa.Value]bool) []externalFuncOrigin {
	// Container reads can expose either the container value itself or function
	// values stored into it earlier. Follow both paths so slices, maps, fields,
	// and pointer dereferences do not hide third-party callables.
	return appendExternalFuncOrigins(
		c.externalFuncOriginsSeen(container, seen),
		c.externalFuncOriginsStoredIn(container, seen),
	)
}

func (c *compiler) externalFuncOriginsStoredIn(container ssa.Value, seen map[ssa.Value]bool) []externalFuncOrigin {
	if container == nil {
		return nil
	}
	refs := container.Referrers()
	if refs == nil {
		return nil
	}

	var origins []externalFuncOrigin
	for _, ref := range *refs {
		origins = appendExternalFuncOrigins(origins, c.externalFuncOriginsFromContainerRef(container, ref, seen))
	}
	return origins
}

func (c *compiler) externalFuncOriginsFromContainerRef(container ssa.Value, ref ssa.Instruction, seen map[ssa.Value]bool) []externalFuncOrigin {
	switch instr := ref.(type) {
	case *ssa.Store:
		if storeTargetsValue(instr.Addr, container) {
			return c.externalFuncOriginsSeen(instr.Val, seen)
		}
	case *ssa.MapUpdate:
		if instr.Map == container {
			return c.externalFuncOriginsSeen(instr.Value, seen)
		}
	case *ssa.IndexAddr:
		if instr.X == container {
			return c.externalFuncOriginsStoredIn(instr, seen)
		}
	case *ssa.FieldAddr:
		if instr.X == container {
			return c.externalFuncOriginsStoredIn(instr, seen)
		}
	}
	return nil
}

func storeTargetsValue(addr, target ssa.Value) bool {
	if addr == target {
		return true
	}
	switch a := addr.(type) {
	case *ssa.IndexAddr:
		return a.X == target
	case *ssa.FieldAddr:
		return a.X == target
	}
	return false
}

func appendExternalFuncOrigins(dst, src []externalFuncOrigin) []externalFuncOrigin {
	for _, candidate := range src {
		if !containsExternalFuncOrigin(dst, candidate) {
			dst = append(dst, candidate)
		}
	}
	return dst
}

func containsExternalFuncOrigin(origins []externalFuncOrigin, candidate externalFuncOrigin) bool {
	for _, existing := range origins {
		if existing.PkgPath == candidate.PkgPath && existing.FuncName == candidate.FuncName {
			return true
		}
	}
	return false
}
