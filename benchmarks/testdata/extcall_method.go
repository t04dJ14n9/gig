package main

import (
	"strconv"
	"strings"
)

// ExtCallMethod benchmarks method calls on external types.
// strings.NewReader().Len() is a method call.
// All method calls currently go through reflect.MethodByName + reflect.Call.
func ExtCallMethod() int {
	sum := 0
	for i := 0; i < 1000; i++ {
		s := strconv.Itoa(i)
		r := strings.NewReader(s)
		sum += r.Len()
	}
	return sum
}
