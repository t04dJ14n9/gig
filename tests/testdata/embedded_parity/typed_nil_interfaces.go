package main

import "fmt"

// Result checks typed nil values stored in empty interfaces.
func Result() string {
	var p *[]int
	var i any = p
	ptrNil := i == nil

	var m map[string]int
	i = m
	mapNil := i == nil

	var s []int
	i = s
	sliceNil := i == nil

	var f func()
	i = f
	funcNil := i == nil

	var ch chan int
	i = ch
	chanNil := i == nil

	i = nil
	nilNil := i == nil

	values := make([]any, 3)
	values[1] = 42
	return fmt.Sprintf("%t:%t:%t:%t:%t:%t:%t:%t:%d",
		ptrNil, mapNil, sliceNil, funcNil, chanNil, nilNil,
		values[0] == nil, values[1] == nil, values[1])
}
