package thirdparty

import "bytes"

// BytesBufferWrite tests bytes.Buffer write operations.
func BytesBufferWrite() int {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte("hello"))
	buf.Write([]byte(" "))
	buf.Write([]byte("world"))
	return buf.Len()
}

// BytesBufferWriteString tests buf.WriteString.
func BytesBufferWriteString() int {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("test")
	return buf.Len()
}

// BytesBufferReadFrom tests reading from a reader.
func BytesBufferReadFrom() int {
	buf := new(bytes.Buffer)
	buf.Write([]byte("source data"))
	result := new(bytes.Buffer)
	result.ReadFrom(buf)
	return result.Len()
}

// BytesBufferString tests Buffer.String().
func BytesBufferString() string {
	buf := new(bytes.Buffer)
	buf.Write([]byte("test string"))
	return buf.String()
}

// BytesBufferLen tests Buffer.Len() and Buffer.Cap().
func BytesBufferLen() int {
	buf := bytes.NewBuffer([]byte("12345"))
	return buf.Len() + buf.Cap()
}

// BytesBufferGrow tests Buffer.Grow().
func BytesBufferGrow() int {
	buf := new(bytes.Buffer)
	buf.Grow(100)
	buf.Write([]byte("x"))
	return buf.Cap()
}

// BytesBufferNext tests reading from buffer after partial read.
func BytesBufferNext() int {
	buf := bytes.NewBuffer([]byte("hello world"))
	prefix := make([]byte, 5)
	buf.Read(prefix)
	return len(buf.Bytes())
}

// BytesBufferReadByte tests Buffer.ReadByte().
func BytesBufferReadByte() int {
	buf := bytes.NewBuffer([]byte("ab"))
	b, _ := buf.ReadByte()
	return int(b)
}

// BytesBufferUnreadByte tests Buffer.UnreadByte().
func BytesBufferUnreadByte() int {
	buf := bytes.NewBuffer([]byte("ab"))
	b1, _ := buf.ReadByte()
	buf.UnreadByte()
	b2, _ := buf.ReadByte()
	if b1 == b2 {
		return 1
	}
	return 0
}

// BytesBufferReadBytes tests Buffer.ReadBytes().
func BytesBufferReadBytes() int {
	buf := bytes.NewBuffer([]byte("hello\nworld\n"))
	line, _ := buf.ReadBytes('\n')
	return len(line)
}

// BytesBufferReadString tests Buffer.ReadString().
func BytesBufferReadString() int {
	buf := bytes.NewBuffer([]byte("hello\nworld\n"))
	line, _ := buf.ReadString('\n')
	return len(line)
}

// BytesNewBuffer tests bytes.NewBuffer.
func BytesNewBuffer() int {
	buf := bytes.NewBuffer([]byte("test"))
	return buf.Len()
}

// BytesNewBufferString tests bytes.NewBufferString.
func BytesNewBufferString() int {
	buf := bytes.NewBufferString("test")
	return buf.Len()
}

// BytesBufferTrim tests trimming buffer bytes.
func BytesBufferTrim() int {
	buf := bytes.NewBuffer([]byte("  hello  "))
	trimmed := bytes.TrimSpace(buf.Bytes())
	return len(trimmed)
}

// BytesSplit tests bytes.Split.
func BytesSplit() int {
	s := []byte("a,b,c")
	parts := bytes.Split(s, []byte(","))
	return len(parts)
}

// BytesSplitN tests bytes.SplitN.
func BytesSplitN() int {
	s := []byte("a,b,c,d")
	parts := bytes.SplitN(s, []byte(","), 2)
	return len(parts)
}

// BytesJoin tests bytes.Join.
func BytesJoin() int {
	parts := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	result := bytes.Join(parts, []byte(","))
	return len(result)
}

// BytesContains tests bytes.Contains.
func BytesContains() int {
	if bytes.Contains([]byte("hello world"), []byte("world")) {
		return 1
	}
	return 0
}

// BytesCount tests bytes.Count.
func BytesCount() int {
	return bytes.Count([]byte("hello hello hello"), []byte("he"))
}

// BytesIndex tests bytes.Index.
func BytesIndex() int {
	return bytes.Index([]byte("hello"), []byte("ll"))
}

// BytesLastIndex tests bytes.LastIndex.
func BytesLastIndex() int {
	return bytes.LastIndex([]byte("hello world hello"), []byte("hello"))
}

// BytesHasPrefix tests bytes.HasPrefix.
func BytesHasPrefix() int {
	if bytes.HasPrefix([]byte("hello world"), []byte("hello")) {
		return 1
	}
	return 0
}

// BytesHasSuffix tests bytes.HasSuffix.
func BytesHasSuffix() int {
	if bytes.HasSuffix([]byte("hello world"), []byte("world")) {
		return 1
	}
	return 0
}

// BytesReplace tests bytes.Replace.
func BytesReplace() int {
	s := []byte("hello world world")
	result := bytes.Replace(s, []byte("world"), []byte("go"), 1)
	return len(result)
}

// BytesReplaceAll tests bytes.ReplaceAll.
func BytesReplaceAll() int {
	s := []byte("hello world world")
	result := bytes.ReplaceAll(s, []byte("world"), []byte("go"))
	return len(result)
}

// BytesFields tests bytes.Fields.
func BytesFields() int {
	s := []byte("hello world test")
	fields := bytes.Fields(s)
	return len(fields)
}

// BytesTrimSpace tests bytes.TrimSpace.
func BytesTrimSpace() int {
	s := []byte("  hello  ")
	return len(bytes.TrimSpace(s))
}

// BytesToUpper tests bytes.ToUpper.
func BytesToUpper() int {
	return len(bytes.ToUpper([]byte("hello")))
}

// BytesToLower tests bytes.ToLower.
func BytesToLower() int {
	return len(bytes.ToLower([]byte("HELLO")))
}

// BytesTrim tests bytes.Trim.
func BytesTrim() int {
	s := []byte("xxhelloxx")
	return len(bytes.Trim(s, "x"))
}

// BytesMap tests bytes.Map.
func BytesMap() int {
	s := []byte("hello")
	mapped := bytes.Map(func(r rune) rune {
		return r + 1
	}, s)
	if string(mapped) == "ifmmp" {
		return 1
	}
	return 0
}
