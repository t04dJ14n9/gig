package vm

import "fmt"

func (v *vm) executePrint(frame *Frame) {
	n := frame.readByte()
	for i := 0; i < int(n); i++ {
		val := v.pop()
		fmt.Print(val.Interface())
	}
}

func (v *vm) executePrintln(frame *Frame) {
	n := frame.readByte()
	args := make([]any, n)
	for i := int(n) - 1; i >= 0; i-- {
		args[i] = v.pop().Interface()
	}
	fmt.Println(args...)
}
