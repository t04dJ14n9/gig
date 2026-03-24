// channel.go implements channel operations: Send, SendContext, Recv, RecvContext, and Close
// with context cancellation support.
package value

import (
	"context"
	"reflect"
)

// Send sends a value on a channel.
func (v Value) Send(val Value) {
	if rv, ok := v.obj.(reflect.Value); ok {
		rv.Send(val.ToReflectValue(rv.Type().Elem()))
		return
	}
	panic("invalid reflect.Value in Send()")
}

// SendContext sends a value on a channel with context cancellation support.
// Returns ctx.Err() if the context is cancelled before the send completes.
//
// NOTE: We must NOT use TrySend as a fast path here. In a concurrent setting
// (e.g. "go func(){ v := <-ch; ch <- v*2 }(); ch <- 10; return <-ch"),
// TrySend can fill the buffer before the goroutine runs, and a subsequent
// TryRecv on the same goroutine steals the value back — causing the child
// goroutine to block forever (deadlock). reflect.Select properly participates
// in Go's channel scheduling and avoids this race.
func (v Value) SendContext(ctx context.Context, val Value) error {
	rv, ok := v.obj.(reflect.Value)
	if !ok {
		panic("invalid reflect.Value in SendContext()")
	}

	sendRV := val.ToReflectValue(rv.Type().Elem())

	done := ctx.Done()
	if done == nil {
		// No cancellation possible — blocking send.
		rv.Send(sendRV)
		return nil
	}

	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectSend,
			Chan: rv,
			Send: sendRV,
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(done),
		},
	}

	chosen, _, _ := reflect.Select(cases)
	if chosen == 1 {
		return ctx.Err()
	}
	return nil
}

// TrySend tries to send a value on a channel (non-blocking).
func (v Value) TrySend(val Value) bool {
	if rv, ok := v.obj.(reflect.Value); ok {
		return rv.TrySend(val.ToReflectValue(rv.Type().Elem()))
	}
	panic("invalid reflect.Value in TrySend()")
}

// Recv receives a value from a channel.
func (v Value) Recv() (Value, bool) {
	if rv, ok := v.obj.(reflect.Value); ok {
		val, ok := rv.Recv()
		return MakeFromReflect(val), ok
	}
	panic("invalid reflect.Value in Recv()")
}

// RecvContext receives a value from a channel with context cancellation support.
// Returns (value, received, error) where error is ctx.Err() if cancelled.
//
// NOTE: We must NOT use TryRecv as a fast path here. See SendContext comment
// for the full explanation — TryRecv races with goroutine scheduling and can
// steal values from the buffer before the intended consumer runs.
func (v Value) RecvContext(ctx context.Context) (Value, bool, error) {
	rv, ok := v.obj.(reflect.Value)
	if !ok {
		panic("invalid reflect.Value in RecvContext()")
	}

	done := ctx.Done()
	if done == nil {
		// No cancellation possible — blocking receive.
		recv, recvOK := rv.Recv()
		return MakeFromReflect(recv), recvOK, nil
	}

	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: rv,
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(done),
		},
	}

	chosen, recv, recvOK := reflect.Select(cases)
	if chosen == 1 {
		return MakeNil(), false, ctx.Err()
	}
	return MakeFromReflect(recv), recvOK, nil
}

// TryRecv tries to receive a value from a channel (non-blocking).
func (v Value) TryRecv() (Value, bool) {
	if rv, ok := v.obj.(reflect.Value); ok {
		val, ok := rv.TryRecv()
		return MakeFromReflect(val), ok
	}
	panic("invalid reflect.Value in TryRecv()")
}

// Close closes a channel.
func (v Value) Close() {
	if rv, ok := v.obj.(reflect.Value); ok {
		rv.Close()
		return
	}
	panic("invalid reflect.Value in Close()")
}

// CanInterface reports whether Interface can be used without panicking.
func (v Value) CanInterface() bool {
	if rv, ok := v.obj.(reflect.Value); ok {
		return rv.CanInterface()
	}
	return true
}
