package main

import "bytes"

func Test(buf *bytes.Buffer) string {
	buf.WriteString("hello ")
	buf.WriteString("from gig")
	return buf.String()
}
