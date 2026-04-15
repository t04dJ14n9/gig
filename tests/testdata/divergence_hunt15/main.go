package divergence_hunt15

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 15: More complex programs, data processing, transformation patterns
// ============================================================================

// WordCount tests word counting pattern
func WordCount() int {
	s := "hello world hello go hello"
	words := strings.Fields(s)
	counts := map[string]int{}
	for _, w := range words {
		counts[w]++
	}
	return counts["hello"]
}

// TopKElements tests finding top K elements
func TopKElements() int {
	data := []int{5, 3, 8, 1, 9, 2, 7, 4, 6}
	sort.Ints(data)
	return data[len(data)-1] + data[len(data)-2] + data[len(data)-3]
}

// FlattenAndSum tests flattening and summing nested data
func FlattenAndSum() int {
	nested := [][]int{{1, 2, 3}, {4, 5}, {6, 7, 8, 9}}
	sum := 0
	for _, inner := range nested {
		for _, v := range inner {
			sum += v
		}
	}
	return sum
}

// FrequencyCount tests frequency counting
func FrequencyCount() int {
	s := "abracadabra"
	counts := map[rune]int{}
	for _, c := range s {
		counts[c]++
	}
	return counts['a']
}

// ReverseString tests string reversal
func ReverseString() string {
	s := "hello"
	result := ""
	for _, c := range s {
		result = string(c) + result
	}
	return result
}

// StringPermutationCheck tests if two strings are permutations
func StringPermutationCheck() bool {
	a, b := "abc", "cab"
	if len(a) != len(b) { return false }
	counts := map[rune]int{}
	for _, c := range a { counts[c]++ }
	for _, c := range b { counts[c]-- }
	for _, v := range counts {
		if v != 0 { return false }
	}
	return true
}

// MatrixSum tests matrix sum
func MatrixSum() int {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	sum := 0
	for _, row := range matrix {
		for _, v := range row {
			sum += v
		}
	}
	return sum
}

// MatrixTranspose tests matrix transpose
func MatrixTranspose() int {
	matrix := [][]int{
		{1, 2},
		{3, 4},
	}
	return matrix[0][1]*10 + matrix[1][0] // original (2,3) -> after transpose still (2,3)
}

// JSONEncodeDecode tests JSON encode/decode round trip
func JSONEncodeDecode() int {
	type Point struct{ X, Y int }
	p1 := Point{X: 10, Y: 20}
	data, _ := json.Marshal(p1)
	var p2 Point
	json.Unmarshal(data, &p2)
	return p2.X + p2.Y
}

// StringCompression tests simple string compression
func StringCompression() string {
	s := "aaabbc"
	result := ""
	i := 0
	for i < len(s) {
		ch := s[i]
		count := 0
		for i < len(s) && s[i] == ch {
			count++
			i++
		}
		result += fmt.Sprintf("%c%d", ch, count)
	}
	return result
}

// UniqueElements tests getting unique elements
func UniqueElements() int {
	s := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	seen := map[int]bool{}
	for _, v := range s { seen[v] = true }
	return len(seen)
}

// IntersectSlices tests finding intersection of two slices
func IntersectSlices() int {
	a := []int{1, 2, 3, 4, 5}
	b := []int{3, 4, 5, 6, 7}
	set := map[int]bool{}
	for _, v := range a { set[v] = true }
	count := 0
	for _, v := range b {
		if set[v] { count++ }
	}
	return count
}

// MergeSortedSlices tests merging two sorted slices
func MergeSortedSlices() int {
	a := []int{1, 3, 5}
	b := []int{2, 4, 6}
	result := []int{}
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result[0] + result[1] + result[5]
}

// MovingAverage tests moving average calculation
func MovingAverage() int {
	data := []int{1, 2, 3, 4, 5}
	window := 3
	sum := 0
	for i := 0; i < window; i++ { sum += data[i] }
	return sum / window
}

// SpiralMatrix tests spiral matrix access
func SpiralMatrix() int {
	m := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	return m[0][0] + m[0][2] + m[2][2] + m[2][0]
}

// FmtStructFormatting tests struct formatting with %+v
func FmtStructFormatting() string {
	type P struct{ X, Y int }
	p := P{X: 1, Y: 2}
	return fmt.Sprintf("%v", p)
}
