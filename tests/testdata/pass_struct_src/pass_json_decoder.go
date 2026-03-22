package main

import "encoding/json"

func Test(dec *json.Decoder) int {
	var data map[string]int
	err := dec.Decode(&data)
	if err != nil {
		return -1
	}
	return data["value"]
}
