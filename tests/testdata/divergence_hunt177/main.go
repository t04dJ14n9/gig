package divergence_hunt177

import (
	"bytes"
	"fmt"
)

// ============================================================================
// Round 177: Byte slice operations (bytes package)
// ============================================================================

func BytesCompare() string {
	a := []byte("abc")
	b := []byte("def")
	return fmt.Sprintf("%d", bytes.Compare(a, b))
}

func BytesEqual() string {
	a := []byte("hello")
	b := []byte("hello")
	c := []byte("world")
	return fmt.Sprintf("%v:%v", bytes.Equal(a, b), bytes.Equal(a, c))
}

func BytesContains() string {
	s := []byte("hello world")
	sub := []byte("world")
	return fmt.Sprintf("%v", bytes.Contains(s, sub))
}

func BytesIndex() string {
	s := []byte("hello world hello")
	sub := []byte("hello")
	return fmt.Sprintf("%d", bytes.Index(s, sub))
}

func BytesLastIndex() string {
	s := []byte("hello world hello")
	sub := []byte("hello")
	return fmt.Sprintf("%d", bytes.LastIndex(s, sub))
}

func BytesCount() string {
	s := []byte("hello world hello")
	sub := []byte("hello")
	return fmt.Sprintf("%d", bytes.Count(s, sub))
}

func BytesReplace() string {
	s := []byte("hello world")
	result := bytes.Replace(s, []byte("world"), []byte("go"), 1)
	return string(result)
}

func BytesReplaceAll() string {
	s := []byte("hello hello hello")
	result := bytes.ReplaceAll(s, []byte("hello"), []byte("hi"))
	return string(result)
}

func BytesRepeat() string {
	result := bytes.Repeat([]byte("ab"), 3)
	return string(result)
}

func BytesToUpper() string {
	s := []byte("hello")
	result := bytes.ToUpper(s)
	return string(result)
}

func BytesToLower() string {
	s := []byte("HELLO")
	result := bytes.ToLower(s)
	return string(result)
}

func BytesTrim() string {
	s := []byte("   hello   ")
	result := bytes.Trim(s, " ")
	return string(result)
}

func BytesTrimSpace() string {
	s := []byte("\t\n hello \r\n")
	result := bytes.TrimSpace(s)
	return string(result)
}

func BytesTrimPrefix() string {
	s := []byte("hello world")
	result := bytes.TrimPrefix(s, []byte("hello "))
	return string(result)
}

func BytesTrimSuffix() string {
	s := []byte("hello world")
	result := bytes.TrimSuffix(s, []byte(" world"))
	return string(result)
}

func BytesSplit() string {
	s := []byte("a,b,c")
	parts := bytes.Split(s, []byte(","))
	return fmt.Sprintf("%d", len(parts))
}

func BytesJoin() string {
	parts := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	result := bytes.Join(parts, []byte("-"))
	return string(result)
}

func BytesHasPrefix() string {
	s := []byte("hello world")
	return fmt.Sprintf("%v", bytes.HasPrefix(s, []byte("hello")))
}

func BytesHasSuffix() string {
	s := []byte("hello world")
	return fmt.Sprintf("%v", bytes.HasSuffix(s, []byte("world")))
}

func BytesFields() string {
	s := []byte("  hello   world  ")
	fields := bytes.Fields(s)
	return fmt.Sprintf("%d:%s:%s", len(fields), fields[0], fields[1])
}
