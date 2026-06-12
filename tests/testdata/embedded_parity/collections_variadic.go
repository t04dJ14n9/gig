package main

import "fmt"

func variadicSum(nums ...int) string {
	total := 0
	for _, n := range nums {
		total += n
	}
	return fmt.Sprintf("%d/%d/%d", total, len(nums), cap(nums))
}

// Result checks collection edge cases and variadic spread behavior.
func Result() string {
	a := []int{1, 2, 3, 4, 5}
	limited := a[1:3:3]
	limited = append(limited, 99)

	shared := a[1:3]
	shared[0] = 77

	var nilMap map[string]int
	nilRead := nilMap["missing"]
	nilWrite := "no-panic"
	func() {
		defer func() {
			if recover() != nil {
				nilWrite = "panic"
			}
		}()
		nilMap["x"] = 1
	}()

	m := map[string]int{"x": 42}
	got, ok := m["x"]
	missing, missingOK := m["missing"]

	nums := []int{10, 20, 30, 40}
	spread := variadicSum(nums[1:3]...)
	var nilNums []int
	nilSpread := variadicSum(nilNums...)

	return fmt.Sprintf("limited=%d/%d/%v:a=%v:nil=%d/%s:map=%d/%t/%d/%t:var=%s:%s",
		len(limited), cap(limited), limited, a,
		nilRead, nilWrite,
		got, ok, missing, missingOK,
		spread, nilSpread)
}
