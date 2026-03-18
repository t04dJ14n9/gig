package algorithms

// InsertionSort tests insertion sort algorithm
func InsertionSort() int {
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
}

// SelectionSort tests selection sort algorithm
func SelectionSort() int {
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
}

// ReverseSlice tests slice reversal
func ReverseSlice() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i + 1
	}
	s = reverse(s)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

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

// IsPalindrome tests palindrome detection
func IsPalindrome() int {
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
}

func isPalindrome(s []int) int {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			return 0
		}
	}
	return 1
}

// PowerFunction tests fast exponentiation
func PowerFunction() int {
	return power(2, 10)
}

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

// MaxSubarraySum tests Kadane's algorithm
func MaxSubarraySum() int {
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
}

// TwoSum tests two sum problem
func TwoSum() int {
	nums := make([]int, 0)
	nums = append(nums, 2)
	nums = append(nums, 7)
	nums = append(nums, 11)
	nums = append(nums, 15)
	target := 9
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return i*10 + j
			}
		}
	}
	return -1
}

// FibMemoized tests memoized fibonacci
func FibMemoized() int {
	memo := make(map[int]int)
	memo[0] = 0
	memo[1] = 1
	for i := 2; i <= 30; i++ {
		memo[i] = memo[i-1] + memo[i-2]
	}
	return memo[30]
}

// CountDigits tests digit counting
func CountDigits() int {
	return countDigits(0)*1000 + countDigits(9)*100 + countDigits(99)*10 + countDigits(12345)
}

func countDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n > 0 {
		count = count + 1
		n = n / 10
	}
	return count
}

// CollatzConjecture tests Collatz sequence
func CollatzConjecture() int {
	return collatzSteps(27)
}

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

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// Reverse returns the reversed slice
func Reverse(s []int) []int { return reverse(s) }

// Power computes base^exp using fast exponentiation
func Power(base, exp int) int { return power(base, exp) }

// CountDigitsN returns the number of digits in n
func CountDigitsN(n int) int { return countDigits(n) }

// CollatzStepsN returns the number of steps to reach 1 from n
func CollatzStepsN(n int) int { return collatzSteps(n) }

// IsPalindromeInt returns 1 if s is a palindrome, 0 otherwise
func IsPalindromeInt(s []int) int { return isPalindrome(s) }
