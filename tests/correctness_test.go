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
	"git.woa.com/youngjin/gig/value"
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
	native   func() any
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

	// Handle value.Value types
	if v, ok := got.(value.Value); ok {
		got = v.Interface()
	}

	// Handle []value.Value (multiple return values)
	if gotSlice, ok := got.([]value.Value); ok {
		if expSlice, ok := expected.([]any); ok && len(gotSlice) == len(expSlice) {
			for i := range gotSlice {
				compareCorrectnessResults(t, gotSlice[i].Interface(), expSlice[i])
			}
			return
		}
		// Single value in slice
		if len(gotSlice) == 1 {
			compareCorrectnessResults(t, gotSlice[0].Interface(), expected)
			return
		}
	}

	// Handle []any (multiple return values from native)
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
	"algorithms/InsertionSort":     {algorithmsSrc, "InsertionSort", func() any { return algorithms.InsertionSort() }},
	"algorithms/SelectionSort":     {algorithmsSrc, "SelectionSort", func() any { return algorithms.SelectionSort() }},
	"algorithms/ReverseSlice":      {algorithmsSrc, "ReverseSlice", func() any { return algorithms.ReverseSlice() }},
	"algorithms/IsPalindrome":      {algorithmsSrc, "IsPalindrome", func() any { return algorithms.IsPalindrome() }},
	"algorithms/PowerFunction":     {algorithmsSrc, "PowerFunction", func() any { return algorithms.PowerFunction() }},
	"algorithms/MaxSubarraySum":    {algorithmsSrc, "MaxSubarraySum", func() any { return algorithms.MaxSubarraySum() }},
	"algorithms/TwoSum":            {algorithmsSrc, "TwoSum", func() any { return algorithms.TwoSum() }},
	"algorithms/FibMemoized":       {algorithmsSrc, "FibMemoized", func() any { return algorithms.FibMemoized() }},
	"algorithms/CountDigits":       {algorithmsSrc, "CountDigits", func() any { return algorithms.CountDigits() }},
	"algorithms/CollatzConjecture": {algorithmsSrc, "CollatzConjecture", func() any { return algorithms.CollatzConjecture() }},

	// ============================================================================
	// advanced
	// ============================================================================
	"advanced/TypeConvertIntIdentity": {advancedSrc, "TypeConvertIntIdentity", func() any { return advanced.TypeConvertIntIdentity() }},
	"advanced/DeepCallChain":          {advancedSrc, "DeepCallChain", func() any { return advanced.DeepCallChain() }},
	"advanced/EarlyReturn":            {advancedSrc, "EarlyReturn", func() any { return advanced.EarlyReturn() }},
	"advanced/NestedIfInLoop":         {advancedSrc, "NestedIfInLoop", func() any { return advanced.NestedIfInLoop() }},
	"advanced/BubbleSort":             {advancedSrc, "BubbleSort", func() any { return advanced.BubbleSort() }},
	"advanced/BinarySearch":           {advancedSrc, "BinarySearch", func() any { return advanced.BinarySearch() }},
	"advanced/GCD":                    {advancedSrc, "GCD", func() any { return advanced.GCD() }},
	"advanced/SieveOfEratosthenes":    {advancedSrc, "SieveOfEratosthenes", func() any { return advanced.SieveOfEratosthenes() }},
	"advanced/MatrixMultiply":         {advancedSrc, "MatrixMultiply", func() any { return advanced.MatrixMultiply() }},
	"advanced/EmptyFunctionReturn":    {advancedSrc, "EmptyFunctionReturn", func() any { return advanced.EmptyFunctionReturn() }},
	"advanced/SingleReturnValue":      {advancedSrc, "SingleReturnValue", func() any { return advanced.SingleReturnValue() }},
	"advanced/ZeroIteration":          {advancedSrc, "ZeroIteration", func() any { return advanced.ZeroIteration() }},
	"advanced/LargeLoop":              {advancedSrc, "LargeLoop", func() any { return advanced.LargeLoop() }},
	"advanced/DeepRecursion":          {advancedSrc, "DeepRecursion", func() any { return advanced.DeepRecursion() }},
	"advanced/MapWithClosure":         {advancedSrc, "MapWithClosure", func() any { return advanced.MapWithClosure() }},
	"advanced/SliceWithMultiReturn":   {advancedSrc, "SliceWithMultiReturn", func() any { return advanced.SliceWithMultiReturn() }},
	"advanced/RecursiveDataBuild":     {advancedSrc, "RecursiveDataBuild", func() any { return advanced.RecursiveDataBuild() }},
	"advanced/FunctionChain":          {advancedSrc, "FunctionChain", func() any { return advanced.FunctionChain() }},
	"advanced/ComplexExpressions":     {advancedSrc, "ComplexExpressions", func() any { return advanced.ComplexExpressions() }},

	// ============================================================================
	// arithmetic
	// ============================================================================
	"arithmetic/Addition":       {arithmeticSrc, "Addition", func() any { return arithmetic.Addition() }},
	"arithmetic/Subtraction":    {arithmeticSrc, "Subtraction", func() any { return arithmetic.Subtraction() }},
	"arithmetic/Multiplication": {arithmeticSrc, "Multiplication", func() any { return arithmetic.Multiplication() }},
	"arithmetic/Division":       {arithmeticSrc, "Division", func() any { return arithmetic.Division() }},
	"arithmetic/Modulo":         {arithmeticSrc, "Modulo", func() any { return arithmetic.Modulo() }},
	"arithmetic/ComplexExpr":    {arithmeticSrc, "ComplexExpr", func() any { return arithmetic.ComplexExpr() }},
	"arithmetic/Negation":       {arithmeticSrc, "Negation", func() any { return arithmetic.Negation() }},
	"arithmetic/ChainedOps":     {arithmeticSrc, "ChainedOps", func() any { return arithmetic.ChainedOps() }},
	"arithmetic/Overflow":       {arithmeticSrc, "Overflow", func() any { return arithmetic.Overflow() }},
	"arithmetic/Precedence":     {arithmeticSrc, "Precedence", func() any { return arithmetic.Precedence() }},

	// ============================================================================
	// autowrap
	// ============================================================================
	"autowrap/WithPackage": {autowrapSrc, "WithPackage", func() any { return autowrap.WithPackage() }},
	"autowrap/WithImport":  {autowrapSrc, "WithImport", func() any { return autowrap.WithImport() }},
	"autowrap/Compute":     {autowrapSrc, "Compute", func() any { return autowrap.Compute() }},

	// ============================================================================
	// bitwise
	// ============================================================================
	"bitwise/And":        {bitwiseSrc, "And", func() any { return bitwise.And() }},
	"bitwise/Or":         {bitwiseSrc, "Or", func() any { return bitwise.Or() }},
	"bitwise/Xor":        {bitwiseSrc, "Xor", func() any { return bitwise.Xor() }},
	"bitwise/LeftShift":  {bitwiseSrc, "LeftShift", func() any { return bitwise.LeftShift() }},
	"bitwise/RightShift": {bitwiseSrc, "RightShift", func() any { return bitwise.RightShift() }},
	"bitwise/Combined":   {bitwiseSrc, "Combined", func() any { return bitwise.Combined() }},
	"bitwise/AndNot":     {bitwiseSrc, "AndNot", func() any { return bitwise.AndNot() }},
	"bitwise/PowerOfTwo": {bitwiseSrc, "PowerOfTwo", func() any { return bitwise.PowerOfTwo() }},

	// ============================================================================
	// closures
	// ============================================================================
	"closures/Counter":           {closuresSrc, "Counter", func() any { return closures.Counter() }},
	"closures/CaptureMutation":   {closuresSrc, "CaptureMutation", func() any { return closures.CaptureMutation() }},
	"closures/Factory":           {closuresSrc, "Factory", func() any { return closures.Factory() }},
	"closures/MultipleInstances": {closuresSrc, "MultipleInstances", func() any { return closures.MultipleInstances() }},
	"closures/OverLoop":          {closuresSrc, "OverLoop", func() any { return closures.OverLoop() }},
	"closures/Chain":             {closuresSrc, "Chain", func() any { return closures.Chain() }},
	"closures/Accumulator":       {closuresSrc, "Accumulator", func() any { return closures.Accumulator() }},

	// ============================================================================
	// closures_advanced
	// ============================================================================
	"closures_advanced/Generator":          {closuresAdvancedSrc, "Generator", func() any { return closures_advanced.Generator() }},
	"closures_advanced/Predicate":          {closuresAdvancedSrc, "Predicate", func() any { return closures_advanced.Predicate() }},
	"closures_advanced/StateMachine":       {closuresAdvancedSrc, "StateMachine", func() any { return closures_advanced.StateMachine() }},
	"closures_advanced/RecursiveHelper":    {closuresAdvancedSrc, "RecursiveHelper", func() any { return closures_advanced.RecursiveHelper() }},
	"closures_advanced/ApplyN":             {closuresAdvancedSrc, "ApplyN", func() any { return closures_advanced.ApplyN() }},
	"closures_advanced/Compose":            {closuresAdvancedSrc, "Compose", func() any { return closures_advanced.Compose() }},
	"closures_advanced/ClosureForLoopTest": {closuresAdvancedSrc, "ClosureForLoopTest", func() any { return closures_advanced.ClosureForLoopTest() }},

	// ============================================================================
	// controlflow
	// ============================================================================
	"controlflow/IfTrue":              {controlflowSrc, "IfTrue", func() any { return controlflow.IfTrue() }},
	"controlflow/IfFalse":             {controlflowSrc, "IfFalse", func() any { return controlflow.IfFalse() }},
	"controlflow/IfElse":              {controlflowSrc, "IfElse", func() any { return controlflow.IfElse() }},
	"controlflow/IfElseChainNegative": {controlflowSrc, "IfElseChainNegative", func() any { return controlflow.IfElseChainNegative() }},
	"controlflow/IfElseChainZero":     {controlflowSrc, "IfElseChainZero", func() any { return controlflow.IfElseChainZero() }},
	"controlflow/IfElseChainPositive": {controlflowSrc, "IfElseChainPositive", func() any { return controlflow.IfElseChainPositive() }},
	"controlflow/ForLoop":             {controlflowSrc, "ForLoop", func() any { return controlflow.ForLoop() }},
	"controlflow/ForConditionOnly":    {controlflowSrc, "ForConditionOnly", func() any { return controlflow.ForConditionOnly() }},
	"controlflow/NestedFor":           {controlflowSrc, "NestedFor", func() any { return controlflow.NestedFor() }},
	"controlflow/ForBreak":            {controlflowSrc, "ForBreak", func() any { return controlflow.ForBreak() }},
	"controlflow/ForContinue":         {controlflowSrc, "ForContinue", func() any { return controlflow.ForContinue() }},
	"controlflow/BooleanAndOr":        {controlflowSrc, "BooleanAndOr", func() any { return controlflow.BooleanAndOr() }},

	// ============================================================================
	// cornercases
	// ============================================================================
	"cornercases/ZeroValue_Int":          {cornercasesSrc, "ZeroValue_Int", func() any { return cornercases.ZeroValue_Int() }},
	"cornercases/ZeroValue_Int64":        {cornercasesSrc, "ZeroValue_Int64", func() any { return cornercases.ZeroValue_Int64() }},
	"cornercases/ZeroValue_Float64":      {cornercasesSrc, "ZeroValue_Float64", func() any { return cornercases.ZeroValue_Float64() }},
	"cornercases/ZeroValue_String":       {cornercasesSrc, "ZeroValue_String", func() any { return cornercases.ZeroValue_String() }},
	"cornercases/ZeroValue_Bool":         {cornercasesSrc, "ZeroValue_Bool", func() any { return cornercases.ZeroValue_Bool() }},
	"cornercases/ZeroValue_Slice":        {cornercasesSrc, "ZeroValue_Slice", func() any { return cornercases.ZeroValue_Slice() }},
	"cornercases/ZeroValue_Map":          {cornercasesSrc, "ZeroValue_Map", func() any { return cornercases.ZeroValue_Map() }},
	"cornercases/IntBoundary_MaxInt32":   {cornercasesSrc, "IntBoundary_MaxInt32", func() any { return cornercases.IntBoundary_MaxInt32() }},
	"cornercases/IntBoundary_MinInt32":   {cornercasesSrc, "IntBoundary_MinInt32", func() any { return cornercases.IntBoundary_MinInt32() }},
	"cornercases/IntBoundary_MaxInt64":   {cornercasesSrc, "IntBoundary_MaxInt64", func() any { return cornercases.IntBoundary_MaxInt64() }},
	"cornercases/IntBoundary_MinInt64":   {cornercasesSrc, "IntBoundary_MinInt64", func() any { return cornercases.IntBoundary_MinInt64() }},
	"cornercases/IntBoundary_MaxUint32":  {cornercasesSrc, "IntBoundary_MaxUint32", func() any { return cornercases.IntBoundary_MaxUint32() }},
	"cornercases/IntBoundary_NearMaxInt": {cornercasesSrc, "IntBoundary_NearMaxInt", func() any { return cornercases.IntBoundary_NearMaxInt() }},
	"cornercases/IntBoundary_NearMinInt": {cornercasesSrc, "IntBoundary_NearMinInt", func() any { return cornercases.IntBoundary_NearMinInt() }},
	// Note: int32 overflow wraps just like native Go (two's complement)
	// We compare the actual native results to verify correctness
	"cornercases/Overflow_Int32Add":           {cornercasesSrc, "Overflow_Int32Add", func() any { return cornercases.Overflow_Int32Add() }},
	"cornercases/Overflow_Int32Sub":           {cornercasesSrc, "Overflow_Int32Sub", func() any { return cornercases.Overflow_Int32Sub() }},
	"cornercases/Overflow_Int32Mul":           {cornercasesSrc, "Overflow_Int32Mul", func() any { return cornercases.Overflow_Int32Mul() }},
	"cornercases/FloatBoundary_SmallPositive": {cornercasesSrc, "FloatBoundary_SmallPositive", func() any { return cornercases.FloatBoundary_SmallPositive() }},
	"cornercases/FloatBoundary_SmallNegative": {cornercasesSrc, "FloatBoundary_SmallNegative", func() any { return cornercases.FloatBoundary_SmallNegative() }},
	"cornercases/FloatBoundary_LargePositive": {cornercasesSrc, "FloatBoundary_LargePositive", func() any { return cornercases.FloatBoundary_LargePositive() }},
	"cornercases/FloatBoundary_LargeNegative": {cornercasesSrc, "FloatBoundary_LargeNegative", func() any { return cornercases.FloatBoundary_LargeNegative() }},
	"cornercases/EmptySlice_Len":              {cornercasesSrc, "EmptySlice_Len", func() any { return cornercases.EmptySlice_Len() }},
	"cornercases/EmptySlice_Cap":              {cornercasesSrc, "EmptySlice_Cap", func() any { return cornercases.EmptySlice_Cap() }},
	"cornercases/EmptySlice_Make":             {cornercasesSrc, "EmptySlice_Make", func() any { return cornercases.EmptySlice_Make() }},
	"cornercases/EmptyMap_Len":                {cornercasesSrc, "EmptyMap_Len", func() any { return cornercases.EmptyMap_Len() }},
	"cornercases/EmptyMap_Make":               {cornercasesSrc, "EmptyMap_Make", func() any { return cornercases.EmptyMap_Make() }},
	"cornercases/EmptyString_Len":             {cornercasesSrc, "EmptyString_Len", func() any { return cornercases.EmptyString_Len() }},
	"cornercases/Slice_ZeroToZero":            {cornercasesSrc, "Slice_ZeroToZero", func() any { return cornercases.Slice_ZeroToZero() }},
	"cornercases/Slice_EndToEnd":              {cornercasesSrc, "Slice_EndToEnd", func() any { return cornercases.Slice_EndToEnd() }},
	"cornercases/Slice_NilSlice":              {cornercasesSrc, "Slice_NilSlice", func() any { return cornercases.Slice_NilSlice() }},
	"cornercases/Slice_AppendToNil":           {cornercasesSrc, "Slice_AppendToNil", func() any { return cornercases.Slice_AppendToNil() }},
	"cornercases/Slice_AppendEmpty":           {cornercasesSrc, "Slice_AppendEmpty", func() any { return cornercases.Slice_AppendEmpty() }},
	"cornercases/Map_NilMap":                  {cornercasesSrc, "Map_NilMap", func() any { return cornercases.Map_NilMap() }},
	"cornercases/Map_AccessMissingKey":        {cornercasesSrc, "Map_AccessMissingKey", func() any { return cornercases.Map_AccessMissingKey() }},
	"cornercases/Map_DeleteMissingKey":        {cornercasesSrc, "Map_DeleteMissingKey", func() any { return cornercases.Map_DeleteMissingKey() }},
	"cornercases/Map_OverwriteKey":            {cornercasesSrc, "Map_OverwriteKey", func() any { return cornercases.Map_OverwriteKey() }},
	"cornercases/Map_NilKeyString":            {cornercasesSrc, "Map_NilKeyString", func() any { return cornercases.Map_NilKeyString() }},
	"cornercases/Map_ZeroIntKey":              {cornercasesSrc, "Map_ZeroIntKey", func() any { return cornercases.Map_ZeroIntKey() }},
	"cornercases/String_Empty":                {cornercasesSrc, "String_Empty", func() any { return cornercases.String_Empty() }},
	"cornercases/String_SingleChar":           {cornercasesSrc, "String_SingleChar", func() any { return cornercases.String_SingleChar() }},
	"cornercases/String_UnicodeMultibyte":     {cornercasesSrc, "String_UnicodeMultibyte", func() any { return cornercases.String_UnicodeMultibyte() }},
	"cornercases/String_Whitespace":           {cornercasesSrc, "String_Whitespace", func() any { return cornercases.String_Whitespace() }},
	"cornercases/String_SingleByteIndex":      {cornercasesSrc, "String_SingleByteIndex", func() any { return cornercases.String_SingleByteIndex() }},
	"cornercases/String_LastByte":             {cornercasesSrc, "String_LastByte", func() any { return cornercases.String_LastByte() }},
	"cornercases/Bool_True":                   {cornercasesSrc, "Bool_True", func() any { return cornercases.Bool_True() }},
	"cornercases/Bool_False":                  {cornercasesSrc, "Bool_False", func() any { return cornercases.Bool_False() }},
	"cornercases/Bool_NotTrue":                {cornercasesSrc, "Bool_NotTrue", func() any { return cornercases.Bool_NotTrue() }},
	"cornercases/Bool_NotFalse":               {cornercasesSrc, "Bool_NotFalse", func() any { return cornercases.Bool_NotFalse() }},
	"cornercases/Bool_DoubleNegation":         {cornercasesSrc, "Bool_DoubleNegation", func() any { return cornercases.Bool_DoubleNegation() }},
	"cornercases/Arith_AddZero":               {cornercasesSrc, "Arith_AddZero", func() any { return cornercases.Arith_AddZero() }},
	"cornercases/Arith_SubZero":               {cornercasesSrc, "Arith_SubZero", func() any { return cornercases.Arith_SubZero() }},
	"cornercases/Arith_MulByOne":              {cornercasesSrc, "Arith_MulByOne", func() any { return cornercases.Arith_MulByOne() }},
	"cornercases/Arith_DivByOne":              {cornercasesSrc, "Arith_DivByOne", func() any { return cornercases.Arith_DivByOne() }},
	"cornercases/Arith_ModByOne":              {cornercasesSrc, "Arith_ModByOne", func() any { return cornercases.Arith_ModByOne() }},
	"cornercases/Arith_MulByZero":             {cornercasesSrc, "Arith_MulByZero", func() any { return cornercases.Arith_MulByZero() }},
	"cornercases/Arith_NegNeg":                {cornercasesSrc, "Arith_NegNeg", func() any { return cornercases.Arith_NegNeg() }},
	"cornercases/Arith_NegAddNeg":             {cornercasesSrc, "Arith_NegAddNeg", func() any { return cornercases.Arith_NegAddNeg() }},
	"cornercases/Compare_IntEqual":            {cornercasesSrc, "Compare_IntEqual", func() any { return cornercases.Compare_IntEqual() }},
	"cornercases/Compare_IntNotEqual":         {cornercasesSrc, "Compare_IntNotEqual", func() any { return cornercases.Compare_IntNotEqual() }},
	"cornercases/Compare_IntGreater":          {cornercasesSrc, "Compare_IntGreater", func() any { return cornercases.Compare_IntGreater() }},
	"cornercases/Compare_IntGreaterEqual":     {cornercasesSrc, "Compare_IntGreaterEqual", func() any { return cornercases.Compare_IntGreaterEqual() }},
	"cornercases/Compare_IntLess":             {cornercasesSrc, "Compare_IntLess", func() any { return cornercases.Compare_IntLess() }},
	"cornercases/Compare_IntLessEqual":        {cornercasesSrc, "Compare_IntLessEqual", func() any { return cornercases.Compare_IntLessEqual() }},
	"cornercases/Compare_StringEqual":         {cornercasesSrc, "Compare_StringEqual", func() any { return cornercases.Compare_StringEqual() }},
	"cornercases/Compare_StringNotEqual":      {cornercasesSrc, "Compare_StringNotEqual", func() any { return cornercases.Compare_StringNotEqual() }},
	"cornercases/Compare_EmptyStringEqual":    {cornercasesSrc, "Compare_EmptyStringEqual", func() any { return cornercases.Compare_EmptyStringEqual() }},
	"cornercases/Logic_TrueAndTrue":           {cornercasesSrc, "Logic_TrueAndTrue", func() any { return cornercases.Logic_TrueAndTrue() }},
	"cornercases/Logic_TrueAndFalse":          {cornercasesSrc, "Logic_TrueAndFalse", func() any { return cornercases.Logic_TrueAndFalse() }},
	"cornercases/Logic_FalseAndTrue":          {cornercasesSrc, "Logic_FalseAndTrue", func() any { return cornercases.Logic_FalseAndTrue() }},
	"cornercases/Logic_TrueOrFalse":           {cornercasesSrc, "Logic_TrueOrFalse", func() any { return cornercases.Logic_TrueOrFalse() }},
	"cornercases/Logic_FalseOrTrue":           {cornercasesSrc, "Logic_FalseOrTrue", func() any { return cornercases.Logic_FalseOrTrue() }},
	"cornercases/Logic_FalseOrFalse":          {cornercasesSrc, "Logic_FalseOrFalse", func() any { return cornercases.Logic_FalseOrFalse() }},
	"cornercases/Control_IfNoElse":            {cornercasesSrc, "Control_IfNoElse", func() any { return cornercases.Control_IfNoElse() }},
	"cornercases/Control_IfFalseNoElse":       {cornercasesSrc, "Control_IfFalseNoElse", func() any { return cornercases.Control_IfFalseNoElse() }},
	"cornercases/Control_ForZeroIter":         {cornercasesSrc, "Control_ForZeroIter", func() any { return cornercases.Control_ForZeroIter() }},
	"cornercases/Control_ForOneIter":          {cornercasesSrc, "Control_ForOneIter", func() any { return cornercases.Control_ForOneIter() }},
	"cornercases/Control_ForBreakFirst":       {cornercasesSrc, "Control_ForBreakFirst", func() any { return cornercases.Control_ForBreakFirst() }},
	"cornercases/Control_ForContinueAll":      {cornercasesSrc, "Control_ForContinueAll", func() any { return cornercases.Control_ForContinueAll() }},
	"cornercases/Control_SwitchNoMatch":       {cornercasesSrc, "Control_SwitchNoMatch", func() any { return cornercases.Control_SwitchNoMatch() }},
	"cornercases/Control_SwitchDefault":       {cornercasesSrc, "Control_SwitchDefault", func() any { return cornercases.Control_SwitchDefault() }},
	"cornercases/Func_NoReturn":               {cornercasesSrc, "Func_NoReturn", func() any { return cornercases.Func_NoReturn() }},
	"cornercases/Func_MultipleReturnAll": {cornercasesSrc, "Func_MultipleReturnAll", func() any {
		a, b := cornercases.Func_MultipleReturnAll()
		return []any{a, b}
	}},
	"cornercases/Func_MultipleReturnIgnore": {cornercasesSrc, "Func_MultipleReturnIgnore", func() any { return cornercases.Func_MultipleReturnIgnore() }},
	"cornercases/Func_NamedReturn":          {cornercasesSrc, "Func_NamedReturn", func() any { return cornercases.Func_NamedReturn() }},
	"cornercases/Func_VariadicEmpty":        {cornercasesSrc, "Func_VariadicEmpty", func() any { return cornercases.Func_VariadicEmpty() }},
	"cornercases/Func_VariadicOne":          {cornercasesSrc, "Func_VariadicOne", func() any { return cornercases.Func_VariadicOne() }},
	"cornercases/Func_VariadicMultiple":     {cornercasesSrc, "Func_VariadicMultiple", func() any { return cornercases.Func_VariadicMultiple() }},
	"cornercases/Func_RecursionBase":        {cornercasesSrc, "Func_RecursionBase", func() any { return cornercases.Func_RecursionBase() }},
	"cornercases/Closure_ReturnClosure":     {cornercasesSrc, "Closure_ReturnClosure", func() any { return cornercases.Closure_ReturnClosure() }},
	"cornercases/Closure_CaptureVariable":   {cornercasesSrc, "Closure_CaptureVariable", func() any { return cornercases.Closure_CaptureVariable() }},
	"cornercases/Closure_ModifyCaptured":    {cornercasesSrc, "Closure_ModifyCaptured", func() any { return cornercases.Closure_ModifyCaptured() }},
	"cornercases/Struct_ZeroValueFields":    {cornercasesSrc, "Struct_ZeroValueFields", func() any { return cornercases.Struct_ZeroValueFields() }},
	"cornercases/Struct_PointerReceiver":    {cornercasesSrc, "Struct_PointerReceiver", func() any { return cornercases.Struct_PointerReceiver() }},
	"cornercases/Struct_NestedStruct":       {cornercasesSrc, "Struct_NestedStruct", func() any { return cornercases.Struct_NestedStruct() }},

	// ============================================================================
	// edgecases
	// ============================================================================
	"edgecases/MaxInt64":           {edgecasesSrc, "MaxInt64", func() any { return edgecases.MaxInt64() }},
	"edgecases/MinInt64":           {edgecasesSrc, "MinInt64", func() any { return edgecases.MinInt64() }},
	"edgecases/DivisionByMinusOne": {edgecasesSrc, "DivisionByMinusOne", func() any { return edgecases.DivisionByMinusOne() }},
	"edgecases/ModuloNegative":     {edgecasesSrc, "ModuloNegative", func() any { return edgecases.ModuloNegative() }},
	"edgecases/EmptyString":        {edgecasesSrc, "EmptyString", func() any { return edgecases.EmptyString() }},
	"edgecases/LargeSlice":         {edgecasesSrc, "LargeSlice", func() any { return edgecases.LargeSlice() }},
	"edgecases/NestedMapLookup":    {edgecasesSrc, "NestedMapLookup", func() any { return edgecases.NestedMapLookup() }},
	"edgecases/ZeroDivisionGuard":  {edgecasesSrc, "ZeroDivisionGuard", func() any { return edgecases.ZeroDivisionGuard() }},
	"edgecases/BooleanComplexExpr": {edgecasesSrc, "BooleanComplexExpr", func() any { return edgecases.BooleanComplexExpr() }},
	"edgecases/SingleElementSlice": {edgecasesSrc, "SingleElementSlice", func() any { return edgecases.SingleElementSlice() }},
	"edgecases/EmptyMap":           {edgecasesSrc, "EmptyMap", func() any { return edgecases.EmptyMap() }},
	"edgecases/TightLoop":          {edgecasesSrc, "TightLoop", func() any { return edgecases.TightLoop() }},

	// ============================================================================
	// external
	// ============================================================================
	"external/FmtSprintf":       {externalSrc, "FmtSprintf", func() any { return external.FmtSprintf() }},
	"external/FmtSprintfMulti":  {externalSrc, "FmtSprintfMulti", func() any { return external.FmtSprintfMulti() }},
	"external/StringsToUpper":   {externalSrc, "StringsToUpper", func() any { return external.StringsToUpper() }},
	"external/StringsToLower":   {externalSrc, "StringsToLower", func() any { return external.StringsToLower() }},
	"external/StringsContains":  {externalSrc, "StringsContains", func() any { return external.StringsContains() }},
	"external/StringsReplace":   {externalSrc, "StringsReplace", func() any { return external.StringsReplace() }},
	"external/StringsHasPrefix": {externalSrc, "StringsHasPrefix", func() any { return external.StringsHasPrefix() }},
	"external/StrconvItoa":      {externalSrc, "StrconvItoa", func() any { return external.StrconvItoa() }},
	"external/StrconvAtoi":      {externalSrc, "StrconvAtoi", func() any { return external.StrconvAtoi() }},

	// ============================================================================
	// functions
	// ============================================================================
	"functions/Call":                 {functionsSrc, "Call", func() any { return functions.Call() }},
	"functions/MultipleReturn":       {functionsSrc, "MultipleReturn", func() any { return functions.MultipleReturn() }},
	"functions/MultipleReturnDivmod": {functionsSrc, "MultipleReturnDivmod", func() any { return functions.MultipleReturnDivmod() }},
	"functions/RecursionFactorial":   {functionsSrc, "RecursionFactorial", func() any { return functions.RecursionFactorial() }},
	"functions/MutualRecursion":      {functionsSrc, "MutualRecursion", func() any { return functions.MutualRecursion() }},
	"functions/FibonacciIterative":   {functionsSrc, "FibonacciIterative", func() any { return functions.FibonacciIterative() }},
	"functions/FibonacciRecursive":   {functionsSrc, "FibonacciRecursive", func() any { return functions.FibonacciRecursive() }},
	"functions/VariadicFunction":     {functionsSrc, "VariadicFunction", func() any { return functions.VariadicFunction() }},
	"functions/FunctionAsValue":      {functionsSrc, "FunctionAsValue", func() any { return functions.FunctionAsValue() }},
	"functions/HigherOrderMap":       {functionsSrc, "HigherOrderMap", func() any { return functions.HigherOrderMap() }},
	"functions/HigherOrderFilter":    {functionsSrc, "HigherOrderFilter", func() any { return functions.HigherOrderFilter() }},
	"functions/HigherOrderReduce":    {functionsSrc, "HigherOrderReduce", func() any { return functions.HigherOrderReduce() }},

	// ============================================================================
	// initialize - Complex initialization tests
	// ============================================================================
	"initialize/ComplexInitTest":     {initializeSrc, "ComplexInitTest", func() any { return initialize.ComplexInitTest() }},
	"initialize/InitOrderTest":       {initializeSrc, "InitOrderTest", func() any { return initialize.InitOrderTest() }},
	"initialize/CacheInitTest":       {initializeSrc, "CacheInitTest", func() any { return initialize.CacheInitTest() }},
	"initialize/LookupTableInitTest": {initializeSrc, "LookupTableInitTest", func() any { return initialize.LookupTableInitTest() }},
	"initialize/FibonacciInitTest":   {initializeSrc, "FibonacciInitTest", func() any { return initialize.FibonacciInitTest() }},
	"initialize/GetA":                {initializeSrc, "GetA", func() any { return initialize.GetA() }},
	"initialize/GetB":                {initializeSrc, "GetB", func() any { return initialize.GetB() }},
	"initialize/GetC":                {initializeSrc, "GetC", func() any { return initialize.GetC() }},
	"initialize/GetCacheSum":         {initializeSrc, "GetCacheSum", func() any { return initialize.GetCacheSum() }},
	"initialize/GetCacheSize":        {initializeSrc, "GetCacheSize", func() any { return initialize.GetCacheSize() }},
	"initialize/GetFibonacciCount":   {initializeSrc, "GetFibonacciCount", func() any { return initialize.GetFibonacciCount() }},

	// ============================================================================
	// leetcode_hard
	// ============================================================================
	"leetcode_hard/TrappingRainWater":           {leetcodeHardSrc, "TrappingRainWater", func() any { return leetcode_hard.TrappingRainWater() }},
	"leetcode_hard/LargestRectangleInHistogram": {leetcodeHardSrc, "LargestRectangleInHistogram", func() any { return leetcode_hard.LargestRectangleInHistogram() }},
	"leetcode_hard/MedianOfTwoSortedArrays":     {leetcodeHardSrc, "MedianOfTwoSortedArrays", func() any { return leetcode_hard.MedianOfTwoSortedArrays() }},
	"leetcode_hard/RegularExpressionMatching":   {leetcodeHardSrc, "RegularExpressionMatching", func() any { return leetcode_hard.RegularExpressionMatching() }},
	"leetcode_hard/NQueens":                     {leetcodeHardSrc, "NQueens", func() any { return leetcode_hard.NQueens() }},
	"leetcode_hard/LongestIncreasingPath":       {leetcodeHardSrc, "LongestIncreasingPath", func() any { return leetcode_hard.LongestIncreasingPath() }},
	"leetcode_hard/WordLadder":                  {leetcodeHardSrc, "WordLadder", func() any { return leetcode_hard.WordLadder() }},
	"leetcode_hard/MergeKSortedLists":           {leetcodeHardSrc, "MergeKSortedLists", func() any { return leetcode_hard.MergeKSortedLists() }},
	"leetcode_hard/EditDistance":                {leetcodeHardSrc, "EditDistance", func() any { return leetcode_hard.EditDistance() }},
	"leetcode_hard/MinimumWindowSubstring":      {leetcodeHardSrc, "MinimumWindowSubstring", func() any { return leetcode_hard.MinimumWindowSubstring() }},

	// ============================================================================
	// maps
	// ============================================================================
	"maps/BasicOps":       {mapsSrc, "BasicOps", func() any { return maps.BasicOps() }},
	"maps/Iteration":      {mapsSrc, "Iteration", func() any { return maps.Iteration() }},
	"maps/Delete":         {mapsSrc, "Delete", func() any { return maps.Delete() }},
	"maps/Len":            {mapsSrc, "Len", func() any { return maps.Len() }},
	"maps/Overwrite":      {mapsSrc, "Overwrite", func() any { return maps.Overwrite() }},
	"maps/IntKeys":        {mapsSrc, "IntKeys", func() any { return maps.IntKeys() }},
	"maps/PassToFunction": {mapsSrc, "PassToFunction", func() any { return maps.PassToFunction() }},

	// ============================================================================
	// mapadvanced
	// ============================================================================
	"mapadvanced/LookupExistingKey": {mapadvancedSrc, "LookupExistingKey", func() any { return mapadvanced.LookupExistingKey() }},
	"mapadvanced/LookupWithDefault": {mapadvancedSrc, "LookupWithDefault", func() any { return mapadvanced.LookupWithDefault() }},
	"mapadvanced/AsCounter":         {mapadvancedSrc, "AsCounter", func() any { return mapadvanced.AsCounter() }},
	"mapadvanced/WithStringValues":  {mapadvancedSrc, "WithStringValues", func() any { return mapadvanced.WithStringValues() }},
	"mapadvanced/BuildFromLoop":     {mapadvancedSrc, "BuildFromLoop", func() any { return mapadvanced.BuildFromLoop() }},
	"mapadvanced/DeleteAndReinsert": {mapadvancedSrc, "DeleteAndReinsert", func() any { return mapadvanced.DeleteAndReinsert() }},

	// ============================================================================
	// multiassign
	// ============================================================================
	"multiassign/Swap":             {multiassignSrc, "Swap", func() any { return multiassign.Swap() }},
	"multiassign/FromFunction":     {multiassignSrc, "FromFunction", func() any { return multiassign.FromFunction() }},
	"multiassign/ThreeValues":      {multiassignSrc, "ThreeValues", func() any { return multiassign.ThreeValues() }},
	"multiassign/InLoop":           {multiassignSrc, "InLoop", func() any { return multiassign.InLoop() }},
	"multiassign/DiscardWithBlank": {multiassignSrc, "DiscardWithBlank", func() any { return multiassign.DiscardWithBlank() }},

	// ============================================================================
	// namedreturn
	// ============================================================================
	"namedreturn/Basic":     {namedreturnSrc, "Basic", func() any { return namedreturn.Basic() }},
	"namedreturn/Multiple":  {namedreturnSrc, "Multiple", func() any { return namedreturn.Multiple() }},
	"namedreturn/ZeroValue": {namedreturnSrc, "ZeroValue", func() any { return namedreturn.ZeroValue() }},

	// ============================================================================
	// recursion
	// ============================================================================
	"recursion/TailRecursionPattern": {recursionSrc, "TailRecursionPattern", func() any { return recursion.TailRecursionPattern() }},
	"recursion/ReverseSlice":         {recursionSrc, "ReverseSlice", func() any { return recursion.ReverseSlice() }},
	"recursion/TowerOfHanoi":         {recursionSrc, "TowerOfHanoi", func() any { return recursion.TowerOfHanoi() }},
	"recursion/MaxSlice":             {recursionSrc, "MaxSlice", func() any { return recursion.MaxSlice() }},
	"recursion/Ackermann":            {recursionSrc, "Ackermann", func() any { return recursion.Ackermann() }},
	"recursion/BinarySearch":         {recursionSrc, "BinarySearch", func() any { return recursion.BinarySearch() }},

	// ============================================================================
	// scope
	// ============================================================================
	"scope/IfInitShortVar":            {scopeSrc, "IfInitShortVar", func() any { return scope.IfInitShortVar() }},
	"scope/IfInitMultiCondition":      {scopeSrc, "IfInitMultiCondition", func() any { return scope.IfInitMultiCondition() }},
	"scope/NestedScopes":              {scopeSrc, "NestedScopes", func() any { return scope.NestedScopes() }},
	"scope/ForScopeIsolation":         {scopeSrc, "ForScopeIsolation", func() any { return scope.ForScopeIsolation() }},
	"scope/MultipleBlockScopes":       {scopeSrc, "MultipleBlockScopes", func() any { return scope.MultipleBlockScopes() }},
	"scope/ClosureCapturesOuterScope": {scopeSrc, "ClosureCapturesOuterScope", func() any { return scope.ClosureCapturesOuterScope() }},

	// ============================================================================
	// slices
	// ============================================================================
	"slices/MakeLen":           {slicesSrc, "MakeLen", func() any { return slices.MakeLen() }},
	"slices/Append":            {slicesSrc, "Append", func() any { return slices.Append() }},
	"slices/ElementAssignment": {slicesSrc, "ElementAssignment", func() any { return slices.ElementAssignment() }},
	"slices/ForRange":          {slicesSrc, "ForRange", func() any { return slices.ForRange() }},
	"slices/ForRangeIndex":     {slicesSrc, "ForRangeIndex", func() any { return slices.ForRangeIndex() }},
	"slices/GrowMultiple":      {slicesSrc, "GrowMultiple", func() any { return slices.GrowMultiple() }},
	"slices/PassToFunction":    {slicesSrc, "PassToFunction", func() any { return slices.PassToFunction() }},
	"slices/LenCap":            {slicesSrc, "LenCap", func() any { return slices.LenCap() }},

	// ============================================================================
	// slicing
	// ============================================================================
	"slicing/SubSliceBasic":            {slicingSrc, "SubSliceBasic", func() any { return slicing.SubSliceBasic() }},
	"slicing/SubSliceLen":              {slicingSrc, "SubSliceLen", func() any { return slicing.SubSliceLen() }},
	"slicing/SubSliceFromStart":        {slicingSrc, "SubSliceFromStart", func() any { return slicing.SubSliceFromStart() }},
	"slicing/SubSliceToEnd":            {slicingSrc, "SubSliceToEnd", func() any { return slicing.SubSliceToEnd() }},
	"slicing/SubSliceCopy":             {slicingSrc, "SubSliceCopy", func() any { return slicing.SubSliceCopy() }},
	"slicing/SubSliceChained":          {slicingSrc, "SubSliceChained", func() any { return slicing.SubSliceChained() }},
	"slicing/SubSliceModifiesOriginal": {slicingSrc, "SubSliceModifiesOriginal", func() any { return slicing.SubSliceModifiesOriginal() }},

	// ============================================================================
	// strings_pkg
	// ============================================================================
	"strings_pkg/Concat":     {stringsPkgSrc, "Concat", func() any { return strings_pkg.Concat() }},
	"strings_pkg/ConcatLoop": {stringsPkgSrc, "ConcatLoop", func() any { return strings_pkg.ConcatLoop() }},
	"strings_pkg/Len":        {stringsPkgSrc, "Len", func() any { return strings_pkg.Len() }},
	"strings_pkg/Index":      {stringsPkgSrc, "Index", func() any { return strings_pkg.Index() }},
	"strings_pkg/Comparison": {stringsPkgSrc, "Comparison", func() any { return strings_pkg.Comparison() }},
	"strings_pkg/Equality":   {stringsPkgSrc, "Equality", func() any { return strings_pkg.Equality() }},
	"strings_pkg/EmptyCheck": {stringsPkgSrc, "EmptyCheck", func() any { return strings_pkg.EmptyCheck() }},

	// ============================================================================
	// structs
	// ============================================================================
	"structs/BasicStruct":                 {structsSrc, "BasicStruct", func() any { return structs.BasicStruct() }},
	"structs/StructPointer":               {structsSrc, "StructPointer", func() any { return structs.StructPointer() }},
	"structs/NestedStruct":                {structsSrc, "NestedStruct", func() any { return structs.NestedStruct() }},
	"structs/EmbeddedField":               {structsSrc, "EmbeddedField", func() any { return structs.EmbeddedField() }},
	"structs/StructInSlice":               {structsSrc, "StructInSlice", func() any { return structs.StructInSlice() }},
	"structs/StructAsParam":               {structsSrc, "StructAsParam", func() any { return structs.StructAsParam() }},
	"structs/StructZeroValue":             {structsSrc, "StructZeroValue", func() any { return structs.StructZeroValue() }},
	"structs/MultipleEmbedded":            {structsSrc, "MultipleEmbedded", func() any { return structs.MultipleEmbedded() }},
	"structs/DeepNesting":                 {structsSrc, "DeepNesting", func() any { return structs.DeepNesting() }},
	"structs/StructFieldMutation":         {structsSrc, "StructFieldMutation", func() any { return structs.StructFieldMutation() }},
	"structs/StructWithBool":              {structsSrc, "StructWithBool", func() any { return structs.StructWithBool() }},
	"structs/StructCopySemantics":         {structsSrc, "StructCopySemantics", func() any { return structs.StructCopySemantics() }},
	"structs/StructPointerSharing":        {structsSrc, "StructPointerSharing", func() any { return structs.StructPointerSharing() }},
	"structs/StructReturnFromFunc":        {structsSrc, "StructReturnFromFunc", func() any { return structs.StructReturnFromFunc() }},
	"structs/StructPointerReturnFromFunc": {structsSrc, "StructPointerReturnFromFunc", func() any { return structs.StructPointerReturnFromFunc() }},
	"structs/StructSliceAppend":           {structsSrc, "StructSliceAppend", func() any { return structs.StructSliceAppend() }},
	"structs/StructPointerSlice":          {structsSrc, "StructPointerSlice", func() any { return structs.StructPointerSlice() }},
	"structs/StructInMap":                 {structsSrc, "StructInMap", func() any { return structs.StructInMap() }},
	"structs/StructConditionalInit":       {structsSrc, "StructConditionalInit", func() any { return structs.StructConditionalInit() }},
	"structs/StructFieldLoop":             {structsSrc, "StructFieldLoop", func() any { return structs.StructFieldLoop() }},
	"structs/StructNestedMutation":        {structsSrc, "StructNestedMutation", func() any { return structs.StructNestedMutation() }},
	"structs/StructEmbeddedOverride":      {structsSrc, "StructEmbeddedOverride", func() any { return structs.StructEmbeddedOverride() }},
	"structs/StructWithClosure":           {structsSrc, "StructWithClosure", func() any { return structs.StructWithClosure() }},
	"structs/StructReassign":              {structsSrc, "StructReassign", func() any { return structs.StructReassign() }},
	"structs/StructSliceOfNested":         {structsSrc, "StructSliceOfNested", func() any { return structs.StructSliceOfNested() }},
	"structs/StructMultiReturn":           {structsSrc, "StructMultiReturn", func() any { return structs.StructMultiReturn() }},
	"structs/StructBuilderPattern":        {structsSrc, "StructBuilderPattern", func() any { return structs.StructBuilderPattern() }},
	"structs/StructArrayField":            {structsSrc, "StructArrayField", func() any { return structs.StructArrayField() }},
	"structs/StructEmbeddedChain":         {structsSrc, "StructEmbeddedChain", func() any { return structs.StructEmbeddedChain() }},

	// ============================================================================
	// switch
	// ============================================================================
	"switch/Simple":      {switchSrc, "Simple", func() any { return switch_pkg.Simple() }},
	"switch/Default":     {switchSrc, "Default", func() any { return switch_pkg.Default() }},
	"switch/MultiCase":   {switchSrc, "MultiCase", func() any { return switch_pkg.MultiCase() }},
	"switch/NoCondition": {switchSrc, "NoCondition", func() any { return switch_pkg.NoCondition() }},
	"switch/WithInit":    {switchSrc, "WithInit", func() any { return switch_pkg.WithInit() }},
	"switch/StringCases": {switchSrc, "StringCases", func() any { return switch_pkg.StringCases() }},
	"switch/Fallthrough": {switchSrc, "Fallthrough", func() any { return switch_pkg.Fallthrough() }},
	"switch/Nested":      {switchSrc, "Nested", func() any { return switch_pkg.Nested() }},

	// ============================================================================
	// tricky - Iteration 1
	// ============================================================================
	"tricky/ShortVarDeclShadow":  {trickySrc, "ShortVarDeclShadow", func() any { return tricky.ShortVarDeclShadow() }},
	"tricky/SliceIndexExpr":      {trickySrc, "SliceIndexExpr", func() any { return tricky.SliceIndexExpr() }},
	"tricky/MapStructKey":        {trickySrc, "MapStructKey", func() any { return tricky.MapStructKey() }},
	"tricky/NestedSliceAppend":   {trickySrc, "NestedSliceAppend", func() any { return tricky.NestedSliceAppend() }},
	"tricky/ClosureCaptureLoop":  {trickySrc, "ClosureCaptureLoop", func() any { return tricky.ClosureCaptureLoop() }},
	"tricky/DeferNamedReturn":    {trickySrc, "DeferNamedReturn", func() any { return tricky.DeferNamedReturn() }},
	"tricky/FullSliceExpr":       {trickySrc, "FullSliceExpr", func() any { return tricky.FullSliceExpr() }},
	"tricky/NestedShadowing":     {trickySrc, "NestedShadowing", func() any { return tricky.NestedShadowing() }},
	"tricky/SliceOfPointers":     {trickySrc, "SliceOfPointers", func() any { return tricky.SliceOfPointers() }},
	"tricky/MapNestedStruct":     {trickySrc, "MapNestedStruct", func() any { return tricky.MapNestedStruct() }},
	"tricky/VariadicEmpty":       {trickySrc, "VariadicEmpty", func() any { return tricky.VariadicEmpty() }},
	"tricky/VariadicOne":         {trickySrc, "VariadicOne", func() any { return tricky.VariadicOne() }},
	"tricky/VariadicMultiple":    {trickySrc, "VariadicMultiple", func() any { return tricky.VariadicMultiple() }},
	"tricky/EmbeddedField":       {trickySrc, "EmbeddedField", func() any { return tricky.EmbeddedField() }},
	"tricky/MapPointerValue":     {trickySrc, "MapPointerValue", func() any { return tricky.MapPointerValue() }},
	"tricky/ComplexBoolExpr":     {trickySrc, "ComplexBoolExpr", func() any { return tricky.ComplexBoolExpr() }},
	"tricky/SwitchFallthrough":   {trickySrc, "SwitchFallthrough", func() any { return tricky.SwitchFallthrough() }},
	"tricky/SliceCopyOperation":  {trickySrc, "SliceCopyOperation", func() any { return tricky.SliceCopyOperation() }},
	"tricky/DeferStackOrder":     {trickySrc, "DeferStackOrder", func() any { return tricky.DeferStackOrder() }},
	"tricky/InterfaceAssertion":  {trickySrc, "InterfaceAssertion", func() any { return tricky.InterfaceAssertion() }},
	"tricky/ChannelBasic":        {trickySrc, "ChannelBasic", func() any { return tricky.ChannelBasic() }},
	"tricky/SelectDefault":       {trickySrc, "SelectDefault", func() any { return tricky.SelectDefault() }},
	"tricky/RecursiveFibMemo":    {trickySrc, "RecursiveFibMemo", func() any { return tricky.RecursiveFibMemo() }},
	"tricky/PanicRecover":        {trickySrc, "PanicRecover", func() any { return tricky.PanicRecover() }},
	"tricky/ClosureWithDefer":    {trickySrc, "ClosureWithDefer", func() any { return tricky.ClosureWithDefer() }},
	"tricky/MethodOnPointer":     {trickySrc, "MethodOnPointer", func() any { return tricky.MethodOnPointer() }},
	"tricky/MultiReturnDiscard":  {trickySrc, "MultiReturnDiscard", func() any { return tricky.MultiReturnDiscard() }},
	"tricky/NilSliceAppend":      {trickySrc, "NilSliceAppend", func() any { return tricky.NilSliceAppend() }},
	"tricky/ShortCircuitEval":    {trickySrc, "ShortCircuitEval", func() any { return tricky.ShortCircuitEval() }},
	"tricky/ShortCircuitEval2":   {trickySrc, "ShortCircuitEval2", func() any { return tricky.ShortCircuitEval2() }},
	"tricky/MapDelete":           {trickySrc, "MapDelete", func() any { return tricky.MapDelete() }},
	"tricky/SliceNil":            {trickySrc, "SliceNil", func() any { return tricky.SliceNil() }},
	"tricky/MapCommaOk":          {trickySrc, "MapCommaOk", func() any { return tricky.MapCommaOk() }},
	"tricky/InterfaceNil":        {trickySrc, "InterfaceNil", func() any { return tricky.InterfaceNil() }},
	"tricky/SliceLenCap":         {trickySrc, "SliceLenCap", func() any { return tricky.SliceLenCap() }},
	"tricky/ComplexArray":        {trickySrc, "ComplexArray", func() any { return tricky.ComplexArray() }},
	"tricky/PointerArithmetic":   {trickySrc, "PointerArithmetic", func() any { return tricky.PointerArithmetic() }},
	"tricky/DoublePointer":       {trickySrc, "DoublePointer", func() any { return tricky.DoublePointer() }},
	"tricky/StructPointerMethod": {trickySrc, "StructPointerMethod", func() any { return tricky.StructPointerMethod() }},
	"tricky/ForRangeWithIndex":   {trickySrc, "ForRangeWithIndex", func() any { return tricky.ForRangeWithIndex() }},
	"tricky/ForRangeKeyValue":    {trickySrc, "ForRangeKeyValue", func() any { return tricky.ForRangeKeyValue() }},
	"tricky/StringIndex":         {trickySrc, "StringIndex", func() any { return tricky.StringIndex() }},
	"tricky/MapAssign":           {trickySrc, "MapAssign", func() any { return tricky.MapAssign() }},
	"tricky/ComplexLiteral":      {trickySrc, "ComplexLiteral", func() any { return tricky.ComplexLiteral() }},
	"tricky/ErrorReturn":         {trickySrc, "ErrorReturn", func() any { return tricky.ErrorReturn() }},
	"tricky/NilPointerCheck":     {trickySrc, "NilPointerCheck", func() any { return tricky.NilPointerCheck() }},
	"tricky/SliceAppendNil":      {trickySrc, "SliceAppendNil", func() any { return tricky.SliceAppendNil() }},
	"tricky/MapLookupNil":        {trickySrc, "MapLookupNil", func() any { return tricky.MapLookupNil() }},
	"tricky/DeferModifyNamed":    {trickySrc, "DeferModifyNamed", func() any { return tricky.DeferModifyNamed() }},
	"tricky/ForRangeMap":         {trickySrc, "ForRangeMap", func() any { return tricky.ForRangeMap() }},
	"tricky/MultipleNamedReturn": {trickySrc, "MultipleNamedReturnCombined", func() any { return tricky.MultipleNamedReturnCombined() }},

	// ============================================================================
	// tricky - Iteration 2
	// ============================================================================
	"tricky/DeferInClosure":         {trickySrc, "DeferInClosure", func() any { return tricky.DeferInClosure() }},
	"tricky/MultipleDeferSameName":  {trickySrc, "MultipleDeferSameName", func() any { return tricky.MultipleDeferSameName() }},
	"tricky/ClosureMutateOuter":     {trickySrc, "ClosureMutateOuter", func() any { return tricky.ClosureMutateOuter() }},
	"tricky/SliceAppendExpand":      {trickySrc, "SliceAppendExpand", func() any { return tricky.SliceAppendExpand() }},
	"tricky/MapIncrement":           {trickySrc, "MapIncrement", func() any { return tricky.MapIncrement() }},
	"tricky/InterfaceTypeSwitch":    {trickySrc, "InterfaceTypeSwitch", func() any { return tricky.InterfaceTypeSwitch() }},
	"tricky/PointerToSlice":         {trickySrc, "PointerToSlice", func() any { return tricky.PointerToSlice() }},
	"tricky/NestedClosure":          {trickySrc, "NestedClosure", func() any { return tricky.NestedClosure() }},
	"tricky/SliceOfSlice":           {trickySrc, "SliceOfSlice", func() any { return tricky.SliceOfSlice() }},
	"tricky/MapOfSlice":             {trickySrc, "MapOfSlice", func() any { return tricky.MapOfSlice() }},
	"tricky/StructWithSlice":        {trickySrc, "StructWithSlice", func() any { return tricky.StructWithSlice() }},
	"tricky/DeferReadAfterAssign":   {trickySrc, "DeferReadAfterAssign", func() any { return tricky.DeferReadAfterAssign() }},
	"tricky/ForRangePointer":        {trickySrc, "ForRangePointer", func() any { return tricky.ForRangePointer() }},
	"tricky/NilInterfaceValue":      {trickySrc, "NilInterfaceValue", func() any { return tricky.NilInterfaceValue() }},
	"tricky/SliceCopyOverlap":       {trickySrc, "SliceCopyOverlap", func() any { return tricky.SliceCopyOverlap() }},
	"tricky/PointerReassign":        {trickySrc, "PointerReassign", func() any { return tricky.PointerReassign() }},
	"tricky/InterfaceNilComparison": {trickySrc, "InterfaceNilComparison", func() any { return tricky.InterfaceNilComparison() }},
	"tricky/DeferClosureCapture":    {trickySrc, "DeferClosureCapture", func() any { return tricky.DeferClosureCapture() }},
	"tricky/MapLookupModify":        {trickySrc, "MapLookupModify", func() any { return tricky.MapLookupModify() }},
	"tricky/SliceZeroLength":        {trickySrc, "SliceZeroLength", func() any { return tricky.SliceZeroLength() }},

	// ============================================================================
	// tricky - Iteration 3
	// ============================================================================
	"tricky/SliceModifyViaSubslice":     {trickySrc, "SliceModifyViaSubslice", func() any { return tricky.SliceModifyViaSubslice() }},
	"tricky/MapDeleteDuringRange":       {trickySrc, "MapDeleteDuringRange", func() any { return tricky.MapDeleteDuringRange() }},
	"tricky/ClosureReturnClosure":       {trickySrc, "ClosureReturnClosure", func() any { return tricky.ClosureReturnClosure() }},
	"tricky/StructMethodOnNil":          {trickySrc, "StructMethodOnNil", func() any { return tricky.StructMethodOnNil() }},
	"tricky/ArrayPointerIndex":          {trickySrc, "ArrayPointerIndex", func() any { return tricky.ArrayPointerIndex() }},
	"tricky/SliceThreeIndex":            {trickySrc, "SliceThreeIndex", func() any { return tricky.SliceThreeIndex() }},
	"tricky/DeferInLoop":                {trickySrc, "DeferInLoop", func() any { return tricky.DeferInLoop() }},
	"tricky/StructCompare":              {trickySrc, "StructCompare", func() any { return tricky.StructCompare() }},
	"tricky/InterfaceSlice":             {trickySrc, "InterfaceSlice", func() any { return tricky.InterfaceSlice() }},
	"tricky/PointerMethodValueReceiver": {trickySrc, "PointerMethodValueReceiver", func() any { return tricky.PointerMethodValueReceiver() }},
	"tricky/SliceOfMaps":                {trickySrc, "SliceOfMaps", func() any { return tricky.SliceOfMaps() }},
	"tricky/MapWithNilValue":            {trickySrc, "MapWithNilValue", func() any { return tricky.MapWithNilValue() }},
	"tricky/SwitchNoCondition":          {trickySrc, "SwitchNoCondition", func() any { return tricky.SwitchNoCondition() }},
	"tricky/DeferModifyReturn":          {trickySrc, "DeferModifyReturn", func() any { return tricky.DeferModifyReturn() }},
	"tricky/SliceAppendToCap":           {trickySrc, "SliceAppendToCap", func() any { return tricky.SliceAppendToCap() }},
	"tricky/ForRangeStringByteIndex":    {trickySrc, "ForRangeStringByteIndex", func() any { return tricky.ForRangeStringByteIndex() }},
	"tricky/StructLiteralEmbedded":      {trickySrc, "StructLiteralEmbedded", func() any { return tricky.StructLiteralEmbedded() }},
	"tricky/MapNilKey":                  {trickySrc, "MapNilKey", func() any { return tricky.MapNilKey() }},
	"tricky/ClosureRecursive":           {trickySrc, "ClosureRecursive", func() any { return tricky.ClosureRecursive() }},

	// ============================================================================
	// tricky - Iteration 4
	// ============================================================================
	"tricky/DeferCallInDefer":       {trickySrc, "DeferCallInDefer", func() any { return tricky.DeferCallInDefer() }},
	"tricky/MapLookupAssign":        {trickySrc, "MapLookupAssign", func() any { return tricky.MapLookupAssign() }},
	"tricky/StructMethodOnValue":    {trickySrc, "StructMethodOnValue", func() any { return tricky.StructMethodOnValue() }},
	"tricky/PointerToMap":           {trickySrc, "PointerToMap", func() any { return tricky.PointerToMap() }},
	"tricky/SliceCapAfterAppend":    {trickySrc, "SliceCapAfterAppend", func() any { return tricky.SliceCapAfterAppend() }},
	"tricky/NestedMaps":             {trickySrc, "NestedMaps", func() any { return tricky.NestedMaps() }},
	"tricky/StructPointerNil":       {trickySrc, "StructPointerNil", func() any { return tricky.StructPointerNil() }},
	"tricky/VariadicWithSlice":      {trickySrc, "VariadicWithSlice", func() any { return tricky.VariadicWithSlice() }},
	"tricky/SliceMakeWithLen":       {trickySrc, "SliceMakeWithLen", func() any { return tricky.SliceMakeWithLen() }},
	"tricky/InterfaceConversion":    {trickySrc, "InterfaceConversion", func() any { return tricky.InterfaceConversion() }},
	"tricky/MapWithEmptyStringKey":  {trickySrc, "MapWithEmptyStringKey", func() any { return tricky.MapWithEmptyStringKey() }},
	"tricky/DeferPanicRecover":      {trickySrc, "DeferPanicRecover", func() any { return tricky.DeferPanicRecover() }},
	"tricky/StructWithMap":          {trickySrc, "StructWithMap", func() any { return tricky.StructWithMap() }},
	"tricky/ForRangeBreak":          {trickySrc, "ForRangeBreak", func() any { return tricky.ForRangeBreak() }},
	"tricky/SliceLiteralNested":     {trickySrc, "SliceLiteralNested", func() any { return tricky.SliceLiteralNested() }},
	"tricky/MapLiteralNested":       {trickySrc, "MapLiteralNested", func() any { return tricky.MapLiteralNested() }},
	"tricky/PointerToStructLiteral": {trickySrc, "PointerToStructLiteral", func() any { return tricky.PointerToStructLiteral() }},
	"tricky/SliceOfStructs":         {trickySrc, "SliceOfStructs", func() any { return tricky.SliceOfStructs() }},
	"tricky/MapIterateModify":       {trickySrc, "MapIterateModify", func() any { return tricky.MapIterateModify() }},

	// ============================================================================
	// tricky - Iteration 5
	// ============================================================================
	"tricky/ChannelBuffered":        {trickySrc, "ChannelBuffered", func() any { return tricky.ChannelBuffered() }},
	"tricky/StructEmbeddedMethod":   {trickySrc, "StructEmbeddedMethod", func() any { return tricky.StructEmbeddedMethod() }},
	"tricky/SliceOfChannels":        {trickySrc, "SliceOfChannels", func() any { return tricky.SliceOfChannels() }},
	"tricky/MapOfChannels":          {trickySrc, "MapOfChannels", func() any { return tricky.MapOfChannels() }},
	"tricky/MultipleAssignment":     {trickySrc, "MultipleAssignment", func() any { return tricky.MultipleAssignment() }},
	"tricky/SliceAssign":            {trickySrc, "SliceAssign", func() any { return tricky.SliceAssign() }},
	"tricky/MapTwoAssign":           {trickySrc, "MapTwoAssign", func() any { return tricky.MapTwoAssign() }},
	"tricky/StructPointerMethodNil": {trickySrc, "StructPointerMethodNil", func() any { return tricky.StructPointerMethodNil() }},
	"tricky/DeferAfterPanic":        {trickySrc, "DeferAfterPanic", func() any { return tricky.DeferAfterPanic() }},
	"tricky/SliceFromArray":         {trickySrc, "SliceFromArray", func() any { return tricky.SliceFromArray() }},
	"tricky/ArrayPointerSlice":      {trickySrc, "ArrayPointerSlice", func() any { return tricky.ArrayPointerSlice() }},
	"tricky/StructFieldPointer":     {trickySrc, "StructFieldPointer", func() any { return tricky.StructFieldPointer() }},
	"tricky/MapLenCap":              {trickySrc, "MapLenCap", func() any { return tricky.MapLenCap() }},
	"tricky/StringConcat":           {trickySrc, "StringConcat", func() any { return tricky.StringConcat() }},
	"tricky/StringLen":              {trickySrc, "StringLen", func() any { return tricky.StringLen() }},

	// ============================================================================
	// tricky - Iteration 6
	// ============================================================================
	"tricky/ComplexMapKey":           {trickySrc, "ComplexMapKey", func() any { return tricky.ComplexMapKey() }},
	"tricky/SliceReverse":            {trickySrc, "SliceReverse", func() any { return tricky.SliceReverse() }},
	"tricky/MapMerge":                {trickySrc, "MapMerge", func() any { return tricky.MapMerge() }},
	"tricky/StructZeroValue":         {trickySrc, "StructZeroValue", func() any { return tricky.StructZeroValue() }},
	"tricky/SliceDeleteByIndex":      {trickySrc, "SliceDeleteByIndex", func() any { return tricky.SliceDeleteByIndex() }},
	"tricky/MapValueOverwrite":       {trickySrc, "MapValueOverwrite", func() any { return tricky.MapValueOverwrite() }},
	"tricky/InterfaceEmbed":          {trickySrc, "InterfaceEmbed", func() any { return tricky.InterfaceEmbed() }},
	"tricky/SliceOfFuncs":            {trickySrc, "SliceOfFuncs", func() any { return tricky.SliceOfFuncs() }},
	"tricky/PointerToSliceElement":   {trickySrc, "PointerToSliceElement", func() any { return tricky.PointerToSliceElement() }},
	"tricky/MapKeyPointer":           {trickySrc, "MapKeyPointer", func() any { return tricky.MapKeyPointer() }},
	"tricky/SliceOfPointersToStruct": {trickySrc, "SliceOfPointersToStruct", func() any { return tricky.SliceOfPointersToStruct() }},
	"tricky/DoubleMapLookup":         {trickySrc, "DoubleMapLookup", func() any { return tricky.DoubleMapLookup() }},
	"tricky/StructSliceLiteral":      {trickySrc, "StructSliceLiteral", func() any { return tricky.StructSliceLiteral() }},
	"tricky/ForRangeModifyValue":     {trickySrc, "ForRangeModifyValue", func() any { return tricky.ForRangeModifyValue() }},
	"tricky/MapWithStructPointerKey": {trickySrc, "MapWithStructPointerKey", func() any { return tricky.MapWithStructPointerKey() }},
	"tricky/SliceCopyDifferentTypes": {trickySrc, "SliceCopyDifferentTypes", func() any { return tricky.SliceCopyDifferentTypes() }},
	"tricky/NestedStructWithPointer": {trickySrc, "NestedStructWithPointer", func() any { return tricky.NestedStructWithPointer() }},
	"tricky/SliceOfSlicesAppend":     {trickySrc, "SliceOfSlicesAppend", func() any { return tricky.SliceOfSlicesAppend() }},
	"tricky/MapDeleteAll":            {trickySrc, "MapDeleteAll", func() any { return tricky.MapDeleteAll() }},

	// ============================================================================
	// tricky - Iteration 7
	// ============================================================================
	"tricky/StructPointerSlice":        {trickySrc, "StructPointerSlice", func() any { return tricky.StructPointerSlice() }},
	"tricky/MapWithInterfaceKey":       {trickySrc, "MapWithInterfaceKey", func() any { return tricky.MapWithInterfaceKey() }},
	"tricky/SliceOfInterfaces":         {trickySrc, "SliceOfInterfaces", func() any { return tricky.SliceOfInterfaces() }},
	"tricky/NestedPointerStruct":       {trickySrc, "NestedPointerStruct", func() any { return tricky.NestedPointerStruct() }},
	"tricky/StructMethodOnNilPointer":  {trickySrc, "StructMethodOnNilPointer", func() any { return tricky.StructMethodOnNilPointer() }},
	"tricky/SliceAppendToSlice":        {trickySrc, "SliceAppendToSlice", func() any { return tricky.SliceAppendToSlice() }},
	"tricky/MapLookupWithDefault":      {trickySrc, "MapLookupWithDefault", func() any { return tricky.MapLookupWithDefault() }},
	"tricky/StructFieldUpdate":         {trickySrc, "StructFieldUpdate", func() any { return tricky.StructFieldUpdate() }},
	"tricky/PointerToNilSlice":         {trickySrc, "PointerToNilSlice", func() any { return tricky.PointerToNilSlice() }},
	"tricky/SliceCopyToSubslice":       {trickySrc, "SliceCopyToSubslice", func() any { return tricky.SliceCopyToSubslice() }},
	"tricky/StructWithMultipleFields":  {trickySrc, "StructWithMultipleFields", func() any { return tricky.StructWithMultipleFields() }},
	"tricky/ForRangeContinue":          {trickySrc, "ForRangeContinue", func() any { return tricky.ForRangeContinue() }},
	"tricky/MapWithBoolKey":            {trickySrc, "MapWithBoolKey", func() any { return tricky.MapWithBoolKey() }},
	"tricky/SliceInsert":               {trickySrc, "SliceInsert", func() any { return tricky.SliceInsert() }},
	"tricky/StructEmbeddedFieldAccess": {trickySrc, "StructEmbeddedFieldAccess", func() any { return tricky.StructEmbeddedFieldAccess() }},
	"tricky/PointerToChannel":          {trickySrc, "PointerToChannel", func() any { return tricky.PointerToChannel() }},
	"tricky/MapKeyModification":        {trickySrc, "MapKeyModification", func() any { return tricky.MapKeyModification() }},
	"tricky/SliceRangeModify":          {trickySrc, "SliceRangeModify", func() any { return tricky.SliceRangeModify() }},
	"tricky/StructLiteralShort":        {trickySrc, "StructLiteralShort", func() any { return tricky.StructLiteralShort() }},

	// ============================================================================
	// tricky - Iteration 8
	// ============================================================================
	"tricky/SliceDrain":        {trickySrc, "SliceDrain", func() any { return tricky.SliceDrain() }},
	"tricky/MapClear":          {trickySrc, "MapClear", func() any { return tricky.MapClear() }},
	"tricky/StructCopy":        {trickySrc, "StructCopy", func() any { return tricky.StructCopy() }},
	"tricky/PointerStructCopy": {trickySrc, "PointerStructCopy", func() any { return tricky.PointerStructCopy() }},
	"tricky/SliceFilter":       {trickySrc, "SliceFilter", func() any { return tricky.SliceFilter() }},
	"tricky/MapTransform":      {trickySrc, "MapTransform", func() any { return tricky.MapTransform() }},
	"tricky/SliceContains":     {trickySrc, "SliceContains", func() any { return tricky.SliceContains() }},
	"tricky/MapKeys":           {trickySrc, "MapKeys", func() any { return tricky.MapKeys() }},

	// ============================================================================
	// typeconv
	// ============================================================================
	"typeconv/IntToFloat64":           {typeconvSrc, "IntToFloat64", func() any { return typeconv.IntToFloat64() }},
	"typeconv/Float64Arithmetic":      {typeconvSrc, "Float64Arithmetic", func() any { return typeconv.Float64Arithmetic() }},
	"typeconv/StringToByteConversion": {typeconvSrc, "StringToByteConversion", func() any { return typeconv.StringToByteConversion() }},
	"typeconv/IntStringConversion":    {typeconvSrc, "IntStringConversion", func() any { return typeconv.IntStringConversion() }},
	"typeconv/StringIntConversion":    {typeconvSrc, "StringIntConversion", func() any { return typeconv.StringIntConversion() }},

	// ============================================================================
	// variables
	// ============================================================================
	"variables/DeclareAndUse":   {variablesSrc, "DeclareAndUse", func() any { return variables.DeclareAndUse() }},
	"variables/Reassignment":    {variablesSrc, "Reassignment", func() any { return variables.Reassignment() }},
	"variables/MultipleDecl":    {variablesSrc, "MultipleDecl", func() any { return variables.MultipleDecl() }},
	"variables/ZeroValues":      {variablesSrc, "ZeroValues", func() any { return variables.ZeroValues() }},
	"variables/StringZeroValue": {variablesSrc, "StringZeroValue", func() any { return variables.StringZeroValue() }},
	"variables/Shadowing":       {variablesSrc, "Shadowing", func() any { return variables.Shadowing() }},

	// ============================================================================
	// init
	// ============================================================================
	"init/GetA": {initSrc, "GetA", func() any { return initialize.GetA() }},
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
			result, err := prog.Run(tc.funcName)
			interpDuration := time.Since(startInterp)
			if err != nil {
				t.Fatalf("Run error: %v", err)
			}

			startNative := time.Now()
			expected := tc.native()
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
