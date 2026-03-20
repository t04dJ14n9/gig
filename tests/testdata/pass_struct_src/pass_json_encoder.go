package main

import "encoding/json"

func Test(enc *json.Encoder) int {
	data := map[string]int{"x": 42}
	err := enc.Encode(data)
	if err != nil {
		return -1
	}
	return 1
}
