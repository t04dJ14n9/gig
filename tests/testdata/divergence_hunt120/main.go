package divergence_hunt120

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 120: Comprehensive integration stress test
// ============================================================================

func Integration1() string {
	// Student grading system
	type Student struct {
		Name  string
		Score float64
	}
	students := []Student{
		{"Alice", 85.5}, {"Bob", 92.0}, {"Charlie", 78.3},
		{"Diana", 95.7}, {"Eve", 88.1},
	}
	sort.Slice(students, func(i, j int) bool { return students[i].Score > students[j].Score })
	return fmt.Sprintf("%s:%.1f", students[0].Name, students[0].Score)
}

func Integration2() int {
	// Matrix operations
	type Matrix struct {
		Data [][]int
		Rows, Cols int
	}
	mk := func(r, c int) Matrix {
		d := make([][]int, r)
		for i := range d {
			d[i] = make([]int, c)
		}
		return Matrix{Data: d, Rows: r, Cols: c}
	}
	a := mk(2, 2)
	a.Data = [][]int{{1, 2}, {3, 4}}
	sum := 0
	for _, row := range a.Data {
		for _, v := range row {
			sum += v
		}
	}
	return sum
}

func Integration3() string {
	// Text processing pipeline
	text := "Hello World Hello Go"
	words := strings.Fields(text)
	counts := map[string]int{}
	for _, w := range words {
		counts[strings.ToLower(w)]++
	}
	keys := []string{}
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result := ""
	for _, k := range keys {
		result += fmt.Sprintf("%s:%d ", k, counts[k])
	}
	return strings.TrimSpace(result)
}

func Integration4() string {
	// JSON round-trip with nested struct
	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}
	p := Person{Name: "Alice", Age: 30, Address: Address{City: "NYC", Country: "US"}}
	b, _ := json.Marshal(p)
	var decoded Person
	json.Unmarshal(b, &decoded)
	return fmt.Sprintf("%s:%d:%s", decoded.Name, decoded.Age, decoded.Address.City)
}

func Integration5() int {
	// Fibonacci with memoization
	memo := map[int]int{}
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		if v, ok := memo[n]; ok {
			return v
		}
		result := fib(n-1) + fib(n-2)
		memo[n] = result
		return result
	}
	return fib(20)
}

func Integration6() string {
	// Tree traversal with multiple orders
	type Node struct {
		Val   int
		Left  *Node
		Right *Node
	}
	root := &Node{
		Val:   4,
		Left:  &Node{Val: 2, Left: &Node{Val: 1}, Right: &Node{Val: 3}},
		Right: &Node{Val: 6, Left: &Node{Val: 5}, Right: &Node{Val: 7}},
	}
	var inorder []int
	var traverse func(*Node)
	traverse = func(n *Node) {
		if n == nil {
			return
		}
		traverse(n.Left)
		inorder = append(inorder, n.Val)
		traverse(n.Right)
	}
	traverse(root)
	return fmt.Sprintf("%v", inorder)
}

func Integration7() int {
	// Pipeline: generate → filter → reduce
	generate := func(n int) []int {
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = i + 1
		}
		return result
	}
	filter := func(data []int, pred func(int) bool) []int {
		result := []int{}
		for _, v := range data {
			if pred(v) {
				result = append(result, v)
			}
		}
		return result
	}
	reduce := func(data []int, fn func(int, int) int, init int) int {
		acc := init
		for _, v := range data {
			acc = fn(acc, v)
		}
		return acc
	}
	data := generate(10)
	evens := filter(data, func(n int) bool { return n%2 == 0 })
	sum := reduce(evens, func(a, b int) int { return a + b }, 0)
	return sum
}

func Integration8() string {
	// Error handling with custom types - use fmt.Errorf
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
	if result, err := process([]int{1, 2, 3}); err == nil {
		return fmt.Sprintf("ok:%d", result)
	}
	return "error"
}

func Integration9() string {
	// String builder pattern
	var b strings.Builder
	words := []string{"The", "quick", "brown", "fox"}
	for i, w := range words {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(w)
	}
	return b.String()
}

func Integration10() string {
	// Map-reduce with struct data
	type Sale struct {
		Product string
		Amount  float64
	}
	sales := []Sale{
		{"A", 100}, {"B", 200}, {"A", 150}, {"C", 300}, {"B", 50},
	}
	totals := map[string]float64{}
	for _, s := range sales {
		totals[s.Product] += s.Amount
	}
	keys := []string{}
	for k := range totals {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result := ""
	for _, k := range keys {
		result += fmt.Sprintf("%s:%.0f ", k, totals[k])
	}
	return strings.TrimSpace(result)
}
