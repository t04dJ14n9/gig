package divergence_hunt229

import "fmt"

// ============================================================================
// Round 229: Nil slice vs empty slice
// ============================================================================

// NilVsEmptyLenCap tests len and cap of nil vs empty
func NilVsEmptyLenCap() string {
	var nilSlice []int
	emptySlice := []int{}
	madeEmpty := make([]int, 0)
	return fmt.Sprintf("nil=%d,%d,empty=%d,%d,made=%d,%d",
		len(nilSlice), cap(nilSlice),
		len(emptySlice), cap(emptySlice),
		len(madeEmpty), cap(madeEmpty))
}

// NilVsEmptyNilCheck tests nil check
func NilVsEmptyNilCheck() string {
	var nilSlice []int
	emptySlice := []int{}
	madeEmpty := make([]int, 0)
	return fmt.Sprintf("nil=%t,empty=%t,made=%t",
		nilSlice == nil,
		emptySlice == nil,
		madeEmpty == nil)
}

// NilVsEmptyAppend tests append behavior
func NilVsEmptyAppend() string {
	var nilSlice []int
	emptySlice := []int{}
	nilSlice = append(nilSlice, 1)
	emptySlice = append(emptySlice, 1)
	return fmt.Sprintf("nil=%v,empty=%v", nilSlice, emptySlice)
}

// NilVsEmptyIteration tests iteration
func NilVsEmptyIteration() string {
	var nilSlice []int
	emptySlice := []int{}
	nilCount, emptyCount := 0, 0
	for range nilSlice {
		nilCount++
	}
	for range emptySlice {
		emptyCount++
	}
	return fmt.Sprintf("nil=%d,empty=%d", nilCount, emptyCount)
}

// NilVsEmptyIndexingPanics tests indexing panics
func NilVsEmptyIndexingPanics() string {
	var nilSlice []int
	emptySlice := []int{}
	nilPanics := false
	emptyPanics := false

	func() {
		defer func() {
			if recover() != nil {
				nilPanics = true
			}
		}()
		_ = nilSlice[0]
	}()

	func() {
		defer func() {
			if recover() != nil {
				emptyPanics = true
			}
		}()
		_ = emptySlice[0]
	}()

	return fmt.Sprintf("nil_panics=%t,empty_panics=%t", nilPanics, emptyPanics)
}

// NilVsEmptyJSONMarshal tests JSON marshaling behavior
func NilVsEmptyJSONMarshal() string {
	var nilSlice []int
	emptySlice := []int{}
	nilJSON := fmt.Sprintf("%v", nilSlice)
	emptyJSON := fmt.Sprintf("%v", emptySlice)
	return fmt.Sprintf("nil=%s,empty=%s", nilJSON, emptyJSON)
}

// NilVsEmptySliceExpr tests slice expressions
func NilVsEmptySliceExpr() string {
	var nilSlice []int
	emptySlice := []int{}
	nilOK := false
	emptyOK := false

	func() {
		defer func() {
			if recover() != nil {
				nilOK = false
			} else {
				nilOK = true
			}
		}()
		_ = nilSlice[0:0]
	}()

	func() {
		defer func() {
			if recover() != nil {
				emptyOK = false
			} else {
				emptyOK = true
			}
		}()
		_ = emptySlice[0:0]
	}()

	return fmt.Sprintf("nil_ok=%t,empty_ok=%t", nilOK, emptyOK)
}

// NilVsEmptySlicingNil slicing nil slice
func NilVsEmptySlicingNil() string {
	var nilSlice []int
	sliced := nilSlice[0:0]
	return fmt.Sprintf("sliced_nil=%t,len=%d", sliced == nil, len(sliced))
}

// NilVsEmptyCopyToNil copy to nil slice
func NilVsEmptyCopyToNil() string {
	var nilSlice []int
	src := []int{1, 2, 3}
	n := copy(nilSlice, src)
	return fmt.Sprintf("copied=%d", n)
}

// NilVsEmptyEqual tests equality checks
func NilVsEmptyEqual() string {
	var nilSlice []int
	emptySlice := []int{}
	madeEmpty := make([]int, 0)
	return fmt.Sprintf("nil==empty:%t,empty==made:%t",
		nilSlice == nil && emptySlice != nil,
		len(emptySlice) == len(madeEmpty))
}

// NilVsEmptyRangeNil range over nil
func NilVsEmptyRangeNil() string {
	var nilSlice []int
	sum := 0
	for _, v := range nilSlice {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// NilVsEmptyLenOnly just len comparison
func NilVsEmptyLenOnly() string {
	var nilSlice []int
	emptySlice := []int{}
	return fmt.Sprintf("both_zero=%t", len(nilSlice) == 0 && len(emptySlice) == 0)
}
