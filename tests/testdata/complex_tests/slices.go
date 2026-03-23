package complex_tests

// ============================================================================
// Complex Slice Operations Tests (25 tests)
// ============================================================================

// SliceFlatten tests flattening nested slices.
func SliceFlatten() int {
	nested := [][]int{{1, 2}, {3, 4, 5}, {6}}
	flat := []int{}
	for _, inner := range nested {
		for _, v := range inner {
			flat = append(flat, v)
		}
	}
	sum := 0
	for _, v := range flat {
		sum += v
	}
	return sum
}

// SliceChunk tests chunking slice.
func SliceChunk() int {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	chunks := [][]int{}
	size := 3
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}
	return len(chunks)
}

// SliceRotateLeft tests left rotation.
func SliceRotateLeft() int {
	arr := []int{1, 2, 3, 4, 5}
	k := 2
	k = k % len(arr)
	result := append(arr[k:], arr[:k]...)
	return result[0]*1000 + result[1]*100 + result[2]*10 + result[3]
}

// SliceRotateRight tests right rotation.
func SliceRotateRight() int {
	arr := []int{1, 2, 3, 4, 5}
	k := 2
	k = k % len(arr)
	result := append(arr[len(arr)-k:], arr[:len(arr)-k]...)
	return result[0]*1000 + result[1]*100 + result[2]*10 + result[3]
}

// SliceUnique tests unique elements.
func SliceUnique() int {
	arr := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	seen := make(map[int]bool)
	result := []int{}
	for _, v := range arr {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceIntersect tests intersection.
func SliceIntersect() int {
	a := []int{1, 2, 3, 4, 5}
	b := []int{3, 4, 5, 6, 7}
	setA := make(map[int]bool)
	for _, v := range a {
		setA[v] = true
	}
	result := []int{}
	for _, v := range b {
		if setA[v] {
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceUnion tests union.
func SliceUnion() int {
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	set := make(map[int]bool)
	for _, v := range a {
		set[v] = true
	}
	for _, v := range b {
		set[v] = true
	}
	return len(set)
}

// SliceDifference tests difference.
func SliceDifference() int {
	a := []int{1, 2, 3, 4, 5}
	b := []int{3, 4, 5, 6, 7}
	setB := make(map[int]bool)
	for _, v := range b {
		setB[v] = true
	}
	result := []int{}
	for _, v := range a {
		if !setB[v] {
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceZip tests zipping slices.
func SliceZip() int {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	result := []int{}
	for i := 0; i < len(a) && i < len(b); i++ {
		result = append(result, a[i], b[i])
	}
	return result[0] + result[1] + result[2] + result[3]
}

// SlicePartition tests partitioning by predicate.
func SlicePartition() int {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := []int{}
	odds := []int{}
	for _, v := range arr {
		if v%2 == 0 {
			evens = append(evens, v)
		} else {
			odds = append(odds, v)
		}
	}
	return len(evens)*10 + len(odds)
}

// SliceTake tests taking first n elements.
func SliceTake() int {
	arr := []int{1, 2, 3, 4, 5}
	n := 3
	if n > len(arr) {
		n = len(arr)
	}
	result := arr[:n]
	return result[0] + result[1] + result[2]
}

// SliceDrop tests dropping first n elements.
func SliceDrop() int {
	arr := []int{1, 2, 3, 4, 5}
	n := 2
	if n > len(arr) {
		n = len(arr)
	}
	result := arr[n:]
	return result[0] + result[1] + result[2]
}

// SliceTakeWhile tests taking while predicate is true.
func SliceTakeWhile() int {
	arr := []int{2, 4, 6, 8, 1, 3, 5}
	result := []int{}
	for _, v := range arr {
		if v%2 != 0 {
			break
		}
		result = append(result, v)
	}
	return len(result)
}

// SliceDropWhile tests dropping while predicate is true.
func SliceDropWhile() int {
	arr := []int{2, 4, 6, 8, 1, 3, 5}
	result := []int{}
	dropping := true
	for _, v := range arr {
		if dropping && v%2 != 0 {
			dropping = false
		}
		if !dropping {
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceGroupBy tests grouping by key.
func SliceGroupBy() int {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	groups := make(map[int][]int)
	for _, v := range arr {
		key := v % 3
		groups[key] = append(groups[key], v)
	}
	return len(groups[0]) + len(groups[1])*10 + len(groups[2])*100
}

// SliceWindow tests sliding window.
func SliceWindow() int {
	arr := []int{1, 2, 3, 4, 5}
	size := 3
	windows := [][]int{}
	for i := 0; i <= len(arr)-size; i++ {
		windows = append(windows, arr[i:i+size])
	}
	return windows[0][0] + windows[1][1] + windows[2][2]
}

// SlicePairwise tests pairwise iteration.
func SlicePairwise() int {
	arr := []int{1, 2, 3, 4, 5}
	sum := 0
	for i := 0; i < len(arr)-1; i++ {
		sum += arr[i] + arr[i+1]
	}
	return sum
}

// SliceCartesian tests cartesian product.
func SliceCartesian() int {
	a := []int{1, 2}
	b := []int{3, 4}
	pairs := [][]int{}
	for _, x := range a {
		for _, y := range b {
			pairs = append(pairs, []int{x, y})
		}
	}
	return len(pairs)
}

// SliceTranspose tests transpose.
func SliceTranspose() int {
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}}
	result := [][]int{}
	for j := 0; j < len(matrix[0]); j++ {
		row := []int{}
		for i := 0; i < len(matrix); i++ {
			row = append(row, matrix[i][j])
		}
		result = append(result, row)
	}
	return result[0][0] + result[0][1] + result[1][0] + result[1][1]
}

// SliceReverseInPlace tests in-place reversal.
func SliceReverseInPlace() int {
	arr := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr[0]*10000 + arr[1]*1000 + arr[2]*100 + arr[3]*10 + arr[4]
}

// SliceCompact tests removing consecutive duplicates.
func SliceCompact() int {
	arr := []int{1, 1, 2, 2, 2, 3, 3, 2, 2, 1}
	result := []int{}
	for i, v := range arr {
		if i == 0 || v != arr[i-1] {
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceProduct tests computing product.
func SliceProduct() int {
	arr := []int{1, 2, 3, 4, 5}
	product := 1
	for _, v := range arr {
		product *= v
	}
	return product
}

// SliceCumSum tests cumulative sum.
func SliceCumSum() int {
	arr := []int{1, 2, 3, 4, 5}
	result := []int{}
	sum := 0
	for _, v := range arr {
		sum += v
		result = append(result, sum)
	}
	return result[len(result)-1]
}

// SliceCumProd tests cumulative product.
func SliceCumProd() int {
	arr := []int{1, 2, 3, 4}
	result := []int{}
	product := 1
	for _, v := range arr {
		product *= v
		result = append(result, product)
	}
	return result[len(result)-1]
}

// SliceFind tests finding element.
func SliceFind() int {
	arr := []int{1, 2, 3, 4, 5}
	target := 3
	for i, v := range arr {
		if v == target {
			return i
		}
	}
	return -1
}

// SliceFindLast tests finding last element.
func SliceFindLast() int {
	arr := []int{1, 2, 3, 2, 1}
	target := 2
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] == target {
			return i
		}
	}
	return -1
}
