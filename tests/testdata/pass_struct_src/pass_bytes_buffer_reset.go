package main

import "bytes"

func Test(buf *bytes.Buffer) string {
	buf.Reset()
	buf.WriteString("new content")
	return buf.String()
}
