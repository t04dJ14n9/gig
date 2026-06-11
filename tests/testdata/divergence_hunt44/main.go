package divergence_hunt44

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 44: Real-world algorithms - search, sort, graph, tree patterns
// ============================================================================

func BinarySearch() int {
	data := []int{1, 3, 5, 7, 9, 11, 13, 15}
	target := 7
	lo, hi := 0, len(data)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if data[mid] == target { return mid }
		if data[mid] < target { lo = mid + 1 } else { hi = mid - 1 }
	}
	return -1
}

func BubbleSort() int {
	s := []int{5, 3, 8, 1, 9, 2, 7}
	for i := 0; i < len(s)-1; i++ {
		for j := 0; j < len(s)-1-i; j++ {
			if s[j] > s[j+1] {
				s[j], s[j+1] = s[j+1], s[j]
			}
		}
	}
	return s[0]*100000 + s[1]*10000 + s[2]*1000 + s[3]*100 + s[4]*10 + s[5]
}

func InsertionSort() int {
	s := []int{5, 3, 8, 1, 9}
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
	}
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func TreeDepth() int {
	type Node struct {
		Value int
		Left  *Node
		Right *Node
	}
	var depth func(*Node) int
	depth = func(n *Node) int {
		if n == nil { return 0 }
		l := depth(n.Left)
		r := depth(n.Right)
		if l > r { return l + 1 }
		return r + 1
	}
	root := &Node{1,
		&Node{2, &Node{4, nil, nil}, &Node{5, nil, nil}},
		&Node{3, nil, &Node{6, nil, nil}},
	}
	return depth(root)
}

func GraphBFS() int {
	// Simple adjacency list BFS
	graph := map[int][]int{
		1: {2, 3},
		2: {4},
		3: {4},
		4: {},
	}
	visited := map[int]bool{}
	queue := []int{1}
	order := []int{}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if visited[node] { continue }
		visited[node] = true
		order = append(order, node)
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}
	return len(order)
}

func LongestCommonSubstrLen() int {
	a, b := "abcdef", "zcdemf"
	maxLen := 0
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			k := 0
			for i+k < len(a) && j+k < len(b) && a[i+k] == b[j+k] {
				k++
			}
			if k > maxLen { maxLen = k }
		}
	}
	return maxLen
}

func TopKFrequent() int {
	nums := []int{1, 1, 1, 2, 2, 3, 3, 3, 3}
	counts := map[int]int{}
	for _, n := range nums { counts[n]++ }
	type kv struct{ K, V int }
	var pairs []kv
	for k, v := range counts { pairs = append(pairs, kv{k, v}) }
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].V > pairs[j].V })
	return pairs[0].K
}

func TwoSum() int {
	nums := []int{2, 7, 11, 15}
	target := 9
	m := map[int]int{}
	for i, n := range nums {
		if j, ok := m[target-n]; ok {
			return j*10 + i
		}
		m[n] = i
	}
	return -1
}

func MergeSort() int {
	merge := func(a, b []int) []int {
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
		return result
	}
	a := []int{1, 3, 5}
	b := []int{2, 4, 6}
	merged := merge(a, b)
	return merged[0] + merged[5]
}

func SlidingWindowMax() int {
	data := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k := 3
	maxSum := 0
	for i := 0; i <= len(data)-k; i++ {
		sum := 0
		for j := 0; j < k; j++ { sum += data[i+j] }
		if sum > maxSum { maxSum = sum }
	}
	return maxSum
}

func JSONDataPipeline() int {
	type Record struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Score int    `json:"score"`
	}
	records := []Record{
		{1, "Alice", 85},
		{2, "Bob", 92},
		{3, "Charlie", 78},
	}
	data, _ := json.Marshal(records)
	var decoded []Record
	json.Unmarshal(data, &decoded)
	sort.Slice(decoded, func(i, j int) bool { return decoded[i].Score > decoded[j].Score })
	return decoded[0].Score
}

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

func WordFrequency() string {
	text := "hello world hello go hello"
	words := strings.Fields(text)
	counts := map[string]int{}
	for _, w := range words { counts[w]++ }
	return fmt.Sprintf("%d:%d:%d", counts["hello"], counts["world"], counts["go"])
}
