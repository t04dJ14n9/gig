package main

import (
	"strings"
	"strconv"
	"math"
)

// ExtCallDirectCall benchmarks external function calls that HAVE DirectCall wrappers.
// These use typed wrappers to avoid reflect.Call.
// Calls: strings.Contains, strings.ToUpper, strconv.Itoa, math.Sqrt — all have DirectCall.
func ExtCallDirectCall() int {
	count := 0
	for i := 0; i < 1000; i++ {
		s := strconv.Itoa(i)
		if strings.Contains(s, "5") {
			count++
		}
		_ = strings.ToUpper(s)
		_ = math.Sqrt(float64(i))
	}
	return count
}
