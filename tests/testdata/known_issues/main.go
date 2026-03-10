package known_issues

// BytesToStringHi converts []byte{104, 105} to a string.
// Expected: "hi"
func BytesToStringHi() string {
	return string([]byte{104, 105})
}

// BytesToStringGo converts []byte{71, 111} to a string.
// Expected: "Go"
func BytesToStringGo() string {
	return string([]byte{71, 111})
}

// BytesToStringSingle converts []byte{65} to a string.
// Expected: "A"
func BytesToStringSingle() string {
	return string([]byte{65})
}

// BytesToStringEmpty converts an empty byte slice to a string.
// Expected: ""
func BytesToStringEmpty() string {
	return string([]byte{})
}

// counter is a package-level type so we can define methods on it.
type counter struct{ n int }

func (c *counter) inc() { c.n++ }

// PointerReceiverMutation calls a pointer-receiver method twice and returns the result.
// Expected: 2
func PointerReceiverMutation() int {
	c := &counter{}
	c.inc()
	c.inc()
	return c.n
}

// box is a helper for the Set/Get test.
type box struct{ val int }

func (b *box) set(v int) { b.val = v }
func (b box) get() int   { return b.val }

// PointerReceiverSetGet calls a pointer-receiver setter then a value-receiver getter.
// Expected: 99
func PointerReceiverSetGet() int {
	b := &box{}
	b.set(99)
	return b.get()
}

var initVal int
var initRegistry []string

func init() {
	initVal = 42
	initRegistry = append(initRegistry, "alpha")
	initRegistry = append(initRegistry, "beta")
}

// InitFuncResult returns the value set by init().
// Expected: 42
func InitFuncResult() int {
	return initVal
}

// InitSliceLen returns the length of the slice populated by init().
// Expected: 2
func InitSliceLen() int {
	return len(initRegistry)
}

// RangeStringRune sums rune values when ranging over an ASCII string.
// 'a'=97, 'b'=98, 'c'=99 -> 294
// Expected: 294
func RangeStringRune() int {
	sum := 0
	for _, r := range "abc" {
		sum += int(r)
	}
	return sum
}

// RangeStringIndex sums byte indices when ranging over a 3-char ASCII string.
// indices: 0, 1, 2 -> 3
// Expected: 3
func RangeStringIndex() int {
	sum := 0
	for i := range "xyz" {
		sum += i
	}
	return sum
}

// RangeStringMultibyte sums rune values for a multi-byte UTF-8 string.
// '中'=0x4E2D=20013, '文'=0x6587=25991 -> 46004
// Expected: 46004
func RangeStringMultibyte() int {
	sum := 0
	for _, r := range "中文" {
		sum += int(r)
	}
	return sum
}
