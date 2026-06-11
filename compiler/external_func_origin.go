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
	if !markExternalFuncOriginSeen(v, seen) {
		return nil
	}
	if origins, ok := c.externalFuncOriginsFromFunctionValue(v, seen); ok {
		return origins
	}
	return c.externalFuncOriginsFromSSAValue(v, seen)
}

func markExternalFuncOriginSeen(v ssa.Value, seen map[ssa.Value]bool) bool {
	if v == nil || seen[v] {
		return false
	}
	seen[v] = true
	return true
}

func (c *compiler) externalFuncOriginsFromFunctionValue(v ssa.Value, seen map[ssa.Value]bool) ([]externalFuncOrigin, bool) {
	if fn, ok := v.(*ssa.Function); ok {
		if origin, ok := c.externalFuncOriginForFunction(fn); ok {
			return []externalFuncOrigin{origin}, true
		}
		return c.externalFuncOriginsFromSyntheticMethodWrapper(fn, seen), true
	}
	return nil, false
}

func (c *compiler) externalFuncOriginsFromSSAValue(v ssa.Value, seen map[ssa.Value]bool) []externalFuncOrigin {
	// Route by SSA shape first; each helper owns one origin-flow domain.
	if origins, ok := c.externalFuncOriginsThroughSSAValueAlias(v, seen); ok {
		return origins
	}
	if origins, ok := c.externalFuncOriginsThroughSSAContainer(v, seen); ok {
		return origins
	}
	if origins, ok := c.externalFuncOriginsFromSSAProducer(v, seen); ok {
		return origins
	}
	return c.externalFuncOriginsStoredIn(v, seen)
}

// externalFuncOriginsThroughSSAValueAlias follows conversions that preserve the callable value.
func (c *compiler) externalFuncOriginsThroughSSAValueAlias(v ssa.Value, seen map[ssa.Value]bool) ([]externalFuncOrigin, bool) {
	switch val := v.(type) {
	case *ssa.ChangeInterface:
		return c.externalFuncOriginsThroughValueAlias(val.X, seen), true
	case *ssa.ChangeType:
		return c.externalFuncOriginsThroughValueAlias(val.X, seen), true
	case *ssa.Convert:
		return c.externalFuncOriginsThroughValueAlias(val.X, seen), true
	case *ssa.MakeInterface:
		return c.externalFuncOriginsThroughValueAlias(val.X, seen), true
	default:
		return nil, false
	}
}

// externalFuncOriginsThroughSSAContainer follows projections that retrieve a callable value from aggregate storage.
func (c *compiler) externalFuncOriginsThroughSSAContainer(v ssa.Value, seen map[ssa.Value]bool) ([]externalFuncOrigin, bool) {
	switch val := v.(type) {
	case *ssa.Slice:
		return c.externalFuncOriginsThroughContainer(val.X, seen), true
	case *ssa.Index:
		return c.externalFuncOriginsThroughContainer(val.X, seen), true
	case *ssa.Lookup:
		return c.externalFuncOriginsThroughContainer(val.X, seen), true
	case *ssa.IndexAddr:
		return c.externalFuncOriginsThroughContainer(val.X, seen), true
	case *ssa.FieldAddr:
		return c.externalFuncOriginsThroughContainer(val.X, seen), true
	case *ssa.UnOp:
		return c.externalFuncOriginsThroughContainer(val.X, seen), true
	case *ssa.Field:
		return c.externalFuncOriginsThroughContainer(val.X, seen), true
	default:
		return nil, false
	}
}

// externalFuncOriginsFromSSAProducer handles SSA values that directly produce callable values.
func (c *compiler) externalFuncOriginsFromSSAProducer(v ssa.Value, seen map[ssa.Value]bool) ([]externalFuncOrigin, bool) {
	switch val := v.(type) {
	case *ssa.MakeClosure:
		return c.externalFuncOriginsSeen(val.Fn, seen), true
	case *ssa.Call:
		return c.externalFuncOriginsFromCallResult(val, 0, seen), true
	case *ssa.Phi:
		return c.externalFuncOriginsFromPhi(val, seen), true
	case *ssa.Extract:
		if call, ok := val.Tuple.(*ssa.Call); ok {
			return c.externalFuncOriginsFromCallResult(call, val.Index, seen), true
		}
		return c.externalFuncOriginsSeen(val.Tuple, seen), true
	default:
		return nil, false
	}
}

func (c *compiler) externalFuncOriginsThroughValueAlias(v ssa.Value, seen map[ssa.Value]bool) []externalFuncOrigin {
	return c.externalFuncOriginsSeen(v, seen)
}
