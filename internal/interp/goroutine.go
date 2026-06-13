// goroutine.go implements *ssa.Go, *ssa.Send, *ssa.Select.
// Goroutines spawn the callee in a new goroutine via callSSA. Send
// uses reflect.Value.Send. Select builds reflect.SelectCase entries
// and dispatches via reflect.Select.
package interp

import (
	"fmt"
	"reflect"

	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/value"
)

func (p *program) runGo(fr *frame, instr *ssa.Go) (continuation, []value.Value, error) {
	common := instr.Common()
	args := make([]value.Value, len(common.Args))
	for i, a := range common.Args {
		v, err := p.readValue(fr, a)
		if err != nil {
			return contNext, nil, err
		}
		args[i] = v
	}
	switch tgt := common.Value.(type) {
	case *ssa.Function:
		go func() {
			defer func() { _ = recover() }()
			_, _ = p.callSSA(fr, tgt, args, nil, 0)
		}()
		return contNext, nil, nil
	case *ssa.Builtin:
		go func() {
			defer func() { _ = recover() }()
			_, _ = p.callBuiltinDirect(fr, tgt, args)
		}()
		return contNext, nil, nil
	}
	// Indirect (closure/function value).
	target, err := p.readValue(fr, common.Value)
	if err != nil {
		return contNext, nil, err
	}
	rv, err := p.reflectOf(target, nil)
	if err != nil {
		return contNext, nil, err
	}
	if rv.Kind() != reflect.Func {
		return contNext, nil, fmt.Errorf("interp: go target not callable (kind=%s)", rv.Kind())
	}
	rargs := make([]reflect.Value, len(args))
	conv := value.DefaultConverter()
	for i, a := range args {
		ra, err := conv.ToReflect(a, rv.Type().In(i))
		if err != nil {
			return contNext, nil, err
		}
		rargs[i] = ra
	}
	go func() {
		defer func() { _ = recover() }()
		rv.Call(rargs)
	}()
	return contNext, nil, nil
}

func (p *program) runSend(fr *frame, instr *ssa.Send) (continuation, []value.Value, error) {
	chV, err := p.readValue(fr, instr.Chan)
	if err != nil {
		return contNext, nil, err
	}
	xV, err := p.readValue(fr, instr.X)
	if err != nil {
		return contNext, nil, err
	}
	rv, err := p.reflectOf(chV, nil)
	if err != nil {
		return contNext, nil, err
	}
	for rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	rx, err := p.reflectOf(xV, rv.Type().Elem())
	if err != nil {
		return contNext, nil, err
	}
	rv.Send(rx)
	return contNext, nil, nil
}

func (p *program) runSelect(fr *frame, instr *ssa.Select) (continuation, []value.Value, error) {
	cases := make([]reflect.SelectCase, 0, len(instr.States)+1)
	if !instr.Blocking {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
	}
	type recvSlot struct {
		idx     int
		elemRT  reflect.Type
		assigned bool
	}
	for _, st := range instr.States {
		ch, err := p.readValue(fr, st.Chan)
		if err != nil {
			return contNext, nil, err
		}
		chRV, err := p.reflectOf(ch, nil)
		if err != nil {
			return contNext, nil, err
		}
		for chRV.Kind() == reflect.Interface {
			chRV = chRV.Elem()
		}
		switch st.Dir {
		case types.SendOnly:
			val, err := p.readValue(fr, st.Send)
			if err != nil {
				return contNext, nil, err
			}
			vRV, err := p.reflectOf(val, chRV.Type().Elem())
			if err != nil {
				return contNext, nil, err
			}
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: chRV,
				Send: vRV,
			})
		case types.RecvOnly:
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: chRV,
			})
		}
	}
	chosen, recv, recvOK := reflect.Select(cases)
	if !instr.Blocking {
		chosen-- // default has index -1 in SSA terms
	}

	tt, ok := instr.Type().(*types.Tuple)
	if !ok {
		return contNext, nil, fmt.Errorf("interp: Select type is not tuple")
	}
	rt, err := p.resolver.ResolveType(tt)
	if err != nil {
		return contNext, nil, err
	}
	holder := reflect.New(rt).Elem()
	holder.Field(0).SetInt(int64(chosen))
	holder.Field(1).SetBool(recvOK)
	// Recv result fields (one per RecvOnly state) follow.
	recvFieldIdx := 2
	for i, st := range instr.States {
		if st.Dir != types.RecvOnly {
			continue
		}
		f := holder.Field(recvFieldIdx)
		if int(i) == chosen && recvOK {
			if f.Kind() == reflect.Interface {
				f.Set(recv)
			} else if recv.Type() != f.Type() && recv.Type().ConvertibleTo(f.Type()) {
				f.Set(recv.Convert(f.Type()))
			} else {
				f.Set(recv)
			}
		}
		recvFieldIdx++
	}
	fr.cells[instr] = &Cell{
		Name:  instr.Name(),
		Type:  instr.Type(),
		Value: reflectValue(holder),
	}
	return contNext, nil, nil
}
