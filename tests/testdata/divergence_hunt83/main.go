package divergence_hunt83

// ============================================================================
// Round 83: Global variable edge cases - package-level vars, init order
// ============================================================================

var GlobalCounter int

func InitGlobal() int {
	GlobalCounter = 10
	return GlobalCounter
}

func IncrementGlobal() int {
	GlobalCounter++
	return GlobalCounter
}

var GlobalSlice = []int{1, 2, 3}

func GlobalSliceAccess() int {
	return len(GlobalSlice)
}

var GlobalMap = map[string]int{"key": 42}

func GlobalMapAccess() int {
	return GlobalMap["key"]
}

var GlobalString = "hello"

func GlobalStringAccess() string {
	return GlobalString
}

func GlobalModifySlice() int {
	GlobalSlice[0] = 99
	return GlobalSlice[0]
}

func GlobalModifyMap() int {
	GlobalMap["new"] = 100
	return len(GlobalMap)
}

const GlobalConst = 42

func GlobalConstAccess() int {
	return GlobalConst
}

const (
	ConstA = iota
	ConstB
	ConstC
)

func GlobalIota() int {
	return ConstA + ConstB + ConstC
}

var GlobalBool = true

func GlobalBoolAccess() bool {
	return GlobalBool
}

var GlobalFloat = 3.14

func GlobalFloatAccess() float64 {
	return GlobalFloat
}

var GlobalPointer *int

func GlobalPointerNil() bool {
	return GlobalPointer == nil
}
