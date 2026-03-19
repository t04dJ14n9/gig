// Package tests - correctness_test.go
//
// Unified correctness test framework that consolidates ALL tests from testdata/
// Compares interpreted execution results with native Go execution results.
package tests

import (
	_ "embed"
	"reflect"
	"testing"
	"time"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
	"git.woa.com/youngjin/gig/tests/testdata/advanced"
	"git.woa.com/youngjin/gig/tests/testdata/algorithms"
	"git.woa.com/youngjin/gig/tests/testdata/arithmetic"
	"git.woa.com/youngjin/gig/tests/testdata/autowrap"
	"git.woa.com/youngjin/gig/tests/testdata/bitwise"
	"git.woa.com/youngjin/gig/tests/testdata/closures"
	"git.woa.com/youngjin/gig/tests/testdata/closures_advanced"
	"git.woa.com/youngjin/gig/tests/testdata/channels"
	"git.woa.com/youngjin/gig/tests/testdata/controlflow"
	"git.woa.com/youngjin/gig/tests/testdata/cornercases"
	"git.woa.com/youngjin/gig/tests/testdata/edgecases"
	"git.woa.com/youngjin/gig/tests/testdata/external"
	"git.woa.com/youngjin/gig/tests/testdata/functions"
	"git.woa.com/youngjin/gig/tests/testdata/goroutine"
	"git.woa.com/youngjin/gig/tests/testdata/initialize"
	"git.woa.com/youngjin/gig/tests/testdata/leetcode_hard"
	"git.woa.com/youngjin/gig/tests/testdata/mapadvanced"
	"git.woa.com/youngjin/gig/tests/testdata/maps"
	"git.woa.com/youngjin/gig/tests/testdata/multiassign"
	"git.woa.com/youngjin/gig/tests/testdata/namedreturn"
	"git.woa.com/youngjin/gig/tests/testdata/recursion"
	"git.woa.com/youngjin/gig/tests/testdata/resolved_issue"
	"git.woa.com/youngjin/gig/tests/testdata/scope"
	"git.woa.com/youngjin/gig/tests/testdata/slices"
	"git.woa.com/youngjin/gig/tests/testdata/slicing"
	"git.woa.com/youngjin/gig/tests/testdata/strings_pkg"
	"git.woa.com/youngjin/gig/tests/testdata/structs"
	switch_pkg "git.woa.com/youngjin/gig/tests/testdata/switch"
	"git.woa.com/youngjin/gig/tests/testdata/tricky"
	"git.woa.com/youngjin/gig/tests/testdata/typeconv"
	"git.woa.com/youngjin/gig/tests/testdata/variables"
)

// Individual test functions for each package - allows running specific suites
func TestCorrectnessAlgorithms(t *testing.T)       { runTestSet(t, testSetsMap["algorithms"]) }
func TestCorrectnessAdvanced(t *testing.T)         { runTestSet(t, testSetsMap["advanced"]) }
func TestCorrectnessArithmetic(t *testing.T)       { runTestSet(t, testSetsMap["arithmetic"]) }
func TestCorrectnessAutowrap(t *testing.T)         { runTestSet(t, testSetsMap["autowrap"]) }
func TestCorrectnessBitwise(t *testing.T)          { runTestSet(t, testSetsMap["bitwise"]) }
func TestCorrectnessClosures(t *testing.T)         { runTestSet(t, testSetsMap["closures"]) }
func TestCorrectnessClosuresAdvanced(t *testing.T) { runTestSet(t, testSetsMap["closures_advanced"]) }
func TestCorrectnessControlflow(t *testing.T)      { runTestSet(t, testSetsMap["controlflow"]) }
func TestCorrectnessCornercases(t *testing.T) {
	runTestSet(t, testSetsMap["cornercases"])
}
func TestCorrectnessEdgecases(t *testing.T)     { runTestSet(t, testSetsMap["edgecases"]) }
func TestCorrectnessExternal(t *testing.T)      { runTestSet(t, testSetsMap["external"]) }
func TestCorrectnessFunctions(t *testing.T)     { runTestSet(t, testSetsMap["functions"]) }
func TestCorrectnessGoroutine(t *testing.T)     { runTestSet(t, testSetsMap["goroutine"]) }
func TestCorrectnessChannels(t *testing.T)       { runTestSet(t, testSetsMap["channels"]) }
func TestCorrectnessInit(t *testing.T)          { runTestSet(t, testSetsMap["init"]) }
func TestCorrectnessLeetcodeHard(t *testing.T)  { runTestSet(t, testSetsMap["leetcode_hard"]) }
func TestCorrectnessMaps(t *testing.T)          { runTestSet(t, testSetsMap["maps"]) }
func TestCorrectnessMapadvanced(t *testing.T)   { runTestSet(t, testSetsMap["mapadvanced"]) }
func TestCorrectnessMultiassign(t *testing.T)   { runTestSet(t, testSetsMap["multiassign"]) }
func TestCorrectnessNamedreturn(t *testing.T)   { runTestSet(t, testSetsMap["namedreturn"]) }
func TestCorrectnessRecursion(t *testing.T)     { runTestSet(t, testSetsMap["recursion"]) }
func TestCorrectnessResolvedIssue(t *testing.T) { runTestSet(t, testSetsMap["resolved_issue"]) }
func TestCorrectnessScope(t *testing.T)         { runTestSet(t, testSetsMap["scope"]) }
func TestCorrectnessSlices(t *testing.T)        { runTestSet(t, testSetsMap["slices"]) }
func TestCorrectnessSlicing(t *testing.T)       { runTestSet(t, testSetsMap["slicing"]) }
func TestCorrectnessStringsPkg(t *testing.T)    { runTestSet(t, testSetsMap["strings_pkg"]) }
func TestCorrectnessStructs(t *testing.T)       { runTestSet(t, testSetsMap["structs"]) }
func TestCorrectnessSwitch(t *testing.T)        { runTestSet(t, testSetsMap["switch"]) }
func TestCorrectnessTypeconv(t *testing.T)      { runTestSet(t, testSetsMap["typeconv"]) }
func TestCorrectnessVariables(t *testing.T)     { runTestSet(t, testSetsMap["variables"]) }

// Tricky tests - split by category
func TestCorrectnessTrickyClosures(t *testing.T)    { runTestSet(t, testSetsMap["tricky/closures"]) }
func TestCorrectnessTrickyDefer(t *testing.T)       { runTestSet(t, testSetsMap["tricky/defer"]) }
func TestCorrectnessTrickyInterfaces(t *testing.T)  { runTestSet(t, testSetsMap["tricky/interfaces"]) }
func TestCorrectnessTrickyMaps(t *testing.T)        { runTestSet(t, testSetsMap["tricky/maps"]) }
func TestCorrectnessTrickyMultiassign(t *testing.T) { runTestSet(t, testSetsMap["tricky/multiassign"]) }
func TestCorrectnessTrickyNested(t *testing.T)      { runTestSet(t, testSetsMap["tricky/nested"]) }
func TestCorrectnessTrickyPointers(t *testing.T)    { runTestSet(t, testSetsMap["tricky/pointers"]) }
func TestCorrectnessTrickySlices(t *testing.T)      { runTestSet(t, testSetsMap["tricky/slices"]) }
func TestCorrectnessTrickyStructs(t *testing.T)     { runTestSet(t, testSetsMap["tricky/structs"]) }

func TestCornerCase(t *testing.T) { runTestSet(t, testSetsMap["cornercases"]) }

// ============================================================================
// Helper Functions
// ============================================================================

// testCase defines a single test case
type testCase struct {
	src      string
	funcName string
	args     []any
	native   any // native function, called via reflection with args
}

// callNative invokes fn with args using reflection and returns the result.
// For multi-return functions, results are wrapped in []any (matching interpreter behavior).
func callNative(fn any, args []any) any {
	v := reflect.ValueOf(fn)
	in := make([]reflect.Value, len(args))
	for i, a := range args {
		in[i] = reflect.ValueOf(a)
	}
	out := v.Call(in)
	if len(out) == 1 {
		return out[0].Interface()
	}
	result := make([]any, len(out))
	for i, o := range out {
		result[i] = o.Interface()
	}
	return result
}

// compareResults compares interpreter result with native result
func compareCorrectnessResults(t *testing.T, got, expected any) {
	t.Helper()

	// Handle nil cases
	if expected == nil {
		if got != nil {
			t.Errorf("expected nil, got %v (%T)", got, got)
		}
		return
	}

	// Handle []any (multiple return values)
	if expSlice, ok := expected.([]any); ok {
		if gotSlice, ok := got.([]any); ok && len(gotSlice) == len(expSlice) {
			for i := range gotSlice {
				compareCorrectnessResults(t, gotSlice[i], expSlice[i])
			}
			return
		}
	}

	// Deep equality check
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("mismatch:\n  got:      %v (%T)\n  expected: %v (%T)", got, got, expected, expected)
	}
}

// ============================================================================
// Test Runner
// ============================================================================

// testSet represents a group of tests that share the same source file
type testSet struct {
	name  string
	src   string
	tests map[string]testCase // key is just funcName, not full path
}

// progCache caches compiled programs by source to avoid recompilation
var progCache = make(map[string]*gig.Program)

// runTestSet runs all tests in a test set, compiling the source once
func runTestSet(t *testing.T, set testSet) {
	t.Helper()
	prog, ok := progCache[set.src]
	if !ok {
		var err error
		prog, err = gig.Build(set.src)
		if err != nil {
			t.Fatalf("Build error: %v", err)
		}
		progCache[set.src] = prog
	}

	for fullKey, tc := range set.tests {
		// Use fullKey for test name (includes package), tc.funcName for actual call
		t.Run(fullKey, func(t *testing.T) {
			startInterp := time.Now()
			result, err := prog.Run(tc.funcName, tc.args...)
			interpDuration := time.Since(startInterp)
			if err != nil {
				t.Fatalf("Run error: %v", err)
			}

			// Skip native comparison when native is nil (e.g., package main sources)
			if tc.native != nil {
				startNative := time.Now()
				expected := callNative(tc.native, tc.args)
				nativeDuration := time.Since(startNative)

				compareCorrectnessResults(t, result, expected)

				ratio := float64(interpDuration) / float64(nativeDuration)
				t.Logf("interp: %v, native(using reflection): %v, ratio: %.1fx", interpDuration, nativeDuration, ratio)
			} else {
				t.Logf("interp: %v (no native comparison for package main source)", interpDuration)
			}
		})
	}
}

// ============================================================================
// Embedded Source Files
// ============================================================================

//go:embed testdata/algorithms/main.go
var algorithmsSrc string

//go:embed testdata/advanced/main.go
var advancedSrc string

//go:embed testdata/arithmetic/main.go
var arithmeticSrc string

//go:embed testdata/autowrap/main.go
var autowrapSrc string

//go:embed testdata/bitwise/main.go
var bitwiseSrc string

//go:embed testdata/closures/main.go
var closuresSrc string

//go:embed testdata/closures_advanced/main.go
var closuresAdvancedSrc string

//go:embed testdata/controlflow/main.go
var controlflowSrc string

//go:embed testdata/cornercases/main.go
var cornercasesSrc string

//go:embed testdata/edgecases/main.go
var edgecasesSrc string

//go:embed testdata/external/main.go
var externalSrc string

//go:embed testdata/functions/main.go
var functionsSrc string

//go:embed testdata/goroutine/main.go
var goroutineSrc string

//go:embed testdata/channels/main.go
var channelsSrc string

//go:embed testdata/initialize/main.go
var initializeSrc string

//go:embed testdata/leetcode_hard/main.go
var leetcodeHardSrc string

//go:embed testdata/maps/main.go
var mapsSrc string

//go:embed testdata/mapadvanced/main.go
var mapadvancedSrc string

//go:embed testdata/multiassign/main.go
var multiassignSrc string

//go:embed testdata/namedreturn/main.go
var namedreturnSrc string

//go:embed testdata/recursion/main.go
var recursionSrc string

//go:embed testdata/resolved_issue/main.go
var resolvedIssueSrc string

//go:embed testdata/scope/main.go
var scopeSrc string

//go:embed testdata/slices/main.go
var slicesSrc string

//go:embed testdata/slicing/main.go
var slicingSrc string

//go:embed testdata/strings_pkg/main.go
var stringsPkgSrc string

//go:embed testdata/structs/main.go
var structsSrc string

//go:embed testdata/switch/main.go
var switchSrc string

//go:embed testdata/tricky/main.go
var trickySrc string

//go:embed testdata/typeconv/main.go
var typeconvSrc string

//go:embed testdata/variables/main.go
var variablesSrc string

//go:embed testdata/initialize/main.go
var initSrc string

// ============================================================================
// All Test Cases - Consolidated from all test files
// ============================================================================

var advancedTests = map[string]testCase{
	"TypeConvertIntIdentity": {advancedSrc, "TypeConvertIntIdentity", nil, advanced.TypeConvertIntIdentity},
	"DeepCallChain":          {advancedSrc, "DeepCallChain", nil, advanced.DeepCallChain},
	"EarlyReturn":            {advancedSrc, "EarlyReturn", nil, advanced.EarlyReturn},
	"NestedIfInLoop":         {advancedSrc, "NestedIfInLoop", nil, advanced.NestedIfInLoop},
	"BubbleSort":             {advancedSrc, "BubbleSort", nil, advanced.BubbleSort},
	"BinarySearch":           {advancedSrc, "BinarySearch", nil, advanced.BinarySearch},
	"GCD":                    {advancedSrc, "GCD", nil, advanced.GCD},
	"SieveOfEratosthenes":    {advancedSrc, "SieveOfEratosthenes", nil, advanced.SieveOfEratosthenes},
	"MatrixMultiply":         {advancedSrc, "MatrixMultiply", nil, advanced.MatrixMultiply},
	"EmptyFunctionReturn":    {advancedSrc, "EmptyFunctionReturn", nil, advanced.EmptyFunctionReturn},
	"SingleReturnValue":      {advancedSrc, "SingleReturnValue", nil, advanced.SingleReturnValue},
	"ZeroIteration":          {advancedSrc, "ZeroIteration", nil, advanced.ZeroIteration},
	"LargeLoop":              {advancedSrc, "LargeLoop", nil, advanced.LargeLoop},
	"DeepRecursion":          {advancedSrc, "DeepRecursion", nil, advanced.DeepRecursion},
	"MapWithClosure":         {advancedSrc, "MapWithClosure", nil, advanced.MapWithClosure},
	"SliceWithMultiReturn":   {advancedSrc, "SliceWithMultiReturn", nil, advanced.SliceWithMultiReturn},
	"RecursiveDataBuild":     {advancedSrc, "RecursiveDataBuild", nil, advanced.RecursiveDataBuild},
	"FunctionChain":          {advancedSrc, "FunctionChain", nil, advanced.FunctionChain},
	"ComplexExpressions":     {advancedSrc, "ComplexExpressions", nil, advanced.ComplexExpressions},
	// Parameterized tests
	"FindFirst": {advancedSrc, "FindFirst", []any{[]int{10, 20, 30}, 20}, advanced.FindFirst},
	"Bsearch":   {advancedSrc, "Bsearch", []any{[]int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}, 50}, advanced.Bsearch},
	"Gcd":       {advancedSrc, "Gcd", []any{48, 18}, advanced.Gcd},
	"Identity":  {advancedSrc, "Identity", []any{42}, advanced.Identity},
	"Minmax":    {advancedSrc, "Minmax", []any{[]int{3, 1, 4, 1, 5}}, advanced.Minmax},
	"Countdown": {advancedSrc, "Countdown", []any{50}, advanced.Countdown},
	"Add":       {advancedSrc, "Add", []any{1, 2}, advanced.Add},
	"Mul":       {advancedSrc, "Mul", []any{3, 4}, advanced.Mul},
	"Sub":       {advancedSrc, "Sub", []any{10, 5}, advanced.Sub},
}

var algorithmsTests = map[string]testCase{
	"InsertionSort":     {algorithmsSrc, "InsertionSort", nil, algorithms.InsertionSort},
	"SelectionSort":     {algorithmsSrc, "SelectionSort", nil, algorithms.SelectionSort},
	"ReverseSlice":      {algorithmsSrc, "ReverseSlice", nil, algorithms.ReverseSlice},
	"IsPalindrome":      {algorithmsSrc, "IsPalindrome", nil, algorithms.IsPalindrome},
	"PowerFunction":     {algorithmsSrc, "PowerFunction", nil, algorithms.PowerFunction},
	"MaxSubarraySum":    {algorithmsSrc, "MaxSubarraySum", nil, algorithms.MaxSubarraySum},
	"TwoSum":            {algorithmsSrc, "TwoSum", nil, algorithms.TwoSum},
	"FibMemoized":       {algorithmsSrc, "FibMemoized", nil, algorithms.FibMemoized},
	"CountDigits":       {algorithmsSrc, "CountDigits", nil, algorithms.CountDigits},
	"CollatzConjecture": {algorithmsSrc, "CollatzConjecture", nil, algorithms.CollatzConjecture},
	// Parameterized tests
	"Reverse":         {algorithmsSrc, "Reverse", []any{[]int{1, 2, 3, 4, 5}}, algorithms.Reverse},
	"Power":           {algorithmsSrc, "Power", []any{2, 10}, algorithms.Power},
	"CountDigitsN":    {algorithmsSrc, "CountDigitsN", []any{12345}, algorithms.CountDigitsN},
	"CollatzStepsN":   {algorithmsSrc, "CollatzStepsN", []any{27}, algorithms.CollatzStepsN},
	"IsPalindromeInt": {algorithmsSrc, "IsPalindromeInt", []any{[]int{1, 2, 3, 2, 1}}, algorithms.IsPalindromeInt},
}

var arithmeticTests = map[string]testCase{
	"Addition":       {arithmeticSrc, "Addition", nil, arithmetic.Addition},
	"Subtraction":    {arithmeticSrc, "Subtraction", nil, arithmetic.Subtraction},
	"Multiplication": {arithmeticSrc, "Multiplication", nil, arithmetic.Multiplication},
	"Division":       {arithmeticSrc, "Division", nil, arithmetic.Division},
	"Modulo":         {arithmeticSrc, "Modulo", nil, arithmetic.Modulo},
	"ComplexExpr":    {arithmeticSrc, "ComplexExpr", nil, arithmetic.ComplexExpr},
	"Negation":       {arithmeticSrc, "Negation", nil, arithmetic.Negation},
	"ChainedOps":     {arithmeticSrc, "ChainedOps", nil, arithmetic.ChainedOps},
	"Overflow":       {arithmeticSrc, "Overflow", nil, arithmetic.Overflow},
	"Precedence":     {arithmeticSrc, "Precedence", nil, arithmetic.Precedence},
	// Parameterized tests
	"Add":          {arithmeticSrc, "Add", []any{10, 32}, arithmetic.Add},
	"Sub":          {arithmeticSrc, "Sub", []any{100, 42}, arithmetic.Sub},
	"Mul":          {arithmeticSrc, "Mul", []any{6, 7}, arithmetic.Mul},
	"Div":          {arithmeticSrc, "Div", []any{100, 4}, arithmetic.Div},
	"Mod":          {arithmeticSrc, "Mod", []any{17, 5}, arithmetic.Mod},
	"ComplexArith": {arithmeticSrc, "ComplexArith", []any{2, 3, 4, 10, 2}, arithmetic.ComplexArith},
}

var autowrapTests = map[string]testCase{
	"WithPackage": {autowrapSrc, "WithPackage", nil, autowrap.WithPackage},
	"WithImport":  {autowrapSrc, "WithImport", nil, autowrap.WithImport},
	"Compute":     {autowrapSrc, "Compute", nil, autowrap.Compute},
}

var bitwiseTests = map[string]testCase{
	"And":        {bitwiseSrc, "And", nil, bitwise.And},
	"Or":         {bitwiseSrc, "Or", nil, bitwise.Or},
	"Xor":        {bitwiseSrc, "Xor", nil, bitwise.Xor},
	"LeftShift":  {bitwiseSrc, "LeftShift", nil, bitwise.LeftShift},
	"RightShift": {bitwiseSrc, "RightShift", nil, bitwise.RightShift},
	"Combined":   {bitwiseSrc, "Combined", nil, bitwise.Combined},
	"AndNot":     {bitwiseSrc, "AndNot", nil, bitwise.AndNot},
	"PowerOfTwo": {bitwiseSrc, "PowerOfTwo", nil, bitwise.PowerOfTwo},
	// Parameterized tests
	"BitAnd":        {bitwiseSrc, "BitAnd", []any{0xFF, 0x0F}, bitwise.BitAnd},
	"BitOr":         {bitwiseSrc, "BitOr", []any{0xFF, 0x100}, bitwise.BitOr},
	"BitXor":        {bitwiseSrc, "BitXor", []any{0xAA, 0x55}, bitwise.BitXor},
	"BitLeftShift":  {bitwiseSrc, "BitLeftShift", []any{1, 10}, bitwise.BitLeftShift},
	"BitRightShift": {bitwiseSrc, "BitRightShift", []any{1024, 5}, bitwise.BitRightShift},
	"IsPowerOfTwo":  {bitwiseSrc, "IsPowerOfTwo", []any{16}, bitwise.IsPowerOfTwo},
}

var closuresTests = map[string]testCase{
	"Counter":           {closuresSrc, "Counter", nil, closures.Counter},
	"CaptureMutation":   {closuresSrc, "CaptureMutation", nil, closures.CaptureMutation},
	"Factory":           {closuresSrc, "Factory", nil, closures.Factory},
	"MultipleInstances": {closuresSrc, "MultipleInstances", nil, closures.MultipleInstances},
	"OverLoop":          {closuresSrc, "OverLoop", nil, closures.OverLoop},
	"Chain":             {closuresSrc, "Chain", nil, closures.Chain},
	"Accumulator":       {closuresSrc, "Accumulator", nil, closures.Accumulator},
	// Parameterized tests (closures return functions, skip for now as we test via Call)
}

var closures_advancedTests = map[string]testCase{
	"Generator":          {closuresAdvancedSrc, "Generator", nil, closures_advanced.Generator},
	"Predicate":          {closuresAdvancedSrc, "Predicate", nil, closures_advanced.Predicate},
	"StateMachine":       {closuresAdvancedSrc, "StateMachine", nil, closures_advanced.StateMachine},
	"RecursiveHelper":    {closuresAdvancedSrc, "RecursiveHelper", nil, closures_advanced.RecursiveHelper},
	"ApplyN":             {closuresAdvancedSrc, "ApplyN", nil, closures_advanced.ApplyN},
	"Compose":            {closuresAdvancedSrc, "Compose", nil, closures_advanced.Compose},
	"ClosureForLoopTest": {closuresAdvancedSrc, "ClosureForLoopTest", nil, closures_advanced.ClosureForLoopTest},
}

var controlflowTests = map[string]testCase{
	"IfTrue":              {controlflowSrc, "IfTrue", nil, controlflow.IfTrue},
	"IfFalse":             {controlflowSrc, "IfFalse", nil, controlflow.IfFalse},
	"IfElse":              {controlflowSrc, "IfElse", nil, controlflow.IfElse},
	"IfElseChainNegative": {controlflowSrc, "IfElseChainNegative", nil, controlflow.IfElseChainNegative},
	"IfElseChainZero":     {controlflowSrc, "IfElseChainZero", nil, controlflow.IfElseChainZero},
	"IfElseChainPositive": {controlflowSrc, "IfElseChainPositive", nil, controlflow.IfElseChainPositive},
	"ForLoop":             {controlflowSrc, "ForLoop", nil, controlflow.ForLoop},
	"ForConditionOnly":    {controlflowSrc, "ForConditionOnly", nil, controlflow.ForConditionOnly},
	"NestedFor":           {controlflowSrc, "NestedFor", nil, controlflow.NestedFor},
	"ForBreak":            {controlflowSrc, "ForBreak", nil, controlflow.ForBreak},
	"ForContinue":         {controlflowSrc, "ForContinue", nil, controlflow.ForContinue},
	"BooleanAndOr":        {controlflowSrc, "BooleanAndOr", nil, controlflow.BooleanAndOr},
	// Parameterized tests
	"ClassifyNegative": {controlflowSrc, "Classify", []any{-5}, controlflow.Classify},
	"ClassifyZero":     {controlflowSrc, "Classify", []any{0}, controlflow.Classify},
	"ClassifyPositive": {controlflowSrc, "Classify", []any{42}, controlflow.Classify},
}

var cornercasesTests = map[string]testCase{
	"ZeroValue_Int":          {cornercasesSrc, "ZeroValue_Int", nil, cornercases.ZeroValue_Int},
	"ZeroValue_Int64":        {cornercasesSrc, "ZeroValue_Int64", nil, cornercases.ZeroValue_Int64},
	"ZeroValue_Float64":      {cornercasesSrc, "ZeroValue_Float64", nil, cornercases.ZeroValue_Float64},
	"ZeroValue_String":       {cornercasesSrc, "ZeroValue_String", nil, cornercases.ZeroValue_String},
	"ZeroValue_Bool":         {cornercasesSrc, "ZeroValue_Bool", nil, cornercases.ZeroValue_Bool},
	"ZeroValue_Slice":        {cornercasesSrc, "ZeroValue_Slice", nil, cornercases.ZeroValue_Slice},
	"ZeroValue_Map":          {cornercasesSrc, "ZeroValue_Map", nil, cornercases.ZeroValue_Map},
	"IntBoundary_MaxInt32":   {cornercasesSrc, "IntBoundary_MaxInt32", nil, cornercases.IntBoundary_MaxInt32},
	"IntBoundary_MinInt32":   {cornercasesSrc, "IntBoundary_MinInt32", nil, cornercases.IntBoundary_MinInt32},
	"IntBoundary_MaxInt64":   {cornercasesSrc, "IntBoundary_MaxInt64", nil, cornercases.IntBoundary_MaxInt64},
	"IntBoundary_MinInt64":   {cornercasesSrc, "IntBoundary_MinInt64", nil, cornercases.IntBoundary_MinInt64},
	"IntBoundary_MaxUint32":  {cornercasesSrc, "IntBoundary_MaxUint32", nil, cornercases.IntBoundary_MaxUint32},
	"IntBoundary_NearMaxInt": {cornercasesSrc, "IntBoundary_NearMaxInt", nil, cornercases.IntBoundary_NearMaxInt},
	"IntBoundary_NearMinInt": {cornercasesSrc, "IntBoundary_NearMinInt", nil, cornercases.IntBoundary_NearMinInt},
	// Note: int32 overflow wraps just like native Go (two's complement)
	// We compare the actual native results to verify correctness
	"Overflow_Int32Add":           {cornercasesSrc, "Overflow_Int32Add", nil, cornercases.Overflow_Int32Add},
	"Overflow_Int32Sub":           {cornercasesSrc, "Overflow_Int32Sub", nil, cornercases.Overflow_Int32Sub},
	"Overflow_Int32Mul":           {cornercasesSrc, "Overflow_Int32Mul", nil, cornercases.Overflow_Int32Mul},
	"FloatBoundary_SmallPositive": {cornercasesSrc, "FloatBoundary_SmallPositive", nil, cornercases.FloatBoundary_SmallPositive},
	"FloatBoundary_SmallNegative": {cornercasesSrc, "FloatBoundary_SmallNegative", nil, cornercases.FloatBoundary_SmallNegative},
	"FloatBoundary_LargePositive": {cornercasesSrc, "FloatBoundary_LargePositive", nil, cornercases.FloatBoundary_LargePositive},
	"FloatBoundary_LargeNegative": {cornercasesSrc, "FloatBoundary_LargeNegative", nil, cornercases.FloatBoundary_LargeNegative},
	"EmptySlice_Len":              {cornercasesSrc, "EmptySlice_Len", nil, cornercases.EmptySlice_Len},
	"EmptySlice_Cap":              {cornercasesSrc, "EmptySlice_Cap", nil, cornercases.EmptySlice_Cap},
	"EmptySlice_Make":             {cornercasesSrc, "EmptySlice_Make", nil, cornercases.EmptySlice_Make},
	"EmptyMap_Len":                {cornercasesSrc, "EmptyMap_Len", nil, cornercases.EmptyMap_Len},
	"EmptyMap_Make":               {cornercasesSrc, "EmptyMap_Make", nil, cornercases.EmptyMap_Make},
	"EmptyString_Len":             {cornercasesSrc, "EmptyString_Len", nil, cornercases.EmptyString_Len},
	"Slice_ZeroToZero":            {cornercasesSrc, "Slice_ZeroToZero", nil, cornercases.Slice_ZeroToZero},
	"Slice_EndToEnd":              {cornercasesSrc, "Slice_EndToEnd", nil, cornercases.Slice_EndToEnd},
	"Slice_NilSlice":              {cornercasesSrc, "Slice_NilSlice", nil, cornercases.Slice_NilSlice},
	"Slice_AppendToNil":           {cornercasesSrc, "Slice_AppendToNil", nil, cornercases.Slice_AppendToNil},
	"Slice_AppendEmpty":           {cornercasesSrc, "Slice_AppendEmpty", nil, cornercases.Slice_AppendEmpty},
	"Map_NilMap":                  {cornercasesSrc, "Map_NilMap", nil, cornercases.Map_NilMap},
	"Map_AccessMissingKey":        {cornercasesSrc, "Map_AccessMissingKey", nil, cornercases.Map_AccessMissingKey},
	"Map_DeleteMissingKey":        {cornercasesSrc, "Map_DeleteMissingKey", nil, cornercases.Map_DeleteMissingKey},
	"Map_OverwriteKey":            {cornercasesSrc, "Map_OverwriteKey", nil, cornercases.Map_OverwriteKey},
	"Map_NilKeyString":            {cornercasesSrc, "Map_NilKeyString", nil, cornercases.Map_NilKeyString},
	"Map_ZeroIntKey":              {cornercasesSrc, "Map_ZeroIntKey", nil, cornercases.Map_ZeroIntKey},
	"String_Empty":                {cornercasesSrc, "String_Empty", nil, cornercases.String_Empty},
	"String_SingleChar":           {cornercasesSrc, "String_SingleChar", nil, cornercases.String_SingleChar},
	"String_UnicodeMultibyte":     {cornercasesSrc, "String_UnicodeMultibyte", nil, cornercases.String_UnicodeMultibyte},
	"String_Whitespace":           {cornercasesSrc, "String_Whitespace", nil, cornercases.String_Whitespace},
	"String_SingleByteIndex":      {cornercasesSrc, "String_SingleByteIndex", nil, cornercases.String_SingleByteIndex},
	"String_LastByte":             {cornercasesSrc, "String_LastByte", nil, cornercases.String_LastByte},
	"Bool_True":                   {cornercasesSrc, "Bool_True", nil, cornercases.Bool_True},
	"Bool_False":                  {cornercasesSrc, "Bool_False", nil, cornercases.Bool_False},
	"Bool_NotTrue":                {cornercasesSrc, "Bool_NotTrue", nil, cornercases.Bool_NotTrue},
	"Bool_NotFalse":               {cornercasesSrc, "Bool_NotFalse", nil, cornercases.Bool_NotFalse},
	"Bool_DoubleNegation":         {cornercasesSrc, "Bool_DoubleNegation", nil, cornercases.Bool_DoubleNegation},
	"Arith_AddZero":               {cornercasesSrc, "Arith_AddZero", nil, cornercases.Arith_AddZero},
	"Arith_SubZero":               {cornercasesSrc, "Arith_SubZero", nil, cornercases.Arith_SubZero},
	"Arith_MulByOne":              {cornercasesSrc, "Arith_MulByOne", nil, cornercases.Arith_MulByOne},
	"Arith_DivByOne":              {cornercasesSrc, "Arith_DivByOne", nil, cornercases.Arith_DivByOne},
	"Arith_ModByOne":              {cornercasesSrc, "Arith_ModByOne", nil, cornercases.Arith_ModByOne},
	"Arith_MulByZero":             {cornercasesSrc, "Arith_MulByZero", nil, cornercases.Arith_MulByZero},
	"Arith_NegNeg":                {cornercasesSrc, "Arith_NegNeg", nil, cornercases.Arith_NegNeg},
	"Arith_NegAddNeg":             {cornercasesSrc, "Arith_NegAddNeg", nil, cornercases.Arith_NegAddNeg},
	"Compare_IntEqual":            {cornercasesSrc, "Compare_IntEqual", nil, cornercases.Compare_IntEqual},
	"Compare_IntNotEqual":         {cornercasesSrc, "Compare_IntNotEqual", nil, cornercases.Compare_IntNotEqual},
	"Compare_IntGreater":          {cornercasesSrc, "Compare_IntGreater", nil, cornercases.Compare_IntGreater},
	"Compare_IntGreaterEqual":     {cornercasesSrc, "Compare_IntGreaterEqual", nil, cornercases.Compare_IntGreaterEqual},
	"Compare_IntLess":             {cornercasesSrc, "Compare_IntLess", nil, cornercases.Compare_IntLess},
	"Compare_IntLessEqual":        {cornercasesSrc, "Compare_IntLessEqual", nil, cornercases.Compare_IntLessEqual},
	"Compare_StringEqual":         {cornercasesSrc, "Compare_StringEqual", nil, cornercases.Compare_StringEqual},
	"Compare_StringNotEqual":      {cornercasesSrc, "Compare_StringNotEqual", nil, cornercases.Compare_StringNotEqual},
	"Compare_EmptyStringEqual":    {cornercasesSrc, "Compare_EmptyStringEqual", nil, cornercases.Compare_EmptyStringEqual},
	"Logic_TrueAndTrue":           {cornercasesSrc, "Logic_TrueAndTrue", nil, cornercases.Logic_TrueAndTrue},
	"Logic_TrueAndFalse":          {cornercasesSrc, "Logic_TrueAndFalse", nil, cornercases.Logic_TrueAndFalse},
	"Logic_FalseAndTrue":          {cornercasesSrc, "Logic_FalseAndTrue", nil, cornercases.Logic_FalseAndTrue},
	"Logic_TrueOrFalse":           {cornercasesSrc, "Logic_TrueOrFalse", nil, cornercases.Logic_TrueOrFalse},
	"Logic_FalseOrTrue":           {cornercasesSrc, "Logic_FalseOrTrue", nil, cornercases.Logic_FalseOrTrue},
	"Logic_FalseOrFalse":          {cornercasesSrc, "Logic_FalseOrFalse", nil, cornercases.Logic_FalseOrFalse},
	"Control_IfNoElse":            {cornercasesSrc, "Control_IfNoElse", nil, cornercases.Control_IfNoElse},
	"Control_IfFalseNoElse":       {cornercasesSrc, "Control_IfFalseNoElse", nil, cornercases.Control_IfFalseNoElse},
	"Control_ForZeroIter":         {cornercasesSrc, "Control_ForZeroIter", nil, cornercases.Control_ForZeroIter},
	"Control_ForOneIter":          {cornercasesSrc, "Control_ForOneIter", nil, cornercases.Control_ForOneIter},
	"Control_ForBreakFirst":       {cornercasesSrc, "Control_ForBreakFirst", nil, cornercases.Control_ForBreakFirst},
	"Control_ForContinueAll":      {cornercasesSrc, "Control_ForContinueAll", nil, cornercases.Control_ForContinueAll},
	"Control_SwitchNoMatch":       {cornercasesSrc, "Control_SwitchNoMatch", nil, cornercases.Control_SwitchNoMatch},
	"Control_SwitchDefault":       {cornercasesSrc, "Control_SwitchDefault", nil, cornercases.Control_SwitchDefault},
	"Func_NoReturn":               {cornercasesSrc, "Func_NoReturn", nil, cornercases.Func_NoReturn},
	"Func_MultipleReturnAll":      {cornercasesSrc, "Func_MultipleReturnAll", nil, cornercases.Func_MultipleReturnAll},
	"Func_MultipleReturnIgnore":   {cornercasesSrc, "Func_MultipleReturnIgnore", nil, cornercases.Func_MultipleReturnIgnore},
	"Func_NamedReturn":            {cornercasesSrc, "Func_NamedReturn", nil, cornercases.Func_NamedReturn},
	"Func_VariadicEmpty":          {cornercasesSrc, "Func_VariadicEmpty", nil, cornercases.Func_VariadicEmpty},
	"Func_VariadicOne":            {cornercasesSrc, "Func_VariadicOne", nil, cornercases.Func_VariadicOne},
	"Func_VariadicMultiple":       {cornercasesSrc, "Func_VariadicMultiple", nil, cornercases.Func_VariadicMultiple},
	"Func_RecursionBase":          {cornercasesSrc, "Func_RecursionBase", nil, cornercases.Func_RecursionBase},
	"Closure_ReturnClosure":       {cornercasesSrc, "Closure_ReturnClosure", nil, cornercases.Closure_ReturnClosure},
	"Closure_CaptureVariable":     {cornercasesSrc, "Closure_CaptureVariable", nil, cornercases.Closure_CaptureVariable},
	"Closure_ModifyCaptured":      {cornercasesSrc, "Closure_ModifyCaptured", nil, cornercases.Closure_ModifyCaptured},
	"Struct_ZeroValueFields":      {cornercasesSrc, "Struct_ZeroValueFields", nil, cornercases.Struct_ZeroValueFields},
	"Struct_PointerReceiver":      {cornercasesSrc, "Struct_PointerReceiver", nil, cornercases.Struct_PointerReceiver},
	"Struct_NestedStruct":         {cornercasesSrc, "Struct_NestedStruct", nil, cornercases.Struct_NestedStruct},
}

var edgecasesTests = map[string]testCase{
	"MaxInt64":           {edgecasesSrc, "MaxInt64", nil, edgecases.MaxInt64},
	"MinInt64":           {edgecasesSrc, "MinInt64", nil, edgecases.MinInt64},
	"DivisionByMinusOne": {edgecasesSrc, "DivisionByMinusOne", nil, edgecases.DivisionByMinusOne},
	"ModuloNegative":     {edgecasesSrc, "ModuloNegative", nil, edgecases.ModuloNegative},
	"EmptyString":        {edgecasesSrc, "EmptyString", nil, edgecases.EmptyString},
	"LargeSlice":         {edgecasesSrc, "LargeSlice", nil, edgecases.LargeSlice},
	"NestedMapLookup":    {edgecasesSrc, "NestedMapLookup", nil, edgecases.NestedMapLookup},
	"ZeroDivisionGuard":  {edgecasesSrc, "ZeroDivisionGuard", nil, edgecases.ZeroDivisionGuard},
	"BooleanComplexExpr": {edgecasesSrc, "BooleanComplexExpr", nil, edgecases.BooleanComplexExpr},
	"SingleElementSlice": {edgecasesSrc, "SingleElementSlice", nil, edgecases.SingleElementSlice},
	"EmptyMap":           {edgecasesSrc, "EmptyMap", nil, edgecases.EmptyMap},
	"TightLoop":          {edgecasesSrc, "TightLoop", nil, edgecases.TightLoop},
}

var externalTests = map[string]testCase{
	"FmtSprintf":       {externalSrc, "FmtSprintf", nil, external.FmtSprintf},
	"FmtSprintfMulti":  {externalSrc, "FmtSprintfMulti", nil, external.FmtSprintfMulti},
	"StringsToUpper":   {externalSrc, "StringsToUpper", nil, external.StringsToUpper},
	"StringsToLower":   {externalSrc, "StringsToLower", nil, external.StringsToLower},
	"StringsContains":  {externalSrc, "StringsContains", nil, external.StringsContains},
	"StringsReplace":   {externalSrc, "StringsReplace", nil, external.StringsReplace},
	"StringsHasPrefix": {externalSrc, "StringsHasPrefix", nil, external.StringsHasPrefix},
	"StrconvItoa":      {externalSrc, "StrconvItoa", nil, external.StrconvItoa},
	"StrconvAtoi":      {externalSrc, "StrconvAtoi", nil, external.StrconvAtoi},
	// Parameterized tests
	"FmtSprintfInt":      {externalSrc, "FmtSprintfInt", []any{42}, external.FmtSprintfInt},
	"StringsToUpperStr":  {externalSrc, "StringsToUpperStr", []any{"hello"}, external.StringsToUpperStr},
	"StringsToLowerStr":  {externalSrc, "StringsToLowerStr", []any{"HELLO"}, external.StringsToLowerStr},
	"StringsContainsStr": {externalSrc, "StringsContainsStr", []any{"hello world", "world"}, external.StringsContainsStr},
	"StrconvItoaN":       {externalSrc, "StrconvItoaN", []any{42}, external.StrconvItoaN},
	"StrconvAtoiStr":     {externalSrc, "StrconvAtoiStr", []any{"123"}, external.StrconvAtoiStr},
}

var functionsTests = map[string]testCase{
	"Call":                 {functionsSrc, "Call", nil, functions.Call},
	"MultipleReturn":       {functionsSrc, "MultipleReturn", nil, functions.MultipleReturn},
	"MultipleReturnDivmod": {functionsSrc, "MultipleReturnDivmod", nil, functions.MultipleReturnDivmod},
	"RecursionFactorial":   {functionsSrc, "RecursionFactorial", nil, functions.RecursionFactorial},
	"MutualRecursion":      {functionsSrc, "MutualRecursion", nil, functions.MutualRecursion},
	"FibonacciIterative":   {functionsSrc, "FibonacciIterative", nil, functions.FibonacciIterative},
	"FibonacciRecursive":   {functionsSrc, "FibonacciRecursive", nil, functions.FibonacciRecursive},
	"VariadicFunction":     {functionsSrc, "VariadicFunction", nil, functions.VariadicFunction},
	"FunctionAsValue":      {functionsSrc, "FunctionAsValue", nil, functions.FunctionAsValue},
	"HigherOrderMap":       {functionsSrc, "HigherOrderMap", nil, functions.HigherOrderMap},
	"HigherOrderFilter":    {functionsSrc, "HigherOrderFilter", nil, functions.HigherOrderFilter},
	"HigherOrderReduce":    {functionsSrc, "HigherOrderReduce", nil, functions.HigherOrderReduce},
	// Parameterized tests
	"Add":        {functionsSrc, "Add", []any{5, 7}, functions.Add},
	"Swap":       {functionsSrc, "Swap", []any{3, 7}, functions.Swap},
	"Divmod":     {functionsSrc, "Divmod", []any{17, 5}, functions.Divmod},
	"FactorialN": {functionsSrc, "FactorialN", []any{5}, functions.FactorialN},
	"FibIterN":   {functionsSrc, "FibIterN", []any{20}, functions.FibIterN},
	"FibRecN":    {functionsSrc, "FibRecN", []any{15}, functions.FibRecN},
	"IsEvenN":    {functionsSrc, "IsEvenN", []any{10}, functions.IsEvenN},
	"IsOddN":     {functionsSrc, "IsOddN", []any{7}, functions.IsOddN},
	// Variadic with args - skip for now (interpreter variadic handling from outside needs work)

	// Multi-return value tests - these functions THEMSELVES return multiple values
	"ThreeReturnValues":               {functionsSrc, "ThreeReturnValues", nil, functions.ThreeReturnValues},
	"FourReturnValues":                {functionsSrc, "FourReturnValues", nil, functions.FourReturnValues},
	"FiveReturnValues":                {functionsSrc, "FiveReturnValues", nil, functions.FiveReturnValues},
	"MixedTypeReturn":                 {functionsSrc, "MixedTypeReturn", nil, functions.MixedTypeReturn},
	"PassMultiReturnToFunc":           {functionsSrc, "PassMultiReturnToFunc", nil, functions.PassMultiReturnToFunc},
	"ChainMultiReturn":                {functionsSrc, "ChainMultiReturn", nil, functions.ChainMultiReturn},
	"NestedMultiReturn":               {functionsSrc, "NestedMultiReturn", nil, functions.NestedMultiReturn},
	"MultiReturnAsSliceIndex":         {functionsSrc, "MultiReturnAsSliceIndex", nil, functions.MultiReturnAsSliceIndex},
	"MultiReturnToMap":                {functionsSrc, "MultiReturnToMap", nil, functions.MultiReturnToMap},
	"MultiReturnAsCondition":          {functionsSrc, "MultiReturnAsCondition", nil, functions.MultiReturnAsCondition},
	"MultiReturnComplexTypes":         {functionsSrc, "MultiReturnComplexTypes", nil, functions.MultiReturnComplexTypes},
	"MultiReturnInClosure":            {functionsSrc, "MultiReturnInClosure", nil, functions.MultiReturnInClosure},
	"AssignMultiReturnToExistingVars": {functionsSrc, "AssignMultiReturnToExistingVars", nil, functions.AssignMultiReturnToExistingVars},
}

var goroutineTests = map[string]testCase{
	"BasicSpawn":                  {goroutineSrc, "BasicSpawn", nil, goroutine.BasicSpawn},
	"ChannelCommunication":        {goroutineSrc, "ChannelCommunication", nil, goroutine.ChannelCommunication},
	"WithArguments":               {goroutineSrc, "WithArguments", nil, goroutine.WithArguments},
	"WithStruct":                  {goroutineSrc, "WithStruct", nil, goroutine.WithStruct},
	"DifferentTypes":              {goroutineSrc, "DifferentTypes", nil, goroutine.DifferentTypes},
	"GlobalsSharing":              {goroutineSrc, "GlobalsSharing", nil, goroutine.GlobalsSharing},
	"MultipleSends":               {goroutineSrc, "MultipleSends", nil, goroutine.MultipleSends},
	"ParallelExecution":           {goroutineSrc, "ParallelExecution", nil, goroutine.ParallelExecution},
	"ClosureCapture":              {goroutineSrc, "ClosureCapture", nil, goroutine.ClosureCapture},
	"ClosureCaptureMultiple":      {goroutineSrc, "ClosureCaptureMultiple", nil, goroutine.ClosureCaptureMultiple},
	"SelectStatement":             {goroutineSrc, "SelectStatement", nil, goroutine.SelectStatement},
	"SelectDefault":               {goroutineSrc, "SelectDefault", nil, goroutine.SelectDefault},
	"SelectSend":                  {goroutineSrc, "SelectSend", nil, goroutine.SelectSend},
	"RangeOverChannel":            {goroutineSrc, "RangeOverChannel", nil, goroutine.RangeOverChannel},
	"RangeOverChannelWithBuiltin": {goroutineSrc, "RangeOverChannelWithBuiltin", nil, goroutine.RangeOverChannelWithBuiltin},
}

var channelsTests = map[string]testCase{
	"ChannelBasic":             {channelsSrc, "ChannelBasic", nil, channels.ChannelBasic},
	"ChannelBuffered":          {channelsSrc, "ChannelBuffered", nil, channels.ChannelBuffered},
	"ChannelUnbuffered":        {channelsSrc, "ChannelUnbuffered", nil, channels.ChannelUnbuffered},
	"ChannelClose":             {channelsSrc, "ChannelClose", nil, channels.ChannelClose},
	"ChannelNil":               {channelsSrc, "ChannelNil", nil, channels.ChannelNil},
	"SelectDefault":            {channelsSrc, "SelectDefault", nil, channels.SelectDefault},
	"SelectSingleCase":         {channelsSrc, "SelectSingleCase", nil, channels.SelectSingleCase},
	"SelectMultiCase":          {channelsSrc, "SelectMultiCase", nil, channels.SelectMultiCase},
	"SelectSendReceive":        {channelsSrc, "SelectSendReceive", nil, channels.SelectSendReceive},
	"SelectLoop":               {channelsSrc, "SelectLoop", nil, channels.SelectLoop},
	"SelectMultipleChannels":   {channelsSrc, "SelectMultipleChannels", nil, channels.SelectMultipleChannels},
	"ChannelDirectionSend":    {channelsSrc, "ChannelDirectionSend", nil, channels.ChannelDirectionSend},
	"ChannelDirectionReceive":   {channelsSrc, "ChannelDirectionReceive", nil, channels.ChannelDirectionReceive},
	"ChannelStruct":            {channelsSrc, "ChannelStruct", nil, channels.ChannelStruct},
	"ChannelStructPointer":     {channelsSrc, "ChannelStructPointer", nil, channels.ChannelStructPointer},
	"SliceOfChannels":          {channelsSrc, "SliceOfChannels", nil, channels.SliceOfChannels},
	"MapOfChannels":             {channelsSrc, "MapOfChannels", nil, channels.MapOfChannels},
	"ChannelDeadlock":          {channelsSrc, "ChannelDeadlock", nil, channels.ChannelDeadlock},
	"SelectAllBlocked":         {channelsSrc, "SelectAllBlocked", nil, channels.SelectAllBlocked},
	"SelectClosedChannel":      {channelsSrc, "SelectClosedChannel", nil, channels.SelectClosedChannel},
	"SelectNilChannel":         {channelsSrc, "SelectNilChannel", nil, channels.SelectNilChannel},
	"ChannelPipeline":          {channelsSrc, "ChannelPipeline", nil, channels.ChannelPipeline},
	"SelectWithAssignment":     {channelsSrc, "SelectWithAssignment", nil, channels.SelectWithAssignment},
	"SelectBreak":              {channelsSrc, "SelectBreak", nil, channels.SelectBreak},
	"SelectContinue":           {channelsSrc, "SelectContinue", nil, channels.SelectContinue},
	"ChannelFullCap":           {channelsSrc, "ChannelFullCap", nil, channels.ChannelFullCap},
	"ChannelEmptyCap":          {channelsSrc, "ChannelEmptyCap", nil, channels.ChannelEmptyCap},
	"SelectMutex":              {channelsSrc, "SelectMutex", nil, channels.SelectMutex},
	"ChannelTwoWay":            {channelsSrc, "ChannelTwoWay", nil, channels.ChannelTwoWay},
	"ChannelFanIn":             {channelsSrc, "ChannelFanIn", nil, channels.ChannelFanIn},
	"ChannelBufferedsize":       {channelsSrc, "ChannelBufferedsize", nil, channels.ChannelBufferedsize},
}

var initTests = map[string]testCase{
	"GetA":                {initSrc, "GetA", nil, initialize.GetA},
	"GetB":                {initSrc, "GetB", nil, initialize.GetB},
	"GetC":                {initSrc, "GetC", nil, initialize.GetC},
	"GetCacheSum":         {initSrc, "GetCacheSum", nil, initialize.GetCacheSum},
	"GetCacheSize":        {initSrc, "GetCacheSize", nil, initialize.GetCacheSize},
	"GetCacheOrder":       {initSrc, "GetCacheOrder", nil, initialize.GetCacheOrder},
	"GetFibonacciCount":   {initSrc, "GetFibonacciCount", nil, initialize.GetFibonacciCount},
	"GetFibonacciSum":     {initSrc, "GetFibonacciSum", nil, initialize.GetFibonacciSum},
	"ComplexInitTest":     {initSrc, "ComplexInitTest", nil, initialize.ComplexInitTest},
	"InitOrderTest":       {initSrc, "InitOrderTest", nil, initialize.InitOrderTest},
	"CacheInitTest":       {initSrc, "CacheInitTest", nil, initialize.CacheInitTest},
	"LookupTableInitTest": {initSrc, "LookupTableInitTest", nil, initialize.LookupTableInitTest},
	"FibonacciInitTest":   {initSrc, "FibonacciInitTest", nil, initialize.FibonacciInitTest},
}

var initializeTests = map[string]testCase{
	"ComplexInitTest":     {initializeSrc, "ComplexInitTest", nil, initialize.ComplexInitTest},
	"InitOrderTest":       {initializeSrc, "InitOrderTest", nil, initialize.InitOrderTest},
	"CacheInitTest":       {initializeSrc, "CacheInitTest", nil, initialize.CacheInitTest},
	"LookupTableInitTest": {initializeSrc, "LookupTableInitTest", nil, initialize.LookupTableInitTest},
	"FibonacciInitTest":   {initializeSrc, "FibonacciInitTest", nil, initialize.FibonacciInitTest},
	"GetA":                {initializeSrc, "GetA", nil, initialize.GetA},
	"GetB":                {initializeSrc, "GetB", nil, initialize.GetB},
	"GetC":                {initializeSrc, "GetC", nil, initialize.GetC},
	"GetCacheSum":         {initializeSrc, "GetCacheSum", nil, initialize.GetCacheSum},
	"GetCacheSize":        {initializeSrc, "GetCacheSize", nil, initialize.GetCacheSize},
	"GetFibonacciCount":   {initializeSrc, "GetFibonacciCount", nil, initialize.GetFibonacciCount},
}

var leetcode_hardTests = map[string]testCase{
	"TrappingRainWater":           {leetcodeHardSrc, "TrappingRainWater", nil, leetcode_hard.TrappingRainWater},
	"LargestRectangleInHistogram": {leetcodeHardSrc, "LargestRectangleInHistogram", nil, leetcode_hard.LargestRectangleInHistogram},
	"MedianOfTwoSortedArrays":     {leetcodeHardSrc, "MedianOfTwoSortedArrays", nil, leetcode_hard.MedianOfTwoSortedArrays},
	"RegularExpressionMatching":   {leetcodeHardSrc, "RegularExpressionMatching", nil, leetcode_hard.RegularExpressionMatching},
	"NQueens":                     {leetcodeHardSrc, "NQueens", nil, leetcode_hard.NQueens},
	"LongestIncreasingPath":       {leetcodeHardSrc, "LongestIncreasingPath", nil, leetcode_hard.LongestIncreasingPath},
	"WordLadder":                  {leetcodeHardSrc, "WordLadder", nil, leetcode_hard.WordLadder},
	"MergeKSortedLists":           {leetcodeHardSrc, "MergeKSortedLists", nil, leetcode_hard.MergeKSortedLists},
	"EditDistance":                {leetcodeHardSrc, "EditDistance", nil, leetcode_hard.EditDistance},
	"MinimumWindowSubstring":      {leetcodeHardSrc, "MinimumWindowSubstring", nil, leetcode_hard.MinimumWindowSubstring},
}

var mapadvancedTests = map[string]testCase{
	"LookupExistingKey": {mapadvancedSrc, "LookupExistingKey", nil, mapadvanced.LookupExistingKey},
	"LookupWithDefault": {mapadvancedSrc, "LookupWithDefault", nil, mapadvanced.LookupWithDefault},
	"AsCounter":         {mapadvancedSrc, "AsCounter", nil, mapadvanced.AsCounter},
	"WithStringValues":  {mapadvancedSrc, "WithStringValues", nil, mapadvanced.WithStringValues},
	"BuildFromLoop":     {mapadvancedSrc, "BuildFromLoop", nil, mapadvanced.BuildFromLoop},
	"DeleteAndReinsert": {mapadvancedSrc, "DeleteAndReinsert", nil, mapadvanced.DeleteAndReinsert},
}

var mapsTests = map[string]testCase{
	"BasicOps":       {mapsSrc, "BasicOps", nil, maps.BasicOps},
	"Iteration":      {mapsSrc, "Iteration", nil, maps.Iteration},
	"Delete":         {mapsSrc, "Delete", nil, maps.Delete},
	"Len":            {mapsSrc, "Len", nil, maps.Len},
	"Overwrite":      {mapsSrc, "Overwrite", nil, maps.Overwrite},
	"IntKeys":        {mapsSrc, "IntKeys", nil, maps.IntKeys},
	"PassToFunction": {mapsSrc, "PassToFunction", nil, maps.PassToFunction},
	// Parameterized tests
	"SumValues": {mapsSrc, "SumValues", []any{map[string]int{"a": 100, "b": 200}}, maps.SumValues},
}

var multiassignTests = map[string]testCase{
	"Swap":             {multiassignSrc, "Swap", nil, multiassign.Swap},
	"FromFunction":     {multiassignSrc, "FromFunction", nil, multiassign.FromFunction},
	"ThreeValues":      {multiassignSrc, "ThreeValues", nil, multiassign.ThreeValues},
	"InLoop":           {multiassignSrc, "InLoop", nil, multiassign.InLoop},
	"DiscardWithBlank": {multiassignSrc, "DiscardWithBlank", nil, multiassign.DiscardWithBlank},
	// Parameterized tests
	"TwoVals":    {multiassignSrc, "TwoVals", nil, multiassign.TwoVals},
	"ThreeValsN": {multiassignSrc, "ThreeValsN", []any{10}, multiassign.ThreeValsN},
	"DivmodAB":   {multiassignSrc, "DivmodAB", []any{17, 5}, multiassign.DivmodAB},
}

var namedreturnTests = map[string]testCase{
	"Basic":     {namedreturnSrc, "Basic", nil, namedreturn.Basic},
	"Multiple":  {namedreturnSrc, "Multiple", nil, namedreturn.Multiple},
	"ZeroValue": {namedreturnSrc, "ZeroValue", nil, namedreturn.ZeroValue},
	"Divmod":    {namedreturnSrc, "Divmod", []any{1000, 7}, namedreturn.Divmod},
}

var recursionTests = map[string]testCase{
	"TailRecursionPattern": {recursionSrc, "TailRecursionPattern", nil, recursion.TailRecursionPattern},
	"ReverseSlice":         {recursionSrc, "ReverseSlice", nil, recursion.ReverseSlice},
	"TowerOfHanoi":         {recursionSrc, "TowerOfHanoi", nil, recursion.TowerOfHanoi},
	"MaxSlice":             {recursionSrc, "MaxSlice", nil, recursion.MaxSlice},
	"Ackermann":            {recursionSrc, "Ackermann", nil, recursion.Ackermann},
	"BinarySearch":         {recursionSrc, "BinarySearch", nil, recursion.BinarySearch},
	// Parameterized tests
	"SumTail": {recursionSrc, "SumTail", []any{50, 0}, recursion.SumTail},
	"HanoiN":  {recursionSrc, "HanoiN", []any{10}, recursion.HanoiN},
	"Ack":     {recursionSrc, "Ack", []any{2, 3}, recursion.Ack},
	"MaxVal":  {recursionSrc, "MaxVal", []any{[]int{3, 7, 1, 9, 4}, 5}, recursion.MaxVal},
}

var resolved_issueTests = map[string]testCase{
	"BytesToString":                      {resolvedIssueSrc, "BytesToString", nil, resolved_issue.BytesToString},
	"BytesToStringHi":                    {resolvedIssueSrc, "BytesToStringHi", nil, resolved_issue.BytesToStringHi},
	"BytesToStringGo":                    {resolvedIssueSrc, "BytesToStringGo", nil, resolved_issue.BytesToStringGo},
	"BytesToStringEmpty":                 {resolvedIssueSrc, "BytesToStringEmpty", nil, resolved_issue.BytesToStringEmpty},
	"BytesToStringSingle":                {resolvedIssueSrc, "BytesToStringSingle", nil, resolved_issue.BytesToStringSingle},
	"PointerReceiverMutation":            {resolvedIssueSrc, "PointerReceiverMutation", nil, resolved_issue.PointerReceiverMutation},
	"PointerReceiverMutationReturnValue": {resolvedIssueSrc, "PointerReceiverMutationReturnValue", nil, resolved_issue.PointerReceiverMutationReturnValue},
	"InitFuncExecuted":                   {resolvedIssueSrc, "InitFuncExecuted", nil, resolved_issue.InitFuncExecuted},
	"InitFuncSideEffect":                 {resolvedIssueSrc, "InitFuncSideEffect", nil, resolved_issue.InitFuncSideEffect},
	"RangeStringRuneValue":               {resolvedIssueSrc, "RangeStringRuneValue", nil, resolved_issue.RangeStringRuneValue},
	"RangeStringIndexValue":              {resolvedIssueSrc, "RangeStringIndexValue", nil, resolved_issue.RangeStringIndexValue},
	"RangeStringMultibyte":               {resolvedIssueSrc, "RangeStringMultibyte", nil, resolved_issue.RangeStringMultibyte},
	"MapWithFuncValue":                   {resolvedIssueSrc, "MapWithFuncValue", nil, resolved_issue.MapWithFuncValue},
	"InterfaceSliceTypeSwitch":           {resolvedIssueSrc, "InterfaceSliceTypeSwitch", nil, resolved_issue.InterfaceSliceTypeSwitch},
	"StructWithFuncField":                {resolvedIssueSrc, "StructWithFuncField", nil, resolved_issue.StructWithFuncField},
	"SliceFlatten":                       {resolvedIssueSrc, "SliceFlatten", nil, resolved_issue.SliceFlatten},
	// Note: resolved_issue/MapUpdateDuringRange removed - non-deterministic map iteration order
	"StructSelfRef":         {resolvedIssueSrc, "StructSelfRef", nil, resolved_issue.StructSelfRef},
	"DeferInClosureWithArg": {resolvedIssueSrc, "DeferInClosureWithArg", nil, resolved_issue.DeferInClosureWithArg},
	"PointerSwapInStruct":   {resolvedIssueSrc, "PointerSwapInStruct", nil, resolved_issue.PointerSwapInStruct},
	"StructWithFuncSlice":   {resolvedIssueSrc, "StructWithFuncSlice", nil, resolved_issue.StructWithFuncSlice},
	"StructAnonymousField":  {resolvedIssueSrc, "StructAnonymousField", nil, resolved_issue.StructAnonymousField},
	// Note: resolved_issue/MapRangeWithBreak removed - non-deterministic map iteration order
	"PointerToInterface":           {resolvedIssueSrc, "PointerToInterface", nil, resolved_issue.PointerToInterface},
	"PointerToSliceElemModify":     {resolvedIssueSrc, "PointerToSliceElemModify", nil, resolved_issue.PointerToSliceElemModify},
	"StructWithFuncPtrTest":        {resolvedIssueSrc, "StructWithFuncPtrTest", nil, resolved_issue.StructWithFuncPtrTest},
	"PointerCompareDiffTest":       {resolvedIssueSrc, "PointerCompareDiffTest", nil, resolved_issue.PointerCompareDiffTest},
	"DeferModifyMultipleNamedTest": {resolvedIssueSrc, "DeferModifyMultipleNamedTest", nil, resolved_issue.DeferModifyMultipleNamedTest},
	"DeferNamedReturnNilTest":      {resolvedIssueSrc, "DeferNamedReturnNilTest", nil, resolved_issue.DeferNamedReturnNilTest},
	"DeferNamedReturnNilPtrTest":   {resolvedIssueSrc, "DeferNamedReturnNilPtrTest", nil, resolved_issue.DeferNamedReturnNilPtrTest},
	"DeferNamedReturnMultiTest":    {resolvedIssueSrc, "DeferNamedReturnMultiTest", nil, resolved_issue.DeferNamedReturnMultiTest},
}

var scopeTests = map[string]testCase{
	"IfInitShortVar":            {scopeSrc, "IfInitShortVar", nil, scope.IfInitShortVar},
	"IfInitMultiCondition":      {scopeSrc, "IfInitMultiCondition", nil, scope.IfInitMultiCondition},
	"NestedScopes":              {scopeSrc, "NestedScopes", nil, scope.NestedScopes},
	"ForScopeIsolation":         {scopeSrc, "ForScopeIsolation", nil, scope.ForScopeIsolation},
	"MultipleBlockScopes":       {scopeSrc, "MultipleBlockScopes", nil, scope.MultipleBlockScopes},
	"ClosureCapturesOuterScope": {scopeSrc, "ClosureCapturesOuterScope", nil, scope.ClosureCapturesOuterScope},
	// Parameterized tests
	"Abs": {scopeSrc, "Abs", []any{-42}, scope.Abs},
}

var slicesTests = map[string]testCase{
	"MakeLen":           {slicesSrc, "MakeLen", nil, slices.MakeLen},
	"Append":            {slicesSrc, "Append", nil, slices.Append},
	"ElementAssignment": {slicesSrc, "ElementAssignment", nil, slices.ElementAssignment},
	"ForRange":          {slicesSrc, "ForRange", nil, slices.ForRange},
	"ForRangeIndex":     {slicesSrc, "ForRangeIndex", nil, slices.ForRangeIndex},
	"GrowMultiple":      {slicesSrc, "GrowMultiple", nil, slices.GrowMultiple},
	"PassToFunction":    {slicesSrc, "PassToFunction", nil, slices.PassToFunction},
	"LenCap":            {slicesSrc, "LenCap", nil, slices.LenCap},
	// Parameterized tests
	"SumSlice": {slicesSrc, "SumSlice", []any{[]int{1, 2, 3, 4, 5}}, slices.SumSlice},
}

var slicingTests = map[string]testCase{
	"SubSliceBasic":            {slicingSrc, "SubSliceBasic", nil, slicing.SubSliceBasic},
	"SubSliceLen":              {slicingSrc, "SubSliceLen", nil, slicing.SubSliceLen},
	"SubSliceFromStart":        {slicingSrc, "SubSliceFromStart", nil, slicing.SubSliceFromStart},
	"SubSliceToEnd":            {slicingSrc, "SubSliceToEnd", nil, slicing.SubSliceToEnd},
	"SubSliceCopy":             {slicingSrc, "SubSliceCopy", nil, slicing.SubSliceCopy},
	"SubSliceChained":          {slicingSrc, "SubSliceChained", nil, slicing.SubSliceChained},
	"SubSliceModifiesOriginal": {slicingSrc, "SubSliceModifiesOriginal", nil, slicing.SubSliceModifiesOriginal},
	// Parameterized tests
	"SliceLen":          {slicingSrc, "SliceLen", []any{[]int{10, 20, 30, 40, 50, 60, 70}, 2, 5}, slicing.SliceLen},
	"SliceSumRange":     {slicingSrc, "SliceSumRange", []any{[]int{10, 20, 30, 40, 50}, 1, 4}, slicing.SliceSumRange},
	"SliceFirstElement": {slicingSrc, "SliceFirstElement", []any{[]int{100, 200, 300}, 0}, slicing.SliceFirstElement},
}

var strings_pkgTests = map[string]testCase{
	"Concat":     {stringsPkgSrc, "Concat", nil, strings_pkg.Concat},
	"ConcatLoop": {stringsPkgSrc, "ConcatLoop", nil, strings_pkg.ConcatLoop},
	"Len":        {stringsPkgSrc, "Len", nil, strings_pkg.Len},
	"Index":      {stringsPkgSrc, "Index", nil, strings_pkg.Index},
	"Comparison": {stringsPkgSrc, "Comparison", nil, strings_pkg.Comparison},
	"Equality":   {stringsPkgSrc, "Equality", nil, strings_pkg.Equality},
	"EmptyCheck": {stringsPkgSrc, "EmptyCheck", nil, strings_pkg.EmptyCheck},
	// Parameterized tests
	"StrConcat":  {stringsPkgSrc, "StrConcat", []any{"hello", " world"}, strings_pkg.StrConcat},
	"StrLen":     {stringsPkgSrc, "StrLen", []any{"hello"}, strings_pkg.StrLen},
	"StrCompare": {stringsPkgSrc, "StrCompare", []any{"abc", "abd"}, strings_pkg.StrCompare},
	"StrEqual":   {stringsPkgSrc, "StrEqual", []any{"hello", "hello"}, strings_pkg.StrEqual},
}

var structsTests = map[string]testCase{
	"BasicStruct":                 {structsSrc, "BasicStruct", nil, structs.BasicStruct},
	"StructPointer":               {structsSrc, "StructPointer", nil, structs.StructPointer},
	"NestedStruct":                {structsSrc, "NestedStruct", nil, structs.NestedStruct},
	"EmbeddedField":               {structsSrc, "EmbeddedField", nil, structs.EmbeddedField},
	"StructInSlice":               {structsSrc, "StructInSlice", nil, structs.StructInSlice},
	"StructAsParam":               {structsSrc, "StructAsParam", nil, structs.StructAsParam},
	"StructZeroValue":             {structsSrc, "StructZeroValue", nil, structs.StructZeroValue},
	"MultipleEmbedded":            {structsSrc, "MultipleEmbedded", nil, structs.MultipleEmbedded},
	"DeepNesting":                 {structsSrc, "DeepNesting", nil, structs.DeepNesting},
	"StructFieldMutation":         {structsSrc, "StructFieldMutation", nil, structs.StructFieldMutation},
	"StructWithBool":              {structsSrc, "StructWithBool", nil, structs.StructWithBool},
	"StructCopySemantics":         {structsSrc, "StructCopySemantics", nil, structs.StructCopySemantics},
	"StructPointerSharing":        {structsSrc, "StructPointerSharing", nil, structs.StructPointerSharing},
	"StructReturnFromFunc":        {structsSrc, "StructReturnFromFunc", nil, structs.StructReturnFromFunc},
	"StructPointerReturnFromFunc": {structsSrc, "StructPointerReturnFromFunc", nil, structs.StructPointerReturnFromFunc},
	"StructSliceAppend":           {structsSrc, "StructSliceAppend", nil, structs.StructSliceAppend},
	"StructPointerSlice":          {structsSrc, "StructPointerSlice", nil, structs.StructPointerSlice},
	"StructInMap":                 {structsSrc, "StructInMap", nil, structs.StructInMap},
	"StructConditionalInit":       {structsSrc, "StructConditionalInit", nil, structs.StructConditionalInit},
	"StructFieldLoop":             {structsSrc, "StructFieldLoop", nil, structs.StructFieldLoop},
	"StructNestedMutation":        {structsSrc, "StructNestedMutation", nil, structs.StructNestedMutation},
	"StructEmbeddedOverride":      {structsSrc, "StructEmbeddedOverride", nil, structs.StructEmbeddedOverride},
	"StructWithClosure":           {structsSrc, "StructWithClosure", nil, structs.StructWithClosure},
	"StructReassign":              {structsSrc, "StructReassign", nil, structs.StructReassign},
	"StructSliceOfNested":         {structsSrc, "StructSliceOfNested", nil, structs.StructSliceOfNested},
	"StructMultiReturn":           {structsSrc, "StructMultiReturn", nil, structs.StructMultiReturn},
	"StructBuilderPattern":        {structsSrc, "StructBuilderPattern", nil, structs.StructBuilderPattern},
	"StructArrayField":            {structsSrc, "StructArrayField", nil, structs.StructArrayField},
	"StructEmbeddedChain":         {structsSrc, "StructEmbeddedChain", nil, structs.StructEmbeddedChain},
}

var switchTests = map[string]testCase{
	"Simple":      {switchSrc, "Simple", nil, switch_pkg.Simple},
	"Default":     {switchSrc, "Default", nil, switch_pkg.Default},
	"MultiCase":   {switchSrc, "MultiCase", nil, switch_pkg.MultiCase},
	"NoCondition": {switchSrc, "NoCondition", nil, switch_pkg.NoCondition},
	"WithInit":    {switchSrc, "WithInit", nil, switch_pkg.WithInit},
	"StringCases": {switchSrc, "StringCases", nil, switch_pkg.StringCases},
	"Fallthrough": {switchSrc, "Fallthrough", nil, switch_pkg.Fallthrough},
	"Nested":      {switchSrc, "Nested", nil, switch_pkg.Nested},
	// Parameterized tests
	"Classify":  {switchSrc, "Classify", []any{2}, switch_pkg.Classify},
	"Weekday":   {switchSrc, "Weekday", []any{3}, switch_pkg.Weekday},
	"Grade":     {switchSrc, "Grade", []any{85}, switch_pkg.Grade},
	"ColorCode": {switchSrc, "ColorCode", []any{"green"}, switch_pkg.ColorCode},
}

var trickyClosuresTests = map[string]testCase{
	"ClosureArgDefaultTest":        {trickySrc, "ClosureArgDefaultTest", nil, tricky.ClosureArgDefaultTest},
	"ClosureAsArg":                 {trickySrc, "ClosureAsArg", nil, tricky.ClosureAsArg},
	"ClosureCaptureAndModifyTest":  {trickySrc, "ClosureCaptureAndModifyTest", nil, tricky.ClosureCaptureAndModifyTest},
	"ClosureCaptureLoopVarTest":    {trickySrc, "ClosureCaptureLoopVarTest", nil, tricky.ClosureCaptureLoopVarTest},
	"ClosureCapturesTwoTest":       {trickySrc, "ClosureCapturesTwoTest", nil, tricky.ClosureCapturesTwoTest},
	"ClosureCapturingPointer":      {trickySrc, "ClosureCapturingPointer", nil, tricky.ClosureCapturingPointer},
	"ClosureCompose":               {trickySrc, "ClosureCompose", nil, tricky.ClosureCompose},
	"ClosureComposeTest":           {trickySrc, "ClosureComposeTest", nil, tricky.ClosureComposeTest},
	"ClosureConstTest":             {trickySrc, "ClosureConstTest", nil, tricky.ClosureConstTest},
	"ClosureCounter":               {trickySrc, "ClosureCounter", nil, tricky.ClosureCounter},
	"ClosureCounterResetTest":      {trickySrc, "ClosureCounterResetTest", nil, tricky.ClosureCounterResetTest},
	"ClosureCounterStateTest":      {trickySrc, "ClosureCounterStateTest", nil, tricky.ClosureCounterStateTest},
	"ClosureCounterTest":           {trickySrc, "ClosureCounterTest", nil, tricky.ClosureCounterTest},
	"ClosureCurry":                 {trickySrc, "ClosureCurry", nil, tricky.ClosureCurry},
	"ClosureCurryMultipleArgTest":  {trickySrc, "ClosureCurryMultipleArgTest", nil, tricky.ClosureCurryMultipleArgTest},
	"ClosureCurryMultipleTest":     {trickySrc, "ClosureCurryMultipleTest", nil, tricky.ClosureCurryMultipleTest},
	"ClosureCurryTest":             {trickySrc, "ClosureCurryTest", nil, tricky.ClosureCurryTest},
	"ClosureEnvCaptureTest":        {trickySrc, "ClosureEnvCaptureTest", nil, tricky.ClosureEnvCaptureTest},
	"ClosureFibonacci":             {trickySrc, "ClosureFibonacci", nil, tricky.ClosureFibonacci},
	"ClosureFlip":                  {trickySrc, "ClosureFlip", nil, tricky.ClosureFlip},
	"ClosureFlipTest":              {trickySrc, "ClosureFlipTest", nil, tricky.ClosureFlipTest},
	"ClosureMap":                   {trickySrc, "ClosureMap", nil, tricky.ClosureMap},
	"ClosureMapBuilderTest":        {trickySrc, "ClosureMapBuilderTest", nil, tricky.ClosureMapBuilderTest},
	"ClosureMemoize":               {trickySrc, "ClosureMemoize", nil, tricky.ClosureMemoize},
	"ClosureMemoizeRecursive":      {trickySrc, "ClosureMemoizeRecursive", nil, tricky.ClosureMemoizeRecursive},
	"ClosureMemoizeRecursiveTest":  {trickySrc, "ClosureMemoizeRecursiveTest", nil, tricky.ClosureMemoizeRecursiveTest},
	"ClosureModifyOuterVarTest":    {trickySrc, "ClosureModifyOuterVarTest", nil, tricky.ClosureModifyOuterVarTest},
	"ClosureMultiCapture":          {trickySrc, "ClosureMultiCapture", nil, tricky.ClosureMultiCapture},
	"ClosureMultipleCallsTest":     {trickySrc, "ClosureMultipleCallsTest", nil, tricky.ClosureMultipleCallsTest},
	"ClosureMultiReturnTest":       {trickySrc, "ClosureMultiReturnTest", nil, tricky.ClosureMultiReturnTest},
	"ClosureMutateCapturedSlice":   {trickySrc, "ClosureMutateCapturedSlice", nil, tricky.ClosureMutateCapturedSlice},
	"ClosureMutateClosureTest":     {trickySrc, "ClosureMutateClosureTest", nil, tricky.ClosureMutateClosureTest},
	"ClosureMutatesOuterTest":      {trickySrc, "ClosureMutatesOuterTest", nil, tricky.ClosureMutatesOuterTest},
	"ClosureOnce":                  {trickySrc, "ClosureOnce", nil, tricky.ClosureOnce},
	"ClosureOnceTest":              {trickySrc, "ClosureOnceTest", nil, tricky.ClosureOnceTest},
	"ClosurePartial":               {trickySrc, "ClosurePartial", nil, tricky.ClosurePartial},
	"ClosurePartialApplyTest":      {trickySrc, "ClosurePartialApplyTest", nil, tricky.ClosurePartialApplyTest},
	"ClosurePartialTest":           {trickySrc, "ClosurePartialTest", nil, tricky.ClosurePartialTest},
	"ClosurePipeline":              {trickySrc, "ClosurePipeline", nil, tricky.ClosurePipeline},
	"ClosurePtrCaptureTest":        {trickySrc, "ClosurePtrCaptureTest", nil, tricky.ClosurePtrCaptureTest},
	"ClosureRecursiveMemoTest":     {trickySrc, "ClosureRecursiveMemoTest", nil, tricky.ClosureRecursiveMemoTest},
	"ClosureRecursiveSimpleTest":   {trickySrc, "ClosureRecursiveSimpleTest", nil, tricky.ClosureRecursiveSimpleTest},
	"ClosureReturningClosure":      {trickySrc, "ClosureReturningClosure", nil, tricky.ClosureReturningClosure},
	"ClosureReturnMultipleTest":    {trickySrc, "ClosureReturnMultipleTest", nil, tricky.ClosureReturnMultipleTest},
	"ClosureReturnsClosureTest":    {trickySrc, "ClosureReturnsClosureTest", nil, tricky.ClosureReturnsClosureTest},
	"ClosureReturnsValueTest":      {trickySrc, "ClosureReturnsValueTest", nil, tricky.ClosureReturnsValueTest},
	"ClosureReturnValueTest":       {trickySrc, "ClosureReturnValueTest", nil, tricky.ClosureReturnValueTest},
	"ClosureSliceAccumTest":        {trickySrc, "ClosureSliceAccumTest", nil, tricky.ClosureSliceAccumTest},
	"ClosureSliceBuilderTest":      {trickySrc, "ClosureSliceBuilderTest", nil, tricky.ClosureSliceBuilderTest},
	"ClosureSliceCaptureTest":      {trickySrc, "ClosureSliceCaptureTest", nil, tricky.ClosureSliceCaptureTest},
	"ClosureTap":                   {trickySrc, "ClosureTap", nil, tricky.ClosureTap},
	"ClosureTapTest":               {trickySrc, "ClosureTapTest", nil, tricky.ClosureTapTest},
	"ClosureVarCaptureTest":        {trickySrc, "ClosureVarCaptureTest", nil, tricky.ClosureVarCaptureTest},
	"ClosureWithDeferAndReturn":    {trickySrc, "ClosureWithDeferAndReturn", nil, tricky.ClosureWithDeferAndReturn},
	"ClosureWithDeferTest":         {trickySrc, "ClosureWithDeferTest", nil, tricky.ClosureWithDeferTest},
	"ClosureWithExternalVar":       {trickySrc, "ClosureWithExternalVar", nil, tricky.ClosureWithExternalVar},
	"ClosureWithLocalVarTest":      {trickySrc, "ClosureWithLocalVarTest", nil, tricky.ClosureWithLocalVarTest},
	"ClosureWithLoopVar":           {trickySrc, "ClosureWithLoopVar", nil, tricky.ClosureWithLoopVar},
	"ClosureWithMapCaptureTest":    {trickySrc, "ClosureWithMapCaptureTest", nil, tricky.ClosureWithMapCaptureTest},
	"ClosureWithMultipleReturns":   {trickySrc, "ClosureWithMultipleReturns", nil, tricky.ClosureWithMultipleReturns},
	"ClosureWithRecursion":         {trickySrc, "ClosureWithRecursion", nil, tricky.ClosureWithRecursion},
	"ClosureWithStructCaptureTest": {trickySrc, "ClosureWithStructCaptureTest", nil, tricky.ClosureWithStructCaptureTest},
	"ClosureWithVarCaptureTest":    {trickySrc, "ClosureWithVarCaptureTest", nil, tricky.ClosureWithVarCaptureTest},
}

var trickyDeferTests = map[string]testCase{
	"DeferAfterReturnTest":                 {trickySrc, "DeferAfterReturnTest", nil, tricky.DeferAfterReturnTest},
	"DeferCaptureMapTest":                  {trickySrc, "DeferCaptureMapTest", nil, tricky.DeferCaptureMapTest},
	"DeferCaptureSliceTest":                {trickySrc, "DeferCaptureSliceTest", nil, tricky.DeferCaptureSliceTest},
	"DeferCaptureValueTest":                {trickySrc, "DeferCaptureValueTest", nil, tricky.DeferCaptureValueTest},
	"DeferClosureArgTest":                  {trickySrc, "DeferClosureArgTest", nil, tricky.DeferClosureArgTest},
	"DeferClosureCaptureModifyTest":        {trickySrc, "DeferClosureCaptureModifyTest", nil, tricky.DeferClosureCaptureModifyTest},
	"DeferClosureModifyingNamed":           {trickySrc, "DeferClosureModifyingNamed", nil, tricky.DeferClosureModifyingNamed},
	"DeferClosureModifyTest":               {trickySrc, "DeferClosureModifyTest", nil, tricky.DeferClosureModifyTest},
	"DeferClosureNestedTest":               {trickySrc, "DeferClosureNestedTest", nil, tricky.DeferClosureNestedTest},
	"DeferConditional":                     {trickySrc, "DeferConditional", nil, tricky.DeferConditional},
	"DeferConditionalModifyTest":           {trickySrc, "DeferConditionalModifyTest", nil, tricky.DeferConditionalModifyTest},
	"DeferInClosureTest":                   {trickySrc, "DeferInClosureTest", nil, tricky.DeferInClosureTest},
	"DeferInClosureWithArg":                {trickySrc, "DeferInClosureWithArg", nil, tricky.DeferInClosureWithArg},
	"DeferInGoroutine":                     {trickySrc, "DeferInGoroutine", nil, tricky.DeferInGoroutine},
	"DeferInMultipleFunctions":             {trickySrc, "DeferInMultipleFunctions", nil, tricky.DeferInMultipleFunctions},
	"DeferInNestedFunction":                {trickySrc, "DeferInNestedFunction", nil, tricky.DeferInNestedFunction},
	"DeferMapModifyTest":                   {trickySrc, "DeferMapModifyTest", nil, tricky.DeferMapModifyTest},
	"DeferModifiesReturnTest":              {trickySrc, "DeferModifiesReturnTest", nil, tricky.DeferModifiesReturnTest},
	"DeferModifyCapture":                   {trickySrc, "DeferModifyCapture", nil, tricky.DeferModifyCapture},
	"DeferModifyMap":                       {trickySrc, "DeferModifyMap", nil, tricky.DeferModifyMap},
	"DeferModifyMapNamedTest":              {trickySrc, "DeferModifyMapNamedTest", nil, tricky.DeferModifyMapNamedTest},
	"DeferModifyMapTest":                   {trickySrc, "DeferModifyMapTest", nil, tricky.DeferModifyMapTest},
	"DeferModifyMultiple":                  {trickySrc, "DeferModifyMultiple", nil, tricky.DeferModifyMultiple},
	"DeferModifyMultipleCombined":          {trickySrc, "DeferModifyMultipleCombined", nil, tricky.DeferModifyMultipleCombined},
	"DeferModifyMultipleNamedTest":         {trickySrc, "DeferModifyMultipleNamedTest", nil, tricky.DeferModifyMultipleNamedTest},
	"DeferModifyNamedReturnTest":           {trickySrc, "DeferModifyNamedReturnTest", nil, tricky.DeferModifyNamedReturnTest},
	"DeferModifyPtrTest":                   {trickySrc, "DeferModifyPtrTest", nil, tricky.DeferModifyPtrTest},
	"DeferModifyReturnValue":               {trickySrc, "DeferModifyReturnValue", nil, tricky.DeferModifyReturnValue},
	"DeferModifyReturnValueTest":           {trickySrc, "DeferModifyReturnValueTest", nil, tricky.DeferModifyReturnValueTest},
	"DeferModifySlice":                     {trickySrc, "DeferModifySlice", nil, tricky.DeferModifySlice},
	"DeferModifySliceTest":                 {trickySrc, "DeferModifySliceTest", nil, tricky.DeferModifySliceTest},
	"DeferMultiNamedReturnTest":            {trickySrc, "DeferMultiNamedReturnTest", nil, tricky.DeferMultiNamedReturnTest},
	"DeferMultipleCalls":                   {trickySrc, "DeferMultipleCalls", nil, tricky.DeferMultipleCalls},
	"DeferMultipleExecTest":                {trickySrc, "DeferMultipleExecTest", nil, tricky.DeferMultipleExecTest},
	"DeferMultipleFuncTest":                {trickySrc, "DeferMultipleFuncTest", nil, tricky.DeferMultipleFuncTest},
	"DeferMultipleNamedTest":               {trickySrc, "DeferMultipleNamedTest", nil, tricky.DeferMultipleNamedTest},
	"DeferMultipleVars":                    {trickySrc, "DeferMultipleVars", nil, tricky.DeferMultipleVars},
	"DeferNamedMultiTest":                  {trickySrc, "DeferNamedMultiTest", nil, tricky.DeferNamedMultiTest},
	"DeferNamedResultChainTest":            {trickySrc, "DeferNamedResultChainTest", nil, tricky.DeferNamedResultChainTest},
	"DeferNamedResultNilTest":              {trickySrc, "DeferNamedResultNilTest", nil, tricky.DeferNamedResultNilTest},
	"DeferNamedResultTest":                 {trickySrc, "DeferNamedResultTest", nil, tricky.DeferNamedResultTest},
	"DeferNamedReturnCaptureTest":          {trickySrc, "DeferNamedReturnCaptureTest", nil, tricky.DeferNamedReturnCaptureTest},
	"DeferNamedReturnCombineTest":          {trickySrc, "DeferNamedReturnCombineTest", nil, tricky.DeferNamedReturnCombineTest},
	"DeferNamedReturnDoubleTest":           {trickySrc, "DeferNamedReturnDoubleTest", nil, tricky.DeferNamedReturnDoubleTest},
	"DeferNamedReturnModifyTest":           {trickySrc, "DeferNamedReturnModifyTest", nil, tricky.DeferNamedReturnModifyTest},
	"DeferNamedReturnMultiTest":            {trickySrc, "DeferNamedReturnMultiTest", nil, tricky.DeferNamedReturnMultiTest},
	"DeferNamedReturnNilPtrTest":           {trickySrc, "DeferNamedReturnNilPtrTest", nil, tricky.DeferNamedReturnNilPtrTest},
	"DeferNamedReturnNilTest":              {trickySrc, "DeferNamedReturnNilTest", nil, tricky.DeferNamedReturnNilTest},
	"DeferNamedReturnOrderTest":            {trickySrc, "DeferNamedReturnOrderTest", nil, tricky.DeferNamedReturnOrderTest},
	"DeferPanicRecoverValueTest":           {trickySrc, "DeferPanicRecoverValueTest", nil, tricky.DeferPanicRecoverValueTest},
	"DeferReadCapture":                     {trickySrc, "DeferReadCapture", nil, tricky.DeferReadCapture},
	"DeferRecoverPanicTest":                {trickySrc, "DeferRecoverPanicTest", nil, tricky.DeferRecoverPanicTest},
	"DeferReturnValue":                     {trickySrc, "DeferReturnValue", nil, tricky.DeferReturnValue},
	"DeferReturnValueModifyTest":           {trickySrc, "DeferReturnValueModifyTest", nil, tricky.DeferReturnValueModifyTest},
	"DeferStackTest":                       {trickySrc, "DeferStackTest", nil, tricky.DeferStackTest},
	"DeferWithCapture":                     {trickySrc, "DeferWithCapture", nil, tricky.DeferWithCapture},
	"DeferWithClosureArg":                  {trickySrc, "DeferWithClosureArg", nil, tricky.DeferWithClosureArg},
	"DeferWithClosureResult":               {trickySrc, "DeferWithClosureResult", nil, tricky.DeferWithClosureResult},
	"DeferWithLoop":                        {trickySrc, "DeferWithLoop", nil, tricky.DeferWithLoop},
	"DeferWithMultipleReturns":             {trickySrc, "DeferWithMultipleReturns", nil, tricky.DeferWithMultipleReturns},
	"DeferWithMultipleReturnsCombined":     {trickySrc, "DeferWithMultipleReturnsCombined", nil, tricky.DeferWithMultipleReturnsCombined},
	"DeferWithNamedResultMultiple":         {trickySrc, "DeferWithNamedResultMultiple", nil, tricky.DeferWithNamedResultMultiple},
	"DeferWithNamedResultMultipleCombined": {trickySrc, "DeferWithNamedResultMultipleCombined", nil, tricky.DeferWithNamedResultMultipleCombined},
	"DeferWithNamedReturn":                 {trickySrc, "DeferWithNamedReturn", nil, tricky.DeferWithNamedReturn},
	"DeferWithRecoveredPanic":              {trickySrc, "DeferWithRecoveredPanic", nil, tricky.DeferWithRecoveredPanic},
	"DeferWithRetFunc":                     {trickySrc, "DeferWithRetFunc", nil, tricky.DeferWithRetFunc},
	"DeferWithReturnFunc":                  {trickySrc, "DeferWithReturnFunc", nil, tricky.DeferWithReturnFunc},
}

var trickyInterfacesTests = map[string]testCase{
	"InterfaceMethod":           {trickySrc, "InterfaceMethod", nil, tricky.InterfaceMethod},
	"InterfaceNilTypeAssertion": {trickySrc, "InterfaceNilTypeAssertion", nil, tricky.InterfaceNilTypeAssertion},
	"InterfaceSliceTypeAssert":  {trickySrc, "InterfaceSliceTypeAssert", nil, tricky.InterfaceSliceTypeAssert},
}

var trickyMapsTests = map[string]testCase{
	"MapAll":                    {trickySrc, "MapAll", nil, tricky.MapAll},
	"MapAllMatch":               {trickySrc, "MapAllMatch", nil, tricky.MapAllMatch},
	"MapAllTest":                {trickySrc, "MapAllTest", nil, tricky.MapAllTest},
	"MapAny":                    {trickySrc, "MapAny", nil, tricky.MapAny},
	"MapAnyMatch":               {trickySrc, "MapAnyMatch", nil, tricky.MapAnyMatch},
	"MapAnyTest":                {trickySrc, "MapAnyTest", nil, tricky.MapAnyTest},
	"MapAnyValueTest":           {trickySrc, "MapAnyValueTest", nil, tricky.MapAnyValueTest},
	"MapApplyToValuesTest":      {trickySrc, "MapApplyToValuesTest", nil, tricky.MapApplyToValuesTest},
	"MapClearMakeTest":          {trickySrc, "MapClearMakeTest", nil, tricky.MapClearMakeTest},
	"MapClearRange":             {trickySrc, "MapClearRange", nil, tricky.MapClearRange},
	"MapCombine":                {trickySrc, "MapCombine", nil, tricky.MapCombine},
	"MapCombineSameKeyTest":     {trickySrc, "MapCombineSameKeyTest", nil, tricky.MapCombineSameKeyTest},
	"MapCombineTest":            {trickySrc, "MapCombineTest", nil, tricky.MapCombineTest},
	"MapCompact":                {trickySrc, "MapCompact", nil, tricky.MapCompact},
	"MapContainsVal":            {trickySrc, "MapContainsVal", nil, tricky.MapContainsVal},
	"MapCopy":                   {trickySrc, "MapCopy", nil, tricky.MapCopy},
	"MapCountByKey":             {trickySrc, "MapCountByKey", nil, tricky.MapCountByKey},
	"MapCountByValueTest":       {trickySrc, "MapCountByValueTest", nil, tricky.MapCountByValueTest},
	"MapCountIfTest":            {trickySrc, "MapCountIfTest", nil, tricky.MapCountIfTest},
	"MapCountPredTest":          {trickySrc, "MapCountPredTest", nil, tricky.MapCountPredTest},
	"MapCountValues":            {trickySrc, "MapCountValues", nil, tricky.MapCountValues},
	"MapDedup":                  {trickySrc, "MapDedup", nil, tricky.MapDedup},
	"MapDeepGet":                {trickySrc, "MapDeepGet", nil, tricky.MapDeepGet},
	"MapDeepMerge":              {trickySrc, "MapDeepMerge", nil, tricky.MapDeepMerge},
	"MapDeepSet":                {trickySrc, "MapDeepSet", nil, tricky.MapDeepSet},
	"MapDefaultPattern":         {trickySrc, "MapDefaultPattern", nil, tricky.MapDefaultPattern},
	"MapDiff":                   {trickySrc, "MapDiff", nil, tricky.MapDiff},
	"MapDiffKeysTest":           {trickySrc, "MapDiffKeysTest", nil, tricky.MapDiffKeysTest},
	"MapDiffTest":               {trickySrc, "MapDiffTest", nil, tricky.MapDiffTest},
	"MapDropKeys":               {trickySrc, "MapDropKeys", nil, tricky.MapDropKeys},
	"MapDropTest":               {trickySrc, "MapDropTest", nil, tricky.MapDropTest},
	"MapEmptyCheck":             {trickySrc, "MapEmptyCheck", nil, tricky.MapEmptyCheck},
	"MapEmptyKey":               {trickySrc, "MapEmptyKey", nil, tricky.MapEmptyKey},
	"MapEvery":                  {trickySrc, "MapEvery", nil, tricky.MapEvery},
	"MapFilter":                 {trickySrc, "MapFilter", nil, tricky.MapFilter},
	"MapFilterByKeyTest":        {trickySrc, "MapFilterByKeyTest", nil, tricky.MapFilterByKeyTest},
	"MapFilterByValueTest":      {trickySrc, "MapFilterByValueTest", nil, tricky.MapFilterByValueTest},
	"MapFilterKeys":             {trickySrc, "MapFilterKeys", nil, tricky.MapFilterKeys},
	"MapFilterKeysTest":         {trickySrc, "MapFilterKeysTest", nil, tricky.MapFilterKeysTest},
	"MapFind":                   {trickySrc, "MapFind", nil, tricky.MapFind},
	"MapFindKeyTest":            {trickySrc, "MapFindKeyTest", nil, tricky.MapFindKeyTest},
	"MapFindValueTest":          {trickySrc, "MapFindValueTest", nil, tricky.MapFindValueTest},
	"MapFirstKey":               {trickySrc, "MapFirstKey", nil, tricky.MapFirstKey},
	"MapFlatten":                {trickySrc, "MapFlatten", nil, tricky.MapFlatten},
	"MapFlattenTest":            {trickySrc, "MapFlattenTest", nil, tricky.MapFlattenTest},
	"MapFlip":                   {trickySrc, "MapFlip", nil, tricky.MapFlip},
	"MapFloatKey":               {trickySrc, "MapFloatKey", nil, tricky.MapFloatKey},
	"MapForEach":                {trickySrc, "MapForEach", nil, tricky.MapForEach},
	"MapGetOrCreate":            {trickySrc, "MapGetOrCreate", nil, tricky.MapGetOrCreate},
	"MapGetOrDefaultTest":       {trickySrc, "MapGetOrDefaultTest", nil, tricky.MapGetOrDefaultTest},
	"MapGetOrElse":              {trickySrc, "MapGetOrElse", nil, tricky.MapGetOrElse},
	"MapGetOrInsertDefaultTest": {trickySrc, "MapGetOrInsertDefaultTest", nil, tricky.MapGetOrInsertDefaultTest},
	"MapGetOrInsertTest":        {trickySrc, "MapGetOrInsertTest", nil, tricky.MapGetOrInsertTest},
	"MapGetOrSet":               {trickySrc, "MapGetOrSet", nil, tricky.MapGetOrSet},
	"MapGetSetTest":             {trickySrc, "MapGetSetTest", nil, tricky.MapGetSetTest},
	"MapGroupBy":                {trickySrc, "MapGroupBy", nil, tricky.MapGroupBy},
	"MapGroupByKey":             {trickySrc, "MapGroupByKey", nil, tricky.MapGroupByKey},
	"MapGroupByValueTest":       {trickySrc, "MapGroupByValueTest", nil, tricky.MapGroupByValueTest},
	"MapHasKey":                 {trickySrc, "MapHasKey", nil, tricky.MapHasKey},
	"MapHasKeyAndValueTest":     {trickySrc, "MapHasKeyAndValueTest", nil, tricky.MapHasKeyAndValueTest},
	"MapHasKeyMultiple":         {trickySrc, "MapHasKeyMultiple", nil, tricky.MapHasKeyMultiple},
	"MapHasKeyMultiTest":        {trickySrc, "MapHasKeyMultiTest", nil, tricky.MapHasKeyMultiTest},
	"MapHasKeyNilTest":          {trickySrc, "MapHasKeyNilTest", nil, tricky.MapHasKeyNilTest},
	"MapHasKeySlice":            {trickySrc, "MapHasKeySlice", nil, tricky.MapHasKeySlice},
	"MapHasKeySliceTest":        {trickySrc, "MapHasKeySliceTest", nil, tricky.MapHasKeySliceTest},
	"MapHasKeyTest":             {trickySrc, "MapHasKeyTest", nil, tricky.MapHasKeyTest},
	"MapHasValueCond":           {trickySrc, "MapHasValueCond", nil, tricky.MapHasValueCond},
	"MapHasValuesTest":          {trickySrc, "MapHasValuesTest", nil, tricky.MapHasValuesTest},
	"MapIncrementAll":           {trickySrc, "MapIncrementAll", nil, tricky.MapIncrementAll},
	"MapIncrementValueTest":     {trickySrc, "MapIncrementValueTest", nil, tricky.MapIncrementValueTest},
	"MapIndexBy":                {trickySrc, "MapIndexBy", nil, tricky.MapIndexBy},
	"MapIntersect":              {trickySrc, "MapIntersect", nil, tricky.MapIntersect},
	"MapIntersectKeysFunc":      {trickySrc, "MapIntersectKeysFunc", nil, tricky.MapIntersectKeysFunc},
	"MapIntKey":                 {trickySrc, "MapIntKey", nil, tricky.MapIntKey},
	"MapInvert":                 {trickySrc, "MapInvert", nil, tricky.MapInvert},
	"MapInvertSlice":            {trickySrc, "MapInvertSlice", nil, tricky.MapInvertSlice},
	"MapIsEmptyTest":            {trickySrc, "MapIsEmptyTest", nil, tricky.MapIsEmptyTest},
	"MapIterateDelete":          {trickySrc, "MapIterateDelete", nil, tricky.MapIterateDelete},
	"MapKeepIfTest":             {trickySrc, "MapKeepIfTest", nil, tricky.MapKeepIfTest},
	"MapKeepKeysTest":           {trickySrc, "MapKeepKeysTest", nil, tricky.MapKeepKeysTest},
	"MapKeyDiffTest":            {trickySrc, "MapKeyDiffTest", nil, tricky.MapKeyDiffTest},
	"MapKeyExistsMultiTest":     {trickySrc, "MapKeyExistsMultiTest", nil, tricky.MapKeyExistsMultiTest},
	"MapKeyExistsTest":          {trickySrc, "MapKeyExistsTest", nil, tricky.MapKeyExistsTest},
	"MapKeyIntersectionTest":    {trickySrc, "MapKeyIntersectionTest", nil, tricky.MapKeyIntersectionTest},
	"MapKeysAsSliceTest":        {trickySrc, "MapKeysAsSliceTest", nil, tricky.MapKeysAsSliceTest},
	"MapKeySetTest":             {trickySrc, "MapKeySetTest", nil, tricky.MapKeySetTest},
	"MapKeyShadowing":           {trickySrc, "MapKeyShadowing", nil, tricky.MapKeyShadowing},
	"MapKeysSliceTest":          {trickySrc, "MapKeysSliceTest", nil, tricky.MapKeysSliceTest},
	"MapKeysSorted":             {trickySrc, "MapKeysSorted", nil, tricky.MapKeysSorted},
	"MapKeysSortedTest":         {trickySrc, "MapKeysSortedTest", nil, tricky.MapKeysSortedTest},
	"MapKeysToSlice":            {trickySrc, "MapKeysToSlice", nil, tricky.MapKeysToSlice},
	"MapLastVal":                {trickySrc, "MapLastVal", nil, tricky.MapLastVal},
	"MapLookupOrInsert":         {trickySrc, "MapLookupOrInsert", nil, tricky.MapLookupOrInsert},
	"MapMapTest":                {trickySrc, "MapMapTest", nil, tricky.MapMapTest},
	"MapMergeConditionalTest":   {trickySrc, "MapMergeConditionalTest", nil, tricky.MapMergeConditionalTest},
	"MapMergeDisjointTest":      {trickySrc, "MapMergeDisjointTest", nil, tricky.MapMergeDisjointTest},
	"MapMergeMultiple":          {trickySrc, "MapMergeMultiple", nil, tricky.MapMergeMultiple},
	"MapMergeMultipleTest":      {trickySrc, "MapMergeMultipleTest", nil, tricky.MapMergeMultipleTest},
	"MapMergeNoOverlapTest":     {trickySrc, "MapMergeNoOverlapTest", nil, tricky.MapMergeNoOverlapTest},
	"MapMergeOverwrite":         {trickySrc, "MapMergeOverwrite", nil, tricky.MapMergeOverwrite},
	"MapMergeOverwriteAllTest":  {trickySrc, "MapMergeOverwriteAllTest", nil, tricky.MapMergeOverwriteAllTest},
	"MapMergePredTest":          {trickySrc, "MapMergePredTest", nil, tricky.MapMergePredTest},
	"MapMergePreserveOrigTest":  {trickySrc, "MapMergePreserveOrigTest", nil, tricky.MapMergePreserveOrigTest},
	"MapMergePreserveTest":      {trickySrc, "MapMergePreserveTest", nil, tricky.MapMergePreserveTest},
	"MapMergeSameTest":          {trickySrc, "MapMergeSameTest", nil, tricky.MapMergeSameTest},
	"MapMergeTwo":               {trickySrc, "MapMergeTwo", nil, tricky.MapMergeTwo},
	"MapMergeWithConflict":      {trickySrc, "MapMergeWithConflict", nil, tricky.MapMergeWithConflict},
	"MapMergeWithConflictTest":  {trickySrc, "MapMergeWithConflictTest", nil, tricky.MapMergeWithConflictTest},
	"MapMergeWithFunc":          {trickySrc, "MapMergeWithFunc", nil, tricky.MapMergeWithFunc},
	"MapMinMaxTest":             {trickySrc, "MapMinMaxTest", nil, tricky.MapMinMaxTest},
	"MapNestedAssign":           {trickySrc, "MapNestedAssign", nil, tricky.MapNestedAssign},
	"MapNestedDelete":           {trickySrc, "MapNestedDelete", nil, tricky.MapNestedDelete},
	"MapNestedUpdate":           {trickySrc, "MapNestedUpdate", nil, tricky.MapNestedUpdate},
	"MapNoneMatch":              {trickySrc, "MapNoneMatch", nil, tricky.MapNoneMatch},
	"MapPartition":              {trickySrc, "MapPartition", nil, tricky.MapPartition},
	"MapPick":                   {trickySrc, "MapPick", nil, tricky.MapPick},
	"MapPickBy":                 {trickySrc, "MapPickBy", nil, tricky.MapPickBy},
	"MapPluck":                  {trickySrc, "MapPluck", nil, tricky.MapPluck},
	"MapRangeSafe":              {trickySrc, "MapRangeSafe", nil, tricky.MapRangeSafe},
	"MapRejectKeys":             {trickySrc, "MapRejectKeys", nil, tricky.MapRejectKeys},
	"MapRemoveKeysTest":         {trickySrc, "MapRemoveKeysTest", nil, tricky.MapRemoveKeysTest},
	"MapReplace":                {trickySrc, "MapReplace", nil, tricky.MapReplace},
	"MapReplaceVals":            {trickySrc, "MapReplaceVals", nil, tricky.MapReplaceVals},
	"MapSameKeyValueTest":       {trickySrc, "MapSameKeyValueTest", nil, tricky.MapSameKeyValueTest},
	"MapSelectKeys":             {trickySrc, "MapSelectKeys", nil, tricky.MapSelectKeys},
	"MapSelectTest":             {trickySrc, "MapSelectTest", nil, tricky.MapSelectTest},
	"MapSize":                   {trickySrc, "MapSize", nil, tricky.MapSize},
	"MapSizeHint":               {trickySrc, "MapSizeHint", nil, tricky.MapSizeHint},
	"MapSizeTest":               {trickySrc, "MapSizeTest", nil, tricky.MapSizeTest},
	"MapSliceKeys":              {trickySrc, "MapSliceKeys", nil, tricky.MapSliceKeys},
	"MapSliceKeysTest":          {trickySrc, "MapSliceKeysTest", nil, tricky.MapSliceKeysTest},
	"MapSliceMap":               {trickySrc, "MapSliceMap", nil, tricky.MapSliceMap},
	"MapSliceReduce":            {trickySrc, "MapSliceReduce", nil, tricky.MapSliceReduce},
	"MapSliceToMap":             {trickySrc, "MapSliceToMap", nil, tricky.MapSliceToMap},
	"MapSliceToMapTest":         {trickySrc, "MapSliceToMapTest", nil, tricky.MapSliceToMapTest},
	"MapSliceValues":            {trickySrc, "MapSliceValues", nil, tricky.MapSliceValues},
	"MapSplitTest":              {trickySrc, "MapSplitTest", nil, tricky.MapSplitTest},
	"MapStructUpdate":           {trickySrc, "MapStructUpdate", nil, tricky.MapStructUpdate},
	"MapSumTest":                {trickySrc, "MapSumTest", nil, tricky.MapSumTest},
	"MapSumVals":                {trickySrc, "MapSumVals", nil, tricky.MapSumVals},
	"MapSwap":                   {trickySrc, "MapSwap", nil, tricky.MapSwap},
	"MapSymDiffTest":            {trickySrc, "MapSymDiffTest", nil, tricky.MapSymDiffTest},
	"MapTakeTest":               {trickySrc, "MapTakeTest", nil, tricky.MapTakeTest},
	"MapTakeWhileTest":          {trickySrc, "MapTakeWhileTest", nil, tricky.MapTakeWhileTest},
	"MapTally":                  {trickySrc, "MapTally", nil, tricky.MapTally},
	"MapToSlice":                {trickySrc, "MapToSlice", nil, tricky.MapToSlice},
	"MapToSliceTest":            {trickySrc, "MapToSliceTest", nil, tricky.MapToSliceTest},
	"MapTransformKeys":          {trickySrc, "MapTransformKeys", nil, tricky.MapTransformKeys},
	"MapTransformKeysToSlice":   {trickySrc, "MapTransformKeysToSlice", nil, tricky.MapTransformKeysToSlice},
	"MapTransformVals":          {trickySrc, "MapTransformVals", nil, tricky.MapTransformVals},
	"MapTransposeTest":          {trickySrc, "MapTransposeTest", nil, tricky.MapTransposeTest},
	"MapTwoKeys":                {trickySrc, "MapTwoKeys", nil, tricky.MapTwoKeys},
	"MapUnion":                  {trickySrc, "MapUnion", nil, tricky.MapUnion},
	"MapUnionKeysTest":          {trickySrc, "MapUnionKeysTest", nil, tricky.MapUnionKeysTest},
	"MapUpdateExistingTest":     {trickySrc, "MapUpdateExistingTest", nil, tricky.MapUpdateExistingTest},
	"MapUpdateIfFunc":           {trickySrc, "MapUpdateIfFunc", nil, tricky.MapUpdateIfFunc},
	"MapUpdateIfKeyExistsTest":  {trickySrc, "MapUpdateIfKeyExistsTest", nil, tricky.MapUpdateIfKeyExistsTest},
	"MapUpdateIfTest":           {trickySrc, "MapUpdateIfTest", nil, tricky.MapUpdateIfTest},
	"MapUpdateNestedMapTest":    {trickySrc, "MapUpdateNestedMapTest", nil, tricky.MapUpdateNestedMapTest},
	"MapUpdateNestedTest":       {trickySrc, "MapUpdateNestedTest", nil, tricky.MapUpdateNestedTest},
	"MapUpdateValueDirect":      {trickySrc, "MapUpdateValueDirect", nil, tricky.MapUpdateValueDirect},
	"MapUpdateWithFunc":         {trickySrc, "MapUpdateWithFunc", nil, tricky.MapUpdateWithFunc},
	"MapVals":                   {trickySrc, "MapVals", nil, tricky.MapVals},
	"MapValueDiffTest":          {trickySrc, "MapValueDiffTest", nil, tricky.MapValueDiffTest},
	"MapValueExistsTest":        {trickySrc, "MapValueExistsTest", nil, tricky.MapValueExistsTest},
	"MapValueMaxTest":           {trickySrc, "MapValueMaxTest", nil, tricky.MapValueMaxTest},
	"MapValueSlice":             {trickySrc, "MapValueSlice", nil, tricky.MapValueSlice},
	"MapValueSliceTest":         {trickySrc, "MapValueSliceTest", nil, tricky.MapValueSliceTest},
	"MapValuesToSlice":          {trickySrc, "MapValuesToSlice", nil, tricky.MapValuesToSlice},
	"MapValueSumKeysTest":       {trickySrc, "MapValueSumKeysTest", nil, tricky.MapValueSumKeysTest},
	"MapValueTypes":             {trickySrc, "MapValueTypes", nil, tricky.MapValueTypes},
	"MapWithFuncValue":          {trickySrc, "MapWithFuncValue", nil, tricky.MapWithFuncValue},
	"MapWithFuncValueDirect":    {trickySrc, "MapWithFuncValueDirect", nil, tricky.MapWithFuncValueDirect},
	"MapWithPointerValue":       {trickySrc, "MapWithPointerValue", nil, tricky.MapWithPointerValue},
	"MapZip":                    {trickySrc, "MapZip", nil, tricky.MapZip},
}

var trickyMultiassignTests = map[string]testCase{
	"MultipleNamedReturnCombined": {trickySrc, "MultipleNamedReturnCombined", nil, tricky.MultipleNamedReturnCombined},
}

var trickyNestedTests = map[string]testCase{
	"NestedClosureWithArg":    {trickySrc, "NestedClosureWithArg", nil, tricky.NestedClosureWithArg},
	"NestedMapDeleteNested":   {trickySrc, "NestedMapDeleteNested", nil, tricky.NestedMapDeleteNested},
	"NestedMapGetOrSet":       {trickySrc, "NestedMapGetOrSet", nil, tricky.NestedMapGetOrSet},
	"NestedMapGetWithDefault": {trickySrc, "NestedMapGetWithDefault", nil, tricky.NestedMapGetWithDefault},
	"NestedMapInit":           {trickySrc, "NestedMapInit", nil, tricky.NestedMapInit},
	"NestedMapIterate":        {trickySrc, "NestedMapIterate", nil, tricky.NestedMapIterate},
	"NestedMapSafeAccess":     {trickySrc, "NestedMapSafeAccess", nil, tricky.NestedMapSafeAccess},
	"NestedMapUpdateNested":   {trickySrc, "NestedMapUpdateNested", nil, tricky.NestedMapUpdateNested},
	"NestedMapWithDelete":     {trickySrc, "NestedMapWithDelete", nil, tricky.NestedMapWithDelete},
	"NestedStructAssign":      {trickySrc, "NestedStructAssign", nil, tricky.NestedStructAssign},
}

var trickyPointersTests = map[string]testCase{
	"PointerAddr":                     {trickySrc, "PointerAddr", nil, tricky.PointerAddr},
	"PointerAlias":                    {trickySrc, "PointerAlias", nil, tricky.PointerAlias},
	"PointerArithSim":                 {trickySrc, "PointerArithSim", nil, tricky.PointerArithSim},
	"PointerArrayElementTest":         {trickySrc, "PointerArrayElementTest", nil, tricky.PointerArrayElementTest},
	"PointerArrayIdx":                 {trickySrc, "PointerArrayIdx", nil, tricky.PointerArrayIdx},
	"PointerArrayIndexTest":           {trickySrc, "PointerArrayIndexTest", nil, tricky.PointerArrayIndexTest},
	"PointerAssignChainTest":          {trickySrc, "PointerAssignChainTest", nil, tricky.PointerAssignChainTest},
	"PointerAssignFromDerefTest":      {trickySrc, "PointerAssignFromDerefTest", nil, tricky.PointerAssignFromDerefTest},
	"PointerAssignFromFuncTest":       {trickySrc, "PointerAssignFromFuncTest", nil, tricky.PointerAssignFromFuncTest},
	"PointerAssignFuncResultTest":     {trickySrc, "PointerAssignFuncResultTest", nil, tricky.PointerAssignFuncResultTest},
	"PointerAssignNilTest":            {trickySrc, "PointerAssignNilTest", nil, tricky.PointerAssignNilTest},
	"PointerAssignSameTest":           {trickySrc, "PointerAssignSameTest", nil, tricky.PointerAssignSameTest},
	"PointerAssignThenNilTest":        {trickySrc, "PointerAssignThenNilTest", nil, tricky.PointerAssignThenNilTest},
	"PointerChainTest":                {trickySrc, "PointerChainTest", nil, tricky.PointerChainTest},
	"PointerCheckNilAfterUseTest":     {trickySrc, "PointerCheckNilAfterUseTest", nil, tricky.PointerCheckNilAfterUseTest},
	"PointerCompare":                  {trickySrc, "PointerCompare", nil, tricky.PointerCompare},
	"PointerCompareDiffTest":          {trickySrc, "PointerCompareDiffTest", nil, tricky.PointerCompareDiffTest},
	"PointerCompareTest":              {trickySrc, "PointerCompareTest", nil, tricky.PointerCompareTest},
	"PointerDeref":                    {trickySrc, "PointerDeref", nil, tricky.PointerDeref},
	"PointerDerefAssignTest":          {trickySrc, "PointerDerefAssignTest", nil, tricky.PointerDerefAssignTest},
	"PointerDerefChain":               {trickySrc, "PointerDerefChain", nil, tricky.PointerDerefChain},
	"PointerDerefChainTest":           {trickySrc, "PointerDerefChainTest", nil, tricky.PointerDerefChainTest},
	"PointerDerefModifyTest":          {trickySrc, "PointerDerefModifyTest", nil, tricky.PointerDerefModifyTest},
	"PointerDerefNilCheckTest":        {trickySrc, "PointerDerefNilCheckTest", nil, tricky.PointerDerefNilCheckTest},
	"PointerDerefNilTest":             {trickySrc, "PointerDerefNilTest", nil, tricky.PointerDerefNilTest},
	"PointerDoubleAssignTest":         {trickySrc, "PointerDoubleAssignTest", nil, tricky.PointerDoubleAssignTest},
	"PointerDoubleDerefTest":          {trickySrc, "PointerDoubleDerefTest", nil, tricky.PointerDoubleDerefTest},
	"PointerLevel":                    {trickySrc, "PointerLevel", nil, tricky.PointerLevel},
	"PointerLevelTest":                {trickySrc, "PointerLevelTest", nil, tricky.PointerLevelTest},
	"PointerNilAssign":                {trickySrc, "PointerNilAssign", nil, tricky.PointerNilAssign},
	"PointerNilAssignAfterUseTest":    {trickySrc, "PointerNilAssignAfterUseTest", nil, tricky.PointerNilAssignAfterUseTest},
	"PointerNilAssignT":               {trickySrc, "PointerNilAssignT", nil, tricky.PointerNilAssignT},
	"PointerNilCheckAfterAssignTest":  {trickySrc, "PointerNilCheckAfterAssignTest", nil, tricky.PointerNilCheckAfterAssignTest},
	"PointerNilCheckChain":            {trickySrc, "PointerNilCheckChain", nil, tricky.PointerNilCheckChain},
	"PointerNilCheckDerefTest":        {trickySrc, "PointerNilCheckDerefTest", nil, tricky.PointerNilCheckDerefTest},
	"PointerNilCompare":               {trickySrc, "PointerNilCompare", nil, tricky.PointerNilCompare},
	"PointerNilCompareTest":           {trickySrc, "PointerNilCompareTest", nil, tricky.PointerNilCompareTest},
	"PointerNilDeref":                 {trickySrc, "PointerNilDeref", nil, tricky.PointerNilDeref},
	"PointerNilReassign":              {trickySrc, "PointerNilReassign", nil, tricky.PointerNilReassign},
	"PointerNilSafe":                  {trickySrc, "PointerNilSafe", nil, tricky.PointerNilSafe},
	"PointerNilSafeDeref":             {trickySrc, "PointerNilSafeDeref", nil, tricky.PointerNilSafeDeref},
	"PointerNilSafeDerefTest":         {trickySrc, "PointerNilSafeDerefTest", nil, tricky.PointerNilSafeDerefTest},
	"PointerNilSafeOpTest":            {trickySrc, "PointerNilSafeOpTest", nil, tricky.PointerNilSafeOpTest},
	"PointerNilThenAssignTest":        {trickySrc, "PointerNilThenAssignTest", nil, tricky.PointerNilThenAssignTest},
	"PointerNullObject":               {trickySrc, "PointerNullObject", nil, tricky.PointerNullObject},
	"PointerReassignChainTest":        {trickySrc, "PointerReassignChainTest", nil, tricky.PointerReassignChainTest},
	"PointerReassignmentChain":        {trickySrc, "PointerReassignmentChain", nil, tricky.PointerReassignmentChain},
	"PointerReassignNil":              {trickySrc, "PointerReassignNil", nil, tricky.PointerReassignNil},
	"PointerReassignNilTest":          {trickySrc, "PointerReassignNilTest", nil, tricky.PointerReassignNilTest},
	"PointerReassignTest":             {trickySrc, "PointerReassignTest", nil, tricky.PointerReassignTest},
	"PointerRotate":                   {trickySrc, "PointerRotate", nil, tricky.PointerRotate},
	"PointerSliceElementModifyTest":   {trickySrc, "PointerSliceElementModifyTest", nil, tricky.PointerSliceElementModifyTest},
	"PointerSliceElementSwap":         {trickySrc, "PointerSliceElementSwap", nil, tricky.PointerSliceElementSwap},
	"PointerSliceIndexTest":           {trickySrc, "PointerSliceIndexTest", nil, tricky.PointerSliceIndexTest},
	"PointerSliceIterateTest":         {trickySrc, "PointerSliceIterateTest", nil, tricky.PointerSliceIterateTest},
	"PointerSliceLenTest":             {trickySrc, "PointerSliceLenTest", nil, tricky.PointerSliceLenTest},
	"PointerSliceModifyTest":          {trickySrc, "PointerSliceModifyTest", nil, tricky.PointerSliceModifyTest},
	"PointerSliceNilTest":             {trickySrc, "PointerSliceNilTest", nil, tricky.PointerSliceNilTest},
	"PointerSliceOfPointers":          {trickySrc, "PointerSliceOfPointers", nil, tricky.PointerSliceOfPointers},
	"PointerSliceOfStructTest":        {trickySrc, "PointerSliceOfStructTest", nil, tricky.PointerSliceOfStructTest},
	"PointerStructFieldNilCheckTest":  {trickySrc, "PointerStructFieldNilCheckTest", nil, tricky.PointerStructFieldNilCheckTest},
	"PointerStructFieldNilTest":       {trickySrc, "PointerStructFieldNilTest", nil, tricky.PointerStructFieldNilTest},
	"PointerStructFieldTest":          {trickySrc, "PointerStructFieldTest", nil, tricky.PointerStructFieldTest},
	"PointerStructFld":                {trickySrc, "PointerStructFld", nil, tricky.PointerStructFld},
	"PointerStructMethodTest":         {trickySrc, "PointerStructMethodTest", nil, tricky.PointerStructMethodTest},
	"PointerStructModifyFieldTest":    {trickySrc, "PointerStructModifyFieldTest", nil, tricky.PointerStructModifyFieldTest},
	"PointerStructModifyTest":         {trickySrc, "PointerStructModifyTest", nil, tricky.PointerStructModifyTest},
	"PointerSwap":                     {trickySrc, "PointerSwap", nil, tricky.PointerSwap},
	"PointerSwapChain":                {trickySrc, "PointerSwapChain", nil, tricky.PointerSwapChain},
	"PointerSwapChainTest":            {trickySrc, "PointerSwapChainTest", nil, tricky.PointerSwapChainTest},
	"PointerSwapInArrayTest":          {trickySrc, "PointerSwapInArrayTest", nil, tricky.PointerSwapInArrayTest},
	"PointerSwapInSlice":              {trickySrc, "PointerSwapInSlice", nil, tricky.PointerSwapInSlice},
	"PointerSwapInStruct":             {trickySrc, "PointerSwapInStruct", nil, tricky.PointerSwapInStruct},
	"PointerSwapInStructTest":         {trickySrc, "PointerSwapInStructTest", nil, tricky.PointerSwapInStructTest},
	"PointerSwapMultipleTest":         {trickySrc, "PointerSwapMultipleTest", nil, tricky.PointerSwapMultipleTest},
	"PointerSwapNilSafe":              {trickySrc, "PointerSwapNilSafe", nil, tricky.PointerSwapNilSafe},
	"PointerSwapSimple":               {trickySrc, "PointerSwapSimple", nil, tricky.PointerSwapSimple},
	"PointerSwapStructFieldsTest":     {trickySrc, "PointerSwapStructFieldsTest", nil, tricky.PointerSwapStructFieldsTest},
	"PointerSwapThroughSliceTest":     {trickySrc, "PointerSwapThroughSliceTest", nil, tricky.PointerSwapThroughSliceTest},
	"PointerSwapVals":                 {trickySrc, "PointerSwapVals", nil, tricky.PointerSwapVals},
	"PointerSwapValues":               {trickySrc, "PointerSwapValues", nil, tricky.PointerSwapValues},
	"PointerSwapValuesTest":           {trickySrc, "PointerSwapValuesTest", nil, tricky.PointerSwapValuesTest},
	"PointerSwapViaSliceTest":         {trickySrc, "PointerSwapViaSliceTest", nil, tricky.PointerSwapViaSliceTest},
	"PointerSwapViaTempTest":          {trickySrc, "PointerSwapViaTempTest", nil, tricky.PointerSwapViaTempTest},
	"PointerToArr":                    {trickySrc, "PointerToArr", nil, tricky.PointerToArr},
	"PointerToArray":                  {trickySrc, "PointerToArray", nil, tricky.PointerToArray},
	"PointerToArrayElement":           {trickySrc, "PointerToArrayElement", nil, tricky.PointerToArrayElement},
	"PointerToArrTest":                {trickySrc, "PointerToArrTest", nil, tricky.PointerToArrTest},
	"PointerToChanTest":               {trickySrc, "PointerToChanTest", nil, tricky.PointerToChanTest},
	"PointerToFunc":                   {trickySrc, "PointerToFunc", nil, tricky.PointerToFunc},
	"PointerToFuncResultTest":         {trickySrc, "PointerToFuncResultTest", nil, tricky.PointerToFuncResultTest},
	"PointerToInterface":              {trickySrc, "PointerToInterface", nil, tricky.PointerToInterface},
	"PointerToMapElement":             {trickySrc, "PointerToMapElement", nil, tricky.PointerToMapElement},
	"PointerToMapKey":                 {trickySrc, "PointerToMapKey", nil, tricky.PointerToMapKey},
	"PointerToMapNilTest":             {trickySrc, "PointerToMapNilTest", nil, tricky.PointerToMapNilTest},
	"PointerToMapTest":                {trickySrc, "PointerToMapTest", nil, tricky.PointerToMapTest},
	"PointerToNilAssignTest":          {trickySrc, "PointerToNilAssignTest", nil, tricky.PointerToNilAssignTest},
	"PointerToNilInterface":           {trickySrc, "PointerToNilInterface", nil, tricky.PointerToNilInterface},
	"PointerToNilMap":                 {trickySrc, "PointerToNilMap", nil, tricky.PointerToNilMap},
	"PointerToNilMapLenTest":          {trickySrc, "PointerToNilMapLenTest", nil, tricky.PointerToNilMapLenTest},
	"PointerToNilSliceLenTest":        {trickySrc, "PointerToNilSliceLenTest", nil, tricky.PointerToNilSliceLenTest},
	"PointerToNilSliceTest":           {trickySrc, "PointerToNilSliceTest", nil, tricky.PointerToNilSliceTest},
	"PointerToNilStruct":              {trickySrc, "PointerToNilStruct", nil, tricky.PointerToNilStruct},
	"PointerToNilStructTest":          {trickySrc, "PointerToNilStructTest", nil, tricky.PointerToNilStructTest},
	"PointerToNilTest":                {trickySrc, "PointerToNilTest", nil, tricky.PointerToNilTest},
	"PointerToPointer":                {trickySrc, "PointerToPointer", nil, tricky.PointerToPointer},
	"PointerToPointerAssign":          {trickySrc, "PointerToPointerAssign", nil, tricky.PointerToPointerAssign},
	"PointerToPointerAssignTest":      {trickySrc, "PointerToPointerAssignTest", nil, tricky.PointerToPointerAssignTest},
	"PointerToPointerDerefTest":       {trickySrc, "PointerToPointerDerefTest", nil, tricky.PointerToPointerDerefTest},
	"PointerToSliceAppend":            {trickySrc, "PointerToSliceAppend", nil, tricky.PointerToSliceAppend},
	"PointerToSliceClear":             {trickySrc, "PointerToSliceClear", nil, tricky.PointerToSliceClear},
	"PointerToSliceClearTest":         {trickySrc, "PointerToSliceClearTest", nil, tricky.PointerToSliceClearTest},
	"PointerToSliceElementModifyTest": {trickySrc, "PointerToSliceElementModifyTest", nil, tricky.PointerToSliceElementModifyTest},
	"PointerToSliceLen":               {trickySrc, "PointerToSliceLen", nil, tricky.PointerToSliceLen},
	"PointerToSliceLenCap":            {trickySrc, "PointerToSliceLenCap", nil, tricky.PointerToSliceLenCap},
	"PointerToSliceModify":            {trickySrc, "PointerToSliceModify", nil, tricky.PointerToSliceModify},
	"PointerToSliceNilTest":           {trickySrc, "PointerToSliceNilTest", nil, tricky.PointerToSliceNilTest},
	"PointerToSliceOfNilTest":         {trickySrc, "PointerToSliceOfNilTest", nil, tricky.PointerToSliceOfNilTest},
	"PointerToSliceOfPtrTest":         {trickySrc, "PointerToSliceOfPtrTest", nil, tricky.PointerToSliceOfPtrTest},
	"PointerToSliceOfStructs":         {trickySrc, "PointerToSliceOfStructs", nil, tricky.PointerToSliceOfStructs},
	"PointerToSliceTest":              {trickySrc, "PointerToSliceTest", nil, tricky.PointerToSliceTest},
	"PointerToStructField":            {trickySrc, "PointerToStructField", nil, tricky.PointerToStructField},
	"PointerToStructMethodTest":       {trickySrc, "PointerToStructMethodTest", nil, tricky.PointerToStructMethodTest},
	"PointerToStructNilMethodTest":    {trickySrc, "PointerToStructNilMethodTest", nil, tricky.PointerToStructNilMethodTest},
	"PointerToStructTest":             {trickySrc, "PointerToStructTest", nil, tricky.PointerToStructTest},
}

var trickySlicesTests = map[string]testCase{
	"SliceAll":                     {trickySrc, "SliceAll", nil, tricky.SliceAll},
	"SliceAppendCapTest":           {trickySrc, "SliceAppendCapTest", nil, tricky.SliceAppendCapTest},
	"SliceAppendFunc":              {trickySrc, "SliceAppendFunc", nil, tricky.SliceAppendFunc},
	"SliceAppendIfTest":            {trickySrc, "SliceAppendIfTest", nil, tricky.SliceAppendIfTest},
	"SliceAppendNilTest":           {trickySrc, "SliceAppendNilTest", nil, tricky.SliceAppendNilTest},
	"SliceAppendSliceTest":         {trickySrc, "SliceAppendSliceTest", nil, tricky.SliceAppendSliceTest},
	"SliceBsearch":                 {trickySrc, "SliceBsearch", nil, tricky.SliceBsearch},
	"SliceCartesianProduct":        {trickySrc, "SliceCartesianProduct", nil, tricky.SliceCartesianProduct},
	"SliceChainedSlice":            {trickySrc, "SliceChainedSlice", nil, tricky.SliceChainedSlice},
	"SliceChunk":                   {trickySrc, "SliceChunk", nil, tricky.SliceChunk},
	"SliceChunkByPredTest":         {trickySrc, "SliceChunkByPredTest", nil, tricky.SliceChunkByPredTest},
	"SliceChunkByTest":             {trickySrc, "SliceChunkByTest", nil, tricky.SliceChunkByTest},
	"SliceChunkEveryTest":          {trickySrc, "SliceChunkEveryTest", nil, tricky.SliceChunkEveryTest},
	"SliceClone":                   {trickySrc, "SliceClone", nil, tricky.SliceClone},
	"SliceCombinations":            {trickySrc, "SliceCombinations", nil, tricky.SliceCombinations},
	"SliceCompact":                 {trickySrc, "SliceCompact", nil, tricky.SliceCompact},
	"SliceCompactMap":              {trickySrc, "SliceCompactMap", nil, tricky.SliceCompactMap},
	"SliceContainsAll":             {trickySrc, "SliceContainsAll", nil, tricky.SliceContainsAll},
	"SliceContainsAllTest":         {trickySrc, "SliceContainsAllTest", nil, tricky.SliceContainsAllTest},
	"SliceContainsAnyTest":         {trickySrc, "SliceContainsAnyTest", nil, tricky.SliceContainsAnyTest},
	"SliceContainsNoneTest":        {trickySrc, "SliceContainsNoneTest", nil, tricky.SliceContainsNoneTest},
	"SliceCopyFromMap":             {trickySrc, "SliceCopyFromMap", nil, tricky.SliceCopyFromMap},
	"SliceCopyModifyTest":          {trickySrc, "SliceCopyModifyTest", nil, tricky.SliceCopyModifyTest},
	"SliceCopyReverseTest":         {trickySrc, "SliceCopyReverseTest", nil, tricky.SliceCopyReverseTest},
	"SliceCopySubsetTest":          {trickySrc, "SliceCopySubsetTest", nil, tricky.SliceCopySubsetTest},
	"SliceCountBy":                 {trickySrc, "SliceCountBy", nil, tricky.SliceCountBy},
	"SliceCountTest":               {trickySrc, "SliceCountTest", nil, tricky.SliceCountTest},
	"SliceCountWhileTest":          {trickySrc, "SliceCountWhileTest", nil, tricky.SliceCountWhileTest},
	"SliceCycleTest":               {trickySrc, "SliceCycleTest", nil, tricky.SliceCycleTest},
	"SliceDedupConsecutive":        {trickySrc, "SliceDedupConsecutive", nil, tricky.SliceDedupConsecutive},
	"SliceDeleteFront":             {trickySrc, "SliceDeleteFront", nil, tricky.SliceDeleteFront},
	"SliceDeleteMiddle":            {trickySrc, "SliceDeleteMiddle", nil, tricky.SliceDeleteMiddle},
	"SliceDetect":                  {trickySrc, "SliceDetect", nil, tricky.SliceDetect},
	"SliceDiff":                    {trickySrc, "SliceDiff", nil, tricky.SliceDiff},
	"SliceDifference":              {trickySrc, "SliceDifference", nil, tricky.SliceDifference},
	"SliceDifferenceBy":            {trickySrc, "SliceDifferenceBy", nil, tricky.SliceDifferenceBy},
	"SliceDrop":                    {trickySrc, "SliceDrop", nil, tricky.SliceDrop},
	"SliceDropN":                   {trickySrc, "SliceDropN", nil, tricky.SliceDropN},
	"SliceDropNFunc":               {trickySrc, "SliceDropNFunc", nil, tricky.SliceDropNFunc},
	"SliceDropTest":                {trickySrc, "SliceDropTest", nil, tricky.SliceDropTest},
	"SliceDropWhile":               {trickySrc, "SliceDropWhile", nil, tricky.SliceDropWhile},
	"SliceDropWhileTest":           {trickySrc, "SliceDropWhileTest", nil, tricky.SliceDropWhileTest},
	"SliceEachWithIndex":           {trickySrc, "SliceEachWithIndex", nil, tricky.SliceEachWithIndex},
	"SliceEqual":                   {trickySrc, "SliceEqual", nil, tricky.SliceEqual},
	"SliceExistsTest":              {trickySrc, "SliceExistsTest", nil, tricky.SliceExistsTest},
	"SliceFill":                    {trickySrc, "SliceFill", nil, tricky.SliceFill},
	"SliceFilterKeepTest":          {trickySrc, "SliceFilterKeepTest", nil, tricky.SliceFilterKeepTest},
	"SliceFilterNotTest":           {trickySrc, "SliceFilterNotTest", nil, tricky.SliceFilterNotTest},
	"SliceFindFirstFunc":           {trickySrc, "SliceFindFirstFunc", nil, tricky.SliceFindFirstFunc},
	"SliceFindFirstTest":           {trickySrc, "SliceFindFirstTest", nil, tricky.SliceFindFirstTest},
	"SliceFindIdx":                 {trickySrc, "SliceFindIdx", nil, tricky.SliceFindIdx},
	"SliceFindIndex":               {trickySrc, "SliceFindIndex", nil, tricky.SliceFindIndex},
	"SliceFindIndexTest":           {trickySrc, "SliceFindIndexTest", nil, tricky.SliceFindIndexTest},
	"SliceFindLastPosTest":         {trickySrc, "SliceFindLastPosTest", nil, tricky.SliceFindLastPosTest},
	"SliceFindLastTest":            {trickySrc, "SliceFindLastTest", nil, tricky.SliceFindLastTest},
	"SliceFirst":                   {trickySrc, "SliceFirst", nil, tricky.SliceFirst},
	"SliceFlatten":                 {trickySrc, "SliceFlatten", nil, tricky.SliceFlatten},
	"SliceFlatten2D":               {trickySrc, "SliceFlatten2D", nil, tricky.SliceFlatten2D},
	"SliceFlattenDeep":             {trickySrc, "SliceFlattenDeep", nil, tricky.SliceFlattenDeep},
	"SliceFlattenLevelTest":        {trickySrc, "SliceFlattenLevelTest", nil, tricky.SliceFlattenLevelTest},
	"SliceFlattenManual":           {trickySrc, "SliceFlattenManual", nil, tricky.SliceFlattenManual},
	"SliceFlattenManualTest":       {trickySrc, "SliceFlattenManualTest", nil, tricky.SliceFlattenManualTest},
	"SliceFoldLeft":                {trickySrc, "SliceFoldLeft", nil, tricky.SliceFoldLeft},
	"SliceFromChan":                {trickySrc, "SliceFromChan", nil, tricky.SliceFromChan},
	"SliceGrep":                    {trickySrc, "SliceGrep", nil, tricky.SliceGrep},
	"SliceGroupBy":                 {trickySrc, "SliceGroupBy", nil, tricky.SliceGroupBy},
	"SliceGroupByFld":              {trickySrc, "SliceGroupByFld", nil, tricky.SliceGroupByFld},
	"SliceGroupByMultiple":         {trickySrc, "SliceGroupByMultiple", nil, tricky.SliceGroupByMultiple},
	"SliceGroupConsecutiveTest":    {trickySrc, "SliceGroupConsecutiveTest", nil, tricky.SliceGroupConsecutiveTest},
	"SliceGrow":                    {trickySrc, "SliceGrow", nil, tricky.SliceGrow},
	"SliceGrowWithAppend":          {trickySrc, "SliceGrowWithAppend", nil, tricky.SliceGrowWithAppend},
	"SliceHeadTest":                {trickySrc, "SliceHeadTest", nil, tricky.SliceHeadTest},
	"SliceIndexOf":                 {trickySrc, "SliceIndexOf", nil, tricky.SliceIndexOf},
	"SliceIndexOfFirstTest":        {trickySrc, "SliceIndexOfFirstTest", nil, tricky.SliceIndexOfFirstTest},
	"SliceIndexOfMaxTest":          {trickySrc, "SliceIndexOfMaxTest", nil, tricky.SliceIndexOfMaxTest},
	"SliceIndexOfTest":             {trickySrc, "SliceIndexOfTest", nil, tricky.SliceIndexOfTest},
	"SliceIndexOutOfRange":         {trickySrc, "SliceIndexOutOfRange", nil, tricky.SliceIndexOutOfRange},
	"SliceInitAndModifyTest":       {trickySrc, "SliceInitAndModifyTest", nil, tricky.SliceInitAndModifyTest},
	"SliceInitCapTest":             {trickySrc, "SliceInitCapTest", nil, tricky.SliceInitCapTest},
	"SliceInsertAt":                {trickySrc, "SliceInsertAt", nil, tricky.SliceInsertAt},
	"SliceInsertAtTest":            {trickySrc, "SliceInsertAtTest", nil, tricky.SliceInsertAtTest},
	"SliceInsertFrontTest":         {trickySrc, "SliceInsertFrontTest", nil, tricky.SliceInsertFrontTest},
	"SliceInsertMultipleTest":      {trickySrc, "SliceInsertMultipleTest", nil, tricky.SliceInsertMultipleTest},
	"SliceInsertSliceTest":         {trickySrc, "SliceInsertSliceTest", nil, tricky.SliceInsertSliceTest},
	"SliceInterleaveTest":          {trickySrc, "SliceInterleaveTest", nil, tricky.SliceInterleaveTest},
	"SliceIntersect":               {trickySrc, "SliceIntersect", nil, tricky.SliceIntersect},
	"SliceIntersectBy":             {trickySrc, "SliceIntersectBy", nil, tricky.SliceIntersectBy},
	"SliceIntersectTest":           {trickySrc, "SliceIntersectTest", nil, tricky.SliceIntersectTest},
	"SliceIntersperseTest":         {trickySrc, "SliceIntersperseTest", nil, tricky.SliceIntersperseTest},
	"SliceIsSortedTest":            {trickySrc, "SliceIsSortedTest", nil, tricky.SliceIsSortedTest},
	"SliceLast":                    {trickySrc, "SliceLast", nil, tricky.SliceLast},
	"SliceLastIndexOfTest":         {trickySrc, "SliceLastIndexOfTest", nil, tricky.SliceLastIndexOfTest},
	"SliceLastIndexOfTest2":        {trickySrc, "SliceLastIndexOfTest2", nil, tricky.SliceLastIndexOfTest2},
	"SliceLastNFunc":               {trickySrc, "SliceLastNFunc", nil, tricky.SliceLastNFunc},
	"SliceMakeFromArr":             {trickySrc, "SliceMakeFromArr", nil, tricky.SliceMakeFromArr},
	"SliceMakeZero":                {trickySrc, "SliceMakeZero", nil, tricky.SliceMakeZero},
	"SliceMapEachTest":             {trickySrc, "SliceMapEachTest", nil, tricky.SliceMapEachTest},
	"SliceMapIndex":                {trickySrc, "SliceMapIndex", nil, tricky.SliceMapIndex},
	"SliceMax":                     {trickySrc, "SliceMax", nil, tricky.SliceMax},
	"SliceMaxVal":                  {trickySrc, "SliceMaxVal", nil, tricky.SliceMaxVal},
	"SliceMinIdx":                  {trickySrc, "SliceMinIdx", nil, tricky.SliceMinIdx},
	"SliceMinMax":                  {trickySrc, "SliceMinMax", nil, tricky.SliceMinMax},
	"SliceMinVal":                  {trickySrc, "SliceMinVal", nil, tricky.SliceMinVal},
	"SliceNegativeIndex":           {trickySrc, "SliceNegativeIndex", nil, tricky.SliceNegativeIndex},
	"SliceNilAppend":               {trickySrc, "SliceNilAppend", nil, tricky.SliceNilAppend},
	"SliceNone":                    {trickySrc, "SliceNone", nil, tricky.SliceNone},
	"SliceOfEmptyInterface":        {trickySrc, "SliceOfEmptyInterface", nil, tricky.SliceOfEmptyInterface},
	"SliceOfInterfacesWithTypes":   {trickySrc, "SliceOfInterfacesWithTypes", nil, tricky.SliceOfInterfacesWithTypes},
	"SlicePad":                     {trickySrc, "SlicePad", nil, tricky.SlicePad},
	"SlicePadLeftTest":             {trickySrc, "SlicePadLeftTest", nil, tricky.SlicePadLeftTest},
	"SlicePadRightTest":            {trickySrc, "SlicePadRightTest", nil, tricky.SlicePadRightTest},
	"SlicePartition":               {trickySrc, "SlicePartition", nil, tricky.SlicePartition},
	"SlicePartitionBy":             {trickySrc, "SlicePartitionBy", nil, tricky.SlicePartitionBy},
	"SlicePartitionPosNeg":         {trickySrc, "SlicePartitionPosNeg", nil, tricky.SlicePartitionPosNeg},
	"SlicePartitionTest":           {trickySrc, "SlicePartitionTest", nil, tricky.SlicePartitionTest},
	"SlicePermutation":             {trickySrc, "SlicePermutation", nil, tricky.SlicePermutation},
	"SlicePermuteSimpleTest":       {trickySrc, "SlicePermuteSimpleTest", nil, tricky.SlicePermuteSimpleTest},
	"SlicePluck":                   {trickySrc, "SlicePluck", nil, tricky.SlicePluck},
	"SlicePluckFld":                {trickySrc, "SlicePluckFld", nil, tricky.SlicePluckFld},
	"SlicePrepend":                 {trickySrc, "SlicePrepend", nil, tricky.SlicePrepend},
	"SlicePrependMultipleTest":     {trickySrc, "SlicePrependMultipleTest", nil, tricky.SlicePrependMultipleTest},
	"SlicePrependValueTest":        {trickySrc, "SlicePrependValueTest", nil, tricky.SlicePrependValueTest},
	"SliceProd":                    {trickySrc, "SliceProd", nil, tricky.SliceProd},
	"SliceProduct":                 {trickySrc, "SliceProduct", nil, tricky.SliceProduct},
	"SliceRandomAccess":            {trickySrc, "SliceRandomAccess", nil, tricky.SliceRandomAccess},
	"SliceReduce":                  {trickySrc, "SliceReduce", nil, tricky.SliceReduce},
	"SliceReduceTest":              {trickySrc, "SliceReduceTest", nil, tricky.SliceReduceTest},
	"SliceReject":                  {trickySrc, "SliceReject", nil, tricky.SliceReject},
	"SliceRemoveAtTest":            {trickySrc, "SliceRemoveAtTest", nil, tricky.SliceRemoveAtTest},
	"SliceRemoveDupes":             {trickySrc, "SliceRemoveDupes", nil, tricky.SliceRemoveDupes},
	"SliceRemoveDupSortedTest":     {trickySrc, "SliceRemoveDupSortedTest", nil, tricky.SliceRemoveDupSortedTest},
	"SliceRemoveDupTest":           {trickySrc, "SliceRemoveDupTest", nil, tricky.SliceRemoveDupTest},
	"SliceRemoveIf":                {trickySrc, "SliceRemoveIf", nil, tricky.SliceRemoveIf},
	"SliceRemoveIfKeepTest":        {trickySrc, "SliceRemoveIfKeepTest", nil, tricky.SliceRemoveIfKeepTest},
	"SliceRemoveIfTest":            {trickySrc, "SliceRemoveIfTest", nil, tricky.SliceRemoveIfTest},
	"SliceRemoveLastTest":          {trickySrc, "SliceRemoveLastTest", nil, tricky.SliceRemoveLastTest},
	"SliceRepeat":                  {trickySrc, "SliceRepeat", nil, tricky.SliceRepeat},
	"SliceRepeatNTest":             {trickySrc, "SliceRepeatNTest", nil, tricky.SliceRepeatNTest},
	"SliceReplaceAtTest":           {trickySrc, "SliceReplaceAtTest", nil, tricky.SliceReplaceAtTest},
	"SliceReverseCopy":             {trickySrc, "SliceReverseCopy", nil, tricky.SliceReverseCopy},
	"SliceReverseCopyTest":         {trickySrc, "SliceReverseCopyTest", nil, tricky.SliceReverseCopyTest},
	"SliceReverseInPlace":          {trickySrc, "SliceReverseInPlace", nil, tricky.SliceReverseInPlace},
	"SliceReverseManualTest":       {trickySrc, "SliceReverseManualTest", nil, tricky.SliceReverseManualTest},
	"SliceReverseRangeTest":        {trickySrc, "SliceReverseRangeTest", nil, tricky.SliceReverseRangeTest},
	"SliceRotate":                  {trickySrc, "SliceRotate", nil, tricky.SliceRotate},
	"SliceRotateByTest":            {trickySrc, "SliceRotateByTest", nil, tricky.SliceRotateByTest},
	"SliceRotateLeft":              {trickySrc, "SliceRotateLeft", nil, tricky.SliceRotateLeft},
	"SliceRotateLeftNTest":         {trickySrc, "SliceRotateLeftNTest", nil, tricky.SliceRotateLeftNTest},
	"SliceRotateLeftTest":          {trickySrc, "SliceRotateLeftTest", nil, tricky.SliceRotateLeftTest},
	"SliceRotateRight":             {trickySrc, "SliceRotateRight", nil, tricky.SliceRotateRight},
	"SliceRotateRightNTest":        {trickySrc, "SliceRotateRightNTest", nil, tricky.SliceRotateRightNTest},
	"SliceRotateRightTest":         {trickySrc, "SliceRotateRightTest", nil, tricky.SliceRotateRightTest},
	"SliceSample":                  {trickySrc, "SliceSample", nil, tricky.SliceSample},
	"SliceScan":                    {trickySrc, "SliceScan", nil, tricky.SliceScan},
	"SliceScanLeftTest":            {trickySrc, "SliceScanLeftTest", nil, tricky.SliceScanLeftTest},
	"SliceSelect":                  {trickySrc, "SliceSelect", nil, tricky.SliceSelect},
	"SliceShiftLeftTest":           {trickySrc, "SliceShiftLeftTest", nil, tricky.SliceShiftLeftTest},
	"SliceSlideTest":               {trickySrc, "SliceSlideTest", nil, tricky.SliceSlideTest},
	"SliceSlidingWindowTest":       {trickySrc, "SliceSlidingWindowTest", nil, tricky.SliceSlidingWindowTest},
	"SliceSortBubble":              {trickySrc, "SliceSortBubble", nil, tricky.SliceSortBubble},
	"SliceSortBy":                  {trickySrc, "SliceSortBy", nil, tricky.SliceSortBy},
	"SliceSortByFld":               {trickySrc, "SliceSortByFld", nil, tricky.SliceSortByFld},
	"SliceSortByMultiple":          {trickySrc, "SliceSortByMultiple", nil, tricky.SliceSortByMultiple},
	"SliceSortStable":              {trickySrc, "SliceSortStable", nil, tricky.SliceSortStable},
	"SliceSplice":                  {trickySrc, "SliceSplice", nil, tricky.SliceSplice},
	"SliceSplit":                   {trickySrc, "SliceSplit", nil, tricky.SliceSplit},
	"SliceSplitAtTest":             {trickySrc, "SliceSplitAtTest", nil, tricky.SliceSplitAtTest},
	"SliceSplitByPredTest":         {trickySrc, "SliceSplitByPredTest", nil, tricky.SliceSplitByPredTest},
	"SliceStride":                  {trickySrc, "SliceStride", nil, tricky.SliceStride},
	"SliceStructIndex":             {trickySrc, "SliceStructIndex", nil, tricky.SliceStructIndex},
	"SliceSubset":                  {trickySrc, "SliceSubset", nil, tricky.SliceSubset},
	"SliceSubsliceTest":            {trickySrc, "SliceSubsliceTest", nil, tricky.SliceSubsliceTest},
	"SliceSum":                     {trickySrc, "SliceSum", nil, tricky.SliceSum},
	"SliceSumOddIdx":               {trickySrc, "SliceSumOddIdx", nil, tricky.SliceSumOddIdx},
	"SliceSumRange":                {trickySrc, "SliceSumRange", nil, tricky.SliceSumRange},
	"SliceSumRangeTest":            {trickySrc, "SliceSumRangeTest", nil, tricky.SliceSumRangeTest},
	"SliceSwapElementsTest":        {trickySrc, "SliceSwapElementsTest", nil, tricky.SliceSwapElementsTest},
	"SliceSymmetricDiff":           {trickySrc, "SliceSymmetricDiff", nil, tricky.SliceSymmetricDiff},
	"SliceSymmetricDiffTest":       {trickySrc, "SliceSymmetricDiffTest", nil, tricky.SliceSymmetricDiffTest},
	"SliceTailTest":                {trickySrc, "SliceTailTest", nil, tricky.SliceTailTest},
	"SliceTake":                    {trickySrc, "SliceTake", nil, tricky.SliceTake},
	"SliceTakeDropTest":            {trickySrc, "SliceTakeDropTest", nil, tricky.SliceTakeDropTest},
	"SliceTakeN":                   {trickySrc, "SliceTakeN", nil, tricky.SliceTakeN},
	"SliceTakeNFunc":               {trickySrc, "SliceTakeNFunc", nil, tricky.SliceTakeNFunc},
	"SliceTakeTest":                {trickySrc, "SliceTakeTest", nil, tricky.SliceTakeTest},
	"SliceTakeWhile":               {trickySrc, "SliceTakeWhile", nil, tricky.SliceTakeWhile},
	"SliceTakeWhileDropWhile":      {trickySrc, "SliceTakeWhileDropWhile", nil, tricky.SliceTakeWhileDropWhile},
	"SliceTakeWhileTest":           {trickySrc, "SliceTakeWhileTest", nil, tricky.SliceTakeWhileTest},
	"SliceTee":                     {trickySrc, "SliceTee", nil, tricky.SliceTee},
	"SliceTranspose":               {trickySrc, "SliceTranspose", nil, tricky.SliceTranspose},
	"SliceTranspose2D":             {trickySrc, "SliceTranspose2D", nil, tricky.SliceTranspose2D},
	"SliceTruncate":                {trickySrc, "SliceTruncate", nil, tricky.SliceTruncate},
	"SliceUnion":                   {trickySrc, "SliceUnion", nil, tricky.SliceUnion},
	"SliceUniqBy":                  {trickySrc, "SliceUniqBy", nil, tricky.SliceUniqBy},
	"SliceUnique":                  {trickySrc, "SliceUnique", nil, tricky.SliceUnique},
	"SliceUniqueCountTest":         {trickySrc, "SliceUniqueCountTest", nil, tricky.SliceUniqueCountTest},
	"SliceUniquePreserveOrderTest": {trickySrc, "SliceUniquePreserveOrderTest", nil, tricky.SliceUniquePreserveOrderTest},
	"SliceUniquePreserveTest":      {trickySrc, "SliceUniquePreserveTest", nil, tricky.SliceUniquePreserveTest},
	"SliceUnzip":                   {trickySrc, "SliceUnzip", nil, tricky.SliceUnzip},
	"SliceWindow":                  {trickySrc, "SliceWindow", nil, tricky.SliceWindow},
	"SliceWithout":                 {trickySrc, "SliceWithout", nil, tricky.SliceWithout},
	"SliceZip":                     {trickySrc, "SliceZip", nil, tricky.SliceZip},
	"SliceZipMap":                  {trickySrc, "SliceZipMap", nil, tricky.SliceZipMap},
	"SliceZipMapTest":              {trickySrc, "SliceZipMapTest", nil, tricky.SliceZipMapTest},
	"SliceZipTest":                 {trickySrc, "SliceZipTest", nil, tricky.SliceZipTest},
	"SliceZipWith":                 {trickySrc, "SliceZipWith", nil, tricky.SliceZipWith},
	"SliceZipWithIndexTest":        {trickySrc, "SliceZipWithIndexTest", nil, tricky.SliceZipWithIndexTest},
}

var trickyStructsTests = map[string]testCase{
	"StructAnon":                         {trickySrc, "StructAnon", nil, tricky.StructAnon},
	"StructAnonymousField":               {trickySrc, "StructAnonymousField", nil, tricky.StructAnonymousField},
	"StructCompareDiff":                  {trickySrc, "StructCompareDiff", nil, tricky.StructCompareDiff},
	"StructCompareDiffTest":              {trickySrc, "StructCompareDiffTest", nil, tricky.StructCompareDiffTest},
	"StructCompareDiffTypeTest":          {trickySrc, "StructCompareDiffTypeTest", nil, tricky.StructCompareDiffTypeTest},
	"StructCompareEqual":                 {trickySrc, "StructCompareEqual", nil, tricky.StructCompareEqual},
	"StructCompareNil":                   {trickySrc, "StructCompareNil", nil, tricky.StructCompareNil},
	"StructCompareNilPtrTest":            {trickySrc, "StructCompareNilPtrTest", nil, tricky.StructCompareNilPtrTest},
	"StructCompareSameTest":              {trickySrc, "StructCompareSameTest", nil, tricky.StructCompareSameTest},
	"StructCopyDeep":                     {trickySrc, "StructCopyDeep", nil, tricky.StructCopyDeep},
	"StructCopyPointerTest":              {trickySrc, "StructCopyPointerTest", nil, tricky.StructCopyPointerTest},
	"StructCopyValueTest":                {trickySrc, "StructCopyValueTest", nil, tricky.StructCopyValueTest},
	"StructEmbeddedAccessTest":           {trickySrc, "StructEmbeddedAccessTest", nil, tricky.StructEmbeddedAccessTest},
	"StructEmbeddedFldAccess":            {trickySrc, "StructEmbeddedFldAccess", nil, tricky.StructEmbeddedFldAccess},
	"StructEmbeddedInterface":            {trickySrc, "StructEmbeddedInterface", nil, tricky.StructEmbeddedInterface},
	"StructEmbeddedMethodOverride":       {trickySrc, "StructEmbeddedMethodOverride", nil, tricky.StructEmbeddedMethodOverride},
	"StructEmbeddedMethodOverrideTest":   {trickySrc, "StructEmbeddedMethodOverrideTest", nil, tricky.StructEmbeddedMethodOverrideTest},
	"StructEmbeddedNil":                  {trickySrc, "StructEmbeddedNil", nil, tricky.StructEmbeddedNil},
	"StructEmbeddedNilCheckTest":         {trickySrc, "StructEmbeddedNilCheckTest", nil, tricky.StructEmbeddedNilCheckTest},
	"StructEmbeddedNilDerefTest":         {trickySrc, "StructEmbeddedNilDerefTest", nil, tricky.StructEmbeddedNilDerefTest},
	"StructEmbeddedNilFld":               {trickySrc, "StructEmbeddedNilFld", nil, tricky.StructEmbeddedNilFld},
	"StructEmbeddedNilMethodTest":        {trickySrc, "StructEmbeddedNilMethodTest", nil, tricky.StructEmbeddedNilMethodTest},
	"StructEmbeddedOverride":             {trickySrc, "StructEmbeddedOverride", nil, tricky.StructEmbeddedOverride},
	"StructEmbeddedPtrInitTest":          {trickySrc, "StructEmbeddedPtrInitTest", nil, tricky.StructEmbeddedPtrInitTest},
	"StructEmbeddedPtrNilTest":           {trickySrc, "StructEmbeddedPtrNilTest", nil, tricky.StructEmbeddedPtrNilTest},
	"StructEmpty":                        {trickySrc, "StructEmpty", nil, tricky.StructEmpty},
	"StructFieldInitTest":                {trickySrc, "StructFieldInitTest", nil, tricky.StructFieldInitTest},
	"StructFieldModifyViaPtrTest":        {trickySrc, "StructFieldModifyViaPtrTest", nil, tricky.StructFieldModifyViaPtrTest},
	"StructFieldPointerModify":           {trickySrc, "StructFieldPointerModify", nil, tricky.StructFieldPointerModify},
	"StructFieldPointerModifyTest":       {trickySrc, "StructFieldPointerModifyTest", nil, tricky.StructFieldPointerModifyTest},
	"StructFieldPtr":                     {trickySrc, "StructFieldPtr", nil, tricky.StructFieldPtr},
	"StructFieldPtrTest":                 {trickySrc, "StructFieldPtrTest", nil, tricky.StructFieldPtrTest},
	"StructFieldShadow":                  {trickySrc, "StructFieldShadow", nil, tricky.StructFieldShadow},
	"StructFieldShadowTest":              {trickySrc, "StructFieldShadowTest", nil, tricky.StructFieldShadowTest},
	"StructFldModify":                    {trickySrc, "StructFldModify", nil, tricky.StructFldModify},
	"StructFldPtrModify":                 {trickySrc, "StructFldPtrModify", nil, tricky.StructFldPtrModify},
	"StructInterface":                    {trickySrc, "StructInterface", nil, tricky.StructInterface},
	"StructMethodChain":                  {trickySrc, "StructMethodChain", nil, tricky.StructMethodChain},
	"StructMethodChainNilTest":           {trickySrc, "StructMethodChainNilTest", nil, tricky.StructMethodChainNilTest},
	"StructMethodChainTest":              {trickySrc, "StructMethodChainTest", nil, tricky.StructMethodChainTest},
	"StructMethodEmbeddedTest":           {trickySrc, "StructMethodEmbeddedTest", nil, tricky.StructMethodEmbeddedTest},
	"StructMethodNilPtrTest":             {trickySrc, "StructMethodNilPtrTest", nil, tricky.StructMethodNilPtrTest},
	"StructMethodOnAddr":                 {trickySrc, "StructMethodOnAddr", nil, tricky.StructMethodOnAddr},
	"StructMethodOnEmbeddedTest":         {trickySrc, "StructMethodOnEmbeddedTest", nil, tricky.StructMethodOnEmbeddedTest},
	"StructMethodOnNilPtrTest":           {trickySrc, "StructMethodOnNilPtrTest", nil, tricky.StructMethodOnNilPtrTest},
	"StructMethodOnNilReceiver":          {trickySrc, "StructMethodOnNilReceiver", nil, tricky.StructMethodOnNilReceiver},
	"StructMethodOnValTest":              {trickySrc, "StructMethodOnValTest", nil, tricky.StructMethodOnValTest},
	"StructMethodOnValueCopy":            {trickySrc, "StructMethodOnValueCopy", nil, tricky.StructMethodOnValueCopy},
	"StructMethodPtrRecTest":             {trickySrc, "StructMethodPtrRecTest", nil, tricky.StructMethodPtrRecTest},
	"StructMethodValRec":                 {trickySrc, "StructMethodValRec", nil, tricky.StructMethodValRec},
	"StructMethodValue":                  {trickySrc, "StructMethodValue", nil, tricky.StructMethodValue},
	"StructMethodValueReceiverTest":      {trickySrc, "StructMethodValueReceiverTest", nil, tricky.StructMethodValueReceiverTest},
	"StructMethodWithPointerReceiver":    {trickySrc, "StructMethodWithPointerReceiver", nil, tricky.StructMethodWithPointerReceiver},
	"StructMethodWithVariadic":           {trickySrc, "StructMethodWithVariadic", nil, tricky.StructMethodWithVariadic},
	"StructModifyViaPointerTest":         {trickySrc, "StructModifyViaPointerTest", nil, tricky.StructModifyViaPointerTest},
	"StructNestedAssign":                 {trickySrc, "StructNestedAssign", nil, tricky.StructNestedAssign},
	"StructNestedInitTest":               {trickySrc, "StructNestedInitTest", nil, tricky.StructNestedInitTest},
	"StructNestedMethodTest":             {trickySrc, "StructNestedMethodTest", nil, tricky.StructNestedMethodTest},
	"StructNestedPtrTest":                {trickySrc, "StructNestedPtrTest", nil, tricky.StructNestedPtrTest},
	"StructNilField":                     {trickySrc, "StructNilField", nil, tricky.StructNilField},
	"StructNilFieldDerefTest":            {trickySrc, "StructNilFieldDerefTest", nil, tricky.StructNilFieldDerefTest},
	"StructNilFieldHolder":               {trickySrc, "StructNilFieldHolder", nil, tricky.StructNilFieldHolder},
	"StructNilFieldInitTest":             {trickySrc, "StructNilFieldInitTest", nil, tricky.StructNilFieldInitTest},
	"StructNilPointerMethod":             {trickySrc, "StructNilPointerMethod", nil, tricky.StructNilPointerMethod},
	"StructNilSafeMethodTest":            {trickySrc, "StructNilSafeMethodTest", nil, tricky.StructNilSafeMethodTest},
	"StructPointerMethodChain":           {trickySrc, "StructPointerMethodChain", nil, tricky.StructPointerMethodChain},
	"StructPtrMethod":                    {trickySrc, "StructPtrMethod", nil, tricky.StructPtrMethod},
	"StructPtrMethodOnNilTest":           {trickySrc, "StructPtrMethodOnNilTest", nil, tricky.StructPtrMethodOnNilTest},
	"StructSelfRef":                      {trickySrc, "StructSelfRef", nil, tricky.StructSelfRef},
	"StructSliceAppend":                  {trickySrc, "StructSliceAppend", nil, tricky.StructSliceAppend},
	"StructSliceFieldAppendTest":         {trickySrc, "StructSliceFieldAppendTest", nil, tricky.StructSliceFieldAppendTest},
	"StructSliceOfPointers":              {trickySrc, "StructSliceOfPointers", nil, tricky.StructSliceOfPointers},
	"StructSliceOfSlices":                {trickySrc, "StructSliceOfSlices", nil, tricky.StructSliceOfSlices},
	"StructValidation":                   {trickySrc, "StructValidation", nil, tricky.StructValidation},
	"StructValidationTest":               {trickySrc, "StructValidationTest", nil, tricky.StructValidationTest},
	"StructWithAnonymousFunc":            {trickySrc, "StructWithAnonymousFunc", nil, tricky.StructWithAnonymousFunc},
	"StructWithArrFieldTest":             {trickySrc, "StructWithArrFieldTest", nil, tricky.StructWithArrFieldTest},
	"StructWithArrInitTest":              {trickySrc, "StructWithArrInitTest", nil, tricky.StructWithArrInitTest},
	"StructWithBoolFld":                  {trickySrc, "StructWithBoolFld", nil, tricky.StructWithBoolFld},
	"StructWithChanFieldTest":            {trickySrc, "StructWithChanFieldTest", nil, tricky.StructWithChanFieldTest},
	"StructWithChanFld":                  {trickySrc, "StructWithChanFld", nil, tricky.StructWithChanFld},
	"StructWithChannel":                  {trickySrc, "StructWithChannel", nil, tricky.StructWithChannel},
	"StructWithChanNilInitTest":          {trickySrc, "StructWithChanNilInitTest", nil, tricky.StructWithChanNilInitTest},
	"StructWithChanOfChan":               {trickySrc, "StructWithChanOfChan", nil, tricky.StructWithChanOfChan},
	"StructWithChanOfChanTest":           {trickySrc, "StructWithChanOfChanTest", nil, tricky.StructWithChanOfChanTest},
	"StructWithChanTest":                 {trickySrc, "StructWithChanTest", nil, tricky.StructWithChanTest},
	"StructWithChanTest2":                {trickySrc, "StructWithChanTest2", nil, tricky.StructWithChanTest2},
	"StructWithComputedField":            {trickySrc, "StructWithComputedField", nil, tricky.StructWithComputedField},
	"StructWithComputedFieldTest":        {trickySrc, "StructWithComputedFieldTest", nil, tricky.StructWithComputedFieldTest},
	"StructWithComputedFld":              {trickySrc, "StructWithComputedFld", nil, tricky.StructWithComputedFld},
	"StructWithDoublePointer":            {trickySrc, "StructWithDoublePointer", nil, tricky.StructWithDoublePointer},
	"StructWithEmbeddedNilPtrTest":       {trickySrc, "StructWithEmbeddedNilPtrTest", nil, tricky.StructWithEmbeddedNilPtrTest},
	"StructWithEmbeddedPointer":          {trickySrc, "StructWithEmbeddedPointer", nil, tricky.StructWithEmbeddedPointer},
	"StructWithEmbeddedPtrTest":          {trickySrc, "StructWithEmbeddedPtrTest", nil, tricky.StructWithEmbeddedPtrTest},
	"StructWithEmbeddedTest":             {trickySrc, "StructWithEmbeddedTest", nil, tricky.StructWithEmbeddedTest},
	"StructWithEmptySliceTest":           {trickySrc, "StructWithEmptySliceTest", nil, tricky.StructWithEmptySliceTest},
	"StructWithFldValidation":            {trickySrc, "StructWithFldValidation", nil, tricky.StructWithFldValidation},
	"StructWithFloatField":               {trickySrc, "StructWithFloatField", nil, tricky.StructWithFloatField},
	"StructWithFloatFldTest":             {trickySrc, "StructWithFloatFldTest", nil, tricky.StructWithFloatFldTest},
	"StructWithFunc":                     {trickySrc, "StructWithFunc", nil, tricky.StructWithFunc},
	"StructWithFuncField":                {trickySrc, "StructWithFuncField", nil, tricky.StructWithFuncField},
	"StructWithFuncFieldCallTest":        {trickySrc, "StructWithFuncFieldCallTest", nil, tricky.StructWithFuncFieldCallTest},
	"StructWithFuncFieldNilTest":         {trickySrc, "StructWithFuncFieldNilTest", nil, tricky.StructWithFuncFieldNilTest},
	"StructWithFuncFieldTest":            {trickySrc, "StructWithFuncFieldTest", nil, tricky.StructWithFuncFieldTest},
	"StructWithFuncFldCall":              {trickySrc, "StructWithFuncFldCall", nil, tricky.StructWithFuncFldCall},
	"StructWithFuncFldExec":              {trickySrc, "StructWithFuncFldExec", nil, tricky.StructWithFuncFldExec},
	"StructWithFuncMap":                  {trickySrc, "StructWithFuncMap", nil, tricky.StructWithFuncMap},
	"StructWithFuncPtrTest":              {trickySrc, "StructWithFuncPtrTest", nil, tricky.StructWithFuncPtrTest},
	"StructWithFuncReturningStruct":      {trickySrc, "StructWithFuncReturningStruct", nil, tricky.StructWithFuncReturningStruct},
	"StructWithFuncSlice":                {trickySrc, "StructWithFuncSlice", nil, tricky.StructWithFuncSlice},
	"StructWithFuncSliceComplex":         {trickySrc, "StructWithFuncSliceComplex", nil, tricky.StructWithFuncSliceComplex},
	"StructWithInitFunc":                 {trickySrc, "StructWithInitFunc", nil, tricky.StructWithInitFunc},
	"StructWithInterfaceFldTest":         {trickySrc, "StructWithInterfaceFldTest", nil, tricky.StructWithInterfaceFldTest},
	"StructWithInterfaceMap":             {trickySrc, "StructWithInterfaceMap", nil, tricky.StructWithInterfaceMap},
	"StructWithInterfaceSlice":           {trickySrc, "StructWithInterfaceSlice", nil, tricky.StructWithInterfaceSlice},
	"StructWithIntField":                 {trickySrc, "StructWithIntField", nil, tricky.StructWithIntField},
	"StructWithIntSliceTest":             {trickySrc, "StructWithIntSliceTest", nil, tricky.StructWithIntSliceTest},
	"StructWithLazyFieldTest":            {trickySrc, "StructWithLazyFieldTest", nil, tricky.StructWithLazyFieldTest},
	"StructWithLazyFld":                  {trickySrc, "StructWithLazyFld", nil, tricky.StructWithLazyFld},
	"StructWithLazyInit":                 {trickySrc, "StructWithLazyInit", nil, tricky.StructWithLazyInit},
	"StructWithMapFld":                   {trickySrc, "StructWithMapFld", nil, tricky.StructWithMapFld},
	"StructWithMapInitTest":              {trickySrc, "StructWithMapInitTest", nil, tricky.StructWithMapInitTest},
	"StructWithMapMakeTest":              {trickySrc, "StructWithMapMakeTest", nil, tricky.StructWithMapMakeTest},
	"StructWithMapNilInit":               {trickySrc, "StructWithMapNilInit", nil, tricky.StructWithMapNilInit},
	"StructWithMapNilInitTest":           {trickySrc, "StructWithMapNilInitTest", nil, tricky.StructWithMapNilInitTest},
	"StructWithMapOfPtrTest":             {trickySrc, "StructWithMapOfPtrTest", nil, tricky.StructWithMapOfPtrTest},
	"StructWithMapOfSlices":              {trickySrc, "StructWithMapOfSlices", nil, tricky.StructWithMapOfSlices},
	"StructWithMapOfStructs":             {trickySrc, "StructWithMapOfStructs", nil, tricky.StructWithMapOfStructs},
	"StructWithMapPointer":               {trickySrc, "StructWithMapPointer", nil, tricky.StructWithMapPointer},
	"StructWithMapRangeDel":              {trickySrc, "StructWithMapRangeDel", nil, tricky.StructWithMapRangeDel},
	"StructWithMethodClosure":            {trickySrc, "StructWithMethodClosure", nil, tricky.StructWithMethodClosure},
	"StructWithMethodPointer":            {trickySrc, "StructWithMethodPointer", nil, tricky.StructWithMethodPointer},
	"StructWithNestedFunc":               {trickySrc, "StructWithNestedFunc", nil, tricky.StructWithNestedFunc},
	"StructWithNestedPointer":            {trickySrc, "StructWithNestedPointer", nil, tricky.StructWithNestedPointer},
	"StructWithNestedSlice":              {trickySrc, "StructWithNestedSlice", nil, tricky.StructWithNestedSlice},
	"StructWithNilChan":                  {trickySrc, "StructWithNilChan", nil, tricky.StructWithNilChan},
	"StructWithNilChanFieldTest":         {trickySrc, "StructWithNilChanFieldTest", nil, tricky.StructWithNilChanFieldTest},
	"StructWithNilChanFld":               {trickySrc, "StructWithNilChanFld", nil, tricky.StructWithNilChanFld},
	"StructWithNilFieldInitTest":         {trickySrc, "StructWithNilFieldInitTest", nil, tricky.StructWithNilFieldInitTest},
	"StructWithNilPtrTest":               {trickySrc, "StructWithNilPtrTest", nil, tricky.StructWithNilPtrTest},
	"StructWithNilSlice":                 {trickySrc, "StructWithNilSlice", nil, tricky.StructWithNilSlice},
	"StructWithNilSliceFieldTest":        {trickySrc, "StructWithNilSliceFieldTest", nil, tricky.StructWithNilSliceFieldTest},
	"StructWithPointerField":             {trickySrc, "StructWithPointerField", nil, tricky.StructWithPointerField},
	"StructWithPointerInterface":         {trickySrc, "StructWithPointerInterface", nil, tricky.StructWithPointerInterface},
	"StructWithPointerMap":               {trickySrc, "StructWithPointerMap", nil, tricky.StructWithPointerMap},
	"StructWithPointerSlice":             {trickySrc, "StructWithPointerSlice", nil, tricky.StructWithPointerSlice},
	"StructWithPointerToInterface":       {trickySrc, "StructWithPointerToInterface", nil, tricky.StructWithPointerToInterface},
	"StructWithPointerToMap":             {trickySrc, "StructWithPointerToMap", nil, tricky.StructWithPointerToMap},
	"StructWithPointerToSelf":            {trickySrc, "StructWithPointerToSelf", nil, tricky.StructWithPointerToSelf},
	"StructWithPtrFld":                   {trickySrc, "StructWithPtrFld", nil, tricky.StructWithPtrFld},
	"StructWithPtrMethodTest":            {trickySrc, "StructWithPtrMethodTest", nil, tricky.StructWithPtrMethodTest},
	"StructWithPtrSliceFieldTest":        {trickySrc, "StructWithPtrSliceFieldTest", nil, tricky.StructWithPtrSliceFieldTest},
	"StructWithPtrToStructTest":          {trickySrc, "StructWithPtrToStructTest", nil, tricky.StructWithPtrToStructTest},
	"StructWithRecursiveType":            {trickySrc, "StructWithRecursiveType", nil, tricky.StructWithRecursiveType},
	"StructWithSelfRefPointer":           {trickySrc, "StructWithSelfRefPointer", nil, tricky.StructWithSelfRefPointer},
	"StructWithSliceAppendMethodTest":    {trickySrc, "StructWithSliceAppendMethodTest", nil, tricky.StructWithSliceAppendMethodTest},
	"StructWithSliceAppendTest":          {trickySrc, "StructWithSliceAppendTest", nil, tricky.StructWithSliceAppendTest},
	"StructWithSliceFieldNamed":          {trickySrc, "StructWithSliceFieldNamed", nil, tricky.StructWithSliceFieldNamed},
	"StructWithSliceFld":                 {trickySrc, "StructWithSliceFld", nil, tricky.StructWithSliceFld},
	"StructWithSliceMakeTest":            {trickySrc, "StructWithSliceMakeTest", nil, tricky.StructWithSliceMakeTest},
	"StructWithSliceMethods":             {trickySrc, "StructWithSliceMethods", nil, tricky.StructWithSliceMethods},
	"StructWithSliceNil":                 {trickySrc, "StructWithSliceNil", nil, tricky.StructWithSliceNil},
	"StructWithSliceNilInitTest":         {trickySrc, "StructWithSliceNilInitTest", nil, tricky.StructWithSliceNilInitTest},
	"StructWithSliceOfMaps":              {trickySrc, "StructWithSliceOfMaps", nil, tricky.StructWithSliceOfMaps},
	"StructWithSliceOfPointersToStructs": {trickySrc, "StructWithSliceOfPointersToStructs", nil, tricky.StructWithSliceOfPointersToStructs},
	"StructWithSliceOfPtrTest":           {trickySrc, "StructWithSliceOfPtrTest", nil, tricky.StructWithSliceOfPtrTest},
	"StructWithSlicePointer":             {trickySrc, "StructWithSlicePointer", nil, tricky.StructWithSlicePointer},
	"StructWithStringFld":                {trickySrc, "StructWithStringFld", nil, tricky.StructWithStringFld},
	"StructWithTag":                      {trickySrc, "StructWithTag", nil, tricky.StructWithTag},
	"StructWithTwoFlds":                  {trickySrc, "StructWithTwoFlds", nil, tricky.StructWithTwoFlds},
	"StructWithUintField":                {trickySrc, "StructWithUintField", nil, tricky.StructWithUintField},
	"StructWithUintFldTest":              {trickySrc, "StructWithUintFldTest", nil, tricky.StructWithUintFldTest},
	"StructWithValidation":               {trickySrc, "StructWithValidation", nil, tricky.StructWithValidation},
	"StructZeroInitTest":                 {trickySrc, "StructZeroInitTest", nil, tricky.StructZeroInitTest},
	"StructZeroValueCheckTest":           {trickySrc, "StructZeroValueCheckTest", nil, tricky.StructZeroValueCheckTest},
}

var typeconvTests = map[string]testCase{
	"IntToFloat64":           {typeconvSrc, "IntToFloat64", nil, typeconv.IntToFloat64},
	"Float64Arithmetic":      {typeconvSrc, "Float64Arithmetic", nil, typeconv.Float64Arithmetic},
	"StringToByteConversion": {typeconvSrc, "StringToByteConversion", nil, typeconv.StringToByteConversion},
	"IntStringConversion":    {typeconvSrc, "IntStringConversion", nil, typeconv.IntStringConversion},
	"StringIntConversion":    {typeconvSrc, "StringIntConversion", nil, typeconv.StringIntConversion},
	// Parameterized tests
	"IntToString":     {typeconvSrc, "IntToString", []any{12345}, typeconv.IntToString},
	"StringToInt":     {typeconvSrc, "StringToInt", []any{"54321"}, typeconv.StringToInt},
	"IntToFloatToInt": {typeconvSrc, "IntToFloatToInt", []any{42}, typeconv.IntToFloatToInt},
}

var variablesTests = map[string]testCase{
	"DeclareAndUse":   {variablesSrc, "DeclareAndUse", nil, variables.DeclareAndUse},
	"Reassignment":    {variablesSrc, "Reassignment", nil, variables.Reassignment},
	"MultipleDecl":    {variablesSrc, "MultipleDecl", nil, variables.MultipleDecl},
	"ZeroValues":      {variablesSrc, "ZeroValues", nil, variables.ZeroValues},
	"StringZeroValue": {variablesSrc, "StringZeroValue", nil, variables.StringZeroValue},
	"Shadowing":       {variablesSrc, "Shadowing", nil, variables.Shadowing},
	// Parameterized tests
	"SumThree":   {variablesSrc, "SumThree", []any{10, 20, 30}, variables.SumThree},
	"Multiply":   {variablesSrc, "Multiply", []any{6, 7}, variables.Multiply},
	"Max":        {variablesSrc, "Max", []any{100, 42}, variables.Max},
	"IsPositive": {variablesSrc, "IsPositive", []any{5}, variables.IsPositive},
}

var testSetsMap = map[string]testSet{
	"advanced":           {name: "advanced", src: advancedSrc, tests: advancedTests},
	"algorithms":         {name: "algorithms", src: algorithmsSrc, tests: algorithmsTests},
	"arithmetic":         {name: "arithmetic", src: arithmeticSrc, tests: arithmeticTests},
	"autowrap":           {name: "autowrap", src: autowrapSrc, tests: autowrapTests},
	"bitwise":            {name: "bitwise", src: bitwiseSrc, tests: bitwiseTests},
	"closures":           {name: "closures", src: closuresSrc, tests: closuresTests},
	"closures_advanced":  {name: "closures_advanced", src: closuresAdvancedSrc, tests: closures_advancedTests},
	"controlflow":        {name: "controlflow", src: controlflowSrc, tests: controlflowTests},
	"cornercases":        {name: "cornercases", src: cornercasesSrc, tests: cornercasesTests},
	"edgecases":          {name: "edgecases", src: edgecasesSrc, tests: edgecasesTests},
	"external":           {name: "external", src: externalSrc, tests: externalTests},
	"functions":          {name: "functions", src: functionsSrc, tests: functionsTests},
	"goroutine":          {name: "goroutine", src: goroutineSrc, tests: goroutineTests},
	"channels":           {name: "channels", src: channelsSrc, tests: channelsTests},
	"init":               {name: "init", src: initSrc, tests: initTests},
	"initialize":         {name: "initialize", src: initializeSrc, tests: initializeTests},
	"leetcode_hard":      {name: "leetcode_hard", src: leetcodeHardSrc, tests: leetcode_hardTests},
	"mapadvanced":        {name: "mapadvanced", src: mapadvancedSrc, tests: mapadvancedTests},
	"maps":               {name: "maps", src: mapsSrc, tests: mapsTests},
	"multiassign":        {name: "multiassign", src: multiassignSrc, tests: multiassignTests},
	"namedreturn":        {name: "namedreturn", src: namedreturnSrc, tests: namedreturnTests},
	"recursion":          {name: "recursion", src: recursionSrc, tests: recursionTests},
	"resolved_issue":     {name: "resolved_issue", src: resolvedIssueSrc, tests: resolved_issueTests},
	"scope":              {name: "scope", src: scopeSrc, tests: scopeTests},
	"slices":             {name: "slices", src: slicesSrc, tests: slicesTests},
	"slicing":            {name: "slicing", src: slicingSrc, tests: slicingTests},
	"strings_pkg":        {name: "strings_pkg", src: stringsPkgSrc, tests: strings_pkgTests},
	"structs":            {name: "structs", src: structsSrc, tests: structsTests},
	"switch":             {name: "switch", src: switchSrc, tests: switchTests},
	"tricky/closures":    {name: "tricky/closures", src: trickySrc, tests: trickyClosuresTests},
	"tricky/defer":       {name: "tricky/defer", src: trickySrc, tests: trickyDeferTests},
	"tricky/interfaces":  {name: "tricky/interfaces", src: trickySrc, tests: trickyInterfacesTests},
	"tricky/maps":        {name: "tricky/maps", src: trickySrc, tests: trickyMapsTests},
	"tricky/multiassign": {name: "tricky/multiassign", src: trickySrc, tests: trickyMultiassignTests},
	"tricky/nested":      {name: "tricky/nested", src: trickySrc, tests: trickyNestedTests},
	"tricky/pointers":    {name: "tricky/pointers", src: trickySrc, tests: trickyPointersTests},
	"tricky/slices":      {name: "tricky/slices", src: trickySrc, tests: trickySlicesTests},
	"tricky/structs":     {name: "tricky/structs", src: trickySrc, tests: trickyStructsTests},
	"typeconv":           {name: "typeconv", src: typeconvSrc, tests: typeconvTests},
	"variables":          {name: "variables", src: variablesSrc, tests: variablesTests},
}
