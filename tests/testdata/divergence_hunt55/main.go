package divergence_hunt55

import "strings"

// ============================================================================
// Round 55: Control flow - switch, fallthrough, labeled break/continue
// ============================================================================

func SwitchBasic() int {
	x := 2
	switch x {
	case 1: return 10
	case 2: return 20
	case 3: return 30
	default: return 40
	}
}

func SwitchDefault() int {
	x := 99
	switch x {
	case 1: return 10
	default: return 99
	}
}

func SwitchMultiCase() int {
	x := 2
	switch x {
	case 1, 2, 3: return 10
	default: return 20
	}
}

func SwitchWithInit() int {
	switch x := 5; x {
	case 5: return 1
	default: return 0
	}
}

func SwitchExpression() string {
	x := 42
	switch {
	case x < 0: return "negative"
	case x == 0: return "zero"
	default: return "positive"
	}
}

func SwitchFallthrough() int {
	x := 1
	result := 0
	switch x {
	case 1: result += 1; fallthrough
	case 2: result += 10; fallthrough
	case 3: result += 100
	}
	return result
}

func SwitchNoMatch() int {
	x := 99
	switch x {
	case 1: return 1
	case 2: return 2
	}
	return 0
}

func NestedSwitch() int {
	x, y := 1, 2
	switch x {
	case 1:
		switch y {
		case 1: return 11
		case 2: return 12
		}
	case 2: return 20
	}
	return 0
}

func LabeledBreak() int {
	sum := 0
outer:
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if i+j > 5 { break outer }
			sum += i + j
		}
	}
	return sum
}

func LabeledContinue() int {
	sum := 0
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 { continue outer }
			sum += i*3 + j
		}
	}
	return sum
}

func ForBreak() int {
	sum := 0
	for i := 0; i < 10; i++ {
		if i == 5 { break }
		sum += i
	}
	return sum
}

func ForContinue() int {
	sum := 0
	for i := 0; i < 10; i++ {
		if i%2 == 0 { continue }
		sum += i
	}
	return sum
}

func InfiniteLoopBreak() int {
	i := 0
	for {
		i++
		if i >= 10 { break }
	}
	return i
}

func RangeBreak() int {
	s := []int{10, 20, 30, 40, 50}
	for _, v := range s {
		if v == 30 { return v }
	}
	return -1
}

func RangeContinue() int {
	s := []int{1, 2, 3, 4, 5}
	sum := 0
	for _, v := range s {
		if v%2 == 0 { continue }
		sum += v
	}
	return sum
}

func FmtSwitch() string {
	x := "hello"
	switch len(x) {
	case 5: return "five"
	default: return "other"
	}
}

func StringsSwitch() string {
	s := "hello"
	switch {
	case strings.HasPrefix(s, "he"): return "starts with he"
	case strings.HasSuffix(s, "lo"): return "ends with lo"
	default: return "other"
	}
}
