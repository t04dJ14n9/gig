// Package tests — pass_struct_test.go
//
// Tests passing real Go structs from the host into gig-interpreted code.
// The interpreted code receives the struct as a parameter, calls methods on it,
// and returns results. This exercises the full cross-boundary interop path:
//
//	host creates *bytes.Buffer → passes via prog.Run → gig calls .Write/.String/etc → returns result
package tests

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
)

// ---------------------------------------------------------------------------
// 1) bytes.Buffer: write into a host-created buffer from interpreted code
// ---------------------------------------------------------------------------

func TestPassBytesBuffer_WriteAndRead(t *testing.T) {
	code := `package main
import "bytes"
// Accept a *bytes.Buffer, write to it, read it back.
func Test(buf *bytes.Buffer) string {
	buf.WriteString("hello ")
	buf.WriteString("from gig")
	return buf.String()
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.WriteString("prefix: ")

	result, err := prog.Run("Test", buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	// The buffer should have host-written prefix + gig-written content
	want := "prefix: hello from gig"
	if result != want {
		t.Errorf("got %q, want %q", result, want)
	}
	// Also verify the host buffer was mutated in-place
	if buf.String() != want {
		t.Errorf("host buffer: got %q, want %q", buf.String(), want)
	}
}

func TestPassBytesBuffer_Len(t *testing.T) {
	code := `package main
import "bytes"
func Test(buf *bytes.Buffer) int {
	return buf.Len()
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	buf := bytes.NewBufferString("12345")
	result, err := prog.Run("Test", buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result != 5 {
		t.Errorf("got %v, want 5", result)
	}
}

func TestPassBytesBuffer_Reset(t *testing.T) {
	code := `package main
import "bytes"
func Test(buf *bytes.Buffer) string {
	buf.Reset()
	buf.WriteString("new content")
	return buf.String()
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	buf := bytes.NewBufferString("old stuff")
	result, err := prog.Run("Test", buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result != "new content" {
		t.Errorf("got %q, want %q", result, "new content")
	}
}

// ---------------------------------------------------------------------------
// 2) strings.Builder: pass host builder, gig writes to it
// ---------------------------------------------------------------------------

func TestPassStringsBuilder(t *testing.T) {
	code := `package main
import "strings"
func Test(b *strings.Builder) string {
	b.WriteString("world")
	return b.String()
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	var b strings.Builder
	b.WriteString("hello ")

	result, err := prog.Run("Test", &b)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result != "hello world" {
		t.Errorf("got %q, want %q", result, "hello world")
	}
}

// ---------------------------------------------------------------------------
// 3) strings.Reader: pass host reader, gig reads from it
// ---------------------------------------------------------------------------

func TestPassStringsReader(t *testing.T) {
	code := `package main
import "strings"
func Test(r *strings.Reader) int {
	return r.Len()
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	r := strings.NewReader("hello gig")
	result, err := prog.Run("Test", r)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result != 9 {
		t.Errorf("got %v, want 9", result)
	}
}

// ---------------------------------------------------------------------------
// 4) json.Encoder: pass host-created encoder, gig writes JSON to host buffer
// ---------------------------------------------------------------------------

func TestPassJsonEncoder(t *testing.T) {
	code := `package main
import "encoding/json"
func Test(enc *json.Encoder) int {
	data := map[string]int{"x": 42}
	err := enc.Encode(data)
	if err != nil {
		return -1
	}
	return 1
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	result, err := prog.Run("Test", enc)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result != 1 {
		t.Errorf("got %v, want 1", result)
	}
	// The host buffer should have the JSON
	var decoded map[string]int
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("decode host buffer: %v", err)
	}
	if decoded["x"] != 42 {
		t.Errorf("decoded x=%v, want 42", decoded["x"])
	}
}

// ---------------------------------------------------------------------------
// 5) json.Decoder: pass host-created decoder, gig decodes JSON
// ---------------------------------------------------------------------------

func TestPassJsonDecoder(t *testing.T) {
	code := `package main
import "encoding/json"
func Test(dec *json.Decoder) int {
	var data map[string]int
	err := dec.Decode(&data)
	if err != nil {
		return -1
	}
	return data["value"]
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	r := strings.NewReader(`{"value": 99}`)
	dec := json.NewDecoder(r)
	result, err := prog.Run("Test", dec)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result != 99 {
		t.Errorf("got %v, want 99", result)
	}
}

// ---------------------------------------------------------------------------
// 6) csv.Writer: pass host-created CSV writer, gig writes rows
// ---------------------------------------------------------------------------

func TestPassCSVWriter(t *testing.T) {
	code := `package main
import "encoding/csv"
func Test(w *csv.Writer) int {
	w.Write([]string{"name", "age"})
	w.Write([]string{"Alice", "30"})
	w.Flush()
	if w.Error() != nil {
		return -1
	}
	return 1
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	result, err := prog.Run("Test", w)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result != 1 {
		t.Errorf("got %v, want 1", result)
	}
	want := "name,age\nAlice,30\n"
	if buf.String() != want {
		t.Errorf("csv output:\ngot  %q\nwant %q", buf.String(), want)
	}
}

// ---------------------------------------------------------------------------
// 7) gzip.Writer: pass host-created gzip writer, gig compresses data
//    This is a critical test because gzip.Writer.Write was the poster child
//    for Bug 8 (Writer.Write collision across 10+ packages).
// ---------------------------------------------------------------------------

func TestPassGzipWriter(t *testing.T) {
	code := `package main
import "compress/gzip"
func Test(w *gzip.Writer) int {
	n, err := w.Write([]byte("compressed by gig"))
	if err != nil {
		return -1
	}
	err = w.Close()
	if err != nil {
		return -2
	}
	return n
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	result, err := prog.Run("Test", w)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result != 17 { // len("compressed by gig") == 17
		t.Errorf("got %v, want 17", result)
	}

	// Verify the host buffer contains valid gzip data
	r, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("gzip.NewReader: %v", err)
	}
	var decoded bytes.Buffer
	decoded.ReadFrom(r)
	r.Close()
	if decoded.String() != "compressed by gig" {
		t.Errorf("decompressed: %q, want %q", decoded.String(), "compressed by gig")
	}
}

// ---------------------------------------------------------------------------
// 8) Multiple struct types in one call: pass both a Buffer and an Encoder
// ---------------------------------------------------------------------------

func TestPassMultipleStructs(t *testing.T) {
	code := `package main
import (
	"bytes"
	"encoding/json"
)
func Test(buf *bytes.Buffer, enc *json.Encoder) string {
	buf.WriteString("raw: ")
	enc.Encode(map[string]string{"k": "v"})
	return buf.String()
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf) // encoder writes to the same buffer
	result, err := prog.Run("Test", buf, enc)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	got, ok := result.(string)
	if !ok {
		t.Fatalf("result type: %T, want string", result)
	}
	// buf should contain "raw: " + JSON output
	if !strings.HasPrefix(got, "raw: ") {
		t.Errorf("got %q, want prefix %q", got, "raw: ")
	}
	if !strings.Contains(got, `"k"`) {
		t.Errorf("got %q, should contain JSON key \"k\"", got)
	}
}

// ---------------------------------------------------------------------------
// 9) Roundtrip: host creates struct, gig modifies it, host reads back
// ---------------------------------------------------------------------------

func TestPassAndMutate_BytesBuffer(t *testing.T) {
	code := `package main
import "bytes"
// Gig appends data and returns nothing meaningful — the mutation is the test.
func Test(buf *bytes.Buffer) int {
	buf.WriteString(" gig was here")
	return buf.Len()
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	buf := bytes.NewBufferString("host says hi,")
	result, err := prog.Run("Test", buf)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	// Verify host can read what gig wrote
	want := "host says hi, gig was here"
	if buf.String() != want {
		t.Errorf("host buffer: got %q, want %q", buf.String(), want)
	}
	if result != len(want) {
		t.Errorf("len: got %v, want %v", result, len(want))
	}
}
