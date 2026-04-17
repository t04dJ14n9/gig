package divergence_hunt106

import "fmt"

// ============================================================================
// Round 106: Switch statements - expression, type, fallthrough
// ============================================================================

func SwitchBasic() string {
	x := 2
	switch x {
	case 1:
		return "one"
	case 2:
		return "two"
	case 3:
		return "three"
	default:
		return "other"
	}
}

func SwitchDefault() string {
	x := 99
	switch x {
	case 1:
		return "one"
	default:
		return "default"
	}
}

func SwitchMultipleValues() string {
	x := 5
	switch x {
	case 1, 2, 3:
		return "low"
	case 4, 5, 6:
		return "mid"
	default:
		return "high"
	}
}

func SwitchNoExpression() string {
	x := 15
	switch {
	case x < 10:
		return "small"
	case x < 20:
		return "medium"
	default:
		return "large"
	}
}

func SwitchFallthrough() string {
	x := 1
	result := ""
	switch x {
	case 1:
		result += "one"
		fallthrough
	case 2:
		result += ":two"
	}
	return result
}

func SwitchInLoop() string {
	items := []string{"a", "bb", "ccc"}
	result := ""
	for _, s := range items {
		switch len(s) {
		case 1:
			result += "S"
		case 2:
			result += "M"
		default:
			result += "L"
		}
	}
	return result
}

func SwitchBreak() string {
	result := ""
loop:
	for i := 0; i < 3; i++ {
		switch i {
		case 1:
			break loop
		default:
			result += fmt.Sprintf("%d", i)
		}
	}
	return result
}

func SwitchString() string {
	s := "hello"
	switch s {
	case "hi":
		return "casual"
	case "hello":
		return "formal"
	default:
		return "unknown"
	}
}

func SwitchWithInit() string {
	switch x := computeValue(); x {
	case 42:
		return "answer"
	default:
		return fmt.Sprintf("other:%d", x)
	}
}

func computeValue() int { return 42 }

func NestedSwitch() string {
	x, y := 1, 2
	switch x {
	case 1:
		switch y {
		case 2:
			return "1,2"
		default:
			return "1,?"
		}
	default:
		return "?,?"
	}
}
