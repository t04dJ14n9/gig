package recursion

// TailRecursionPattern tests tail recursion
func TailRecursionPattern() int { return sumTail(50, 0) }

func sumTail(n, acc int) int {
	if n <= 0 {
		return acc
	}
	return sumTail(n-1, acc+n)
}

// ReverseSlice tests recursive slice reversal
func ReverseSlice() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i + 1
	}
	_ = reverseHelper(s, 0, 4)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func reverseHelper(s []int, lo, hi int) int {
	if lo >= hi {
		return 0
	}
	tmp := s[lo]
	s[lo] = s[hi]
	s[hi] = tmp
	return reverseHelper(s, lo+1, hi-1)
}

// TowerOfHanoi tests Tower of Hanoi
func TowerOfHanoi() int { return hanoi(10) }

func hanoi(n int) int {
	if n == 1 {
		return 1
	}
	return 2*hanoi(n-1) + 1
}

// MaxSlice tests recursive max in slice
func MaxSlice() int {
	s := make([]int, 0)
	s = append(s, 3)
	s = append(s, 7)
	s = append(s, 1)
	s = append(s, 9)
	s = append(s, 4)
	return maxVal(s, len(s))
}

func maxVal(s []int, n int) int {
	if n == 1 {
		return s[0]
	}
	rest := maxVal(s, n-1)
	if s[n-1] > rest {
		return s[n-1]
	}
	return rest
}

// Ackermann tests Ackermann function
func Ackermann() int { return ack(2, 3) }

func ack(m, n int) int {
	if m == 0 {
		return n + 1
	}
	if n == 0 {
		return ack(m-1, 1)
	}
	return ack(m-1, ack(m, n-1))
}

// BinarySearch tests recursive binary search
func BinarySearch() int {
	s := make([]int, 10)
	for i := 0; i < 10; i++ {
		s[i] = i * 5
	}
	return bsearch(s, 25, 0, 9)
}

func bsearch(s []int, target, lo, hi int) int {
	if lo > hi {
		return -1
	}
	mid := (lo + hi) / 2
	if s[mid] == target {
		return mid
	}
	if s[mid] < target {
		return bsearch(s, target, mid+1, hi)
	}
	return bsearch(s, target, lo, mid-1)
}
