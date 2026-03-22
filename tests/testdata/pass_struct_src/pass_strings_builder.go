package main

import "strings"

func Test(b *strings.Builder) string {
	b.WriteString("world")
	return b.String()
}
