package tests

import (
	_ "embed"
	"testing"

	thirdparty "git.woa.com/youngjin/gig/tests/testdata/thirdparty"
)

// Embed each category source file for interpreter testing.
var (
	//go:embed testdata/thirdparty/bytes.go
	srcBytes string
	//go:embed testdata/thirdparty/strings.go
	srcStrings string
	//go:embed testdata/thirdparty/strconv.go
	srcStrconv string
	//go:embed testdata/thirdparty/math.go
	srcMath string
	//go:embed testdata/thirdparty/time.go
	srcTime string
	//go:embed testdata/thirdparty/context.go
	srcContext string
	//go:embed testdata/thirdparty/sync.go
	srcSync string
	//go:embed testdata/thirdparty/sort.go
	srcSort string
	//go:embed testdata/thirdparty/encoding.go
	srcEncoding string
	//go:embed testdata/thirdparty/io.go
	srcIO string
	// path/filepath removed from sandbox stdlib (OS-specific, uses filesystem)
	// //go:embed testdata/thirdparty/filepath.go
	// srcFilepath string
	//go:embed testdata/thirdparty/regexp.go
	srcRegexp string
	//go:embed testdata/thirdparty/errors.go
	srcErrors string
	//go:embed testdata/thirdparty/fmt.go
	srcFmt string
	//go:embed testdata/thirdparty/patterns.go
	srcPatterns string
)

// TestCorrectnessThirdparty tests third-party library calls through the interpreter,
// comparing interpreted results with native Go execution.
func TestCorrectnessThirdparty(t *testing.T) {
	for name, set := range thirdpartyTestSets {
		t.Run(name, func(t *testing.T) {
			runTestSet(t, set)
		})
	}
}

// ============================================================================
// Test registrations — native function references, no hardcoded values
// ============================================================================

var bytesTests = map[string]testCase{
	"BytesBufferWrite":       {srcBytes, "BytesBufferWrite", nil, thirdparty.BytesBufferWrite},
	"BytesBufferWriteString": {srcBytes, "BytesBufferWriteString", nil, thirdparty.BytesBufferWriteString},
	"BytesBufferReadFrom":    {srcBytes, "BytesBufferReadFrom", nil, thirdparty.BytesBufferReadFrom},
	"BytesBufferString":      {srcBytes, "BytesBufferString", nil, thirdparty.BytesBufferString},
	"BytesBufferLen":         {srcBytes, "BytesBufferLen", nil, thirdparty.BytesBufferLen},
	"BytesBufferGrow":        {srcBytes, "BytesBufferGrow", nil, thirdparty.BytesBufferGrow},
	"BytesBufferNext":        {srcBytes, "BytesBufferNext", nil, thirdparty.BytesBufferNext},
	"BytesBufferReadByte":    {srcBytes, "BytesBufferReadByte", nil, thirdparty.BytesBufferReadByte},
	"BytesBufferUnreadByte":  {srcBytes, "BytesBufferUnreadByte", nil, thirdparty.BytesBufferUnreadByte},
	"BytesBufferReadBytes":   {srcBytes, "BytesBufferReadBytes", nil, thirdparty.BytesBufferReadBytes},
	"BytesBufferReadString":  {srcBytes, "BytesBufferReadString", nil, thirdparty.BytesBufferReadString},
	"BytesNewBuffer":         {srcBytes, "BytesNewBuffer", nil, thirdparty.BytesNewBuffer},
	"BytesNewBufferString":  {srcBytes, "BytesNewBufferString", nil, thirdparty.BytesNewBufferString},
	"BytesBufferTrim":       {srcBytes, "BytesBufferTrim", nil, thirdparty.BytesBufferTrim},
	"BytesSplit":            {srcBytes, "BytesSplit", nil, thirdparty.BytesSplit},
	"BytesSplitN":           {srcBytes, "BytesSplitN", nil, thirdparty.BytesSplitN},
	"BytesJoin":             {srcBytes, "BytesJoin", nil, thirdparty.BytesJoin},
	"BytesContains":         {srcBytes, "BytesContains", nil, thirdparty.BytesContains},
	"BytesCount":            {srcBytes, "BytesCount", nil, thirdparty.BytesCount},
	"BytesIndex":            {srcBytes, "BytesIndex", nil, thirdparty.BytesIndex},
	"BytesLastIndex":        {srcBytes, "BytesLastIndex", nil, thirdparty.BytesLastIndex},
	"BytesHasPrefix":        {srcBytes, "BytesHasPrefix", nil, thirdparty.BytesHasPrefix},
	"BytesHasSuffix":        {srcBytes, "BytesHasSuffix", nil, thirdparty.BytesHasSuffix},
	"BytesReplace":          {srcBytes, "BytesReplace", nil, thirdparty.BytesReplace},
	"BytesReplaceAll":       {srcBytes, "BytesReplaceAll", nil, thirdparty.BytesReplaceAll},
	"BytesFields":           {srcBytes, "BytesFields", nil, thirdparty.BytesFields},
	"BytesTrimSpace":        {srcBytes, "BytesTrimSpace", nil, thirdparty.BytesTrimSpace},
	"BytesToUpper":          {srcBytes, "BytesToUpper", nil, thirdparty.BytesToUpper},
	"BytesToLower":          {srcBytes, "BytesToLower", nil, thirdparty.BytesToLower},
	"BytesTrim":             {srcBytes, "BytesTrim", nil, thirdparty.BytesTrim},
	"BytesMap":              {srcBytes, "BytesMap", nil, thirdparty.BytesMap},
}

var stringsTests = map[string]testCase{
	"StringsBuilder":        {srcStrings, "StringsBuilder", nil, thirdparty.StringsBuilder},
	"StringsBuilderString":  {srcStrings, "StringsBuilderString", nil, thirdparty.StringsBuilderString},
	"StringsBuilderGrow":   {srcStrings, "StringsBuilderGrow", nil, thirdparty.StringsBuilderGrow},
	"StringsMap":           {srcStrings, "StringsMap", nil, thirdparty.StringsMap},
	"StringsRepeat":        {srcStrings, "StringsRepeat", nil, thirdparty.StringsRepeat},
	"StringsRepeatCount":   {srcStrings, "StringsRepeatCount", nil, thirdparty.StringsRepeatCount},
	"StringsIndexAny":      {srcStrings, "StringsIndexAny", nil, thirdparty.StringsIndexAny},
	"StringsIndexFunc":     {srcStrings, "StringsIndexFunc", nil, thirdparty.StringsIndexFunc},
	"StringsTitle":         {srcStrings, "StringsTitle", nil, thirdparty.StringsTitle},
	"StringsToTitle":       {srcStrings, "StringsToTitle", nil, thirdparty.StringsToTitle},
	"StringsToValidUTF8":   {srcStrings, "StringsToValidUTF8", nil, thirdparty.StringsToValidUTF8},
	"StringsTrimLeft":      {srcStrings, "StringsTrimLeft", nil, thirdparty.StringsTrimLeft},
	"StringsTrimRight":     {srcStrings, "StringsTrimRight", nil, thirdparty.StringsTrimRight},
	"StringsTrimFunc":      {srcStrings, "StringsTrimFunc", nil, thirdparty.StringsTrimFunc},
	"StringsIndexFuncTest": {srcStrings, "StringsIndexFuncTest", nil, thirdparty.StringsIndexFuncTest},
	"StringsCut":           {srcStrings, "StringsCut", nil, thirdparty.StringsCut},
	"StringsCutPrefix":     {srcStrings, "StringsCutPrefix", nil, thirdparty.StringsCutPrefix},
	"StringsCutSuffix":     {srcStrings, "StringsCutSuffix", nil, thirdparty.StringsCutSuffix},
}

var strconvTests = map[string]testCase{
	"StrconvParseBool":    {srcStrconv, "StrconvParseBool", nil, thirdparty.StrconvParseBool},
	"StrconvFormatBool":   {srcStrconv, "StrconvFormatBool", nil, thirdparty.StrconvFormatBool},
	"StrconvParseInt":     {srcStrconv, "StrconvParseInt", nil, thirdparty.StrconvParseInt},
	"StrconvParseUint":    {srcStrconv, "StrconvParseUint", nil, thirdparty.StrconvParseUint},
	"StrconvFormatInt":    {srcStrconv, "StrconvFormatInt", nil, thirdparty.StrconvFormatInt},
	"StrconvFormatUint":   {srcStrconv, "StrconvFormatUint", nil, thirdparty.StrconvFormatUint},
	"StrconvParseFloat":   {srcStrconv, "StrconvParseFloat", nil, thirdparty.StrconvParseFloat},
	"StrconvFormatFloat":  {srcStrconv, "StrconvFormatFloat", nil, thirdparty.StrconvFormatFloat},
	"StrconvQuote":        {srcStrconv, "StrconvQuote", nil, thirdparty.StrconvQuote},
	"StrconvQuoteToASCII": {srcStrconv, "StrconvQuoteToASCII", nil, thirdparty.StrconvQuoteToASCII},
	"StrconvUnquote":      {srcStrconv, "StrconvUnquote", nil, thirdparty.StrconvUnquote},
	"StrconvAppendInt":    {srcStrconv, "StrconvAppendInt", nil, thirdparty.StrconvAppendInt},
	"StrconvAppendFloat":  {srcStrconv, "StrconvAppendFloat", nil, thirdparty.StrconvAppendFloat},
}

var mathTests = map[string]testCase{
	"MathAbs":      {srcMath, "MathAbs", nil, thirdparty.MathAbs},
	"MathMax":      {srcMath, "MathMax", nil, thirdparty.MathMax},
	"MathMin":      {srcMath, "MathMin", nil, thirdparty.MathMin},
	"MathFloor":    {srcMath, "MathFloor", nil, thirdparty.MathFloor},
	"MathCeil":     {srcMath, "MathCeil", nil, thirdparty.MathCeil},
	"MathRound":    {srcMath, "MathRound", nil, thirdparty.MathRound},
	"MathPow":      {srcMath, "MathPow", nil, thirdparty.MathPow},
	"MathSqrt":     {srcMath, "MathSqrt", nil, thirdparty.MathSqrt},
	"MathMod":      {srcMath, "MathMod", nil, thirdparty.MathMod},
	"MathSin":      {srcMath, "MathSin", nil, thirdparty.MathSin},
	"MathCos":      {srcMath, "MathCos", nil, thirdparty.MathCos},
	"MathTan":      {srcMath, "MathTan", nil, thirdparty.MathTan},
	"MathLog":      {srcMath, "MathLog", nil, thirdparty.MathLog},
	"MathLog10":    {srcMath, "MathLog10", nil, thirdparty.MathLog10},
	"MathExp":      {srcMath, "MathExp", nil, thirdparty.MathExp},
	"MathInf":      {srcMath, "MathInf", nil, thirdparty.MathInf},
	"MathNaN":      {srcMath, "MathNaN", nil, thirdparty.MathNaN},
	"MathCopysign": {srcMath, "MathCopysign", nil, thirdparty.MathCopysign},
}

var timeTests = map[string]testCase{
	"TimeNow":          {srcTime, "TimeNow", nil, thirdparty.TimeNow},
	"TimeParse":        {srcTime, "TimeParse", nil, thirdparty.TimeParse},
	"TimeFormat":       {srcTime, "TimeFormat", nil, thirdparty.TimeFormat},
	"TimeAdd":          {srcTime, "TimeAdd", nil, thirdparty.TimeAdd},
	"TimeSub":          {srcTime, "TimeSub", nil, thirdparty.TimeSub},
	"TimeBefore":       {srcTime, "TimeBefore", nil, thirdparty.TimeBefore},
	"TimeAfter":        {srcTime, "TimeAfter", nil, thirdparty.TimeAfter},
	"TimeEqual":        {srcTime, "TimeEqual", nil, thirdparty.TimeEqual},
	"TimeUnix":         {srcTime, "TimeUnix", nil, thirdparty.TimeUnix},
	"TimeUnixMilli":    {srcTime, "TimeUnixMilli", nil, thirdparty.TimeUnixMilli},
	"TimeDuration":     {srcTime, "TimeDuration", nil, thirdparty.TimeDuration},
	"TimeSleep":        {srcTime, "TimeSleep", nil, thirdparty.TimeSleep},
}

var contextTests = map[string]testCase{
	"ContextBackground":      {srcContext, "ContextBackground", nil, thirdparty.ContextBackground},
	"ContextTODO":            {srcContext, "ContextTODO", nil, thirdparty.ContextTODO},
	"ContextWithValue":       {srcContext, "ContextWithValue", nil, thirdparty.ContextWithValue},
	"ContextWithCancel":      {srcContext, "ContextWithCancel", nil, thirdparty.ContextWithCancel},
	"ContextWithTimeout":     {srcContext, "ContextWithTimeout", nil, thirdparty.ContextWithTimeout},
	"ContextWithCancelParent": {srcContext, "ContextWithCancelParent", nil, thirdparty.ContextWithCancelParent},
}

var syncTests = map[string]testCase{
	"SyncMutex":          {srcSync, "SyncMutex", nil, thirdparty.SyncMutex},
	"SyncMutexCounter":   {srcSync, "SyncMutexCounter", nil, thirdparty.SyncMutexCounter},
	"SyncRWMutex":        {srcSync, "SyncRWMutex", nil, thirdparty.SyncRWMutex},
	"SyncWaitGroup":      {srcSync, "SyncWaitGroup", nil, thirdparty.SyncWaitGroup},
	"SyncOnce":           {srcSync, "SyncOnce", nil, thirdparty.SyncOnce},
	"SyncOnceFunc":       {srcSync, "SyncOnceFunc", nil, thirdparty.SyncOnceFunc},
	"SyncMap":            {srcSync, "SyncMap", nil, thirdparty.SyncMap},
	"SyncMapLoadOrStore": {srcSync, "SyncMapLoadOrStore", nil, thirdparty.SyncMapLoadOrStore},
}

var sortTests = map[string]testCase{
	"SortStrings":       {srcSort, "SortStrings", nil, thirdparty.SortStrings},
	"SortInts":          {srcSort, "SortInts", nil, thirdparty.SortInts},
	"SortFloat64s":      {srcSort, "SortFloat64s", nil, thirdparty.SortFloat64s},
	"SortSearchInts":     {srcSort, "SortSearchInts", nil, thirdparty.SortSearchInts},
	"SortSearchStrings": {srcSort, "SortSearchStrings", nil, thirdparty.SortSearchStrings},
	"SortSlice":         {srcSort, "SortSlice", nil, thirdparty.SortSlice},
	"SortSliceStable":   {srcSort, "SortSliceStable", nil, thirdparty.SortSliceStable},
	"SortIsSorted":      {srcSort, "SortIsSorted", nil, thirdparty.SortIsSorted},
}

var encodingTests = map[string]testCase{
	"JsonMarshal":       {srcEncoding, "JsonMarshal", nil, thirdparty.JsonMarshal},
	"JsonUnmarshal":     {srcEncoding, "JsonUnmarshal", nil, thirdparty.JsonUnmarshal},
	"JsonMarshalIndent": {srcEncoding, "JsonMarshalIndent", nil, thirdparty.JsonMarshalIndent},
	// JsonDecode skipped — Bug 8: json.Decoder.Decode vs xml.Decoder.Decode type collision
	"JsonNumber": {srcEncoding, "JsonNumber", nil, thirdparty.JsonNumber},
	"Base64Encode":      {srcEncoding, "Base64Encode", nil, thirdparty.Base64Encode},
	"Base64Decode":      {srcEncoding, "Base64Decode", nil, thirdparty.Base64Decode},
	"Base64URLEncode":   {srcEncoding, "Base64URLEncode", nil, thirdparty.Base64URLEncode},
	"Base64URLDecode":   {srcEncoding, "Base64URLDecode", nil, thirdparty.Base64URLDecode},
	"Base64NewEncoder":  {srcEncoding, "Base64NewEncoder", nil, thirdparty.Base64NewEncoder},
	"Base64NewDecoder":  {srcEncoding, "Base64NewDecoder", nil, thirdparty.Base64NewDecoder},
	"HexEncodeToString": {srcEncoding, "HexEncodeToString", nil, thirdparty.HexEncodeToString},
	"HexDecodeString":   {srcEncoding, "HexDecodeString", nil, thirdparty.HexDecodeString},
	"HexNewEncoder":     {srcEncoding, "HexNewEncoder", nil, thirdparty.HexNewEncoder},
	"HexNewDecoder":     {srcEncoding, "HexNewDecoder", nil, thirdparty.HexNewDecoder},
}

var ioTests = map[string]testCase{
	"IoReadAll":        {srcIO, "IoReadAll", nil, thirdparty.IoReadAll},
	"IoCopy":           {srcIO, "IoCopy", nil, thirdparty.IoCopy},
	"IoReadFull":       {srcIO, "IoReadFull", nil, thirdparty.IoReadFull},
	"IoWriteString":    {srcIO, "IoWriteString", nil, thirdparty.IoWriteString},
	"IoSectionReader":  {srcIO, "IoSectionReader", nil, thirdparty.IoSectionReader},
	"IoLimitedReader":  {srcIO, "IoLimitedReader", nil, thirdparty.IoLimitedReader},
	"IoTeeReader":      {srcIO, "IoTeeReader", nil, thirdparty.IoTeeReader},
}

var regexpTests = map[string]testCase{
	"RegexpMatch":             {srcRegexp, "RegexpMatch", nil, thirdparty.RegexpMatch},
	"RegexpCompile":           {srcRegexp, "RegexpCompile", nil, thirdparty.RegexpCompile},
	"RegexpMustCompile":       {srcRegexp, "RegexpMustCompile", nil, thirdparty.RegexpMustCompile},
	"RegexpFindString":        {srcRegexp, "RegexpFindString", nil, thirdparty.RegexpFindString},
	"RegexpFindStringSubmatch": {srcRegexp, "RegexpFindStringSubmatch", nil, thirdparty.RegexpFindStringSubmatch},
	"RegexpFindAllString":    {srcRegexp, "RegexpFindAllString", nil, thirdparty.RegexpFindAllString},
	"RegexpReplaceAllString": {srcRegexp, "RegexpReplaceAllString", nil, thirdparty.RegexpReplaceAllString},
	"RegexpSplit":             {srcRegexp, "RegexpSplit", nil, thirdparty.RegexpSplit},
	"RegexpNumSubexp":         {srcRegexp, "RegexpNumSubexp", nil, thirdparty.RegexpNumSubexp},
	"RegexpLongest":           {srcRegexp, "RegexpLongest", nil, thirdparty.RegexpLongest},
}

var errorsTests = map[string]testCase{
	"ErrorsNew":  {srcErrors, "ErrorsNew", nil, thirdparty.ErrorsNew},
	"ErrorsIs":   {srcErrors, "ErrorsIs", nil, thirdparty.ErrorsIs},
	"ErrorsAs":   {srcErrors, "ErrorsAs", nil, thirdparty.ErrorsAs},
	"ErrorsJoin": {srcErrors, "ErrorsJoin", nil, thirdparty.ErrorsJoin},
}

var fmtTests = map[string]testCase{
	"FmtSprintfVarious": {srcFmt, "FmtSprintfVarious", nil, thirdparty.FmtSprintfVarious},
	"FmtSprintfStruct":  {srcFmt, "FmtSprintfStruct", nil, thirdparty.FmtSprintfStruct},
	"FmtSprintfPointer": {srcFmt, "FmtSprintfPointer", nil, thirdparty.FmtSprintfPointer},
	"FmtSprintfBool":    {srcFmt, "FmtSprintfBool", nil, thirdparty.FmtSprintfBool},
	"FmtSprintfHex":     {srcFmt, "FmtSprintfHex", nil, thirdparty.FmtSprintfHex},
	"FmtFprintf":        {srcFmt, "FmtFprintf", nil, thirdparty.FmtFprintf},
	"FmtSprint":          {srcFmt, "FmtSprint", nil, thirdparty.FmtSprint},
	"FmtSprintln":        {srcFmt, "FmtSprintln", nil, thirdparty.FmtSprintln},
	"FmtErrorf":         {srcFmt, "FmtErrorf", nil, thirdparty.FmtErrorf},
}

var patternsTests = map[string]testCase{
	"ChainBytesToStringToBase64":   {srcPatterns, "ChainBytesToStringToBase64", nil, thirdparty.ChainBytesToStringToBase64},
	"ChainStringsBuilderToBuffer":  {srcPatterns, "ChainStringsBuilderToBuffer", nil, thirdparty.ChainStringsBuilderToBuffer},
	"ChainSortSearch":              {srcPatterns, "ChainSortSearch", nil, thirdparty.ChainSortSearch},
	"ChainBufferWriteRead":         {srcPatterns, "ChainBufferWriteRead", nil, thirdparty.ChainBufferWriteRead},
	"ChainContextWithValueChain":   {srcPatterns, "ChainContextWithValueChain", nil, thirdparty.ChainContextWithValueChain},
	"InterfaceWithPointerReceiver": {srcPatterns, "InterfaceWithPointerReceiver", nil, thirdparty.InterfaceWithPointerReceiver},
	"InterfaceSliceOfPointers":     {srcPatterns, "InterfaceSliceOfPointers", nil, thirdparty.InterfaceSliceOfPointers},
	"InterfaceMap":                 {srcPatterns, "InterfaceMap", nil, thirdparty.InterfaceMap},
	"VariadicAppend":               {srcPatterns, "VariadicAppend", nil, thirdparty.VariadicAppend},
	"VariadicStringsJoin":          {srcPatterns, "VariadicStringsJoin", nil, thirdparty.VariadicStringsJoin},
	"VariadicAppendSlice":          {srcPatterns, "VariadicAppendSlice", nil, thirdparty.VariadicAppendSlice},
	"MethodChainingBuilder":        {srcPatterns, "MethodChainingBuilder", nil, thirdparty.MethodChainingBuilder},
	"TableDrivenOp":                {srcPatterns, "TableDrivenOp", nil, thirdparty.TableDrivenOp},
	"FunctionValueFromMap":         {srcPatterns, "FunctionValueFromMap", nil, thirdparty.FunctionValueFromMap},
	"DeferWithMutex":               {srcPatterns, "DeferWithMutex", nil, thirdparty.DeferWithMutex},
	"SelectWithChannels":           {srcPatterns, "SelectWithChannels", nil, thirdparty.SelectWithChannels},
}

var thirdpartyTestSets = map[string]testSet{
	"Bytes":    {name: "Bytes", src: srcBytes, tests: bytesTests},
	"Strings":  {name: "Strings", src: srcStrings, tests: stringsTests},
	"Strconv":  {name: "Strconv", src: srcStrconv, tests: strconvTests},
	"Math":     {name: "Math", src: srcMath, tests: mathTests},
	"Time":     {name: "Time", src: srcTime, tests: timeTests},
	"Context":  {name: "Context", src: srcContext, tests: contextTests},
	"Sync":     {name: "Sync", src: srcSync, tests: syncTests},
	"Sort":     {name: "Sort", src: srcSort, tests: sortTests},
	"Encoding": {name: "Encoding", src: srcEncoding, tests: encodingTests},
	"IO":       {name: "IO", src: srcIO, tests: ioTests},
	"Regexp":   {name: "Regexp", src: srcRegexp, tests: regexpTests},
	"Errors":   {name: "Errors", src: srcErrors, tests: errorsTests},
	"Fmt":      {name: "Fmt", src: srcFmt, tests: fmtTests},
	"Patterns": {name: "Patterns", src: srcPatterns, tests: patternsTests},
}
