// Package tests — known_issues_test.go
//
// These tests definitively verify whether reported Gig limitations are real bugs.
// Each test FAILS if the bug is present, so CI catches regressions and the
// failure message tells you exactly what Gig returned vs. what Go requires.
//
// When a bug is fixed, the test will pass automatically — no changes needed here.
package tests

import (
	"testing"

	gig "git.woa.com/youngjin/gig"
)

// ── helper ────────────────────────────────────────────────────────────────────

// buildAndRun compiles source, calls funcName with no arguments, and returns the
// raw result. It fails the test immediately on Build or Run errors.
func buildAndRun(t *testing.T, source, funcName string) any {
	t.Helper()
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run(funcName)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	return result
}

// ── Issue 1: string([]byte{...}) conversion ───────────────────────────────────

// TestKnownIssue_BytesToString verifies that string([]byte{104, 105}) returns "hi".
//
// Go spec: string([]byte{104, 105}) must produce the UTF-8 string "hi" (h=104, i=105).
// Reported Gig behaviour: returns "[104 105]" (fmt-style slice rendering).
func TestKnownIssue_BytesToString(t *testing.T) {
	got := buildAndRun(t, `
func Compute() string {
	return string([]byte{104, 105})
}
`, "Compute")

	want := "hi"
	s, ok := got.(string)
	if !ok {
		t.Fatalf("string([]byte{104,105}): expected string result, got %T: %v", got, got)
	}
	if s != want {
		t.Errorf("string([]byte{104,105}): got %q, want %q", s, want)
	}
}

// TestKnownIssue_BytesToStringMulti checks several byte-to-string conversions to
// confirm the issue is general, not specific to {104, 105}.
func TestKnownIssue_BytesToStringMulti(t *testing.T) {
	cases := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "ascii_hi",
			source: `func Compute() string { return string([]byte{104, 105}) }`,
			want:   "hi",
		},
		{
			name:   "ascii_go",
			source: `func Compute() string { return string([]byte{71, 111}) }`,
			want:   "Go",
		},
		{
			name:   "empty",
			source: `func Compute() string { return string([]byte{}) }`,
			want:   "",
		},
		{
			name:   "single_byte",
			source: `func Compute() string { return string([]byte{65}) }`,
			want:   "A",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildAndRun(t, tc.source, "Compute")
			s, ok := got.(string)
			if !ok {
				t.Fatalf("expected string, got %T: %v", got, got)
			}
			if s != tc.want {
				t.Errorf("got %q, want %q", s, tc.want)
			}
		})
	}
}

// ── Issue 2: Pointer-receiver method mutation ─────────────────────────────────

// TestKnownIssue_PointerReceiverMutation verifies that a pointer-receiver method
// that increments a field actually persists the mutation.
//
// Go spec: (c *Counter).Inc() must modify the value c points to; after two calls
// c.n must equal 2.
// Reported Gig behaviour: the mutation does not persist, c.n remains 0.
func TestKnownIssue_PointerReceiverMutation(t *testing.T) {
	got := buildAndRun(t, `
type Counter struct{ n int }
func (c *Counter) Inc() { c.n++ }
func Compute() int {
	c := &Counter{}
	c.Inc()
	c.Inc()
	return c.n
}
`, "Compute")

	if got == nil {
		t.Fatal("pointer-receiver mutation: got nil, want 2")
	}
	n := toInt64(t, got)
	if n != 2 {
		t.Errorf("pointer-receiver mutation: got %d, want 2", n)
	}
}

// TestKnownIssue_PointerReceiverMutationReturnValue verifies that calling a
// pointer-receiver method and then reading a field via value receiver is consistent.
func TestKnownIssue_PointerReceiverMutationReturnValue(t *testing.T) {
	got := buildAndRun(t, `
type Box struct{ val int }
func (b *Box) Set(v int) { b.val = v }
func (b Box) Get() int   { return b.val }
func Compute() int {
	b := &Box{}
	b.Set(99)
	return b.Get()
}
`, "Compute")

	if got == nil {
		t.Fatal("pointer-receiver Set+Get: got nil, want 99")
	}
	n := toInt64(t, got)
	if n != 99 {
		t.Errorf("pointer-receiver Set+Get: got %d, want 99", n)
	}
}

// ── Issue 3: init() execution ─────────────────────────────────────────────────

// TestKnownIssue_InitFuncExecuted verifies that package-level init() runs before
// the entry function is called.
//
// Go spec: init() is always executed once before any other code in the package.
// Reported Gig behaviour: init() is not executed; package variable stays at zero.
func TestKnownIssue_InitFuncExecuted(t *testing.T) {
	got := buildAndRun(t, `
var initVal int

func init() { initVal = 42 }

func Compute() int { return initVal }
`, "Compute")

	if got == nil {
		t.Fatal("init(): got nil result, want 42 (init() not executed)")
	}
	n := toInt64(t, got)
	if n != 42 {
		t.Errorf("init(): got %d, want 42 (init() not executed)", n)
	}
}

// TestKnownIssue_InitFuncSideEffect verifies that init() can perform an
// operation whose result is visible to the entry function.
func TestKnownIssue_InitFuncSideEffect(t *testing.T) {
	got := buildAndRun(t, `
var registry []string

func init() {
	registry = append(registry, "alpha")
	registry = append(registry, "beta")
}

func Compute() int { return len(registry) }
`, "Compute")

	if got == nil {
		t.Fatal("init() slice append: got nil result, want 2")
	}
	n := toInt64(t, got)
	if n != 2 {
		t.Errorf("init() slice append: got %d, want 2", n)
	}
}

// ── Issue 4: range-over-string rune values ────────────────────────────────────

// TestKnownIssue_RangeStringRuneValue verifies that the second variable in
// "for _, r := range str" receives the actual Unicode code point, not 0.
//
// Go spec: range over a string yields (byteIndex int, runeValue rune).
// 'a'=97, 'b'=98, 'c'=99 → sum should be 294.
// Reported Gig behaviour: rune variable is always 0, so sum = 0.
func TestKnownIssue_RangeStringRuneValue(t *testing.T) {
	got := buildAndRun(t, `
func Compute() int {
	sum := 0
	for _, r := range "abc" { sum += int(r) }
	return sum
}
`, "Compute")

	if got == nil {
		t.Fatal("range-over-string rune: got nil, want 294")
	}
	n := toInt64(t, got)
	const want = int64(97 + 98 + 99) // 'a'+'b'+'c' = 294
	if n != want {
		t.Errorf("range-over-string rune: got %d, want %d ('a'+'b'+'c')", n, want)
	}
}

// TestKnownIssue_RangeStringIndexValue verifies the byte index produced by
// range-over-string. For ASCII strings the index should be 0, 1, 2, ...
// Sum of indices over "xyz" (3 chars) should be 0+1+2 = 3.
func TestKnownIssue_RangeStringIndexValue(t *testing.T) {
	got := buildAndRun(t, `
func Compute() int {
	sum := 0
	for i := range "xyz" { sum += i }
	return sum
}
`, "Compute")

	if got == nil {
		t.Fatal("range-over-string index: got nil, want 3")
	}
	n := toInt64(t, got)
	if n != 3 {
		t.Errorf("range-over-string index: got %d, want 3 (0+1+2)", n)
	}
}

// TestKnownIssue_RangeStringMultibyte verifies that multi-byte UTF-8 codepoints
// are decoded correctly: '中'=0x4E2D=20013, '文'=0x6587=25991; sum=46004.
func TestKnownIssue_RangeStringMultibyte(t *testing.T) {
	// Build the test string from rune literals to avoid gosmopolitan lint on Han script.
	// string([]rune{0x4E2D, 0x6587}) == "中文"
	zhStr := string([]rune{0x4E2D, 0x6587})
	src := `func Compute() int {
	sum := 0
	for _, r := range "` + zhStr + `" { sum += int(r) }
	return sum
}`
	got := buildAndRun(t, src, "Compute")

	if got == nil {
		t.Fatal("range-over-string multibyte rune: got nil, want 46004")
	}
	n := toInt64(t, got)
	const want = int64(0x4E2D + 0x6587) // 20013 + 25991 = 46004
	if n != want {
		t.Errorf("range-over-string multibyte rune: got %d, want %d", n, want)
	}
}
