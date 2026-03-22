package main

import (
	"bytes"
	"encoding/json"
)

func Test(buf *bytes.Buffer, enc *json.Encoder) string {
	buf.WriteString("raw: ")
	enc.Encode(map[string]string{"k": "v"})
	return buf.String()
}
