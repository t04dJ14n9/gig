package divergence_hunt20

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 20: Comprehensive integration tests, real-world-like patterns
// ============================================================================

// StudentGradeSystem tests a grade management system
func StudentGradeSystem() int {
	type Student struct {
		Name  string
		Grade int
	}
	students := []Student{
		{"Alice", 85},
		{"Bob", 92},
		{"Charlie", 78},
		{"Diana", 95},
	}
	sum := 0
	for _, s := range students { sum += s.Grade }
	avg := sum / len(students)
	aboveAvg := 0
	for _, s := range students {
		if s.Grade > avg { aboveAvg++ }
	}
	return avg*10 + aboveAvg
}

// TextProcessing tests text processing
func TextProcessing() int {
	text := "The quick brown fox jumps over the lazy dog"
	words := strings.Fields(text)
	return len(words)
}

// DataTransform tests data transformation
func DataTransform() int {
	input := []int{1, 2, 3, 4, 5}
	doubled := make([]int, len(input))
	for i, v := range input { doubled[i] = v * 2 }
	return doubled[0] + doubled[4]
}

// InventorySystem tests inventory management
func InventorySystem() int {
	type Item struct {
		Name     string
		Quantity int
		Price    float64
	}
	inventory := []Item{
		{"apple", 10, 1.5},
		{"banana", 20, 0.75},
		{"cherry", 15, 2.0},
	}
	totalValue := 0.0
	for _, item := range inventory {
		totalValue += float64(item.Quantity) * item.Price
	}
	return int(totalValue)
}

// JSONProcessing tests JSON processing pipeline
func JSONProcessing() int {
	type Record struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	records := []Record{
		{1, "a", 10},
		{2, "b", 20},
		{3, "c", 30},
	}
	data, _ := json.Marshal(records)
	var decoded []Record
	json.Unmarshal(data, &decoded)
	sum := 0
	for _, r := range decoded { sum += r.Value }
	return sum
}

// StringProcessing tests string processing
func StringProcessing() string {
	lines := []string{"hello", "world", "foo"}
	return strings.Join(lines, "\n")
}

// SortAndSearch tests sort and search
func SortAndSearch() int {
	data := []int{42, 17, 23, 8, 99, 56, 31}
	sort.Ints(data)
	i := sort.SearchInts(data, 23)
	return data[i]
}

// MatrixOperations tests matrix operations
func MatrixOperations() int {
	a := [][]int{{1, 2}, {3, 4}}
	b := [][]int{{5, 6}, {7, 8}}
	result := make([][]int, 2)
	for i := range result {
		result[i] = make([]int, 2)
		for j := range result[i] {
			result[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j]
		}
	}
	return result[0][0] + result[1][1]
}

// FmtTable tests formatted table output
func FmtTable() string {
	rows := []struct{ Name string; Value int }{
		{"a", 1},
		{"bb", 22},
		{"ccc", 333},
	}
	var b strings.Builder
	for _, row := range rows {
		b.WriteString(fmt.Sprintf("%-3s %4d\n", row.Name, row.Value))
	}
	return b.String()
}

// Histogram tests histogram generation
func Histogram() int {
	data := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	hist := map[int]int{}
	for _, v := range data { hist[v]++ }
	return hist[3]
}

// ParseAndCompute tests parsing and computation
func ParseAndCompute() int {
	data := []string{"10", "20", "30", "40", "50"}
	sum := 0
	for _, s := range data {
		n, _ := fmt.Sscanf(s, "%d", &sum)
		_ = n
	}
	return sum
}

// SetOperations tests set operations
func SetOperations() int {
	a := map[int]bool{1: true, 2: true, 3: true}
	b := map[int]bool{2: true, 3: true, 4: true}
	union := len(a) + len(b)
	intersection := 0
	for k := range a {
		if b[k] { intersection++ }
	}
	return union - intersection
}

// GroupBy tests grouping pattern
func GroupBy() int {
	type Item struct{ Category string; Value int }
	items := []Item{
		{"a", 1},
		{"b", 2},
		{"a", 3},
		{"b", 4},
		{"a", 5},
	}
	groups := map[string][]int{}
	for _, item := range items {
		groups[item.Category] = append(groups[item.Category], item.Value)
	}
	return len(groups["a"])
}

// RunningSum tests running sum pattern
func RunningSum() int {
	nums := []int{1, 2, 3, 4, 5}
	sum := 0
	result := make([]int, len(nums))
	for i, v := range nums {
		sum += v
		result[i] = sum
	}
	return result[4]
}

// SlidingWindow tests sliding window pattern
func SlidingWindow() int {
	data := []int{1, 3, 5, 7, 9, 2, 4, 6, 8, 10}
	windowSize := 3
	maxSum := 0
	for i := 0; i <= len(data)-windowSize; i++ {
		sum := 0
		for j := 0; j < windowSize; j++ {
			sum += data[i+j]
		}
		if sum > maxSum { maxSum = sum }
	}
	return maxSum
}
