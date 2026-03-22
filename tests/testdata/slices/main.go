package slices

// MakeLen tests slice make with length
func MakeLen() int {
	nums := make([]int, 3)
	return len(nums)
}

// Append tests slice append
func Append() int {
	s := make([]int, 0)
	s = append(s, 10)
	s = append(s, 20)
	s = append(s, 30)
	return s[0] + s[1] + s[2] + len(s)
}

// ElementAssignment tests slice element assignment
func ElementAssignment() int {
	s := make([]int, 3)
	s[0] = 100
	s[1] = 200
	s[2] = 300
	return s[0] + s[1] + s[2]
}

// ForRange tests slice for range
func ForRange() int {
	nums := make([]int, 0)
	nums = append(nums, 10)
	nums = append(nums, 20)
	nums = append(nums, 30)
	sum := 0
	for _, v := range nums {
		sum = sum + v
	}
	return sum
}

// ForRangeIndex tests slice for range with index
func ForRangeIndex() int {
	s := make([]int, 3)
	s[0] = 10
	s[1] = 20
	s[2] = 30
	sum := 0
	for i, v := range s {
		sum = sum + i*100 + v
	}
	return sum
}

// GrowMultiple tests slice grow multiple times
func GrowMultiple() int {
	s := make([]int, 0)
	for i := 0; i < 20; i++ {
		s = append(s, i)
	}
	sum := 0
	for _, v := range s {
		sum = sum + v
	}
	return sum
}

// PassToFunction tests passing slice to function
func PassToFunction() int {
	s := make([]int, 0)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	return sumSlice(s)
}

func sumSlice(s []int) int {
	total := 0
	for _, v := range s {
		total = total + v
	}
	return total
}

// LenCap tests slice len and cap
func LenCap() int {
	s := make([]int, 3, 10)
	return len(s)*100 + cap(s)
}

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// SumSlice returns the sum of all elements in s
func SumSlice(s []int) int { return sumSlice(s) }
