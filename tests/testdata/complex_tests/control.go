package complex_tests

func ControlNestedIf() int {
	x, y, z := 1, 2, 3
	result := 0
	if x > 0 {
		if y > 1 {
			if z > 2 { result = 100 } else { result = 10 }
		} else { result = 1 }
	}
	return result
}

func ControlNestedLoop() int {
	sum := 0
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ { sum += i*j }
	}
	return sum
}

func ControlTripleLoop() int {
	count := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				if i+j+k == 6 { count++ }
			}
		}
	}
	return count
}
