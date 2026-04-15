package divergence_hunt22

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// Round 22: More encoding, format patterns, error wrapping
// ============================================================================

func JSONMarshalInt() string {
	b, _ := json.Marshal(42)
	return string(b)
}

func JSONMarshalString() string {
	b, _ := json.Marshal("hello")
	return string(b)
}

func JSONMarshalFloat() string {
	b, _ := json.Marshal(3.14)
	return string(b)
}

func JSONUnmarshalInt() int {
	var x int
	json.Unmarshal([]byte("42"), &x)
	return x
}

func JSONUnmarshalString() string {
	var x string
	json.Unmarshal([]byte(`"hello"`), &x)
	return x
}

func JSONUnmarshalFloat() float64 {
	var x float64
	json.Unmarshal([]byte("3.14"), &x)
	return x
}

func JSONUnmarshalArray() int {
	var x []int
	json.Unmarshal([]byte("[1,2,3]"), &x)
	return x[0] + x[1] + x[2]
}

func FmtVerbP() string {
	// Pointer addresses vary between runs - just test that it formats
	x := 42
	p := &x
	_ = p
	return fmt.Sprintf("%d", x)
}

func FmtVerbT() string {
	x := 42
	return fmt.Sprintf("%T", x)
}

func FmtVerbV() string {
	type S struct{ X int }
	s := S{X: 42}
	return fmt.Sprintf("%v", s)
}

func FmtVerbPlusV() string {
	type S struct{ X int }
	s := S{X: 42}
	return fmt.Sprintf("%+v", s)
}

func FmtVerbHashV() string {
	type S struct{ X int }
	s := S{X: 42}
	return fmt.Sprintf("%#v", s)
}

func FmtSprintfPointer() string {
	x := 42
	p := &x
	return fmt.Sprintf("%d", *p)
}

func ErrorWrap() string {
	err := fmt.Errorf("outer: %w", fmt.Errorf("inner"))
	return err.Error()
}

func ErrorIs() bool {
	inner := fmt.Errorf("inner")
	outer := fmt.Errorf("outer: %w", inner)
	return fmt.Sprintf("%v", outer) != ""
}

func JSONNestedMap() int {
	data := `{"a":{"x":1},"b":{"y":2}}`
	var m map[string]map[string]int
	json.Unmarshal([]byte(data), &m)
	return m["a"]["x"] + m["b"]["y"]
}

func JSONStructTag() int {
	type Item struct {
		Name  string `json:"name"`
		Value int    `json:"val"`
	}
	data := `{"name":"test","val":42}`
	var item Item
	json.Unmarshal([]byte(data), &item)
	return item.Value
}

func JSONOmitEmpty() string {
	type S struct {
		X int    `json:"x"`
		Y string `json:"y,omitempty"`
	}
	s := S{X: 1}
	b, _ := json.Marshal(s)
	return string(b)
}

func FmtWidthInt() string {
	return fmt.Sprintf("[%5d][%-5d]", 42, 42)
}

func FmtFloatScientific() string {
	return fmt.Sprintf("%e", 1234.5678)
}
