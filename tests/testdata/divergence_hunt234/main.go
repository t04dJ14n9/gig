package divergence_hunt234

import "fmt"

// ============================================================================
// Round 234: String concatenation patterns
// ============================================================================

func BasicConcatenation() string {
	s1 := "Hello"
	s2 := "World"
	return s1 + " " + s2
}

func ConcatenationInLoop() string {
	result := ""
	for i := 0; i < 5; i++ {
		result += fmt.Sprintf("%d", i)
	}
	return result
}

func ConcatenationWithNumbers() string {
	n := 42
	f := 3.14
	return fmt.Sprintf("int=%d,float=%f", n, f)
}

func ConcatenationWithBooleans() string {
	b1, b2 := true, false
	return fmt.Sprintf("%v:%v", b1, b2)
}

func MixedTypeConcatenation() string {
	name := "count"
	value := 10
	return fmt.Sprintf("%s=%d", name, value)
}

func StringBuilderPattern() string {
	parts := []string{"a", "b", "c", "d"}
	result := ""
	for _, p := range parts {
		result += p
	}
	return result
}

func JoinWithSeparator() string {
	items := []string{"apple", "banana", "cherry"}
	result := ""
	for i, item := range items {
		if i > 0 {
			result += ", "
		}
		result += item
	}
	return result
}

func ConcatenationWithRunes() string {
	r1, r2 := 'A', 'B'
	return string(r1) + string(r2)
}

func EmptyStringConcatenation() string {
	s := "start"
	s += ""
	s += "end"
	return s
}

func UnicodeConcatenation() string {
	s1 := "Hello"
	s2 := "世界"
	return s1 + " " + s2
}

func RepeatedConcatenation() string {
	s := "ha"
	result := ""
	for i := 0; i < 3; i++ {
		result += s
	}
	return result
}
