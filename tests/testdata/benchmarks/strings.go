package benchmarks

// ============================================================================
// String Operations
// ============================================================================

func StringConcat() int {
	s := ""
	for i := 0; i < 100; i++ {
		s = s + "abcdefghij"
	}
	return len(s)
}
