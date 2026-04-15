package divergence_hunt21

// ============================================================================
// Round 21: Map iteration order independence, slice tricks, complex data flows
// ============================================================================

func MapIterateSum() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	sum := 0
	for _, v := range m { sum += v }
	return sum
}

func SliceRotateLeft() int {
	s := []int{1, 2, 3, 4, 5}
	n := 2
	s = append(s[n:], s[:n]...)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func SliceRotateRight() int {
	s := []int{1, 2, 3, 4, 5}
	n := 2
	s = append(s[len(s)-n:], s[:len(s)-n]...)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func SliceChunk() int {
	s := []int{1, 2, 3, 4, 5, 6}
	chunkSize := 2
	chunks := 0
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) { end = len(s) }
		chunks++
		_ = s[i:end]
	}
	return chunks
}

func MapFilterSlice() int {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := []int{}
	for _, v := range s {
		if v%2 == 0 { evens = append(evens, v) }
	}
	doubled := make([]int, len(evens))
	for i, v := range evens { doubled[i] = v * 2 }
	sum := 0
	for _, v := range doubled { sum += v }
	return sum
}

func ReducePattern() int {
	s := []int{1, 2, 3, 4, 5}
	acc := 0
	for _, v := range s { acc += v }
	return acc
}

func ZipSlices() int {
	keys := []string{"a", "b", "c"}
	vals := []int{1, 2, 3}
	m := map[string]int{}
	for i := range keys { m[keys[i]] = vals[i] }
	return m["a"] + m["b"] + m["c"]
}

func SliceCompact() int {
	s := []int{1, 1, 2, 2, 3, 3}
	result := []int{}
	prev := -1
	first := true
	for _, v := range s {
		if first || v != prev {
			result = append(result, v)
			prev = v
			first = false
		}
	}
	return len(result)
}

func MapMergeOverwrite() int {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	for k, v := range m2 { m1[k] = v }
	return m1["b"]
}

func SlicePartition() int {
	s := []int{1, 2, 3, 4, 5, 6}
	evens, odds := []int{}, []int{}
	for _, v := range s {
		if v%2 == 0 { evens = append(evens, v) } else { odds = append(odds, v) }
	}
	return len(evens)*10 + len(odds)
}

func NestedMapAccess() int {
	m := map[string]map[string]int{}
	m["x"] = map[string]int{"a": 1}
	m["x"]["b"] = 2
	return m["x"]["a"] + m["x"]["b"]
}

func FlattenMap() int {
	m := map[string][]int{
		"a": {1, 2},
		"b": {3, 4},
	}
	total := 0
	for _, v := range m {
		for _, n := range v { total += n }
	}
	return total
}

func MapKeySlice() int {
	m := map[int]bool{1: true, 2: true, 3: true}
	keys := []int{}
	for k := range m { keys = append(keys, k) }
	return len(keys)
}

func SliceSlidingWindow() int {
	s := []int{1, 2, 3, 4, 5}
	windows := [][]int{}
	for i := 0; i <= len(s)-3; i++ {
		windows = append(windows, s[i:i+3])
	}
	return len(windows)
}

func MultiLevelSlice() int {
	data := [][][]int{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}
	return data[0][1][0] + data[1][0][1]
}
