package divergence_hunt73

// ============================================================================
// Round 73: Switch and type switch edge cases
// ============================================================================

func SwitchNoExpression() int {
	x := 42
	switch {
	case x < 0:
		return -1
	case x == 0:
		return 0
	default:
		return 1
	}
}

func SwitchMultiCase() int {
	x := 3
	switch x {
	case 1, 2, 3:
		return 10
	case 4, 5:
		return 20
	default:
		return 30
	}
}

func SwitchFallthrough() int {
	x := 1
	result := 0
	switch x {
	case 1:
		result += 1
		fallthrough
	case 2:
		result += 10
		fallthrough
	case 3:
		result += 100
	}
	return result
}

func SwitchWithInit() int {
	switch x := 5; x {
	case 5:
		return 1
	default:
		return 0
	}
}

func SwitchString() int {
	s := "hello"
	switch s {
	case "world":
		return 1
	case "hello":
		return 2
	default:
		return 3
	}
}

func SwitchEmpty() int {
	// Empty switch without cases - just return
	return 0
}

func SwitchOnlyDefault() int {
	switch 42 {
	default:
		return 99
	}
}

func TypeSwitchWithDefault() string {
	var x any = 3.14
	switch v := x.(type) {
	case int:
		return "int"
	case float64:
		return "float64"
	default:
		_ = v
		return "other"
	}
}

func TypeSwitchNil() string {
	var x any
	switch x.(type) {
	case nil:
		return "nil"
	case int:
		return "int"
	default:
		return "other"
	}
}

func SwitchBreak() int {
	result := 0
	switch 1 {
	case 1:
		result = 10
		break
		result = 20 // unreachable
	}
	return result
}

func SwitchNested() int {
	x, y := 1, 2
	switch x {
	case 1:
		switch y {
		case 1:
			return 11
		case 2:
			return 12
		}
	case 2:
		return 20
	}
	return 0
}

func SwitchBool() int {
	x := true
	switch x {
	case true:
		return 1
	case false:
		return 0
	default:
		return -1
	}
}

func SwitchInterface() string {
	var x any = "hello"
	switch v := x.(type) {
	case string:
		return v
	default:
		return ""
	}
}
