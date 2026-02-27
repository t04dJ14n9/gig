package value

import "reflect"

// Send sends a value on a channel.
func (v Value) Send(val Value) {
	if rv, ok := v.obj.(reflect.Value); ok {
		rv.Send(val.ToReflectValue(rv.Type().Elem()))
		return
	}
	panic("invalid reflect.Value in Send()")
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
