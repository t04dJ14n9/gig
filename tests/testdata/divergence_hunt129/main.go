package divergence_hunt129

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// Round 129: Struct tags, JSON marshal/unmarshal
// ============================================================================

type Tagged struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Skip  string `json:"-"`
}

func StructTagJSON() string {
	t := Tagged{Name: "test", Value: 42, Skip: "hidden"}
	b, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(b)
}

func StructTagUnmarshal() string {
	data := `{"name":"hello","value":99,"Skip":"ignored"}`
	var t Tagged
	err := json.Unmarshal([]byte(data), &t)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("name=%s-value=%d-skip=%s", t.Name, t.Value, t.Skip)
}

type OmitEmpty struct {
	Name  string `json:"name,omitempty"`
	Value int    `json:"value,omitempty"`
	Empty string `json:"empty,omitempty"`
}

func StructTagOmitEmpty() string {
	t := OmitEmpty{Name: "test", Value: 0, Empty: ""}
	b, _ := json.Marshal(t)
	return string(b)
}

type NestedJSON struct {
	Outer string     `json:"outer"`
	Inner Tagged     `json:"inner"`
	Items []string   `json:"items"`
}

func StructNestedJSON() string {
	t := NestedJSON{
		Outer: "hello",
		Inner: Tagged{Name: "inner", Value: 7},
		Items: []string{"a", "b"},
	}
	b, _ := json.Marshal(t)
	return string(b)
}

func StructMapJSON() string {
	m := map[string]int{"x": 1, "y": 2}
	b, _ := json.Marshal(m)
	return string(b)
}

func StructSliceJSON() string {
	s := []int{10, 20, 30}
	b, _ := json.Marshal(s)
	return string(b)
}

type StringInt struct {
	Key string
	Val int
}

func StructSliceOfStructs() string {
	s := []StringInt{{Key: "a", Val: 1}, {Key: "b", Val: 2}}
	b, _ := json.Marshal(s)
	return string(b)
}

func StructBoolJSON() string {
	b, _ := json.Marshal(true)
	return string(b)
}

func StructNilJSON() string {
	var s []string
	b, _ := json.Marshal(s)
	return string(b)
}
