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
	"git.woa.com/youngjin/gig/tests/testdata/initialize"
	"git.woa.com/youngjin/gig/tests/testdata/leetcode_hard"
	"git.woa.com/youngjin/gig/tests/testdata/mapadvanced"
	"git.woa.com/youngjin/gig/tests/testdata/maps"
	"git.woa.com/youngjin/gig/tests/testdata/multiassign"
	"git.woa.com/youngjin/gig/tests/testdata/namedreturn"
	"git.woa.com/youngjin/gig/tests/testdata/recursion"
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
	"algorithms/Reverse":        {algorithmsSrc, "Reverse", []any{[]int{1, 2, 3, 4, 5}}, algorithms.Reverse},
	"algorithms/Power":          {algorithmsSrc, "Power", []any{2, 10}, algorithms.Power},
	"algorithms/CountDigitsN":   {algorithmsSrc, "CountDigitsN", []any{12345}, algorithms.CountDigitsN},
	"algorithms/CollatzStepsN":  {algorithmsSrc, "CollatzStepsN", []any{27}, algorithms.CollatzStepsN},
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
	"advanced/FindFirst":    {advancedSrc, "FindFirst", []any{[]int{10, 20, 30}, 20}, advanced.FindFirst},
	"advanced/Bsearch":      {advancedSrc, "Bsearch", []any{[]int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}, 50}, advanced.Bsearch},
	"advanced/Gcd":          {advancedSrc, "Gcd", []any{48, 18}, advanced.Gcd},
	"advanced/Identity":     {advancedSrc, "Identity", []any{42}, advanced.Identity},
	"advanced/Minmax":       {advancedSrc, "Minmax", []any{[]int{3, 1, 4, 1, 5}}, advanced.Minmax},
	"advanced/Countdown":    {advancedSrc, "Countdown", []any{50}, advanced.Countdown},
	"advanced/Add":          {advancedSrc, "Add", []any{1, 2}, advanced.Add},
	"advanced/Mul":          {advancedSrc, "Mul", []any{3, 4}, advanced.Mul},
	"advanced/Sub":          {advancedSrc, "Sub", []any{10, 5}, advanced.Sub},

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
	"arithmetic/Add":        {arithmeticSrc, "Add", []any{10, 32}, arithmetic.Add},
	"arithmetic/Sub":        {arithmeticSrc, "Sub", []any{100, 42}, arithmetic.Sub},
	"arithmetic/Mul":        {arithmeticSrc, "Mul", []any{6, 7}, arithmetic.Mul},
	"arithmetic/Div":        {arithmeticSrc, "Div", []any{100, 4}, arithmetic.Div},
	"arithmetic/Mod":        {arithmeticSrc, "Mod", []any{17, 5}, arithmetic.Mod},
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
	"cornercases/Func_VariadicMultiple":     {cornercasesSrc, "Func_VariadicMultiple", nil, cornercases.Func_VariadicMultiple},
	"cornercases/Func_RecursionBase":        {cornercasesSrc, "Func_RecursionBase", nil, cornercases.Func_RecursionBase},
	"cornercases/Closure_ReturnClosure":     {cornercasesSrc, "Closure_ReturnClosure", nil, cornercases.Closure_ReturnClosure},
	"cornercases/Closure_CaptureVariable":   {cornercasesSrc, "Closure_CaptureVariable", nil, cornercases.Closure_CaptureVariable},
	"cornercases/Closure_ModifyCaptured":    {cornercasesSrc, "Closure_ModifyCaptured", nil, cornercases.Closure_ModifyCaptured},
	"cornercases/Struct_ZeroValueFields":    {cornercasesSrc, "Struct_ZeroValueFields", nil, cornercases.Struct_ZeroValueFields},
	"cornercases/Struct_PointerReceiver":    {cornercasesSrc, "Struct_PointerReceiver", nil, cornercases.Struct_PointerReceiver},
	"cornercases/Struct_NestedStruct":       {cornercasesSrc, "Struct_NestedStruct", nil, cornercases.Struct_NestedStruct},

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
	"external/FmtSprintfInt":        {externalSrc, "FmtSprintfInt", []any{42}, external.FmtSprintfInt},
	"external/StringsToUpperStr":    {externalSrc, "StringsToUpperStr", []any{"hello"}, external.StringsToUpperStr},
	"external/StringsToLowerStr":    {externalSrc, "StringsToLowerStr", []any{"HELLO"}, external.StringsToLowerStr},
	"external/StringsContainsStr":   {externalSrc, "StringsContainsStr", []any{"hello world", "world"}, external.StringsContainsStr},
	"external/StrconvItoaN":         {externalSrc, "StrconvItoaN", []any{42}, external.StrconvItoaN},
	"external/StrconvAtoiStr":       {externalSrc, "StrconvAtoiStr", []any{"123"}, external.StrconvAtoiStr},

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
	"functions/Add":             {functionsSrc, "Add", []any{5, 7}, functions.Add},
	"functions/Swap":            {functionsSrc, "Swap", []any{3, 7}, functions.Swap},
	"functions/Divmod":          {functionsSrc, "Divmod", []any{17, 5}, functions.Divmod},
	"functions/FactorialN":      {functionsSrc, "FactorialN", []any{5}, functions.FactorialN},
	"functions/FibIterN":        {functionsSrc, "FibIterN", []any{20}, functions.FibIterN},
	"functions/FibRecN":         {functionsSrc, "FibRecN", []any{15}, functions.FibRecN},
	"functions/IsEvenN":         {functionsSrc, "IsEvenN", []any{10}, functions.IsEvenN},
	"functions/IsOddN":          {functionsSrc, "IsOddN", []any{7}, functions.IsOddN},
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
	"recursion/SumTail":   {recursionSrc, "SumTail", []any{50, 0}, recursion.SumTail},
	"recursion/HanoiN":    {recursionSrc, "HanoiN", []any{10}, recursion.HanoiN},
	"recursion/Ack":       {recursionSrc, "Ack", []any{2, 3}, recursion.Ack},
	"recursion/MaxVal":    {recursionSrc, "MaxVal", []any{[]int{3, 7, 1, 9, 4}, 5}, recursion.MaxVal},

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
	"strings_pkg/StrConcat":   {stringsPkgSrc, "StrConcat", []any{"hello", " world"}, strings_pkg.StrConcat},
	"strings_pkg/StrLen":      {stringsPkgSrc, "StrLen", []any{"hello"}, strings_pkg.StrLen},
	"strings_pkg/StrCompare":  {stringsPkgSrc, "StrCompare", []any{"abc", "abd"}, strings_pkg.StrCompare},
	"strings_pkg/StrEqual":    {stringsPkgSrc, "StrEqual", []any{"hello", "hello"}, strings_pkg.StrEqual},

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
	"switch/Classify":    {switchSrc, "Classify", []any{2}, switch_pkg.Classify},
	"switch/Weekday":     {switchSrc, "Weekday", []any{3}, switch_pkg.Weekday},
	"switch/Grade":       {switchSrc, "Grade", []any{85}, switch_pkg.Grade},
	"switch/ColorCode":   {switchSrc, "ColorCode", []any{"green"}, switch_pkg.ColorCode},

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
	"tricky/ModuloWithNegative":    {trickySrc, "ModuloWithNegative", []any{-17, 5}, tricky.ModuloWithNegative},
	"tricky/PowerRecursive":        {trickySrc, "PowerRecursive", []any{2, 10}, tricky.PowerRecursive},
	"tricky/GCD":                   {trickySrc, "GCD", []any{48, 18}, tricky.GCD},
	"tricky/LCM":                   {trickySrc, "LCM", []any{4, 6}, tricky.LCM},
	"tricky/SumOfDigits":           {trickySrc, "SumOfDigits", []any{12345}, tricky.SumOfDigits},
	"tricky/MapEqualParam":         {trickySrc, "MapEqualParam", []any{map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2}}, tricky.MapEqualParam},
	"tricky/StringJoinParam":       {trickySrc, "StringJoinParam", []any{[]string{"a", "b", "c"}, ","}, tricky.StringJoinParam},
	"tricky/SliceEqualParam":       {trickySrc, "SliceEqualParam", []any{[]int{1, 2, 3}, []int{1, 2, 3}}, tricky.SliceEqualParam},
	"tricky/MapInvertParam":        {trickySrc, "MapInvertParam", []any{map[string]int{"a": 1, "b": 2, "c": 3}}, tricky.MapInvertParam},
	// NOTE: Known failing tests moved to known_issue_test.go

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
	"variables/SumThree":    {variablesSrc, "SumThree", []any{10, 20, 30}, variables.SumThree},
	"variables/Multiply":    {variablesSrc, "Multiply", []any{6, 7}, variables.Multiply},
	"variables/Max":         {variablesSrc, "Max", []any{100, 42}, variables.Max},
	"variables/IsPositive":  {variablesSrc, "IsPositive", []any{5}, variables.IsPositive},

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
