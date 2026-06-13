// frame.go is the per-call execution record and dispatcher loop. It
// mirrors gofun's frame model: walk the SSA basic blocks one
// instruction at a time, branching on the SSA node type, with Phi
// nodes resolved at block entry.
package interp

import (
	"fmt"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/value"
)

// continuation is the next-action signal returned by every per-instruction
// runner. It mirrors gofun's _NEXT/_JUMP/_RETURN tri-state.
type continuation int

const (
	contNext   continuation = iota // continue with the next instruction
	contJump                       // block changed; restart the outer loop
	contReturn                     // function is returning
)

// frame is the per-call activation record. It carries the function's
// SSA blocks, the cell map for SSA-value → Cell, free-var cells for
// closures, and an iterator side-channel for ssa.Range/Next pairs.
type frame struct {
	fn        *ssa.Function
	block     *ssa.BasicBlock
	prevBlock *ssa.BasicBlock
	cells     map[ssa.Value]*Cell
	freeVars  []*Cell
	iters     map[ssa.Value]*rangeIter

	// defer / panic / recover state.
	defers    []*deferRecord
	panicking bool
	panicVal  any
}

// callSSA invokes an SSA function with the given args. Returns the
// function's result tuple (zero, one, or many). depth is the current
// call depth, bumped on every entry to catch runaway recursion.
//
// caller is used for diagnostics and recover() once defer/panic land.
// freeVars is for closures (Phase 6.3); pass nil for plain functions.
func (p *program) callSSA(caller *frame, fn *ssa.Function, args []value.Value, freeVars []*Cell, depth int) (results []value.Value, err error) {
	if depth >= p.maxDepth {
		return nil, fmt.Errorf("interp: max call depth %d exceeded calling %s", p.maxDepth, fn.Name())
	}
	if len(fn.Blocks) == 0 {
		return nil, fmt.Errorf("interp: function %s has no body", fn.Name())
	}

	fr := &frame{
		fn:       fn,
		block:    fn.Blocks[0],
		cells:    make(map[ssa.Value]*Cell, 16),
		freeVars: freeVars,
	}

	// Bind parameters.
	for i, param := range fn.Params {
		fr.cells[param] = &Cell{Name: param.Name(), Type: param.Type(), Value: args[i]}
	}

	// Bind free variables (closures). Empty for plain functions.
	for i, fv := range fn.FreeVars {
		if i >= len(freeVars) {
			break
		}
		fr.cells[fv] = freeVars[i]
	}

	// Pre-allocate Cells for every Local. Locals are pointer-typed in
	// SSA and the interpreter models them as addressable
	// reflect.Values: see runAlloc for the same treatment of heap
	// Allocs.
	for _, local := range fn.Locals {
		ptr := derefSSAType(local.Type())
		addr, err := p.makeAddressable(ptr)
		if err != nil {
			return nil, fmt.Errorf("interp: %s: alloc local %s: %w", fn.Name(), local.Name(), err)
		}
		fr.cells[local] = &Cell{
			Name:  local.Name(),
			Type:  local.Type(),
			Value: reflectValue(addr.Addr()),
		}
	}

	// Install a panic handler so deferred functions can run and
	// recover() can take effect. If the panic isn't recovered here,
	// re-panic so an outer interpreted callSSA can run its own defers
	// and (potentially) consume the panic. Only the top-level Call
	// surfaces the panic as a returned error; intermediate frames
	// must propagate it as a panic so chained recover() works.
	defer func() {
		if re := recover(); re != nil {
			fr.panicking = true
			fr.panicVal = re
			prev := p.panicFrame
			p.panicFrame = fr
			for i := len(fr.defers) - 1; i >= 0; i-- {
				_ = p.runDeferRec(fr, fr.defers[i])
			}
			p.panicFrame = prev
			fr.defers = nil
			if fr.panicking {
				// Not recovered: propagate as a panic so the caller
				// frame can engage its own defers/recover.
				panic(fr.panicVal)
			}
			// Recovered: jump into the function's Recover block (if
			// SSA emitted one) so any named-return reads land in the
			// caller's results.
			if fr.fn.Recover != nil {
				fr.block = fr.fn.Recover
				fr.prevBlock = nil
				rs, rerr := p.runFrame(caller, fr, depth)
				if rerr != nil {
					err = rerr
					results = nil
					return
				}
				results = rs
				err = nil
				return
			}
			results, _ = p.zeroResultsFor(fn)
			err = nil
		}
	}()

	results, err = p.runFrame(caller, fr, depth)
	return results, err
}

// zeroResultsFor returns the function's zero result tuple (one
// value.Value per Results entry). Used by the panic-recover path
// when a deferred recover() consumes a panic and the function would
// otherwise leak nil results to the caller.
func (p *program) zeroResultsFor(fn *ssa.Function) ([]value.Value, error) {
	sig := fn.Signature
	if sig == nil {
		return nil, nil
	}
	res := sig.Results()
	if res == nil || res.Len() == 0 {
		return nil, nil
	}
	out := make([]value.Value, res.Len())
	for i := 0; i < res.Len(); i++ {
		v, err := p.converter.Zero(res.At(i).Type(), p.resolver)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

// runFrame is the dispatch loop. It walks blocks until a Return is
// hit or an error escapes.
func (p *program) runFrame(caller *frame, fr *frame, depth int) ([]value.Value, error) {
	for fr.block != nil {
		// Phi nodes are read at block entry, BEFORE any other instruction
		// of the block runs. The semantic is "pick the edge from
		// prevBlock". We compute all Phis from a snapshot of the
		// current cell map so simultaneous updates don't see each
		// other.
		if err := p.runBlockPhis(fr); err != nil {
			return nil, err
		}

		// Step through the rest of the instructions.
		var ret []value.Value
		var contState continuation
		var err error
	instrs:
		for _, instr := range fr.block.Instrs {
			if _, isPhi := instr.(*ssa.Phi); isPhi {
				continue // already handled
			}
			contState, ret, err = p.visitInstr(caller, fr, instr, depth)
			if err != nil {
				return nil, err
			}
			switch contState {
			case contNext:
				continue
			case contJump:
				break instrs // restart outer loop with new fr.block
			case contReturn:
				return ret, nil
			}
		}
	}
	return nil, fmt.Errorf("interp: %s: ran off the end of a block", fr.fn.Name())
}

// runBlockPhis evaluates every Phi at the start of the current block,
// using the prevBlock to choose which edge applies. All Phis are
// computed from a snapshot first, then assigned together — that
// matches Go SSA semantics where Phis run "simultaneously".
func (p *program) runBlockPhis(fr *frame) error {
	type pending struct {
		instr *ssa.Phi
		val   value.Value
	}
	var todo []pending
	for _, instr := range fr.block.Instrs {
		phi, ok := instr.(*ssa.Phi)
		if !ok {
			break // Phis are always at block start
		}
		// Find which predecessor edge we came in from.
		var picked ssa.Value
		for i, pred := range fr.block.Preds {
			if pred == fr.prevBlock {
				picked = phi.Edges[i]
				break
			}
		}
		if picked == nil {
			// Entry block has no Phis to resolve in practice; if we hit
			// this, the SSA is malformed.
			return fmt.Errorf("interp: %s: phi at block %d has no matching predecessor edge", fr.fn.Name(), fr.block.Index)
		}
		v, err := p.readValue(fr, picked)
		if err != nil {
			return err
		}
		todo = append(todo, pending{phi, v})
	}
	for _, t := range todo {
		fr.cells[t.instr] = &Cell{Name: t.instr.Name(), Type: t.instr.Type(), Value: t.val}
	}
	return nil
}
