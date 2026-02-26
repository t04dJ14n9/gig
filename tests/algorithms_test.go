package tests

import "testing"

// Additional algorithm tests for advanced language feature coverage.

func TestInsertionSort(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 0)
	s = append(s, 5)
	s = append(s, 2)
	s = append(s, 8)
	s = append(s, 1)
	s = append(s, 9)
	s = append(s, 3)
	n := len(s)
	for i := 1; i < n; i++ {
		key := s[i]
		j := i - 1
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j = j - 1
		}
		s[j+1] = key
	}
	return s[0]*100000 + s[1]*10000 + s[2]*1000 + s[3]*100 + s[4]*10 + s[5]
}`, 123589)
}

func TestSelectionSort(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 0)
	s = append(s, 4)
	s = append(s, 1)
	s = append(s, 3)
	s = append(s, 2)
	n := len(s)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if s[j] < s[minIdx] {
				minIdx = j
			}
		}
		tmp := s[i]
		s[i] = s[minIdx]
		s[minIdx] = tmp
	}
	return s[0]*1000 + s[1]*100 + s[2]*10 + s[3]
}`, 1234)
}

func TestReverseSlice(t *testing.T) {
	runInt(t, `package main
func reverse(s []int) []int {
	n := len(s)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		tmp := s[i]
		s[i] = s[j]
		s[j] = tmp
	}
	return s
}
func Compute() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i + 1
	}
	s = reverse(s)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}`, 54321)
}

func TestIsPalindrome(t *testing.T) {
	runInt(t, `package main
func isPalindrome(s []int) int {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			return 0
		}
	}
	return 1
}
func Compute() int {
	s1 := make([]int, 0)
	s1 = append(s1, 1)
	s1 = append(s1, 2)
	s1 = append(s1, 3)
	s1 = append(s1, 2)
	s1 = append(s1, 1)

	s2 := make([]int, 0)
	s2 = append(s2, 1)
	s2 = append(s2, 2)
	s2 = append(s2, 3)

	return isPalindrome(s1)*10 + isPalindrome(s2)
}`, 10)
}

func TestPowerFunction(t *testing.T) {
	runInt(t, `package main
func power(base, exp int) int {
	result := 1
	for exp > 0 {
		if exp%2 == 1 {
			result = result * base
		}
		base = base * base
		exp = exp / 2
	}
	return result
}
func Compute() int { return power(2, 10) }`, 1024)
}

func TestMaxSubarraySum(t *testing.T) {
	// Kadane's algorithm
	runInt(t, `package main
func Compute() int {
	s := make([]int, 0)
	s = append(s, -2)
	s = append(s, 1)
	s = append(s, -3)
	s = append(s, 4)
	s = append(s, -1)
	s = append(s, 2)
	s = append(s, 1)
	s = append(s, -5)
	s = append(s, 4)

	maxSoFar := s[0]
	maxEndingHere := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] > maxEndingHere+s[i] {
			maxEndingHere = s[i]
		} else {
			maxEndingHere = maxEndingHere + s[i]
		}
		if maxEndingHere > maxSoFar {
			maxSoFar = maxEndingHere
		}
	}
	return maxSoFar
}`, 6) // subarray [4, -1, 2, 1]
}

func TestTwoSum(t *testing.T) {
	// Find two indices whose values sum to target using brute force
	runInt(t, `package main
func Compute() int {
	nums := make([]int, 0)
	nums = append(nums, 2)
	nums = append(nums, 7)
	nums = append(nums, 11)
	nums = append(nums, 15)
	target := 9
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] + nums[j] == target {
				return i*10 + j
			}
		}
	}
	return -1
}`, 1) // nums[0]=2, nums[1]=7, 2+7=9 -> 0*10+1=1
}

func TestFibMemoized(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	memo := make(map[int]int)
	memo[0] = 0
	memo[1] = 1
	// iterative fill for memoization pattern
	for i := 2; i <= 30; i++ {
		memo[i] = memo[i-1] + memo[i-2]
	}
	return memo[30]
}`, 832040)
}

func TestCountDigits(t *testing.T) {
	runInt(t, `package main
func countDigits(n int) int {
	if n == 0 { return 1 }
	count := 0
	for n > 0 {
		count = count + 1
		n = n / 10
	}
	return count
}
func Compute() int {
	return countDigits(0)*1000 + countDigits(9)*100 + countDigits(99)*10 + countDigits(12345)
}`, 1125)
}

func TestCollatzConjecture(t *testing.T) {
	runInt(t, `package main
func collatzSteps(n int) int {
	steps := 0
	for n != 1 {
		if n%2 == 0 {
			n = n / 2
		} else {
			n = 3*n + 1
		}
		steps = steps + 1
	}
	return steps
}
func Compute() int { return collatzSteps(27) }`, 111)
}
