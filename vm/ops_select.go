package vm

import (
	"fmt"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeSelect(frame *Frame) error {
	// OpSelect performs a select statement using reflect.Select.
	// Operands: [meta_idx:2]
	// Stack (bottom to top): for each state, Chan; if send, also SendVal.
	// Result pushed: tuple (index, recvOk, recv_0, ..., recv_{n-1})
	metaIdx := frame.readUint16()
	meta, ok := v.program.Constants[metaIdx].(bytecode.SelectMeta)
	if !ok {
		return fmt.Errorf("OpSelect: invalid meta at index %d", metaIdx)
	}

	// Pop channels and send values from stack (they were pushed in order,
	// so we need to pop in reverse).
	type stateData struct {
		ch      value.Value
		sendVal value.Value
		isSend  bool
	}
	states := make([]stateData, meta.NumStates)
	// Pop in reverse order
	for i := meta.NumStates - 1; i >= 0; i-- {
		if meta.Dirs[i] { // send
			states[i].sendVal = v.pop()
			states[i].ch = v.pop()
			states[i].isSend = true
		} else { // recv
			states[i].ch = v.pop()
		}
	}

	// Build reflect.SelectCase slice.
	// Add 1 for default case (non-blocking) or context cancellation case (blocking).
	numCases := meta.NumStates + 1
	cases := make([]reflect.SelectCase, numCases)
	for i := 0; i < meta.NumStates; i++ {
		rv, _ := states[i].ch.ReflectValue()
		if states[i].isSend {
			sendRV := states[i].sendVal.ToReflectValue(rv.Type().Elem())
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: rv,
				Send: sendRV,
			}
		} else {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: rv,
			}
		}
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

	chosen, recv, recvOK := reflect.Select(cases)

	// Check if context was cancelled (chosen == meta.NumStates in blocking mode).
	if meta.Blocking && chosen == meta.NumStates {
		return v.ctx.Err()
	}

	// Adjust chosen index: if default was selected, chosen == meta.NumStates maps to -1.
	if !meta.Blocking && chosen == meta.NumStates {
		chosen = -1
	}

	// Build result tuple: (index, recvOk, recv_0, ..., recv_{n-1})
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

	v.push(value.FromInterface(tuple))
	return nil
}
