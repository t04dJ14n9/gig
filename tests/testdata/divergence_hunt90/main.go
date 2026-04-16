package divergence_hunt90

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 90: Final comprehensive integration test
// ============================================================================

func Comprehensive1() int {
	// Student grade analyzer with multiple aggregation
	type Student struct {
		Name  string
		Score int
	}
	students := []Student{
		{"Alice", 85}, {"Bob", 92}, {"Charlie", 78},
		{"Diana", 95}, {"Eve", 88}, {"Frank", 72},
	}
	total := 0
	for _, s := range students {
		total += s.Score
	}
	avg := total / len(students)
	aboveAvg := 0
	for _, s := range students {
		if s.Score > avg {
			aboveAvg++
		}
	}
	return aboveAvg*1000 + avg
}

func Comprehensive2() string {
	// Text processing with word frequency
	text := "Hello World Hello Go Hello Code"
	words := strings.Fields(text)
	counts := map[string]int{}
	for _, w := range words {
		counts[strings.ToLower(w)]++
	}
	return fmt.Sprintf("%d:%d:%d", counts["hello"], counts["world"], counts["go"])
}

func Comprehensive3() int {
	// JSON round-trip with nested data
	type Data struct {
		Items []int `json:"items"`
		Total int   `json:"total"`
	}
	d := Data{Items: []int{10, 20, 30}, Total: 60}
	b, _ := json.Marshal(d)
	var decoded Data
	json.Unmarshal(b, &decoded)
	return decoded.Total
}

func Comprehensive4() int {
	// Matrix multiplication
	a := [][]int{{1, 2}, {3, 4}}
	b := [][]int{{5, 6}, {7, 8}}
	c := make([][]int, 2)
	for i := range c {
		c[i] = make([]int, 2)
		for j := range c[i] {
			for k := 0; k < 2; k++ {
				c[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return c[0][0]*1000 + c[0][1]*100 + c[1][0]*10 + c[1][1]
}

func Comprehensive5() string {
	// String manipulation pipeline
	s := "  The Quick Brown Fox  "
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	words := strings.Split(s, " ")
	sort.Strings(words)
	return strings.Join(words, "-")
}

func Comprehensive6() int {
	// Recursive tree traversal with sum
	type Node struct {
		Value int
		Left  *Node
		Right *Node
	}
	var sumTree func(*Node) int
	sumTree = func(n *Node) int {
		if n == nil {
			return 0
		}
		return n.Value + sumTree(n.Left) + sumTree(n.Right)
	}
	root := &Node{
		Value: 1,
		Left:  &Node{Value: 2, Left: &Node{Value: 4}, Right: &Node{Value: 5}},
		Right: &Node{Value: 3, Left: &Node{Value: 6}, Right: &Node{Value: 7}},
	}
	return sumTree(root)
}

func Comprehensive7() int {
	// Data pipeline: filter, map, reduce
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	filtered := []int{}
	for _, v := range data {
		if v%2 == 0 {
			filtered = append(filtered, v)
		}
	}
	doubled := make([]int, len(filtered))
	for i, v := range filtered {
		doubled[i] = v * 2
	}
	sum := 0
	for _, v := range doubled {
		sum += v
	}
	return sum
}

func Comprehensive8() int {
	// Error handling chain
	process := func(data []int) (int, error) {
		if len(data) == 0 {
			return 0, fmt.Errorf("empty data")
		}
		sum := 0
		for _, v := range data {
			sum += v
		}
		return sum, nil
	}
	validate := func(data []int) error {
		for _, v := range data {
			if v < 0 {
				return fmt.Errorf("negative value")
			}
		}
		return nil
	}
	data := []int{1, 2, 3, 4, 5}
	if err := validate(data); err != nil {
		return -1
	}
	result, err := process(data)
	if err != nil {
		return -2
	}
	return result
}

func Comprehensive9() int {
	// Priority queue simulation
	type Item struct{ Value, Priority int }
	items := []Item{
		{10, 3}, {20, 1}, {30, 2},
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Priority < items[j].Priority })
	return items[0].Value
}

func Comprehensive10() string {
	// Map-reduce pattern
	data := []string{"apple", "banana", "avocado", "blueberry", "cherry"}
	mapped := map[string]int{}
	for _, s := range data {
		mapped[s] = len(s)
	}
	total := 0
	for _, v := range mapped {
		total += v
	}
	return fmt.Sprintf("%d", total)
}
