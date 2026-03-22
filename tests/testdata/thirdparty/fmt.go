package thirdparty

import (
	"bytes"
	"fmt"
)

// FmtSprintfVarious tests various Sprintf verbs.
func FmtSprintfVarious() int {
	s := fmt.Sprintf("%d %.2f %s %v", 42, 3.14, "hello", []int{1, 2})
	return len(s)
}

// FmtSprintfStruct tests Sprintf with struct.
func FmtSprintfStruct() int {
	type Point struct{ X, Y int }
	p := Point{X: 10, Y: 20}
	s := fmt.Sprintf("%v", p)
	return len(s)
}

// FmtSprintfPointer tests Sprintf with pointer.
func FmtSprintfPointer() int {
	x := 42
	s := fmt.Sprintf("%p", &x)
	return len(s)
}

// FmtSprintfBool tests Sprintf with bool.
func FmtSprintfBool() string {
	return fmt.Sprintf("%t", true)
}

// FmtSprintfHex tests Sprintf with hex.
func FmtSprintfHex() string {
	return fmt.Sprintf("%x", 255)
}

// FmtFprintf tests fmt.Fprintf.
func FmtFprintf() int {
	buf := new(bytes.Buffer)
	n, _ := fmt.Fprintf(buf, "value: %d", 42)
	return n
}

// FmtSprint tests fmt.Sprint.
func FmtSprint() string {
	return fmt.Sprint(1, 2, 3)
}

// FmtSprintln tests fmt.Sprintln.
func FmtSprintln() string {
	return fmt.Sprint(1, 2, 3)
}

// FmtErrorf tests fmt.Errorf.
func FmtErrorf() string {
	return fmt.Errorf("error: %d", 42).Error()
}
