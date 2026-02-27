package benchmarks

// ============================================================================
// Concurrency Operations
// ============================================================================

func Defer() int {
	sum := 0
	for i := 0; i < 10; i++ {
		defer func() { sum = sum + 1 }()
	}
	return sum
}

func Select() int {
	ch := make(chan int, 1)
	sum := 0
	for i := 0; i < 100; i++ {
		select {
		case ch <- i:
		default:
			sum = sum + 1
		}
		select {
		case v := <-ch:
			sum = sum + v
		default:
			sum = sum + 1
		}
	}
	return sum
}
