package tests

import (
	_ "embed"
	"strings"
	"testing"

	"gig"
	_ "gig/stdlib/packages"

	"gig/tests/testdata/advanced"
	"gig/tests/testdata/algorithms"
	"gig/tests/testdata/arithmetic"
	"gig/tests/testdata/autowrap"
	"gig/tests/testdata/bitwise"
	"gig/tests/testdata/closures"
	"gig/tests/testdata/closures_advanced"
	"gig/tests/testdata/controlflow"
	"gig/tests/testdata/edgecases"
	"gig/tests/testdata/external"
	"gig/tests/testdata/functions"
	"gig/tests/testdata/mapadvanced"
	"gig/tests/testdata/maps"
	"gig/tests/testdata/multiassign"
	"gig/tests/testdata/namedreturn"
	"gig/tests/testdata/recursion"
	"gig/tests/testdata/scope"
	"gig/tests/testdata/slices"
	"gig/tests/testdata/slicing"
	"gig/tests/testdata/strings_pkg"
	"gig/tests/testdata/switch"
	"gig/tests/testdata/typeconv"
	"gig/tests/testdata/variables"
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

//go:embed testdata/switch/main.go
var switchSrc string

//go:embed testdata/typeconv/main.go
var typeconvSrc string

//go:embed testdata/variables/main.go
var variablesSrc string

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
	native   func() interface{}
}

// testSuite groups tests by source file
type testSuite struct {
	src   string
	tests map[string]func() interface{}
}

// allTests contains all test definitions
var allTests = map[string]testCase{
	// algorithms
	"algorithms/InsertionSort":    {algorithmsSrc, "InsertionSort", func() interface{} { return algorithms.InsertionSort() }},
	"algorithms/SelectionSort":    {algorithmsSrc, "SelectionSort", func() interface{} { return algorithms.SelectionSort() }},
	"algorithms/ReverseSlice":     {algorithmsSrc, "ReverseSlice", func() interface{} { return algorithms.ReverseSlice() }},
	"algorithms/IsPalindrome":     {algorithmsSrc, "IsPalindrome", func() interface{} { return algorithms.IsPalindrome() }},
	"algorithms/PowerFunction":    {algorithmsSrc, "PowerFunction", func() interface{} { return algorithms.PowerFunction() }},
	"algorithms/MaxSubarraySum":   {algorithmsSrc, "MaxSubarraySum", func() interface{} { return algorithms.MaxSubarraySum() }},
	"algorithms/TwoSum":           {algorithmsSrc, "TwoSum", func() interface{} { return algorithms.TwoSum() }},
	"algorithms/FibMemoized":      {algorithmsSrc, "FibMemoized", func() interface{} { return algorithms.FibMemoized() }},
	"algorithms/CountDigits":      {algorithmsSrc, "CountDigits", func() interface{} { return algorithms.CountDigits() }},
	"algorithms/CollatzConjecture": {algorithmsSrc, "CollatzConjecture", func() interface{} { return algorithms.CollatzConjecture() }},

	// advanced
	"advanced/TypeConvertIntIdentity": {advancedSrc, "TypeConvertIntIdentity", func() interface{} { return advanced.TypeConvertIntIdentity() }},
	"advanced/DeepCallChain":          {advancedSrc, "DeepCallChain", func() interface{} { return advanced.DeepCallChain() }},
	"advanced/EarlyReturn":            {advancedSrc, "EarlyReturn", func() interface{} { return advanced.EarlyReturn() }},
	"advanced/NestedIfInLoop":         {advancedSrc, "NestedIfInLoop", func() interface{} { return advanced.NestedIfInLoop() }},
	"advanced/BubbleSort":             {advancedSrc, "BubbleSort", func() interface{} { return advanced.BubbleSort() }},
	"advanced/BinarySearch":           {advancedSrc, "BinarySearch", func() interface{} { return advanced.BinarySearch() }},
	"advanced/GCD":                    {advancedSrc, "GCD", func() interface{} { return advanced.GCD() }},
	"advanced/SieveOfEratosthenes":    {advancedSrc, "SieveOfEratosthenes", func() interface{} { return advanced.SieveOfEratosthenes() }},
	"advanced/MatrixMultiply":         {advancedSrc, "MatrixMultiply", func() interface{} { return advanced.MatrixMultiply() }},
	"advanced/EmptyFunctionReturn":    {advancedSrc, "EmptyFunctionReturn", func() interface{} { return advanced.EmptyFunctionReturn() }},
	"advanced/SingleReturnValue":      {advancedSrc, "SingleReturnValue", func() interface{} { return advanced.SingleReturnValue() }},
	"advanced/ZeroIteration":          {advancedSrc, "ZeroIteration", func() interface{} { return advanced.ZeroIteration() }},
	"advanced/LargeLoop":              {advancedSrc, "LargeLoop", func() interface{} { return advanced.LargeLoop() }},
	"advanced/DeepRecursion":          {advancedSrc, "DeepRecursion", func() interface{} { return advanced.DeepRecursion() }},
	"advanced/MapWithClosure":         {advancedSrc, "MapWithClosure", func() interface{} { return advanced.MapWithClosure() }},
	"advanced/SliceWithMultiReturn":   {advancedSrc, "SliceWithMultiReturn", func() interface{} { return advanced.SliceWithMultiReturn() }},
	"advanced/RecursiveDataBuild":     {advancedSrc, "RecursiveDataBuild", func() interface{} { return advanced.RecursiveDataBuild() }},
	"advanced/FunctionChain":          {advancedSrc, "FunctionChain", func() interface{} { return advanced.FunctionChain() }},
	"advanced/ComplexExpressions":     {advancedSrc, "ComplexExpressions", func() interface{} { return advanced.ComplexExpressions() }},

	// arithmetic
	"arithmetic/Addition":       {arithmeticSrc, "Addition", func() interface{} { return arithmetic.Addition() }},
	"arithmetic/Subtraction":    {arithmeticSrc, "Subtraction", func() interface{} { return arithmetic.Subtraction() }},
	"arithmetic/Multiplication": {arithmeticSrc, "Multiplication", func() interface{} { return arithmetic.Multiplication() }},
	"arithmetic/Division":       {arithmeticSrc, "Division", func() interface{} { return arithmetic.Division() }},
	"arithmetic/Modulo":         {arithmeticSrc, "Modulo", func() interface{} { return arithmetic.Modulo() }},
	"arithmetic/ComplexExpr":    {arithmeticSrc, "ComplexExpr", func() interface{} { return arithmetic.ComplexExpr() }},
	"arithmetic/Negation":       {arithmeticSrc, "Negation", func() interface{} { return arithmetic.Negation() }},
	"arithmetic/ChainedOps":     {arithmeticSrc, "ChainedOps", func() interface{} { return arithmetic.ChainedOps() }},
	"arithmetic/Overflow":       {arithmeticSrc, "Overflow", func() interface{} { return arithmetic.Overflow() }},
	"arithmetic/Precedence":     {arithmeticSrc, "Precedence", func() interface{} { return arithmetic.Precedence() }},

	// autowrap
	"autowrap/WithPackage": {autowrapSrc, "WithPackage", func() interface{} { return autowrap.WithPackage() }},
	"autowrap/WithImport":  {autowrapSrc, "WithImport", func() interface{} { return autowrap.WithImport() }},
	"autowrap/Compute":     {autowrapSrc, "Compute", func() interface{} { return autowrap.Compute() }},

	// bitwise
	"bitwise/And":         {bitwiseSrc, "And", func() interface{} { return bitwise.And() }},
	"bitwise/Or":          {bitwiseSrc, "Or", func() interface{} { return bitwise.Or() }},
	"bitwise/Xor":         {bitwiseSrc, "Xor", func() interface{} { return bitwise.Xor() }},
	"bitwise/LeftShift":   {bitwiseSrc, "LeftShift", func() interface{} { return bitwise.LeftShift() }},
	"bitwise/RightShift":  {bitwiseSrc, "RightShift", func() interface{} { return bitwise.RightShift() }},
	"bitwise/Combined":    {bitwiseSrc, "Combined", func() interface{} { return bitwise.Combined() }},
	"bitwise/AndNot":      {bitwiseSrc, "AndNot", func() interface{} { return bitwise.AndNot() }},
	"bitwise/PowerOfTwo":  {bitwiseSrc, "PowerOfTwo", func() interface{} { return bitwise.PowerOfTwo() }},

	// closures
	"closures/Counter":           {closuresSrc, "Counter", func() interface{} { return closures.Counter() }},
	"closures/CaptureMutation":   {closuresSrc, "CaptureMutation", func() interface{} { return closures.CaptureMutation() }},
	"closures/Factory":           {closuresSrc, "Factory", func() interface{} { return closures.Factory() }},
	"closures/MultipleInstances": {closuresSrc, "MultipleInstances", func() interface{} { return closures.MultipleInstances() }},
	"closures/OverLoop":          {closuresSrc, "OverLoop", func() interface{} { return closures.OverLoop() }},
	"closures/Chain":             {closuresSrc, "Chain", func() interface{} { return closures.Chain() }},
	"closures/Accumulator":       {closuresSrc, "Accumulator", func() interface{} { return closures.Accumulator() }},

	// closures_advanced
	"closures_advanced/Generator":        {closuresAdvancedSrc, "Generator", func() interface{} { return closures_advanced.Generator() }},
	"closures_advanced/Predicate":        {closuresAdvancedSrc, "Predicate", func() interface{} { return closures_advanced.Predicate() }},
	"closures_advanced/StateMachine":     {closuresAdvancedSrc, "StateMachine", func() interface{} { return closures_advanced.StateMachine() }},
	"closures_advanced/RecursiveHelper":  {closuresAdvancedSrc, "RecursiveHelper", func() interface{} { return closures_advanced.RecursiveHelper() }},
	"closures_advanced/ApplyN":           {closuresAdvancedSrc, "ApplyN", func() interface{} { return closures_advanced.ApplyN() }},
	"closures_advanced/Compose":          {closuresAdvancedSrc, "Compose", func() interface{} { return closures_advanced.Compose() }},

	// controlflow
	"controlflow/IfTrue":              {controlflowSrc, "IfTrue", func() interface{} { return controlflow.IfTrue() }},
	"controlflow/IfFalse":             {controlflowSrc, "IfFalse", func() interface{} { return controlflow.IfFalse() }},
	"controlflow/IfElse":              {controlflowSrc, "IfElse", func() interface{} { return controlflow.IfElse() }},
	"controlflow/IfElseChainNegative": {controlflowSrc, "IfElseChainNegative", func() interface{} { return controlflow.IfElseChainNegative() }},
	"controlflow/IfElseChainZero":     {controlflowSrc, "IfElseChainZero", func() interface{} { return controlflow.IfElseChainZero() }},
	"controlflow/IfElseChainPositive": {controlflowSrc, "IfElseChainPositive", func() interface{} { return controlflow.IfElseChainPositive() }},
	"controlflow/ForLoop":             {controlflowSrc, "ForLoop", func() interface{} { return controlflow.ForLoop() }},
	"controlflow/ForConditionOnly":    {controlflowSrc, "ForConditionOnly", func() interface{} { return controlflow.ForConditionOnly() }},
	"controlflow/NestedFor":           {controlflowSrc, "NestedFor", func() interface{} { return controlflow.NestedFor() }},
	"controlflow/ForBreak":            {controlflowSrc, "ForBreak", func() interface{} { return controlflow.ForBreak() }},
	"controlflow/ForContinue":         {controlflowSrc, "ForContinue", func() interface{} { return controlflow.ForContinue() }},
	"controlflow/BooleanAndOr":        {controlflowSrc, "BooleanAndOr", func() interface{} { return controlflow.BooleanAndOr() }},

	// edgecases
	"edgecases/MaxInt64":           {edgecasesSrc, "MaxInt64", func() interface{} { return edgecases.MaxInt64() }},
	"edgecases/MinInt64":           {edgecasesSrc, "MinInt64", func() interface{} { return edgecases.MinInt64() }},
	"edgecases/DivisionByMinusOne": {edgecasesSrc, "DivisionByMinusOne", func() interface{} { return edgecases.DivisionByMinusOne() }},
	"edgecases/ModuloNegative":     {edgecasesSrc, "ModuloNegative", func() interface{} { return edgecases.ModuloNegative() }},
	"edgecases/EmptyString":        {edgecasesSrc, "EmptyString", func() interface{} { return edgecases.EmptyString() }},
	"edgecases/LargeSlice":         {edgecasesSrc, "LargeSlice", func() interface{} { return edgecases.LargeSlice() }},
	"edgecases/NestedMapLookup":    {edgecasesSrc, "NestedMapLookup", func() interface{} { return edgecases.NestedMapLookup() }},
	"edgecases/ZeroDivisionGuard":  {edgecasesSrc, "ZeroDivisionGuard", func() interface{} { return edgecases.ZeroDivisionGuard() }},
	"edgecases/BooleanComplexExpr": {edgecasesSrc, "BooleanComplexExpr", func() interface{} { return edgecases.BooleanComplexExpr() }},
	"edgecases/SingleElementSlice": {edgecasesSrc, "SingleElementSlice", func() interface{} { return edgecases.SingleElementSlice() }},
	"edgecases/EmptyMap":           {edgecasesSrc, "EmptyMap", func() interface{} { return edgecases.EmptyMap() }},
	"edgecases/TightLoop":          {edgecasesSrc, "TightLoop", func() interface{} { return edgecases.TightLoop() }},

	// external
	"external/FmtSprintf":      {externalSrc, "FmtSprintf", func() interface{} { return external.FmtSprintf() }},
	"external/FmtSprintfMulti": {externalSrc, "FmtSprintfMulti", func() interface{} { return external.FmtSprintfMulti() }},
	"external/StringsToUpper":  {externalSrc, "StringsToUpper", func() interface{} { return external.StringsToUpper() }},
	"external/StringsToLower":  {externalSrc, "StringsToLower", func() interface{} { return external.StringsToLower() }},
	"external/StringsContains": {externalSrc, "StringsContains", func() interface{} { return external.StringsContains() }},
	"external/StringsReplace":  {externalSrc, "StringsReplace", func() interface{} { return external.StringsReplace() }},
	"external/StringsHasPrefix": {externalSrc, "StringsHasPrefix", func() interface{} { return external.StringsHasPrefix() }},
	"external/StrconvItoa":     {externalSrc, "StrconvItoa", func() interface{} { return external.StrconvItoa() }},
	"external/StrconvAtoi":     {externalSrc, "StrconvAtoi", func() interface{} { return external.StrconvAtoi() }},

	// functions
	"functions/Call":                {functionsSrc, "Call", func() interface{} { return functions.Call() }},
	"functions/MultipleReturn":      {functionsSrc, "MultipleReturn", func() interface{} { return functions.MultipleReturn() }},
	"functions/MultipleReturnDivmod": {functionsSrc, "MultipleReturnDivmod", func() interface{} { return functions.MultipleReturnDivmod() }},
	"functions/RecursionFactorial":  {functionsSrc, "RecursionFactorial", func() interface{} { return functions.RecursionFactorial() }},
	"functions/MutualRecursion":     {functionsSrc, "MutualRecursion", func() interface{} { return functions.MutualRecursion() }},
	"functions/FibonacciIterative":  {functionsSrc, "FibonacciIterative", func() interface{} { return functions.FibonacciIterative() }},
	"functions/FibonacciRecursive":  {functionsSrc, "FibonacciRecursive", func() interface{} { return functions.FibonacciRecursive() }},
	"functions/VariadicFunction":    {functionsSrc, "VariadicFunction", func() interface{} { return functions.VariadicFunction() }},
	"functions/FunctionAsValue":     {functionsSrc, "FunctionAsValue", func() interface{} { return functions.FunctionAsValue() }},
	"functions/HigherOrderMap":      {functionsSrc, "HigherOrderMap", func() interface{} { return functions.HigherOrderMap() }},
	"functions/HigherOrderFilter":   {functionsSrc, "HigherOrderFilter", func() interface{} { return functions.HigherOrderFilter() }},
	"functions/HigherOrderReduce":   {functionsSrc, "HigherOrderReduce", func() interface{} { return functions.HigherOrderReduce() }},

	// maps
	"maps/BasicOps":      {mapsSrc, "BasicOps", func() interface{} { return maps.BasicOps() }},
	"maps/Iteration":     {mapsSrc, "Iteration", func() interface{} { return maps.Iteration() }},
	"maps/Delete":        {mapsSrc, "Delete", func() interface{} { return maps.Delete() }},
	"maps/Len":           {mapsSrc, "Len", func() interface{} { return maps.Len() }},
	"maps/Overwrite":     {mapsSrc, "Overwrite", func() interface{} { return maps.Overwrite() }},
	"maps/IntKeys":       {mapsSrc, "IntKeys", func() interface{} { return maps.IntKeys() }},
	"maps/PassToFunction": {mapsSrc, "PassToFunction", func() interface{} { return maps.PassToFunction() }},

	// mapadvanced
	"mapadvanced/LookupExistingKey": {mapadvancedSrc, "LookupExistingKey", func() interface{} { return mapadvanced.LookupExistingKey() }},
	"mapadvanced/LookupWithDefault": {mapadvancedSrc, "LookupWithDefault", func() interface{} { return mapadvanced.LookupWithDefault() }},
	"mapadvanced/AsCounter":         {mapadvancedSrc, "AsCounter", func() interface{} { return mapadvanced.AsCounter() }},
	"mapadvanced/WithStringValues":  {mapadvancedSrc, "WithStringValues", func() interface{} { return mapadvanced.WithStringValues() }},
	"mapadvanced/BuildFromLoop":     {mapadvancedSrc, "BuildFromLoop", func() interface{} { return mapadvanced.BuildFromLoop() }},
	"mapadvanced/DeleteAndReinsert": {mapadvancedSrc, "DeleteAndReinsert", func() interface{} { return mapadvanced.DeleteAndReinsert() }},

	// multiassign
	"multiassign/Swap":            {multiassignSrc, "Swap", func() interface{} { return multiassign.Swap() }},
	"multiassign/FromFunction":    {multiassignSrc, "FromFunction", func() interface{} { return multiassign.FromFunction() }},
	"multiassign/ThreeValues":     {multiassignSrc, "ThreeValues", func() interface{} { return multiassign.ThreeValues() }},
	"multiassign/InLoop":          {multiassignSrc, "InLoop", func() interface{} { return multiassign.InLoop() }},
	"multiassign/DiscardWithBlank": {multiassignSrc, "DiscardWithBlank", func() interface{} { return multiassign.DiscardWithBlank() }},

	// namedreturn
	"namedreturn/Basic":     {namedreturnSrc, "Basic", func() interface{} { return namedreturn.Basic() }},
	"namedreturn/Multiple":  {namedreturnSrc, "Multiple", func() interface{} { return namedreturn.Multiple() }},
	"namedreturn/ZeroValue": {namedreturnSrc, "ZeroValue", func() interface{} { return namedreturn.ZeroValue() }},

	// recursion
	"recursion/TailRecursionPattern": {recursionSrc, "TailRecursionPattern", func() interface{} { return recursion.TailRecursionPattern() }},
	"recursion/ReverseSlice":         {recursionSrc, "ReverseSlice", func() interface{} { return recursion.ReverseSlice() }},
	"recursion/TowerOfHanoi":         {recursionSrc, "TowerOfHanoi", func() interface{} { return recursion.TowerOfHanoi() }},
	"recursion/MaxSlice":             {recursionSrc, "MaxSlice", func() interface{} { return recursion.MaxSlice() }},
	"recursion/Ackermann":            {recursionSrc, "Ackermann", func() interface{} { return recursion.Ackermann() }},
	"recursion/BinarySearch":         {recursionSrc, "BinarySearch", func() interface{} { return recursion.BinarySearch() }},

	// scope
	"scope/IfInitShortVar":          {scopeSrc, "IfInitShortVar", func() interface{} { return scope.IfInitShortVar() }},
	"scope/IfInitMultiCondition":    {scopeSrc, "IfInitMultiCondition", func() interface{} { return scope.IfInitMultiCondition() }},
	"scope/NestedScopes":            {scopeSrc, "NestedScopes", func() interface{} { return scope.NestedScopes() }},
	"scope/ForScopeIsolation":       {scopeSrc, "ForScopeIsolation", func() interface{} { return scope.ForScopeIsolation() }},
	"scope/MultipleBlockScopes":     {scopeSrc, "MultipleBlockScopes", func() interface{} { return scope.MultipleBlockScopes() }},
	"scope/ClosureCapturesOuterScope": {scopeSrc, "ClosureCapturesOuterScope", func() interface{} { return scope.ClosureCapturesOuterScope() }},

	// slices
	"slices/MakeLen":           {slicesSrc, "MakeLen", func() interface{} { return slices.MakeLen() }},
	"slices/Append":            {slicesSrc, "Append", func() interface{} { return slices.Append() }},
	"slices/ElementAssignment": {slicesSrc, "ElementAssignment", func() interface{} { return slices.ElementAssignment() }},
	"slices/ForRange":          {slicesSrc, "ForRange", func() interface{} { return slices.ForRange() }},
	"slices/ForRangeIndex":     {slicesSrc, "ForRangeIndex", func() interface{} { return slices.ForRangeIndex() }},
	"slices/GrowMultiple":      {slicesSrc, "GrowMultiple", func() interface{} { return slices.GrowMultiple() }},
	"slices/PassToFunction":    {slicesSrc, "PassToFunction", func() interface{} { return slices.PassToFunction() }},
	"slices/LenCap":            {slicesSrc, "LenCap", func() interface{} { return slices.LenCap() }},

	// slicing
	"slicing/SubSliceBasic":            {slicingSrc, "SubSliceBasic", func() interface{} { return slicing.SubSliceBasic() }},
	"slicing/SubSliceLen":              {slicingSrc, "SubSliceLen", func() interface{} { return slicing.SubSliceLen() }},
	"slicing/SubSliceFromStart":        {slicingSrc, "SubSliceFromStart", func() interface{} { return slicing.SubSliceFromStart() }},
	"slicing/SubSliceToEnd":            {slicingSrc, "SubSliceToEnd", func() interface{} { return slicing.SubSliceToEnd() }},
	"slicing/SubSliceCopy":             {slicingSrc, "SubSliceCopy", func() interface{} { return slicing.SubSliceCopy() }},
	"slicing/SubSliceChained":          {slicingSrc, "SubSliceChained", func() interface{} { return slicing.SubSliceChained() }},
	"slicing/SubSliceModifiesOriginal": {slicingSrc, "SubSliceModifiesOriginal", func() interface{} { return slicing.SubSliceModifiesOriginal() }},

	// strings_pkg
	"strings_pkg/Concat":       {stringsPkgSrc, "Concat", func() interface{} { return strings_pkg.Concat() }},
	"strings_pkg/ConcatLoop":   {stringsPkgSrc, "ConcatLoop", func() interface{} { return strings_pkg.ConcatLoop() }},
	"strings_pkg/Len":          {stringsPkgSrc, "Len", func() interface{} { return strings_pkg.Len() }},
	"strings_pkg/Index":        {stringsPkgSrc, "Index", func() interface{} { return strings_pkg.Index() }},
	"strings_pkg/Comparison":   {stringsPkgSrc, "Comparison", func() interface{} { return strings_pkg.Comparison() }},
	"strings_pkg/Equality":     {stringsPkgSrc, "Equality", func() interface{} { return strings_pkg.Equality() }},
	"strings_pkg/EmptyCheck":   {stringsPkgSrc, "EmptyCheck", func() interface{} { return strings_pkg.EmptyCheck() }},

	// switch
	"switch/Simple":       {switchSrc, "Simple", func() interface{} { return switch_pkg.Simple() }},
	"switch/Default":      {switchSrc, "Default", func() interface{} { return switch_pkg.Default() }},
	"switch/MultiCase":    {switchSrc, "MultiCase", func() interface{} { return switch_pkg.MultiCase() }},
	"switch/NoCondition":  {switchSrc, "NoCondition", func() interface{} { return switch_pkg.NoCondition() }},
	"switch/WithInit":     {switchSrc, "WithInit", func() interface{} { return switch_pkg.WithInit() }},
	"switch/StringCases":  {switchSrc, "StringCases", func() interface{} { return switch_pkg.StringCases() }},
	"switch/Fallthrough":  {switchSrc, "Fallthrough", func() interface{} { return switch_pkg.Fallthrough() }},
	"switch/Nested":       {switchSrc, "Nested", func() interface{} { return switch_pkg.Nested() }},

	// typeconv
	"typeconv/IntToFloat64":         {typeconvSrc, "IntToFloat64", func() interface{} { return typeconv.IntToFloat64() }},
	"typeconv/Float64Arithmetic":    {typeconvSrc, "Float64Arithmetic", func() interface{} { return typeconv.Float64Arithmetic() }},
	"typeconv/StringToByteConversion": {typeconvSrc, "StringToByteConversion", func() interface{} { return typeconv.StringToByteConversion() }},
	"typeconv/IntStringConversion":  {typeconvSrc, "IntStringConversion", func() interface{} { return typeconv.IntStringConversion() }},
	"typeconv/StringIntConversion":  {typeconvSrc, "StringIntConversion", func() interface{} { return typeconv.StringIntConversion() }},

	// variables
	"variables/DeclareAndUse":  {variablesSrc, "DeclareAndUse", func() interface{} { return variables.DeclareAndUse() }},
	"variables/Reassignment":   {variablesSrc, "Reassignment", func() interface{} { return variables.Reassignment() }},
	"variables/MultipleDecl":   {variablesSrc, "MultipleDecl", func() interface{} { return variables.MultipleDecl() }},
	"variables/ZeroValues":     {variablesSrc, "ZeroValues", func() interface{} { return variables.ZeroValues() }},
	"variables/StringZeroValue": {variablesSrc, "StringZeroValue", func() interface{} { return variables.StringZeroValue() }},
	"variables/Shadowing":      {variablesSrc, "Shadowing", func() interface{} { return variables.Shadowing() }},
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
			result, err := prog.Run(tc.funcName)
			if err != nil {
				t.Fatalf("Run error: %v", err)
			}
			expected := tc.native()
			compareResults(t, result, expected)
		})
	}
}

// compareResults compares interpreter result with native result
func compareResults(t *testing.T, result, expected interface{}) {
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
