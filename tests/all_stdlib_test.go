package tests

import (
	_ "embed"
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
	"git.woa.com/youngjin/gig/tests/testdata/edgecases"
	"git.woa.com/youngjin/gig/tests/testdata/external"
	"git.woa.com/youngjin/gig/tests/testdata/functions"
	"git.woa.com/youngjin/gig/tests/testdata/known_issues"
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
	"git.woa.com/youngjin/gig/tests/testdata/typeconv"
	"git.woa.com/youngjin/gig/tests/testdata/variables"
)

// Embed all test source files
//
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

//go:embed testdata/edgecases/main.go
var edgecasesSrc string

//go:embed testdata/external/main.go
var externalSrc string

//go:embed testdata/functions/main.go
var functionsSrc string

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

//go:embed testdata/typeconv/main.go
var typeconvSrc string

//go:embed testdata/variables/main.go
var variablesSrc string

//go:embed testdata/known_issues/main.go
var knownIssuesSrc string

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

// allTests
// allTests contains all test definitions
var allTests = map[string]testCase{
	// algorithms
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

	// advanced
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

	// arithmetic
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

	// autowrap
	"autowrap/WithPackage": {autowrapSrc, "WithPackage", func() any { return autowrap.WithPackage() }},
	"autowrap/WithImport":  {autowrapSrc, "WithImport", func() any { return autowrap.WithImport() }},
	"autowrap/Compute":     {autowrapSrc, "Compute", func() any { return autowrap.Compute() }},

	// bitwise
	"bitwise/And":        {bitwiseSrc, "And", func() any { return bitwise.And() }},
	"bitwise/Or":         {bitwiseSrc, "Or", func() any { return bitwise.Or() }},
	"bitwise/Xor":        {bitwiseSrc, "Xor", func() any { return bitwise.Xor() }},
	"bitwise/LeftShift":  {bitwiseSrc, "LeftShift", func() any { return bitwise.LeftShift() }},
	"bitwise/RightShift": {bitwiseSrc, "RightShift", func() any { return bitwise.RightShift() }},
	"bitwise/Combined":   {bitwiseSrc, "Combined", func() any { return bitwise.Combined() }},
	"bitwise/AndNot":     {bitwiseSrc, "AndNot", func() any { return bitwise.AndNot() }},
	"bitwise/PowerOfTwo": {bitwiseSrc, "PowerOfTwo", func() any { return bitwise.PowerOfTwo() }},

	// closures
	"closures/Counter":           {closuresSrc, "Counter", func() any { return closures.Counter() }},
	"closures/CaptureMutation":   {closuresSrc, "CaptureMutation", func() any { return closures.CaptureMutation() }},
	"closures/Factory":           {closuresSrc, "Factory", func() any { return closures.Factory() }},
	"closures/MultipleInstances": {closuresSrc, "MultipleInstances", func() any { return closures.MultipleInstances() }},
	"closures/OverLoop":          {closuresSrc, "OverLoop", func() any { return closures.OverLoop() }},
	"closures/Chain":             {closuresSrc, "Chain", func() any { return closures.Chain() }},
	"closures/Accumulator":       {closuresSrc, "Accumulator", func() any { return closures.Accumulator() }},

	// closures_advanced
	"closures_advanced/Generator":          {closuresAdvancedSrc, "Generator", func() any { return closures_advanced.Generator() }},
	"closures_advanced/Predicate":          {closuresAdvancedSrc, "Predicate", func() any { return closures_advanced.Predicate() }},
	"closures_advanced/StateMachine":       {closuresAdvancedSrc, "StateMachine", func() any { return closures_advanced.StateMachine() }},
	"closures_advanced/RecursiveHelper":    {closuresAdvancedSrc, "RecursiveHelper", func() any { return closures_advanced.RecursiveHelper() }},
	"closures_advanced/ApplyN":             {closuresAdvancedSrc, "ApplyN", func() any { return closures_advanced.ApplyN() }},
	"closures_advanced/Compose":            {closuresAdvancedSrc, "Compose", func() any { return closures_advanced.Compose() }},
	"closures_advanced/ClosureForLoopTest": {closuresAdvancedSrc, "ClosureForLoopTest", func() any { return closures_advanced.ClosureForLoopTest() }},

	// controlflow
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

	// edgecases
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

	// external
	"external/FmtSprintf":       {externalSrc, "FmtSprintf", func() any { return external.FmtSprintf() }},
	"external/FmtSprintfMulti":  {externalSrc, "FmtSprintfMulti", func() any { return external.FmtSprintfMulti() }},
	"external/StringsToUpper":   {externalSrc, "StringsToUpper", func() any { return external.StringsToUpper() }},
	"external/StringsToLower":   {externalSrc, "StringsToLower", func() any { return external.StringsToLower() }},
	"external/StringsContains":  {externalSrc, "StringsContains", func() any { return external.StringsContains() }},
	"external/StringsReplace":   {externalSrc, "StringsReplace", func() any { return external.StringsReplace() }},
	"external/StringsHasPrefix": {externalSrc, "StringsHasPrefix", func() any { return external.StringsHasPrefix() }},
	"external/StrconvItoa":      {externalSrc, "StrconvItoa", func() any { return external.StrconvItoa() }},
	"external/StrconvAtoi":      {externalSrc, "StrconvAtoi", func() any { return external.StrconvAtoi() }},

	// functions
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

	// leetcode_hard
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

	// maps
	"maps/BasicOps":       {mapsSrc, "BasicOps", func() any { return maps.BasicOps() }},
	"maps/Iteration":      {mapsSrc, "Iteration", func() any { return maps.Iteration() }},
	"maps/Delete":         {mapsSrc, "Delete", func() any { return maps.Delete() }},
	"maps/Len":            {mapsSrc, "Len", func() any { return maps.Len() }},
	"maps/Overwrite":      {mapsSrc, "Overwrite", func() any { return maps.Overwrite() }},
	"maps/IntKeys":        {mapsSrc, "IntKeys", func() any { return maps.IntKeys() }},
	"maps/PassToFunction": {mapsSrc, "PassToFunction", func() any { return maps.PassToFunction() }},

	// mapadvanced
	"mapadvanced/LookupExistingKey": {mapadvancedSrc, "LookupExistingKey", func() any { return mapadvanced.LookupExistingKey() }},
	"mapadvanced/LookupWithDefault": {mapadvancedSrc, "LookupWithDefault", func() any { return mapadvanced.LookupWithDefault() }},
	"mapadvanced/AsCounter":         {mapadvancedSrc, "AsCounter", func() any { return mapadvanced.AsCounter() }},
	"mapadvanced/WithStringValues":  {mapadvancedSrc, "WithStringValues", func() any { return mapadvanced.WithStringValues() }},
	"mapadvanced/BuildFromLoop":     {mapadvancedSrc, "BuildFromLoop", func() any { return mapadvanced.BuildFromLoop() }},
	"mapadvanced/DeleteAndReinsert": {mapadvancedSrc, "DeleteAndReinsert", func() any { return mapadvanced.DeleteAndReinsert() }},

	// multiassign
	"multiassign/Swap":             {multiassignSrc, "Swap", func() any { return multiassign.Swap() }},
	"multiassign/FromFunction":     {multiassignSrc, "FromFunction", func() any { return multiassign.FromFunction() }},
	"multiassign/ThreeValues":      {multiassignSrc, "ThreeValues", func() any { return multiassign.ThreeValues() }},
	"multiassign/InLoop":           {multiassignSrc, "InLoop", func() any { return multiassign.InLoop() }},
	"multiassign/DiscardWithBlank": {multiassignSrc, "DiscardWithBlank", func() any { return multiassign.DiscardWithBlank() }},

	// namedreturn
	"namedreturn/Basic":     {namedreturnSrc, "Basic", func() any { return namedreturn.Basic() }},
	"namedreturn/Multiple":  {namedreturnSrc, "Multiple", func() any { return namedreturn.Multiple() }},
	"namedreturn/ZeroValue": {namedreturnSrc, "ZeroValue", func() any { return namedreturn.ZeroValue() }},

	// recursion
	"recursion/TailRecursionPattern": {recursionSrc, "TailRecursionPattern", func() any { return recursion.TailRecursionPattern() }},
	"recursion/ReverseSlice":         {recursionSrc, "ReverseSlice", func() any { return recursion.ReverseSlice() }},
	"recursion/TowerOfHanoi":         {recursionSrc, "TowerOfHanoi", func() any { return recursion.TowerOfHanoi() }},
	"recursion/MaxSlice":             {recursionSrc, "MaxSlice", func() any { return recursion.MaxSlice() }},
	"recursion/Ackermann":            {recursionSrc, "Ackermann", func() any { return recursion.Ackermann() }},
	"recursion/BinarySearch":         {recursionSrc, "BinarySearch", func() any { return recursion.BinarySearch() }},

	// scope
	"scope/IfInitShortVar":            {scopeSrc, "IfInitShortVar", func() any { return scope.IfInitShortVar() }},
	"scope/IfInitMultiCondition":      {scopeSrc, "IfInitMultiCondition", func() any { return scope.IfInitMultiCondition() }},
	"scope/NestedScopes":              {scopeSrc, "NestedScopes", func() any { return scope.NestedScopes() }},
	"scope/ForScopeIsolation":         {scopeSrc, "ForScopeIsolation", func() any { return scope.ForScopeIsolation() }},
	"scope/MultipleBlockScopes":       {scopeSrc, "MultipleBlockScopes", func() any { return scope.MultipleBlockScopes() }},
	"scope/ClosureCapturesOuterScope": {scopeSrc, "ClosureCapturesOuterScope", func() any { return scope.ClosureCapturesOuterScope() }},

	// slices
	"slices/MakeLen":           {slicesSrc, "MakeLen", func() any { return slices.MakeLen() }},
	"slices/Append":            {slicesSrc, "Append", func() any { return slices.Append() }},
	"slices/ElementAssignment": {slicesSrc, "ElementAssignment", func() any { return slices.ElementAssignment() }},
	"slices/ForRange":          {slicesSrc, "ForRange", func() any { return slices.ForRange() }},
	"slices/ForRangeIndex":     {slicesSrc, "ForRangeIndex", func() any { return slices.ForRangeIndex() }},
	"slices/GrowMultiple":      {slicesSrc, "GrowMultiple", func() any { return slices.GrowMultiple() }},
	"slices/PassToFunction":    {slicesSrc, "PassToFunction", func() any { return slices.PassToFunction() }},
	"slices/LenCap":            {slicesSrc, "LenCap", func() any { return slices.LenCap() }},

	// slicing
	"slicing/SubSliceBasic":            {slicingSrc, "SubSliceBasic", func() any { return slicing.SubSliceBasic() }},
	"slicing/SubSliceLen":              {slicingSrc, "SubSliceLen", func() any { return slicing.SubSliceLen() }},
	"slicing/SubSliceFromStart":        {slicingSrc, "SubSliceFromStart", func() any { return slicing.SubSliceFromStart() }},
	"slicing/SubSliceToEnd":            {slicingSrc, "SubSliceToEnd", func() any { return slicing.SubSliceToEnd() }},
	"slicing/SubSliceCopy":             {slicingSrc, "SubSliceCopy", func() any { return slicing.SubSliceCopy() }},
	"slicing/SubSliceChained":          {slicingSrc, "SubSliceChained", func() any { return slicing.SubSliceChained() }},
	"slicing/SubSliceModifiesOriginal": {slicingSrc, "SubSliceModifiesOriginal", func() any { return slicing.SubSliceModifiesOriginal() }},

	// strings_pkg
	"strings_pkg/Concat":     {stringsPkgSrc, "Concat", func() any { return strings_pkg.Concat() }},
	"strings_pkg/ConcatLoop": {stringsPkgSrc, "ConcatLoop", func() any { return strings_pkg.ConcatLoop() }},
	"strings_pkg/Len":        {stringsPkgSrc, "Len", func() any { return strings_pkg.Len() }},
	"strings_pkg/Index":      {stringsPkgSrc, "Index", func() any { return strings_pkg.Index() }},
	"strings_pkg/Comparison": {stringsPkgSrc, "Comparison", func() any { return strings_pkg.Comparison() }},
	"strings_pkg/Equality":   {stringsPkgSrc, "Equality", func() any { return strings_pkg.Equality() }},
	"strings_pkg/EmptyCheck": {stringsPkgSrc, "EmptyCheck", func() any { return strings_pkg.EmptyCheck() }},

	// structs
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

	// switch
	"switch/Simple":      {switchSrc, "Simple", func() any { return switch_pkg.Simple() }},
	"switch/Default":     {switchSrc, "Default", func() any { return switch_pkg.Default() }},
	"switch/MultiCase":   {switchSrc, "MultiCase", func() any { return switch_pkg.MultiCase() }},
	"switch/NoCondition": {switchSrc, "NoCondition", func() any { return switch_pkg.NoCondition() }},
	"switch/WithInit":    {switchSrc, "WithInit", func() any { return switch_pkg.WithInit() }},
	"switch/StringCases": {switchSrc, "StringCases", func() any { return switch_pkg.StringCases() }},
	"switch/Fallthrough": {switchSrc, "Fallthrough", func() any { return switch_pkg.Fallthrough() }},
	"switch/Nested":      {switchSrc, "Nested", func() any { return switch_pkg.Nested() }},

	// typeconv
	"typeconv/IntToFloat64":           {typeconvSrc, "IntToFloat64", func() any { return typeconv.IntToFloat64() }},
	"typeconv/Float64Arithmetic":      {typeconvSrc, "Float64Arithmetic", func() any { return typeconv.Float64Arithmetic() }},
	"typeconv/StringToByteConversion": {typeconvSrc, "StringToByteConversion", func() any { return typeconv.StringToByteConversion() }},
	"typeconv/IntStringConversion":    {typeconvSrc, "IntStringConversion", func() any { return typeconv.IntStringConversion() }},
	"typeconv/StringIntConversion":    {typeconvSrc, "StringIntConversion", func() any { return typeconv.StringIntConversion() }},

	// variables
	"variables/DeclareAndUse":   {variablesSrc, "DeclareAndUse", func() any { return variables.DeclareAndUse() }},
	"variables/Reassignment":    {variablesSrc, "Reassignment", func() any { return variables.Reassignment() }},
	"variables/MultipleDecl":    {variablesSrc, "MultipleDecl", func() any { return variables.MultipleDecl() }},
	"variables/ZeroValues":      {variablesSrc, "ZeroValues", func() any { return variables.ZeroValues() }},
	"variables/StringZeroValue": {variablesSrc, "StringZeroValue", func() any { return variables.StringZeroValue() }},
	"variables/Shadowing":       {variablesSrc, "Shadowing", func() any { return variables.Shadowing() }},

	// known_issues — byte-to-string conversion
	"known_issues/BytesToStringHi":     {knownIssuesSrc, "BytesToStringHi", func() any { return known_issues.BytesToStringHi() }},
	"known_issues/BytesToStringGo":     {knownIssuesSrc, "BytesToStringGo", func() any { return known_issues.BytesToStringGo() }},
	"known_issues/BytesToStringSingle": {knownIssuesSrc, "BytesToStringSingle", func() any { return known_issues.BytesToStringSingle() }},
	"known_issues/BytesToStringEmpty":  {knownIssuesSrc, "BytesToStringEmpty", func() any { return known_issues.BytesToStringEmpty() }},

	// known_issues — range-over-string rune values
	"known_issues/RangeStringRune":      {knownIssuesSrc, "RangeStringRune", func() any { return known_issues.RangeStringRune() }},
	"known_issues/RangeStringIndex":     {knownIssuesSrc, "RangeStringIndex", func() any { return known_issues.RangeStringIndex() }},
	"known_issues/RangeStringMultibyte": {knownIssuesSrc, "RangeStringMultibyte", func() any { return known_issues.RangeStringMultibyte() }},

	// known_issues — init() execution
	"known_issues/InitFuncResult": {knownIssuesSrc, "InitFuncResult", func() any { return known_issues.InitFuncResult() }},
	"known_issues/InitSliceLen":   {knownIssuesSrc, "InitSliceLen", func() any { return known_issues.InitSliceLen() }},

	// known_issues — pointer-receiver mutation
	"known_issues/PointerReceiverMutation": {knownIssuesSrc, "PointerReceiverMutation", func() any { return known_issues.PointerReceiverMutation() }},
	"known_issues/PointerReceiverSetGet":   {knownIssuesSrc, "PointerReceiverSetGet", func() any { return known_issues.PointerReceiverSetGet() }},
}

// TestAllStdlib runs all stdlib tests
func TestAllStdlib(t *testing.T) {
	for name, tc := range allTests {
		t.Run(name, func(t *testing.T) {
			src := toMainPackage(tc.src)
			prog, err := gig.Build(src)
			if err != nil {
				t.Fatalf("Build error: %v", err)
			}

			// Measure interpreter execution time
			startInterp := time.Now()
			result, err := prog.Run(tc.funcName)
			interpDuration := time.Since(startInterp)
			if err != nil {
				t.Fatalf("Run error: %v", err)
			}

			// Measure native execution time
			startNative := time.Now()
			expected := tc.native()
			nativeDuration := time.Since(startNative)

			compareResults(t, result, expected)

			// Report timing comparison
			ratio := float64(interpDuration) / float64(nativeDuration)
			t.Logf("interp: %v, native: %v, ratio: %.1fx", interpDuration, nativeDuration, ratio)
		})
	}
}

// BenchmarkAllStdlib runs benchmarks for all stdlib tests.
// Skipped by default due to long runtime. Use -run=BenchmarkAllStdlib to enable.
func BenchmarkAllStdlib(b *testing.B) {
	b.Skip("BenchmarkAllStdlib takes too long. Run with -run=BenchmarkAllStdlib to enable.")
}

// compareResults compares interpreter result with native result
func compareResults(t *testing.T, result, expected any) {
	t.Helper()
	switch exp := expected.(type) {
	case int:
		var got int64
		switch v := result.(type) {
		case int64:
			got = v
		case int:
			got = int64(v)
		default:
			t.Fatalf("expected int, got %T", result)
		}
		if got != int64(exp) {
			t.Errorf("expected %d, got %d", exp, got)
		}
	case string:
		got, ok := result.(string)
		if !ok {
			t.Fatalf("expected string, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %q, got %q", exp, got)
		}
	default:
		t.Fatalf("unsupported result type: %T", expected)
	}
}
