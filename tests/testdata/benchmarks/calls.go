package benchmarks

// ============================================================================
// Function Call Overhead
// ============================================================================

func inc(x int) int { return x + 1 }

func CallOverhead() int {
	x := 0
	for i := 0; i < 1000; i++ {
		x = inc(x)
	}
	return x
}
