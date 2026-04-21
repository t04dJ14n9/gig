package divergence_hunt178

import (
	"bytes"
	"fmt"
)

// ============================================================================
// Round 178: Buffer operations
// ============================================================================

func BufferNew() string {
	var b bytes.Buffer
	return fmt.Sprintf("len=%d", b.Len())
}

func BufferNewString() string {
	b := bytes.NewBufferString("hello")
	return b.String()
}

func BufferWriteString() string {
	var b bytes.Buffer
	b.WriteString("hello")
	b.WriteString(" ")
	b.WriteString("world")
	return b.String()
}

func BufferWriteByte() string {
	var b bytes.Buffer
	b.WriteByte('h')
	b.WriteByte('i')
	return b.String()
}

func BufferWriteRune() string {
	var b bytes.Buffer
	b.WriteRune('h')
	b.WriteRune('i')
	return b.String()
}

func BufferWrite() string {
	var b bytes.Buffer
	b.Write([]byte("hello"))
	return b.String()
}

func BufferLen() string {
	var b bytes.Buffer
	b.WriteString("hello")
	return fmt.Sprintf("len=%d", b.Len())
}

func BufferCap() string {
	var b bytes.Buffer
	b.WriteString("hello world this is a test")
	return fmt.Sprintf("cap>0=%v", b.Cap() > 0)
}

func BufferReset() string {
	var b bytes.Buffer
	b.WriteString("hello")
	b.Reset()
	return fmt.Sprintf("len=%d", b.Len())
}

func BufferBytes() string {
	var b bytes.Buffer
	b.WriteString("hello")
	return string(b.Bytes())
}

func BufferGrow() string {
	var b bytes.Buffer
	b.Grow(100)
	return fmt.Sprintf("cap>0=%v", b.Cap() > 0)
}

func BufferRead() string {
	b := bytes.NewBufferString("hello")
	buf := make([]byte, 3)
	n, _ := b.Read(buf)
	return fmt.Sprintf("%d:%s", n, string(buf))
}

func BufferReadString() string {
	b := bytes.NewBufferString("hello\nworld")
	line, _ := b.ReadString('\n')
	return line
}

func BufferReadByte() string {
	b := bytes.NewBufferString("abc")
	c, _ := b.ReadByte()
	return fmt.Sprintf("%c", c)
}

func BufferNext() string {
	b := bytes.NewBufferString("hello")
	next := b.Next(3)
	return string(next)
}

func BufferUnreadByte() string {
	b := bytes.NewBufferString("abc")
	b.ReadByte()
	b.UnreadByte()
	c, _ := b.ReadByte()
	return fmt.Sprintf("%c", c)
}

func BufferTruncate() string {
	var b bytes.Buffer
	b.WriteString("hello")
	b.Truncate(3)
	return b.String()
}

func BufferAvailable() string {
	var b bytes.Buffer
	b.Grow(100)
	return fmt.Sprintf("avail>0=%v", b.Available() > 0)
}

func BufferAvailableBuffer() string {
	var b bytes.Buffer
	b.Grow(10)
	buf := b.AvailableBuffer()
	return fmt.Sprintf("len=%d", len(buf))
}

func BufferStringEmpty() string {
	var b bytes.Buffer
	return b.String()
}
