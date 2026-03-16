package tests

import (
	_ "embed"
	"strings"
	"testing"
	"time"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
	"git.woa.com/youngjin/gig/tests/testdata/tricky"
)

//go:embed testdata/tricky/main.go
var trickySrc string

var trickyTests = map[string]testCase{
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
}

// toMainPackage converts a source file to package main
func toMainPackageTricky(src string) string {
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "package ") {
			lines[i] = "package main"
			break
		}
	}
	return strings.Join(lines, "\n")
}

// TestTrickyCases runs all tricky test cases
func TestTrickyCases(t *testing.T) {
	for name, tc := range trickyTests {
		t.Run(name, func(t *testing.T) {
			src := toMainPackageTricky(tc.src)
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

			compareResults(t, result, expected)

			ratio := float64(interpDuration) / float64(nativeDuration)
			t.Logf("interp: %v, native: %v, ratio: %.1fx", interpDuration, nativeDuration, ratio)
		})
	}
}
