// Package tests - divergence_hunt_test.go
//
// Divergence hunt tests: compare interpreted execution with native Go results.
// Uses //go:embed to load source from testdata/ directories, same pattern as correctness_test.go.
package tests

import (
	_ "embed"
	"reflect"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt1"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt2"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt3"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt4"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt5"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt6"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt7"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt8"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt9"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt10"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt11"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt12"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt13"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt14"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt15"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt16"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt17"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt18"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt19"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt20"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt21"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt22"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt23"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt24"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt25"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt26"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt27"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt28"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt29"
	"git.woa.com/youngjin/gig/tests/testdata/divergence_hunt30"
)

//go:embed testdata/divergence_hunt1/main.go
var divergenceHunt1Src string

//go:embed testdata/divergence_hunt2/main.go
var divergenceHunt2Src string

//go:embed testdata/divergence_hunt3/main.go
var divergenceHunt3Src string

//go:embed testdata/divergence_hunt4/main.go
var divergenceHunt4Src string

//go:embed testdata/divergence_hunt5/main.go
var divergenceHunt5Src string

//go:embed testdata/divergence_hunt6/main.go
var divergenceHunt6Src string

//go:embed testdata/divergence_hunt7/main.go
var divergenceHunt7Src string

//go:embed testdata/divergence_hunt8/main.go
var divergenceHunt8Src string

//go:embed testdata/divergence_hunt9/main.go
var divergenceHunt9Src string

//go:embed testdata/divergence_hunt10/main.go
var divergenceHunt10Src string

//go:embed testdata/divergence_hunt11/main.go
var divergenceHunt11Src string

//go:embed testdata/divergence_hunt12/main.go
var divergenceHunt12Src string

//go:embed testdata/divergence_hunt13/main.go
var divergenceHunt13Src string

//go:embed testdata/divergence_hunt14/main.go
var divergenceHunt14Src string

//go:embed testdata/divergence_hunt15/main.go
var divergenceHunt15Src string

//go:embed testdata/divergence_hunt16/main.go
var divergenceHunt16Src string

//go:embed testdata/divergence_hunt17/main.go
var divergenceHunt17Src string

//go:embed testdata/divergence_hunt18/main.go
var divergenceHunt18Src string

//go:embed testdata/divergence_hunt19/main.go
var divergenceHunt19Src string

//go:embed testdata/divergence_hunt20/main.go
var divergenceHunt20Src string

//go:embed testdata/divergence_hunt21/main.go
var divergenceHunt21Src string

//go:embed testdata/divergence_hunt22/main.go
var divergenceHunt22Src string

//go:embed testdata/divergence_hunt23/main.go
var divergenceHunt23Src string

//go:embed testdata/divergence_hunt24/main.go
var divergenceHunt24Src string

//go:embed testdata/divergence_hunt25/main.go
var divergenceHunt25Src string

//go:embed testdata/divergence_hunt26/main.go
var divergenceHunt26Src string

//go:embed testdata/divergence_hunt27/main.go
var divergenceHunt27Src string

//go:embed testdata/divergence_hunt28/main.go
var divergenceHunt28Src string

//go:embed testdata/divergence_hunt29/main.go
var divergenceHunt29Src string

//go:embed testdata/divergence_hunt30/main.go
var divergenceHunt30Src string

// divergenceTestCase is like testCase but with explicit expected value.
// This is used for divergence hunting where we compare interpreter output
// against native Go execution.
type divergenceTestCase struct {
	funcName string
	args     []any
	native   any // native function, called via reflection
}

// divergenceTestSet is a set of divergence test cases sharing one source file.
type divergenceTestSet struct {
	src       string
	tests     map[string]divergenceTestCase
	buildOpts []gig.BuildOption
}

// runDivergenceTestSet compiles the source once and runs each test,
// comparing interpreter output with native Go execution.
func runDivergenceTestSet(t *testing.T, set divergenceTestSet) {
	t.Helper()
	prog, err := gig.Build(set.src, set.buildOpts...)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	defer prog.Close()

	for name, tc := range set.tests {
		t.Run(name, func(t *testing.T) {
			// Run interpreter
			interpResult, interpErr := prog.Run(tc.funcName, tc.args...)
			if interpErr != nil {
				t.Errorf("DIVERGENCE (error): %v", interpErr)
				return
			}

			// Run native
			if tc.native == nil {
				t.Fatalf("native function is nil for %s", name)
			}
			nativeResult := callNative(tc.native, tc.args)

			// Compare
			if !reflect.DeepEqual(interpResult, nativeResult) {
				t.Errorf("DIVERGENCE (mismatch): interp=%v (%T), native=%v (%T)",
					interpResult, interpResult, nativeResult, nativeResult)
			}
		})
	}
}

func TestDivergenceHunt1(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt1Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"NilSliceCompare":    {funcName: "NilSliceCompare", native: divergence_hunt1.NilSliceCompare},
			"NilMapCompare":      {funcName: "NilMapCompare", native: divergence_hunt1.NilMapCompare},
			"NilChanCompare":     {funcName: "NilChanCompare", native: divergence_hunt1.NilChanCompare},
			"ComplexArith":       {funcName: "ComplexArith", native: divergence_hunt1.ComplexArith},
			"StringIndexByte":    {funcName: "StringIndexByte", native: divergence_hunt1.StringIndexByte},
			"IntOverflow":        {funcName: "IntOverflow", native: divergence_hunt1.IntOverflow},
			"DeferModify":        {funcName: "DeferModify", native: divergence_hunt1.DeferModify},
			"TypeAssertPanic":    {funcName: "TypeAssertPanic", native: divergence_hunt1.TypeAssertPanic},
			"Complex64Arith":     {funcName: "Complex64Arith", native: divergence_hunt1.Complex64Arith},
			"SliceBoundsPanic":   {funcName: "SliceBoundsPanic", native: divergence_hunt1.SliceBoundsPanic},
			"NilPointerDeref":    {funcName: "NilPointerDeref", native: divergence_hunt1.NilPointerDeref},
			"NilMapWrite":        {funcName: "NilMapWrite", native: divergence_hunt1.NilMapWrite},
			"DivZeroPanicTest":   {funcName: "DivZeroPanicTest", native: divergence_hunt1.DivZeroPanicTest},
			"UintOverflow":       {funcName: "UintOverflow", native: divergence_hunt1.UintOverflow},
			"Int8Negative":       {funcName: "Int8Negative", native: divergence_hunt1.Int8Negative},
			"NaNCompare":         {funcName: "NaNCompare", native: divergence_hunt1.NaNCompare},
			"MapNilLookup":       {funcName: "MapNilLookup", native: divergence_hunt1.MapNilLookup},
			"SliceCopy":          {funcName: "SliceCopy", native: divergence_hunt1.SliceCopy},
			"RuneLiteral":        {funcName: "RuneLiteral", native: divergence_hunt1.RuneLiteral},
			"NilInterfaceAssert": {funcName: "NilInterfaceAssert", native: divergence_hunt1.NilInterfaceAssert},
			"SortInts":           {funcName: "SortInts", native: divergence_hunt1.SortInts},
			"StringsJoin":        {funcName: "StringsJoin", native: divergence_hunt1.StringsJoin},
			"StringsSplit":       {funcName: "StringsSplit", native: divergence_hunt1.StringsSplit},
			"StringsContains":    {funcName: "StringsContains", native: divergence_hunt1.StringsContains},
			"StrconvRoundTrip":   {funcName: "StrconvRoundTrip", native: divergence_hunt1.StrconvRoundTrip},
			"FmtSprintf":         {funcName: "FmtSprintf", native: divergence_hunt1.FmtSprintf},
			"PanicInDefer":       {funcName: "PanicInDefer", native: divergence_hunt1.PanicInDefer},
			"MultipleRecoverCalls": {funcName: "MultipleRecoverCalls", native: divergence_hunt1.MultipleRecoverCalls},
			"BoolToStrconv":      {funcName: "BoolToStrconv", native: divergence_hunt1.BoolToStrconv},
			"FloatToStrconv":     {funcName: "FloatToStrconv", native: divergence_hunt1.FloatToStrconv},
			"StringsReplace":     {funcName: "StringsReplace", native: divergence_hunt1.StringsReplace},
			"StringsHasPrefix":   {funcName: "StringsHasPrefix", native: divergence_hunt1.StringsHasPrefix},
			"StringsTrim":        {funcName: "StringsTrim", native: divergence_hunt1.StringsTrim},
			"MapIntKey":          {funcName: "MapIntKey", native: divergence_hunt1.MapIntKey},
			"CapSlice":           {funcName: "CapSlice", native: divergence_hunt1.CapSlice},
			"ByteSliceIndex":     {funcName: "ByteSliceIndex", native: divergence_hunt1.ByteSliceIndex},
			"DeferMultipleOrder": {funcName: "DeferMultipleOrder", native: divergence_hunt1.DeferMultipleOrder},
			"ErrorTypeAssertion": {funcName: "ErrorTypeAssertion", native: divergence_hunt1.ErrorTypeAssertion},
			"RecursiveFactorial": {funcName: "RecursiveFactorial", native: divergence_hunt1.RecursiveFactorial},
			"ClosureCounter":     {funcName: "ClosureCounter", native: divergence_hunt1.ClosureCounter},
			"BitwiseAnd":         {funcName: "BitwiseAnd", native: divergence_hunt1.BitwiseAnd},
			"BitwiseOr":          {funcName: "BitwiseOr", native: divergence_hunt1.BitwiseOr},
			"BitwiseXor":         {funcName: "BitwiseXor", native: divergence_hunt1.BitwiseXor},
			"BitwiseShift":       {funcName: "BitwiseShift", native: divergence_hunt1.BitwiseShift},
			"Float64Arith":       {funcName: "Float64Arith", native: divergence_hunt1.Float64Arith},
			"PanicIntValue":      {funcName: "PanicIntValue", native: divergence_hunt1.PanicIntValue},
			"DoublePanic":        {funcName: "DoublePanic", native: divergence_hunt1.DoublePanic},
			"DeferModifyAfterPanic": {funcName: "DeferModifyAfterPanic", native: divergence_hunt1.DeferModifyAfterPanic},
			"SliceOfStructs":     {funcName: "SliceOfStructs", native: divergence_hunt1.SliceOfStructs},
			"ForBreak":           {funcName: "ForBreak", native: divergence_hunt1.ForBreak},
			"NestedLoop":         {funcName: "NestedLoop", native: divergence_hunt1.NestedLoop},
			"StringCompareOps":   {funcName: "StringCompareOps", native: divergence_hunt1.StringCompareOps},
			"MapCommaOkMissing":  {funcName: "MapCommaOkMissing", native: divergence_hunt1.MapCommaOkMissing},
			"SwitchDefault":      {funcName: "SwitchDefault", native: divergence_hunt1.SwitchDefault},
			"VariadicFunc":       {funcName: "VariadicFunc", native: divergence_hunt1.VariadicFunc},
			"TypeSwitch":         {funcName: "TypeSwitch", native: divergence_hunt1.TypeSwitch},
			"StructEmbedding":    {funcName: "StructEmbedding", native: divergence_hunt1.StructEmbedding},
			"ChannelBuffered":    {funcName: "ChannelBuffered", native: divergence_hunt1.ChannelBuffered},
		},
	})
}

func TestDivergenceHunt2(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt2Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"MapLen":                    {funcName: "MapLen", native: divergence_hunt2.MapLen},
			"MapDelete":                 {funcName: "MapDelete", native: divergence_hunt2.MapDelete},
			"MapOverwrite":              {funcName: "MapOverwrite", native: divergence_hunt2.MapOverwrite},
			"SliceNilAppend":            {funcName: "SliceNilAppend", native: divergence_hunt2.SliceNilAppend},
			"SliceGrow":                 {funcName: "SliceGrow", native: divergence_hunt2.SliceGrow},
			"StringLen":                 {funcName: "StringLen", native: divergence_hunt2.StringLen},
			"StringConcat":              {funcName: "StringConcat", native: divergence_hunt2.StringConcat},
			"IntConversion":             {funcName: "IntConversion", native: divergence_hunt2.IntConversion},
			"UintConversion":            {funcName: "UintConversion", native: divergence_hunt2.UintConversion},
			"MultiReturnSwap":           {funcName: "MultiReturnSwap", native: divergence_hunt2.MultiReturnSwap},
			"BlankIdentifier":           {funcName: "BlankIdentifier", native: divergence_hunt2.BlankIdentifier},
			"NilSliceLen":               {funcName: "NilSliceLen", native: divergence_hunt2.NilSliceLen},
			"NilMapLen":                 {funcName: "NilMapLen", native: divergence_hunt2.NilMapLen},
			"PointerDeref":              {funcName: "PointerDeref", native: divergence_hunt2.PointerDeref},
			"PointerAssign":             {funcName: "PointerAssign", native: divergence_hunt2.PointerAssign},
			"SliceOfPointers":           {funcName: "SliceOfPointers", native: divergence_hunt2.SliceOfPointers},
			"MapIteration":              {funcName: "MapIteration", native: divergence_hunt2.MapIteration},
			"StringRange":               {funcName: "StringRange", native: divergence_hunt2.StringRange},
			"FloatConversion":           {funcName: "FloatConversion", native: divergence_hunt2.FloatConversion},
			"ByteSliceAppend":           {funcName: "ByteSliceAppend", native: divergence_hunt2.ByteSliceAppend},
			"ByteSliceWrite":            {funcName: "ByteSliceWrite", native: divergence_hunt2.ByteSliceWrite},
			"StructCompare":             {funcName: "StructCompare", native: divergence_hunt2.StructCompare},
			"ArrayLen":                  {funcName: "ArrayLen", native: divergence_hunt2.ArrayLen},
			"ArrayValue":                {funcName: "ArrayValue", native: divergence_hunt2.ArrayValue},
			"StringIndexOutOfRange":     {funcName: "StringIndexOutOfRange", native: divergence_hunt2.StringIndexOutOfRange},
			"MapKeyIntFloat":            {funcName: "MapKeyIntFloat", native: divergence_hunt2.MapKeyIntFloat},
			"ShortVarDecl":              {funcName: "ShortVarDecl", native: divergence_hunt2.ShortVarDecl},
			"MultipleShortVar":          {funcName: "MultipleShortVar", native: divergence_hunt2.MultipleShortVar},
			"SliceThreeIndex":           {funcName: "SliceThreeIndex", native: divergence_hunt2.SliceThreeIndex},
			"NilFuncCall":               {funcName: "NilFuncCall", native: divergence_hunt2.NilFuncCall},
			"StringByteSliceConversion": {funcName: "StringByteSliceConversion", native: divergence_hunt2.StringByteSliceConversion},
		},
	})
}

func TestDivergenceHunt3(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt3Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"StringBuilder":     {funcName: "StringBuilder", native: divergence_hunt3.StringBuilder},
			"ConstBlock":        {funcName: "ConstBlock", native: divergence_hunt3.ConstBlock},
			"IotaEnum":          {funcName: "IotaEnum", native: divergence_hunt3.IotaEnum},
			"MultipleAssign":    {funcName: "MultipleAssign", native: divergence_hunt3.MultipleAssign},
			"NestedMap":         {funcName: "NestedMap", native: divergence_hunt3.NestedMap},
			"RuneIteration":     {funcName: "RuneIteration", native: divergence_hunt3.RuneIteration},
			"StringIndexRune":   {funcName: "StringIndexRune", native: divergence_hunt3.StringIndexRune},
			"StringCount":       {funcName: "StringCount", native: divergence_hunt3.StringCount},
			"MapBoolKey":        {funcName: "MapBoolKey", native: divergence_hunt3.MapBoolKey},
			"SliceReverse":      {funcName: "SliceReverse", native: divergence_hunt3.SliceReverse},
			"StructMethod":      {funcName: "StructMethod", native: divergence_hunt3.StructMethod},
			"InterfaceEmpty":    {funcName: "InterfaceEmpty", native: divergence_hunt3.InterfaceEmpty},
			"InterfaceNil":      {funcName: "InterfaceNil", native: divergence_hunt3.InterfaceNil},
			"SliceOfInterface":  {funcName: "SliceOfInterface", native: divergence_hunt3.SliceOfInterface},
			"MapWithStructValue": {funcName: "MapWithStructValue", native: divergence_hunt3.MapWithStructValue},
			"StringFields":      {funcName: "StringFields", native: divergence_hunt3.StringFields},
			"StringRepeat":      {funcName: "StringRepeat", native: divergence_hunt3.StringRepeat},
			"StringMap":         {funcName: "StringMap", native: divergence_hunt3.StringMap},
			"MapStructKey":      {funcName: "MapStructKey", native: divergence_hunt3.MapStructKey},
			"SliceMinMax":       {funcName: "SliceMinMax", native: divergence_hunt3.SliceMinMax},
			"NestedIf":          {funcName: "NestedIf", native: divergence_hunt3.NestedIf},
			"StringToLower":     {funcName: "StringToLower", native: divergence_hunt3.StringToLower},
			"StringToUpper":     {funcName: "StringToUpper", native: divergence_hunt3.StringToUpper},
			"ContinueLoop":      {funcName: "ContinueLoop", native: divergence_hunt3.ContinueLoop},
			"LabeledBreak":      {funcName: "LabeledBreak", native: divergence_hunt3.LabeledBreak},
			"SliceMakeZero":     {funcName: "SliceMakeZero", native: divergence_hunt3.SliceMakeZero},
			"ArrayIteration":    {funcName: "ArrayIteration", native: divergence_hunt3.ArrayIteration},
			"Float32Arith":      {funcName: "Float32Arith", native: divergence_hunt3.Float32Arith},
			"Int8Arith":         {funcName: "Int8Arith", native: divergence_hunt3.Int8Arith},
			"Uint16Arith":       {funcName: "Uint16Arith", native: divergence_hunt3.Uint16Arith},
		},
	})
}

func TestDivergenceHunt4(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt4Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"Float64NaN":         {funcName: "Float64NaN", native: divergence_hunt4.Float64NaN},
			"Float64Inf":         {funcName: "Float64Inf", native: divergence_hunt4.Float64Inf},
			"Float64NegZero":     {funcName: "Float64NegZero", native: divergence_hunt4.Float64NegZero},
			"Int16Conversion":    {funcName: "Int16Conversion", native: divergence_hunt4.Int16Conversion},
			"Uint32Conversion":   {funcName: "Uint32Conversion", native: divergence_hunt4.Uint32Conversion},
			"FloatToIntTruncation": {funcName: "FloatToIntTruncation", native: divergence_hunt4.FloatToIntTruncation},
			"NegativeFloatToInt": {funcName: "NegativeFloatToInt", native: divergence_hunt4.NegativeFloatToInt},
			"StrconvAtoi":        {funcName: "StrconvAtoi", native: divergence_hunt4.StrconvAtoi},
			"StrconvItoa":        {funcName: "StrconvItoa", native: divergence_hunt4.StrconvItoa},
			"StrconvFormatInt":   {funcName: "StrconvFormatInt", native: divergence_hunt4.StrconvFormatInt},
			"StrconvParseFloat":  {funcName: "StrconvParseFloat", native: divergence_hunt4.StrconvParseFloat},
			"MathAbs":            {funcName: "MathAbs", native: divergence_hunt4.MathAbs},
			"MathMax":            {funcName: "MathMax", native: divergence_hunt4.MathMax},
			"MathMin":            {funcName: "MathMin", native: divergence_hunt4.MathMin},
			"MathPow":            {funcName: "MathPow", native: divergence_hunt4.MathPow},
			"MathSqrt":           {funcName: "MathSqrt", native: divergence_hunt4.MathSqrt},
			"MathCeil":           {funcName: "MathCeil", native: divergence_hunt4.MathCeil},
			"MathFloor":          {funcName: "MathFloor", native: divergence_hunt4.MathFloor},
			"IntMin":             {funcName: "IntMin", native: divergence_hunt4.IntMin},
			"IntMax":             {funcName: "IntMax", native: divergence_hunt4.IntMax},
			"UintptrSize":        {funcName: "UintptrSize", native: divergence_hunt4.UintptrSize},
			"ByteArith":          {funcName: "ByteArith", native: divergence_hunt4.ByteArith},
			"Int32Overflow":      {funcName: "Int32Overflow", native: divergence_hunt4.Int32Overflow},
			"Uint8Wrap":          {funcName: "Uint8Wrap", native: divergence_hunt4.Uint8Wrap},
			"ComplexConj":        {funcName: "ComplexConj", native: divergence_hunt4.ComplexConj},
			"Float32Precision":   {funcName: "Float32Precision", native: divergence_hunt4.Float32Precision},
			"MapLenAfterDelete":  {funcName: "MapLenAfterDelete", native: divergence_hunt4.MapLenAfterDelete},
			"SliceCapAfterAppend": {funcName: "SliceCapAfterAppend", native: divergence_hunt4.SliceCapAfterAppend},
			"StringFromRunes":    {funcName: "StringFromRunes", native: divergence_hunt4.StringFromRunes},
			"RuneToInt":          {funcName: "RuneToInt", native: divergence_hunt4.RuneToInt},
			"BoolToInt":          {funcName: "BoolToInt", native: divergence_hunt4.BoolToInt},
		},
	})
}

func TestDivergenceHunt5(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt5Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"ErrorNew":            {funcName: "ErrorNew", native: divergence_hunt5.ErrorNew},
			"Errorf":              {funcName: "Errorf", native: divergence_hunt5.Errorf},
			"FmtPrintln":          {funcName: "FmtPrintln", native: divergence_hunt5.FmtPrintln},
			"FmtIntWidth":         {funcName: "FmtIntWidth", native: divergence_hunt5.FmtIntWidth},
			"FmtFloat":            {funcName: "FmtFloat", native: divergence_hunt5.FmtFloat},
			"FmtBool":             {funcName: "FmtBool", native: divergence_hunt5.FmtBool},
			"FmtHex":              {funcName: "FmtHex", native: divergence_hunt5.FmtHex},
			"FmtOctal":            {funcName: "FmtOctal", native: divergence_hunt5.FmtOctal},
			"FmtBinary":           {funcName: "FmtBinary", native: divergence_hunt5.FmtBinary},
			"FmtChar":             {funcName: "FmtChar", native: divergence_hunt5.FmtChar},
			"FmtStringWidth":      {funcName: "FmtStringWidth", native: divergence_hunt5.FmtStringWidth},
			"SliceFilter":         {funcName: "SliceFilter", native: divergence_hunt5.SliceFilter},
			"SliceMap":            {funcName: "SliceMap", native: divergence_hunt5.SliceMap},
			"ClosureSum":          {funcName: "ClosureSum", native: divergence_hunt5.ClosureSum},
			"ClosureCapture":      {funcName: "ClosureCapture", native: divergence_hunt5.ClosureCapture},
			"InterfaceSlice":      {funcName: "InterfaceSlice", native: divergence_hunt5.InterfaceSlice},
			"MultipleReturnIgnore": {funcName: "MultipleReturnIgnore", native: divergence_hunt5.MultipleReturnIgnore},
			"NamedReturn":         {funcName: "NamedReturn", native: divergence_hunt5.NamedReturn},
			"NamedReturnBare":     {funcName: "NamedReturnBare", native: divergence_hunt5.NamedReturnBare},
			"StringJoinInts":      {funcName: "StringJoinInts", native: divergence_hunt5.StringJoinInts},
			"MapStringSlice":      {funcName: "MapStringSlice", native: divergence_hunt5.MapStringSlice},
			"NestedStruct":        {funcName: "NestedStruct", native: divergence_hunt5.NestedStruct},
			"StructLiteral":       {funcName: "StructLiteral", native: divergence_hunt5.StructLiteral},
			"StructPointer":       {funcName: "StructPointer", native: divergence_hunt5.StructPointer},
			"DeferReturn":         {funcName: "DeferReturn", native: divergence_hunt5.DeferReturn},
			"DeferClosure":        {funcName: "DeferClosure", native: divergence_hunt5.DeferClosure},
			"StringEqual":         {funcName: "StringEqual", native: divergence_hunt5.StringEqual},
			"MapLookup":           {funcName: "MapLookup", native: divergence_hunt5.MapLookup},
		},
	})
}

func TestDivergenceHunt6(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt6Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"ChannelClose":           {funcName: "ChannelClose", native: divergence_hunt6.ChannelClose},
			"ChannelSelect":          {funcName: "ChannelSelect", native: divergence_hunt6.ChannelSelect},
			"ChannelNilBlock":        {funcName: "ChannelNilBlock", native: divergence_hunt6.ChannelNilBlock},
			"FuncAsValue":            {funcName: "FuncAsValue", native: divergence_hunt6.FuncAsValue},
			"HigherOrderFunc":        {funcName: "HigherOrderFunc", native: divergence_hunt6.HigherOrderFunc},
			"ClosureOverLoop":        {funcName: "ClosureOverLoop", native: divergence_hunt6.ClosureOverLoop},
			"RecursiveFib":           {funcName: "RecursiveFib", native: divergence_hunt6.RecursiveFib},
			"PartialApplication":     {funcName: "PartialApplication", native: divergence_hunt6.PartialApplication},
			"FunctionSlice":          {funcName: "FunctionSlice", native: divergence_hunt6.FunctionSlice},
			"MapFunc":                {funcName: "MapFunc", native: divergence_hunt6.MapFunc},
			"ChannelBufferLen":       {funcName: "ChannelBufferLen", native: divergence_hunt6.ChannelBufferLen},
			"ChannelBufferCap":       {funcName: "ChannelBufferCap", native: divergence_hunt6.ChannelBufferCap},
			"SelectDefault":          {funcName: "SelectDefault", native: divergence_hunt6.SelectDefault},
			"MultiReturnFunc":        {funcName: "MultiReturnFunc", native: divergence_hunt6.MultiReturnFunc},
			"NestedClosure":          {funcName: "NestedClosure", native: divergence_hunt6.NestedClosure},
			"ClosureReturnFunc":      {funcName: "ClosureReturnFunc", native: divergence_hunt6.ClosureReturnFunc},
			"ChannelReceiveOnClosed": {funcName: "ChannelReceiveOnClosed", native: divergence_hunt6.ChannelReceiveOnClosed},
			"FuncTypeDeclaration":    {funcName: "FuncTypeDeclaration", native: divergence_hunt6.FuncTypeDeclaration},
			"VariadicSpread":         {funcName: "VariadicSpread", native: divergence_hunt6.VariadicSpread},
			"InterfaceMethod":        {funcName: "InterfaceMethod", native: divergence_hunt6.InterfaceMethod},
			"StringConversion":       {funcName: "StringConversion", native: divergence_hunt6.StringConversion},
		},
	})
}

func TestDivergenceHunt7(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt7Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"SortInts":            {funcName: "SortInts", native: divergence_hunt7.SortInts},
			"SortStrings":         {funcName: "SortStrings", native: divergence_hunt7.SortStrings},
			"SortFloat64s":        {funcName: "SortFloat64s", native: divergence_hunt7.SortFloat64s},
			"SliceDelete":         {funcName: "SliceDelete", native: divergence_hunt7.SliceDelete},
			"SliceInsert":         {funcName: "SliceInsert", native: divergence_hunt7.SliceInsert},
			"SliceContains":       {funcName: "SliceContains", native: divergence_hunt7.SliceContains},
			"MapKeys":             {funcName: "MapKeys", native: divergence_hunt7.MapKeys},
			"MapValues":           {funcName: "MapValues", native: divergence_hunt7.MapValues},
			"StructWithMethods":   {funcName: "StructWithMethods", native: divergence_hunt7.StructWithMethods},
			"PointerReceiverMethod": {funcName: "PointerReceiverMethod", native: divergence_hunt7.PointerReceiverMethod},
			"TypeAssertion":       {funcName: "TypeAssertion", native: divergence_hunt7.TypeAssertion},
			"TypeAssertionString": {funcName: "TypeAssertionString", native: divergence_hunt7.TypeAssertionString},
			"TypeAssertionFail":   {funcName: "TypeAssertionFail", native: divergence_hunt7.TypeAssertionFail},
			"InterfaceTypeSwitch": {funcName: "InterfaceTypeSwitch", native: divergence_hunt7.InterfaceTypeSwitch},
			"SliceDedupe":         {funcName: "SliceDedupe", native: divergence_hunt7.SliceDedupe},
			"MapMerge":            {funcName: "MapMerge", native: divergence_hunt7.MapMerge},
			"StructSliceSort":     {funcName: "StructSliceSort", native: divergence_hunt7.StructSliceSort},
			"MapInvert":           {funcName: "MapInvert", native: divergence_hunt7.MapInvert},
			"NestedInterface":     {funcName: "NestedInterface", native: divergence_hunt7.NestedInterface},
			"SliceFlatten":        {funcName: "SliceFlatten", native: divergence_hunt7.SliceFlatten},
			"IntSliceSortCustom":  {funcName: "IntSliceSortCustom", native: divergence_hunt7.IntSliceSortCustom},
			"MapCountValues":      {funcName: "MapCountValues", native: divergence_hunt7.MapCountValues},
		},
	})
}

func TestDivergenceHunt8(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt8Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"MutexBasic":            {funcName: "MutexBasic", native: divergence_hunt8.MutexBasic},
			"OnceBasic":             {funcName: "OnceBasic", native: divergence_hunt8.OnceBasic},
			"SliceOfSlice":          {funcName: "SliceOfSlice", native: divergence_hunt8.SliceOfSlice},
			"MapOfMap":              {funcName: "MapOfMap", native: divergence_hunt8.MapOfMap},
			"StructWithSlice":       {funcName: "StructWithSlice", native: divergence_hunt8.StructWithSlice},
			"StructWithMap":         {funcName: "StructWithMap", native: divergence_hunt8.StructWithMap},
			"NestedSliceAppend":     {funcName: "NestedSliceAppend", native: divergence_hunt8.NestedSliceAppend},
			"DeepStruct":            {funcName: "DeepStruct", native: divergence_hunt8.DeepStruct},
			"SliceOfStructAppend":   {funcName: "SliceOfStructAppend", native: divergence_hunt8.SliceOfStructAppend},
			"MapWithSliceValue":     {funcName: "MapWithSliceValue", native: divergence_hunt8.MapWithSliceValue},
			"MutexInDefer":          {funcName: "MutexInDefer", native: divergence_hunt8.MutexInDefer},
			"RWMutexBasic":          {funcName: "RWMutexBasic", native: divergence_hunt8.RWMutexBasic},
			"StructWithFunc":        {funcName: "StructWithFunc", native: divergence_hunt8.StructWithFunc},
			"StructWithPointer":     {funcName: "StructWithPointer", native: divergence_hunt8.StructWithPointer},
			"SliceGrowPattern":      {funcName: "SliceGrowPattern", native: divergence_hunt8.SliceGrowPattern},
			"MapGrowPattern":        {funcName: "MapGrowPattern", native: divergence_hunt8.MapGrowPattern},
			"CompositeLiteralNested": {funcName: "CompositeLiteralNested", native: divergence_hunt8.CompositeLiteralNested},
		},
	})
}

func TestDivergenceHunt9(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt9Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"JSONMarshal":       {funcName: "JSONMarshal", native: divergence_hunt9.JSONMarshal},
			"JSONUnmarshal":     {funcName: "JSONUnmarshal", native: divergence_hunt9.JSONUnmarshal},
			"JSONMarshalMap":    {funcName: "JSONMarshalMap", native: divergence_hunt9.JSONMarshalMap},
			"RegexMatch":        {funcName: "RegexMatch", native: divergence_hunt9.RegexMatch},
			"RegexFind":         {funcName: "RegexFind", native: divergence_hunt9.RegexFind},
			"RegexFindAll":      {funcName: "RegexFindAll", native: divergence_hunt9.RegexFindAll},
			"RegexReplace":      {funcName: "RegexReplace", native: divergence_hunt9.RegexReplace},
			"MathMod":           {funcName: "MathMod", native: divergence_hunt9.MathMod},
			"MathLog":           {funcName: "MathLog", native: divergence_hunt9.MathLog},
			"MathExp":           {funcName: "MathExp", native: divergence_hunt9.MathExp},
			"MathRound":         {funcName: "MathRound", native: divergence_hunt9.MathRound},
			"MathTrunc":         {funcName: "MathTrunc", native: divergence_hunt9.MathTrunc},
			"MathRemainder":     {funcName: "MathRemainder", native: divergence_hunt9.MathRemainder},
			"MathCopysign":      {funcName: "MathCopysign", native: divergence_hunt9.MathCopysign},
			"JSONMarshalSlice":  {funcName: "JSONMarshalSlice", native: divergence_hunt9.JSONMarshalSlice},
			"JSONUnmarshalSlice": {funcName: "JSONUnmarshalSlice", native: divergence_hunt9.JSONUnmarshalSlice},
			"RegexSplit":        {funcName: "RegexSplit", native: divergence_hunt9.RegexSplit},
			"RegexSubmatch":     {funcName: "RegexSubmatch", native: divergence_hunt9.RegexSubmatch},
			"MathHypot":         {funcName: "MathHypot", native: divergence_hunt9.MathHypot},
			"MathPow10":         {funcName: "MathPow10", native: divergence_hunt9.MathPow10},
			"MathSignbit":       {funcName: "MathSignbit", native: divergence_hunt9.MathSignbit},
		},
	})
}

func TestDivergenceHunt10(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt10Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"BinarySearch":       {funcName: "BinarySearch", native: divergence_hunt10.BinarySearch},
			"StackPattern":       {funcName: "StackPattern", native: divergence_hunt10.StackPattern},
			"QueuePattern":       {funcName: "QueuePattern", native: divergence_hunt10.QueuePattern},
			"TwoSum":             {funcName: "TwoSum", native: divergence_hunt10.TwoSum},
			"IsPalindrome":       {funcName: "IsPalindrome", native: divergence_hunt10.IsPalindrome},
			"FizzBuzz":           {funcName: "FizzBuzz", native: divergence_hunt10.FizzBuzz},
			"FmtVerb":            {funcName: "FmtVerb", native: divergence_hunt10.FmtVerb},
			"FmtWidthPrecision":  {funcName: "FmtWidthPrecision", native: divergence_hunt10.FmtWidthPrecision},
			"NestedMapLookup":    {funcName: "NestedMapLookup", native: divergence_hunt10.NestedMapLookup},
			"StructSliceFilter":  {funcName: "StructSliceFilter", native: divergence_hunt10.StructSliceFilter},
			"GCD":                {funcName: "GCD", native: divergence_hunt10.GCD},
			"LCM":                {funcName: "LCM", native: divergence_hunt10.LCM},
			"Power":              {funcName: "Power", native: divergence_hunt10.Power},
			"CountDigits":        {funcName: "CountDigits", native: divergence_hunt10.CountDigits},
			"ReverseInt":         {funcName: "ReverseInt", native: divergence_hunt10.ReverseInt},
			"FibIterative":       {funcName: "FibIterative", native: divergence_hunt10.FibIterative},
			"PrimeCheck":         {funcName: "PrimeCheck", native: divergence_hunt10.PrimeCheck},
			"FactorialIterative": {funcName: "FactorialIterative", native: divergence_hunt10.FactorialIterative},
			"CountingSort":       {funcName: "CountingSort", native: divergence_hunt10.CountingSort},
			"PrefixSum":          {funcName: "PrefixSum", native: divergence_hunt10.PrefixSum},
			"StringAnagram":      {funcName: "StringAnagram", native: divergence_hunt10.StringAnagram},
		},
	})
}

func TestDivergenceHunt11(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt11Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"DeferInLoop":          {funcName: "DeferInLoop", native: divergence_hunt11.DeferInLoop},
			"DeferAndPanicOrder":   {funcName: "DeferAndPanicOrder", native: divergence_hunt11.DeferAndPanicOrder},
			"RecoverInFunction":    {funcName: "RecoverInFunction", native: divergence_hunt11.RecoverInFunction},
			"PanicWithStruct":      {funcName: "PanicWithStruct", native: divergence_hunt11.PanicWithStruct},
			"NamedReturnWithDefer": {funcName: "NamedReturnWithDefer", native: divergence_hunt11.NamedReturnWithDefer},
			"MultipleDeferModify":  {funcName: "MultipleDeferModify", native: divergence_hunt11.MultipleDeferModify},
			"DeferWithArgument":    {funcName: "DeferWithArgument", native: divergence_hunt11.DeferWithArgument},
			"PanicNilValue":        {funcName: "PanicNilValue", native: divergence_hunt11.PanicNilValue},
			"ClosureReturnFunc":    {funcName: "ClosureReturnFunc", native: divergence_hunt11.ClosureReturnFunc},
			"FmtSprintfMulti":      {funcName: "FmtSprintfMulti", native: divergence_hunt11.FmtSprintfMulti},
			"FmtErrorf":            {funcName: "FmtErrorf", native: divergence_hunt11.FmtErrorf},
			"NestedDeferRecover":   {funcName: "NestedDeferRecover", native: divergence_hunt11.NestedDeferRecover},
			"DeferWithMethod":      {funcName: "DeferWithMethod", native: divergence_hunt11.DeferWithMethod},
			"ClosureCaptureSlice":  {funcName: "ClosureCaptureSlice", native: divergence_hunt11.ClosureCaptureSlice},
			"ClosureCaptureMap":    {funcName: "ClosureCaptureMap", native: divergence_hunt11.ClosureCaptureMap},
			"MultiplePanicRecover": {funcName: "MultiplePanicRecover", native: divergence_hunt11.MultiplePanicRecover},
			"DeferRecoverReturnsValue": {funcName: "DeferRecoverReturnsValue", native: divergence_hunt11.DeferRecoverReturnsValue},
			"SliceAppendInClosure": {funcName: "SliceAppendInClosure", native: divergence_hunt11.SliceAppendInClosure},
			"MapWriteInClosure":    {funcName: "MapWriteInClosure", native: divergence_hunt11.MapWriteInClosure},
			"DeferChain":           {funcName: "DeferChain", native: divergence_hunt11.DeferChain},
			"RecoverReturnsNilAfter": {funcName: "RecoverReturnsNilAfter", native: divergence_hunt11.RecoverReturnsNilAfter},
		},
	})
}

func TestDivergenceHunt12(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt12Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"JSONNestedStruct":     {funcName: "JSONNestedStruct", native: divergence_hunt12.JSONNestedStruct},
			"JSONSliceOfStructs":   {funcName: "JSONSliceOfStructs", native: divergence_hunt12.JSONSliceOfStructs},
			"JSONUnmarshalIntoMap": {funcName: "JSONUnmarshalIntoMap", native: divergence_hunt12.JSONUnmarshalIntoMap},
			"StringTitle":          {funcName: "StringTitle", native: divergence_hunt12.StringTitle},
			"StringEqualFold":      {funcName: "StringEqualFold", native: divergence_hunt12.StringEqualFold},
			"StringIndex":          {funcName: "StringIndex", native: divergence_hunt12.StringIndex},
			"StringLastIndex":      {funcName: "StringLastIndex", native: divergence_hunt12.StringLastIndex},
			"StringIndexAny":       {funcName: "StringIndexAny", native: divergence_hunt12.StringIndexAny},
			"StringNewReplacer":    {funcName: "StringNewReplacer", native: divergence_hunt12.StringNewReplacer},
			"StringBuilderGrow":    {funcName: "StringBuilderGrow", native: divergence_hunt12.StringBuilderGrow},
			"SortSliceStable":      {funcName: "SortSliceStable", native: divergence_hunt12.SortSliceStable},
			"SortSearch":           {funcName: "SortSearch", native: divergence_hunt12.SortSearch},
			"FmtSprintfBoolean":    {funcName: "FmtSprintfBoolean", native: divergence_hunt12.FmtSprintfBoolean},
			"FmtSprintfFloat":      {funcName: "FmtSprintfFloat", native: divergence_hunt12.FmtSprintfFloat},
			"FmtSprintfInt":        {funcName: "FmtSprintfInt", native: divergence_hunt12.FmtSprintfInt},
			"FmtSprintfString":     {funcName: "FmtSprintfString", native: divergence_hunt12.FmtSprintfString},
			"JSONMarshalBool":      {funcName: "JSONMarshalBool", native: divergence_hunt12.JSONMarshalBool},
			"JSONUnmarshalBool":    {funcName: "JSONUnmarshalBool", native: divergence_hunt12.JSONUnmarshalBool},
			"JSONMarshalNil":       {funcName: "JSONMarshalNil", native: divergence_hunt12.JSONMarshalNil},
			"SliceMinMaxInt":       {funcName: "SliceMinMaxInt", native: divergence_hunt12.SliceMinMaxInt},
			"StringCountSubstring": {funcName: "StringCountSubstring", native: divergence_hunt12.StringCountSubstring},
			"MapHasKey":            {funcName: "MapHasKey", native: divergence_hunt12.MapHasKey},
		},
	})
}

func TestDivergenceHunt13(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt13Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"StructZeroValue":              {funcName: "StructZeroValue", native: divergence_hunt13.StructZeroValue},
			"StructPointerNil":             {funcName: "StructPointerNil", native: divergence_hunt13.StructPointerNil},
			"StructCopyOnAssign":           {funcName: "StructCopyOnAssign", native: divergence_hunt13.StructCopyOnAssign},
			"StructFieldAccess":            {funcName: "StructFieldAccess", native: divergence_hunt13.StructFieldAccess},
			"InterfaceNilComparison":       {funcName: "InterfaceNilComparison", native: divergence_hunt13.InterfaceNilComparison},
			"InterfaceTypedNil":            {funcName: "InterfaceTypedNil", native: divergence_hunt13.InterfaceTypedNil},
			"TypeAssertionWithBool":        {funcName: "TypeAssertionWithBool", native: divergence_hunt13.TypeAssertionWithBool},
			"MultipleTypeAssertions":       {funcName: "MultipleTypeAssertions", native: divergence_hunt13.MultipleTypeAssertions},
			"PointerToStruct":              {funcName: "PointerToStruct", native: divergence_hunt13.PointerToStruct},
			"PointerToStructModify":        {funcName: "PointerToStructModify", native: divergence_hunt13.PointerToStructModify},
			"StructAsMapValue":             {funcName: "StructAsMapValue", native: divergence_hunt13.StructAsMapValue},
			"StructInSlice":                {funcName: "StructInSlice", native: divergence_hunt13.StructInSlice},
			"IntTypeAlias":                 {funcName: "IntTypeAlias", native: divergence_hunt13.IntTypeAlias},
			"StringTypeAlias":              {funcName: "StringTypeAlias", native: divergence_hunt13.StringTypeAlias},
			"SliceOfAlias":                 {funcName: "SliceOfAlias", native: divergence_hunt13.SliceOfAlias},
			"NestedTypeDefinitions":        {funcName: "NestedTypeDefinitions", native: divergence_hunt13.NestedTypeDefinitions},
			"FmtStruct":                    {funcName: "FmtStruct", native: divergence_hunt13.FmtStruct},
			"FmtPointer":                   {funcName: "FmtPointer", native: divergence_hunt13.FmtPointer},
			"ConversionBetweenNumericTypes": {funcName: "ConversionBetweenNumericTypes", native: divergence_hunt13.ConversionBetweenNumericTypes},
			"UnsignedToSigned":             {funcName: "UnsignedToSigned", native: divergence_hunt13.UnsignedToSigned},
		},
	})
}

func TestDivergenceHunt14(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt14Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"FloatAddPrecision":       {funcName: "FloatAddPrecision", native: divergence_hunt14.FloatAddPrecision},
			"FloatMultiplyPrecision":  {funcName: "FloatMultiplyPrecision", native: divergence_hunt14.FloatMultiplyPrecision},
			"FloatDivPrecision":       {funcName: "FloatDivPrecision", native: divergence_hunt14.FloatDivPrecision},
			"FloatNegative":           {funcName: "FloatNegative", native: divergence_hunt14.FloatNegative},
			"FloatZeroDivision":       {funcName: "FloatZeroDivision", native: divergence_hunt14.FloatZeroDivision},
			"FloatNaNArithmetic":      {funcName: "FloatNaNArithmetic", native: divergence_hunt14.FloatNaNArithmetic},
			"FloatInfArithmetic":      {funcName: "FloatInfArithmetic", native: divergence_hunt14.FloatInfArithmetic},
			"FloatComparisonPrecision": {funcName: "FloatComparisonPrecision", native: divergence_hunt14.FloatComparisonPrecision},
			"IntDivisionTruncation":   {funcName: "IntDivisionTruncation", native: divergence_hunt14.IntDivisionTruncation},
			"IntModulo":               {funcName: "IntModulo", native: divergence_hunt14.IntModulo},
			"NegativeModulo":          {funcName: "NegativeModulo", native: divergence_hunt14.NegativeModulo},
			"Float32NaN":              {funcName: "Float32NaN", native: divergence_hunt14.Float32NaN},
			"Float32Inf":              {funcName: "Float32Inf", native: divergence_hunt14.Float32Inf},
			"MathSin":                 {funcName: "MathSin", native: divergence_hunt14.MathSin},
			"MathCos":                 {funcName: "MathCos", native: divergence_hunt14.MathCos},
			"MathTan":                 {funcName: "MathTan", native: divergence_hunt14.MathTan},
			"MathAtan2":               {funcName: "MathAtan2", native: divergence_hunt14.MathAtan2},
			"MathLog2":                {funcName: "MathLog2", native: divergence_hunt14.MathLog2},
			"MathLog10":               {funcName: "MathLog10", native: divergence_hunt14.MathLog10},
			"FmtFloatFormat":          {funcName: "FmtFloatFormat", native: divergence_hunt14.FmtFloatFormat},
			"FmtIntFormat":            {funcName: "FmtIntFormat", native: divergence_hunt14.FmtIntFormat},
			"FloatMaxMin":             {funcName: "FloatMaxMin", native: divergence_hunt14.FloatMaxMin},
			"Float32Limits":           {funcName: "Float32Limits", native: divergence_hunt14.Float32Limits},
			"ComplexMagnitude":        {funcName: "ComplexMagnitude", native: divergence_hunt14.ComplexMagnitude},
		},
	})
}

func TestDivergenceHunt15(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt15Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"WordCount":              {funcName: "WordCount", native: divergence_hunt15.WordCount},
			"TopKElements":           {funcName: "TopKElements", native: divergence_hunt15.TopKElements},
			"FlattenAndSum":          {funcName: "FlattenAndSum", native: divergence_hunt15.FlattenAndSum},
			"FrequencyCount":         {funcName: "FrequencyCount", native: divergence_hunt15.FrequencyCount},
			"ReverseString":          {funcName: "ReverseString", native: divergence_hunt15.ReverseString},
			"StringPermutationCheck": {funcName: "StringPermutationCheck", native: divergence_hunt15.StringPermutationCheck},
			"MatrixSum":              {funcName: "MatrixSum", native: divergence_hunt15.MatrixSum},
			"MatrixTranspose":        {funcName: "MatrixTranspose", native: divergence_hunt15.MatrixTranspose},
			"JSONEncodeDecode":       {funcName: "JSONEncodeDecode", native: divergence_hunt15.JSONEncodeDecode},
			"StringCompression":      {funcName: "StringCompression", native: divergence_hunt15.StringCompression},
			"UniqueElements":         {funcName: "UniqueElements", native: divergence_hunt15.UniqueElements},
			"IntersectSlices":        {funcName: "IntersectSlices", native: divergence_hunt15.IntersectSlices},
			"MergeSortedSlices":      {funcName: "MergeSortedSlices", native: divergence_hunt15.MergeSortedSlices},
			"MovingAverage":          {funcName: "MovingAverage", native: divergence_hunt15.MovingAverage},
			"SpiralMatrix":           {funcName: "SpiralMatrix", native: divergence_hunt15.SpiralMatrix},
			"FmtStructFormatting":    {funcName: "FmtStructFormatting", native: divergence_hunt15.FmtStructFormatting},
		},
	})
}

func TestDivergenceHunt16(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt16Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"SwitchNoCase":             {funcName: "SwitchNoCase", native: divergence_hunt16.SwitchNoCase},
			"SwitchMultipleCases":      {funcName: "SwitchMultipleCases", native: divergence_hunt16.SwitchMultipleCases},
			"SwitchWithInit":           {funcName: "SwitchWithInit", native: divergence_hunt16.SwitchWithInit},
			"NestedSwitch":             {funcName: "NestedSwitch", native: divergence_hunt16.NestedSwitch},
			"ForRangeWithIndex":        {funcName: "ForRangeWithIndex", native: divergence_hunt16.ForRangeWithIndex},
			"ForRangeWithValue":        {funcName: "ForRangeWithValue", native: divergence_hunt16.ForRangeWithValue},
			"ForRangeMap":              {funcName: "ForRangeMap", native: divergence_hunt16.ForRangeMap},
			"IfElseChain":              {funcName: "IfElseChain", native: divergence_hunt16.IfElseChain},
			"NestedIfElse":             {funcName: "NestedIfElse", native: divergence_hunt16.NestedIfElse},
			"InfiniteLoopBreak":        {funcName: "InfiniteLoopBreak", native: divergence_hunt16.InfiniteLoopBreak},
			"ForLoopContinue":          {funcName: "ForLoopContinue", native: divergence_hunt16.ForLoopContinue},
			"LoopWithMultipleBreaks":   {funcName: "LoopWithMultipleBreaks", native: divergence_hunt16.LoopWithMultipleBreaks},
			"SwitchExpression":         {funcName: "SwitchExpression", native: divergence_hunt16.SwitchExpression},
			"ForRangeString":           {funcName: "ForRangeString", native: divergence_hunt16.ForRangeString},
			"ForRangeEmptySlice":       {funcName: "ForRangeEmptySlice", native: divergence_hunt16.ForRangeEmptySlice},
			"DoubleLoop":               {funcName: "DoubleLoop", native: divergence_hunt16.DoubleLoop},
			"LoopAccumulator":          {funcName: "LoopAccumulator", native: divergence_hunt16.LoopAccumulator},
			"SwitchFallthroughSimulated": {funcName: "SwitchFallthroughSimulated", native: divergence_hunt16.SwitchFallthroughSimulated},
			"EarlyReturn":              {funcName: "EarlyReturn", native: divergence_hunt16.EarlyReturn},
			"LoopWithEarlyReturn":       {funcName: "LoopWithEarlyReturn", native: divergence_hunt16.LoopWithEarlyReturn},
		},
	})
}

func TestDivergenceHunt17(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt17Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"InterfaceComposition":    {funcName: "InterfaceComposition", native: divergence_hunt17.InterfaceComposition},
			"InterfaceEmpty":          {funcName: "InterfaceEmpty", native: divergence_hunt17.InterfaceEmpty},
			"InterfaceSlice":          {funcName: "InterfaceSlice", native: divergence_hunt17.InterfaceSlice},
			"InterfaceMap":            {funcName: "InterfaceMap", native: divergence_hunt17.InterfaceMap},
			"StructMethodOnPointer":   {funcName: "StructMethodOnPointer", native: divergence_hunt17.StructMethodOnPointer},
			"StructMethodOnValue":     {funcName: "StructMethodOnValue", native: divergence_hunt17.StructMethodOnValue},
			"MethodChain":             {funcName: "MethodChain", native: divergence_hunt17.MethodChain},
			"PolymorphismPattern":     {funcName: "PolymorphismPattern", native: divergence_hunt17.PolymorphismPattern},
			"NullableInterface":       {funcName: "NullableInterface", native: divergence_hunt17.NullableInterface},
			"InterfaceTypeAssertion":  {funcName: "InterfaceTypeAssertion", native: divergence_hunt17.InterfaceTypeAssertion},
			"EmbeddedStructAccess":    {funcName: "EmbeddedStructAccess", native: divergence_hunt17.EmbeddedStructAccess},
			"NestedStructAccess":      {funcName: "NestedStructAccess", native: divergence_hunt17.NestedStructAccess},
			"StructSliceMethods":      {funcName: "StructSliceMethods", native: divergence_hunt17.StructSliceMethods},
			"FmtInterface":            {funcName: "FmtInterface", native: divergence_hunt17.FmtInterface},
			"FmtNilInterface":         {funcName: "FmtNilInterface", native: divergence_hunt17.FmtNilInterface},
			"StructComparison":        {funcName: "StructComparison", native: divergence_hunt17.StructComparison},
			"InterfaceEquality":       {funcName: "InterfaceEquality", native: divergence_hunt17.InterfaceEquality},
			"InterfaceInequality":     {funcName: "InterfaceInequality", native: divergence_hunt17.InterfaceInequality},
		},
	})
}

func TestDivergenceHunt18(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt18Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"StringToIntConversion":  {funcName: "StringToIntConversion", native: divergence_hunt18.StringToIntConversion},
			"IntToStringConversion":  {funcName: "IntToStringConversion", native: divergence_hunt18.IntToStringConversion},
			"FloatToStringConversion": {funcName: "FloatToStringConversion", native: divergence_hunt18.FloatToStringConversion},
			"StringToFloatConversion": {funcName: "StringToFloatConversion", native: divergence_hunt18.StringToFloatConversion},
			"BoolToStringConversion": {funcName: "BoolToStringConversion", native: divergence_hunt18.BoolToStringConversion},
			"StringToBoolConversion": {funcName: "StringToBoolConversion", native: divergence_hunt18.StringToBoolConversion},
			"StringSplitJoin":        {funcName: "StringSplitJoin", native: divergence_hunt18.StringSplitJoin},
			"StringTrimSpace":        {funcName: "StringTrimSpace", native: divergence_hunt18.StringTrimSpace},
			"StringTrimPrefix":       {funcName: "StringTrimPrefix", native: divergence_hunt18.StringTrimPrefix},
			"StringTrimSuffix":       {funcName: "StringTrimSuffix", native: divergence_hunt18.StringTrimSuffix},
			"StringReplaceAll":       {funcName: "StringReplaceAll", native: divergence_hunt18.StringReplaceAll},
			"StringBuilderPattern":   {funcName: "StringBuilderPattern", native: divergence_hunt18.StringBuilderPattern},
			"StringRuneConversion":   {funcName: "StringRuneConversion", native: divergence_hunt18.StringRuneConversion},
			"RuneToStringConversion": {funcName: "RuneToStringConversion", native: divergence_hunt18.RuneToStringConversion},
			"StringByteConversion":   {funcName: "StringByteConversion", native: divergence_hunt18.StringByteConversion},
			"ByteToStringConversion": {funcName: "ByteToStringConversion", native: divergence_hunt18.ByteToStringConversion},
			"FmtSprintfComplex":      {funcName: "FmtSprintfComplex", native: divergence_hunt18.FmtSprintfComplex},
			"FmtSprintfPadding":      {funcName: "FmtSprintfPadding", native: divergence_hunt18.FmtSprintfPadding},
			"StringPadLeft":          {funcName: "StringPadLeft", native: divergence_hunt18.StringPadLeft},
			"StringPadRight":         {funcName: "StringPadRight", native: divergence_hunt18.StringPadRight},
			"CamelCaseSplit":         {funcName: "CamelCaseSplit", native: divergence_hunt18.CamelCaseSplit},
			"StringReverseWords":     {funcName: "StringReverseWords", native: divergence_hunt18.StringReverseWords},
		},
	})
}

func TestDivergenceHunt19(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt19Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"EmptySliceOperations":   {funcName: "EmptySliceOperations", native: divergence_hunt19.EmptySliceOperations},
			"EmptyMapOperations":     {funcName: "EmptyMapOperations", native: divergence_hunt19.EmptyMapOperations},
			"EmptyStringOperations":  {funcName: "EmptyStringOperations", native: divergence_hunt19.EmptyStringOperations},
			"ZeroValueInt":           {funcName: "ZeroValueInt", native: divergence_hunt19.ZeroValueInt},
			"ZeroValueFloat":         {funcName: "ZeroValueFloat", native: divergence_hunt19.ZeroValueFloat},
			"ZeroValueBool":          {funcName: "ZeroValueBool", native: divergence_hunt19.ZeroValueBool},
			"ZeroValueString":        {funcName: "ZeroValueString", native: divergence_hunt19.ZeroValueString},
			"ZeroValueSlice":         {funcName: "ZeroValueSlice", native: divergence_hunt19.ZeroValueSlice},
			"ZeroValueMap":           {funcName: "ZeroValueMap", native: divergence_hunt19.ZeroValueMap},
			"ZeroValuePointer":       {funcName: "ZeroValuePointer", native: divergence_hunt19.ZeroValuePointer},
			"NilSliceAppend":         {funcName: "NilSliceAppend", native: divergence_hunt19.NilSliceAppend},
			"NilMapRead":             {funcName: "NilMapRead", native: divergence_hunt19.NilMapRead},
			"NilSliceRange":          {funcName: "NilSliceRange", native: divergence_hunt19.NilSliceRange},
			"NilMapRange":            {funcName: "NilMapRange", native: divergence_hunt19.NilMapRange},
			"NilChannelRead":         {funcName: "NilChannelRead", native: divergence_hunt19.NilChannelRead},
			"SliceBoundary":          {funcName: "SliceBoundary", native: divergence_hunt19.SliceBoundary},
			"MapBoundary":            {funcName: "MapBoundary", native: divergence_hunt19.MapBoundary},
			"ErrorHandlingPattern":   {funcName: "ErrorHandlingPattern", native: divergence_hunt19.ErrorHandlingPattern},
			"MultipleErrorCheck":     {funcName: "MultipleErrorCheck", native: divergence_hunt19.MultipleErrorCheck},
			"NilFuncVariable":        {funcName: "NilFuncVariable", native: divergence_hunt19.NilFuncVariable},
			"EmptyInterfaceContains": {funcName: "EmptyInterfaceContains", native: divergence_hunt19.EmptyInterfaceContains},
			"StructZeroValueFields":  {funcName: "StructZeroValueFields", native: divergence_hunt19.StructZeroValueFields},
		},
	})
}

func TestDivergenceHunt20(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt20Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"StudentGradeSystem": {funcName: "StudentGradeSystem", native: divergence_hunt20.StudentGradeSystem},
			"TextProcessing":     {funcName: "TextProcessing", native: divergence_hunt20.TextProcessing},
			"DataTransform":      {funcName: "DataTransform", native: divergence_hunt20.DataTransform},
			"InventorySystem":    {funcName: "InventorySystem", native: divergence_hunt20.InventorySystem},
			"JSONProcessing":     {funcName: "JSONProcessing", native: divergence_hunt20.JSONProcessing},
			"StringProcessing":   {funcName: "StringProcessing", native: divergence_hunt20.StringProcessing},
			"SortAndSearch":      {funcName: "SortAndSearch", native: divergence_hunt20.SortAndSearch},
			"MatrixOperations":   {funcName: "MatrixOperations", native: divergence_hunt20.MatrixOperations},
			"FmtTable":           {funcName: "FmtTable", native: divergence_hunt20.FmtTable},
			"Histogram":          {funcName: "Histogram", native: divergence_hunt20.Histogram},
			"ParseAndCompute":    {funcName: "ParseAndCompute", native: divergence_hunt20.ParseAndCompute},
			"SetOperations":      {funcName: "SetOperations", native: divergence_hunt20.SetOperations},
			"GroupBy":            {funcName: "GroupBy", native: divergence_hunt20.GroupBy},
			"RunningSum":         {funcName: "RunningSum", native: divergence_hunt20.RunningSum},
			"SlidingWindow":      {funcName: "SlidingWindow", native: divergence_hunt20.SlidingWindow},
		},
	})
}

func TestDivergenceHunt21(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt21Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"MapIterateSum":     {funcName: "MapIterateSum", native: divergence_hunt21.MapIterateSum},
			"SliceRotateLeft":   {funcName: "SliceRotateLeft", native: divergence_hunt21.SliceRotateLeft},
			"SliceRotateRight":  {funcName: "SliceRotateRight", native: divergence_hunt21.SliceRotateRight},
			"SliceChunk":        {funcName: "SliceChunk", native: divergence_hunt21.SliceChunk},
			"MapFilterSlice":    {funcName: "MapFilterSlice", native: divergence_hunt21.MapFilterSlice},
			"ReducePattern":     {funcName: "ReducePattern", native: divergence_hunt21.ReducePattern},
			"ZipSlices":         {funcName: "ZipSlices", native: divergence_hunt21.ZipSlices},
			"SliceCompact":      {funcName: "SliceCompact", native: divergence_hunt21.SliceCompact},
			"MapMergeOverwrite": {funcName: "MapMergeOverwrite", native: divergence_hunt21.MapMergeOverwrite},
			"SlicePartition":    {funcName: "SlicePartition", native: divergence_hunt21.SlicePartition},
			"NestedMapAccess":   {funcName: "NestedMapAccess", native: divergence_hunt21.NestedMapAccess},
			"FlattenMap":        {funcName: "FlattenMap", native: divergence_hunt21.FlattenMap},
			"MapKeySlice":       {funcName: "MapKeySlice", native: divergence_hunt21.MapKeySlice},
			"SliceSlidingWindow": {funcName: "SliceSlidingWindow", native: divergence_hunt21.SliceSlidingWindow},
			"MultiLevelSlice":   {funcName: "MultiLevelSlice", native: divergence_hunt21.MultiLevelSlice},
		},
	})
}

func TestDivergenceHunt22(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{
		src:       divergenceHunt22Src,
		buildOpts: []gig.BuildOption{gig.WithAllowPanic()},
		tests: map[string]divergenceTestCase{
			"JSONMarshalInt":     {funcName: "JSONMarshalInt", native: divergence_hunt22.JSONMarshalInt},
			"JSONMarshalString":  {funcName: "JSONMarshalString", native: divergence_hunt22.JSONMarshalString},
			"JSONMarshalFloat":   {funcName: "JSONMarshalFloat", native: divergence_hunt22.JSONMarshalFloat},
			"JSONUnmarshalInt":   {funcName: "JSONUnmarshalInt", native: divergence_hunt22.JSONUnmarshalInt},
			"JSONUnmarshalString": {funcName: "JSONUnmarshalString", native: divergence_hunt22.JSONUnmarshalString},
			"JSONUnmarshalFloat": {funcName: "JSONUnmarshalFloat", native: divergence_hunt22.JSONUnmarshalFloat},
			"JSONUnmarshalArray": {funcName: "JSONUnmarshalArray", native: divergence_hunt22.JSONUnmarshalArray},
			"FmtVerbP":           {funcName: "FmtVerbP", native: divergence_hunt22.FmtVerbP},
			"FmtVerbT":           {funcName: "FmtVerbT", native: divergence_hunt22.FmtVerbT},
			"FmtVerbV":           {funcName: "FmtVerbV", native: divergence_hunt22.FmtVerbV},
			"FmtVerbPlusV":       {funcName: "FmtVerbPlusV", native: divergence_hunt22.FmtVerbPlusV},
			"FmtVerbHashV":       {funcName: "FmtVerbHashV", native: divergence_hunt22.FmtVerbHashV},
			"FmtSprintfPointer":  {funcName: "FmtSprintfPointer", native: divergence_hunt22.FmtSprintfPointer},
			"ErrorWrap":          {funcName: "ErrorWrap", native: divergence_hunt22.ErrorWrap},
			"ErrorIs":            {funcName: "ErrorIs", native: divergence_hunt22.ErrorIs},
			"JSONNestedMap":      {funcName: "JSONNestedMap", native: divergence_hunt22.JSONNestedMap},
			"JSONStructTag":      {funcName: "JSONStructTag", native: divergence_hunt22.JSONStructTag},
			"JSONOmitEmpty":      {funcName: "JSONOmitEmpty", native: divergence_hunt22.JSONOmitEmpty},
			"FmtWidthInt":        {funcName: "FmtWidthInt", native: divergence_hunt22.FmtWidthInt},
			"FmtFloatScientific": {funcName: "FmtFloatScientific", native: divergence_hunt22.FmtFloatScientific},
		},
	})
}

func TestDivergenceHunt23(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt23Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"NewInt": {funcName: "NewInt", native: divergence_hunt23.NewInt}, "NewStruct": {funcName: "NewStruct", native: divergence_hunt23.NewStruct}, "MakeSliceLen": {funcName: "MakeSliceLen", native: divergence_hunt23.MakeSliceLen}, "MakeSliceLenCap": {funcName: "MakeSliceLenCap", native: divergence_hunt23.MakeSliceLenCap}, "MakeMapSize": {funcName: "MakeMapSize", native: divergence_hunt23.MakeMapSize}, "PointerSwap": {funcName: "PointerSwap", native: divergence_hunt23.PointerSwap}, "StructPointerNew": {funcName: "StructPointerNew", native: divergence_hunt23.StructPointerNew}, "SliceOfNew": {funcName: "SliceOfNew", native: divergence_hunt23.SliceOfNew}, "PointerToSlice": {funcName: "PointerToSlice", native: divergence_hunt23.PointerToSlice}, "PointerToMap": {funcName: "PointerToMap", native: divergence_hunt23.PointerToMap}, "DoublePointer": {funcName: "DoublePointer", native: divergence_hunt23.DoublePointer}, "PointerArithmeticSim": {funcName: "PointerArithmeticSim", native: divergence_hunt23.PointerArithmeticSim}, "NewArray": {funcName: "NewArray", native: divergence_hunt23.NewArray}, "SliceFromArray": {funcName: "SliceFromArray", native: divergence_hunt23.SliceFromArray}, "SliceFromArrayPointer": {funcName: "SliceFromArrayPointer", native: divergence_hunt23.SliceFromArrayPointer}, "MapPointer": {funcName: "MapPointer", native: divergence_hunt23.MapPointer}, "StructPointerMethod": {funcName: "StructPointerMethod", native: divergence_hunt23.StructPointerMethod}, "PointerComparison": {funcName: "PointerComparison", native: divergence_hunt23.PointerComparison}, "NilPointerComparison": {funcName: "NilPointerComparison", native: divergence_hunt23.NilPointerComparison},
	}})
}
func TestDivergenceHunt24(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt24Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SortAndDedupe": {funcName: "SortAndDedupe", native: divergence_hunt24.SortAndDedupe}, "WordFrequency": {funcName: "WordFrequency", native: divergence_hunt24.WordFrequency}, "CSVLikeParsing": {funcName: "CSVLikeParsing", native: divergence_hunt24.CSVLikeParsing}, "HistogramFromData": {funcName: "HistogramFromData", native: divergence_hunt24.HistogramFromData}, "FlattenJSON": {funcName: "FlattenJSON", native: divergence_hunt24.FlattenJSON}, "StringTokenize": {funcName: "StringTokenize", native: divergence_hunt24.StringTokenize}, "MatrixRowColSum": {funcName: "MatrixRowColSum", native: divergence_hunt24.MatrixRowColSum}, "StringTemplate": {funcName: "StringTemplate", native: divergence_hunt24.StringTemplate}, "MapTransformKeys": {funcName: "MapTransformKeys", native: divergence_hunt24.MapTransformKeys}, "SlicePartitionPoint": {funcName: "SlicePartitionPoint", native: divergence_hunt24.SlicePartitionPoint}, "NestedLoopBreak": {funcName: "NestedLoopBreak", native: divergence_hunt24.NestedLoopBreak}, "RecursiveSum": {funcName: "RecursiveSum", native: divergence_hunt24.RecursiveSum}, "ReverseSliceInPlace": {funcName: "ReverseSliceInPlace", native: divergence_hunt24.ReverseSliceInPlace}, "MapToSlice": {funcName: "MapToSlice", native: divergence_hunt24.MapToSlice}, "StringDiff": {funcName: "StringDiff", native: divergence_hunt24.StringDiff}, "FmtSlice": {funcName: "FmtSlice", native: divergence_hunt24.FmtSlice},
	}})
}
func TestDivergenceHunt25(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt25Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"DeferStack": {funcName: "DeferStack", native: divergence_hunt25.DeferStack}, "DeferInClosure": {funcName: "DeferInClosure", native: divergence_hunt25.DeferInClosure}, "RecoverInNestedDefer": {funcName: "RecoverInNestedDefer", native: divergence_hunt25.RecoverInNestedDefer}, "MultipleRecover": {funcName: "MultipleRecover", native: divergence_hunt25.MultipleRecover}, "DeferClosureCapture": {funcName: "DeferClosureCapture", native: divergence_hunt25.DeferClosureCapture}, "DeferClosureCopy": {funcName: "DeferClosureCopy", native: divergence_hunt25.DeferClosureCopy}, "PanicInDeferRecover": {funcName: "PanicInDeferRecover", native: divergence_hunt25.PanicInDeferRecover}, "DeferModifyNamedReturn": {funcName: "DeferModifyNamedReturn", native: divergence_hunt25.DeferModifyNamedReturn}, "NestedPanicRecover": {funcName: "NestedPanicRecover", native: divergence_hunt25.NestedPanicRecover}, "ClosureWithDefer": {funcName: "ClosureWithDefer", native: divergence_hunt25.ClosureWithDefer}, "RecursiveWithDefer": {funcName: "RecursiveWithDefer", native: divergence_hunt25.RecursiveWithDefer}, "PanicRecoverTypeSwitch": {funcName: "PanicRecoverTypeSwitch", native: divergence_hunt25.PanicRecoverTypeSwitch}, "DeferMultipleModifies": {funcName: "DeferMultipleModifies", native: divergence_hunt25.DeferMultipleModifies}, "RecoverReturnsPanicValue": {funcName: "RecoverReturnsPanicValue", native: divergence_hunt25.RecoverReturnsPanicValue}, "DeferInMethod": {funcName: "DeferInMethod", native: divergence_hunt25.DeferInMethod}, "ClosureState": {funcName: "ClosureState", native: divergence_hunt25.ClosureState}, "ClosureSharedState": {funcName: "ClosureSharedState", native: divergence_hunt25.ClosureSharedState}, "FmtDefer": {funcName: "FmtDefer", native: divergence_hunt25.FmtDefer},
	}})
}
func TestDivergenceHunt26(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt26Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Int8Range": {funcName: "Int8Range", native: divergence_hunt26.Int8Range}, "Int8MinRange": {funcName: "Int8MinRange", native: divergence_hunt26.Int8MinRange}, "Uint8Max": {funcName: "Uint8Max", native: divergence_hunt26.Uint8Max}, "Int16Range": {funcName: "Int16Range", native: divergence_hunt26.Int16Range}, "Uint16Max": {funcName: "Uint16Max", native: divergence_hunt26.Uint16Max}, "Float32Smallest": {funcName: "Float32Smallest", native: divergence_hunt26.Float32Smallest}, "Complex64Basic": {funcName: "Complex64Basic", native: divergence_hunt26.Complex64Basic}, "Complex128Basic": {funcName: "Complex128Basic", native: divergence_hunt26.Complex128Basic}, "RuneType": {funcName: "RuneType", native: divergence_hunt26.RuneType}, "ByteType": {funcName: "ByteType", native: divergence_hunt26.ByteType}, "StringType": {funcName: "StringType", native: divergence_hunt26.StringType}, "BoolType": {funcName: "BoolType", native: divergence_hunt26.BoolType}, "IntType": {funcName: "IntType", native: divergence_hunt26.IntType}, "Int64Type": {funcName: "Int64Type", native: divergence_hunt26.Int64Type}, "UintType": {funcName: "UintType", native: divergence_hunt26.UintType}, "Uint64Type": {funcName: "Uint64Type", native: divergence_hunt26.Uint64Type}, "Float64Type": {funcName: "Float64Type", native: divergence_hunt26.Float64Type}, "Float32Type": {funcName: "Float32Type", native: divergence_hunt26.Float32Type}, "TypeConversionChain": {funcName: "TypeConversionChain", native: divergence_hunt26.TypeConversionChain}, "UnsignedConversion": {funcName: "UnsignedConversion", native: divergence_hunt26.UnsignedConversion}, "SignedToUnsigned": {funcName: "SignedToUnsigned", native: divergence_hunt26.SignedToUnsigned}, "FloatToIntTrunc": {funcName: "FloatToIntTrunc", native: divergence_hunt26.FloatToIntTrunc}, "IntToFloatPrecise": {funcName: "IntToFloatPrecise", native: divergence_hunt26.IntToFloatPrecise}, "StringToSlice": {funcName: "StringToSlice", native: divergence_hunt26.StringToSlice}, "SliceToString": {funcName: "SliceToString", native: divergence_hunt26.SliceToString},
	}})
}
func TestDivergenceHunt27(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt27Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"StringSort": {funcName: "StringSort", native: divergence_hunt27.StringSort}, "StringUnique": {funcName: "StringUnique", native: divergence_hunt27.StringUnique}, "StringIsDigit": {funcName: "StringIsDigit", native: divergence_hunt27.StringIsDigit}, "StringIsAlpha": {funcName: "StringIsAlpha", native: divergence_hunt27.StringIsAlpha}, "StringToUpperLower": {funcName: "StringToUpperLower", native: divergence_hunt27.StringToUpperLower}, "StringCapitalize": {funcName: "StringCapitalize", native: divergence_hunt27.StringCapitalize}, "StringCountWords": {funcName: "StringCountWords", native: divergence_hunt27.StringCountWords}, "StringReverseWords": {funcName: "StringReverseWords", native: divergence_hunt27.StringReverseWords}, "FmtInteger": {funcName: "FmtInteger", native: divergence_hunt27.FmtInteger}, "FmtHexInt": {funcName: "FmtHexInt", native: divergence_hunt27.FmtHexInt}, "FmtOctalInt": {funcName: "FmtOctalInt", native: divergence_hunt27.FmtOctalInt}, "FmtBinaryInt": {funcName: "FmtBinaryInt", native: divergence_hunt27.FmtBinaryInt}, "FmtCharFromInt": {funcName: "FmtCharFromInt", native: divergence_hunt27.FmtCharFromInt}, "FmtUnicode": {funcName: "FmtUnicode", native: divergence_hunt27.FmtUnicode}, "SortIntSliceDesc": {funcName: "SortIntSliceDesc", native: divergence_hunt27.SortIntSliceDesc}, "SortFloatSliceDesc": {funcName: "SortFloatSliceDesc", native: divergence_hunt27.SortFloatSliceDesc}, "StringJoinWithSep": {funcName: "StringJoinWithSep", native: divergence_hunt27.StringJoinWithSep}, "StringSplitN": {funcName: "StringSplitN", native: divergence_hunt27.StringSplitN}, "StringRepeatN": {funcName: "StringRepeatN", native: divergence_hunt27.StringRepeatN}, "StringMapFunc": {funcName: "StringMapFunc", native: divergence_hunt27.StringMapFunc},
	}})
}
func TestDivergenceHunt28(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt28Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"ChannelSendRecv": {funcName: "ChannelSendRecv", native: divergence_hunt28.ChannelSendRecv}, "ChannelBuffered": {funcName: "ChannelBuffered", native: divergence_hunt28.ChannelBuffered}, "ChannelCloseRange": {funcName: "ChannelCloseRange", native: divergence_hunt28.ChannelCloseRange}, "ChannelSelectTwo": {funcName: "ChannelSelectTwo", native: divergence_hunt28.ChannelSelectTwo}, "ChannelSelectDefault2": {funcName: "ChannelSelectDefault2", native: divergence_hunt28.ChannelSelectDefault2}, "ChannelNilSelect": {funcName: "ChannelNilSelect", native: divergence_hunt28.ChannelNilSelect}, "ChannelLen": {funcName: "ChannelLen", native: divergence_hunt28.ChannelLen}, "ChannelCap2": {funcName: "ChannelCap2", native: divergence_hunt28.ChannelCap2}, "ChannelRecvAfterClose": {funcName: "ChannelRecvAfterClose", native: divergence_hunt28.ChannelRecvAfterClose}, "ChannelDirection": {funcName: "ChannelDirection", native: divergence_hunt28.ChannelDirection}, "SelectMultipleReady": {funcName: "SelectMultipleReady", native: divergence_hunt28.SelectMultipleReady}, "ChannelAsSignal": {funcName: "ChannelAsSignal", native: divergence_hunt28.ChannelAsSignal},
	}})
}
func TestDivergenceHunt29(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt29Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"SimpleError": {funcName: "SimpleError", native: divergence_hunt29.SimpleError}, "ErrorWithFormat": {funcName: "ErrorWithFormat", native: divergence_hunt29.ErrorWithFormat}, "ValidatePositive": {funcName: "ValidatePositive", native: divergence_hunt29.ValidatePositive}, "ValidateRange": {funcName: "ValidateRange", native: divergence_hunt29.ValidateRange}, "ErrorPropagation": {funcName: "ErrorPropagation", native: divergence_hunt29.ErrorPropagation}, "ErrorInDefer": {funcName: "ErrorInDefer", native: divergence_hunt29.ErrorInDefer}, "MultiErrorCollect": {funcName: "MultiErrorCollect", native: divergence_hunt29.MultiErrorCollect}, "ErrorTypeAssertion": {funcName: "ErrorTypeAssertion", native: divergence_hunt29.ErrorTypeAssertion}, "PanicWithFmtError": {funcName: "PanicWithFmtError", native: divergence_hunt29.PanicWithFmtError}, "NilErrorCheck": {funcName: "NilErrorCheck", native: divergence_hunt29.NilErrorCheck}, "ErrorStringMethods": {funcName: "ErrorStringMethods", native: divergence_hunt29.ErrorStringMethods}, "ValidateStruct": {funcName: "ValidateStruct", native: divergence_hunt29.ValidateStruct}, "ErrorInClosure": {funcName: "ErrorInClosure", native: divergence_hunt29.ErrorInClosure}, "FmtErrorfWrap": {funcName: "FmtErrorfWrap", native: divergence_hunt29.FmtErrorfWrap},
	}})
}
func TestDivergenceHunt30(t *testing.T) {
	runDivergenceTestSet(t, divergenceTestSet{src: divergenceHunt30Src, buildOpts: []gig.BuildOption{gig.WithAllowPanic()}, tests: map[string]divergenceTestCase{
		"Comprehensive1": {funcName: "Comprehensive1", native: divergence_hunt30.Comprehensive1}, "Comprehensive2": {funcName: "Comprehensive2", native: divergence_hunt30.Comprehensive2}, "Comprehensive3": {funcName: "Comprehensive3", native: divergence_hunt30.Comprehensive3}, "Comprehensive4": {funcName: "Comprehensive4", native: divergence_hunt30.Comprehensive4}, "Comprehensive5": {funcName: "Comprehensive5", native: divergence_hunt30.Comprehensive5}, "Comprehensive6": {funcName: "Comprehensive6", native: divergence_hunt30.Comprehensive6}, "Comprehensive7": {funcName: "Comprehensive7", native: divergence_hunt30.Comprehensive7}, "Comprehensive8": {funcName: "Comprehensive8", native: divergence_hunt30.Comprehensive8}, "Comprehensive9": {funcName: "Comprehensive9", native: divergence_hunt30.Comprehensive9}, "Comprehensive10": {funcName: "Comprehensive10", native: divergence_hunt30.Comprehensive10},
	}})
}
