package divergence_hunt39

import (
	"encoding/json"
	"strings"
)

// ============================================================================
// Round 39: Slice edge cases - reslice, aliasing, append, copy, nil, empty
// ============================================================================

func ResliceAlias() int {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3] // shares backing array
	sub[0] = 99   // modifies s[1]
	return s[1]
}

func ResliceCap() int {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3]
	return cap(sub) // cap = 4 (from index 1 to end)
}

func ThreeIndexSlice() int {
	s := []int{0, 1, 2, 3, 4}
	sub := s[1:3:3] // len=2, cap=2
	return len(sub)*10 + cap(sub)
}

func ThreeIndexSliceNoAlias() int {
	s := []int{0, 1, 2, 3, 4}
	sub := s[1:3:3]
	_ = sub
	return len(s) // original unchanged
}

func AppendNil() int {
	var s []int
	s = append(s, 1, 2, 3)
	return len(s)
}

func AppendToEmpty() int {
	s := make([]int, 0)
	s = append(s, 10)
	return s[0]
}

func AppendSliceSpread() int {
	s := []int{1, 2}
	rest := []int{3, 4, 5}
	s = append(s, rest...)
	return len(s)
}

func CopySlice() int {
	src := []int{1, 2, 3}
	dst := make([]int, 3)
	n := copy(dst, src)
	return n*10 + dst[0] + dst[2]
}

func CopyPartial() int {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src[1:4])
	return n*10 + dst[0] + dst[2]
}

func NilSliceLenCap() int {
	var s []int
	return len(s)*10 + cap(s)
}

func EmptySliceLenCap() int {
	s := make([]int, 0)
	return len(s)*10 + cap(s)
}

func NilSliceCompare() bool {
	var s []int
	return s == nil
}

func EmptySliceNotNIl() bool {
	s := make([]int, 0)
	return s != nil
}

func SliceMakeWithCap() int {
	s := make([]int, 0, 10)
	return len(s)*100 + cap(s)
}

func SliceMakeWithLen() int {
	s := make([]int, 5)
	return len(s)*100 + s[0] // s[0] is zero
}

func SliceOfString() int {
	s := []string{"hello", "world"}
	return len(s[0]) + len(s[1])
}

func SliceOfBool() int {
	s := []bool{true, false, true}
	count := 0
	for _, v := range s {
		if v { count++ }
	}
	return count
}

func ByteSliceOperations() int {
	b := []byte("hello")
	b = append(b, ' ', 'w', 'o', 'r', 'l', 'd')
	return len(b)
}

func SliceDeletePattern() int {
	s := []int{1, 2, 3, 4, 5}
	i := 2
	s = append(s[:i], s[i+1:]...)
	return len(s)*10 + s[2] // len=4, s[2]=4
}

func JSONRoundTripSlice() int {
	s := []int{10, 20, 30}
	data, _ := json.Marshal(s)
	var decoded []int
	json.Unmarshal(data, &decoded)
	return decoded[0] + decoded[1] + decoded[2]
}

func StringSliceJoin() string {
	parts := []string{"a", "b", "c"}
	return strings.Join(parts, ",")
}
