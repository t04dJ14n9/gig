package divergence_hunt179

import (
	"bytes"
	"fmt"
	"strings"
)

// ============================================================================
// Round 179: Reader/Writer interfaces
// ============================================================================

func ReaderFromString() string {
	r := strings.NewReader("hello world")
	buf := make([]byte, 5)
	n, _ := r.Read(buf)
	return fmt.Sprintf("%d:%s", n, string(buf))
}

func ReaderReadAt() string {
	r := strings.NewReader("hello world")
	buf := make([]byte, 5)
	n, _ := r.ReadAt(buf, 6)
	return fmt.Sprintf("%d:%s", n, string(buf))
}

func ReaderSeek() string {
	r := strings.NewReader("hello world")
	r.Seek(6, 0)
	buf := make([]byte, 5)
	n, _ := r.Read(buf)
	return fmt.Sprintf("%d:%s", n, string(buf))
}

func ReaderSize() string {
	r := strings.NewReader("hello world")
	return fmt.Sprintf("%d", r.Size())
}

func ReaderLen() string {
	r := strings.NewReader("hello world")
	r.Read(make([]byte, 5))
	return fmt.Sprintf("%d", r.Len())
}

func WriterToBuffer() string {
	var buf bytes.Buffer
	w := &buf
	w.WriteString("hello")
	w.Write([]byte(" world"))
	return buf.String()
}

func WriterWriteByte() string {
	var buf bytes.Buffer
	w := &buf
	w.WriteByte('a')
	w.WriteByte('b')
	w.WriteByte('c')
	return buf.String()
}

func WriterWriteRune() string {
	var buf bytes.Buffer
	w := &buf
	w.WriteRune('h')
	w.WriteRune('i')
	return buf.String()
}

func ReaderWriterCombined() string {
	var buf bytes.Buffer
	// Write
	buf.WriteString("hello")
	// Read
	result := make([]byte, 5)
	buf.Read(result)
	return string(result)
}

func ReadFull() string {
	r := strings.NewReader("hello")
	buf := make([]byte, 3)
	// Using our own read
	n, _ := r.Read(buf)
	return fmt.Sprintf("%d:%s", n, string(buf))
}

func WriteString() string {
	var buf bytes.Buffer
	n, _ := buf.WriteString("hello world")
	return fmt.Sprintf("%d:%s", n, buf.String())
}

func CopyReaderToWriter() string {
	r := strings.NewReader("hello")
	var buf bytes.Buffer
	// Simulate copy
	tmp := make([]byte, 5)
	n, _ := r.Read(tmp)
	buf.Write(tmp[:n])
	return buf.String()
}

func MultiWrite() string {
	var buf bytes.Buffer
	data := []byte("hello")
	buf.Write(data)
	buf.Write(data)
	return buf.String()
}

func LimitedRead() string {
	r := strings.NewReader("hello world")
	// Read only first 5 bytes
	buf := make([]byte, 5)
	n, _ := r.Read(buf)
	return fmt.Sprintf("%d:%s", n, string(buf))
}

func DiscardRead() string {
	r := strings.NewReader("hello world")
	// Discard first 6 bytes
	buf := make([]byte, 6)
	r.Read(buf)
	// Read rest
	rest := make([]byte, 5)
	n, _ := r.Read(rest)
	return fmt.Sprintf("%d:%s", n, string(rest))
}

func PipeSimulation() string {
	// Simulate a simple pipe using buffer
	var buf bytes.Buffer
	// Write end
	buf.WriteString("hello")
	// Read end
	result := make([]byte, 5)
	buf.Read(result)
	return string(result)
}

func TeeReaderSimulation() string {
	// Simulate tee reader - read and also capture
	source := strings.NewReader("hello")
	var copy bytes.Buffer
	buf := make([]byte, 5)
	n, _ := source.Read(buf)
	copy.Write(buf[:n])
	return fmt.Sprintf("original=%s copy=%s", string(buf), copy.String())
}

func SectionReaderSimulation() string {
	// Simulate section reader - read from offset
	source := strings.NewReader("0123456789")
	// Seek to offset 3
	source.Seek(3, 0)
	buf := make([]byte, 4)
	n, _ := source.Read(buf)
	return fmt.Sprintf("%d:%s", n, string(buf))
}
