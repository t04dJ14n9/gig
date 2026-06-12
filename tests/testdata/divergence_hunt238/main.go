package divergence_hunt238

import "fmt"

// ============================================================================
// Round 238: fmt.Sprintf formatting
// ============================================================================

func SprintfBasicVerbs() string {
	return fmt.Sprintf("int=%d string=%s bool=%t", 42, "hello", true)
}

func SprintfFloatVerbs() string {
	f := 3.14159
	return fmt.Sprintf("float=%f exp=%e", f, f)
}

func SprintfWidthAndPrecision() string {
	return fmt.Sprintf("|%5d|%-5d|%10s|%-10s|", 42, 42, "hi", "hi")
}

func SprintfFloatPrecision() string {
	return fmt.Sprintf("%.2f %.4f %.0f", 3.14159, 3.14159, 3.14159)
}

func SprintfVerbWidthPrecision() string {
	return fmt.Sprintf("|%8.2f|%-8.2f|%08.2f|", 3.14, 3.14, 3.14)
}

func SprintfHexOctalBinary() string {
	n := 255
	return fmt.Sprintf("hex=%x oct=%o bin(not_direct)=%d", n, n, n)
}

func SprintfPointerVerb() string {
	x := 42
	// Pointer addresses are non-deterministic; return a deterministic marker.
	s := fmt.Sprintf("ptr=%p", &x)
	if len(s) > 6 && s[:6] == "ptr=0x" {
		return "ptr=0x-ok"
	}
	return s
}

func SprintfGenericVerb() string {
	return fmt.Sprintf("int=%v string=%v bool=%v", 42, "hi", true)
}

func SprintfTypeVerb() string {
	return fmt.Sprintf("int=%T string=%T bool=%T", 42, "hi", true)
}

func SprintfPercentEscape() string {
	return fmt.Sprintf("100%% complete")
}

func SprintfMultipleValues() string {
	return fmt.Sprintf("a=%d b=%d c=%d d=%d e=%d", 1, 2, 3, 4, 5)
}

func SprintfStringPrecision() string {
	return fmt.Sprintf("%.3s %.5s %.10s", "hello", "world", "hi")
}
