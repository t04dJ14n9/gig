package divergence_hunt10

import "fmt"

// ============================================================================
// Round 10: Complex data structures, algorithm patterns, edge cases in type
// system, fmt edge cases
// ============================================================================

// BinarySearch tests binary search pattern
func BinarySearch() int {
	s := []int{1, 3, 5, 7, 9, 11, 13}
	target := 7
	lo, hi := 0, len(s)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if s[mid] == target { return mid }
		if s[mid] < target { lo = mid + 1 } else { hi = mid - 1 }
	}
	return -1
}

// StackPattern tests stack via slice
func StackPattern() int {
	stack := []int{}
	push := func(v int) { stack = append(stack, v) }
	pop := func() int {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v
	}
	push(1)
	push(2)
	push(3)
	return pop() * 100 + pop() * 10 + pop()
}

// QueuePattern tests queue via slice
func QueuePattern() int {
	queue := []int{}
	enqueue := func(v int) { queue = append(queue, v) }
	dequeue := func() int {
		v := queue[0]
		queue = queue[1:]
		return v
	}
	enqueue(1)
	enqueue(2)
	enqueue(3)
	return dequeue()*100 + dequeue()*10 + dequeue()
}

// TwoSum tests two-sum pattern
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

// IsPalindrome tests palindrome check
func IsPalindrome() bool {
	s := "racecar"
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] { return false }
	}
	return true
}

// FizzBuzz tests FizzBuzz pattern
func FizzBuzz() int {
	count := 0
	for i := 1; i <= 15; i++ {
		if i%3 == 0 && i%5 == 0 {
			count += 3
		} else if i%3 == 0 {
			count += 1
		} else if i%5 == 0 {
			count += 2
		}
	}
	return count
}

// FmtVerb tests various fmt verbs
func FmtVerb() string {
	return fmt.Sprintf("%d %s %t %f", 42, "hi", true, 3.14)
}

// FmtWidthPrecision tests fmt width/precision
func FmtWidthPrecision() string {
	return fmt.Sprintf("|%5d|%-5d|%.2f|%10s|", 42, 42, 3.14159, "hi")
}

// NestedMapLookup tests nested map lookup with ok
func NestedMapLookup() int {
	m := map[string]map[string]int{
		"a": {"x": 1, "y": 2},
		"b": {"x": 3},
	}
	if inner, ok := m["a"]; ok {
		if v, ok2 := inner["y"]; ok2 {
			return v
		}
	}
	return -1
}

// StructSliceFilter tests filtering structs
func StructSliceFilter() int {
	type Item struct {
		Name  string
		Value int
	}
	items := []Item{
		{"a", 10},
		{"b", 20},
		{"c", 30},
		{"d", 5},
	}
	var filtered []Item
	for _, item := range items {
		if item.Value >= 20 {
			filtered = append(filtered, item)
		}
	}
	return len(filtered)
}

// GCD computes greatest common divisor
func GCD() int {
	a, b := 48, 18
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM computes least common multiple
func LCM() int {
	gcd := func(a, b int) int {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}
	a, b := 12, 18
	return a * b / gcd(a, b)
}

// Power tests power calculation
func Power() int {
	result := 1
	base, exp := 2, 10
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

// CountDigits tests digit counting
func CountDigits() int {
	n := 12345
	count := 0
	for n > 0 {
		n /= 10
		count++
	}
	return count
}

// ReverseInt tests integer reversal
func ReverseInt() int {
	n := 12345
	result := 0
	for n > 0 {
		result = result*10 + n%10
		n /= 10
	}
	return result
}

// FibIterative tests iterative fibonacci
func FibIterative() int {
	a, b := 0, 1
	for i := 0; i < 10; i++ {
		a, b = b, a+b
	}
	return a
}

// PrimeCheck tests prime checking
func PrimeCheck() bool {
	n := 17
	if n < 2 { return false }
	for i := 2; i*i <= n; i++ {
		if n%i == 0 { return false }
	}
	return true
}

// FactorialIterative tests iterative factorial
func FactorialIterative() int {
	result := 1
	for i := 2; i <= 10; i++ {
		result *= i
	}
	return result
}

// CountingSort tests counting sort
func CountingSort() int {
	s := []int{5, 3, 1, 4, 2, 5, 3}
	count := make([]int, 6)
	for _, v := range s { count[v]++ }
	sorted := []int{}
	for i := 1; i <= 5; i++ {
		for j := 0; j < count[i]; j++ {
			sorted = append(sorted, i)
		}
	}
	return sorted[0]*100000 + sorted[1]*10000 + sorted[2]*1000 + sorted[3]*100 + sorted[4]*10 + sorted[5]
}

// PrefixSum tests prefix sum
func PrefixSum() int {
	s := []int{1, 2, 3, 4, 5}
	prefix := make([]int, len(s)+1)
	for i, v := range s { prefix[i+1] = prefix[i] + v }
	return prefix[5]
}

// StringAnagram tests anagram check
func StringAnagram() bool {
	a, b := "listen", "silent"
	if len(a) != len(b) { return false }
	count := map[rune]int{}
	for _, c := range a { count[c]++ }
	for _, c := range b { count[c]-- }
	for _, v := range count {
		if v != 0 { return false }
	}
	return true
}
