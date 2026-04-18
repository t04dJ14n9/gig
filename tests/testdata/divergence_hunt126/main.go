package divergence_hunt126

import "fmt"

// ============================================================================
// Round 126: Switch fallthrough and default behavior
// ============================================================================

func SwitchFallthrough() string {
	x := 1
	var result string
	switch x {
	case 1:
		result = "one"
		fallthrough
	case 2:
		result += "-two"
		fallthrough
	case 3:
		result += "-three"
	}
	return result
}

func SwitchNoFallthrough() string {
	x := 2
	var result string
	switch x {
	case 1:
		result = "one"
	case 2:
		result = "two"
	case 3:
		result = "three"
	}
	return result
}

func SwitchDefaultOnly() string {
	x := 99
	var result string
	switch x {
	case 1:
		result = "one"
	default:
		result = "default"
	}
	return result
}

func SwitchCaseOrder() string {
	x := 5
	var result string
	switch {
	case x > 10:
		result = "big"
	case x > 3:
		result = "medium"
	case x > 0:
		result = "small"
	}
	return result
}

func SwitchStringCase() string {
	s := "banana"
	var result string
	switch s {
	case "apple":
		result = "fruit-a"
	case "banana":
		result = "fruit-b"
	case "cherry":
		result = "fruit-c"
	}
	return result
}

func SwitchNoCaseNoDefault() string {
	x := 42
	var result string
	switch x {
	case 1:
		result = "one"
	}
	if result == "" {
		result = "none"
	}
	return result
}

func SwitchBreakExplicit() string {
	x := 1
	var result string
	switch x {
	case 1:
		result = "one"
		break
		// This code is unreachable but valid
		result += "-unreachable"
	}
	return result
}

func SwitchMultiCase() string {
	x := 2
	var result string
	switch x {
	case 1, 2, 3:
		result = "small"
	case 4, 5, 6:
		result = "medium"
	}
	return result
}

func SwitchInLoop() string {
	for i := 0; i < 5; i++ {
		switch i {
		case 2:
			return fmt.Sprintf("found-%d", i)
		}
	}
	return "not-found"
}
