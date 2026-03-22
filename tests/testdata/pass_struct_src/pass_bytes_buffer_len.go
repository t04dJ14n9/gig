package main

import "bytes"

func Test(buf *bytes.Buffer) int {
	return buf.Len()
}
