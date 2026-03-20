package thirdparty

import (
	"bytes"
	"io"
	"strings"
)

// IoReadAll tests io.ReadAll.
func IoReadAll() int {
	reader := strings.NewReader("hello world")
	result, _ := io.ReadAll(reader)
	return len(result)
}

// IoCopy tests io.Copy.
func IoCopy() int {
	reader := strings.NewReader("hello")
	writer := new(bytes.Buffer)
	io.Copy(writer, reader)
	return writer.Len()
}

// IoReadFull tests io.ReadFull.
func IoReadFull() int {
	reader := strings.NewReader("hello world")
	buf := make([]byte, 5)
	n, _ := io.ReadFull(reader, buf)
	return n
}

// IoWriteString tests io.WriteString.
func IoWriteString() int {
	writer := new(bytes.Buffer)
	n, _ := io.WriteString(writer, "test")
	return n
}

// IoSectionReader tests io.NewSectionReader.
func IoSectionReader() int {
	reader := strings.NewReader("hello world")
	section := io.NewSectionReader(reader, 0, 5)
	buf := make([]byte, 5)
	n, _ := section.Read(buf)
	return n
}

// IoLimitedReader tests io.LimitedReader.
func IoLimitedReader() int {
	reader := strings.NewReader("hello world")
	limited := io.LimitedReader{R: reader, N: 5}
	result, _ := io.ReadAll(&limited)
	return len(result)
}

// IoTeeReader tests io.TeeReader.
func IoTeeReader() int {
	reader := strings.NewReader("hello")
	writer := new(bytes.Buffer)
	tee := io.TeeReader(reader, writer)
	io.ReadAll(tee)
	return writer.Len()
}
