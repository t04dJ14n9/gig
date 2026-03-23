package complex_tests

// ============================================================================
// Complex Recursion Tests
// ============================================================================

// RecursionFib tests Fibonacci recursion.
func RecursionFib(n int) int {
	if n <= 1 {
		return n
	}
	return RecursionFib(n-1) + RecursionFib(n-2)
}

// RecursionFibCheck tests Fibonacci.
func RecursionFibCheck() int {
	return RecursionFib(15)
}

// RecursionAckermann tests Ackermann function.
func RecursionAckermann(m, n int) int {
	if m == 0 {
		return n + 1
	}
	if n == 0 {
		return RecursionAckermann(m-1, 1)
	}
	return RecursionAckermann(m-1, RecursionAckermann(m, n-1))
}

// RecursionAckermannCheck tests Ackermann.
func RecursionAckermannCheck() int {
	return RecursionAckermann(3, 4)
}

// RecursionGCD tests GCD via recursion.
func RecursionGCD(a, b int) int {
	if b == 0 {
		return a
	}
	return RecursionGCD(b, a%b)
}

// RecursionGCDCheck tests GCD.
func RecursionGCDCheck() int {
	return RecursionGCD(48, 18)
}

// RecursionSum tests sum via recursion.
func RecursionSum(n int) int {
	if n <= 0 {
		return 0
	}
	return n + RecursionSum(n-1)
}

// RecursionSumCheck tests sum recursion.
func RecursionSumCheck() int {
	return RecursionSum(100)
}

// RecursionPower tests power via recursion.
func RecursionPower(base, exp int) int {
	if exp == 0 {
		return 1
	}
	if exp%2 == 0 {
		half := RecursionPower(base, exp/2)
		return half * half
	}
	return base * RecursionPower(base, exp-1)
}

// RecursionPowerCheck tests power.
func RecursionPowerCheck() int {
	return RecursionPower(2, 10)
}

// RecursionBinarySearch tests binary search recursion.
func RecursionBinarySearch(arr []int, target, low, high int) int {
	if low > high {
		return -1
	}
	mid := (low + high) / 2
	if arr[mid] == target {
		return mid
	}
	if arr[mid] > target {
		return RecursionBinarySearch(arr, target, low, mid-1)
	}
	return RecursionBinarySearch(arr, target, mid+1, high)
}

// RecursionBinarySearchCheck tests binary search.
func RecursionBinarySearchCheck() int {
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	return RecursionBinarySearch(arr, 11, 0, len(arr)-1)
}

// RecursionReverse tests string reversal via recursion.
func RecursionReverse(s string) string {
	if len(s) <= 1 {
		return s
	}
	return RecursionReverse(s[1:]) + string(s[0])
}

// RecursionReverseCheck tests reversal.
func RecursionReverseCheck() int {
	r := RecursionReverse("hello")
	if r == "olleh" {
		return 1
	}
	return 0
}

// RecursionPalindrome tests palindrome check via recursion.
func RecursionPalindrome(s string) bool {
	if len(s) <= 1 {
		return true
	}
	if s[0] != s[len(s)-1] {
		return false
	}
	return RecursionPalindrome(s[1 : len(s)-1])
}

// RecursionPalindromeCheck tests palindrome.
func RecursionPalindromeCheck() int {
	if RecursionPalindrome("racecar") && !RecursionPalindrome("hello") {
		return 1
	}
	return 0
}

// RecursionMergeSort tests merge sort recursion.
func RecursionMergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := RecursionMergeSort(arr[:mid])
	right := RecursionMergeSort(arr[mid:])
	return RecursionMerge(left, right)
}

// RecursionMerge merges two sorted arrays.
func RecursionMerge(left, right []int) []int {
	result := []int{}
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// RecursionMergeSortCheck tests merge sort.
func RecursionMergeSortCheck() int {
	arr := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	sorted := RecursionMergeSort(arr)
	for i := 1; i < len(sorted); i++ {
		if sorted[i] < sorted[i-1] {
			return 0
		}
	}
	return 1
}

// RecursionQuickSort tests quick sort recursion.
func RecursionQuickSort(arr []int, low, high int) {
	if low < high {
		pivot := RecursionPartition(arr, low, high)
		RecursionQuickSort(arr, low, pivot-1)
		RecursionQuickSort(arr, pivot+1, high)
	}
}

// RecursionPartition partitions array for quick sort.
func RecursionPartition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// RecursionQuickSortCheck tests quick sort.
func RecursionQuickSortCheck() int {
	arr := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	RecursionQuickSort(arr, 0, len(arr)-1)
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return 0
		}
	}
	return 1
}

// RecursionCountDigits tests digit counting via recursion.
func RecursionCountDigits(n int) int {
	if n < 10 {
		return 1
	}
	return 1 + RecursionCountDigits(n/10)
}

// RecursionCountDigitsCheck tests digit counting.
func RecursionCountDigitsCheck() int {
	return RecursionCountDigits(12345)
}

// RecursionSumDigits tests digit sum via recursion.
func RecursionSumDigits(n int) int {
	if n == 0 {
		return 0
	}
	return n%10 + RecursionSumDigits(n/10)
}

// RecursionSumDigitsCheck tests digit sum.
func RecursionSumDigitsCheck() int {
	return RecursionSumDigits(12345)
}

// RecursionHanoi tests Tower of Hanoi move count.
func RecursionHanoi(n int) int {
	if n == 1 {
		return 1
	}
	return 2*RecursionHanoi(n-1) + 1
}

// RecursionHanoiCheck tests Hanoi.
func RecursionHanoiCheck() int {
	return RecursionHanoi(10)
}

// RecursionStaircase tests ways to climb staircase.
func RecursionStaircase(n int) int {
	if n <= 1 {
		return 1
	}
	return RecursionStaircase(n-1) + RecursionStaircase(n-2)
}

// RecursionStaircaseCheck tests staircase.
func RecursionStaircaseCheck() int {
	return RecursionStaircase(10)
}

// RecursionSubsetSum tests subset sum recursion.
func RecursionSubsetSum(arr []int, target, index int) bool {
	if target == 0 {
		return true
	}
	if index >= len(arr) || target < 0 {
		return false
	}
	return RecursionSubsetSum(arr, target-arr[index], index+1) ||
		RecursionSubsetSum(arr, target, index+1)
}

// RecursionSubsetSumCheck tests subset sum.
func RecursionSubsetSumCheck() int {
	arr := []int{3, 34, 4, 12, 5, 2}
	if RecursionSubsetSum(arr, 9, 0) {
		return 1
	}
	return 0
}

// RecursionPermuteCount counts permutations.
func RecursionPermuteCount(n int) int {
	if n <= 1 {
		return 1
	}
	return n * RecursionPermuteCount(n-1)
}

// RecursionPermuteCheck tests permutation count.
func RecursionPermuteCheck() int {
	return RecursionPermuteCount(7)
}

// RecursionCombination counts combinations.
func RecursionCombination(n, r int) int {
	if r == 0 || r == n {
		return 1
	}
	return RecursionCombination(n-1, r-1) + RecursionCombination(n-1, r)
}

// RecursionCombinationCheck tests combinations.
func RecursionCombinationCheck() int {
	return RecursionCombination(10, 4)
}
