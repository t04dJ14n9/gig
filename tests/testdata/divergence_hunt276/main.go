package divergence_hunt276

import (
	"fmt"
)

// ============================================================================
// Round 276: Switch/select edge cases — fallthrough, default, type switch with assignment

// SwitchFallthrough tests fallthrough in switch
func SwitchFallthrough() string {
	result := ""
	switch 1 {
	case 1:
		result += "one"
		fallthrough
	case 2:
		result += "two"
		fallthrough
	case 3:
		result += "three"
	}
	return result
}

// SwitchNoDefault tests switch without default
func SwitchNoDefault() string {
	x := 99
	result := "none"
	switch x {
	case 1:
		result = "one"
	case 2:
		result = "two"
	}
	return result
}

// SwitchDefaultOnly tests switch that only matches default
func SwitchDefaultOnly() string {
	x := 99
	result := ""
	switch x {
	case 1:
		result = "one"
	default:
		result = "default"
	}
	return result
}

// SwitchWithBreak tests break in switch
func SwitchWithBreak() string {
	result := ""
	switch 1 {
	case 1:
		result = "found"
		break
		result += "after_break"
	}
	return result
}

// SwitchMultipleValues tests case with multiple values
func SwitchMultipleValues() string {
	v := 5
	result := ""
	switch v {
	case 1, 3, 5, 7:
		result = "odd_single"
	case 2, 4, 6, 8:
		result = "even_single"
	default:
		result = "other"
	}
	return result
}

// SwitchExpression tests switch with expression
func SwitchExpression() string {
	x := 15
	result := ""
	switch {
	case x < 10:
		result = "small"
	case x < 20:
		result = "medium"
	case x < 30:
		result = "large"
	}
	return result
}

// TypeSwitchWithShortDecl tests type switch with short variable declaration
func TypeSwitchWithShortDecl() string {
	var val interface{} = "hello"
	switch v := val.(type) {
	case int:
		return fmt.Sprintf("int:%d", v)
	case string:
		return fmt.Sprintf("string:%s,len=%d", v, len(v))
	default:
		return fmt.Sprintf("other:%T", v)
	}
}

// SwitchInLoop tests switch with break in loop
func SwitchInLoop() string {
	result := ""
	for i := 0; i < 3; i++ {
		switch i {
		case 1:
			continue
		}
		result += fmt.Sprintf("%d", i)
	}
	return result
}

// SelectWithDefault tests select with default (non-blocking)
func SelectWithDefault() string {
	ch := make(chan int, 1)
	result := ""
	select {
	case v := <-ch:
		result = fmt.Sprintf("got:%d", v)
	default:
		result = "nothing"
	}
	return result
}

// SelectSend tests select with send
func SelectSend() string {
	ch := make(chan int, 1)
	ch <- 42
	result := ""
	select {
	case v := <-ch:
		result = fmt.Sprintf("received:%d", v)
	default:
		result = "blocked"
	}
	return result
}

// NestedSwitch tests nested switch
func NestedSwitch() string {
	result := ""
outer:
	for i := 0; i < 2; i++ {
		switch i {
		case 0:
			switch i * 10 {
			case 0:
				result += "zero"
			}
		case 1:
			result += "one"
			break outer
		}
		result += fmt.Sprintf(".%d.", i)
	}
	return result
}

// SwitchString tests switch on string
func SwitchString() string {
	s := "banana"
	result := ""
	switch s {
	case "apple":
		result = "red"
	case "banana":
		result = "yellow"
	case "cherry":
		result = "red_too"
	}
	return result
}
