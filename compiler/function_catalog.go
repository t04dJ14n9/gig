package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"
)

func collectPackageFunctions(mainPkg *ssa.Package) []*ssa.Function {
	// Function order is the source of bytecode function indices. The collector
	// intentionally walks package functions, receiver methods, and synthetic
	// wrappers before indices are assigned so every internal call can use a
	// stable direct index.
	collector := newFunctionCollector()
	collector.collectMemberFunctions(mainPkg)
	collector.collectMethodFunctions(mainPkg)
	collector.collectSyntheticWrappers()
	return collector.functions
}

type functionCollector struct {
	functions []*ssa.Function
	seen      map[*ssa.Function]bool
}

func newFunctionCollector() *functionCollector {
	return &functionCollector{
		seen: make(map[*ssa.Function]bool),
	}
}

func (c *functionCollector) collectFunction(fn *ssa.Function) {
	// Anonymous functions are nested under their parent in SSA. Add them at the
	// point their parent is discovered so closure calls are compiled with the
	// rest of the package-local function graph.
	if c.seen[fn] {
		return
	}
	c.seen[fn] = true
	c.functions = append(c.functions, fn)
	for _, anon := range fn.AnonFuncs {
		c.collectFunction(anon)
	}
}

func (c *functionCollector) collectMemberFunctions(mainPkg *ssa.Package) {
	for _, member := range mainPkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			c.collectFunction(fn)
		}
	}
}

func (c *functionCollector) collectMethodFunctions(mainPkg *ssa.Package) {
	// SSA type members do not directly contain methods. Query both value and
	// pointer receiver method sets so compiled method dispatch can resolve
	// methods regardless of receiver shape.
	for _, member := range mainPkg.Members {
		t, ok := member.(*ssa.Type)
		if !ok {
			continue
		}
		for _, recv := range []types.Type{t.Type(), types.NewPointer(t.Type())} {
			c.collectMethodsForReceiver(mainPkg, recv)
		}
	}
}

func (c *functionCollector) collectMethodsForReceiver(mainPkg *ssa.Package, recv types.Type) {
	mset := mainPkg.Prog.MethodSets.MethodSet(recv)
	for i := 0; i < mset.Len(); i++ {
		fn := mainPkg.Prog.MethodValue(mset.At(i))
		if fn != nil && fn.Package() == mainPkg {
			c.collectFunction(fn)
		}
	}
}

func (c *functionCollector) collectSyntheticWrappers() {
	// Method values and method expressions can reference synthetic $bound and
	// $thunk functions. They are not package members, so keep scanning until no
	// newly discovered wrapper can reveal another wrapper.
	changed := true
	for changed {
		changed = false
		for _, fn := range c.functions {
			if c.collectSyntheticWrappersFromFunction(fn) {
				changed = true
			}
		}
	}
}

func (c *functionCollector) collectSyntheticWrappersFromFunction(fn *ssa.Function) bool {
	if fn.Blocks == nil {
		return false
	}
	changed := false
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if c.collectSyntheticWrapperFromInstruction(instr) {
				changed = true
			}
		}
	}
	return changed
}

func (c *functionCollector) collectSyntheticWrapperFromInstruction(instr ssa.Instruction) bool {
	if mc, ok := instr.(*ssa.MakeClosure); ok {
		return c.collectSyntheticWrapperValue(mc.Fn)
	}
	if call, ok := instr.(*ssa.Call); ok {
		return c.collectSyntheticWrapperValue(call.Call.Value)
	}
	return false
}

func (c *functionCollector) collectSyntheticWrapperValue(v ssa.Value) bool {
	fn, ok := v.(*ssa.Function)
	if !ok || c.seen[fn] || fn.Blocks == nil {
		return false
	}
	c.collectFunction(fn)
	return true
}
