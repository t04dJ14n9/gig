package divergence_hunt226

import "fmt"

// ============================================================================
// Round 226: Slice append growth behavior
// ============================================================================

// SliceAppendGrowth tracks capacity growth through appends
func SliceAppendGrowth() string {
	s := make([]int, 0, 2)
	caps := []int{cap(s)}
	for i := 0; i < 20; i++ {
		s = append(s, i)
		if cap(s) != caps[len(caps)-1] {
			caps = append(caps, cap(s))
		}
	}
	return fmt.Sprintf("growth_points=%d", len(caps))
}

// SliceAppendFromNil appends to nil slice
func SliceAppendFromNil() string {
	var s []int
	s = append(s, 1)
	s = append(s, 2, 3)
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// SliceAppendManyElements appends many elements at once
func SliceAppendManyElements() string {
	s := []int{1, 2}
	s = append(s, []int{3, 4, 5, 6, 7, 8, 9, 10}...)
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// SliceAppendToFull appends to a full slice
func SliceAppendToFull() string {
	s := make([]int, 2, 2)
	s[0], s[1] = 1, 2
	originalCap := cap(s)
	s = append(s, 3)
	newCap := cap(s)
	return fmt.Sprintf("orig_cap=%d,new_cap=%d,grew=%t", originalCap, newCap, newCap > originalCap)
}

// SliceAppendEmptySlice appends empty slice
func SliceAppendEmptySlice() string {
	s := []int{1, 2, 3}
	originalLen := len(s)
	s = append(s, []int{}...)
	return fmt.Sprintf("orig_len=%d,new_len=%d", originalLen, len(s))
}

// SliceAppendByteGrowth tracks byte slice growth
func SliceAppendByteGrowth() string {
	b := make([]byte, 0, 4)
	caps := []int{cap(b)}
	for i := 0; i < 50; i++ {
		b = append(b, byte(i))
		if cap(b) != caps[len(caps)-1] {
			caps = append(caps, cap(b))
		}
	}
	return fmt.Sprintf("growth_points=%d", len(caps))
}

// SliceAppendInLoop appends in a loop pattern
func SliceAppendInLoop() string {
	var s []int
	for i := 0; i < 10; i++ {
		s = append(s, i*i)
	}
	return fmt.Sprintf("len=%d,sum=%d", len(s), sumSlice(s))
}

func sumSlice(s []int) int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

// SliceAppendStringToBytes appends string to byte slice
func SliceAppendStringToBytes() string {
	b := []byte("hello")
	b = append(b, " world"...)
	return fmt.Sprintf("len=%d", len(b))
}

// SliceAppendSelf appends slice to itself
func SliceAppendSelf() string {
	s := []int{1, 2}
	s = append(s, s...)
	return fmt.Sprintf("len=%d,vals=%d,%d,%d,%d", len(s), s[0], s[1], s[2], s[3])
}

// SliceAppendCapacityReuse checks if capacity is reused
func SliceAppendCapacityReuse() string {
	s := make([]int, 0, 10)
	s = append(s, 1, 2, 3)
	originalCap := cap(s)
	s = s[:0]
	s = append(s, 4, 5)
	return fmt.Sprintf("cap_unchanged=%t", cap(s) == originalCap)
}

// SliceAppendLargeBatch appends large batch at once
func SliceAppendLargeBatch() string {
	s := make([]int, 0)
	large := make([]int, 100)
	for i := range large {
		large[i] = i
	}
	s = append(s, large...)
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}
