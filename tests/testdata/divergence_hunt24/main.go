package divergence_hunt24

import (
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 24: Comprehensive real-world patterns
// ============================================================================

func SortAndDedupe() int {
	s := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	sort.Ints(s)
	result := []int{s[0]}
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] { result = append(result, s[i]) }
	}
	return len(result)
}

func WordFrequency() int {
	text := "the cat sat on the mat the cat"
	words := strings.Fields(text)
	freq := map[string]int{}
	for _, w := range words { freq[w]++ }
	return freq["the"]
}

func CSVLikeParsing() int {
	line := "name,age,score"
	fields := strings.Split(line, ",")
	return len(fields)
}

func HistogramFromData() string {
	data := []int{1, 2, 2, 3, 3, 3, 4}
	hist := map[int]int{}
	for _, v := range data { hist[v]++ }
	return fmt.Sprintf("%d:%d:%d:%d", hist[1], hist[2], hist[3], hist[4])
}

func FlattenJSON() int {
	type Nested struct {
		A int
		B struct{ C int }
	}
	n := Nested{A: 1}
	n.B.C = 2
	return n.A + n.B.C
}

func StringTokenize() int {
	s := "hello,world;foo|bar"
	count := 1
	for _, c := range s {
		if c == ',' || c == ';' || c == '|' { count++ }
	}
	return count
}

func MatrixRowColSum() int {
	m := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	rowSum := 0
	for _, v := range m[0] { rowSum += v }
	colSum := 0
	for _, row := range m { colSum += row[0] }
	return rowSum + colSum
}

func StringTemplate() string {
	template := "Hello, {name}!"
	result := strings.ReplaceAll(template, "{name}", "World")
	return result
}

func MapTransformKeys() int {
	m := map[string]int{"1": 10, "2": 20, "3": 30}
	result := map[int]int{}
	for k, v := range m {
		n := 0
		for _, c := range k { n = n*10 + int(c-'0') }
		result[n] = v
	}
	return result[1] + result[2] + result[3]
}

func SlicePartitionPoint() int {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	pivot := 5
	less, greater := []int{}, []int{}
	for _, v := range s {
		if v < pivot { less = append(less, v) } else { greater = append(greater, v) }
	}
	return len(less)*10 + len(greater)
}

func NestedLoopBreak() int {
	found := false
	result := 0
	for i := 0; i < 5 && !found; i++ {
		for j := 0; j < 5; j++ {
			if i*5+j == 17 {
				found = true
				result = i*10 + j
				break
			}
		}
	}
	return result
}

func RecursiveSum() int {
	var sum func(s []int) int
	sum = func(s []int) int {
		if len(s) == 0 { return 0 }
		return s[0] + sum(s[1:])
	}
	return sum([]int{1, 2, 3, 4, 5})
}

func ReverseSliceInPlace() int {
	s := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func MapToSlice() int {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	keys := make([]string, 0, len(m))
	for k := range m { keys = append(keys, k) }
	sort.Strings(keys)
	return m[keys[0]] + m[keys[1]] + m[keys[2]]
}

func StringDiff() int {
	a, b := "abcde", "abXde"
	diff := 0
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] { diff++ }
	}
	return diff
}

func FmtSlice() string {
	s := []int{1, 2, 3}
	return fmt.Sprintf("%v", s)
}
