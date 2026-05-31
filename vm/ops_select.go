package vm

import (
	"fmt"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// selectState is the decoded stack payload for one source-level select arm.
// Send arms carry both the channel and value; receive arms only need channel.
type selectState struct {
	ch      value.Value
	sendVal value.Value
	isSend  bool
}

func (v *vm) executeSelect(frame *Frame) error {
	// OpSelect performs a select statement using reflect.Select.
	// Operands: [meta_idx:2]
	// Stack (bottom to top): for each state, Chan; if send, also SendVal.
	// Result pushed: tuple (index, recvOk, recv_0, ..., recv_{n-1})
	meta, err := v.selectMeta(frame)
	if err != nil {
		return err
	}

	states := v.popSelectStates(meta)
	cases := v.selectCases(meta, states)
	chosen, recv, recvOK := reflect.Select(cases)
	chosen, err = v.normalizeSelectChoice(meta, chosen)
	if err != nil {
		return err
	}

	v.push(value.FromInterface(selectResultTuple(meta, chosen, recv, recvOK)))
	return nil
}

func (v *vm) selectMeta(frame *Frame) (bytecode.SelectMeta, error) {
	// The compiler stores select layout in constants so the VM instruction can
	// stay compact: a single operand points to channel count, directions, and
	// receive-slot metadata shared by the helpers below.
	metaIdx := frame.readUint16()
	meta, ok := v.program.Constants[metaIdx].(bytecode.SelectMeta)
	if !ok {
		return bytecode.SelectMeta{}, fmt.Errorf("OpSelect: invalid meta at index %d", metaIdx)
	}
	return meta, nil
}

func (v *vm) popSelectStates(meta bytecode.SelectMeta) []selectState {
	// Pop channels and send values from stack (they were pushed in order,
	// so we need to pop in reverse).
	states := make([]selectState, meta.NumStates)
	for i := meta.NumStates - 1; i >= 0; i-- {
		if meta.Dirs[i] { // send
			states[i].sendVal = v.pop()
			states[i].ch = v.pop()
			states[i].isSend = true
		} else { // recv
			states[i].ch = v.pop()
		}
	}
	return states
}

func (v *vm) selectCases(meta bytecode.SelectMeta, states []selectState) []reflect.SelectCase {
	// Build reflect.SelectCase slice.
	// Add 1 for default case (non-blocking) or context cancellation case (blocking).
	numCases := meta.NumStates + 1
	cases := make([]reflect.SelectCase, numCases)
	for i := 0; i < meta.NumStates; i++ {
		cases[i] = selectCase(states[i])
	}
	if !meta.Blocking {
		cases[meta.NumStates] = reflect.SelectCase{Dir: reflect.SelectDefault}
	} else {
		// Inject context cancellation case for blocking select.
		cases[meta.NumStates] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(v.ctx.Done()),
		}
	}
	return cases
}

func selectCase(state selectState) reflect.SelectCase {
	// reflect.Select requires native reflect.Value channels, while Gig values
	// keep the interpreter's conversion rules in one place.
	rv, _ := state.ch.ReflectValue()
	if state.isSend {
		sendRV := state.sendVal.ToReflectValue(rv.Type().Elem())
		return reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: rv,
			Send: sendRV,
		}
	}
	return reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: rv,
	}
}

func (v *vm) normalizeSelectChoice(meta bytecode.SelectMeta, chosen int) (int, error) {
	// Check if context was cancelled (chosen == meta.NumStates in blocking mode).
	if meta.Blocking && chosen == meta.NumStates {
		return 0, v.ctx.Err()
	}

	// Adjust chosen index: if default was selected, chosen == meta.NumStates maps to -1.
	if !meta.Blocking && chosen == meta.NumStates {
		return -1, nil
	}
	return chosen, nil
}

func selectResultTuple(meta bytecode.SelectMeta, chosen int, recv reflect.Value, recvOK bool) []value.Value {
	// Build result tuple: (index, recvOk, recv_0, ..., recv_{n-1})
	// Receive values are aligned with compiler-provided receive slots, so
	// unchosen receive arms still get their zero Value placeholders.
	tupleLen := 2 + meta.NumRecv
	tuple := make([]value.Value, tupleLen)
	tuple[0] = value.MakeInt(int64(chosen))
	tuple[1] = value.MakeBool(recvOK)

	// Fill recv values: for each recv state, if it was chosen, set the value.
	recvIdx := 0
	for i := 0; i < meta.NumStates; i++ {
		if !meta.Dirs[i] {
			if i == chosen {
				tuple[2+recvIdx] = value.MakeFromReflect(recv)
			} else {
				tuple[2+recvIdx] = value.MakeNil()
			}
			recvIdx++
		}
	}
	return tuple
}
