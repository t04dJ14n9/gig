package main

import (
	"strings"
	"strconv"
)

// ExtCallMixed benchmarks a realistic mix of external calls:
// - DirectCall functions (strings.Contains, strconv.Itoa, strings.Join)
// - Method calls (strings.NewReader().Len())
// This represents a typical string-processing workload.
func ExtCallMixed() int {
	sum := 0
	for i := 0; i < 500; i++ {
		s := strconv.Itoa(i)
		if strings.Contains(s, "3") {
			sum += len(strings.ToUpper(s))
		}
		r := strings.NewReader(s)
		sum += r.Len()
	}
	return sum
}
