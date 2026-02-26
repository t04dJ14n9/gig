package switch_pkg

// Simple tests simple switch
func Simple() int { return classify(2) }

func classify(x int) int {
	switch x {
	case 1:
		return 10
	case 2:
		return 20
	case 3:
		return 30
	default:
		return -1
	}
}

// Default tests switch default
func Default() int { return classify(99) }

// MultiCase tests switch with multiple cases
func MultiCase() int { return weekday(3)*10 + weekday(6) }

func weekday(d int) int {
	switch d {
	case 1, 2, 3, 4, 5:
		return 1
	case 6, 7:
		return 0
	default:
		return -1
	}
}

// NoCondition tests switch without condition
func NoCondition() int { return grade(85)*10 + grade(55) }

func grade(score int) int {
	switch {
	case score >= 90:
		return 4
	case score >= 80:
		return 3
	case score >= 70:
		return 2
	case score >= 60:
		return 1
	default:
		return 0
	}
}

// WithInit tests switch with init
func WithInit() int {
	x := 42
	switch v := x % 10; v {
	case 0:
		return 0
	case 2:
		return 2
	default:
		return -1
	}
}

// StringCases tests switch with string cases
func StringCases() int { return colorCode("green")*10 + colorCode("purple") }

func colorCode(name string) int {
	switch name {
	case "red":
		return 1
	case "green":
		return 2
	case "blue":
		return 3
	default:
		return 0
	}
}

// Fallthrough tests switch fallthrough behavior
func Fallthrough() int {
	x := 1
	result := 0
	switch x {
	case 1:
		result = 10
	case 2:
		result = 20
	}
	return result
}

// Nested tests nested switch
func Nested() int {
	a := 1
	b := 2
	switch a {
	case 1:
		switch b {
		case 1:
			return 11
		case 2:
			return 12
		default:
			return 10
		}
	case 2:
		return 20
	default:
		return 0
	}
}
