package main

import "strings"

func Test(r *strings.Reader) int {
	return r.Len()
}
