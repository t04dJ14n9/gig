package main

func Sieve() int {
	n := 1000
	sieve := make([]int, n+1)
	for i := 2; i <= n; i++ {
		sieve[i] = 1
	}
	for i := 2; i*i <= n; i++ {
		if sieve[i] == 1 {
			for j := i * i; j <= n; j = j + i {
				sieve[j] = 0
			}
		}
	}
	count := 0
	for i := 2; i <= n; i++ {
		if sieve[i] == 1 {
			count = count + 1
		}
	}
	return count
}
