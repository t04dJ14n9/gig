package main

import "fmt"

type payload struct {
	id int
}

// Result checks interface-typed slices and maps plus comma-ok map lookup.
func Result() string {
	s := []any{1, "two", payload{id: 3}, []int{4, 5}}
	s0 := s[0].(int)
	s1 := s[1].(string)
	p := s[2].(payload).id
	i := s[3].([]int)[0]

	m := map[string]any{"x": true, "y": "ok", "z": payload{id: 7}, "arr": []string{"a", "b"}}
	mv, okX := m["x"].(bool)
	_, okMissing := m["not-exist"]
	arr := m["arr"].([]string)

	return fmt.Sprintf("%d:%s:%d:%d:%v:%v:%d", s0, s1, p, i, mv, okX, len(arr)+btoi(okMissing))
}

func btoi(v bool) int {
	if v {
		return 1
	}
	return 0
}
