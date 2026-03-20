package main

import "bytes"

func Test(buf *bytes.Buffer) int {
	buf.WriteString(" gig was here")
	return buf.Len()
}
