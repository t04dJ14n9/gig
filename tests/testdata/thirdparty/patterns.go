package thirdparty

import (
	"bytes"
	"context"
	"encoding/base64"
	"sort"
	"strings"
	"sync"
)

// ============================================================================
// CHAINED CALLS - Multiple packages in one function
// ============================================================================

// ChainBytesToStringToBase64 tests chained encoding calls.
func ChainBytesToStringToBase64() string {
	data := []byte("hello")
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded
}

// ChainStringsBuilderToBuffer tests strings.Builder length.
func ChainStringsBuilderToBuffer() int {
	var b strings.Builder
	b.WriteString("hello")
	b.WriteString(" ")
	b.WriteString("world")
	return b.Len()
}

// ChainSortSearch tests combined sort and search.
func ChainSortSearch() int {
	s := []int{10, 20, 30, 40, 50}
	sort.Ints(s)
	idx := sort.SearchInts(s, 30)
	return s[idx]
}

// ChainBufferWriteRead tests bytes.Buffer write then read.
func ChainBufferWriteRead() int {
	buf := new(bytes.Buffer)
	buf.Write([]byte("hello world"))
	prefix := make([]byte, 5)
	buf.Read(prefix)
	return len(buf.Bytes())
}

// ChainContextWithValueChain tests nested context values.
func ChainContextWithValueChain() int {
	ctx := context.Background()
	ctx1 := context.WithValue(ctx, "level1", "value1")
	ctx2 := context.WithValue(ctx1, "level2", "value2")
	ctx3 := context.WithValue(ctx2, "level3", "value3")

	v1 := ctx3.Value("level1")
	v2 := ctx3.Value("level2")
	v3 := ctx3.Value("level3")

	if v1 == "value1" && v2 == "value2" && v3 == "value3" {
		return 1
	}
	return 0
}

// ============================================================================
// INTERFACE PATTERNS - Pointer receivers, slices, maps
// ============================================================================

// DataProcessor interface for testing pointer receiver patterns.
type DataProcessor interface {
	Process(data []byte) []byte
}

type uppercaseProcessor struct{}

func (u *uppercaseProcessor) Process(data []byte) []byte {
	return bytes.ToUpper(data)
}

type prefixProcessor struct {
	prefix string
}

func (p *prefixProcessor) Process(data []byte) []byte {
	return bytes.Join([][]byte{[]byte(p.prefix), data}, []byte{})
}

// InterfaceWithPointerReceiver tests interface with pointer receiver.
func InterfaceWithPointerReceiver() int {
	var processor DataProcessor = &uppercaseProcessor{}
	result := processor.Process([]byte("hello"))
	if string(result) == "HELLO" {
		return 1
	}
	return 0
}

// InterfaceSliceOfPointers tests slice of interface with pointer impls.
func InterfaceSliceOfPointers() int {
	processors := []DataProcessor{
		&uppercaseProcessor{},
		&prefixProcessor{prefix: "PREFIX_"},
	}
	sum := 0
	for _, p := range processors {
		result := p.Process([]byte("x"))
		sum += len(result)
	}
	return sum
}

// InterfaceMap tests map with interface values.
func InterfaceMap() int {
	m := make(map[string]DataProcessor)
	m["upper"] = &uppercaseProcessor{}
	m["prefix"] = &prefixProcessor{prefix: "X_"}

	result := m["upper"].Process([]byte("test"))
	if string(result) == "TEST" {
		return 1
	}
	return 0
}

// ============================================================================
// VARIADIC FUNCTIONS
// ============================================================================

// VariadicAppend tests append with variadic.
func VariadicAppend() int {
	s := []int{1, 2}
	s = append(s, 3, 4, 5)
	return len(s)
}

// VariadicStringsJoin tests strings.Join.
func VariadicStringsJoin() string {
	return strings.Join([]string{"a", "b", "c"}, ",")
}

// VariadicAppendSlice tests append with slice expansion.
func VariadicAppendSlice() int {
	s := []int{1, 2}
	t := []int{3, 4, 5}
	s = append(s, t...)
	return len(s)
}

// ============================================================================
// METHOD CHAINING - Builder pattern
// ============================================================================

// QueryBuilder simulates a builder pattern.
type QueryBuilder struct {
	parts []string
}

func (q *QueryBuilder) Select(fields ...string) *QueryBuilder {
	q.parts = append(q.parts, "SELECT "+strings.Join(fields, ","))
	return q
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
	q.parts = append(q.parts, "FROM "+table)
	return q
}

func (q *QueryBuilder) Where(condition string) *QueryBuilder {
	q.parts = append(q.parts, "WHERE "+condition)
	return q
}

func (q *QueryBuilder) Build() string {
	return strings.Join(q.parts, " ")
}

// MethodChainingBuilder tests builder pattern.
func MethodChainingBuilder() string {
	query := &QueryBuilder{}
	query.Select("id", "name").
		From("users").
		Where("active = true")
	return query.Build()
}

// ============================================================================
// TABLE-DRIVEN AND FUNCTION VALUES
// ============================================================================

// TableDrivenOp simulates table-driven test style operations.
func TableDrivenOp() int {
	type op struct {
		name string
		a, b int
		fn   func(int, int) int
	}

	ops := []op{
		{"add", 2, 3, func(a, b int) int { return a + b }},
		{"sub", 5, 3, func(a, b int) int { return a - b }},
		{"mul", 4, 3, func(a, b int) int { return a * b }},
	}

	sum := 0
	for _, op := range ops {
		sum += op.fn(op.a, op.b)
	}
	return sum
}

// FunctionValueFromMap tests function values stored in map.
func FunctionValueFromMap() int {
	m := map[string]func(int, int) int{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
	}

	result := m["add"](10, 5)
	result += m["sub"](10, 5)
	result += m["mul"](10, 5)
	return result
}

// ============================================================================
// DEFER WITH EXTERNAL PACKAGES
// ============================================================================

// DeferWithMutex tests defer with mutex unlock.
func DeferWithMutex() int {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	return 42
}

// ============================================================================
// SELECT WITH CHANNELS
// ============================================================================

// SelectWithChannels tests select with sync channels.
func SelectWithChannels() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 42

	select {
	case v := <-ch1:
		return v
	case v := <-ch2:
		return v
	default:
		return -1
	}
}

// ============================================================================
// GENERIC-LIKE CONTAINER PATTERNS
// ============================================================================

// GenericStack simulation using interface{}.
type GenericStack struct {
	items []interface{}
}

func (s *GenericStack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *GenericStack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *GenericStack) Len() int {
	return len(s.items)
}

// GenericContainerOps tests generic-like container operations.
func GenericContainerOps() int {
	stack := &GenericStack{}
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	sum := 0
	for stack.Len() > 0 {
		sum += stack.Pop().(int)
	}
	return sum
}

// ============================================================================
// INTERFACE-BASED IO PATTERNS
// ============================================================================

// ReaderInterface simulates io.Reader interface.
type ReaderInterface interface {
	Read(p []byte) (n int, err error)
}

// BytesReaderAdapter adapts bytes.Reader to our interface.
type BytesReaderAdapter struct {
	reader *bytes.Reader
}

func (a *BytesReaderAdapter) Read(p []byte) (n int, err error) {
	return a.reader.Read(p)
}

// ReadWithInterface tests interface-based reading.
func ReadWithInterface() int {
	adapter := &BytesReaderAdapter{reader: bytes.NewReader([]byte("hello"))}
	buf := make([]byte, 10)
	n, _ := adapter.Read(buf)
	return n
}

// WriterInterface simulates io.Writer interface.
type WriterInterface interface {
	Write(p []byte) (n int, err error)
}

// BytesWriterAdapter adapts bytes.Buffer to our interface.
type BytesWriterAdapter struct {
	buf *bytes.Buffer
}

func (a *BytesWriterAdapter) Write(p []byte) (n int, err error) {
	return a.buf.Write(p)
}

// WriteWithInterface tests interface-based writing.
func WriteWithInterface() int {
	adapter := &BytesWriterAdapter{buf: new(bytes.Buffer)}
	n, _ := adapter.Write([]byte("world"))
	return n
}
