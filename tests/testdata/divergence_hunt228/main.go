package divergence_hunt228

import "fmt"

// ============================================================================
// Round 228: Slice sharing and independence
// ============================================================================

// SliceSharingModify tests shared slice modification
func SliceSharingModify() string {
	original := []int{1, 2, 3, 4, 5}
	shared := original[1:4]
	shared[0] = 99
	return fmt.Sprintf("original[1]=%d", original[1])
}

// SliceSharingIndependence tests slice independence after append
func SliceSharingIndependence() string {
	original := make([]int, 3, 6)
	original[0], original[1], original[2] = 1, 2, 3
	independent := append(original, 4)
	independent[0] = 99
	return fmt.Sprintf("original[0]=%d,independent[0]=%d", original[0], independent[0])
}

// SliceSharingCapacityCheck tests capacity affects sharing
func SliceSharingCapacityCheck() string {
	original := make([]int, 3, 10)
	original[0], original[1], original[2] = 1, 2, 3
	shared := original[:]
	shared = append(shared, 4)
	shared[0] = 99
	return fmt.Sprintf("orig[0]=%d,shared_len=%d", original[0], len(shared))
}

// SliceSharingCopy creates independent copy
func SliceSharingCopy() string {
	original := []int{1, 2, 3}
	independent := make([]int, len(original))
	copy(independent, original)
	independent[0] = 99
	return fmt.Sprintf("original[0]=%d", original[0])
}

// SliceSharingMultipleRefs multiple references to same backing
func SliceSharingMultipleRefs() string {
	base := []int{1, 2, 3, 4, 5}
	s1 := base[0:3]
	s2 := base[1:4]
	s3 := base[2:5]
	s2[0] = 99
	return fmt.Sprintf("s1[1]=%d,s2[0]=%d,s3[0]=%d", s1[1], s2[0], s3[0])
}

// SliceSharingResliceIndependence tests if reslicing creates independence
func SliceSharingResliceIndependence() string {
	original := []int{1, 2, 3, 4, 5}
	resliced := original[1:3]
	resliced = append(resliced, 99)
	return fmt.Sprintf("original[3]=%d", original[3])
}

// SliceSharingFullReslice tests full reslice behavior
func SliceSharingFullReslice() string {
	original := make([]int, 5, 10)
	for i := range original {
		original[i] = i + 1
	}
	limited := original[0:3:3]
	limited = append(limited, 99)
	return fmt.Sprintf("orig[3]=%d,limited_len=%d,limited_cap=%d", original[3], len(limited), cap(limited))
}

// SliceSharingPointerSharing tests sharing via pointer
func SliceSharingPointerSharing() string {
	original := []int{1, 2, 3}
	ptr := &original[1]
	original[1] = 99
	return fmt.Sprintf("*ptr=%d", *ptr)
}

// SliceSharingNestedSlice tests nested slice sharing
func SliceSharingNestedSlice() string {
	outer := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	inner := outer[0]
	inner[0] = 99
	return fmt.Sprintf("outer[0][0]=%d", outer[0][0])
}

// SliceSharingAppendAndModify tests append then modify behavior
func SliceSharingAppendAndModify() string {
	original := make([]int, 2, 4)
	original[0], original[1] = 1, 2
	appended := append(original, 3)
	appended[0] = 99
	return fmt.Sprintf("orig[0]=%d,app_len=%d", original[0], len(appended))
}

// SliceSharingMakeCopy explicit make and copy
func SliceSharingMakeCopy() string {
	original := []int{1, 2, 3, 4, 5}
	copy1 := make([]int, len(original))
	copy(copy1, original)
	copy1[0] = 99
	return fmt.Sprintf("orig[0]=%d,copy1[0]=%d", original[0], copy1[0])
}
