// closure.go contains MakeClosure and the helper makeFuncValue. Both
// produce reflect-backed callable values that, when invoked, re-enter
// the interpreter with the captured free variables bound to their
// cells.
package interp

import (
	"fmt"
	"reflect"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/value"
)

// makeFuncValue wraps an *ssa.Function as a reflect-callable Value.
// freeVars is non-nil only for closures (MakeClosure); plain function
// references use nil. The reflect.MakeFunc path goes through ToReflect
// for arguments and FromReflect for results, so host code can call
// interpreted functions transparently.
func (p *program) makeFuncValue(fn *ssa.Function, freeVars []*Cell) (value.Value, error) {
	rt, err := p.resolver.ResolveType(fn.Signature)
	if err != nil {
		return value.Value{}, err
	}
	wrapper := reflect.MakeFunc(rt, func(rargs []reflect.Value) []reflect.Value {
		args := make([]value.Value, len(rargs))
		for i, ra := range rargs {
			v, err := p.converter.FromReflect(ra)
			if err != nil {
				panic(fmt.Sprintf("interp: convert closure arg %d: %v", i, err))
			}
			args[i] = v
		}
		results, err := p.callSSA(nil, fn, args, freeVars, 0)
		if err != nil {
			panic(err)
		}
		out := make([]reflect.Value, rt.NumOut())
		for i := range out {
			ot := rt.Out(i)
			if i < len(results) {
				rv, err := p.converter.ToReflect(results[i], ot)
				if err != nil {
					panic(fmt.Sprintf("interp: convert closure result %d: %v", i, err))
				}
				out[i] = rv
			} else {
				out[i] = reflect.Zero(ot)
			}
		}
		return out
	})
	v, err := p.converter.FromReflect(wrapper)
	if err != nil {
		return value.Value{}, err
	}
	return v, nil
}

// runMakeClosure handles ssa.MakeClosure: build a free-vars list from
// the binding cells in the surrounding frame, then wrap the inner
// function so subsequent Call instructions see a callable value.
func (p *program) runMakeClosure(fr *frame, instr *ssa.MakeClosure) (continuation, []value.Value, error) {
	fn, ok := instr.Fn.(*ssa.Function)
	if !ok {
		return contNext, nil, fmt.Errorf("interp: MakeClosure target %T not a function", instr.Fn)
	}
	freeVars := make([]*Cell, len(instr.Bindings))
	for i, b := range instr.Bindings {
		// Each binding is some SSA value in the current frame; we
		// reference the *cell* (not its value) so closures can mutate
		// captured locals.
		if cell, ok := fr.cells[b]; ok {
			freeVars[i] = cell
			continue
		}
		// Bindings can also be globals or constants — read once.
		v, err := p.readValue(fr, b)
		if err != nil {
			return contNext, nil, err
		}
		freeVars[i] = &Cell{Name: b.Name(), Type: b.Type(), Value: v}
	}
	v, err := p.makeFuncValue(fn, freeVars)
	if err != nil {
		return contNext, nil, err
	}
	fr.cells[instr] = &Cell{Name: instr.Name(), Type: instr.Type(), Value: v}
	return contNext, nil, nil
}
