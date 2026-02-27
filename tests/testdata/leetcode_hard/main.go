package leetcode_hard

// Classic LeetCode Hard Problems for gig interpreter testing.
// These problems demonstrate advanced algorithmic patterns including
// dynamic programming, backtracking, and graph algorithms.

// Problem 1: Trapping Rain Water
// Given n non-negative integers representing an elevation map where the width of each bar is 1,
// compute how much water it can trap after raining.
func TrappingRainWater() int {
	height := []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}
	n := len(height)
	if n == 0 {
		return 0
	}

	leftMax := make([]int, n)
	rightMax := make([]int, n)

	leftMax[0] = height[0]
	for i := 1; i < n; i++ {
		if height[i] > leftMax[i-1] {
			leftMax[i] = height[i]
		} else {
			leftMax[i] = leftMax[i-1]
		}
	}

	rightMax[n-1] = height[n-1]
	for i := n - 2; i >= 0; i-- {
		if height[i] > rightMax[i+1] {
			rightMax[i] = height[i]
		} else {
			rightMax[i] = rightMax[i+1]
		}
	}

	water := 0
	for i := 0; i < n; i++ {
		minHeight := leftMax[i]
		if rightMax[i] < minHeight {
			minHeight = rightMax[i]
		}
		water += minHeight - height[i]
	}
	return water
}

// Problem 2: Largest Rectangle in Histogram
// Given an array of integers heights representing the histogram's bar height where the width of each bar is 1,
// return the area of the largest rectangle in the histogram.
func LargestRectangleInHistogram() int {
	heights := []int{2, 1, 5, 6, 2, 3}
	n := len(heights)
	if n == 0 {
		return 0
	}

	maxArea := 0
	for i := 0; i < n; i++ {
		minHeight := heights[i]
		width := 1
		area := minHeight * width
		if area > maxArea {
			maxArea = area
		}

		for j := i + 1; j < n; j++ {
			if heights[j] < minHeight {
				minHeight = heights[j]
			}
			width = j - i + 1
			area = minHeight * width
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

// Problem 3: Median of Two Sorted Arrays
// Given two sorted arrays nums1 and nums2 of size m and n respectively,
// return the median of the two sorted arrays.
func MedianOfTwoSortedArrays() int {
	nums1 := []int{1, 3}
	nums2 := []int{2}

	// Merge arrays
	merged := make([]int, 0)
	i, j := 0, 0
	for i < len(nums1) && j < len(nums2) {
		if nums1[i] < nums2[j] {
			merged = append(merged, nums1[i])
			i++
		} else {
			merged = append(merged, nums2[j])
			j++
		}
	}
	for i < len(nums1) {
		merged = append(merged, nums1[i])
		i++
	}
	for j < len(nums2) {
		merged = append(merged, nums2[j])
		j++
	}

	n := len(merged)
	if n%2 == 1 {
		return merged[n/2]
	}
	return (merged[n/2-1] + merged[n/2]) / 2
}

// Problem 4: Regular Expression Matching (simplified)
// Check if string matches pattern with '.' and '*'
func RegularExpressionMatching() int {
	s := "aab"
	p := "c*a*b"

	// Simplified DP approach
	m, n := len(s), len(p)
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}

	dp[0][0] = true
	for j := 2; j <= n; j += 2 {
		if p[j-1] == '*' {
			dp[0][j] = dp[0][j-2]
		}
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if p[j-1] == '*' {
				dp[i][j] = dp[i][j-2]
				if p[j-2] == '.' || p[j-2] == s[i-1] {
					dp[i][j] = dp[i][j] || dp[i-1][j]
				}
			} else if p[j-1] == '.' || p[j-1] == s[i-1] {
				dp[i][j] = dp[i-1][j-1]
			}
		}
	}

	if dp[m][n] {
		return 1
	}
	return 0
}

// Problem 5: N-Queens (count solutions)
// Return the number of distinct solutions to the n-queens puzzle.
func NQueens() int {
	n := 8
	count := 0

	// Use iterative backtracking with stack
	queens := make([]int, n)
	for i := range queens {
		queens[i] = -1
	}

	row := 0
	for row >= 0 {
		// Try next column in current row
		queens[row]++

		// Check if we found valid position
		valid := false
		for queens[row] < n {
			col := queens[row]
			// Check if this position is valid
			conflict := false
			for r := 0; r < row; r++ {
				c := queens[r]
				if c == col || c-r == col-row || c+r == col+row {
					conflict = true
					break
				}
			}
			if !conflict {
				valid = true
				break
			}
			queens[row]++
		}

		if valid {
			if row == n-1 {
				// Found a solution
				count++
				// Continue searching in this row
				queens[row]++
			} else {
				// Move to next row
				row++
				queens[row] = -1
			}
		} else {
			// Backtrack
			row--
		}
	}

	return count
}

// Problem 6: Longest Increasing Path in Matrix
// Given an m x n integers matrix, return the length of the longest increasing path in matrix.
func LongestIncreasingPath() int {
	matrix := [][]int{
		{9, 9, 4},
		{6, 6, 8},
		{2, 1, 1},
	}

	m := len(matrix)
	n := len(matrix[0])
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n)
	}

	maxPath := 0

	// Use iterative DP instead of recursive closure
	for iteration := 0; iteration < m*n; iteration++ {
		updated := false
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				// Check 4 directions
				maxLen := 1
				dirs := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
				for _, d := range dirs {
					ni := i + d[0]
					nj := j + d[1]
					if ni >= 0 && ni < m && nj >= 0 && nj < n && matrix[ni][nj] < matrix[i][j] {
						if memo[ni][nj]+1 > maxLen {
							maxLen = memo[ni][nj] + 1
						}
					}
				}
				if maxLen > memo[i][j] {
					memo[i][j] = maxLen
					updated = true
					if maxLen > maxPath {
						maxPath = maxLen
					}
				}
			}
		}
		if !updated {
			break
		}
	}

	return maxPath
}

// Problem 7: Word Ladder
// Given two words and a word list, find the length of shortest transformation sequence.
func WordLadder() int {
	beginWord := "hit"
	endWord := "cog"
	wordList := []string{"hot", "dot", "dog", "lot", "log", "cog"}

	// BFS approach
	wordSet := make(map[string]bool)
	for _, w := range wordList {
		wordSet[w] = true
	}

	if !wordSet[endWord] {
		return 0
	}

	visited := make(map[string]bool)
	queue := []string{beginWord}
	visited[beginWord] = true
	level := 1

	// Helper: check if two words differ by one character
	oneDiff := func(a, b string) bool {
		diff := 0
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				diff++
			}
		}
		return diff == 1
	}

	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			word := queue[0]
			queue = queue[1:]

			if word == endWord {
				return level
			}

			// Try all words in wordSet
			for nextWord := range wordSet {
				if !visited[nextWord] && oneDiff(word, nextWord) {
					visited[nextWord] = true
					queue = append(queue, nextWord)
				}
			}
		}
		level++
	}

	return 0
}

// Problem 8: Merge k Sorted Lists (simulated with slices)
// Merge all sorted slices into one sorted slice and return median.
func MergeKSortedLists() int {
	lists := [][]int{
		{1, 4, 5},
		{1, 3, 4},
		{2, 6},
	}

	// Simple merge approach
	result := make([]int, 0)
	indices := make([]int, len(lists))

	for {
		minVal := -1
		minIdx := -1

		for i, lst := range lists {
			if indices[i] < len(lst) {
				if minIdx == -1 || lst[indices[i]] < minVal {
					minVal = lst[indices[i]]
					minIdx = i
				}
			}
		}

		if minIdx == -1 {
			break
		}

		result = append(result, minVal)
		indices[minIdx]++
	}

	// Return sum as result
	sum := 0
	for _, v := range result {
		sum += v
	}
	return sum
}

// Problem 9: Edit Distance
// Given two strings word1 and word2, return the minimum number of operations required to convert word1 to word2.
func EditDistance() int {
	word1 := "horse"
	word2 := "ros"

	m, n := len(word1), len(word2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				// min of insert, delete, replace
				minOps := dp[i][j-1]
				if dp[i-1][j] < minOps {
					minOps = dp[i-1][j]
				}
				if dp[i-1][j-1] < minOps {
					minOps = dp[i-1][j-1]
				}
				dp[i][j] = 1 + minOps
			}
		}
	}

	return dp[m][n]
}

// Problem 10: Minimum Window Substring
// Given two strings s and t, return the minimum window substring of s such that every character in t is included.
// Return the length of the minimum window.
func MinimumWindowSubstring() int {
	s := "ADOBECODEBANC"
	t := "ABC"

	if len(s) < len(t) {
		return 0
	}

	// Count characters in t
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}

	have := make(map[byte]int)
	required := len(need)
	formed := 0
	minLen := len(s) + 1

	left := 0
	for right := 0; right < len(s); right++ {
		char := s[right]
		have[char]++

		if need[char] > 0 && have[char] == need[char] {
			formed++
		}

		for formed == required {
			windowLen := right - left + 1
			if windowLen < minLen {
				minLen = windowLen
			}

			leftChar := s[left]
			have[leftChar]--
			if need[leftChar] > 0 && have[leftChar] < need[leftChar] {
				formed--
			}
			left++
		}
	}

	if minLen > len(s) {
		return 0
	}
	return minLen
}
