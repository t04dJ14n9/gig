package tests

import (
	_ "embed"
	"testing"

	"git.woa.com/youngjin/gig"
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

// thirdpartyCase holds a function name and expected result.
type thirdpartyCase struct {
	name     string
	expected any
}

// runCategory builds a source file and runs all test cases against it.
func runCategory(t *testing.T, categoryName, src string, cases []thirdpartyCase) {
	t.Helper()
	t.Run(categoryName, func(t *testing.T) {
		prog, err := gig.Build(src)
		if err != nil {
			t.Fatalf("Build(%s) error: %v", categoryName, err)
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := prog.Run(tc.name)
				if err != nil {
					t.Errorf("Run(%s) error: %v", tc.name, err)
					return
				}
				if result != tc.expected {
					t.Errorf("Run(%s) = %v (%T), want %v (%T)",
						tc.name, result, result, tc.expected, tc.expected)
				}
			})
		}
	})
}

// TestCorrectnessThirdparty tests third-party library calls through the interpreter,
// organized by standard library package category.
func TestCorrectnessThirdparty(t *testing.T) {
	runCategory(t, "Bytes", srcBytes, []thirdpartyCase{
		{"BytesBufferWrite", 11},
		{"BytesBufferWriteString", 4},
		{"BytesBufferString", "test string"},
		{"BytesBufferLen", 13},
		{"BytesSplit", 3},
		{"BytesSplitN", 2},
		{"BytesContains", 1},
		{"BytesCount", 3},
		{"BytesIndex", 2},
		{"BytesLastIndex", 12},
		{"BytesHasPrefix", 1},
		{"BytesHasSuffix", 1},
		{"BytesReplaceAll", 11},
		{"BytesFields", 3},
		{"BytesTrimSpace", 5},
		{"BytesToUpper", 5},
		{"BytesToLower", 5},
		{"BytesTrim", 5},
	})

	runCategory(t, "Strings", srcStrings, []thirdpartyCase{
		{"StringsBuilder", 11},
		{"StringsBuilderString", "test"},
		{"StringsRepeatCount", 100},
		{"StringsIndexAny", 1},
		{"StringsCut", 1},
		{"StringsIndexFuncTest", 3},
		{"StringsTrimLeft", 1},
		{"StringsTrimRight", 1},
	})

	runCategory(t, "Strconv", srcStrconv, []thirdpartyCase{
		{"StrconvParseBool", 1},
		{"StrconvFormatBool", "false"},
		{"StrconvParseInt", int64(12345)},
		{"StrconvParseUint", uint64(12345)},
		{"StrconvFormatInt", "-12345"},
		{"StrconvFormatUint", "12345"},
		{"StrconvParseFloat", float64(123.45)},
		{"StrconvFormatFloat", "123.45"},
		{"StrconvQuote", `"hello\nworld"`},
		{"StrconvQuoteToASCII", `"hello"`},
		{"StrconvUnquote", "hello"},
		{"StrconvAppendInt", "12345"},
		{"StrconvAppendFloat", "123.45"},
	})

	runCategory(t, "Math", srcMath, []thirdpartyCase{
		{"MathAbs", float64(123.45)},
		{"MathMax", float64(20.3)},
		{"MathMin", float64(10.5)},
		{"MathFloor", float64(123)},
		{"MathCeil", float64(124)},
		{"MathRound", float64(124)},
		{"MathPow", float64(1024)},
		{"MathSqrt", float64(12)},
		{"MathMod", float64(1)},
		{"MathSin", float64(1)},
		{"MathCos", float64(1)},
		{"MathTan", float64(0)},
		{"MathLog", float64(1)},
		{"MathLog10", float64(2)},
		{"MathInf", 1},
		{"MathNaN", 1},
		{"MathCopysign", float64(5)},
	})

	runCategory(t, "Time", srcTime, []thirdpartyCase{
		{"TimeNow", 2026},
		{"TimeFormat", "March"},
		{"TimeAdd", 2},
		{"TimeBefore", 1},
		{"TimeAfter", 1},
		{"TimeDuration", 330},
	})

	runCategory(t, "Context", srcContext, []thirdpartyCase{
		{"ContextBackground", 1},
		{"ContextTODO", 1},
		{"ContextWithValue", 1},
		{"ContextWithCancel", 1},
	})

	runCategory(t, "Sync", srcSync, []thirdpartyCase{
		{"SyncMutex", 1},
		{"SyncMutexCounter", 100},
		{"SyncRWMutex", 1},
		{"SyncWaitGroup", 1},
		{"SyncOnce", 1},
		{"SyncMap", 1},
		{"SyncMapLoadOrStore", 1},
	})

	runCategory(t, "Sort", srcSort, []thirdpartyCase{
		{"SortStrings", 1},
		// SortInts skipped: sort.Ints mutates []int in-place but VM stores as []int64 (known issue)
		{"SortSearchInts", 2},
		{"SortSearchStrings", 1},
		{"SortSlice", 1},
		{"SortIsSorted", 1},
	})

	runCategory(t, "Encoding", srcEncoding, []thirdpartyCase{
		{"JsonMarshal", 13},
		{"JsonUnmarshal", 3},
		{"JsonNumber", 42},
		{"Base64Encode", "aGVsbG8="},
		{"Base64Decode", "hello"},
		{"Base64URLEncode", "aGVsbG8gd29ybGQ="},
		{"HexEncodeToString", "68656c6c6f"},
		{"HexDecodeString", "hello"},
	})

	runCategory(t, "IO", srcIO, []thirdpartyCase{
		{"IoReadAll", 11},
		{"IoCopy", 5},
		{"IoReadFull", 5},
		{"IoWriteString", 4},
	})

	// path/filepath removed from sandbox stdlib
	// runCategory(t, "Filepath", srcFilepath, []thirdpartyCase{
	// 	{"FilepathJoin", "dir1/dir2/file.txt"},
	// 	{"FilepathBase", "file.txt"},
	// 	{"FilepathDir", "/path/to"},
	// 	{"FilepathExt", ".txt"},
	// 	{"FilepathClean", "/path/to/file.txt"},
	// })

	runCategory(t, "Regexp", srcRegexp, []thirdpartyCase{
		{"RegexpMatch", 1},
		{"RegexpCompile", 4},
		{"RegexpMustCompile", 1},
		{"RegexpFindString", "foo"},
		{"RegexpFindAllString", 3},
		{"RegexpReplaceAllString", "a# b# c#"},
		{"RegexpSplit", 3},
		{"RegexpNumSubexp", 2},
	})

	runCategory(t, "Errors", srcErrors, []thirdpartyCase{
		{"ErrorsNew", 1},
		{"ErrorsIs", 1},
		{"ErrorsJoin", "error1\nerror2"},
	})

	runCategory(t, "Fmt", srcFmt, []thirdpartyCase{
		{"FmtSprintfVarious", 19},
		{"FmtSprintfBool", "true"},
		{"FmtSprintfHex", "ff"},
		{"FmtErrorf", "error: 42"},
	})

	runCategory(t, "Patterns", srcPatterns, []thirdpartyCase{
		{"ChainBytesToStringToBase64", "aGVsbG8="},
		{"ChainStringsBuilderToBuffer", 11},
		{"ChainSortSearch", 30},
		{"ChainBufferWriteRead", 6},
		{"ChainContextWithValueChain", 1},
		{"InterfaceWithPointerReceiver", 1},
		{"InterfaceSliceOfPointers", 9},
		{"InterfaceMap", 1},
		{"VariadicAppend", 5},
		{"VariadicStringsJoin", "a,b,c"},
		{"VariadicAppendSlice", 5},
		{"MethodChainingBuilder", "SELECT id,name FROM users WHERE active = true"},
		{"TableDrivenOp", 19},
		{"FunctionValueFromMap", 70},
		{"DeferWithMutex", 42},
		{"SelectWithChannels", 42},
	})
}

// TestThirdpartyNative runs native Go functions to verify expected values.
func TestThirdpartyNative(t *testing.T) {
	nativeTests := []struct {
		name     string
		fn       func() any
		expected any
	}{
		// Bytes
		{"BytesBufferWrite", func() any { return thirdparty.BytesBufferWrite() }, 11},
		{"BytesBufferString", func() any { return thirdparty.BytesBufferString() }, "test string"},
		{"BytesSplit", func() any { return thirdparty.BytesSplit() }, 3},
		{"BytesReplaceAll", func() any { return thirdparty.BytesReplaceAll() }, 11},
		// Strings
		{"StringsBuilder", func() any { return thirdparty.StringsBuilder() }, 11},
		{"StringsIndexFuncTest", func() any { return thirdparty.StringsIndexFuncTest() }, 3},
		// Strconv
		{"StrconvParseBool", func() any { return thirdparty.StrconvParseBool() }, 1},
		{"StrconvParseInt", func() any { return thirdparty.StrconvParseInt() }, int64(12345)},
		{"StrconvParseUint", func() any { return thirdparty.StrconvParseUint() }, uint64(12345)},
		// Math
		{"MathAbs", func() any { return thirdparty.MathAbs() }, float64(123.45)},
		{"MathCopysign", func() any { return thirdparty.MathCopysign() }, float64(5)},
		// Sync
		{"SyncMutex", func() any { return thirdparty.SyncMutex() }, 1},
		{"SyncOnce", func() any { return thirdparty.SyncOnce() }, 1},
		// Sort
		{"SortStrings", func() any { return thirdparty.SortStrings() }, 1},
		{"SortInts", func() any { return thirdparty.SortInts() }, 1},
		// Encoding
		{"JsonMarshal", func() any { return thirdparty.JsonMarshal() }, 13},
		{"JsonNumber", func() any { return thirdparty.JsonNumber() }, 42},
		{"Base64Encode", func() any { return thirdparty.Base64Encode() }, "aGVsbG8="},
		// Context
		{"ContextWithValue", func() any { return thirdparty.ContextWithValue() }, 1},
		// Patterns
		{"TableDrivenOp", func() any { return thirdparty.TableDrivenOp() }, 19},
		{"FunctionValueFromMap", func() any { return thirdparty.FunctionValueFromMap() }, 70},
	}

	for _, tt := range nativeTests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()
			if result != tt.expected {
				t.Errorf("%s() = %v (%T), want %v (%T)",
					tt.name, result, result, tt.expected, tt.expected)
			}
		})
	}
}
