package tests

import "testing"

func TestSliceMakeLen(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	nums := make([]int, 3)
	return len(nums)
}`, 3)
}

func TestSliceAppend(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 0)
	s = append(s, 10)
	s = append(s, 20)
	s = append(s, 30)
	return s[0] + s[1] + s[2] + len(s)
}`, 63)
}

func TestSliceElementAssignment(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 3)
	s[0] = 100
	s[1] = 200
	s[2] = 300
	return s[0] + s[1] + s[2]
}`, 600)
}

func TestSliceForRange(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	nums := make([]int, 0)
	nums = append(nums, 10)
	nums = append(nums, 20)
	nums = append(nums, 30)
	sum := 0
	for _, v := range nums {
		sum = sum + v
	}
	return sum
}`, 60)
}

func TestSliceForRangeIndex(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 3)
	s[0] = 10
	s[1] = 20
	s[2] = 30
	sum := 0
	for i, v := range s {
		sum = sum + i*100 + v
	}
	return sum
}`, 360)
}

func TestSliceGrowMultiple(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 0)
	for i := 0; i < 20; i++ {
		s = append(s, i)
	}
	sum := 0
	for _, v := range s {
		sum = sum + v
	}
	return sum
}`, 190)
}

func TestSlicePassToFunction(t *testing.T) {
	runInt(t, `package main
func sumSlice(s []int) int {
	total := 0
	for _, v := range s {
		total = total + v
	}
	return total
}
func Compute() int {
	s := make([]int, 0)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	return sumSlice(s)
}`, 6)
}

func TestSliceLenCap(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 3, 10)
	return len(s)*100 + cap(s)
}`, 310)
}
