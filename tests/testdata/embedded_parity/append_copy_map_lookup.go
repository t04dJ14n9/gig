package main

import "fmt"

// Result checks append/copy semantics and map comma-ok behavior.
func Result() string {
	nums := []int{1, 2, 3, 4}
	sub := nums[1:3]
	copied := make([]int, 2)
	n := copy(copied, sub)

	appended := append(sub, 9, 10)
	sub[0] = 99

	m := map[string]int{"a": 1}
	va, okA := m["a"]
	_, okMissing := m["missing"]

	return fmt.Sprintf("%d:%v:%v:%d:%v:%v", n, copied, appended, va, okA, okMissing)
}
