package divergence_hunt163

import "fmt"

// ============================================================================
// Round 163: Type parameters edge cases and generic-like patterns
// ============================================================================

// MinInt finds minimum of two ints
func MinInt() string {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	return fmt.Sprintf("min(3,5)=%d,min(10,7)=%d", min(3, 5), min(10, 7))
}

// MinFloat64 finds minimum of two float64s
func MinFloat64() string {
	min := func(a, b float64) float64 {
		if a < b {
			return a
		}
		return b
	}
	return fmt.Sprintf("min(3.14,2.71)=%.2f,min(1.5,2.5)=%.2f", min(3.14, 2.71), min(1.5, 2.5))
}

// MinString finds minimum of two strings (lexicographically)
func MinString() string {
	min := func(a, b string) string {
		if a < b {
			return a
		}
		return b
	}
	return fmt.Sprintf("min(apple,banana)=%s,min(zebra,apple)=%s", min("apple", "banana"), min("zebra", "apple"))
}

// GenericStackPattern simulates a generic stack for int
func GenericStackPattern() string {
	type IntStack struct {
		items []int
	}
	push := func(s *IntStack, v int) {
		s.items = append(s.items, v)
	}
	pop := func(s *IntStack) (int, bool) {
		if len(s.items) == 0 {
			return 0, false
		}
		v := s.items[len(s.items)-1]
		s.items = s.items[:len(s.items)-1]
		return v, true
	}
	s := IntStack{}
	push(&s, 10)
	push(&s, 20)
	push(&s, 30)
	v1, _ := pop(&s)
	v2, _ := pop(&s)
	return fmt.Sprintf("pop1=%d,pop2=%d,size=%d", v1, v2, len(s.items))
}

// GenericStackPatternString simulates a generic stack for string
func GenericStackPatternString() string {
	type StringStack struct {
		items []string
	}
	push := func(s *StringStack, v string) {
		s.items = append(s.items, v)
	}
	pop := func(s *StringStack) (string, bool) {
		if len(s.items) == 0 {
			return "", false
		}
		v := s.items[len(s.items)-1]
		s.items = s.items[:len(s.items)-1]
		return v, true
	}
	s := StringStack{}
	push(&s, "hello")
	push(&s, "world")
	v1, _ := pop(&s)
	v2, _ := pop(&s)
	return fmt.Sprintf("pop1=%s,pop2=%s,size=%d", v1, v2, len(s.items))
}

// GenericMapPattern simulates generic map functions
func GenericMapPattern() string {
	// Map int to int
	mapInt := func(arr []int, fn func(int) int) []int {
		result := make([]int, len(arr))
		for i, v := range arr {
			result[i] = fn(v)
		}
		return result
	}
	arr := []int{1, 2, 3, 4, 5}
	doubled := mapInt(arr, func(x int) int { return x * 2 })
	return fmt.Sprintf("doubled=%v", doubled)
}

// GenericFilterPattern simulates generic filter
func GenericFilterPattern() string {
	filterInt := func(arr []int, fn func(int) bool) []int {
		result := []int{}
		for _, v := range arr {
			if fn(v) {
				result = append(result, v)
			}
		}
		return result
	}
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := filterInt(arr, func(x int) bool { return x%2 == 0 })
	return fmt.Sprintf("evens=%v", evens)
}

// GenericReducePattern simulates generic reduce
func GenericReducePattern() string {
	reduceInt := func(arr []int, fn func(int, int) int, init int) int {
		result := init
		for _, v := range arr {
			result = fn(result, v)
		}
		return result
	}
	arr := []int{1, 2, 3, 4, 5}
	sum := reduceInt(arr, func(a, b int) int { return a + b }, 0)
	product := reduceInt(arr, func(a, b int) int { return a * b }, 1)
	return fmt.Sprintf("sum=%d,product=%d", sum, product)
}

// ComparableConstraintPattern simulates comparable constraint
func ComparableConstraintPattern() string {
	findIndex := func(arr []int, target int) int {
		for i, v := range arr {
			if v == target {
				return i
			}
		}
		return -1
	}
	arr := []int{10, 20, 30, 40, 50}
	return fmt.Sprintf("find(30)=%d,find(99)=%d", findIndex(arr, 30), findIndex(arr, 99))
}

// OrderedConstraintPattern simulates ordered constraint
func OrderedConstraintPattern() string {
	isSorted := func(arr []int) bool {
		for i := 1; i < len(arr); i++ {
			if arr[i] < arr[i-1] {
				return false
			}
		}
		return true
	}
	sorted := []int{1, 2, 3, 4, 5}
	unsorted := []int{1, 3, 2, 4, 5}
	return fmt.Sprintf("sorted=%t,unsorted=%t", isSorted(sorted), isSorted(unsorted))
}

// GenericCachePattern simulates generic cache pattern
func GenericCachePattern() string {
	type IntCache struct {
		data map[string]int
	}
	get := func(c *IntCache, key string, compute func() int) int {
		if c.data == nil {
			c.data = make(map[string]int)
		}
		if v, ok := c.data[key]; ok {
			return v
		}
		v := compute()
		c.data[key] = v
		return v
	}
	cache := IntCache{}
	counter := 0
	compute := func() int {
		counter++
		return counter * 10
	}
	v1 := get(&cache, "a", compute)
	v2 := get(&cache, "a", compute)
	v3 := get(&cache, "b", compute)
	return fmt.Sprintf("v1=%d,v2=%d,v3=%d,counter=%d", v1, v2, v3, counter)
}
