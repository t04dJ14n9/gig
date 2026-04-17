package divergence_hunt115

import "fmt"

// ============================================================================
// Round 115: Global variable patterns and init
// Note: Each test function is run independently (not sequentially),
// so GlobalCompute only uses init-time values, not values from other functions.
// ============================================================================

var GlobalCounter int = 10

var GlobalSlice = []int{1, 2, 3}

var GlobalMap = map[string]int{"a": 1, "b": 2}

var GlobalString = "hello"

func init() {
	GlobalCounter += 5
}

func GlobalRead() int {
	return GlobalCounter
}

func GlobalModify() int {
	GlobalCounter += 10
	return GlobalCounter
}

func GlobalSliceRead() string {
	return fmt.Sprintf("%v", GlobalSlice)
}

func GlobalSliceLen() int {
	return len(GlobalSlice)
}

func GlobalMapRead() string {
	return fmt.Sprintf("%d", GlobalMap["a"])
}

func GlobalMapLen() int {
	return len(GlobalMap)
}

func GlobalStringRead() string {
	return GlobalString
}

// GlobalInitValues tests that init was run correctly
func GlobalInitValues() string {
	return fmt.Sprintf("%d:%d:%s", GlobalCounter, len(GlobalSlice), GlobalString)
}
