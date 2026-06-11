package divergence_hunt79

// ============================================================================
// Round 79: For-range edge cases - string runes, map order, channel, index only
// ============================================================================

func RangeSlice() int {
	s := []int{10, 20, 30}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func RangeSliceIndex() int {
	s := []int{10, 20, 30}
	sum := 0
	for i := range s {
		sum += i
	}
	return sum
}

func RangeStringRunes() int {
	s := "Hello"
	count := 0
	for _ = range s {
		count++
	}
	return count
}

func RangeMapKeys() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	keys := 0
	for k := range m {
		keys += k
	}
	return keys
}

func RangeModifySlice() []int {
	s := []int{1, 2, 3, 4, 5}
	for i := range s {
		s[i] *= 2
	}
	return s
}

func RangeWithBreak() int {
	s := []int{1, 2, 3, 4, 5}
	sum := 0
	for _, v := range s {
		if v > 3 {
			break
		}
		sum += v
	}
	return sum
}

func RangeWithContinue() int {
	s := []int{1, 2, 3, 4, 5}
	sum := 0
	for _, v := range s {
		if v%2 == 0 {
			continue
		}
		sum += v
	}
	return sum
}

func RangeEmptySlice() int {
	s := []int{}
	count := 0
	for _ = range s {
		count++
	}
	return count
}

func RangeNilSlice() int {
	var s []int
	count := 0
	for _ = range s {
		count++
	}
	return count
}

func RangeNilMap() int {
	var m map[string]int
	count := 0
	for _ = range m {
		count++
	}
	return count
}

func RangeChannel() int {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

func RangeArray() int {
	arr := [3]int{100, 200, 300}
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

func RangeMultiByteString() int {
	s := "Hello世界"
	count := 0
	for _ = range s {
		count++
	}
	return count
}

func RangeStringIndexRune() string {
	s := "abc"
	result := ""
	for i, r := range s {
		if i > 0 {
			result += ","
		}
		result += string(r)
	}
	return result
}
