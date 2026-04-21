package divergence_hunt259

import (
	"fmt"
)

// ============================================================================
// Round 259: Init functions
// ============================================================================

var initOrder string

func init() {
	initOrder += "first;"
}

func init() {
	initOrder += "second;"
}

// InitOrder tests the execution order of init functions
func InitOrder() string {
	return fmt.Sprintf("order=%s", initOrder)
}

// InitWithVariable tests init with package-level variable
var initVar int

func init() {
	initVar = 42
}

func InitWithVariable() string {
	return fmt.Sprintf("initVar=%d", initVar)
}

// InitWithMap tests init with map initialization
var initMap map[string]int

func init() {
	initMap = make(map[string]int)
	initMap["key"] = 100
}

func InitWithMap() string {
	return fmt.Sprintf("initMap[key]=%d", initMap["key"])
}

// InitWithSlice tests init with slice
var initSlice []int

func init() {
	initSlice = []int{1, 2, 3}
}

func InitWithSlice() string {
	return fmt.Sprintf("len=%d,cap=%d", len(initSlice), cap(initSlice))
}

// InitWithStruct tests init with struct
var initStruct struct {
	Name string
	Age  int
}

func init() {
	initStruct.Name = "test"
	initStruct.Age = 25
}

func InitWithStruct() string {
	return fmt.Sprintf("name=%s,age=%d", initStruct.Name, initStruct.Age)
}

// InitMultipleVars tests init with multiple variables
var (
	varA int
	varB string
	varC bool
)

func init() {
	varA = 10
	varB = "initialized"
	varC = true
}

func InitMultipleVars() string {
	return fmt.Sprintf("A=%d,B=%s,C=%v", varA, varB, varC)
}

// InitComplexSetup tests init with complex setup
var config map[string]interface{}

func init() {
	config = map[string]interface{}{
		"port":    8080,
		"debug":   false,
		"timeout": 30,
	}
}

func InitComplexSetup() string {
	return fmt.Sprintf("port=%v,debug=%v", config["port"], config["debug"])
}

// InitCounter tests init side effects
var initCounter int

func init() {
	initCounter++
}

func init() {
	initCounter += 10
}

func init() {
	initCounter += 100
}

func InitCounter() string {
	return fmt.Sprintf("counter=%d", initCounter)
}

// InitWithArray tests init with array
var initArray [3]int

func init() {
	initArray[0] = 100
	initArray[1] = 200
	initArray[2] = 300
}

func InitWithArray() string {
	return fmt.Sprintf("array=[%d,%d,%d]", initArray[0], initArray[1], initArray[2])
}

// InitConditional tests init with conditional logic
var conditionalValue int

func init() {
	if true {
		conditionalValue = 1
	} else {
		conditionalValue = 2
	}
}

func InitConditional() string {
	return fmt.Sprintf("value=%d", conditionalValue)
}

// InitLoop tests init with loop
var sumFromInit int

func init() {
	for i := 1; i <= 5; i++ {
		sumFromInit += i
	}
}

func InitLoop() string {
	return fmt.Sprintf("sum=%d", sumFromInit)
}
