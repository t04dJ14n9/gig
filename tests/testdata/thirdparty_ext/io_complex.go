package thirdparty

import (
	"bytes"
	"io"
)

// ============================================================================
// BUFFERED IO WITH EXTERNAL PACKAGES
// ============================================================================

// BufferedReader tests reading from bytes.Reader.
func BufferedReader() int {
	reader := bytes.NewReader([]byte("hello world"))
	buf := make([]byte, 5)
	n, _ := reader.Read(buf)
	return n
}

// BufferedWriter tests writing to bytes.Buffer.
func BufferedWriter() int {
	buf := new(bytes.Buffer)
	buf.Write([]byte("hello"))
	buf.Write([]byte(" world"))
	return buf.Len()
}

// MultiWriter tests io.MultiWriter.
func MultiWriter() int {
	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)
	mw := io.MultiWriter(buf1, buf2)
	mw.Write([]byte("test"))
	return buf1.Len() + buf2.Len()
}

// TeeReaderWrapper wraps io.TeeReader for testing.
func TeeReaderWrapper() int {
	reader := bytes.NewReader([]byte("hello"))
	writer := new(bytes.Buffer)
	tee := io.TeeReader(reader, writer)
	result := make([]byte, 5)
	n, _ := tee.Read(result)
	return n
}

// ============================================================================
// COMPLEX INTERFACE PATTERNS - Interface embedding and composition
// ============================================================================

// ReadWriter interface composition.
type ReadWriter interface {
	io.Reader
	io.Writer
}

// BytesBufferAsReadWriter tests bytes.Buffer as ReadWriter interface.
func BytesBufferAsReadWriter() int {
	var buf bytes.Buffer
	var rw ReadWriter = &buf
	rw.Write([]byte("hello"))
	rw.Read(make([]byte, 10))
	return buf.Len()
}

// InterfaceValueNil tests nil interface handling.
func InterfaceValueNil() int {
	var reader io.Reader
	if reader == nil {
		return 1
	}
	return 0
}
