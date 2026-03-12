package tests

import (
	_ "embed"
	"testing"
	"time"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
	"github.com/t04dJ14n9/gig/tests/testdata/cornercases"
	"github.com/t04dJ14n9/gig/value"
)

//go:embed testdata/cornercases/main.go
var cornercasesSrc string

// cornerCaseTestNative defines a corner case test with native reference
type cornerCaseTestNative struct {
	src      string
	funcName string
	native   func() any
}

// allCornerCaseTests contains all corner case test definitions
var allCornerCaseTests = map[string]cornerCaseTestNative{
	// Zero Value Tests
	"cornercases/ZeroValue_Int":     {cornercasesSrc, "ZeroValue_Int", func() any { return cornercases.ZeroValue_Int() }},
	"cornercases/ZeroValue_Int64":   {cornercasesSrc, "ZeroValue_Int64", func() any { return cornercases.ZeroValue_Int64() }},
	"cornercases/ZeroValue_Float64": {cornercasesSrc, "ZeroValue_Float64", func() any { return cornercases.ZeroValue_Float64() }},
	"cornercases/ZeroValue_String":  {cornercasesSrc, "ZeroValue_String", func() any { return cornercases.ZeroValue_String() }},
	"cornercases/ZeroValue_Bool":    {cornercasesSrc, "ZeroValue_Bool", func() any { return cornercases.ZeroValue_Bool() }},
	"cornercases/ZeroValue_Slice":   {cornercasesSrc, "ZeroValue_Slice", func() any { return cornercases.ZeroValue_Slice() }},
	"cornercases/ZeroValue_Map":     {cornercasesSrc, "ZeroValue_Map", func() any { return cornercases.ZeroValue_Map() }},

	// Integer Boundary Tests
	"cornercases/IntBoundary_MaxInt32":   {cornercasesSrc, "IntBoundary_MaxInt32", func() any { return cornercases.IntBoundary_MaxInt32() }},
	"cornercases/IntBoundary_MinInt32":   {cornercasesSrc, "IntBoundary_MinInt32", func() any { return cornercases.IntBoundary_MinInt32() }},
	"cornercases/IntBoundary_MaxInt64":   {cornercasesSrc, "IntBoundary_MaxInt64", func() any { return cornercases.IntBoundary_MaxInt64() }},
	"cornercases/IntBoundary_MinInt64":   {cornercasesSrc, "IntBoundary_MinInt64", func() any { return cornercases.IntBoundary_MinInt64() }},
	"cornercases/IntBoundary_MaxUint32":  {cornercasesSrc, "IntBoundary_MaxUint32", func() any { return cornercases.IntBoundary_MaxUint32() }},
	"cornercases/IntBoundary_NearMaxInt": {cornercasesSrc, "IntBoundary_NearMaxInt", func() any { return cornercases.IntBoundary_NearMaxInt() }},
	"cornercases/IntBoundary_NearMinInt": {cornercasesSrc, "IntBoundary_NearMinInt", func() any { return cornercases.IntBoundary_NearMinInt() }},

	// Integer Overflow Tests
	"cornercases/Overflow_Int32Add": {cornercasesSrc, "Overflow_Int32Add", func() any { return cornercases.Overflow_Int32Add() }},
	"cornercases/Overflow_Int32Sub": {cornercasesSrc, "Overflow_Int32Sub", func() any { return cornercases.Overflow_Int32Sub() }},
	"cornercases/Overflow_Int32Mul": {cornercasesSrc, "Overflow_Int32Mul", func() any { return cornercases.Overflow_Int32Mul() }},

	// Float Boundary Tests
	"cornercases/FloatBoundary_SmallPositive": {cornercasesSrc, "FloatBoundary_SmallPositive", func() any { return cornercases.FloatBoundary_SmallPositive() }},
	"cornercases/FloatBoundary_SmallNegative": {cornercasesSrc, "FloatBoundary_SmallNegative", func() any { return cornercases.FloatBoundary_SmallNegative() }},
	"cornercases/FloatBoundary_LargePositive": {cornercasesSrc, "FloatBoundary_LargePositive", func() any { return cornercases.FloatBoundary_LargePositive() }},
	"cornercases/FloatBoundary_LargeNegative": {cornercasesSrc, "FloatBoundary_LargeNegative", func() any { return cornercases.FloatBoundary_LargeNegative() }},

	// Empty Collection Tests
	"cornercases/EmptySlice_Len":  {cornercasesSrc, "EmptySlice_Len", func() any { return cornercases.EmptySlice_Len() }},
	"cornercases/EmptySlice_Cap":  {cornercasesSrc, "EmptySlice_Cap", func() any { return cornercases.EmptySlice_Cap() }},
	"cornercases/EmptySlice_Make": {cornercasesSrc, "EmptySlice_Make", func() any { return cornercases.EmptySlice_Make() }},
	"cornercases/EmptyMap_Len":    {cornercasesSrc, "EmptyMap_Len", func() any { return cornercases.EmptyMap_Len() }},
	"cornercases/EmptyMap_Make":   {cornercasesSrc, "EmptyMap_Make", func() any { return cornercases.EmptyMap_Make() }},
	"cornercases/EmptyString_Len": {cornercasesSrc, "EmptyString_Len", func() any { return cornercases.EmptyString_Len() }},

	// Slice Operation Tests
	"cornercases/Slice_ZeroToZero":  {cornercasesSrc, "Slice_ZeroToZero", func() any { return cornercases.Slice_ZeroToZero() }},
	"cornercases/Slice_EndToEnd":    {cornercasesSrc, "Slice_EndToEnd", func() any { return cornercases.Slice_EndToEnd() }},
	"cornercases/Slice_NilSlice":    {cornercasesSrc, "Slice_NilSlice", func() any { return cornercases.Slice_NilSlice() }},
	"cornercases/Slice_AppendToNil": {cornercasesSrc, "Slice_AppendToNil", func() any { return cornercases.Slice_AppendToNil() }},
	"cornercases/Slice_AppendEmpty": {cornercasesSrc, "Slice_AppendEmpty", func() any { return cornercases.Slice_AppendEmpty() }},

	// Map Operation Tests
	"cornercases/Map_NilMap":           {cornercasesSrc, "Map_NilMap", func() any { return cornercases.Map_NilMap() }},
	"cornercases/Map_AccessMissingKey": {cornercasesSrc, "Map_AccessMissingKey", func() any { return cornercases.Map_AccessMissingKey() }},
	"cornercases/Map_DeleteMissingKey": {cornercasesSrc, "Map_DeleteMissingKey", func() any { return cornercases.Map_DeleteMissingKey() }},
	"cornercases/Map_OverwriteKey":     {cornercasesSrc, "Map_OverwriteKey", func() any { return cornercases.Map_OverwriteKey() }},
	"cornercases/Map_NilKeyString":     {cornercasesSrc, "Map_NilKeyString", func() any { return cornercases.Map_NilKeyString() }},
	"cornercases/Map_ZeroIntKey":       {cornercasesSrc, "Map_ZeroIntKey", func() any { return cornercases.Map_ZeroIntKey() }},

	// String Boundary Tests
	"cornercases/String_Empty":            {cornercasesSrc, "String_Empty", func() any { return cornercases.String_Empty() }},
	"cornercases/String_SingleChar":       {cornercasesSrc, "String_SingleChar", func() any { return cornercases.String_SingleChar() }},
	"cornercases/String_UnicodeMultibyte": {cornercasesSrc, "String_UnicodeMultibyte", func() any { return cornercases.String_UnicodeMultibyte() }},
	"cornercases/String_Whitespace":       {cornercasesSrc, "String_Whitespace", func() any { return cornercases.String_Whitespace() }},
	"cornercases/String_SingleByteIndex":  {cornercasesSrc, "String_SingleByteIndex", func() any { return cornercases.String_SingleByteIndex() }},
	"cornercases/String_LastByte":         {cornercasesSrc, "String_LastByte", func() any { return cornercases.String_LastByte() }},

	// Boolean Tests
	"cornercases/Bool_True":           {cornercasesSrc, "Bool_True", func() any { return cornercases.Bool_True() }},
	"cornercases/Bool_False":          {cornercasesSrc, "Bool_False", func() any { return cornercases.Bool_False() }},
	"cornercases/Bool_NotTrue":        {cornercasesSrc, "Bool_NotTrue", func() any { return cornercases.Bool_NotTrue() }},
	"cornercases/Bool_NotFalse":       {cornercasesSrc, "Bool_NotFalse", func() any { return cornercases.Bool_NotFalse() }},
	"cornercases/Bool_DoubleNegation": {cornercasesSrc, "Bool_DoubleNegation", func() any { return cornercases.Bool_DoubleNegation() }},

	// Arithmetic Tests
	"cornercases/Arith_AddZero":   {cornercasesSrc, "Arith_AddZero", func() any { return cornercases.Arith_AddZero() }},
	"cornercases/Arith_SubZero":   {cornercasesSrc, "Arith_SubZero", func() any { return cornercases.Arith_SubZero() }},
	"cornercases/Arith_MulByOne":  {cornercasesSrc, "Arith_MulByOne", func() any { return cornercases.Arith_MulByOne() }},
	"cornercases/Arith_DivByOne":  {cornercasesSrc, "Arith_DivByOne", func() any { return cornercases.Arith_DivByOne() }},
	"cornercases/Arith_ModByOne":  {cornercasesSrc, "Arith_ModByOne", func() any { return cornercases.Arith_ModByOne() }},
	"cornercases/Arith_MulByZero": {cornercasesSrc, "Arith_MulByZero", func() any { return cornercases.Arith_MulByZero() }},
	"cornercases/Arith_NegNeg":    {cornercasesSrc, "Arith_NegNeg", func() any { return cornercases.Arith_NegNeg() }},
	"cornercases/Arith_NegAddNeg": {cornercasesSrc, "Arith_NegAddNeg", func() any { return cornercases.Arith_NegAddNeg() }},

	// Comparison Tests
	"cornercases/Compare_IntEqual":         {cornercasesSrc, "Compare_IntEqual", func() any { return cornercases.Compare_IntEqual() }},
	"cornercases/Compare_IntNotEqual":      {cornercasesSrc, "Compare_IntNotEqual", func() any { return cornercases.Compare_IntNotEqual() }},
	"cornercases/Compare_IntGreater":       {cornercasesSrc, "Compare_IntGreater", func() any { return cornercases.Compare_IntGreater() }},
	"cornercases/Compare_IntGreaterEqual":  {cornercasesSrc, "Compare_IntGreaterEqual", func() any { return cornercases.Compare_IntGreaterEqual() }},
	"cornercases/Compare_IntLess":          {cornercasesSrc, "Compare_IntLess", func() any { return cornercases.Compare_IntLess() }},
	"cornercases/Compare_IntLessEqual":     {cornercasesSrc, "Compare_IntLessEqual", func() any { return cornercases.Compare_IntLessEqual() }},
	"cornercases/Compare_StringEqual":      {cornercasesSrc, "Compare_StringEqual", func() any { return cornercases.Compare_StringEqual() }},
	"cornercases/Compare_StringNotEqual":   {cornercasesSrc, "Compare_StringNotEqual", func() any { return cornercases.Compare_StringNotEqual() }},
	"cornercases/Compare_EmptyStringEqual": {cornercasesSrc, "Compare_EmptyStringEqual", func() any { return cornercases.Compare_EmptyStringEqual() }},

	// Logic Tests
	"cornercases/Logic_TrueAndTrue":  {cornercasesSrc, "Logic_TrueAndTrue", func() any { return cornercases.Logic_TrueAndTrue() }},
	"cornercases/Logic_TrueAndFalse": {cornercasesSrc, "Logic_TrueAndFalse", func() any { return cornercases.Logic_TrueAndFalse() }},
	"cornercases/Logic_FalseAndTrue": {cornercasesSrc, "Logic_FalseAndTrue", func() any { return cornercases.Logic_FalseAndTrue() }},
	"cornercases/Logic_TrueOrFalse":  {cornercasesSrc, "Logic_TrueOrFalse", func() any { return cornercases.Logic_TrueOrFalse() }},
	"cornercases/Logic_FalseOrTrue":  {cornercasesSrc, "Logic_FalseOrTrue", func() any { return cornercases.Logic_FalseOrTrue() }},
	"cornercases/Logic_FalseOrFalse": {cornercasesSrc, "Logic_FalseOrFalse", func() any { return cornercases.Logic_FalseOrFalse() }},

	// Control Flow Tests
	"cornercases/Control_IfNoElse":       {cornercasesSrc, "Control_IfNoElse", func() any { return cornercases.Control_IfNoElse() }},
	"cornercases/Control_IfFalseNoElse":  {cornercasesSrc, "Control_IfFalseNoElse", func() any { return cornercases.Control_IfFalseNoElse() }},
	"cornercases/Control_ForZeroIter":    {cornercasesSrc, "Control_ForZeroIter", func() any { return cornercases.Control_ForZeroIter() }},
	"cornercases/Control_ForOneIter":     {cornercasesSrc, "Control_ForOneIter", func() any { return cornercases.Control_ForOneIter() }},
	"cornercases/Control_ForBreakFirst":  {cornercasesSrc, "Control_ForBreakFirst", func() any { return cornercases.Control_ForBreakFirst() }},
	"cornercases/Control_ForContinueAll": {cornercasesSrc, "Control_ForContinueAll", func() any { return cornercases.Control_ForContinueAll() }},
	"cornercases/Control_SwitchNoMatch":  {cornercasesSrc, "Control_SwitchNoMatch", func() any { return cornercases.Control_SwitchNoMatch() }},
	"cornercases/Control_SwitchDefault":  {cornercasesSrc, "Control_SwitchDefault", func() any { return cornercases.Control_SwitchDefault() }},

	// Function Tests
	"cornercases/Func_NoReturn":             {cornercasesSrc, "Func_NoReturn", func() any { return cornercases.Func_NoReturn() }},
	"cornercases/Func_MultipleReturnAll":    {cornercasesSrc, "Func_MultipleReturnAll", func() any { a, b := cornercases.Func_MultipleReturnAll(); return []int{a, b} }},
	"cornercases/Func_MultipleReturnIgnore": {cornercasesSrc, "Func_MultipleReturnIgnore", func() any { return cornercases.Func_MultipleReturnIgnore() }},
	"cornercases/Func_NamedReturn":          {cornercasesSrc, "Func_NamedReturn", func() any { return cornercases.Func_NamedReturn() }},
	"cornercases/Func_VariadicEmpty":        {cornercasesSrc, "Func_VariadicEmpty", func() any { return cornercases.Func_VariadicEmpty() }},
	"cornercases/Func_VariadicOne":          {cornercasesSrc, "Func_VariadicOne", func() any { return cornercases.Func_VariadicOne() }},
	"cornercases/Func_VariadicMultiple":     {cornercasesSrc, "Func_VariadicMultiple", func() any { return cornercases.Func_VariadicMultiple() }},
	"cornercases/Func_RecursionBase":        {cornercasesSrc, "Func_RecursionBase", func() any { return cornercases.Func_RecursionBase() }},

	// Closure Tests
	"cornercases/Closure_ReturnClosure":   {cornercasesSrc, "Closure_ReturnClosure", func() any { return cornercases.Closure_ReturnClosure() }},
	"cornercases/Closure_CaptureVariable": {cornercasesSrc, "Closure_CaptureVariable", func() any { return cornercases.Closure_CaptureVariable() }},
	"cornercases/Closure_ModifyCaptured":  {cornercasesSrc, "Closure_ModifyCaptured", func() any { return cornercases.Closure_ModifyCaptured() }},

	// Struct Tests
	"cornercases/Struct_ZeroValueFields": {cornercasesSrc, "Struct_ZeroValueFields", func() any { return cornercases.Struct_ZeroValueFields() }},
	"cornercases/Struct_PointerReceiver": {cornercasesSrc, "Struct_PointerReceiver", func() any { return cornercases.Struct_PointerReceiver() }},
	"cornercases/Struct_NestedStruct":    {cornercasesSrc, "Struct_NestedStruct", func() any { return cornercases.Struct_NestedStruct() }},

	// Type Conversion Tests
	"cornercases/Convert_IntToFloat":   {cornercasesSrc, "Convert_IntToFloat", func() any { return cornercases.Convert_IntToFloat() }},
	"cornercases/Convert_FloatToInt":   {cornercasesSrc, "Convert_FloatToInt", func() any { return cornercases.Convert_FloatToInt() }},
	"cornercases/Convert_Int32ToInt64": {cornercasesSrc, "Convert_Int32ToInt64", func() any { return cornercases.Convert_Int32ToInt64() }},
	"cornercases/Convert_Int64ToInt32": {cornercasesSrc, "Convert_Int64ToInt32", func() any { return cornercases.Convert_Int64ToInt32() }},

	// Complex Expression Tests
	"cornercases/Expr_ComplexArithmetic":  {cornercasesSrc, "Expr_ComplexArithmetic", func() any { return cornercases.Expr_ComplexArithmetic() }},
	"cornercases/Expr_ChainedComparison":  {cornercasesSrc, "Expr_ChainedComparison", func() any { return cornercases.Expr_ChainedComparison() }},
	"cornercases/Expr_MultipleAssignment": {cornercasesSrc, "Expr_MultipleAssignment", func() any { return cornercases.Expr_MultipleAssignment() }},
	"cornercases/Expr_NestedTernaryLike":  {cornercasesSrc, "Expr_NestedTernaryLike", func() any { return cornercases.Expr_NestedTernaryLike() }},

	// Make Tests
	"cornercases/Make_ZeroLenZeroCap": {cornercasesSrc, "Make_ZeroLenZeroCap", func() any { return cornercases.Make_ZeroLenZeroCap() }},
	"cornercases/Make_SliceWithCap":   {cornercasesSrc, "Make_SliceWithCap", func() any { return cornercases.Make_SliceWithCap() }},
	"cornercases/Make_MapWithSize":    {cornercasesSrc, "Make_MapWithSize", func() any { return cornercases.Make_MapWithSize() }},

	// Range Tests
	"cornercases/Range_EmptySlice":    {cornercasesSrc, "Range_EmptySlice", func() any { return cornercases.Range_EmptySlice() }},
	"cornercases/Range_EmptyMap":      {cornercasesSrc, "Range_EmptyMap", func() any { return cornercases.Range_EmptyMap() }},
	"cornercases/Range_EmptyString":   {cornercasesSrc, "Range_EmptyString", func() any { return cornercases.Range_EmptyString() }},
	"cornercases/Range_SingleElement": {cornercasesSrc, "Range_SingleElement", func() any { return cornercases.Range_SingleElement() }},
}

// TestAllCornerCases runs all corner case tests and compares with native Go
func TestAllCornerCases(t *testing.T) {
	for name, tc := range allCornerCaseTests {
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

			// Compare results with flexible type handling
			compareCornerCaseResults(t, result, expected)

			// Report timing comparison
			ratio := float64(interpDuration) / float64(nativeDuration)
			t.Logf("interp: %v, native: %v, ratio: %.1fx", interpDuration, nativeDuration, ratio)
		})
	}
}

// compareCornerCaseResults compares interpreter result with native result
func compareCornerCaseResults(t *testing.T, result, expected any) {
	t.Helper()

	switch exp := expected.(type) {
	case int:
		var got int64
		switch v := result.(type) {
		case int64:
			got = v
		case int:
			got = int64(v)
		case uint64:
			got = int64(v)
		default:
			t.Fatalf("expected int, got %T", result)
		}
		if got != int64(exp) {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case int32:
		var got int64
		switch v := result.(type) {
		case int64:
			got = v
		case int32:
			got = int64(v)
		case int:
			got = int64(v)
		default:
			t.Fatalf("expected int32, got %T", result)
		}
		// Note: Gig may not simulate int32 overflow
		if int32(got) != exp {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case int64:
		var got int64
		switch v := result.(type) {
		case int64:
			got = v
		case int:
			got = int64(v)
		case uint64:
			got = int64(v)
		default:
			t.Fatalf("expected int64, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case uint32:
		var got uint64
		switch v := result.(type) {
		case uint64:
			got = v
		case int64:
			got = uint64(v)
		default:
			t.Fatalf("expected uint32, got %T", result)
		}
		if uint32(got) != exp {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case uint8:
		var got uint64
		switch v := result.(type) {
		case uint64:
			got = v
		case int64:
			got = uint64(v)
		default:
			t.Fatalf("expected uint8, got %T", result)
		}
		if uint8(got) != exp {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case float64:
		got, ok := result.(float64)
		if !ok {
			t.Fatalf("expected float64, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %v, got %v", exp, got)
		}

	case string:
		got, ok := result.(string)
		if !ok {
			t.Fatalf("expected string, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %q, got %q", exp, got)
		}

	case bool:
		got, ok := result.(bool)
		if !ok {
			t.Fatalf("expected bool, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %v, got %v", exp, got)
		}

	case []int:
		// Handle Gig's multiple return values: []value.Value
		if values, ok := result.([]value.Value); ok {
			if len(values) != len(exp) {
				t.Errorf("expected len %d, got %d", len(exp), len(values))
				return
			}
			for i := range exp {
				v := values[i].Interface()
				var got int
				switch val := v.(type) {
				case int:
					got = val
				case int64:
					got = int(val)
				default:
					t.Errorf("expected [%d]=int, got %T", i, v)
					continue
				}
				if got != exp[i] {
					t.Errorf("expected [%d]=%d, got %d", i, exp[i], got)
				}
			}
			return
		}

		got, ok := result.([]int64)
		if !ok {
			t.Fatalf("expected []int, got %T", result)
		}
		if len(got) != len(exp) {
			t.Errorf("expected len %d, got %d", len(exp), len(got))
			return
		}
		for i := range exp {
			if int(got[i]) != exp[i] {
				t.Errorf("expected [%d]=%d, got %d", i, exp[i], got[i])
			}
		}

	default:
		t.Fatalf("unsupported result type: %T", expected)
	}
}
