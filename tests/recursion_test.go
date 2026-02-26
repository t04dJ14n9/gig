package tests

import "testing"

// Advanced recursion tests.

func TestTailRecursionPattern(t *testing.T) {
	runInt(t, `package main
func sumTail(n, acc int) int {
	if n <= 0 { return acc }
	return sumTail(n-1, acc+n)
}
func Compute() int { return sumTail(50, 0) }`, 1275)
}

func TestRecursiveReverseSlice(t *testing.T) {
	runInt(t, `package main
func reverseHelper(s []int, lo, hi int) int {
	if lo >= hi { return 0 }
	tmp := s[lo]
	s[lo] = s[hi]
	s[hi] = tmp
	return reverseHelper(s, lo+1, hi-1)
}
func Compute() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ { s[i] = i + 1 }
	_ = reverseHelper(s, 0, 4)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}`, 54321)
}

func TestTowerOfHanoi(t *testing.T) {
	// Count number of moves for Tower of Hanoi
	runInt(t, `package main
func hanoi(n int) int {
	if n == 1 { return 1 }
	return 2*hanoi(n-1) + 1
}
func Compute() int { return hanoi(10) }`, 1023)
}

func TestRecursiveMaxSlice(t *testing.T) {
	runInt(t, `package main
func maxVal(s []int, n int) int {
	if n == 1 { return s[0] }
	rest := maxVal(s, n-1)
	if s[n-1] > rest { return s[n-1] }
	return rest
}
func Compute() int {
	s := make([]int, 0)
	s = append(s, 3)
	s = append(s, 7)
	s = append(s, 1)
	s = append(s, 9)
	s = append(s, 4)
	return maxVal(s, len(s))
}`, 9)
}

func TestAckermann(t *testing.T) {
	// A(2,3) = 9 — kept small to avoid stack overflow
	runInt(t, `package main
func ack(m, n int) int {
	if m == 0 { return n + 1 }
	if n == 0 { return ack(m-1, 1) }
	return ack(m-1, ack(m, n-1))
}
func Compute() int { return ack(2, 3) }`, 9)
}

func TestRecursiveBinarySearch(t *testing.T) {
	runInt(t, `package main
func bsearch(s []int, target, lo, hi int) int {
	if lo > hi { return -1 }
	mid := (lo + hi) / 2
	if s[mid] == target { return mid }
	if s[mid] < target {
		return bsearch(s, target, mid+1, hi)
	}
	return bsearch(s, target, lo, mid-1)
}
func Compute() int {
	s := make([]int, 10)
	for i := 0; i < 10; i++ {
		s[i] = i * 5
	}
	return bsearch(s, 25, 0, 9)
}`, 5)
}
