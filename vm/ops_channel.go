package vm

import (
	"fmt"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeSend() error {
	val := v.pop()
	ch := v.pop()
	// Use a Go-level recover to catch "send on closed channel" panic
	// and convert it to a guest-level panic (recoverable by defer/recover).
	var sendErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				sendErr = fmt.Errorf("%v", r)
			}
		}()
		sendErr = ch.SendContext(v.ctx, val)
	}()
	if sendErr != nil {
		// Trigger guest-level panic so defer/recover can handle it.
		v.panicking = true
		v.panicVal = value.FromInterface(sendErr.Error())
	}
	return nil
}

func (v *vm) executeRecv() error {
	ch := v.pop()
	val, _, err := ch.RecvContext(v.ctx)
	if err != nil {
		return err
	}
	v.push(val)
	return nil
}

func (v *vm) executeRecvOk() error {
	ch := v.pop()
	val, recvOK, err := ch.RecvContext(v.ctx)
	if err != nil {
		return err
	}
	v.pushCommaOk(val, recvOK)
	return nil
}

func (v *vm) executeClose() {
	ch := v.pop()
	ch.Close()
}
