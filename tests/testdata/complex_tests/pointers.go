package complex_tests

func PointerBasic() int { x := 42; p := &x; return *p }

func PointerModify() int { x := 10; p := &x; *p = 20; return x }

func PointerSwap() int { a, b := 1, 2; pa, pb := &a, &b; *pa, *pb = *pb, *pa; return a*10+b }

func PointerToPointer() int { x := 42; p := &x; pp := &p; return **pp }

func PointerSlice() int { arr := []int{1,2,3}; p := &arr; *p = append(*p, 4); return len(arr) }
