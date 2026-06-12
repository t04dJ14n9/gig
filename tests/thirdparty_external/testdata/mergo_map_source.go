package main

import "dario.cat/mergo"

func MergoMap() string {
	dst := map[string]string{"a": "1", "b": "2"}
	src := map[string]string{"c": "3"}
	mergo.Map(&dst, src)
	return dst["c"]
}
