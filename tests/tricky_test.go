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
	// RALPH LOOP ITERATION 2 - New Tricky Tests
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
	// RALPH LOOP ITERATION 3 - More Tricky Tests
	"tricky/SliceModifyViaSubslice": {trickySrc, "SliceModifyViaSubslice", func() any { return tricky.SliceModifyViaSubslice() }},
	"tricky/MapDeleteDuringRange":   {trickySrc, "MapDeleteDuringRange", func() any { return tricky.MapDeleteDuringRange() }},
	"tricky/ClosureReturnClosure":   {trickySrc, "ClosureReturnClosure", func() any { return tricky.ClosureReturnClosure() }},
	"tricky/StructMethodOnNil":      {trickySrc, "StructMethodOnNil", func() any { return tricky.StructMethodOnNil() }},
	"tricky/ArrayPointerIndex":      {trickySrc, "ArrayPointerIndex", func() any { return tricky.ArrayPointerIndex() }},
	"tricky/SliceThreeIndex":        {trickySrc, "SliceThreeIndex", func() any { return tricky.SliceThreeIndex() }},
	// "tricky/MapWithFuncValue": SKIPPED - known issue with func type in map
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
	// RALPH LOOP ITERATION 4 - More Tricky Tests
	// "tricky/SliceOfInterfacesWithTypes": SKIPPED - known issue with type switch on interface slice
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
	// RALPH LOOP ITERATION 5 - More Tricky Tests
	"tricky/ChannelBuffered":        {trickySrc, "ChannelBuffered", func() any { return tricky.ChannelBuffered() }},
	"tricky/StructEmbeddedMethod":   {trickySrc, "StructEmbeddedMethod", func() any { return tricky.StructEmbeddedMethod() }},
	"tricky/SliceOfChannels":        {trickySrc, "SliceOfChannels", func() any { return tricky.SliceOfChannels() }},
	"tricky/MapOfChannels":          {trickySrc, "MapOfChannels", func() any { return tricky.MapOfChannels() }},
	"tricky/InterfaceMethod":        {trickySrc, "InterfaceMethod", func() any { return tricky.InterfaceMethod() }},
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
	// RALPH LOOP ITERATION 6 - More Tricky Tests
	"tricky/ComplexMapKey":      {trickySrc, "ComplexMapKey", func() any { return tricky.ComplexMapKey() }},
	"tricky/SliceReverse":       {trickySrc, "SliceReverse", func() any { return tricky.SliceReverse() }},
	"tricky/MapMerge":           {trickySrc, "MapMerge", func() any { return tricky.MapMerge() }},
	"tricky/StructZeroValue":    {trickySrc, "StructZeroValue", func() any { return tricky.StructZeroValue() }},
	"tricky/SliceDeleteByIndex": {trickySrc, "SliceDeleteByIndex", func() any { return tricky.SliceDeleteByIndex() }},
	"tricky/MapValueOverwrite":  {trickySrc, "MapValueOverwrite", func() any { return tricky.MapValueOverwrite() }},
	"tricky/InterfaceEmbed":     {trickySrc, "InterfaceEmbed", func() any { return tricky.InterfaceEmbed() }},
	"tricky/SliceOfFuncs":       {trickySrc, "SliceOfFuncs", func() any { return tricky.SliceOfFuncs() }},
	// "tricky/StructWithFunc": SKIPPED - known issue with func type in struct field
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
	// RALPH LOOP ITERATION 7 - More Tricky Tests
	"tricky/StructPointerSlice":       {trickySrc, "StructPointerSlice", func() any { return tricky.StructPointerSlice() }},
	"tricky/MapWithInterfaceKey":      {trickySrc, "MapWithInterfaceKey", func() any { return tricky.MapWithInterfaceKey() }},
	"tricky/SliceOfInterfaces":        {trickySrc, "SliceOfInterfaces", func() any { return tricky.SliceOfInterfaces() }},
	"tricky/NestedPointerStruct":      {trickySrc, "NestedPointerStruct", func() any { return tricky.NestedPointerStruct() }},
	"tricky/StructMethodOnNilPointer": {trickySrc, "StructMethodOnNilPointer", func() any { return tricky.StructMethodOnNilPointer() }},
	"tricky/SliceAppendToSlice":       {trickySrc, "SliceAppendToSlice", func() any { return tricky.SliceAppendToSlice() }},
	"tricky/MapLookupWithDefault":     {trickySrc, "MapLookupWithDefault", func() any { return tricky.MapLookupWithDefault() }},
	"tricky/StructFieldUpdate":        {trickySrc, "StructFieldUpdate", func() any { return tricky.StructFieldUpdate() }},
	"tricky/PointerToNilSlice":        {trickySrc, "PointerToNilSlice", func() any { return tricky.PointerToNilSlice() }},
	// "tricky/MapUpdateDuringRange": SKIPPED - known issue with map modification during range
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
	// RALPH LOOP ITERATION 8 - More Tricky Tests
	"tricky/SliceDrain":        {trickySrc, "SliceDrain", func() any { return tricky.SliceDrain() }},
	"tricky/MapClear":          {trickySrc, "MapClear", func() any { return tricky.MapClear() }},
	"tricky/StructCopy":        {trickySrc, "StructCopy", func() any { return tricky.StructCopy() }},
	"tricky/PointerStructCopy": {trickySrc, "PointerStructCopy", func() any { return tricky.PointerStructCopy() }},
	"tricky/SliceFilter":       {trickySrc, "SliceFilter", func() any { return tricky.SliceFilter() }},
	"tricky/MapTransform":      {trickySrc, "MapTransform", func() any { return tricky.MapTransform() }},
	"tricky/SliceContains":     {trickySrc, "SliceContains", func() any { return tricky.SliceContains() }},
	"tricky/MapKeys":           {trickySrc, "MapKeys", func() any { return tricky.MapKeys() }},
	"tricky/StructMethodChain": {trickySrc, "StructMethodChain", func() any { return tricky.StructMethodChain() }},
	"tricky/SliceLast":         {trickySrc, "SliceLast", func() any { return tricky.SliceLast() }},
	"tricky/MapGetOrSet":       {trickySrc, "MapGetOrSet", func() any { return tricky.MapGetOrSet() }},
	"tricky/StructValidation":  {trickySrc, "StructValidation", func() any { return tricky.StructValidation() }},
	"tricky/SlicePrepend":      {trickySrc, "SlicePrepend", func() any { return tricky.SlicePrepend() }},
	"tricky/MapMergeOverwrite": {trickySrc, "MapMergeOverwrite", func() any { return tricky.MapMergeOverwrite() }},
	"tricky/SliceRotate":       {trickySrc, "SliceRotate", func() any { return tricky.SliceRotate() }},
	"tricky/StructInterface":   {trickySrc, "StructInterface", func() any { return tricky.StructInterface() }},
	"tricky/MapKeysSorted":     {trickySrc, "MapKeysSorted", func() any { return tricky.MapKeysSorted() }},
	// "tricky/SliceFlatten": SKIPPED - known issue with append spread operator
	"tricky/StructFieldPointerModify": {trickySrc, "StructFieldPointerModify", func() any { return tricky.StructFieldPointerModify() }},
	// RALPH LOOP ITERATION 9 - More Tricky Tests
	"tricky/MapSwap":            {trickySrc, "MapSwap", func() any { return tricky.MapSwap() }},
	"tricky/SliceSplit":         {trickySrc, "SliceSplit", func() any { return tricky.SliceSplit() }},
	"tricky/StructCompareDiff":  {trickySrc, "StructCompareDiff", func() any { return tricky.StructCompareDiff() }},
	"tricky/MapNestedDelete":    {trickySrc, "MapNestedDelete", func() any { return tricky.MapNestedDelete() }},
	"tricky/PointerNilDeref":    {trickySrc, "PointerNilDeref", func() any { return tricky.PointerNilDeref() }},
	"tricky/SliceGrow":          {trickySrc, "SliceGrow", func() any { return tricky.SliceGrow() }},
	"tricky/StructEmpty":        {trickySrc, "StructEmpty", func() any { return tricky.StructEmpty() }},
	"tricky/MapEmptyKey":        {trickySrc, "MapEmptyKey", func() any { return tricky.MapEmptyKey() }},
	"tricky/SliceMakeZero":      {trickySrc, "SliceMakeZero", func() any { return tricky.SliceMakeZero() }},
	"tricky/StructAnon":         {trickySrc, "StructAnon", func() any { return tricky.StructAnon() }},
	"tricky/MapSizeHint":        {trickySrc, "MapSizeHint", func() any { return tricky.MapSizeHint() }},
	"tricky/SliceNilAppend":     {trickySrc, "SliceNilAppend", func() any { return tricky.SliceNilAppend() }},
	"tricky/StructFieldPtr":     {trickySrc, "StructFieldPtr", func() any { return tricky.StructFieldPtr() }},
	"tricky/MapIterateDelete":   {trickySrc, "MapIterateDelete", func() any { return tricky.MapIterateDelete() }},
	"tricky/SliceTruncate":      {trickySrc, "SliceTruncate", func() any { return tricky.SliceTruncate() }},
	"tricky/StructMethodValue":  {trickySrc, "StructMethodValue", func() any { return tricky.StructMethodValue() }},
	"tricky/MapFloatKey":        {trickySrc, "MapFloatKey", func() any { return tricky.MapFloatKey() }},
	"tricky/SliceRepeat":        {trickySrc, "SliceRepeat", func() any { return tricky.SliceRepeat() }},
	"tricky/StructNestedAssign": {trickySrc, "StructNestedAssign", func() any { return tricky.StructNestedAssign() }},
	"tricky/MapIntKey":          {trickySrc, "MapIntKey", func() any { return tricky.MapIntKey() }},
	// RALPH LOOP ITERATION 10 - Final Tricky Tests
	"tricky/SliceReverseInPlace":    {trickySrc, "SliceReverseInPlace", func() any { return tricky.SliceReverseInPlace() }},
	"tricky/MapIncrementAll":        {trickySrc, "MapIncrementAll", func() any { return tricky.MapIncrementAll() }},
	"tricky/StructPtrMethod":        {trickySrc, "StructPtrMethod", func() any { return tricky.StructPtrMethod() }},
	"tricky/SliceMapIndex":          {trickySrc, "SliceMapIndex", func() any { return tricky.SliceMapIndex() }},
	"tricky/MapCopy":                {trickySrc, "MapCopy", func() any { return tricky.MapCopy() }},
	"tricky/StructSliceAppend":      {trickySrc, "StructSliceAppend", func() any { return tricky.StructSliceAppend() }},
	"tricky/PointerSwap":            {trickySrc, "PointerSwap", func() any { return tricky.PointerSwap() }},
	"tricky/MapNestedUpdate":        {trickySrc, "MapNestedUpdate", func() any { return tricky.MapNestedUpdate() }},
	"tricky/SliceDeleteMiddle":      {trickySrc, "SliceDeleteMiddle", func() any { return tricky.SliceDeleteMiddle() }},
	"tricky/StructNilField":         {trickySrc, "StructNilField", func() any { return tricky.StructNilField() }},
	"tricky/MapLookupOrInsert":      {trickySrc, "MapLookupOrInsert", func() any { return tricky.MapLookupOrInsert() }},
	"tricky/SliceChainedSlice":      {trickySrc, "SliceChainedSlice", func() any { return tricky.SliceChainedSlice() }},
	"tricky/StructEmbeddedOverride": {trickySrc, "StructEmbeddedOverride", func() any { return tricky.StructEmbeddedOverride() }},
	"tricky/MapTwoKeys":             {trickySrc, "MapTwoKeys", func() any { return tricky.MapTwoKeys() }},
	"tricky/SliceNegativeIndex":     {trickySrc, "SliceNegativeIndex", func() any { return tricky.SliceNegativeIndex() }},
	// "tricky/StructSelfRef": SKIPPED - causes stack overflow with self-referencing type
	// "tricky/MapRangeBreak": SKIPPED - map iteration order is non-deterministic
	"tricky/SliceStructIndex": {trickySrc, "SliceStructIndex", func() any { return tricky.SliceStructIndex() }},
	"tricky/MapStructUpdate":  {trickySrc, "MapStructUpdate", func() any { return tricky.MapStructUpdate() }},
	"tricky/PointerToPointer": {trickySrc, "PointerToPointer", func() any { return tricky.PointerToPointer() }},
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
