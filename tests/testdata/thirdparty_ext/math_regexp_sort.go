package thirdparty

import (
	"math"
	"regexp"
	"sort"
	"strconv"
)

// ============================================================================
// MATH COMPLEX OPERATIONS
// ============================================================================

// MathHypot tests math.Hypot.
func MathHypot() float64 {
	return math.Hypot(3, 4)
}

// MathAtan2 tests math.Atan2.
func MathAtan2() float64 {
	return math.Atan2(1, 1)
}

// MathModf tests math.Modf.
func MathModf() int {
	intPart, fracPart := math.Modf(3.14)
	if intPart == 3 && fracPart > 0.13 && fracPart < 0.15 {
		return 1
	}
	return 0
}

// MathIsNaNCheck tests checking for NaN.
func MathIsNaNCheck() int {
	if math.IsNaN(math.NaN()) {
		return 1
	}
	return 0
}

// MathIsInfCheck tests checking for Inf.
func MathIsInfCheck() int {
	inf := math.Inf(1)
	if math.IsInf(inf, 1) {
		return 1
	}
	return 0
}

// ============================================================================
// REGEXP ADVANCED PATTERNS
// ============================================================================

// RegexpLongestMatch tests regexp longest match.
func RegexpLongestMatch() string {
	re := regexp.MustCompile(`a(b+|c+)d`)
	re.Longest()
	return re.FindString("abbd")
}

// RegexpFindAllSubmatch tests FindAllStringSubmatch.
func RegexpFindAllSubmatch() int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllStringSubmatch("a1 b2 c3", -1)
	return len(matches)
}

// RegexpReplaceAllFunc tests ReplaceAllStringFunc.
func RegexpReplaceAllFunc() string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllStringFunc("a1 b2 c3", func(s string) string {
		n, _ := strconv.Atoi(s)
		return strconv.Itoa(n * 2)
	})
}

// ============================================================================
// SORTING WITH CUSTOM COMPARATOR
// ============================================================================

// SortWithFunc tests Sort.Slice with custom func.
func SortWithFunc() int {
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Bob", 30},
		{"Alice", 25},
		{"Eve", 35},
	}
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	return people[0].Age
}

// SortStablePreservingOrder tests Sort.SliceStable.
func SortStablePreservingOrder() int {
	type Item struct {
		Value int
		Id    int
	}
	items := []Item{
		{Value: 3, Id: 1},
		{Value: 1, Id: 2},
		{Value: 3, Id: 3},
		{Value: 2, Id: 4},
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Value < items[j].Value
	})
	if items[0].Id == 2 && items[1].Id == 4 && items[2].Id == 1 {
		return 1
	}
	return 0
}

// SortFloat64sWithNaN tests Float64s with NaN.
func SortFloat64sWithNaN() int {
	s := []float64{3.0, 1.0, math.NaN(), 2.0}
	sort.Float64s(s)
	if sort.IsSorted(sort.Float64Slice(s)) {
		return 1
	}
	return 0
}
