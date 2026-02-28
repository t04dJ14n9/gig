package main

import (
	"strings"
	"strconv"
)

// ExtCallReflect benchmarks external function calls that use reflect.Call (no DirectCall).
// strings.NewReplacer returns *Replacer, and Replace is a method call → reflect path.
// Also exercises strconv.FormatFloat which has DirectCall, mixed with method calls.
func ExtCallReflect() int {
	r := strings.NewReplacer("a", "b", "c", "d")
	sum := 0
	for i := 0; i < 1000; i++ {
		s := strconv.Itoa(i)
		result := r.Replace(s)
		sum += len(result)
	}
	return sum
}
