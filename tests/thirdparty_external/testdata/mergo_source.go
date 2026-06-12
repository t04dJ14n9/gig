package main

import "dario.cat/mergo"

func MergoMerge() string {
	type Dest struct {
		Name string
		Age  int
	}
	dst := Dest{Name: "Alice", Age: 25}
	src := Dest{Name: "Bob"}
	mergo.Merge(&dst, src)
	return dst.Name
}

func MergoMergeWithOverride() string {
	type Dest struct {
		Name string
		Age  int
	}
	dst := Dest{Name: "Alice", Age: 25}
	src := Dest{Name: "Bob"}
	mergo.Merge(&dst, src, mergo.WithOverride)
	return dst.Name
}

func MergoMap() string {
	dst := map[string]string{"a": "1", "b": "2"}
	src := map[string]string{"c": "3"}
	mergo.Map(&dst, src)
	return dst["c"]
}
