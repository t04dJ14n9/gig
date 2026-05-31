package compiler

import (
	"go/types"

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
	// This tracer follows SSA value flow from the callable value used at a call
	// site back to the concrete external function. The seen map is shared across
	// recursive branches so phi cycles, container aliases, and wrapper calls do
	// not loop forever.
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
		return c.externalFuncOriginsFromPhi(val, seen)
	case *ssa.Slice:
		return c.externalFuncOriginsThroughContainer(val.X, seen)
	case *ssa.Index:
		return c.externalFuncOriginsThroughContainer(val.X, seen)
	case *ssa.Lookup:
		return c.externalFuncOriginsThroughContainer(val.X, seen)
	case *ssa.IndexAddr:
		return c.externalFuncOriginsThroughContainer(val.X, seen)
	case *ssa.FieldAddr:
		return c.externalFuncOriginsThroughContainer(val.X, seen)
	case *ssa.UnOp:
		return c.externalFuncOriginsThroughContainer(val.X, seen)
	case *ssa.Field:
		return c.externalFuncOriginsThroughContainer(val.X, seen)
	case *ssa.Extract:
		if call, ok := val.Tuple.(*ssa.Call); ok {
			return c.externalFuncOriginsFromCallResult(call, val.Index, seen)
		}
		return c.externalFuncOriginsSeen(val.Tuple, seen)
	}

	return c.externalFuncOriginsStoredIn(v, seen)
}
