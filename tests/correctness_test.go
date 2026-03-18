// Package tests - correctness_test.go
//
// Unified correctness test framework that consolidates ALL tests from testdata/
// Compares interpreted execution results with native Go execution results.
package tests

import (
	_ "embed"
	"reflect"
	"strings"
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
// Helper Functions
// ============================================================================

// toMainPackage converts a source file to package main for interpretation
func toMainPackage(src string) string {
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "package ") {
			lines[i] = "package main"
			break
		}
	}
	return strings.Join(lines, "\n")
}

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
// All Test Cases - Consolidated from all test files
// ============================================================================

var allCorrectnessTests = map[string]testCase{
	// ============================================================================
	// algorithms
	// ============================================================================
	"algorithms/InsertionSort":     {algorithmsSrc, "InsertionSort", nil, algorithms.InsertionSort},
	"algorithms/SelectionSort":     {algorithmsSrc, "SelectionSort", nil, algorithms.SelectionSort},
	"algorithms/ReverseSlice":      {algorithmsSrc, "ReverseSlice", nil, algorithms.ReverseSlice},
	"algorithms/IsPalindrome":      {algorithmsSrc, "IsPalindrome", nil, algorithms.IsPalindrome},
	"algorithms/PowerFunction":     {algorithmsSrc, "PowerFunction", nil, algorithms.PowerFunction},
	"algorithms/MaxSubarraySum":    {algorithmsSrc, "MaxSubarraySum", nil, algorithms.MaxSubarraySum},
	"algorithms/TwoSum":            {algorithmsSrc, "TwoSum", nil, algorithms.TwoSum},
	"algorithms/FibMemoized":       {algorithmsSrc, "FibMemoized", nil, algorithms.FibMemoized},
	"algorithms/CountDigits":       {algorithmsSrc, "CountDigits", nil, algorithms.CountDigits},
	"algorithms/CollatzConjecture": {algorithmsSrc, "CollatzConjecture", nil, algorithms.CollatzConjecture},
	// Parameterized tests
	"algorithms/Reverse":         {algorithmsSrc, "Reverse", []any{[]int{1, 2, 3, 4, 5}}, algorithms.Reverse},
	"algorithms/Power":           {algorithmsSrc, "Power", []any{2, 10}, algorithms.Power},
	"algorithms/CountDigitsN":    {algorithmsSrc, "CountDigitsN", []any{12345}, algorithms.CountDigitsN},
	"algorithms/CollatzStepsN":   {algorithmsSrc, "CollatzStepsN", []any{27}, algorithms.CollatzStepsN},
	"algorithms/IsPalindromeInt": {algorithmsSrc, "IsPalindromeInt", []any{[]int{1, 2, 3, 2, 1}}, algorithms.IsPalindromeInt},

	// ============================================================================
	// advanced
	// ============================================================================
	"advanced/TypeConvertIntIdentity": {advancedSrc, "TypeConvertIntIdentity", nil, advanced.TypeConvertIntIdentity},
	"advanced/DeepCallChain":          {advancedSrc, "DeepCallChain", nil, advanced.DeepCallChain},
	"advanced/EarlyReturn":            {advancedSrc, "EarlyReturn", nil, advanced.EarlyReturn},
	"advanced/NestedIfInLoop":         {advancedSrc, "NestedIfInLoop", nil, advanced.NestedIfInLoop},
	"advanced/BubbleSort":             {advancedSrc, "BubbleSort", nil, advanced.BubbleSort},
	"advanced/BinarySearch":           {advancedSrc, "BinarySearch", nil, advanced.BinarySearch},
	"advanced/GCD":                    {advancedSrc, "GCD", nil, advanced.GCD},
	"advanced/SieveOfEratosthenes":    {advancedSrc, "SieveOfEratosthenes", nil, advanced.SieveOfEratosthenes},
	"advanced/MatrixMultiply":         {advancedSrc, "MatrixMultiply", nil, advanced.MatrixMultiply},
	"advanced/EmptyFunctionReturn":    {advancedSrc, "EmptyFunctionReturn", nil, advanced.EmptyFunctionReturn},
	"advanced/SingleReturnValue":      {advancedSrc, "SingleReturnValue", nil, advanced.SingleReturnValue},
	"advanced/ZeroIteration":          {advancedSrc, "ZeroIteration", nil, advanced.ZeroIteration},
	"advanced/LargeLoop":              {advancedSrc, "LargeLoop", nil, advanced.LargeLoop},
	"advanced/DeepRecursion":          {advancedSrc, "DeepRecursion", nil, advanced.DeepRecursion},
	"advanced/MapWithClosure":         {advancedSrc, "MapWithClosure", nil, advanced.MapWithClosure},
	"advanced/SliceWithMultiReturn":   {advancedSrc, "SliceWithMultiReturn", nil, advanced.SliceWithMultiReturn},
	"advanced/RecursiveDataBuild":     {advancedSrc, "RecursiveDataBuild", nil, advanced.RecursiveDataBuild},
	"advanced/FunctionChain":          {advancedSrc, "FunctionChain", nil, advanced.FunctionChain},
	"advanced/ComplexExpressions":     {advancedSrc, "ComplexExpressions", nil, advanced.ComplexExpressions},
	// Parameterized tests
	"advanced/FindFirst": {advancedSrc, "FindFirst", []any{[]int{10, 20, 30}, 20}, advanced.FindFirst},
	"advanced/Bsearch":   {advancedSrc, "Bsearch", []any{[]int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}, 50}, advanced.Bsearch},
	"advanced/Gcd":       {advancedSrc, "Gcd", []any{48, 18}, advanced.Gcd},
	"advanced/Identity":  {advancedSrc, "Identity", []any{42}, advanced.Identity},
	"advanced/Minmax":    {advancedSrc, "Minmax", []any{[]int{3, 1, 4, 1, 5}}, advanced.Minmax},
	"advanced/Countdown": {advancedSrc, "Countdown", []any{50}, advanced.Countdown},
	"advanced/Add":       {advancedSrc, "Add", []any{1, 2}, advanced.Add},
	"advanced/Mul":       {advancedSrc, "Mul", []any{3, 4}, advanced.Mul},
	"advanced/Sub":       {advancedSrc, "Sub", []any{10, 5}, advanced.Sub},

	// ============================================================================
	// arithmetic
	// ============================================================================
	"arithmetic/Addition":       {arithmeticSrc, "Addition", nil, arithmetic.Addition},
	"arithmetic/Subtraction":    {arithmeticSrc, "Subtraction", nil, arithmetic.Subtraction},
	"arithmetic/Multiplication": {arithmeticSrc, "Multiplication", nil, arithmetic.Multiplication},
	"arithmetic/Division":       {arithmeticSrc, "Division", nil, arithmetic.Division},
	"arithmetic/Modulo":         {arithmeticSrc, "Modulo", nil, arithmetic.Modulo},
	"arithmetic/ComplexExpr":    {arithmeticSrc, "ComplexExpr", nil, arithmetic.ComplexExpr},
	"arithmetic/Negation":       {arithmeticSrc, "Negation", nil, arithmetic.Negation},
	"arithmetic/ChainedOps":     {arithmeticSrc, "ChainedOps", nil, arithmetic.ChainedOps},
	"arithmetic/Overflow":       {arithmeticSrc, "Overflow", nil, arithmetic.Overflow},
	"arithmetic/Precedence":     {arithmeticSrc, "Precedence", nil, arithmetic.Precedence},
	// Parameterized tests
	"arithmetic/Add":          {arithmeticSrc, "Add", []any{10, 32}, arithmetic.Add},
	"arithmetic/Sub":          {arithmeticSrc, "Sub", []any{100, 42}, arithmetic.Sub},
	"arithmetic/Mul":          {arithmeticSrc, "Mul", []any{6, 7}, arithmetic.Mul},
	"arithmetic/Div":          {arithmeticSrc, "Div", []any{100, 4}, arithmetic.Div},
	"arithmetic/Mod":          {arithmeticSrc, "Mod", []any{17, 5}, arithmetic.Mod},
	"arithmetic/ComplexArith": {arithmeticSrc, "ComplexArith", []any{2, 3, 4, 10, 2}, arithmetic.ComplexArith},

	// ============================================================================
	// autowrap
	// ============================================================================
	"autowrap/WithPackage": {autowrapSrc, "WithPackage", nil, autowrap.WithPackage},
	"autowrap/WithImport":  {autowrapSrc, "WithImport", nil, autowrap.WithImport},
	"autowrap/Compute":     {autowrapSrc, "Compute", nil, autowrap.Compute},

	// ============================================================================
	// bitwise
	// ============================================================================
	"bitwise/And":        {bitwiseSrc, "And", nil, bitwise.And},
	"bitwise/Or":         {bitwiseSrc, "Or", nil, bitwise.Or},
	"bitwise/Xor":        {bitwiseSrc, "Xor", nil, bitwise.Xor},
	"bitwise/LeftShift":  {bitwiseSrc, "LeftShift", nil, bitwise.LeftShift},
	"bitwise/RightShift": {bitwiseSrc, "RightShift", nil, bitwise.RightShift},
	"bitwise/Combined":   {bitwiseSrc, "Combined", nil, bitwise.Combined},
	"bitwise/AndNot":     {bitwiseSrc, "AndNot", nil, bitwise.AndNot},
	"bitwise/PowerOfTwo": {bitwiseSrc, "PowerOfTwo", nil, bitwise.PowerOfTwo},
	// Parameterized tests
	"bitwise/BitAnd":        {bitwiseSrc, "BitAnd", []any{0xFF, 0x0F}, bitwise.BitAnd},
	"bitwise/BitOr":         {bitwiseSrc, "BitOr", []any{0xFF, 0x100}, bitwise.BitOr},
	"bitwise/BitXor":        {bitwiseSrc, "BitXor", []any{0xAA, 0x55}, bitwise.BitXor},
	"bitwise/BitLeftShift":  {bitwiseSrc, "BitLeftShift", []any{1, 10}, bitwise.BitLeftShift},
	"bitwise/BitRightShift": {bitwiseSrc, "BitRightShift", []any{1024, 5}, bitwise.BitRightShift},
	"bitwise/IsPowerOfTwo":  {bitwiseSrc, "IsPowerOfTwo", []any{16}, bitwise.IsPowerOfTwo},

	// ============================================================================
	// closures
	// ============================================================================
	"closures/Counter":           {closuresSrc, "Counter", nil, closures.Counter},
	"closures/CaptureMutation":   {closuresSrc, "CaptureMutation", nil, closures.CaptureMutation},
	"closures/Factory":           {closuresSrc, "Factory", nil, closures.Factory},
	"closures/MultipleInstances": {closuresSrc, "MultipleInstances", nil, closures.MultipleInstances},
	"closures/OverLoop":          {closuresSrc, "OverLoop", nil, closures.OverLoop},
	"closures/Chain":             {closuresSrc, "Chain", nil, closures.Chain},
	"closures/Accumulator":       {closuresSrc, "Accumulator", nil, closures.Accumulator},
	// Parameterized tests (closures return functions, skip for now as we test via Call)

	// ============================================================================
	// closures_advanced
	// ============================================================================
	"closures_advanced/Generator":          {closuresAdvancedSrc, "Generator", nil, closures_advanced.Generator},
	"closures_advanced/Predicate":          {closuresAdvancedSrc, "Predicate", nil, closures_advanced.Predicate},
	"closures_advanced/StateMachine":       {closuresAdvancedSrc, "StateMachine", nil, closures_advanced.StateMachine},
	"closures_advanced/RecursiveHelper":    {closuresAdvancedSrc, "RecursiveHelper", nil, closures_advanced.RecursiveHelper},
	"closures_advanced/ApplyN":             {closuresAdvancedSrc, "ApplyN", nil, closures_advanced.ApplyN},
	"closures_advanced/Compose":            {closuresAdvancedSrc, "Compose", nil, closures_advanced.Compose},
	"closures_advanced/ClosureForLoopTest": {closuresAdvancedSrc, "ClosureForLoopTest", nil, closures_advanced.ClosureForLoopTest},

	// ============================================================================
	// controlflow
	// ============================================================================
	"controlflow/IfTrue":              {controlflowSrc, "IfTrue", nil, controlflow.IfTrue},
	"controlflow/IfFalse":             {controlflowSrc, "IfFalse", nil, controlflow.IfFalse},
	"controlflow/IfElse":              {controlflowSrc, "IfElse", nil, controlflow.IfElse},
	"controlflow/IfElseChainNegative": {controlflowSrc, "IfElseChainNegative", nil, controlflow.IfElseChainNegative},
	"controlflow/IfElseChainZero":     {controlflowSrc, "IfElseChainZero", nil, controlflow.IfElseChainZero},
	"controlflow/IfElseChainPositive": {controlflowSrc, "IfElseChainPositive", nil, controlflow.IfElseChainPositive},
	"controlflow/ForLoop":             {controlflowSrc, "ForLoop", nil, controlflow.ForLoop},
	"controlflow/ForConditionOnly":    {controlflowSrc, "ForConditionOnly", nil, controlflow.ForConditionOnly},
	"controlflow/NestedFor":           {controlflowSrc, "NestedFor", nil, controlflow.NestedFor},
	"controlflow/ForBreak":            {controlflowSrc, "ForBreak", nil, controlflow.ForBreak},
	"controlflow/ForContinue":         {controlflowSrc, "ForContinue", nil, controlflow.ForContinue},
	"controlflow/BooleanAndOr":        {controlflowSrc, "BooleanAndOr", nil, controlflow.BooleanAndOr},
	// Parameterized tests
	"controlflow/ClassifyNegative": {controlflowSrc, "Classify", []any{-5}, controlflow.Classify},
	"controlflow/ClassifyZero":     {controlflowSrc, "Classify", []any{0}, controlflow.Classify},
	"controlflow/ClassifyPositive": {controlflowSrc, "Classify", []any{42}, controlflow.Classify},

	// ============================================================================
	// cornercases
	// ============================================================================
	"cornercases/ZeroValue_Int":          {cornercasesSrc, "ZeroValue_Int", nil, cornercases.ZeroValue_Int},
	"cornercases/ZeroValue_Int64":        {cornercasesSrc, "ZeroValue_Int64", nil, cornercases.ZeroValue_Int64},
	"cornercases/ZeroValue_Float64":      {cornercasesSrc, "ZeroValue_Float64", nil, cornercases.ZeroValue_Float64},
	"cornercases/ZeroValue_String":       {cornercasesSrc, "ZeroValue_String", nil, cornercases.ZeroValue_String},
	"cornercases/ZeroValue_Bool":         {cornercasesSrc, "ZeroValue_Bool", nil, cornercases.ZeroValue_Bool},
	"cornercases/ZeroValue_Slice":        {cornercasesSrc, "ZeroValue_Slice", nil, cornercases.ZeroValue_Slice},
	"cornercases/ZeroValue_Map":          {cornercasesSrc, "ZeroValue_Map", nil, cornercases.ZeroValue_Map},
	"cornercases/IntBoundary_MaxInt32":   {cornercasesSrc, "IntBoundary_MaxInt32", nil, cornercases.IntBoundary_MaxInt32},
	"cornercases/IntBoundary_MinInt32":   {cornercasesSrc, "IntBoundary_MinInt32", nil, cornercases.IntBoundary_MinInt32},
	"cornercases/IntBoundary_MaxInt64":   {cornercasesSrc, "IntBoundary_MaxInt64", nil, cornercases.IntBoundary_MaxInt64},
	"cornercases/IntBoundary_MinInt64":   {cornercasesSrc, "IntBoundary_MinInt64", nil, cornercases.IntBoundary_MinInt64},
	"cornercases/IntBoundary_MaxUint32":  {cornercasesSrc, "IntBoundary_MaxUint32", nil, cornercases.IntBoundary_MaxUint32},
	"cornercases/IntBoundary_NearMaxInt": {cornercasesSrc, "IntBoundary_NearMaxInt", nil, cornercases.IntBoundary_NearMaxInt},
	"cornercases/IntBoundary_NearMinInt": {cornercasesSrc, "IntBoundary_NearMinInt", nil, cornercases.IntBoundary_NearMinInt},
	// Note: int32 overflow wraps just like native Go (two's complement)
	// We compare the actual native results to verify correctness
	"cornercases/Overflow_Int32Add":           {cornercasesSrc, "Overflow_Int32Add", nil, cornercases.Overflow_Int32Add},
	"cornercases/Overflow_Int32Sub":           {cornercasesSrc, "Overflow_Int32Sub", nil, cornercases.Overflow_Int32Sub},
	"cornercases/Overflow_Int32Mul":           {cornercasesSrc, "Overflow_Int32Mul", nil, cornercases.Overflow_Int32Mul},
	"cornercases/FloatBoundary_SmallPositive": {cornercasesSrc, "FloatBoundary_SmallPositive", nil, cornercases.FloatBoundary_SmallPositive},
	"cornercases/FloatBoundary_SmallNegative": {cornercasesSrc, "FloatBoundary_SmallNegative", nil, cornercases.FloatBoundary_SmallNegative},
	"cornercases/FloatBoundary_LargePositive": {cornercasesSrc, "FloatBoundary_LargePositive", nil, cornercases.FloatBoundary_LargePositive},
	"cornercases/FloatBoundary_LargeNegative": {cornercasesSrc, "FloatBoundary_LargeNegative", nil, cornercases.FloatBoundary_LargeNegative},
	"cornercases/EmptySlice_Len":              {cornercasesSrc, "EmptySlice_Len", nil, cornercases.EmptySlice_Len},
	"cornercases/EmptySlice_Cap":              {cornercasesSrc, "EmptySlice_Cap", nil, cornercases.EmptySlice_Cap},
	"cornercases/EmptySlice_Make":             {cornercasesSrc, "EmptySlice_Make", nil, cornercases.EmptySlice_Make},
	"cornercases/EmptyMap_Len":                {cornercasesSrc, "EmptyMap_Len", nil, cornercases.EmptyMap_Len},
	"cornercases/EmptyMap_Make":               {cornercasesSrc, "EmptyMap_Make", nil, cornercases.EmptyMap_Make},
	"cornercases/EmptyString_Len":             {cornercasesSrc, "EmptyString_Len", nil, cornercases.EmptyString_Len},
	"cornercases/Slice_ZeroToZero":            {cornercasesSrc, "Slice_ZeroToZero", nil, cornercases.Slice_ZeroToZero},
	"cornercases/Slice_EndToEnd":              {cornercasesSrc, "Slice_EndToEnd", nil, cornercases.Slice_EndToEnd},
	"cornercases/Slice_NilSlice":              {cornercasesSrc, "Slice_NilSlice", nil, cornercases.Slice_NilSlice},
	"cornercases/Slice_AppendToNil":           {cornercasesSrc, "Slice_AppendToNil", nil, cornercases.Slice_AppendToNil},
	"cornercases/Slice_AppendEmpty":           {cornercasesSrc, "Slice_AppendEmpty", nil, cornercases.Slice_AppendEmpty},
	"cornercases/Map_NilMap":                  {cornercasesSrc, "Map_NilMap", nil, cornercases.Map_NilMap},
	"cornercases/Map_AccessMissingKey":        {cornercasesSrc, "Map_AccessMissingKey", nil, cornercases.Map_AccessMissingKey},
	"cornercases/Map_DeleteMissingKey":        {cornercasesSrc, "Map_DeleteMissingKey", nil, cornercases.Map_DeleteMissingKey},
	"cornercases/Map_OverwriteKey":            {cornercasesSrc, "Map_OverwriteKey", nil, cornercases.Map_OverwriteKey},
	"cornercases/Map_NilKeyString":            {cornercasesSrc, "Map_NilKeyString", nil, cornercases.Map_NilKeyString},
	"cornercases/Map_ZeroIntKey":              {cornercasesSrc, "Map_ZeroIntKey", nil, cornercases.Map_ZeroIntKey},
	"cornercases/String_Empty":                {cornercasesSrc, "String_Empty", nil, cornercases.String_Empty},
	"cornercases/String_SingleChar":           {cornercasesSrc, "String_SingleChar", nil, cornercases.String_SingleChar},
	"cornercases/String_UnicodeMultibyte":     {cornercasesSrc, "String_UnicodeMultibyte", nil, cornercases.String_UnicodeMultibyte},
	"cornercases/String_Whitespace":           {cornercasesSrc, "String_Whitespace", nil, cornercases.String_Whitespace},
	"cornercases/String_SingleByteIndex":      {cornercasesSrc, "String_SingleByteIndex", nil, cornercases.String_SingleByteIndex},
	"cornercases/String_LastByte":             {cornercasesSrc, "String_LastByte", nil, cornercases.String_LastByte},
	"cornercases/Bool_True":                   {cornercasesSrc, "Bool_True", nil, cornercases.Bool_True},
	"cornercases/Bool_False":                  {cornercasesSrc, "Bool_False", nil, cornercases.Bool_False},
	"cornercases/Bool_NotTrue":                {cornercasesSrc, "Bool_NotTrue", nil, cornercases.Bool_NotTrue},
	"cornercases/Bool_NotFalse":               {cornercasesSrc, "Bool_NotFalse", nil, cornercases.Bool_NotFalse},
	"cornercases/Bool_DoubleNegation":         {cornercasesSrc, "Bool_DoubleNegation", nil, cornercases.Bool_DoubleNegation},
	"cornercases/Arith_AddZero":               {cornercasesSrc, "Arith_AddZero", nil, cornercases.Arith_AddZero},
	"cornercases/Arith_SubZero":               {cornercasesSrc, "Arith_SubZero", nil, cornercases.Arith_SubZero},
	"cornercases/Arith_MulByOne":              {cornercasesSrc, "Arith_MulByOne", nil, cornercases.Arith_MulByOne},
	"cornercases/Arith_DivByOne":              {cornercasesSrc, "Arith_DivByOne", nil, cornercases.Arith_DivByOne},
	"cornercases/Arith_ModByOne":              {cornercasesSrc, "Arith_ModByOne", nil, cornercases.Arith_ModByOne},
	"cornercases/Arith_MulByZero":             {cornercasesSrc, "Arith_MulByZero", nil, cornercases.Arith_MulByZero},
	"cornercases/Arith_NegNeg":                {cornercasesSrc, "Arith_NegNeg", nil, cornercases.Arith_NegNeg},
	"cornercases/Arith_NegAddNeg":             {cornercasesSrc, "Arith_NegAddNeg", nil, cornercases.Arith_NegAddNeg},
	"cornercases/Compare_IntEqual":            {cornercasesSrc, "Compare_IntEqual", nil, cornercases.Compare_IntEqual},
	"cornercases/Compare_IntNotEqual":         {cornercasesSrc, "Compare_IntNotEqual", nil, cornercases.Compare_IntNotEqual},
	"cornercases/Compare_IntGreater":          {cornercasesSrc, "Compare_IntGreater", nil, cornercases.Compare_IntGreater},
	"cornercases/Compare_IntGreaterEqual":     {cornercasesSrc, "Compare_IntGreaterEqual", nil, cornercases.Compare_IntGreaterEqual},
	"cornercases/Compare_IntLess":             {cornercasesSrc, "Compare_IntLess", nil, cornercases.Compare_IntLess},
	"cornercases/Compare_IntLessEqual":        {cornercasesSrc, "Compare_IntLessEqual", nil, cornercases.Compare_IntLessEqual},
	"cornercases/Compare_StringEqual":         {cornercasesSrc, "Compare_StringEqual", nil, cornercases.Compare_StringEqual},
	"cornercases/Compare_StringNotEqual":      {cornercasesSrc, "Compare_StringNotEqual", nil, cornercases.Compare_StringNotEqual},
	"cornercases/Compare_EmptyStringEqual":    {cornercasesSrc, "Compare_EmptyStringEqual", nil, cornercases.Compare_EmptyStringEqual},
	"cornercases/Logic_TrueAndTrue":           {cornercasesSrc, "Logic_TrueAndTrue", nil, cornercases.Logic_TrueAndTrue},
	"cornercases/Logic_TrueAndFalse":          {cornercasesSrc, "Logic_TrueAndFalse", nil, cornercases.Logic_TrueAndFalse},
	"cornercases/Logic_FalseAndTrue":          {cornercasesSrc, "Logic_FalseAndTrue", nil, cornercases.Logic_FalseAndTrue},
	"cornercases/Logic_TrueOrFalse":           {cornercasesSrc, "Logic_TrueOrFalse", nil, cornercases.Logic_TrueOrFalse},
	"cornercases/Logic_FalseOrTrue":           {cornercasesSrc, "Logic_FalseOrTrue", nil, cornercases.Logic_FalseOrTrue},
	"cornercases/Logic_FalseOrFalse":          {cornercasesSrc, "Logic_FalseOrFalse", nil, cornercases.Logic_FalseOrFalse},
	"cornercases/Control_IfNoElse":            {cornercasesSrc, "Control_IfNoElse", nil, cornercases.Control_IfNoElse},
	"cornercases/Control_IfFalseNoElse":       {cornercasesSrc, "Control_IfFalseNoElse", nil, cornercases.Control_IfFalseNoElse},
	"cornercases/Control_ForZeroIter":         {cornercasesSrc, "Control_ForZeroIter", nil, cornercases.Control_ForZeroIter},
	"cornercases/Control_ForOneIter":          {cornercasesSrc, "Control_ForOneIter", nil, cornercases.Control_ForOneIter},
	"cornercases/Control_ForBreakFirst":       {cornercasesSrc, "Control_ForBreakFirst", nil, cornercases.Control_ForBreakFirst},
	"cornercases/Control_ForContinueAll":      {cornercasesSrc, "Control_ForContinueAll", nil, cornercases.Control_ForContinueAll},
	"cornercases/Control_SwitchNoMatch":       {cornercasesSrc, "Control_SwitchNoMatch", nil, cornercases.Control_SwitchNoMatch},
	"cornercases/Control_SwitchDefault":       {cornercasesSrc, "Control_SwitchDefault", nil, cornercases.Control_SwitchDefault},
	"cornercases/Func_NoReturn":               {cornercasesSrc, "Func_NoReturn", nil, cornercases.Func_NoReturn},
	"cornercases/Func_MultipleReturnAll":      {cornercasesSrc, "Func_MultipleReturnAll", nil, cornercases.Func_MultipleReturnAll},
	"cornercases/Func_MultipleReturnIgnore":   {cornercasesSrc, "Func_MultipleReturnIgnore", nil, cornercases.Func_MultipleReturnIgnore},
	"cornercases/Func_NamedReturn":            {cornercasesSrc, "Func_NamedReturn", nil, cornercases.Func_NamedReturn},
	"cornercases/Func_VariadicEmpty":          {cornercasesSrc, "Func_VariadicEmpty", nil, cornercases.Func_VariadicEmpty},
	"cornercases/Func_VariadicOne":            {cornercasesSrc, "Func_VariadicOne", nil, cornercases.Func_VariadicOne},
	"cornercases/Func_VariadicMultiple":       {cornercasesSrc, "Func_VariadicMultiple", nil, cornercases.Func_VariadicMultiple},
	"cornercases/Func_RecursionBase":          {cornercasesSrc, "Func_RecursionBase", nil, cornercases.Func_RecursionBase},
	"cornercases/Closure_ReturnClosure":       {cornercasesSrc, "Closure_ReturnClosure", nil, cornercases.Closure_ReturnClosure},
	"cornercases/Closure_CaptureVariable":     {cornercasesSrc, "Closure_CaptureVariable", nil, cornercases.Closure_CaptureVariable},
	"cornercases/Closure_ModifyCaptured":      {cornercasesSrc, "Closure_ModifyCaptured", nil, cornercases.Closure_ModifyCaptured},
	"cornercases/Struct_ZeroValueFields":      {cornercasesSrc, "Struct_ZeroValueFields", nil, cornercases.Struct_ZeroValueFields},
	"cornercases/Struct_PointerReceiver":      {cornercasesSrc, "Struct_PointerReceiver", nil, cornercases.Struct_PointerReceiver},
	"cornercases/Struct_NestedStruct":         {cornercasesSrc, "Struct_NestedStruct", nil, cornercases.Struct_NestedStruct},

	// ============================================================================
	// edgecases
	// ============================================================================
	"edgecases/MaxInt64":           {edgecasesSrc, "MaxInt64", nil, edgecases.MaxInt64},
	"edgecases/MinInt64":           {edgecasesSrc, "MinInt64", nil, edgecases.MinInt64},
	"edgecases/DivisionByMinusOne": {edgecasesSrc, "DivisionByMinusOne", nil, edgecases.DivisionByMinusOne},
	"edgecases/ModuloNegative":     {edgecasesSrc, "ModuloNegative", nil, edgecases.ModuloNegative},
	"edgecases/EmptyString":        {edgecasesSrc, "EmptyString", nil, edgecases.EmptyString},
	"edgecases/LargeSlice":         {edgecasesSrc, "LargeSlice", nil, edgecases.LargeSlice},
	"edgecases/NestedMapLookup":    {edgecasesSrc, "NestedMapLookup", nil, edgecases.NestedMapLookup},
	"edgecases/ZeroDivisionGuard":  {edgecasesSrc, "ZeroDivisionGuard", nil, edgecases.ZeroDivisionGuard},
	"edgecases/BooleanComplexExpr": {edgecasesSrc, "BooleanComplexExpr", nil, edgecases.BooleanComplexExpr},
	"edgecases/SingleElementSlice": {edgecasesSrc, "SingleElementSlice", nil, edgecases.SingleElementSlice},
	"edgecases/EmptyMap":           {edgecasesSrc, "EmptyMap", nil, edgecases.EmptyMap},
	"edgecases/TightLoop":          {edgecasesSrc, "TightLoop", nil, edgecases.TightLoop},

	// ============================================================================
	// external
	// ============================================================================
	"external/FmtSprintf":       {externalSrc, "FmtSprintf", nil, external.FmtSprintf},
	"external/FmtSprintfMulti":  {externalSrc, "FmtSprintfMulti", nil, external.FmtSprintfMulti},
	"external/StringsToUpper":   {externalSrc, "StringsToUpper", nil, external.StringsToUpper},
	"external/StringsToLower":   {externalSrc, "StringsToLower", nil, external.StringsToLower},
	"external/StringsContains":  {externalSrc, "StringsContains", nil, external.StringsContains},
	"external/StringsReplace":   {externalSrc, "StringsReplace", nil, external.StringsReplace},
	"external/StringsHasPrefix": {externalSrc, "StringsHasPrefix", nil, external.StringsHasPrefix},
	"external/StrconvItoa":      {externalSrc, "StrconvItoa", nil, external.StrconvItoa},
	"external/StrconvAtoi":      {externalSrc, "StrconvAtoi", nil, external.StrconvAtoi},
	// Parameterized tests
	"external/FmtSprintfInt":      {externalSrc, "FmtSprintfInt", []any{42}, external.FmtSprintfInt},
	"external/StringsToUpperStr":  {externalSrc, "StringsToUpperStr", []any{"hello"}, external.StringsToUpperStr},
	"external/StringsToLowerStr":  {externalSrc, "StringsToLowerStr", []any{"HELLO"}, external.StringsToLowerStr},
	"external/StringsContainsStr": {externalSrc, "StringsContainsStr", []any{"hello world", "world"}, external.StringsContainsStr},
	"external/StrconvItoaN":       {externalSrc, "StrconvItoaN", []any{42}, external.StrconvItoaN},
	"external/StrconvAtoiStr":     {externalSrc, "StrconvAtoiStr", []any{"123"}, external.StrconvAtoiStr},

	// ============================================================================
	// functions
	// ============================================================================
	"functions/Call":                 {functionsSrc, "Call", nil, functions.Call},
	"functions/MultipleReturn":       {functionsSrc, "MultipleReturn", nil, functions.MultipleReturn},
	"functions/MultipleReturnDivmod": {functionsSrc, "MultipleReturnDivmod", nil, functions.MultipleReturnDivmod},
	"functions/RecursionFactorial":   {functionsSrc, "RecursionFactorial", nil, functions.RecursionFactorial},
	"functions/MutualRecursion":      {functionsSrc, "MutualRecursion", nil, functions.MutualRecursion},
	"functions/FibonacciIterative":   {functionsSrc, "FibonacciIterative", nil, functions.FibonacciIterative},
	"functions/FibonacciRecursive":   {functionsSrc, "FibonacciRecursive", nil, functions.FibonacciRecursive},
	"functions/VariadicFunction":     {functionsSrc, "VariadicFunction", nil, functions.VariadicFunction},
	"functions/FunctionAsValue":      {functionsSrc, "FunctionAsValue", nil, functions.FunctionAsValue},
	"functions/HigherOrderMap":       {functionsSrc, "HigherOrderMap", nil, functions.HigherOrderMap},
	"functions/HigherOrderFilter":    {functionsSrc, "HigherOrderFilter", nil, functions.HigherOrderFilter},
	"functions/HigherOrderReduce":    {functionsSrc, "HigherOrderReduce", nil, functions.HigherOrderReduce},
	// Parameterized tests
	"functions/Add":        {functionsSrc, "Add", []any{5, 7}, functions.Add},
	"functions/Swap":       {functionsSrc, "Swap", []any{3, 7}, functions.Swap},
	"functions/Divmod":     {functionsSrc, "Divmod", []any{17, 5}, functions.Divmod},
	"functions/FactorialN": {functionsSrc, "FactorialN", []any{5}, functions.FactorialN},
	"functions/FibIterN":   {functionsSrc, "FibIterN", []any{20}, functions.FibIterN},
	"functions/FibRecN":    {functionsSrc, "FibRecN", []any{15}, functions.FibRecN},
	"functions/IsEvenN":    {functionsSrc, "IsEvenN", []any{10}, functions.IsEvenN},
	"functions/IsOddN":     {functionsSrc, "IsOddN", []any{7}, functions.IsOddN},
	// Variadic with args - skip for now (interpreter variadic handling from outside needs work)

	// Multi-return value tests - these functions THEMSELVES return multiple values
	"functions/ThreeReturnValues":               {functionsSrc, "ThreeReturnValues", nil, functions.ThreeReturnValues},
	"functions/FourReturnValues":                {functionsSrc, "FourReturnValues", nil, functions.FourReturnValues},
	"functions/FiveReturnValues":                {functionsSrc, "FiveReturnValues", nil, functions.FiveReturnValues},
	"functions/MixedTypeReturn":                 {functionsSrc, "MixedTypeReturn", nil, functions.MixedTypeReturn},
	"functions/PassMultiReturnToFunc":           {functionsSrc, "PassMultiReturnToFunc", nil, functions.PassMultiReturnToFunc},
	"functions/ChainMultiReturn":                {functionsSrc, "ChainMultiReturn", nil, functions.ChainMultiReturn},
	"functions/NestedMultiReturn":               {functionsSrc, "NestedMultiReturn", nil, functions.NestedMultiReturn},
	"functions/MultiReturnAsSliceIndex":         {functionsSrc, "MultiReturnAsSliceIndex", nil, functions.MultiReturnAsSliceIndex},
	"functions/MultiReturnToMap":                {functionsSrc, "MultiReturnToMap", nil, functions.MultiReturnToMap},
	"functions/MultiReturnAsCondition":          {functionsSrc, "MultiReturnAsCondition", nil, functions.MultiReturnAsCondition},
	"functions/MultiReturnComplexTypes":         {functionsSrc, "MultiReturnComplexTypes", nil, functions.MultiReturnComplexTypes},
	"functions/MultiReturnInClosure":            {functionsSrc, "MultiReturnInClosure", nil, functions.MultiReturnInClosure},
	"functions/AssignMultiReturnToExistingVars": {functionsSrc, "AssignMultiReturnToExistingVars", nil, functions.AssignMultiReturnToExistingVars},

	// ============================================================================
	// initialize - Complex initialization tests
	// ============================================================================
	"initialize/ComplexInitTest":     {initializeSrc, "ComplexInitTest", nil, initialize.ComplexInitTest},
	"initialize/InitOrderTest":       {initializeSrc, "InitOrderTest", nil, initialize.InitOrderTest},
	"initialize/CacheInitTest":       {initializeSrc, "CacheInitTest", nil, initialize.CacheInitTest},
	"initialize/LookupTableInitTest": {initializeSrc, "LookupTableInitTest", nil, initialize.LookupTableInitTest},
	"initialize/FibonacciInitTest":   {initializeSrc, "FibonacciInitTest", nil, initialize.FibonacciInitTest},
	"initialize/GetA":                {initializeSrc, "GetA", nil, initialize.GetA},
	"initialize/GetB":                {initializeSrc, "GetB", nil, initialize.GetB},
	"initialize/GetC":                {initializeSrc, "GetC", nil, initialize.GetC},
	"initialize/GetCacheSum":         {initializeSrc, "GetCacheSum", nil, initialize.GetCacheSum},
	"initialize/GetCacheSize":        {initializeSrc, "GetCacheSize", nil, initialize.GetCacheSize},
	"initialize/GetFibonacciCount":   {initializeSrc, "GetFibonacciCount", nil, initialize.GetFibonacciCount},

	// ============================================================================
	// leetcode_hard
	// ============================================================================
	"leetcode_hard/TrappingRainWater":           {leetcodeHardSrc, "TrappingRainWater", nil, leetcode_hard.TrappingRainWater},
	"leetcode_hard/LargestRectangleInHistogram": {leetcodeHardSrc, "LargestRectangleInHistogram", nil, leetcode_hard.LargestRectangleInHistogram},
	"leetcode_hard/MedianOfTwoSortedArrays":     {leetcodeHardSrc, "MedianOfTwoSortedArrays", nil, leetcode_hard.MedianOfTwoSortedArrays},
	"leetcode_hard/RegularExpressionMatching":   {leetcodeHardSrc, "RegularExpressionMatching", nil, leetcode_hard.RegularExpressionMatching},
	"leetcode_hard/NQueens":                     {leetcodeHardSrc, "NQueens", nil, leetcode_hard.NQueens},
	"leetcode_hard/LongestIncreasingPath":       {leetcodeHardSrc, "LongestIncreasingPath", nil, leetcode_hard.LongestIncreasingPath},
	"leetcode_hard/WordLadder":                  {leetcodeHardSrc, "WordLadder", nil, leetcode_hard.WordLadder},
	"leetcode_hard/MergeKSortedLists":           {leetcodeHardSrc, "MergeKSortedLists", nil, leetcode_hard.MergeKSortedLists},
	"leetcode_hard/EditDistance":                {leetcodeHardSrc, "EditDistance", nil, leetcode_hard.EditDistance},
	"leetcode_hard/MinimumWindowSubstring":      {leetcodeHardSrc, "MinimumWindowSubstring", nil, leetcode_hard.MinimumWindowSubstring},

	// ============================================================================
	// maps
	// ============================================================================
	"maps/BasicOps":       {mapsSrc, "BasicOps", nil, maps.BasicOps},
	"maps/Iteration":      {mapsSrc, "Iteration", nil, maps.Iteration},
	"maps/Delete":         {mapsSrc, "Delete", nil, maps.Delete},
	"maps/Len":            {mapsSrc, "Len", nil, maps.Len},
	"maps/Overwrite":      {mapsSrc, "Overwrite", nil, maps.Overwrite},
	"maps/IntKeys":        {mapsSrc, "IntKeys", nil, maps.IntKeys},
	"maps/PassToFunction": {mapsSrc, "PassToFunction", nil, maps.PassToFunction},
	// Parameterized tests
	"maps/SumValues": {mapsSrc, "SumValues", []any{map[string]int{"a": 100, "b": 200}}, maps.SumValues},

	// ============================================================================
	// mapadvanced
	// ============================================================================
	"mapadvanced/LookupExistingKey": {mapadvancedSrc, "LookupExistingKey", nil, mapadvanced.LookupExistingKey},
	"mapadvanced/LookupWithDefault": {mapadvancedSrc, "LookupWithDefault", nil, mapadvanced.LookupWithDefault},
	"mapadvanced/AsCounter":         {mapadvancedSrc, "AsCounter", nil, mapadvanced.AsCounter},
	"mapadvanced/WithStringValues":  {mapadvancedSrc, "WithStringValues", nil, mapadvanced.WithStringValues},
	"mapadvanced/BuildFromLoop":     {mapadvancedSrc, "BuildFromLoop", nil, mapadvanced.BuildFromLoop},
	"mapadvanced/DeleteAndReinsert": {mapadvancedSrc, "DeleteAndReinsert", nil, mapadvanced.DeleteAndReinsert},

	// ============================================================================
	// multiassign
	// ============================================================================
	"multiassign/Swap":             {multiassignSrc, "Swap", nil, multiassign.Swap},
	"multiassign/FromFunction":     {multiassignSrc, "FromFunction", nil, multiassign.FromFunction},
	"multiassign/ThreeValues":      {multiassignSrc, "ThreeValues", nil, multiassign.ThreeValues},
	"multiassign/InLoop":           {multiassignSrc, "InLoop", nil, multiassign.InLoop},
	"multiassign/DiscardWithBlank": {multiassignSrc, "DiscardWithBlank", nil, multiassign.DiscardWithBlank},
	// Parameterized tests
	"multiassign/TwoVals":    {multiassignSrc, "TwoVals", nil, multiassign.TwoVals},
	"multiassign/ThreeValsN": {multiassignSrc, "ThreeValsN", []any{10}, multiassign.ThreeValsN},
	"multiassign/DivmodAB":   {multiassignSrc, "DivmodAB", []any{17, 5}, multiassign.DivmodAB},

	// ============================================================================
	// namedreturn
	// ============================================================================
	"namedreturn/Basic":     {namedreturnSrc, "Basic", nil, namedreturn.Basic},
	"namedreturn/Multiple":  {namedreturnSrc, "Multiple", nil, namedreturn.Multiple},
	"namedreturn/ZeroValue": {namedreturnSrc, "ZeroValue", nil, namedreturn.ZeroValue},
	"namedreturn/DivMod":    {namedreturnSrc, "Divmod", []any{1000, 7}, namedreturn.Divmod},

	// ============================================================================
	// recursion
	// ============================================================================
	"recursion/TailRecursionPattern": {recursionSrc, "TailRecursionPattern", nil, recursion.TailRecursionPattern},
	"recursion/ReverseSlice":         {recursionSrc, "ReverseSlice", nil, recursion.ReverseSlice},
	"recursion/TowerOfHanoi":         {recursionSrc, "TowerOfHanoi", nil, recursion.TowerOfHanoi},
	"recursion/MaxSlice":             {recursionSrc, "MaxSlice", nil, recursion.MaxSlice},
	"recursion/Ackermann":            {recursionSrc, "Ackermann", nil, recursion.Ackermann},
	"recursion/BinarySearch":         {recursionSrc, "BinarySearch", nil, recursion.BinarySearch},
	// Parameterized tests
	"recursion/SumTail": {recursionSrc, "SumTail", []any{50, 0}, recursion.SumTail},
	"recursion/HanoiN":  {recursionSrc, "HanoiN", []any{10}, recursion.HanoiN},
	"recursion/Ack":     {recursionSrc, "Ack", []any{2, 3}, recursion.Ack},
	"recursion/MaxVal":  {recursionSrc, "MaxVal", []any{[]int{3, 7, 1, 9, 4}, 5}, recursion.MaxVal},

	// ============================================================================
	// scope
	// ============================================================================
	"scope/IfInitShortVar":            {scopeSrc, "IfInitShortVar", nil, scope.IfInitShortVar},
	"scope/IfInitMultiCondition":      {scopeSrc, "IfInitMultiCondition", nil, scope.IfInitMultiCondition},
	"scope/NestedScopes":              {scopeSrc, "NestedScopes", nil, scope.NestedScopes},
	"scope/ForScopeIsolation":         {scopeSrc, "ForScopeIsolation", nil, scope.ForScopeIsolation},
	"scope/MultipleBlockScopes":       {scopeSrc, "MultipleBlockScopes", nil, scope.MultipleBlockScopes},
	"scope/ClosureCapturesOuterScope": {scopeSrc, "ClosureCapturesOuterScope", nil, scope.ClosureCapturesOuterScope},
	// Parameterized tests
	"scope/Abs": {scopeSrc, "Abs", []any{-42}, scope.Abs},

	// ============================================================================
	// slices
	// ============================================================================
	"slices/MakeLen":           {slicesSrc, "MakeLen", nil, slices.MakeLen},
	"slices/Append":            {slicesSrc, "Append", nil, slices.Append},
	"slices/ElementAssignment": {slicesSrc, "ElementAssignment", nil, slices.ElementAssignment},
	"slices/ForRange":          {slicesSrc, "ForRange", nil, slices.ForRange},
	"slices/ForRangeIndex":     {slicesSrc, "ForRangeIndex", nil, slices.ForRangeIndex},
	"slices/GrowMultiple":      {slicesSrc, "GrowMultiple", nil, slices.GrowMultiple},
	"slices/PassToFunction":    {slicesSrc, "PassToFunction", nil, slices.PassToFunction},
	"slices/LenCap":            {slicesSrc, "LenCap", nil, slices.LenCap},
	// Parameterized tests
	"slices/SumSlice": {slicesSrc, "SumSlice", []any{[]int{1, 2, 3, 4, 5}}, slices.SumSlice},

	// ============================================================================
	// slicing
	// ============================================================================
	"slicing/SubSliceBasic":            {slicingSrc, "SubSliceBasic", nil, slicing.SubSliceBasic},
	"slicing/SubSliceLen":              {slicingSrc, "SubSliceLen", nil, slicing.SubSliceLen},
	"slicing/SubSliceFromStart":        {slicingSrc, "SubSliceFromStart", nil, slicing.SubSliceFromStart},
	"slicing/SubSliceToEnd":            {slicingSrc, "SubSliceToEnd", nil, slicing.SubSliceToEnd},
	"slicing/SubSliceCopy":             {slicingSrc, "SubSliceCopy", nil, slicing.SubSliceCopy},
	"slicing/SubSliceChained":          {slicingSrc, "SubSliceChained", nil, slicing.SubSliceChained},
	"slicing/SubSliceModifiesOriginal": {slicingSrc, "SubSliceModifiesOriginal", nil, slicing.SubSliceModifiesOriginal},
	// Parameterized tests
	"slicing/SliceLen":          {slicingSrc, "SliceLen", []any{[]int{10, 20, 30, 40, 50, 60, 70}, 2, 5}, slicing.SliceLen},
	"slicing/SliceSumRange":     {slicingSrc, "SliceSumRange", []any{[]int{10, 20, 30, 40, 50}, 1, 4}, slicing.SliceSumRange},
	"slicing/SliceFirstElement": {slicingSrc, "SliceFirstElement", []any{[]int{100, 200, 300}, 0}, slicing.SliceFirstElement},

	// ============================================================================
	// strings_pkg
	// ============================================================================
	"strings_pkg/Concat":     {stringsPkgSrc, "Concat", nil, strings_pkg.Concat},
	"strings_pkg/ConcatLoop": {stringsPkgSrc, "ConcatLoop", nil, strings_pkg.ConcatLoop},
	"strings_pkg/Len":        {stringsPkgSrc, "Len", nil, strings_pkg.Len},
	"strings_pkg/Index":      {stringsPkgSrc, "Index", nil, strings_pkg.Index},
	"strings_pkg/Comparison": {stringsPkgSrc, "Comparison", nil, strings_pkg.Comparison},
	"strings_pkg/Equality":   {stringsPkgSrc, "Equality", nil, strings_pkg.Equality},
	"strings_pkg/EmptyCheck": {stringsPkgSrc, "EmptyCheck", nil, strings_pkg.EmptyCheck},
	// Parameterized tests
	"strings_pkg/StrConcat":  {stringsPkgSrc, "StrConcat", []any{"hello", " world"}, strings_pkg.StrConcat},
	"strings_pkg/StrLen":     {stringsPkgSrc, "StrLen", []any{"hello"}, strings_pkg.StrLen},
	"strings_pkg/StrCompare": {stringsPkgSrc, "StrCompare", []any{"abc", "abd"}, strings_pkg.StrCompare},
	"strings_pkg/StrEqual":   {stringsPkgSrc, "StrEqual", []any{"hello", "hello"}, strings_pkg.StrEqual},

	// ============================================================================
	// structs
	// ============================================================================
	"structs/BasicStruct":                 {structsSrc, "BasicStruct", nil, structs.BasicStruct},
	"structs/StructPointer":               {structsSrc, "StructPointer", nil, structs.StructPointer},
	"structs/NestedStruct":                {structsSrc, "NestedStruct", nil, structs.NestedStruct},
	"structs/EmbeddedField":               {structsSrc, "EmbeddedField", nil, structs.EmbeddedField},
	"structs/StructInSlice":               {structsSrc, "StructInSlice", nil, structs.StructInSlice},
	"structs/StructAsParam":               {structsSrc, "StructAsParam", nil, structs.StructAsParam},
	"structs/StructZeroValue":             {structsSrc, "StructZeroValue", nil, structs.StructZeroValue},
	"structs/MultipleEmbedded":            {structsSrc, "MultipleEmbedded", nil, structs.MultipleEmbedded},
	"structs/DeepNesting":                 {structsSrc, "DeepNesting", nil, structs.DeepNesting},
	"structs/StructFieldMutation":         {structsSrc, "StructFieldMutation", nil, structs.StructFieldMutation},
	"structs/StructWithBool":              {structsSrc, "StructWithBool", nil, structs.StructWithBool},
	"structs/StructCopySemantics":         {structsSrc, "StructCopySemantics", nil, structs.StructCopySemantics},
	"structs/StructPointerSharing":        {structsSrc, "StructPointerSharing", nil, structs.StructPointerSharing},
	"structs/StructReturnFromFunc":        {structsSrc, "StructReturnFromFunc", nil, structs.StructReturnFromFunc},
	"structs/StructPointerReturnFromFunc": {structsSrc, "StructPointerReturnFromFunc", nil, structs.StructPointerReturnFromFunc},
	"structs/StructSliceAppend":           {structsSrc, "StructSliceAppend", nil, structs.StructSliceAppend},
	"structs/StructPointerSlice":          {structsSrc, "StructPointerSlice", nil, structs.StructPointerSlice},
	"structs/StructInMap":                 {structsSrc, "StructInMap", nil, structs.StructInMap},
	"structs/StructConditionalInit":       {structsSrc, "StructConditionalInit", nil, structs.StructConditionalInit},
	"structs/StructFieldLoop":             {structsSrc, "StructFieldLoop", nil, structs.StructFieldLoop},
	"structs/StructNestedMutation":        {structsSrc, "StructNestedMutation", nil, structs.StructNestedMutation},
	"structs/StructEmbeddedOverride":      {structsSrc, "StructEmbeddedOverride", nil, structs.StructEmbeddedOverride},
	"structs/StructWithClosure":           {structsSrc, "StructWithClosure", nil, structs.StructWithClosure},
	"structs/StructReassign":              {structsSrc, "StructReassign", nil, structs.StructReassign},
	"structs/StructSliceOfNested":         {structsSrc, "StructSliceOfNested", nil, structs.StructSliceOfNested},
	"structs/StructMultiReturn":           {structsSrc, "StructMultiReturn", nil, structs.StructMultiReturn},
	"structs/StructBuilderPattern":        {structsSrc, "StructBuilderPattern", nil, structs.StructBuilderPattern},
	"structs/StructArrayField":            {structsSrc, "StructArrayField", nil, structs.StructArrayField},
	"structs/StructEmbeddedChain":         {structsSrc, "StructEmbeddedChain", nil, structs.StructEmbeddedChain},

	// ============================================================================
	// switch
	// ============================================================================
	"switch/Simple":      {switchSrc, "Simple", nil, switch_pkg.Simple},
	"switch/Default":     {switchSrc, "Default", nil, switch_pkg.Default},
	"switch/MultiCase":   {switchSrc, "MultiCase", nil, switch_pkg.MultiCase},
	"switch/NoCondition": {switchSrc, "NoCondition", nil, switch_pkg.NoCondition},
	"switch/WithInit":    {switchSrc, "WithInit", nil, switch_pkg.WithInit},
	"switch/StringCases": {switchSrc, "StringCases", nil, switch_pkg.StringCases},
	"switch/Fallthrough": {switchSrc, "Fallthrough", nil, switch_pkg.Fallthrough},
	"switch/Nested":      {switchSrc, "Nested", nil, switch_pkg.Nested},
	// Parameterized tests
	"switch/Classify":  {switchSrc, "Classify", []any{2}, switch_pkg.Classify},
	"switch/Weekday":   {switchSrc, "Weekday", []any{3}, switch_pkg.Weekday},
	"switch/Grade":     {switchSrc, "Grade", []any{85}, switch_pkg.Grade},
	"switch/ColorCode": {switchSrc, "ColorCode", []any{"green"}, switch_pkg.ColorCode},

	// ============================================================================
	// tricky - Iteration 1
	// ============================================================================
	"tricky/ShortVarDeclShadow":  {trickySrc, "ShortVarDeclShadow", nil, tricky.ShortVarDeclShadow},
	"tricky/SliceIndexExpr":      {trickySrc, "SliceIndexExpr", nil, tricky.SliceIndexExpr},
	"tricky/MapStructKey":        {trickySrc, "MapStructKey", nil, tricky.MapStructKey},
	"tricky/NestedSliceAppend":   {trickySrc, "NestedSliceAppend", nil, tricky.NestedSliceAppend},
	"tricky/ClosureCaptureLoop":  {trickySrc, "ClosureCaptureLoop", nil, tricky.ClosureCaptureLoop},
	"tricky/DeferNamedReturn":    {trickySrc, "DeferNamedReturn", nil, tricky.DeferNamedReturn},
	"tricky/FullSliceExpr":       {trickySrc, "FullSliceExpr", nil, tricky.FullSliceExpr},
	"tricky/NestedShadowing":     {trickySrc, "NestedShadowing", nil, tricky.NestedShadowing},
	"tricky/SliceOfPointers":     {trickySrc, "SliceOfPointers", nil, tricky.SliceOfPointers},
	"tricky/MapNestedStruct":     {trickySrc, "MapNestedStruct", nil, tricky.MapNestedStruct},
	"tricky/VariadicEmpty":       {trickySrc, "VariadicEmpty", nil, tricky.VariadicEmpty},
	"tricky/VariadicOne":         {trickySrc, "VariadicOne", nil, tricky.VariadicOne},
	"tricky/VariadicMultiple":    {trickySrc, "VariadicMultiple", nil, tricky.VariadicMultiple},
	"tricky/EmbeddedField":       {trickySrc, "EmbeddedField", nil, tricky.EmbeddedField},
	"tricky/MapPointerValue":     {trickySrc, "MapPointerValue", nil, tricky.MapPointerValue},
	"tricky/ComplexBoolExpr":     {trickySrc, "ComplexBoolExpr", nil, tricky.ComplexBoolExpr},
	"tricky/SwitchFallthrough":   {trickySrc, "SwitchFallthrough", nil, tricky.SwitchFallthrough},
	"tricky/SliceCopyOperation":  {trickySrc, "SliceCopyOperation", nil, tricky.SliceCopyOperation},
	"tricky/DeferStackOrder":     {trickySrc, "DeferStackOrder", nil, tricky.DeferStackOrder},
	"tricky/InterfaceAssertion":  {trickySrc, "InterfaceAssertion", nil, tricky.InterfaceAssertion},
	"tricky/ChannelBasic":        {trickySrc, "ChannelBasic", nil, tricky.ChannelBasic},
	"tricky/SelectDefault":       {trickySrc, "SelectDefault", nil, tricky.SelectDefault},
	"tricky/RecursiveFibMemo":    {trickySrc, "RecursiveFibMemo", nil, tricky.RecursiveFibMemo},
	"tricky/PanicRecover":        {trickySrc, "PanicRecover", nil, tricky.PanicRecover},
	"tricky/ClosureWithDefer":    {trickySrc, "ClosureWithDefer", nil, tricky.ClosureWithDefer},
	"tricky/MethodOnPointer":     {trickySrc, "MethodOnPointer", nil, tricky.MethodOnPointer},
	"tricky/MultiReturnDiscard":  {trickySrc, "MultiReturnDiscard", nil, tricky.MultiReturnDiscard},
	"tricky/NilSliceAppend":      {trickySrc, "NilSliceAppend", nil, tricky.NilSliceAppend},
	"tricky/ShortCircuitEval":    {trickySrc, "ShortCircuitEval", nil, tricky.ShortCircuitEval},
	"tricky/ShortCircuitEval2":   {trickySrc, "ShortCircuitEval2", nil, tricky.ShortCircuitEval2},
	"tricky/MapDelete":           {trickySrc, "MapDelete", nil, tricky.MapDelete},
	"tricky/SliceNil":            {trickySrc, "SliceNil", nil, tricky.SliceNil},
	"tricky/MapCommaOk":          {trickySrc, "MapCommaOk", nil, tricky.MapCommaOk},
	"tricky/InterfaceNil":        {trickySrc, "InterfaceNil", nil, tricky.InterfaceNil},
	"tricky/SliceLenCap":         {trickySrc, "SliceLenCap", nil, tricky.SliceLenCap},
	"tricky/ComplexArray":        {trickySrc, "ComplexArray", nil, tricky.ComplexArray},
	"tricky/PointerArithmetic":   {trickySrc, "PointerArithmetic", nil, tricky.PointerArithmetic},
	"tricky/DoublePointer":       {trickySrc, "DoublePointer", nil, tricky.DoublePointer},
	"tricky/StructPointerMethod": {trickySrc, "StructPointerMethod", nil, tricky.StructPointerMethod},
	"tricky/ForRangeWithIndex":   {trickySrc, "ForRangeWithIndex", nil, tricky.ForRangeWithIndex},
	"tricky/ForRangeKeyValue":    {trickySrc, "ForRangeKeyValue", nil, tricky.ForRangeKeyValue},
	"tricky/StringIndex":         {trickySrc, "StringIndex", nil, tricky.StringIndex},
	"tricky/MapAssign":           {trickySrc, "MapAssign", nil, tricky.MapAssign},
	"tricky/ComplexLiteral":      {trickySrc, "ComplexLiteral", nil, tricky.ComplexLiteral},
	"tricky/ErrorReturn":         {trickySrc, "ErrorReturn", nil, tricky.ErrorReturn},
	"tricky/NilPointerCheck":     {trickySrc, "NilPointerCheck", nil, tricky.NilPointerCheck},
	"tricky/SliceAppendNil":      {trickySrc, "SliceAppendNil", nil, tricky.SliceAppendNil},
	"tricky/MapLookupNil":        {trickySrc, "MapLookupNil", nil, tricky.MapLookupNil},
	"tricky/DeferModifyNamed":    {trickySrc, "DeferModifyNamed", nil, tricky.DeferModifyNamed},
	"tricky/ForRangeMap":         {trickySrc, "ForRangeMap", nil, tricky.ForRangeMap},
	"tricky/MultipleNamedReturn": {trickySrc, "MultipleNamedReturnCombined", nil, tricky.MultipleNamedReturnCombined},

	// ============================================================================
	// tricky - Iteration 2
	// ============================================================================
	"tricky/DeferInClosure":         {trickySrc, "DeferInClosure", nil, tricky.DeferInClosure},
	"tricky/MultipleDeferSameName":  {trickySrc, "MultipleDeferSameName", nil, tricky.MultipleDeferSameName},
	"tricky/ClosureMutateOuter":     {trickySrc, "ClosureMutateOuter", nil, tricky.ClosureMutateOuter},
	"tricky/SliceAppendExpand":      {trickySrc, "SliceAppendExpand", nil, tricky.SliceAppendExpand},
	"tricky/MapIncrement":           {trickySrc, "MapIncrement", nil, tricky.MapIncrement},
	"tricky/InterfaceTypeSwitch":    {trickySrc, "InterfaceTypeSwitch", nil, tricky.InterfaceTypeSwitch},
	"tricky/PointerToSlice":         {trickySrc, "PointerToSlice", nil, tricky.PointerToSlice},
	"tricky/NestedClosure":          {trickySrc, "NestedClosure", nil, tricky.NestedClosure},
	"tricky/SliceOfSlice":           {trickySrc, "SliceOfSlice", nil, tricky.SliceOfSlice},
	"tricky/MapOfSlice":             {trickySrc, "MapOfSlice", nil, tricky.MapOfSlice},
	"tricky/StructWithSlice":        {trickySrc, "StructWithSlice", nil, tricky.StructWithSlice},
	"tricky/DeferReadAfterAssign":   {trickySrc, "DeferReadAfterAssign", nil, tricky.DeferReadAfterAssign},
	"tricky/ForRangePointer":        {trickySrc, "ForRangePointer", nil, tricky.ForRangePointer},
	"tricky/NilInterfaceValue":      {trickySrc, "NilInterfaceValue", nil, tricky.NilInterfaceValue},
	"tricky/SliceCopyOverlap":       {trickySrc, "SliceCopyOverlap", nil, tricky.SliceCopyOverlap},
	"tricky/PointerReassign":        {trickySrc, "PointerReassign", nil, tricky.PointerReassign},
	"tricky/InterfaceNilComparison": {trickySrc, "InterfaceNilComparison", nil, tricky.InterfaceNilComparison},
	"tricky/DeferClosureCapture":    {trickySrc, "DeferClosureCapture", nil, tricky.DeferClosureCapture},
	"tricky/MapLookupModify":        {trickySrc, "MapLookupModify", nil, tricky.MapLookupModify},
	"tricky/SliceZeroLength":        {trickySrc, "SliceZeroLength", nil, tricky.SliceZeroLength},

	// ============================================================================
	// tricky - Iteration 3
	// ============================================================================
	"tricky/SliceModifyViaSubslice":     {trickySrc, "SliceModifyViaSubslice", nil, tricky.SliceModifyViaSubslice},
	"tricky/MapDeleteDuringRange":       {trickySrc, "MapDeleteDuringRange", nil, tricky.MapDeleteDuringRange},
	"tricky/ClosureReturnClosure":       {trickySrc, "ClosureReturnClosure", nil, tricky.ClosureReturnClosure},
	"tricky/StructMethodOnNil":          {trickySrc, "StructMethodOnNil", nil, tricky.StructMethodOnNil},
	"tricky/ArrayPointerIndex":          {trickySrc, "ArrayPointerIndex", nil, tricky.ArrayPointerIndex},
	"tricky/SliceThreeIndex":            {trickySrc, "SliceThreeIndex", nil, tricky.SliceThreeIndex},
	"tricky/DeferInLoop":                {trickySrc, "DeferInLoop", nil, tricky.DeferInLoop},
	"tricky/StructCompare":              {trickySrc, "StructCompare", nil, tricky.StructCompare},
	"tricky/InterfaceSlice":             {trickySrc, "InterfaceSlice", nil, tricky.InterfaceSlice},
	"tricky/PointerMethodValueReceiver": {trickySrc, "PointerMethodValueReceiver", nil, tricky.PointerMethodValueReceiver},
	"tricky/SliceOfMaps":                {trickySrc, "SliceOfMaps", nil, tricky.SliceOfMaps},
	"tricky/MapWithNilValue":            {trickySrc, "MapWithNilValue", nil, tricky.MapWithNilValue},
	"tricky/SwitchNoCondition":          {trickySrc, "SwitchNoCondition", nil, tricky.SwitchNoCondition},
	"tricky/DeferModifyReturn":          {trickySrc, "DeferModifyReturn", nil, tricky.DeferModifyReturn},
	"tricky/SliceAppendToCap":           {trickySrc, "SliceAppendToCap", nil, tricky.SliceAppendToCap},
	"tricky/ForRangeStringByteIndex":    {trickySrc, "ForRangeStringByteIndex", nil, tricky.ForRangeStringByteIndex},
	"tricky/StructLiteralEmbedded":      {trickySrc, "StructLiteralEmbedded", nil, tricky.StructLiteralEmbedded},
	"tricky/MapNilKey":                  {trickySrc, "MapNilKey", nil, tricky.MapNilKey},
	"tricky/ClosureRecursive":           {trickySrc, "ClosureRecursive", nil, tricky.ClosureRecursive},

	// ============================================================================
	// tricky - Iteration 4
	// ============================================================================
	"tricky/DeferCallInDefer":       {trickySrc, "DeferCallInDefer", nil, tricky.DeferCallInDefer},
	"tricky/MapLookupAssign":        {trickySrc, "MapLookupAssign", nil, tricky.MapLookupAssign},
	"tricky/StructMethodOnValue":    {trickySrc, "StructMethodOnValue", nil, tricky.StructMethodOnValue},
	"tricky/PointerToMap":           {trickySrc, "PointerToMap", nil, tricky.PointerToMap},
	"tricky/SliceCapAfterAppend":    {trickySrc, "SliceCapAfterAppend", nil, tricky.SliceCapAfterAppend},
	"tricky/NestedMaps":             {trickySrc, "NestedMaps", nil, tricky.NestedMaps},
	"tricky/StructPointerNil":       {trickySrc, "StructPointerNil", nil, tricky.StructPointerNil},
	"tricky/VariadicWithSlice":      {trickySrc, "VariadicWithSlice", nil, tricky.VariadicWithSlice},
	"tricky/SliceMakeWithLen":       {trickySrc, "SliceMakeWithLen", nil, tricky.SliceMakeWithLen},
	"tricky/InterfaceConversion":    {trickySrc, "InterfaceConversion", nil, tricky.InterfaceConversion},
	"tricky/MapWithEmptyStringKey":  {trickySrc, "MapWithEmptyStringKey", nil, tricky.MapWithEmptyStringKey},
	"tricky/DeferPanicRecover":      {trickySrc, "DeferPanicRecover", nil, tricky.DeferPanicRecover},
	"tricky/StructWithMap":          {trickySrc, "StructWithMap", nil, tricky.StructWithMap},
	"tricky/ForRangeBreak":          {trickySrc, "ForRangeBreak", nil, tricky.ForRangeBreak},
	"tricky/SliceLiteralNested":     {trickySrc, "SliceLiteralNested", nil, tricky.SliceLiteralNested},
	"tricky/MapLiteralNested":       {trickySrc, "MapLiteralNested", nil, tricky.MapLiteralNested},
	"tricky/PointerToStructLiteral": {trickySrc, "PointerToStructLiteral", nil, tricky.PointerToStructLiteral},
	"tricky/SliceOfStructs":         {trickySrc, "SliceOfStructs", nil, tricky.SliceOfStructs},
	"tricky/MapIterateModify":       {trickySrc, "MapIterateModify", nil, tricky.MapIterateModify},

	// ============================================================================
	// tricky - Iteration 5
	// ============================================================================
	"tricky/ChannelBuffered":        {trickySrc, "ChannelBuffered", nil, tricky.ChannelBuffered},
	"tricky/StructEmbeddedMethod":   {trickySrc, "StructEmbeddedMethod", nil, tricky.StructEmbeddedMethod},
	"tricky/SliceOfChannels":        {trickySrc, "SliceOfChannels", nil, tricky.SliceOfChannels},
	"tricky/MapOfChannels":          {trickySrc, "MapOfChannels", nil, tricky.MapOfChannels},
	"tricky/MultipleAssignment":     {trickySrc, "MultipleAssignment", nil, tricky.MultipleAssignment},
	"tricky/SliceAssign":            {trickySrc, "SliceAssign", nil, tricky.SliceAssign},
	"tricky/MapTwoAssign":           {trickySrc, "MapTwoAssign", nil, tricky.MapTwoAssign},
	"tricky/StructPointerMethodNil": {trickySrc, "StructPointerMethodNil", nil, tricky.StructPointerMethodNil},
	"tricky/DeferAfterPanic":        {trickySrc, "DeferAfterPanic", nil, tricky.DeferAfterPanic},
	"tricky/SliceFromArray":         {trickySrc, "SliceFromArray", nil, tricky.SliceFromArray},
	"tricky/ArrayPointerSlice":      {trickySrc, "ArrayPointerSlice", nil, tricky.ArrayPointerSlice},
	"tricky/StructFieldPointer":     {trickySrc, "StructFieldPointer", nil, tricky.StructFieldPointer},
	"tricky/MapLenCap":              {trickySrc, "MapLenCap", nil, tricky.MapLenCap},
	"tricky/StringConcat":           {trickySrc, "StringConcat", nil, tricky.StringConcat},
	"tricky/StringLen":              {trickySrc, "StringLen", nil, tricky.StringLen},

	// ============================================================================
	// tricky - Iteration 6
	// ============================================================================
	"tricky/ComplexMapKey":           {trickySrc, "ComplexMapKey", nil, tricky.ComplexMapKey},
	"tricky/SliceReverse":            {trickySrc, "SliceReverse", nil, tricky.SliceReverse},
	"tricky/MapMerge":                {trickySrc, "MapMerge", nil, tricky.MapMerge},
	"tricky/StructZeroValue":         {trickySrc, "StructZeroValue", nil, tricky.StructZeroValue},
	"tricky/SliceDeleteByIndex":      {trickySrc, "SliceDeleteByIndex", nil, tricky.SliceDeleteByIndex},
	"tricky/MapValueOverwrite":       {trickySrc, "MapValueOverwrite", nil, tricky.MapValueOverwrite},
	"tricky/InterfaceEmbed":          {trickySrc, "InterfaceEmbed", nil, tricky.InterfaceEmbed},
	"tricky/SliceOfFuncs":            {trickySrc, "SliceOfFuncs", nil, tricky.SliceOfFuncs},
	"tricky/PointerToSliceElement":   {trickySrc, "PointerToSliceElement", nil, tricky.PointerToSliceElement},
	"tricky/MapKeyPointer":           {trickySrc, "MapKeyPointer", nil, tricky.MapKeyPointer},
	"tricky/SliceOfPointersToStruct": {trickySrc, "SliceOfPointersToStruct", nil, tricky.SliceOfPointersToStruct},
	"tricky/DoubleMapLookup":         {trickySrc, "DoubleMapLookup", nil, tricky.DoubleMapLookup},
	"tricky/StructSliceLiteral":      {trickySrc, "StructSliceLiteral", nil, tricky.StructSliceLiteral},
	"tricky/ForRangeModifyValue":     {trickySrc, "ForRangeModifyValue", nil, tricky.ForRangeModifyValue},
	"tricky/MapWithStructPointerKey": {trickySrc, "MapWithStructPointerKey", nil, tricky.MapWithStructPointerKey},
	"tricky/SliceCopyDifferentTypes": {trickySrc, "SliceCopyDifferentTypes", nil, tricky.SliceCopyDifferentTypes},
	"tricky/NestedStructWithPointer": {trickySrc, "NestedStructWithPointer", nil, tricky.NestedStructWithPointer},
	"tricky/SliceOfSlicesAppend":     {trickySrc, "SliceOfSlicesAppend", nil, tricky.SliceOfSlicesAppend},
	"tricky/MapDeleteAll":            {trickySrc, "MapDeleteAll", nil, tricky.MapDeleteAll},

	// ============================================================================
	// tricky - Iteration 7
	// ============================================================================
	"tricky/StructPointerSlice":        {trickySrc, "StructPointerSlice", nil, tricky.StructPointerSlice},
	"tricky/MapWithInterfaceKey":       {trickySrc, "MapWithInterfaceKey", nil, tricky.MapWithInterfaceKey},
	"tricky/SliceOfInterfaces":         {trickySrc, "SliceOfInterfaces", nil, tricky.SliceOfInterfaces},
	"tricky/NestedPointerStruct":       {trickySrc, "NestedPointerStruct", nil, tricky.NestedPointerStruct},
	"tricky/StructMethodOnNilPointer":  {trickySrc, "StructMethodOnNilPointer", nil, tricky.StructMethodOnNilPointer},
	"tricky/SliceAppendToSlice":        {trickySrc, "SliceAppendToSlice", nil, tricky.SliceAppendToSlice},
	"tricky/MapLookupWithDefault":      {trickySrc, "MapLookupWithDefault", nil, tricky.MapLookupWithDefault},
	"tricky/StructFieldUpdate":         {trickySrc, "StructFieldUpdate", nil, tricky.StructFieldUpdate},
	"tricky/PointerToNilSlice":         {trickySrc, "PointerToNilSlice", nil, tricky.PointerToNilSlice},
	"tricky/SliceCopyToSubslice":       {trickySrc, "SliceCopyToSubslice", nil, tricky.SliceCopyToSubslice},
	"tricky/StructWithMultipleFields":  {trickySrc, "StructWithMultipleFields", nil, tricky.StructWithMultipleFields},
	"tricky/ForRangeContinue":          {trickySrc, "ForRangeContinue", nil, tricky.ForRangeContinue},
	"tricky/MapWithBoolKey":            {trickySrc, "MapWithBoolKey", nil, tricky.MapWithBoolKey},
	"tricky/SliceInsert":               {trickySrc, "SliceInsert", nil, tricky.SliceInsert},
	"tricky/StructEmbeddedFieldAccess": {trickySrc, "StructEmbeddedFieldAccess", nil, tricky.StructEmbeddedFieldAccess},
	"tricky/PointerToChannel":          {trickySrc, "PointerToChannel", nil, tricky.PointerToChannel},
	"tricky/MapKeyModification":        {trickySrc, "MapKeyModification", nil, tricky.MapKeyModification},
	"tricky/SliceRangeModify":          {trickySrc, "SliceRangeModify", nil, tricky.SliceRangeModify},
	"tricky/StructLiteralShort":        {trickySrc, "StructLiteralShort", nil, tricky.StructLiteralShort},

	// ============================================================================
	// tricky - Iteration 8
	// ============================================================================
	"tricky/SliceDrain":        {trickySrc, "SliceDrain", nil, tricky.SliceDrain},
	"tricky/MapClear":          {trickySrc, "MapClear", nil, tricky.MapClear},
	"tricky/StructCopy":        {trickySrc, "StructCopy", nil, tricky.StructCopy},
	"tricky/PointerStructCopy": {trickySrc, "PointerStructCopy", nil, tricky.PointerStructCopy},
	"tricky/SliceFilter":       {trickySrc, "SliceFilter", nil, tricky.SliceFilter},
	"tricky/MapTransform":      {trickySrc, "MapTransform", nil, tricky.MapTransform},
	"tricky/SliceContains":     {trickySrc, "SliceContains", nil, tricky.SliceContains},
	"tricky/MapKeys":           {trickySrc, "MapKeys", nil, tricky.MapKeys},
	// Parameterized tricky tests - PASSING
	"tricky/ModuloWithNegative": {trickySrc, "ModuloWithNegative", []any{-17, 5}, tricky.ModuloWithNegative},
	"tricky/PowerRecursive":     {trickySrc, "PowerRecursive", []any{2, 10}, tricky.PowerRecursive},
	"tricky/GCD":                {trickySrc, "GCD", []any{48, 18}, tricky.GCD},
	"tricky/LCM":                {trickySrc, "LCM", []any{4, 6}, tricky.LCM},
	"tricky/SumOfDigits":        {trickySrc, "SumOfDigits", []any{12345}, tricky.SumOfDigits},
	"tricky/MapEqualParam":      {trickySrc, "MapEqualParam", []any{map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2}}, tricky.MapEqualParam},
	"tricky/StringJoinParam":    {trickySrc, "StringJoinParam", []any{[]string{"a", "b", "c"}, ","}, tricky.StringJoinParam},
	"tricky/SliceEqualParam":    {trickySrc, "SliceEqualParam", []any{[]int{1, 2, 3}, []int{1, 2, 3}}, tricky.SliceEqualParam},
	"tricky/MapInvertParam":     {trickySrc, "MapInvertParam", []any{map[string]int{"a": 1, "b": 2, "c": 3}}, tricky.MapInvertParam},
	// Previously known issues — now fixed
	"tricky/StringReverse":        {trickySrc, "StringReverse", []any{"hello"}, tricky.StringReverse},
	"tricky/Clamp":                {trickySrc, "Clamp", []any{150, 0, 100}, tricky.Clamp},
	"tricky/Sign":                 {trickySrc, "Sign", []any{-42}, tricky.Sign},
	"tricky/SliceUniqueParam":     {trickySrc, "SliceUniqueParam", []any{[]int{1, 2, 2, 3, 3, 3}}, tricky.SliceUniqueParam},
	"tricky/SliceInterleave":      {trickySrc, "SliceInterleave", []any{[]int{1, 3, 5}, []int{2, 4, 6}}, tricky.SliceInterleave},
	"tricky/SliceRotateLeftParam": {trickySrc, "SliceRotateLeftParam", []any{[]int{1, 2, 3, 4, 5}, 2}, tricky.SliceRotateLeftParam},
	"tricky/BitCountOnes":         {trickySrc, "BitCountOnes", []any{255}, tricky.BitCountOnes},
	"tricky/BinomialCoefficient":  {trickySrc, "BinomialCoefficient", []any{5, 2}, tricky.BinomialCoefficient},
	"tricky/FibonacciNth":         {trickySrc, "FibonacciNth", []any{20}, tricky.FibonacciNth},
	"tricky/IsPrime":              {trickySrc, "IsPrime", []any{17}, tricky.IsPrime},
	"tricky/FactorialIterative":   {trickySrc, "FactorialIterative", []any{5}, tricky.FactorialIterative},
	"tricky/MapDeepCopy":          {trickySrc, "MapDeepCopy", []any{map[int][]int{1: {1, 2}, 2: {3, 4}}}, tricky.MapDeepCopy},

	// ============================================================================
	// tricky - Additional tests (previously missing)
	// ============================================================================
	"tricky/ClosureArgDefaultTest":                {trickySrc, "ClosureArgDefaultTest", nil, tricky.ClosureArgDefaultTest},
	"tricky/ClosureAsArg":                         {trickySrc, "ClosureAsArg", nil, tricky.ClosureAsArg},
	"tricky/ClosureCaptureAndModifyTest":          {trickySrc, "ClosureCaptureAndModifyTest", nil, tricky.ClosureCaptureAndModifyTest},
	"tricky/ClosureCaptureLoopVarTest":            {trickySrc, "ClosureCaptureLoopVarTest", nil, tricky.ClosureCaptureLoopVarTest},
	"tricky/ClosureCapturesTwoTest":               {trickySrc, "ClosureCapturesTwoTest", nil, tricky.ClosureCapturesTwoTest},
	"tricky/ClosureCapturingPointer":              {trickySrc, "ClosureCapturingPointer", nil, tricky.ClosureCapturingPointer},
	"tricky/ClosureCompose":                       {trickySrc, "ClosureCompose", nil, tricky.ClosureCompose},
	"tricky/ClosureComposeTest":                   {trickySrc, "ClosureComposeTest", nil, tricky.ClosureComposeTest},
	"tricky/ClosureConstTest":                     {trickySrc, "ClosureConstTest", nil, tricky.ClosureConstTest},
	"tricky/ClosureCounter":                       {trickySrc, "ClosureCounter", nil, tricky.ClosureCounter},
	"tricky/ClosureCounterResetTest":              {trickySrc, "ClosureCounterResetTest", nil, tricky.ClosureCounterResetTest},
	"tricky/ClosureCounterStateTest":              {trickySrc, "ClosureCounterStateTest", nil, tricky.ClosureCounterStateTest},
	"tricky/ClosureCounterTest":                   {trickySrc, "ClosureCounterTest", nil, tricky.ClosureCounterTest},
	"tricky/ClosureCurry":                         {trickySrc, "ClosureCurry", nil, tricky.ClosureCurry},
	"tricky/ClosureCurryMultipleArgTest":          {trickySrc, "ClosureCurryMultipleArgTest", nil, tricky.ClosureCurryMultipleArgTest},
	"tricky/ClosureCurryMultipleTest":             {trickySrc, "ClosureCurryMultipleTest", nil, tricky.ClosureCurryMultipleTest},
	"tricky/ClosureCurryTest":                     {trickySrc, "ClosureCurryTest", nil, tricky.ClosureCurryTest},
	"tricky/ClosureEnvCaptureTest":                {trickySrc, "ClosureEnvCaptureTest", nil, tricky.ClosureEnvCaptureTest},
	"tricky/ClosureFibonacci":                     {trickySrc, "ClosureFibonacci", nil, tricky.ClosureFibonacci},
	"tricky/ClosureFlip":                          {trickySrc, "ClosureFlip", nil, tricky.ClosureFlip},
	"tricky/ClosureFlipTest":                      {trickySrc, "ClosureFlipTest", nil, tricky.ClosureFlipTest},
	"tricky/ClosureMap":                           {trickySrc, "ClosureMap", nil, tricky.ClosureMap},
	"tricky/ClosureMapBuilderTest":                {trickySrc, "ClosureMapBuilderTest", nil, tricky.ClosureMapBuilderTest},
	"tricky/ClosureMemoize":                       {trickySrc, "ClosureMemoize", nil, tricky.ClosureMemoize},
	"tricky/ClosureMemoizeRecursive":              {trickySrc, "ClosureMemoizeRecursive", nil, tricky.ClosureMemoizeRecursive},
	"tricky/ClosureMemoizeRecursiveTest":          {trickySrc, "ClosureMemoizeRecursiveTest", nil, tricky.ClosureMemoizeRecursiveTest},
	"tricky/ClosureModifyOuterVarTest":            {trickySrc, "ClosureModifyOuterVarTest", nil, tricky.ClosureModifyOuterVarTest},
	"tricky/ClosureMultiCapture":                  {trickySrc, "ClosureMultiCapture", nil, tricky.ClosureMultiCapture},
	"tricky/ClosureMultipleCallsTest":             {trickySrc, "ClosureMultipleCallsTest", nil, tricky.ClosureMultipleCallsTest},
	"tricky/ClosureMultiReturnTest":               {trickySrc, "ClosureMultiReturnTest", nil, tricky.ClosureMultiReturnTest},
	"tricky/ClosureMutateCapturedSlice":           {trickySrc, "ClosureMutateCapturedSlice", nil, tricky.ClosureMutateCapturedSlice},
	"tricky/ClosureMutateClosureTest":             {trickySrc, "ClosureMutateClosureTest", nil, tricky.ClosureMutateClosureTest},
	"tricky/ClosureMutatesOuterTest":              {trickySrc, "ClosureMutatesOuterTest", nil, tricky.ClosureMutatesOuterTest},
	"tricky/ClosureOnce":                          {trickySrc, "ClosureOnce", nil, tricky.ClosureOnce},
	"tricky/ClosureOnceTest":                      {trickySrc, "ClosureOnceTest", nil, tricky.ClosureOnceTest},
	"tricky/ClosurePartial":                       {trickySrc, "ClosurePartial", nil, tricky.ClosurePartial},
	"tricky/ClosurePartialApplyTest":              {trickySrc, "ClosurePartialApplyTest", nil, tricky.ClosurePartialApplyTest},
	"tricky/ClosurePartialTest":                   {trickySrc, "ClosurePartialTest", nil, tricky.ClosurePartialTest},
	"tricky/ClosurePipeline":                      {trickySrc, "ClosurePipeline", nil, tricky.ClosurePipeline},
	"tricky/ClosurePtrCaptureTest":                {trickySrc, "ClosurePtrCaptureTest", nil, tricky.ClosurePtrCaptureTest},
	"tricky/ClosureRecursiveMemoTest":             {trickySrc, "ClosureRecursiveMemoTest", nil, tricky.ClosureRecursiveMemoTest},
	"tricky/ClosureRecursiveSimpleTest":           {trickySrc, "ClosureRecursiveSimpleTest", nil, tricky.ClosureRecursiveSimpleTest},
	"tricky/ClosureReturningClosure":              {trickySrc, "ClosureReturningClosure", nil, tricky.ClosureReturningClosure},
	"tricky/ClosureReturnMultipleTest":            {trickySrc, "ClosureReturnMultipleTest", nil, tricky.ClosureReturnMultipleTest},
	"tricky/ClosureReturnsClosureTest":            {trickySrc, "ClosureReturnsClosureTest", nil, tricky.ClosureReturnsClosureTest},
	"tricky/ClosureReturnsValueTest":              {trickySrc, "ClosureReturnsValueTest", nil, tricky.ClosureReturnsValueTest},
	"tricky/ClosureReturnValueTest":               {trickySrc, "ClosureReturnValueTest", nil, tricky.ClosureReturnValueTest},
	"tricky/ClosureSliceAccumTest":                {trickySrc, "ClosureSliceAccumTest", nil, tricky.ClosureSliceAccumTest},
	"tricky/ClosureSliceBuilderTest":              {trickySrc, "ClosureSliceBuilderTest", nil, tricky.ClosureSliceBuilderTest},
	"tricky/ClosureSliceCaptureTest":              {trickySrc, "ClosureSliceCaptureTest", nil, tricky.ClosureSliceCaptureTest},
	"tricky/ClosureTap":                           {trickySrc, "ClosureTap", nil, tricky.ClosureTap},
	"tricky/ClosureTapTest":                       {trickySrc, "ClosureTapTest", nil, tricky.ClosureTapTest},
	"tricky/ClosureVarCaptureTest":                {trickySrc, "ClosureVarCaptureTest", nil, tricky.ClosureVarCaptureTest},
	"tricky/ClosureWithDeferAndReturn":            {trickySrc, "ClosureWithDeferAndReturn", nil, tricky.ClosureWithDeferAndReturn},
	"tricky/ClosureWithDeferTest":                 {trickySrc, "ClosureWithDeferTest", nil, tricky.ClosureWithDeferTest},
	"tricky/ClosureWithExternalVar":               {trickySrc, "ClosureWithExternalVar", nil, tricky.ClosureWithExternalVar},
	"tricky/ClosureWithLocalVarTest":              {trickySrc, "ClosureWithLocalVarTest", nil, tricky.ClosureWithLocalVarTest},
	"tricky/ClosureWithLoopVar":                   {trickySrc, "ClosureWithLoopVar", nil, tricky.ClosureWithLoopVar},
	"tricky/ClosureWithMapCaptureTest":            {trickySrc, "ClosureWithMapCaptureTest", nil, tricky.ClosureWithMapCaptureTest},
	"tricky/ClosureWithMultipleReturns":           {trickySrc, "ClosureWithMultipleReturns", nil, tricky.ClosureWithMultipleReturns},
	"tricky/ClosureWithRecursion":                 {trickySrc, "ClosureWithRecursion", nil, tricky.ClosureWithRecursion},
	"tricky/ClosureWithStructCaptureTest":         {trickySrc, "ClosureWithStructCaptureTest", nil, tricky.ClosureWithStructCaptureTest},
	"tricky/ClosureWithVarCaptureTest":            {trickySrc, "ClosureWithVarCaptureTest", nil, tricky.ClosureWithVarCaptureTest},
	"tricky/DeferAfterReturnTest":                 {trickySrc, "DeferAfterReturnTest", nil, tricky.DeferAfterReturnTest},
	"tricky/DeferCaptureMapTest":                  {trickySrc, "DeferCaptureMapTest", nil, tricky.DeferCaptureMapTest},
	"tricky/DeferCaptureSliceTest":                {trickySrc, "DeferCaptureSliceTest", nil, tricky.DeferCaptureSliceTest},
	"tricky/DeferCaptureValueTest":                {trickySrc, "DeferCaptureValueTest", nil, tricky.DeferCaptureValueTest},
	"tricky/DeferClosureArgTest":                  {trickySrc, "DeferClosureArgTest", nil, tricky.DeferClosureArgTest},
	"tricky/DeferClosureCaptureModifyTest":        {trickySrc, "DeferClosureCaptureModifyTest", nil, tricky.DeferClosureCaptureModifyTest},
	"tricky/DeferClosureModifyingNamed":           {trickySrc, "DeferClosureModifyingNamed", nil, tricky.DeferClosureModifyingNamed},
	"tricky/DeferClosureModifyTest":               {trickySrc, "DeferClosureModifyTest", nil, tricky.DeferClosureModifyTest},
	"tricky/DeferClosureNestedTest":               {trickySrc, "DeferClosureNestedTest", nil, tricky.DeferClosureNestedTest},
	"tricky/DeferConditional":                     {trickySrc, "DeferConditional", nil, tricky.DeferConditional},
	"tricky/DeferConditionalModifyTest":           {trickySrc, "DeferConditionalModifyTest", nil, tricky.DeferConditionalModifyTest},
	"tricky/DeferInClosureTest":                   {trickySrc, "DeferInClosureTest", nil, tricky.DeferInClosureTest},
	"tricky/DeferInClosureWithArg":                {trickySrc, "DeferInClosureWithArg", nil, tricky.DeferInClosureWithArg},
	"tricky/DeferInGoroutine":                     {trickySrc, "DeferInGoroutine", nil, tricky.DeferInGoroutine},
	"tricky/DeferInMultipleFunctions":             {trickySrc, "DeferInMultipleFunctions", nil, tricky.DeferInMultipleFunctions},
	"tricky/DeferInNestedFunction":                {trickySrc, "DeferInNestedFunction", nil, tricky.DeferInNestedFunction},
	"tricky/DeferMapModifyTest":                   {trickySrc, "DeferMapModifyTest", nil, tricky.DeferMapModifyTest},
	"tricky/DeferModifiesReturnTest":              {trickySrc, "DeferModifiesReturnTest", nil, tricky.DeferModifiesReturnTest},
	"tricky/DeferModifyCapture":                   {trickySrc, "DeferModifyCapture", nil, tricky.DeferModifyCapture},
	"tricky/DeferModifyMap":                       {trickySrc, "DeferModifyMap", nil, tricky.DeferModifyMap},
	"tricky/DeferModifyMapNamedTest":              {trickySrc, "DeferModifyMapNamedTest", nil, tricky.DeferModifyMapNamedTest},
	"tricky/DeferModifyMapTest":                   {trickySrc, "DeferModifyMapTest", nil, tricky.DeferModifyMapTest},
	"tricky/DeferModifyMultiple":                  {trickySrc, "DeferModifyMultiple", nil, tricky.DeferModifyMultiple},
	"tricky/DeferModifyMultipleCombined":          {trickySrc, "DeferModifyMultipleCombined", nil, tricky.DeferModifyMultipleCombined},
	"tricky/DeferModifyMultipleNamedTest":         {trickySrc, "DeferModifyMultipleNamedTest", nil, tricky.DeferModifyMultipleNamedTest},
	"tricky/DeferModifyNamedReturnTest":           {trickySrc, "DeferModifyNamedReturnTest", nil, tricky.DeferModifyNamedReturnTest},
	"tricky/DeferModifyPtrTest":                   {trickySrc, "DeferModifyPtrTest", nil, tricky.DeferModifyPtrTest},
	"tricky/DeferModifyReturnValue":               {trickySrc, "DeferModifyReturnValue", nil, tricky.DeferModifyReturnValue},
	"tricky/DeferModifyReturnValueTest":           {trickySrc, "DeferModifyReturnValueTest", nil, tricky.DeferModifyReturnValueTest},
	"tricky/DeferModifySlice":                     {trickySrc, "DeferModifySlice", nil, tricky.DeferModifySlice},
	"tricky/DeferModifySliceTest":                 {trickySrc, "DeferModifySliceTest", nil, tricky.DeferModifySliceTest},
	"tricky/DeferMultiNamedReturnTest":            {trickySrc, "DeferMultiNamedReturnTest", nil, tricky.DeferMultiNamedReturnTest},
	"tricky/DeferMultipleCalls":                   {trickySrc, "DeferMultipleCalls", nil, tricky.DeferMultipleCalls},
	"tricky/DeferMultipleExecTest":                {trickySrc, "DeferMultipleExecTest", nil, tricky.DeferMultipleExecTest},
	"tricky/DeferMultipleFuncTest":                {trickySrc, "DeferMultipleFuncTest", nil, tricky.DeferMultipleFuncTest},
	"tricky/DeferMultipleNamedTest":               {trickySrc, "DeferMultipleNamedTest", nil, tricky.DeferMultipleNamedTest},
	"tricky/DeferMultipleVars":                    {trickySrc, "DeferMultipleVars", nil, tricky.DeferMultipleVars},
	"tricky/DeferNamedMultiTest":                  {trickySrc, "DeferNamedMultiTest", nil, tricky.DeferNamedMultiTest},
	"tricky/DeferNamedResultChainTest":            {trickySrc, "DeferNamedResultChainTest", nil, tricky.DeferNamedResultChainTest},
	"tricky/DeferNamedResultNilTest":              {trickySrc, "DeferNamedResultNilTest", nil, tricky.DeferNamedResultNilTest},
	"tricky/DeferNamedResultTest":                 {trickySrc, "DeferNamedResultTest", nil, tricky.DeferNamedResultTest},
	"tricky/DeferNamedReturnCaptureTest":          {trickySrc, "DeferNamedReturnCaptureTest", nil, tricky.DeferNamedReturnCaptureTest},
	"tricky/DeferNamedReturnCombineTest":          {trickySrc, "DeferNamedReturnCombineTest", nil, tricky.DeferNamedReturnCombineTest},
	"tricky/DeferNamedReturnDoubleTest":           {trickySrc, "DeferNamedReturnDoubleTest", nil, tricky.DeferNamedReturnDoubleTest},
	"tricky/DeferNamedReturnModifyTest":           {trickySrc, "DeferNamedReturnModifyTest", nil, tricky.DeferNamedReturnModifyTest},
	"tricky/DeferNamedReturnMultiTest":            {trickySrc, "DeferNamedReturnMultiTest", nil, tricky.DeferNamedReturnMultiTest},
	"tricky/DeferNamedReturnNilPtrTest":           {trickySrc, "DeferNamedReturnNilPtrTest", nil, tricky.DeferNamedReturnNilPtrTest},
	"tricky/DeferNamedReturnNilTest":              {trickySrc, "DeferNamedReturnNilTest", nil, tricky.DeferNamedReturnNilTest},
	"tricky/DeferNamedReturnOrderTest":            {trickySrc, "DeferNamedReturnOrderTest", nil, tricky.DeferNamedReturnOrderTest},
	"tricky/DeferPanicRecoverValueTest":           {trickySrc, "DeferPanicRecoverValueTest", nil, tricky.DeferPanicRecoverValueTest},
	"tricky/DeferReadCapture":                     {trickySrc, "DeferReadCapture", nil, tricky.DeferReadCapture},
	"tricky/DeferRecoverPanicTest":                {trickySrc, "DeferRecoverPanicTest", nil, tricky.DeferRecoverPanicTest},
	"tricky/DeferReturnValue":                     {trickySrc, "DeferReturnValue", nil, tricky.DeferReturnValue},
	"tricky/DeferReturnValueModifyTest":           {trickySrc, "DeferReturnValueModifyTest", nil, tricky.DeferReturnValueModifyTest},
	"tricky/DeferStackTest":                       {trickySrc, "DeferStackTest", nil, tricky.DeferStackTest},
	"tricky/DeferWithCapture":                     {trickySrc, "DeferWithCapture", nil, tricky.DeferWithCapture},
	"tricky/DeferWithClosureArg":                  {trickySrc, "DeferWithClosureArg", nil, tricky.DeferWithClosureArg},
	"tricky/DeferWithClosureResult":               {trickySrc, "DeferWithClosureResult", nil, tricky.DeferWithClosureResult},
	"tricky/DeferWithLoop":                        {trickySrc, "DeferWithLoop", nil, tricky.DeferWithLoop},
	"tricky/DeferWithMultipleReturns":             {trickySrc, "DeferWithMultipleReturns", nil, tricky.DeferWithMultipleReturns},
	"tricky/DeferWithMultipleReturnsCombined":     {trickySrc, "DeferWithMultipleReturnsCombined", nil, tricky.DeferWithMultipleReturnsCombined},
	"tricky/DeferWithNamedResultMultiple":         {trickySrc, "DeferWithNamedResultMultiple", nil, tricky.DeferWithNamedResultMultiple},
	"tricky/DeferWithNamedResultMultipleCombined": {trickySrc, "DeferWithNamedResultMultipleCombined", nil, tricky.DeferWithNamedResultMultipleCombined},
	"tricky/DeferWithNamedReturn":                 {trickySrc, "DeferWithNamedReturn", nil, tricky.DeferWithNamedReturn},
	"tricky/DeferWithRecoveredPanic":              {trickySrc, "DeferWithRecoveredPanic", nil, tricky.DeferWithRecoveredPanic},
	"tricky/DeferWithRetFunc":                     {trickySrc, "DeferWithRetFunc", nil, tricky.DeferWithRetFunc},
	"tricky/DeferWithReturnFunc":                  {trickySrc, "DeferWithReturnFunc", nil, tricky.DeferWithReturnFunc},
	"tricky/InterfaceMethod":                      {trickySrc, "InterfaceMethod", nil, tricky.InterfaceMethod},
	"tricky/InterfaceNilTypeAssertion":            {trickySrc, "InterfaceNilTypeAssertion", nil, tricky.InterfaceNilTypeAssertion},
	"tricky/InterfaceSliceTypeAssert":             {trickySrc, "InterfaceSliceTypeAssert", nil, tricky.InterfaceSliceTypeAssert},
	"tricky/MapAll":                               {trickySrc, "MapAll", nil, tricky.MapAll},
	"tricky/MapAllMatch":                          {trickySrc, "MapAllMatch", nil, tricky.MapAllMatch},
	"tricky/MapAllTest":                           {trickySrc, "MapAllTest", nil, tricky.MapAllTest},
	"tricky/MapAny":                               {trickySrc, "MapAny", nil, tricky.MapAny},
	"tricky/MapAnyMatch":                          {trickySrc, "MapAnyMatch", nil, tricky.MapAnyMatch},
	"tricky/MapAnyTest":                           {trickySrc, "MapAnyTest", nil, tricky.MapAnyTest},
	"tricky/MapAnyValueTest":                      {trickySrc, "MapAnyValueTest", nil, tricky.MapAnyValueTest},
	"tricky/MapApplyToValuesTest":                 {trickySrc, "MapApplyToValuesTest", nil, tricky.MapApplyToValuesTest},
	"tricky/MapClearMakeTest":                     {trickySrc, "MapClearMakeTest", nil, tricky.MapClearMakeTest},
	"tricky/MapClearRange":                        {trickySrc, "MapClearRange", nil, tricky.MapClearRange},
	"tricky/MapCombine":                           {trickySrc, "MapCombine", nil, tricky.MapCombine},
	"tricky/MapCombineSameKeyTest":                {trickySrc, "MapCombineSameKeyTest", nil, tricky.MapCombineSameKeyTest},
	"tricky/MapCombineTest":                       {trickySrc, "MapCombineTest", nil, tricky.MapCombineTest},
	"tricky/MapCompact":                           {trickySrc, "MapCompact", nil, tricky.MapCompact},
	"tricky/MapContainsVal":                       {trickySrc, "MapContainsVal", nil, tricky.MapContainsVal},
	"tricky/MapCopy":                              {trickySrc, "MapCopy", nil, tricky.MapCopy},
	"tricky/MapCountByKey":                        {trickySrc, "MapCountByKey", nil, tricky.MapCountByKey},
	"tricky/MapCountByValueTest":                  {trickySrc, "MapCountByValueTest", nil, tricky.MapCountByValueTest},
	"tricky/MapCountIfTest":                       {trickySrc, "MapCountIfTest", nil, tricky.MapCountIfTest},
	"tricky/MapCountPredTest":                     {trickySrc, "MapCountPredTest", nil, tricky.MapCountPredTest},
	"tricky/MapCountValues":                       {trickySrc, "MapCountValues", nil, tricky.MapCountValues},
	"tricky/MapDedup":                             {trickySrc, "MapDedup", nil, tricky.MapDedup},
	"tricky/MapDeepGet":                           {trickySrc, "MapDeepGet", nil, tricky.MapDeepGet},
	"tricky/MapDeepMerge":                         {trickySrc, "MapDeepMerge", nil, tricky.MapDeepMerge},
	"tricky/MapDeepSet":                           {trickySrc, "MapDeepSet", nil, tricky.MapDeepSet},
	"tricky/MapDefaultPattern":                    {trickySrc, "MapDefaultPattern", nil, tricky.MapDefaultPattern},
	"tricky/MapDiff":                              {trickySrc, "MapDiff", nil, tricky.MapDiff},
	"tricky/MapDiffKeysTest":                      {trickySrc, "MapDiffKeysTest", nil, tricky.MapDiffKeysTest},
	"tricky/MapDiffTest":                          {trickySrc, "MapDiffTest", nil, tricky.MapDiffTest},
	"tricky/MapDropKeys":                          {trickySrc, "MapDropKeys", nil, tricky.MapDropKeys},
	"tricky/MapDropTest":                          {trickySrc, "MapDropTest", nil, tricky.MapDropTest},
	"tricky/MapDropWhileTest":                     {trickySrc, "MapDropWhileTest", nil, tricky.MapDropWhileTest},
	"tricky/MapEmptyCheck":                        {trickySrc, "MapEmptyCheck", nil, tricky.MapEmptyCheck},
	"tricky/MapEmptyKey":                          {trickySrc, "MapEmptyKey", nil, tricky.MapEmptyKey},
	"tricky/MapEvery":                             {trickySrc, "MapEvery", nil, tricky.MapEvery},
	"tricky/MapFilter":                            {trickySrc, "MapFilter", nil, tricky.MapFilter},
	"tricky/MapFilterByKeyTest":                   {trickySrc, "MapFilterByKeyTest", nil, tricky.MapFilterByKeyTest},
	"tricky/MapFilterByValueTest":                 {trickySrc, "MapFilterByValueTest", nil, tricky.MapFilterByValueTest},
	"tricky/MapFilterKeys":                        {trickySrc, "MapFilterKeys", nil, tricky.MapFilterKeys},
	"tricky/MapFilterKeysTest":                    {trickySrc, "MapFilterKeysTest", nil, tricky.MapFilterKeysTest},
	"tricky/MapFind":                              {trickySrc, "MapFind", nil, tricky.MapFind},
	"tricky/MapFindKeyTest":                       {trickySrc, "MapFindKeyTest", nil, tricky.MapFindKeyTest},
	"tricky/MapFindValueTest":                     {trickySrc, "MapFindValueTest", nil, tricky.MapFindValueTest},
	"tricky/MapFirstKey":                          {trickySrc, "MapFirstKey", nil, tricky.MapFirstKey},
	"tricky/MapFlatten":                           {trickySrc, "MapFlatten", nil, tricky.MapFlatten},
	"tricky/MapFlattenTest":                       {trickySrc, "MapFlattenTest", nil, tricky.MapFlattenTest},
	"tricky/MapFlip":                              {trickySrc, "MapFlip", nil, tricky.MapFlip},
	"tricky/MapFloatKey":                          {trickySrc, "MapFloatKey", nil, tricky.MapFloatKey},
	"tricky/MapForEach":                           {trickySrc, "MapForEach", nil, tricky.MapForEach},
	"tricky/MapGetOrCreate":                       {trickySrc, "MapGetOrCreate", nil, tricky.MapGetOrCreate},
	"tricky/MapGetOrDefaultTest":                  {trickySrc, "MapGetOrDefaultTest", nil, tricky.MapGetOrDefaultTest},
	"tricky/MapGetOrElse":                         {trickySrc, "MapGetOrElse", nil, tricky.MapGetOrElse},
	"tricky/MapGetOrInsertDefaultTest":            {trickySrc, "MapGetOrInsertDefaultTest", nil, tricky.MapGetOrInsertDefaultTest},
	"tricky/MapGetOrInsertTest":                   {trickySrc, "MapGetOrInsertTest", nil, tricky.MapGetOrInsertTest},
	"tricky/MapGetOrSet":                          {trickySrc, "MapGetOrSet", nil, tricky.MapGetOrSet},
	"tricky/MapGetSetTest":                        {trickySrc, "MapGetSetTest", nil, tricky.MapGetSetTest},
	"tricky/MapGroupBy":                           {trickySrc, "MapGroupBy", nil, tricky.MapGroupBy},
	"tricky/MapGroupByKey":                        {trickySrc, "MapGroupByKey", nil, tricky.MapGroupByKey},
	"tricky/MapGroupByValueTest":                  {trickySrc, "MapGroupByValueTest", nil, tricky.MapGroupByValueTest},
	"tricky/MapHasKey":                            {trickySrc, "MapHasKey", nil, tricky.MapHasKey},
	"tricky/MapHasKeyAndValueTest":                {trickySrc, "MapHasKeyAndValueTest", nil, tricky.MapHasKeyAndValueTest},
	"tricky/MapHasKeyMultiple":                    {trickySrc, "MapHasKeyMultiple", nil, tricky.MapHasKeyMultiple},
	"tricky/MapHasKeyMultiTest":                   {trickySrc, "MapHasKeyMultiTest", nil, tricky.MapHasKeyMultiTest},
	"tricky/MapHasKeyNilTest":                     {trickySrc, "MapHasKeyNilTest", nil, tricky.MapHasKeyNilTest},
	"tricky/MapHasKeySlice":                       {trickySrc, "MapHasKeySlice", nil, tricky.MapHasKeySlice},
	"tricky/MapHasKeySliceTest":                   {trickySrc, "MapHasKeySliceTest", nil, tricky.MapHasKeySliceTest},
	"tricky/MapHasKeyTest":                        {trickySrc, "MapHasKeyTest", nil, tricky.MapHasKeyTest},
	"tricky/MapHasValueCond":                      {trickySrc, "MapHasValueCond", nil, tricky.MapHasValueCond},
	"tricky/MapHasValuesTest":                     {trickySrc, "MapHasValuesTest", nil, tricky.MapHasValuesTest},
	"tricky/MapIncrementAll":                      {trickySrc, "MapIncrementAll", nil, tricky.MapIncrementAll},
	"tricky/MapIncrementValueTest":                {trickySrc, "MapIncrementValueTest", nil, tricky.MapIncrementValueTest},
	"tricky/MapIndexBy":                           {trickySrc, "MapIndexBy", nil, tricky.MapIndexBy},
	"tricky/MapIntersect":                         {trickySrc, "MapIntersect", nil, tricky.MapIntersect},
	"tricky/MapIntersectKeysFunc":                 {trickySrc, "MapIntersectKeysFunc", nil, tricky.MapIntersectKeysFunc},
	"tricky/MapIntKey":                            {trickySrc, "MapIntKey", nil, tricky.MapIntKey},
	"tricky/MapInvert":                            {trickySrc, "MapInvert", nil, tricky.MapInvert},
	"tricky/MapInvertSlice":                       {trickySrc, "MapInvertSlice", nil, tricky.MapInvertSlice},
	"tricky/MapIsEmptyTest":                       {trickySrc, "MapIsEmptyTest", nil, tricky.MapIsEmptyTest},
	"tricky/MapIterateDelete":                     {trickySrc, "MapIterateDelete", nil, tricky.MapIterateDelete},
	"tricky/MapKeepIfTest":                        {trickySrc, "MapKeepIfTest", nil, tricky.MapKeepIfTest},
	"tricky/MapKeepKeysTest":                      {trickySrc, "MapKeepKeysTest", nil, tricky.MapKeepKeysTest},
	"tricky/MapKeyDiffTest":                       {trickySrc, "MapKeyDiffTest", nil, tricky.MapKeyDiffTest},
	"tricky/MapKeyExistsMultiTest":                {trickySrc, "MapKeyExistsMultiTest", nil, tricky.MapKeyExistsMultiTest},
	"tricky/MapKeyExistsTest":                     {trickySrc, "MapKeyExistsTest", nil, tricky.MapKeyExistsTest},
	"tricky/MapKeyIntersectionTest":               {trickySrc, "MapKeyIntersectionTest", nil, tricky.MapKeyIntersectionTest},
	"tricky/MapKeysAsSliceTest":                   {trickySrc, "MapKeysAsSliceTest", nil, tricky.MapKeysAsSliceTest},
	"tricky/MapKeySetTest":                        {trickySrc, "MapKeySetTest", nil, tricky.MapKeySetTest},
	"tricky/MapKeyShadowing":                      {trickySrc, "MapKeyShadowing", nil, tricky.MapKeyShadowing},
	"tricky/MapKeysSliceTest":                     {trickySrc, "MapKeysSliceTest", nil, tricky.MapKeysSliceTest},
	"tricky/MapKeysSorted":                        {trickySrc, "MapKeysSorted", nil, tricky.MapKeysSorted},
	"tricky/MapKeysSortedTest":                    {trickySrc, "MapKeysSortedTest", nil, tricky.MapKeysSortedTest},
	"tricky/MapKeysToSlice":                       {trickySrc, "MapKeysToSlice", nil, tricky.MapKeysToSlice},
	"tricky/MapLastVal":                           {trickySrc, "MapLastVal", nil, tricky.MapLastVal},
	"tricky/MapLookupOrInsert":                    {trickySrc, "MapLookupOrInsert", nil, tricky.MapLookupOrInsert},
	"tricky/MapMapTest":                           {trickySrc, "MapMapTest", nil, tricky.MapMapTest},
	"tricky/MapMergeConditionalTest":              {trickySrc, "MapMergeConditionalTest", nil, tricky.MapMergeConditionalTest},
	"tricky/MapMergeDisjointTest":                 {trickySrc, "MapMergeDisjointTest", nil, tricky.MapMergeDisjointTest},
	"tricky/MapMergeMultiple":                     {trickySrc, "MapMergeMultiple", nil, tricky.MapMergeMultiple},
	"tricky/MapMergeMultipleTest":                 {trickySrc, "MapMergeMultipleTest", nil, tricky.MapMergeMultipleTest},
	"tricky/MapMergeNoOverlapTest":                {trickySrc, "MapMergeNoOverlapTest", nil, tricky.MapMergeNoOverlapTest},
	"tricky/MapMergeOverwrite":                    {trickySrc, "MapMergeOverwrite", nil, tricky.MapMergeOverwrite},
	"tricky/MapMergeOverwriteAllTest":             {trickySrc, "MapMergeOverwriteAllTest", nil, tricky.MapMergeOverwriteAllTest},
	"tricky/MapMergePredTest":                     {trickySrc, "MapMergePredTest", nil, tricky.MapMergePredTest},
	"tricky/MapMergePreserveOrigTest":             {trickySrc, "MapMergePreserveOrigTest", nil, tricky.MapMergePreserveOrigTest},
	"tricky/MapMergePreserveTest":                 {trickySrc, "MapMergePreserveTest", nil, tricky.MapMergePreserveTest},
	"tricky/MapMergeSameTest":                     {trickySrc, "MapMergeSameTest", nil, tricky.MapMergeSameTest},
	"tricky/MapMergeTwo":                          {trickySrc, "MapMergeTwo", nil, tricky.MapMergeTwo},
	"tricky/MapMergeWithConflict":                 {trickySrc, "MapMergeWithConflict", nil, tricky.MapMergeWithConflict},
	"tricky/MapMergeWithConflictTest":             {trickySrc, "MapMergeWithConflictTest", nil, tricky.MapMergeWithConflictTest},
	"tricky/MapMergeWithFunc":                     {trickySrc, "MapMergeWithFunc", nil, tricky.MapMergeWithFunc},
	"tricky/MapMinMaxTest":                        {trickySrc, "MapMinMaxTest", nil, tricky.MapMinMaxTest},
	"tricky/MapNestedAssign":                      {trickySrc, "MapNestedAssign", nil, tricky.MapNestedAssign},
	"tricky/MapNestedDelete":                      {trickySrc, "MapNestedDelete", nil, tricky.MapNestedDelete},
	"tricky/MapNestedUpdate":                      {trickySrc, "MapNestedUpdate", nil, tricky.MapNestedUpdate},
	"tricky/MapNoneMatch":                         {trickySrc, "MapNoneMatch", nil, tricky.MapNoneMatch},
	"tricky/MapPartition":                         {trickySrc, "MapPartition", nil, tricky.MapPartition},
	"tricky/MapPick":                              {trickySrc, "MapPick", nil, tricky.MapPick},
	"tricky/MapPickBy":                            {trickySrc, "MapPickBy", nil, tricky.MapPickBy},
	"tricky/MapPluck":                             {trickySrc, "MapPluck", nil, tricky.MapPluck},
	"tricky/MapRangeBreak":                        {trickySrc, "MapRangeBreak", nil, tricky.MapRangeBreak},
	"tricky/MapRangeSafe":                         {trickySrc, "MapRangeSafe", nil, tricky.MapRangeSafe},
	// Note: tricky/MapRangeWithBreak removed - non-deterministic map iteration order
	"tricky/MapRejectKeys":           {trickySrc, "MapRejectKeys", nil, tricky.MapRejectKeys},
	"tricky/MapRemoveKeysTest":       {trickySrc, "MapRemoveKeysTest", nil, tricky.MapRemoveKeysTest},
	"tricky/MapReplace":              {trickySrc, "MapReplace", nil, tricky.MapReplace},
	"tricky/MapReplaceVals":          {trickySrc, "MapReplaceVals", nil, tricky.MapReplaceVals},
	"tricky/MapSameKeyValueTest":     {trickySrc, "MapSameKeyValueTest", nil, tricky.MapSameKeyValueTest},
	"tricky/MapSelectKeys":           {trickySrc, "MapSelectKeys", nil, tricky.MapSelectKeys},
	"tricky/MapSelectTest":           {trickySrc, "MapSelectTest", nil, tricky.MapSelectTest},
	"tricky/MapSize":                 {trickySrc, "MapSize", nil, tricky.MapSize},
	"tricky/MapSizeHint":             {trickySrc, "MapSizeHint", nil, tricky.MapSizeHint},
	"tricky/MapSizeTest":             {trickySrc, "MapSizeTest", nil, tricky.MapSizeTest},
	"tricky/MapSliceKeys":            {trickySrc, "MapSliceKeys", nil, tricky.MapSliceKeys},
	"tricky/MapSliceKeysTest":        {trickySrc, "MapSliceKeysTest", nil, tricky.MapSliceKeysTest},
	"tricky/MapSliceMap":             {trickySrc, "MapSliceMap", nil, tricky.MapSliceMap},
	"tricky/MapSliceReduce":          {trickySrc, "MapSliceReduce", nil, tricky.MapSliceReduce},
	"tricky/MapSliceToMap":           {trickySrc, "MapSliceToMap", nil, tricky.MapSliceToMap},
	"tricky/MapSliceToMapTest":       {trickySrc, "MapSliceToMapTest", nil, tricky.MapSliceToMapTest},
	"tricky/MapSliceValues":          {trickySrc, "MapSliceValues", nil, tricky.MapSliceValues},
	"tricky/MapSplitTest":            {trickySrc, "MapSplitTest", nil, tricky.MapSplitTest},
	"tricky/MapStructUpdate":         {trickySrc, "MapStructUpdate", nil, tricky.MapStructUpdate},
	"tricky/MapSumTest":              {trickySrc, "MapSumTest", nil, tricky.MapSumTest},
	"tricky/MapSumVals":              {trickySrc, "MapSumVals", nil, tricky.MapSumVals},
	"tricky/MapSwap":                 {trickySrc, "MapSwap", nil, tricky.MapSwap},
	"tricky/MapSymDiffTest":          {trickySrc, "MapSymDiffTest", nil, tricky.MapSymDiffTest},
	"tricky/MapTakeTest":             {trickySrc, "MapTakeTest", nil, tricky.MapTakeTest},
	"tricky/MapTakeWhileTest":        {trickySrc, "MapTakeWhileTest", nil, tricky.MapTakeWhileTest},
	"tricky/MapTally":                {trickySrc, "MapTally", nil, tricky.MapTally},
	"tricky/MapToSlice":              {trickySrc, "MapToSlice", nil, tricky.MapToSlice},
	"tricky/MapToSliceTest":          {trickySrc, "MapToSliceTest", nil, tricky.MapToSliceTest},
	"tricky/MapTransformKeys":        {trickySrc, "MapTransformKeys", nil, tricky.MapTransformKeys},
	"tricky/MapTransformKeysToSlice": {trickySrc, "MapTransformKeysToSlice", nil, tricky.MapTransformKeysToSlice},
	"tricky/MapTransformVals":        {trickySrc, "MapTransformVals", nil, tricky.MapTransformVals},
	"tricky/MapTransposeTest":        {trickySrc, "MapTransposeTest", nil, tricky.MapTransposeTest},
	"tricky/MapTwoKeys":              {trickySrc, "MapTwoKeys", nil, tricky.MapTwoKeys},
	"tricky/MapUnion":                {trickySrc, "MapUnion", nil, tricky.MapUnion},
	"tricky/MapUnionKeysTest":        {trickySrc, "MapUnionKeysTest", nil, tricky.MapUnionKeysTest},
	// Note: tricky/MapUpdateDuringRange removed - non-deterministic map iteration order
	"tricky/MapUpdateExistingTest":       {trickySrc, "MapUpdateExistingTest", nil, tricky.MapUpdateExistingTest},
	"tricky/MapUpdateIfFunc":             {trickySrc, "MapUpdateIfFunc", nil, tricky.MapUpdateIfFunc},
	"tricky/MapUpdateIfKeyExistsTest":    {trickySrc, "MapUpdateIfKeyExistsTest", nil, tricky.MapUpdateIfKeyExistsTest},
	"tricky/MapUpdateIfTest":             {trickySrc, "MapUpdateIfTest", nil, tricky.MapUpdateIfTest},
	"tricky/MapUpdateNestedMapTest":      {trickySrc, "MapUpdateNestedMapTest", nil, tricky.MapUpdateNestedMapTest},
	"tricky/MapUpdateNestedTest":         {trickySrc, "MapUpdateNestedTest", nil, tricky.MapUpdateNestedTest},
	"tricky/MapUpdateValueDirect":        {trickySrc, "MapUpdateValueDirect", nil, tricky.MapUpdateValueDirect},
	"tricky/MapUpdateWithFunc":           {trickySrc, "MapUpdateWithFunc", nil, tricky.MapUpdateWithFunc},
	"tricky/MapVals":                     {trickySrc, "MapVals", nil, tricky.MapVals},
	"tricky/MapValueDiffTest":            {trickySrc, "MapValueDiffTest", nil, tricky.MapValueDiffTest},
	"tricky/MapValueExistsTest":          {trickySrc, "MapValueExistsTest", nil, tricky.MapValueExistsTest},
	"tricky/MapValueMaxTest":             {trickySrc, "MapValueMaxTest", nil, tricky.MapValueMaxTest},
	"tricky/MapValueSlice":               {trickySrc, "MapValueSlice", nil, tricky.MapValueSlice},
	"tricky/MapValueSliceTest":           {trickySrc, "MapValueSliceTest", nil, tricky.MapValueSliceTest},
	"tricky/MapValuesToSlice":            {trickySrc, "MapValuesToSlice", nil, tricky.MapValuesToSlice},
	"tricky/MapValueSumKeysTest":         {trickySrc, "MapValueSumKeysTest", nil, tricky.MapValueSumKeysTest},
	"tricky/MapValueTypes":               {trickySrc, "MapValueTypes", nil, tricky.MapValueTypes},
	"tricky/MapWithFuncValue":            {trickySrc, "MapWithFuncValue", nil, tricky.MapWithFuncValue},
	"tricky/MapWithFuncValueDirect":      {trickySrc, "MapWithFuncValueDirect", nil, tricky.MapWithFuncValueDirect},
	"tricky/MapWithPointerValue":         {trickySrc, "MapWithPointerValue", nil, tricky.MapWithPointerValue},
	"tricky/MapZip":                      {trickySrc, "MapZip", nil, tricky.MapZip},
	"tricky/MultipleNamedReturnCombined": {trickySrc, "MultipleNamedReturnCombined", nil, tricky.MultipleNamedReturnCombined},
	"tricky/NestedClosureWithArg":        {trickySrc, "NestedClosureWithArg", nil, tricky.NestedClosureWithArg},
	"tricky/NestedMapDeleteNested":       {trickySrc, "NestedMapDeleteNested", nil, tricky.NestedMapDeleteNested},
	"tricky/NestedMapGetOrSet":           {trickySrc, "NestedMapGetOrSet", nil, tricky.NestedMapGetOrSet},
	"tricky/NestedMapGetWithDefault":     {trickySrc, "NestedMapGetWithDefault", nil, tricky.NestedMapGetWithDefault},
	"tricky/NestedMapInit":               {trickySrc, "NestedMapInit", nil, tricky.NestedMapInit},
	"tricky/NestedMapIterate":            {trickySrc, "NestedMapIterate", nil, tricky.NestedMapIterate},
	"tricky/NestedMapSafeAccess":         {trickySrc, "NestedMapSafeAccess", nil, tricky.NestedMapSafeAccess},
	"tricky/NestedMapUpdateNested":       {trickySrc, "NestedMapUpdateNested", nil, tricky.NestedMapUpdateNested},
	"tricky/NestedMapWithDelete":         {trickySrc, "NestedMapWithDelete", nil, tricky.NestedMapWithDelete},
	"tricky/NestedStructAssign":          {trickySrc, "NestedStructAssign", nil, tricky.NestedStructAssign},
	// Note: tricky/NextPermutation removed - needs slice argument
	"tricky/PointerAddr":                        {trickySrc, "PointerAddr", nil, tricky.PointerAddr},
	"tricky/PointerAlias":                       {trickySrc, "PointerAlias", nil, tricky.PointerAlias},
	"tricky/PointerArithSim":                    {trickySrc, "PointerArithSim", nil, tricky.PointerArithSim},
	"tricky/PointerArrayElementTest":            {trickySrc, "PointerArrayElementTest", nil, tricky.PointerArrayElementTest},
	"tricky/PointerArrayIdx":                    {trickySrc, "PointerArrayIdx", nil, tricky.PointerArrayIdx},
	"tricky/PointerArrayIndexTest":              {trickySrc, "PointerArrayIndexTest", nil, tricky.PointerArrayIndexTest},
	"tricky/PointerAssignChainTest":             {trickySrc, "PointerAssignChainTest", nil, tricky.PointerAssignChainTest},
	"tricky/PointerAssignFromDerefTest":         {trickySrc, "PointerAssignFromDerefTest", nil, tricky.PointerAssignFromDerefTest},
	"tricky/PointerAssignFromFuncTest":          {trickySrc, "PointerAssignFromFuncTest", nil, tricky.PointerAssignFromFuncTest},
	"tricky/PointerAssignFuncResultTest":        {trickySrc, "PointerAssignFuncResultTest", nil, tricky.PointerAssignFuncResultTest},
	"tricky/PointerAssignNilTest":               {trickySrc, "PointerAssignNilTest", nil, tricky.PointerAssignNilTest},
	"tricky/PointerAssignSameTest":              {trickySrc, "PointerAssignSameTest", nil, tricky.PointerAssignSameTest},
	"tricky/PointerAssignThenNilTest":           {trickySrc, "PointerAssignThenNilTest", nil, tricky.PointerAssignThenNilTest},
	"tricky/PointerChainTest":                   {trickySrc, "PointerChainTest", nil, tricky.PointerChainTest},
	"tricky/PointerCheckNilAfterUseTest":        {trickySrc, "PointerCheckNilAfterUseTest", nil, tricky.PointerCheckNilAfterUseTest},
	"tricky/PointerCompare":                     {trickySrc, "PointerCompare", nil, tricky.PointerCompare},
	"tricky/PointerCompareDiffTest":             {trickySrc, "PointerCompareDiffTest", nil, tricky.PointerCompareDiffTest},
	"tricky/PointerCompareTest":                 {trickySrc, "PointerCompareTest", nil, tricky.PointerCompareTest},
	"tricky/PointerDeref":                       {trickySrc, "PointerDeref", nil, tricky.PointerDeref},
	"tricky/PointerDerefAssignTest":             {trickySrc, "PointerDerefAssignTest", nil, tricky.PointerDerefAssignTest},
	"tricky/PointerDerefChain":                  {trickySrc, "PointerDerefChain", nil, tricky.PointerDerefChain},
	"tricky/PointerDerefChainTest":              {trickySrc, "PointerDerefChainTest", nil, tricky.PointerDerefChainTest},
	"tricky/PointerDerefModifyTest":             {trickySrc, "PointerDerefModifyTest", nil, tricky.PointerDerefModifyTest},
	"tricky/PointerDerefNilCheckTest":           {trickySrc, "PointerDerefNilCheckTest", nil, tricky.PointerDerefNilCheckTest},
	"tricky/PointerDerefNilTest":                {trickySrc, "PointerDerefNilTest", nil, tricky.PointerDerefNilTest},
	"tricky/PointerDoubleAssignTest":            {trickySrc, "PointerDoubleAssignTest", nil, tricky.PointerDoubleAssignTest},
	"tricky/PointerDoubleDerefTest":             {trickySrc, "PointerDoubleDerefTest", nil, tricky.PointerDoubleDerefTest},
	"tricky/PointerLevel":                       {trickySrc, "PointerLevel", nil, tricky.PointerLevel},
	"tricky/PointerLevelTest":                   {trickySrc, "PointerLevelTest", nil, tricky.PointerLevelTest},
	"tricky/PointerNilAssign":                   {trickySrc, "PointerNilAssign", nil, tricky.PointerNilAssign},
	"tricky/PointerNilAssignAfterUseTest":       {trickySrc, "PointerNilAssignAfterUseTest", nil, tricky.PointerNilAssignAfterUseTest},
	"tricky/PointerNilAssignT":                  {trickySrc, "PointerNilAssignT", nil, tricky.PointerNilAssignT},
	"tricky/PointerNilCheckAfterAssignTest":     {trickySrc, "PointerNilCheckAfterAssignTest", nil, tricky.PointerNilCheckAfterAssignTest},
	"tricky/PointerNilCheckChain":               {trickySrc, "PointerNilCheckChain", nil, tricky.PointerNilCheckChain},
	"tricky/PointerNilCheckDerefTest":           {trickySrc, "PointerNilCheckDerefTest", nil, tricky.PointerNilCheckDerefTest},
	"tricky/PointerNilCompare":                  {trickySrc, "PointerNilCompare", nil, tricky.PointerNilCompare},
	"tricky/PointerNilCompareTest":              {trickySrc, "PointerNilCompareTest", nil, tricky.PointerNilCompareTest},
	"tricky/PointerNilDeref":                    {trickySrc, "PointerNilDeref", nil, tricky.PointerNilDeref},
	"tricky/PointerNilReassign":                 {trickySrc, "PointerNilReassign", nil, tricky.PointerNilReassign},
	"tricky/PointerNilSafe":                     {trickySrc, "PointerNilSafe", nil, tricky.PointerNilSafe},
	"tricky/PointerNilSafeDeref":                {trickySrc, "PointerNilSafeDeref", nil, tricky.PointerNilSafeDeref},
	"tricky/PointerNilSafeDerefTest":            {trickySrc, "PointerNilSafeDerefTest", nil, tricky.PointerNilSafeDerefTest},
	"tricky/PointerNilSafeOpTest":               {trickySrc, "PointerNilSafeOpTest", nil, tricky.PointerNilSafeOpTest},
	"tricky/PointerNilThenAssignTest":           {trickySrc, "PointerNilThenAssignTest", nil, tricky.PointerNilThenAssignTest},
	"tricky/PointerNullObject":                  {trickySrc, "PointerNullObject", nil, tricky.PointerNullObject},
	"tricky/PointerReassignChainTest":           {trickySrc, "PointerReassignChainTest", nil, tricky.PointerReassignChainTest},
	"tricky/PointerReassignmentChain":           {trickySrc, "PointerReassignmentChain", nil, tricky.PointerReassignmentChain},
	"tricky/PointerReassignNil":                 {trickySrc, "PointerReassignNil", nil, tricky.PointerReassignNil},
	"tricky/PointerReassignNilTest":             {trickySrc, "PointerReassignNilTest", nil, tricky.PointerReassignNilTest},
	"tricky/PointerReassignTest":                {trickySrc, "PointerReassignTest", nil, tricky.PointerReassignTest},
	"tricky/PointerRotate":                      {trickySrc, "PointerRotate", nil, tricky.PointerRotate},
	"tricky/PointerSliceElementModifyTest":      {trickySrc, "PointerSliceElementModifyTest", nil, tricky.PointerSliceElementModifyTest},
	"tricky/PointerSliceElementSwap":            {trickySrc, "PointerSliceElementSwap", nil, tricky.PointerSliceElementSwap},
	"tricky/PointerSliceIndexTest":              {trickySrc, "PointerSliceIndexTest", nil, tricky.PointerSliceIndexTest},
	"tricky/PointerSliceIterateTest":            {trickySrc, "PointerSliceIterateTest", nil, tricky.PointerSliceIterateTest},
	"tricky/PointerSliceLenTest":                {trickySrc, "PointerSliceLenTest", nil, tricky.PointerSliceLenTest},
	"tricky/PointerSliceModifyTest":             {trickySrc, "PointerSliceModifyTest", nil, tricky.PointerSliceModifyTest},
	"tricky/PointerSliceNilTest":                {trickySrc, "PointerSliceNilTest", nil, tricky.PointerSliceNilTest},
	"tricky/PointerSliceOfPointers":             {trickySrc, "PointerSliceOfPointers", nil, tricky.PointerSliceOfPointers},
	"tricky/PointerSliceOfStructTest":           {trickySrc, "PointerSliceOfStructTest", nil, tricky.PointerSliceOfStructTest},
	"tricky/PointerStructFieldNilCheckTest":     {trickySrc, "PointerStructFieldNilCheckTest", nil, tricky.PointerStructFieldNilCheckTest},
	"tricky/PointerStructFieldNilTest":          {trickySrc, "PointerStructFieldNilTest", nil, tricky.PointerStructFieldNilTest},
	"tricky/PointerStructFieldTest":             {trickySrc, "PointerStructFieldTest", nil, tricky.PointerStructFieldTest},
	"tricky/PointerStructFld":                   {trickySrc, "PointerStructFld", nil, tricky.PointerStructFld},
	"tricky/PointerStructMethodTest":            {trickySrc, "PointerStructMethodTest", nil, tricky.PointerStructMethodTest},
	"tricky/PointerStructModifyFieldTest":       {trickySrc, "PointerStructModifyFieldTest", nil, tricky.PointerStructModifyFieldTest},
	"tricky/PointerStructModifyTest":            {trickySrc, "PointerStructModifyTest", nil, tricky.PointerStructModifyTest},
	"tricky/PointerSwap":                        {trickySrc, "PointerSwap", nil, tricky.PointerSwap},
	"tricky/PointerSwapChain":                   {trickySrc, "PointerSwapChain", nil, tricky.PointerSwapChain},
	"tricky/PointerSwapChainTest":               {trickySrc, "PointerSwapChainTest", nil, tricky.PointerSwapChainTest},
	"tricky/PointerSwapInArrayTest":             {trickySrc, "PointerSwapInArrayTest", nil, tricky.PointerSwapInArrayTest},
	"tricky/PointerSwapInSlice":                 {trickySrc, "PointerSwapInSlice", nil, tricky.PointerSwapInSlice},
	"tricky/PointerSwapInStruct":                {trickySrc, "PointerSwapInStruct", nil, tricky.PointerSwapInStruct},
	"tricky/PointerSwapInStructTest":            {trickySrc, "PointerSwapInStructTest", nil, tricky.PointerSwapInStructTest},
	"tricky/PointerSwapMultipleTest":            {trickySrc, "PointerSwapMultipleTest", nil, tricky.PointerSwapMultipleTest},
	"tricky/PointerSwapNilSafe":                 {trickySrc, "PointerSwapNilSafe", nil, tricky.PointerSwapNilSafe},
	"tricky/PointerSwapSimple":                  {trickySrc, "PointerSwapSimple", nil, tricky.PointerSwapSimple},
	"tricky/PointerSwapStructFieldsTest":        {trickySrc, "PointerSwapStructFieldsTest", nil, tricky.PointerSwapStructFieldsTest},
	"tricky/PointerSwapThroughSliceTest":        {trickySrc, "PointerSwapThroughSliceTest", nil, tricky.PointerSwapThroughSliceTest},
	"tricky/PointerSwapVals":                    {trickySrc, "PointerSwapVals", nil, tricky.PointerSwapVals},
	"tricky/PointerSwapValues":                  {trickySrc, "PointerSwapValues", nil, tricky.PointerSwapValues},
	"tricky/PointerSwapValuesTest":              {trickySrc, "PointerSwapValuesTest", nil, tricky.PointerSwapValuesTest},
	"tricky/PointerSwapViaSliceTest":            {trickySrc, "PointerSwapViaSliceTest", nil, tricky.PointerSwapViaSliceTest},
	"tricky/PointerSwapViaTempTest":             {trickySrc, "PointerSwapViaTempTest", nil, tricky.PointerSwapViaTempTest},
	"tricky/PointerToArr":                       {trickySrc, "PointerToArr", nil, tricky.PointerToArr},
	"tricky/PointerToArray":                     {trickySrc, "PointerToArray", nil, tricky.PointerToArray},
	"tricky/PointerToArrayElement":              {trickySrc, "PointerToArrayElement", nil, tricky.PointerToArrayElement},
	"tricky/PointerToArrTest":                   {trickySrc, "PointerToArrTest", nil, tricky.PointerToArrTest},
	"tricky/PointerToChanTest":                  {trickySrc, "PointerToChanTest", nil, tricky.PointerToChanTest},
	"tricky/PointerToFunc":                      {trickySrc, "PointerToFunc", nil, tricky.PointerToFunc},
	"tricky/PointerToFuncResultTest":            {trickySrc, "PointerToFuncResultTest", nil, tricky.PointerToFuncResultTest},
	"tricky/PointerToInterface":                 {trickySrc, "PointerToInterface", nil, tricky.PointerToInterface},
	"tricky/PointerToMapElement":                {trickySrc, "PointerToMapElement", nil, tricky.PointerToMapElement},
	"tricky/PointerToMapKey":                    {trickySrc, "PointerToMapKey", nil, tricky.PointerToMapKey},
	"tricky/PointerToMapNilTest":                {trickySrc, "PointerToMapNilTest", nil, tricky.PointerToMapNilTest},
	"tricky/PointerToMapTest":                   {trickySrc, "PointerToMapTest", nil, tricky.PointerToMapTest},
	"tricky/PointerToNilAssignTest":             {trickySrc, "PointerToNilAssignTest", nil, tricky.PointerToNilAssignTest},
	"tricky/PointerToNilInterface":              {trickySrc, "PointerToNilInterface", nil, tricky.PointerToNilInterface},
	"tricky/PointerToNilMap":                    {trickySrc, "PointerToNilMap", nil, tricky.PointerToNilMap},
	"tricky/PointerToNilMapLenTest":             {trickySrc, "PointerToNilMapLenTest", nil, tricky.PointerToNilMapLenTest},
	"tricky/PointerToNilSliceLenTest":           {trickySrc, "PointerToNilSliceLenTest", nil, tricky.PointerToNilSliceLenTest},
	"tricky/PointerToNilSliceTest":              {trickySrc, "PointerToNilSliceTest", nil, tricky.PointerToNilSliceTest},
	"tricky/PointerToNilStruct":                 {trickySrc, "PointerToNilStruct", nil, tricky.PointerToNilStruct},
	"tricky/PointerToNilStructTest":             {trickySrc, "PointerToNilStructTest", nil, tricky.PointerToNilStructTest},
	"tricky/PointerToNilTest":                   {trickySrc, "PointerToNilTest", nil, tricky.PointerToNilTest},
	"tricky/PointerToPointer":                   {trickySrc, "PointerToPointer", nil, tricky.PointerToPointer},
	"tricky/PointerToPointerAssign":             {trickySrc, "PointerToPointerAssign", nil, tricky.PointerToPointerAssign},
	"tricky/PointerToPointerAssignTest":         {trickySrc, "PointerToPointerAssignTest", nil, tricky.PointerToPointerAssignTest},
	"tricky/PointerToPointerDerefTest":          {trickySrc, "PointerToPointerDerefTest", nil, tricky.PointerToPointerDerefTest},
	"tricky/PointerToSliceAppend":               {trickySrc, "PointerToSliceAppend", nil, tricky.PointerToSliceAppend},
	"tricky/PointerToSliceClear":                {trickySrc, "PointerToSliceClear", nil, tricky.PointerToSliceClear},
	"tricky/PointerToSliceClearTest":            {trickySrc, "PointerToSliceClearTest", nil, tricky.PointerToSliceClearTest},
	"tricky/PointerToSliceElementModifyTest":    {trickySrc, "PointerToSliceElementModifyTest", nil, tricky.PointerToSliceElementModifyTest},
	"tricky/PointerToSliceLen":                  {trickySrc, "PointerToSliceLen", nil, tricky.PointerToSliceLen},
	"tricky/PointerToSliceLenCap":               {trickySrc, "PointerToSliceLenCap", nil, tricky.PointerToSliceLenCap},
	"tricky/PointerToSliceModify":               {trickySrc, "PointerToSliceModify", nil, tricky.PointerToSliceModify},
	"tricky/PointerToSliceNilTest":              {trickySrc, "PointerToSliceNilTest", nil, tricky.PointerToSliceNilTest},
	"tricky/PointerToSliceOfNilTest":            {trickySrc, "PointerToSliceOfNilTest", nil, tricky.PointerToSliceOfNilTest},
	"tricky/PointerToSliceOfPtrTest":            {trickySrc, "PointerToSliceOfPtrTest", nil, tricky.PointerToSliceOfPtrTest},
	"tricky/PointerToSliceOfStructs":            {trickySrc, "PointerToSliceOfStructs", nil, tricky.PointerToSliceOfStructs},
	"tricky/PointerToSliceTest":                 {trickySrc, "PointerToSliceTest", nil, tricky.PointerToSliceTest},
	"tricky/PointerToStructField":               {trickySrc, "PointerToStructField", nil, tricky.PointerToStructField},
	"tricky/PointerToStructMethodTest":          {trickySrc, "PointerToStructMethodTest", nil, tricky.PointerToStructMethodTest},
	"tricky/PointerToStructNilMethodTest":       {trickySrc, "PointerToStructNilMethodTest", nil, tricky.PointerToStructNilMethodTest},
	"tricky/PointerToStructTest":                {trickySrc, "PointerToStructTest", nil, tricky.PointerToStructTest},
	"tricky/SliceAll":                           {trickySrc, "SliceAll", nil, tricky.SliceAll},
	"tricky/SliceAppendCapTest":                 {trickySrc, "SliceAppendCapTest", nil, tricky.SliceAppendCapTest},
	"tricky/SliceAppendFunc":                    {trickySrc, "SliceAppendFunc", nil, tricky.SliceAppendFunc},
	"tricky/SliceAppendIfTest":                  {trickySrc, "SliceAppendIfTest", nil, tricky.SliceAppendIfTest},
	"tricky/SliceAppendNilTest":                 {trickySrc, "SliceAppendNilTest", nil, tricky.SliceAppendNilTest},
	"tricky/SliceAppendSliceTest":               {trickySrc, "SliceAppendSliceTest", nil, tricky.SliceAppendSliceTest},
	"tricky/SliceBsearch":                       {trickySrc, "SliceBsearch", nil, tricky.SliceBsearch},
	"tricky/SliceCartesianProduct":              {trickySrc, "SliceCartesianProduct", nil, tricky.SliceCartesianProduct},
	"tricky/SliceChainedSlice":                  {trickySrc, "SliceChainedSlice", nil, tricky.SliceChainedSlice},
	"tricky/SliceChunk":                         {trickySrc, "SliceChunk", nil, tricky.SliceChunk},
	"tricky/SliceChunkByPredTest":               {trickySrc, "SliceChunkByPredTest", nil, tricky.SliceChunkByPredTest},
	"tricky/SliceChunkByTest":                   {trickySrc, "SliceChunkByTest", nil, tricky.SliceChunkByTest},
	"tricky/SliceChunkEveryTest":                {trickySrc, "SliceChunkEveryTest", nil, tricky.SliceChunkEveryTest},
	"tricky/SliceClone":                         {trickySrc, "SliceClone", nil, tricky.SliceClone},
	"tricky/SliceCombinations":                  {trickySrc, "SliceCombinations", nil, tricky.SliceCombinations},
	"tricky/SliceCompact":                       {trickySrc, "SliceCompact", nil, tricky.SliceCompact},
	"tricky/SliceCompactMap":                    {trickySrc, "SliceCompactMap", nil, tricky.SliceCompactMap},
	"tricky/SliceContainsAll":                   {trickySrc, "SliceContainsAll", nil, tricky.SliceContainsAll},
	"tricky/SliceContainsAllTest":               {trickySrc, "SliceContainsAllTest", nil, tricky.SliceContainsAllTest},
	"tricky/SliceContainsAnyTest":               {trickySrc, "SliceContainsAnyTest", nil, tricky.SliceContainsAnyTest},
	"tricky/SliceContainsNoneTest":              {trickySrc, "SliceContainsNoneTest", nil, tricky.SliceContainsNoneTest},
	"tricky/SliceCopyFromMap":                   {trickySrc, "SliceCopyFromMap", nil, tricky.SliceCopyFromMap},
	"tricky/SliceCopyModifyTest":                {trickySrc, "SliceCopyModifyTest", nil, tricky.SliceCopyModifyTest},
	"tricky/SliceCopyReverseTest":               {trickySrc, "SliceCopyReverseTest", nil, tricky.SliceCopyReverseTest},
	"tricky/SliceCopySubsetTest":                {trickySrc, "SliceCopySubsetTest", nil, tricky.SliceCopySubsetTest},
	"tricky/SliceCountBy":                       {trickySrc, "SliceCountBy", nil, tricky.SliceCountBy},
	"tricky/SliceCountTest":                     {trickySrc, "SliceCountTest", nil, tricky.SliceCountTest},
	"tricky/SliceCountWhileTest":                {trickySrc, "SliceCountWhileTest", nil, tricky.SliceCountWhileTest},
	"tricky/SliceCycleTest":                     {trickySrc, "SliceCycleTest", nil, tricky.SliceCycleTest},
	"tricky/SliceDedupConsecutive":              {trickySrc, "SliceDedupConsecutive", nil, tricky.SliceDedupConsecutive},
	"tricky/SliceDeleteFront":                   {trickySrc, "SliceDeleteFront", nil, tricky.SliceDeleteFront},
	"tricky/SliceDeleteMiddle":                  {trickySrc, "SliceDeleteMiddle", nil, tricky.SliceDeleteMiddle},
	"tricky/SliceDetect":                        {trickySrc, "SliceDetect", nil, tricky.SliceDetect},
	"tricky/SliceDiff":                          {trickySrc, "SliceDiff", nil, tricky.SliceDiff},
	"tricky/SliceDifference":                    {trickySrc, "SliceDifference", nil, tricky.SliceDifference},
	"tricky/SliceDifferenceBy":                  {trickySrc, "SliceDifferenceBy", nil, tricky.SliceDifferenceBy},
	"tricky/SliceDrop":                          {trickySrc, "SliceDrop", nil, tricky.SliceDrop},
	"tricky/SliceDropN":                         {trickySrc, "SliceDropN", nil, tricky.SliceDropN},
	"tricky/SliceDropNFunc":                     {trickySrc, "SliceDropNFunc", nil, tricky.SliceDropNFunc},
	"tricky/SliceDropTest":                      {trickySrc, "SliceDropTest", nil, tricky.SliceDropTest},
	"tricky/SliceDropWhile":                     {trickySrc, "SliceDropWhile", nil, tricky.SliceDropWhile},
	"tricky/SliceDropWhileTest":                 {trickySrc, "SliceDropWhileTest", nil, tricky.SliceDropWhileTest},
	"tricky/SliceEachWithIndex":                 {trickySrc, "SliceEachWithIndex", nil, tricky.SliceEachWithIndex},
	"tricky/SliceEqual":                         {trickySrc, "SliceEqual", nil, tricky.SliceEqual},
	"tricky/SliceExistsTest":                    {trickySrc, "SliceExistsTest", nil, tricky.SliceExistsTest},
	"tricky/SliceFill":                          {trickySrc, "SliceFill", nil, tricky.SliceFill},
	"tricky/SliceFilterKeepTest":                {trickySrc, "SliceFilterKeepTest", nil, tricky.SliceFilterKeepTest},
	"tricky/SliceFilterNotTest":                 {trickySrc, "SliceFilterNotTest", nil, tricky.SliceFilterNotTest},
	"tricky/SliceFindFirstFunc":                 {trickySrc, "SliceFindFirstFunc", nil, tricky.SliceFindFirstFunc},
	"tricky/SliceFindFirstTest":                 {trickySrc, "SliceFindFirstTest", nil, tricky.SliceFindFirstTest},
	"tricky/SliceFindIdx":                       {trickySrc, "SliceFindIdx", nil, tricky.SliceFindIdx},
	"tricky/SliceFindIndex":                     {trickySrc, "SliceFindIndex", nil, tricky.SliceFindIndex},
	"tricky/SliceFindIndexTest":                 {trickySrc, "SliceFindIndexTest", nil, tricky.SliceFindIndexTest},
	"tricky/SliceFindLastPosTest":               {trickySrc, "SliceFindLastPosTest", nil, tricky.SliceFindLastPosTest},
	"tricky/SliceFindLastTest":                  {trickySrc, "SliceFindLastTest", nil, tricky.SliceFindLastTest},
	"tricky/SliceFirst":                         {trickySrc, "SliceFirst", nil, tricky.SliceFirst},
	"tricky/SliceFlatten":                       {trickySrc, "SliceFlatten", nil, tricky.SliceFlatten},
	"tricky/SliceFlatten2D":                     {trickySrc, "SliceFlatten2D", nil, tricky.SliceFlatten2D},
	"tricky/SliceFlattenDeep":                   {trickySrc, "SliceFlattenDeep", nil, tricky.SliceFlattenDeep},
	"tricky/SliceFlattenLevelTest":              {trickySrc, "SliceFlattenLevelTest", nil, tricky.SliceFlattenLevelTest},
	"tricky/SliceFlattenManual":                 {trickySrc, "SliceFlattenManual", nil, tricky.SliceFlattenManual},
	"tricky/SliceFlattenManualTest":             {trickySrc, "SliceFlattenManualTest", nil, tricky.SliceFlattenManualTest},
	"tricky/SliceFoldLeft":                      {trickySrc, "SliceFoldLeft", nil, tricky.SliceFoldLeft},
	"tricky/SliceFromChan":                      {trickySrc, "SliceFromChan", nil, tricky.SliceFromChan},
	"tricky/SliceGrep":                          {trickySrc, "SliceGrep", nil, tricky.SliceGrep},
	"tricky/SliceGroupBy":                       {trickySrc, "SliceGroupBy", nil, tricky.SliceGroupBy},
	"tricky/SliceGroupByFld":                    {trickySrc, "SliceGroupByFld", nil, tricky.SliceGroupByFld},
	"tricky/SliceGroupByMultiple":               {trickySrc, "SliceGroupByMultiple", nil, tricky.SliceGroupByMultiple},
	"tricky/SliceGroupConsecutiveTest":          {trickySrc, "SliceGroupConsecutiveTest", nil, tricky.SliceGroupConsecutiveTest},
	"tricky/SliceGrow":                          {trickySrc, "SliceGrow", nil, tricky.SliceGrow},
	"tricky/SliceGrowWithAppend":                {trickySrc, "SliceGrowWithAppend", nil, tricky.SliceGrowWithAppend},
	"tricky/SliceHeadTest":                      {trickySrc, "SliceHeadTest", nil, tricky.SliceHeadTest},
	"tricky/SliceIndexOf":                       {trickySrc, "SliceIndexOf", nil, tricky.SliceIndexOf},
	"tricky/SliceIndexOfFirstTest":              {trickySrc, "SliceIndexOfFirstTest", nil, tricky.SliceIndexOfFirstTest},
	"tricky/SliceIndexOfMaxTest":                {trickySrc, "SliceIndexOfMaxTest", nil, tricky.SliceIndexOfMaxTest},
	"tricky/SliceIndexOfTest":                   {trickySrc, "SliceIndexOfTest", nil, tricky.SliceIndexOfTest},
	"tricky/SliceIndexOutOfRange":               {trickySrc, "SliceIndexOutOfRange", nil, tricky.SliceIndexOutOfRange},
	"tricky/SliceInitAndModifyTest":             {trickySrc, "SliceInitAndModifyTest", nil, tricky.SliceInitAndModifyTest},
	"tricky/SliceInitCapTest":                   {trickySrc, "SliceInitCapTest", nil, tricky.SliceInitCapTest},
	"tricky/SliceInsertAt":                      {trickySrc, "SliceInsertAt", nil, tricky.SliceInsertAt},
	"tricky/SliceInsertAtTest":                  {trickySrc, "SliceInsertAtTest", nil, tricky.SliceInsertAtTest},
	"tricky/SliceInsertFrontTest":               {trickySrc, "SliceInsertFrontTest", nil, tricky.SliceInsertFrontTest},
	"tricky/SliceInsertMultipleTest":            {trickySrc, "SliceInsertMultipleTest", nil, tricky.SliceInsertMultipleTest},
	"tricky/SliceInsertSliceTest":               {trickySrc, "SliceInsertSliceTest", nil, tricky.SliceInsertSliceTest},
	"tricky/SliceInterleaveTest":                {trickySrc, "SliceInterleaveTest", nil, tricky.SliceInterleaveTest},
	"tricky/SliceIntersect":                     {trickySrc, "SliceIntersect", nil, tricky.SliceIntersect},
	"tricky/SliceIntersectBy":                   {trickySrc, "SliceIntersectBy", nil, tricky.SliceIntersectBy},
	"tricky/SliceIntersectTest":                 {trickySrc, "SliceIntersectTest", nil, tricky.SliceIntersectTest},
	"tricky/SliceIntersperseTest":               {trickySrc, "SliceIntersperseTest", nil, tricky.SliceIntersperseTest},
	"tricky/SliceIsSortedTest":                  {trickySrc, "SliceIsSortedTest", nil, tricky.SliceIsSortedTest},
	"tricky/SliceLast":                          {trickySrc, "SliceLast", nil, tricky.SliceLast},
	"tricky/SliceLastIndexOfTest":               {trickySrc, "SliceLastIndexOfTest", nil, tricky.SliceLastIndexOfTest},
	"tricky/SliceLastIndexOfTest2":              {trickySrc, "SliceLastIndexOfTest2", nil, tricky.SliceLastIndexOfTest2},
	"tricky/SliceLastNFunc":                     {trickySrc, "SliceLastNFunc", nil, tricky.SliceLastNFunc},
	"tricky/SliceMakeFromArr":                   {trickySrc, "SliceMakeFromArr", nil, tricky.SliceMakeFromArr},
	"tricky/SliceMakeZero":                      {trickySrc, "SliceMakeZero", nil, tricky.SliceMakeZero},
	"tricky/SliceMapEachTest":                   {trickySrc, "SliceMapEachTest", nil, tricky.SliceMapEachTest},
	"tricky/SliceMapIndex":                      {trickySrc, "SliceMapIndex", nil, tricky.SliceMapIndex},
	"tricky/SliceMax":                           {trickySrc, "SliceMax", nil, tricky.SliceMax},
	"tricky/SliceMaxVal":                        {trickySrc, "SliceMaxVal", nil, tricky.SliceMaxVal},
	"tricky/SliceMinIdx":                        {trickySrc, "SliceMinIdx", nil, tricky.SliceMinIdx},
	"tricky/SliceMinMax":                        {trickySrc, "SliceMinMax", nil, tricky.SliceMinMax},
	"tricky/SliceMinVal":                        {trickySrc, "SliceMinVal", nil, tricky.SliceMinVal},
	"tricky/SliceNegativeIndex":                 {trickySrc, "SliceNegativeIndex", nil, tricky.SliceNegativeIndex},
	"tricky/SliceNilAppend":                     {trickySrc, "SliceNilAppend", nil, tricky.SliceNilAppend},
	"tricky/SliceNone":                          {trickySrc, "SliceNone", nil, tricky.SliceNone},
	"tricky/SliceOfEmptyInterface":              {trickySrc, "SliceOfEmptyInterface", nil, tricky.SliceOfEmptyInterface},
	"tricky/SliceOfInterfacesWithTypes":         {trickySrc, "SliceOfInterfacesWithTypes", nil, tricky.SliceOfInterfacesWithTypes},
	"tricky/SlicePad":                           {trickySrc, "SlicePad", nil, tricky.SlicePad},
	"tricky/SlicePadLeftTest":                   {trickySrc, "SlicePadLeftTest", nil, tricky.SlicePadLeftTest},
	"tricky/SlicePadRightTest":                  {trickySrc, "SlicePadRightTest", nil, tricky.SlicePadRightTest},
	"tricky/SlicePartition":                     {trickySrc, "SlicePartition", nil, tricky.SlicePartition},
	"tricky/SlicePartitionBy":                   {trickySrc, "SlicePartitionBy", nil, tricky.SlicePartitionBy},
	"tricky/SlicePartitionPosNeg":               {trickySrc, "SlicePartitionPosNeg", nil, tricky.SlicePartitionPosNeg},
	"tricky/SlicePartitionTest":                 {trickySrc, "SlicePartitionTest", nil, tricky.SlicePartitionTest},
	"tricky/SlicePermutation":                   {trickySrc, "SlicePermutation", nil, tricky.SlicePermutation},
	"tricky/SlicePermuteSimpleTest":             {trickySrc, "SlicePermuteSimpleTest", nil, tricky.SlicePermuteSimpleTest},
	"tricky/SlicePluck":                         {trickySrc, "SlicePluck", nil, tricky.SlicePluck},
	"tricky/SlicePluckFld":                      {trickySrc, "SlicePluckFld", nil, tricky.SlicePluckFld},
	"tricky/SlicePrepend":                       {trickySrc, "SlicePrepend", nil, tricky.SlicePrepend},
	"tricky/SlicePrependMultipleTest":           {trickySrc, "SlicePrependMultipleTest", nil, tricky.SlicePrependMultipleTest},
	"tricky/SlicePrependValueTest":              {trickySrc, "SlicePrependValueTest", nil, tricky.SlicePrependValueTest},
	"tricky/SliceProd":                          {trickySrc, "SliceProd", nil, tricky.SliceProd},
	"tricky/SliceProduct":                       {trickySrc, "SliceProduct", nil, tricky.SliceProduct},
	"tricky/SliceRandomAccess":                  {trickySrc, "SliceRandomAccess", nil, tricky.SliceRandomAccess},
	"tricky/SliceReduce":                        {trickySrc, "SliceReduce", nil, tricky.SliceReduce},
	"tricky/SliceReduceTest":                    {trickySrc, "SliceReduceTest", nil, tricky.SliceReduceTest},
	"tricky/SliceReject":                        {trickySrc, "SliceReject", nil, tricky.SliceReject},
	"tricky/SliceRemoveAtTest":                  {trickySrc, "SliceRemoveAtTest", nil, tricky.SliceRemoveAtTest},
	"tricky/SliceRemoveDupes":                   {trickySrc, "SliceRemoveDupes", nil, tricky.SliceRemoveDupes},
	"tricky/SliceRemoveDupSortedTest":           {trickySrc, "SliceRemoveDupSortedTest", nil, tricky.SliceRemoveDupSortedTest},
	"tricky/SliceRemoveDupTest":                 {trickySrc, "SliceRemoveDupTest", nil, tricky.SliceRemoveDupTest},
	"tricky/SliceRemoveIf":                      {trickySrc, "SliceRemoveIf", nil, tricky.SliceRemoveIf},
	"tricky/SliceRemoveIfKeepTest":              {trickySrc, "SliceRemoveIfKeepTest", nil, tricky.SliceRemoveIfKeepTest},
	"tricky/SliceRemoveIfTest":                  {trickySrc, "SliceRemoveIfTest", nil, tricky.SliceRemoveIfTest},
	"tricky/SliceRemoveLastTest":                {trickySrc, "SliceRemoveLastTest", nil, tricky.SliceRemoveLastTest},
	"tricky/SliceRepeat":                        {trickySrc, "SliceRepeat", nil, tricky.SliceRepeat},
	"tricky/SliceRepeatNTest":                   {trickySrc, "SliceRepeatNTest", nil, tricky.SliceRepeatNTest},
	"tricky/SliceReplaceAtTest":                 {trickySrc, "SliceReplaceAtTest", nil, tricky.SliceReplaceAtTest},
	"tricky/SliceReverseCopy":                   {trickySrc, "SliceReverseCopy", nil, tricky.SliceReverseCopy},
	"tricky/SliceReverseCopyTest":               {trickySrc, "SliceReverseCopyTest", nil, tricky.SliceReverseCopyTest},
	"tricky/SliceReverseInPlace":                {trickySrc, "SliceReverseInPlace", nil, tricky.SliceReverseInPlace},
	"tricky/SliceReverseManualTest":             {trickySrc, "SliceReverseManualTest", nil, tricky.SliceReverseManualTest},
	"tricky/SliceReverseRangeTest":              {trickySrc, "SliceReverseRangeTest", nil, tricky.SliceReverseRangeTest},
	"tricky/SliceRotate":                        {trickySrc, "SliceRotate", nil, tricky.SliceRotate},
	"tricky/SliceRotateByTest":                  {trickySrc, "SliceRotateByTest", nil, tricky.SliceRotateByTest},
	"tricky/SliceRotateLeft":                    {trickySrc, "SliceRotateLeft", nil, tricky.SliceRotateLeft},
	"tricky/SliceRotateLeftNTest":               {trickySrc, "SliceRotateLeftNTest", nil, tricky.SliceRotateLeftNTest},
	"tricky/SliceRotateLeftTest":                {trickySrc, "SliceRotateLeftTest", nil, tricky.SliceRotateLeftTest},
	"tricky/SliceRotateRight":                   {trickySrc, "SliceRotateRight", nil, tricky.SliceRotateRight},
	"tricky/SliceRotateRightNTest":              {trickySrc, "SliceRotateRightNTest", nil, tricky.SliceRotateRightNTest},
	"tricky/SliceRotateRightTest":               {trickySrc, "SliceRotateRightTest", nil, tricky.SliceRotateRightTest},
	"tricky/SliceSample":                        {trickySrc, "SliceSample", nil, tricky.SliceSample},
	"tricky/SliceScan":                          {trickySrc, "SliceScan", nil, tricky.SliceScan},
	"tricky/SliceScanLeftTest":                  {trickySrc, "SliceScanLeftTest", nil, tricky.SliceScanLeftTest},
	"tricky/SliceSelect":                        {trickySrc, "SliceSelect", nil, tricky.SliceSelect},
	"tricky/SliceShiftLeftTest":                 {trickySrc, "SliceShiftLeftTest", nil, tricky.SliceShiftLeftTest},
	"tricky/SliceSlideTest":                     {trickySrc, "SliceSlideTest", nil, tricky.SliceSlideTest},
	"tricky/SliceSlidingWindowTest":             {trickySrc, "SliceSlidingWindowTest", nil, tricky.SliceSlidingWindowTest},
	"tricky/SliceSortBubble":                    {trickySrc, "SliceSortBubble", nil, tricky.SliceSortBubble},
	"tricky/SliceSortBy":                        {trickySrc, "SliceSortBy", nil, tricky.SliceSortBy},
	"tricky/SliceSortByFld":                     {trickySrc, "SliceSortByFld", nil, tricky.SliceSortByFld},
	"tricky/SliceSortByMultiple":                {trickySrc, "SliceSortByMultiple", nil, tricky.SliceSortByMultiple},
	"tricky/SliceSortStable":                    {trickySrc, "SliceSortStable", nil, tricky.SliceSortStable},
	"tricky/SliceSplice":                        {trickySrc, "SliceSplice", nil, tricky.SliceSplice},
	"tricky/SliceSplit":                         {trickySrc, "SliceSplit", nil, tricky.SliceSplit},
	"tricky/SliceSplitAtTest":                   {trickySrc, "SliceSplitAtTest", nil, tricky.SliceSplitAtTest},
	"tricky/SliceSplitByPredTest":               {trickySrc, "SliceSplitByPredTest", nil, tricky.SliceSplitByPredTest},
	"tricky/SliceStride":                        {trickySrc, "SliceStride", nil, tricky.SliceStride},
	"tricky/SliceStructIndex":                   {trickySrc, "SliceStructIndex", nil, tricky.SliceStructIndex},
	"tricky/SliceSubset":                        {trickySrc, "SliceSubset", nil, tricky.SliceSubset},
	"tricky/SliceSubsliceTest":                  {trickySrc, "SliceSubsliceTest", nil, tricky.SliceSubsliceTest},
	"tricky/SliceSum":                           {trickySrc, "SliceSum", nil, tricky.SliceSum},
	"tricky/SliceSumOddIdx":                     {trickySrc, "SliceSumOddIdx", nil, tricky.SliceSumOddIdx},
	"tricky/SliceSumRange":                      {trickySrc, "SliceSumRange", nil, tricky.SliceSumRange},
	"tricky/SliceSumRangeTest":                  {trickySrc, "SliceSumRangeTest", nil, tricky.SliceSumRangeTest},
	"tricky/SliceSwapElementsTest":              {trickySrc, "SliceSwapElementsTest", nil, tricky.SliceSwapElementsTest},
	"tricky/SliceSymmetricDiff":                 {trickySrc, "SliceSymmetricDiff", nil, tricky.SliceSymmetricDiff},
	"tricky/SliceSymmetricDiffTest":             {trickySrc, "SliceSymmetricDiffTest", nil, tricky.SliceSymmetricDiffTest},
	"tricky/SliceTailTest":                      {trickySrc, "SliceTailTest", nil, tricky.SliceTailTest},
	"tricky/SliceTake":                          {trickySrc, "SliceTake", nil, tricky.SliceTake},
	"tricky/SliceTakeDropTest":                  {trickySrc, "SliceTakeDropTest", nil, tricky.SliceTakeDropTest},
	"tricky/SliceTakeN":                         {trickySrc, "SliceTakeN", nil, tricky.SliceTakeN},
	"tricky/SliceTakeNFunc":                     {trickySrc, "SliceTakeNFunc", nil, tricky.SliceTakeNFunc},
	"tricky/SliceTakeTest":                      {trickySrc, "SliceTakeTest", nil, tricky.SliceTakeTest},
	"tricky/SliceTakeWhile":                     {trickySrc, "SliceTakeWhile", nil, tricky.SliceTakeWhile},
	"tricky/SliceTakeWhileDropWhile":            {trickySrc, "SliceTakeWhileDropWhile", nil, tricky.SliceTakeWhileDropWhile},
	"tricky/SliceTakeWhileTest":                 {trickySrc, "SliceTakeWhileTest", nil, tricky.SliceTakeWhileTest},
	"tricky/SliceTee":                           {trickySrc, "SliceTee", nil, tricky.SliceTee},
	"tricky/SliceTranspose":                     {trickySrc, "SliceTranspose", nil, tricky.SliceTranspose},
	"tricky/SliceTranspose2D":                   {trickySrc, "SliceTranspose2D", nil, tricky.SliceTranspose2D},
	"tricky/SliceTruncate":                      {trickySrc, "SliceTruncate", nil, tricky.SliceTruncate},
	"tricky/SliceUnion":                         {trickySrc, "SliceUnion", nil, tricky.SliceUnion},
	"tricky/SliceUniqBy":                        {trickySrc, "SliceUniqBy", nil, tricky.SliceUniqBy},
	"tricky/SliceUnique":                        {trickySrc, "SliceUnique", nil, tricky.SliceUnique},
	"tricky/SliceUniqueCountTest":               {trickySrc, "SliceUniqueCountTest", nil, tricky.SliceUniqueCountTest},
	"tricky/SliceUniquePreserveOrderTest":       {trickySrc, "SliceUniquePreserveOrderTest", nil, tricky.SliceUniquePreserveOrderTest},
	"tricky/SliceUniquePreserveTest":            {trickySrc, "SliceUniquePreserveTest", nil, tricky.SliceUniquePreserveTest},
	"tricky/SliceUnzip":                         {trickySrc, "SliceUnzip", nil, tricky.SliceUnzip},
	"tricky/SliceWindow":                        {trickySrc, "SliceWindow", nil, tricky.SliceWindow},
	"tricky/SliceWithout":                       {trickySrc, "SliceWithout", nil, tricky.SliceWithout},
	"tricky/SliceZip":                           {trickySrc, "SliceZip", nil, tricky.SliceZip},
	"tricky/SliceZipMap":                        {trickySrc, "SliceZipMap", nil, tricky.SliceZipMap},
	"tricky/SliceZipMapTest":                    {trickySrc, "SliceZipMapTest", nil, tricky.SliceZipMapTest},
	"tricky/SliceZipTest":                       {trickySrc, "SliceZipTest", nil, tricky.SliceZipTest},
	"tricky/SliceZipWith":                       {trickySrc, "SliceZipWith", nil, tricky.SliceZipWith},
	"tricky/SliceZipWithIndexTest":              {trickySrc, "SliceZipWithIndexTest", nil, tricky.SliceZipWithIndexTest},
	"tricky/StructAnon":                         {trickySrc, "StructAnon", nil, tricky.StructAnon},
	"tricky/StructAnonymousField":               {trickySrc, "StructAnonymousField", nil, tricky.StructAnonymousField},
	"tricky/StructCompareDiff":                  {trickySrc, "StructCompareDiff", nil, tricky.StructCompareDiff},
	"tricky/StructCompareDiffTest":              {trickySrc, "StructCompareDiffTest", nil, tricky.StructCompareDiffTest},
	"tricky/StructCompareDiffTypeTest":          {trickySrc, "StructCompareDiffTypeTest", nil, tricky.StructCompareDiffTypeTest},
	"tricky/StructCompareEqual":                 {trickySrc, "StructCompareEqual", nil, tricky.StructCompareEqual},
	"tricky/StructCompareNil":                   {trickySrc, "StructCompareNil", nil, tricky.StructCompareNil},
	"tricky/StructCompareNilPtrTest":            {trickySrc, "StructCompareNilPtrTest", nil, tricky.StructCompareNilPtrTest},
	"tricky/StructCompareSameTest":              {trickySrc, "StructCompareSameTest", nil, tricky.StructCompareSameTest},
	"tricky/StructCopyDeep":                     {trickySrc, "StructCopyDeep", nil, tricky.StructCopyDeep},
	"tricky/StructCopyPointerTest":              {trickySrc, "StructCopyPointerTest", nil, tricky.StructCopyPointerTest},
	"tricky/StructCopyValueTest":                {trickySrc, "StructCopyValueTest", nil, tricky.StructCopyValueTest},
	"tricky/StructEmbeddedAccessTest":           {trickySrc, "StructEmbeddedAccessTest", nil, tricky.StructEmbeddedAccessTest},
	"tricky/StructEmbeddedFldAccess":            {trickySrc, "StructEmbeddedFldAccess", nil, tricky.StructEmbeddedFldAccess},
	// Note: tricky/StructEmbeddedInterface moved to known_issue_test.go - complex interface/embedding issue
	"tricky/StructEmbeddedMethodOverride":       {trickySrc, "StructEmbeddedMethodOverride", nil, tricky.StructEmbeddedMethodOverride},
	"tricky/StructEmbeddedMethodOverrideTest":   {trickySrc, "StructEmbeddedMethodOverrideTest", nil, tricky.StructEmbeddedMethodOverrideTest},
	"tricky/StructEmbeddedNil":                  {trickySrc, "StructEmbeddedNil", nil, tricky.StructEmbeddedNil},
	"tricky/StructEmbeddedNilCheckTest":         {trickySrc, "StructEmbeddedNilCheckTest", nil, tricky.StructEmbeddedNilCheckTest},
	"tricky/StructEmbeddedNilDerefTest":         {trickySrc, "StructEmbeddedNilDerefTest", nil, tricky.StructEmbeddedNilDerefTest},
	"tricky/StructEmbeddedNilFld":               {trickySrc, "StructEmbeddedNilFld", nil, tricky.StructEmbeddedNilFld},
	"tricky/StructEmbeddedNilMethodTest":        {trickySrc, "StructEmbeddedNilMethodTest", nil, tricky.StructEmbeddedNilMethodTest},
	"tricky/StructEmbeddedOverride":             {trickySrc, "StructEmbeddedOverride", nil, tricky.StructEmbeddedOverride},
	"tricky/StructEmbeddedPtrInitTest":          {trickySrc, "StructEmbeddedPtrInitTest", nil, tricky.StructEmbeddedPtrInitTest},
	"tricky/StructEmbeddedPtrNilTest":           {trickySrc, "StructEmbeddedPtrNilTest", nil, tricky.StructEmbeddedPtrNilTest},
	"tricky/StructEmpty":                        {trickySrc, "StructEmpty", nil, tricky.StructEmpty},
	"tricky/StructFieldInitTest":                {trickySrc, "StructFieldInitTest", nil, tricky.StructFieldInitTest},
	"tricky/StructFieldModifyViaPtrTest":        {trickySrc, "StructFieldModifyViaPtrTest", nil, tricky.StructFieldModifyViaPtrTest},
	"tricky/StructFieldPointerModify":           {trickySrc, "StructFieldPointerModify", nil, tricky.StructFieldPointerModify},
	"tricky/StructFieldPointerModifyTest":       {trickySrc, "StructFieldPointerModifyTest", nil, tricky.StructFieldPointerModifyTest},
	"tricky/StructFieldPtr":                     {trickySrc, "StructFieldPtr", nil, tricky.StructFieldPtr},
	"tricky/StructFieldPtrTest":                 {trickySrc, "StructFieldPtrTest", nil, tricky.StructFieldPtrTest},
	"tricky/StructFieldShadow":                  {trickySrc, "StructFieldShadow", nil, tricky.StructFieldShadow},
	"tricky/StructFieldShadowTest":              {trickySrc, "StructFieldShadowTest", nil, tricky.StructFieldShadowTest},
	"tricky/StructFldModify":                    {trickySrc, "StructFldModify", nil, tricky.StructFldModify},
	"tricky/StructFldPtrModify":                 {trickySrc, "StructFldPtrModify", nil, tricky.StructFldPtrModify},
	"tricky/StructInterface":                    {trickySrc, "StructInterface", nil, tricky.StructInterface},
	"tricky/StructMethodChain":                  {trickySrc, "StructMethodChain", nil, tricky.StructMethodChain},
	"tricky/StructMethodChainNilTest":           {trickySrc, "StructMethodChainNilTest", nil, tricky.StructMethodChainNilTest},
	"tricky/StructMethodChainTest":              {trickySrc, "StructMethodChainTest", nil, tricky.StructMethodChainTest},
	"tricky/StructMethodEmbeddedTest":           {trickySrc, "StructMethodEmbeddedTest", nil, tricky.StructMethodEmbeddedTest},
	"tricky/StructMethodNilPtrTest":             {trickySrc, "StructMethodNilPtrTest", nil, tricky.StructMethodNilPtrTest},
	"tricky/StructMethodOnAddr":                 {trickySrc, "StructMethodOnAddr", nil, tricky.StructMethodOnAddr},
	"tricky/StructMethodOnEmbeddedTest":         {trickySrc, "StructMethodOnEmbeddedTest", nil, tricky.StructMethodOnEmbeddedTest},
	"tricky/StructMethodOnNilPtrTest":           {trickySrc, "StructMethodOnNilPtrTest", nil, tricky.StructMethodOnNilPtrTest},
	"tricky/StructMethodOnNilReceiver":          {trickySrc, "StructMethodOnNilReceiver", nil, tricky.StructMethodOnNilReceiver},
	"tricky/StructMethodOnValTest":              {trickySrc, "StructMethodOnValTest", nil, tricky.StructMethodOnValTest},
	"tricky/StructMethodOnValueCopy":            {trickySrc, "StructMethodOnValueCopy", nil, tricky.StructMethodOnValueCopy},
	"tricky/StructMethodPtrRecTest":             {trickySrc, "StructMethodPtrRecTest", nil, tricky.StructMethodPtrRecTest},
	"tricky/StructMethodValRec":                 {trickySrc, "StructMethodValRec", nil, tricky.StructMethodValRec},
	"tricky/StructMethodValue":                  {trickySrc, "StructMethodValue", nil, tricky.StructMethodValue},
	"tricky/StructMethodValueReceiverTest":      {trickySrc, "StructMethodValueReceiverTest", nil, tricky.StructMethodValueReceiverTest},
	"tricky/StructMethodWithPointerReceiver":    {trickySrc, "StructMethodWithPointerReceiver", nil, tricky.StructMethodWithPointerReceiver},
	"tricky/StructMethodWithVariadic":           {trickySrc, "StructMethodWithVariadic", nil, tricky.StructMethodWithVariadic},
	"tricky/StructModifyViaPointerTest":         {trickySrc, "StructModifyViaPointerTest", nil, tricky.StructModifyViaPointerTest},
	"tricky/StructNestedAssign":                 {trickySrc, "StructNestedAssign", nil, tricky.StructNestedAssign},
	"tricky/StructNestedInitTest":               {trickySrc, "StructNestedInitTest", nil, tricky.StructNestedInitTest},
	"tricky/StructNestedMethodTest":             {trickySrc, "StructNestedMethodTest", nil, tricky.StructNestedMethodTest},
	"tricky/StructNestedPtrTest":                {trickySrc, "StructNestedPtrTest", nil, tricky.StructNestedPtrTest},
	"tricky/StructNilField":                     {trickySrc, "StructNilField", nil, tricky.StructNilField},
	"tricky/StructNilFieldDerefTest":            {trickySrc, "StructNilFieldDerefTest", nil, tricky.StructNilFieldDerefTest},
	"tricky/StructNilFieldHolder":               {trickySrc, "StructNilFieldHolder", nil, tricky.StructNilFieldHolder},
	"tricky/StructNilFieldInitTest":             {trickySrc, "StructNilFieldInitTest", nil, tricky.StructNilFieldInitTest},
	"tricky/StructNilPointerMethod":             {trickySrc, "StructNilPointerMethod", nil, tricky.StructNilPointerMethod},
	"tricky/StructNilSafeMethodTest":            {trickySrc, "StructNilSafeMethodTest", nil, tricky.StructNilSafeMethodTest},
	"tricky/StructPointerMethodChain":           {trickySrc, "StructPointerMethodChain", nil, tricky.StructPointerMethodChain},
	"tricky/StructPtrMethod":                    {trickySrc, "StructPtrMethod", nil, tricky.StructPtrMethod},
	"tricky/StructPtrMethodOnNilTest":           {trickySrc, "StructPtrMethodOnNilTest", nil, tricky.StructPtrMethodOnNilTest},
	"tricky/StructSelfRef":                      {trickySrc, "StructSelfRef", nil, tricky.StructSelfRef},
	"tricky/StructSliceAppend":                  {trickySrc, "StructSliceAppend", nil, tricky.StructSliceAppend},
	"tricky/StructSliceFieldAppendTest":         {trickySrc, "StructSliceFieldAppendTest", nil, tricky.StructSliceFieldAppendTest},
	"tricky/StructSliceOfPointers":              {trickySrc, "StructSliceOfPointers", nil, tricky.StructSliceOfPointers},
	"tricky/StructSliceOfSlices":                {trickySrc, "StructSliceOfSlices", nil, tricky.StructSliceOfSlices},
	"tricky/StructValidation":                   {trickySrc, "StructValidation", nil, tricky.StructValidation},
	"tricky/StructValidationTest":               {trickySrc, "StructValidationTest", nil, tricky.StructValidationTest},
	"tricky/StructWithAnonymousFunc":            {trickySrc, "StructWithAnonymousFunc", nil, tricky.StructWithAnonymousFunc},
	"tricky/StructWithArrFieldTest":             {trickySrc, "StructWithArrFieldTest", nil, tricky.StructWithArrFieldTest},
	"tricky/StructWithArrInitTest":              {trickySrc, "StructWithArrInitTest", nil, tricky.StructWithArrInitTest},
	"tricky/StructWithBoolFld":                  {trickySrc, "StructWithBoolFld", nil, tricky.StructWithBoolFld},
	"tricky/StructWithChanFieldTest":            {trickySrc, "StructWithChanFieldTest", nil, tricky.StructWithChanFieldTest},
	"tricky/StructWithChanFld":                  {trickySrc, "StructWithChanFld", nil, tricky.StructWithChanFld},
	"tricky/StructWithChannel":                  {trickySrc, "StructWithChannel", nil, tricky.StructWithChannel},
	"tricky/StructWithChanNilInitTest":          {trickySrc, "StructWithChanNilInitTest", nil, tricky.StructWithChanNilInitTest},
	"tricky/StructWithChanOfChan":               {trickySrc, "StructWithChanOfChan", nil, tricky.StructWithChanOfChan},
	"tricky/StructWithChanOfChanTest":           {trickySrc, "StructWithChanOfChanTest", nil, tricky.StructWithChanOfChanTest},
	"tricky/StructWithChanTest":                 {trickySrc, "StructWithChanTest", nil, tricky.StructWithChanTest},
	"tricky/StructWithChanTest2":                {trickySrc, "StructWithChanTest2", nil, tricky.StructWithChanTest2},
	"tricky/StructWithComputedField":            {trickySrc, "StructWithComputedField", nil, tricky.StructWithComputedField},
	"tricky/StructWithComputedFieldTest":        {trickySrc, "StructWithComputedFieldTest", nil, tricky.StructWithComputedFieldTest},
	"tricky/StructWithComputedFld":              {trickySrc, "StructWithComputedFld", nil, tricky.StructWithComputedFld},
	"tricky/StructWithDoublePointer":            {trickySrc, "StructWithDoublePointer", nil, tricky.StructWithDoublePointer},
	"tricky/StructWithEmbeddedNilPtrTest":       {trickySrc, "StructWithEmbeddedNilPtrTest", nil, tricky.StructWithEmbeddedNilPtrTest},
	"tricky/StructWithEmbeddedPointer":          {trickySrc, "StructWithEmbeddedPointer", nil, tricky.StructWithEmbeddedPointer},
	"tricky/StructWithEmbeddedPtrTest":          {trickySrc, "StructWithEmbeddedPtrTest", nil, tricky.StructWithEmbeddedPtrTest},
	"tricky/StructWithEmbeddedTest":             {trickySrc, "StructWithEmbeddedTest", nil, tricky.StructWithEmbeddedTest},
	"tricky/StructWithEmptySliceTest":           {trickySrc, "StructWithEmptySliceTest", nil, tricky.StructWithEmptySliceTest},
	"tricky/StructWithFldValidation":            {trickySrc, "StructWithFldValidation", nil, tricky.StructWithFldValidation},
	"tricky/StructWithFloatField":               {trickySrc, "StructWithFloatField", nil, tricky.StructWithFloatField},
	"tricky/StructWithFloatFldTest":             {trickySrc, "StructWithFloatFldTest", nil, tricky.StructWithFloatFldTest},
	"tricky/StructWithFunc":                     {trickySrc, "StructWithFunc", nil, tricky.StructWithFunc},
	"tricky/StructWithFuncField":                {trickySrc, "StructWithFuncField", nil, tricky.StructWithFuncField},
	"tricky/StructWithFuncFieldCallTest":        {trickySrc, "StructWithFuncFieldCallTest", nil, tricky.StructWithFuncFieldCallTest},
	"tricky/StructWithFuncFieldNilTest":         {trickySrc, "StructWithFuncFieldNilTest", nil, tricky.StructWithFuncFieldNilTest},
	"tricky/StructWithFuncFieldTest":            {trickySrc, "StructWithFuncFieldTest", nil, tricky.StructWithFuncFieldTest},
	"tricky/StructWithFuncFldCall":              {trickySrc, "StructWithFuncFldCall", nil, tricky.StructWithFuncFldCall},
	"tricky/StructWithFuncFldExec":              {trickySrc, "StructWithFuncFldExec", nil, tricky.StructWithFuncFldExec},
	"tricky/StructWithFuncMap":                  {trickySrc, "StructWithFuncMap", nil, tricky.StructWithFuncMap},
	"tricky/StructWithFuncPtrTest":              {trickySrc, "StructWithFuncPtrTest", nil, tricky.StructWithFuncPtrTest},
	"tricky/StructWithFuncReturningStruct":      {trickySrc, "StructWithFuncReturningStruct", nil, tricky.StructWithFuncReturningStruct},
	"tricky/StructWithFuncSlice":                {trickySrc, "StructWithFuncSlice", nil, tricky.StructWithFuncSlice},
	"tricky/StructWithFuncSliceComplex":         {trickySrc, "StructWithFuncSliceComplex", nil, tricky.StructWithFuncSliceComplex},
	"tricky/StructWithInitFunc":                 {trickySrc, "StructWithInitFunc", nil, tricky.StructWithInitFunc},
	"tricky/StructWithInterfaceFldTest":         {trickySrc, "StructWithInterfaceFldTest", nil, tricky.StructWithInterfaceFldTest},
	"tricky/StructWithInterfaceMap":             {trickySrc, "StructWithInterfaceMap", nil, tricky.StructWithInterfaceMap},
	"tricky/StructWithInterfaceSlice":           {trickySrc, "StructWithInterfaceSlice", nil, tricky.StructWithInterfaceSlice},
	"tricky/StructWithIntField":                 {trickySrc, "StructWithIntField", nil, tricky.StructWithIntField},
	"tricky/StructWithIntSliceTest":             {trickySrc, "StructWithIntSliceTest", nil, tricky.StructWithIntSliceTest},
	"tricky/StructWithLazyFieldTest":            {trickySrc, "StructWithLazyFieldTest", nil, tricky.StructWithLazyFieldTest},
	"tricky/StructWithLazyFld":                  {trickySrc, "StructWithLazyFld", nil, tricky.StructWithLazyFld},
	"tricky/StructWithLazyInit":                 {trickySrc, "StructWithLazyInit", nil, tricky.StructWithLazyInit},
	"tricky/StructWithMapFld":                   {trickySrc, "StructWithMapFld", nil, tricky.StructWithMapFld},
	"tricky/StructWithMapInitTest":              {trickySrc, "StructWithMapInitTest", nil, tricky.StructWithMapInitTest},
	"tricky/StructWithMapMakeTest":              {trickySrc, "StructWithMapMakeTest", nil, tricky.StructWithMapMakeTest},
	"tricky/StructWithMapNilInit":               {trickySrc, "StructWithMapNilInit", nil, tricky.StructWithMapNilInit},
	"tricky/StructWithMapNilInitTest":           {trickySrc, "StructWithMapNilInitTest", nil, tricky.StructWithMapNilInitTest},
	"tricky/StructWithMapOfPtrTest":             {trickySrc, "StructWithMapOfPtrTest", nil, tricky.StructWithMapOfPtrTest},
	"tricky/StructWithMapOfSlices":              {trickySrc, "StructWithMapOfSlices", nil, tricky.StructWithMapOfSlices},
	"tricky/StructWithMapOfStructs":             {trickySrc, "StructWithMapOfStructs", nil, tricky.StructWithMapOfStructs},
	"tricky/StructWithMapPointer":               {trickySrc, "StructWithMapPointer", nil, tricky.StructWithMapPointer},
	"tricky/StructWithMapRangeDel":              {trickySrc, "StructWithMapRangeDel", nil, tricky.StructWithMapRangeDel},
	"tricky/StructWithMethodClosure":            {trickySrc, "StructWithMethodClosure", nil, tricky.StructWithMethodClosure},
	"tricky/StructWithMethodPointer":            {trickySrc, "StructWithMethodPointer", nil, tricky.StructWithMethodPointer},
	"tricky/StructWithNestedFunc":               {trickySrc, "StructWithNestedFunc", nil, tricky.StructWithNestedFunc},
	"tricky/StructWithNestedPointer":            {trickySrc, "StructWithNestedPointer", nil, tricky.StructWithNestedPointer},
	"tricky/StructWithNestedSlice":              {trickySrc, "StructWithNestedSlice", nil, tricky.StructWithNestedSlice},
	"tricky/StructWithNilChan":                  {trickySrc, "StructWithNilChan", nil, tricky.StructWithNilChan},
	"tricky/StructWithNilChanFieldTest":         {trickySrc, "StructWithNilChanFieldTest", nil, tricky.StructWithNilChanFieldTest},
	"tricky/StructWithNilChanFld":               {trickySrc, "StructWithNilChanFld", nil, tricky.StructWithNilChanFld},
	"tricky/StructWithNilFieldInitTest":         {trickySrc, "StructWithNilFieldInitTest", nil, tricky.StructWithNilFieldInitTest},
	"tricky/StructWithNilPtrTest":               {trickySrc, "StructWithNilPtrTest", nil, tricky.StructWithNilPtrTest},
	"tricky/StructWithNilSlice":                 {trickySrc, "StructWithNilSlice", nil, tricky.StructWithNilSlice},
	"tricky/StructWithNilSliceFieldTest":        {trickySrc, "StructWithNilSliceFieldTest", nil, tricky.StructWithNilSliceFieldTest},
	"tricky/StructWithPointerField":             {trickySrc, "StructWithPointerField", nil, tricky.StructWithPointerField},
	"tricky/StructWithPointerInterface":         {trickySrc, "StructWithPointerInterface", nil, tricky.StructWithPointerInterface},
	"tricky/StructWithPointerMap":               {trickySrc, "StructWithPointerMap", nil, tricky.StructWithPointerMap},
	"tricky/StructWithPointerSlice":             {trickySrc, "StructWithPointerSlice", nil, tricky.StructWithPointerSlice},
	"tricky/StructWithPointerToInterface":       {trickySrc, "StructWithPointerToInterface", nil, tricky.StructWithPointerToInterface},
	"tricky/StructWithPointerToMap":             {trickySrc, "StructWithPointerToMap", nil, tricky.StructWithPointerToMap},
	"tricky/StructWithPointerToSelf":            {trickySrc, "StructWithPointerToSelf", nil, tricky.StructWithPointerToSelf},
	"tricky/StructWithPtrFld":                   {trickySrc, "StructWithPtrFld", nil, tricky.StructWithPtrFld},
	"tricky/StructWithPtrMethodTest":            {trickySrc, "StructWithPtrMethodTest", nil, tricky.StructWithPtrMethodTest},
	"tricky/StructWithPtrSliceFieldTest":        {trickySrc, "StructWithPtrSliceFieldTest", nil, tricky.StructWithPtrSliceFieldTest},
	"tricky/StructWithPtrToStructTest":          {trickySrc, "StructWithPtrToStructTest", nil, tricky.StructWithPtrToStructTest},
	"tricky/StructWithRecursiveType":            {trickySrc, "StructWithRecursiveType", nil, tricky.StructWithRecursiveType},
	"tricky/StructWithSelfRefPointer":           {trickySrc, "StructWithSelfRefPointer", nil, tricky.StructWithSelfRefPointer},
	"tricky/StructWithSliceAppendMethodTest":    {trickySrc, "StructWithSliceAppendMethodTest", nil, tricky.StructWithSliceAppendMethodTest},
	"tricky/StructWithSliceAppendTest":          {trickySrc, "StructWithSliceAppendTest", nil, tricky.StructWithSliceAppendTest},
	"tricky/StructWithSliceFieldNamed":          {trickySrc, "StructWithSliceFieldNamed", nil, tricky.StructWithSliceFieldNamed},
	"tricky/StructWithSliceFld":                 {trickySrc, "StructWithSliceFld", nil, tricky.StructWithSliceFld},
	"tricky/StructWithSliceMakeTest":            {trickySrc, "StructWithSliceMakeTest", nil, tricky.StructWithSliceMakeTest},
	"tricky/StructWithSliceMethods":             {trickySrc, "StructWithSliceMethods", nil, tricky.StructWithSliceMethods},
	"tricky/StructWithSliceNil":                 {trickySrc, "StructWithSliceNil", nil, tricky.StructWithSliceNil},
	"tricky/StructWithSliceNilInitTest":         {trickySrc, "StructWithSliceNilInitTest", nil, tricky.StructWithSliceNilInitTest},
	"tricky/StructWithSliceOfMaps":              {trickySrc, "StructWithSliceOfMaps", nil, tricky.StructWithSliceOfMaps},
	"tricky/StructWithSliceOfPointersToStructs": {trickySrc, "StructWithSliceOfPointersToStructs", nil, tricky.StructWithSliceOfPointersToStructs},
	"tricky/StructWithSliceOfPtrTest":           {trickySrc, "StructWithSliceOfPtrTest", nil, tricky.StructWithSliceOfPtrTest},
	"tricky/StructWithSlicePointer":             {trickySrc, "StructWithSlicePointer", nil, tricky.StructWithSlicePointer},
	"tricky/StructWithStringFld":                {trickySrc, "StructWithStringFld", nil, tricky.StructWithStringFld},
	"tricky/StructWithTag":                      {trickySrc, "StructWithTag", nil, tricky.StructWithTag},
	"tricky/StructWithTwoFlds":                  {trickySrc, "StructWithTwoFlds", nil, tricky.StructWithTwoFlds},
	"tricky/StructWithUintField":                {trickySrc, "StructWithUintField", nil, tricky.StructWithUintField},
	"tricky/StructWithUintFldTest":              {trickySrc, "StructWithUintFldTest", nil, tricky.StructWithUintFldTest},
	"tricky/StructWithValidation":               {trickySrc, "StructWithValidation", nil, tricky.StructWithValidation},
	"tricky/StructZeroInitTest":                 {trickySrc, "StructZeroInitTest", nil, tricky.StructZeroInitTest},
	"tricky/StructZeroValueCheckTest":           {trickySrc, "StructZeroValueCheckTest", nil, tricky.StructZeroValueCheckTest},

	// ============================================================================
	// goroutine
	// ============================================================================
	"goroutine/BasicSpawn":                  {goroutineSrc, "BasicSpawn", nil, goroutine.BasicSpawn},
	"goroutine/ChannelCommunication":        {goroutineSrc, "ChannelCommunication", nil, goroutine.ChannelCommunication},
	"goroutine/WithArguments":               {goroutineSrc, "WithArguments", nil, goroutine.WithArguments},
	"goroutine/WithStruct":                  {goroutineSrc, "WithStruct", nil, goroutine.WithStruct},
	"goroutine/DifferentTypes":              {goroutineSrc, "DifferentTypes", nil, goroutine.DifferentTypes},
	"goroutine/GlobalsSharing":              {goroutineSrc, "GlobalsSharing", nil, goroutine.GlobalsSharing},
	"goroutine/MultipleSends":               {goroutineSrc, "MultipleSends", nil, goroutine.MultipleSends},
	"goroutine/ParallelExecution":           {goroutineSrc, "ParallelExecution", nil, goroutine.ParallelExecution},
	"goroutine/ClosureCapture":              {goroutineSrc, "ClosureCapture", nil, goroutine.ClosureCapture},
	"goroutine/ClosureCaptureMultiple":      {goroutineSrc, "ClosureCaptureMultiple", nil, goroutine.ClosureCaptureMultiple},
	"goroutine/SelectStatement":             {goroutineSrc, "SelectStatement", nil, goroutine.SelectStatement},
	"goroutine/SelectDefault":               {goroutineSrc, "SelectDefault", nil, goroutine.SelectDefault},
	"goroutine/SelectSend":                  {goroutineSrc, "SelectSend", nil, goroutine.SelectSend},
	"goroutine/RangeOverChannel":            {goroutineSrc, "RangeOverChannel", nil, goroutine.RangeOverChannel},
	"goroutine/RangeOverChannelWithBuiltin": {goroutineSrc, "RangeOverChannelWithBuiltin", nil, goroutine.RangeOverChannelWithBuiltin},

	// ============================================================================
	// resolved_issue
	// ============================================================================
	"resolved_issue/BytesToString":                      {resolvedIssueSrc, "BytesToString", nil, resolved_issue.BytesToString},
	"resolved_issue/BytesToStringHi":                    {resolvedIssueSrc, "BytesToStringHi", nil, resolved_issue.BytesToStringHi},
	"resolved_issue/BytesToStringGo":                    {resolvedIssueSrc, "BytesToStringGo", nil, resolved_issue.BytesToStringGo},
	"resolved_issue/BytesToStringEmpty":                 {resolvedIssueSrc, "BytesToStringEmpty", nil, resolved_issue.BytesToStringEmpty},
	"resolved_issue/BytesToStringSingle":                {resolvedIssueSrc, "BytesToStringSingle", nil, resolved_issue.BytesToStringSingle},
	"resolved_issue/PointerReceiverMutation":            {resolvedIssueSrc, "PointerReceiverMutation", nil, resolved_issue.PointerReceiverMutation},
	"resolved_issue/PointerReceiverMutationReturnValue": {resolvedIssueSrc, "PointerReceiverMutationReturnValue", nil, resolved_issue.PointerReceiverMutationReturnValue},
	"resolved_issue/InitFuncExecuted":                   {resolvedIssueSrc, "InitFuncExecuted", nil, resolved_issue.InitFuncExecuted},
	"resolved_issue/InitFuncSideEffect":                 {resolvedIssueSrc, "InitFuncSideEffect", nil, resolved_issue.InitFuncSideEffect},
	"resolved_issue/RangeStringRuneValue":               {resolvedIssueSrc, "RangeStringRuneValue", nil, resolved_issue.RangeStringRuneValue},
	"resolved_issue/RangeStringIndexValue":              {resolvedIssueSrc, "RangeStringIndexValue", nil, resolved_issue.RangeStringIndexValue},
	"resolved_issue/RangeStringMultibyte":               {resolvedIssueSrc, "RangeStringMultibyte", nil, resolved_issue.RangeStringMultibyte},
	"resolved_issue/MapWithFuncValue":                   {resolvedIssueSrc, "MapWithFuncValue", nil, resolved_issue.MapWithFuncValue},
	"resolved_issue/InterfaceSliceTypeSwitch":           {resolvedIssueSrc, "InterfaceSliceTypeSwitch", nil, resolved_issue.InterfaceSliceTypeSwitch},
	"resolved_issue/StructWithFuncField":                {resolvedIssueSrc, "StructWithFuncField", nil, resolved_issue.StructWithFuncField},
	"resolved_issue/SliceFlatten":                       {resolvedIssueSrc, "SliceFlatten", nil, resolved_issue.SliceFlatten},
	// Note: resolved_issue/MapUpdateDuringRange removed - non-deterministic map iteration order
	"resolved_issue/StructSelfRef":         {resolvedIssueSrc, "StructSelfRef", nil, resolved_issue.StructSelfRef},
	"resolved_issue/DeferInClosureWithArg": {resolvedIssueSrc, "DeferInClosureWithArg", nil, resolved_issue.DeferInClosureWithArg},
	"resolved_issue/PointerSwapInStruct":   {resolvedIssueSrc, "PointerSwapInStruct", nil, resolved_issue.PointerSwapInStruct},
	"resolved_issue/StructWithFuncSlice":   {resolvedIssueSrc, "StructWithFuncSlice", nil, resolved_issue.StructWithFuncSlice},
	"resolved_issue/StructAnonymousField":  {resolvedIssueSrc, "StructAnonymousField", nil, resolved_issue.StructAnonymousField},
	// Note: resolved_issue/MapRangeWithBreak removed - non-deterministic map iteration order
	"resolved_issue/PointerToInterface":           {resolvedIssueSrc, "PointerToInterface", nil, resolved_issue.PointerToInterface},
	"resolved_issue/PointerToSliceElemModify":     {resolvedIssueSrc, "PointerToSliceElemModify", nil, resolved_issue.PointerToSliceElemModify},
	"resolved_issue/StructWithFuncPtrTest":        {resolvedIssueSrc, "StructWithFuncPtrTest", nil, resolved_issue.StructWithFuncPtrTest},
	"resolved_issue/PointerCompareDiffTest":       {resolvedIssueSrc, "PointerCompareDiffTest", nil, resolved_issue.PointerCompareDiffTest},
	"resolved_issue/DeferModifyMultipleNamedTest": {resolvedIssueSrc, "DeferModifyMultipleNamedTest", nil, resolved_issue.DeferModifyMultipleNamedTest},
	"resolved_issue/DeferNamedReturnNilTest":      {resolvedIssueSrc, "DeferNamedReturnNilTest", nil, resolved_issue.DeferNamedReturnNilTest},
	"resolved_issue/DeferNamedReturnNilPtrTest":   {resolvedIssueSrc, "DeferNamedReturnNilPtrTest", nil, resolved_issue.DeferNamedReturnNilPtrTest},
	"resolved_issue/DeferNamedReturnMultiTest":    {resolvedIssueSrc, "DeferNamedReturnMultiTest", nil, resolved_issue.DeferNamedReturnMultiTest},
	// ============================================================================
	// typeconv
	// ============================================================================
	"typeconv/IntToFloat64":           {typeconvSrc, "IntToFloat64", nil, typeconv.IntToFloat64},
	"typeconv/Float64Arithmetic":      {typeconvSrc, "Float64Arithmetic", nil, typeconv.Float64Arithmetic},
	"typeconv/StringToByteConversion": {typeconvSrc, "StringToByteConversion", nil, typeconv.StringToByteConversion},
	"typeconv/IntStringConversion":    {typeconvSrc, "IntStringConversion", nil, typeconv.IntStringConversion},
	"typeconv/StringIntConversion":    {typeconvSrc, "StringIntConversion", nil, typeconv.StringIntConversion},
	// Parameterized tests
	"typeconv/IntToString":     {typeconvSrc, "IntToString", []any{12345}, typeconv.IntToString},
	"typeconv/StringToInt":     {typeconvSrc, "StringToInt", []any{"54321"}, typeconv.StringToInt},
	"typeconv/IntToFloatToInt": {typeconvSrc, "IntToFloatToInt", []any{42}, typeconv.IntToFloatToInt},

	// ============================================================================
	// variables
	// ============================================================================
	"variables/DeclareAndUse":   {variablesSrc, "DeclareAndUse", nil, variables.DeclareAndUse},
	"variables/Reassignment":    {variablesSrc, "Reassignment", nil, variables.Reassignment},
	"variables/MultipleDecl":    {variablesSrc, "MultipleDecl", nil, variables.MultipleDecl},
	"variables/ZeroValues":      {variablesSrc, "ZeroValues", nil, variables.ZeroValues},
	"variables/StringZeroValue": {variablesSrc, "StringZeroValue", nil, variables.StringZeroValue},
	"variables/Shadowing":       {variablesSrc, "Shadowing", nil, variables.Shadowing},
	// Parameterized tests
	"variables/SumThree":   {variablesSrc, "SumThree", []any{10, 20, 30}, variables.SumThree},
	"variables/Multiply":   {variablesSrc, "Multiply", []any{6, 7}, variables.Multiply},
	"variables/Max":        {variablesSrc, "Max", []any{100, 42}, variables.Max},
	"variables/IsPositive": {variablesSrc, "IsPositive", []any{5}, variables.IsPositive},

	// ============================================================================
	// init
	// ============================================================================
	"init/GetA": {initSrc, "GetA", nil, initialize.GetA},
}

// ============================================================================
// Test Runner
// ============================================================================

// TestCorrectness runs all correctness tests
func TestCorrectness(t *testing.T) {
	if len(allCorrectnessTests) == 0 {
		t.Skip("No correctness tests defined")
	}

	passed := 0
	failed := 0

	for name, tc := range allCorrectnessTests {
		t.Run(name, func(t *testing.T) {
			src := toMainPackage(tc.src)
			prog, err := gig.Build(src)
			if err != nil {
				t.Fatalf("Build error: %v", err)
			}

			startInterp := time.Now()
			result, err := prog.Run(tc.funcName, tc.args...)
			interpDuration := time.Since(startInterp)
			if err != nil {
				t.Fatalf("Run error: %v", err)
			}

			startNative := time.Now()
			expected := callNative(tc.native, tc.args)
			nativeDuration := time.Since(startNative)

			compareCorrectnessResults(t, result, expected)

			passed++
			ratio := float64(interpDuration) / float64(nativeDuration)
			t.Logf("interp: %v, native: %v, ratio: %.1fx", interpDuration, nativeDuration, ratio)
		})
	}

	t.Logf("\n=== Correctness Test Summary ===")
	t.Logf("Total:   %d", passed+failed)
	t.Logf("Passed:  %d", passed)
	t.Logf("Failed:  %d", failed)
}
