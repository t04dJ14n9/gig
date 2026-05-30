package compiler

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/ssa"
)

type externalFuncOrigin struct {
	PkgPath  string
	FuncName string
	Sig      *types.Signature
}

func (c *compiler) externalFuncOrigins(v ssa.Value) []externalFuncOrigin {
	return c.externalFuncOriginsSeen(v, make(map[ssa.Value]bool))
}

func (c *compiler) externalFuncOriginsSeen(v ssa.Value, seen map[ssa.Value]bool) []externalFuncOrigin {
	if v == nil || seen[v] {
		return nil
	}
	seen[v] = true

	if fn, ok := v.(*ssa.Function); ok {
		if origin, ok := c.externalFuncOriginForFunction(fn); ok {
			return []externalFuncOrigin{origin}
		}
		return c.externalFuncOriginsFromSyntheticMethodWrapper(fn, seen)
	}

	switch val := v.(type) {
	case *ssa.ChangeInterface:
		return c.externalFuncOriginsSeen(val.X, seen)
	case *ssa.ChangeType:
		return c.externalFuncOriginsSeen(val.X, seen)
	case *ssa.Convert:
		return c.externalFuncOriginsSeen(val.X, seen)
	case *ssa.MakeInterface:
		return c.externalFuncOriginsSeen(val.X, seen)
	case *ssa.MakeClosure:
		return c.externalFuncOriginsSeen(val.Fn, seen)
	case *ssa.Call:
		return c.externalFuncOriginsFromCallResult(val, 0, seen)
	case *ssa.Phi:
		var origins []externalFuncOrigin
		for _, edge := range val.Edges {
			origins = appendExternalFuncOrigins(origins, c.externalFuncOriginsSeen(edge, seen))
		}
		return origins
	case *ssa.Slice:
		return appendExternalFuncOrigins(
			c.externalFuncOriginsSeen(val.X, seen),
			c.externalFuncOriginsStoredIn(val.X, seen),
		)
	case *ssa.Index:
		return appendExternalFuncOrigins(
			c.externalFuncOriginsSeen(val.X, seen),
			c.externalFuncOriginsStoredIn(val.X, seen),
		)
	case *ssa.Lookup:
		return appendExternalFuncOrigins(
			c.externalFuncOriginsSeen(val.X, seen),
			c.externalFuncOriginsStoredIn(val.X, seen),
		)
	case *ssa.IndexAddr:
		return appendExternalFuncOrigins(
			c.externalFuncOriginsSeen(val.X, seen),
			c.externalFuncOriginsStoredIn(val.X, seen),
		)
	case *ssa.FieldAddr:
		return appendExternalFuncOrigins(
			c.externalFuncOriginsSeen(val.X, seen),
			c.externalFuncOriginsStoredIn(val.X, seen),
		)
	case *ssa.UnOp:
		return appendExternalFuncOrigins(
			c.externalFuncOriginsSeen(val.X, seen),
			c.externalFuncOriginsStoredIn(val.X, seen),
		)
	case *ssa.Field:
		return appendExternalFuncOrigins(
			c.externalFuncOriginsSeen(val.X, seen),
			c.externalFuncOriginsStoredIn(val.X, seen),
		)
	case *ssa.Extract:
		if call, ok := val.Tuple.(*ssa.Call); ok {
			return c.externalFuncOriginsFromCallResult(call, val.Index, seen)
		}
		return c.externalFuncOriginsSeen(val.Tuple, seen)
	}

	return c.externalFuncOriginsStoredIn(v, seen)
}

func (c *compiler) externalFuncOriginsFromCallResult(call *ssa.Call, resultIndex int, seen map[ssa.Value]bool) []externalFuncOrigin {
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

func (c *compiler) externalFuncOriginsFromSyntheticMethodWrapper(fn *ssa.Function, seen map[ssa.Value]bool) []externalFuncOrigin {
	if !isSyntheticMethodWrapper(fn) {
		return nil
	}

	var origins []externalFuncOrigin
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			call, ok := instr.(ssa.CallInstruction)
			if !ok {
				continue
			}
			common := call.Common()
			if common == nil {
				continue
			}
			callValue := common.Value
			if static := common.StaticCallee(); static != nil {
				callValue = static
			}
			for _, origin := range c.externalFuncOriginsSeen(callValue, seen) {
				if fn.Signature != nil {
					origin.Sig = fn.Signature
				}
				origins = appendExternalFuncOrigins(origins, []externalFuncOrigin{origin})
			}
		}
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

func (c *compiler) externalFuncOriginForFunction(fn *ssa.Function) (externalFuncOrigin, bool) {
	if fn == nil || fn.Signature == nil {
		return externalFuncOrigin{}, false
	}
	if _, known := c.funcIndex[fn]; known {
		return externalFuncOrigin{}, false
	}

	if fn.Signature.Recv() != nil {
		pkgPath := methodOwnerPkgPath(fn)
		if pkgPath == "" || pkgPath == "main" || pkgPath == "command-line-arguments" {
			return externalFuncOrigin{}, false
		}
		return externalFuncOrigin{
			PkgPath:  pkgPath,
			FuncName: extractMethodName(fn.Name()),
			Sig:      fn.Signature,
		}, true
	}

	if fn.Pkg == nil || fn.Pkg.Pkg == nil {
		return externalFuncOrigin{}, false
	}
	pkgPath := fn.Pkg.Pkg.Path()
	if pkgPath == "" || pkgPath == "main" || pkgPath == "command-line-arguments" {
		return externalFuncOrigin{}, false
	}
	return externalFuncOrigin{
		PkgPath:  pkgPath,
		FuncName: fn.Name(),
		Sig:      fn.Signature,
	}, true
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
		switch instr := ref.(type) {
		case *ssa.Store:
			if storeTargetsValue(instr.Addr, container) {
				origins = appendExternalFuncOrigins(origins, c.externalFuncOriginsSeen(instr.Val, seen))
			}
		case *ssa.MapUpdate:
			if instr.Map == container {
				origins = appendExternalFuncOrigins(origins, c.externalFuncOriginsSeen(instr.Value, seen))
			}
		case *ssa.IndexAddr:
			if instr.X == container {
				origins = appendExternalFuncOrigins(origins, c.externalFuncOriginsStoredIn(instr, seen))
			}
		case *ssa.FieldAddr:
			if instr.X == container {
				origins = appendExternalFuncOrigins(origins, c.externalFuncOriginsStoredIn(instr, seen))
			}
		}
	}
	return origins
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
		duplicate := false
		for _, existing := range dst {
			if existing.PkgPath == candidate.PkgPath && existing.FuncName == candidate.FuncName {
				duplicate = true
				break
			}
		}
		if !duplicate {
			dst = append(dst, candidate)
		}
	}
	return dst
}
