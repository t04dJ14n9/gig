package divergence_hunt161

import "fmt"

// ============================================================================
// Round 161: Recursive function patterns (factorial, fibonacci, tree traversal)
// ============================================================================

// RecursiveFactorial calculates factorial recursively
func RecursiveFactorial() string {
	var fact func(n int) int
	fact = func(n int) int {
		if n <= 1 {
			return 1
		}
		return n * fact(n-1)
	}
	return fmt.Sprintf("5!=%d,7!=%d", fact(5), fact(7))
}

// RecursiveFibonacci calculates fibonacci recursively
func RecursiveFibonacci() string {
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}
	return fmt.Sprintf("fib(10)=%d,fib(15)=%d", fib(10), fib(15))
}

// RecursiveSum sums a slice recursively
func RecursiveSum() string {
	var sum func([]int, int) int
	sum = func(arr []int, idx int) int {
		if idx >= len(arr) {
			return 0
		}
		return arr[idx] + sum(arr, idx+1)
	}
	arr := []int{1, 2, 3, 4, 5}
	return fmt.Sprintf("sum=%d", sum(arr, 0))
}

// RecursiveMax finds maximum recursively
func RecursiveMax() string {
	var max func([]int, int) int
	max = func(arr []int, idx int) int {
		if idx == len(arr)-1 {
			return arr[idx]
		}
		next := max(arr, idx+1)
		if arr[idx] > next {
			return arr[idx]
		}
		return next
	}
	arr := []int{3, 7, 2, 9, 1}
	return fmt.Sprintf("max=%d", max(arr, 0))
}

// RecursiveReverse reverses a string recursively
func RecursiveReverse() string {
	var rev func(string) string
	rev = func(s string) string {
		if len(s) <= 1 {
			return s
		}
		return rev(s[1:]) + s[:1]
	}
	return fmt.Sprintf("%s->%s", "hello", rev("hello"))
}

// RecursiveTreeTraversal tests tree traversal
func RecursiveTreeTraversal() string {
	type TreeNode struct {
		Val   int
		Left  *TreeNode
		Right *TreeNode
	}
	var inorder func(*TreeNode, []int) []int
	inorder = func(node *TreeNode, acc []int) []int {
		if node == nil {
			return acc
		}
		acc = inorder(node.Left, acc)
		acc = append(acc, node.Val)
		acc = inorder(node.Right, acc)
		return acc
	}
	root := &TreeNode{Val: 4,
		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 1}, Right: &TreeNode{Val: 3}},
		Right: &TreeNode{Val: 6, Left: &TreeNode{Val: 5}, Right: &TreeNode{Val: 7}},
	}
	result := inorder(root, []int{})
	return fmt.Sprintf("inorder=%v", result)
}

// RecursivePower calculates power recursively
func RecursivePower() string {
	var power func(base, exp int) int
	power = func(base, exp int) int {
		if exp == 0 {
			return 1
		}
		return base * power(base, exp-1)
	}
	return fmt.Sprintf("2^10=%d,3^5=%d", power(2, 10), power(3, 5))
}

// RecursiveGCD calculates GCD using Euclidean algorithm
func RecursiveGCD() string {
	var gcd func(a, b int) int
	gcd = func(a, b int) int {
		if b == 0 {
			return a
		}
		return gcd(b, a%b)
	}
	return fmt.Sprintf("gcd(48,18)=%d,gcd(100,35)=%d", gcd(48, 18), gcd(100, 35))
}

// RecursivePalindrome checks if string is palindrome
func RecursivePalindrome() string {
	var isPal func(string, int, int) bool
	isPal = func(s string, left, right int) bool {
		if left >= right {
			return true
		}
		if s[left] != s[right] {
			return false
		}
		return isPal(s, left+1, right-1)
	}
	return fmt.Sprintf("radar=%t,hello=%t", isPal("radar", 0, 4), isPal("hello", 0, 4))
}

// RecursiveBinarySearch implements binary search recursively
func RecursiveBinarySearch() string {
	var search func([]int, int, int, int) int
	search = func(arr []int, target, left, right int) int {
		if left > right {
			return -1
		}
		mid := (left + right) / 2
		if arr[mid] == target {
			return mid
		}
		if arr[mid] > target {
			return search(arr, target, left, mid-1)
		}
		return search(arr, target, mid+1, right)
	}
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15}
	return fmt.Sprintf("found(7)=%d,found(4)=%d", search(arr, 7, 0, len(arr)-1), search(arr, 4, 0, len(arr)-1))
}

// RecursiveMergeSort implements merge sort
func RecursiveMergeSort() string {
	var mergeSort func([]int) []int
	var merge func([]int, []int) []int
	merge = func(left, right []int) []int {
		result := []int{}
		i, j := 0, 0
		for i < len(left) && j < len(right) {
			if left[i] < right[j] {
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
	mergeSort = func(arr []int) []int {
		if len(arr) <= 1 {
			return arr
		}
		mid := len(arr) / 2
		left := mergeSort(arr[:mid])
		right := mergeSort(arr[mid:])
		return merge(left, right)
	}
	arr := []int{64, 34, 25, 12, 22, 11, 90}
	sorted := mergeSort(arr)
	return fmt.Sprintf("sorted=%v", sorted)
}

// TailRecursiveFactorial uses tail recursion with accumulator
func TailRecursiveFactorial() string {
	var factHelper func(n, acc int) int
	factHelper = func(n, acc int) int {
		if n <= 1 {
			return acc
		}
		return factHelper(n-1, n*acc)
	}
	return fmt.Sprintf("5!=%d,10!=%d", factHelper(5, 1), factHelper(10, 1))
}
