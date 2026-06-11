package divergence_hunt50

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 50: Comprehensive integration - real-world-like scenarios
// ============================================================================

func StudentRanking() int {
	type Student struct {
		Name  string
		Score int
	}
	students := []Student{
		{"Alice", 85}, {"Bob", 92}, {"Charlie", 78},
		{"Diana", 95}, {"Eve", 88},
	}
	sort.Slice(students, func(i, j int) bool { return students[i].Score > students[j].Score })
	// Top student: Diana(95), Second: Bob(92)
	return students[0].Score + students[1].Score
}

func TextAnalyzer() string {
	text := "The quick brown fox jumps over the lazy dog"
	words := strings.Fields(text)
	return fmt.Sprintf("%d:%d", len(words), len(text))
}

func ShoppingCart() float64 {
	type Item struct {
		Name  string
		Price float64
		Qty   int
	}
	cart := []Item{
		{"apple", 1.5, 3},
		{"banana", 0.75, 5},
		{"cherry", 2.0, 2},
	}
	total := 0.0
	for _, item := range cart {
		total += item.Price * float64(item.Qty)
	}
	return total
}

func JSONAPIResponse() int {
	type Response struct {
		Status  string `json:"status"`
		Count   int    `json:"count"`
		Results []int  `json:"results"`
	}
	resp := Response{
		Status:  "ok",
		Count:   3,
		Results: []int{10, 20, 30},
	}
	data, _ := json.Marshal(resp)
	var decoded Response
	json.Unmarshal(data, &decoded)
	return decoded.Count + decoded.Results[0]
}

func MatrixRotate() int {
	// Rotate matrix 90 degrees clockwise
	m := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	n := len(m)
	rotated := make([][]int, n)
	for i := range rotated {
		rotated[i] = make([]int, n)
		for j := range rotated[i] {
			rotated[i][j] = m[n-1-j][i]
		}
	}
	return rotated[0][0]*100 + rotated[0][1]*10 + rotated[0][2]
}

func DataDedup() int {
	data := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	seen := map[int]bool{}
	unique := []int{}
	for _, v := range data {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}
	sort.Ints(unique)
	return len(unique)
}

func StringTemplate() string {
	template := "Hello {name}, your order #{id} is {status}."
	s := strings.ReplaceAll(template, "{name}", "Alice")
	s = strings.ReplaceAll(s, "{id}", "123")
	s = strings.ReplaceAll(s, "{status}", "shipped")
	return s
}

func GroupByCategory() int {
	type Item struct {
		Category string
		Value    int
	}
	items := []Item{
		{"A", 10}, {"B", 20}, {"A", 30},
		{"B", 40}, {"A", 50},
	}
	groups := map[string][]int{}
	for _, item := range items {
		groups[item.Category] = append(groups[item.Category], item.Value)
	}
	sumA := 0
	for _, v := range groups["A"] { sumA += v }
	return sumA
}

func LRUPrototype() int {
	// Simplified LRU: just track access order
	keys := []string{}
	cache := map[string]int{"a": 1, "b": 2, "c": 3}
	access := func(key string) int {
		v := cache[key]
		// Move to front (simplified: just append if not present)
		newKeys := []string{key}
		for _, k := range keys {
			if k != key { newKeys = append(newKeys, k) }
		}
		keys = newKeys
		return v
	}
	access("b")
	access("a")
	return len(keys)
}

func FibonacciMemo() int {
	memo := map[int]int{}
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 { return n }
		if v, ok := memo[n]; ok { return v }
		result := fib(n-1) + fib(n-2)
		memo[n] = result
		return result
	}
	return fib(20)
}

func FmtTable() string {
	rows := []struct {
		Name string
		Age  int
	}{
		{"Alice", 30},
		{"Bob", 25},
	}
	var b strings.Builder
	for _, row := range rows {
		b.WriteString(fmt.Sprintf("%-5s %3d\n", row.Name, row.Age))
	}
	return b.String()
}

func Pipeline() int {
	// Data pipeline: filter -> transform -> aggregate
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// Filter even
	filtered := []int{}
	for _, v := range data {
		if v%2 == 0 { filtered = append(filtered, v) }
	}
	// Transform: square
	transformed := make([]int, len(filtered))
	for i, v := range filtered { transformed[i] = v * v }
	// Aggregate: sum
	sum := 0
	for _, v := range transformed { sum += v }
	return sum
}
