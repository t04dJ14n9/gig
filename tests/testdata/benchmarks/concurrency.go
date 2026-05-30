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

func PanicRecover() int {
	sum := 0
	for j := 0; j < 10; j++ {
		func() {
			defer func() { recover() }()
			if j == 5 {
				panic("test")
			}
			sum += j
		}()
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
