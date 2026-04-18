package divergence_hunt145

import "fmt"

// ============================================================================
// Round 145: Break, continue, and labeled statements
// ============================================================================

func BreakBasic() string {
	for i := 0; i < 10; i++ {
		if i == 5 {
			break
		}
		_ = i
	}
	return "done"
}

func ContinueBasic() string {
	result := 0
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		result += i
	}
	return fmt.Sprintf("sum=%d", result)
}

func LabeledBreak() string {
	result := 0
outer:
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if i+j == 7 {
				result = i*10 + j
				break outer
			}
		}
	}
	return fmt.Sprintf("result=%d", result)
}

func LabeledContinue() string {
	result := 0
outer:
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == j {
				continue outer
			}
			result++
		}
	}
	return fmt.Sprintf("result=%d", result)
}

func RangeBreak() string {
	s := []int{1, 2, 3, 4, 5}
	sum := 0
	for _, v := range s {
		sum += v
		if sum > 6 {
			break
		}
	}
	return fmt.Sprintf("sum=%d", sum)
}

func RangeContinue() string {
	s := []int{1, 2, 3, 4, 5}
	var result []int
	for _, v := range s {
		if v%2 == 0 {
			continue
		}
		result = append(result, v)
	}
	return fmt.Sprintf("%v", result)
}

func NestedLoopBreak() string {
	found := false
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	return fmt.Sprintf("found=%t", found)
}

func SwitchBreakInLoop() string {
	for i := 0; i < 5; i++ {
		switch i {
		case 3:
			return fmt.Sprintf("found-%d", i)
		}
	}
	return "not-found"
}
