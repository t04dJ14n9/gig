package tests

import "testing"

// Tests for slice sub-slicing operations s[lo:hi] and s[lo:hi:max].

func TestSubSliceBasic(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = (i + 1) * 10
	}
	sub := s[1:4]
	return sub[0] + sub[1] + sub[2]
}`, 90) // 20+30+40
}

func TestSubSliceLen(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 10)
	return len(s[2:7])
}`, 5)
}

func TestSubSliceFromStart(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i
	}
	sub := s[:3]
	sum := 0
	for _, v := range sub {
		sum = sum + v
	}
	return sum
}`, 3) // 0+1+2
}

func TestSubSliceToEnd(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i
	}
	sub := s[3:]
	sum := 0
	for _, v := range sub {
		sum = sum + v
	}
	return sum
}`, 7) // 3+4
}

func TestSubSliceCopy(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i
	}
	sub := s[:]
	return len(sub)
}`, 5)
}

func TestSubSliceChained(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 10)
	for i := 0; i < 10; i++ {
		s[i] = i
	}
	sub := s[2:8]
	sub2 := sub[1:4]
	return sub2[0] + sub2[1] + sub2[2]
}`, 12) // 3+4+5
}

func TestSubSliceModifiesOriginal(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i
	}
	sub := s[1:4]
	sub[0] = 99
	return s[1]
}`, 99)
}
