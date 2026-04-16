package divergence_hunt38

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// Round 38: Complex struct patterns - embedding, field access, nested types
// ============================================================================

type Inner38 struct {
	Value int
}

type Middle38 struct {
	Inner38
	Label string
}

type Outer38 struct {
	Middle38
	Name string
}

func DeepEmbedding() int {
	o := Outer38{
		Middle38: Middle38{
			Inner38: Inner38{Value: 42},
			Label:   "test",
		},
		Name: "outer",
	}
	return o.Value + len(o.Label) + len(o.Name) // 42 + 4 + 5
}

func EmbeddingFieldAccess() int {
	o := Outer38{Name: "hello"}
	o.Value = 10
	o.Label = "world"
	return o.Value + len(o.Label) + len(o.Name)
}

type Point38 struct {
	X, Y int
}

type Circle38 struct {
	Point38
	Radius int
}

func EmbeddedMethodAccess() int {
	c := Circle38{Point38: Point38{X: 3, Y: 4}, Radius: 5}
	return c.X + c.Y + c.Radius
}

type Address38 struct {
	City    string
	Country string
}

type Person38 struct {
	Name    string
	Age     int
	Address Address38
}

func NestedStructField() string {
	p := Person38{
		Name: "Alice",
		Age:  30,
		Address: Address38{
			City:    "NYC",
			Country: "US",
		},
	}
	return p.Address.City
}

func StructWithSliceField() int {
	type Team struct {
		Name    string
		Members []string
	}
	t := Team{
		Name:    "Alpha",
		Members: []string{"alice", "bob", "charlie"},
	}
	return len(t.Members)
}

func StructWithMapField() int {
	type Config struct {
		Values map[string]int
	}
	c := Config{
		Values: map[string]int{"x": 10, "y": 20},
	}
	return c.Values["x"] + c.Values["y"]
}

func StructWithPointerField() int {
	type Node struct {
		Value int
		Next  *Node
	}
	n2 := &Node{Value: 20}
	n1 := &Node{Value: 10, Next: n2}
	return n1.Value + n1.Next.Value
}

func StructWithFuncField() int {
	type Op struct {
		Apply func(int) int
	}
	double := Op{Apply: func(x int) int { return x * 2 }}
	return double.Apply(21)
}

func StructWithArrayField() int {
	type Vector struct {
		Data [3]float64
	}
	v := Vector{Data: [3]float64{1.0, 2.0, 3.0}}
	return int(v.Data[0] + v.Data[1] + v.Data[2])
}

func StructWithChanField() int {
	type Pipeline struct {
		Out chan int
	}
	p := Pipeline{Out: make(chan int, 1)}
	p.Out <- 42
	return <-p.Out
}

func StructWithInterfaceField() int {
	type Holder struct {
		Value any
	}
	h := Holder{Value: 42}
	if v, ok := h.Value.(int); ok {
		return v
	}
	return -1
}

func StructJSONRoundTrip() int {
	type Item struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	original := Item{Name: "test", Value: 42}
	data, _ := json.Marshal(original)
	var decoded Item
	json.Unmarshal(data, &decoded)
	return decoded.Value
}

func StructSliceJSONRoundTrip() int {
	type Item struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	items := []Item{{"a", 1}, {"b", 2}, {"c", 3}}
	data, _ := json.Marshal(items)
	var decoded []Item
	json.Unmarshal(data, &decoded)
	sum := 0
	for _, item := range decoded {
		sum += item.Value
	}
	return sum
}

func FmtNestedStruct() string {
	p := Person38{Name: "Bob", Age: 25, Address: Address38{City: "LA", Country: "US"}}
	return fmt.Sprintf("%s-%s", p.Name, p.Address.City)
}

func StructComparisonEqual() bool {
	p1 := Point38{X: 1, Y: 2}
	p2 := Point38{X: 1, Y: 2}
	return p1 == p2
}

func StructComparisonNotEqual() bool {
	p1 := Point38{X: 1, Y: 2}
	p2 := Point38{X: 1, Y: 3}
	return p1 != p2
}
