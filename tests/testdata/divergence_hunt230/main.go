package divergence_hunt230

import "fmt"

// ============================================================================
// Round 230: Three-index slicing
// ============================================================================

// ThreeIndexBasic tests basic three-index slicing
func ThreeIndexBasic() string {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s2 := s[2:5:7]
	return fmt.Sprintf("len=%d,cap=%d", len(s2), cap(s2))
}

// ThreeIndexFullCapacity tests three-index with max capacity
func ThreeIndexFullCapacity() string {
	s := make([]int, 5, 10)
	s2 := s[1:4:10]
	return fmt.Sprintf("len=%d,cap=%d", len(s2), cap(s2))
}

// ThreeIndexLimits tests three-index at limits
func ThreeIndexLimits() string {
	s := []int{1, 2, 3, 4, 5}
	s2 := s[0:0:5]
	return fmt.Sprintf("len=%d,cap=%d", len(s2), cap(s2))
}

// ThreeIndexAppendIndependence tests append independence with three-index
func ThreeIndexAppendIndependence() string {
	s := []int{1, 2, 3, 4, 5}
	s2 := s[1:3:3]
	s2 = append(s2, 99)
	return fmt.Sprintf("s[3]=%d,s2=%v", s[3], s2)
}

// ThreeIndexPanicLow tests panic on low > high
func ThreeIndexPanicLow() string {
	panics := false
	func() {
		defer func() {
			if recover() != nil {
				panics = true
			}
		}()
		s := []int{1, 2, 3}
		low, high := 3, 2
		_ = s[low:high:3]
	}()
	return fmt.Sprintf("panics=%t", panics)
}

// ThreeIndexPanicMax tests panic on high > max
func ThreeIndexPanicMax() string {
	panics := false
	func() {
		defer func() {
			if recover() != nil {
				panics = true
			}
		}()
		s := make([]int, 3, 5)
		_ = s[0:3:6]
	}()
	return fmt.Sprintf("panics=%t", panics)
}

// ThreeIndexByteSlice tests three-index with byte slice
func ThreeIndexByteSlice() string {
	b := []byte("hello world")
	b2 := b[0:5:11]
	return fmt.Sprintf("len=%d,cap=%d", len(b2), cap(b2))
}

// ThreeIndexStringSlice tests three-index with string slice
func ThreeIndexStringSlice() string {
	s := []string{"a", "b", "c", "d", "e"}
	s2 := s[1:4:5]
	return fmt.Sprintf("len=%d,cap=%d,first=%s", len(s2), cap(s2), s2[0])
}

// ThreeIndexNestedSlice tests three-index with nested slices
func ThreeIndexNestedSlice() string {
	s := [][]int{
		{1, 2}, {3, 4}, {5, 6}, {7, 8},
	}
	s2 := s[1:3:4]
	return fmt.Sprintf("len=%d,cap=%d", len(s2), cap(s2))
}

// ThreeIndexThenReslice tests reslicing after three-index
func ThreeIndexThenReslice() string {
	s := make([]int, 10, 20)
	s2 := s[2:5:8]
	s3 := s2[0:3]
	return fmt.Sprintf("len=%d,cap=%d", len(s3), cap(s3))
}

// ThreeIndexPreserveOriginal tests that original is preserved
func ThreeIndexPreserveOriginal() string {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s2 := s[2:4:6]
	s2[0] = 99
	return fmt.Sprintf("s[2]=%d,s2[0]=%d", s[2], s2[0])
}

// ThreeIndexCapacityControl tests capacity control
func ThreeIndexCapacityControl() string {
	s := make([]int, 5, 20)
	limited := s[0:5:5]
	full := s[0:5:20]
	return fmt.Sprintf("limited_cap=%d,full_cap=%d", cap(limited), cap(full))
}
