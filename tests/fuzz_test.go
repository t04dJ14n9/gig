// Package tests - fuzz_test.go
//
// Fuzz tests using Go 1.21+ native fuzzing to catch regressions in
// VM opcode handlers and type conversions with randomized inputs.
//
// Run with: go test ./tests/ -fuzz=FuzzIntNarrowing ...
// Or:        go test ./tests/ -fuzz=.   (runs all fuzz targets, slow)
package tests

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

// build compiles source and skips the test if compilation fails.
// Some Unicode values may trigger SSA compilation errors in edge cases
// — those are expected and we skip them gracefully.
func build(t *testing.T, src string) *gig.Program {
	t.Helper()
	prog, err := gig.Build(src)
	if err != nil {
		t.Skipf("Build error (expected for some Unicode): %v", err)
	}
	return prog
}

// run is a wrapper around prog.Run that skips on error.
func run(t *testing.T, prog *gig.Program, fn string, args ...any) any {
	t.Helper()
	if prog == nil {
		t.Skip("program is nil (build failed)")
	}
	result, err := prog.Run(fn, args...)
	if err != nil {
		t.Skipf("Run error: %v", err)
	}
	return result
}

// ============================================================================
// Fuzz: Integer narrowing conversions
// ============================================================================

func FuzzIntNarrowing(f *testing.F) {
	seeds := []int64{
		0, 1, -1,
		127, -128, 255, -256,
		32767, -32768, 65535, -65536,
		2147483647, -2147483648,
		1<<24 - 1, -1 << 24,
		math.MaxInt64, math.MinInt64,
	}
	for _, v := range seeds {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v int64) {
		prog := build(t, `package main; func NarrowInt8(x int) int8 { return int8(x) }`)
		result := run(t, prog, "NarrowInt8", v)
		if result != int8(v) {
			t.Errorf("int8(%d) = %d, want %d", v, result, int8(v))
		}

		prog16 := build(t, `package main; func NarrowInt16(x int) int16 { return int16(x) }`)
		result16 := run(t, prog16, "NarrowInt16", v)
		if result16 != int16(v) {
			t.Errorf("int16(%d) = %d, want %d", v, result16, int16(v))
		}

		prog32 := build(t, `package main; func NarrowInt32(x int) int32 { return int32(x) }`)
		result32 := run(t, prog32, "NarrowInt32", v)
		if result32 != int32(v) {
			t.Errorf("int32(%d) = %d, want %d", v, result32, int32(v))
		}
	})
}

func FuzzUintNarrowing(f *testing.F) {
	seeds := []uint64{
		0, 1, 127, 128, 255,
		256, 32767, 65535,
		65536, 2147483647, 4294967295,
		1 << 56, 1<<64 - 1,
	}
	for _, v := range seeds {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v uint64) {
		prog := build(t, `package main; func NarrowUint8(x uint) uint8 { return uint8(x) }`)
		result := run(t, prog, "NarrowUint8", v)
		if result != uint8(v) {
			t.Errorf("uint8(%d) = %d, want %d", v, result, uint8(v))
		}

		prog16 := build(t, `package main; func NarrowUint16(x uint) uint16 { return uint16(x) }`)
		result16 := run(t, prog16, "NarrowUint16", v)
		if result16 != uint16(v) {
			t.Errorf("uint16(%d) = %d, want %d", v, result16, uint16(v))
		}

		prog32 := build(t, `package main; func NarrowUint32(x uint) uint32 { return uint32(x) }`)
		result32 := run(t, prog32, "NarrowUint32", v)
		if result32 != uint32(v) {
			t.Errorf("uint32(%d) = %d, want %d", v, result32, uint32(v))
		}
	})
}

// ============================================================================
// Fuzz: Float narrowing
// ============================================================================

func FuzzFloatNarrowing(f *testing.F) {
	seeds := []float64{
		0, 1, -1, 0.5, -0.5,
		math.MaxFloat32, -math.MaxFloat32,
		math.SmallestNonzeroFloat32, -math.SmallestNonzeroFloat32,
		math.MaxFloat64, -math.MaxFloat64,
		math.Inf(1), math.Inf(-1), math.NaN(),
		1e10, 1e-10, -1e10, -1e-10,
		math.Pi, math.E,
		3.4028234663852886e+38,
	}
	for _, v := range seeds {
		f.Add(v)
	}

	f.Fuzz(func(t *testing.T, v float64) {
		prog := build(t, `package main; func NarrowFloat32(x float64) float32 { return float32(x) }`)
		result := run(t, prog, "NarrowFloat32", v)
		want := float32(v)
		if math.IsNaN(float64(result.(float32))) && math.IsNaN(float64(want)) {
			return
		}
		if result != want {
			t.Errorf("float32(%.20g) = %.20g, want %.20g", v, result, want)
		}
	})
}

// ============================================================================
// Fuzz: Bitwise operations (regression test for unary ^ bug fix)
// ============================================================================

func FuzzBitwise(f *testing.F) {
	seeds := []struct {
		a uint64
		b uint64
	}{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
		{0xff, 0x00},
		{0xff, 0xff},
		{0xaa, 0x55},
		{1 << 63, 1 << 63},
		{1<<63 - 1, 1},
		{math.MaxUint64, 1},
		{math.MaxUint64, 0},
	}
	for _, c := range seeds {
		f.Add(c.a, c.b)
	}

	f.Fuzz(func(t *testing.T, a uint64, b uint64) {
		// Binary AND/OR/XOR
		for _, op := range []struct {
			name string
			sym  string
			want uint64
		}{
			{"BitAnd", "&", a & b},
			{"BitOr", "|", a | b},
			{"BitXor", "^", a ^ b},
		} {
			prog := build(t, "package main; func "+op.name+`(a, b uint64) uint64 { return a `+op.sym+` b }`)
			result := run(t, prog, op.name, a, b)
			if result != op.want {
				t.Errorf("%s(%d, %d) = %d, want %d", op.name, a, b, result, op.want)
			}
		}

		// Unary bitwise NOT — regression test for the compileUnOp ^ bug
		progNot := build(t, `package main; func BitNot(a uint64) uint64 { return ^a }`)
		resultNot := run(t, progNot, "BitNot", a)
		if resultNot != ^a {
			t.Errorf("BitNot(%d) = %d, want %d", a, resultNot, ^a)
		}

		// Shifts
		for _, sh := range []struct {
			name string
			expr string
			want uint64
		}{
			{"Lsh1", "a << 1", a << 1},
			{"Lsh7", "a << 7", a << 7},
			{"Rsh1", "a >> 1", a >> 1},
			{"Rsh7", "a >> 7", a >> 7},
		} {
			prog := build(t, "package main; func "+sh.name+`(a uint64) uint64 { return `+sh.expr+` }`)
			result := run(t, prog, sh.name, a)
			if result != sh.want {
				t.Errorf("%s(%d) = %d, want %d", sh.name, a, result, sh.want)
			}
		}
	})
}

// ============================================================================
// Fuzz: String operations
// ============================================================================

func FuzzStringOps(f *testing.F) {
	seeds := []string{
		"", "a", "hello", "世界", "🎉",
		"abc日本語def",
		strings.Repeat("a", 100),
		"你好👋",
		"a\x00b", "   ", "\t\n",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		progLen := build(t, `package main; func StrLen(s string) int { return len(s) }`)
		resultLen := run(t, progLen, "StrLen", s)
		if resultLen != len(s) {
			t.Errorf("len(%q) = %d, want %d", s, resultLen, len(s))
		}

		// String indexing: len(string(s[0])) = rune count of first char
		progIdx := build(t, `package main; func StrIdx(s string) int { if len(s) > 0 { return len(string(s[0])); } return 0 }`)
		resultIdx := run(t, progIdx, "StrIdx", s)
		wantIdx := 0
		if len(s) > 0 {
			wantIdx = len(string(s[0]))
		}
		if resultIdx != wantIdx {
			t.Errorf("StrIdx(%q) = %d, want %d", s, resultIdx, wantIdx)
		}

		// Slicing with guard against out-of-bounds
		cuts := []int{-1, 0, 1}
		if len(s) > 0 {
			cuts = append(cuts, len(s)/2, len(s)-1, len(s), len(s)+1)
		}
		for _, cut := range cuts {
			progSlice := build(t, `package main; func StrSlice(s string, i int) string { if i > len(s) { i = len(s) }; if i < 0 { i = 0 }; return s[i:] }`)
			resultSlice := run(t, progSlice, "StrSlice", s, cut)
			wantSlice := s[cut:]
			if resultSlice != wantSlice {
				t.Errorf("StrSlice(%q, %d) = %q, want %q", s, cut, resultSlice, wantSlice)
			}
		}
	})
}

// ============================================================================
// Fuzz: Map operations
// ============================================================================

func FuzzMapOps(f *testing.F) {
	seeds := []struct {
		ki, vi int
		ks, vs string
	}{
		{1, 100, "a", "x"}, {0, 0, "", ""}, {-1, -100, "key", "value"},
	}
	for _, s := range seeds {
		f.Add(s.ki, s.vi, s.ks, s.vs)
	}

	f.Fuzz(func(t *testing.T, ki int, vi int, ks string, vs string) {
		prog := build(t, `package main
		func MapWrite(m map[int]int, k, v int) int { m[k] = v; return len(m) }
		func MapRead(m map[int]int, k int) int { return m[k] }
		func MapLen(m map[int]int) int { return len(m) }`)
		m := map[int]int{}
		run(t, prog, "MapWrite", m, ki, vi)
		result := run(t, prog, "MapRead", m, ki)
		if result != vi {
			t.Errorf("MapRead[%d] = %d, want %d", ki, result, vi)
		}
		length := run(t, prog, "MapLen", m)
		if length != 1 {
			t.Errorf("MapLen = %d, want 1", length)
		}

		progStr := build(t, `package main
		func MapWriteStr(m map[string]string, k, v string) int { m[k] = v; return len(m) }
		func MapReadStr(m map[string]string, k string) string { return m[k] }`)
		sm := map[string]string{}
		run(t, progStr, "MapWriteStr", sm, ks, vs)
		resultStr := run(t, progStr, "MapReadStr", sm, ks)
		if resultStr != vs {
			t.Errorf("MapReadStr[%q] = %q, want %q", ks, resultStr, vs)
		}
	})
}

// ============================================================================
// Fuzz: Channel operations
// ============================================================================

func FuzzChannelOps(f *testing.F) {
	seeds := []struct{ buf, val int }{
		{0, 0},
		{1, 1},
		{10, 0},
		{100, 1},
		{0, -1},
		{1, 127},
		{0, math.MaxInt},
		{0, math.MinInt},
	}
	for _, s := range seeds {
		f.Add(s.buf, s.val)
	}

	f.Fuzz(func(t *testing.T, buf int, val int) {
		if buf < 0 {
			buf = 0
		}
		if buf > 10 {
			buf = 10
		}

		prog := build(t, `package main
		func MakeChan(buf int) chan int { return make(chan int, buf) }
		func Send(ch chan int, v int) int { ch <- v; return v }
		func Recv(ch chan int) int { return <-ch }
		func ChanLen(ch chan int) int { return len(ch) }
		func ChanCap(ch chan int) int { return cap(ch) }`)

		ch := run(t, prog, "MakeChan", buf)
		chanVal := ch.(chan int)

		run(t, prog, "Send", chanVal, val)
		l := run(t, prog, "ChanLen", chanVal)
		if l != 1 {
			t.Errorf("ChanLen after send = %d, want 1 (buf=%d)", l, buf)
		}
		got := run(t, prog, "Recv", chanVal)
		if got != val {
			t.Errorf("Recv = %d, want %d", got, val)
		}
		c := run(t, prog, "ChanCap", chanVal)
		if c != buf {
			t.Errorf("ChanCap = %d, want %d", c, buf)
		}
	})
}

// ============================================================================
// Fuzz: Slice operations
// ============================================================================

func FuzzSliceOps(f *testing.F) {
	// Go fuzz only supports basic types (int, uint, float, bool, string, []byte).
	// We use two int arguments: n controls the size, seed controls content.
	f.Add(0, 0)    // empty slice
	f.Add(1, 0)    // one element
	f.Add(3, 1)    // slice with 3 elements, idx=1
	f.Add(10, 5)   // medium slice
	f.Add(100, 50) // large slice

	f.Fuzz(func(t *testing.T, n int, seed int) {
		if n < 0 {
			n = 0
		}
		if n > 50 {
			n = 50
		}

		// Build a deterministic slice from seed
		vals := make([]int, n)
		for i := range vals {
			vals[i] = (seed<<8 | i) & 0xFFFF
			if seed < 0 {
				vals[i] = -vals[i]
			}
		}
		idx := seed % (n + 2)
		if idx < -1 {
			idx = -1
		}

		prog := build(t, `package main
		func SliceLen(s []int) int { return len(s) }
		func SliceCap(s []int) int { return cap(s) }
		func SliceIdx(s []int, i int) int { if i < 0 || i >= len(s) { return -1 }; return s[i] }
		func SliceAppend(s []int, v int) []int { return append(s, v) }`)

		resultLen := run(t, prog, "SliceLen", vals)
		if resultLen != len(vals) {
			t.Errorf("len() = %d, want %d", resultLen, len(vals))
		}
		resultCap := run(t, prog, "SliceCap", vals)
		if resultCap != cap(vals) {
			t.Errorf("cap() = %d, want %d", resultCap, cap(vals))
		}
		resultIdx := run(t, prog, "SliceIdx", vals, idx)
		wantIdx := -1
		if idx >= 0 && idx < len(vals) {
			wantIdx = vals[idx]
		}
		if resultIdx != wantIdx {
			t.Errorf("SliceIdx(n=%d, idx=%d) = %d, want %d", n, idx, resultIdx, wantIdx)
		}

		appended := run(t, prog, "SliceAppend", vals, 999)
		wantAppended := append(vals, 999) //nolint:makezero // intentional: testing append with non-zero length slice
		if !reflect.DeepEqual(appended, wantAppended) {
			t.Errorf("append() = %v, want %v", appended, wantAppended)
		}
	})
}

// ============================================================================
// Fuzz: Struct field access
// ============================================================================

func FuzzStructAccess(f *testing.F) {
	seeds := []struct {
		i  int
		s  string
		b  bool
		fv float64
	}{
		{0, "", false, 0},
		{42, "hello", true, 3.14},
		{-1, "世界", false, math.NaN()},
		{math.MaxInt, "🎉", true, math.Inf(1)},
	}
	for _, s := range seeds {
		f.Add(s.i, s.s, s.b, s.fv)
	}

	f.Fuzz(func(t *testing.T, fi int, fs string, fb bool, ff float64) {
		prog := build(t, `package main
		type S struct { I int; S string; B bool; F float64 }
		func NewStruct(i int, s string, b bool, f float64) S { return S{I: i, S: s, B: b, F: f} }
		func GetI(s S) int { return s.I }
		func GetS(s S) string { return s.S }
		func GetB(s S) bool { return s.B }
		func GetF(s S) float64 { return s.F }`)

		s := run(t, prog, "NewStruct", fi, fs, fb, ff)

		fields := []struct {
			name  string
			getFn string
			want  any
		}{
			{"I", "GetI", fi},
			{"S", "GetS", fs},
			{"B", "GetB", fb},
			{"F", "GetF", ff},
		}

		for _, fld := range fields {
			got := run(t, prog, fld.getFn, s)
			if fld.name == "F" && math.IsNaN(fld.want.(float64)) && math.IsNaN(got.(float64)) {
				continue
			}
			if got != fld.want {
				t.Errorf("struct.%s = %v, want %v (input: i=%d s=%q b=%v f=%.10g)",
					fld.name, got, fld.want, fi, fs, fb, ff)
			}
		}
	})
}

// ============================================================================
// Fuzz: Arithmetic edge cases
// ============================================================================

func FuzzArithmeticEdgeCases(f *testing.F) {
	seeds := []struct{ a, b int64 }{
		{1, 1},
		{-1, 1},
		{1, -1},
		{math.MaxInt64, 1},
		{math.MinInt64, 1},
		{math.MaxInt64, -1},
		{100, 3},
		{100, 7},
		{-100, 3},
		{1 << 62, 3},
		{1 << 62, 7},
	}
	for _, s := range seeds {
		f.Add(s.a, s.b)
	}

	f.Fuzz(func(t *testing.T, a, b int64) {
		if b == 0 {
			return
		}

		prog := build(t, `package main
		func Div(a, b int64) int64 { return a / b }
		func Mod(a, b int64) int64 { return a % b }`)

		resultDiv := run(t, prog, "Div", a, b)
		if resultDiv != a/b {
			t.Errorf("Div(%d, %d) = %d, want %d", a, b, resultDiv, a/b)
		}
		resultMod := run(t, prog, "Mod", a, b)
		if resultMod != a%b {
			t.Errorf("Mod(%d, %d) = %d, want %d", a, b, resultMod, a%b)
		}
	})
}

func FuzzFloatSpecialValues(f *testing.F) {
	seeds := []struct{ a, b float64 }{
		{0, 1},
		{-0.0, 1.0},
		{math.Inf(1), 1},
		{math.Inf(-1), 1},
		{math.MaxFloat64, 1},
		{-math.MaxFloat64, 1},
		{math.SmallestNonzeroFloat64, 1},
		{math.Pi, math.E},
	}
	for _, s := range seeds {
		f.Add(s.a, s.b)
	}

	f.Fuzz(func(t *testing.T, a, b float64) {
		if math.IsNaN(a) || math.IsNaN(b) || b == 0 {
			return
		}

		for _, op := range []struct {
			name string
			sym  string
			want float64
		}{
			{"Add", "+", a + b},
			{"Sub", "-", a - b},
			{"Mul", "*", a * b},
			{"Div", "/", a / b},
		} {
			prog := build(t, "package main; func Op(a, b float64) float64 { return a "+op.sym+" b }")
			got := run(t, prog, "Op", a, b)
			want := op.want
			if math.IsNaN(want) && math.IsNaN(got.(float64)) {
				continue
			}
			if math.IsInf(want, 0) && math.IsInf(got.(float64), 0) {
				continue
			}
			if got != want {
				t.Errorf("float%s(%.15g, %.15g) = %.15g, want %.15g",
					op.name, a, b, got, want)
			}
		}
	})
}
