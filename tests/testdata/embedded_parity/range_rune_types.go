package main

import "fmt"

// Result checks that range over strings yields byte offsets and rune-typed values.
func Result() string {
	total := int32(0)
	lastOffset := -1
	okCount := 0
	for off, r := range "a世" {
		lastOffset = off
		if _, ok := any(r).(rune); ok {
			okCount++
		}
		total += r
	}
	return fmt.Sprintf("%d:%d:%d", lastOffset, okCount, total)
}
