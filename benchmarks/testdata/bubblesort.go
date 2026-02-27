package main

func BubbleSort() int {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = 100 - i
	}
	n := len(s)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if s[j] > s[j+1] {
				tmp := s[j]
				s[j] = s[j+1]
				s[j+1] = tmp
			}
		}
	}
	return s[0] + s[99]
}
