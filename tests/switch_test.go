package tests

import "testing"

// Switch statements are decomposed by SSA into if/else chains.

func TestSwitchSimple(t *testing.T) {
	runInt(t, `package main
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
func Compute() int { return classify(2) }`, 20)
}

func TestSwitchDefault(t *testing.T) {
	runInt(t, `package main
func classify(x int) int {
	switch x {
	case 1:
		return 10
	default:
		return -1
	}
}
func Compute() int { return classify(99) }`, -1)
}

func TestSwitchMultiCase(t *testing.T) {
	runInt(t, `package main
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
func Compute() int { return weekday(3)*10 + weekday(6)*1 }`, 10)
}

func TestSwitchNoCondition(t *testing.T) {
	// switch without condition acts like switch true
	runInt(t, `package main
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
func Compute() int { return grade(85)*10 + grade(55) }`, 30)
}

func TestSwitchWithInit(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 42
	switch v := x % 10; v {
	case 0:
		return 0
	case 2:
		return 2
	default:
		return -1
	}
}`, 2)
}

func TestSwitchStringCases(t *testing.T) {
	runInt(t, `package main
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
func Compute() int { return colorCode("green")*10 + colorCode("purple") }`, 20)
}

func TestSwitchFallthrough(t *testing.T) {
	// Go's switch doesn't fall through by default
	runInt(t, `package main
func Compute() int {
	x := 1
	result := 0
	switch x {
	case 1:
		result = 10
	case 2:
		result = 20
	}
	return result
}`, 10)
}

func TestSwitchNested(t *testing.T) {
	runInt(t, `package main
func Compute() int {
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
}`, 12)
}
