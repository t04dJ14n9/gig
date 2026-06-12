package divergence_hunt180

import (
	"fmt"
)

// ============================================================================
// Round 180: fmt package formatting verbs
// ============================================================================

func FmtVerbInt() string {
	return fmt.Sprintf("%d %v %+d", 42, 42, 42)
}

func FmtVerbFloat() string {
	return fmt.Sprintf("%f %.2f %e %E", 3.14159, 3.14159, 123.45, 123.45)
}

func FmtVerbString() string {
	return fmt.Sprintf("%s %q %x", "hello", "hello", "hello")
}

func FmtVerbBool() string {
	return fmt.Sprintf("%t %v", true, false)
}

func FmtVerbChar() string {
	return fmt.Sprintf("%c %q", 65, 65)
}

func FmtVerbBase() string {
	return fmt.Sprintf("%b %o %O %x %X", 42, 42, 42, 255, 255)
}

func FmtVerbWidth() string {
	return fmt.Sprintf("|%5d|%-5d|%05d|", 42, 42, 42)
}

func FmtVerbPrecision() string {
	return fmt.Sprintf("|%.2f|%.4f|%8.2f|", 3.14159, 3.14159, 3.14159)
}

func FmtVerbPointer() string {
	x := 42
	// Pointer addresses are non-deterministic; compare only the format prefix.
	s := fmt.Sprintf("%p", &x)
	if len(s) > 2 && s[:2] == "0x" {
		return "ok:0x-prefix"
	}
	return s
}

func FmtVerbType() string {
	return fmt.Sprintf("%T %T %T", 42, "hello", 3.14)
}

func FmtVerbSlice() string {
	nums := []int{1, 2, 3}
	return fmt.Sprintf("%v %+v %#v", nums, nums, nums)
}

func FmtVerbStruct() string {
	type Point struct{ X, Y int }
	p := Point{1, 2}
	return fmt.Sprintf("%v %+v %#v", p, p, p)
}

func FmtSprintVsSprintln() string {
	return fmt.Sprintf("%q:%q", fmt.Sprint("a", "b"), fmt.Sprintln("a", "b"))
}

func FmtFprintfSimulation() string {
	// Simulate fprintf with Sprint
	return fmt.Sprintf("name=%s age=%d", "Alice", 30)
}

func FmtErrorf() string {
	err := fmt.Errorf("error code: %d", 404)
	return err.Error()
}

func FmtComplex() string {
	c := 3 + 4i
	return fmt.Sprintf("%v %g %.1f", c, c, c)
}

func FmtUint() string {
	return fmt.Sprintf("%d %u", 42, uint(42))
}

func FmtHexPointer() string {
	x := 42
	// Pointer addresses are non-deterministic; return a deterministic marker.
	s := fmt.Sprintf("%#p", &x)
	if len(s) > 0 {
		return "ok"
	}
	return s
}

func FmtPadding() string {
	return fmt.Sprintf("|%10s|%-10s|%010d|", "hi", "hi", 42)
}

func FmtEmptyInterface() string {
	var i any = "hello"
	return fmt.Sprintf("%v %T", i, i)
}
